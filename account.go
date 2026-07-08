package vynvpn

import "context"

// AccountService handles user account/profile endpoints.
type AccountService struct {
	client *Client
}

// GetMe returns the authenticated user's profile.
// Routes to /api/v1/account (API key) or /api/me (JWT).
func (s *AccountService) GetMe(ctx context.Context) (*User, error) {
	if s.client.http.IsAPIKeyAuth() {
		var resp struct {
			Data User `json:"data"`
		}
		if err := s.client.http.get(ctx, "/api/v1/account", nil, &resp); err != nil {
			return nil, err
		}
		return &resp.Data, nil
	}
	var user User
	if err := s.client.http.get(ctx, "/api/me", nil, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateProfileRequest contains fields to update on the user's profile.
type UpdateProfileRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Language  *string `json:"language,omitempty"`
	Country   *string `json:"country,omitempty"`
}

// UpdateProfile updates the authenticated user's profile.
// Requires JWT auth.
func (s *AccountService) UpdateProfile(ctx context.Context, req *UpdateProfileRequest) (*User, error) {
	var user User
	if err := s.client.http.patch(ctx, "/api/profile", req, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// RequestAccountDeletion initiates the account deletion process.
// A confirmation code will be sent to the user's email.
func (s *AccountService) RequestAccountDeletion(ctx context.Context) error {
	return s.client.http.post(ctx, "/api/account/delete/request", nil, nil)
}

// ConfirmAccountDeletion confirms account deletion with the received code.
func (s *AccountService) ConfirmAccountDeletion(ctx context.Context, code string) error {
	return s.client.http.post(ctx, "/api/account/delete/confirm",
		map[string]string{"code": code}, nil)
}
