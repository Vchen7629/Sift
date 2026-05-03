package app.repository;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

import org.opensearch.client.opensearch.OpenSearchClient;
import org.opensearch.client.opensearch.core.BulkResponse;
import org.opensearch.client.opensearch.core.bulk.BulkOperation;
import org.springframework.stereotype.Repository;
import org.springframework.validation.annotation.Validated;

import app.service.TextEmbeddingService.IndexableDocument;
import jakarta.annotation.PostConstruct;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import lombok.extern.slf4j.Slf4j;
import static net.logstash.logback.argument.StructuredArguments.kv;

@Repository
@Validated
@Slf4j
public class DependencyRepository {
    public final static String changeLogIndexName = "dependency-changelog";
    public final static String issuesIndexName = "dependency-issue";

    private final OpenSearchClient openSearchClient;

    public DependencyRepository(OpenSearchClient openSearchClient) {
        this.openSearchClient = openSearchClient;
    }

    @PostConstruct
    private void init() throws IOException {
        createChangeLogIndexIfNotExist();
        createIssueIndexIfNotExist();
    }

    public <T extends IndexableDocument> void bulkInsertDocuments(
        @NotEmpty List<T> documents,
        @NotBlank String indexName,
        @NotBlank String requestId
    ) throws IOException {
        List<BulkOperation> operations = new ArrayList<>();

        for (T doc : documents) {
            operations.add(new BulkOperation.Builder()
                .index(i -> i.id(doc.url()).document(doc))
                .build()
            );
        }

        int BATCH_SIZE = 500;

        for (int i = 0; i < operations.size(); i += BATCH_SIZE) {
            List<BulkOperation> batch = operations.subList(i, Math.min(i + BATCH_SIZE, operations.size()));

            BulkResponse bulkRes = openSearchClient.bulk(r -> r
                .index(indexName)
                .operations(batch)
            );

            log.debug("bulk added {} documents",
                batch.size(),
                kv("index", indexName),
                kv("requestId", requestId));

            if (bulkRes.errors()) {
                List<String> failures = bulkRes.items().stream()
                    .filter(item -> item.error() != null)
                    .map(item -> {
                        var error = item.error();
                        String reason = error != null ? error.reason() : null;
                        return item.id() + ": " + (reason != null ? reason : "unknown error");
                    })
                    .toList();

                log.error("failed to bulk insert documents",
                    kv("index", indexName),
                    kv("error", failures),
                    kv("requestId", requestId));

                throw new RuntimeException("Bulk index had failures: " + failures);
            }
        }
    }

    // dim number from: https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2
    private static final Integer embeddingDim = 384;

    /**
     * stores dependency version changelog metadata. One entry per changelog change
     * @throws IOException
     */
    private void createChangeLogIndexIfNotExist() throws IOException {
        boolean exists = openSearchClient.indices().exists(r -> r.index(changeLogIndexName)).value();

        if (!exists) {
            openSearchClient.indices().create(r -> r
                .index(changeLogIndexName)
                .settings(s -> s.knn(true))
                .mappings(m -> m
                    .properties("dependency_name", p -> p.keyword(k -> k))
                    .properties("version", p -> p.keyword(k -> k))
                    .properties("changes", p -> p.text(t -> t))
                    .properties("url", p -> p.keyword(k -> k))
                    .properties("changeEmbedding", p -> p.knnVector(k -> k
                        .dimension(embeddingDim)
                        .method(met -> met
                            .name("hnsw")
                            .spaceType("cosinesimil")
                            .engine("lucene")
                        )
                    ))
                )
            );

            log.info("created index", kv("index", changeLogIndexName));
        }
    }

    /**
     * stores dependency version issue metadata. Only issues with labels containing bug, breaking-change, breaking
     * deprecation, deprecated, regression are indexed to reduce noise
     * @throws IOException
     */
    private void createIssueIndexIfNotExist() throws IOException {
        boolean exists = openSearchClient.indices().exists(r -> r.index(issuesIndexName)).value();

        if (!exists) {
            openSearchClient.indices().create(r -> r
                .index(issuesIndexName)
                .settings(s -> s.knn(true))
                .mappings(m -> m
                    .properties("dependency_name", p -> p.keyword(k -> k))
                    .properties("version", p -> p.keyword(k -> k))
                    .properties("title", p -> p.text(t -> t))
                    .properties("body", p -> p.text(t -> t))
                    .properties("url", p -> p.keyword(k -> k))
                    .properties("labelList", p -> p.keyword(k -> k))
                    .properties("createdOn", p -> p.date(d -> d))
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

            log.info("created index", kv("index", issuesIndexName));
        }
    }
}
