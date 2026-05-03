package app.repository;

import java.io.IOException;
import java.time.Instant;
import java.util.Map;

import org.opensearch.client.opensearch.OpenSearchClient;
import org.opensearch.client.opensearch._types.Result;
import org.opensearch.client.opensearch.core.IndexResponse;
import org.springframework.stereotype.Repository;
import org.springframework.validation.annotation.Validated;

import jakarta.annotation.PostConstruct;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import lombok.extern.slf4j.Slf4j;
import static net.logstash.logback.argument.StructuredArguments.kv;


@Repository
@Validated
@Slf4j
public class UserRepoRepository {
    private final OpenSearchClient openSearchClient;

    private final static String indexName = "user-repo";

    public UserRepoRepository(OpenSearchClient openSearchClient) {
        this.openSearchClient = openSearchClient;
    }

    @PostConstruct
    private void init() throws IOException {
        createIndexIfNotExist();
    }

    public static record UserRepo(
        @NotBlank String userId,
        @NotBlank String repoName,
        @NotEmpty Map<String, String> dependencies,
        Instant lastIndexed
    ) {};

    public void insertDocument(
        @Valid UserRepo document
    ) throws IOException {
        IndexResponse insertRes = openSearchClient.index(r -> r
            .index(indexName).id(document.repoName).document(document)
        );

        if (insertRes.result() == Result.Created || insertRes.result() == Result.Updated) {
            log.debug("indexed repo metadata for user", kv("repoName", document.repoName));
        } else {
            throw new RuntimeException(
                "Failed to index document: " + document.repoName + ", result: " + insertRes.result()
            );
        }
    }
    
    /**
     * index unique per user, shows their indexed repo info like name, libraries/dependencies
     * and when they last indexed it
     * @throws IOException
     */
    private void createIndexIfNotExist() throws IOException {
        boolean exists = openSearchClient.indices().exists(r -> r.index(indexName)).value();

        if (!exists) {
            openSearchClient.indices().create(r -> r
                .index(indexName)
                .mappings(m -> m
                    .properties("user_id", p -> p.keyword(k -> k))
                    .properties("repo_name", p -> p.keyword(k -> k))
                    .properties("dependencies", p -> p.flatObject(f -> f))
                    .properties("last_indexed", p -> p.date(d -> d))
                )
            );
        }
    }
}
