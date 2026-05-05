package app.repository;

import java.io.IOException;

import org.opensearch.client.opensearch.OpenSearchClient;
import org.springframework.stereotype.Repository;
import org.springframework.validation.annotation.Validated;

import app.dto.JobStatusDocument;
import io.micrometer.observation.annotation.Observed;

import static net.logstash.logback.argument.StructuredArguments.kv;

import jakarta.validation.Valid;
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

    @Observed(name="jobstatus.upsert.repository")
    public void upsert(@Valid JobStatusDocument jobStatus) {
        try {
            openSearchClient.update(r -> r
                .index(jobStatusIndexName)
                .id(jobStatus.repoName())
                .doc(new JobStatusDocument(jobStatus.repoName(), jobStatus.status()))
                .docAsUpsert(true)
                , JobStatusDocument.class
            );
            log.debug("upserted job status to {}", 
                jobStatus.status(), 
                kv("index", jobStatusIndexName),
                kv("repoName", jobStatus.repoName()));

        } catch (IOException e) { 
            log.error("failed to upsert job status to {}", 
                jobStatus.status(), 
                kv("index", jobStatusIndexName),
                kv("repoName", jobStatus.repoName()));
        }
    }

}
