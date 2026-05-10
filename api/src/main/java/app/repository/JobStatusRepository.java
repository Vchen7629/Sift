package app.repository;

import static net.logstash.logback.argument.StructuredArguments.kv;

import org.springframework.stereotype.Repository;
import org.springframework.validation.annotation.Validated;

import io.lettuce.core.RedisException;
import io.lettuce.core.api.StatefulRedisConnection;
import io.lettuce.core.api.sync.RedisStringCommands;
import io.micrometer.observation.annotation.Observed;
import jakarta.validation.constraints.NotBlank;
import lombok.extern.slf4j.Slf4j;

@Repository
@Validated
@Slf4j
public class JobStatusRepository {
    private final RedisStringCommands<String, String> commands;

    public JobStatusRepository(StatefulRedisConnection<String, String> connection) {
        this.commands = connection.sync();
    }

    @Observed(name="jobstatus.findstatus.repository")
    public String findStatus(@NotBlank String userId, @NotBlank String repoName) {
        String key = "job:" + userId + ":" + repoName;
        
        try {
            return commands.get(key);
        } catch (RedisException e) {
            log.error("Failed to get job status", kv("key", key), e);
            return null;
        }
    }
}
