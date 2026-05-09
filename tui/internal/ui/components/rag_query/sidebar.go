package rag_query

import (
	"fmt"
	"strings"
	"tui/internal/api"
	"tui/internal/service"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"
	"tui/internal/types"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type SidebarModel struct {
	ctx          *context.App
	indexedRepos []types.IndexedRepo
	viewport     viewport.Model
	focused      types.IndexedRepo
}

type SelectRepoMsg struct{ RepoName string }

func NewSidebar(ctx *context.App) *SidebarModel {
	return &SidebarModel{
		ctx: ctx,
		indexedRepos: []types.IndexedRepo{},
	}
}

func (m *SidebarModel) Init() tea.Cmd {
	return m.fetchIndexedRepo
}

func (m *SidebarModel) Update(msg tea.Msg, isSidebarFocused bool) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if !isSidebarFocused {
			break
		}
		switch msg.String() {
		case "down":
			if m.focused.Id < len(m.indexedRepos) - 1 {
				m.focused.Id++
				m.focused = m.indexedRepos[m.focused.Id]
				service.ScrollToFocused(&m.viewport, m.focused.Id, 1)
			}
		case "up": 
			if m.focused.Id > 0 {
				m.focused.Id--
				m.focused = m.indexedRepos[m.focused.Id]
				service.ScrollToFocused(&m.viewport, m.focused.Id, 1)
			}
		
		case "enter":
			return func() tea.Msg { 
				return SelectRepoMsg{ RepoName: m.focused.Name } 
			}
		}

	case fetchIndexedRepoMsg:
		m.indexedRepos = msg.repos
		return nil
	}
	return nil
}

func (m *SidebarModel) View() tea.View {
	sidebarContent := m.sideBarList()
	if m.sideBarList() == "" {
		sidebarContent = lipgloss.NewStyle().Padding(1, 1).Render("No indexed repos, go to Your Repositories to index your repos")
	}

	content := lipgloss.JoinVertical(lipgloss.Top, m.header(), sidebarContent)

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

	for _, repo := range m.indexedRepos {
		textColor := m.ctx.SelectedTheme.AccentMid
		if repo.Id == m.focused.Id {
			textColor = m.ctx.SelectedTheme.AccentBright
		}

		repoName := lipgloss.NewStyle().PaddingLeft(2).Width(22).Foreground(textColor).Render(repo.Name)
		totalDependencies := lipgloss.NewStyle().Width(10).Foreground(textColor).Render(fmt.Sprintf("%d libs", repo.TotalDependencies))
		lastIndexed := lipgloss.NewStyle().Width(18).Foreground(textColor).Render(repo.LastIndexed)

		spaceBelow := lipgloss.NewStyle().MarginBottom(0)

		row := spaceBelow.Render(lipgloss.JoinHorizontal(lipgloss.Left, repoName, totalDependencies, lastIndexed))

		indexedRepoList = append(indexedRepoList, row)
	}

	m.viewport.SetHeight(m.ctx.MainHeight)
	m.viewport.SetWidth(m.ctx.SidebarWidth)
	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, indexedRepoList...))

	return m.viewport.View()
}

type fetchIndexedRepoMsg struct { repos []types.IndexedRepo }

func (m SidebarModel) fetchIndexedRepo() tea.Msg {
	indexRepos, err := api.GetAllIndexedRepos(m.ctx.Username)
	if err != nil {
		return err
	}

	return fetchIndexedRepoMsg{ repos: indexRepos }
}