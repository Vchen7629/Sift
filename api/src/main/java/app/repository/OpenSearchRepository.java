package app.repository;

import org.opensearch.client.opensearch.OpenSearchClient;
import org.springframework.stereotype.Repository;

@Repository
public class OpenSearchRepository {
    private final OpenSearchClient openSearchClient;

    public OpenSearchRepository(OpenSearchClient openSearchClient) {
        this.openSearchClient = openSearchClient;
    }

    public void indexGithubIssue() {

    }
}
