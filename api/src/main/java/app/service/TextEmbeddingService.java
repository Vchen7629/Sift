package app.service;

import java.util.ArrayList;
import java.util.List;

import org.springframework.stereotype.Service;
import org.springframework.validation.annotation.Validated;

import ai.djl.repository.zoo.ZooModel;
import ai.djl.translate.TranslateException;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;
import jakarta.validation.constraints.NotBlank;

@Service
@Validated
public class TextEmbeddingService {
    private final ZooModel<String, float[]> embeddingModel;

    public TextEmbeddingService(ZooModel<String, float[]> embeddingModel) {
        this.embeddingModel = embeddingModel;
    }

    public static record embeddingDocument(
        @NotBlank String repoName, 
        @NotBlank String url, 
        @NotBlank String title, 
        @NotBlank String body,
        @NotEmpty float[] titleEmbedding,
        @NotEmpty float[] bodyEmbedding
    ) {} 

    // pass in a list of <url, text> and return issue <url, embedding>
    // does batch embedding generation with a sentence transformer model
    public List<embeddingDocument> generateEmbeddings(
        @NotEmpty @Valid List<GithubApiService.IssueDocument> issueDocuments
    ) throws TranslateException {
        List<List<GithubApiService.IssueDocument>> batches = partition(issueDocuments, 32);
        List<embeddingDocument> embeddingDocuments = new ArrayList<>();

        try (var predictor = embeddingModel.newPredictor()) {
            for (List<GithubApiService.IssueDocument> batch : batches) {
                List<String> repoNames = batch.stream().map(doc -> doc.repoName()).toList();
                List<String> urls = batch.stream().map(doc -> doc.url()).toList();
                List<String> titles = batch.stream().map(doc -> doc.title()).toList();
                List<String> bodies = batch.stream().map(doc -> doc.body()).toList();
                
                List<float[]> titleEmbeddings = predictor.batchPredict(titles);
                List<float[]> bodyEmbeddings = predictor.batchPredict(bodies);

                for (int i = 0; i < urls.size(); i++) {
                    embeddingDocuments.add(
                        new embeddingDocument(
                            repoNames.get(i), urls.get(i), titles.get(i), bodies.get(i), 
                            titleEmbeddings.get(i), bodyEmbeddings.get(i)
                        )
                    );
                }
                //System.out.println("hi" + java.util.Arrays.toString(embeddings.get(0)));
            }
        }
    
        return embeddingDocuments;
    }

    // partition helper so we dont pass a large amount of issues at a time 
    // (500+) to embeddings and overwhelm the resources
    private static List<List<GithubApiService.IssueDocument>> partition(
        @NotEmpty @Valid List<GithubApiService.IssueDocument> issues, 
        @NotNull @Positive int batchSize
    ) {
        List<List<GithubApiService.IssueDocument>> batches = new ArrayList<>();

        for (int i = 0; i < issues.size(); i += batchSize) {
            batches.add(issues.subList(i, Math.min(i + batchSize, issues.size())));
        }

        return batches;
    }
}
