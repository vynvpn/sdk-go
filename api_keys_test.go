package vynvpn

import (
	"context"
	"testing"
)

func TestAPIKeys_List(t *testing.T) {
	client, req, srv := newTestServer(t, 200, []map[string]any{
		{"id": "00000000-0000-0000-0000-000000000001", "name": "Test Key"},
	}, WithToken("jwt"))
	defer srv.Close()

	keys, err := client.APIKeys.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("got %d keys, want 1", len(keys))
	}
	if keys[0].Name != "Test Key" {
		t.Errorf("Name = %q", keys[0].Name)
	}
	if req.Path != "/api/keys" {
		t.Errorf("path = %q", req.Path)
	}
}

func TestAPIKeys_Revoke(t *testing.T) {
	client, req, srv := newTestServer(t, 200, nil, WithToken("jwt"))
	defer srv.Close()

	err := client.APIKeys.Revoke(context.Background(), "key-id")
	if err != nil {
		t.Fatalf("Revoke failed: %v", err)
	}
	if req.Method != "DELETE" || req.Path != "/api/keys/key-id" {
		t.Errorf("request = %s %s", req.Method, req.Path)
	}
}
