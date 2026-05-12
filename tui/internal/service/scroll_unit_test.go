//go:build unit

package service

import (
	"testing"

	"charm.land/bubbles/v2/viewport"
	"github.com/stretchr/testify/assert"
)

func TestNavigateUp(t *testing.T) {
	tc := []struct {
		name    string
		start   int
		wantIdx int
	}{
		{"moves up normally", 2, 1},
		{"moves up from middle", 1, 0},
		{"does not go negative", 0, 0},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			vp := viewport.New()
			idx := tt.start
			NavigateUp(&idx, &vp, 1)
			assert.Equal(t, tt.wantIdx, idx, "should be this index after navigating")
		})
	}
}

func TestNavigateDown(t *testing.T) {
	tc := []struct {
		name    string
		start   int
		listLen int
		wantIdx int
	}{
		{"moves down normally", 0, 3, 1},
		{"moves down from middle", 1, 3, 2},
		{"does not move past end", 2, 3, 2},
		{"does not move on single item list", 0, 1, 0},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			vp := viewport.New()
			idx := tt.start
			NavigateDown(&idx, tt.listLen, &vp, 1)
			assert.Equal(t, tt.wantIdx, idx, "should be this index after navigating")
		})
	}
}
