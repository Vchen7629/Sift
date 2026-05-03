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
import app.repository.UserRepoRepository;
import app.repository.JobStatusRepository;
import app.service.ProducerService;
import io.nats.client.JetStreamApiException;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;
import lombok.extern.slf4j.Slf4j;
import static net.logstash.logback.argument.StructuredArguments.kv;


@Controller
@RequestMapping("/index_repo")
@Validated
@Slf4j
public class IndexRepoController {
    private final GitHub githubClient;
    private final UserRepoRepository indexedRepoRepository;
    private final JobStatusRepository jobStatusRepository;
    private final ProducerService producerService;

    public IndexRepoController(
        GitHub githubClient,
        UserRepoRepository indexedRepoRepository,
        JobStatusRepository jobStatusRepository,
        ProducerService producerService
    ) {
        this.githubClient = githubClient;
        this.indexedRepoRepository = indexedRepoRepository;
        this.jobStatusRepository = jobStatusRepository;
        this.producerService = producerService;
    }
    
    private record IndexRepoRequest(@NotBlank String repoName) {}

    @PostMapping("/add")
    public ResponseEntity<String> addNewRepo(
        @RequestBody @Valid IndexRepoRequest request
    ) throws IOException, JetStreamApiException, TranslateException { 
        String requestId = UUID.randomUUID().toString();

        log.info("recieved add new repo request", 
            kv("repoName", request.repoName), 
            kv("requestId", requestId));

        githubClient.getRepository(request.repoName); 

        if (indexedRepoRepository.isRepoIndexed(request.repoName)) {
            log.warn("repo already indexed, cant add it again", kv("repoName", request.repoName));
            return ResponseEntity.status(HttpStatus.CONFLICT).body("Repository already indexed");
        }      

        producerService.PublishIndexRepoJobRequest(new ProducerService.RepoIndexMsg(request.repoName, requestId));
        
        return ResponseEntity.accepted().body("added " + request.repoName + " to processing");
    }

    @GetMapping("/get_status/{repoName}")
    public ResponseEntity<String> getStatus(@NotBlank @PathVariable String repoName) throws IOException {
        String requestId = UUID.randomUUID().toString();

        log.info("recieved get repo index job status request", 
            kv("repoName", repoName),
            kv("requestId", requestId));

        String jobStatus = jobStatusRepository.findJobStatus(repoName);

        if (jobStatus.equals(null)) {
            log.warn("repo hasn't been added to db yet", 
                kv("repoName", repoName),
                kv("requestId", requestId));
                
            return ResponseEntity.status(404).body("repo status not found, add it for processing first");
        }

        return ResponseEntity.ok().body(jobStatus);
    }
}
