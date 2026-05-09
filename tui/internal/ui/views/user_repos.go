package views

import (
	"strings"
	"tui/internal/api"
	"tui/internal/types"
	"tui/internal/ui/components/user_repo"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type UserRepoModel struct {
	ctx 		      *context.App
	ActionBar         *user_repo.ActionBarModel
	SearchBar 	      *user_repo.SearchBarModel
	RepoList	      *user_repo.ListModel
	Sidebar 	      *user_repo.Sidebar
	repos			  []types.Repository
	focusedIdx        int
	isSidebarFocused  bool
}

func NewUserRepo(ctx *context.App) *UserRepoModel {
	return &UserRepoModel{
		ctx: ctx,
		ActionBar: user_repo.NewActionBar(ctx),
		SearchBar: user_repo.NewUserRepoSearchBar(ctx),
		RepoList: user_repo.NewUserRepoList(ctx),
		Sidebar: user_repo.NewSidebar(ctx),
		repos: []types.Repository{},
		isSidebarFocused: false,
	}
}

// called by update, populates the struct data once data is fetched from gh api
func (m *UserRepoModel) SetRepos(repos []types.Repository) {
	byId := make(map[int]types.Repository, len(dummyData))
	for _, data := range dummyData {
		byId[data.GithubId] = data
	}

	for i, repo := range repos {
		data, ok := byId[repo.GithubId]
		if ok {
			repos[i].Status = data.Status
			repos[i].LastIndexed = data.LastIndexed
			repos[i].TotalDependencies = data.TotalDependencies
			repos[i].Dependencies = data.Dependencies
		}
	}

	m.repos = repos
	m.RepoList.FetchedRepos = repos
	m.RepoList.FocusedIdx = 0
	if len(repos) > 0 {
		m.Sidebar.FocusedRepo = &m.repos[0]
	}
}

func (m *UserRepoModel) Init() tea.Cmd {
	return m.fetchRepoList
}

func (m *UserRepoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case user_repo.ToggleFocusMsg:
		if !m.SearchBar.IsSearching() {
			m.isSidebarFocused = !m.isSidebarFocused
		}
		return m, nil
	case githubRepoFetchedMsg:
		repos := make([]types.Repository, len(msg.repoList))
		for i, r := range msg.repoList {
			repos[i] = types.Repository{Name: r.Name, Description: r.Description, LastUpdated: r.LastCommit}
		}
		m.repos = repos
		m.SetRepos(m.repos)
		return m, nil
	}

	actionBarCmd := m.ActionBar.Update(msg)
	repoListCmd := m.RepoList.Update(msg, m.isSidebarFocused)
	searchBarCmd := m.SearchBar.Update(msg, m.isSidebarFocused)
	sidebarCmd := m.Sidebar.Update(msg, m.isSidebarFocused)

	if len(m.RepoList.FetchedRepos) > 0 {
		newIdx := m.RepoList.FocusedIdx
		if newIdx != m.focusedIdx {
			m.focusedIdx = newIdx
			m.Sidebar.ResetFocus()
		}
		m.Sidebar.FocusedRepo = &m.repos[m.focusedIdx]
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

type githubRepoFetchedMsg struct { repoList []api.RepoList }

// fetches user's repositories from github
func (m *UserRepoModel) fetchRepoList() tea.Msg {
	repos, err := m.ctx.GithubApiClient.GithubUserRepositories()
	if err != nil {
		return err
	}

	return githubRepoFetchedMsg{ repoList: repos }
}