package app.repository;

import java.util.Map;

import org.opensearch.client.opensearch.OpenSearchClient;
import org.springframework.stereotype.Repository;

@Repository
public class OpenSearchRepository {
    private final OpenSearchClient openSearchClient;

    public OpenSearchRepository(OpenSearchClient openSearchClient) {
        this.openSearchClient = openSearchClient;
    }

    public void indexGithubIssue(Map<String, float[]> issueUrlTexts) {
        openSearchClient.bulk()
    }
}
