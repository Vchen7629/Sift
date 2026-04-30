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

import ai.djl.translate.TranslateException;
import app.repository.OpenSearchRepository;
import io.nats.client.JetStreamApiException;
import io.nats.client.JetStreamSubscription;
import io.nats.client.Message;

@Service
@Validated
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
    private static final int POLL_DELAY_MS = 1000;

    @Scheduled(fixedDelay = POLL_DELAY_MS)
    public void pollIndexJobs() throws StreamReadException, DatabindException, IOException {
        List<Message> messages = jetStreamSubscription.fetch(10, MAX_FETCH_WAIT);

        for (Message msg : messages) {
            RepoIndexMsg payload = objectMapper.readValue(msg.getData(), RepoIndexMsg.class);

            try {
                processMsg(payload);

                msg.ack();
                openSearchRepository.upsertJobStatus(new OpenSearchRepository.JobStatus(payload.repoName, "processed"));
            } catch (Exception e) {
                msg.nak();
            }
        }
    }

    private void processMsg(@Valid RepoIndexMsg payload) throws JetStreamApiException, TranslateException, IOException {
        List<GithubApiService.IssueDocument> githubIssues = githubApiService.fetchRepoIssues(payload.repoName).join();

        if (githubIssues.isEmpty()) {
            openSearchRepository.upsertJobStatus(new OpenSearchRepository.JobStatus(
                payload.repoName, "Issues Not Found"
            ));
        }

        openSearchRepository.upsertJobStatus(new OpenSearchRepository.JobStatus(payload.repoName, "processing"));
        
        List<TextEmbeddingService.embeddingDocument> embeddings = textEmbeddingService.generateEmbeddings(githubIssues);

        openSearchRepository.indexGithubIssue(embeddings);
    }
}
