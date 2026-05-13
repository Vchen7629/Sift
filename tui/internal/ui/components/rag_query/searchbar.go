package rag_query

import (
	"tui/internal/api"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type SearchBarModel struct {
	ctx          *context.App
	textInput 	 textinput.Model
	focused      bool
}

func NewSearchBar(ctx *context.App, placeholderText string) *SearchBarModel {
	ti := textinput.New()
	ti.Placeholder = placeholderText

	return &SearchBarModel{
		ctx:       ctx,
		textInput: ti,
	}
}

func (m *SearchBarModel) Init() tea.Cmd {
	return nil
}

func (m *SearchBarModel) Update(msg tea.Msg, isSidebarFocused bool, selectedRepo string) tea.Cmd {
	if isSidebarFocused {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "/":
			m.focused = !m.focused
			if m.focused {
				return m.textInput.Focus()
			}
			m.textInput.Blur()
			return nil
		case "esc":
			if m.focused {
				m.textInput.Reset()
			}
			return nil

		case "enter":
			if !m.focused || selectedRepo == "" {
				return nil
			}

			return m.newSearchQuery(selectedRepo)
		}
	case tea.WindowSizeMsg:
		m.textInput.SetWidth(m.ctx.MainWidth - 10)
	}

	if m.focused {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return cmd
	}

	return nil
}

func (m *SearchBarModel) View() string {
	s := m.textInput.Styles()
	s.Focused.Text = lipgloss.NewStyle().Foreground(m.ctx.SelectedTheme.AccentBright)
	m.textInput.SetStyles(s)

	borderColor := styles.Divider
	if m.focused {
		borderColor = m.ctx.SelectedTheme.AccentMid
	}

	style := lipgloss.NewStyle().
		MarginLeft(2).Width(m.ctx.MainWidth-4).Padding(0, 1).
		Border(lipgloss.RoundedBorder()).BorderForeground(borderColor)

	return style.Render(m.textInput.View())
}

// used to disable panel focus swap if user is searching
func (m *SearchBarModel) IsSearching() bool {
	return m.focused
}

type NewSearchQueryMsg struct { Res api.SearchRes }
type NewSearchQueryErr struct { RepoName, Err string }

func (m *SearchBarModel) newSearchQuery(repoName string) tea.Cmd {
	return func() tea.Msg {
		searchRes, err := api.Search(m.ctx.Username, repoName, m.textInput.Value())
		if err != nil {
			return NewSearchQueryErr{RepoName: repoName, Err: err.Error()}
		}

		return NewSearchQueryMsg{Res: searchRes}
	}
}