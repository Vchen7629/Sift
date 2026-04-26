package app.controller;

import java.io.IOException;

import org.kohsuke.github.GitHub;
import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import ai.djl.translate.TranslateException;
import app.service.GithubApiService;
import app.service.TextEmbeddingService;
import jakarta.validation.Valid;                                                                                                                                                       
import jakarta.validation.constraints.NotBlank;

@RestController
@RequestMapping("/tracked_repository")
@Validated
public class TrackedRepositoryController {
    private final GitHub githubClient;
    private final GithubApiService githubApiService;
    private final TextEmbeddingService textEmbService;

    public TrackedRepositoryController(GitHub githubClient, GithubApiService githubApiService, TextEmbeddingService textEmbService) {
        this.githubClient = githubClient;
        this.githubApiService = githubApiService;
        this.textEmbService = textEmbService;
    }

    private record AddRepoRequest(@NotBlank String repositoryUrl) {}

    @PostMapping("/add")
    public ResponseEntity<String> addNewRepo(@RequestBody @Valid AddRepoRequest request) throws IOException { 
        githubClient.getRepository(request.repositoryUrl);   
        githubApiService.fetchRepoIssues(request.repositoryUrl)
            .thenApply(issueUrlTexts -> {
                try {
                    return textEmbService.generateEmbeddings(issueUrlTexts);
                } catch (TranslateException e) {
                    throw new RuntimeException(e);
                }
            })
            .thenAccept(embeddings -> {});
        
        return ResponseEntity.accepted().body("added");
    }

    private record DeleteRepoRequest(@NotBlank String repositoryName) {}

    @DeleteMapping("/delete")
    public ResponseEntity<Void> deleteRepo(@RequestBody @Valid DeleteRepoRequest request) {

        return ResponseEntity.noContent().build();
    }

    @GetMapping("/list")
    public ResponseEntity<String> listTrackedRepos() throws IOException {

        return ResponseEntity.ok().body("dicking robin down");
    }
}
