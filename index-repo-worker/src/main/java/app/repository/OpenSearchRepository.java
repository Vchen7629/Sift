package app.repository;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

import org.opensearch.client.opensearch.OpenSearchClient;
import org.opensearch.client.opensearch.core.BulkResponse;
import org.opensearch.client.opensearch.core.bulk.BulkOperation;
import org.springframework.stereotype.Repository;
import org.springframework.validation.annotation.Validated;

import app.service.TextEmbeddingService;
import jakarta.annotation.PostConstruct;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;

@Repository
@Validated
public class OpenSearchRepository {
    private final OpenSearchClient openSearchClient;
    private final static String issuesIndexName = "github-issues";
    private final static String jobStatusIndexName = "job-status";

    public OpenSearchRepository(OpenSearchClient openSearchClient) {
        this.openSearchClient = openSearchClient;
    }

    @PostConstruct
    private void init() throws IOException {
        createIndexIfNotExist();
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
            .index(issuesIndexName)
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

    public record JobStatus(@NotBlank String repoName, @NotBlank String status) {}

    public void upsertJobStatus(@Valid JobStatus jobStatus) throws IOException {
        openSearchClient.update(r -> r
            .index(jobStatusIndexName)
            .id(jobStatus.repoName)
            .doc(new JobStatus(jobStatus.repoName, jobStatus.status))
            .docAsUpsert(true)
            , JobStatus.class
        );
    }

    // dim number from: https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2
    private static final Integer embeddingDim = 384;

    private void createIndexIfNotExist() throws IOException {
        boolean exists = openSearchClient.indices().exists(r -> r.index(issuesIndexName)).value();

        if (!exists) {
            openSearchClient.indices().create(r -> r
                .index(issuesIndexName)
                .settings(s -> s.knn(true))
                .mappings(m -> m // using keyword instead of text since we only need it for displaying/filtering
                    .properties("title", p -> p.keyword(k -> k))
                    .properties("body", p -> p.keyword(k -> k))
                    .properties("repoName", p -> p.keyword(k -> k))
                    .properties("url", p -> p.keyword(k -> k))
                    .properties("titleEmbedding", p -> p.knnVector(k -> k
                        .dimension(embeddingDim)
                        .method(met -> met
                            .name("hnsw")
                            .spaceType("cosinesimil")
                            .engine("lucene")
                        )
                    ))
                    .properties("bodyEmbedding", p -> p.knnVector(k -> k
                        .dimension(embeddingDim)
                        .method(met -> met
                            .name("hnsw")
                            .spaceType("cosinesimil")
                            .engine("lucene")
                        )
                    ))
                )
            );
        }
    }

}
