package footer

import (
	"fmt"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type BaseModel struct {
	height, width int
	ctx           *context.App
	NavButtons    *NavButtonsModel
	ThemeSelector *ThemeSelectorModel
}

func NewFooterBaseModel(ctx *context.App) *BaseModel {
	return &BaseModel{
		ctx:           ctx,
		NavButtons:    NewNavButtons(ctx),
		ThemeSelector: NewThemeSelector(ctx),
	}
}

func (m *BaseModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *BaseModel) Init() tea.Cmd {
	return nil
}

func (m *BaseModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "3":
			m.ctx.ThemeSelectorOpen = !m.ctx.ThemeSelectorOpen
			return nil
		}
	}

	return tea.Batch(m.NavButtons.Update(msg), m.ThemeSelector.Update(msg))
}

func (m *BaseModel) View() tea.View {
	title := "Sift · "
	if m.ctx.Username != "" {
		title = fmt.Sprintf("Sift · @%s", m.ctx.Username)
	}
	titleText := lipgloss.NewStyle().
		PaddingLeft(2).PaddingRight(1).
		Background(styles.Footer).Foreground(styles.TextPrimary).
		Render(title)

	background := lipgloss.NewStyle().Background(styles.Footer).Width(m.width)

	var content string
	if m.ctx.ThemeSelectorOpen {
		content = background.Render(lipgloss.JoinHorizontal(
			lipgloss.Left, titleText, m.NavButtons.View(), m.themeBtns(), m.ThemeSelector.View().Content,
		))
	} else {
		content = background.Render(lipgloss.JoinHorizontal(lipgloss.Left, titleText, m.NavButtons.View(), m.themeBtns()))
	}

	return tea.NewView(content)
}

func (m *BaseModel) themeBtns() string {
	themeLabel := "[3] theme"
	if m.ctx.ThemeSelectorOpen {
		themeLabel = lipgloss.NewStyle().
			PaddingRight(2).
			Background(styles.Footer).Foreground(m.ctx.SelectedTheme.AccentBright).Bold(true).
			Render(themeLabel)
	}

	return lipgloss.NewStyle().
		PaddingLeft(1).PaddingRight(1).
		Background(styles.Footer).
		Render(themeLabel)
}
