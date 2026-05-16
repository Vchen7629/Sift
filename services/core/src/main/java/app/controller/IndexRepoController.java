package app.controller;

import java.io.IOException;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.stereotype.Controller;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;

import ai.djl.translate.TranslateException;
import app.repository.UserRepoRepository;
import app.service.messaging.ProducerService;
import app.dto.IndexRepoMsg;
import app.repository.JobStatusRepository;
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
    private final UserRepoRepository indexedRepoRepository;
    private final JobStatusRepository jobStatusRepository;
    private final ProducerService producerService;

    public IndexRepoController(
        UserRepoRepository indexedRepoRepository,
        JobStatusRepository jobStatusRepository,
        ProducerService producerService
    ) {
        this.indexedRepoRepository = indexedRepoRepository;
        this.jobStatusRepository = jobStatusRepository;
        this.producerService = producerService;
    }
    
    private record ReqBody(@NotBlank String repoName) {}

    @PostMapping("/add")
    @Observed(name="indexrepo.addNewRepo.controller")
    public ResponseEntity<String> addNewRepo(
        @RequestBody @Valid ReqBody reqBody, @AuthenticationPrincipal String username
    ) throws IOException, JetStreamApiException, TranslateException { 
        log.info("recieved add new repo request", kv("repoName", reqBody.repoName()));

        if (indexedRepoRepository.isRepoIndexed(reqBody.repoName())) {
            log.warn("repo already indexed, cant add it again", kv("repoName", reqBody.repoName()));
            return ResponseEntity.status(HttpStatus.CONFLICT).body("Repository already indexed");
        }      

        producerService.PublishIndexRepoJobRequest(new IndexRepoMsg(username, reqBody.repoName()));
        
        return ResponseEntity.accepted().body("added " + reqBody.repoName() + " to processing");
    }

    @PostMapping("/job_status")
    public ResponseEntity<String> getStatus(
        @RequestBody @Valid ReqBody reqBody, @AuthenticationPrincipal String username
    ) throws IOException {
        log.info("recieved get repo index job status request", kv("repoName", reqBody.repoName()));

        String jobStatus = jobStatusRepository.findStatus(username, reqBody.repoName());

        if (jobStatus == null) {
            log.warn("repo hasn't been added to db yet", kv("username", username), kv("repoName", reqBody.repoName()));
                
            return ResponseEntity.status(404).body("repo status not found, add it for processing first");
        }

        return ResponseEntity.ok().body(jobStatus);
    }
}
