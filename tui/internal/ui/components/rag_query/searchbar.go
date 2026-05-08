// todo: refactor this into a reusable component for both rag_query and user_repo
package rag_query

import (
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

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

func (m *SearchBarModel) Update(msg tea.Msg, isSidebarFocused bool) tea.Cmd {
	if isSidebarFocused {
		return nil
	}

	key, ok := msg.(tea.KeyPressMsg)	
	if ok && key.String() == "/" {
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

// Todo: Disable searchbar left right movement when theme is open
func (m *SearchBarModel) View() string {
	m.textInput.SetWidth(m.ctx.MainWidth - 10)
	s := m.textInput.Styles()
	s.Focused.Text = lipgloss.NewStyle().Foreground(m.ctx.SelectedTheme.AccentBright)
	m.textInput.SetStyles(s)

	borderColor := styles.Divider
	if m.focused {
		borderColor = m.ctx.SelectedTheme.AccentBright
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