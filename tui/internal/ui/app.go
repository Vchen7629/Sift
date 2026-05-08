package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tui/internal/ui/components/footer"
	"tui/internal/ui/components/rag_query"
	"tui/internal/ui/components/user_repo"
	"tui/internal/ui/context"
	"tui/internal/ui/views"
)

type model struct {
	ctx 		  *context.App
	pages		  map[context.Page]tea.Model
	footer 	     footer.BaseModel
}

// constructor to initialize the pages map
func New() model {
	ctx := context.NewApp()
	return model{ 
		ctx: ctx,
		pages: map[context.Page]tea.Model{
			context.QueryPage: 	   views.RagQueryModel{
				Ctx: ctx,
				SelectedRepo: "Sift",
				ActionBar: rag_query.NewActionBar(ctx),
				Searchbar: rag_query.NewRagQuerySearchBar(ctx),
				ResponseDisplay: rag_query.NewRagQueryResponse(ctx),
				Sidebar: rag_query.NewSidebar(ctx),
			},
			context.UserReposPage: &views.UserRepoModel{
				Ctx: ctx,
				SearchBar: user_repo.NewUserRepoSearchBar(ctx),
				RepoList: user_repo.NewUserRepoList(ctx),
				Sidebar: user_repo.NewSidebar(ctx),
			},
		},
		footer:  footer.BaseModel{
			Ctx: ctx,
			NavButtons: footer.NewNavButtons(ctx),
			ThemeSelector: footer.NewThemeSelector(ctx),
		},
	}
}

func (m model) Init() tea.Cmd {
	return m.pages[m.ctx.CurrentPage].Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ctx.WindowWidth = msg.Width
		m.ctx.WindowHeight = msg.Height
		m.ctx.MainWidth = msg.Width - 55
		m.ctx.SidebarWidth = m.ctx.WindowWidth - m.ctx.MainWidth - 2
		m.ctx.MainHeight = msg.Height - 3
		m.footer.SetSize(msg.Width - 2, 1)

		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
            return m, tea.Quit
		}
	}

	sbCmd := m.footer.Update(msg)

	updated, cmd := m.pages[m.ctx.CurrentPage].Update(msg)
	m.pages[m.ctx.CurrentPage] = updated

	return m, tea.Batch(sbCmd, cmd)
}

func (m model) View() tea.View {
	pageContent := m.pages[m.ctx.CurrentPage].View()
	footer := m.footer.View()

	footerHeight := lipgloss.Height(footer.Content)
	constrainedPage := lipgloss.NewStyle().
		Height(m.ctx.WindowHeight - footerHeight).
		MaxHeight(m.ctx.WindowHeight - footerHeight).
		Render(pageContent.Content)

	content := lipgloss.JoinVertical(lipgloss.Left, constrainedPage, footer.Content)

	padding := lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Width(m.ctx.WindowWidth)

	v := tea.NewView(padding.Render(content))
	v.AltScreen = true
	return v
}