package vynvpn

import (
	"context"
	"testing"
)

func TestHTTP_APIKeyHeader(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{"data": []any{}}, WithAPIKey("vyn_test123"))
	defer srv.Close()

	_, _ = client.Nodes.List(context.Background())
	if got := req.Header.Get("X-Api-Key"); got != "vyn_test123" {
		t.Errorf("X-API-Key = %q, want %q", got, "vyn_test123")
	}
	if got := req.Header.Get("Authorization"); got != "" {
		t.Errorf("Authorization should be empty when API key is set, got %q", got)
	}
}

func TestHTTP_BearerHeader(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{"data": []any{}}, WithToken("jwt.here"))
	defer srv.Close()

	_, _ = client.Nodes.List(context.Background())
	if got := req.Header.Get("Authorization"); got != "Bearer jwt.here" {
		t.Errorf("Authorization = %q, want %q", got, "Bearer jwt.here")
	}
	if got := req.Header.Get("X-Api-Key"); got != "" {
		t.Errorf("X-API-Key should be empty when using JWT, got %q", got)
	}
}

func TestHTTP_UserAgent(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{"data": []any{}})
	defer srv.Close()

	_, _ = client.Nodes.List(context.Background())
	want := "vynvpn-sdk-go/" + Version
	if got := req.Header.Get("User-Agent"); got != want {
		t.Errorf("User-Agent = %q, want %q", got, want)
	}
}

func TestHTTP_ErrorParsing(t *testing.T) {
	client, srv := newErrorServer(t, 403, "account is suspended")
	defer srv.Close()

	_, err := client.Nodes.List(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 403 {
		t.Errorf("StatusCode = %d, want 403", apiErr.StatusCode)
	}
	if apiErr.Message != "account is suspended" {
		t.Errorf("Message = %q, want %q", apiErr.Message, "account is suspended")
	}
}
