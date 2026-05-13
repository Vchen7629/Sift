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

func searchbarCtx() *context.App {
	return &context.App{SelectedTheme: styles.Warm, MainWidth: 80}
}

func newUserRepoTestSearchBar(repos []api.RepoApiRes, indexedMap map[string]*types.IndexedRepo) *SearchBarModel {
	sb := NewSearchBar(searchbarCtx(), "")
	sb.OriginalGHRepoList = repos
	sb.OriginalIndexedRepoList = indexedMap

	return sb
}

func TestSearchbar_SidebarFocusedGuard(t *testing.T) {
	tt := []struct {
		name string
		msg  tea.Msg
	}{
		{"blocks / toggle", tea.KeyPressMsg{Code: '/'}},
		{"blocks esc", tea.KeyPressMsg{Code: tea.KeyEscape}},
		{"blocks letter keys", tea.KeyPressMsg{Code: 'a'}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sb := newUserRepoTestSearchBar([]api.RepoApiRes{}, map[string]*types.IndexedRepo{})
			cmd := sb.Update(tc.msg, true)

			assert.Nil(t, cmd)
			assert.False(t, sb.IsFocused)
		})
	}
}

func TestSearchBar_SlashTogglesFocus(t *testing.T) {
	tt := []struct {
		name          string
		toggleCount   int
		expectedFocus bool
		wantCmd       bool
	}{
		{"1 slash input focuses searchbar", 1, true, true},
		{"2 slash input unfocuses searchbar", 2, false, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sb := newUserRepoTestSearchBar([]api.RepoApiRes{}, map[string]*types.IndexedRepo{})

			var cmd tea.Cmd
			for range tc.toggleCount {
				cmd = sb.Update(tea.KeyPressMsg{Code: '/'}, false)
			}

			assert.Equal(t, tc.expectedFocus, sb.IsFocused)
			if tc.wantCmd {
				assert.NotNil(t, cmd)
			} else {
				assert.Nil(t, cmd)
			}
		})
	}
}

func TestSearchBar_EscResetsInput(t *testing.T) {
	tt := []struct {
		name      string
		focused   bool
		wantReset bool
	}{
		{"esc when focused resets input", true, true},
		{"esc when not focused does nothing", false, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sb := newUserRepoTestSearchBar([]api.RepoApiRes{}, map[string]*types.IndexedRepo{})
			if tc.focused {
				sb.ToggleFocus()
			}
			sb.Update(tea.KeyPressMsg{Code: tea.KeyEscape}, false)

			if tc.wantReset {
				assert.Equal(t, "", sb.TextInput.Value())
			}
		})
	}
}

func TestSearchBar_EscEmitsOriginalList(t *testing.T) {
	repos := []api.RepoApiRes{{Name: "sift"}, {Name: "other"}}
	sb := newUserRepoTestSearchBar(repos, nil)
	sb.ToggleFocus()

	cmd := sb.Update(tea.KeyPressMsg{Code: tea.KeyEscape}, false)
	msg, ok := cmd().(searchQueryMsg)

	assert.True(t, ok)
	assert.Equal(t, repos, msg.filteredGHRepos)
}

func TestSearchBar_FocusGuard(t *testing.T) {
	tt := []struct {
		name    string
		focused bool
		wantCmd bool
	}{
		{"focused emits cmd", true, true},
		{"unfocused no cmd", false, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sb := newUserRepoTestSearchBar([]api.RepoApiRes{{Name: "sift"}}, nil)
			if tc.focused {
				sb.ToggleFocus()
			}

			cmd := sb.Update(tea.KeyPressMsg{Code: 's', Text: "s"}, false)

			if tc.wantCmd {
				assert.NotNil(t, cmd)
			} else {
				assert.Nil(t, cmd)
			}
		})
	}
}

func TestFilterNameMatch(t *testing.T) {
	tt := []struct {
		name          string
		query         string
		expectedOrder []string
	}{
		{"orders exact then startsWith then contains", "sift", []string{"sift", "sift-api", "my-sift-fork"}},
		{"no match returns empty", "zzz", []string{}},
		{"case insensitive", "SIFT", []string{"sift", "sift-api", "my-sift-fork"}},
		{"only includes matching name", "starbright", []string{"starbright"}},
	}

	repos := []api.RepoApiRes{{Name: "my-sift-fork"}, {Name: "sift"}, {Name: "sift-api"}, {Name: "starbright"}}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sb := newUserRepoTestSearchBar(repos, nil)
			result, _ := sb.filterNameMatch(tc.query)

			names := make([]string, len(result))
			for i, r := range result {
				names[i] = r.Name
			}

			assert.Equal(t, tc.expectedOrder, names)
		})
	}
}
