//go:build unit

package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteIndexedRepo_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)

		cookie, err := r.Cookie("JSESSIONID")
		require.NoError(t, err, "expected JSESSIONID cookie")
		assert.Equal(t, "mytoken", cookie.Value)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	orig := indexedRepoBaseUrl
	indexedRepoBaseUrl = ts.URL

	t.Cleanup(func() { indexedRepoBaseUrl = orig })

	err := DeleteIndexedRepo("mytoken", "some-repo")

	require.NoError(t, err)
}

func TestGetAllIndexedRepos_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		cookie, err := r.Cookie("JSESSIONID")
		require.NoError(t, err, "expected JSESSIONID cookie")
		assert.Equal(t, "mytoken", cookie.Value)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"repoName":"user/some-repo","lastIndexed":"2024-01-01","totalDependencies":1,"dependencies":[]}]`))
	}))
	defer ts.Close()

	orig := indexedRepoBaseUrl
	indexedRepoBaseUrl = ts.URL

	t.Cleanup(func() { indexedRepoBaseUrl = orig })

	_, err := GetAllIndexedRepos("mytoken")

	require.NoError(t, err)
}

func TestDeleteIndexedRepos_InvalidResponse(t *testing.T) {
	tt := []struct {
		name        string
		statusCode  int
		expectedErr error
	}{
		{"403 Forbidden returns ErrUnauthorized", 403, ErrUnauthorized},
		{"Unexpected status returns error", 500, errors.New("unexpected error sending req: 500")},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "DELETE", r.Method)

				cookie, err := r.Cookie("JSESSIONID")
				require.NoError(t, err, "expected JSESSIONID cookie")
				assert.Equal(t, "mytoken", cookie.Value)

				w.WriteHeader(tc.statusCode)
			}))
			defer ts.Close()

			orig := indexedRepoBaseUrl
			indexedRepoBaseUrl = ts.URL
			t.Cleanup(func() { indexedRepoBaseUrl = orig })

			err := DeleteIndexedRepo("mytoken", "some-repo")

			require.Equal(t, err, tc.expectedErr)
		})
	}
}

func TestGetAllIndexedRepos_InvalidResponse(t *testing.T) {
	tt := []struct {
		name        string
		statusCode  int
		expectedErr error
	}{
		{"403 Forbidden returns ErrUnauthorized", 403, ErrUnauthorized},
		{"Unexpected status returns error", 500, errors.New("unexpected error sending req: 500")},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)

				cookie, err := r.Cookie("JSESSIONID")
				require.NoError(t, err, "expected JSESSIONID cookie")
				assert.Equal(t, "mytoken", cookie.Value)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.statusCode)
				w.Write([]byte(`[{"repoName":"user/some-repo","lastIndexed":"2024-01-01","totalDependencies":1,"dependencies":[]}]`))
			}))
			defer ts.Close()

			orig := indexedRepoBaseUrl
			indexedRepoBaseUrl = ts.URL
			t.Cleanup(func() { indexedRepoBaseUrl = orig })

			_, err := GetAllIndexedRepos("mytoken")

			require.Equal(t, err, tc.expectedErr)
		})
	}
}

func TestDeleteIndexedRepo_NetworkError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(0)
	}))

	orig := indexedRepoBaseUrl
	indexedRepoBaseUrl = ts.URL
	t.Cleanup(func() { indexedRepoBaseUrl = orig })

	ts.Close()

	_, err := GetAllIndexedRepos("mytoken")

	require.Error(t, err)
}

func TestGetAllIndexedRepos_NetworkError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(0)
	}))

	orig := indexedRepoBaseUrl
	indexedRepoBaseUrl = ts.URL
	t.Cleanup(func() { indexedRepoBaseUrl = orig })

	ts.Close()

	_, err := GetAllIndexedRepos("mytoken")

	require.Error(t, err)
}
