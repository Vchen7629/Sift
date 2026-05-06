package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tui/internal/ui/components"
	"tui/internal/ui/styles"
	"tui/internal/ui/views"
)

type page int

const (
	authPage page = iota
	queryPage
	userReposPage
)

type model struct {
	width, height int
	currentPage   page
	pages 		  map[page]tea.Model
	statusBar 	  components.StatusBarModel
}

// constructor to initialize the pages map
func New() model {
	return model{
		currentPage: authPage,
		pages: map[page]tea.Model{
			authPage: 	   views.AuthModel{},
			queryPage:     views.QueryModel{},
			userReposPage: views.UserRepoModel{},
		},
		statusBar: components.StatusBarModel{},
	}
}

func (m model) Init() tea.Cmd {
	return m.pages[m.currentPage].Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.statusBar.SetSize(msg.Width, 2)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
            return m, tea.Quit
		}
	}

	updated, cmd := m.pages[m.currentPage].Update(msg)
	m.pages[m.currentPage] = updated

	return m, cmd
}

func (m model) View() tea.View {
	pageContent := m.pages[m.currentPage].View()

	main := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height - 1).
		Background(styles.Background).
		Render(pageContent.Content)

	statusBar := m.statusBar.View()

	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, main, statusBar.Content))
}