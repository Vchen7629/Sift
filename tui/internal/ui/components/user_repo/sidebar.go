package user_repo

import (
	"fmt"
	"tui/internal/service"
	"tui/internal/types"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Sidebar struct {
	ctx 		 *context.App
	viewport 	 viewport.Model
	FocusedRepo  *types.Repository
	FocusedIdx   int
}

func NewSidebar(ctx *context.App) *Sidebar {
	return &Sidebar{
		ctx: ctx,
		FocusedRepo: nil,

		FocusedIdx: 0,
	}
}

func (m Sidebar) Init() tea.Cmd {
	return nil
}

func (m *Sidebar) Update(msg tea.Msg, isSidebarFocused bool) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if !isSidebarFocused {
			break
		}
		switch msg.String() {
		case "down":
			m.viewport.ScrollDown(2)
		case "up":
			m.viewport.ScrollUp(2)
		}
	}
	return nil
}

func (m *Sidebar) View() tea.View {
	if m.FocusedRepo == nil {
		return tea.NewView(lipgloss.NewStyle().Padding(1, 2).Render("Loading Repo Data..."))
	}

	description := m.FocusedRepo.Description
	if m.FocusedRepo.Description == "" {
		description = "No description found for this repository"
	}
	repoDesc := lipgloss.NewStyle().MarginBottom(2).Render(description)

	content := lipgloss.JoinVertical(lipgloss.Top, m.sidebarHeader(), repoDesc, m.repoDependencyList().Content)

	padding := lipgloss.NewStyle().PaddingLeft(2).PaddingRight(2).PaddingTop(1).Width(m.ctx.SidebarWidth)

	return tea.NewView(padding.Render(content))
}

func (m *Sidebar) sidebarHeader() string {
	repoName := lipgloss.NewStyle().Foreground(m.ctx.SelectedTheme.AccentBright).Render(m.FocusedRepo.Name)
	lastUpdate := lipgloss.NewStyle().
		Foreground(styles.TextDim).
		Render(fmt.Sprintf("Updated %s", service.FormatRelativeDate(m.FocusedRepo.LastUpdated)))

	spaceBetween := lipgloss.NewStyle().
		Width(m.ctx.SidebarWidth - 4 - lipgloss.Width(repoName) - lipgloss.Width(lastUpdate)).
		Render("")

	topBlock := lipgloss.JoinHorizontal(lipgloss.Left, repoName, spaceBetween, lastUpdate)

	totalLibs := lipgloss.NewStyle().
		Foreground(styles.TextDim).
		MarginRight(1).
		Render(fmt.Sprintf("%d total dependencies", m.FocusedRepo.TotalDependencies))

	lastIndexed := lipgloss.NewStyle().
		Foreground(styles.TextDim).
		Render(fmt.Sprintf("· indexed %s ago", m.FocusedRepo.LastIndexed))
		
	botBlock := lipgloss.JoinHorizontal(lipgloss.Left, totalLibs, lastIndexed)
	
	marginBottom := lipgloss.NewStyle().MarginBottom(1)

	return marginBottom.Render(lipgloss.JoinVertical(lipgloss.Top, topBlock, botBlock))
}

// todo: refactor this into reusable func since both list and this file use same style to create viewport
func (m *Sidebar) repoDependencyList() tea.View {
	var dependencyCards []string

	if len(m.FocusedRepo.Dependencies) == 0 {
		text := lipgloss.NewStyle().Foreground(styles.TextMuted).Render("No dependencies indexed for this repo")

		return tea.NewView(text)
	}

	for _, dependency := range m.FocusedRepo.Dependencies {
		dependencyCards = append(dependencyCards, m.dependencyCard(dependency))
	}

	m.viewport.SetWidth(m.ctx.MainWidth)
	m.viewport.SetHeight(m.ctx.MainHeight - 8)
	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, dependencyCards...))
	return tea.NewView(m.viewport.View())
}

func (m *Sidebar) dependencyCard(dependency types.DependencyStatus) string {
	name := lipgloss.NewStyle().Render(dependency.Name)

	version := lipgloss.NewStyle().Width(10).MarginRight(2).Align(lipgloss.Center).
		Background(lipgloss.Blue).Render(dependency.Version)

	statusText := lipgloss.White
	switch dependency.Status {
	case "healthy":
		statusText = lipgloss.Green
	case "deprecated":
		statusText = lipgloss.Red
	case "archived":
		statusText = lipgloss.Yellow
	}
	
	status := lipgloss.NewStyle().Foreground(statusText).Width(10).Render(dependency.Status)
	rightBlock := lipgloss.JoinHorizontal(lipgloss.Left, version, status)

	spaceBetween := lipgloss.NewStyle().
		Width(m.ctx.SidebarWidth - 4 - lipgloss.Width(name) - lipgloss.Width(rightBlock)).
		Render("")
	
	marginBottom := lipgloss.NewStyle().MarginBottom(1)

	return marginBottom.Render(lipgloss.JoinHorizontal(lipgloss.Left, name, spaceBetween, rightBlock))
}