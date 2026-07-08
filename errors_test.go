package vynvpn

import (
	"errors"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	e := &APIError{StatusCode: 404, Message: "not found"}
	got := e.Error()
	want := "vynvpn: HTTP 404: not found"
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestIsHelpers(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		isNot  func(error) bool
		isYes  func(error) bool
		status int
	}{
		{"not found", &APIError{StatusCode: 404}, IsUnauthorized, IsNotFound, 404},
		{"unauthorized", &APIError{StatusCode: 401}, IsNotFound, IsUnauthorized, 401},
		{"forbidden", &APIError{StatusCode: 403}, IsUnauthorized, IsForbidden, 403},
		{"conflict", &APIError{StatusCode: 409}, IsNotFound, IsConflict, 409},
		{"rate limited", &APIError{StatusCode: 429}, IsNotFound, IsRateLimited, 429},
		{"server error", &APIError{StatusCode: 500}, IsNotFound, IsServerError, 500},
		{"server error 503", &APIError{StatusCode: 503}, IsNotFound, IsServerError, 503},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.isYes(tc.err) {
				t.Errorf("expected true for status %d", tc.status)
			}
			if tc.isNot(tc.err) {
				t.Errorf("expected false for status %d", tc.status)
			}
		})
	}

	// Non-APIError should always return false
	plain := errors.New("plain error")
	if IsNotFound(plain) || IsUnauthorized(plain) || IsForbidden(plain) ||
		IsConflict(plain) || IsRateLimited(plain) || IsServerError(plain) {
		t.Error("non-APIError should return false for all helpers")
	}
}
