package views

import (
	"strings"
	"tui/internal/types"
	"tui/internal/service"
	"tui/internal/ui/common"
	"tui/internal/ui/components/user_repo"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type UserRepoModel struct {
	ctx 		      *context.App
	ActionBar         *user_repo.ActionBarModel
	SearchBar 	      *common.SearchBarModel
	RepoList	      *user_repo.ListModel
	Sidebar 	      *user_repo.Sidebar
	ghRepos		      []types.GHRepository
	indexedRepos 	  []types.IndexedRepo
	focusedIdx        int
	isSidebarFocused  bool
}

func NewUserRepo(ctx *context.App) *UserRepoModel {
	return &UserRepoModel{
		ctx: ctx,
		ActionBar: user_repo.NewActionBar(ctx),
		SearchBar: common.NewSearchBar(ctx, "Search Your Repositories..."),
		RepoList: user_repo.NewUserRepoList(ctx),
		Sidebar: user_repo.NewSidebar(ctx),
		ghRepos: []types.GHRepository{},
		isSidebarFocused: false,
	}
}

func (m *UserRepoModel) Init() tea.Cmd {
	return tea.Batch(m.fetchRepoList, common.FetchIndexedRepo(m.ctx.Username), m.RepoList.Init())
}

func (m *UserRepoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case user_repo.ToggleFocusMsg:
		if !m.SearchBar.IsSearching() {
			m.isSidebarFocused = !m.isSidebarFocused
		}
		return m, nil
	case githubRepoFetchedMsg:
		m.ghRepos = msg.repoList
		m.RepoList.GHRepos = msg.repoList
		m.RepoList.FocusedIdx = 0
		if len(msg.repoList) > 0 {
			m.Sidebar.FocusedGHRepo = &m.ghRepos[0]
			m.Sidebar.FocusedIndexedRepo = service.FindIndexedRepo(m.ghRepos[0].Name, m.indexedRepos)
			m.populateIndexRepoStatus()
		}
		return m, nil

	case common.FetchIndexedRepoMsg:
		m.indexedRepos = msg.IndexedRepos
		m.RepoList.IndexedRepos = msg.IndexedRepos
		m.populateIndexRepoStatus()

		if len(m.ghRepos) > 0 {
			m.Sidebar.FocusedIndexedRepo = service.FindIndexedRepo(m.ghRepos[m.focusedIdx].Name, m.indexedRepos)
		}
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
		m.Sidebar.FocusedGHRepo = &m.ghRepos[m.focusedIdx]
		m.Sidebar.FocusedIndexedRepo = service.FindIndexedRepo(m.ghRepos[m.focusedIdx].Name, m.indexedRepos)
	}

	return m, tea.Batch(actionBarCmd, repoListCmd, searchBarCmd, sidebarCmd)
}

func (m *UserRepoModel) View() tea.View {
	if m.ctx.WindowWidth == 0 {
		return tea.NewView("")
	}

	dividerLine := strings.Repeat("│\n", m.ctx.MainHeight - 1) + "│"
	divider := lipgloss.NewStyle().Foreground(styles.Divider).Render(dividerLine)

	repoListContent := lipgloss.JoinVertical(lipgloss.Top, m.SearchBar.View(), m.RepoList.View().Content)
	content := lipgloss.JoinHorizontal(lipgloss.Left, repoListContent, divider, m.Sidebar.View().Content)

	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, m.ActionBar.View(m.isSidebarFocused).Content, content))
}

type githubRepoFetchedMsg struct { repoList []types.GHRepository }

// fetches user's repositories from github
func (m *UserRepoModel) fetchRepoList() tea.Msg {
	repos, err := m.ctx.GithubApiClient.GithubUserRepositories()
	if err != nil {
		return err
	}

	return githubRepoFetchedMsg{ repoList: repos }
}

func (m *UserRepoModel) populateIndexRepoStatus() {
	for i, ghRepo := range m.RepoList.GHRepos {
		if service.FindIndexedRepo(ghRepo.Name, m.RepoList.IndexedRepos) != nil {
			m.RepoList.ProcessingStatus[i] = "Indexed"
		}
	}
}