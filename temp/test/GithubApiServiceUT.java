package app.service;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.argThat;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.never;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import java.io.ByteArrayInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ExecutionException;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.kohsuke.github.GHContent;
import org.kohsuke.github.GHRepository;
import org.kohsuke.github.GitHub;

import app.model.DependencyFileEnum;
import app.parser.DependencyParserStrategy;
import app.parser.DependencyParserStrategy.Dependency;

public class GithubApiServiceUT {

    private GitHub githubClient;
    private GHRepository repo;
    private DependencyParserStrategy parserStrategy;
    private GithubApiService service;

    @BeforeEach
    void setup() throws IOException {
        githubClient = mock(GitHub.class);
        repo = mock(GHRepository.class);
        parserStrategy = mock(DependencyParserStrategy.class);

        when(githubClient.getRepository("owner/repo")).thenReturn(repo);

        service = new GithubApiService(githubClient, parserStrategy);
    }

    // fetchRepoDependencies tests

    @Test
    void callsMultipleParsers_whenMultipleFilesFound() throws IOException, ExecutionException, InterruptedException {
        String goModContent = "module test\n\ngo 1.21\n\nrequire (\n\tgithub.com/foo/bar v1.0.0\n)\n";
        GHContent goModFile = mock(GHContent.class);

        
        
        when(repo.getFileContent("go.mod")).thenReturn(goModFile);
        when(repo.getFileContent(argThat(f -> !f.equals("go.mod")))).thenThrow(new FileNotFoundException());
        when(goModFile.read()).thenReturn(new ByteArrayInputStream(goModContent.getBytes(StandardCharsets.UTF_8)));

        List<Dependency> expected = List.of(new Dependency("github.com/foo/bar", "v1.0.0"));
        when(parserStrategy.parse(eq(DependencyFileEnum.GO_MOD), any())).thenReturn(expected);

        Map<String, List<Dependency>> result = service.fetchRepoDependencies("owner/repo", "req-1").get();

        verify(parserStrategy).parse(eq(DependencyFileEnum.GO_MOD), eq(goModContent));
        assertThat(result.get("go")).isEqualTo(expected);
    }


    @Test
    void skipsLanguage_whenNoFileFound() throws IOException, ExecutionException, InterruptedException {
        when(repo.getFileContent(any())).thenThrow(new FileNotFoundException());

        Map<String, List<Dependency>> result = service.fetchRepoDependencies("owner/repo", "req-1").get();

        verify(parserStrategy, never()).parse(any(), any());
        assertThat(result).isEmpty();
    }
}
