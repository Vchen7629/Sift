package app.repository;

import org.springframework.stereotype.Repository;
import org.springframework.validation.annotation.Validated;

import app.dto.JobStatusDocument;
import io.lettuce.core.RedisException;
import io.lettuce.core.api.StatefulRedisConnection;
import io.lettuce.core.api.sync.RedisStringCommands;
import io.micrometer.observation.annotation.Observed;

import static net.logstash.logback.argument.StructuredArguments.kv;

import jakarta.validation.Valid;
import lombok.extern.slf4j.Slf4j;

@Repository
@Validated
@Slf4j
public class JobStatusRepository {
    private final RedisStringCommands<String, String> commands;
    private static final int JOB_STATUS_TTL_SECONDS = 3600;

    public JobStatusRepository(StatefulRedisConnection<String, String> connection) {
        this.commands = connection.sync();
    }

    @Observed(name="jobstatus.upsert.repository")
    public void upsert(@Valid JobStatusDocument jobStatus) {
        String key = "job:" + jobStatus.userId() + ":" + jobStatus.repoName();

        try {
            commands.setex(key, JOB_STATUS_TTL_SECONDS, jobStatus.status());
        } catch (RedisException e) {
            log.error("Failed to update job status", kv("key", key), e);
        }

        log.debug("upserted job status to {}", 
            jobStatus.status(), kv("userId", jobStatus.userId()), kv("repoName", jobStatus.repoName()));
    }
}
