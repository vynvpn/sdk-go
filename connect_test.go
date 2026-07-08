package vynvpn

import (
	"context"
	"testing"
)

func TestConnect_JWT(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "ready", "session_id": "sess-123", "config_link": "vless://...",
	}, WithToken("jwt"))
	defer srv.Close()

	resp, err := client.Connect.Connect(context.Background(), &ConnectRequest{
		LocationSlug: "de001",
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	if resp.Status != "ready" {
		t.Errorf("Status = %q, want %q", resp.Status, "ready")
	}
	if resp.ConfigLink != "vless://..." {
		t.Errorf("ConfigLink = %q", resp.ConfigLink)
	}
	if req.Path != "/v2/connect" {
		t.Errorf("path = %q, want /v2/connect", req.Path)
	}
}

func TestConnect_APIKey(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "provisioning", "session_id": "sess-456",
	}, WithAPIKey("vyn_key"))
	defer srv.Close()

	resp, err := client.Connect.Connect(context.Background(), &ConnectRequest{
		LocationSlug: "us001",
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	if resp.Status != "provisioning" {
		t.Errorf("Status = %q, want %q", resp.Status, "provisioning")
	}
	if req.Path != "/api/v1/connect" {
		t.Errorf("path = %q, want /api/v1/connect", req.Path)
	}
}

func TestConnect_Disconnect(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{"status": "disconnected"}, WithToken("jwt"))
	defer srv.Close()

	resp, err := client.Connect.Disconnect(context.Background(), &DisconnectRequest{
		LocationSlug: "de001",
	})
	if err != nil {
		t.Fatalf("Disconnect failed: %v", err)
	}
	if resp.Status != "disconnected" {
		t.Errorf("Status = %q, want %q", resp.Status, "disconnected")
	}
	if req.Path != "/v2/disconnect" {
		t.Errorf("path = %q, want /v2/disconnect", req.Path)
	}
}

func TestConnect_Status(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "ready", "config_link": "vless://ok",
	}, WithToken("jwt"))
	defer srv.Close()

	resp, err := client.Connect.Status(context.Background(), "sess-789")
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}
	if resp.Status != "ready" {
		t.Errorf("Status = %q, want %q", resp.Status, "ready")
	}
	if req.Path != "/v2/status/sess-789" {
		t.Errorf("path = %q", req.Path)
	}
}
