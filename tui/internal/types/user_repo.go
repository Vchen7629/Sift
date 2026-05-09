package types

type Repository struct {
	GithubId, TotalDependencies int
	Name, Description, Status   string 
	LastUpdated, LastIndexed    string
	Dependencies                []DependencyStatus
}

type DependencyStatus struct {
	Id					  int
	Name, Version, Status string
}