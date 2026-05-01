package app.service;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CompletableFuture;

import org.kohsuke.github.GHFileNotFoundException;
import org.kohsuke.github.GHIssue;
import org.kohsuke.github.GHIssueState;
import org.kohsuke.github.GHRepository;
import org.kohsuke.github.GitHub;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;
import org.springframework.validation.annotation.Validated;

import jakarta.validation.constraints.NotBlank;
import lombok.extern.slf4j.Slf4j;
import static net.logstash.logback.argument.StructuredArguments.kv;

@Service
@Validated
@Slf4j
public class GithubApiService {
    private final GitHub githubClient;

    // constructor
    public GithubApiService(GitHub githubClient) {
        this.githubClient = githubClient;
    }

    public static record IssueDocument(
        @NotBlank String repoName, 
        @NotBlank String url, 
        @NotBlank String title, 
        @NotBlank String body
    ) {}

    @Async
    public CompletableFuture<List<IssueDocument>> fetchRepoIssues(@NotBlank String repoName) {
        try {
            GHRepository repo = githubClient.getRepository(repoName);
            
            List<GHIssue> issues = repo.queryIssues()
                .state(GHIssueState.OPEN)
                .list()
                .toList();

            log.debug("fetched {} issues for repo {}", issues.size(), repoName);
            
            List<IssueDocument> issueDocuments = new ArrayList<>();
            for (GHIssue issue: issues) {
                String body = issue.getBody();
                if (body == null || body.isBlank()) continue; // skip over issues will blank body

                issueDocuments.add(new IssueDocument(
                    repoName,
                    issue.getHtmlUrl().toString(), 
                    issue.getTitle(), 
                    body
                ));
            }

            return CompletableFuture.completedFuture(issueDocuments);
        } catch (GHFileNotFoundException e) {
            log.error("Github repo {} not found", repoName, kv("error", e.getMessage()));
            return CompletableFuture.failedFuture(e);
        } catch (IOException e) {
            log.error("Unknown error fetching repo {} issues", repoName, kv("error", e.getMessage()));
            return CompletableFuture.failedFuture(e);
        }
    }
}
