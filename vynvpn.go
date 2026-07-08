// Package vynvpn provides a Go SDK for the VynVPN API.
//
// It supports two authentication modes:
//   - API Key: for programmatic integrations (CLI tools, third-party apps)
//   - JWT Token: for user-session based access (web/mobile apps)
//
// Quick start with API key:
//
//	client := vynvpn.New(
//	    vynvpn.WithAPIKey("vyn_..."),
//	    vynvpn.WithBaseURL("https://api.vynvpn.com"),
//	)
//	nodes, err := client.Nodes.List(ctx)
//
// Quick start with email login:
//
//	client := vynvpn.New(vynvpn.WithBaseURL("https://api.vynvpn.com"))
//	auth, err := client.Auth.Login(ctx, &vynvpn.LoginRequest{
//	    Email:    "user@example.com",
//	    Password: "secret",
//	})
//	// client automatically stores the JWT for subsequent requests
package vynvpn

import (
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the production API endpoint.
	DefaultBaseURL = "https://api.vynvpn.com"

	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second

	// Version is the SDK version.
	Version = "0.1.0"
)

// Client is the top-level VynVPN SDK client.
// Use New() to create one with the desired options.
type Client struct {
	// Service clients — each covers a logical API area.
	Auth          *AuthService
	Health        *HealthService
	Nodes         *NodesService
	Subscriptions *SubscriptionsService
	Connect       *ConnectService
	Plans         *PlansService
	Usage         *UsageService
	Payments      *PaymentsService
	Account       *AccountService
	Billing       *BillingService
	Tickets       *TicketsService
	Sessions      *SessionsService
	APIKeys       *APIKeysService
	TwoFA         *TwoFAService

	// internal transport
	http    *httpClient
	baseURL string
}

// New creates a new VynVPN SDK client with the given options.
func New(opts ...Option) *Client {
	cfg := &config{
		baseURL:    DefaultBaseURL,
		timeout:    DefaultTimeout,
		httpClient: &http.Client{Timeout: DefaultTimeout},
	}
	for _, opt := range opts {
		opt(cfg)
	}

	c := &Client{
		baseURL: cfg.baseURL,
	}

	c.http = &httpClient{
		client:    cfg.httpClient,
		baseURL:   cfg.baseURL,
		apiKey:    cfg.apiKey,
		token:     cfg.token,
		userAgent: "vynvpn-sdk-go/" + Version,
	}

	// Wire up service clients.
	c.Auth = &AuthService{client: c}
	c.Health = &HealthService{client: c}
	c.Nodes = &NodesService{client: c}
	c.Subscriptions = &SubscriptionsService{client: c}
	c.Connect = &ConnectService{client: c}
	c.Plans = &PlansService{client: c}
	c.Usage = &UsageService{client: c}
	c.Payments = &PaymentsService{client: c}
	c.Account = &AccountService{client: c}
	c.Billing = &BillingService{client: c}
	c.Tickets = &TicketsService{client: c}
	c.Sessions = &SessionsService{client: c}
	c.APIKeys = &APIKeysService{client: c}
	c.TwoFA = &TwoFAService{client: c}

	return c
}

// SetToken sets the JWT token for authenticated requests.
// This is called automatically after a successful login.
func (c *Client) SetToken(token string) {
	c.http.SetToken(token)
}

// Token returns the current JWT token, or empty string if not authenticated.
func (c *Client) Token() string {
	return c.http.token
}

// BaseURL returns the configured API base URL.
func (c *Client) BaseURL() string {
	return c.baseURL
}
