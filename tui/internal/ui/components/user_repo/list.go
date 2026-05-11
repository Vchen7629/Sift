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
	progressBars     map[int]*ProgressBarModel
}

func NewUserRepoList(ctx *context.App) *ListModel {
	m := &ListModel{
		ctx: ctx,
		GHRepos: []types.GHRepository{},
		ProcessingStatus: map[int]string{},
		progressBars: map[int]*ProgressBarModel{},
	}
	m.FocusedIdx = 0
	return m
}

func (m ListModel) Init() tea.Cmd {
	return nil
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

		pb := NewProgressBar(m.ctx, msg.idx, fmt.Sprintf("%s/%s", m.ctx.Username, m.GHRepos[m.FocusedIdx].Name))
		m.progressBars[msg.idx] = pb
		return pb.Init()

	// for the progress bar
	case tickMsg:
		if pb, ok := m.progressBars[msg.idx]; ok {
			return pb.Update(msg)
		}

	case progress.FrameMsg:
		var cmds []tea.Cmd
		for _, pb := range m.progressBars {
			cmds = append(cmds, pb.Update(msg))
		}
		return tea.Batch(cmds...)

	case statusUpdateMsg:
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
	borderColor, textColor := m.focusedStyle(idx)
	repoName := lipgloss.NewStyle().Foreground(textColor).Render(ghRepo.Name)

	// right text section of the top row of the card
	var rightText string
	status, exists := m.ProcessingStatus[idx]
	switch {
	case !exists:
		rightText = lipgloss.NewStyle().Width(12).Align(lipgloss.Right).Render("Unindexed")
	case status == "Indexed":
		indexStatus := lipgloss.NewStyle().Width(16).MarginRight(1).Align(lipgloss.Right).Render(status)
		lastIndexed := lipgloss.NewStyle().Width(16).Align(lipgloss.Right).Render(indexedRepo.LastIndexed)
		totalDependencies := lipgloss.NewStyle().
			Width(17).PaddingLeft(1).Align(lipgloss.Right).
			Render(fmt.Sprintf("· %s dependencies", strconv.Itoa(indexedRepo.TotalDependencies)))

		rightText = lipgloss.JoinHorizontal(lipgloss.Top, indexStatus, lastIndexed, totalDependencies)
	default:
		rightText = lipgloss.NewStyle().Width(16).Align(lipgloss.Right).Render(status)
	}

	spacer := lipgloss.NewStyle().
		Width(m.ctx.MainWidth - lipgloss.Width(repoName) - lipgloss.Width(rightText) - 4).
		Render("")
	header := lipgloss.JoinHorizontal(lipgloss.Top, repoName, spacer, rightText)

	var content string
	if pb, ok := m.progressBars[idx]; ok {
		content = lipgloss.JoinVertical(lipgloss.Top, header, pb.View().Content)
	} else {
		content = header
	}
		
	card := lipgloss.NewStyle().
		Width(m.ctx.MainWidth).PaddingLeft(2).Padding(0,1).
		Border(lipgloss.RoundedBorder()).BorderForeground(borderColor).
		Render(content)

	return card
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
		return indexRepoMsg{idx: m.FocusedIdx, status: "error indexing"}
	}

	return indexRepoMsg{idx: m.FocusedIdx, status: "processing"}
}