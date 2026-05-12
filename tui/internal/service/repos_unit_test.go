//go:build unit

package service

import (
	"testing"
	"tui/internal/types"

	"github.com/stretchr/testify/assert"
)

var mockRepos = []types.IndexedRepo{
	{Name: "coco"},
	{Name: "cotton"},
	{Name: "maple"},
}

func TestFindIndexedRepo_Found(t *testing.T) {
	tc := []struct {
		name   string
		search string
	}{
		{"finds first item", "coco"},
		{"finds middle item", "cotton"},
		{"finds last item", "maple"},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			got := FindIndexedRepo(tt.search, mockRepos)

			assert.NotNil(t, got)
			assert.Equal(t, tt.search, got.Name)
		})
	}
}

func TestFindIndexedRepo_NotFound(t *testing.T) {
	tc := []struct {
		name   string
		search string
		list   []types.IndexedRepo
	}{
		{"not in list", "shit", mockRepos},
		{"empty list", "coco", []types.IndexedRepo{}},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			got := FindIndexedRepo(tt.search, tt.list)
			assert.Nil(t, got)
		})
	}
}
