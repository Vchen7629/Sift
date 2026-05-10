package app.controller;

import java.io.IOException;

import org.kohsuke.github.GitHub;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;

import ai.djl.translate.TranslateException;
import app.repository.UserRepoRepository;
import app.dto.IndexRepoMsg;
import app.repository.JobStatusRepository;
import app.service.ProducerService;
import io.micrometer.observation.annotation.Observed;
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
    
    private record Request(@NotBlank String userId, @NotBlank String repoName) {}

    @PostMapping("/add")
    @Observed(name="indexrepo.addNewRepo.controller")
    public ResponseEntity<String> addNewRepo(
        @RequestBody @Valid Request request
    ) throws IOException, JetStreamApiException, TranslateException { 
        log.info("recieved add new repo request", kv("repoName", request.repoName()));

        githubClient.getRepository(request.repoName()); 

        if (indexedRepoRepository.isRepoIndexed(request.repoName())) {
            log.warn("repo already indexed, cant add it again", kv("repoName", request.repoName()));
            return ResponseEntity.status(HttpStatus.CONFLICT).body("Repository already indexed");
        }      

        producerService.PublishIndexRepoJobRequest(new IndexRepoMsg(request.userId(), request.repoName()));
        
        return ResponseEntity.accepted().body("added " + request.repoName() + " to processing");
    }

    @GetMapping("/get_status")
    @Observed(name="indexrepo.getStatus.controller")
    public ResponseEntity<String> getStatus(
        @RequestBody @Valid Request request
    ) throws IOException {
        log.info("recieved get repo index job status request", kv("userId", request.userId), kv("repoName", request.repoName));

        String jobStatus = jobStatusRepository.findStatus(request.userId, request.repoName);

        if (jobStatus == null) {
            log.warn("repo hasn't been added to db yet", kv("userId", request.userId), kv("repoName", request.repoName));
                
            return ResponseEntity.status(404).body("repo status not found, add it for processing first");
        }

        return ResponseEntity.ok().body(jobStatus);
    }
}
