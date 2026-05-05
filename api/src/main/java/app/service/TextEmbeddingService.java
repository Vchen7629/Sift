package app.service;

import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.stereotype.Service;
import org.springframework.validation.annotation.Validated;

import ai.djl.repository.zoo.ZooModel;
import ai.djl.translate.TranslateException;
import io.micrometer.observation.annotation.Observed;

@Service
@Validated
public class TextEmbeddingService {
    private final ZooModel<String, float[]> embeddingModel;

    public TextEmbeddingService(@Qualifier("embeddingModel") ZooModel<String, float[]> embeddingModel) {
        this.embeddingModel = embeddingModel;
    }

    @Observed(name="textembedding.embedtext.service")
    public float[] embedText(String text) throws TranslateException {
        try (var predictor = embeddingModel.newPredictor()) {
            float[] textEmbedding = predictor.predict(text);

            return textEmbedding;
        }
    }
}
