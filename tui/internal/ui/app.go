package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tui/internal/ui/components/footer"
	"tui/internal/ui/context"
	"tui/internal/ui/views"
)

type model struct {
	ctx    *context.App
	pages  map[context.Page]tea.Model
	footer *footer.BaseModel
}

func New() (model, error) {
	ctx, err := context.NewApp()
	if err != nil {
		return model{}, err
	}
	return model{
		ctx: ctx,
		pages: map[context.Page]tea.Model{
			context.QueryPage:     views.NewRagQuery(ctx),
			context.UserReposPage: views.NewUserRepo(ctx),
		},
		footer: footer.NewFooterBaseModel(ctx),
	}, nil
}

func (m model) Init() tea.Cmd {
	footerCmd := m.footer.Init()
	pageCmd := m.pages[m.ctx.CurrentPage].Init()

	return tea.Batch(pageCmd, footerCmd)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ctx.WindowWidth = msg.Width
		m.ctx.WindowHeight = msg.Height
		m.ctx.MainWidth = msg.Width - 55
		m.ctx.SidebarWidth = m.ctx.WindowWidth - m.ctx.MainWidth - 2
		m.ctx.MainHeight = msg.Height - 3
		m.footer.SetSize(msg.Width-2, 1)
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	prevPage := m.ctx.CurrentPage
	sbCmd := m.footer.Update(msg)

	var initCmd tea.Cmd
	if m.ctx.CurrentPage != prevPage {
		initCmd = m.pages[m.ctx.CurrentPage].Init()
	}

	updatedModel, cmd := m.pages[m.ctx.CurrentPage].Update(msg)
	m.pages[m.ctx.CurrentPage] = updatedModel

	return m, tea.Batch(sbCmd, cmd, initCmd)
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