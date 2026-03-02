package opencode

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrNotFound       = errors.New("resource not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrRateLimited    = errors.New("rate limited")
	ErrInvalidRequest = errors.New("invalid request")
	// ErrMissingRequiredParameter matches input validation failures where a
	// required API parameter is not provided.
	ErrMissingRequiredParameter = errors.New("missing required parameter")
	// ErrParamsRequired matches validation failures where a request params
	// object is required but nil was passed.
	ErrParamsRequired = errors.New("params is required")
	// ErrRequiredField matches validation failures where a required field on a
	// provided struct is empty.
	ErrRequiredField = errors.New("required field missing")
	// ErrContextRequired is returned when a nil context.Context is passed to
	// an API method that performs HTTP requests.
	ErrContextRequired = errors.New("context is required")
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

const maxAPIErrorMessageLength = 512

type MissingRequiredParameterError struct {
	Parameter string
}

func (e *MissingRequiredParameterError) Error() string {
	return fmt.Sprintf("missing required %s parameter", e.Parameter)
}

func (e *MissingRequiredParameterError) Is(target error) bool {
	if target == ErrMissingRequiredParameter {
		return true
	}
	t, ok := target.(*MissingRequiredParameterError)
	if !ok {
		return false
	}
	return t.Parameter == "" || t.Parameter == e.Parameter
}

func missingRequiredParameterError(parameter string) error {
	return &MissingRequiredParameterError{Parameter: parameter}
}

type RequiredFieldError struct {
	Field string
}

func (e *RequiredFieldError) Error() string {
	return fmt.Sprintf("%s is required", e.Field)
}

func (e *RequiredFieldError) Is(target error) bool {
	if target == ErrRequiredField {
		return true
	}
	t, ok := target.(*RequiredFieldError)
	if !ok {
		return false
	}
	return t.Field == "" || t.Field == e.Field
}

func requiredFieldError(field string) error {
	return &RequiredFieldError{Field: field}
}

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
	msg := apiErrorMessage(resp.StatusCode, body)

	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    msg,
		RequestID:  resp.Header.Get("X-Request-Id"),
		Body:       body,
		Truncated:  truncated,
		ReadErr:    readErr,
	}
}

func apiErrorMessage(statusCode int, body string) string {
	if candidate := apiErrorMessageFromBody(body); candidate != "" {
		return candidate
	}
	return apiErrorStatusText(statusCode)
}

func apiErrorStatusText(statusCode int) string {
	msg := http.StatusText(statusCode)
	if msg == "" {
		return fmt.Sprintf("http %d", statusCode)
	}
	return msg
}

func apiErrorMessageFromBody(body string) string {
	trimmed := strings.TrimSpace(body)
	if trimmed == "" {
		return ""
	}

	if fromJSON := apiErrorMessageFromJSON(trimmed); fromJSON != "" {
		return sanitizeAPIErrorMessage(fromJSON)
	}

	return sanitizeAPIErrorMessage(trimmed)
}

func apiErrorMessageFromJSON(raw string) string {
	var payload any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return ""
	}
	return findMessageString(payload)
}

func findMessageString(value any) string {
	switch v := value.(type) {
	case map[string]any:
		keys := []string{"message", "error", "detail", "title", "reason", "description"}
		checkedKeys := make(map[string]struct{}, len(keys))
		for _, key := range keys {
			checkedKeys[key] = struct{}{}
			if field, ok := v[key]; ok {
				if msg := findMessageString(field); msg != "" {
					return msg
				}
			}
		}

		fallbackKeys := make([]string, 0, len(v))
		for key := range v {
			if _, isPreferred := checkedKeys[key]; !isPreferred {
				fallbackKeys = append(fallbackKeys, key)
			}
		}
		sort.Strings(fallbackKeys)

		for _, key := range fallbackKeys {
			if msg := findMessageString(v[key]); msg != "" {
				return msg
			}
		}
	case []any:
		for _, item := range v {
			if msg := findMessageString(item); msg != "" {
				return msg
			}
		}
	case string:
		return strings.TrimSpace(v)
	}

	return ""
}

func sanitizeAPIErrorMessage(msg string) string {
	if msg == "" {
		return ""
	}

	var cleaned strings.Builder
	for _, r := range msg {
		if !unicode.IsPrint(r) {
			continue
		}
		if unicode.IsSpace(r) {
			cleaned.WriteByte(' ')
			continue
		}
		cleaned.WriteRune(r)
	}

	normalized := strings.Join(strings.Fields(cleaned.String()), " ")
	if normalized == "" {
		return ""
	}
	if utf8.RuneCountInString(normalized) <= maxAPIErrorMessageLength {
		return normalized
	}

	runes := []rune(normalized)
	return string(runes[:maxAPIErrorMessageLength-3]) + "..."
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
