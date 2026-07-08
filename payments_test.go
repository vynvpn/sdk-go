package vynvpn

import (
	"context"
	"testing"
)

func TestPayments_List_JWT(t *testing.T) {
	client, req, srv := newTestServer(t, 200, []map[string]any{
		{"id": "00000000-0000-0000-0000-000000000001", "amount_usd": 9.99, "status": "confirmed"},
	}, WithToken("jwt"))
	defer srv.Close()

	payments, err := client.Payments.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(payments) != 1 {
		t.Fatalf("got %d payments, want 1", len(payments))
	}
	if req.Path != "/api/payments" {
		t.Errorf("path = %q, want /api/payments", req.Path)
	}
}

func TestPayments_List_APIKey(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"data": []map[string]any{
			{"id": "00000000-0000-0000-0000-000000000001", "amount_usd": 5.00},
		},
	}, WithAPIKey("vyn_key"))
	defer srv.Close()

	payments, err := client.Payments.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(payments) != 1 {
		t.Fatalf("got %d payments, want 1", len(payments))
	}
	if req.Path != "/api/v1/payments" {
		t.Errorf("path = %q, want /api/v1/payments", req.Path)
	}
}
