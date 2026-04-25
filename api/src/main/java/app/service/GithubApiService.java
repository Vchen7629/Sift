package app.service;

import java.io.IOException;
import java.util.List;
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
        githubClient.getRepository(repoName);
    }

    @Async
    public CompletableFuture<Void> indexRepoIssues(String repoName) {
        try {
            GHRepository repo = githubClient.getRepository(repoName);
            
            List<GHIssue> issues = repo.queryIssues()
                .state(GHIssueState.OPEN)
                .list()
                .withPageSize(20) // todo: remove this 20 issue limit later, testing for now                                                                                                                                                                
                .iterator()
                .nextPage(); 

            for (GHIssue issue: issues) {
                // todo: add issues to opensearch or smth

                String message = String.format("fetched issue #" + issue.getNumber() + ": " + issue.getTitle() + "\n");
                System.out.println(message);
            }

            return CompletableFuture.completedFuture(null);
        } catch (GHFileNotFoundException e) {
            return CompletableFuture.failedFuture(e);
        } catch (IOException e) {
            return CompletableFuture.failedFuture(e);
        }
    }
}
