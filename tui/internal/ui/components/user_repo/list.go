package user_repo

import (
	"fmt"
	"image/color"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tui/internal/api"
	"tui/internal/service"
	"tui/internal/types"
	"tui/internal/ui/common"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"
)

type ListModel struct {
	ctx              *context.App
	GHRepos          []types.GHRepository
	IndexedRepos     []types.IndexedRepo
	FocusedIdx       int
	ProcessingStatus map[int]string
	viewport         viewport.Model
	progressBars     map[int]*ProgressBarModel
}

func NewUserRepoList(ctx *context.App) *ListModel {
	m := &ListModel{
		ctx:              ctx,
		GHRepos:          []types.GHRepository{},
		ProcessingStatus: map[int]string{},
		progressBars:     map[int]*ProgressBarModel{},
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

		// progress bar adds 2 lines of height vs without
		cardHeight := 3
		if _, ok := m.progressBars[m.FocusedIdx]; ok {
			cardHeight = 5
		}

		switch msg.String() {
		case "down":
			if m.FocusedIdx < len(m.GHRepos)-1 {
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
		m.viewport.SetWidth(m.ctx.MainWidth)
		m.viewport.SetHeight(m.ctx.MainHeight - 4)

	// trigger from user pressing r
	case IndexRepoRequestMsg:
		if service.FindIndexedRepo(m.GHRepos[m.FocusedIdx].Name, m.IndexedRepos) != nil {
			return nil // prevents an already indexed repo from sending a new index repo request
		}
		return m.IndexRepo

	// response from api call
	case indexRepoMsg:
		pb := NewProgressBar(m.ctx, msg.idx, fmt.Sprintf("%s/%s", m.ctx.Username, m.GHRepos[m.FocusedIdx].Name))
		m.progressBars[msg.idx] = pb
		return pb.Init()

	case indexRepoErrMsg:
		m.ProcessingStatus[msg.idx] = fmt.Sprintf("error: %s", msg.err.Error())

	// for the progress bar
	case tickMsg:
		statusText := "new index job request"
		if msg.status != "" {
			statusText = msg.status
		}

		m.ProcessingStatus[msg.idx] = statusText
		if pb, ok := m.progressBars[msg.idx]; ok {
			return pb.Update(msg)
		}

	case progress.FrameMsg:
		var cmds []tea.Cmd
		for _, pb := range m.progressBars {
			cmds = append(cmds, pb.Update(msg))
		}
		return tea.Batch(cmds...)
	}

	return nil
}

func (m *ListModel) View() tea.View {
	var cards []string

	if len(m.GHRepos) <= 0 {
		return tea.NewView(lipgloss.NewStyle().Padding(1, 2).Render("Loading your repos..."))
	}

	for i, repo := range m.GHRepos {
		ir := service.FindIndexedRepo(repo.Name, m.IndexedRepos)
		var indexedRepo types.IndexedRepo
		if ir != nil {
			indexedRepo = *ir
		}

		// color switch to alternate between two colors for focused/unfocused
		borderColor, textColor := styles.Divider, lipgloss.Color("#ffffff")
		if m.FocusedIdx == i {
			borderColor, textColor = m.ctx.SelectedTheme.AccentMid, m.ctx.SelectedTheme.AccentBright
		}
		header := m.cardHeader(i, indexedRepo, repo, textColor)

		// just render header text if progress bar doesnt exist
		var content string
		if pb, ok := m.progressBars[i]; ok {
			content = lipgloss.JoinVertical(lipgloss.Top, header, pb.View().Content)
		} else {
			content = header
		}

		card := lipgloss.NewStyle().
			Width(m.ctx.MainWidth).PaddingLeft(2).Padding(0, 1).
			Border(lipgloss.RoundedBorder()).BorderForeground(borderColor).
			Render(content)

		cards = append(cards, card)
	}
	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, cards...))

	return tea.NewView(m.viewport.View())
}

// This is the top row in each list card, shows things like name, index status, total deps, and last indexed time
func (m *ListModel) cardHeader(
	idx int, indexedRepo types.IndexedRepo, ghRepo types.GHRepository, textColor color.Color,
) string {
	repoName := lipgloss.NewStyle().Foreground(textColor).Render(ghRepo.Name)

	// index metadata like index status, lastIndexed and  totalDependencies
	var indexMetadata string
	status, exists := m.ProcessingStatus[idx]
	switch {
	case !exists:
		indexMetadata = lipgloss.NewStyle().Render("Unindexed")
	case status == "processed" || status == "Indexed":
		indexStatus := lipgloss.NewStyle().MarginRight(1).Render(status)
		lastIndexed := lipgloss.NewStyle().Render(indexedRepo.LastIndexed)
		totalDependencies := lipgloss.NewStyle().PaddingLeft(1).Render(fmt.Sprintf("· %d dependencies", indexedRepo.TotalDependencies))

		indexMetadata = lipgloss.JoinHorizontal(lipgloss.Top, indexStatus, lastIndexed, totalDependencies)
	default:
		indexMetadata = lipgloss.NewStyle().Render(status)
	}

	spacer := common.SpaceBetween(m.ctx.MainWidth, lipgloss.Width(repoName), lipgloss.Width(indexMetadata), 4)

	return lipgloss.JoinHorizontal(lipgloss.Top, repoName, spacer, indexMetadata)
}

type indexRepoMsg struct{ idx int }
type indexRepoErrMsg struct {
	idx int
	err error
}

func (m *ListModel) IndexRepo() tea.Msg {
	gitUser := m.ctx.Username

	err := api.IndexRepo(gitUser, fmt.Sprintf("%s/%s", gitUser, m.GHRepos[m.FocusedIdx].Name))
	if err != nil {
		return indexRepoErrMsg{idx: m.FocusedIdx, err: err}
	}

	return indexRepoMsg{idx: m.FocusedIdx}
}
