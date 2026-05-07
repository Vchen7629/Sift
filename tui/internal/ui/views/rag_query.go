package views

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tui/internal/ui/context"
	"tui/internal/ui/components/rag_query"
)

type RagQueryModel struct {
	Ctx *context.App
	Searchbar *rag_query.SearchBarModel
}

func (m RagQueryModel) Init() tea.Cmd {
	return nil
}

// user actions
func (m RagQueryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	searchBarCmd := m.Searchbar.Update(msg)

	return m, tea.Batch(searchBarCmd)
}

func (m RagQueryModel) View() tea.View {
	content := lipgloss.JoinVertical(lipgloss.Top, m.ragQueryActionBar().Content, m.Searchbar.View())

	return tea.NewView(content)
}

func (m RagQueryModel) ragQueryActionBar() tea.View {
	navBtnStyle := lipgloss.NewStyle().PaddingLeft(2)

	navBtnTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444444")).
		Bold(true)

	searchBtn := navBtnStyle.Render(navBtnTextStyle.Render("[↵] search"))
	switchRepoBtn := navBtnStyle.Render(navBtnTextStyle.Render("[s] switch repo"))
	clearSearchBtn := navBtnStyle.Render(navBtnTextStyle.Render("[esc] clear"))

	return tea.NewView(lipgloss.NewStyle().
		BorderBottom(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderBottomForeground(lipgloss.Color("#444444")).
		Width(m.Ctx.Width - 2).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, searchBtn, switchRepoBtn, clearSearchBtn)))
}