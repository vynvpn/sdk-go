package vynvpn

// testutil_test.go — shared test helpers used across all _test.go files.
//
// newTestServer spins up an httptest.Server that returns a fixed JSON body
// for every request. Tests can inspect the captured request to verify the
// SDK sent the right method, path, and headers.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// testRequest captures the last request received by the test server.
type testRequest struct {
	Method string
	Path   string
	Header http.Header
	Body   map[string]any
}

// newTestServer starts a test HTTP server that responds with statusCode and
// body (JSON-encoded). It records every incoming request in *testRequest.
// The returned Client is pre-configured to hit this server.
func newTestServer(t *testing.T, statusCode int, body any, opts ...Option) (*Client, *testRequest, *httptest.Server) {
	t.Helper()
	captured := &testRequest{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured.Method = r.Method
		captured.Path = r.URL.Path
		captured.Header = r.Header.Clone()
		if r.Body != nil {
			_ = json.NewDecoder(r.Body).Decode(&captured.Body)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(body)
	}))

	baseOpts := []Option{WithBaseURL(srv.URL)}
	client := New(append(baseOpts, opts...)...)
	return client, captured, srv
}

// newErrorServer returns a server that always responds with the given HTTP
// status and an {"error": msg} body.
func newErrorServer(t *testing.T, statusCode int, msg string, opts ...Option) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
	}))
	baseOpts := []Option{WithBaseURL(srv.URL)}
	client := New(append(baseOpts, opts...)...)
	return client, srv
}

// ptr returns a pointer to the given value — handy in test literals.
func ptr[T any](v T) *T { return &v }
