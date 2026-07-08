package vynvpn

import (
	"context"
	"testing"
)

func TestSessions_List(t *testing.T) {
	client, req, srv := newTestServer(t, 200, []map[string]any{
		{"id": "00000000-0000-0000-0000-000000000001", "ip": "1.2.3.4", "client": "web"},
	}, WithToken("jwt"))
	defer srv.Close()

	sessions, err := client.Sessions.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("got %d sessions, want 1", len(sessions))
	}
	if req.Path != "/api/sessions" {
		t.Errorf("path = %q", req.Path)
	}
}

func TestSessions_Revoke(t *testing.T) {
	client, req, srv := newTestServer(t, 200, nil, WithToken("jwt"))
	defer srv.Close()

	err := client.Sessions.Revoke(context.Background(), "sess-id")
	if err != nil {
		t.Fatalf("Revoke failed: %v", err)
	}
	if req.Method != "DELETE" || req.Path != "/api/sessions/sess-id" {
		t.Errorf("request = %s %s", req.Method, req.Path)
	}
}
