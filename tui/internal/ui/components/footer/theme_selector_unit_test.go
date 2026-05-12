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

func testThemeCtx(selectorOpen bool) *context.App {
	return &context.App{
		SelectedTheme:     styles.Mono,
		ThemeSelectorOpen: selectorOpen,
	}
}

func TestThemeSelector_Update(t *testing.T) {
	tt := []struct {
		name           string
		key            rune
		isSelectorOpen bool
		expected       styles.Theme
	}{
		{"pressing right updates the selected theme", tea.KeyRight, true, styles.Warm},
		{"pressing left updates the selected theme", tea.KeyLeft, true, styles.Cold},
		{"pressing some other key doesnt update theme", tea.KeyEnter, true, styles.Mono},
		{"pressing right doesnt change theme when selector is not open", tea.KeyRight, false, styles.Mono},
		{"pressing left doesnt change theme when selector is not open", tea.KeyLeft, false, styles.Mono},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testThemeCtx(tc.isSelectorOpen)
			m := footer.NewThemeSelector(ctx)
			m.Update(tea.KeyPressMsg{Code: tc.key})

			assert.Equal(t, tc.expected, ctx.SelectedTheme)
		})
	}
}
