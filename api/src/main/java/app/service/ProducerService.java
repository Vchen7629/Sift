package app.service;

import java.io.IOException;

import org.springframework.stereotype.Service;
import org.springframework.validation.annotation.Validated;

import com.fasterxml.jackson.databind.ObjectMapper;

import io.nats.client.JetStream;
import io.nats.client.JetStreamApiException;
import io.nats.client.Message;
import io.nats.client.api.PublishAck;
import io.nats.client.impl.Headers;
import io.nats.client.impl.NatsMessage;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import lombok.extern.slf4j.Slf4j;
import static net.logstash.logback.argument.StructuredArguments.kv;


@Service
@Validated
@Slf4j
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

        Message msg = indexRepoJobReqMsg(repoIndexMsg, data);

        PublishAck ack = js.publish(msg);
        if (ack.hasError()) {
            log.error("failed to pub index repo job request to jetstream", 
                kv("repoName", repoIndexMsg.repoName),
                kv("requestId", repoIndexMsg.requestId),
                kv("error", ack.getError()));

            throw new RuntimeException("Failed to publish: " + ack.getError());
        }

        log.debug("published index repo job request", 
            kv("repoName", repoIndexMsg.repoName),
            kv("requestId", repoIndexMsg.requestId));
    }

    // build the index repo job request nats message with headers so it doesnt publish duplicate messages
    // using the header (repo name)
    private static Message indexRepoJobReqMsg(@Valid RepoIndexMsg repoIndexMsg, @NotEmpty byte[] data) {
        Headers headers = new Headers();
        headers.add("Nats-Msg-Id", repoIndexMsg.repoName());

        return NatsMessage.builder()
            .subject("index-repo.subject.request")
            .headers(headers)
            .data(data)
            .build();
    }
}
