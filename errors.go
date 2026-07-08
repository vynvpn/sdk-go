package vynvpn

import "fmt"

// APIError represents an error response from the VynVPN API.
type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	RawBody    string `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("vynvpn: HTTP %d: %s", e.StatusCode, e.Message)
}

// IsNotFound returns true if the error is a 404 Not Found.
func IsNotFound(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.StatusCode == 404
	}
	return false
}

// IsUnauthorized returns true if the error is a 401 Unauthorized.
func IsUnauthorized(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.StatusCode == 401
	}
	return false
}

// IsForbidden returns true if the error is a 403 Forbidden.
func IsForbidden(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.StatusCode == 403
	}
	return false
}

// IsConflict returns true if the error is a 409 Conflict.
func IsConflict(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.StatusCode == 409
	}
	return false
}

// IsRateLimited returns true if the error is a 429 Too Many Requests.
func IsRateLimited(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.StatusCode == 429
	}
	return false
}

// IsServerError returns true if the error is a 5xx server error.
func IsServerError(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.StatusCode >= 500
	}
	return false
}
