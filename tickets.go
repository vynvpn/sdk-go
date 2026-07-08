package vynvpn

import "context"

// TicketsService handles support ticket operations.
type TicketsService struct {
	client *Client
}

// List returns all tickets for the authenticated user.
// Requires JWT auth.
func (s *TicketsService) List(ctx context.Context) ([]Ticket, error) {
	var tickets []Ticket
	if err := s.client.http.get(ctx, "/api/tickets", nil, &tickets); err != nil {
		return nil, err
	}
	return tickets, nil
}

// Get returns a single ticket with its messages.
// Requires JWT auth.
func (s *TicketsService) Get(ctx context.Context, id string) (*Ticket, error) {
	var ticket Ticket
	if err := s.client.http.get(ctx, "/api/tickets/"+id, nil, &ticket); err != nil {
		return nil, err
	}
	return &ticket, nil
}

// CreateTicketRequest is the body for creating a support ticket.
type CreateTicketRequest struct {
	Subject    string `json:"subject"`
	Message    string `json:"message"`
	Priority   string `json:"priority,omitempty"`   // low, medium, high
	Department string `json:"department,omitempty"` // general, technical, billing, sales
}

// Create opens a new support ticket.
// Requires JWT auth.
func (s *TicketsService) Create(ctx context.Context, req *CreateTicketRequest) (*Ticket, error) {
	var ticket Ticket
	if err := s.client.http.post(ctx, "/api/tickets", req, &ticket); err != nil {
		return nil, err
	}
	return &ticket, nil
}

// ReplyRequest is the body for replying to a ticket.
type ReplyRequest struct {
	Message string `json:"message"`
}

// Reply adds a message to an existing ticket.
// Requires JWT auth.
func (s *TicketsService) Reply(ctx context.Context, ticketID string, message string) error {
	return s.client.http.post(ctx, "/api/tickets/"+ticketID+"/reply",
		&ReplyRequest{Message: message}, nil)
}

// Close closes a ticket.
// Requires JWT auth.
func (s *TicketsService) Close(ctx context.Context, ticketID string) error {
	return s.client.http.patch(ctx, "/api/tickets/"+ticketID+"/close", nil, nil)
}

// SendSupportEmailRequest is the body for the contact-form endpoint.
type SendSupportEmailRequest struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
	Email   string `json:"email,omitempty"` // optional if authenticated
}

// SendSupportEmail sends an email to the support inbox.
// Requires JWT auth.
func (s *TicketsService) SendSupportEmail(ctx context.Context, req *SendSupportEmailRequest) error {
	return s.client.http.post(ctx, "/api/support/email", req, nil)
}
