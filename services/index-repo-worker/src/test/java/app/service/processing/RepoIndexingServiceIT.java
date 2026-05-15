package app.service.processing;

import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyList;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.ArgumentMatchers.argThat;
import static org.mockito.Mockito.doThrow;
import static org.mockito.Mockito.inOrder;
import static org.mockito.Mockito.never;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import java.io.IOException;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.ExecutionException;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.InOrder;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import app.component.parser.DependencyParserStrategy.Dependency;
import app.dto.DependencyDocument;
import app.dto.GithubChangeLogResponse;
import app.dto.IndexableDocuments;
import app.dto.ProcessedGithubIssue;
import app.dto.UserRepoDocument;
import app.repository.DependencyRepository;
import app.repository.JobStatusRepository;
import app.repository.UserRepoRepository;
import app.service.github.ChangelogService;
import app.service.github.IssueService;
import io.micrometer.tracing.Span;
import io.micrometer.tracing.Tracer;

public class RepoIndexingServiceIT {
    
    @Mock private DependencyRepository dependencyRepository;
    @Mock private JobStatusRepository jobStatusRepository;
    @Mock private UserRepoRepository userRepoRepository;
    @Mock private IssueService issueService;
    @Mock private ChangelogService changelogService;
    @Mock private TextEmbeddingService textEmbeddingService;
    @Mock private Tracer tracer;
    @Mock private Span parentSpan;
    @Mock private Tracer.SpanInScope spanInScope;

    private RepoIndexingService service;

    @BeforeEach
    void setup() {
        MockitoAnnotations.openMocks(this);
        
        when(tracer.currentSpan()).thenReturn(parentSpan);
        when(tracer.withSpan(any())).thenReturn(spanInScope);

        service = new RepoIndexingService(
            dependencyRepository, jobStatusRepository, userRepoRepository, 
            issueService, changelogService, textEmbeddingService, tracer
        );
    }

    @Test
    void processAll_HappyPath() throws Exception {
        Dependency dep = new Dependency("dep", "1.0", "org/dep");
        ProcessedGithubIssue issue = new ProcessedGithubIssue("org/dep", "no version", "title", "body", "http://issue", List.of(), "2024-01-01");
        GithubChangeLogResponse changelog = new GithubChangeLogResponse("org/dep", "1.0", "fixes", "http://changelog");
        IndexableDocuments.Issue embeddedIssue = new IndexableDocuments.Issue("org/dep", "no version", "title", "body", "http://issue", List.of(),
"2024-01-01", new float[]{0.1f}, new float[]{0.2f});
        IndexableDocuments.ChangeLog embeddedChangelog = new IndexableDocuments.ChangeLog("org/dep", "1.0", "fixes", "http://changelog", new float[]{0.3f});

        when(dependencyRepository.list()).thenReturn(Set.of());
        when(issueService.fetch("org/dep")).thenReturn(List.of(issue));
        when(changelogService.fetchForVersion("org/dep", "1.0")).thenReturn(changelog);
        when(textEmbeddingService.githubIssue(anyList())).thenReturn(List.of(embeddedIssue));
        when(textEmbeddingService.githubChangelog(anyList())).thenReturn(List.of(embeddedChangelog));

        service.processAll(Map.of("go", List.of(dep)), "my-repo", "user-1");

        verify(dependencyRepository).bulkInsertDocuments(List.of(embeddedIssue), DependencyRepository.issuesIndexName);
        verify(dependencyRepository).bulkInsertDocuments(List.of(embeddedChangelog), DependencyRepository.changeLogIndexName);
        verify(userRepoRepository).insert(any(UserRepoDocument.class));
    }

    @Test
    void processAll_WritesJobStatusesInOrder() throws Exception {
        Dependency dep = new Dependency("dep", "1.0", "org/dep");
        when(dependencyRepository.list()).thenReturn(Set.of());
        when(issueService.fetch(any())).thenReturn(List.of());
        when(changelogService.fetchForVersion(any(), any()))
            .thenReturn(new GithubChangeLogResponse("org/dep", "1.0", "fixes", "http://url"));
        when(textEmbeddingService.githubChangelog(anyList()))
            .thenReturn(List.of(new IndexableDocuments.ChangeLog("org/dep", "1.0", "fixes", "http://url", new float[]{0.1f})));

        service.processAll(Map.of("go", List.of(dep)), "my-repo", "user-1");

        InOrder inOrder = inOrder(jobStatusRepository);
        inOrder.verify(jobStatusRepository).upsert(argThat(d -> d.status().equals("processing:fetched_all_issues_changelogs")));
        inOrder.verify(jobStatusRepository).upsert(argThat(d -> d.status().equals("processing:created_embeddings")));
        inOrder.verify(jobStatusRepository).upsert(argThat(d -> d.status().equals("processing:inserted_all_changelogs")));
        inOrder.verify(jobStatusRepository).upsert(argThat(d -> d.status().equals("processing:inserted_indexed_repo")));
    }

    @Test
    void processAll_SkipsEmbeddingWhenNoNewContent() throws Exception {
        Dependency dep = new Dependency("dep", "1.0", "org/dep");
        when(dependencyRepository.list()).thenReturn(Set.of(new DependencyDocument("org/dep", "1.0")));

        service.processAll(Map.of("go", List.of(dep)), "my-repo", "user-1");

        verify(textEmbeddingService, never()).githubIssue(anyList());
        verify(textEmbeddingService, never()).githubChangelog(anyList());
    }

    @Test
    void processAll_WrapsIOExceptionFromBulkInsert() throws Exception {
        Dependency dep = new Dependency("dep", "1.0", "org/dep");
        IndexableDocuments.ChangeLog embeddedChangelog =
            new IndexableDocuments.ChangeLog("org/dep", "1.0", "fixes", "http://url", new float[]{0.1f});

        when(dependencyRepository.list()).thenReturn(Set.of());
        when(issueService.fetch(any())).thenReturn(List.of());
        when(changelogService.fetchForVersion(any(), any()))
            .thenReturn(new GithubChangeLogResponse("org/dep", "1.0", "fixes", "http://url"));
        when(textEmbeddingService.githubChangelog(anyList())).thenReturn(List.of(embeddedChangelog));
        doThrow(new IOException("index error"))
            .when(dependencyRepository).bulkInsertDocuments(anyList(), anyString());

        assertThrows(ExecutionException.class, () ->
            service.processAll(Map.of("go", List.of(dep)), "my-repo", "user-1")
        );
    }
}
