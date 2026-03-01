package opencode

import (
	"context"
	"encoding/json"
	"errors"
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
		return false, errors.New("missing required id parameter")
	}
	if params == nil {
		return false, errors.New("params is required")
	}
	if params.Auth == nil {
		return false, errors.New("missing required Auth field")
	}

	var result bool
	err := s.client.do(ctx, http.MethodPut, "auth/"+id, params, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

// AuthSetParamsAuthUnion is satisfied by [OAuth], [ApiAuth], and [WellKnownAuth].
type AuthSetParamsAuthUnion interface {
	implementsAuthSetParamsAuthUnion()
}

type AuthSetParams struct {
	Auth      AuthSetParamsAuthUnion `json:"-"`
	Directory *string                `json:"-" query:"directory,omitempty"`
}

// URLQuery serializes [AuthSetParams]'s query parameters as `url.Values`.
func (r AuthSetParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

// MarshalJSON marshals the Auth field for the request body.
// It sets the Type discriminator automatically based on the concrete type,
// so callers don't need to set it manually.
func (r AuthSetParams) MarshalJSON() ([]byte, error) {
	switch v := r.Auth.(type) {
	case OAuth:
		v.Type = AuthTypeOAuth
		return json.Marshal(v)
	case *OAuth:
		v.Type = AuthTypeOAuth
		return json.Marshal(v)
	case ApiAuth:
		v.Type = AuthTypeAPI
		return json.Marshal(v)
	case *ApiAuth:
		v.Type = AuthTypeAPI
		return json.Marshal(v)
	case WellKnownAuth:
		v.Type = AuthTypeWellKnown
		return json.Marshal(v)
	case *WellKnownAuth:
		v.Type = AuthTypeWellKnown
		return json.Marshal(v)
	default:
		return json.Marshal(r.Auth)
	}
}
