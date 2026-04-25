package app.controller;

import java.io.IOException;

import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import app.service.GithubApiService;
import jakarta.validation.Valid;                                                                                                                                                       
import jakarta.validation.constraints.NotBlank;

@RestController
@RequestMapping("/tracked_repository")
@Validated
public class TrackedRepositoryController {
    private final GithubApiService githubApiService;

    public TrackedRepositoryController(GithubApiService githubApiService) {
        this.githubApiService = githubApiService;
    }

    private record AddRepoRequest(@NotBlank String repositoryUrl) {}

    @PostMapping("/add")
    public ResponseEntity<String> addNewRepo(@RequestBody @Valid AddRepoRequest request) throws IOException { 
        githubApiService.validateRepoExist(request.repositoryUrl);   
        githubApiService.indexRepoIssues(request.repositoryUrl);
        
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
