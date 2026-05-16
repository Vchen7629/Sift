package app.service.internal;

import org.opensearch.client.opensearch.OpenSearchClient;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;

import lombok.extern.slf4j.Slf4j;
import static net.logstash.logback.argument.StructuredArguments.kv;


@Service
@Slf4j
public class ConnectionWarmerService {
    private final OpenSearchClient client;

    public ConnectionWarmerService(OpenSearchClient client) {
        this.client = client;
    }

    @Scheduled(fixedDelay = 300000) // 5 minutes
    public void openSearch() {
        try {
            client.ping();
            log.debug("pinged opensearch to keep connection alive");
        } catch (Exception e) {
            log.error("failed to ping the opensearch client to keep connection warm", 
                kv("error", e.getMessage())
            );
        }
    }
}
