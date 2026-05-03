package app.component.parser;

import java.util.ArrayList;
import java.util.List;
import java.util.Set;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import org.springframework.stereotype.Component;

import app.model.DependencyFileEnum;
import edu.umd.cs.findbugs.annotations.Nullable;

@Component
public class GoDependency implements DependencyParserStrategy {
    
    @Override
    public Set<DependencyFileEnum> supports() {
        return Set.of(DependencyFileEnum.GO_MOD);
    }

    @Override
    public List<Dependency> parse(
        DependencyFileEnum fileType, String nonLockFileContent, @Nullable String lockFileContent
    ) {
        return switch (fileType) {
            case DependencyFileEnum.GO_MOD -> parseGoMod(nonLockFileContent);
            default -> throw new UnsupportedOperationException();
        };
    }

    private static List<Dependency> parseGoMod(String content) {
        List<Dependency> deps = new ArrayList<>();

        Pattern block = Pattern.compile("^require \\(([^)]+)\\)", Pattern.MULTILINE); // only matches first direct block
        Matcher matcher = block.matcher(content);

        if (!matcher.find()) {
            return deps;
        }

        for (String line : matcher.group(1).split("\n")) {
            line = line.trim();
            if (line.isEmpty() || line.contains("// indirect")) continue;

            String[] parts = line.split("\\s+");
            
            String modulePath = parts[0];
            if (!modulePath.startsWith("github.com/")) continue;

            String[] segments = modulePath.substring("github.com/".length()).split("/");
            if (segments.length < 2) continue;

            String repoName = segments[0] + "/" + segments[1];
            String version = parts[1];

            deps.add(new Dependency(repoName, version, repoName));
        }

        return deps;
    }
}
