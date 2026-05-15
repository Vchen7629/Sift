package app.service.processing;

import static org.junit.jupiter.api.Assertions.assertArrayEquals;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.mockito.ArgumentMatchers.anyList;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Set;
import java.util.stream.Stream;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import ai.djl.inference.Predictor;
import ai.djl.repository.zoo.ZooModel;
import ai.djl.translate.TranslateException;
import app.dto.GithubChangeLogResponse;
import app.dto.IndexableDocuments;
import app.dto.ProcessedGithubIssue;
import jakarta.validation.ConstraintViolation;
import jakarta.validation.Validation;
import jakarta.validation.executable.ExecutableValidator;


public class TextEmbeddingServiceUT {

    @Mock private ZooModel<String, float[]> embeddingModel;
    @Mock private Predictor<String, float[]> predictor;

    private TextEmbeddingService service;
    private ExecutableValidator executableValidator;

    @BeforeEach
    void setup() throws Exception {
        MockitoAnnotations.openMocks(this);
        
        when(embeddingModel.newPredictor()).thenReturn(predictor);
        service = new TextEmbeddingService(embeddingModel);

        executableValidator = Validation.buildDefaultValidatorFactory()
            .getValidator().forExecutables();
    }
    
    static Stream<List<GithubChangeLogResponse>> invalidChangelogInputs() {
        return Stream.of(null, Collections.emptyList());
    }

    @ParameterizedTest
    @MethodSource("invalidChangelogInputs")
    void githubChangeLog_InvalidInputs(List<GithubChangeLogResponse> input) throws NoSuchMethodException {
        Set<ConstraintViolation<TextEmbeddingService>> violations = executableValidator.validateParameters(
            service, 
            TextEmbeddingService.class.getMethod("githubChangelog", List.class), 
            new Object[]{input}
        );

        assertFalse(violations.isEmpty());
    }

    static Stream<List<ProcessedGithubIssue>> invalidIssueInputs() {
        return Stream.of(null, Collections.emptyList());
    }

    @ParameterizedTest
    @MethodSource("invalidIssueInputs")
    void githubIssue_InvalidInputs(List<ProcessedGithubIssue> input) throws NoSuchMethodException {
        Set<ConstraintViolation<TextEmbeddingService>> violations = executableValidator.validateParameters(
            service, 
            TextEmbeddingService.class.getMethod("githubIssue", List.class), 
            new Object[]{input}
        );

        assertFalse(violations.isEmpty());
    }

    static Stream<Arguments> changelogFieldMappingCases() {
        return Stream.of(
            Arguments.of("coco", "3.0.0", "coco ate her treats", "http://coco.com"),
            Arguments.of("cotton", "5.0.0", "cotton ate his treats", "http://cottom.com"),
            Arguments.of("maple", "2.15.0", "maple ate her treats", "http://mapled.com")
        );
    }

    @ParameterizedTest
    @MethodSource("changelogFieldMappingCases")
    void githubChangelog_MapsAllFields(String dep, String version, String changes, String url) throws TranslateException {
        float[] embedding = {0.5f, 0.6f};
        when(predictor.batchPredict(anyList())).thenReturn(List.of(embedding));

        IndexableDocuments.ChangeLog doc = service.githubChangelog(
            List.of(new GithubChangeLogResponse(dep, version, changes, url))
        ).get(0);

        assertEquals(dep,     doc.dependencyName());
        assertEquals(version, doc.version());
        assertEquals(changes, doc.changes());
        assertEquals(url,     doc.url());
        assertArrayEquals(embedding, doc.changeEmbedding());
    }

    static Stream<Arguments> issueFieldMappingCases() {
        return Stream.of(
            Arguments.of("coco", "3.0.0", "coco treats", "coco ate her treats", "http://coco.com", List.of("puppy"), "2024-06-01"),
            Arguments.of("cotton", "5.0.0", "cotton treats", "cotton ate his treats", "http://cottom.com", List.of("dog"), "2025-06-01"),
            Arguments.of("maple", "2.15.0", "maple treats", "maple ate her treats", "http://mapled.com", List.of(), "2024-09-01")
        );
    }

    @ParameterizedTest
    @MethodSource("issueFieldMappingCases")
    void githubIssue_MapsAllFields(
        String dep, String version, String title, String body, String url, List<String> labels, String createdOn
    ) throws TranslateException {
        float[] titleEmb = {0.1f}, bodyEmb = {0.2f};
        when(predictor.batchPredict(anyList()))
            .thenReturn(List.of(titleEmb))
            .thenReturn(List.of(bodyEmb));

        IndexableDocuments.Issue doc = service.githubIssue(
            List.of(new ProcessedGithubIssue(dep, version, title, body, url, labels, createdOn))
        ).get(0);
        
        assertEquals(dep,       doc.dependencyName());
        assertEquals(version,   doc.version());
        assertEquals(title,     doc.title());
        assertEquals(body,      doc.body());
        assertEquals(url,       doc.url());
        assertEquals(labels,    doc.labelList());
        assertEquals(createdOn, doc.createdOn());
        assertArrayEquals(titleEmb, doc.titleEmbedding());
        assertArrayEquals(bodyEmb,  doc.bodyEmbedding());
    }

    @Test
    void githubChangelog_PreservesOrdering() throws TranslateException {
        float[] embA = {0.1f}, embB = {0.3f};
        when(predictor.batchPredict(anyList())).thenReturn(List.of(embA, embB));

        List<IndexableDocuments.ChangeLog> result = service.githubChangelog(List.of(
            new GithubChangeLogResponse("dep-a", "1.0", "fix bug",     "http://a.com"),
            new GithubChangeLogResponse("dep-b", "2.0", "add feature", "http://b.com")
        ));

        assertEquals("dep-a", result.get(0).dependencyName());
        assertEquals("dep-b", result.get(1).dependencyName());
        assertArrayEquals(embA, result.get(0).changeEmbedding());
        assertArrayEquals(embB, result.get(1).changeEmbedding());
    }

    @Test
    void githubIssue_PreservesOrdering() throws TranslateException {
        float[] titleEmbA = {0.1f}, titleEmbB = {0.2f}, bodyEmbA = {0.3f}, bodyEmbB = {0.4f};
        when(predictor.batchPredict(anyList()))
            .thenReturn(List.of(titleEmbA, titleEmbB))
            .thenReturn(List.of(bodyEmbA, bodyEmbB));

        List<IndexableDocuments.Issue> result = service.githubIssue(List.of(
            new ProcessedGithubIssue("dep-a", "1.0", "Issue A", "body a", "http://a.com", List.of(), "2024-01-01"),
            new ProcessedGithubIssue("dep-b", "2.0", "Issue B", "body b", "http://b.com", List.of(), "2024-01-02"))
        );

        assertEquals("dep-a", result.get(0).dependencyName());
        assertEquals("dep-b", result.get(1).dependencyName());
        assertArrayEquals(titleEmbA, result.get(0).titleEmbedding());
        assertArrayEquals(titleEmbB, result.get(1).titleEmbedding());
        assertArrayEquals(bodyEmbA, result.get(0).bodyEmbedding());
        assertArrayEquals(bodyEmbB, result.get(1).bodyEmbedding());
    }

    static Stream<Arguments> batchingCases() {
        return Stream.of(
            Arguments.of(1,  1),  // single item
            Arguments.of(32, 1),  // exactly one full batch
            Arguments.of(33, 2),  // one full + one partial
            Arguments.of(64, 2),  // exactly two full batches
            Arguments.of(65, 3)   // two full + one partial
        );
    }

    @ParameterizedTest
    @MethodSource("batchingCases")
    void githubChangelog_BatchesCorrectly(int inputSize, int expectedBatches) throws TranslateException {
        List<GithubChangeLogResponse> input = new ArrayList<>();
        for (int i = 0; i < inputSize; i++) {
            input.add(new GithubChangeLogResponse("dep-" + i, "1.0", "change", "http://" + i));
        }
        when(predictor.batchPredict(anyList())).thenAnswer(inv -> {
            List<?> batch = inv.getArgument(0);
            return batch.stream().map(s -> new float[]{0.1f}).toList();
        });

        service.githubChangelog(input);

        verify(predictor, times(expectedBatches)).batchPredict(anyList());
    }
}
