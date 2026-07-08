package vynvpn

import (
	"net/http"
	"testing"
	"time"
)

func TestWithBaseURL(t *testing.T) {
	c := New(WithBaseURL("https://custom.api.com"))
	if c.baseURL != "https://custom.api.com" {
		t.Errorf("baseURL = %q, want %q", c.baseURL, "https://custom.api.com")
	}
}

func TestWithAPIKey(t *testing.T) {
	c := New(WithAPIKey("vyn_testkey"))
	if !c.http.IsAPIKeyAuth() {
		t.Error("expected IsAPIKeyAuth() = true")
	}
	if c.http.apiKey != "vyn_testkey" {
		t.Errorf("apiKey = %q, want %q", c.http.apiKey, "vyn_testkey")
	}
}

func TestWithToken(t *testing.T) {
	c := New(WithToken("jwt.token.here"))
	if c.http.IsAPIKeyAuth() {
		t.Error("expected IsAPIKeyAuth() = false for JWT client")
	}
	if c.http.token != "jwt.token.here" {
		t.Errorf("token = %q, want %q", c.http.token, "jwt.token.here")
	}
}

func TestWithTimeout(t *testing.T) {
	c := New(WithTimeout(60 * time.Second))
	// httpClient.Timeout is on the inner http.Client
	if c.http.client.Timeout != 60*time.Second {
		t.Errorf("timeout = %v, want 60s", c.http.client.Timeout)
	}
}

func TestWithHTTPClient(t *testing.T) {
	custom := &http.Client{Timeout: 5 * time.Second}
	c := New(WithHTTPClient(custom))
	if c.http.client != custom {
		t.Error("expected custom http.Client to be used")
	}
}

func TestDefaultBaseURL(t *testing.T) {
	c := New()
	if c.baseURL != DefaultBaseURL {
		t.Errorf("baseURL = %q, want %q", c.baseURL, DefaultBaseURL)
	}
}

func TestAPIKeyTakesPrecedenceOverToken(t *testing.T) {
	c := New(WithAPIKey("vyn_key"), WithToken("jwt"))
	if !c.http.IsAPIKeyAuth() {
		t.Error("API key should take precedence")
	}
}
