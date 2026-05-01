package app.config;

import java.io.IOException;
import java.time.Duration;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import io.nats.client.Connection;
import io.nats.client.ErrorListener;
import io.nats.client.JetStream;
import io.nats.client.JetStreamOptions;
import io.nats.client.Nats;
import io.nats.client.Options;

@Configuration
public class NatsConfig {
    @Value("${nats.connection_url}")
    private String connURL;

    private static final int MAX_RECONNECTS = -1;
    private static final Duration RECONNECT_WAIT = Duration.ofSeconds(2);

    @Bean(destroyMethod = "close")
    public Connection natsConnection() throws IOException, InterruptedException {
        Options options = new Options.Builder()
            .server(connURL)
            .maxReconnects(MAX_RECONNECTS)
            .reconnectWait(RECONNECT_WAIT)
            .errorListener(new ErrorListener() {}) // todo: implement a error listener
            .build();

        return Nats.connect(options);
    }

    @Bean
    public JetStream jetStream(Connection connection) throws IOException {
        JetStreamOptions options = JetStreamOptions.defaultOptions();

        return connection.jetStream(options);
    }
}
