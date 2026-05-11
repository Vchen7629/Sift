package types

type IndexedRepo struct {
	TotalDependencies int
	Name, LastIndexed string
	Dependencies      []Dependency
}

type Dependency struct {
	Id      int
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"`
}
