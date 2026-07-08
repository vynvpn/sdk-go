package vynvpn

import "context"

// UsageService handles free-tier usage reporting and status.
type UsageService struct {
	client *Client
}

// GetStatus returns today's usage status and 7-day history.
// Requires JWT auth.
// GET /v2/usage
func (s *UsageService) GetStatus(ctx context.Context) (*UsageStatus, error) {
	var resp UsageStatus
	if err := s.client.http.get(ctx, "/v2/usage", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Report reports bytes consumed while connected.
// Returns updated usage status. Requires JWT auth.
func (s *UsageService) Report(ctx context.Context, req *UsageReport) (*UsageStatus, error) {
	var resp UsageStatus
	if err := s.client.http.post(ctx, "/v2/usage", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}


