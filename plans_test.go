package vynvpn

import (
	"context"
	"testing"
)

func TestPlans_List_JWT(t *testing.T) {
	client, req, srv := newTestServer(t, 200, []map[string]any{
		{"id": "00000000-0000-0000-0000-000000000001", "name": "Pro", "price_usd": 9.99},
	}, WithToken("jwt"))
	defer srv.Close()

	plans, err := client.Plans.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(plans) != 1 {
		t.Fatalf("got %d plans, want 1", len(plans))
	}
	if req.Path != "/api/plans" {
		t.Errorf("path = %q, want /api/plans", req.Path)
	}
}

func TestPlans_List_APIKey(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"data": []map[string]any{
			{"id": "00000000-0000-0000-0000-000000000001", "name": "Pro"},
		},
	}, WithAPIKey("vyn_key"))
	defer srv.Close()

	plans, err := client.Plans.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(plans) != 1 {
		t.Fatalf("got %d plans, want 1", len(plans))
	}
	if req.Path != "/api/v1/plans" {
		t.Errorf("path = %q, want /api/v1/plans", req.Path)
	}
}
