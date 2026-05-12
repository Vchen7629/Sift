//go:build unit

package user_repo_test

import (
	"testing"
	"tui/internal/api"
	"tui/internal/ui/components/user_repo"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
)

func testListCtx() *context.App {
	return &context.App{Username: "testuser", SelectedTheme: styles.Warm}
}

func listWithRepos(repos []api.RepoApiRes) *user_repo.ListModel {
	m := user_repo.NewUserRepoList(testListCtx())
	m.GHRepos = repos

	return m
}

var fakeRepos = []api.RepoApiRes{{Name: "Coco"}, {Name: "Cottom"}, {Name: "Maple"}}

func TestListUpdate_SidebarFocusedGuard(t *testing.T) {
	tt := []struct {
		name             string
		isSidebarFocused bool
		key              rune
		wantIdx          int
	}{
		{"down blocked when sidebar focused", true, tea.KeyDown, 0},
		{"up blocked when sidebar focused", true, tea.KeyUp, 0},
		{"down works when sidebar not focused", false, tea.KeyDown, 1},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := listWithRepos(fakeRepos)
			m.Update(tea.KeyPressMsg{Code: tc.key}, tc.isSidebarFocused)

			assert.Equal(t, tc.wantIdx, m.FocusedIdx)
		})
	}
}

func TestListUpdate_IndexRepoRequest(t *testing.T) {
	tt := []struct {
		name    string
		repos   []api.RepoApiRes
		wantCmd bool
	}{
		{"no-op when repos empty", nil, false},
		{"returns cmd when repo not indexed", fakeRepos, true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := listWithRepos(tc.repos)
			cmd := m.Update(user_repo.IndexRepoRequestMsg{}, false)

			if tc.wantCmd {
				assert.NotNil(t, cmd)
			} else {
				assert.Nil(t, cmd)
			}
		})
	}
}
