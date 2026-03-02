package opencode

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dominicnunez/opencode-sdk-go/internal/queryparams"
)

type AuthService struct {
	client *Client
}

// Set configures authentication credentials for a provider
// Endpoint: PUT /auth/{id}
func (s *AuthService) Set(ctx context.Context, id string, params *AuthSetParams) (bool, error) {
	if strings.TrimSpace(id) == "" {
		return false, missingRequiredParameterError("id")
	}
	if params == nil {
		return false, ErrParamsRequired
	}
	if params.Auth == nil {
		return false, requiredFieldError("AuthSetParams: Auth field")
	}
	if err := validateAuthCredentials(params.Auth); err != nil {
		return false, err
	}

	var result bool
	err := s.client.do(ctx, http.MethodPut, "auth/"+url.PathEscape(id), params, &result)
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
	auth, err := normalizeAuthSetParamsAuth(r.Auth)
	if err != nil {
		return nil, err
	}
	if err := validateAuthCredentials(auth); err != nil {
		return nil, err
	}

	switch v := auth.(type) {
	case OAuth:
		v.Type = AuthTypeOAuth
		return json.Marshal(v)
	case ApiAuth:
		v.Type = AuthTypeAPI
		return json.Marshal(v)
	case WellKnownAuth:
		v.Type = AuthTypeWellKnown
		return json.Marshal(v)
	default:
		return nil, fmt.Errorf("auth type %T: %w", r.Auth, ErrUnknownAuthType)
	}
}

func normalizeAuthSetParamsAuth(auth AuthSetParamsAuthUnion) (AuthSetParamsAuthUnion, error) {
	if auth == nil {
		return nil, fmt.Errorf("AuthSetParams: Auth field is required: %w", ErrNilAuth)
	}

	switch v := auth.(type) {
	case *OAuth:
		if v == nil {
			return nil, fmt.Errorf("AuthSetParams: Auth contains typed nil *OAuth: %w", ErrNilAuth)
		}
		copy := *v
		return copy, nil
	case *ApiAuth:
		if v == nil {
			return nil, fmt.Errorf("AuthSetParams: Auth contains typed nil *ApiAuth: %w", ErrNilAuth)
		}
		copy := *v
		return copy, nil
	case *WellKnownAuth:
		if v == nil {
			return nil, fmt.Errorf("AuthSetParams: Auth contains typed nil *WellKnownAuth: %w", ErrNilAuth)
		}
		copy := *v
		return copy, nil
	default:
		return auth, nil
	}
}

func validateAuthCredentials(auth AuthSetParamsAuthUnion) error {
	normalized, err := normalizeAuthSetParamsAuth(auth)
	if err != nil {
		return err
	}

	switch v := normalized.(type) {
	case OAuth:
		if strings.TrimSpace(v.Access) == "" {
			return requiredFieldError("oauth access")
		}
		if strings.TrimSpace(v.Refresh) == "" {
			return requiredFieldError("oauth refresh")
		}
	case ApiAuth:
		if strings.TrimSpace(v.Key) == "" {
			return requiredFieldError("api auth key")
		}
	case WellKnownAuth:
		if strings.TrimSpace(v.Key) == "" {
			return requiredFieldError("well-known auth key")
		}
		if strings.TrimSpace(v.Token) == "" {
			return requiredFieldError("well-known auth token")
		}
	default:
		return fmt.Errorf("auth type %T: %w", normalized, ErrUnknownAuthType)
	}

	return nil
}
