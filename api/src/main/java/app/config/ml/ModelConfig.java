package app.config.ml;

import ai.djl.MalformedModelException;
import ai.djl.huggingface.translator.TextEmbeddingTranslatorFactory;
import ai.djl.huggingface.tokenizers.HuggingFaceTokenizer;
import ai.djl.huggingface.translator.CrossEncoderTranslator;
import ai.djl.repository.zoo.Criteria;
import ai.djl.repository.zoo.ModelNotFoundException;
import ai.djl.repository.zoo.ZooModel;
import ai.djl.training.util.ProgressBar;
import ai.djl.util.StringPair;

import java.io.IOException;
import java.util.Map;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class ModelConfig {
    static {
        System.setProperty("ai.djl.pytorch.num_intraop_threads", "2");
        System.setProperty("ai.djl.pytorch.num_interop_threads", "1");
    }

    private static final String EMB_MODEL_NAME = "BAAI/bge-base-en-v1.5";

    @Bean(name="embeddingModel", destroyMethod = "close")
    public static ZooModel<String, float[]> embeddingModel() throws ModelNotFoundException, MalformedModelException, IOException {
        Criteria<String, float[]> criteria = Criteria.builder()
            .setTypes(String.class, float[].class)
            .optModelUrls("djl://ai.djl.huggingface.pytorch/" + EMB_MODEL_NAME)
            .optEngine("PyTorch")
            .optTranslatorFactory(new TextEmbeddingTranslatorFactory())
            .optProgress(new ProgressBar())
            .build();

        return criteria.loadModel();
    }

    private static final String CROSS_ENC_MODEL_NAME = "BAAI/bge-reranker-v2-m3";

    @Bean(name="crossEncoderModel", destroyMethod = "close")
    public static ZooModel<StringPair, float[]> crossEncoderModel() throws ModelNotFoundException, MalformedModelException, IOException {
        HuggingFaceTokenizer tokenizer = HuggingFaceTokenizer.newInstance(
            CROSS_ENC_MODEL_NAME, Map.of("maxLength", "512", "truncation", "true")
        );
        
        CrossEncoderTranslator translator = CrossEncoderTranslator.builder(tokenizer)
            .optSigmoid(true)
            .build();

        Criteria<StringPair, float[]> criteria = Criteria.builder()
            .setTypes(StringPair.class, float[].class)
            .optModelUrls("djl://ai.djl.huggingface.pytorch/" + CROSS_ENC_MODEL_NAME)
            .optEngine("PyTorch")
            .optTranslator(translator)
            .optProgress(new ProgressBar())
            .build();

        return criteria.loadModel();
    }
}
