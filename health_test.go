package vynvpn

import (
	"context"
	"testing"
)

func TestHealth_Get_JWT(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"user":            map[string]any{"id": "u1"},
		"subscription":    map[string]any{"active": true},
		"connections":     []any{},
		"usage":           map[string]any{"bytes_used_today": 1024},
		"nodes_available": 5,
		"alerts":          []string{"expiring soon"},
		"_hints":          map[string]any{"suggested_action": "renew"},
	}, WithToken("jwt"))
	defer srv.Close()

	resp, err := client.Health.Get(context.Background())
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if resp.NodesAvailable != 5 {
		t.Errorf("NodesAvailable = %d, want 5", resp.NodesAvailable)
	}
	if len(resp.Alerts) != 1 || resp.Alerts[0] != "expiring soon" {
		t.Errorf("Alerts = %v", resp.Alerts)
	}
	if resp.Hints == nil || resp.Hints.SuggestedAction != "renew" {
		t.Errorf("Hints.SuggestedAction = %v", resp.Hints)
	}
	if req.Path != "/v2/health" {
		t.Errorf("path = %q, want /v2/health", req.Path)
	}
}

func TestHealth_Get_APIKey(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"nodes_available": 3,
		"alerts":          []string{},
	}, WithAPIKey("vyn_key"))
	defer srv.Close()

	resp, err := client.Health.Get(context.Background())
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if resp.NodesAvailable != 3 {
		t.Errorf("NodesAvailable = %d, want 3", resp.NodesAvailable)
	}
	if req.Path != "/api/v1/health" {
		t.Errorf("path = %q, want /api/v1/health", req.Path)
	}
}
