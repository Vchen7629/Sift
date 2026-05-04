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
- [ ] Swap sentence transformer embedding model with one thats purpose built for RAG
- [ ] Implement search result reranker with cross encoding, add a threshold so it drops really unrelated issue results
- [ ] Issue body text chunking so each chunk of the issue has its own embedding, solves token limit of
      embedding model. 
      - Need to figure out what my chunk target is, could be dumb like chunking on paragraphs (easy) or smart where it takes into account
        context (harder)
- [ ] Valkey caching
- [x] Add check when indexing a new user repo's dependencies to see if the dependency name + version is already indexed so they dont try to refetch issues/changelog again for an already indexed dependency
- [ ] filter out markdown in text ('''Markdown) issue body before inserting to database
- [ ] investigate if search pipeline is used properly in the search query

### Edge cases
- Detect if the dependency repo is deprecated/archived and save it in the database as a "status field"
- Issues may not be marked with correct label or any labels at all from what i observed
- Changelog may not be found for the dependency repository
- Github api rate limits

### Tests/Metrics
- [x] create an evaluation dataset where i sample issues (200+) from a library and maybe use a LLM to generate 
symptom queries per issue (3 - 5 ), manually label the dataset, and benchmark:
      - MRR (Mean Reciprocal Rank)
      - Recall@10
      - NDCG@10
- [ ] Unit/Integration tests for code
- [ ] Prometheus metrics

### Deployment
- [ ] Kubernetes deployment
- [ ] Finish Readme