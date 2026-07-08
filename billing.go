package vynvpn

import (
	"context"

	"github.com/google/uuid"
)

// BillingService handles Stripe billing operations.
// All methods require JWT auth.
type BillingService struct {
	client *Client
}

// CreateCheckoutSession creates a Stripe Checkout session for the given plan.
// Returns a URL to redirect the user to for payment.
func (s *BillingService) CreateCheckoutSession(ctx context.Context, planID uuid.UUID) (*CheckoutSession, error) {
	var resp CheckoutSession
	if err := s.client.http.post(ctx, "/api/billing/checkout-session",
		map[string]string{"plan_id": planID.String()}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreatePortalSession creates a Stripe Customer Portal session.
// Returns a URL the user can use to manage their payment method, invoices, or cancel.
func (s *BillingService) CreatePortalSession(ctx context.Context) (*PortalSession, error) {
	var resp PortalSession
	if err := s.client.http.post(ctx, "/api/billing/portal-session", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelSubscription requests cancellation at the end of the current billing period.
func (s *BillingService) CancelSubscription(ctx context.Context, subscriptionID uuid.UUID) error {
	return s.client.http.post(ctx,
		"/api/billing/subscriptions/"+subscriptionID.String()+"/cancel", nil, nil)
}

// ListInvoices returns the user's Stripe invoices.
func (s *BillingService) ListInvoices(ctx context.Context) ([]Invoice, error) {
	var invoices []Invoice
	if err := s.client.http.get(ctx, "/api/billing/invoices", nil, &invoices); err != nil {
		return nil, err
	}
	return invoices, nil
}

// ActivateTrial activates the one-time free trial plan.
func (s *BillingService) ActivateTrial(ctx context.Context) (*Subscription, error) {
	var sub Subscription
	if err := s.client.http.post(ctx, "/api/billing/trial", nil, &sub); err != nil {
		return nil, err
	}
	return &sub, nil
}

// GetTrialStatus returns the trial plan availability/status.
func (s *BillingService) GetTrialStatus(ctx context.Context) (*TrialStatus, error) {
	var status TrialStatus
	if err := s.client.http.get(ctx, "/api/billing/trial", nil, &status); err != nil {
		return nil, err
	}
	return &status, nil
}
