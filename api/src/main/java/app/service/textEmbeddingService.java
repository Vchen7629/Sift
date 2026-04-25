package app.service;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.springframework.stereotype.Service;

import ai.djl.repository.zoo.ZooModel;
import ai.djl.translate.TranslateException;

@Service
public class textEmbeddingService {
    private final ZooModel<String, float[]> embeddingModel;

    public textEmbeddingService(ZooModel<String, float[]> embeddingModel) {
        this.embeddingModel = embeddingModel;
    }

    // Combine issue title and body into one single string for better embeddings
    public static String combineDescBody(String title, String body) {
        if (title == null || title.trim().isEmpty()) {
            throw new IllegalArgumentException("issue title must not be null or empty");
        }

        if (body == null || body.trim().isEmpty()) {
            throw new IllegalArgumentException("issue body must not be null or empty");
        }

        return title + "\n" + body;
    }

    // pass in a map of issue <url, text> and return issue <url, embedding>
    // does batch embedding generation with a sentence transformer model
    public Map<String, float[]> generateEmbeddings(Map<String, String> issueUrlTexts) throws TranslateException {
        if (issueUrlTexts == null || issueUrlTexts.isEmpty()) {
            throw new IllegalArgumentException("issueTexts must not be null or empty");
        }

        List<Map<String, String>> batches = partition(issueUrlTexts, 32);
        Map<String, float[]> allEmbeddings = new HashMap<>();

        try (var predictor = embeddingModel.newPredictor()) {
            for (Map<String, String> batch : batches) {
                List<String> urls = new ArrayList<>(batch.keySet());
                List<String> texts = urls.stream().map(batch::get).toList();

                List<float[]> embeddings = predictor.batchPredict(texts);
                for (int i = 0; i < urls.size(); i++) {
                    allEmbeddings.put(urls.get(i), embeddings.get(i));
                }
                System.out.println("hi" + embeddings);
            }
        }
    
        return allEmbeddings;
    }

    // partition helper so we dont pass a large amount of issues at a time (500+) to embeddings
    // and overwhelm the resources
    private static List<Map<String, String>> partition(Map<String, String> allIssues, int batchSize) {
        List<Map.Entry<String, String>> entries = new ArrayList<>(allIssues.entrySet());
        List<Map<String, String>> batches = new ArrayList<>();

        for (int i = 0; i < entries.size(); i++) {
            Map<String,String> batch = new HashMap<>();

            entries.subList(i, Math.min(i + batchSize, entries.size()))
                .forEach(e -> batch.put(e.getKey(), e.getValue()));
            
            batches.add(batch);
        }

        return batches;
    }
}
