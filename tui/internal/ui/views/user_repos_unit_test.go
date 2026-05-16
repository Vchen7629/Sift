//go:build unit

package views

import (
	"testing"
	"tui/internal/api"
	"tui/internal/types"
	"tui/internal/ui/common"
	"tui/internal/ui/context"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
)

func newUserReposCtx() *context.App {
	return &context.App{}
}

func newUserRepoModel(isSeaching, isSidebarFocused bool) *UserRepoModel {
	m := NewUserRepo(newUserReposCtx())
	m.isSidebarFocused = isSidebarFocused
	m.SearchBar.IsFocused = isSeaching

	return m
}

func TestUserRepo_Update_ToggleFocusMsg(t *testing.T) {
	tt := []struct {
		name                 string
		isSearching          bool
		numToggles           int
		originalSidebarState bool
		expectedSidebarState bool
	}{
		{"it doesnt toggle sidebarfocused when isSearching is true", true, 1, false, false},
		{"Toggles sidebarfocused to true", false, 1, false, true},
		{"Toggles sidebarfocused to false", false, 1, true, false},
		{"2 toggles toggles it back to original state", false, 2, true, true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := newUserRepoModel(tc.isSearching, tc.originalSidebarState)

			var model tea.Model
			var cmd tea.Cmd
			for range tc.numToggles {
				model, cmd = m.Update(common.ToggleFocusMsg{})
			}

			assert.Equal(t, tc.expectedSidebarState, m.isSidebarFocused)
			assert.Equal(t, m, model)
			assert.Nil(t, cmd)
		})
	}
}

func TestUserRepo_Update_GithubRepoFetchedMsg(t *testing.T) {
	tt := []struct {
		name                   string
		ghRepoList             []api.RepoApiRes
		indexedRepos           []types.IndexedRepo
		wantGHRepoCount        int
		wantFocusedRepo        *api.RepoApiRes
		wantFocusedIndexedRepo *types.IndexedRepo
	}{
		{
			name:                   "empty repo list sets no focused repo",
			ghRepoList:             []api.RepoApiRes{},
			wantGHRepoCount:        0,
			wantFocusedRepo:        nil,
			wantFocusedIndexedRepo: nil,
		},
		{
			name:                   "focuses first repo, not indexed",
			ghRepoList:             []api.RepoApiRes{{Name: "Sift"}, {Name: "idk"}},
			wantGHRepoCount:        2,
			wantFocusedRepo:        &api.RepoApiRes{Name: "Sift"},
			wantFocusedIndexedRepo: nil,
		},
		{
			name:                   "focuses first repo, indexed",
			ghRepoList:             []api.RepoApiRes{{Name: "Sift"}, {Name: "idk"}},
			indexedRepos:           []types.IndexedRepo{{Name: "Sift"}, {Name: "idk"}},
			wantGHRepoCount:        2,
			wantFocusedRepo:        &api.RepoApiRes{Name: "Sift"},
			wantFocusedIndexedRepo: &types.IndexedRepo{Name: "Sift"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := newUserRepoModel(false, false)

			indexedMap := make(map[string]*types.IndexedRepo)
			for i := range tc.indexedRepos {
				indexedMap[tc.indexedRepos[i].Name] = &tc.indexedRepos[i]
			}
			m.RepoList.IndexedRepoMap = indexedMap

			m.Update(githubRepoFetchedMsg{repoList: tc.ghRepoList})

			assert.Equal(t, tc.wantGHRepoCount, m.ActionBar.GHRepoCount)
			assert.Equal(t, tc.wantFocusedRepo, m.Sidebar.FocusedGHRepo)
			assert.Equal(t, tc.wantFocusedIndexedRepo, m.Sidebar.FocusedIndexedRepo)
		})
	}
}
func TestUserRepoUpdate_FetchIndexedRepoMsg(t *testing.T) {
	indexedSift := types.IndexedRepo{Name: "Sift"}
	indexedBudget := types.IndexedRepo{Name: "Budget.me"}

	tt := []struct {
		name                   string
		ghRepos                []api.RepoApiRes
		indexedRepos           []types.IndexedRepo
		isReauthed             bool
		newSessToken           string
		expectedToken          string
		expectedIndexedMap     map[string]*types.IndexedRepo
		expectedFocusedIndexed *types.IndexedRepo
		expectedIndexedStatus  map[int]string
	}{
		{
			name:                   "empty indexed repos clears map, no focused repo",
			ghRepos:                []api.RepoApiRes{{Name: "Sift"}},
			indexedRepos:           []types.IndexedRepo{},
			isReauthed:             false,
			newSessToken:           "new-token",
			expectedToken:          "old-token",
			expectedIndexedMap:     map[string]*types.IndexedRepo{},
			expectedFocusedIndexed: nil,
			expectedIndexedStatus:  map[int]string{},
		},
		{
			name:          "focused repo set to first GH repo when it is indexed",
			ghRepos:       []api.RepoApiRes{{Name: "Sift"}, {Name: "Budget.me"}},
			indexedRepos:  []types.IndexedRepo{indexedSift, indexedBudget},
			isReauthed:    false,
			newSessToken:  "new-token",
			expectedToken: "old-token",
			expectedIndexedMap: map[string]*types.IndexedRepo{
				"Sift": &indexedSift, "Budget.me": &indexedBudget,
			},
			expectedFocusedIndexed: &indexedSift,
			expectedIndexedStatus:  map[int]string{0: "Indexed", 1: "Indexed"},
		},
		{
			name:                   "focused repo nil when first GH repo not indexed",
			ghRepos:                []api.RepoApiRes{{Name: "Sift"}},
			indexedRepos:           []types.IndexedRepo{indexedBudget},
			isReauthed:             false,
			newSessToken:           "new-token",
			expectedToken:          "old-token",
			expectedIndexedMap:     map[string]*types.IndexedRepo{"Budget.me": &indexedBudget},
			expectedFocusedIndexed: nil,
			expectedIndexedStatus:  map[int]string{},
		},
		{
			name:                   "no GH repos means no focused repo regardless of indexed repos",
			ghRepos:                []api.RepoApiRes{},
			indexedRepos:           []types.IndexedRepo{indexedSift},
			isReauthed:             false,
			newSessToken:           "new-token",
			expectedToken:          "old-token",
			expectedIndexedMap:     map[string]*types.IndexedRepo{"Sift": &indexedSift},
			expectedFocusedIndexed: nil,
			expectedIndexedStatus:  map[int]string{},
		},
		{
			name:                   "updates session token when reauthed",
			ghRepos:                []api.RepoApiRes{},
			indexedRepos:           []types.IndexedRepo{},
			isReauthed:             true,
			newSessToken:           "new-token",
			expectedToken:          "new-token",
			expectedIndexedMap:     map[string]*types.IndexedRepo{},
			expectedFocusedIndexed: nil,
			expectedIndexedStatus:  map[int]string{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := newUserRepoModel(false, false)
			m.ctx.SessionToken = "old-token"
			m.RepoList.GHRepos = tc.ghRepos

			m.Update(common.FetchIndexedRepoMsg{
				IndexedRepos:    tc.indexedRepos,
				NewSessionToken: tc.newSessToken,
				IsReauthed:      tc.isReauthed,
			})

			assert.Equal(t, tc.expectedToken, m.ctx.SessionToken)
			assert.Equal(t, tc.expectedIndexedMap, m.RepoList.IndexedRepoMap)
			assert.Equal(t, tc.expectedIndexedMap, m.SearchBar.OriginalIndexedRepoList)
			assert.False(t, m.ActionBar.IndexRepoApiDown)
			assert.Equal(t, tc.expectedFocusedIndexed, m.Sidebar.FocusedIndexedRepo)
			for idx, wantStatus := range tc.expectedIndexedStatus {
				status, ok := m.RepoList.IndexCoord.StatusFor(idx)
				assert.True(t, ok)
				assert.Equal(t, wantStatus, status)
			}
		})
	}
}
