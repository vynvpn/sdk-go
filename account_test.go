package vynvpn

import (
	"context"
	"testing"
)

func TestAccount_GetMe_JWT(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"id": "00000000-0000-0000-0000-000000000001", "first_name": "Test",
	}, WithToken("jwt"))
	defer srv.Close()

	user, err := client.Account.GetMe(context.Background())
	if err != nil {
		t.Fatalf("GetMe failed: %v", err)
	}
	if user.FirstName != "Test" {
		t.Errorf("FirstName = %q, want %q", user.FirstName, "Test")
	}
	if req.Path != "/api/me" {
		t.Errorf("path = %q, want /api/me", req.Path)
	}
}

func TestAccount_GetMe_APIKey(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"data": map[string]any{"id": "00000000-0000-0000-0000-000000000001", "first_name": "API"},
	}, WithAPIKey("vyn_key"))
	defer srv.Close()

	user, err := client.Account.GetMe(context.Background())
	if err != nil {
		t.Fatalf("GetMe failed: %v", err)
	}
	if user.FirstName != "API" {
		t.Errorf("FirstName = %q, want %q", user.FirstName, "API")
	}
	if req.Path != "/api/v1/account" {
		t.Errorf("path = %q, want /api/v1/account", req.Path)
	}
}

func TestAccount_UpdateProfile(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"id": "00000000-0000-0000-0000-000000000001", "first_name": "Updated",
	}, WithToken("jwt"))
	defer srv.Close()

	user, err := client.Account.UpdateProfile(context.Background(), &UpdateProfileRequest{
		FirstName: ptr("Updated"),
	})
	if err != nil {
		t.Fatalf("UpdateProfile failed: %v", err)
	}
	if user.FirstName != "Updated" {
		t.Errorf("FirstName = %q, want %q", user.FirstName, "Updated")
	}
	if req.Method != "PATCH" || req.Path != "/api/profile" {
		t.Errorf("request = %s %s, want PATCH /api/profile", req.Method, req.Path)
	}
}
