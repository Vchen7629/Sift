package app.parser;

import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import org.springframework.stereotype.Component;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.dataformat.toml.TomlMapper;

import app.model.DependencyFileEnum;
import edu.umd.cs.findbugs.annotations.Nullable;

@Component
public class RustDependencyParser implements DependencyParserStrategy {
    
    @Override
    public Set<DependencyFileEnum> supports() {
        return Set.of(
            DependencyFileEnum.CARGO_TOML
        );
    }

    @Override
    public List<Dependency> parse(
        DependencyFileEnum fileType, String nonLockFileContent, @Nullable String lockFileContent
    ) throws JsonProcessingException {
        return switch (fileType) {
            case DependencyFileEnum.CARGO_TOML -> parseCargoToml(nonLockFileContent, lockFileContent);
            default -> throw new UnsupportedOperationException();
       };
    }

    private static List<Dependency> parseCargoToml(
        String nonLockFileContent, @Nullable String lockFileContent
    ) throws JsonProcessingException {
        Set<String> depNames = extractDepNames(nonLockFileContent);

        if (lockFileContent == null || lockFileContent.isBlank()) {
            return depNames.stream()
                .map(name -> new Dependency(name, null))
                .toList();
        }

        Map<String, String> depVersions = extractLockFileDepVersions(lockFileContent, depNames);
        
        return depVersions.entrySet().stream()
            .map(e -> new Dependency(e.getKey(), e.getValue()))
            .toList();
    }

    private static Set<String> extractDepNames(String nonLockFileContent) throws JsonProcessingException {
        JsonNode deps = new TomlMapper().readTree(nonLockFileContent).path("dependencies");
        
        Set<String> depNames = new HashSet<>();
        deps.fieldNames().forEachRemaining(depNames::add);

        return depNames;
    }

    private static Map<String, String> extractLockFileDepVersions(String lockFileContent, Set<String> depNames) {
        Map<String, String> depVersions = new HashMap<>();

        Pattern name = Pattern.compile("^name = \"([^\"]+)\"", Pattern.MULTILINE);
        Pattern version = Pattern.compile("^version = \"([^\"]+)\"", Pattern.MULTILINE);

        for (String block : lockFileContent.split("\\[\\[package\\]\\]")) {
            Matcher nameMatcher = name.matcher(block);
            Matcher versionMatcher = version.matcher(block);

            if (nameMatcher.find() && versionMatcher.find()) {
                String depName = nameMatcher.group(1);

                if (depNames.contains(depName)) {
                    depVersions.put(nameMatcher.group(1), versionMatcher.group(1));
                }
            }
        }

        return depVersions;
    }
}
