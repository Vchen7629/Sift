package app.service;

import java.util.ArrayList;
import java.util.List;

import org.springframework.stereotype.Service;
import org.springframework.validation.annotation.Validated;

import ai.djl.repository.zoo.ZooModel;
import ai.djl.translate.TranslateException;
import app.dto.GithubChangeLogResponse;
import app.dto.IndexableDocuments;
import app.dto.IndexableDocuments.ChangeLog;
import app.dto.IndexableDocuments.Issue;
import app.dto.ProcessedGithubIssue;
import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;
import lombok.extern.slf4j.Slf4j;
import jakarta.validation.Valid;
import static net.logstash.logback.argument.StructuredArguments.kv;
import jakarta.validation.constraints.NotBlank;

@Service
@Validated
@Slf4j
public class TextEmbeddingService {
    private final ZooModel<String, float[]> embeddingModel;

    public TextEmbeddingService(ZooModel<String, float[]> embeddingModel) {
        this.embeddingModel = embeddingModel;
    }

    public List<IndexableDocuments.ChangeLog> generateChangeLogEmbeddings(
        @NotEmpty @Valid List<GithubChangeLogResponse> changeLogs,
        @NotBlank String requestId
    ) throws TranslateException {
        List<List<GithubChangeLogResponse>> batches = partition(changeLogs, 32);
        List<IndexableDocuments.ChangeLog> changeLogDocuments = new ArrayList<>();

        long start = System.currentTimeMillis();

        try (var predictor = embeddingModel.newPredictor()) {
            for (List<GithubChangeLogResponse> batch : batches) {
                List<String> dependencyNames = batch.stream().map(doc -> doc.dependencyName()).toList();
                List<String> versions = batch.stream().map(doc -> doc.version()).toList();
                List<String> changeList = batch.stream().map(doc -> doc.changes()).toList();
                List<String> urls = batch.stream().map(doc -> doc.url()).toList();

                List<float[]> changeEmbeddings = predictor.batchPredict(changeList);

                log.debug("changelog embeddings for batch", kv("requestId", requestId));
            
                for (int i = 0; i < urls.size(); i++) {
                    changeLogDocuments.add(
                        new ChangeLog(
                            dependencyNames.get(i), versions.get(i), changeList.get(i), urls.get(i),
                            changeEmbeddings.get(i)
                        )
                    );
                }

                log.debug("added changelog document to result list", kv("requestId", requestId));
            }
        }

        long elapsed = System.currentTimeMillis() - start;
        log.debug("processed {} changelogs in {}ms ({}s)", 
                changeLogDocuments.size(), elapsed, elapsed / 1000.0,
                kv("requestId", requestId));
    
        return changeLogDocuments;
    }
    
    public List<IndexableDocuments.Issue> generateIssueEmbeddings(
        @NotEmpty @Valid List<ProcessedGithubIssue> issueDocuments,
        @NotBlank String requestId
    ) throws TranslateException {
        List<List<ProcessedGithubIssue>> batches = partition(issueDocuments, 32);
        List<IndexableDocuments.Issue> embeddingDocuments = new ArrayList<>();

        log.debug("created {} batches", batches.size(), kv("requestId", requestId));

        long start = System.currentTimeMillis();

        try (var predictor = embeddingModel.newPredictor()) {
            for (List<ProcessedGithubIssue> batch : batches) {
                List<String> dependencyNames = batch.stream().map(doc -> doc.dependencyName()).toList();
                List<String> versions = batch.stream().map(doc -> doc.version()).toList();
                List<String> titles = batch.stream().map(doc -> doc.title()).toList();
                List<String> bodies = batch.stream().map(doc -> doc.body()).toList();
                List<String> urls = batch.stream().map(doc -> doc.url()).toList();
                List<List<String>> labelLists = batch.stream().map(doc -> doc.labelList()).toList();
                List<String> createdOnList = batch.stream().map(doc -> doc.createdOn()).toList();

                List<float[]> titleEmbeddings = predictor.batchPredict(titles);
                List<float[]> chunkEmbeddings = predictor.batchPredict(bodies);

                for (int i = 0; i < urls.size(); i++) {
                    embeddingDocuments.add(
                        new Issue(
                            dependencyNames.get(i), versions.get(i), titles.get(i), bodies.get(i),
                            urls.get(i), labelLists.get(i), createdOnList.get(i),
                            titleEmbeddings.get(i), chunkEmbeddings.get(i)
                        )
                    );
                }
            }
        }

        long elapsed = System.currentTimeMillis() - start;
        log.debug("processed {} issues in {}ms ({}s)", 
                embeddingDocuments.size(), elapsed, elapsed / 1000.0,
                kv("requestId", requestId));
    
        return embeddingDocuments;
    }

    // partition helper so we dont pass a large amount of items at a time 
    // (500+) to embeddings and overwhelm the resources
    private static <T> List<List<T>> partition(
        @NotEmpty @Valid List<T> items, 
        @NotNull @Positive int batchSize
    ) {
        List<List<T>> batches = new ArrayList<>();

        for (int i = 0; i < items.size(); i += batchSize) {
            batches.add(items.subList(i, Math.min(i + batchSize, items.size())));
        }

        return batches;
    }
}
