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

public class RustDependencyParserUT {

    private static final Path TEST_DATA_DIR = Path.of("src/test/java/app/parser_test_data");

    private RustDependencyParser parser;

    @BeforeEach
    void setup() {
        parser = new RustDependencyParser();
    }

    @Test
    void parseCargoToml_withLockFile_returnsDirectDepsWithExactVersions() throws IOException {
        String nonLockFileContent = Files.readString(TEST_DATA_DIR.resolve("cargo.toml"));
        String lockFileContent = Files.readString(TEST_DATA_DIR.resolve("cargo.lock"));

        List<Dependency> deps = parser.parse(DependencyFileEnum.CARGO_TOML, nonLockFileContent, lockFileContent);

        assertThat(deps).containsExactlyInAnyOrder(
            new Dependency("serde_json", "1.0.145"),
            new Dependency("serde", "1.0.228"),
            new Dependency("log", "0.4.28"),
            new Dependency("tauri", "2.9.1"),
            new Dependency("tauri-plugin-log", "2.7.0"),
            new Dependency("tauri-plugin-shell", "2.3.1"),
            new Dependency("tauri-plugin-http", "2.5.2"),
            new Dependency("tokio", "1.48.0"),
            new Dependency("tauri-plugin-store", "2.4.1")
        );
    }

    @Test
    void parseCargoToml_withoutLockFile_returnsDirectDepsWithNullVersions() throws IOException {
        String nonLockFileContent = Files.readString(TEST_DATA_DIR.resolve("cargo.toml"));

        List<Dependency> deps = parser.parse(DependencyFileEnum.CARGO_TOML, nonLockFileContent, null);

        assertThat(deps).extracting(Dependency::name)
            .containsExactlyInAnyOrder(
                "serde_json", "serde", "log", "tauri",
                "tauri-plugin-log", "tauri-plugin-shell", "tauri-plugin-http",
                "tokio", "tauri-plugin-store"
            );

        assertThat(deps).extracting(Dependency::version).containsOnlyNulls();
    }
}
