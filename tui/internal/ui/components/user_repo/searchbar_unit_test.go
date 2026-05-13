//go:build unit

package user_repo_test

import (
	"testing"
	"tui/internal/ui/components/user_repo"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
)

func searchbarCtx() *context.App {
	return &context.App{SelectedTheme: styles.Warm, MainWidth: 80}
}

func newUserRepoSearchBar() *user_repo.SearchBarModel {
	return user_repo.NewSearchBar(searchbarCtx(), "")
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
			sb := newUserRepoSearchBar()
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
			sb := newUserRepoSearchBar()

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
			sb := newUserRepoSearchBar()
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
