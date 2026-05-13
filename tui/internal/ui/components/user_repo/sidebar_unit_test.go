//go:build unit

package user_repo

import (
	"testing"
	"tui/internal/api"
	"tui/internal/types"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
)

func testSidebarCtx() *context.App {
	return &context.App{
		SelectedTheme: styles.Warm,
	}
}

func sidebarWithDeps(n int) *Sidebar {
	deps := make([]types.Dependency, n)
	for i := range deps {
		deps[i] = types.Dependency{Name: "dep", Version: "1.0", Status: "healthy"}
	}
	m := NewSidebar(testSidebarCtx())
	m.FocusedIndexedRepo = &types.IndexedRepo{Dependencies: deps, TotalDependencies: n}
	return m
}

func TestSidebarUpdate_FocusGuard(t *testing.T) {
	tt := []struct {
		name             string
		isSidebarFocused bool
		focusedRepo      *types.IndexedRepo
	}{
		{"blocked when sidebar not focused", false, &types.IndexedRepo{Dependencies: make([]types.Dependency, 3)}},
		{"blocked when no indexed repo", true, nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := NewSidebar(testSidebarCtx())
			m.FocusedIndexedRepo = tc.focusedRepo
			m.FocusedIdx = 1

			m.Update(tea.KeyPressMsg{Code: tea.KeyDown}, tc.isSidebarFocused)

			assert.Equal(t, 1, m.FocusedIdx)
		})
	}
}

func TestSidebarUpdate_Navigation(t *testing.T) {
	tt := []struct {
		name       string
		key        rune
		initialIdx int
		wantIdx    int
	}{
		{"up decrements FocusedIdx", tea.KeyUp, 2, 1},
		{"down increments FocusedIdx", tea.KeyDown, 0, 1},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := sidebarWithDeps(3)
			m.FocusedIdx = tc.initialIdx

			m.Update(tea.KeyPressMsg{Code: tc.key}, true)

			assert.Equal(t, tc.wantIdx, m.FocusedIdx)
		})
	}
}

func TestSidebarUpdate_SearchQueryMsg(t *testing.T) {
	repos := []api.RepoApiRes{{Name: "sift"}, {Name: "other"}}
	indexed := map[string]*types.IndexedRepo{"sift": {Name: "sift"}}

	tt := []struct {
		name            string
		filtered        []api.RepoApiRes
		wantFocusedRepo string
		wantNil         bool
		wantNoResults   bool
	}{
		{"sets first filtered repo on sidebar", repos[:1], "sift", false, false},
		{"empty results clears sidebar", []api.RepoApiRes{}, "", true, true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := NewSidebar(testSidebarCtx())
			m.Update(searchQueryMsg{filteredGHRepos: tc.filtered, filteredIndexedRepos: indexed}, false)

			if tc.wantNil {
				assert.Nil(t, m.FocusedGHRepo)
				assert.Nil(t, m.FocusedIndexedRepo)
			} else {
				assert.Equal(t, tc.wantFocusedRepo, m.FocusedGHRepo.Name)
			}
			assert.Equal(t, tc.wantNoResults, m.noSearchResults)
		})
	}
}
