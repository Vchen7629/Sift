package app.repository;

import java.io.IOException;

import org.opensearch.client.opensearch.OpenSearchClient;
import org.opensearch.client.opensearch.core.GetResponse;
import org.springframework.stereotype.Repository;
import org.springframework.validation.annotation.Validated;

import io.micrometer.observation.annotation.Observed;
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

    private record JobStatus(@NotBlank String status) {}

    @Observed(name="jobstatus.findstatus.repository")
    public String findStatus(
        @NotBlank String userId,
        @NotBlank String repoName
    ) throws IOException {
        GetResponse<JobStatus> getRes = openSearchClient.get(r -> r
            .index(jobStatusIndexName)
            .id(userId + ":" + repoName)
            , JobStatus.class
        );

        if (!getRes.found()) {
            return null;
        }

        JobStatus dbRes = getRes.source();
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
                    .properties("id", p -> p.keyword(k -> k))
                    .properties("status", p -> p.keyword(k -> k))
                )
            );
        }
    }
}
