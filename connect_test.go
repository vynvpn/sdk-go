package vynvpn

import (
	"context"
	"testing"
)

func TestConnect_JWT(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "ready", "session_id": "sess-123", "config_link": "vless://...",
	}, WithToken("jwt"))
	defer srv.Close()

	resp, err := client.Connect.Connect(context.Background(), &ConnectRequest{
		LocationSlug: "de001",
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	if resp.Status != "ready" {
		t.Errorf("Status = %q, want %q", resp.Status, "ready")
	}
	if resp.ConfigLink != "vless://..." {
		t.Errorf("ConfigLink = %q", resp.ConfigLink)
	}
	if req.Path != "/v2/connect" {
		t.Errorf("path = %q, want /v2/connect", req.Path)
	}
}

func TestConnect_APIKey(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "provisioning", "session_id": "sess-456",
	}, WithAPIKey("vyn_key"))
	defer srv.Close()

	resp, err := client.Connect.Connect(context.Background(), &ConnectRequest{
		LocationSlug: "us001",
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	if resp.Status != "provisioning" {
		t.Errorf("Status = %q, want %q", resp.Status, "provisioning")
	}
	if req.Path != "/api/v1/connect" {
		t.Errorf("path = %q, want /api/v1/connect", req.Path)
	}
}

func TestConnect_Disconnect(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{"status": "disconnected"}, WithToken("jwt"))
	defer srv.Close()

	resp, err := client.Connect.Disconnect(context.Background(), &DisconnectRequest{
		LocationSlug: "de001",
	})
	if err != nil {
		t.Fatalf("Disconnect failed: %v", err)
	}
	if resp.Status != "disconnected" {
		t.Errorf("Status = %q, want %q", resp.Status, "disconnected")
	}
	if req.Path != "/v2/disconnect" {
		t.Errorf("path = %q, want /v2/disconnect", req.Path)
	}
}

func TestConnect_Status(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "ready", "config_link": "vless://ok",
	}, WithToken("jwt"))
	defer srv.Close()

	resp, err := client.Connect.Status(context.Background(), "sess-789")
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}
	if resp.Status != "ready" {
		t.Errorf("Status = %q, want %q", resp.Status, "ready")
	}
	if req.Path != "/v2/status/sess-789" {
		t.Errorf("path = %q", req.Path)
	}
}

// ── Inbound Configuration Tests ──────────────────────────────────────────────

func TestConnect_WithInboundConfig_Socks(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "ready", "session_id": "sess-s5", "config_link": "socks5://user:pass@node.example.com:1080#de001",
	}, WithToken("jwt"))
	defer srv.Close()

	resp, err := client.Connect.Connect(context.Background(), &ConnectRequest{
		LocationSlug: "de001",
		InboundConfig: &InboundConfig{
			Protocol: "socks",
			Port:     1080,
			Network:  "tcp",
			Security: "none",
		},
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	if resp.Status != "ready" {
		t.Errorf("Status = %q, want %q", resp.Status, "ready")
	}
	if resp.ConfigLink != "socks5://user:pass@node.example.com:1080#de001" {
		t.Errorf("ConfigLink = %q", resp.ConfigLink)
	}
	// Verify inbound_config was sent in the request body
	if req.Body == nil {
		t.Fatal("request body is nil")
	}
	ic, ok := req.Body["inbound_config"].(map[string]any)
	if !ok {
		t.Fatal("inbound_config not found in request body")
	}
	if ic["protocol"] != "socks" {
		t.Errorf("protocol = %v, want socks", ic["protocol"])
	}
	if ic["network"] != "tcp" {
		t.Errorf("network = %v, want tcp", ic["network"])
	}
}

func TestConnect_WithInboundConfig_VLESS_Reality(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "ready", "config_link": "vless://uuid@host:443?security=reality#remark",
	}, WithToken("jwt"))
	defer srv.Close()

	resp, err := client.Connect.Connect(context.Background(), &ConnectRequest{
		LocationSlug: "de001",
		InboundConfig: &InboundConfig{
			Protocol: "vless",
			Network:  "tcp",
			Security: "reality",
			Flow:     "xtls-rprx-vision",
			RealitySettings: &RealitySettings{
				Dest:        "www.google.com:443",
				ServerNames: []string{"www.google.com"},
				ShortIds:    []string{"abcdef"},
				Settings: RealityPublicSettings{
					PublicKey:   "test-pub-key",
					Fingerprint: "chrome",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	if resp.Status != "ready" {
		t.Errorf("Status = %q", resp.Status)
	}

	ic := req.Body["inbound_config"].(map[string]any)
	if ic["protocol"] != "vless" {
		t.Errorf("protocol = %v", ic["protocol"])
	}
	if ic["security"] != "reality" {
		t.Errorf("security = %v", ic["security"])
	}
	if ic["flow"] != "xtls-rprx-vision" {
		t.Errorf("flow = %v", ic["flow"])
	}
	rs, ok := ic["realitySettings"].(map[string]any)
	if !ok {
		t.Fatal("realitySettings not found")
	}
	if rs["dest"] != "www.google.com:443" {
		t.Errorf("dest = %v", rs["dest"])
	}
}

func TestConnect_WithInboundConfig_WS_TLS(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "ready", "config_link": "vless://uuid@cdn.example.com:443?type=ws#remark",
	}, WithToken("jwt"))
	defer srv.Close()

	_, err := client.Connect.Connect(context.Background(), &ConnectRequest{
		LocationSlug: "de001",
		InboundConfig: &InboundConfig{
			Protocol: "vless",
			Network:  "ws",
			Security: "tls",
			WSSettings: &WSSettings{
				Path: "/ws",
				Host: "cdn.example.com",
			},
			TLSSettings: &TLSSettings{
				ServerName:  "cdn.example.com",
				Fingerprint: "chrome",
			},
		},
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	ic := req.Body["inbound_config"].(map[string]any)
	if ic["network"] != "ws" {
		t.Errorf("network = %v", ic["network"])
	}
	ws := ic["wsSettings"].(map[string]any)
	if ws["path"] != "/ws" {
		t.Errorf("ws path = %v", ws["path"])
	}
	if ws["host"] != "cdn.example.com" {
		t.Errorf("ws host = %v", ws["host"])
	}
	tls := ic["tlsSettings"].(map[string]any)
	if tls["serverName"] != "cdn.example.com" {
		t.Errorf("tls serverName = %v", tls["serverName"])
	}
}

func TestConnect_WithInboundConfig_HTTP_Proxy(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "ready", "config_link": "https://user:pass@proxy.example.com:8080#de001",
	}, WithToken("jwt"))
	defer srv.Close()

	resp, err := client.Connect.Connect(context.Background(), &ConnectRequest{
		LocationSlug: "de001",
		InboundConfig: &InboundConfig{
			Protocol: "http",
			Port:     8080,
			Network:  "tcp",
			Security: "tls",
			TLSSettings: &TLSSettings{
				ServerName: "proxy.example.com",
			},
		},
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	if resp.ConfigLink != "https://user:pass@proxy.example.com:8080#de001" {
		t.Errorf("ConfigLink = %q", resp.ConfigLink)
	}

	ic := req.Body["inbound_config"].(map[string]any)
	if ic["protocol"] != "http" {
		t.Errorf("protocol = %v", ic["protocol"])
	}
	if ic["port"] != float64(8080) {
		t.Errorf("port = %v", ic["port"])
	}
}

func TestConnect_WithInboundConfig_gRPC(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "ready", "config_link": "vmess://...",
	}, WithToken("jwt"))
	defer srv.Close()

	_, err := client.Connect.Connect(context.Background(), &ConnectRequest{
		LocationSlug: "de001",
		InboundConfig: &InboundConfig{
			Protocol: "vmess",
			Network:  "grpc",
			Security: "tls",
			GRPCSettings: &GRPCSettings{
				ServiceName: "tunnel",
				MultiMode:   true,
			},
		},
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	ic := req.Body["inbound_config"].(map[string]any)
	grpc := ic["grpcSettings"].(map[string]any)
	if grpc["serviceName"] != "tunnel" {
		t.Errorf("serviceName = %v", grpc["serviceName"])
	}
	if grpc["multiMode"] != true {
		t.Errorf("multiMode = %v", grpc["multiMode"])
	}
}

func TestConnect_WithInboundConfig_APIKey(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "ready", "config_link": "socks5://user:pass@host:1080#de001",
	}, WithAPIKey("vyn_test"))
	defer srv.Close()

	_, err := client.Connect.Connect(context.Background(), &ConnectRequest{
		LocationSlug: "de001",
		InboundConfig: &InboundConfig{
			Protocol: "socks",
			Network:  "tcp",
			Security: "none",
		},
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	if req.Path != "/api/v1/connect" {
		t.Errorf("path = %q, want /api/v1/connect", req.Path)
	}
}

// ── Content Filter Tests ─────────────────────────────────────────────────────

func TestConnect_WithContentFilter_FamilySafe(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "ready", "config_link": "vless://...",
	}, WithToken("jwt"))
	defer srv.Close()

	_, err := client.Connect.Connect(context.Background(), &ConnectRequest{
		LocationSlug:  "de001",
		ContentFilter: ContentFilterFamilySafe(),
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	cf, ok := req.Body["content_filter"].(map[string]any)
	if !ok {
		t.Fatal("content_filter not found in request body")
	}
	cats, ok := cf["block_categories"].([]any)
	if !ok || len(cats) != 3 {
		t.Fatalf("block_categories = %v, want 3 items", cf["block_categories"])
	}
	if cats[0] != "porn" || cats[1] != "gambling" || cats[2] != "malware" {
		t.Errorf("block_categories = %v", cats)
	}
	dns, ok := cf["dns"].(map[string]any)
	if !ok {
		t.Fatal("dns not found in content_filter")
	}
	servers := dns["servers"].([]any)
	if len(servers) != 1 || servers[0] != "family-cloudflare" {
		t.Errorf("dns servers = %v", servers)
	}
}

func TestConnect_WithContentFilter_KidSafe(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "ready", "config_link": "vless://...",
	}, WithToken("jwt"))
	defer srv.Close()

	_, err := client.Connect.Connect(context.Background(), &ConnectRequest{
		LocationSlug:  "de001",
		ContentFilter: ContentFilterKidSafe(),
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	cf := req.Body["content_filter"].(map[string]any)
	cats := cf["block_categories"].([]any)
	if len(cats) != 6 {
		t.Fatalf("block_categories should have 6 items, got %d: %v", len(cats), cats)
	}
	// Verify it contains social and gaming (strict mode)
	hasSocial, hasGaming := false, false
	for _, c := range cats {
		if c == "social" {
			hasSocial = true
		}
		if c == "gaming" {
			hasGaming = true
		}
	}
	if !hasSocial || !hasGaming {
		t.Errorf("KidSafe should include social and gaming: %v", cats)
	}
	dns := cf["dns"].(map[string]any)
	servers := dns["servers"].([]any)
	if servers[0] != "family-cleanbrowsing" {
		t.Errorf("dns servers = %v, want family-cleanbrowsing", servers)
	}
}

func TestConnect_WithContentFilter_AdsOnly(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "ready", "config_link": "vless://...",
	}, WithToken("jwt"))
	defer srv.Close()

	_, err := client.Connect.Connect(context.Background(), &ConnectRequest{
		LocationSlug:  "de001",
		ContentFilter: ContentFilterAdsOnly(),
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	cf := req.Body["content_filter"].(map[string]any)
	cats := cf["block_categories"].([]any)
	if len(cats) != 1 || cats[0] != "ads" {
		t.Errorf("block_categories = %v, want [ads]", cats)
	}
	dns := cf["dns"].(map[string]any)
	servers := dns["servers"].([]any)
	if servers[0] != "family-adguard" {
		t.Errorf("dns servers = %v, want family-adguard", servers)
	}
}

func TestConnect_WithContentFilter_Custom(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "ready", "config_link": "vless://...",
	}, WithToken("jwt"))
	defer srv.Close()

	_, err := client.Connect.Connect(context.Background(), &ConnectRequest{
		LocationSlug: "de001",
		ContentFilter: &ContentFilter{
			BlockCategories: []string{"porn", "drugs"},
			BlockDomains:    []string{"domain:onlyfans.com", "keyword:adult"},
			DNS: &DNSFilter{
				Servers: []string{"1.1.1.3", "9.9.9.9"},
			},
		},
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	cf := req.Body["content_filter"].(map[string]any)
	cats := cf["block_categories"].([]any)
	if len(cats) != 2 {
		t.Fatalf("block_categories = %v, want 2 items", cats)
	}
	domains := cf["block_domains"].([]any)
	if len(domains) != 2 {
		t.Fatalf("block_domains = %v, want 2 items", domains)
	}
	if domains[0] != "domain:onlyfans.com" {
		t.Errorf("block_domains[0] = %v", domains[0])
	}
	dns := cf["dns"].(map[string]any)
	servers := dns["servers"].([]any)
	if len(servers) != 2 || servers[0] != "1.1.1.3" || servers[1] != "9.9.9.9" {
		t.Errorf("dns servers = %v", servers)
	}
}

func TestConnect_WithContentFilter_CategoriesOnly(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "ready", "config_link": "vless://...",
	}, WithToken("jwt"))
	defer srv.Close()

	_, err := client.Connect.Connect(context.Background(), &ConnectRequest{
		LocationSlug: "de001",
		ContentFilter: &ContentFilter{
			BlockCategories: []string{"gambling", "piracy"},
		},
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	cf := req.Body["content_filter"].(map[string]any)
	cats := cf["block_categories"].([]any)
	if len(cats) != 2 || cats[0] != "gambling" || cats[1] != "piracy" {
		t.Errorf("block_categories = %v", cats)
	}
	// DNS should not be present
	if cf["dns"] != nil {
		t.Errorf("dns should be nil when not set, got %v", cf["dns"])
	}
}

func TestConnect_WithContentFilter_DNSOnly(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "ready", "config_link": "vless://...",
	}, WithToken("jwt"))
	defer srv.Close()

	_, err := client.Connect.Connect(context.Background(), &ConnectRequest{
		LocationSlug: "de001",
		ContentFilter: &ContentFilter{
			DNS: &DNSFilter{
				Servers: []string{"family-opendns"},
			},
		},
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	cf := req.Body["content_filter"].(map[string]any)
	dns := cf["dns"].(map[string]any)
	servers := dns["servers"].([]any)
	if len(servers) != 1 || servers[0] != "family-opendns" {
		t.Errorf("dns servers = %v, want [family-opendns]", servers)
	}
}

// ── Combined: InboundConfig + ContentFilter ──────────────────────────────────

func TestConnect_InboundConfig_And_ContentFilter(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "ready", "config_link": "socks5://user:pass@host:1080#de001",
	}, WithToken("jwt"))
	defer srv.Close()

	_, err := client.Connect.Connect(context.Background(), &ConnectRequest{
		LocationSlug: "de001",
		InboundConfig: &InboundConfig{
			Protocol: "socks",
			Port:     1080,
			Network:  "tcp",
			Security: "none",
		},
		ContentFilter: ContentFilterKidSafe(),
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	// Verify both are present in the body
	if req.Body["inbound_config"] == nil {
		t.Error("inbound_config should be present")
	}
	if req.Body["content_filter"] == nil {
		t.Error("content_filter should be present")
	}

	ic := req.Body["inbound_config"].(map[string]any)
	if ic["protocol"] != "socks" {
		t.Errorf("protocol = %v", ic["protocol"])
	}

	cf := req.Body["content_filter"].(map[string]any)
	cats := cf["block_categories"].([]any)
	if len(cats) != 6 {
		t.Errorf("block_categories should have 6 items for KidSafe, got %d", len(cats))
	}
}

func TestConnect_NilConfig_DoesNotSendFields(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"status": "ready", "config_link": "vless://...",
	}, WithToken("jwt"))
	defer srv.Close()

	_, err := client.Connect.Connect(context.Background(), &ConnectRequest{
		LocationSlug: "de001",
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	// When InboundConfig and ContentFilter are nil, they should not be in the JSON
	if req.Body["inbound_config"] != nil {
		t.Error("inbound_config should be omitted when nil")
	}
	if req.Body["content_filter"] != nil {
		t.Error("content_filter should be omitted when nil")
	}
}

// ── Content Filter Preset Unit Tests ─────────────────────────────────────────

func TestContentFilterFamilySafe(t *testing.T) {
	f := ContentFilterFamilySafe()
	if len(f.BlockCategories) != 3 {
		t.Errorf("FamilySafe should block 3 categories, got %d", len(f.BlockCategories))
	}
	if f.DNS == nil || len(f.DNS.Servers) != 1 {
		t.Fatal("FamilySafe should have 1 DNS server")
	}
	if f.DNS.Servers[0] != "family-cloudflare" {
		t.Errorf("DNS server = %q", f.DNS.Servers[0])
	}
}

func TestContentFilterKidSafe(t *testing.T) {
	f := ContentFilterKidSafe()
	if len(f.BlockCategories) != 6 {
		t.Errorf("KidSafe should block 6 categories, got %d: %v", len(f.BlockCategories), f.BlockCategories)
	}
	if f.DNS == nil || f.DNS.Servers[0] != "family-cleanbrowsing" {
		t.Error("KidSafe should use family-cleanbrowsing")
	}
}

func TestContentFilterAdsOnly(t *testing.T) {
	f := ContentFilterAdsOnly()
	if len(f.BlockCategories) != 1 || f.BlockCategories[0] != "ads" {
		t.Errorf("AdsOnly should block [ads], got %v", f.BlockCategories)
	}
	if f.DNS == nil || f.DNS.Servers[0] != "family-adguard" {
		t.Error("AdsOnly should use family-adguard")
	}
}
