package app.service;

import java.util.ArrayList;
import java.util.List;

import org.springframework.stereotype.Service;
import org.springframework.validation.annotation.Validated;

import ai.djl.repository.zoo.ZooModel;
import ai.djl.translate.TranslateException;
import app.service.githubRepo.ChangelogService;
import app.service.githubRepo.IssueService;
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

    public interface IndexableDocument {
        @NotBlank String url();
    }

    public static record ChangeLogDocument(
        @NotBlank String library_name,
        @NotBlank String version,
        @NotBlank String changes,
        @NotBlank String url,
        @NotEmpty float[] changeEmbedding
    ) implements IndexableDocument {}

    public List<ChangeLogDocument> generateChangeLogEmbeddings(
        @NotEmpty @Valid List<ChangelogService.Result> changeLogs,
        @NotBlank String requestId
    ) throws TranslateException {
        List<List<ChangelogService.Result>> batches = partition(changeLogs, 32);
        List<ChangeLogDocument> changeLogDocuments = new ArrayList<>();

        long start = System.currentTimeMillis();

        try (var predictor = embeddingModel.newPredictor()) {
            for (List<ChangelogService.Result> batch : batches) {
                List<String> libraryNames = batch.stream().map(doc -> doc.libraryName()).toList();
                List<String> versions = batch.stream().map(doc -> doc.version()).toList();
                List<String> changeList = batch.stream().map(doc -> doc.changes()).toList();
                List<String> urls = batch.stream().map(doc -> doc.url()).toList();

                List<float[]> changeEmbeddings = predictor.batchPredict(changeList);

                log.debug("changelog embeddings for batch", kv("requestId", requestId));
            
                for (int i = 0; i < urls.size(); i++) {
                    changeLogDocuments.add(
                        new ChangeLogDocument(
                            libraryNames.get(i), versions.get(i), changeList.get(i), urls.get(i),
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


    public static record IssueDocument(
        @NotBlank String repoName, 
        @NotBlank String url, 
        @NotBlank String title, 
        @NotBlank String body,
        @NotNull List<String> labelList,
        @NotEmpty float[] titleEmbedding,
        @NotEmpty float[] bodyEmbedding
    ) implements IndexableDocument {} 

    
    public List<IssueDocument> generateIssueEmbeddings(
        @NotEmpty @Valid List<IssueService.Result> issueDocuments,
        @NotBlank String requestId
    ) throws TranslateException {
        List<List<IssueService.Result>> batches = partition(issueDocuments, 32);
        List<IssueDocument> embeddingDocuments = new ArrayList<>();

        log.debug("created {} batches", batches.size(), kv("requestId", requestId));

        long start = System.currentTimeMillis();

        try (var predictor = embeddingModel.newPredictor()) {
            for (List<IssueService.Result> batch : batches) {
                List<String> repoNames = batch.stream().map(doc -> doc.repoName()).toList();
                List<String> urls = batch.stream().map(doc -> doc.url()).toList();
                List<String> titles = batch.stream().map(doc -> doc.title()).toList();
                List<String> bodies = batch.stream().map(doc -> doc.body()).toList();
                List<List<String>> labelLists = batch.stream().map(doc -> doc.labelList()).toList();
                
                List<float[]> titleEmbeddings = predictor.batchPredict(titles);
                List<float[]> bodyEmbeddings = predictor.batchPredict(bodies);

                log.debug("embeddings for batch", kv("requestId", requestId));

                for (int i = 0; i < urls.size(); i++) {
                    embeddingDocuments.add(
                        new IssueDocument(
                            repoNames.get(i), urls.get(i), titles.get(i), bodies.get(i), labelLists.get(i),
                            titleEmbeddings.get(i), bodyEmbeddings.get(i)
                        )
                    );
                }

                log.debug("added embedding document to result list", kv("requestId", requestId));
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
