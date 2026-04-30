package app.controller;

import java.io.IOException;
import java.util.List;

import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import ai.djl.translate.TranslateException;
import app.repository.OpenSearchRepository;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;

@RestController
@RequestMapping("/issue")
@Validated
public class GitIssuesController {
    private final OpenSearchRepository openSearchRepository;

    public GitIssuesController(OpenSearchRepository openSearchRepository) {
        this.openSearchRepository = openSearchRepository;
    }

    private record searchIssueRequest(@NotBlank String repoUrl, @NotBlank String searchQuery) {}
    
    @PostMapping("/search")
    public ResponseEntity<?> searchRelevantIssues(@RequestBody @Valid searchIssueRequest query) throws TranslateException, IOException {
        if (!openSearchRepository.isRepoIndexed(query.repoUrl)) {
            return ResponseEntity.status(404).body("Repository isn't indexed, no relevant issues");
        }

        List<OpenSearchRepository.IssueSearchResult> issueResult =  openSearchRepository.findRelevantIssues(
            query.repoUrl,
            query.searchQuery
        );

        return ResponseEntity.ok().body(issueResult);
    }
}
