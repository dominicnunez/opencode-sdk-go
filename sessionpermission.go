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
		params = &SessionPermissionRespondParams{}
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
	p.raw = data
	return nil
}

// AsString returns the pattern as a string if it is a string, or ("", false) if it is an array.
func (p PermissionPattern) AsString() (string, bool) {
	var s string
	if err := json.Unmarshal(p.raw, &s); err != nil {
		return "", false
	}
	return s, true
}

// AsArray returns the pattern as an array of strings if it is an array, or (nil, false) if it is a string.
func (p PermissionPattern) AsArray() ([]string, bool) {
	var arr []string
	if err := json.Unmarshal(p.raw, &arr); err != nil {
		return nil, false
	}
	return arr, true
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
