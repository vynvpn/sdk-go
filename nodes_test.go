package vynvpn

import (
	"context"
	"testing"
)

func TestNodes_List(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"data": []map[string]any{
			{"location_slug": "de001", "label": "Germany 1", "country": "Germany", "available": true},
			{"location_slug": "us001", "label": "US East", "country": "US", "available": false},
		},
	})
	defer srv.Close()

	nodes, err := client.Nodes.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(nodes) != 2 {
		t.Fatalf("got %d nodes, want 2", len(nodes))
	}
	if nodes[0].LocationSlug != "de001" {
		t.Errorf("nodes[0].LocationSlug = %q, want %q", nodes[0].LocationSlug, "de001")
	}
	if nodes[0].Available != true {
		t.Error("nodes[0] should be available")
	}
	if req.Path != "/v2/nodes" {
		t.Errorf("path = %q, want /v2/nodes", req.Path)
	}
}

func TestNodes_ListWithAPIKey(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"data": []map[string]any{
			{"location_slug": "de001", "country": "Germany", "available": true},
		},
	}, WithAPIKey("vyn_test"))
	defer srv.Close()

	nodes, err := client.Nodes.ListWithAPIKey(context.Background())
	if err != nil {
		t.Fatalf("ListWithAPIKey failed: %v", err)
	}
	if len(nodes) != 1 {
		t.Fatalf("got %d nodes, want 1", len(nodes))
	}
	if req.Path != "/api/v1/nodes" {
		t.Errorf("path = %q, want /api/v1/nodes", req.Path)
	}
}
