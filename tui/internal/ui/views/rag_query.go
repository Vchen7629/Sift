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
	Searchbar        *common.SearchBarModel
	ResponseDisplay  *rag_query.RagQueryResponseModel
	Sidebar          *rag_query.SidebarModel
	SelectedRepo     string
	isSidebarFocused bool
}

func NewRagQuery(ctx *context.App) *RagQueryModel {
	return &RagQueryModel{
		ctx:             ctx,
		SelectedRepo:    "Sift",
		ActionBar:       rag_query.NewActionBar(ctx),
		Searchbar:       common.NewSearchBar(ctx, "Describe Your Issue..."),
		ResponseDisplay: rag_query.NewRagQueryResponse(ctx),
		Sidebar:         rag_query.NewSidebar(ctx),
	}
}

func (m *RagQueryModel) Init() tea.Cmd {
	fetchIndexedRepoCmd := m.Sidebar.Init()

	return tea.Batch(fetchIndexedRepoCmd)
}

// user actions
func (m *RagQueryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ActionBar.Update(msg)
		m.Searchbar.Update(msg, false)
		m.ResponseDisplay.Update(msg, false)
		m.Sidebar.Update(msg, false)
		return m, nil

	case rag_query.ToggleFocusMsg:
		if !m.Searchbar.IsSearching() {
			m.isSidebarFocused = !m.isSidebarFocused
		}
		return m, nil

	case rag_query.SelectRepoMsg:
		m.SelectedRepo = msg.RepoName

		return m, nil
	}

	actionBarCmd := m.ActionBar.Update(msg)
	searchBarCmd := m.Searchbar.Update(msg, m.isSidebarFocused)
	queryResCmd := m.ResponseDisplay.Update(msg, m.isSidebarFocused)
	sidebarCmd := m.Sidebar.Update(msg, m.isSidebarFocused)

	return m, tea.Batch(actionBarCmd, searchBarCmd, queryResCmd, sidebarCmd)
}

func (m *RagQueryModel) View() tea.View {
	leftPanel := lipgloss.JoinVertical(lipgloss.Top, m.Searchbar.View(), m.ResponseDisplay.View().Content)

	divider := common.VerticalDivider(m.ctx.MainHeight)

	mainContent := lipgloss.JoinHorizontal(lipgloss.Left, leftPanel, divider, m.Sidebar.View().Content)

	screen := lipgloss.JoinVertical(lipgloss.Top, m.ActionBar.View(m.isSidebarFocused, m.SelectedRepo).Content, mainContent)

	return tea.NewView(screen)
}
