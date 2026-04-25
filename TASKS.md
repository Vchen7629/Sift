# Sift - RAG Pipeline Tasks

## Goal
Semantic search over GitHub issues for a repo using embeddings + OpenSearch.

## Architecture Flow
```
fetchRepoIssues() → generateEmbeddings() → indexGithubIssues()
GithubApiService   TextEmbeddingService    OpenSearchRepository
```

---

## Completed
- [x] Spring Boot project setup
- [x] GitHub API integration (`GithubApiService`)
- [x] OpenSearch client config (`OpenSearchConfig`)
- [x] DJL + sentence-transformers model config (`ModelConfig`)
- [x] `TextEmbeddingService` skeleton with batching + `partition()` helper
- [x] `OpenSearchRepository` skeleton with index mapping (knn_vector, 384 dims)
- [x] Exception handlers for model loading errors

---

## In Progress

### Model layer
- [ ] Create `app/model/GithubIssueDocument.java` record (url, title, text)
- [ ] Create `app/model/GithubIssueEmbedding.java` record (url, title, embedding) — OR decide on nullable embedding approach
- [ ] Decide: two records vs one record with nullable embedding

### GithubApiService
- [ ] Update `fetchRepoIssues()` to return `CompletableFuture<List<GithubIssueDocument>>` instead of `Void`
- [ ] Populate url (`issue.getHtmlUrl()`), title, and body text per issue
- [ ] Remove 20 issue page limit (currently hardcoded for testing)

### TextEmbeddingService
- [ ] Update `generateEmbeddings()` to accept `List<GithubIssueDocument>` and return `List<GithubIssueEmbedding>`
- [ ] Thread url + title through alongside embeddings so they aren't lost after batching

### OpenSearchRepository
- [ ] Implement `createIndexIfNotExists()` with knn_vector mapping
- [ ] Call `createIndexIfNotExists()` via `@PostConstruct`
- [ ] Implement `indexGithubIssues(List<GithubIssueEmbedding>)` using bulk API

### Controller / Orchestration
- [ ] Wire up the async chain in the route:
  ```java
  fetchRepoIssues(repoName)
      .thenApply(issues -> embeddingService.generateEmbeddings(issues))
      .thenAccept(embeddings -> openSearchRepository.indexGithubIssues(embeddings));
  ```

---

## TODO Later
- [ ] Search endpoint — accept query string, embed it, run knn query against OpenSearch, return list of `{url, title}`
- [ ] Remove 20 issue hardcoded limit
- [ ] Pagination for large repos
- [ ] Hybrid search (BM25 + embeddings) if semantic-only isn't accurate enough
