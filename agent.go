package opencode

import (
	"context"
	"net/http"
	"net/url"

	"github.com/dominicnunez/opencode-sdk-go/internal/apijson"
	"github.com/dominicnunez/opencode-sdk-go/internal/apiquery"
)

type AgentService struct {
	client *Client
}

func (s *AgentService) List(ctx context.Context, params *AgentListParams) ([]Agent, error) {
	if params == nil {
		params = &AgentListParams{}
	}
	var result []Agent
	err := s.client.do(ctx, http.MethodGet, "agent", params, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type Agent struct {
	BuiltIn     bool                   `json:"builtIn"`
	Mode        AgentMode              `json:"mode"`
	Name        string                 `json:"name"`
	Options     map[string]interface{} `json:"options"`
	Permission  AgentPermission        `json:"permission"`
	Tools       map[string]bool        `json:"tools"`
	Description string                 `json:"description,omitempty"`
	Model       *AgentModel            `json:"model,omitempty"`
	Prompt      string                 `json:"prompt,omitempty"`
	Temperature float64                `json:"temperature,omitempty"`
	TopP        float64                `json:"topP,omitempty"`
}

func (r *Agent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

type AgentMode string

const (
	AgentModeSubagent AgentMode = "subagent"
	AgentModePrimary  AgentMode = "primary"
	AgentModeAll      AgentMode = "all"
)

func (r AgentMode) IsKnown() bool {
	switch r {
	case AgentModeSubagent, AgentModePrimary, AgentModeAll:
		return true
	}
	return false
}

type AgentPermission struct {
	Bash     map[string]AgentPermissionBash `json:"bash"`
	Edit     AgentPermissionEdit            `json:"edit"`
	Webfetch *AgentPermissionWebfetch       `json:"webfetch,omitempty"`
}

func (r *AgentPermission) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

type AgentPermissionBash string

const (
	AgentPermissionBashAsk   AgentPermissionBash = "ask"
	AgentPermissionBashAllow AgentPermissionBash = "allow"
	AgentPermissionBashDeny  AgentPermissionBash = "deny"
)

func (r AgentPermissionBash) IsKnown() bool {
	switch r {
	case AgentPermissionBashAsk, AgentPermissionBashAllow, AgentPermissionBashDeny:
		return true
	}
	return false
}

type AgentPermissionEdit string

const (
	AgentPermissionEditAsk   AgentPermissionEdit = "ask"
	AgentPermissionEditAllow AgentPermissionEdit = "allow"
	AgentPermissionEditDeny  AgentPermissionEdit = "deny"
)

func (r AgentPermissionEdit) IsKnown() bool {
	switch r {
	case AgentPermissionEditAsk, AgentPermissionEditAllow, AgentPermissionEditDeny:
		return true
	}
	return false
}

type AgentPermissionWebfetch string

const (
	AgentPermissionWebfetchAsk   AgentPermissionWebfetch = "ask"
	AgentPermissionWebfetchAllow AgentPermissionWebfetch = "allow"
	AgentPermissionWebfetchDeny  AgentPermissionWebfetch = "deny"
)

func (r AgentPermissionWebfetch) IsKnown() bool {
	switch r {
	case AgentPermissionWebfetchAsk, AgentPermissionWebfetchAllow, AgentPermissionWebfetchDeny:
		return true
	}
	return false
}

type AgentModel struct {
	ModelID    string `json:"modelID"`
	ProviderID string `json:"providerID"`
}

func (r *AgentModel) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

type AgentListParams struct {
	Directory *string `query:"directory,omitempty"`
}

func (r AgentListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
