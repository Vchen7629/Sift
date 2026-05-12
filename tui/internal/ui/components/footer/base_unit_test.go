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

func testBaseCtx(isThemeSelectorOpen bool) *context.App {
	return &context.App{
		SelectedTheme:     styles.Mono,
		ThemeSelectorOpen: isThemeSelectorOpen,
	}
}

func TestBase_Update(t *testing.T) {
	tt := []struct {
		name              string
		key               rune
		originalOpenState bool
		expectedOpenState bool
	}{
		{"pressing 3 opens theme selector when its closed", tea.KeyKp3, false, true},
		{"pressing 3 closes theme selector when its open", tea.KeyKp3, true, false},
		{"pressing enter does not affect theme selector", tea.KeyEnter, false, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testBaseCtx(tc.originalOpenState)
			m := footer.NewFooterBaseModel(ctx)
			m.Update(tea.KeyPressMsg{Code: tc.key})

			assert.Equal(t, tc.expectedOpenState, ctx.ThemeSelectorOpen)
		})
	}
}
