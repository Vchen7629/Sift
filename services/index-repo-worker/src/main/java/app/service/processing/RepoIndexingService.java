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

    public void processAll(
        Map<String, List<Dependency>> dependenciesByLanguage, String repoName, String userId
    ) throws TranslateException, IOException, InterruptedException, ExecutionException {
        Set<DependencyDocument> indexedDependencies = dependencyRepository.list();
    
        List<ProcessedGithubIssue> issueList = fetchDependencyIssues(dependenciesByLanguage, indexedDependencies);
        ChangeLogResult changeLogResult = fetchDependencyChangelogs(dependenciesByLanguage, indexedDependencies);
    
        embedAndInsert(issueList, changeLogResult, userId, repoName);
        userRepoRepository.insert(new UserRepoDocument(userId, repoName, changeLogResult.libraryMap(), Instant.now()));
        jobStatusRepository.upsert(new JobStatusDocument(userId, repoName, "processing:inserted_indexed_repo"));
    }

    private List<ProcessedGithubIssue> fetchDependencyIssues(
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

    private ChangeLogResult fetchDependencyChangelogs(
        Map<String, List<Dependency>> dependenciesByLanguage,
        Set<DependencyDocument> indexedDependencies
    ) {
        // name+version pairs already indexed (to avoid false positives from independent sets)
        Set<String> indexedDepNameVersionPairs = indexedDependencies.stream()
            .map(d -> d.dependencyName() + "@" + d.version())
            .collect(Collectors.toSet());

        List<GithubChangeLogResponse> changeLogs = new ArrayList<>();
        Map<String, String> libraryMap = new HashMap<>();

        for (List<Dependency> dependencies : dependenciesByLanguage.values()) {
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

        return new ChangeLogResult(changeLogs, libraryMap);
    }

    private void embedAndInsert(
        List<ProcessedGithubIssue> issueList,
        ChangeLogResult changeLogResult,
        String userId, String repoName
    ) throws TranslateException, IOException {
        int BATCH_SIZE = 750;

        jobStatusRepository.upsert(new JobStatusDocument(userId, repoName, "processing:fetched_all_issues_changelogs"));
        
        if (!issueList.isEmpty()) {
            List<IndexableDocuments.Issue> issueDocuments = textEmbeddingService.githubIssue(issueList);
            dependencyRepository.bulkInsertDocuments(issueDocuments, DependencyRepository.issuesIndexName, BATCH_SIZE);

            jobStatusRepository.upsert(new JobStatusDocument(userId, repoName, "processing:inserted_all_issues"));
            log.info("inserted new issue documents into openSearch successfully!", kv("repoName", repoName));
        }

        if (!changeLogResult.changelogs().isEmpty()) {
            List<IndexableDocuments.ChangeLog> changeLogDocuments = textEmbeddingService.githubChangelog(changeLogResult.changelogs());
            dependencyRepository.bulkInsertDocuments(changeLogDocuments, DependencyRepository.changeLogIndexName, BATCH_SIZE);

            jobStatusRepository.upsert(new JobStatusDocument(userId, repoName, "processing:inserted_all_changelogs"));
            log.info("inserted new changelog documents into openSearch successfully!", kv("repoName", repoName));
        }
    }
}
