package app.service.githubRepo;

import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.UncheckedIOException;
import java.nio.charset.StandardCharsets;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;

import org.kohsuke.github.GHContent;
import org.kohsuke.github.GHFileNotFoundException;
import org.kohsuke.github.GHRepository;
import org.kohsuke.github.GHTree;
import org.kohsuke.github.GHTreeEntry;
import org.kohsuke.github.GitHub;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;
import org.springframework.validation.annotation.Validated;

import com.fasterxml.jackson.core.JsonProcessingException;

import app.model.DependencyFileEnum;
import app.component.parser.DependencyParserStrategy;
import app.component.parser.DependencyParserStrategy.Dependency;
import jakarta.validation.constraints.NotBlank;
import lombok.extern.slf4j.Slf4j;
import static net.logstash.logback.argument.StructuredArguments.kv;


@Service
@Validated
@Slf4j
public class DependencyService {
    private final GitHub githubClient;
    private final DependencyParserStrategy dependencyParserStrategy;

    public DependencyService(
        GitHub githubClient,
        DependencyParserStrategy dependencyParserStrategy
    ) {
        this.githubClient = githubClient;
        this.dependencyParserStrategy = dependencyParserStrategy;
    }


    @Async
    public CompletableFuture<Map<String, List<Dependency>>> fetchRepoDependencies(
        @NotBlank String repoName, @NotBlank String requestId
    ) {
        try {
            GHRepository repo = githubClient.getRepository(repoName);
            
            long start = System.currentTimeMillis();

            Map<String, List<Dependency>> dependenciesByLanguage = new HashMap<>();

            GHTree tree = repo.getTreeRecursive("HEAD", 1);
            List<String> allPaths = tree.getTree().stream()
                .map(GHTreeEntry::getPath)
                .toList();

            for (String language : DependencyFileEnum.getUniqueLanguages()) {
                Optional<List<Dependency>> languageDependencies = fetchLanguageDeps(repo, language, allPaths);
                if (languageDependencies.isEmpty()) continue;

                dependenciesByLanguage.put(language, languageDependencies.get());
            }

            long elapsed = System.currentTimeMillis() - start;

            log.debug("fetched dependencies in {}ms ({}s)", 
                elapsed, elapsed / 1000.0,
                kv("repoName", repoName),
                kv("requestId", requestId));

            return CompletableFuture.completedFuture(dependenciesByLanguage);
        } catch (GHFileNotFoundException e) {
            log.error("Repo not found", 
                kv("repoName", repoName), 
                kv("error", e.getMessage()),
                kv("requestId", requestId));

            return CompletableFuture.failedFuture(e);
        } catch (IOException e) {
            log.error("Unknown error fetching issues", 
                kv("repoName", repoName), 
                kv("error", e.getMessage()),
                kv("requestId", requestId));

            return CompletableFuture.failedFuture(e);
        }
    }

    private Optional<List<Dependency>> fetchLanguageDeps(
        GHRepository repo, String language, List<String> allPaths
    ) throws JsonProcessingException {
        Map<String, Dependency> allDependencies = new HashMap<>();

        for (DependencyFileEnum nonLockFile : DependencyFileEnum.getNonLockFilesForLanguage(language)) {
            List<String> matchedPaths = allPaths.stream()
                .filter(p -> p.equals(nonLockFile.path) || p.endsWith("/" + nonLockFile.path))
                .toList();

            for (String matchedPath : matchedPaths) {
                Optional<String> nonLockFileContent = tryFetchFileContent(repo, matchedPath);
                if (nonLockFileContent.isEmpty()) continue;

                String dir = matchedPath.contains("/")
                    ? matchedPath.substring(0, matchedPath.lastIndexOf('/') + 1)
                    : "";

                Optional<String> lockFileContent = nonLockFile.getLockFiles().stream()
                    .flatMap(lockfile -> allPaths.stream()
                        .filter(p -> p.equals(dir + lockfile.path) || p.endsWith("/" + lockfile.path)))
                    .findFirst()
                    .flatMap(lp -> tryFetchFileContent(repo, lp));

                dependencyParserStrategy.parse(nonLockFile, nonLockFileContent.get(), lockFileContent.orElse(null))                                                                
                    .forEach(dep -> allDependencies.putIfAbsent(dep.name(), dep));
            }
        }

        return allDependencies.isEmpty() ? Optional.empty() : Optional.of(List.copyOf(allDependencies.values()));
    }

    private Optional<String> tryFetchFileContent(GHRepository repo, String path) {
        try {
            GHContent ghContent = repo.getFileContent(path);

            String fileContent = new String(ghContent.read().readAllBytes(), StandardCharsets.UTF_8);

            return Optional.of(fileContent);
        } catch (FileNotFoundException e) {
            return Optional.empty();
        } catch (IOException e) {
            throw new UncheckedIOException(e);
        }
    }
}
