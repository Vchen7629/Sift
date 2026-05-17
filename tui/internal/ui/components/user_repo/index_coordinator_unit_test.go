//go:build unit

package user_repo

import (
	"fmt"
	"testing"
	"tui/internal/api"
	"tui/internal/types"
	"tui/internal/ui/context"

	"charm.land/bubbles/v2/spinner"
	"github.com/stretchr/testify/assert"
)

func indexCoordCtx() *context.App {
	return &context.App{SessionToken: "old-token"}
}

func TestIndexCoordinatorUpdate_IndexRepoMsg(t *testing.T) {
	tt := []struct {
		name          string
		isReauthed    bool
		newSessToken  string
		expectedToken string
	}{
		{"updates the token on reauth", true, "new-token", "new-token"},
		{"does not overwrite token when not reauthed", false, "old-token-2", "old-token"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := indexCoordCtx()
			c := newIndexCoordinator(ctx)
			c.ctx.SessionToken = "old-token"

			c.Update(indexRepoMsg{idx: 0, repoName: "repo", NewSessionToken: tc.newSessToken, isReauthed: tc.isReauthed})

			assert.Equal(t, tc.expectedToken, ctx.SessionToken)
			assert.Contains(t, c.progressBars, 0)
		})
	}
}

func TestIndexCoordinatorUpdate_IndexRepoErrMsg(t *testing.T) {
	c := newIndexCoordinator(indexCoordCtx())
	c.Update(indexRepoErrMsg{idx: 0, err: fmt.Errorf("something broke")})

	assert.Equal(t, "error: something broke", c.statuses[0])
}

func TestIndexCoordinatorUpdate_TickMsg(t *testing.T) {
	tt := []struct {
		name           string
		isReauthed     bool
		newSessToken   string
		expectedToken  string
		status         string
		wantStatusText string
	}{
		{"updates the token on reauth", true, "new-token", "new-token", "", "new index job request"},
		{"does not overwrite token when not reauthed", false, "old-token-2", "old-token", "", "new index job request"},
		{"non-empty status is used as-is", true, "new-token", "new-token", "processing:created_job", "processing:created_job"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := indexCoordCtx()
			c := newIndexCoordinator(ctx)
			c.ctx.SessionToken = "old-token"
			c.progressBars[0] = NewProgressBar(ctx, 0, "repo")

			c.Update(tickMsg{idx: 0, status: tc.status, newSessionToken: tc.newSessToken, isReauthed: tc.isReauthed})

			assert.Equal(t, tc.expectedToken, ctx.SessionToken)
			assert.Equal(t, tc.wantStatusText, c.statuses[0])
		})
	}
}

func TestIndexCoordinatorUpdate_DoneProcessingMsg(t *testing.T) {
	tt := []struct {
		name                 string
		status               string
		wantProgressBarGone  bool
		wantInPendingCleanup bool
		wantCmd              bool
		wantStatusContains   string
	}{
		{
			name:                 "skipped:no dependencies found deletes progress bar",
			status:               "skipped:no dependencies found",
			wantProgressBarGone:  true,
			wantInPendingCleanup: false,
			wantCmd:              false,
			wantStatusContains:   "skipped:no dependencies found",
		},
		{
			name:                 "other status adds to pendingCleanup",
			status:               "completed",
			wantProgressBarGone:  false,
			wantInPendingCleanup: true,
			wantCmd:              true,
			wantStatusContains:   "fetching indexed repo",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := indexCoordCtx()
			c := newIndexCoordinator(ctx)
			c.progressBars[0] = NewProgressBar(ctx, 0, "repo")

			cmd := c.Update(doneProcessingMsg{idx: 0, status: tc.status})

			if tc.wantProgressBarGone {
				assert.NotContains(t, c.progressBars, 0)
			} else {
				assert.Contains(t, c.progressBars, 0)
			}

			if tc.wantInPendingCleanup {
				assert.True(t, c.pendingCleanup[0])
			} else {
				assert.NotContains(t, c.pendingCleanup, 0)
			}

			assert.Contains(t, c.statuses[0], tc.wantStatusContains)

			if tc.wantCmd {
				assert.NotNil(t, cmd)
			} else {
				assert.Nil(t, cmd)
			}
		})
	}
}

func TestIndexCoordinatorUpdate_SpinnerTickMsg_UpdatesPendingCleanupStatuses(t *testing.T) {
	ctx := indexCoordCtx()
	c := newIndexCoordinator(ctx)
	c.pendingCleanup[0] = true
	c.pendingCleanup[2] = true

	c.Update(spinner.TickMsg{})

	assert.Contains(t, c.statuses[0], "fetching indexed repo")
	assert.Contains(t, c.statuses[2], "fetching indexed repo")
}

func TestIndexCoordinatorUpdate_GetJobStatusErr(t *testing.T) {
	c := newIndexCoordinator(indexCoordCtx())
	c.Update(getJobStatusErr{idx: 1, err: "job not found"})
	assert.Equal(t, "job not found", c.statuses[1])
}

func TestIndexCoordinatorUpdate_RetryFetchMsg_ReturnsCmd(t *testing.T) {
	c := newIndexCoordinator(indexCoordCtx())
	cmd := c.Update(retryFetchMsg{})
	assert.NotNil(t, cmd)
}

func TestIndexCoordinator_CleanupProgressBars(t *testing.T) {
	ctx := indexCoordCtx()

	tt := []struct {
		name               string
		pendingCleanup     map[int]bool
		progressBars       map[int]*ProgressBarModel
		ghRepos            []api.RepoApiRes
		indexedRepoMap     map[string]*types.IndexedRepo
		wantProgressBars   []int
		wantPendingCleanup []int
		wantCmd            bool
	}{
		{
			name:               "empty pending cleanup returns nil",
			pendingCleanup:     map[int]bool{},
			progressBars:       map[int]*ProgressBarModel{},
			ghRepos:            []api.RepoApiRes{{Name: "repo-a"}},
			indexedRepoMap:     map[string]*types.IndexedRepo{"repo-a": {}},
			wantProgressBars:   []int{},
			wantPendingCleanup: []int{},
			wantCmd:            false,
		},
		{
			name:               "deletes progress bar and pending cleanup when conditions met",
			pendingCleanup:     map[int]bool{0: true},
			progressBars:       map[int]*ProgressBarModel{0: NewProgressBar(ctx, 0, "repo-a")},
			ghRepos:            []api.RepoApiRes{{Name: "repo-a"}},
			indexedRepoMap:     map[string]*types.IndexedRepo{"repo-a": {}},
			wantProgressBars:   []int{},
			wantPendingCleanup: []int{},
			wantCmd:            false,
		},
		{
			name:               "does not delete when idx is out of bounds of ghRepos",
			pendingCleanup:     map[int]bool{5: true},
			progressBars:       map[int]*ProgressBarModel{5: NewProgressBar(ctx, 5, "repo-a")},
			ghRepos:            []api.RepoApiRes{{Name: "repo-a"}},
			indexedRepoMap:     map[string]*types.IndexedRepo{"repo-a": {}},
			wantProgressBars:   []int{5},
			wantPendingCleanup: []int{5},
			wantCmd:            true,
		},
		{
			name:               "does not delete when repo not found in indexedRepoMap",
			pendingCleanup:     map[int]bool{0: true},
			progressBars:       map[int]*ProgressBarModel{0: NewProgressBar(ctx, 0, "repo-a")},
			ghRepos:            []api.RepoApiRes{{Name: "repo-a"}},
			indexedRepoMap:     map[string]*types.IndexedRepo{},
			wantProgressBars:   []int{0},
			wantPendingCleanup: []int{0},
			wantCmd:            true,
		},
		{
			name:           "deletes only repos present in indexedRepoMap, retries for remaining",
			pendingCleanup: map[int]bool{0: true, 1: true},
			progressBars: map[int]*ProgressBarModel{
				0: NewProgressBar(ctx, 0, "repo-a"),
				1: NewProgressBar(ctx, 1, "repo-b"),
			},
			ghRepos:            []api.RepoApiRes{{Name: "repo-a"}, {Name: "repo-b"}},
			indexedRepoMap:     map[string]*types.IndexedRepo{"repo-a": {}},
			wantProgressBars:   []int{1},
			wantPendingCleanup: []int{1},
			wantCmd:            true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := newIndexCoordinator(ctx)
			c.pendingCleanup = tc.pendingCleanup
			c.progressBars = tc.progressBars

			cmd := c.CleanupProgressBars(tc.ghRepos, tc.indexedRepoMap)

			assert.Len(t, c.progressBars, len(tc.wantProgressBars))
			assert.Len(t, c.pendingCleanup, len(tc.wantPendingCleanup))

			for _, idx := range tc.wantProgressBars {
				assert.Contains(t, c.progressBars, idx)
			}
			for _, idx := range tc.wantPendingCleanup {
				assert.Contains(t, c.pendingCleanup, idx)
			}

			if tc.wantCmd {
				assert.NotNil(t, cmd)
			} else {
				assert.Nil(t, cmd)
			}
		})
	}
}

func TestIndexCoordinator_StatusFor(t *testing.T) {
	tt := []struct {
		name       string
		statusMap  map[int]string
		idx        int
		wantStatus string
	}{
		{
			name:       "fetches status for nonexistant idx properly",
			statusMap:  map[int]string{},
			idx:        0,
			wantStatus: "",
		},
		{
			name:       "fetches status for idx properly",
			statusMap:  map[int]string{0: "old-1", 1: "old-2", 2: "old-3"},
			idx:        0,
			wantStatus: "old-1",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := newIndexCoordinator(indexCoordCtx())
			c.statuses = tc.statusMap
			status, _ := c.StatusFor(tc.idx)

			assert.Equal(t, tc.wantStatus, status)
		})
	}
}

func TestIndexCoordinator_SetStatus(t *testing.T) {
	tt := []struct {
		name       string
		statusMap  map[int]string
		idx        int
		newStatus  string
		wantStatus string
	}{
		{
			name:       "creates a new status with correct status properly",
			statusMap:  map[int]string{},
			idx:        0,
			newStatus:  "new",
			wantStatus: "new",
		},
		{
			name:       "setting status on idx updates correct status",
			statusMap:  map[int]string{0: "old-1", 1: "old-2", 2: "old-3"},
			idx:        0,
			newStatus:  "new",
			wantStatus: "new",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := newIndexCoordinator(indexCoordCtx())
			c.statuses = tc.statusMap
			c.SetStatus(tc.idx, tc.newStatus)

			assert.Equal(t, tc.wantStatus, c.statuses[tc.idx])
		})
	}
}
