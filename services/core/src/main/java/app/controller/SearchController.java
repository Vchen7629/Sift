package app.controller;

import java.io.IOException;
import java.util.List;

import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import ai.djl.translate.TranslateException;
import app.dto.IssueSearchResponse;
import app.service.external.OllamaService;
import app.service.internal.SearchResponseService;
import app.service.ml.RerankingService;
import io.micrometer.observation.annotation.Observed;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import lombok.extern.slf4j.Slf4j;
import static net.logstash.logback.argument.StructuredArguments.kv;

@RestController
@RequestMapping("/search")
@Validated
@Slf4j
public class SearchController {
    private final SearchResponseService searchResponseService;
    private final OllamaService ollamaService;

    public SearchController(
        SearchResponseService searchResponseService,
        OllamaService ollamaService
    ) {
        this.searchResponseService = searchResponseService;
        this.ollamaService = ollamaService;
    }

    private record searchIssueRequest(@NotBlank String searchQuery, @NotBlank String repoName) {}
    private record searchQueryResponse(
        @NotBlank String repoName,
        Number numSources,
        @NotEmpty List<IssueSearchResponse> issues, 
        @NotBlank String summary
    ) {}
    
    /**
     * rag search query endpoint
     * @param request - contains the userId and searchQuery text
     * @return a response message json with all the sources used by the llm and the 2-3 sentence summary
     * @throws TranslateException from searcRepository.findRelevantIssues
     * @throws IOException from both userRepoRepository.findRepoDependency and searchRepository.findRelevantIssues
     */
    @PostMapping("/new")
    @Observed(name="search.new.controller")
    public ResponseEntity<?> searchRelevantIssues(
        @RequestBody @Valid searchIssueRequest reqBody, @AuthenticationPrincipal String username
    ) throws TranslateException, IOException {
        log.info("recieved hybrid search request", kv("query", reqBody.searchQuery));

        List<IssueSearchResponse> issueCandidates = searchResponseService.generateIssueCandidates(
            username, reqBody.searchQuery(), reqBody.repoName()
        );
        if (issueCandidates.isEmpty()) {
            return ResponseEntity.status(404).body("No dependencies found for the userId");
        }

        List<IssueSearchResponse> rerankedResults = searchResponseService.rerankCandidates(reqBody.searchQuery(), issueCandidates);
        List<IssueSearchResponse> refilteredResults = RerankingService.refilterRelevance(rerankedResults);
        if (refilteredResults.size() == 0) {
            log.debug("0 issue results after refiltering");
            return ResponseEntity.ok().body(new searchQueryResponse(
                reqBody.repoName(), refilteredResults.size(), refilteredResults, "No relevant issues were found for your query."
            ));
        }
        
        String finalResponse = ollamaService.generateFinalResponse(reqBody.searchQuery(), refilteredResults);

        return ResponseEntity.ok().body(new searchQueryResponse(reqBody.repoName(), refilteredResults.size(), refilteredResults, finalResponse));
    }
}
