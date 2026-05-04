## TODO Later

### Features
- [ ] Support for Python dependency parsing
- [ ] Support for JS/TS dependency parsing
- [ ] Support for Java dependency parsing (maybe depending on how cancer maven is to parse)
- [ ] Golang TUI client with Bubbletea library (poggers)
      - show depency status like deprecated or not found
      - Could be cool to display text that show that the search results are strong match vs weak match
- [x] add a last updated field to user Repo

### Optimizations
- [x] Swap sentence transformer embedding model with one thats purpose built for RAG
- [ ] Implement search result reranker with cross encoding, add a threshold so it drops really unrelated issue results
- [ ] Issue body text chunking so each chunk of the issue has its own embedding, solves token limit of
      embedding model. 
      - Need to figure out what my chunk target is, could be dumb like chunking on paragraphs (easy) or smart where it takes into account
        context (harder)
- [ ] Valkey caching
- [x] Add check when indexing a new user repo's dependencies to see if the dependency name + version is already indexed so they dont try to refetch issues/changelog again for an already indexed dependency
- [x] filter out markdown in text ('''Markdown) issue body before inserting to database
- [x] investigate if search pipeline is used properly in the search query
- [ ] look into parallelizing Embedding cpu processing since it takes 64665ms (64.665s) to create text embeddings for 3128 issues
- [ ] Look into parallelizing Fetch issue query page calls into seperate virtual threads to speed up I/O bottleneck even more
- [ ] could look into rate limiting the public api

### Edge cases
- Detect if the dependency repo is deprecated/archived and save it in the database as a "status field"
- Issues may not be marked with correct label or any labels at all from what i observed
- Changelog may not be found for the dependency repository
- Github api rate limits

### Tests/Metrics
- [ ] Unit/Integration tests for code
- [ ] Performance testing with K6s or similar tool
- [ ] e2e tests (maybe can drop this)

### Metrics
- [x] create an evaluation dataset where i sample issues (200+) from a library and maybe use a LLM to generate 
symptom queries per issue (3 - 5 ), manually label the dataset, and benchmark:
      - MRR (Mean Reciprocal Rank)
      - Recall@10
      - NDCG@10
- [ ] Prometheus metrics
- [ ] Replace all timing logs with Micrometer spans
- [ ] Add otel bridge in micrometer to convert tracing api calls into OTEL SDK calls
- [ ] Add a exporter library like opentelemetry-exporter-otlp to send the spans to the storage
- [ ] Add grafana tempo to store spans
- [ ] Grafana metrics and span visualization

### Deployment
- [ ] look into spring security to lock down sensitive routes
- [ ] Kubernetes deployment
- [ ] Finish Readme