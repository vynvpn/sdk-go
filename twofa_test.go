package vynvpn

import (
	"context"
	"testing"
)

func TestTwoFA_Setup(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"secret": "JBSWY3DPEHPK3PXP", "qr_code": "base64png",
	}, WithToken("jwt"))
	defer srv.Close()

	resp, err := client.TwoFA.Setup(context.Background())
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}
	if resp.Secret != "JBSWY3DPEHPK3PXP" {
		t.Errorf("Secret = %q", resp.Secret)
	}
	if req.Path != "/api/auth/2fa/setup" {
		t.Errorf("path = %q", req.Path)
	}
}

func TestTwoFA_Enable(t *testing.T) {
	client, req, srv := newTestServer(t, 200, nil, WithToken("jwt"))
	defer srv.Close()

	err := client.TwoFA.Enable(context.Background(), "123456")
	if err != nil {
		t.Fatalf("Enable failed: %v", err)
	}
	if req.Method != "POST" || req.Path != "/api/auth/2fa/enable" {
		t.Errorf("request = %s %s", req.Method, req.Path)
	}
}

func TestTwoFA_Status(t *testing.T) {
	client, _, srv := newTestServer(t, 200, map[string]any{
		"enabled": true,
	}, WithToken("jwt"))
	defer srv.Close()

	status, err := client.TwoFA.Status(context.Background())
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}
	if !status.Enabled {
		t.Error("expected Enabled = true")
	}
}
