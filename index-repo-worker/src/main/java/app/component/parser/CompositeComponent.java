package app.component.parser;

import java.util.List;
import java.util.Set;
import java.util.stream.Collectors;

import org.springframework.context.annotation.Primary;
import org.springframework.stereotype.Component;

import com.fasterxml.jackson.core.JsonProcessingException;

import app.model.DependencyFileEnum;
import edu.umd.cs.findbugs.annotations.Nullable;
import io.micrometer.observation.annotation.Observed;

@Primary
@Component
public class CompositeComponent implements DependencyParserStrategy {

    private final List<DependencyParserStrategy> parsers;

    public CompositeComponent(List<DependencyParserStrategy> parsers) {
        this.parsers = parsers;
    }

    @Override
    public Set<DependencyFileEnum> supports() {
        return parsers.stream()
            .flatMap(p -> p.supports().stream())
            .collect(Collectors.toSet());
    }

    // Todo: Investigate if i can observe the individual parsers
    @Override
    @Observed(name="composite.parse.component")
    public List<Dependency> parse(
        DependencyFileEnum fileType, String nonLockFileContent, @Nullable String lockFileContent
    ) throws JsonProcessingException {
        return parsers.stream()
            .filter(p -> p.supports().contains(fileType)) // finds right parser for the filetype
            .findFirst()
            .orElseThrow(() -> new UnsupportedOperationException("No parser found for " + fileType))
            .parse(fileType, nonLockFileContent, lockFileContent);
    }
}
