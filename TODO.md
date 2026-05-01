## TODO Later

### Features
- [ ] Golang TUI client with Bubbletea library (poggers)
      - Could be cool to display text that show that the search results are strong match vs weak match

### Optimizations
- [ ] Issue label detection, detect if search query contains an issue label like "good first issue" and apply an optional filter to only fetch
      issues that are labeled as that
      - Edge case: user only types in query with issue label, reranker will probably rerank results to show the ones with the most comments/most recent
- [ ] Swap sentence transformer embedding model with one thats purpose built for RAG
- [ ] Implement search result reranker with cross encoding, add a threshold so it drops really unrelated issue results
- [ ] Issue body text chunking so each chunk of the issue has its own embedding, solves token limit of
      embedding model. 
      - Need to figure out what my chunk target is, could be dumb like chunking on paragraphs (easy) or smart where it takes into account
        context (harder)
- [ ] Valkey caching
- [ ] Idk if this is a good fit for this project but maybe levenshtein distance for fuzzy search

### Tests/Metrics
- [ ] Unit/Integration tests for code
- [ ] Prometheus metrics

### Deployment
- [ ] Kubernetes deployment
- [ ] Finish Readme