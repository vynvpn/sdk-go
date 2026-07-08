package vynvpn

import "context"

// SessionsService handles active login session management.
type SessionsService struct {
	client *Client
}

// List returns all active sessions for the authenticated user.
// Requires JWT auth.
func (s *SessionsService) List(ctx context.Context) ([]UserSession, error) {
	var sessions []UserSession
	if err := s.client.http.get(ctx, "/api/sessions", nil, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

// Revoke terminates a specific session by ID.
// Requires JWT auth.
func (s *SessionsService) Revoke(ctx context.Context, sessionID string) error {
	return s.client.http.delete(ctx, "/api/sessions/"+sessionID, nil)
}

// RevokeOthers revokes all sessions except the current one.
// Requires JWT auth.
func (s *SessionsService) RevokeOthers(ctx context.Context) error {
	return s.client.http.post(ctx, "/api/sessions/revoke-others", nil, nil)
}
