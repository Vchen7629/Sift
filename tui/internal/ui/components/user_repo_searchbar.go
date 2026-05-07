package components

import (
	"tui/internal/ui/context"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type UserRepoSearchBarModel struct {
	ctx 	  *context.App
	textInput textinput.Model
	focused   bool
}

func NewUserRepoSearchBar(ctx *context.App) *UserRepoSearchBarModel {
	ti := textinput.New()
	ti.Placeholder = "Search Your Repositories..."

	return &UserRepoSearchBarModel{
		ctx: 	   ctx,
		textInput: ti,
	}
}

func (m *UserRepoSearchBarModel) Init() tea.Cmd {
	return nil
}

func (m *UserRepoSearchBarModel) Update(msg tea.Msg) tea.Cmd {
	key, ok := msg.(tea.KeyPressMsg)
	if ok && key.String() == "enter" {
		m.focused = !m.focused
		if m.focused {
			return m.textInput.Focus()
		}
		m.textInput.Blur()
		return nil
	}

	if !m.focused {
		return nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return cmd
}

func (m *UserRepoSearchBarModel) View() string {
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