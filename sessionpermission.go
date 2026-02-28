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

type SessionPermissionService struct {
	client *Client
}

func (s *SessionPermissionService) Respond(ctx context.Context, id string, permissionID string, params *SessionPermissionRespondParams) (bool, error) {
	if id == "" {
		return false, errors.New("missing required id parameter")
	}
	if permissionID == "" {
		return false, errors.New("missing required permissionID parameter")
	}
	if params == nil {
		return false, errors.New("params is required")
	}
	var result bool
	err := s.client.do(ctx, http.MethodPost, fmt.Sprintf("session/%s/permissions/%s", id, permissionID), params, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

type PermissionResponse string

const (
	PermissionResponseOnce   PermissionResponse = "once"
	PermissionResponseAlways PermissionResponse = "always"
	PermissionResponseReject PermissionResponse = "reject"
)

func (r PermissionResponse) IsKnown() bool {
	switch r {
	case PermissionResponseOnce, PermissionResponseAlways, PermissionResponseReject:
		return true
	}
	return false
}

type Permission struct {
	ID        string                 `json:"id"`
	MessageID string                 `json:"messageID"`
	Metadata  map[string]interface{} `json:"metadata"`
	SessionID string                 `json:"sessionID"`
	Time      PermissionTime         `json:"time"`
	Title     string                 `json:"title"`
	Type      string                 `json:"type"`
	CallID    string                 `json:"callID,omitempty"`
	Pattern   *PermissionPattern     `json:"pattern,omitempty"`
}

type PermissionTime struct {
	Created float64 `json:"created"`
}

// PermissionPattern can be either a string or an array of strings.
// Use AsString() or AsArray() to access the value.
type PermissionPattern struct {
	raw json.RawMessage
}

func (p *PermissionPattern) UnmarshalJSON(data []byte) error {
	if !json.Valid(data) {
		return fmt.Errorf("invalid JSON for PermissionPattern")
	}
	p.raw = append(json.RawMessage(nil), data...)
	return nil
}

// AsString returns the pattern as a string if it is a string, or ("", error) if it is an array.
func (p PermissionPattern) AsString() (string, error) {
	if p.raw == nil {
		return "", ErrWrongVariant
	}
	var s string
	if err := json.Unmarshal(p.raw, &s); err != nil {
		return "", err
	}
	return s, nil
}

// AsArray returns the pattern as an array of strings if it is an array, or (nil, error) if it is a string.
func (p PermissionPattern) AsArray() ([]string, error) {
	if p.raw == nil {
		return nil, ErrWrongVariant
	}
	var arr []string
	if err := json.Unmarshal(p.raw, &arr); err != nil {
		return nil, err
	}
	return arr, nil
}

func (p PermissionPattern) MarshalJSON() ([]byte, error) {
	if p.raw == nil {
		return []byte("null"), nil
	}
	return p.raw, nil
}

type SessionPermissionRespondParams struct {
	Response  PermissionResponse `json:"response"`
	Directory *string            `query:"directory,omitempty"`
}

func (r SessionPermissionRespondParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

var _ json.Marshaler = (*SessionPermissionRespondParams)(nil)

func (r SessionPermissionRespondParams) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Response PermissionResponse `json:"response"`
	}{Response: r.Response})
}
