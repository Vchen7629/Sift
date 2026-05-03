package app.repository;

import java.io.IOException;

import org.opensearch.client.opensearch.OpenSearchClient;
import org.springframework.stereotype.Repository;
import org.springframework.validation.annotation.Validated;
import static net.logstash.logback.argument.StructuredArguments.kv;

import jakarta.validation.Valid;
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

    public record JobStatus(@NotBlank String repoName, @NotBlank String status) {}

    public void upsertJobStatus(@Valid JobStatus jobStatus, @NotBlank String requestId) {
        try {
            openSearchClient.update(r -> r
                .index(jobStatusIndexName)
                .id(jobStatus.repoName)
                .doc(new JobStatus(jobStatus.repoName, jobStatus.status))
                .docAsUpsert(true)
                , JobStatus.class
            );
            log.debug("upserted job status to {}", 
                jobStatus.status, 
                kv("index", jobStatusIndexName),
                kv("repoName", jobStatus.repoName),
                kv("requestId", requestId));

        } catch (IOException e) { 
            log.error("failed to upsert job status to {}", 
                jobStatus.status, 
                kv("index", jobStatusIndexName),
                kv("repoName", jobStatus.repoName),
                kv("requestId", requestId));
        }
    }

}
