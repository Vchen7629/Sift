package user_repo

import (
	"errors"
	"fmt"
	"image/color"

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
	ctx             *context.App
	GHRepos         []api.RepoApiRes
	IndexedRepoMap  map[string]*types.IndexedRepo
	FocusedIdx      int
	FetchError      string
	viewport        viewport.Model
	IndexCoord      *indexCoordinator
	noSearchResults bool
}

func NewUserRepoList(ctx *context.App) *ListModel {
	return &ListModel{
		ctx:        ctx,
		GHRepos:    []api.RepoApiRes{},
		FetchError: "",
		IndexCoord: newIndexCoordinator(ctx),
	}
}

func (m *ListModel) Init() tea.Cmd {
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
		if _, ok := m.IndexCoord.progressBars[m.FocusedIdx]; ok {
			cardHeight = 5
		}

		switch msg.String() {
		case "up":
			service.NavigateUp(&m.FocusedIdx, &m.viewport, cardHeight)
		case "down":
			service.NavigateDown(&m.FocusedIdx, len(m.GHRepos), &m.viewport, cardHeight)
		case "enter":
			url := fmt.Sprintf("https://github.com/%s/%s", m.ctx.Username, m.GHRepos[m.FocusedIdx].Name)

			return common.OpenInBrowser(url)
		}
	case tea.WindowSizeMsg:
		m.viewport.SetWidth(m.ctx.MainWidth)
		m.viewport.SetHeight(m.ctx.MainHeight - 4)

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

		return IndexRepo(idx, m.ctx.SessionToken, repoName)

	case searchQueryMsg:
		m.GHRepos = msg.filteredGHRepos
		m.IndexedRepoMap = msg.filteredIndexedRepos
		m.noSearchResults = len(msg.filteredGHRepos) == 0
		return nil
	}

	return m.IndexCoord.Update(msg)
}

var repoCardStyle = lipgloss.NewStyle().PaddingLeft(2).Padding(0, 1).Border(lipgloss.RoundedBorder())

func (m *ListModel) View() tea.View {
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

	var cards []string
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
		if pb, ok := m.IndexCoord.progressBars[i]; ok {
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
	status, exists := m.IndexCoord.StatusFor(idx)
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
	idx                       int
	repoName, NewSessionToken string
	isReauthed                bool
}

type indexRepoErrMsg struct {
	idx int
	err error
}

func IndexRepo(idx int, sessionToken, repoName string) tea.Cmd {
	return func() tea.Msg {
		err := api.IndexRepo(sessionToken, repoName)

		if errors.Is(err, api.ErrUnauthorized) {
			ghToken := api.GithubPatToken()

			newSessionToken, err := api.NewSession(ghToken)
			if err != nil {
				return indexRepoErrMsg{idx: idx, err: err}
			}

			err = api.IndexRepo(newSessionToken, repoName)
			if err != nil {
				return indexRepoErrMsg{idx: idx, err: err}
			}

			return indexRepoMsg{idx: idx, repoName: repoName, NewSessionToken: newSessionToken, isReauthed: true}
		}
		if err != nil {
			return indexRepoErrMsg{idx: idx, err: err}
		}

		return indexRepoMsg{idx: idx, repoName: repoName, NewSessionToken: sessionToken, isReauthed: false}
	}
}
