package app.repository;

import java.io.IOException;
import java.util.List;

import org.opensearch.client.opensearch.OpenSearchClient;
import org.springframework.stereotype.Repository;

import app.service.TextEmbeddingService;
import jakarta.annotation.PostConstruct;

@Repository
public class OpenSearchRepository {
    private final OpenSearchClient openSearchClient;

    public OpenSearchRepository(OpenSearchClient openSearchClient) {
        this.openSearchClient = openSearchClient;
    }

    @PostConstruct
    private void init() throws IOException {
        createIndexIfNotExist("github-issues");
    }

    public void indexGithubIssue(List<TextEmbeddingService.embeddingDocument> issueUrlTexts) throws IOException {
        openSearchClient.bulk();
    }

    private static final Integer embeddingDim = 384;

    private void createIndexIfNotExist(String indexName) throws IOException {
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
