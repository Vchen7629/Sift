package app.repository;

import static net.logstash.logback.argument.StructuredArguments.kv;

import java.io.IOException;
import java.util.concurrent.CompletableFuture;

import org.springframework.stereotype.Repository;
import org.springframework.validation.annotation.Validated;

import io.lettuce.core.api.StatefulRedisConnection;
import io.lettuce.core.api.async.RedisAsyncCommands;
import io.micrometer.observation.annotation.Observed;
import jakarta.validation.constraints.NotBlank;
import lombok.extern.slf4j.Slf4j;

@Repository
@Validated
@Slf4j
public class JobStatusRepository {
    private final RedisAsyncCommands<String, String> asyncCommands;

    public JobStatusRepository(StatefulRedisConnection<String, String> connection) {
        this.asyncCommands = connection.async();
    }

    @Observed(name="jobstatus.findstatus.repository")
    public CompletableFuture<String> findStatus(@NotBlank String userId, @NotBlank String repoName) throws IOException {
        String key = "job:" + userId + ":" + repoName;
        
        return asyncCommands.get(key)
            .exceptionally(e -> {
                log.error("Failed to get job status", kv("key", key), e);
                return null;
            })
            .toCompletableFuture();
        }
}
