package app.service.githubRepo;

import java.io.IOException;
import java.util.concurrent.CompletableFuture;

import org.kohsuke.github.GHFileNotFoundException;
import org.kohsuke.github.GHRelease;
import org.kohsuke.github.GHRepository;
import org.kohsuke.github.GitHub;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;
import org.springframework.validation.annotation.Validated;

import jakarta.validation.constraints.NotBlank;

@Service
@Validated
public class ChangelogService {
    private final GitHub githubClient;

    public ChangelogService(GitHub githubClient) {
        this.githubClient = githubClient;
    }

    public record Result(
        @NotBlank String dependencyName,
        @NotBlank String version, 
        @NotBlank String changes,
        @NotBlank String url
    ) {};

    @Async
    public CompletableFuture<Result> fetchChangeLogForVersion(
        @NotBlank String repoName, @NotBlank String version
    ) {
        try {
            GHRepository repo = githubClient.getRepository(repoName);
            GHRelease release = repo.getReleaseByTagName(version);

            if (release == null) {
                return CompletableFuture.completedFuture(
                    new Result(repoName, version, "no-release", "no-url")
                );
            }

            // todo: chunk the release notes into a list of individual updates instead of one giant string for all updates
            return CompletableFuture.completedFuture(
                new Result(repoName, version, release.getBody(), release.getHtmlUrl().toString())
            );
        } catch (GHFileNotFoundException e) {
            return CompletableFuture.failedFuture(e);
        } catch (IOException e) {
            return CompletableFuture.failedFuture(e);
        }
    }
}
