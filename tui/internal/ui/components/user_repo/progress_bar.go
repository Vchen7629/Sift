package user_repo

import (
	"time"
	"tui/internal/ui/context"

	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ProgressBarModel struct {
	ctx 	 *context.App
	progress progress.Model
}

func NewProgressBar(ctx *context.App) *ProgressBarModel {
	return &ProgressBarModel{
		ctx:      ctx,
		progress: progress.New(),
	}
}

func (m *ProgressBarModel) Init() tea.Cmd {
	return m.checkProgress()
}

func (m *ProgressBarModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tickMsg:
		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := m.progress.IncrPercent(0.05)
		return tea.Batch(m.checkProgress(), cmd)
	
	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		var cmd tea.Cmd
		m.progress, cmd = m.progress.Update(msg)
		return cmd
	}
	return nil
}

func (m *ProgressBarModel) View() tea.View {
	progress.WithColors(
		m.ctx.SelectedTheme.GradientDim, 
		m.ctx.SelectedTheme.GradientMid,
		m.ctx.SelectedTheme.GradientBright,
	)(&m.progress)

	m.progress.SetWidth(m.ctx.MainWidth - 6)
	paddingTop := lipgloss.NewStyle().PaddingTop(1)

	return tea.NewView(paddingTop.Render(m.progress.View()))
}

type tickMsg time.Time

func (m ProgressBarModel) checkProgress() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}