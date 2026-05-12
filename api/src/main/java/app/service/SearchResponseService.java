package app.service;

import java.io.IOException;
import java.util.List;
import java.util.Map;

import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;

import ai.djl.translate.TranslateException;
import app.dto.IssueSearchResponse;
import app.repository.SearchRepository;
import app.repository.UserRepoRepository;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import lombok.extern.slf4j.Slf4j;
import static net.logstash.logback.argument.StructuredArguments.kv;

@Service
@Slf4j
public class SearchResponseService {
    private final UserRepoRepository userRepoRepository;
    private final SearchRepository searchRepository;
    private final RerankingService rerankingService;
    private final RestClient ollamaRestClient;

    public SearchResponseService(
        UserRepoRepository userRepoRepository,
        SearchRepository searchRepository,
        RerankingService rerankingService,
        RestClient ollamaRestClient
    ) {
        this.userRepoRepository = userRepoRepository;
        this.searchRepository = searchRepository;
        this.rerankingService = rerankingService;
        this.ollamaRestClient = ollamaRestClient;
    }

    public List<IssueSearchResponse> generateIssueCandidates(
        @NotBlank String userId, @NotBlank String searchQuery
    ) throws TranslateException, IOException {
        Map<String, String> userRepoDependencies = userRepoRepository.listAllDependencies(userId);
        if (userRepoDependencies.isEmpty()) {
            log.debug("no dependencies found for user", kv("userId", userId));
            return List.of();
            //return ResponseEntity.status(404).body("No dependencies found for the userId");
        }

        return searchRepository.findRelevantIssues(userRepoDependencies, searchQuery);
    }

    public List<IssueSearchResponse> rerankCandidates(
        @NotBlank String searchQuery, @NotEmpty List<IssueSearchResponse> issueResults
    ) {
        long start = System.currentTimeMillis();
        List<IssueSearchResponse> rerankedResults = rerankingService.rerank(searchQuery, issueResults);
        long elapsed = System.currentTimeMillis() - start;

        log.debug("reranking took {} ms", elapsed);

        return rerankedResults;
    }

    public void generateFinalResponse() {
        
    }
}
