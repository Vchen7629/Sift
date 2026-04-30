package app.controller;

import java.io.IOException;
import java.util.List;

import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import app.repository.OpenSearchRepository;
import jakarta.validation.Valid;                                                                                                                                                       
import jakarta.validation.constraints.NotBlank;

@RestController
@RequestMapping("/indexed_repo")
@Validated
public class IndexedRepoController {
    private final OpenSearchRepository openSearchRepository;

    public IndexedRepoController(OpenSearchRepository openSearchRepository) {
        this.openSearchRepository = openSearchRepository;
    }

    private record DeleteRepoRequest(@NotBlank String repoName) {}

    @DeleteMapping("/delete")
    public ResponseEntity<Void> deleteRepo(@RequestBody @Valid DeleteRepoRequest request) throws IOException {
        openSearchRepository.deleteTrackedRepo(request.repoName);

        return ResponseEntity.noContent().build();
    }

    @GetMapping("/list")
    public ResponseEntity<List<String>> listTrackedRepos() throws IOException {
        List<String> trackedRepos = openSearchRepository.findAllIndexedRepoNames();

        return ResponseEntity.ok().body(trackedRepos);
    }
}
