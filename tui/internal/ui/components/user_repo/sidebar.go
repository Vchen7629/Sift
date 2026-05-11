package user_repo

import (
	"fmt"
	"tui/internal/service"
	"tui/internal/types"
	"tui/internal/ui/common"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Sidebar struct {
	ctx                *context.App
	viewport           viewport.Model
	FocusedGHRepo      *types.GHRepository
	FocusedIndexedRepo *types.IndexedRepo
	FocusedIdx         int
}

func NewSidebar(ctx *context.App) *Sidebar {
	return &Sidebar{
		ctx:                ctx,
		FocusedGHRepo:      nil,
		FocusedIndexedRepo: nil,
		FocusedIdx:         0,
	}
}

func (m *Sidebar) ResetFocus() {
	m.FocusedIdx = 0
	m.viewport.GotoTop()
}

func (m Sidebar) Init() tea.Cmd {
	return nil
}

func (m *Sidebar) Update(msg tea.Msg, isSidebarFocused bool) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if !isSidebarFocused || m.FocusedIndexedRepo == nil {
			break
		}
		cardHeight := lipgloss.Height(m.dependencyCard(m.FocusedIdx, m.FocusedIndexedRepo.Dependencies[m.FocusedIdx]))
		switch msg.String() {
		case "down":
			if m.FocusedIdx < len(m.FocusedIndexedRepo.Dependencies)-1 {
				m.FocusedIdx++
				service.ScrollToFocused(&m.viewport, m.FocusedIdx, cardHeight)
			}
		case "up":
			if m.FocusedIdx > 0 {
				m.FocusedIdx--
				service.ScrollToFocused(&m.viewport, m.FocusedIdx, cardHeight)
			}
		}

	case tea.WindowSizeMsg:
		m.viewport.SetWidth(m.ctx.SidebarWidth)
		m.viewport.SetHeight(m.ctx.MainHeight - 8)
	}
	return nil
}

func (m *Sidebar) View() tea.View {
	if m.FocusedGHRepo == nil {
		return tea.NewView(lipgloss.NewStyle().Padding(1, 2).Render("loading repo info..."))
	}

	description := m.FocusedGHRepo.Description
	if m.FocusedGHRepo.Description == "" {
		description = "No description found for this repository"
	}
	repoDesc := lipgloss.NewStyle().MarginBottom(2).Render(description)

	content := lipgloss.JoinVertical(lipgloss.Top, m.sidebarHeader(), repoDesc, m.repoDependencyList().Content)

	padding := lipgloss.NewStyle().PaddingLeft(2).PaddingRight(2).PaddingTop(1).Width(m.ctx.SidebarWidth)

	return tea.NewView(padding.Render(content))
}

func (m *Sidebar) sidebarHeader() string {
	repoName := lipgloss.NewStyle().Foreground(m.ctx.SelectedTheme.AccentBright).Render(m.FocusedGHRepo.Name)
	lastUpdate := lipgloss.NewStyle().
		Foreground(styles.TextDim).
		Render(fmt.Sprintf("Updated %s", service.FormatRelativeDate(m.FocusedGHRepo.LastCommit)))

	spaceBetween := common.SpaceBetween(m.ctx.SidebarWidth, lipgloss.Width(repoName), lipgloss.Width(lastUpdate), 4)

	topBlock := lipgloss.JoinHorizontal(lipgloss.Left, repoName, spaceBetween, lastUpdate)
	marginBottom := lipgloss.NewStyle().MarginBottom(1)

	if m.FocusedIndexedRepo != nil {
		totalLibs := lipgloss.NewStyle().
			Foreground(styles.TextDim).
			MarginRight(1).
			Render(fmt.Sprintf("%d total dependencies", m.FocusedIndexedRepo.TotalDependencies))

		lastIndexed := lipgloss.NewStyle().
			Foreground(styles.TextDim).
			Render(fmt.Sprintf("· indexed %s", m.FocusedIndexedRepo.LastIndexed))

		botBlock := lipgloss.JoinHorizontal(lipgloss.Left, totalLibs, lastIndexed)

		return marginBottom.Render(lipgloss.JoinVertical(lipgloss.Top, topBlock, botBlock))
	}

	return marginBottom.Render(topBlock)
}

// todo: refactor this into reusable func since both list and this file use same style to create viewport
func (m *Sidebar) repoDependencyList() tea.View {
	var dependencyCards []string

	if m.FocusedIndexedRepo == nil {
		text := lipgloss.NewStyle().Foreground(styles.TextMuted).Render("This repo isn't indexed yet, press r to index the repo")

		return tea.NewView(text)
	}

	for i, dependency := range m.FocusedIndexedRepo.Dependencies {
		dependencyCards = append(dependencyCards, m.dependencyCard(i, dependency))
	}

	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, dependencyCards...))
	return tea.NewView(m.viewport.View())
}

func (m *Sidebar) dependencyCard(idx int, dependency types.Dependency) string {
	textColor := m.ctx.SelectedTheme.AccentMid
	if idx == m.FocusedIdx {
		textColor = m.ctx.SelectedTheme.AccentBright
	}

	name := lipgloss.NewStyle().Foreground(textColor).Render(dependency.Name)

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
	spaceBetween := common.SpaceBetween(m.ctx.SidebarWidth, lipgloss.Width(name), lipgloss.Width(rightBlock), 4)
	marginBottom := lipgloss.NewStyle().MarginBottom(1)

	return marginBottom.Render(lipgloss.JoinHorizontal(lipgloss.Left, name, spaceBetween, rightBlock))
}
