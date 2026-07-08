package vynvpn

import "testing"

func TestNew_ServicesWired(t *testing.T) {
	c := New()
	if c.Auth == nil { t.Error("Auth is nil") }
	if c.Health == nil { t.Error("Health is nil") }
	if c.Nodes == nil { t.Error("Nodes is nil") }
	if c.Subscriptions == nil { t.Error("Subscriptions is nil") }
	if c.Connect == nil { t.Error("Connect is nil") }
	if c.Plans == nil { t.Error("Plans is nil") }
	if c.Usage == nil { t.Error("Usage is nil") }
	if c.Payments == nil { t.Error("Payments is nil") }
	if c.Account == nil { t.Error("Account is nil") }
	if c.Billing == nil { t.Error("Billing is nil") }
	if c.Tickets == nil { t.Error("Tickets is nil") }
	if c.Sessions == nil { t.Error("Sessions is nil") }
	if c.APIKeys == nil { t.Error("APIKeys is nil") }
	if c.TwoFA == nil { t.Error("TwoFA is nil") }
}

func TestClient_SetAndGetToken(t *testing.T) {
	c := New()
	if c.Token() != "" {
		t.Error("expected empty token initially")
	}
	c.SetToken("abc.def.ghi")
	if c.Token() != "abc.def.ghi" {
		t.Errorf("Token() = %q, want %q", c.Token(), "abc.def.ghi")
	}
}

func TestClient_BaseURL(t *testing.T) {
	c := New(WithBaseURL("https://api.example.com"))
	if c.BaseURL() != "https://api.example.com" {
		t.Errorf("BaseURL() = %q, want %q", c.BaseURL(), "https://api.example.com")
	}
}
