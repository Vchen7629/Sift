package app.service.processing;

import java.io.IOException;
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

import org.springframework.stereotype.Service;
import org.springframework.validation.annotation.Validated;

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
import app.service.github.ChangelogService;
import app.service.github.IssueService;
import io.micrometer.observation.annotation.Observed;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import lombok.extern.slf4j.Slf4j;
import static net.logstash.logback.argument.StructuredArguments.kv;


@Service
@Slf4j
@Validated
public class RepoIndexingService {
    private final DependencyRepository dependencyRepository;
    private final JobStatusRepository jobStatusRepository;
    private final UserRepoRepository userRepoRepository;
    private final IssueService issueService;
    private final ChangelogService changelogService;
    private final TextEmbeddingService textEmbeddingService;
    private final io.micrometer.tracing.Tracer tracer;

    public RepoIndexingService(
        DependencyRepository dependencyRepository,
        JobStatusRepository jobStatusRepository,
        UserRepoRepository userRepoRepository,
        IssueService issueService,
        ChangelogService changeLogService,
        TextEmbeddingService textEmbeddingService,
        io.micrometer.tracing.Tracer tracer
    ) {
        this.dependencyRepository = dependencyRepository;
        this.jobStatusRepository = jobStatusRepository;
        this.userRepoRepository = userRepoRepository;
        this.issueService = issueService;
        this.changelogService = changeLogService;
        this.textEmbeddingService = textEmbeddingService;
        this.tracer = tracer;
    }

    @Observed(name = "repoindexing.processall.service")
    public void processAll(
        @NotEmpty Map<String, List<Dependency>> dependenciesByLanguage, 
        @NotBlank String repoName, @NotBlank String userId
    ) throws TranslateException, IOException, InterruptedException, ExecutionException {
        Set<DependencyDocument> indexedDependencies = dependencyRepository.list();
    
        List<ProcessedGithubIssue> issueList = fetchDependencyIssues(dependenciesByLanguage, indexedDependencies);
        ChangeLogResult changeLogResult = fetchDependencyChangelogs(dependenciesByLanguage, indexedDependencies);
    
        EmbeddingRes embeddings = createEmbeddings(issueList, changeLogResult, userId, repoName);
        upsertIssueChangelogs(embeddings.issueDocuments(), embeddings.changeLogDocuments(), userId, repoName);

        userRepoRepository.insert(new UserRepoDocument(userId, repoName, changeLogResult.libraryMap(), Instant.now()));
        jobStatusRepository.upsert(new JobStatusDocument(userId, repoName, "processing:inserted_indexed_repo"));
    }

    List<ProcessedGithubIssue> fetchDependencyIssues(
        Map<String, List<Dependency>> dependenciesByLanguage,
        Set<DependencyDocument> indexedDependencies
    ) throws InterruptedException, ExecutionException {
        // dependency names whose issues are already fetched
        Set<String> indexedDepNames = indexedDependencies.stream()
            .map(DependencyDocument::dependencyName)
            .collect(Collectors.toSet());

        List<ProcessedGithubIssue> issueList = new ArrayList<>();
        io.micrometer.tracing.Span parentSpan = tracer.currentSpan();

        for (List<Dependency> dependencies : dependenciesByLanguage.values()) {
            // the amount of threads to spawn in for fetching issues, one per repo
            Set<String> uniqueRepos = dependencies.stream()
                .map(Dependency -> Dependency.repoName())
                .collect(Collectors.toSet());

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
        }

        log.debug("fetched all {} issue chunks", issueList.size());
        return issueList;
    }

    private record ChangeLogResult(List<GithubChangeLogResponse> changelogs, Map<String, String> libraryMap) {}

    ChangeLogResult fetchDependencyChangelogs(
        Map<String, List<Dependency>> dependenciesByLanguage,
        Set<DependencyDocument> indexedDependencies
    ) throws InterruptedException, ExecutionException {
        // name+version pairs already indexed (to avoid false positives from independent sets)
        Set<String> indexedDepNameVersionPairs = indexedDependencies.stream()
            .map(d -> d.dependencyName() + "@" + d.version())
            .collect(Collectors.toSet());

        List<GithubChangeLogResponse> changeLogs = new ArrayList<>();
        Map<String, String> libraryMap = new HashMap<>();

        for (List<Dependency> dependencies : dependenciesByLanguage.values()) {
            List<Dependency> toFetch = new ArrayList<>();

            // build the list of changelogs to process (not processed before) before using virtual threads to 
            // prevent concurrency issues
            for (Dependency dependency : dependencies) {
                libraryMap.put(dependency.name(), dependency.version());
                if (indexedDepNameVersionPairs.contains(dependency.repoName() + "@" + dependency.version())) {
                    log.debug("changelog already indexed, skipping...", 
                        kv("name", dependency.name()),
                        kv("version", dependency.version()));
                } else {
                    toFetch.add(dependency);
                }
            }

            io.micrometer.tracing.Span parentSpan = tracer.currentSpan();
            try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
                List<Future<GithubChangeLogResponse>> futures = toFetch.stream()
                    .map(dep -> executor.submit(() -> {
                        try (io.micrometer.tracing.Tracer.SpanInScope scope = tracer.withSpan(parentSpan)) {
                            return changelogService.fetchForVersion(dep.repoName(), dep.version());
                        }
                    }))
                    .toList();

                for (var future : futures) {
                    GithubChangeLogResponse changeLog = future.get();
                    if (!changeLog.changes().equals("no-release")) {
                        changeLogs.add(changeLog);
                    }
                }
            }
        }

        return new ChangeLogResult(changeLogs, libraryMap);
    }

    private record EmbeddingRes(
        List<IndexableDocuments.Issue> issueDocuments, 
        List<IndexableDocuments.ChangeLog> changeLogDocuments
    ) {}

    EmbeddingRes createEmbeddings(
        List<ProcessedGithubIssue> issueList, 
        ChangeLogResult changeLogResult,
        String userId, String repoName
    ) throws TranslateException, IOException {
        jobStatusRepository.upsert(new JobStatusDocument(userId, repoName, "processing:fetched_all_issues_changelogs"));

        var issueDocuments = issueList.isEmpty() 
            ? List.<IndexableDocuments.Issue>of() 
            : textEmbeddingService.githubIssue(issueList);
        var changeLogDocuments = changeLogResult.changelogs().isEmpty()
            ? List.<IndexableDocuments.ChangeLog>of()
            : textEmbeddingService.githubChangelog(changeLogResult.changelogs());

        return new EmbeddingRes(issueDocuments, changeLogDocuments);
    }

    void upsertIssueChangelogs(
        List<IndexableDocuments.Issue> issueEmbeddings,
        List<IndexableDocuments.ChangeLog> changeLogEmbeddings,
        String userId, String repoName
    ) throws TranslateException, InterruptedException, ExecutionException {
        int BATCH_SIZE = 750;

        jobStatusRepository.upsert(new JobStatusDocument(userId, repoName, "processing:created_embeddings"));

        io.micrometer.tracing.Span parentSpan = tracer.currentSpan();
        try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
            List<Future<?>> futures = new ArrayList<>();

            for (int i = 0; i < issueEmbeddings.size(); i += BATCH_SIZE) {
                // creating a copy of the sublist for thread safety
                List<IndexableDocuments.Issue> batch = new ArrayList<>(
                    issueEmbeddings.subList(i, Math.min(i + BATCH_SIZE, issueEmbeddings.size()))
                );
                futures.add(executor.submit(() -> insertBatch(batch, DependencyRepository.issuesIndexName, parentSpan)));
            }

            for (int i = 0; i < changeLogEmbeddings.size(); i += BATCH_SIZE) {
                List<IndexableDocuments.ChangeLog> batch = new ArrayList<>(
                    changeLogEmbeddings.subList(i, Math.min(i + BATCH_SIZE, changeLogEmbeddings.size()))
                );
                futures.add(executor.submit(() -> insertBatch(batch, DependencyRepository.changeLogIndexName, parentSpan)));
            }

            for (Future<?> future : futures) {
                future.get();
            }

            if (!issueEmbeddings.isEmpty()) {
                jobStatusRepository.upsert(new JobStatusDocument(userId, repoName, "processing:inserted_all_issues"));
                log.info("inserted new issue documents into openSearch successfully!", kv("repoName", repoName));
            }

            if (!changeLogEmbeddings.isEmpty()) {
                jobStatusRepository.upsert(new JobStatusDocument(userId, repoName, "processing:inserted_all_changelogs"));
                log.info("inserted new changelog documents into openSearch successfully!", kv("repoName", repoName));
            }
        }
    }

    private <T extends IndexableDocuments.Base> void insertBatch(
        List<T> batch, String indexName, io.micrometer.tracing.Span parentSpan
    ) {
        try (io.micrometer.tracing.Tracer.SpanInScope scope = tracer.withSpan(parentSpan)) {
            try {
                dependencyRepository.bulkInsertDocuments(batch, indexName);
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        }
    }
}
