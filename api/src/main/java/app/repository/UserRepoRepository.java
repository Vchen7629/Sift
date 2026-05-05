package app.repository;

import java.io.IOException;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.NoSuchElementException;
import java.util.Objects;

import org.opensearch.client.opensearch.OpenSearchClient;
import org.opensearch.client.opensearch._types.aggregations.Aggregate;
import org.opensearch.client.opensearch._types.aggregations.StringTermsBucket;
import org.opensearch.client.opensearch.core.DeleteByQueryResponse;
import org.opensearch.client.opensearch.core.SearchResponse;
import org.springframework.stereotype.Repository;
import org.springframework.validation.annotation.Validated;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import io.micrometer.observation.annotation.Observed;
import jakarta.validation.constraints.NotBlank;
import lombok.extern.slf4j.Slf4j;
import static net.logstash.logback.argument.StructuredArguments.kv;


@Repository
@Validated
@Slf4j
public class UserRepoRepository {
    private final OpenSearchClient openSearchClient;
    private final static String indexedRepoIndexName = "user-repo";

    public UserRepoRepository(OpenSearchClient openSearchClient) {
        this.openSearchClient = openSearchClient;
    }

    @Observed(name="userrepo.delete.repository")
    public void delete(@NotBlank String userId, @NotBlank String repoName) throws IOException {
        DeleteByQueryResponse deleteRes =  openSearchClient.deleteByQuery(d -> d
            .index(indexedRepoIndexName)
            .query(q -> q
                .bool(b -> b
                    .must(m -> m.term(t -> t.field("userId").value(v -> v.stringValue(userId))))
                    .must(m -> m.term(t -> t.field("repoName").value(v -> v.stringValue(repoName))))
                )
            )
        );

        log.debug("successfully deleted track Repo", kv("repoName", repoName));

        if (deleteRes.deleted() == 0) {
            log.warn("No repo found in db to delete", kv("repoName", repoName));

            throw new NoSuchElementException("No repo found to delete");
        }
    }

    @JsonIgnoreProperties(ignoreUnknown = true)
    private record UserRepoDepDoc(Map<String, String> dependencies) {}

    /**
     * fetches all dependencies across the user's indexed repo for search filtering on
     * @param userId used to decide which user to fetch for
     * @return a map containing the dependency name and version pairs
     * @throws IOException from openSearchClient.search
     */
    @Observed(name="userrepo.listalldependencies.repository")
    public Map<String, String> listAllDependencies(@NotBlank String userId) throws IOException {
        final int maxUniqueDependencies = 1000;

        SearchResponse<UserRepoDepDoc> searchRes = openSearchClient.search(r -> r
            .index(indexedRepoIndexName)
            .size(maxUniqueDependencies)
            .query(q -> q
                .term(t -> t.field("userId").value(v -> v.stringValue(userId)))
            ),
            UserRepoDepDoc.class
        );

        Map<String, String> repoDependencies = new HashMap<>();

        searchRes.hits().hits().stream()
            .map(hit -> hit.source())
            .filter(Objects::nonNull)
            .forEach(doc -> repoDependencies.putAll(doc.dependencies()));


        log.debug("fetched {} dependencies", repoDependencies.size(), kv("userId", userId));

        return repoDependencies;
    }

    /**
     * fetch all indexed repos belonging to the user
     * @param userId the user to fetch the indexed repos for
     * @return a list of repos indexed for the user
     * @throws IOException from openSearchClient.search
     */
    @Observed(name="userrepo.listall.repository")
    public List<String> listAll(@NotBlank String userId) throws IOException {
        final int maxUniqueRepos = 1000;

        SearchResponse<Void> searchRes = openSearchClient.search(r -> r
            .index(indexedRepoIndexName)
            .size(0) // need this so we dont return the actual document, use aggregations to return the strings instead
            .timeout("30s")
            .query(q -> q
                .term(t -> t.field("userId").value(v -> v.stringValue(userId)))
            )
            .aggregations("repoNames", a -> a
                .terms(t -> t.field("repoName").size(maxUniqueRepos))
            ),
            Void.class
        );

        Aggregate repoNameAgg = searchRes.aggregations().get("repoNames");
        if (repoNameAgg == null || !repoNameAgg.isSterms()) {
            log.debug("found no indexed repos in the db");
            return List.of(); // return empty list instead of throwing an exception
        }

        List<String> indexedRepos = repoNameAgg
            .sterms()
            .buckets().array()
            .stream()
            .map(StringTermsBucket::key)
            .toList();

        log.debug("Successfully fetched {} indexed repos", indexedRepos.size());

        return indexedRepos;
    }

    @Observed(name="userrepo.isrepoindexed.repository")
    public boolean isRepoIndexed(String repoName) throws IOException{
        SearchResponse<Void> searchRes = openSearchClient.search(r -> r
            .index(indexedRepoIndexName)
            .timeout("30s")
            .query(q -> q
                .term(t -> t.field("repoName").value(v -> v.stringValue(repoName))
                )
            ), 
            Void.class
        );

        var total = searchRes.hits().total();

        return total != null && total.value() > 0;
    }
}
