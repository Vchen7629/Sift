package app.repository;

import java.io.IOException;

import org.opensearch.client.opensearch.OpenSearchClient;
import org.opensearch.client.opensearch.core.SearchResponse;
import org.springframework.stereotype.Repository;
import org.springframework.validation.annotation.Validated;

import jakarta.annotation.PostConstruct;
import jakarta.validation.constraints.NotBlank;
import lombok.extern.slf4j.Slf4j;

@Repository
@Validated
@Slf4j
public class JobStatusRepository {
    private final OpenSearchClient openSearchClient;
    private final static String jobStatusIndexName = "job-status";

    public JobStatusRepository(OpenSearchClient openSearchClient) {
        this.openSearchClient = openSearchClient;
    }

    @PostConstruct
    private void init() throws IOException {
        createIndexIfNotExist();
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
}
