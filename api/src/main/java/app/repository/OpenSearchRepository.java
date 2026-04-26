package app.repository;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.NoSuchElementException;

import org.opensearch.client.opensearch.OpenSearchClient;
import org.opensearch.client.opensearch._types.aggregations.Aggregate;
import org.opensearch.client.opensearch._types.aggregations.StringTermsBucket;
import org.opensearch.client.opensearch.core.BulkResponse;
import org.opensearch.client.opensearch.core.DeleteByQueryResponse;
import org.opensearch.client.opensearch.core.SearchResponse;
import org.opensearch.client.opensearch.core.bulk.BulkOperation;
import org.springframework.stereotype.Repository;
import org.springframework.validation.annotation.Validated;

import app.service.TextEmbeddingService;
import jakarta.annotation.PostConstruct;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;

@Repository
@Validated
public class OpenSearchRepository {
    private final OpenSearchClient openSearchClient;

    public OpenSearchRepository(OpenSearchClient openSearchClient) {
        this.openSearchClient = openSearchClient;
    }

    @PostConstruct
    private void init() throws IOException {
        createIndexIfNotExist("github-issues");
    }

    public void indexGithubIssue(@NotEmpty List<TextEmbeddingService.embeddingDocument> issueDocuments) throws IOException {
        List<BulkOperation> operations = new ArrayList<>();

        for (TextEmbeddingService.embeddingDocument doc : issueDocuments) {
            operations.add(new BulkOperation.Builder()
                .index(i -> i.document(doc))
                .build()
            );
        }
                
        BulkResponse bulkRes = openSearchClient.bulk(r -> r
            .index("github-issues")
            .operations(operations)
        );

        if (bulkRes.errors()) {
            List<String> failures = bulkRes.items().stream()
                .filter(i -> i.error() != null)
                .map(i -> {                                                                                                                                                            
                    var error = i.error();                                                                                                                                             
                    String reason = error != null ? error.reason() : null;                                                                                                             
                    return i.id() + ": " + (reason != null ? reason : "unknown error");                                                                                                
                })
                .toList();
            
            throw new RuntimeException("Bulk index had failures: " + failures);
        }
    }

    public List<String> listAllIndexedRepoNames() throws IOException {
        final int maxUniqueRepos = 1000;

        SearchResponse<Void> searchRes = openSearchClient.search(r -> r
            .index("github-issues")
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

    public void deleteTrackedRepo(@NotBlank String repoName) throws IOException {
        DeleteByQueryResponse deleteRes =  openSearchClient.deleteByQuery(d -> d
            .index("github-issues")
            .query(q -> q
                .term(t -> t
                    .field("github-issues")
                    .value(v -> v.stringValue(repoName))
                )
            )
        );

        if (deleteRes.deleted() == 0) {
            throw new NoSuchElementException("No repo found to delete");
        }
    }

    // dim number from: https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2
    private static final Integer embeddingDim = 384;

    private void createIndexIfNotExist(@NotBlank String indexName) throws IOException {
        boolean exists = openSearchClient.indices().exists(r -> r.index(indexName)).value();

        if (!exists) {
            openSearchClient.indices().create(r -> r
                .index(indexName)
                .settings(s -> s.knn(true))
                .mappings(m -> m // using keyword instead of text since we only need it for displaying/filtering
                    .properties("title", p -> p.keyword(k -> k))
                    .properties("repoName", p -> p.keyword(k -> k))
                    .properties("url", p -> p.keyword(k -> k))
                    .properties("embedding", p -> p.knnVector(k -> k
                        .dimension(embeddingDim)
                        .method(met -> met
                            .name("hnsw")
                            .spaceType("cosinesimil")
                            .engine("faiss")
                        )
                    ))
                )
            );
        }
    }
}
