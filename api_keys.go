package vynvpn

import (
	"context"
	"time"
)

// APIKeysService handles API key management.
// Requires JWT auth — keys are managed by humans from the dashboard or CLI.
type APIKeysService struct {
	client *Client
}

// List returns all API keys for the authenticated user.
func (s *APIKeysService) List(ctx context.Context) ([]APIKey, error) {
	var keys []APIKey
	if err := s.client.http.get(ctx, "/api/keys", nil, &keys); err != nil {
		return nil, err
	}
	return keys, nil
}

// Get returns a single API key by ID.
func (s *APIKeysService) Get(ctx context.Context, id string) (*APIKey, error) {
	var key APIKey
	if err := s.client.http.get(ctx, "/api/keys/"+id, nil, &key); err != nil {
		return nil, err
	}
	return &key, nil
}

// CreateKeyRequest is the body for creating a new API key.
type CreateKeyRequest struct {
	Name      string     `json:"name"`
	Scopes    []string   `json:"scopes"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// Create generates a new API key.
// The raw key in CreateKeyResponse.RawKey is only returned once — store it securely.
func (s *APIKeysService) Create(ctx context.Context, req *CreateKeyRequest) (*CreateKeyResponse, error) {
	var resp CreateKeyResponse
	if err := s.client.http.post(ctx, "/api/keys", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Revoke permanently disables an API key.
func (s *APIKeysService) Revoke(ctx context.Context, id string) error {
	return s.client.http.delete(ctx, "/api/keys/"+id, nil)
}

// ListScopes returns all valid scope strings with descriptions.
func (s *APIKeysService) ListScopes(ctx context.Context) ([]string, error) {
	var scopes []string
	if err := s.client.http.get(ctx, "/api/keys/scopes", nil, &scopes); err != nil {
		return nil, err
	}
	return scopes, nil
}

// GetKeyUsage returns usage stats for a specific key.
func (s *APIKeysService) GetKeyUsage(ctx context.Context, id string) (*KeyUsage, error) {
	var usage KeyUsage
	if err := s.client.http.get(ctx, "/api/keys/"+id+"/usage", nil, &usage); err != nil {
		return nil, err
	}
	return &usage, nil
}

// GetOverallUsage returns aggregate usage across all keys for the user.
func (s *APIKeysService) GetOverallUsage(ctx context.Context) (*KeyUsage, error) {
	var usage KeyUsage
	if err := s.client.http.get(ctx, "/api/keys/usage", nil, &usage); err != nil {
		return nil, err
	}
	return &usage, nil
}

// GetRequestLogs returns recent API request logs.
func (s *APIKeysService) GetRequestLogs(ctx context.Context) ([]RequestLog, error) {
	var logs []RequestLog
	if err := s.client.http.get(ctx, "/api/keys/logs", nil, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}
