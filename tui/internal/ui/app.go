package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tui/internal/ui/components"
	"tui/internal/ui/context"
	"tui/internal/ui/views"
)

type model struct {
	ctx 		  *context.App
	pages		  map[context.Page]tea.Model
	statusBar 	  components.StatusBarModel
}

// constructor to initialize the pages map
func New() model {
	ctx := context.NewApp()
	return model{ 
		ctx: ctx,
		pages: map[context.Page]tea.Model{
			context.AuthPage:	   views.AuthModel{Ctx: ctx},
			context.QueryPage: 	   views.QueryModel{Ctx: ctx},
			context.UserReposPage: &views.UserRepoModel{
				Ctx: ctx,
				RepoList: &components.UserRepoListModel{
					Ctx: ctx,
					FetchedRepos: []components.UserRepo{
						{Name: "react", Status: "indexed", LastIndexed: "1746", TotalLibs: "42"},
						{Name: "next.js", Status: "indexed", LastIndexed: "1745", TotalLibs: "18"},
						{Name: "tailwindcss", Status: "pending", LastIndexed: "0", TotalLibs: "5"},
						{Name: "react", Status: "indexed", LastIndexed: "1746", TotalLibs: "420"},
						{Name: "next.js", Status: "indexed", LastIndexed: "1745", TotalLibs: "18"},
						{Name: "tailwindcss", Status: "pending", LastIndexed: "0", TotalLibs: "5"},
						{Name: "react", Status: "indexed", LastIndexed: "1746", TotalLibs: "42"},
						{Name: "next.js", Status: "indexed", LastIndexed: "1745", TotalLibs: "18"},
						{Name: "tailwindcss", Status: "pending", LastIndexed: "0", TotalLibs: "5"},
						{Name: "react", Status: "indexed", LastIndexed: "1746", TotalLibs: "42"},
						{Name: "next.js", Status: "indexed", LastIndexed: "1745", TotalLibs: "18"},
						{Name: "tailwindcss", Status: "pending", LastIndexed: "0", TotalLibs: "5"},
						{Name: "react", Status: "indexed", LastIndexed: "1746", TotalLibs: "42"},
						{Name: "next.js", Status: "indexed", LastIndexed: "1745", TotalLibs: "18"},
						{Name: "tailwindcss", Status: "pending", LastIndexed: "0", TotalLibs: "5"},
					},
				},
			},
		},
		statusBar:  components.StatusBarModel{Ctx: ctx},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ctx.Width = msg.Width
		m.ctx.Height = msg.Height
		m.ctx.ViewPortWidth = msg.Width - 55
		m.ctx.ViewPortHeight = msg.Height - 3
		m.statusBar.SetSize(msg.Width - 2, 1)

		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
            return m, tea.Quit
		}
	}

	sbCmd := m.statusBar.Update(msg)

	updated, cmd := m.pages[m.ctx.CurrentPage].Update(msg)
	m.pages[m.ctx.CurrentPage] = updated

	return m, tea.Batch(sbCmd, cmd)
}

func (m model) View() tea.View {
	pageContent := m.pages[m.ctx.CurrentPage].View()
	statusBar := m.statusBar.View()

	content := lipgloss.JoinVertical(lipgloss.Left, pageContent.Content, statusBar.Content)

	padding := lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Width(m.ctx.Width)

	v := tea.NewView(padding.Render(content))
	v.AltScreen = true
	return v
}