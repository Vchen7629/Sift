package views

import (
	"fmt"
	"tui/internal/api"
	"tui/internal/types"
	"tui/internal/ui/common"
	"tui/internal/ui/components/user_repo"
	"tui/internal/ui/context"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type UserRepoModel struct {
	ctx              *context.App
	ActionBar        *user_repo.ActionBarModel
	SearchBar        *user_repo.SearchBarModel
	RepoList         *user_repo.ListModel
	Sidebar          *user_repo.Sidebar
	ghRepos          []api.RepoApiRes
	indexedRepoMap   map[string]*types.IndexedRepo
	focusedIdx       int
	isSidebarFocused bool
}

func NewUserRepo(ctx *context.App) *UserRepoModel {
	return &UserRepoModel{
		ctx:              ctx,
		ActionBar:        user_repo.NewActionBar(ctx),
		SearchBar:        user_repo.NewSearchBar(ctx, "Search Your Repositories..."),
		RepoList:         user_repo.NewUserRepoList(ctx),
		Sidebar:          user_repo.NewSidebar(ctx),
		ghRepos:          []api.RepoApiRes{},
		isSidebarFocused: false,
	}
}

func (m *UserRepoModel) Init() tea.Cmd {
	m.focusedIdx = 0
	m.isSidebarFocused = false
	m.RepoList.Reset()

	return tea.Batch(m.fetchRepoList, common.FetchIndexedRepo(m.ctx.Username))
}

func (m *UserRepoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.ToggleFocusMsg:
		if !m.SearchBar.IsSearching() {
			m.isSidebarFocused = !m.isSidebarFocused
		}
		return m, nil
	case githubRepoFetchedMsg:
		m.ghRepos = msg.repoList
		m.RepoList.GHRepos = msg.repoList
		m.RepoList.FocusedIdx = 0
		if len(msg.repoList) > 0 {
			focused := m.ghRepos[0]
			m.Sidebar.FocusedGHRepo = &focused
			m.Sidebar.FocusedIndexedRepo = m.indexedRepoMap[m.ghRepos[0].Name]
			m.populateIndexRepoStatus()
		}
		return m, nil

	case githubRepoFetchedErr:
		m.RepoList.FetchError = fmt.Sprintf("failed to fetch repos from github: %s", msg.Err.Error())
		return m, nil

	case common.FetchIndexedRepoMsg:
		m.indexedRepoMap = make(map[string]*types.IndexedRepo, len(msg.IndexedRepos))
		for i := range msg.IndexedRepos {
			m.indexedRepoMap[msg.IndexedRepos[i].Name] = &msg.IndexedRepos[i]
		}
		m.RepoList.IndexedRepoMap = m.indexedRepoMap
		m.populateIndexRepoStatus()
		m.ActionBar.IndexRepoApiDown = false

		if len(m.ghRepos) > 0 {
			m.Sidebar.FocusedIndexedRepo = m.indexedRepoMap[m.ghRepos[m.focusedIdx].Name]
		}
		return m, nil

	case common.FetchIndexedRepoErr:
		m.ActionBar.IndexRepoApiDown = true
		return m, nil
	}

	actionBarCmd := m.ActionBar.Update(msg)
	repoListCmd := m.RepoList.Update(msg, m.isSidebarFocused)
	searchBarCmd := m.SearchBar.Update(msg, m.isSidebarFocused)
	sidebarCmd := m.Sidebar.Update(msg, m.isSidebarFocused)

	if len(m.RepoList.GHRepos) > 0 {
		newIdx := m.RepoList.FocusedIdx
		if newIdx != m.focusedIdx {
			m.focusedIdx = newIdx
			m.Sidebar.ResetFocus()
		}
		focused := &m.ghRepos[m.focusedIdx]
		m.Sidebar.FocusedGHRepo = focused
		m.Sidebar.FocusedIndexedRepo = m.indexedRepoMap[m.ghRepos[m.focusedIdx].Name]
	}

	return m, tea.Batch(actionBarCmd, repoListCmd, searchBarCmd, sidebarCmd)
}

func (m *UserRepoModel) View() tea.View {
	if m.ctx.WindowWidth == 0 {
		return tea.NewView("")
	}

	divider := common.VerticalDivider(m.ctx.MainHeight)

	repoListContent := lipgloss.JoinVertical(lipgloss.Top, m.SearchBar.View(), m.RepoList.View().Content)
	content := lipgloss.JoinHorizontal(lipgloss.Left, repoListContent, divider, m.Sidebar.View().Content)

	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, m.ActionBar.View(m.isSidebarFocused, m.Sidebar.FocusedIndexedRepo != nil).Content, content))
}

type githubRepoFetchedMsg struct{ repoList []api.RepoApiRes }
type githubRepoFetchedErr struct{ Err error }

// fetches user's repositories from github
func (m *UserRepoModel) fetchRepoList() tea.Msg {
	repos, err := m.ctx.GithubApiClient.GithubUserRepositories()
	if err != nil {
		return githubRepoFetchedErr{Err: err}
	}

	return githubRepoFetchedMsg{repoList: repos}
}

func (m *UserRepoModel) populateIndexRepoStatus() {
	for i, ghRepo := range m.RepoList.GHRepos {
		if m.indexedRepoMap[ghRepo.Name] != nil {
			m.RepoList.ProcessingStatus[i] = "Indexed"
		}
	}
}
