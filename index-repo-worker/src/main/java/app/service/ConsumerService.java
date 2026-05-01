package app.service;

import java.io.IOException;
import java.time.Duration;
import java.util.List;

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
import app.repository.OpenSearchRepository;
import io.nats.client.JetStreamApiException;
import io.nats.client.JetStreamSubscription;
import io.nats.client.Message;

@Service
@Validated
@Slf4j
public class ConsumerService {
    private final GithubApiService githubApiService;
    private final TextEmbeddingService textEmbeddingService;
    private final OpenSearchRepository openSearchRepository;
    private final JetStreamSubscription jetStreamSubscription; 
    private final ObjectMapper objectMapper;

    public ConsumerService(
        GithubApiService githubApiService,
        TextEmbeddingService textEmbeddingService,
        OpenSearchRepository openSearchRepository,
        JetStreamSubscription jetStreamSubscription,
        ObjectMapper objectMapper
    ) {
        this.githubApiService = githubApiService;
        this.textEmbeddingService = textEmbeddingService;
        this.openSearchRepository = openSearchRepository;
        this.jetStreamSubscription = jetStreamSubscription;
        this.objectMapper = objectMapper;
    }

    private record RepoIndexMsg(@NotBlank String repoName, @NotBlank String requestId) {};
    private static final Duration MAX_FETCH_WAIT = Duration.ofSeconds(5);
    private static final int POLL_DELAY_MS = 100;

    @Scheduled(fixedDelay = POLL_DELAY_MS)
    public void pollIndexJobs() throws StreamReadException, DatabindException, IOException {
        List<Message> messages = jetStreamSubscription.fetch(10, MAX_FETCH_WAIT);

        for (Message msg : messages) {
            RepoIndexMsg payload = objectMapper.readValue(msg.getData(), RepoIndexMsg.class);
            log.debug("recived index repo msg for processing", kv("requestId", payload.requestId));

            try {
                processMsg(msg, payload);
            } catch (Exception e) {
                log.error("failed to process index repo msg", kv("requestId", payload.requestId), e);
                openSearchRepository.upsertJobStatus(new OpenSearchRepository.JobStatus(payload.repoName, "failed"), payload.requestId);
                msg.nak();
            }
        }
    }

    private void processMsg(Message msg, @Valid RepoIndexMsg payload) throws JetStreamApiException, TranslateException, IOException {
        List<GithubApiService.IssueDocument> githubIssues = githubApiService.fetchRepoIssues(payload.repoName).join();

        if (githubIssues.isEmpty()) {
            openSearchRepository.upsertJobStatus(
                new OpenSearchRepository.JobStatus(payload.repoName, "Issues Not Found"), 
                payload.requestId);
                
            log.warn("no issues found for the repo: {}", payload.repoName, kv("requestId", payload.requestId));
            msg.ack();
            return;
        }

        openSearchRepository.upsertJobStatus(
            new OpenSearchRepository.JobStatus(payload.repoName, "processing"),
            payload.requestId);
        
        List<TextEmbeddingService.embeddingDocument> embeddings = textEmbeddingService.generateEmbeddings(githubIssues, payload.requestId);

        openSearchRepository.indexGithubIssue(embeddings, payload.requestId);

        msg.ack();
        openSearchRepository.upsertJobStatus(
            new OpenSearchRepository.JobStatus(payload.repoName, "processed"), 
            payload.requestId);

        log.debug("fully processed all issues for repo: {}", payload.repoName);
    }
}
