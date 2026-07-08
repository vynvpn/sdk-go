package vynvpn

import (
	"context"

	"github.com/google/uuid"
)

// SubscriptionsService handles subscription management.
type SubscriptionsService struct {
	client *Client
}

// List returns all active subscriptions for the authenticated user.
// Routes to /api/v1/subscriptions (API key) or /api/subscriptions (JWT).
func (s *SubscriptionsService) List(ctx context.Context) ([]Subscription, error) {
	if s.client.http.IsAPIKeyAuth() {
		var resp struct {
			Data []Subscription `json:"data"`
		}
		if err := s.client.http.get(ctx, "/api/v1/subscriptions", nil, &resp); err != nil {
			return nil, err
		}
		return resp.Data, nil
	}
	var subs []Subscription
	if err := s.client.http.get(ctx, "/api/subscriptions", nil, &subs); err != nil {
		return nil, err
	}
	return subs, nil
}

// Get returns a single subscription by ID.
// Routes to /api/v1/subscriptions/{id} (API key) or /api/subscriptions/{id} (JWT).
func (s *SubscriptionsService) Get(ctx context.Context, id uuid.UUID) (*Subscription, error) {
	if s.client.http.IsAPIKeyAuth() {
		var resp struct {
			Data Subscription `json:"data"`
		}
		if err := s.client.http.get(ctx, "/api/v1/subscriptions/"+id.String(), nil, &resp); err != nil {
			return nil, err
		}
		return &resp.Data, nil
	}
	var sub Subscription
	if err := s.client.http.get(ctx, "/api/subscriptions/"+id.String(), nil, &sub); err != nil {
		return nil, err
	}
	return &sub, nil
}

// GetConfigs returns all provisioned VPN configs for a subscription token.
// GET /v2/config/{token}
func (s *SubscriptionsService) GetConfigs(ctx context.Context, token string) ([]Config, error) {
	var resp struct {
		Active  bool     `json:"active"`
		Configs []Config `json:"configs"`
	}
	if err := s.client.http.get(ctx, "/v2/config/"+token, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Configs, nil
}

// GetProfile returns subscription profile/usage for a v2 subscription token.
func (s *SubscriptionsService) GetProfile(ctx context.Context, token string) (*Profile, error) {
	var resp Profile
	if err := s.client.http.get(ctx, "/v2/profile/"+token, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// PurchaseRequest is the body for purchasing a subscription via JWT auth.
type PurchaseRequest struct {
	PlanID string `json:"plan_id"`
}

// Purchase initiates a subscription purchase (JWT-authenticated).
func (s *SubscriptionsService) Purchase(ctx context.Context, req *PurchaseRequest) (*Subscription, error) {
	var resp Subscription
	if err := s.client.http.post(ctx, "/api/subscriptions/purchase", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// EnsureFreeSub ensures the user has an active free-tier subscription.
func (s *SubscriptionsService) EnsureFreeSub(ctx context.Context) error {
	return s.client.http.post(ctx, "/api/ensure-free-sub", nil, nil)
}
