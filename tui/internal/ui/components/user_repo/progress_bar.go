package user_repo

import (
	"time"
	"tui/internal/api"
	"tui/internal/ui/context"

	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ProgressBarModel struct {
	ctx      *context.App
	progress progress.Model
	idx      int
	repoName string
}

func NewProgressBar(ctx *context.App, idx int, repoName string) *ProgressBarModel {
	return &ProgressBarModel{
		ctx:      ctx,
		progress: progress.New(),
		idx:      idx,
		repoName: repoName,
	}
}

func (m *ProgressBarModel) Init() tea.Cmd {
	m.progress.SetWidth(max(0, m.ctx.MainWidth-6))
	return m.checkProgress()
}

var statusProgress = map[string]float64{
	"processing:created_job":                   0.05,
	"processing:fetched_repo":                  0.075,
	"processing:fetched_dependency_list":       0.10,
	"processing:fetched_all_issues_changelogs": 0.80,
	"processing:inserted_all_issues":           0.85,
	"processing:inserted_all_changelogs":       0.90,
	"processing:inserted_indexed_repo":         0.95,
	"processed":                                1.0,
	"skipped:no dependencies found":            1.0,
}

type doneProcessingMsg struct{ idx int }

func (m *ProgressBarModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.SetWidth(max(0, m.ctx.MainWidth-6))

	case tickMsg:
		var cmd tea.Cmd
		if pct, ok := statusProgress[msg.status]; ok {
			cmd = m.progress.SetPercent(pct)
			if pct == 1.0 {
				return tea.Batch(cmd, func() tea.Msg { return doneProcessingMsg{idx: m.idx} })
			}
		}
		// this is to stop polling after its done processing
		if msg.status == "processed" || msg.status == "skipped:no dependencies found" {
			return cmd
		}
		return tea.Batch(m.checkProgress(), cmd)

	case getJobStatusErr:
		return nil

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

	paddingTop := lipgloss.NewStyle().PaddingTop(1)

	return tea.NewView(paddingTop.Render(m.progress.View()))
}

type tickMsg struct {
	idx    int
	t      time.Time
	status string
}

type getJobStatusErr struct{ Err error }

func (m *ProgressBarModel) checkProgress() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		status, err := api.GetJobStatus(m.ctx.Username, m.repoName)
		if err != nil {
			return getJobStatusErr{Err: err}
		}

		return tickMsg{t: t, idx: m.idx, status: status}
	})
}
