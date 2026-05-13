package app.service.githubRepo;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.stream.Collectors;

import org.kohsuke.github.GHIssue;
import org.kohsuke.github.GHIssueState;
import org.kohsuke.github.GHLabel;
import org.kohsuke.github.GHRepository;
import org.kohsuke.github.GitHub;
import org.springframework.stereotype.Service;
import org.springframework.validation.annotation.Validated;

import app.dto.ProcessedGithubIssue;
import io.micrometer.observation.annotation.Observed;
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

    @Observed(name="issue.fetch.service")
    public List<ProcessedGithubIssue> fetch(@NotBlank String dependencyName) throws IOException {
        GHRepository repo = githubClient.getRepository(dependencyName);
        
        List<GHIssue> issues = repo.queryIssues()
            .state(GHIssueState.ALL).list().toList()
            .stream().filter(issue -> !issue.isPullRequest()).toList();

        log.debug("fetched {} issues", issues.size(), kv("dependencyName", dependencyName));
        
        List<ProcessedGithubIssue> issueDocuments = new ArrayList<>();
        for (GHIssue issue: issues) {
            addIssueDocument(issue, issueDocuments, dependencyName);
        }

        return issueDocuments;
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
