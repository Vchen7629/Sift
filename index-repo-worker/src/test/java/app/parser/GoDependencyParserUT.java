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


public class GoDependencyParserUT {

    private static final Path TEST_DATA_DIR = Path.of("src/test/java/app/parser_test_data");

    private GoDependencyParser parser;

    @BeforeEach
    void setup() {
        parser = new GoDependencyParser();
    }

    @Test
    void parseGoMod_returnsOnlyDirectDependencies() throws IOException {
        String content = Files.readString(TEST_DATA_DIR.resolve("go.mod"));

        List<Dependency> deps = parser.parse(DependencyFileEnum.GO_MOD, content, null);

        assertThat(deps).hasSize(4);
        assertThat(deps).containsExactlyInAnyOrder(
            new Dependency("github.com/jackc/pgconn", "v1.14.3"),
	        new Dependency("github.com/jackc/pgx/v4", "v4.18.3"),
	        new Dependency("github.com/joho/godotenv", "v1.5.1"),
	        new Dependency("github.com/kelseyhightower/envconfig", "v1.4.0")
        );
        assertThat(deps).noneMatch(d -> d.name().equals("dario.cat/mergo"));
    }
}
