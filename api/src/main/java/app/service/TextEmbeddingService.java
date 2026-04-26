package app.service;

import java.util.ArrayList;
import java.util.List;

import org.springframework.stereotype.Service;

import ai.djl.repository.zoo.ZooModel;
import ai.djl.translate.TranslateException;

@Service
public class TextEmbeddingService {
    private final ZooModel<String, float[]> embeddingModel;

    public TextEmbeddingService(ZooModel<String, float[]> embeddingModel) {
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

    public static record embeddingDocument(String repoName, String url, String title, float[] embedding) {} 

    // pass in a list of <url, text> and return issue <url, embedding>
    // does batch embedding generation with a sentence transformer model
    public List<embeddingDocument> generateEmbeddings(List<GithubApiService.IssueDocument> issueDocuments) throws TranslateException {
        if (issueDocuments == null || issueDocuments.isEmpty()) {
            throw new IllegalArgumentException("issueDocuments must not be null or empty");
        }

        List<List<GithubApiService.IssueDocument>> batches = partition(issueDocuments, 32);
        List<embeddingDocument> embeddingDocuments = new ArrayList<>();

        try (var predictor = embeddingModel.newPredictor()) {
            for (List<GithubApiService.IssueDocument> batch : batches) {
                List<String> repoNames = batch.stream().map(doc -> doc.repoName()).toList();
                List<String> urls = batch.stream().map(doc -> doc.url()).toList();
                List<String> titles = batch.stream().map(doc -> doc.title()).toList();
                List<String> texts = batch.stream().map(doc -> doc.text()).toList();
                
                List<float[]> embeddings = predictor.batchPredict(texts);
                for (int i = 0; i < urls.size(); i++) {
                    embeddingDocuments.add(new embeddingDocument(repoNames.get(i), urls.get(i), titles.get(i), embeddings.get(i)));
                }
                System.out.println("hi" + java.util.Arrays.toString(embeddings.get(0)));
            }
        }
    
        return embeddingDocuments;
    }

    // partition helper so we dont pass a large amount of issues at a time (500+) to embeddings and overwhelm the resources
    private static List<List<GithubApiService.IssueDocument>> partition(List<GithubApiService.IssueDocument> issues, int batchSize) {
        List<List<GithubApiService.IssueDocument>> batches = new ArrayList<>();

        for (int i = 0; i < issues.size(); i += batchSize) {
            batches.add(issues.subList(i, Math.min(i + batchSize, issues.size())));
        }

        return batches;
    }
}
