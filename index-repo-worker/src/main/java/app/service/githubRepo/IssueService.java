package app.service.githubRepo;

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

import app.dto.ProcessedGithubIssue;
import jakarta.validation.constraints.NotBlank;
import lombok.extern.slf4j.Slf4j;
import static net.logstash.logback.argument.StructuredArguments.kv;

@Service
@Validated
@Slf4j
public class IssueService {
    private final GitHub githubClient;

    // constructor
    public IssueService(GitHub githubClient) {
        this.githubClient = githubClient;
    }

    @Async
    public CompletableFuture<List<ProcessedGithubIssue>> fetchDependencyIssues(
        @NotBlank String dependencyName, @NotBlank String requestId
    ) {
        try {
            GHRepository repo = githubClient.getRepository(dependencyName);
            
            long start = System.currentTimeMillis();

            List<GHIssue> issues = repo.queryIssues()
                .state(GHIssueState.ALL).list().toList()
                .stream().filter(issue -> !issue.isPullRequest()).toList();

            long elapsed = System.currentTimeMillis() - start;

            log.debug("fetched {} issues in {}ms ({}s)", 
                issues.size(), elapsed, elapsed / 1000.0,
                kv("dependencyName", dependencyName),
                kv("requestId", requestId));
            
            List<ProcessedGithubIssue> issueDocuments = new ArrayList<>();
            for (GHIssue issue: issues) {
                addIssueDocument(issue, issueDocuments, dependencyName);
            }

            return CompletableFuture.completedFuture(issueDocuments);
        } catch (GHFileNotFoundException e) {
            log.error("Repo not found", 
                kv("dependencyName", dependencyName), 
                kv("error", e.getMessage()),
                kv("requestId", requestId));

            return CompletableFuture.failedFuture(e);
        } catch (IOException e) {
            log.error("Unknown error fetching issues", 
                kv("dependencyName", dependencyName),
                kv("error", e.getMessage()),
                kv("requestId", requestId));

            return CompletableFuture.failedFuture(e);
        }
    }

    private void addIssueDocument(
        GHIssue issue, List<ProcessedGithubIssue> issueDocuments, @NotBlank String repoName
    ) throws IOException {
        String rawBody = issue.getBody();
        if (rawBody == null || rawBody.isBlank()) return; // skip over issues with blank body

        String body = cleanIssueBody(rawBody);
        if (body.isBlank()) return; // skip issues whose body was only code blocks

        List<String> labelList = issue.getLabels().stream()
            .map(GHLabel::getName)
            .collect(Collectors.toList());
        
        String version = "no version";

        issueDocuments.add(new ProcessedGithubIssue(
            repoName,
            version,
            issue.getTitle(),
            body,
            issue.getHtmlUrl().toString(),
            labelList,
            issue.getCreatedAt().toString()
        ));
    }

    protected String cleanIssueBody(String body) {
        return body.replaceAll("(?s)```.*?```", "").trim();
    }
}
