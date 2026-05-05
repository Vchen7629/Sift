package app.service.githubRepo;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;

import java.io.IOException;
import java.net.URISyntaxException;
import java.nio.file.Files;
import java.nio.file.Path;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.kohsuke.github.GitHub;
import org.mockito.Mockito;

public class IssueServiceUT {
    private IssueService issueService;

    @BeforeEach
    void setup() {
        issueService = new IssueService(Mockito.mock(GitHub.class));
    }
    
    @Test
    void cleanIssueBody_RemovesMarkdownText() throws IOException, URISyntaxException {
        String unProcessedText = Files.readString(
            Path.of(getClass().getClassLoader().getResource("service/clean_issue_body_raw.txt").toURI())
        );
        String processedText = Files.readString(
            Path.of(getClass().getClassLoader().getResource("service/clean_issue_body_processed.txt").toURI())
        );
        String result = issueService.cleanIssueBody(unProcessedText);


        assertEquals(result, processedText);
        assertFalse(result.contains("```"));
    }
}
