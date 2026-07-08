package vynvpn

import "context"

// PaymentsService handles payment listing and creation.
type PaymentsService struct {
	client *Client
}

// List returns the user's payment history.
// Routes to /api/v1/payments (API key) or /api/payments (JWT).
func (s *PaymentsService) List(ctx context.Context) ([]Payment, error) {
	if s.client.http.IsAPIKeyAuth() {
		var resp struct {
			Data []Payment `json:"data"`
		}
		if err := s.client.http.get(ctx, "/api/v1/payments", nil, &resp); err != nil {
			return nil, err
		}
		return resp.Data, nil
	}
	var payments []Payment
	if err := s.client.http.get(ctx, "/api/payments", nil, &payments); err != nil {
		return nil, err
	}
	return payments, nil
}

// Get returns a single payment by ID. Requires JWT auth.
func (s *PaymentsService) Get(ctx context.Context, id string) (*Payment, error) {
	var payment Payment
	if err := s.client.http.get(ctx, "/api/payments/"+id, nil, &payment); err != nil {
		return nil, err
	}
	return &payment, nil
}

// CreatePaymentRequest is the body for creating a new payment.
type CreatePaymentRequest struct {
	PlanID        string `json:"plan_id"`
	PaymentMethod string `json:"payment_method,omitempty"`
}

// Create initiates a new payment. Requires JWT auth.
func (s *PaymentsService) Create(ctx context.Context, req *CreatePaymentRequest) (*Payment, error) {
	var payment Payment
	if err := s.client.http.post(ctx, "/api/payments", req, &payment); err != nil {
		return nil, err
	}
	return &payment, nil
}

// ApplyDiscountRequest is the body for applying a discount code to a payment.
type ApplyDiscountRequest struct {
	Code string `json:"code"`
}

// ApplyDiscount applies a discount code to a pending payment. Requires JWT auth.
func (s *PaymentsService) ApplyDiscount(ctx context.Context, paymentID string, code string) error {
	return s.client.http.post(ctx, "/api/payments/"+paymentID+"/apply-discount",
		&ApplyDiscountRequest{Code: code}, nil)
}

// GetMethods returns available payment methods. Requires JWT auth.
func (s *PaymentsService) GetMethods(ctx context.Context) ([]string, error) {
	var methods []string
	if err := s.client.http.get(ctx, "/api/payments/methods", nil, &methods); err != nil {
		return nil, err
	}
	return methods, nil
}
