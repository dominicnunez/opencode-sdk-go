package opencode

import (
	"errors"
	"fmt"
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

func IsNotFoundError(err error) bool       { return errors.Is(err, ErrNotFound) }
func IsUnauthorizedError(err error) bool   { return errors.Is(err, ErrUnauthorized) }
func IsForbiddenError(err error) bool      { return errors.Is(err, ErrForbidden) }
func IsRateLimitedError(err error) bool    { return errors.Is(err, ErrRateLimited) }
func IsInvalidRequestError(err error) bool { return errors.Is(err, ErrInvalidRequest) }
func IsInternalError(err error) bool       { return errors.Is(err, ErrInternal) }
