package user_repo

import (
	"fmt"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type FocusedSidebar struct {
	ctx 		 *context.App
	viewport 	 viewport.Model
	FocusedRepo  FocusedRepo
}

type DependencyStatus struct {
	name, version, status string
}

func NewFocusedSidebar(ctx *context.App) *FocusedSidebar {
	return &FocusedSidebar{
		ctx: ctx,
		FocusedRepo: FocusedRepo{},
	}
}

func (m FocusedSidebar) Init() tea.Cmd {
	return nil
}

func (m *FocusedSidebar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "down":
			m.viewport.ScrollDown(2)
		case "up":
			m.viewport.ScrollUp(2)
		}
	}
	return m, nil
}

func (m *FocusedSidebar) View() tea.View {
	repoDesc := lipgloss.NewStyle().MarginBottom(2).Render(m.FocusedRepo.userRepo.Description)

	content := lipgloss.JoinVertical(lipgloss.Top, m.focusedHeader(), repoDesc, m.repoDependencyList().Content)

	padding := lipgloss.NewStyle().PaddingLeft(2).PaddingRight(2).PaddingTop(1).Width(m.ctx.SidebarWidth)

	return tea.NewView(padding.Render(content))
}

func (m *FocusedSidebar) focusedHeader() string {
	repoName := lipgloss.NewStyle().Foreground(styles.Warm.AccentBright).Render(m.FocusedRepo.userRepo.Name)

	totalLibs := lipgloss.NewStyle().Foreground(styles.TextDim).MarginRight(1).
		Render(fmt.Sprintf("%s total dependencies", m.FocusedRepo.userRepo.TotalDependencies))
	lastIndexed := lipgloss.NewStyle().Foreground(styles.TextDim).
		Render(fmt.Sprintf("· %s ago", m.FocusedRepo.userRepo.LastIndexed))
	rightBlock := lipgloss.JoinHorizontal(lipgloss.Left, totalLibs, lastIndexed)
	
	spaceBetween := lipgloss.NewStyle().
		Width(m.ctx.SidebarWidth - 4 - lipgloss.Width(repoName) - lipgloss.Width(rightBlock)).
		Render("")
	
	marginBottom := lipgloss.NewStyle().MarginBottom(1)

	return marginBottom.Render(lipgloss.JoinHorizontal(lipgloss.Left, repoName, spaceBetween, rightBlock))
}

// todo: refactor this into reusable func since both list and this file use same style to create viewport
func (m *FocusedSidebar) repoDependencyList() tea.View {
	var dependencyCards []string

	if len(m.FocusedRepo.userRepo.Dependencies) == 0 {
		text := lipgloss.NewStyle().Foreground(styles.TextMuted).Render("No dependencies indexed for this repo")

		return tea.NewView(text)
	}

	for _, dependency := range m.FocusedRepo.userRepo.Dependencies {
		dependencyCards = append(dependencyCards, m.dependencyCard(dependency))
	}

	m.viewport.SetWidth(m.ctx.RepoListWidth)
	m.viewport.SetHeight(m.ctx.RepoListHeight - 8)
	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, dependencyCards...))
	return tea.NewView(m.viewport.View())
}

func (m *FocusedSidebar) dependencyCard(dependency DependencyStatus) string {
	name := lipgloss.NewStyle().Render(dependency.name)

	version := lipgloss.NewStyle().Width(10).MarginRight(2).Align(lipgloss.Center).
		Background(lipgloss.Blue).Render(dependency.version)

	statusText := lipgloss.White
	switch dependency.status {
	case "healthy":
		statusText = lipgloss.Green
	case "deprecated":
		statusText = lipgloss.Red
	case "archived":
		statusText = lipgloss.Yellow
	}
	
	status := lipgloss.NewStyle().Foreground(statusText).Width(10).Render(dependency.status)
	rightBlock := lipgloss.JoinHorizontal(lipgloss.Left, version, status)

	spaceBetween := lipgloss.NewStyle().
		Width(m.ctx.SidebarWidth - 4 - lipgloss.Width(name) - lipgloss.Width(rightBlock)).
		Render("")
	
	marginBottom := lipgloss.NewStyle().MarginBottom(1)

	return marginBottom.Render(lipgloss.JoinHorizontal(lipgloss.Left, name, spaceBetween, rightBlock))
}