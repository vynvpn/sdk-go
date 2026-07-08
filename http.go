package vynvpn

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// httpClient wraps net/http with auth handling and JSON serialization.
type httpClient struct {
	client    *http.Client
	baseURL   string
	apiKey    string
	token     string
	userAgent string
	mu        sync.RWMutex
}

// IsAPIKeyAuth reports whether the client is using API key authentication.
func (h *httpClient) IsAPIKeyAuth() bool { return h.apiKey != "" }

// SetToken stores a JWT token (thread-safe).
func (h *httpClient) SetToken(token string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.token = token
}

// request builds and executes an HTTP request.
func (h *httpClient) request(ctx context.Context, method, path string, body any, query url.Values) (*http.Response, error) {
	u := h.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("vynvpn: marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("vynvpn: create request: %w", err)
	}

	req.Header.Set("User-Agent", h.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	// Auth: prefer API key, fall back to JWT Bearer.
	if h.apiKey != "" {
		req.Header.Set("X-API-Key", h.apiKey)
	} else {
		h.mu.RLock()
		token := h.token
		h.mu.RUnlock()
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vynvpn: %s %s: %w", method, path, err)
	}
	return resp, nil
}

// do executes a request and decodes the JSON response into dest.
// If dest is nil, the body is discarded (useful for 204 No Content).
func (h *httpClient) do(ctx context.Context, method, path string, body any, query url.Values, dest any) error {
	resp, err := h.request(ctx, method, path, body, query)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read body for error messages.
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("vynvpn: read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return parseAPIError(resp.StatusCode, respBody)
	}

	if dest != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, dest); err != nil {
			return fmt.Errorf("vynvpn: decode response: %w (body: %s)", err, truncate(string(respBody), 200))
		}
	}
	return nil
}

// get is a convenience for GET requests.
func (h *httpClient) get(ctx context.Context, path string, query url.Values, dest any) error {
	return h.do(ctx, http.MethodGet, path, nil, query, dest)
}

// post is a convenience for POST requests.
func (h *httpClient) post(ctx context.Context, path string, body any, dest any) error {
	return h.do(ctx, http.MethodPost, path, body, nil, dest)
}

// patch is a convenience for PATCH requests.
func (h *httpClient) patch(ctx context.Context, path string, body any, dest any) error {
	return h.do(ctx, http.MethodPatch, path, body, nil, dest)
}

// put is a convenience for PUT requests.
func (h *httpClient) put(ctx context.Context, path string, body any, dest any) error {
	return h.do(ctx, http.MethodPut, path, body, nil, dest)
}

// delete is a convenience for DELETE requests.
func (h *httpClient) delete(ctx context.Context, path string, dest any) error {
	return h.do(ctx, http.MethodDelete, path, nil, nil, dest)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

// parseAPIError attempts to decode the API error response.
func parseAPIError(status int, body []byte) *APIError {
	apiErr := &APIError{
		StatusCode: status,
		RawBody:    string(body),
	}

	// Try parsing {"error": "message"} shape.
	var errResp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error != "" {
		apiErr.Message = errResp.Error
		return apiErr
	}

	// Fallback: use status text + raw body.
	apiErr.Message = strings.TrimSpace(string(body))
	if apiErr.Message == "" {
		apiErr.Message = http.StatusText(status)
	}
	return apiErr
}
