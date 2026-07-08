package vynvpn

import "context"

// AuthService handles authentication endpoints.
type AuthService struct {
	client *Client
}

// ── Request types ─────────────────────────────────────────────────────────────

// LoginRequest is the body for email+password login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest is the body for email+password registration.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login2FARequest completes a 2FA-protected login.
type Login2FARequest struct {
	LoginToken string `json:"login_token"`
	Code       string `json:"code"`
}

// TelegramLoginRequest authenticates via Telegram.
type TelegramLoginRequest struct {
	InitData   string            `json:"initData,omitempty"`
	WidgetData map[string]string `json:"widgetData,omitempty"`
}

// ForgotPasswordRequest initiates a password reset.
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// ResetPasswordRequest completes a password reset.
type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// NativeTokenRequest verifies an OAuth ID token (mobile SDKs).
type NativeTokenRequest struct {
	Provider string `json:"provider"` // "google" or "microsoft"
	IDToken  string `json:"id_token"`
}

// ── Methods ───────────────────────────────────────────────────────────────────

// Register creates a new email+password account.
// On success the JWT is stored on the client automatically.
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	var resp RegisterResponse
	if err := s.client.http.post(ctx, "/auth/register", req, &resp); err != nil {
		return nil, err
	}
	if resp.Token != "" {
		s.client.SetToken(resp.Token)
	}
	return &resp, nil
}

// Login authenticates with email+password.
// If 2FA is required, LoginResponse.Requires2FA is true — call Login2FA next.
// On success the JWT is stored on the client automatically.
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	var resp LoginResponse
	if err := s.client.http.post(ctx, "/auth/login", req, &resp); err != nil {
		return nil, err
	}
	if resp.Token != "" {
		s.client.SetToken(resp.Token)
	}
	return &resp, nil
}

// Login2FA completes a 2FA-protected login.
// On success the JWT is stored on the client automatically.
func (s *AuthService) Login2FA(ctx context.Context, req *Login2FARequest) (*LoginResponse, error) {
	var resp LoginResponse
	if err := s.client.http.post(ctx, "/auth/login/2fa", req, &resp); err != nil {
		return nil, err
	}
	if resp.Token != "" {
		s.client.SetToken(resp.Token)
	}
	return &resp, nil
}

// LoginTelegram authenticates via Telegram initData or widget data.
// On success the JWT is stored on the client automatically.
func (s *AuthService) LoginTelegram(ctx context.Context, req *TelegramLoginRequest) (*LoginResponse, error) {
	var resp LoginResponse
	if err := s.client.http.post(ctx, "/auth/telegram", req, &resp); err != nil {
		return nil, err
	}
	if resp.Token != "" {
		s.client.SetToken(resp.Token)
	}
	return &resp, nil
}

// Refresh exchanges the current JWT for a fresh one.
// The new token is stored on the client automatically.
func (s *AuthService) Refresh(ctx context.Context) (*RefreshResponse, error) {
	var resp RefreshResponse
	if err := s.client.http.post(ctx, "/auth/refresh", nil, &resp); err != nil {
		return nil, err
	}
	if resp.Token != "" {
		s.client.SetToken(resp.Token)
	}
	return &resp, nil
}

// ForgotPassword sends a password-reset email.
func (s *AuthService) ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) error {
	return s.client.http.post(ctx, "/auth/forgot-password", req, nil)
}

// ResetPassword sets a new password using a reset token.
func (s *AuthService) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	return s.client.http.post(ctx, "/auth/reset-password", req, nil)
}

// ResendVerification resends the email verification link.
func (s *AuthService) ResendVerification(ctx context.Context) error {
	return s.client.http.post(ctx, "/auth/verify-email/resend", nil, nil)
}

// VerifyEmailCode verifies the email using a 6-digit in-app code.
func (s *AuthService) VerifyEmailCode(ctx context.Context, code string) error {
	return s.client.http.post(ctx, "/auth/verify-email/code", map[string]string{"code": code}, nil)
}

// Logout invalidates the current session.
func (s *AuthService) Logout(ctx context.Context) error {
	err := s.client.http.post(ctx, "/auth/logout", nil, nil)
	s.client.SetToken("")
	return err
}

// LinkTelegram attaches a Telegram identity to the authenticated account.
func (s *AuthService) LinkTelegram(ctx context.Context, req *TelegramLoginRequest) error {
	return s.client.http.post(ctx, "/api/auth/link/telegram", req, nil)
}

// OAuthGoogleURL returns the Google OAuth initiation URL.
// Redirect the user's browser to this URL to begin the flow.
func (s *AuthService) OAuthGoogleURL() string {
	return s.client.baseURL + "/auth/oauth/google"
}

// OAuthMicrosoftURL returns the Microsoft OAuth initiation URL.
func (s *AuthService) OAuthMicrosoftURL() string {
	return s.client.baseURL + "/auth/oauth/microsoft"
}

// OAuthNativeToken verifies a native OAuth ID token (Google/Microsoft mobile SDKs).
// On success the JWT is stored on the client automatically.
func (s *AuthService) OAuthNativeToken(ctx context.Context, req *NativeTokenRequest) (*LoginResponse, error) {
	var resp LoginResponse
	if err := s.client.http.post(ctx, "/auth/oauth/token", req, &resp); err != nil {
		return nil, err
	}
	if resp.Token != "" {
		s.client.SetToken(resp.Token)
	}
	return &resp, nil
}

// CreateOAuthSession creates a polling OAuth session for desktop apps.
func (s *AuthService) CreateOAuthSession(ctx context.Context) (*OAuthSessionResponse, error) {
	var resp OAuthSessionResponse
	if err := s.client.http.post(ctx, "/auth/oauth/session", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// PollOAuthSession checks the status of an OAuth polling session.
// When status is "consumed" the JWT is stored on the client automatically.
func (s *AuthService) PollOAuthSession(ctx context.Context, sessionID string) (*OAuthSessionStatus, error) {
	var resp OAuthSessionStatus
	q := makeQuery("session_id", sessionID)
	if err := s.client.http.get(ctx, "/auth/oauth/session/status", q, &resp); err != nil {
		return nil, err
	}
	if resp.Token != "" {
		s.client.SetToken(resp.Token)
	}
	return &resp, nil
}

// RequestReactivation requests reactivation of a soft-deleted account.
func (s *AuthService) RequestReactivation(ctx context.Context, email string) error {
	return s.client.http.post(ctx, "/auth/reactivate/request", map[string]string{"email": email}, nil)
}

// ConfirmReactivation confirms account reactivation with the received token.
func (s *AuthService) ConfirmReactivation(ctx context.Context, token string) error {
	return s.client.http.post(ctx, "/auth/reactivate/confirm", map[string]string{"token": token}, nil)
}
