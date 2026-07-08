package vynvpn

import "context"

// PlansService handles plan listing.
type PlansService struct {
	client *Client
}

// List returns all active plans available for purchase.
// Routes to /api/v1/plans (API key) or /api/plans (JWT).
func (s *PlansService) List(ctx context.Context) ([]Plan, error) {
	if s.client.http.IsAPIKeyAuth() {
		var resp struct {
			Data []Plan `json:"data"`
		}
		if err := s.client.http.get(ctx, "/api/v1/plans", nil, &resp); err != nil {
			return nil, err
		}
		return resp.Data, nil
	}
	var plans []Plan
	if err := s.client.http.get(ctx, "/api/plans", nil, &plans); err != nil {
		return nil, err
	}
	return plans, nil
}
