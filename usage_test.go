package vynvpn

import (
	"context"
	"testing"
)

func TestUsage_GetStatus(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"date_utc":          "2026-07-08",
		"bytes_used_today":  512000000,
		"daily_limit_bytes": 1073741824,
		"bytes_remaining":   561741824,
		"limit_exceeded":    false,
		"has_premium":       false,
	}, WithToken("jwt"))
	defer srv.Close()

	status, err := client.Usage.GetStatus(context.Background())
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.BytesUsedToday != 512000000 {
		t.Errorf("BytesUsedToday = %d", status.BytesUsedToday)
	}
	if status.LimitExceeded {
		t.Error("expected limit_exceeded = false")
	}
	if req.Path != "/v2/usage" {
		t.Errorf("path = %q, want /v2/usage", req.Path)
	}
}

func TestUsage_Report(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"bytes_used_today": 600000000,
		"limit_exceeded":   false,
	}, WithToken("jwt"))
	defer srv.Close()

	status, err := client.Usage.Report(context.Background(), &UsageReport{
		BytesDelta:   100000000,
		LocationSlug: "de001",
	})
	if err != nil {
		t.Fatalf("Report failed: %v", err)
	}
	if status.BytesUsedToday != 600000000 {
		t.Errorf("BytesUsedToday = %d", status.BytesUsedToday)
	}
	if req.Method != "POST" || req.Path != "/v2/usage" {
		t.Errorf("request = %s %s", req.Method, req.Path)
	}
}
