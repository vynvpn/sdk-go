package vynvpn

import "context"

// NodesService handles VPN node/location listing.
type NodesService struct {
	client *Client
}

// List returns all available VPN server locations.
// Public endpoint — no auth required.
// GET /v2/nodes
func (s *NodesService) List(ctx context.Context) ([]Node, error) {
	var resp struct {
		Data []Node `json:"data"`
	}
	if err := s.client.http.get(ctx, "/v2/nodes", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// ListWithAPIKey returns nodes via the API-key path which includes node IDs.
// Only needed when you specifically require node IDs for advanced use.
// GET /api/v1/nodes  (API key required)
func (s *NodesService) ListWithAPIKey(ctx context.Context) ([]Node, error) {
	var resp struct {
		Data []Node `json:"data"`
	}
	if err := s.client.http.get(ctx, "/api/v1/nodes", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}
