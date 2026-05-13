//go:build unit

package common_test

import (
	"testing"
	"tui/internal/ui/common"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
)

func newSearchBarCtx() *context.App {
	return &context.App{
		SelectedTheme: styles.Warm,
		MainWidth:     80,
	}
}

func newSearchBar(placeholder string) *common.SearchBar {
	return common.NewSearchBar(newSearchBarCtx(), placeholder)
}

func TestNewSearchBar_InitialState(t *testing.T) {
	sb := newSearchBar("Search...")
	assert.False(t, sb.IsFocused)
	assert.Equal(t, "Search...", sb.TextInput.Placeholder)
}

func TestSearchBar_ToggleFocus(t *testing.T) {
	tt := []struct {
		name          string
		toggleCount   int
		expectedFocus bool
	}{
		{"unfocused after 0 toggles", 0, false},
		{"focused after 1 toggle", 1, true},
		{"unfocused after 2 toggles", 2, false},
		{"focused after 3 toggles", 3, true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sb := newSearchBar("")
			for range tc.toggleCount {
				sb.ToggleFocus()
			}

			assert.Equal(t, tc.expectedFocus, sb.IsFocused)
		})
	}
}

func TestSearchBar_UpdateInput_GuardWhenNotFocused(t *testing.T) {
	tt := []struct {
		name string
		msg  tea.Msg
	}{
		{"letter key ignored", tea.KeyPressMsg{Code: 'a'}},
		{"enter ignored", tea.KeyPressMsg{Code: tea.KeyEnter}},
		{"backspace ignored", tea.KeyPressMsg{Code: tea.KeyBackspace}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sb := newSearchBar("")
			assert.Nil(t, sb.UpdateInput(tc.msg))
		})
	}
}

func TestSearchBar_IsSearching(t *testing.T) {
	tt := []struct {
		name          string
		toggleCount   int
		expectedValue bool
	}{
		{"not searching initially", 0, false},
		{"searching after focus", 1, true},
		{"not searching after blur", 2, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sb := newSearchBar("")
			for range tc.toggleCount {
				sb.ToggleFocus()
			}

			assert.Equal(t, tc.expectedValue, sb.IsSearching())
		})
	}
}
