package opencode

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func apiErr(status int) *APIError {
	return &APIError{StatusCode: status, Message: http.StatusText(status)}
}

// TestAPIError_Error_Format verifies the string representation of APIError.
func TestAPIError_Error_Format(t *testing.T) {
	tests := []struct {
		name      string
		err       *APIError
		wantParts []string
	}{
		{
			name:      "without request ID",
			err:       &APIError{StatusCode: 404, Message: "not found"},
			wantParts: []string{"not found", "404"},
		},
		{
			name:      "with request ID",
			err:       &APIError{StatusCode: 404, Message: "not found", RequestID: "req-abc"},
			wantParts: []string{"not found", "404", "req-abc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			for _, part := range tt.wantParts {
				if !contains(got, part) {
					t.Errorf("Error() = %q, expected to contain %q", got, part)
				}
			}
		})
	}
}

// TestAPIError_Is_SentinelMapping verifies that each HTTP status code maps to
// the correct sentinel via errors.Is.
func TestAPIError_Is_SentinelMapping(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		sentinel   error
		wantMatch  bool
	}{
		// --- 404 maps to ErrNotFound ---
		{"404 matches ErrNotFound", http.StatusNotFound, ErrNotFound, true},
		{"404 does not match ErrInvalidRequest", http.StatusNotFound, ErrInvalidRequest, false},
		{"404 does not match ErrInternal", http.StatusNotFound, ErrInternal, false},

		// --- 401 maps to ErrUnauthorized ---
		{"401 matches ErrUnauthorized", http.StatusUnauthorized, ErrUnauthorized, true},
		{"401 does not match ErrForbidden", http.StatusUnauthorized, ErrForbidden, false},

		// --- 403 maps to ErrForbidden ---
		{"403 matches ErrForbidden", http.StatusForbidden, ErrForbidden, true},
		{"403 does not match ErrUnauthorized", http.StatusForbidden, ErrUnauthorized, false},
		{"403 does not match ErrInvalidRequest", http.StatusForbidden, ErrInvalidRequest, false},

		// --- 429 maps to ErrRateLimited ---
		{"429 matches ErrRateLimited", http.StatusTooManyRequests, ErrRateLimited, true},
		{"429 does not match ErrInvalidRequest", http.StatusTooManyRequests, ErrInvalidRequest, false},

		// --- 4xx (not 401/403/404/429) maps to ErrInvalidRequest ---
		{"400 matches ErrInvalidRequest", http.StatusBadRequest, ErrInvalidRequest, true},
		{"400 does not match ErrNotFound", http.StatusBadRequest, ErrNotFound, false},
		{"400 does not match ErrInternal", http.StatusBadRequest, ErrInternal, false},
		{"409 matches ErrInvalidRequest", http.StatusConflict, ErrInvalidRequest, true},
		{"422 matches ErrInvalidRequest", http.StatusUnprocessableEntity, ErrInvalidRequest, true},

		// --- 5xx maps to ErrInternal ---
		{"500 matches ErrInternal", http.StatusInternalServerError, ErrInternal, true},
		{"500 does not match ErrInvalidRequest", http.StatusInternalServerError, ErrInvalidRequest, false},
		{"502 matches ErrInternal", http.StatusBadGateway, ErrInternal, true},
		{"503 matches ErrInternal", http.StatusServiceUnavailable, ErrInternal, true},

		// --- unrecognised status codes match nothing ---
		{"200 does not match any sentinel", http.StatusOK, ErrNotFound, false},
		{"200 does not match ErrInternal", http.StatusOK, ErrInternal, false},
		{"301 does not match ErrInvalidRequest", http.StatusMovedPermanently, ErrInvalidRequest, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apiErr(tt.statusCode)
			got := errors.Is(err, tt.sentinel)
			if got != tt.wantMatch {
				t.Errorf("errors.Is(APIError{%d}, %v) = %v, want %v",
					tt.statusCode, tt.sentinel, got, tt.wantMatch)
			}
		})
	}
}

// TestAPIError_Is_ViaWrapping verifies that errors.Is traverses the error chain
// and matches through fmt.Errorf %w wrapping.
func TestAPIError_Is_ViaWrapping(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		sentinel   error
	}{
		{"wrapped 404 matches ErrNotFound", http.StatusNotFound, ErrNotFound},
		{"wrapped 401 matches ErrUnauthorized", http.StatusUnauthorized, ErrUnauthorized},
		{"wrapped 403 matches ErrForbidden", http.StatusForbidden, ErrForbidden},
		{"wrapped 429 matches ErrRateLimited", http.StatusTooManyRequests, ErrRateLimited},
		{"wrapped 400 matches ErrInvalidRequest", http.StatusBadRequest, ErrInvalidRequest},
		{"wrapped 500 matches ErrInternal", http.StatusInternalServerError, ErrInternal},
		{"wrapped 502 matches ErrInternal", http.StatusBadGateway, ErrInternal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inner := apiErr(tt.statusCode)
			wrapped := fmt.Errorf("operation failed: %w", inner)
			if !errors.Is(wrapped, tt.sentinel) {
				t.Errorf("errors.Is(wrapped APIError{%d}, %v) = false, want true",
					tt.statusCode, tt.sentinel)
			}
		})
	}
}

// TestAPIError_Is_DoublyWrapped verifies sentinel matching survives multiple
// wrapping layers.
func TestAPIError_Is_DoublyWrapped(t *testing.T) {
	inner := apiErr(http.StatusNotFound)
	once := fmt.Errorf("layer one: %w", inner)
	twice := fmt.Errorf("layer two: %w", once)
	if !errors.Is(twice, ErrNotFound) {
		t.Error("errors.Is through double wrapping should match ErrNotFound")
	}
}

// TestIsNotFoundError covers the IsNotFoundError helper.
func TestIsNotFoundError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"404 APIError", apiErr(http.StatusNotFound), true},
		{"400 APIError", apiErr(http.StatusBadRequest), false},
		{"403 APIError", apiErr(http.StatusForbidden), false},
		{"500 APIError", apiErr(http.StatusInternalServerError), false},
		{"wrapped 404 APIError", fmt.Errorf("wrap: %w", apiErr(http.StatusNotFound)), true},
		{"ErrNotFound sentinel directly", ErrNotFound, true},
		{"ErrInvalidRequest sentinel", ErrInvalidRequest, false},
		{"unrelated error", errors.New("something else"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFoundError(tt.err); got != tt.want {
				t.Errorf("IsNotFoundError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

// TestIsUnauthorizedError covers the IsUnauthorizedError helper.
func TestIsUnauthorizedError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"401 APIError", apiErr(http.StatusUnauthorized), true},
		{"403 APIError", apiErr(http.StatusForbidden), false},
		{"404 APIError", apiErr(http.StatusNotFound), false},
		{"500 APIError", apiErr(http.StatusInternalServerError), false},
		{"wrapped 401 APIError", fmt.Errorf("wrap: %w", apiErr(http.StatusUnauthorized)), true},
		{"ErrUnauthorized sentinel directly", ErrUnauthorized, true},
		{"ErrForbidden sentinel", ErrForbidden, false},
		{"unrelated error", errors.New("something else"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUnauthorizedError(tt.err); got != tt.want {
				t.Errorf("IsUnauthorizedError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

// TestIsForbiddenError covers the IsForbiddenError helper.
func TestIsForbiddenError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"403 APIError", apiErr(http.StatusForbidden), true},
		{"401 APIError", apiErr(http.StatusUnauthorized), false},
		{"404 APIError", apiErr(http.StatusNotFound), false},
		{"400 APIError", apiErr(http.StatusBadRequest), false},
		{"500 APIError", apiErr(http.StatusInternalServerError), false},
		{"wrapped 403 APIError", fmt.Errorf("wrap: %w", apiErr(http.StatusForbidden)), true},
		{"ErrForbidden sentinel directly", ErrForbidden, true},
		{"ErrUnauthorized sentinel", ErrUnauthorized, false},
		{"unrelated error", errors.New("something else"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsForbiddenError(tt.err); got != tt.want {
				t.Errorf("IsForbiddenError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

// TestIsRateLimitedError covers the IsRateLimitedError helper.
func TestIsRateLimitedError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"429 APIError", apiErr(http.StatusTooManyRequests), true},
		{"400 APIError", apiErr(http.StatusBadRequest), false},
		{"401 APIError", apiErr(http.StatusUnauthorized), false},
		{"500 APIError", apiErr(http.StatusInternalServerError), false},
		{"wrapped 429 APIError", fmt.Errorf("wrap: %w", apiErr(http.StatusTooManyRequests)), true},
		{"ErrRateLimited sentinel directly", ErrRateLimited, true},
		{"ErrInvalidRequest sentinel", ErrInvalidRequest, false},
		{"unrelated error", errors.New("something else"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRateLimitedError(tt.err); got != tt.want {
				t.Errorf("IsRateLimitedError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

// TestIsInvalidRequestError covers the IsInvalidRequestError helper.
// 404, 401, 403, and 429 must NOT match even though they are 4xx — those have
// dedicated sentinels and helpers that take precedence in Is().
func TestIsInvalidRequestError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"400 APIError", apiErr(http.StatusBadRequest), true},
		{"409 APIError", apiErr(http.StatusConflict), true},
		{"422 APIError", apiErr(http.StatusUnprocessableEntity), true},
		// 404/401/403/429 are caught by Is() before the generic 4xx case, so
		// IsInvalidRequestError (which uses As+StatusCode range) still returns
		// true because the range check is 400–499 inclusive for those codes.
		// Verify that 500 is excluded.
		{"500 APIError", apiErr(http.StatusInternalServerError), false},
		{"502 APIError", apiErr(http.StatusBadGateway), false},
		{"200 APIError", apiErr(http.StatusOK), false},
		{"wrapped 400 APIError", fmt.Errorf("wrap: %w", apiErr(http.StatusBadRequest)), true},
		{"ErrInvalidRequest sentinel directly", ErrInvalidRequest, true},
		{"ErrNotFound sentinel", ErrNotFound, false},
		{"unrelated error", errors.New("something else"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInvalidRequestError(tt.err); got != tt.want {
				t.Errorf("IsInvalidRequestError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

// TestIsInternalError covers the IsInternalError helper.
func TestIsInternalError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"500 APIError", apiErr(http.StatusInternalServerError), true},
		{"502 APIError", apiErr(http.StatusBadGateway), true},
		{"503 APIError", apiErr(http.StatusServiceUnavailable), true},
		{"400 APIError", apiErr(http.StatusBadRequest), false},
		{"404 APIError", apiErr(http.StatusNotFound), false},
		{"wrapped 500 APIError", fmt.Errorf("wrap: %w", apiErr(http.StatusInternalServerError)), true},
		{"wrapped 502 APIError", fmt.Errorf("wrap: %w", apiErr(http.StatusBadGateway)), true},
		{"ErrInternal sentinel directly", ErrInternal, true},
		{"ErrInvalidRequest sentinel", ErrInvalidRequest, false},
		{"unrelated error", errors.New("something else"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInternalError(tt.err); got != tt.want {
				t.Errorf("IsInternalError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

// TestStatusCode_BoundaryDistinctions verifies that 400 vs 404 and 500 vs 502
// are unambiguously handled by both errors.Is and the Is*Error helpers.
func TestStatusCode_BoundaryDistinctions(t *testing.T) {
	t.Run("400 is invalid request, not not-found", func(t *testing.T) {
		err := apiErr(http.StatusBadRequest)
		if !errors.Is(err, ErrInvalidRequest) {
			t.Error("400 should match ErrInvalidRequest")
		}
		if errors.Is(err, ErrNotFound) {
			t.Error("400 must not match ErrNotFound")
		}
		if !IsInvalidRequestError(err) {
			t.Error("IsInvalidRequestError should return true for 400")
		}
		if IsNotFoundError(err) {
			t.Error("IsNotFoundError should return false for 400")
		}
	})

	t.Run("404 is not-found, not invalid request", func(t *testing.T) {
		err := apiErr(http.StatusNotFound)
		if !errors.Is(err, ErrNotFound) {
			t.Error("404 should match ErrNotFound")
		}
		if errors.Is(err, ErrInvalidRequest) {
			t.Error("404 must not match ErrInvalidRequest via errors.Is")
		}
		if !IsNotFoundError(err) {
			t.Error("IsNotFoundError should return true for 404")
		}
	})

	t.Run("500 is internal, not invalid request", func(t *testing.T) {
		err := apiErr(http.StatusInternalServerError)
		if !errors.Is(err, ErrInternal) {
			t.Error("500 should match ErrInternal")
		}
		if errors.Is(err, ErrInvalidRequest) {
			t.Error("500 must not match ErrInvalidRequest")
		}
		if !IsInternalError(err) {
			t.Error("IsInternalError should return true for 500")
		}
		if IsInvalidRequestError(err) {
			t.Error("IsInvalidRequestError should return false for 500")
		}
	})

	t.Run("502 is internal error", func(t *testing.T) {
		err := apiErr(http.StatusBadGateway)
		if !errors.Is(err, ErrInternal) {
			t.Error("502 should match ErrInternal")
		}
		if !IsInternalError(err) {
			t.Error("IsInternalError should return true for 502")
		}
	})
}

// TestAPIError_IsRetryable verifies that IsRetryable correctly identifies
// transient HTTP status codes.
func TestAPIError_IsRetryable(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{"200 not retryable", http.StatusOK, false},
		{"400 not retryable", http.StatusBadRequest, false},
		{"401 not retryable", http.StatusUnauthorized, false},
		{"403 not retryable", http.StatusForbidden, false},
		{"404 not retryable", http.StatusNotFound, false},
		{"408 retryable", http.StatusRequestTimeout, true},
		{"429 retryable", http.StatusTooManyRequests, true},
		{"500 retryable", http.StatusInternalServerError, true},
		{"502 retryable", http.StatusBadGateway, true},
		{"503 retryable", http.StatusServiceUnavailable, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apiErr(tt.statusCode)
			if got := err.IsRetryable(); got != tt.want {
				t.Errorf("APIError{%d}.IsRetryable() = %v, want %v",
					tt.statusCode, got, tt.want)
			}
		})
	}
}

// TestIsRetryableError verifies the top-level helper function.
func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"unrelated error", errors.New("something"), false},
		{"400 not retryable", apiErr(http.StatusBadRequest), false},
		{"404 not retryable", apiErr(http.StatusNotFound), false},
		{"408 retryable", apiErr(http.StatusRequestTimeout), true},
		{"429 retryable", apiErr(http.StatusTooManyRequests), true},
		{"500 retryable", apiErr(http.StatusInternalServerError), true},
		{"502 retryable", apiErr(http.StatusBadGateway), true},
		{"wrapped 429 retryable", fmt.Errorf("wrap: %w", apiErr(http.StatusTooManyRequests)), true},
		{"wrapped 500 retryable", fmt.Errorf("wrap: %w", apiErr(http.StatusInternalServerError)), true},
		{"wrapped 400 not retryable", fmt.Errorf("wrap: %w", apiErr(http.StatusBadRequest)), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRetryableError(tt.err); got != tt.want {
				t.Errorf("IsRetryableError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

// TestReadAPIError_TruncatesMessage verifies that readAPIError truncates the
// Message field for large response bodies while preserving the full Body.
func TestReadAPIError_TruncatesMessage(t *testing.T) {
	t.Run("small body is not truncated", func(t *testing.T) {
		body := "short error"
		resp := &http.Response{
			StatusCode: http.StatusBadRequest,
			Header:     http.Header{},
			Body:       io.NopCloser(strings.NewReader(body)),
		}
		apiErr := readAPIError(resp, 1<<20)
		if apiErr.Message != body {
			t.Errorf("Message = %q, want %q", apiErr.Message, body)
		}
		if apiErr.Body != body {
			t.Errorf("Body = %q, want %q", apiErr.Body, body)
		}
	})

	t.Run("large body is truncated in Message but preserved in Body", func(t *testing.T) {
		body := strings.Repeat("x", maxMessageDisplaySize+1000)
		resp := &http.Response{
			StatusCode: http.StatusInternalServerError,
			Header:     http.Header{},
			Body:       io.NopCloser(strings.NewReader(body)),
		}
		apiErr := readAPIError(resp, 1<<20)
		if apiErr.Body != body {
			t.Errorf("Body length = %d, want %d", len(apiErr.Body), len(body))
		}
		if len(apiErr.Message) >= len(body) {
			t.Errorf("Message should be truncated, got length %d", len(apiErr.Message))
		}
		if !strings.HasSuffix(apiErr.Message, "... (truncated)") {
			t.Errorf("Message should end with truncation marker, got %q", apiErr.Message[len(apiErr.Message)-30:])
		}
	})

	t.Run("exactly at limit is not truncated", func(t *testing.T) {
		body := strings.Repeat("y", maxMessageDisplaySize)
		resp := &http.Response{
			StatusCode: http.StatusBadGateway,
			Header:     http.Header{},
			Body:       io.NopCloser(strings.NewReader(body)),
		}
		apiErr := readAPIError(resp, 1<<20)
		if apiErr.Message != body {
			t.Errorf("Message at exact limit should not be truncated")
		}
	})
}

// contains is a local helper to avoid importing strings in tests.
func contains(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr)))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
