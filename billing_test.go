package vynvpn

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestBilling_CreateCheckout(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"session_id": "cs_123", "url": "https://checkout.stripe.com/pay",
	}, WithToken("jwt"))
	defer srv.Close()

	planID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	sess, err := client.Billing.CreateCheckoutSession(context.Background(), planID)
	if err != nil {
		t.Fatalf("CreateCheckoutSession failed: %v", err)
	}
	if sess.URL != "https://checkout.stripe.com/pay" {
		t.Errorf("URL = %q", sess.URL)
	}
	if req.Method != "POST" || req.Path != "/api/billing/checkout-session" {
		t.Errorf("request = %s %s", req.Method, req.Path)
	}
}

func TestBilling_ActivateTrial(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"id": "00000000-0000-0000-0000-000000000001", "active": true,
	}, WithToken("jwt"))
	defer srv.Close()

	sub, err := client.Billing.ActivateTrial(context.Background())
	if err != nil {
		t.Fatalf("ActivateTrial failed: %v", err)
	}
	if !sub.Active {
		t.Error("expected active trial subscription")
	}
	if req.Path != "/api/billing/trial" {
		t.Errorf("path = %q, want /api/billing/trial", req.Path)
	}
}
