package vynvpn

import "context"

// TwoFAService handles two-factor authentication setup and management.
type TwoFAService struct {
	client *Client
}

// Setup initiates 2FA enrollment.
// Returns a TOTP secret and QR code — the user scans the QR code in their
// authenticator app, then calls Enable with a valid code to confirm.
// Requires JWT auth.
func (s *TwoFAService) Setup(ctx context.Context) (*TwoFASetupResponse, error) {
	var resp TwoFASetupResponse
	if err := s.client.http.post(ctx, "/api/auth/2fa/setup", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Enable confirms 2FA enrollment with a valid TOTP code from the authenticator app.
// Requires JWT auth.
func (s *TwoFAService) Enable(ctx context.Context, code string) error {
	return s.client.http.post(ctx, "/api/auth/2fa/enable", map[string]string{"code": code}, nil)
}

// Disable turns off 2FA for the account.
// Requires JWT auth.
func (s *TwoFAService) Disable(ctx context.Context, code string) error {
	return s.client.http.post(ctx, "/api/auth/2fa/disable", map[string]string{"code": code}, nil)
}

// Status returns whether 2FA is enabled for the authenticated user.
// Requires JWT auth.
func (s *TwoFAService) Status(ctx context.Context) (*TwoFAStatus, error) {
	var status TwoFAStatus
	if err := s.client.http.get(ctx, "/api/auth/2fa/status", nil, &status); err != nil {
		return nil, err
	}
	return &status, nil
}
