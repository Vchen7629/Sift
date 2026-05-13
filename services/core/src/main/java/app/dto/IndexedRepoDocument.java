package app.dto;

import java.util.List;
import java.util.Map;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

@JsonIgnoreProperties(ignoreUnknown = true)
public record IndexedRepoDocument(
    String repoName, String lastIndexed, int totalDependencies, List<DependencyDocument> dependencies
) {
    @JsonIgnoreProperties(ignoreUnknown = true)
    public record Source(String repoName, String lastIndexed, Map<String, String> dependencies) {}

    public static IndexedRepoDocument from(Source source) {
        List<DependencyDocument> deps = source.dependencies() != null 
            ? source.dependencies().entrySet().stream()
                .map(e -> new DependencyDocument(e.getKey(), e.getValue(), "healthy"))
                .toList()
            : List.of();

        return new IndexedRepoDocument(source.repoName(), source.lastIndexed(), deps.size(), deps);
    }
}