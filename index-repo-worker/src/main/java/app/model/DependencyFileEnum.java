package app.model;

import java.util.Arrays;
import java.util.List;
import java.util.Set;
import java.util.stream.Collectors;

public enum DependencyFileEnum {
    // Java support
    POM_XML("pom.xml", "java", false),
    BUILD_GRADLE("build.gradle", "java", false),
    BUILD_GRADLE_KTS("build.gradle.kts", "java", false),

    // Python support
    POETRY_LOCK("poetry.lock", "python", true),
    UV_LOCK("uv.lock", "python", true),
    PYPROJECT_TOML("pyproject.toml", "python", false),
    REQUIREMENTS_TXT("requirements.txt", "python", false),

    // Go Support
    GO_MOD("go.mod", "go", false),

    // JS/TS Support
    PACKAGE_LOCK("package-lock.json", "javascript", true),
    PACKAGE_JSON("package.json", "javascript", false),

    // Rust Support
    CARGO_LOCK("Cargo.lock", "rust", true),
    CARGO_TOML("Cargo.toml", "rust", false);

    public final String path;
    public final String language;
    public final boolean isLockFile;

    DependencyFileEnum(String path, String language, boolean isLockFile) {
        this.path = path;
        this.language = language;
        this.isLockFile = isLockFile;
    }

    /**
     * fetch all languages to check dependency files for
     * @return set of language strings
     */
    public static Set<String> getUniqueLanguages() {
        return Arrays.stream(values())
            .map(f -> f.language)
            .collect(Collectors.toSet());
    }

    /**
        Fetch non lock file dependency files to check for language
        @return list of files to check
    */
    public static List<DependencyFileEnum> getNonLockFilesForLanguage(String language) {
        return Arrays.stream(values())
            .filter(f -> f.language.equals(language))
            .filter(f -> !f.isLockFile)
            .toList();
    }

    /**
     * Maps the non lock file to it's lock file equivalent for exact version matching
     * @return a list of the lock file equivalent to try and fetch for exact dep versions
     */
    public List<DependencyFileEnum> getLockFiles() {
        return switch (this) {
            case PYPROJECT_TOML -> List.of(POETRY_LOCK, UV_LOCK); 
            case PACKAGE_JSON   -> List.of(PACKAGE_LOCK);
            case CARGO_TOML     -> List.of(CARGO_LOCK);
            default             -> List.of();
        };
    }
}
