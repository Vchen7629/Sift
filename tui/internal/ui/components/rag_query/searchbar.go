// todo: refactor this into a reusable component for both rag_query and user_repo
package rag_query

import (
	"tui/internal/ui/context"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type SearchBarModel struct {
	ctx 	  *context.App
	textInput textinput.Model
	focused   bool
}

func NewRagQuerySearchBar(ctx *context.App) *SearchBarModel {
	ti := textinput.New()
	ti.Placeholder = "Describe Your Issue..."

	return &SearchBarModel{
		ctx: 	   ctx,
		textInput: ti,
	}
}

func (m *SearchBarModel) Init() tea.Cmd {
	return nil
}

func (m *SearchBarModel) Update(msg tea.Msg) tea.Cmd {
	key, ok := msg.(tea.KeyPressMsg)	
	if ok && key.String() == "enter" {
		m.focused = !m.focused
		if m.focused {
			return m.textInput.Focus()
		}
		m.textInput.Blur()
		return nil
	}

	if m.focused && key.String() == "esc" {
		m.textInput.Reset()
	}

	if !m.focused {
		return nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return cmd
}

func (m *SearchBarModel) View() string {
	// border (2) + padding (2) + margin (2) = 6
	m.textInput.SetWidth(m.ctx.RepoListWidth - 6)

	style := lipgloss.NewStyle().
		MarginLeft(2).
		Width(m.ctx.RepoListWidth - 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Padding(0, 1)

	return style.Render(m.textInput.View())
}

// used to disable panel focus swap if user is searching
func (m *SearchBarModel) IsSearching() bool {
	return m.focused
}