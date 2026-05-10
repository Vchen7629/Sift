package user_repo

import (
	"fmt"
	"image/color"
	"strconv"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tui/internal/api"
	"tui/internal/service"
	"tui/internal/types"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"
)

type ListModel struct {
	ctx 		     *context.App
	GHRepos     	 []types.GHRepository
	IndexedRepos     []types.IndexedRepo
	FocusedIdx       int
	ProcessingStatus map[int]string
	viewport 	 	 viewport.Model
	progressBar      *ProgressBarModel
}

func NewUserRepoList(ctx *context.App) *ListModel {
	m := &ListModel{
		ctx: ctx,
		GHRepos: []types.GHRepository{},
		ProcessingStatus: map[int]string{},
		progressBar: NewProgressBar(ctx),
	}
	m.FocusedIdx = 0
	return m
}

func (m ListModel) Init() tea.Cmd {
	return m.progressBar.Init()
}

func (m *ListModel) Update(msg tea.Msg, isSidebarFocused bool) tea.Cmd {                                                                                                            
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if isSidebarFocused || len(m.GHRepos) == 0 {
			break
		}

		ir := service.FindIndexedRepo(m.GHRepos[m.FocusedIdx].Name, m.IndexedRepos)
		var indexedRepo types.IndexedRepo
		if ir != nil {
			indexedRepo = *ir
		}
		cardHeight := lipgloss.Height(m.repoCard(m.FocusedIdx, m.GHRepos[m.FocusedIdx], indexedRepo))

		switch msg.String() {
		case "down":
			if m.FocusedIdx < len(m.GHRepos) - 1 {
				m.FocusedIdx++
    			service.ScrollToFocused(&m.viewport, m.FocusedIdx, cardHeight)
			}
		case "up":
			if m.FocusedIdx > 0 {
				m.FocusedIdx--
				service.ScrollToFocused(&m.viewport, m.FocusedIdx, cardHeight)
			}
		}
	// trigger from user pressing r
	case IndexRepoRequestMsg:                                                                                                                                   
      	return m.IndexRepo
	
	// response from api call
	case indexRepoMsg:                                                                                                                                          
		m.ProcessingStatus[msg.idx] = msg.status
		return nil  

	// for the progress bar
	case tickMsg:
		return m.progressBar.Update(msg)

	case progress.FrameMsg:
		return m.progressBar.Update(msg)
	}

	return nil
}

func (m *ListModel) View() tea.View {
	var cards []string
	
	if len(m.GHRepos) > 0 {
		for i, repo := range m.GHRepos {
			ir := service.FindIndexedRepo(repo.Name, m.IndexedRepos)
			var indexedRepo types.IndexedRepo
			if ir != nil {
				indexedRepo = *ir
			}
			cards = append(cards, m.repoCard(i, repo, indexedRepo))
		}

		m.viewport.SetWidth(m.ctx.MainWidth)
		m.viewport.SetHeight(m.ctx.MainHeight - 4)
		m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, cards...))

		return tea.NewView(m.viewport.View())
	}

	fetchingPlaceholder := lipgloss.NewStyle().Padding(1, 2).Render("Loading your repos...")
	return tea.NewView(fetchingPlaceholder)
}

func (m *ListModel) repoCard(idx int, ghRepo types.GHRepository, indexedRepo types.IndexedRepo) string {
	borderColor, _ := m.focusedStyle(idx)

	content := lipgloss.JoinVertical(lipgloss.Top, m.header(idx, ghRepo, indexedRepo), m.progressBar.View().Content)
		
	card := lipgloss.NewStyle().
		Width(m.ctx.MainWidth).PaddingLeft(2).Padding(0,1).
		Border(lipgloss.RoundedBorder()).BorderForeground(borderColor).
		Render(content)

	return card
}

func (m *ListModel) header(idx int, ghRepo types.GHRepository, indexedRepo types.IndexedRepo) string {
	_, textColor := m.focusedStyle(idx)

	repoName := lipgloss.NewStyle().Foreground(textColor).Render(ghRepo.Name)

	if status, exists := m.ProcessingStatus[idx]; exists {
		indexStatus := lipgloss.NewStyle().Width(12).Align(lipgloss.Right).Render(status)                                                                       
		lastIndexed := lipgloss.NewStyle().Width(12).Align(lipgloss.Right).Render(indexedRepo.LastIndexed)
		totalDependencies := lipgloss.NewStyle().Width(5).Align(lipgloss.Right).Render(strconv.Itoa(indexedRepo.TotalDependencies))

		right := lipgloss.JoinHorizontal(lipgloss.Top, indexStatus, lastIndexed, totalDependencies)

		spacer := lipgloss.NewStyle().
			Width(m.ctx.MainWidth - lipgloss.Width(repoName) - lipgloss.Width(right) - 4).
			Render("")

		return lipgloss.JoinHorizontal(lipgloss.Top, repoName, spacer, right)
	} 

	indexStatus := lipgloss.NewStyle().Width(12).Align(lipgloss.Right).Render("Unindexed")       

	spacer := lipgloss.NewStyle().
		Width(m.ctx.MainWidth - lipgloss.Width(repoName) - lipgloss.Width(indexStatus) - 4).
		Render("")

	return lipgloss.JoinHorizontal(lipgloss.Top, repoName, spacer, indexStatus)
}

func (m *ListModel) focusedStyle(idx int) (color.Color, color.Color) {
	if m.FocusedIdx == idx {
		return m.ctx.SelectedTheme.AccentMid, m.ctx.SelectedTheme.AccentBright
	}

	return styles.Divider, lipgloss.Color("#ffffff")
}

type indexRepoMsg struct { idx int; status string }

func (m *ListModel) IndexRepo() tea.Msg {
	gitUser := m.ctx.Username

	err := api.IndexRepo(gitUser, fmt.Sprintf("%s/%s", gitUser, m.GHRepos[m.FocusedIdx].Name))
	if err != nil {
		return indexRepoMsg{idx: m.FocusedIdx, status: err.Error()}
	}

	return indexRepoMsg{idx: m.FocusedIdx, status: "processing"}
}