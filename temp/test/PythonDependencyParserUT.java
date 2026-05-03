package app.parser;

import static org.assertj.core.api.Assertions.assertThat;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import app.model.DependencyFileEnum;
import app.parser.DependencyParserStrategy.Dependency;

public class PythonDependencyParserUT {
    
    private static final Path TEST_DATA_DIR = Path.of("src/test/java/app/parser_test_data");

    private PythonDependencyParser parser;

    @BeforeEach
    void setup() {
        parser = new PythonDependencyParser();
    }

    @Test
    void parsePyProjectTOML_withLockFile_returnsDirectDepsWithExactVersions() throws IOException {
        String nonLockFileContent = Files.readString(TEST_DATA_DIR.resolve("pyproject.toml"));
        String lockFileContent = Files.readString(TEST_DATA_DIR.resolve("uv.lock"));

        List<Dependency> deps = parser.parse(DependencyFileEnum.PYPROJECT_TOML, nonLockFileContent, lockFileContent);

        assertThat(deps).containsExactlyInAnyOrder(
            new Dependency("asyncpg", "0.31.0"),
            new Dependency("bcrypt", "5.0.0"),
            new Dependency("fastapi", "0.135.3"),
            new Dependency("pydantic-settings", "2.13.1"),
            new Dependency("scikit-learn", "1.8.0"),
            new Dependency("sqlalchemy", "2.0.49"),
            new Dependency("structlog", "25.5.0"),
            new Dependency("xgboost", "3.2.0")
        );
    }

    @Test
    void parsePyProjectTOML_withoutLockFile_returnsDirectDepsWithNullVersions() throws IOException {
        String nonLockFileContent = Files.readString(TEST_DATA_DIR.resolve("pyproject.toml"));

        List<Dependency> deps = parser.parse(DependencyFileEnum.PYPROJECT_TOML, nonLockFileContent, null);

        assertThat(deps).extracting(Dependency::name)
            .containsExactlyInAnyOrder(
                "asyncpg", "bcrypt", "fastapi", "pydantic-settings",
                 "scikit-learn", "sqlalchemy", "structlog", "xgboost"
            );
        
        assertThat(deps).extracting(Dependency::version).containsOnlyNulls();
    }

    @Test
    void parseRequirementsTXT_returnsDepsWithPinnedVersions() throws IOException {
        String content = Files.readString(TEST_DATA_DIR.resolve("requirements.txt"));

        List<Dependency> deps = parser.parse(DependencyFileEnum.REQUIREMENTS_TXT, content, null);

        assertThat(deps).contains(
            new Dependency("asyncpg", "0.31.0"),
            new Dependency("fastapi", "0.135.3"),
            new Dependency("sqlalchemy", "2.0.49")
        );

        assertThat(deps).noneMatch(d -> d.name().startsWith("#"));
        assertThat(deps).noneMatch(d -> d.name().isBlank());
    }
}
