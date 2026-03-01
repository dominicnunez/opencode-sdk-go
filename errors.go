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
	ErrInternal       = errors.New("internal server error")
	ErrWrongVariant   = errors.New("wrong union variant")
)

type APIError struct {
	StatusCode int
	Message    string
	RequestID  string
	Body       string
}

func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("%s (status %d, request %s)", e.Message, e.StatusCode, e.RequestID)
	}
	return fmt.Sprintf("%s (status %d)", e.Message, e.StatusCode)
}

// IsRetryable reports whether the error represents a transient failure
// (408 Request Timeout, 429 Too Many Requests, or 5xx) that may succeed
// on retry. Callers implementing SSE reconnection or custom retry logic
// can use this to distinguish transient from permanent failures.
func (e *APIError) IsRetryable() bool {
	return e.StatusCode == http.StatusRequestTimeout ||
		e.StatusCode == http.StatusTooManyRequests ||
		e.StatusCode >= http.StatusInternalServerError
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
	case e.StatusCode >= http.StatusBadRequest && e.StatusCode < http.StatusInternalServerError:
		return target == ErrInvalidRequest
	case e.StatusCode >= http.StatusInternalServerError:
		return target == ErrInternal
	}
	return false
}

// IsRetryableError reports whether err wraps an *APIError with a retryable
// status code (408, 429, or 5xx).
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
	bodyBytes, readErr := io.ReadAll(io.LimitReader(resp.Body, bodyLimit))
	_ = resp.Body.Close()

	body := string(bodyBytes)
	if body == "" {
		body = http.StatusText(resp.StatusCode)
	}
	if readErr != nil {
		body += fmt.Sprintf(" (read error: %v)", readErr)
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
	}
}

func IsNotFoundError(err error) bool       { return errors.Is(err, ErrNotFound) }
func IsUnauthorizedError(err error) bool   { return errors.Is(err, ErrUnauthorized) }
func IsForbiddenError(err error) bool      { return errors.Is(err, ErrForbidden) }
func IsRateLimitedError(err error) bool    { return errors.Is(err, ErrRateLimited) }
func IsInvalidRequestError(err error) bool { return errors.Is(err, ErrInvalidRequest) }
func IsInternalError(err error) bool       { return errors.Is(err, ErrInternal) }
