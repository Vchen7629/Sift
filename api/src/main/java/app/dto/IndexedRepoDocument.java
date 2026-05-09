package app.dto;

import java.util.Map;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

@JsonIgnoreProperties(ignoreUnknown = true)
public record IndexedRepoDocument(
    String repoName, String lastIndexed, int totalDependencies, Map<String, String> dependencies
) {
    @JsonIgnoreProperties(ignoreUnknown = true)
    public record Source(String repoName, String lastIndexed, Map<String, String> dependencies) {}

    public static IndexedRepoDocument from(Source source) {
        return new IndexedRepoDocument(
            source.repoName(),
            source.lastIndexed(),
            source.dependencies() != null ? source.dependencies().size() : 0,
            source.dependencies()
        );
    }
}