package app.controller;

import java.io.IOException;
import java.util.UUID;

import org.kohsuke.github.GitHub;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;

import ai.djl.translate.TranslateException;
import app.repository.OpenSearchRepository;
import app.service.ProducerService;
import io.nats.client.JetStreamApiException;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;

@Controller
@RequestMapping("/index_repo")
@Validated
public class IndexRepoController {
    private final GitHub githubClient;
    private final OpenSearchRepository openSearchRepository;
    private final ProducerService producerService;

    public IndexRepoController(
        GitHub githubClient,
        OpenSearchRepository openSearchRepository,
        ProducerService producerService
    ) {
        this.githubClient = githubClient;
        this.openSearchRepository = openSearchRepository;
        this.producerService = producerService;
    }
    
    private record IndexRepoRequest(@NotBlank String repoName) {}

    @PostMapping("/add")
    public ResponseEntity<String> addNewRepo(
        @RequestBody @Valid IndexRepoRequest request
    ) throws IOException, JetStreamApiException, TranslateException { 
        githubClient.getRepository(request.repoName); 

        if (openSearchRepository.isRepoIndexed(request.repoName)) {
            return ResponseEntity.status(HttpStatus.CONFLICT).body("Repository already indexed");
        }      

        String requestId = UUID.randomUUID().toString();

        producerService.PublishIndexRepoJobRequest(new ProducerService.RepoIndexMsg(request.repoName, requestId));
        
        return ResponseEntity.accepted().body("added " + request.repoName + "to processing");
    }

    @GetMapping("/get_status/{repoName}")
    public ResponseEntity<String> getStatus(@NotBlank @PathVariable String repoName) throws IOException {
        String jobStatus = openSearchRepository.findJobStatus(repoName);

        if (jobStatus.equals(null)) {
            return ResponseEntity.status(404).body("repo status not found, add it for processing first");
        }

        return ResponseEntity.ok().body(jobStatus);
    }
}
