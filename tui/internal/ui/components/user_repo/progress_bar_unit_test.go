//go:build unit

package user_repo

import (
	"testing"
	"tui/internal/ui/context"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
)

func progressBarCtx() *context.App {
	return &context.App{}
}

func TestProgressBar_Init(t *testing.T) {
	tt := []struct {
		name          string
		width         int
		expectedWidth int
	}{
		{"sets window properly", 50, 44},
		{"doesnt set width to negative", 4, 0},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := progressBarCtx()
			ctx.MainWidth = tc.width

			m := NewProgressBar(ctx, 0, "repo-1")
			cmd := m.Init()

			assert.Equal(t, tc.expectedWidth, m.progress.Width())
			assert.NotNil(t, cmd)
		})
	}
}

func TestProgressBar_WindowSizeMsg(t *testing.T) {
	tt := []struct {
		name          string
		width         int
		expectedWidth int
	}{
		{"resizing window properly", 50, 44},
		{"doesnt resize width negative", 4, 0},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := progressBarCtx()
			ctx.MainWidth = tc.width

			m := NewProgressBar(ctx, 0, "repo-1")
			m.Update(tea.WindowSizeMsg{Width: 33, Height: 22})

			assert.Equal(t, tc.expectedWidth, m.progress.Width())
		})
	}
}

func TestProgressBar_TickMsg_SetsPercent(t *testing.T) {
	for status, wantPct := range statusProgress {
		t.Run(status, func(t *testing.T) {
			m := NewProgressBar(progressBarCtx(), 0, "repo-1")
			m.Update(tickMsg{status: status})
			assert.InDelta(t, wantPct, m.progress.Percent(), 0.001)
		})
	}
}

func TestProgressBar_TickMsg_EmitsDoneMsgAt100Pct(t *testing.T) {
	tt := []string{"processed", "skipped:no dependencies found"}

	for _, status := range tt {
		t.Run(status, func(t *testing.T) {
			m := NewProgressBar(progressBarCtx(), 2, "repo-1")
			cmd := m.Update(tickMsg{status: status})

			msgs := runBatch(cmd)
			assert.Contains(t, msgs, doneProcessingMsg{idx: 2, status: status})
		})
	}
}

func TestProgressBar_TickMsg_UnknownStatusDoesNotSetPercent(t *testing.T) {
	m := NewProgressBar(progressBarCtx(), 0, "repo-1")
	m.Update(tickMsg{status: "some-unknown-status"})
	assert.Equal(t, 0.0, m.progress.Percent())
}

func TestProgressBar_GetJobStatusErr_ReturnsNil(t *testing.T) {
	m := NewProgressBar(progressBarCtx(), 0, "repo-1")
	cmd := m.Update(getJobStatusErr{})
	assert.Nil(t, cmd)
}

func runBatch(cmd tea.Cmd) []tea.Msg {
	msg := cmd()
	batch, ok := msg.(tea.BatchMsg)
	if !ok {
		return []tea.Msg{msg}
	}
	var msgs []tea.Msg
	for _, c := range batch {
		if c != nil {
			msgs = append(msgs, c())
		}
	}
	return msgs
}
