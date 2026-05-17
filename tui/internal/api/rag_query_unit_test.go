//go:build unit

package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)

		cookie, err := r.Cookie("JSESSIONID")
		require.NoError(t, err, "expected JSESSIONID cookie")
		assert.Equal(t, "mytoken", cookie.Value)

		var body SearchReq
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "username/some-repo", body.RepoName)
		assert.Equal(t, "search query", body.SearchQuery)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"repoName":"username/some-repo","numSources":0,"issues":[],"summary":""}`))
	}))
	defer ts.Close()

	orig := searchBaseUrl
	searchBaseUrl = ts.URL

	t.Cleanup(func() { searchBaseUrl = orig })

	_, err := Search("mytoken", "username", "some-repo", "search query")

	require.NoError(t, err)
}

func TestSearch_InvalidResponse(t *testing.T) {
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
				assert.Equal(t, "POST", r.Method)

				cookie, err := r.Cookie("JSESSIONID")
				require.NoError(t, err, "expected JSESSIONID cookie")
				assert.Equal(t, "mytoken", cookie.Value)

				var body SearchReq
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				assert.Equal(t, "username/some-repo", body.RepoName)
				assert.Equal(t, "search query", body.SearchQuery)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.statusCode)
				w.Write([]byte(`{"repoName":"username/some-repo","numSources":0,"issues":[],"summary":""}`))
			}))
			defer ts.Close()

			orig := searchBaseUrl
			searchBaseUrl = ts.URL

			t.Cleanup(func() { searchBaseUrl = orig })

			_, err := Search("mytoken", "username", "some-repo", "search query")

			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSearch_NetworkError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(0)
	}))

	orig := searchBaseUrl
	searchBaseUrl = ts.URL
	t.Cleanup(func() { searchBaseUrl = orig })

	ts.Close()

	_, err := Search("mytoken", "username", "some-repo", "search query")

	require.Error(t, err)
}
