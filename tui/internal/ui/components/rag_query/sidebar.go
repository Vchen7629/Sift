package rag_query

import (
	"fmt"
	"strings"
	"tui/internal/service"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type SidebarModel struct {
	ctx      *context.App
	repos    []repoStatus
	viewport viewport.Model
	focused  repoStatus
}

type repoStatus struct {
	id, totalDep      int
	name, lastIndexed string
}

type SelectRepoMsg struct{ RepoName string }

func NewSidebar(ctx *context.App) *SidebarModel {
	return &SidebarModel{
		ctx: ctx,
		repos: []repoStatus{
			{id: 0, name: "Sift", totalDep: 17, lastIndexed: "19"},
			{id: 1, name: "Atlaxiom", totalDep: 9, lastIndexed: "2"},
			{id: 2, name: "Cyphria", totalDep: 12, lastIndexed: "1"},
			{id: 3, name: "Kubernetes", totalDep: 7, lastIndexed: "5"},
			{id: 4, name: "Docker", totalDep: 26, lastIndexed: "10"},
			{id: 5, name: "Kafka", totalDep: 3, lastIndexed: "2"},
			{id: 6, name: "Splice", totalDep: 45, lastIndexed: "3"},
		},
	}
}

func (m *SidebarModel) Init() tea.Cmd {
	return nil
}

func (m *SidebarModel) Update(msg tea.Msg, isSidebarFocused bool) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if !isSidebarFocused {
			break
		}
		switch msg.String() {
		case "down":
			if m.focused.id < len(m.repos) - 1 {
				m.focused.id++
				m.focused = m.repos[m.focused.id]
				service.ScrollToFocused(&m.viewport, m.focused.id, 1)
			}
		case "up": 
			if m.focused.id > 0 {
				m.focused.id--
				m.focused = m.repos[m.focused.id]
				service.ScrollToFocused(&m.viewport, m.focused.id, 1)
			}
		
		case "enter":
			return func() tea.Msg { 
				return SelectRepoMsg{ RepoName: m.focused.name } 
			}
		}
	}
	return nil
}

func (m *SidebarModel) View() tea.View {
	content := lipgloss.JoinVertical(lipgloss.Top, m.header(), m.sideBarList())

	return tea.NewView(content)
}

func (m *SidebarModel) header() string {
	name := lipgloss.NewStyle().PaddingLeft(2).Width(22).Render("Repo Name")
	totalDep := lipgloss.NewStyle().Width(10).Render("Total")
	lastIndexed := lipgloss.NewStyle().Width(18).Render("Last Indexed")

	titleText := lipgloss.JoinHorizontal(lipgloss.Left, name, totalDep, lastIndexed)

	divider := lipgloss.NewStyle().
		Foreground(styles.Divider).
		Render(strings.Repeat("─", m.ctx.SidebarWidth - 1))

	return lipgloss.JoinVertical(lipgloss.Left, titleText, divider)
}

func (m *SidebarModel) sideBarList() string {
	var indexedRepoList []string

	for _, repo := range m.repos {
		textColor := m.ctx.SelectedTheme.AccentMid
		if repo.id == m.focused.id {
			textColor = m.ctx.SelectedTheme.AccentBright
		}

		repoName := lipgloss.NewStyle().PaddingLeft(2).Width(22).Foreground(textColor).Render(repo.name)
		totalDependencies := lipgloss.NewStyle().Width(10).Foreground(textColor).Render(fmt.Sprintf("%d libs", repo.totalDep))
		lastIndexed := lipgloss.NewStyle().Width(18).Foreground(textColor).Render(fmt.Sprintf("%s days ago", repo.lastIndexed))

		spaceBelow := lipgloss.NewStyle().MarginBottom(0)

		row := spaceBelow.Render(lipgloss.JoinHorizontal(lipgloss.Left, repoName, totalDependencies, lastIndexed))

		indexedRepoList = append(indexedRepoList, row)
	}

	m.viewport.SetHeight(m.ctx.MainHeight)
	m.viewport.SetWidth(m.ctx.SidebarWidth)
	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, indexedRepoList...))

	return m.viewport.View()
}