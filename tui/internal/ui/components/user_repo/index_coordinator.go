package user_repo

import (
	"fmt"
	"time"
	"tui/internal/api"
	"tui/internal/types"
	"tui/internal/ui/common"
	"tui/internal/ui/context"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// This owns all the state and logic for repo-indexing lifecycle
// progress bars, pending cleanup, spinner, and per-repo status text
type indexCoordinator struct {
	ctx            *context.App
	progressBars   map[int]*ProgressBarModel
	pendingCleanup map[int]bool
	spinner        spinner.Model
	statuses       map[int]string
}

func newIndexCoordinator(ctx *context.App) *indexCoordinator {
	s := spinner.New()
	s.Spinner = spinner.Points
	return &indexCoordinator{
		ctx:            ctx,
		progressBars:   map[int]*ProgressBarModel{},
		pendingCleanup: map[int]bool{},
		spinner:        s,
		statuses:       map[int]string{},
	}
}

func (c *indexCoordinator) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		for _, pb := range c.progressBars {
			pb.Update(msg)
		}

	case indexRepoMsg:
		c.ctx.SessionToken = msg.NewSessionToken
		pb := NewProgressBar(c.ctx, msg.idx, msg.repoName)
		c.progressBars[msg.idx] = pb

		return pb.Init()
	case indexRepoErrMsg:
		c.statuses[msg.idx] = fmt.Sprintf("error: %s", msg.err.Error())

	case tickMsg:
		c.ctx.SessionToken = msg.newSessionToken
		statusText := "new index job request"
		if msg.status != "" {
			statusText = msg.status
		}

		c.statuses[msg.idx] = statusText
		if pb, ok := c.progressBars[msg.idx]; ok {
			return pb.Update(msg)
		}

	case progress.FrameMsg:
		var cmds []tea.Cmd
		for _, pb := range c.progressBars {
			cmds = append(cmds, pb.Update(msg))
		}

		return tea.Batch(cmds...)

	case doneProcessingMsg:
		if msg.status == "skipped:no dependencies found" {
			delete(c.progressBars, msg.idx)
			c.statuses[msg.idx] = msg.status

			return nil
		}

		c.pendingCleanup[msg.idx] = true
		c.statuses[msg.idx] = "fetching indexed repo " + c.spinner.View()

		return tea.Batch(c.spinner.Tick, common.FetchIndexedRepo(c.ctx.SessionToken))

	case spinner.TickMsg:
		var cmd tea.Cmd
		c.spinner, cmd = c.spinner.Update(msg)
		c.spinner.Style = lipgloss.NewStyle().Foreground(c.ctx.SelectedTheme.AccentBright)

		for idx := range c.pendingCleanup {
			c.statuses[idx] = "fetching indexed repo " + c.spinner.View()
		}

		return cmd

	case retryFetchMsg:
		return common.FetchIndexedRepo(c.ctx.SessionToken)

	case getJobStatusErr:
		c.statuses[msg.idx] = msg.err
		return nil

	}

	return nil
}

// display status string for given repo index
func (c *indexCoordinator) StatusFor(idx int) (string, bool) {
	s, ok := c.statuses[idx]

	return s, ok
}

// used by external callers like populate indexrepostatus write a status
func (c *indexCoordinator) SetStatus(idx int, status string) {
	c.statuses[idx] = status
}

// triggers a re-fetch when a processed repo hasn't appeared in the index yet
type retryFetchMsg struct{}

// called after FetchIndexedRepoMsg for progress bar cleanup for repos
// that finished indexing
func (c *indexCoordinator) CleanupProgressBars(
	ghRepos []api.RepoApiRes, indexedRepoMap map[string]*types.IndexedRepo,
) tea.Cmd {
	for idx := range c.pendingCleanup {
		if idx < len(ghRepos) && indexedRepoMap[ghRepos[idx].Name] != nil {
			delete(c.progressBars, idx)
			delete(c.pendingCleanup, idx)
		}
	}

	if len(c.pendingCleanup) > 0 {
		return tea.Tick(time.Second, func(t time.Time) tea.Msg { return retryFetchMsg{} })
	}

	return nil
}
