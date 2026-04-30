package app.config;

import ai.djl.MalformedModelException;
import ai.djl.huggingface.translator.TextEmbeddingTranslatorFactory;
import ai.djl.repository.zoo.Criteria;
import ai.djl.repository.zoo.ModelNotFoundException;
import ai.djl.repository.zoo.ZooModel;
import ai.djl.training.util.ProgressBar;

import java.io.IOException;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class ModelConfig {
    private static final String MODEL_NAME = "sentence-transformers/all-MiniLM-L6-v2";

    @Bean
    public static ZooModel<String, float[]> embeddingModel() throws ModelNotFoundException, MalformedModelException, IOException {
        Criteria<String, float[]> criteria = Criteria.builder()
            .setTypes(String.class, float[].class)
            .optModelUrls("djl://ai.djl.huggingface.pytorch/" + MODEL_NAME)
            .optEngine("PyTorch")
            .optTranslatorFactory(new TextEmbeddingTranslatorFactory())
            .optProgress(new ProgressBar())
            .build();

        return criteria.loadModel();
    }
}
