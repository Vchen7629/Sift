package app.controller;

import java.io.IOException;
import java.util.List;

import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import app.dto.IndexedRepoDocument;
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

    private record DeleteRepoRequest(@NotBlank String repoName) {}

    @DeleteMapping("/delete")
    @Observed(name="userrepo.delete.controller")
    public ResponseEntity<?> deleteRepo(
        @RequestBody @Valid DeleteRepoRequest reqBody, @AuthenticationPrincipal String username
    ) throws IOException {
        log.info("recieved delete repo request", kv("userId", username), kv("repoUrl", reqBody.repoName()));
        indexedRepoRepository.delete(username, reqBody.repoName());

        return ResponseEntity.noContent().build();
    }

    @GetMapping("/list")
    @Observed(name="userrepo.list.controller")
    public ResponseEntity<?> listTrackedRepos(@AuthenticationPrincipal String username) throws IOException {
        log.info("recieved list all tracked repos request");
        List<IndexedRepoDocument> trackedRepos = indexedRepoRepository.listAll(username);

        return ResponseEntity.ok().body(trackedRepos);
    }
}
