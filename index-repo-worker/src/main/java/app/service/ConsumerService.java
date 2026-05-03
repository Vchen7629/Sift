package app.service;

import java.io.IOException;
import java.time.Duration;
import java.time.Instant;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;

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
import app.repository.DependencyRepository;
import app.repository.JobStatusRepository;
import app.repository.UserRepoRepository;
import app.service.githubRepo.ChangelogService;
import app.service.githubRepo.DependencyService;
import app.service.githubRepo.IssueService;
import io.nats.client.JetStreamApiException;
import io.nats.client.JetStreamSubscription;
import io.nats.client.Message;

@Service
@Validated
@Slf4j
public class ConsumerService {
    private final DependencyService dependencyService;
    private final ChangelogService changelogService;
    private final IssueService issueService;
    private final TextEmbeddingService textEmbeddingService;
    private final JobStatusRepository jobStatusRepository;
    private final DependencyRepository dependencyRepository;
    private final UserRepoRepository userRepoRepository;
    private final JetStreamSubscription jetStreamSubscription; 
    private final ObjectMapper objectMapper;

    public ConsumerService(
        DependencyService dependencyService,
        ChangelogService changelogService,
        IssueService issueService,
        TextEmbeddingService textEmbeddingService,
        JobStatusRepository jobStatusRepository,
        DependencyRepository dependencyRepository,
        UserRepoRepository userRepoRepository,
        JetStreamSubscription jetStreamSubscription,
        ObjectMapper objectMapper
    ) {
        this.dependencyService = dependencyService;
        this.changelogService = changelogService;
        this.issueService = issueService;
        this.textEmbeddingService = textEmbeddingService;
        this.jobStatusRepository = jobStatusRepository;
        this.dependencyRepository = dependencyRepository;
        this.userRepoRepository = userRepoRepository;
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
                Map<String, List<Dependency>> dependenciesByLanguage = fetchAllRepoDependencies(msg, payload);
                if (dependenciesByLanguage.isEmpty()) continue;

                processRepoDependencies(dependenciesByLanguage, payload.repoName, payload.requestId);

                msg.ack();
                jobStatusRepository.upsertJobStatus(
                    new JobStatusRepository.JobStatus(payload.repoName, "processed"), 
                    payload.requestId);

                log.debug("fully processed all issues for repo: {}", payload.repoName);
            } catch (Exception e) {
                log.error("failed to process index repo msg", kv("requestId", payload.requestId), e);
                jobStatusRepository.upsertJobStatus(new JobStatusRepository.JobStatus(payload.repoName, "failed"), payload.requestId);
                msg.nak();
            }
        }
    }

    private Map<String, List<Dependency>> fetchAllRepoDependencies(
        Message msg, @Valid RepoIndexMsg payload
    ) throws JetStreamApiException, TranslateException, IOException {
        log.info("processMsg called", kv("requestId", payload.requestId));

        Map<String, List<Dependency>> dependenciesByLanguage = dependencyService.fetchRepoDependencies(payload.repoName, payload.requestId).join();

        if (dependenciesByLanguage.isEmpty()) {
            jobStatusRepository.upsertJobStatus(
                new JobStatusRepository.JobStatus(payload.repoName, "Dependencies Not Found"), 
                payload.requestId);
                
            log.warn("no dependencies found for the repo: {}", payload.repoName, kv("requestId", payload.requestId));
            msg.ack();
            return Map.of();
        }

        return dependenciesByLanguage;
    }

    private void processRepoDependencies(
        Map<String, List<Dependency>> dependenciesByLanguage,
        @NotBlank String repoName,
        @NotBlank String requestId
    ) throws TranslateException, IOException {
        List<IssueService.Result> issueList = new ArrayList<>();
        List<ChangelogService.Result> changeLogs = new ArrayList<>();
        Map<String, String> libraryMap = new HashMap<>();
        Set<String> fetchedIssueRepos = new HashSet<>();

        jobStatusRepository.upsertJobStatus(
            new JobStatusRepository.JobStatus(repoName, "processing"),
            requestId
        );

        for (Map.Entry<String, List<Dependency>> entry : dependenciesByLanguage.entrySet()) {
            List<Dependency> dependencies = entry.getValue();

            for (Dependency dependency : dependencies) {
                if (fetchedIssueRepos.add(dependency.repoName())) {
                    List<IssueService.Result> repoIssues = issueService.fetchDependencyIssues(dependency.repoName(), requestId).join();
                    issueList.addAll(repoIssues);
                }

                ChangelogService.Result changeLog = changelogService.fetchChangeLogForVersion(
                    dependency.repoName(), dependency.version()
                ).join();

                if (!changeLog.changes().equals("no-release")) {
                    changeLogs.add(changeLog);
                }
                
                libraryMap.put(dependency.name(), dependency.version());
            }
        }

        List<TextEmbeddingService.IssueDocument> issueDocuments = textEmbeddingService.generateIssueEmbeddings(issueList, requestId);
        List<TextEmbeddingService.ChangeLogDocument> changeLogDocuments = textEmbeddingService.generateChangeLogEmbeddings(changeLogs, requestId);

        dependencyRepository.bulkInsertDocuments(issueDocuments, DependencyRepository.issuesIndexName, requestId);
        dependencyRepository.bulkInsertDocuments(changeLogDocuments, DependencyRepository.changeLogIndexName, requestId);
        userRepoRepository.insertDocument(new UserRepoRepository.UserRepo("test-user-1", repoName, libraryMap, Instant.now()));
    }
}
