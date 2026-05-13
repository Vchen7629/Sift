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
import io.micrometer.observation.annotation.Observed;
import jakarta.annotation.PostConstruct;
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

    @PostConstruct
    private void ollamaHealthCheck() throws IOException {
        try {
            ollamaRestClient.get().uri("/").retrieve().toBodilessEntity();
            log.info("Ollama reachable at startup");
        } catch (Exception e) {
            log.warn("ollama not reachable at startup", kv("err", e.getMessage()));
        }
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

    private static final String OLLAMA_MODEL = "qwen3:4b-q4_K_M";
    private static final float LLM_RELEVANCE_THRESHOLD = 0.3f;

    private record OllamaMessage(String role, String content) {}
    private record OllamaChatRequest(String model, List<OllamaMessage> messages, boolean stream, Map<String, Object> options) {}
    private record OllamaChatResponseMessage(String role, String content) {}
    private record OllamaChatResponse(OllamaChatResponseMessage message) {}

    @Observed(name = "searchresponse.generatefinalresponse.service")
    public String generateFinalResponse(@NotBlank String searchQuery, @NotEmpty List<IssueSearchResponse> rerankedIssues) {
        if (rerankedIssues.isEmpty()) {
            log.debug("no issues above relevance threshold, skipping llm call",
                kv("threshold", LLM_RELEVANCE_THRESHOLD),
                kv("maxScore", rerankedIssues.stream().mapToDouble(IssueSearchResponse::rerankScore).max().orElse(0)));
            return "No relevant issues were found for your query.";
        }
        
        String prompt = buildLLMPrompt(searchQuery, rerankedIssues);
        OllamaChatRequest request = new OllamaChatRequest(
            OLLAMA_MODEL,
            List.of(new OllamaMessage("user", prompt)),
            false,
            Map.of("temperature", 0.1, "think", false)
        );

        long start = System.currentTimeMillis();
        log.debug("sending LLM request", kv("model", OLLAMA_MODEL), kv("issues", rerankedIssues.size()));

        OllamaChatResponse response = ollamaRestClient
            .post()
            .uri("/api/chat")
            .body(request)
            .retrieve()
            .body(OllamaChatResponse.class);

        if (response == null || response.message() == null) {
            log.warn("ollama returned null response");
            return "";
        }

        long elapsed = System.currentTimeMillis() - start;
        log.debug("llm response took {}s", elapsed / 1000);

        return response.message().content();
    }

    private static final int BODY_TRUNCATION_CHARS = 1200;

    private String buildLLMPrompt(String query, List<IssueSearchResponse> issues) {
        StringBuilder sb = new StringBuilder();

        sb.append("A developer has this problem:\n");
        sb.append(query).append("\n\n");
        sb.append("Relevant GitHub issues:\n\n");

        for (int i = 0; i < issues.size(); i++) {
            IssueSearchResponse issue = issues.get(i);
            String body = issue.body().length() > BODY_TRUNCATION_CHARS
                ? issue.body().substring(0, BODY_TRUNCATION_CHARS) + "..."
                : issue.body();

            sb.append("Issue ").append(i + 1).append(" (relevance: ").append(String.format("%.2f", issue.rerankScore())).append("): ").append(issue.title()).append("\n");
            sb.append(body).append("\n");
            sb.append("---\n\n");
        }

        sb.append("Based only on the issues above, directly answer the developer's question. Address them directly using 'you'.\n");
        sb.append("Prioritize issues with higher relevance scores but do not mention scores in your response. Write as direct advice in 2-3 sentences. Do not use any knowledge outside the issues above.\n");
        sb.append("If the issues do not contain enough information to give a specific cause or fix, say: \"I don't know, the retrieved issues describe similar symptoms but contain no stated fix.\"");

        return sb.toString();
    }
}
