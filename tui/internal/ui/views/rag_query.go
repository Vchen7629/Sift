package views

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tui/internal/ui/common"
	"tui/internal/ui/components/rag_query"
	"tui/internal/ui/context"
)

type RagQueryModel struct {
	ctx              *context.App
	ActionBar        *rag_query.ActionBarModel
	Searchbar        *rag_query.SearchBarModel
	ResponseDisplay  *rag_query.RagQueryResponseModel
	Sidebar          *rag_query.SidebarModel
	SelectedRepo     string
	isSidebarFocused bool
}

func NewRagQuery(ctx *context.App) *RagQueryModel {
	return &RagQueryModel{
		ctx:             ctx,
		SelectedRepo:    "",
		ActionBar:       rag_query.NewActionBar(ctx),
		Searchbar:       rag_query.NewSearchBar(ctx, "Describe Your Issue..."),
		ResponseDisplay: rag_query.NewRagQueryResponse(ctx),
		Sidebar:         rag_query.NewSidebar(ctx),
	}
}

func (m *RagQueryModel) Init() tea.Cmd {
	return m.Sidebar.Init()
}

// user actions
func (m *RagQueryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.ToggleFocusMsg:
		if !m.Searchbar.IsSearching() {
			m.isSidebarFocused = !m.isSidebarFocused
		}
		return m, nil

	case rag_query.SelectRepoMsg:
		selectedRepo := msg.RepoName

		m.SelectedRepo = selectedRepo

		return m, nil
	case rag_query.NewSearchQueryMsg:
		m.isSidebarFocused = false
		m.Searchbar.IsFocused = false
	}

	actionBarCmd := m.ActionBar.Update(msg)
	searchBarCmd := m.Searchbar.Update(msg, m.SelectedRepo)
	queryResCmd := m.ResponseDisplay.Update(msg, m.isSidebarFocused, m.Searchbar.IsSearching())
	sidebarCmd := m.Sidebar.Update(msg, m.isSidebarFocused, m.Searchbar.IsSearching())

	return m, tea.Batch(actionBarCmd, searchBarCmd, queryResCmd, sidebarCmd)
}

func (m *RagQueryModel) View() tea.View {
	leftPanel := lipgloss.JoinVertical(lipgloss.Top, m.Searchbar.View(), m.ResponseDisplay.View().Content)

	divider := common.VerticalDivider(m.ctx.MainHeight)

	mainContent := lipgloss.JoinHorizontal(lipgloss.Left, leftPanel, divider, m.Sidebar.View().Content)

	screen := lipgloss.JoinVertical(
		lipgloss.Top,
		m.ActionBar.View(m.isSidebarFocused, m.Searchbar.IsSearching(), m.SelectedRepo).Content,
		mainContent,
	)

	return tea.NewView(screen)
}
