//go:build unit

package common_test

import (
	"strings"
	"testing"
	"tui/internal/ui/common"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

func TestSpaceBetween(t *testing.T) {
	tt := []struct {
		name                        string
		total, left, right, padding int
		wantWidth                   int
	}{
		{"normal spacing", 100, 30, 20, 5, 45},
		{"zero padding", 100, 40, 40, 0, 20},
		{"overflow clamps to zero", 10, 8, 8, 0, 0},
		{"exactly zero", 10, 5, 5, 0, 0},
		{"negative clamps to zero", 10, 6, 6, 0, 0},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := common.SpaceBetween(tc.total, tc.left, tc.right, tc.padding)
			assert.Equal(t, tc.wantWidth, lipgloss.Width(result))
		})
	}
}

func TestVerticalDivider(t *testing.T) {
	tt := []struct {
		name      string
		height    int
		wantCount int
	}{
		{"height 1", 1, 1},
		{"height 3", 3, 3},
		{"height 5", 5, 5},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := common.VerticalDivider(tc.height)
			assert.Equal(t, tc.wantCount, strings.Count(result, "│"))
		})
	}
}
