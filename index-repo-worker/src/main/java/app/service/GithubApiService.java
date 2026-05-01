package app.service;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.stream.Collectors;

import org.kohsuke.github.GHFileNotFoundException;
import org.kohsuke.github.GHIssue;
import org.kohsuke.github.GHIssueState;
import org.kohsuke.github.GHLabel;
import org.kohsuke.github.GHRepository;
import org.kohsuke.github.GitHub;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;
import org.springframework.validation.annotation.Validated;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
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
        @NotBlank String body,
        @NotNull List<String> labelList
    ) {}

    @Async
    public CompletableFuture<List<IssueDocument>> fetchRepoIssues(
        @NotBlank String repoName, @NotBlank String requestId
    ) {
        try {
            GHRepository repo = githubClient.getRepository(repoName);
            
            long start = System.currentTimeMillis();

            List<GHIssue> issues = repo.queryIssues()
                .state(GHIssueState.OPEN)
                .list()
                .toList();

            long elapsed = System.currentTimeMillis() - start;

            log.debug("fetched {} issues in {}ms ({}s)", 
                issues.size(), elapsed, elapsed / 1000.0,
                kv("repoName", repoName),
                kv("requestId", requestId));
            
            List<IssueDocument> issueDocuments = new ArrayList<>();
            for (GHIssue issue: issues) {
                addIssueDocument(issue, issueDocuments, repoName);
            }

            return CompletableFuture.completedFuture(issueDocuments);
        } catch (GHFileNotFoundException e) {
            log.error("Repo not found", 
                kv("repoName", repoName), 
                kv("error", e.getMessage()),
                kv("requestId", requestId));

            return CompletableFuture.failedFuture(e);
        } catch (IOException e) {
            log.error("Unknown error fetching issues", 
                kv("repoName", repoName), 
                kv("error", e.getMessage()),
                kv("requestId", requestId));

            return CompletableFuture.failedFuture(e);
        }
    }

    private void addIssueDocument(GHIssue issue, List<IssueDocument> issueDocuments, @NotBlank String repoName) {
        String body = issue.getBody();
        if (body == null || body.isBlank()) return; // skip over issues with blank body

        int charCountLimit = 32000;
        if (body.length() > charCountLimit) {
            // todo: remove truncation bandaid fix with proper issue body chunking in the future
            body = body.substring(0, charCountLimit); 
        }

        List<String> labelList = issue.getLabels().stream()
            .map(GHLabel::getName)
            .collect(Collectors.toList());

        issueDocuments.add(new IssueDocument(
            repoName,
            issue.getHtmlUrl().toString(), 
            issue.getTitle(), 
            body,
            labelList
        ));
    }
}
