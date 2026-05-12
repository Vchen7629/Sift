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
	tt := []struct {
		name   string
		search string
	}{
		{"finds first item", "coco"},
		{"finds middle item", "cotton"},
		{"finds last item", "maple"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := FindIndexedRepo(tc.search, mockRepos)

			assert.NotNil(t, got)
			assert.Equal(t, tc.search, got.Name)
		})
	}
}

func TestFindIndexedRepo_NotFound(t *testing.T) {
	tt := []struct {
		name   string
		search string
		list   []types.IndexedRepo
	}{
		{"not in list", "shit", mockRepos},
		{"empty list", "coco", []types.IndexedRepo{}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := FindIndexedRepo(tc.search, tc.list)
			assert.Nil(t, got)
		})
	}
}
