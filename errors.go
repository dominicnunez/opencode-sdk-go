package opencode

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound       = errors.New("resource not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrRateLimited    = errors.New("rate limited")
	ErrInvalidRequest = errors.New("invalid request")
	ErrInternal       = errors.New("internal server error")
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
	case e.StatusCode == 404:
		return errors.Is(target, ErrNotFound)
	case e.StatusCode == 401 || e.StatusCode == 403:
		return errors.Is(target, ErrUnauthorized)
	case e.StatusCode == 429:
		return errors.Is(target, ErrRateLimited)
	case e.StatusCode >= 400 && e.StatusCode < 500:
		return errors.Is(target, ErrInvalidRequest)
	case e.StatusCode >= 500:
		return errors.Is(target, ErrInternal)
	}
	return false
}

func IsNotFoundError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 404
	}
	return errors.Is(err, ErrNotFound)
}

func IsUnauthorizedError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 401 || apiErr.StatusCode == 403
	}
	return errors.Is(err, ErrUnauthorized)
}

func IsRateLimitedError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 429
	}
	return errors.Is(err, ErrRateLimited)
}
