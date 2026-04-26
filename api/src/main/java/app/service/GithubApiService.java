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

@Service
@Validated
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
                .withPageSize(20) // todo: remove this 20 issue limit later, testing for now                                                                                                                                                                
                .iterator()
                .nextPage(); 
            
            List<IssueDocument> issueDocuments = new ArrayList<>();
            for (GHIssue issue: issues) {
                issueDocuments.add(new IssueDocument(
                    repoName,
                    issue.getHtmlUrl().toString(), 
                    issue.getTitle(), 
                    issue.getBody()
                ));
            }

            return CompletableFuture.completedFuture(issueDocuments);
        } catch (GHFileNotFoundException e) {
            return CompletableFuture.failedFuture(e);
        } catch (IOException e) {
            return CompletableFuture.failedFuture(e);
        }
    }
}
