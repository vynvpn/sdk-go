package vynvpn

import (
	"time"

	"github.com/google/uuid"
)

// ── User ─────────────────────────────────────────────────────────────────────

// User represents a VynVPN user account.
type User struct {
	ID              uuid.UUID  `json:"id"`
	Email           *string    `json:"email,omitempty"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	TelegramID      int64      `json:"telegram_id"`
	Username        string     `json:"username"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	Language        string     `json:"language"`
	Country         string     `json:"country"`
	AcceptedToS     bool       `json:"accepted_tos"`
	CreatedAt       time.Time  `json:"created_at"`
	TrialUsed       bool       `json:"trial_used"`
	IsReseller      bool       `json:"is_reseller"`
	IsBanned        bool       `json:"is_banned"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

// ── Auth responses ────────────────────────────────────────────────────────────

// LoginResponse is returned on successful login or 2FA challenge.
type LoginResponse struct {
	Token       string `json:"token,omitempty"`
	User        *User  `json:"user,omitempty"`
	Requires2FA bool   `json:"requires_2fa,omitempty"`
	LoginToken  string `json:"login_token,omitempty"` // present when Requires2FA is true
}

// RegisterResponse is returned on successful registration.
type RegisterResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// RefreshResponse is returned on token refresh.
type RefreshResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user,omitempty"`
}

// OAuthSessionResponse is returned when creating an OAuth polling session.
type OAuthSessionResponse struct {
	SessionID string `json:"session_id"`
	URL       string `json:"url,omitempty"`
}

// OAuthSessionStatus is the polling response for OAuth sessions.
type OAuthSessionStatus struct {
	Status string `json:"status"`
	Token  string `json:"token,omitempty"`
	User   *User  `json:"user,omitempty"`
}

// ── Plan ─────────────────────────────────────────────────────────────────────

// Plan represents a VPN subscription plan.
type Plan struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Features       []string  `json:"features"`
	DataLimitBytes int64     `json:"data_limit_bytes"`
	PriceUSD       float64   `json:"price_usd"`
	Currency       string    `json:"currency"`
	DurationDays   int       `json:"duration_days"`
	Active         bool      `json:"active"`
	IsTrial        bool      `json:"is_trial"`
	CreatedAt      time.Time `json:"created_at"`
}

// DataLimitGB returns the data limit in gigabytes.
func (p *Plan) DataLimitGB() float64 {
	return float64(p.DataLimitBytes) / 1_073_741_824
}

// ── Node ─────────────────────────────────────────────────────────────────────

// Node represents a VPN server location.
type Node struct {
	LocationSlug string `json:"location_slug"`
	Label        string `json:"label"`
	NodeID       string `json:"node_id,omitempty"`
	Country      string `json:"country"`
	Address      string `json:"address,omitempty"`
	Port         int    `json:"port,omitempty"`
	Available    bool   `json:"available"`
}

// ── Subscription ─────────────────────────────────────────────────────────────

// Subscription represents a VPN subscription.
type Subscription struct {
	ID                   uuid.UUID  `json:"id"`
	UserID               uuid.UUID  `json:"user_id"`
	PlanID               *uuid.UUID `json:"plan_id,omitempty"`
	PlanName             string     `json:"plan_name,omitempty"`
	Status               string     `json:"status"`
	Active               bool       `json:"active"`
	DataLimitBytes       int64      `json:"data_limit_bytes"`
	DataUsedBytes        int64      `json:"data_used_bytes"`
	DataRemainingBytes   int64      `json:"data_remaining_bytes"`
	Token                string     `json:"token,omitempty"`
	ExpiresAt            *time.Time `json:"expires_at,omitempty"`
	CurrentPeriodStart   *time.Time `json:"current_period_start,omitempty"`
	CurrentPeriodEnd     *time.Time `json:"current_period_end,omitempty"`
	CancelAtPeriodEnd    bool       `json:"cancel_at_period_end"`
	CanceledAt           *time.Time `json:"canceled_at,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	Plan                 *Plan      `json:"plan,omitempty"`
}

// DataUsedGB returns data used in gigabytes.
func (s *Subscription) DataUsedGB() float64 {
	return float64(s.DataUsedBytes) / 1_073_741_824
}

// DataRemainingGB returns data remaining in gigabytes.
func (s *Subscription) DataRemainingGB() float64 {
	return float64(s.DataRemainingBytes) / 1_073_741_824
}

// ── Connection ───────────────────────────────────────────────────────────────

// Connection statuses.
const (
	StatusReady         = "ready"
	StatusProvisioning  = "provisioning"
	StatusFailed        = "failed"
	StatusLimitExceeded = "limit_exceeded"
	StatusDisconnected  = "disconnected"
	StatusNotConnected  = "not_connected"
)

// ConnectResponse is returned by Connect and status-polling endpoints.
type ConnectResponse struct {
	Status       string          `json:"status"`
	SessionID    string          `json:"session_id,omitempty"`
	LocationSlug string          `json:"location_slug,omitempty"`
	Label        string          `json:"label,omitempty"`
	ConfigLink   string          `json:"config_link,omitempty"`
	TvpnURL      string          `json:"tvpn_url,omitempty"`
	Online       bool            `json:"online,omitempty"`
	Message      string          `json:"message,omitempty"`
	Node         *ConnectionNode `json:"node,omitempty"`
	// Free tier limit fields
	BytesUsedToday  int64 `json:"bytes_used_today,omitempty"`
	DailyLimitBytes int64 `json:"daily_limit_bytes,omitempty"`
	UpgradeRequired bool  `json:"upgrade_required,omitempty"`
}

// ConnectionNode holds server details for an active connection.
type ConnectionNode struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	Country string `json:"country"`
	Name    string `json:"name"`
}

// Connection represents an active VPN connection session.
type Connection struct {
	ID           uuid.UUID  `json:"id"`
	LocationSlug string     `json:"location_slug"`
	NodeID       uuid.UUID  `json:"node_id"`
	Status       string     `json:"status"`
	ConfigLink   *string    `json:"config_link,omitempty"`
	Online       bool       `json:"online,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	ExpiresAt    time.Time  `json:"expires_at"`
}

// Config represents a provisioned VPN config for a location.
type Config struct {
	LocationSlug string          `json:"location_slug"`
	Label        string          `json:"label"`
	Status       string          `json:"status"`
	Online       bool            `json:"online"`
	TvpnURL      string          `json:"tvpn_url,omitempty"`
	ConfigURL    string          `json:"config_url,omitempty"`
	Node         *ConnectionNode `json:"node,omitempty"`
}

// ── Profile ──────────────────────────────────────────────────────────────────

// Profile represents the v2 subscription profile response.
type Profile struct {
	Active               bool       `json:"active"`
	Status               string     `json:"status"`
	Plan                 string     `json:"plan"`
	DataLimitGB          float64    `json:"data_limit_gb"`
	DataUsedGB           float64    `json:"data_used_gb"`
	DataRemainingGB      string     `json:"data_remaining_gb"`
	ExpiresAt            *time.Time `json:"expires_at"`
	ProvisionedLocations int        `json:"provisioned_locations"`
	OnlineLocations      int        `json:"online_locations"`
	CreatedAt            time.Time  `json:"created_at"`
}

// ── Health ────────────────────────────────────────────────────────────────────

// HealthResponse is the complete user state snapshot from GET /v2/health.
type HealthResponse struct {
	User           map[string]any   `json:"user"`
	Subscription   map[string]any   `json:"subscription"`
	Connections    []map[string]any `json:"connections"`
	Usage          map[string]any   `json:"usage"`
	NodesAvailable int              `json:"nodes_available"`
	Alerts         []string         `json:"alerts"`
	Hints          *HealthHints     `json:"_hints,omitempty"`
}

// HealthHints provides actionable suggestions for agents and CLIs.
type HealthHints struct {
	SuggestedAction   string   `json:"suggested_action"`
	AvailableActions  []string `json:"available_actions"`
	SubscriptionToken string   `json:"subscription_token,omitempty"`
}

// ── Payment ──────────────────────────────────────────────────────────────────

// Payment represents a payment record.
type Payment struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	PlanID         *uuid.UUID `json:"plan_id,omitempty"`
	AmountUSD      float64    `json:"amount_usd"`
	OriginalAmount float64    `json:"original_amount"`
	DiscountAmount float64    `json:"discount_amount"`
	Status         string     `json:"status"`
	PaymentMethod  string     `json:"payment_method"`
	ReceiptTxID    *string    `json:"receipt_tx_id,omitempty"`
	NetworkHash    *string    `json:"network_hash,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	Plan           *Plan      `json:"plan,omitempty"`
}

// Payment statuses.
const (
	PaymentPending   = "pending"
	PaymentConfirmed = "confirmed"
	PaymentCancelled = "cancelled"
)

// ── Usage ────────────────────────────────────────────────────────────────────

// UsageStatus is the response from GET /v2/usage.
type UsageStatus struct {
	DateUTC         string         `json:"date_utc"`
	BytesUsedToday  int64          `json:"bytes_used_today"`
	DailyLimitBytes int64          `json:"daily_limit_bytes"`
	BytesRemaining  int64          `json:"bytes_remaining"`
	LimitExceeded   bool           `json:"limit_exceeded"`
	HasPremium      bool           `json:"has_premium"`
	History         []UsageHistory `json:"history,omitempty"`
	Premium         *PremiumUsage  `json:"premium,omitempty"`
}

// UsageHistory represents one day's usage data.
type UsageHistory struct {
	DateUTC   string `json:"date_utc"`
	BytesUsed int64  `json:"bytes_used"`
}

// PremiumUsage holds premium subscription usage info.
type PremiumUsage struct {
	PlanName       string     `json:"plan_name"`
	DataLimitBytes int64      `json:"data_limit_bytes"`
	DataUsedBytes  int64      `json:"data_used_bytes"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
}

// UsageReport is the request body for reporting usage.
type UsageReport struct {
	BytesDelta   int64  `json:"bytes_delta"`
	LocationSlug string `json:"location_slug,omitempty"`
}

// ── Billing ──────────────────────────────────────────────────────────────────

// CheckoutSession is returned when creating a Stripe checkout session.
type CheckoutSession struct {
	SessionID string `json:"session_id"`
	URL       string `json:"url"`
}

// PortalSession is returned when creating a Stripe customer portal session.
type PortalSession struct {
	URL string `json:"url"`
}

// TrialStatus represents the trial plan status for a user.
type TrialStatus struct {
	Available bool       `json:"available"`
	Active    bool       `json:"active"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	PlanName  string     `json:"plan_name,omitempty"`
}

// Invoice represents a billing invoice.
type Invoice struct {
	ID               uuid.UUID  `json:"id"`
	StripeInvoiceID  string     `json:"stripe_invoice_id"`
	AmountTotal      int64      `json:"amount_total"`
	AmountPaid       int64      `json:"amount_paid"`
	Currency         string     `json:"currency"`
	Status           string     `json:"status"`
	HostedInvoiceURL *string    `json:"hosted_invoice_url,omitempty"`
	InvoicePDFURL    *string    `json:"invoice_pdf_url,omitempty"`
	PeriodStart      *time.Time `json:"period_start,omitempty"`
	PeriodEnd        *time.Time `json:"period_end,omitempty"`
	PaidAt           *time.Time `json:"paid_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

// ── Ticket ───────────────────────────────────────────────────────────────────

// Ticket represents a support ticket.
type Ticket struct {
	ID         uuid.UUID      `json:"id"`
	Subject    string         `json:"subject"`
	Status     string         `json:"status"`
	Priority   string         `json:"priority"`
	Department string         `json:"department"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	Messages   []TicketMessage `json:"messages,omitempty"`
}

// TicketMessage represents a single message in a ticket thread.
type TicketMessage struct {
	ID         uuid.UUID `json:"id"`
	SenderRole string    `json:"sender_role"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}

// Ticket statuses.
const (
	TicketStatusOpen         = "open"
	TicketStatusWaitingUser  = "waiting_user"
	TicketStatusWaitingAdmin = "waiting_admin"
	TicketStatusClosed       = "closed"
)

// ── Session ──────────────────────────────────────────────────────────────────

// UserSession represents an active login session.
type UserSession struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	LastSeenAt  time.Time  `json:"last_seen_at"`
	ExpiresAt   time.Time  `json:"expires_at"`
	RevokedAt   *time.Time `json:"revoked_at,omitempty"`
	IP          string     `json:"ip"`
	UserAgent   string     `json:"user_agent"`
	LoginMethod string     `json:"login_method"`
	Provider    string     `json:"provider"`
	Client      string     `json:"client"`
	Country     string     `json:"country"`
	City        string     `json:"city"`
	Region      string     `json:"region"`
}

// ── API Key ──────────────────────────────────────────────────────────────────

// APIKey represents a personal API access token.
type APIKey struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	KeyPrefix  string     `json:"key_prefix"`
	Scopes     []string   `json:"scopes"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// CreateKeyResponse is returned when creating a new API key.
// The raw key is only shown once — store it securely.
type CreateKeyResponse struct {
	RawKey string `json:"raw_key"`
	APIKey APIKey `json:"key"`
}

// KeyUsage represents usage stats for an API key.
type KeyUsage struct {
	RequestsTotal int     `json:"requests_total"`
	RequestsOK    int     `json:"requests_ok"`
	RequestsErr   int     `json:"requests_err"`
	AvgLatencyMs  float64 `json:"avg_latency_ms"`
}

// RequestLog represents a single logged API request.
type RequestLog struct {
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	LatencyMs int    `json:"latency_ms"`
	IP        string `json:"ip"`
	CreatedAt string `json:"created_at"`
}

// API key scopes.
const (
	ScopeReadAccount       = "read:account"
	ScopeReadPlans         = "read:plans"
	ScopeReadNodes         = "read:nodes"
	ScopeReadSubscriptions = "read:subscriptions"
	ScopeReadConfig        = "read:config"
	ScopeWriteConnect      = "write:connect"
	ScopeReadUsage         = "read:usage"
	ScopeReadPayments      = "read:payments"
)

// AllScopes returns all available API key scopes.
func AllScopes() []string {
	return []string{
		ScopeReadAccount,
		ScopeReadPlans,
		ScopeReadNodes,
		ScopeReadSubscriptions,
		ScopeReadConfig,
		ScopeWriteConnect,
		ScopeReadUsage,
		ScopeReadPayments,
	}
}

// ── 2FA ───────────────────────────────────────────────────────────────────────

// TwoFASetupResponse is returned when initiating 2FA setup.
type TwoFASetupResponse struct {
	Secret  string `json:"secret"`
	QRCode  string `json:"qr_code"` // base64-encoded PNG
	OTPAuth string `json:"otpauth_url,omitempty"`
}

// TwoFAStatus holds the 2FA enrollment status.
type TwoFAStatus struct {
	Enabled bool `json:"enabled"`
}
