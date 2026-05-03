package app.controller;

import java.io.IOException;
import java.util.List;
import java.util.UUID;

import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import app.repository.UserRepoRepository;
import jakarta.validation.Valid;                                                                                                                                                       
import jakarta.validation.constraints.NotBlank;
import lombok.extern.slf4j.Slf4j;
import static net.logstash.logback.argument.StructuredArguments.kv;


@RestController
@RequestMapping("/user_repo")
@Validated
@Slf4j
public class UserRepoController {
    private final UserRepoRepository indexedRepoRepository;

    public UserRepoController(UserRepoRepository indexedRepoRepository) {
        this.indexedRepoRepository = indexedRepoRepository;
    }

    private record DeleteRepoRequest(@NotBlank String repoName) {}

    @DeleteMapping("/delete")
    public ResponseEntity<Void> deleteRepo(@RequestBody @Valid DeleteRepoRequest request) throws IOException {
        String requestId = UUID.randomUUID().toString();

        log.info("recieved delete repo request", 
                kv("repoUrl", request.repoName), 
                kv("requestId", requestId));

        indexedRepoRepository.delete(request.repoName, requestId);

        return ResponseEntity.noContent().build();
    }

    @GetMapping("/list/{userId}")
    public ResponseEntity<List<String>> listTrackedRepos(@NotBlank @PathVariable String userId) throws IOException {
        String requestId = UUID.randomUUID().toString();
        log.info("recieved list all tracked repos request", kv("requestId", requestId));

        List<String> trackedRepos = indexedRepoRepository.findAll(requestId, userId);

        return ResponseEntity.ok().body(trackedRepos);
    }
}
