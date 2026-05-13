package user_repo

import (
	"tui/internal/ui/common"
	"tui/internal/ui/context"

	tea "charm.land/bubbletea/v2"
)

type SearchBarModel struct {
	*common.SearchBar
}

func NewSearchBar(ctx *context.App, placeholderText string) *SearchBarModel {
	return &SearchBarModel{SearchBar: common.NewSearchBar(ctx, placeholderText)}
}

func (m *SearchBarModel) Update(msg tea.Msg, isSidebarFocused bool) tea.Cmd {
	if isSidebarFocused {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "/":
			return m.ToggleFocus()
		case "esc":
			if m.IsFocused {
				m.TextInput.Reset()
			}
			return nil
		}
	case tea.WindowSizeMsg:
		m.TextInput.SetWidth(m.Ctx.MainWidth - 10)
	}

	return m.UpdateInput(msg)
}
