package opencode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dominicnunez/opencode-sdk-go/internal/queryparams"
)

type AuthService struct {
	client *Client
}

// Set configures authentication credentials for a provider
// Endpoint: PUT /auth/{id}
func (s *AuthService) Set(ctx context.Context, id string, params *AuthSetParams) (bool, error) {
	if id == "" {
		return false, errors.New("id is required")
	}
	if params == nil {
		return false, errors.New("params is required")
	}

	path := fmt.Sprintf("auth/%s", id)
	var result bool
	err := s.client.do(ctx, http.MethodPut, path, params, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

type AuthSetParams struct {
	// Auth union: one of OAuth, ApiAuth, or WellKnownAuth
	// Use the concrete types when calling Set
	Auth      interface{} `json:"-"`
	Directory *string     `query:"directory,omitempty"`
}

// URLQuery serializes [AuthSetParams]'s query parameters as `url.Values`.
func (r AuthSetParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

// MarshalJSON marshals the Auth field for the request body
func (r AuthSetParams) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Auth)
}
