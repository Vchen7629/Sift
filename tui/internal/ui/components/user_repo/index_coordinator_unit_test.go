//go:build unit

package user_repo

import (
	"testing"
	"tui/internal/ui/context"

	"github.com/stretchr/testify/assert"
)

func indexCoordCtx() *context.App {
	return &context.App{SessionToken: "old-token"}
}

func TestIndexCoordinator_IndexRepoMsg(t *testing.T) {
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
		})
	}
}

func TestIndexCoordinator_TickMsg(t *testing.T) {
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
			c.progressBars[0] = NewProgressBar(ctx, 0, "repo")

			c.Update(tickMsg{idx: 0, status: "processing:created_job", newSessionToken: tc.newSessToken, isReauthed: tc.isReauthed})

			assert.Equal(t, tc.expectedToken, ctx.SessionToken)
		})
	}
}
