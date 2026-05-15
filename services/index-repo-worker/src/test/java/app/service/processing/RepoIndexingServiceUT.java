package app.service.processing;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.atLeastOnce;
import static org.mockito.Mockito.never;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.stream.Stream;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import app.component.parser.DependencyParserStrategy.Dependency;
import app.dto.DependencyDocument;
import app.dto.GithubChangeLogResponse;
import app.repository.DependencyRepository;
import app.repository.JobStatusRepository;
import app.repository.UserRepoRepository;
import app.service.github.ChangelogService;
import app.service.github.IssueService;
import io.micrometer.tracing.Span;
import io.micrometer.tracing.Tracer;
import jakarta.validation.ConstraintViolation;
import jakarta.validation.Validation;
import jakarta.validation.executable.ExecutableValidator;

public class RepoIndexingServiceUT {
    
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
    private ExecutableValidator executableValidator;

    @BeforeEach
    void setup() throws Exception {
        MockitoAnnotations.openMocks(this);

        when(tracer.currentSpan()).thenReturn(parentSpan);
        when(tracer.withSpan(any())).thenReturn(spanInScope);

        service = new RepoIndexingService(
            dependencyRepository, jobStatusRepository, userRepoRepository, 
            issueService, changelogService, textEmbeddingService, tracer
        );

        executableValidator = Validation.buildDefaultValidatorFactory()
            .getValidator().forExecutables();
    }

    static Stream<Arguments> invalidProcessAllArgs() {
        Map<String, List<Dependency>> validDeps = Map.of("go", List.of(new Dependency("dep", "1.0", "org/dep")));

        return Stream.of(
            Arguments.of(null,      "repo", "user"),
            Arguments.of(Map.of(),      "repo", "user"),
            Arguments.of(validDeps, "", "user"),
            Arguments.of(validDeps, "repo", "")
        );
    }

    @ParameterizedTest
    @MethodSource("invalidProcessAllArgs")
    void processAll_InvalidInput(Map<String, List<Dependency>> deps, String repoName, String userId) throws NoSuchMethodException {
        Set<ConstraintViolation<RepoIndexingService>> violations = executableValidator.validateParameters(
            service,
            RepoIndexingService.class.getMethod("processAll", Map.class, String.class, String.class),
            new Object[]{deps, repoName, userId}
        );

        assertFalse(violations.isEmpty());
    }

    @Test
    void fetchDependencyIssues_PropagatesParentSpan() throws Exception {
        Dependency dep = new Dependency("dep", "1.0", "org/dep");
        Map<String, List<Dependency>> deps = Map.of("go", List.of(dep));
        when(issueService.fetch("org/dep")).thenReturn(List.of());

        service.fetchDependencyIssues(deps, Set.of());

        verify(tracer).currentSpan();
        verify(tracer, atLeastOnce()).withSpan(parentSpan);
    }

    @Test
    void fetchDependencyChangelogs_PropagatesParentSpan() throws Exception {
        Dependency dep = new Dependency("dep", "1.0", "org/dep");
        Map<String, List<Dependency>> deps = Map.of("go", List.of(dep));
        when(changelogService.fetchForVersion("org/dep", "1.0"))
            .thenReturn(new GithubChangeLogResponse("org/dep", "1.0", "fixes", "http://url"));

        service.fetchDependencyChangelogs(deps, Set.of());

        verify(tracer).currentSpan();
        verify(tracer, atLeastOnce()).withSpan(parentSpan);
    }

    @Test
    void fetchDependencyIssues_SkipsAlreadyIndexedRepos() throws Exception {
        Dependency dep = new Dependency("dep", "1.0", "org/dep");
        Map<String, List<Dependency>> deps = Map.of("go", List.of(dep));
        Set<DependencyDocument> indexed = Set.of(new DependencyDocument("org/dep", "1.0"));

        service.fetchDependencyIssues(deps, indexed);

        verify(issueService, never()).fetch(anyString());
    }

    @Test
    void fetchDependencyIssues_FetchesNonIndexedRepos() throws Exception {
        Dependency depA = new Dependency("dep-a", "1.0", "org/dep-a");
        Dependency depB = new Dependency("dep-b", "2.0", "org/dep-b");
        Map<String, List<Dependency>> deps = Map.of("go", List.of(depA, depB));
        Set<DependencyDocument> indexed = Set.of(new DependencyDocument("org/dep-a", "1.0"));
        when(issueService.fetch("org/dep-b")).thenReturn(List.of());

        service.fetchDependencyIssues(deps, indexed);

        verify(issueService, never()).fetch("org/dep-a");
        verify(issueService).fetch("org/dep-b");
    }

    @Test
    void fetchDependencyChangelogs_SkipsAlreadyIndexedVersions() throws Exception {
        Dependency dep = new Dependency("dep", "1.0", "org/dep");
        Map<String, List<Dependency>> deps = Map.of("go", List.of(dep));
        Set<DependencyDocument> indexed = Set.of(new DependencyDocument("org/dep", "1.0"));

        service.fetchDependencyChangelogs(deps, indexed);

        verify(changelogService, never()).fetchForVersion(anyString(), anyString());
    }

    @Test
    void fetchDependencyChangelogs_FetchesNewVersion() throws Exception {
        Dependency dep = new Dependency("dep", "2.0", "org/dep");
        Map<String, List<Dependency>> deps = Map.of("go", List.of(dep));
        when(changelogService.fetchForVersion("org/dep", "2.0"))
            .thenReturn(new GithubChangeLogResponse("org/dep", "2.0", "new stuff", "http://url"));

        service.fetchDependencyChangelogs(deps, Set.of(new DependencyDocument("org/dep", "1.0")));

        verify(changelogService).fetchForVersion("org/dep", "2.0");
    }   

    @Test
    void fetchDependencyChangelogs_FiltersNoReleaseResults() throws Exception {
        Dependency dep = new Dependency("dep", "1.0", "org/dep");
        Map<String, List<Dependency>> deps = Map.of("go", List.of(dep));
        when(changelogService.fetchForVersion("org/dep", "1.0"))
            .thenReturn(new GithubChangeLogResponse("org/dep", "1.0", "no-release", "no-url"));

        RepoIndexingService.ChangeLogResult result = service.fetchDependencyChangelogs(deps, Set.of());

        assertTrue(result.changelogs().isEmpty());
    }
}
