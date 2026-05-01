package app.model;

public enum DependencyFileEnum {
    // Java support
    POM_XML("pom.xml", "java", true),
    BUILD_GRADLE("build.gradle", "java", true),
    BUILD_GRADLE_KTS("build.gradle.kts", "java", true),

    // Python support
    POETRY_LOCK("poetry.lock", "python", true),
    UV_LOCK("uv.lock", "python", true),
    PYPROJECT_TOML("pyproject.toml", "python", false),
    REQUIREMENT_TXT("requirements.txt", "python", false),
    SETUP_PY("setup.py", "python", false),

    // Go Support
    GO_SUM("go.sum", "go", true),
    GO_MOD("go.mod", "go", false),

    // JS/TS Support
    PACKAGE_LOCK("package-lock.json", "javascript", true),
    PACKAGE_JSON("package.json", "javascript", false),

    // Rust Support
    CARGO_LOCK("Cargo.lock", "rust", true),
    CARGO_TOML("Cargo.toml", "rust", false);

    public final String path;
    public final String ecosystem;
    public final Boolean isLockFile;

    DependencyFileEnum(String path, String ecosystem, Boolean isLockFile) {
        this.path = path;
        this.ecosystem = ecosystem;
        this.isLockFile = isLockFile;
    }
}
