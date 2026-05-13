//go:build unit

package rag_query

import (
	"testing"
	"tui/internal/types"
	"tui/internal/ui/common"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
)

func testCtx() *context.App {
	return &context.App{
		SelectedTheme: styles.Warm,
	}
}

var fakeIndexedRepos = []types.IndexedRepo{{Name: "Coco"}, {Name: "Cottom"}, {Name: "Maple"}}

func TestSidebarUpdate_SidebarFocusedGuard(t *testing.T) {
	tt := []struct {
		name             string
		isSidebarFocused bool
		isSearching	     bool
		key              rune
		wantIdx          int
	}{
		{"down blocked when not focused", false, false, tea.KeyDown, 0},
		{"up blocked when not focused", false, false, tea.KeyUp, 0},
		{"down blocked when searching", true, true, tea.KeyDown, 0},
		{"up blocked when searching", true, true, tea.KeyUp, 0},
		{"down works when focused and not searching", true, false, tea.KeyDown, 1},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := NewSidebar(testCtx())

			m.Update(common.FetchIndexedRepoMsg{IndexedRepos: fakeIndexedRepos}, false, tc.isSearching)
			m.Update(tea.KeyPressMsg{Code: tc.key}, tc.isSidebarFocused, tc.isSearching)

			assert.Equal(t, m.focusedIdx, tc.wantIdx)
		})
	}
}

func TestSidebarUpdate_Enter(t *testing.T) {
	tc := []struct {
		name    string
		repos   []types.IndexedRepo
		wantCmd bool
	}{
		{"no-op when repos empty", nil, false},
		{"returns cmd when repos present", fakeIndexedRepos, true},
	}

	for _, tc := range tc {
		t.Run(tc.name, func(t *testing.T) {
			m := NewSidebar(testCtx())
			m.Update(common.FetchIndexedRepoMsg{IndexedRepos: tc.repos}, false, false)
			cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter}, true, false)

			if tc.wantCmd {
				assert.NotNil(t, cmd)
			} else {
				assert.Nil(t, cmd)
			}
		})
	}
}
