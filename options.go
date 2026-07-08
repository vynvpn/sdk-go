package vynvpn

import (
	"net/http"
	"time"
)

// config holds internal SDK configuration.
type config struct {
	baseURL    string
	apiKey     string
	token      string
	timeout    time.Duration
	httpClient *http.Client
}

// Option configures the SDK client.
type Option func(*config)

// WithBaseURL sets the API base URL. Defaults to DefaultBaseURL.
func WithBaseURL(url string) Option {
	return func(c *config) {
		c.baseURL = url
	}
}

// WithAPIKey sets an API key for authentication.
// API keys are long-lived and scoped — ideal for CLI tools, agents, and
// server-side integrations. Takes precedence over WithToken when both are set.
// Generate keys from the dashboard at /api/keys.
func WithAPIKey(key string) Option {
	return func(c *config) {
		c.apiKey = key
	}
}

// WithToken sets an existing JWT token for authentication.
// Useful when restoring a session from persistent storage.
// Prefer WithAPIKey for non-interactive / agent use cases.
func WithToken(token string) Option {
	return func(c *config) {
		c.token = token
	}
}

// WithTimeout sets the HTTP client timeout. Defaults to 30s.
func WithTimeout(d time.Duration) Option {
	return func(c *config) {
		c.timeout = d
		c.httpClient.Timeout = d
	}
}

// WithHTTPClient sets a custom http.Client for all requests.
// Useful for proxies, custom TLS, or testing.
func WithHTTPClient(client *http.Client) Option {
	return func(c *config) {
		c.httpClient = client
	}
}
