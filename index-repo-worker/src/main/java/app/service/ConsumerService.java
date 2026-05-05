package app.service;

import java.io.IOException;
import java.time.Duration;
import java.time.Instant;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;
import java.util.stream.Collectors;

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
import app.dto.GithubChangeLogResponse;
import app.dto.IndexableDocuments;
import app.dto.JobStatusDocument;
import app.dto.ProcessedGithubIssue;
import app.dto.UserRepoDocument;
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
                    new JobStatusDocument(payload.repoName, "processed"), 
                    payload.requestId);

                log.debug("fully processed all issues for repo: {}", payload.repoName);
            } catch (Exception e) {
                log.error("failed to process index repo msg", kv("requestId", payload.requestId), e);
                jobStatusRepository.upsertJobStatus(new JobStatusDocument(payload.repoName, "failed"), payload.requestId);
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
                new JobStatusDocument(payload.repoName, "Dependencies Not Found"), 
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
    ) throws TranslateException, IOException, InterruptedException, ExecutionException {
        List<ProcessedGithubIssue> issueList = new ArrayList<>();
        List<GithubChangeLogResponse> changeLogs = new ArrayList<>();
        Map<String, String> libraryMap = new HashMap<>();

        jobStatusRepository.upsertJobStatus(
            new JobStatusDocument(repoName, "processing"),
            requestId
        );

        Set<DependencyRepository.DependencyNameVersion> indexedDependencies = dependencyRepository.listDependencyNameVersion(requestId);

        for (Map.Entry<String, List<Dependency>> entry : dependenciesByLanguage.entrySet()) {
            List<Dependency> dependencies = entry.getValue();

            long start = System.currentTimeMillis();

            // the amount of threads to spawn in for fetching issues, one per repo
            Set<String> uniqueRepos = dependencies.stream()
                .map(Dependency -> Dependency.repoName())
                .collect(Collectors.toSet());

            // dependency names whose issues are already fetched
            Set<String> indexedDepNames = indexedDependencies.stream()
                .map(DependencyRepository.DependencyNameVersion::dependencyName)
                .collect(Collectors.toSet());

            try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
                List<Future<List<ProcessedGithubIssue>>> futures = uniqueRepos.stream()
                    .filter(dep -> !indexedDepNames.contains(dep))
                    .map(dependencyName -> executor.submit(() -> 
                        issueService.fetchDependencyIssues(dependencyName, requestId).join()))
                    .toList();
                
                for (Future<List<ProcessedGithubIssue>> future : futures) {
                    issueList.addAll(future.get());
                }
            }

            long elapsed = System.currentTimeMillis() - start;
            log.debug("fetched all {} issue chunks in {}ms ({}s)", 
                issueList.size(), elapsed, elapsed / 1000.0,
                kv("requestId", requestId));
            
            // dependency versions whose changelogs are already fetched
            Set<String> indexedDepVersions = indexedDependencies.stream()
                .map(DependencyRepository.DependencyNameVersion::version)
                .collect(Collectors.toSet());

            for (Dependency dependency : dependencies) {
                if (indexedDepNames.contains(dependency.repoName()) && indexedDepVersions.contains(dependency.version())) {
                    log.debug("changelog already indexed, skipping...", 
                        kv("name", dependency.name()),
                        kv("version", dependency.version()),
                        kv("requestId", requestId));

                    libraryMap.put(dependency.name(), dependency.version());

                    continue;
                }

                GithubChangeLogResponse changeLog = changelogService.fetchChangeLogForVersion(
                    dependency.repoName(), dependency.version()
                ).join();

                if (!changeLog.changes().equals("no-release")) {
                    changeLogs.add(changeLog);
                }
                
                libraryMap.put(dependency.name(), dependency.version());
            }
        }

        if (!issueList.isEmpty()) {
            List<IndexableDocuments.Issue> issueDocuments = textEmbeddingService.generateIssueEmbeddings(issueList, requestId);
            dependencyRepository.bulkInsertDocuments(issueDocuments, DependencyRepository.issuesIndexName, requestId);

            log.info("inserted new issue documents into openSearch successfully!", 
                kv("requestId", requestId), kv("repoName", repoName)
            );
        }

        if (!changeLogs.isEmpty()) {
            List<IndexableDocuments.ChangeLog> changeLogDocuments = textEmbeddingService.generateChangeLogEmbeddings(changeLogs, requestId);
            dependencyRepository.bulkInsertDocuments(changeLogDocuments, DependencyRepository.changeLogIndexName, requestId);

            log.info("inserted new changelog documents into openSearch successfully!", 
                kv("requestId", requestId), kv("repoName", repoName)
            );
        }

        userRepoRepository.insertDocument(new UserRepoDocument("test-user-1", repoName, libraryMap, Instant.now()));
    }
}
