package types

type Repository struct {
	GithubId, TotalDependencies int
	Name, Description, Status   string 
	LastUpdated, LastIndexed    string
	Dependencies                []DependencyStatus
}

type DependencyStatus struct {
	Name, Version, Status string
}