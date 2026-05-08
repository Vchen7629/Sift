package user_repo

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

func NewUserRepoSearchBar(ctx *context.App) *SearchBarModel {
	ti := textinput.New()
	ti.Placeholder = "Search Your Repositories..."

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
	m.textInput.SetWidth(m.ctx.MainWidth - 10)
	s := m.textInput.Styles()
	s.Focused.Text = lipgloss.NewStyle().Foreground(m.ctx.SelectedTheme.AccentBright)
	m.textInput.SetStyles(s)

	borderColor := lipgloss.Color("#444444")
	if m.focused {
		borderColor = m.ctx.SelectedTheme.AccentMid
	}

	style := lipgloss.NewStyle().
		MarginLeft(2).Width(m.ctx.MainWidth - 4).Padding(0, 1).
		Border(lipgloss.RoundedBorder()).BorderForeground(borderColor)

	return style.Render(m.textInput.View())
}

// used to disable panel focus swap if user is searching
func (m *SearchBarModel) IsSearching() bool {
	return m.focused
}