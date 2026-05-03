package app.component.parser;

import java.util.List;
import java.util.Set;

import com.fasterxml.jackson.core.JsonProcessingException;

import app.model.DependencyFileEnum;
import edu.umd.cs.findbugs.annotations.Nullable;

public interface DependencyParserStrategy {
    public record Dependency(String name, String version, String repoName) {};

    List<Dependency> parse(
        DependencyFileEnum fileType, 
        String nonLockFileContent, 
        @Nullable String lockFileContent
    ) throws JsonProcessingException;

    Set<DependencyFileEnum> supports();
}
