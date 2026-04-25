package app.service;

import java.io.IOException;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;

import org.kohsuke.github.GHFileNotFoundException;
import org.kohsuke.github.GHIssue;
import org.kohsuke.github.GHIssueState;
import org.kohsuke.github.GHRepository;
import org.kohsuke.github.GitHub;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;

@Service
public class GithubApiService {
    private final GitHub githubClient;

    // constructor
    public GithubApiService(GitHub githubClient) {
        this.githubClient = githubClient;
    }

    public void validateRepoExist(String repoName) throws IOException {
        if (repoName == null || repoName.trim().isEmpty()) {
            throw new IllegalArgumentException("repoName must not be null or empty");
        }

        githubClient.getRepository(repoName);
    }

    @Async
    public CompletableFuture<Map<String, String>> fetchRepoIssues(String repoName) {
        try {
            GHRepository repo = githubClient.getRepository(repoName);
            
            List<GHIssue> issues = repo.queryIssues()
                .state(GHIssueState.OPEN)
                .list()
                .withPageSize(20) // todo: remove this 20 issue limit later, testing for now                                                                                                                                                                
                .iterator()
                .nextPage(); 

            Map<String, String> issueUrlTexts = new HashMap<>();
            for (GHIssue issue: issues) {
                issueUrlTexts.put(
                    issue.getHtmlUrl().toString(), 
                    textEmbeddingService.combineDescBody(issue.getTitle(), issue.getBody())
                );
            }

            return CompletableFuture.completedFuture(issueUrlTexts);
        } catch (GHFileNotFoundException e) {
            return CompletableFuture.failedFuture(e);
        } catch (IOException e) {
            return CompletableFuture.failedFuture(e);
        }
    }
}
