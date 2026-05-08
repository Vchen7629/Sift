package views

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tui/internal/ui/components/rag_query"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"
)

type RagQueryModel struct {
	Ctx 		    *context.App
	ActionBar 	    *rag_query.ActionBarModel
	Searchbar 	    *rag_query.SearchBarModel
	ResponseDisplay *rag_query.RagQueryResponseModel
	Sidebar			*rag_query.SidebarModel
	SelectedRepo    string
	FocusRepoList   bool
}

func (m RagQueryModel) Init() tea.Cmd {
	return nil
}

// user actions
func (m RagQueryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	searchBarCmd := m.Searchbar.Update(msg)
	_, queryResCmd := m.ResponseDisplay.Update(msg)

	return m, tea.Batch(searchBarCmd, queryResCmd)
}

func (m RagQueryModel) View() tea.View {
	leftPanel := lipgloss.JoinVertical(lipgloss.Top, m.Searchbar.View(), m.ResponseDisplay.View().Content)

	dividerLine := strings.Repeat("│\n", m.Ctx.MainHeight - 1) + "│"
	divider := lipgloss.NewStyle().Foreground(styles.Divider).Render(dividerLine)

	mainContent := lipgloss.JoinHorizontal(lipgloss.Left, leftPanel, divider, m.Sidebar.View().Content)

	screen := lipgloss.JoinVertical(lipgloss.Top, m.ActionBar.View().Content, mainContent)

	return tea.NewView(screen)
}