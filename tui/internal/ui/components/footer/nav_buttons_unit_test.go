//go:build unit

package footer_test

import (
	"testing"
	"tui/internal/ui/components/footer"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
)

func testNavButtonsCtx() *context.App {
	return &context.App{
		SelectedTheme:     styles.Mono,
		CurrentPage:       context.UserReposPage,
		ThemeSelectorOpen: true,
	}
}

func TestNavButtons_Update(t *testing.T) {
	tt := []struct {
		name                string
		key                 rune
		expectedPage        context.Page
		isThemeSelectorOpen bool
	}{
		{"pressing 1 switches current page to user repo", tea.KeyKp1, context.UserReposPage, false},
		{"pressing 2 switches current page to query", tea.KeyKp2, context.QueryPage, false},
		{"pressing some other key doesnt current page", tea.KeyKp3, context.UserReposPage, true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testNavButtonsCtx()
			m := footer.NewNavButtons(ctx)
			m.Update(tea.KeyPressMsg{Code: tc.key})

			assert.Equal(t, ctx.ThemeSelectorOpen, tc.isThemeSelectorOpen)
			assert.Equal(t, tc.expectedPage, ctx.CurrentPage)
		})
	}
}
