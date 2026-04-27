package app.controller;

import java.io.IOException;
import java.util.List;

import org.kohsuke.github.GitHub;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import ai.djl.translate.TranslateException;
import app.repository.OpenSearchRepository;
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
    private final OpenSearchRepository openSearchRepository;

    public TrackedRepositoryController(
        GitHub githubClient, 
        GithubApiService githubApiService, 
        TextEmbeddingService textEmbService,
        OpenSearchRepository openSearchRepository
    ) {
        this.githubClient = githubClient;
        this.githubApiService = githubApiService;
        this.textEmbService = textEmbService;
        this.openSearchRepository = openSearchRepository;
    }

    private record AddRepoRequest(@NotBlank String repositoryUrl) {}

    @PostMapping("/add")
    public ResponseEntity<String> addNewRepo(@RequestBody @Valid AddRepoRequest request) throws IOException, TranslateException { 
        githubClient.getRepository(request.repositoryUrl); 

        if (openSearchRepository.isRepoIndexed(request.repositoryUrl)) {
            return ResponseEntity.status(HttpStatus.CONFLICT).body("Repository already indexed");
        }       

        List<GithubApiService.IssueDocument> githubIssues = githubApiService.fetchRepoIssues(request.repositoryUrl).join();
        if (githubIssues.isEmpty()) {
            return ResponseEntity.ok().body("No open issues found for repository");
        }
        
        List<TextEmbeddingService.embeddingDocument> embeddings = textEmbService.generateEmbeddings(githubIssues);
        openSearchRepository.indexGithubIssue(embeddings);

        return ResponseEntity.accepted().body("added");
    }

    private record DeleteRepoRequest(@NotBlank String repositoryName) {}

    @DeleteMapping("/delete")
    public ResponseEntity<Void> deleteRepo(@RequestBody @Valid DeleteRepoRequest request) throws IOException {
        openSearchRepository.deleteTrackedRepo(request.repositoryName);

        return ResponseEntity.noContent().build();
    }

    @GetMapping("/list")
    public ResponseEntity<List<String>> listTrackedRepos() throws IOException {
        List<String> trackedRepos = openSearchRepository.findAllIndexedRepoNames();

        return ResponseEntity.ok().body(trackedRepos);
    }
}
