package vynvpn

import (
	"context"
	"testing"
)

func TestAuth_Login(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"token": "jwt.abc.def",
		"user":  map[string]any{"id": "00000000-0000-0000-0000-000000000001", "first_name": "Test"},
	})
	defer srv.Close()

	resp, err := client.Auth.Login(context.Background(), &LoginRequest{
		Email: "test@example.com", Password: "pass123",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if resp.Token != "jwt.abc.def" {
		t.Errorf("Token = %q, want %q", resp.Token, "jwt.abc.def")
	}
	if client.Token() != "jwt.abc.def" {
		t.Error("token should be auto-stored on client")
	}
	if req.Method != "POST" || req.Path != "/auth/login" {
		t.Errorf("request = %s %s, want POST /auth/login", req.Method, req.Path)
	}
}

func TestAuth_Login2FA(t *testing.T) {
	client, _, srv := newTestServer(t, 200, map[string]any{
		"requires_2fa": true,
		"login_token":  "pending.token",
	})
	defer srv.Close()

	resp, err := client.Auth.Login(context.Background(), &LoginRequest{
		Email: "test@example.com", Password: "pass123",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if !resp.Requires2FA {
		t.Error("expected Requires2FA = true")
	}
	if resp.LoginToken != "pending.token" {
		t.Errorf("LoginToken = %q, want %q", resp.LoginToken, "pending.token")
	}
	if client.Token() != "" {
		t.Error("token should NOT be stored when 2FA is required")
	}
}

func TestAuth_Register(t *testing.T) {
	client, req, srv := newTestServer(t, 201, map[string]any{
		"token": "new.jwt",
		"user":  map[string]any{"id": "00000000-0000-0000-0000-000000000002"},
	})
	defer srv.Close()

	resp, err := client.Auth.Register(context.Background(), &RegisterRequest{
		Email: "new@example.com", Password: "secret123",
	})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if resp.Token != "new.jwt" {
		t.Errorf("Token = %q, want %q", resp.Token, "new.jwt")
	}
	if client.Token() != "new.jwt" {
		t.Error("token should be auto-stored after register")
	}
	if req.Path != "/auth/register" {
		t.Errorf("path = %q, want /auth/register", req.Path)
	}
}

func TestAuth_Logout(t *testing.T) {
	client, _, srv := newTestServer(t, 200, nil, WithToken("existing.jwt"))
	defer srv.Close()

	_ = client.Auth.Logout(context.Background())
	if client.Token() != "" {
		t.Error("token should be cleared after logout")
	}
}

func TestAuth_Refresh(t *testing.T) {
	client, _, srv := newTestServer(t, 200, map[string]any{"token": "refreshed.jwt"}, WithToken("old.jwt"))
	defer srv.Close()

	resp, err := client.Auth.Refresh(context.Background())
	if err != nil {
		t.Fatalf("Refresh failed: %v", err)
	}
	if resp.Token != "refreshed.jwt" {
		t.Errorf("Token = %q, want %q", resp.Token, "refreshed.jwt")
	}
	if client.Token() != "refreshed.jwt" {
		t.Error("token should be updated after refresh")
	}
}
