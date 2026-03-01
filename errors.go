package opencode

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrNotFound       = errors.New("resource not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrRateLimited    = errors.New("rate limited")
	ErrInvalidRequest = errors.New("invalid request")
	// ErrTimeout matches HTTP 408 (Request Timeout) responses from the server.
	// Client-side timeouts from context.WithTimeout or context.WithDeadline
	// surface as context.DeadlineExceeded and are NOT matched by this sentinel.
	// Use errors.Is(err, context.DeadlineExceeded) to detect client-side timeouts.
	ErrTimeout  = errors.New("request timeout")
	ErrInternal = errors.New("internal server error")
	// ErrWrongVariant is returned when a union type accessor is called with
	// a discriminator value that does not match the requested variant.
	ErrWrongVariant = errors.New("wrong union variant")

	// ErrNilAuth is returned when AuthSetParams.MarshalJSON is called with a nil
	// Auth field or a non-nil interface holding a nil pointer.
	ErrNilAuth = errors.New("nil auth value")
	// ErrUnknownAuthType is returned when AuthSetParams.MarshalJSON encounters
	// an Auth implementation that is not one of OAuth, ApiAuth, or WellKnownAuth.
	ErrUnknownAuthType = errors.New("unknown auth union type")
)

func wrongVariant(expected, actual string) error {
	return fmt.Errorf("%s, got %s: %w", expected, actual, ErrWrongVariant)
}

type APIError struct {
	StatusCode int
	Message    string
	RequestID  string
	Body       string
	// Truncated is true when the response body exceeded the read limit
	// and Body contains only the first portion of the original response.
	Truncated bool
	// ReadErr is non-nil when the response body could not be fully read
	// (e.g. connection dropped mid-transfer). Body contains whatever
	// partial data was received before the error.
	ReadErr error
}

func (e *APIError) Error() string {
	var msg string
	if e.RequestID != "" {
		msg = fmt.Sprintf("%s (status %d, request %s)", e.Message, e.StatusCode, e.RequestID)
	} else {
		msg = fmt.Sprintf("%s (status %d)", e.Message, e.StatusCode)
	}
	if e.ReadErr != nil {
		msg += fmt.Sprintf(" (body read error: %v)", e.ReadErr)
	}
	return msg
}

func isRetryableStatus(code int) bool {
	return code == http.StatusRequestTimeout ||
		code == http.StatusTooManyRequests ||
		code >= http.StatusInternalServerError
}

// IsRetryable reports whether the error represents a transient failure
// (408 Request Timeout, 429 Too Many Requests, or 5xx) that may succeed
// on retry. Callers implementing SSE reconnection or custom retry logic
// can use this to distinguish transient from permanent failures.
func (e *APIError) IsRetryable() bool {
	return isRetryableStatus(e.StatusCode)
}

func (e *APIError) Is(target error) bool {
	switch {
	case e.StatusCode == http.StatusNotFound:
		return target == ErrNotFound
	case e.StatusCode == http.StatusUnauthorized:
		return target == ErrUnauthorized
	case e.StatusCode == http.StatusForbidden:
		return target == ErrForbidden
	case e.StatusCode == http.StatusTooManyRequests:
		return target == ErrRateLimited
	case e.StatusCode == http.StatusRequestTimeout:
		return target == ErrTimeout
	case e.StatusCode >= http.StatusBadRequest && e.StatusCode < http.StatusInternalServerError:
		return target == ErrInvalidRequest
	case e.StatusCode >= http.StatusInternalServerError:
		return target == ErrInternal
	}
	return false
}

// IsRetryableError reports whether err wraps an *APIError with a retryable
// status code (408, 429, or 5xx).
//
// Transport-level failures (DNS resolution, connection refused, TLS handshake
// errors, etc.) are NOT wrapped as *APIError, so this function returns false
// for them. To detect transport errors after retries are exhausted, unwrap the
// underlying net.Error or check for specific error types directly.
func IsRetryableError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.IsRetryable()
	}
	return false
}

// maxMessageDisplaySize caps the Message field to avoid large error bodies
// appearing in log lines and error wrapping chains. The full body is always
// available in the Body field.
const maxMessageDisplaySize = 4096

// readAPIError reads the response body (up to limit bytes), constructs an
// *APIError, and closes the body. The caller should not use resp.Body after.
func readAPIError(resp *http.Response, bodyLimit int64) *APIError {
	// Read one extra byte beyond the limit to detect truncation.
	bodyBytes, readErr := io.ReadAll(io.LimitReader(resp.Body, bodyLimit+1))
	_ = resp.Body.Close()

	truncated := int64(len(bodyBytes)) > bodyLimit
	if truncated {
		bodyBytes = bodyBytes[:bodyLimit]
	}

	body := string(bodyBytes)
	if body == "" {
		body = http.StatusText(resp.StatusCode)
	}

	msg := body
	if len(msg) > maxMessageDisplaySize {
		msg = msg[:maxMessageDisplaySize] + "... (truncated)"
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    msg,
		RequestID:  resp.Header.Get("X-Request-Id"),
		Body:       body,
		Truncated:  truncated,
		ReadErr:    readErr,
	}
}

// IsTimeoutError reports whether err matches an HTTP 408 (Request Timeout)
// response. It does not match client-side timeouts caused by context
// cancellation or deadlines — use errors.Is(err, context.DeadlineExceeded)
// for those.
//
// All Is*Error helpers only match errors that wrap *APIError (HTTP responses).
// Transport-level failures (DNS, connection refused, TLS errors) are returned
// as plain errors and will not match any of these helpers. Use errors.As with
// net.Error or unwrap the error directly to classify transport failures.
func IsTimeoutError(err error) bool        { return errors.Is(err, ErrTimeout) }
func IsNotFoundError(err error) bool       { return errors.Is(err, ErrNotFound) }
func IsUnauthorizedError(err error) bool   { return errors.Is(err, ErrUnauthorized) }
func IsForbiddenError(err error) bool      { return errors.Is(err, ErrForbidden) }
func IsRateLimitedError(err error) bool    { return errors.Is(err, ErrRateLimited) }
func IsInvalidRequestError(err error) bool { return errors.Is(err, ErrInvalidRequest) }
func IsInternalError(err error) bool       { return errors.Is(err, ErrInternal) }
