package app.service.github;

import java.io.IOException;
import org.kohsuke.github.GHRelease;
import org.kohsuke.github.GHRepository;
import org.kohsuke.github.GitHub;
import org.springframework.stereotype.Service;
import org.springframework.validation.annotation.Validated;

import app.dto.GithubChangeLogResponse;
import io.micrometer.observation.annotation.Observed;
import jakarta.validation.constraints.NotBlank;

@Service
@Validated
public class ChangelogService {
    private final GitHub githubClient;

    public ChangelogService(GitHub githubClient) {
        this.githubClient = githubClient;
    }

    @Observed(name="changelog.fetchforversion.service")
    public GithubChangeLogResponse fetchForVersion(@NotBlank String repoName, @NotBlank String version) throws IOException {
        GHRepository repo = githubClient.getRepository(repoName);
        GHRelease release = repo.getReleaseByTagName(version);

        if (release == null) {
            return new GithubChangeLogResponse(repoName, version, "no-release", "no-url");
        }

        // todo: chunk the release notes into a list of individual updates instead of one giant string for all updates
        return new GithubChangeLogResponse(repoName, version, release.getBody(), release.getHtmlUrl().toString());
    }
}
