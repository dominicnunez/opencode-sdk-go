// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package opencode

import (
	"context"
	"net/http"
	"net/url"

	"github.com/dominicnunez/opencode-sdk-go/internal/apijson"
	"github.com/dominicnunez/opencode-sdk-go/internal/apiquery"
	"github.com/dominicnunez/opencode-sdk-go/internal/param"
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
	BuiltIn     bool                   `json:"builtIn,required"`
	Mode        AgentMode              `json:"mode,required"`
	Name        string                 `json:"name,required"`
	Options     map[string]interface{} `json:"options,required"`
	Permission  AgentPermission        `json:"permission,required"`
	Tools       map[string]bool        `json:"tools,required"`
	Description string                 `json:"description"`
	Model       AgentModel             `json:"model"`
	Prompt      string                 `json:"prompt"`
	Temperature float64                `json:"temperature"`
	TopP        float64                `json:"topP"`
	JSON        agentJSON              `json:"-"`
}

type agentJSON struct {
	BuiltIn     apijson.Field
	Mode        apijson.Field
	Name        apijson.Field
	Options     apijson.Field
	Permission  apijson.Field
	Tools       apijson.Field
	Description apijson.Field
	Model       apijson.Field
	Prompt      apijson.Field
	Temperature apijson.Field
	TopP        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *Agent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r agentJSON) RawJSON() string {
	return r.raw
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
	Bash     map[string]AgentPermissionBash `json:"bash,required"`
	Edit     AgentPermissionEdit            `json:"edit,required"`
	Webfetch AgentPermissionWebfetch        `json:"webfetch"`
	JSON     agentPermissionJSON            `json:"-"`
}

type agentPermissionJSON struct {
	Bash        apijson.Field
	Edit        apijson.Field
	Webfetch    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AgentPermission) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r agentPermissionJSON) RawJSON() string {
	return r.raw
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
	ModelID    string         `json:"modelID,required"`
	ProviderID string         `json:"providerID,required"`
	JSON       agentModelJSON `json:"-"`
}

type agentModelJSON struct {
	ModelID     apijson.Field
	ProviderID  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AgentModel) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r agentModelJSON) RawJSON() string {
	return r.raw
}

type AgentListParams struct {
	Directory param.Field[string] `query:"directory"`
}

func (r AgentListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
