package user_repo

import (
	"fmt"
	"image/color"
	"time"

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
	GHRepos          []api.RepoApiRes
	IndexedRepoMap   map[string]*types.IndexedRepo
	FocusedIdx       int
	ProcessingStatus map[int]string
	FetchError       string
	viewport         viewport.Model
	progressBars     map[int]*ProgressBarModel
	pendingCleanup   map[int]bool
	noSearchResults  bool
}

func NewUserRepoList(ctx *context.App) *ListModel {
	m := &ListModel{
		ctx:              ctx,
		GHRepos:          []api.RepoApiRes{},
		ProcessingStatus: map[int]string{},
		progressBars:     map[int]*ProgressBarModel{},
		pendingCleanup:   map[int]bool{},
		FetchError:       "",
	}
	return m
}

func (m *ListModel) Init() tea.Cmd {
	return nil
}

// retryFetchMsg triggers a re-fetch when a processed repo hasn't appeared in the index yet
type retryFetchMsg struct{}

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
		case "up":
			service.NavigateUp(&m.FocusedIdx, &m.viewport, cardHeight)
		case "down":
			service.NavigateDown(&m.FocusedIdx, len(m.GHRepos), &m.viewport, cardHeight)
		}

	case tea.WindowSizeMsg:
		m.viewport.SetWidth(m.ctx.MainWidth)
		m.viewport.SetHeight(m.ctx.MainHeight - 4)

		for _, pb := range m.progressBars {
			pb.Update(msg)
		}

	// trigger from user pressing r
	case IndexRepoRequestMsg:
		// prevents case where r is pressed before gh fetch completes
		if len(m.GHRepos) == 0 {
			return nil
		}
		if m.IndexedRepoMap[m.GHRepos[m.FocusedIdx].Name] != nil {
			return nil // prevents an already indexed repo from sending a new index repo request
		}

		idx := m.FocusedIdx
		repoName := fmt.Sprintf("%s/%s", m.ctx.Username, m.GHRepos[idx].Name)

		return IndexRepo(idx, m.ctx.Username, repoName)
	// response from api call
	case indexRepoMsg:
		pb := NewProgressBar(m.ctx, msg.idx, msg.repoName)
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

	case common.FetchIndexedRepoMsg:
		return m.cleanupProgressBars()

	case retryFetchMsg:
		return common.FetchIndexedRepo(m.ctx.Username)
	case doneProcessingMsg:
		m.pendingCleanup[msg.idx] = true
		m.ProcessingStatus[msg.idx] = "fetching indexed repo..."
		return common.FetchIndexedRepo(m.ctx.Username)

	case searchQueryMsg:
		if len(msg.filteredGHRepos) == 0 {
			m.GHRepos = msg.filteredGHRepos
			m.IndexedRepoMap = msg.filteredIndexedRepos
			m.noSearchResults = true
			return nil
		}

		m.GHRepos = msg.filteredGHRepos
		m.IndexedRepoMap = msg.filteredIndexedRepos
		m.noSearchResults = false

		return nil
	}

	return nil
}

var repoCardStyle = lipgloss.NewStyle().PaddingLeft(2).Padding(0, 1).Border(lipgloss.RoundedBorder())

func (m *ListModel) View() tea.View {
	var cards []string

	if len(m.GHRepos) <= 0 {
		text := "Loading your repos..."
		if m.noSearchResults {
			text = "No repositories match your search"
		}
		return tea.NewView(lipgloss.NewStyle().Padding(1, 2).Render(text))
	}

	if m.FetchError != "" {
		return tea.NewView(lipgloss.NewStyle().Padding(1, 2).Render("error fetching from github"))
	}

	for i, repo := range m.GHRepos {
		ir := m.IndexedRepoMap[repo.Name]
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

		cards = append(cards, repoCardStyle.Width(m.ctx.MainWidth).BorderForeground(borderColor).Render(content))
	}
	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, cards...))

	return tea.NewView(m.viewport.View())
}

func (m *ListModel) Reset() {
	m.FocusedIdx = 0
	m.viewport.GotoTop()
}

// This is the top row in each list card, shows things like name, index status, total deps, and last indexed time
func (m *ListModel) cardHeader(idx int, indexedRepo types.IndexedRepo, ghRepo api.RepoApiRes, textColor color.Color) string {
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

type indexRepoMsg struct {
	idx      int
	repoName string
}

type indexRepoErrMsg struct {
	idx int
	err error
}

func IndexRepo(idx int, username, repoName string) tea.Cmd {
	return func() tea.Msg {
		err := api.IndexRepo(username, repoName)
		if err != nil {
			return indexRepoErrMsg{idx: idx, err: err}
		}

		return indexRepoMsg{idx: idx, repoName: repoName}
	}
}

func (m *ListModel) cleanupProgressBars() tea.Cmd {
	for idx := range m.pendingCleanup {
		if idx < len(m.GHRepos) && m.IndexedRepoMap[m.GHRepos[idx].Name] != nil {
			delete(m.progressBars, idx)
			delete(m.pendingCleanup, idx)
		}
	}

	// retry fetching for any remaining pending
	if len(m.pendingCleanup) > 0 {
		return tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return retryFetchMsg{}
		})
	}

	return nil
}
