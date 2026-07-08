package vynvpn

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestSubscriptions_List_JWT(t *testing.T) {
	client, req, srv := newTestServer(t, 200, []map[string]any{
		{"id": "00000000-0000-0000-0000-000000000001", "active": true, "status": "active"},
	}, WithToken("jwt"))
	defer srv.Close()

	subs, err := client.Subscriptions.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(subs) != 1 {
		t.Fatalf("got %d subs, want 1", len(subs))
	}
	if req.Path != "/api/subscriptions" {
		t.Errorf("path = %q, want /api/subscriptions", req.Path)
	}
}

func TestSubscriptions_List_APIKey(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"data": []map[string]any{
			{"id": "00000000-0000-0000-0000-000000000001", "active": true},
		},
	}, WithAPIKey("vyn_key"))
	defer srv.Close()

	subs, err := client.Subscriptions.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(subs) != 1 {
		t.Fatalf("got %d subs, want 1", len(subs))
	}
	if req.Path != "/api/v1/subscriptions" {
		t.Errorf("path = %q, want /api/v1/subscriptions", req.Path)
	}
}

func TestSubscriptions_Get(t *testing.T) {
	id := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	client, req, srv := newTestServer(t, 200, map[string]any{
		"id": id.String(), "active": true,
	}, WithToken("jwt"))
	defer srv.Close()

	sub, err := client.Subscriptions.Get(context.Background(), id)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !sub.Active {
		t.Error("expected active subscription")
	}
	if req.Path != "/api/subscriptions/"+id.String() {
		t.Errorf("path = %q", req.Path)
	}
}
