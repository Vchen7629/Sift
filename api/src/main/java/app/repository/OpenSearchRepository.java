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
import lombok.extern.slf4j.Slf4j;
import static net.logstash.logback.argument.StructuredArguments.kv;


@Repository
@Validated
@Slf4j
public class OpenSearchRepository {
    private final OpenSearchClient openSearchClient;
    private final TextEmbeddingService textEmbeddingService;
    private final static String issuesIndexName = "github-issues";
    private final static String jobStatusIndexName = "job-status";

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
        createIndexIfNotExist();
    }

    public void deleteTrackedRepo(@NotBlank String repoName, @NotBlank String requestId) throws IOException {
        DeleteByQueryResponse deleteRes =  openSearchClient.deleteByQuery(d -> d
            .index(issuesIndexName)
            .query(q -> q
                .term(t -> t
                    .field("repoName")
                    .value(v -> v.stringValue(repoName))
                )
            )
        );

        log.debug("successfully deleted track Repo", 
            kv("repoName", repoName), 
            kv("requestId", requestId));

        if (deleteRes.deleted() == 0) {
            log.warn("No repo found in db to delete", 
                kv("repoName", repoName),
                kv("requestId", requestId));

            throw new NoSuchElementException("No repo found to delete");
        }
    }

    @JsonIgnoreProperties(ignoreUnknown = true)
    public static record IssueSearchResult(String url, String title, String body) {};

    public List<IssueSearchResult> findRelevantIssues(
        @NotBlank String repoName,
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
            .index(issuesIndexName)
            .query(hybridQuery)
            .postFilter(repoFilter)
            .size(resultAmount)
            .timeout("30s")
        );

        SearchResponse<IssueSearchResult> searchRes = openSearchClient.search(searchRequest, IssueSearchResult.class);
    
        return deduplicateSearchResults((searchRes.hits().hits().stream()
            .map(hit -> hit.source())
            .filter(Objects::nonNull)
            .toList()
        ), requestId);
    }

    // opensearch returns duplicate issues sometimes, using url to deduplicate
    // since each issue has a unique url.
    private final List<IssueSearchResult> deduplicateSearchResults(
        @NotEmpty @Valid List<IssueSearchResult> searchResults,
        @NotBlank String requestId
    ) {
        List<IssueSearchResult> deduplicatedList = new ArrayList<>();
        Set<String> seenIds = new HashSet<>();

        int deduplicatedCount = 0;

        for (IssueSearchResult issue : searchResults) {
            boolean notSeen = seenIds.add(issue.url());

            if (notSeen) {
                deduplicatedList.add(issue);
            } else {
                deduplicatedCount++;
            }
        }

        log.debug("deduplicated search results", 
            kv("requestId", requestId), 
            kv("duplicatesRemoved", deduplicatedCount));

        return deduplicatedList;
    }

    public List<String> findAllIndexedRepoNames(@NotBlank String requestId) throws IOException {
        final int maxUniqueRepos = 1000;

        SearchResponse<Void> searchRes = openSearchClient.search(r -> r
            .index(issuesIndexName)
            .size(0) // need this so we dont return the actual document, use aggregations to return the strings instead
            .timeout("30s")
            .aggregations("repoNames", a -> a
                .terms(t -> t.field("repoName").size(maxUniqueRepos))
            ),
            Void.class
        );

        Aggregate repoNameAgg = searchRes.aggregations().get("repoNames");
        if (repoNameAgg == null || !repoNameAgg.isSterms()) {
            log.debug("found no indexed repos in the db", kv("requestId", requestId));
            return List.of(); // return empty list instead of throwing an exception
        }

        List<String> indexedRepos = repoNameAgg
            .sterms()
            .buckets().array()
            .stream()
            .map(StringTermsBucket::key)
            .toList();

        log.debug("Successfully fetched {} indexed repos", indexedRepos.size(), 
            kv("requestId", requestId));

        return indexedRepos;
    }

    public boolean isRepoIndexed(String repoName) throws IOException{
        SearchResponse<Void> searchRes = openSearchClient.search(r -> r
            .index(issuesIndexName)
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

    private record JobStatus(@NotBlank String repoName, @NotBlank String status) {}

    public String findJobStatus(@NotBlank String repoName) throws IOException {
        SearchResponse<JobStatus> searchRes = openSearchClient.search(r -> r
            .index(jobStatusIndexName)
            .timeout("30s")
            .query(q -> q
                .term(t -> t
                    .field("repoName")
                    .value(v -> v.stringValue(repoName))
                )
            ), JobStatus.class
        );
        
        JobStatus dbRes = searchRes.hits().hits().get(0).source();
        
        return dbRes != null ? dbRes.status() : null;
    }

    private void createIndexIfNotExist() throws IOException {
        boolean exists = openSearchClient.indices().exists(r -> r
            .index(jobStatusIndexName)
        ).value();

        if (!exists) {
            openSearchClient.indices().create(r -> r
                .index(jobStatusIndexName)
                .mappings(m -> m
                    .properties("repoName", p -> p.keyword(k -> k))
                    .properties("status", p -> p.keyword(k -> k))
                )
            );
        }
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
