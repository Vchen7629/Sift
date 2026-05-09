package types

type IndexedRepo struct { 
	Id, TotalDependencies int
	Name, LastIndexed     string
}