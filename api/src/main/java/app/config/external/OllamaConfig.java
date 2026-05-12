package app.config.external;

import org.apache.hc.client5.http.impl.io.PoolingHttpClientConnectionManager;
import org.apache.hc.client5.http.impl.io.PoolingHttpClientConnectionManagerBuilder;
import org.apache.hc.client5.http.config.ConnectionConfig;
import org.apache.hc.client5.http.impl.classic.CloseableHttpClient;
import org.apache.hc.client5.http.impl.classic.HttpClients;
import org.apache.hc.core5.util.Timeout;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.client.HttpComponentsClientHttpRequestFactory;
import org.springframework.web.client.RestClient;

@Configuration
public class OllamaConfig {

    @Value("${ollama.host}")
    private String hostname;

    @Value("${ollama.port}")
    private int port;

    @Bean
    public CloseableHttpClient ollamaHttpClient() {
        PoolingHttpClientConnectionManager cm = PoolingHttpClientConnectionManagerBuilder.create()
            .setDefaultConnectionConfig(ConnectionConfig.custom()
                .setConnectTimeout(Timeout.ofSeconds(2))
                .setSocketTimeout(Timeout.ofSeconds(10))
                .build())
            .build();

        return HttpClients.custom()
            .setConnectionManager(cm)
            .build();
    }

    @Bean
    public HttpComponentsClientHttpRequestFactory ollamaRequestFactory(CloseableHttpClient ollamaHttpClient) {
        return new HttpComponentsClientHttpRequestFactory(ollamaHttpClient);
    }

    @Bean
    public RestClient ollamaRestClient(HttpComponentsClientHttpRequestFactory ollamaRequestFactory) {
        return RestClient.builder()
            .baseUrl(String.format("http://%s:%d", hostname, port))
            .defaultHeader("Accept", "application/json")
            .requestFactory(ollamaRequestFactory)
            .build();
    }
}
