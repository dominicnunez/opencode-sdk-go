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
		return errors.Is(target, ErrNotFound)
	case e.StatusCode == http.StatusUnauthorized:
		return errors.Is(target, ErrUnauthorized)
	case e.StatusCode == http.StatusForbidden:
		return errors.Is(target, ErrForbidden)
	case e.StatusCode == http.StatusTooManyRequests:
		return errors.Is(target, ErrRateLimited)
	case e.StatusCode >= http.StatusBadRequest && e.StatusCode < http.StatusInternalServerError:
		return errors.Is(target, ErrInvalidRequest)
	case e.StatusCode >= http.StatusInternalServerError:
		return errors.Is(target, ErrInternal)
	}
	return false
}

func IsNotFoundError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return errors.Is(err, ErrNotFound)
}

func IsUnauthorizedError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusUnauthorized
	}
	return errors.Is(err, ErrUnauthorized)
}

func IsForbiddenError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusForbidden
	}
	return errors.Is(err, ErrForbidden)
}

func IsRateLimitedError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusTooManyRequests
	}
	return errors.Is(err, ErrRateLimited)
}

func IsInvalidRequestError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode >= http.StatusBadRequest && apiErr.StatusCode < http.StatusInternalServerError
	}
	return errors.Is(err, ErrInvalidRequest)
}

func IsInternalError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode >= http.StatusInternalServerError
	}
	return errors.Is(err, ErrInternal)
}
