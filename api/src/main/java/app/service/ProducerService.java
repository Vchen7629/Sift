package app.service;

import java.io.IOException;

import org.springframework.stereotype.Service;
import org.springframework.validation.annotation.Validated;

import com.fasterxml.jackson.databind.ObjectMapper;

import io.nats.client.JetStream;
import io.nats.client.JetStreamApiException;
import io.nats.client.api.PublishAck;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;

@Service
@Validated
public class ProducerService {
    private final JetStream js;
    private final ObjectMapper objectMapper;

    public ProducerService(JetStream js, ObjectMapper objectMapper) {
        this.js = js;
        this.objectMapper = objectMapper;
    }

    public record RepoIndexMsg(@NotBlank String repoName, @NotBlank String requestId) {};

    public void PublishIndexRepoJobRequest(@Valid RepoIndexMsg repoIndexMsg) throws JetStreamApiException, IOException {
        byte[] data = objectMapper.writeValueAsBytes(repoIndexMsg);

        PublishAck ack = js.publish("index-repo.subject.request", data);
        if (ack.hasError()) {
            throw new RuntimeException("Failed to publish: " + ack.getError());
        }
    }
}
