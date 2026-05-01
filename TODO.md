## TODO Later
- [ ] Issue label detection, detect if search query contains an issue label like "good first issue" and apply an optional filter to only fetch
      issues that are labeled as that
      - Edge case: user only types in query with issue label, reranker will probably rerank results to show the ones with the most comments/most recent
- [ ] Idk if this is a good fit for this project but maybe levenshtein distance for fuzzy search
- [ ] Swap sentence transformer embedding model with one thats purpose built for RAG
- [ ] Implement search result reranker with cross encoding
- [ ] Issue body text chunking so each chunk of the issue has its own embedding, solves token limit of
      embedding model. 
      - Need to figure out what my chunk target is, could be dumb like chunking on paragraphs (easy) or smart where it takes into account
        context (harder)
- [ ] Golang TUI client with Bubbletea library (poggers)
- [ ] Valkey caching for optimizations
- [ ] Unit/Integration tests for code
- [ ] Prometheus metrics
- [ ] Kubernetes deployment
- [ ] Finish Readme