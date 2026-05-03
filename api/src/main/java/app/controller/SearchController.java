package app.controller;

import java.io.IOException;
import java.util.List;
import java.util.UUID;

import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import ai.djl.translate.TranslateException;
import app.repository.IndexedRepoRepository;
import app.repository.SearchRepository;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;
import lombok.extern.slf4j.Slf4j;
import static net.logstash.logback.argument.StructuredArguments.kv;

@RestController
@RequestMapping("/search")
@Validated
@Slf4j
public class SearchController {
    private final IndexedRepoRepository indexedRepoRepository;
    private final SearchRepository searchRepository;

    public SearchController(
        IndexedRepoRepository indexedRepoRepository,
        SearchRepository searchRepository
    ) {
        this.indexedRepoRepository = indexedRepoRepository;
        this.searchRepository = searchRepository;
    }

    private record searchIssueRequest(
        @NotBlank String repoName, 
        @NotBlank String searchQuery
    ) {}
    
    @PostMapping("/new")
    public ResponseEntity<?> searchRelevantIssues(@RequestBody @Valid searchIssueRequest request) throws TranslateException, IOException {
        String requestId = UUID.randomUUID().toString();

        log.info("recieved hybrid search request", 
                kv("repoName", request.repoName),
                kv("query", request.searchQuery),
                kv("requestId", requestId));

        if (!indexedRepoRepository.isRepoIndexed(request.repoName)) {
            log.warn("no results for unindexed repo", kv("repoName", request.repoName), kv("requestId", requestId));

            return ResponseEntity.status(404).body("Repository isn't indexed, no relevant issues");
        }

        List<SearchRepository.IssueSearchResult> issueResult =  searchRepository.findRelevantIssues(
            request.repoName, request.searchQuery, requestId);

        return ResponseEntity.ok().body(issueResult);
    }
}
