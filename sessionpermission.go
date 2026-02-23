package opencode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/dominicnunez/opencode-sdk-go/internal/apijson"
	"github.com/dominicnunez/opencode-sdk-go/internal/apiquery"
	"github.com/dominicnunez/opencode-sdk-go/shared"
	"github.com/tidwall/gjson"
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
	Pattern   PermissionPatternUnion `json:"pattern,omitempty"`
}

func (r *Permission) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type PermissionTime struct {
	Created float64 `json:"created"`
}

type PermissionPatternUnion interface {
	ImplementsPermissionPatternUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*PermissionPatternUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(shared.UnionString("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(PermissionPatternArray{}),
		},
	)
}

type PermissionPatternArray []string

func (r PermissionPatternArray) ImplementsPermissionPatternUnion() {}

type SessionPermissionRespondParams struct {
	Response  PermissionResponse `json:"response"`
	Directory *string            `query:"directory,omitempty"`
}

func (r SessionPermissionRespondParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

var _ json.Marshaler = (*SessionPermissionRespondParams)(nil)

func (r SessionPermissionRespondParams) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Response PermissionResponse `json:"response"`
	}{Response: r.Response})
}
