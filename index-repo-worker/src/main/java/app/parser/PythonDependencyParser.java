package app.parser;

import java.util.ArrayList;
import java.util.HashMap;
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
public class PythonDependencyParser implements DependencyParserStrategy {
    
    @Override
    public Set<DependencyFileEnum> supports() {
        return Set.of(
            DependencyFileEnum.POETRY_LOCK, 
            DependencyFileEnum.UV_LOCK,
            DependencyFileEnum.PYPROJECT_TOML,
            DependencyFileEnum.REQUIREMENTS_TXT
        );
    }

    @Override
    public List<Dependency> parse(
        DependencyFileEnum fileType, String nonLockFileContent, @Nullable String lockFileContent
    ) throws JsonProcessingException {
        return switch (fileType) {
            case DependencyFileEnum.PYPROJECT_TOML   -> parsePyProjectTOML(nonLockFileContent, lockFileContent);
            case DependencyFileEnum.REQUIREMENTS_TXT -> parseRequirementsTXT(nonLockFileContent);
            default -> throw new UnsupportedOperationException();
        };
    }

    private static List<Dependency> parsePyProjectTOML(
        String nonLockFileContent, @Nullable String lockFileContent
    ) throws JsonProcessingException {
        List<String> depNames = extractDepNames(nonLockFileContent);

        if (lockFileContent == null || lockFileContent.isBlank()) {
            return depNames.stream()
                .map(name -> new Dependency(name, null))
                .toList();
        }

        Map<String, String> depVersions = extractLockFileDepVersions(lockFileContent);

        return depNames.stream()
            .map(name -> new Dependency(name, depVersions.get(name)))
            .toList();
    }

    private static List<String> extractDepNames(String nonLockFileContent) throws JsonProcessingException {
        JsonNode deps = new TomlMapper().readTree(nonLockFileContent)
            .path("project").path("dependencies");
        
        List<String> depNames = new ArrayList<>();
        for (JsonNode dep : deps) {
            String raw = dep.asText(); // parses it with version included "requests>=2.28.0"
            String name = raw.split("[><=!~\\s\\[]")[0];

            depNames.add(name);
        }

        return depNames;
    }

    private static Map<String, String> extractLockFileDepVersions(String lockFileContent) {
        Map<String, String> depVersions = new HashMap<>();

        Pattern name = Pattern.compile("^name = \"([^\"]+)\"", Pattern.MULTILINE);
        Pattern version = Pattern.compile("^version = \"([^\"]+)\"", Pattern.MULTILINE);

        for (String block : lockFileContent.split("\\[\\[package\\]\\]")) {
            Matcher nameMatcher = name.matcher(block);
            Matcher versionMatcher = version.matcher(block);

            if (nameMatcher.find() && versionMatcher.find()) {
                depVersions.put(nameMatcher.group(1), versionMatcher.group(1));
            }
        }

        return depVersions;
    }

    private static List<Dependency> parseRequirementsTXT(String nonLockFileContent) {
        List<Dependency> depList = new ArrayList<>();

        Pattern line = Pattern.compile("^([A-Za-z0-9]([A-Za-z0-9._-]*[A-Za-z0-9])?)[^\\d\\n#]*([\\d][^\\s,#\\n]*)?", Pattern.MULTILINE);
        Matcher matcher = line.matcher(nonLockFileContent);

        while (matcher.find()) {
            depList.add(new Dependency(matcher.group(1), matcher.group(3)));
        }
        
        return depList;
    }
}
