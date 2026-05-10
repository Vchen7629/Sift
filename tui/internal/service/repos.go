package service

import "tui/internal/types"

func FindIndexedRepo(name string, indexedRepos []types.IndexedRepo) *types.IndexedRepo {
	for i := range indexedRepos {
		if indexedRepos[i].Name == name {
			return &indexedRepos[i]
		}
	}
	return nil
}
