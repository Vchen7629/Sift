package app.service;

import java.util.Comparator;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.stream.IntStream;

import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.stereotype.Service;

import ai.djl.repository.zoo.ZooModel;
import ai.djl.translate.TranslateException;
import ai.djl.util.StringPair;
import app.dto.IssueSearchResponse;

@Service
public class RerankingService {
    private final ZooModel<StringPair, float[]> crossEncoderModel;

    public RerankingService(
        @Qualifier("crossEncoderModel") ZooModel<StringPair, float[]> crossEncoderModel
    ) {
        this.crossEncoderModel = crossEncoderModel;
    }

    private final ExecutorService rerankPool = Executors.newFixedThreadPool(
        Runtime.getRuntime().availableProcessors() / 2
    );

    public List<IssueSearchResponse> rerank(String searchQuery, List<IssueSearchResponse> searchRes) {
        List<CompletableFuture<float[]>> futures = searchRes.stream()
            .map(issue -> CompletableFuture.supplyAsync(() -> {
                try (var predictor = crossEncoderModel.newPredictor()) {
                    return predictor.predict(new StringPair(searchQuery, issue.title() + " " + issue.body()));
                } catch (TranslateException e) {
                    throw new RuntimeException(e);
                }
            }, rerankPool))
            .toList();
        
        List<float[]> scores = futures.stream().map(CompletableFuture::join).toList();

        float filterThreshold = 0.05f;
        
        return IntStream.range(0, searchRes.size())
            .boxed()
            .filter(i -> scores.get(i)[0] >= filterThreshold)
            .sorted(Comparator.comparingDouble(i -> -scores.get(i)[0]))
            .limit(10)
            .map(i -> searchRes.get(i))
            .toList();
    }
}
