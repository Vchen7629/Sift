package app.parser;

import static org.assertj.core.api.Assertions.assertThat;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import app.model.DependencyFileEnum;
import app.component.parser.DependencyParserStrategy.Dependency;
import app.component.parser.GoDependency;


public class GoDependencyParserUT {

    private static final Path TEST_DATA_DIR = Path.of("src/test/java/app/parser_test_data");

    private GoDependency parser;

    @BeforeEach
    void setup() {
        parser = new GoDependency();
    }

    @Test
    void parseGoMod_returnsOnlyDirectDependencies() throws IOException {
        String content = Files.readString(TEST_DATA_DIR.resolve("go.mod"));

        List<Dependency> deps = parser.parse(DependencyFileEnum.GO_MOD, content, null);

        assertThat(deps).hasSize(4);
        assertThat(deps).containsExactlyInAnyOrder(
            new Dependency("jackc/pgconn", "v1.14.3", "jackc/pgconn"),
	        new Dependency("pgx/v4", "v4.18.3", "jackc/pgx/v4"),
	        new Dependency("joho/godotenv", "v1.5.1", "joho/godotenv"),
	        new Dependency("kelseyhightower/envconfig", "v1.4.0", "kelseyhightower/envconfig")
        );
        assertThat(deps).noneMatch(d -> d.name().equals("dario.cat/mergo"));
    }
}
