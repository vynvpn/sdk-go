package vynvpn

import (
	"context"
	"testing"
)

func TestTickets_List(t *testing.T) {
	client, req, srv := newTestServer(t, 200, []map[string]any{
		{"id": "00000000-0000-0000-0000-000000000001", "subject": "Help", "status": "open"},
	}, WithToken("jwt"))
	defer srv.Close()

	tickets, err := client.Tickets.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(tickets) != 1 {
		t.Fatalf("got %d tickets, want 1", len(tickets))
	}
	if tickets[0].Subject != "Help" {
		t.Errorf("Subject = %q", tickets[0].Subject)
	}
	if req.Path != "/api/tickets" {
		t.Errorf("path = %q", req.Path)
	}
}

func TestTickets_Create(t *testing.T) {
	client, req, srv := newTestServer(t, 200, map[string]any{
		"id": "00000000-0000-0000-0000-000000000002", "subject": "New Issue",
	}, WithToken("jwt"))
	defer srv.Close()

	ticket, err := client.Tickets.Create(context.Background(), &CreateTicketRequest{
		Subject: "New Issue", Message: "Help me",
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if ticket.Subject != "New Issue" {
		t.Errorf("Subject = %q", ticket.Subject)
	}
	if req.Method != "POST" || req.Path != "/api/tickets" {
		t.Errorf("request = %s %s", req.Method, req.Path)
	}
}
