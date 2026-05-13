package app.config.repository;

import java.time.Duration;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import io.lettuce.core.ClientOptions;
import io.lettuce.core.RedisClient;
import io.lettuce.core.RedisURI;
import io.lettuce.core.SocketOptions;
import io.lettuce.core.TimeoutOptions;
import io.lettuce.core.api.StatefulRedisConnection;

@Configuration
public class ValkeyConfig {
    @Value("${redis.hostname}")
    private String hostName;

    @Value("${redis.port}")
    private int port;

    private int TIMEOUT_S = 5;
    
    // Todo: configure authentication, password, ssl/tls in prod
    @Bean
    public RedisClient client() {
        RedisURI uri = RedisURI.builder()
            .withHost(hostName)
            .withPort(port)
            .withTimeout(Duration.ofSeconds(TIMEOUT_S))
            .build();
        
        RedisClient client = RedisClient.create(uri);

        client.setOptions(ClientOptions.builder()
            .socketOptions(SocketOptions.builder()
                .connectTimeout(Duration.ofSeconds(TIMEOUT_S))
                .keepAlive(true)
                .build())
            .timeoutOptions(TimeoutOptions.enabled())
            .disconnectedBehavior(ClientOptions.DisconnectedBehavior.REJECT_COMMANDS)
            .build());

        return client;
    }

    @Bean
    public StatefulRedisConnection<String, String> connection(RedisClient client) {
        return client.connect();
    }
}
