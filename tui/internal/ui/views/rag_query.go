package views

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tui/internal/ui/components/rag_query"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"
)

type RagQueryModel struct {
	Ctx 		  *context.App
	Searchbar 	  *rag_query.SearchBarModel
	SelectedRepo  string
	QueryResponse *rag_query.RagQueryResponseModel
}

func (m RagQueryModel) Init() tea.Cmd {
	return nil
}

// user actions
func (m RagQueryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	searchBarCmd := m.Searchbar.Update(msg)
	_, queryResCmd := m.QueryResponse.Update(msg)

	return m, tea.Batch(searchBarCmd, queryResCmd)
}

func (m RagQueryModel) View() tea.View {
	content := lipgloss.JoinVertical(
		lipgloss.Top, m.ragQueryActionBar().Content, m.Searchbar.View(), m.QueryResponse.View().Content,
	)

	return tea.NewView(content)
}

func (m RagQueryModel) ragQueryActionBar() tea.View {
	navBtnStyle := lipgloss.NewStyle().PaddingLeft(2)

	navBtnTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444444")).
		Bold(true)

	selectedRepoText := lipgloss.NewStyle().
		PaddingRight(1).Foreground(styles.Warm.AccentBright).
		Render(fmt.Sprintf("Selected Repo: %s", m.SelectedRepo))

	selectedRepo := navBtnStyle.Render(selectedRepoText)
	searchBtn := navBtnStyle.Render(navBtnTextStyle.Render("[/] new query"))
	scrollBtn := navBtnStyle.Render(navBtnTextStyle.Render("[↑↓] navigate"))
	openBrowserBtn := navBtnStyle.Render(navBtnTextStyle.Render("[↵] open in browser"))
	switchRepoBtn := navBtnStyle.Render(navBtnTextStyle.Render("[s] switch repo"))

	return tea.NewView(lipgloss.NewStyle().
		BorderBottom(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderBottomForeground(lipgloss.Color("#444444")).
		Width(m.Ctx.WindowWidth - 2).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, selectedRepo, searchBtn, scrollBtn, switchRepoBtn, openBrowserBtn)))
}