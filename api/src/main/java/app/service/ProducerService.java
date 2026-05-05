package app.service;

import java.io.IOException;
import java.util.Objects;

import org.springframework.stereotype.Service;
import org.springframework.validation.annotation.Validated;

import com.fasterxml.jackson.databind.ObjectMapper;

import app.dto.IndexRepoMsg;
import app.exception.JobAlreadyProcessingException;
import app.repository.JobStatusRepository;
import io.micrometer.observation.annotation.Observed;
import io.nats.client.JetStream;
import io.nats.client.JetStreamApiException;
import io.nats.client.Message;
import io.nats.client.api.PublishAck;
import io.nats.client.impl.Headers;
import io.nats.client.impl.NatsMessage;
import io.opentelemetry.api.OpenTelemetry;
import io.opentelemetry.context.Context;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.NotNull;
import lombok.extern.slf4j.Slf4j;
import static net.logstash.logback.argument.StructuredArguments.kv;


@Service
@Validated
@Slf4j
public class ProducerService {
    private final JetStream js;
    private final ObjectMapper objectMapper;
    private final OpenTelemetry openTelemetry;
    private final JobStatusRepository jobStatusRepository;

    public ProducerService(
        JetStream js, 
        ObjectMapper objectMapper, 
        OpenTelemetry openTelemetry,
        JobStatusRepository jobStatusRepository
    ) {
        this.js = js;
        this.objectMapper = objectMapper;
        this.openTelemetry = openTelemetry;
        this.jobStatusRepository = jobStatusRepository;
    }

    // todo: pass span trace id and pass it into nats message as header so it travels across service boundaries
    @Observed(name="producer.publishindexrepojobrequest.service")
    public void PublishIndexRepoJobRequest(@Valid IndexRepoMsg indexRepoMsg) throws JetStreamApiException, IOException {
        String jobStatus = jobStatusRepository.findStatus(indexRepoMsg.userId(), indexRepoMsg.repoName());
        if ("processing".equals(jobStatus)) {
            log.warn("repo is already processing, skipping...", 
                kv("userId", indexRepoMsg.userId()), kv("repoName", indexRepoMsg.repoName())
            );
            throw new JobAlreadyProcessingException();
        }

        byte[] data = objectMapper.writeValueAsBytes(indexRepoMsg);

        Headers headers = new Headers();
        Context context = Objects.requireNonNull(Context.current());
        openTelemetry.getPropagators().getTextMapPropagator().inject(
            context, 
            headers, 
            (carrier, key, value) -> { if (carrier != null) carrier.add(key, value); }
        );

        log.debug("injected trace context into nats headers", kv("headers", headers.toString()));

        Message msg = indexRepoJobReqMsg(indexRepoMsg, data, headers);

        PublishAck ack = js.publish(msg);
        if (ack.hasError()) {
            log.error("failed to pub index repo job request to jetstream", 
                kv("repoName", indexRepoMsg.repoName()),
                kv("error", ack.getError()));

            throw new RuntimeException("Failed to publish: " + ack.getError());
        }

        log.debug("published index repo job request", kv("repoName", indexRepoMsg.repoName()));
    }

    // build the index repo job request nats message with headers so it doesnt publish duplicate messages
    // using the header (repo name)
    private static Message indexRepoJobReqMsg(
        @Valid IndexRepoMsg indexRepoMsg, 
        @NotEmpty byte[] data,
        @NotNull Headers headers
    ) {
        return NatsMessage.builder()
            .subject("index-repo.subject.request")
            .headers(headers)
            .data(data)
            .build();
    }
}
