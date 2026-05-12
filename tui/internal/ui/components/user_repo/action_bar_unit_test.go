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

func actionBarCtx() *context.App {
	return &context.App{
		SelectedTheme: styles.Warm,
	}
}

func TestActionBar_Update(t *testing.T) {
	tt := []struct {
		name        string
		key         rune
		expectedMsg tea.Msg
	}{
		{"pressing s returns ToggleFocusMsg", 's', user_repo.ToggleFocusMsg{}},
		{"pressing r returns IndexRepoRequestMsg", 'r', user_repo.IndexRepoRequestMsg{}},
		{"pressing other key (enter) returns nil", tea.KeyEnter, nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := actionBarCtx()
			m := user_repo.NewActionBar(ctx)
			cmd := m.Update(tea.KeyPressMsg{Code: tc.key})

			if tc.expectedMsg == nil {
				assert.Nil(t, cmd)
			} else {
				assert.NotNil(t, cmd)
				assert.Equal(t, tc.expectedMsg, cmd())
			}
		})
	}
}
