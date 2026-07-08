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
| `client.Connect` | Connect, disconnect, status polling |
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
