# VynVPN Go SDK

Official Go SDK for the VynVPN API. Supports both API key auth (for CLIs, agents, and server-side integrations) and JWT auth (for user-session apps).

## Installation

```bash
go get github.com/vynvpn/sdk-go
```

Requires Go 1.22+.

---

## Authentication

Two modes, pick one:

**API key** — recommended for CLIs, agents, automation. Long-lived, scoped, no expiry.
Generate from the VynVPN dashboard → Settings → API Keys (verified email required).

```go
client := vynvpn.New(
    vynvpn.WithAPIKey("vyn_your_key_here"),
)
```

**JWT** — for interactive apps (mobile, desktop, web). Login once, store the token.

```go
client := vynvpn.New(vynvpn.WithBaseURL("https://api.vynvpn.com"))
auth, err := client.Auth.Login(ctx, &vynvpn.LoginRequest{
    Email:    "user@example.com",
    Password: "secret",
})
// JWT is stored automatically on the client after login
```

---

## Quick Start

### API key (agent / CLI)

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    vynvpn "github.com/vynvpn/sdk-go"
)

func main() {
    client := vynvpn.New(
        vynvpn.WithAPIKey("vyn_your_key_here"),
        vynvpn.WithBaseURL("https://api.vynvpn.com"),
    )
    ctx := context.Background()

    // Full account state in one call — great starting point for agents
    health, err := client.Health.Get(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Alerts: %v\n", health.Alerts)
    fmt.Printf("Suggested action: %s\n", health.Hints.SuggestedAction)

    // Connect — subscription auto-detected, best server auto-picked
    resp, err := client.Connect.Connect(ctx, &vynvpn.ConnectRequest{})
    if err != nil {
        log.Fatal(err)
    }

    // Poll if provisioning
    if resp.Status == vynvpn.StatusProvisioning {
        resp, _ = client.Connect.WaitForReady(ctx, resp.SessionID, 30*time.Second)
    }
    fmt.Printf("Config: %s\n", resp.ConfigLink)

    // Disconnect
    client.Connect.Disconnect(ctx, &vynvpn.DisconnectRequest{
        LocationSlug: resp.LocationSlug,
    })
}
```

### JWT (interactive app)

```go
client := vynvpn.New(vynvpn.WithBaseURL("https://api.vynvpn.com"))
ctx := context.Background()

auth, err := client.Auth.Login(ctx, &vynvpn.LoginRequest{
    Email:    "user@example.com",
    Password: "secret",
})
if err != nil {
    log.Fatal(err)
}

// Handle 2FA if required
if auth.Requires2FA {
    auth, err = client.Auth.Login2FA(ctx, &vynvpn.Login2FARequest{
        LoginToken: auth.LoginToken,
        Code:       "123456",
    })
}

// All subsequent calls are authenticated automatically
me, _ := client.Account.GetMe(ctx)
fmt.Printf("Hello %s\n", me.FirstName)
```

---

## Services

| Service | Description |
|---------|-------------|
| `client.Health` | Full account state snapshot — best first call for agents |
| `client.Auth` | Login, register, 2FA, OAuth, Telegram, password reset, verify email |
| `client.Account` | Profile, account deletion |
| `client.Nodes` | List VPN server locations |
| `client.Plans` | List available plans |
| `client.Subscriptions` | List, get, configs, profile |
| `client.Connect` | Connect, disconnect, status polling, content filtering |
| `client.Usage` | Data usage status and reporting |
| `client.Billing` | Stripe checkout, portal, cancel, invoices, trial |
| `client.Payments` | Payment history, create payments |
| `client.Tickets` | Support tickets |
| `client.Sessions` | Active login session management |
| `client.APIKeys` | API key CRUD and usage stats |
| `client.TwoFA` | 2FA setup, enable, disable |

---

## VPN Connect

All fields are optional. The API fills in what's missing:

```go
// Fully automatic — detects subscription, picks best server by geo-IP + load
client.Connect.Connect(ctx, &vynvpn.ConnectRequest{})

// Country preference
client.Connect.Connect(ctx, &vynvpn.ConnectRequest{
    PreferredCountry: "Germany",
})

// Explicit location
client.Connect.Connect(ctx, &vynvpn.ConnectRequest{
    LocationSlug: "de001",
})

// With subscription token (JWT mode)
client.Connect.Connect(ctx, &vynvpn.ConnectRequest{
    Token:        "sub-token",
    LocationSlug: "de001",
})
```

Handle provisioning:

```go
resp, _ := client.Connect.Connect(ctx, &vynvpn.ConnectRequest{})

if resp.Status == vynvpn.StatusProvisioning {
    // Option 1: block until ready (up to 30s)
    resp, _ = client.Connect.WaitForReady(ctx, resp.SessionID, 30*time.Second)

    // Option 2: poll manually
    for resp.Status == vynvpn.StatusProvisioning {
        time.Sleep(2 * time.Second)
        resp, _ = client.Connect.Status(ctx, resp.SessionID)
    }
}

if resp.Status == vynvpn.StatusReady {
    fmt.Println(resp.ConfigLink)
}
```

---

## Dynamic Inbound Configuration

Override the VPN protocol, transport, and security per-connection:

```go
// SOCKS5 proxy (for Chrome extensions, browsers)
resp, _ := client.Connect.Connect(ctx, &vynvpn.ConnectRequest{
    LocationSlug: "de001",
    InboundConfig: &vynvpn.InboundConfig{
        Protocol: "socks",
        Port:     1080,
        Network:  "tcp",
        Security: "none",
    },
})
// resp.ConfigLink = "socks5://user:pass@node.example.com:1080#de001-abc12345"

// VLESS + Reality (maximum stealth)
resp, _ := client.Connect.Connect(ctx, &vynvpn.ConnectRequest{
    LocationSlug: "de001",
    InboundConfig: &vynvpn.InboundConfig{
        Protocol: "vless",
        Network:  "tcp",
        Security: "reality",
        Flow:     "xtls-rprx-vision",
        RealitySettings: &vynvpn.RealitySettings{
            Dest:        "www.google.com:443",
            ServerNames: []string{"www.google.com"},
            ShortIds:    []string{"abcdef"},
            Settings: vynvpn.RealityPublicSettings{
                PublicKey:   "your-public-key",
                Fingerprint: "chrome",
            },
        },
    },
})

// VLESS + WebSocket (CDN compatible)
resp, _ := client.Connect.Connect(ctx, &vynvpn.ConnectRequest{
    LocationSlug: "de001",
    InboundConfig: &vynvpn.InboundConfig{
        Protocol: "vless",
        Network:  "ws",
        Security: "tls",
        WSSettings: &vynvpn.WSSettings{
            Path: "/ws",
            Host: "cdn.example.com",
        },
        TLSSettings: &vynvpn.TLSSettings{
            ServerName:  "cdn.example.com",
            Fingerprint: "chrome",
        },
    },
})

// VMess + gRPC
resp, _ := client.Connect.Connect(ctx, &vynvpn.ConnectRequest{
    LocationSlug: "de001",
    InboundConfig: &vynvpn.InboundConfig{
        Protocol: "vmess",
        Network:  "grpc",
        Security: "tls",
        GRPCSettings: &vynvpn.GRPCSettings{
            ServiceName: "tunnel",
            MultiMode:   true,
        },
    },
})

// HTTP/HTTPS proxy (for Chrome extensions)
resp, _ := client.Connect.Connect(ctx, &vynvpn.ConnectRequest{
    LocationSlug: "de001",
    InboundConfig: &vynvpn.InboundConfig{
        Protocol: "http",
        Port:     8080,
        Network:  "tcp",
        Security: "tls",
        TLSSettings: &vynvpn.TLSSettings{
            ServerName: "proxy.example.com",
        },
    },
})
// resp.ConfigLink = "https://user:pass@proxy.example.com:8080#de001-abc12345"
```

### Supported protocols

| Protocol | Config link format | Use case |
|----------|-------------------|----------|
| `vless` | `vless://uuid@host:port?params#remark` | Desktop/mobile apps (xray, sing-box) |
| `vmess` | `vmess://base64(json)` | Desktop/mobile apps (v2ray) |
| `trojan` | `trojan://pass@host:port?params#remark` | Desktop/mobile apps |
| `shadowsocks` | `ss://base64(method:pass)@host:port#remark` | Desktop/mobile apps |
| `socks` | `socks5://user:pass@host:port#remark` | Browser extensions, system proxy |
| `http` | `http://user:pass@host:port#remark` | Browser extensions, system proxy |

### Supported transports

| Transport | Struct | Notes |
|-----------|--------|-------|
| `tcp` | `TCPSettings` | HTTP obfuscation header available |
| `ws` | `WSSettings` | CDN-compatible, path + host |
| `grpc` | `GRPCSettings` | Multi-mode, authority |
| `h2` | `HTTPSettings` | HTTP/2 |
| `quic` | `QUICSettings` | UDP-based, with header types |
| `kcp` | `KCPSettings` | mKCP with seed, capacity tuning |
| `httpupgrade` | `HTTPUpgradeSettings` | HTTP Upgrade based |
| `splithttp` | `SplitHTTPSettings` | Split HTTP based |

### Supported security layers

| Security | Struct | Notes |
|----------|--------|-------|
| `none` | — | No encryption |
| `tls` | `TLSSettings` | Standard TLS with fingerprint, ALPN |
| `reality` | `RealitySettings` | Stealth — mimics real HTTPS sites |

---

## Content Filtering (Family Control)

Block adult content, gambling, ads, and more. Two methods available independently or combined:

### Quick presets

```go
// Family safe — blocks porn, gambling, malware + Cloudflare family DNS
resp, _ := client.Connect.Connect(ctx, &vynvpn.ConnectRequest{
    LocationSlug:  "de001",
    ContentFilter: vynvpn.ContentFilterFamilySafe(),
})

// Kid safe — strict: blocks porn, gambling, malware, drugs, social, gaming
resp, _ := client.Connect.Connect(ctx, &vynvpn.ConnectRequest{
    LocationSlug:  "de001",
    ContentFilter: vynvpn.ContentFilterKidSafe(),
})

// Ads only — blocks ads and trackers
resp, _ := client.Connect.Connect(ctx, &vynvpn.ConnectRequest{
    LocationSlug:  "de001",
    ContentFilter: vynvpn.ContentFilterAdsOnly(),
})
```

### Custom filter

```go
resp, _ := client.Connect.Connect(ctx, &vynvpn.ConnectRequest{
    LocationSlug: "de001",
    ContentFilter: &vynvpn.ContentFilter{
        // Block by category (routing rules → blackhole)
        BlockCategories: []string{"porn", "gambling", "malware", "drugs"},

        // Block specific domains
        BlockDomains: []string{
            "domain:onlyfans.com",
            "keyword:adult",
            "regexp:.*xxx.*",
        },

        // DNS-level filtering (broader coverage)
        DNS: &vynvpn.DNSFilter{
            Servers: []string{"family-cloudflare"},
        },
    },
})
```

### Block categories

| Category | What it blocks |
|----------|---------------|
| `porn` | Adult/pornographic content |
| `gambling` | Gambling and betting sites |
| `ads` | Advertisements and trackers |
| `malware` | Known malware and phishing domains |
| `drugs` | Drug-related content |
| `piracy` | Piracy and torrent sites |
| `social` | Facebook, Instagram, TikTok, Twitter |
| `gaming` | Gaming sites and platforms |

### DNS filter presets

| Preset | Provider | What it filters |
|--------|----------|-----------------|
| `family-cloudflare` | Cloudflare 1.1.1.3 | Malware + adult |
| `family-opendns` | OpenDNS FamilyShield | Adult + phishing |
| `family-cleanbrowsing` | CleanBrowsing | Adult + phishing + mixed |
| `family-adguard` | AdGuard Family | Ads + adult + trackers |
| `safe-google` | Google 8.8.8.8 | SafeSearch enforced |

Or use raw IPs/DoH URLs: `[]string{"1.1.1.3", "https://dns.google/dns-query"}`

### How it works

- **Block categories** → Xray routing rules that send matching traffic to a blackhole outbound. Blocks at the connection level.
- **Block domains** → Same as categories but for custom domain patterns.
- **DNS filtering** → Sets the upstream DNS resolver to a family-safe provider. Blocks at DNS resolution before any connection is made.
- **Combined** → Both layers active. DNS catches broad domains, routing rules catch anything that slips through.

---

## Health Check

Single call to assess full account state — ideal for agents and CLIs:

```go
health, err := client.Health.Get(ctx)

fmt.Println(health.Alerts)                   // e.g. ["Subscription expires in 2 days"]
fmt.Println(health.Hints.SuggestedAction)    // e.g. "renew"
fmt.Println(health.Hints.SubscriptionToken) // ready to pass to Connect
fmt.Println(health.NodesAvailable)          // number of available servers
fmt.Println(health.Usage)                   // today's bytes used/remaining
```

---

## Error Handling

```go
_, err := client.Nodes.List(ctx)
if err != nil {
    switch {
    case vynvpn.IsUnauthorized(err):
        // invalid or expired credentials
    case vynvpn.IsRateLimited(err):
        // back off and retry
    case vynvpn.IsNotFound(err):
        // resource doesn't exist
    case vynvpn.IsServerError(err):
        // 5xx — retry with backoff
    }

    // Access raw details
    if apiErr, ok := err.(*vynvpn.APIError); ok {
        fmt.Printf("HTTP %d: %s\n", apiErr.StatusCode, apiErr.Message)
    }
}
```

---

## Configuration

```go
client := vynvpn.New(
    vynvpn.WithBaseURL("https://api.vynvpn.com"), // default, can omit
    vynvpn.WithAPIKey("vyn_..."),                  // API key auth
    vynvpn.WithToken("existing-jwt"),              // restore saved JWT session
    vynvpn.WithTimeout(60 * time.Second),          // default: 30s
    vynvpn.WithHTTPClient(customClient),           // custom transport, proxy, etc.
)

// Set/update token after construction (e.g. after login)
client.SetToken(auth.Token)
```

---

## Contributing

- **Bug reports / feature requests** → open an issue
- **Code changes** → open a pull request against `next`

Branch structure:
- `next` — active development, all PRs go here
- `main` — stable, merged from `next` when ready
- `release/*` — official rollouts cut from `main`

### Running tests

```bash
# Run all tests
go test ./... -v

# Run a specific test file
go test -run TestConnect ./...

# With race detector
go test -race ./...
```

Tests use `net/http/httptest` — no external dependencies, no network calls, no API key needed. Every service file has a corresponding `_test.go`.

### Pull request guidelines

1. Fork the repo and branch from `next`
2. Write tests for any new or changed behavior
3. Run `go test ./...` and `go vet ./...` — both must pass
4. Keep PRs focused — one feature or fix per PR
5. For significant changes, open an issue first to discuss the approach

### Requesting a feature

Open an issue with:
- What you want to do (use case)
- What endpoint or behavior is missing
- Any API response examples if relevant

---

## Security

Found a vulnerability? Please **do not** open a public issue.

Report to **security@vyntech.com.au** with a description and reproduction steps. We respond within 48 hours.

We run a **bug bounty program** — responsible disclosures are eligible for rewards depending on severity.

---

## License

MIT
