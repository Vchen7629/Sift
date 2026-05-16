package common

import (
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type SearchBar struct {
	Ctx       *context.App
	TextInput textinput.Model
	IsFocused bool
}

func NewSearchBar(ctx *context.App, placeholderText string) *SearchBar {
	ti := textinput.New()
	ti.Placeholder = placeholderText

	return &SearchBar{
		Ctx:       ctx,
		TextInput: ti,
	}
}

func (m *SearchBar) Init() tea.Cmd {
	return nil
}

func (m *SearchBar) ToggleFocus() tea.Cmd {
	m.IsFocused = !m.IsFocused
	if m.IsFocused {
		return m.TextInput.Focus()
	}
	m.TextInput.Blur()
	return nil
}

func (m *SearchBar) UpdateInput(msg tea.Msg) tea.Cmd {
	if !m.IsFocused {
		return nil
	}

	var cmd tea.Cmd
	m.TextInput, cmd = m.TextInput.Update(msg)
	return cmd
}

func (m *SearchBar) View() string {
	s := m.TextInput.Styles()
	s.Focused.Text = lipgloss.NewStyle().Foreground(m.Ctx.SelectedTheme.AccentBright)
	m.TextInput.SetStyles(s)

	borderColor := styles.Divider
	if m.IsFocused {
		borderColor = m.Ctx.SelectedTheme.AccentMid
	}

	style := lipgloss.NewStyle().
		MarginLeft(2).Width(m.Ctx.MainWidth-4).Padding(0, 1).
		Border(lipgloss.RoundedBorder()).BorderForeground(borderColor)

	return style.Render(m.TextInput.View())
}

func (m *SearchBar) IsSearching() bool {
	return m.IsFocused
}
