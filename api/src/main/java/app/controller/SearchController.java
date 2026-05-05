package app.controller;

import java.io.IOException;
import java.util.List;
import java.util.Map;
import java.util.UUID;

import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import ai.djl.translate.TranslateException;
import app.dto.IssueSearchResponse;
import app.repository.SearchRepository;
import app.repository.UserRepoRepository;
import app.service.RerankingService;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;
import lombok.extern.slf4j.Slf4j;
import static net.logstash.logback.argument.StructuredArguments.kv;

@RestController
@RequestMapping("/search")
@Validated
@Slf4j
public class SearchController {
    private final SearchRepository searchRepository;
    private final UserRepoRepository userRepoRepository;
    private final RerankingService rerankingService;

    public SearchController(
        SearchRepository searchRepository,
        UserRepoRepository userRepoRepository,
        RerankingService rerankingService
    ) {
        this.searchRepository = searchRepository;
        this.userRepoRepository = userRepoRepository;
        this.rerankingService = rerankingService;
    }

    private record searchIssueRequest(
        @NotBlank String userId,
        @NotBlank String searchQuery
    ) {}
    
    /**
     * semantic search endpoint controller
     * @param request 
     * @return a list of the top 10 search results
     * @throws TranslateException from searcRepository.findRelevantIssues
     * @throws IOException from both userRepoRepository.findRepoDependency and searchRepository.findRelevantIssues
     */
    @PostMapping("/new")
    public ResponseEntity<?> searchRelevantIssues(
        @RequestBody @Valid searchIssueRequest request
    ) throws TranslateException, IOException {
        String requestId = UUID.randomUUID().toString();

        log.info("recieved hybrid search request", kv("query", request.searchQuery), kv("requestId", requestId));

        Map<String, String> userRepoDependencies = userRepoRepository.listAllDependencies(requestId, request.userId);
        if (userRepoDependencies.isEmpty()) {
            log.debug("no dependencies found for user", kv("userId", request.userId), kv("requestId", requestId));
            return ResponseEntity.status(404).body("No dependencies found for the userId");
        }

        List<IssueSearchResponse> issueResults = searchRepository.findRelevantIssues(
            userRepoDependencies, request.searchQuery, requestId
        );

        long start = System.currentTimeMillis();
        List<IssueSearchResponse> rerankedResults = rerankingService.rerank(request.searchQuery, issueResults);
        long elapsed = System.currentTimeMillis() - start;

        log.debug("reranking took {} ms", elapsed);

        return ResponseEntity.ok().body(rerankedResults);
    }
}
