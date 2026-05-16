package app.service.internal;

import java.io.IOException;
import java.util.List;
import org.springframework.stereotype.Service;

import ai.djl.translate.TranslateException;
import app.dto.IssueSearchResponse;
import app.repository.SearchRepository;
import app.repository.UserRepoRepository;
import app.service.ml.RerankingService;
import io.micrometer.observation.annotation.Observed;
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

    public SearchResponseService(
        UserRepoRepository userRepoRepository,
        SearchRepository searchRepository,
        RerankingService rerankingService
    ) {
        this.userRepoRepository = userRepoRepository;
        this.searchRepository = searchRepository;
        this.rerankingService = rerankingService;
    }

    @Observed(name = "searchresponse.generateissuecandidates.service")
    public List<IssueSearchResponse> generateIssueCandidates(
        @NotBlank String userId, @NotBlank String searchQuery, @NotBlank String repoName
    ) throws TranslateException, IOException {
        List<String> repoDependencies = userRepoRepository.listAllRepoDependencies(userId, repoName);

        if (repoDependencies.isEmpty()) {
            log.debug("no dependencies found for repo", kv("userId", userId), kv("repoName", repoName));
            return List.of();
        }

        return searchRepository.findRelevantIssues(repoDependencies, searchQuery);
    }

    @Observed(name = "searchresponse.rerankcandidates.service")
    public List<IssueSearchResponse> rerankCandidates(
        @NotBlank String searchQuery, @NotEmpty List<IssueSearchResponse> issueResults
    ) {
        long start = System.currentTimeMillis();
        List<IssueSearchResponse> rerankedResults = rerankingService.rerank(searchQuery, issueResults);
        long elapsed = System.currentTimeMillis() - start;

        log.debug("reranking took {} ms", elapsed);

        return rerankedResults;
    }
}
