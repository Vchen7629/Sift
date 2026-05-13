package app.dto;

import java.util.List;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.NotNull;

public class IndexableDocuments {
    public interface Base {
        @NotBlank String url();
    }

    public static record ChangeLog (
        @NotBlank String dependencyName,
        @NotBlank String version,
        @NotBlank String changes,
        @NotBlank String url,
        @NotEmpty float[] changeEmbedding
    ) implements Base {}

    public static record Issue (
        @NotBlank String dependencyName,
        @NotBlank String version,
        @NotBlank String title,
        @NotBlank String body,
        @NotBlank String url,
        @NotNull List<String> labelList,
        @NotNull String createdOn,
        @NotEmpty float[] titleEmbedding,
        @NotEmpty float[] bodyEmbedding
    ) implements Base {}
}
