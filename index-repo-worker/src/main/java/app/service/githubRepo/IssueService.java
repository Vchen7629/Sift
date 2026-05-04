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

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
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

    public static record Result(
        @NotBlank String dependencyName, 
        @NotBlank String version,
        @NotBlank String title, 
        @NotBlank String body,
        @NotBlank String url, 
        @NotNull List<String> labelList,
        @NotBlank String createdOn
    ) {}

    @Async
    public CompletableFuture<List<Result>> fetchDependencyIssues(
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
            
            List<Result> issueDocuments = new ArrayList<>();
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

    private void addIssueDocument(GHIssue issue, List<Result> issueDocuments, @NotBlank String repoName) throws IOException {
        //if (!hasRelevantLabel(issue)) return;

        String body = issue.getBody();
        if (body == null || body.isBlank()) return; // skip over issues with blank body

        List<String> labelList = issue.getLabels().stream()
            .map(GHLabel::getName)
            .collect(Collectors.toList());
        
        String version = "no version";

        issueDocuments.add(new Result(
            repoName,
            version,
            issue.getTitle(),
            cleanIssueBody(body),
            issue.getHtmlUrl().toString(),
            labelList,
            issue.getCreatedAt().toString()
        ));
    }

    String cleanIssueBody(String body) {
        return body.replaceAll("(?s)```.*?```", "").trim();
    }
}
