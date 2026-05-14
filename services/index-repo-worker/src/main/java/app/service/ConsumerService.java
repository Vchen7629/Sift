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
import app.dto.DependencyDocument;
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
import io.micrometer.observation.Observation;
import io.micrometer.observation.ObservationRegistry;
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
    private final ObservationRegistry observationRegistry;
    private final io.micrometer.tracing.Tracer tracer;
    private final io.micrometer.tracing.propagation.Propagator propagator;

    public ConsumerService(
        DependencyService dependencyService,
        ChangelogService changelogService,
        IssueService issueService,
        TextEmbeddingService textEmbeddingService,
        JobStatusRepository jobStatusRepository,
        DependencyRepository dependencyRepository,
        UserRepoRepository userRepoRepository,
        JetStreamSubscription jetStreamSubscription,
        ObjectMapper objectMapper,
        ObservationRegistry observationRegistry,
        io.micrometer.tracing.Tracer tracer,
        io.micrometer.tracing.propagation.Propagator propagator
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
        this.observationRegistry = observationRegistry;
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

                    processRepoDependencies(dependenciesByLanguage, payload.repoName, payload.userId);

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

    private void processRepoDependencies(
        Map<String, List<Dependency>> dependenciesByLanguage,
        @NotBlank String repoName,
        @NotBlank String userId
    ) throws TranslateException, IOException, InterruptedException, ExecutionException {
        List<ProcessedGithubIssue> issueList = new ArrayList<>();
        List<GithubChangeLogResponse> changeLogs = new ArrayList<>();
        Map<String, String> libraryMap = new HashMap<>();

        Set<DependencyDocument> indexedDependencies = dependencyRepository.list();

        for (Map.Entry<String, List<Dependency>> entry : dependenciesByLanguage.entrySet()) {
            List<Dependency> dependencies = entry.getValue();

            // the amount of threads to spawn in for fetching issues, one per repo
            Set<String> uniqueRepos = dependencies.stream()
                .map(Dependency -> Dependency.repoName())
                .collect(Collectors.toSet());

            // dependency names whose issues are already fetched
            Set<String> indexedDepNames = indexedDependencies.stream()
                .map(DependencyDocument::dependencyName)
                .collect(Collectors.toSet());

            io.micrometer.tracing.Span parentSpan = tracer.currentSpan();

            try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
                List<Future<List<ProcessedGithubIssue>>> futures = uniqueRepos.stream()
                    .filter(dep -> !indexedDepNames.contains(dep))
                    .map(dependencyName -> executor.submit(() -> {
                        try (io.micrometer.tracing.Tracer.SpanInScope scope = tracer.withSpan(parentSpan)) {
                            return issueService.fetch(dependencyName);
                        }
                    }))
                    .toList();
                
                for (Future<List<ProcessedGithubIssue>> future : futures) {
                    issueList.addAll(future.get());
                }
            }

            log.debug("fetched all {} issue chunks", issueList.size());

            // name+version pairs already indexed (to avoid false positives from independent sets)
            Set<String> indexedDepNameVersionPairs = indexedDependencies.stream()
                .map(d -> d.dependencyName() + "@" + d.version())
                .collect(Collectors.toSet());
            
            for (Dependency dependency : dependencies) {
                if (indexedDepNameVersionPairs.contains(dependency.repoName() + "@" + dependency.version())) {
                    log.debug("changelog already indexed, skipping...", 
                        kv("name", dependency.name()),
                        kv("version", dependency.version()));

                    libraryMap.put(dependency.name(), dependency.version());

                    continue;
                }

                GithubChangeLogResponse changeLog = changelogService.fetchForVersion(
                    dependency.repoName(), dependency.version()
                ).join();

                if (!changeLog.changes().equals("no-release")) {
                    changeLogs.add(changeLog);
                }
                
                libraryMap.put(dependency.name(), dependency.version());
            }
        }

        jobStatusRepository.upsert(new JobStatusDocument(userId, repoName, "processing:fetched_all_issues_changelogs"));
        if (!issueList.isEmpty()) {
            List<IndexableDocuments.Issue> issueDocuments = textEmbeddingService.generateIssue(issueList);
            dependencyRepository.bulkInsertDocuments(issueDocuments, DependencyRepository.issuesIndexName);

            jobStatusRepository.upsert(new JobStatusDocument(userId, repoName, "processing:inserted_all_issues"));
            log.info("inserted new issue documents into openSearch successfully!", kv("repoName", repoName));
        }

        if (!changeLogs.isEmpty()) {
            List<IndexableDocuments.ChangeLog> changeLogDocuments = textEmbeddingService.generateChangeLog(changeLogs);
            dependencyRepository.bulkInsertDocuments(changeLogDocuments, DependencyRepository.changeLogIndexName);

            jobStatusRepository.upsert(new JobStatusDocument(userId, repoName, "processing:inserted_all_changelogs"));
            log.info("inserted new changelog documents into openSearch successfully!", kv("repoName", repoName));
        }
        
        userRepoRepository.insert(new UserRepoDocument(userId, repoName, libraryMap, Instant.now()));
        jobStatusRepository.upsert(new JobStatusDocument(userId, repoName, "processing:inserted_indexed_repo"));
    }
}
