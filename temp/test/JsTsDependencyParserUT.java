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

public class JsTsDependencyParserUT {

    private static final Path TEST_DATA_DIR = Path.of("src/test/java/app/parser_test_data");

    private JsTsDependencyParser parser;

    @BeforeEach
    void setup() {
        parser = new JsTsDependencyParser();
    }

    @Test
    void parsePackageJson_withLockFile_returnsDirectDepsWithExactVersions() throws IOException {
        String nonLockFileContent = Files.readString(TEST_DATA_DIR.resolve("package.json"));
        String lockFileContent = Files.readString(TEST_DATA_DIR.resolve("package-lock.json"));

        List<Dependency> deps = parser.parse(DependencyFileEnum.PACKAGE_JSON, nonLockFileContent, lockFileContent);

        assertThat(deps).containsExactlyInAnyOrder(
            new Dependency("@react-native-picker/picker", "2.11.1"),
            new Dependency("@react-navigation/native", "7.2.2"),
            new Dependency("@react-navigation/native-stack", "7.14.10"),
            new Dependency("@tanstack/react-query", "5.97.0"),
            new Dependency("axios", "1.15.0"),
            new Dependency("expo", "54.0.34"),
            new Dependency("expo-secure-store", "15.0.8"),
            new Dependency("expo-sensors", "15.0.8"),
            new Dependency("expo-status-bar", "3.0.9"),
            new Dependency("lucide-react-native", "1.8.0"),
            new Dependency("nativewind", "4.2.3"),
            new Dependency("react", "19.1.0"),
            new Dependency("react-native", "0.81.5"),
            new Dependency("react-native-calendars", "1.1314.0"),
            new Dependency("react-native-safe-area-context", "5.6.2"),
            new Dependency("react-native-screens", "4.16.0"),
            new Dependency("tailwindcss", "3.4.19")
        );
    }

    @Test
    void parsePackageJson_withoutLockFile_returnsDirectDepsWithNullVersions() throws IOException {
        String nonLockFileContent = Files.readString(TEST_DATA_DIR.resolve("package.json"));

        List<Dependency> deps = parser.parse(DependencyFileEnum.PACKAGE_JSON, nonLockFileContent, null);

        assertThat(deps).extracting(Dependency::name)
            .containsExactlyInAnyOrder(
                "@react-native-picker/picker", "@react-navigation/native", "@react-navigation/native-stack",
                "@tanstack/react-query", "axios", "expo", "expo-secure-store", "expo-sensors",
                "expo-status-bar", "lucide-react-native", "nativewind", "react", "react-native",
                "react-native-calendars", "react-native-safe-area-context", "react-native-screens", "tailwindcss"
            );

        assertThat(deps).extracting(Dependency::version).containsOnlyNulls();
    }
}
