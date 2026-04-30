package app.service;

import org.springframework.stereotype.Service;
import org.springframework.validation.annotation.Validated;

import ai.djl.repository.zoo.ZooModel;
import ai.djl.translate.TranslateException;

@Service
@Validated
public class TextEmbeddingService {
    private final ZooModel<String, float[]> embeddingModel;

    public TextEmbeddingService(ZooModel<String, float[]> embeddingModel) {
        this.embeddingModel = embeddingModel;
    }

    public float[] embedText(String text) throws TranslateException {
        try (var predictor = embeddingModel.newPredictor()) {
            float[] textEmbedding = predictor.predict(text);

            return textEmbedding;
        }
    }
}
