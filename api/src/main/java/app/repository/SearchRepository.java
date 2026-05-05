package app.repository;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Objects;

import org.opensearch.client.opensearch.OpenSearchClient;
import org.opensearch.client.opensearch._types.FieldValue;
import org.opensearch.client.opensearch._types.query_dsl.Query;
import org.opensearch.client.opensearch._types.query_dsl.TextQueryType;
import org.opensearch.client.opensearch.core.SearchRequest;
import org.opensearch.client.opensearch.core.SearchResponse;
import org.opensearch.client.opensearch.search_pipeline.ScoreCombinationTechnique;
import org.opensearch.client.opensearch.search_pipeline.ScoreNormalizationTechnique;
import org.springframework.stereotype.Repository;
import org.springframework.validation.annotation.Validated;

import ai.djl.translate.TranslateException;
import app.dto.IssueSearchResponse;
import app.service.TextEmbeddingService;
import jakarta.annotation.PostConstruct;
import jakarta.validation.constraints.NotBlank;
import lombok.extern.slf4j.Slf4j;
import static net.logstash.logback.argument.StructuredArguments.kv;


@Repository
@Validated
@Slf4j
public class SearchRepository {
    private final OpenSearchClient openSearchClient;
    private final TextEmbeddingService textEmbeddingService;
    private final static String issueIndexName = "dependency-issues";

    public SearchRepository(
        OpenSearchClient openSearchClient,
        TextEmbeddingService textEmbeddingService
    ) {
        this.openSearchClient = openSearchClient;
        this.textEmbeddingService = textEmbeddingService;
    }

    @PostConstruct
    private void init() throws IOException {
        createSearchPipelineIfNotExist();
    }

    public List<IssueSearchResponse> findRelevantIssues(
        Map<String, String> dependencyFilterList,
        @NotBlank String searchQuery,
        @NotBlank String requestId
    ) throws TranslateException, IOException {
        Query keywordQuery = Query.of(q -> q
            .multiMatch(m -> m
                .fields("title^1.1", "body")
                .query(searchQuery)
                .type(TextQueryType.BestFields)
            )
        );

        float[] searchEmbedding = textEmbeddingService.embedText(searchQuery);
        List<Float> searchVector = new ArrayList<>();
        for (float f : searchEmbedding) searchVector.add(f); 

        Query dependencyNameFilter = Query.of(q -> q
            .terms(t -> t
                .field("dependencyName")
                .terms(tv -> tv
                    .value(dependencyFilterList.keySet().stream()
                        .map(FieldValue::of)
                        .toList()
                    )
                )
            )
        );

        Query titleSemQuery = Query.of(q -> q
            .knn(k -> k
                .field("titleEmbedding")
                .vector(searchVector)
                .k(10)
                .filter(dependencyNameFilter)
            )
        );

        Query bodySemQuery = Query.of(q -> q
            .knn(k -> k
                .field("bodyEmbedding")
                .vector(searchVector)
                .k(10)
                .filter(dependencyNameFilter)
            )
        );

        Query hybridQuery = Query.of(q -> q
            .hybrid(h -> h
                .queries(keywordQuery, titleSemQuery, bodySemQuery)
            )
        );

        var resultAmount = 10;
        SearchRequest searchRequest = SearchRequest.of(s -> s
            .index(issueIndexName)
            .query(hybridQuery)
            .size(resultAmount)
            .timeout("30s")
            .searchPipeline("hybrid-search-pipeline")
            .collapse(c -> c.field("url"))
            .postFilter(dependencyNameFilter)
        );

        SearchResponse<IssueSearchResponse> searchRes = openSearchClient.search(searchRequest, IssueSearchResponse.class);

        log.debug("search returned {} hits", searchRes.hits().hits().size(), kv("requestId", requestId));

        return searchRes.hits().hits().stream()
            .map(hit -> hit.source())
            .filter(Objects::nonNull)
            .toList();
    }
    
    private void createSearchPipelineIfNotExist() throws IOException {
        boolean exists = !openSearchClient.searchPipeline()
            .get(r -> r.id("hybrid-search-pipeline"))
            .result()
            .isEmpty();

        if (!exists) {
            openSearchClient.searchPipeline().put(r -> r
                .id("hybrid-search-pipeline")
                .phaseResultsProcessors(p -> p
                    .normalizationProcessor(n -> n
                        .normalization(norm -> norm
                            .technique(ScoreNormalizationTechnique.MinMax)
                        )
                        .combination(comb -> comb
                            .technique(ScoreCombinationTechnique.ArithmeticMean)
                            .parameters(params -> params
                                .weights(List.of(0.3f, 0.35f, 0.35f)) // 30% BM25 35% titleEmb KNN 35% bodyEmb KNN
                            )
                        )
                    )
                )
            );
            log.debug("Created Hybrid search pipeline");
        }
    }
}
