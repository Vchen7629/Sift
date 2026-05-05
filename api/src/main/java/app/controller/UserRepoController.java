package app.controller;

import java.io.IOException;
import java.util.List;

import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import app.repository.UserRepoRepository;
import io.micrometer.observation.annotation.Observed;
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

    private record DeleteRepoRequest(@NotBlank String userId, @NotBlank String repoName) {}

    @DeleteMapping("/delete")
    @Observed(name="userrepo.delete.controller")
    public ResponseEntity<Void> deleteRepo(@RequestBody @Valid DeleteRepoRequest request) throws IOException {
        log.info("recieved delete repo request", kv("userId", request.userId), kv("repoUrl", request.repoName));

        indexedRepoRepository.delete(request.userId, request.repoName);

        return ResponseEntity.noContent().build();
    }

    @GetMapping("/list/{userId}")
    @Observed(name="userrepo.list.controller")
    public ResponseEntity<List<String>> listTrackedRepos(@NotBlank @PathVariable String userId) throws IOException {
        log.info("recieved list all tracked repos request");

        List<String> trackedRepos = indexedRepoRepository.listAll(userId);

        return ResponseEntity.ok().body(trackedRepos);
    }
}
