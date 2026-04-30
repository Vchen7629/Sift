package app.repository;

import java.io.IOException;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.NoSuchElementException;
import java.util.Objects;
import java.util.Set;

import org.opensearch.client.opensearch.OpenSearchClient;
import org.opensearch.client.opensearch._types.aggregations.Aggregate;
import org.opensearch.client.opensearch._types.aggregations.StringTermsBucket;
import org.opensearch.client.opensearch._types.query_dsl.Query;
import org.opensearch.client.opensearch._types.query_dsl.TextQueryType;
import org.opensearch.client.opensearch.core.DeleteByQueryResponse;
import org.opensearch.client.opensearch.core.SearchRequest;
import org.opensearch.client.opensearch.core.SearchResponse;
import org.opensearch.client.opensearch.search_pipeline.ScoreCombinationTechnique;
import org.opensearch.client.opensearch.search_pipeline.ScoreNormalizationTechnique;
import org.springframework.stereotype.Repository;
import org.springframework.validation.annotation.Validated;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import ai.djl.translate.TranslateException;
import app.service.TextEmbeddingService;
import jakarta.annotation.PostConstruct;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;

@Repository
@Validated
public class OpenSearchRepository {
    private final OpenSearchClient openSearchClient;
    private final TextEmbeddingService textEmbeddingService;
    private final static String indexName = "github-issues";

    public OpenSearchRepository(
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

    public void deleteTrackedRepo(@NotBlank String repoName) throws IOException {
        DeleteByQueryResponse deleteRes =  openSearchClient.deleteByQuery(d -> d
            .index(indexName)
            .query(q -> q
                .term(t -> t
                    .field("repoName")
                    .value(v -> v.stringValue(repoName))
                )
            )
        );

        if (deleteRes.deleted() == 0) {
            throw new NoSuchElementException("No repo found to delete");
        }
    }

    @JsonIgnoreProperties(ignoreUnknown = true)
    public static record IssueSearchResult(String url, String title, String body) {};

    public List<IssueSearchResult> findRelevantIssues(
        @NotBlank String repoName,
        @NotBlank String searchQuery
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

        Query repoFilter = Query.of(q -> q
            .term(t -> t
                .field("repoName")
                .value(v -> v.stringValue(repoName))
            )
        );

        Query titleSemQuery = Query.of(q -> q
            .knn(k -> k
                .field("titleEmbedding")
                .vector(searchVector)
                .k(10)
                .filter(repoFilter)
            )
        );

        Query bodySemQuery = Query.of(q -> q
            .knn(k -> k
                .field("bodyEmbedding")
                .vector(searchVector)
                .k(10)
                .filter(repoFilter)
            )
        );

        Query hybridQuery = Query.of(q -> q
            .hybrid(h -> h
                .queries(keywordQuery, titleSemQuery, bodySemQuery)
            )
        );

        var resultAmount = 10;
        SearchRequest searchRequest = SearchRequest.of(s -> s
            .index(indexName)
            .query(hybridQuery)
            .postFilter(repoFilter)
            .size(resultAmount)
            .timeout("30s")
        );

        SearchResponse<IssueSearchResult> searchRes = openSearchClient.search(searchRequest, IssueSearchResult.class);
    
        return deduplicateSearchResults(searchRes.hits().hits().stream()
            .map(hit -> hit.source())
            .filter(Objects::nonNull)
            .toList());
    }

    // opensearch returns duplicate issues sometimes, using url to deduplicate
    // since each issue has a unique url.
    private final List<IssueSearchResult> deduplicateSearchResults(
        @NotEmpty @Valid List<IssueSearchResult> searchResults
    ) {
        List<IssueSearchResult> deduplicatedList = new ArrayList<>();
        Set<String> seenIds = new HashSet<>();

        for (IssueSearchResult issue : searchResults) {
            boolean notSeen = seenIds.add(issue.url());

            if (notSeen) {
                deduplicatedList.add(issue);
            }
        }

        return deduplicatedList;
    }

    public List<String> findAllIndexedRepoNames() throws IOException {
        final int maxUniqueRepos = 1000;

        SearchResponse<Void> searchRes = openSearchClient.search(r -> r
            .index(indexName)
            .size(0) // need this so we dont return the actual document, use aggregations to return the strings instead
            .timeout("30s")
            .aggregations("repoNames", a -> a
                .terms(t -> t.field("repoName").size(maxUniqueRepos))
            ),
            Void.class
        );

        Aggregate repoNameAgg = searchRes.aggregations().get("repoNames");
        if (repoNameAgg == null || !repoNameAgg.isSterms()) {
            return List.of(); // return empty list instead of throwing an exception
        }

        return repoNameAgg
            .sterms()
            .buckets().array()
            .stream()
            .map(StringTermsBucket::key)
            .toList();
    }

    public boolean isRepoIndexed(String repoName) throws IOException{
        SearchResponse<Void> searchRes = openSearchClient.search(r -> r
            .index(indexName)
            .timeout("30s")
            .query(q -> q
                .term(t -> t
                    .field("repoName")
                    .value(v -> v.stringValue(repoName))
                )
            ), 
            Void.class
        );

        var total = searchRes.hits().total();

        return total != null && total.value() > 0;
    }

    private void createSearchPipelineIfNotExist() throws IOException {
        boolean exists = openSearchClient.searchPipeline()
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
        }
    }
}
