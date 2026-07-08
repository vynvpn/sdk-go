package vynvpn

import "context"

// HealthService serves the health/state snapshot endpoint.
type HealthService struct {
	client *Client
}

// Get returns the full user state in one call:
// subscription, active connections, data usage, available nodes, alerts,
// and actionable hints for the next step.
//
// Uses /api/v1/health (API key) or /v2/health (JWT).
func (s *HealthService) Get(ctx context.Context) (*HealthResponse, error) {
	var resp HealthResponse
	path := "/v2/health"
	if s.client.http.IsAPIKeyAuth() {
		path = "/api/v1/health"
	}
	if err := s.client.http.get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
