package vynvpn

import (
	"context"
	"net/url"
	"time"

	"github.com/google/uuid"
)

// ConnectService handles VPN connect/disconnect operations.
type ConnectService struct {
	client *Client
}

// ConnectRequest is the body for connecting to a VPN location.
//
// JWT mode — all fields optional:
//   - Token omitted → subscription auto-detected from JWT
//   - LocationSlug omitted → API picks best server via geo-IP + load
//   - PreferredCountry hints the auto-picker
//
// API key mode — set SubscriptionID to target a specific subscription.
type ConnectRequest struct {
	Token            string     `json:"token,omitempty"`
	LocationSlug     string     `json:"location_slug,omitempty"`
	PreferredCountry string     `json:"preferred_country,omitempty"`
	SubscriptionID   *uuid.UUID `json:"-"` // API key mode only
}

// DisconnectRequest is the body for disconnecting from a VPN location.
type DisconnectRequest struct {
	Token          string     `json:"token,omitempty"`
	LocationSlug   string     `json:"location_slug"`
	SubscriptionID *uuid.UUID `json:"-"` // API key mode only
}

// Connect provisions a VPN config for a location on-demand.
//
// All fields on ConnectRequest are optional:
//   - Call with an empty &ConnectRequest{} to auto-detect subscription and
//     auto-pick the best server for the caller's location.
//   - Set PreferredCountry to hint toward a region (e.g. "Germany", "US").
//   - Set LocationSlug for an exact server choice.
//
// API key mode uses /api/v1/connect with SubscriptionID.
// JWT mode uses /v2/connect with optional Token.
//
// If status is "provisioning", poll with WaitForReady or Status.
func (s *ConnectService) Connect(ctx context.Context, req *ConnectRequest) (*ConnectResponse, error) {
	var resp ConnectResponse
	if s.client.http.IsAPIKeyAuth() {
		body := map[string]any{
			"location_slug": req.LocationSlug,
		}
		if req.SubscriptionID != nil {
			body["subscription_id"] = req.SubscriptionID.String()
		}
		if err := s.client.http.post(ctx, "/api/v1/connect", body, &resp); err != nil {
			return nil, err
		}
		return &resp, nil
	}
	if err := s.client.http.post(ctx, "/v2/connect", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Disconnect removes the VPN config for a location.
// The slot can be re-provisioned later with Connect.
func (s *ConnectService) Disconnect(ctx context.Context, req *DisconnectRequest) (*ConnectResponse, error) {
	var resp ConnectResponse
	if s.client.http.IsAPIKeyAuth() {
		body := map[string]any{
			"location_slug": req.LocationSlug,
		}
		if req.SubscriptionID != nil {
			body["subscription_id"] = req.SubscriptionID.String()
		}
		if err := s.client.http.post(ctx, "/api/v1/disconnect", body, &resp); err != nil {
			return nil, err
		}
		return &resp, nil
	}
	if err := s.client.http.post(ctx, "/v2/disconnect", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Status polls a provisioning session for completion.
// Pass the SessionID returned by Connect when status is "provisioning".
//
// GET /v2/status/{id}
func (s *ConnectService) Status(ctx context.Context, sessionID string) (*ConnectResponse, error) {
	var resp ConnectResponse
	if err := s.client.http.get(ctx, "/v2/status/"+sessionID, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListConnections returns all active connection sessions.
// token is optional when the client has a JWT set.
//
// GET /v2/connections
func (s *ConnectService) ListConnections(ctx context.Context, token string) ([]Connection, error) {
	var q url.Values
	if token != "" {
		q = url.Values{"token": {token}}
	}
	var resp struct {
		Data []Connection `json:"data"`
	}
	if err := s.client.http.get(ctx, "/v2/connections", q, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// WaitForReady polls Status until the connection is ready, failed, or the
// context/timeout is exceeded. Returns the final ConnectResponse.
func (s *ConnectService) WaitForReady(ctx context.Context, sessionID string, timeout time.Duration) (*ConnectResponse, error) {
	deadline := time.Now().Add(timeout)
	pollInterval := 2 * time.Second

	for {
		resp, err := s.Status(ctx, sessionID)
		if err != nil {
			return nil, err
		}
		switch resp.Status {
		case StatusReady, StatusFailed, StatusLimitExceeded:
			return resp, nil
		}
		if time.Now().After(deadline) {
			return resp, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(pollInterval):
		}
	}
}
