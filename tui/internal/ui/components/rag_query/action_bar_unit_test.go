//go:build unit

package rag_query_test

import (
	"testing"
	"tui/internal/ui/common"
	"tui/internal/ui/components/rag_query"
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
		name     string
		key      rune
		expected func() tea.Msg
	}{
		{"pressing s returns the tea msg", 's', func() tea.Msg { return common.ToggleFocusMsg{} }},
		{"pressing enter does not return tea.msg", tea.KeyEnter, nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := rag_query.NewActionBar(actionBarCtx())
			cmd := m.Update(tea.KeyPressMsg{Code: tc.key})

			if tc.expected == nil {
				assert.Nil(t, cmd)
			} else {
				assert.Equal(t, tc.expected(), cmd())
			}
		})
	}
}
