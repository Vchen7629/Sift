package app.parser;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

import org.springframework.stereotype.Component;
import org.springframework.validation.annotation.Validated;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

import app.model.DependencyFileEnum;
import edu.umd.cs.findbugs.annotations.Nullable;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;

@Component
@Validated
public class JsTsDependencyParser implements DependencyParserStrategy {
    
    @Override
    public Set<DependencyFileEnum> supports() {
        return Set.of(
            DependencyFileEnum.PACKAGE_LOCK,
            DependencyFileEnum.PACKAGE_JSON
        );
    }

    @Override
    public List<Dependency> parse(
        DependencyFileEnum fileType, @NotBlank String nonLockFileContent, @Nullable String lockFileContent
    ) throws JsonProcessingException {
        return switch (fileType) {
            case DependencyFileEnum.PACKAGE_JSON -> parsePackageJson(nonLockFileContent, lockFileContent);
            default -> throw new UnsupportedOperationException();
        };
    }

    private static List<Dependency> parsePackageJson(
        String nonLockFileContent, @Nullable String lockFileContent
    ) throws JsonProcessingException {
        List<String> depNames = extractDepNames(nonLockFileContent);

        if (lockFileContent == null || lockFileContent.isBlank()) {
            return depNames.stream()
                .map(name -> new Dependency(name, null))
                .toList();
        }

        Map<String, String> depVersions = extractLockFileDepVersions(lockFileContent, depNames);

        return depNames.stream()
            .map(name -> new Dependency(name, depVersions.get(name)))
            .toList();
    }

    private static List<String> extractDepNames(String nonLockFileContent) throws JsonProcessingException {
        ObjectMapper mapper = new ObjectMapper();
        JsonNode deps = mapper.readTree(nonLockFileContent).path("dependencies");
        
        List<String> depNames = new ArrayList<>();
        deps.fieldNames().forEachRemaining(depNames::add);

        return depNames;
    }

    private static Map<String, String> extractLockFileDepVersions(
        String lockFileContent, @NotEmpty List<String> depNames
    ) throws JsonProcessingException {
        Map<String, String> depVersions = new HashMap<>();

        ObjectMapper mapper = new ObjectMapper();
        JsonNode packages = mapper.readTree(lockFileContent).path("packages");

        for (var entry : packages.properties()) {
            String key = entry.getKey();
            
            if (key.startsWith("node_modules/")) {
                String name = key.substring("node_modules/".length());
                if (depNames.contains(name)) {
                    String version = entry.getValue().path("version").asText(null);
                    depVersions.put(name, version);
                }
            }
        }

        return depVersions;
    }
}
