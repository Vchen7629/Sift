package app.service;

import java.io.IOException;
import java.time.Duration;
import java.util.List;
import java.util.Map;

import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import org.springframework.validation.annotation.Validated;

import com.fasterxml.jackson.core.exc.StreamReadException;
import com.fasterxml.jackson.databind.DatabindException;
import com.fasterxml.jackson.databind.ObjectMapper;

import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;
import lombok.extern.slf4j.Slf4j;
import static net.logstash.logback.argument.StructuredArguments.kv;
import ai.djl.translate.TranslateException;
import app.component.parser.DependencyParserStrategy.Dependency;
import app.dto.JobStatusDocument;
import app.repository.JobStatusRepository;
import app.service.github.DependencyService;
import app.service.processing.RepoIndexingService;
import io.nats.client.JetStreamApiException;
import io.nats.client.JetStreamSubscription;
import io.nats.client.Message;

@Service
@Validated
@Slf4j
public class ConsumerService {
    private final DependencyService dependencyService;
    private final RepoIndexingService repoIndexingService;
    private final JobStatusRepository jobStatusRepository;
    private final JetStreamSubscription jetStreamSubscription; 
    private final ObjectMapper objectMapper;
    private final io.micrometer.tracing.Tracer tracer;
    private final io.micrometer.tracing.propagation.Propagator propagator;

    public ConsumerService(
        DependencyService dependencyService,
        RepoIndexingService repoIndexingService,
        JobStatusRepository jobStatusRepository,
        JetStreamSubscription jetStreamSubscription,
        ObjectMapper objectMapper,
        io.micrometer.tracing.Tracer tracer,
        io.micrometer.tracing.propagation.Propagator propagator
    ) {
        this.dependencyService = dependencyService;
        this.repoIndexingService = repoIndexingService;
        this.jobStatusRepository = jobStatusRepository;
        this.jetStreamSubscription = jetStreamSubscription;
        this.objectMapper = objectMapper;
        this.tracer = tracer;
        this.propagator = propagator;
    }

    private record RepoIndexMsg(@NotBlank String userId, @NotBlank String repoName) {};
    private static final Duration MAX_FETCH_WAIT = Duration.ofSeconds(5);
    private static final int POLL_DELAY_MS = 100;

    @Scheduled(fixedDelay = POLL_DELAY_MS)
    public void pollIndexJobs() throws StreamReadException, DatabindException, IOException {
        List<Message> messages = jetStreamSubscription.fetch(10, MAX_FETCH_WAIT);

        for (Message msg : messages) {
            io.micrometer.tracing.Span span = propagator
                .extract(msg.getHeaders(), (carrier, key) -> carrier == null ? null : carrier.getFirst(key))
                .name("consumer.pollindexjobs.service")
                .kind(io.micrometer.tracing.Span.Kind.CONSUMER)
                .start();

            try (io.micrometer.tracing.Tracer.SpanInScope scope = tracer.withSpan(span)) {
                RepoIndexMsg payload = objectMapper.readValue(msg.getData(), RepoIndexMsg.class);
                log.debug("recived index repo msg for processing");

                try {
                    jobStatusRepository.upsert(new JobStatusDocument(payload.userId, payload.repoName, "processing:created_job"));

                    Map<String, List<Dependency>> dependenciesByLanguage = fetchAllRepoDependencies(msg, payload);
                    if (dependenciesByLanguage.isEmpty()) continue;

                    repoIndexingService.processAll(dependenciesByLanguage, payload.repoName, payload.userId);

                    msg.ack();
                    jobStatusRepository.upsert(new JobStatusDocument(payload.userId, payload.repoName, "processed"));

                    log.debug("fully processed all issues", kv("repoName", payload.repoName), kv("userId", payload.userId));
                } catch (Exception e) {
                    log.error("failed to process index repo msg", e);
                    jobStatusRepository.upsert(new JobStatusDocument(payload.userId, payload.repoName, "failed"));
                    msg.nak();
                }
            } finally {
                span.end();
            }
        }
    }

    private Map<String, List<Dependency>> fetchAllRepoDependencies(
        Message msg, @Valid RepoIndexMsg payload
    ) throws JetStreamApiException, TranslateException, IOException {
        log.info("processMsg called");

        Map<String, List<Dependency>> dependenciesByLanguage = dependencyService.fetchRepoDependencies(
            payload.repoName, payload.userId
        ).join();

        if (dependenciesByLanguage.isEmpty()) {
            jobStatusRepository.upsert(new JobStatusDocument(payload.userId, payload.repoName, "skipped:no dependencies found"));
                
            log.warn("no dependencies found for the repo", kv("repoName", payload.repoName), kv("userId", payload.userId));
            msg.ack();
            return Map.of();
        }

        return dependenciesByLanguage;
    }
}
