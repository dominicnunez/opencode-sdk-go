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

type ConfigService struct {
	client *Client
}

func (s *ConfigService) Get(ctx context.Context, params *ConfigGetParams) (*Config, error) {
	if params == nil {
		params = &ConfigGetParams{}
	}
	var result Config
	err := s.client.do(ctx, http.MethodGet, "config", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ConfigService) Update(ctx context.Context, params *ConfigUpdateParams) (*Config, error) {
	if params == nil {
		return nil, errors.New("params is required")
	}
	var result Config
	err := s.client.do(ctx, http.MethodPatch, "config", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ConfigService) Providers(ctx context.Context, params *ConfigProviderListParams) (*ConfigProviderListResponse, error) {
	if params == nil {
		params = &ConfigProviderListParams{}
	}
	var result ConfigProviderListResponse
	err := s.client.do(ctx, http.MethodGet, "config/providers", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type Config struct {
	// JSON schema reference for configuration validation
	Schema string `json:"$schema,omitempty"`
	// Agent configuration, see https://opencode.ai/docs/agent
	Agent ConfigAgent `json:"agent,omitempty"`
	// Deprecated: Use the Share field instead. Share newly created sessions
	// automatically.
	Autoshare bool `json:"autoshare,omitempty"`
	// Automatically update to the latest version
	Autoupdate bool `json:"autoupdate,omitempty"`
	// Command configuration, see https://opencode.ai/docs/commands
	Command map[string]ConfigCommand `json:"command,omitempty"`
	// Disable providers that are loaded automatically
	DisabledProviders []string                   `json:"disabled_providers,omitempty"`
	Experimental      ConfigExperimental         `json:"experimental,omitempty"`
	Formatter         map[string]ConfigFormatter `json:"formatter,omitempty"`
	// Additional instruction files or patterns to include
	Instructions []string `json:"instructions,omitempty"`
	// Custom keybind configurations
	Keybinds KeybindsConfig `json:"keybinds,omitempty"`
	// Deprecated: Always uses stretch layout.
	Layout ConfigLayout         `json:"layout,omitempty"`
	Lsp    map[string]ConfigLsp `json:"lsp,omitempty"`
	// MCP (Model Context Protocol) server configurations
	Mcp map[string]ConfigMcp `json:"mcp,omitempty"`
	// Deprecated: Use the Agent field instead.
	Mode ConfigMode `json:"mode,omitempty"`
	// Model to use in the format of provider/model, eg anthropic/claude-2
	Model      string           `json:"model,omitempty"`
	Permission ConfigPermission `json:"permission,omitempty"`
	Plugin     []string         `json:"plugin,omitempty"`
	// Custom provider configurations and model overrides
	Provider map[string]ConfigProvider `json:"provider,omitempty"`
	// Control sharing behavior:'manual' allows manual sharing via commands, 'auto'
	// enables automatic sharing, 'disabled' disables all sharing
	Share ConfigShare `json:"share,omitempty"`
	// Small model to use for tasks like title generation in the format of
	// provider/model
	SmallModel string `json:"small_model,omitempty"`
	Snapshot   bool   `json:"snapshot,omitempty"`
	// Theme name to use for the interface
	Theme string          `json:"theme,omitempty"`
	Tools map[string]bool `json:"tools,omitempty"`
	// TUI specific settings
	Tui ConfigTui `json:"tui,omitempty"`
	// Custom username to display in conversations instead of system username
	Username string        `json:"username,omitempty"`
	Watcher  ConfigWatcher `json:"watcher,omitempty"`
}

// Agent configuration, see https://opencode.ai/docs/agent
type ConfigAgent struct {
	Build   ConfigAgentBuild   `json:"build"`
	General ConfigAgentGeneral `json:"general"`
	Plan    ConfigAgentPlan    `json:"plan"`
}

type ConfigAgentBuild struct {
	// Description of when to use the agent
	Description string                     `json:"description"`
	Disable     bool                       `json:"disable"`
	Mode        ConfigAgentBuildMode       `json:"mode"`
	Model       string                     `json:"model"`
	Permission  ConfigAgentBuildPermission `json:"permission"`
	Prompt      string                     `json:"prompt"`
	Temperature float64                    `json:"temperature"`
	Tools       map[string]bool            `json:"tools"`
	TopP        float64                    `json:"top_p"`
}

type ConfigAgentBuildMode string

const (
	ConfigAgentBuildModeSubagent ConfigAgentBuildMode = "subagent"
	ConfigAgentBuildModePrimary  ConfigAgentBuildMode = "primary"
	ConfigAgentBuildModeAll      ConfigAgentBuildMode = "all"
)

func (r ConfigAgentBuildMode) IsKnown() bool {
	switch r {
	case ConfigAgentBuildModeSubagent, ConfigAgentBuildModePrimary, ConfigAgentBuildModeAll:
		return true
	}
	return false
}

type ConfigAgentBuildPermission struct {
	Bash     ConfigAgentBuildPermissionBashUnion `json:"bash"`
	Edit     ConfigAgentBuildPermissionEdit      `json:"edit"`
	Webfetch ConfigAgentBuildPermissionWebfetch  `json:"webfetch"`
}

// ConfigAgentBuildPermissionBashUnion can be either a string or a map.
// Use AsString() or AsMap() to access the value.
type ConfigAgentBuildPermissionBashUnion struct {
	raw json.RawMessage
}

func (p *ConfigAgentBuildPermissionBashUnion) UnmarshalJSON(data []byte) error {
	p.raw = append(json.RawMessage(nil), data...)
	return nil
}

// AsString returns the value as a string if it is a string, or ("", ErrWrongVariant) otherwise.
func (p ConfigAgentBuildPermissionBashUnion) AsString() (ConfigAgentBuildPermissionBashString, error) {
	if len(p.raw) == 0 || p.raw[0] != '"' {
		return "", wrongVariant("string variant", "non-string JSON value")
	}
	var s ConfigAgentBuildPermissionBashString
	if err := json.Unmarshal(p.raw, &s); err != nil {
		return "", err
	}
	return s, nil
}

// AsMap returns the value as a map if it is a map, or (nil, ErrWrongVariant) otherwise.
func (p ConfigAgentBuildPermissionBashUnion) AsMap() (ConfigAgentBuildPermissionBashMap, error) {
	if len(p.raw) == 0 || p.raw[0] != '{' {
		return nil, wrongVariant("map variant", "non-object JSON value")
	}
	var m ConfigAgentBuildPermissionBashMap
	if err := json.Unmarshal(p.raw, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (p ConfigAgentBuildPermissionBashUnion) MarshalJSON() ([]byte, error) {
	if p.raw == nil {
		return []byte("null"), nil
	}
	return p.raw, nil
}

type ConfigAgentBuildPermissionBashString string

const (
	ConfigAgentBuildPermissionBashStringAsk   ConfigAgentBuildPermissionBashString = "ask"
	ConfigAgentBuildPermissionBashStringAllow ConfigAgentBuildPermissionBashString = "allow"
	ConfigAgentBuildPermissionBashStringDeny  ConfigAgentBuildPermissionBashString = "deny"
)

func (r ConfigAgentBuildPermissionBashString) IsKnown() bool {
	switch r {
	case ConfigAgentBuildPermissionBashStringAsk, ConfigAgentBuildPermissionBashStringAllow, ConfigAgentBuildPermissionBashStringDeny:
		return true
	}
	return false
}

type ConfigAgentBuildPermissionBashMap map[string]ConfigAgentBuildPermissionBashMapItem

type ConfigAgentBuildPermissionBashMapItem string

const (
	ConfigAgentBuildPermissionBashMapAsk   ConfigAgentBuildPermissionBashMapItem = "ask"
	ConfigAgentBuildPermissionBashMapAllow ConfigAgentBuildPermissionBashMapItem = "allow"
	ConfigAgentBuildPermissionBashMapDeny  ConfigAgentBuildPermissionBashMapItem = "deny"
)

func (r ConfigAgentBuildPermissionBashMapItem) IsKnown() bool {
	switch r {
	case ConfigAgentBuildPermissionBashMapAsk, ConfigAgentBuildPermissionBashMapAllow, ConfigAgentBuildPermissionBashMapDeny:
		return true
	}
	return false
}

type ConfigAgentBuildPermissionEdit string

const (
	ConfigAgentBuildPermissionEditAsk   ConfigAgentBuildPermissionEdit = "ask"
	ConfigAgentBuildPermissionEditAllow ConfigAgentBuildPermissionEdit = "allow"
	ConfigAgentBuildPermissionEditDeny  ConfigAgentBuildPermissionEdit = "deny"
)

func (r ConfigAgentBuildPermissionEdit) IsKnown() bool {
	switch r {
	case ConfigAgentBuildPermissionEditAsk, ConfigAgentBuildPermissionEditAllow, ConfigAgentBuildPermissionEditDeny:
		return true
	}
	return false
}

type ConfigAgentBuildPermissionWebfetch string

const (
	ConfigAgentBuildPermissionWebfetchAsk   ConfigAgentBuildPermissionWebfetch = "ask"
	ConfigAgentBuildPermissionWebfetchAllow ConfigAgentBuildPermissionWebfetch = "allow"
	ConfigAgentBuildPermissionWebfetchDeny  ConfigAgentBuildPermissionWebfetch = "deny"
)

func (r ConfigAgentBuildPermissionWebfetch) IsKnown() bool {
	switch r {
	case ConfigAgentBuildPermissionWebfetchAsk, ConfigAgentBuildPermissionWebfetchAllow, ConfigAgentBuildPermissionWebfetchDeny:
		return true
	}
	return false
}

type ConfigAgentGeneral struct {
	// Description of when to use the agent
	Description string                       `json:"description"`
	Disable     bool                         `json:"disable"`
	Mode        ConfigAgentGeneralMode       `json:"mode"`
	Model       string                       `json:"model"`
	Permission  ConfigAgentGeneralPermission `json:"permission"`
	Prompt      string                       `json:"prompt"`
	Temperature float64                      `json:"temperature"`
	Tools       map[string]bool              `json:"tools"`
	TopP        float64                      `json:"top_p"`
}

type ConfigAgentGeneralMode string

const (
	ConfigAgentGeneralModeSubagent ConfigAgentGeneralMode = "subagent"
	ConfigAgentGeneralModePrimary  ConfigAgentGeneralMode = "primary"
	ConfigAgentGeneralModeAll      ConfigAgentGeneralMode = "all"
)

func (r ConfigAgentGeneralMode) IsKnown() bool {
	switch r {
	case ConfigAgentGeneralModeSubagent, ConfigAgentGeneralModePrimary, ConfigAgentGeneralModeAll:
		return true
	}
	return false
}

type ConfigAgentGeneralPermission struct {
	Bash     ConfigAgentGeneralPermissionBashUnion `json:"bash"`
	Edit     ConfigAgentGeneralPermissionEdit      `json:"edit"`
	Webfetch ConfigAgentGeneralPermissionWebfetch  `json:"webfetch"`
}

// ConfigAgentGeneralPermissionBashUnion can be either a string or a map.
// Use AsString() or AsMap() to access the value.
type ConfigAgentGeneralPermissionBashUnion struct {
	raw json.RawMessage
}

func (p *ConfigAgentGeneralPermissionBashUnion) UnmarshalJSON(data []byte) error {
	p.raw = append(json.RawMessage(nil), data...)
	return nil
}

// AsString returns the value as a string if it is a string, or ("", ErrWrongVariant) otherwise.
func (p ConfigAgentGeneralPermissionBashUnion) AsString() (ConfigAgentGeneralPermissionBashString, error) {
	if len(p.raw) == 0 || p.raw[0] != '"' {
		return "", wrongVariant("string variant", "non-string JSON value")
	}
	var s ConfigAgentGeneralPermissionBashString
	if err := json.Unmarshal(p.raw, &s); err != nil {
		return "", err
	}
	return s, nil
}

// AsMap returns the value as a map if it is a map, or (nil, ErrWrongVariant) otherwise.
func (p ConfigAgentGeneralPermissionBashUnion) AsMap() (ConfigAgentGeneralPermissionBashMap, error) {
	if len(p.raw) == 0 || p.raw[0] != '{' {
		return nil, wrongVariant("map variant", "non-object JSON value")
	}
	var m ConfigAgentGeneralPermissionBashMap
	if err := json.Unmarshal(p.raw, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (p ConfigAgentGeneralPermissionBashUnion) MarshalJSON() ([]byte, error) {
	if p.raw == nil {
		return []byte("null"), nil
	}
	return p.raw, nil
}

type ConfigAgentGeneralPermissionBashString string

const (
	ConfigAgentGeneralPermissionBashStringAsk   ConfigAgentGeneralPermissionBashString = "ask"
	ConfigAgentGeneralPermissionBashStringAllow ConfigAgentGeneralPermissionBashString = "allow"
	ConfigAgentGeneralPermissionBashStringDeny  ConfigAgentGeneralPermissionBashString = "deny"
)

func (r ConfigAgentGeneralPermissionBashString) IsKnown() bool {
	switch r {
	case ConfigAgentGeneralPermissionBashStringAsk, ConfigAgentGeneralPermissionBashStringAllow, ConfigAgentGeneralPermissionBashStringDeny:
		return true
	}
	return false
}

type ConfigAgentGeneralPermissionBashMap map[string]ConfigAgentGeneralPermissionBashMapItem

type ConfigAgentGeneralPermissionBashMapItem string

const (
	ConfigAgentGeneralPermissionBashMapAsk   ConfigAgentGeneralPermissionBashMapItem = "ask"
	ConfigAgentGeneralPermissionBashMapAllow ConfigAgentGeneralPermissionBashMapItem = "allow"
	ConfigAgentGeneralPermissionBashMapDeny  ConfigAgentGeneralPermissionBashMapItem = "deny"
)

func (r ConfigAgentGeneralPermissionBashMapItem) IsKnown() bool {
	switch r {
	case ConfigAgentGeneralPermissionBashMapAsk, ConfigAgentGeneralPermissionBashMapAllow, ConfigAgentGeneralPermissionBashMapDeny:
		return true
	}
	return false
}

type ConfigAgentGeneralPermissionEdit string

const (
	ConfigAgentGeneralPermissionEditAsk   ConfigAgentGeneralPermissionEdit = "ask"
	ConfigAgentGeneralPermissionEditAllow ConfigAgentGeneralPermissionEdit = "allow"
	ConfigAgentGeneralPermissionEditDeny  ConfigAgentGeneralPermissionEdit = "deny"
)

func (r ConfigAgentGeneralPermissionEdit) IsKnown() bool {
	switch r {
	case ConfigAgentGeneralPermissionEditAsk, ConfigAgentGeneralPermissionEditAllow, ConfigAgentGeneralPermissionEditDeny:
		return true
	}
	return false
}

type ConfigAgentGeneralPermissionWebfetch string

const (
	ConfigAgentGeneralPermissionWebfetchAsk   ConfigAgentGeneralPermissionWebfetch = "ask"
	ConfigAgentGeneralPermissionWebfetchAllow ConfigAgentGeneralPermissionWebfetch = "allow"
	ConfigAgentGeneralPermissionWebfetchDeny  ConfigAgentGeneralPermissionWebfetch = "deny"
)

func (r ConfigAgentGeneralPermissionWebfetch) IsKnown() bool {
	switch r {
	case ConfigAgentGeneralPermissionWebfetchAsk, ConfigAgentGeneralPermissionWebfetchAllow, ConfigAgentGeneralPermissionWebfetchDeny:
		return true
	}
	return false
}

type ConfigAgentPlan struct {
	// Description of when to use the agent
	Description string                    `json:"description"`
	Disable     bool                      `json:"disable"`
	Mode        ConfigAgentPlanMode       `json:"mode"`
	Model       string                    `json:"model"`
	Permission  ConfigAgentPlanPermission `json:"permission"`
	Prompt      string                    `json:"prompt"`
	Temperature float64                   `json:"temperature"`
	Tools       map[string]bool           `json:"tools"`
	TopP        float64                   `json:"top_p"`
}

type ConfigAgentPlanMode string

const (
	ConfigAgentPlanModeSubagent ConfigAgentPlanMode = "subagent"
	ConfigAgentPlanModePrimary  ConfigAgentPlanMode = "primary"
	ConfigAgentPlanModeAll      ConfigAgentPlanMode = "all"
)

func (r ConfigAgentPlanMode) IsKnown() bool {
	switch r {
	case ConfigAgentPlanModeSubagent, ConfigAgentPlanModePrimary, ConfigAgentPlanModeAll:
		return true
	}
	return false
}

type ConfigAgentPlanPermission struct {
	Bash     ConfigAgentPlanPermissionBashUnion `json:"bash"`
	Edit     ConfigAgentPlanPermissionEdit      `json:"edit"`
	Webfetch ConfigAgentPlanPermissionWebfetch  `json:"webfetch"`
}

// ConfigAgentPlanPermissionBashUnion can be either a string or a map.
// Use AsString() or AsMap() to access the value.
type ConfigAgentPlanPermissionBashUnion struct {
	raw json.RawMessage
}

func (p *ConfigAgentPlanPermissionBashUnion) UnmarshalJSON(data []byte) error {
	p.raw = append(json.RawMessage(nil), data...)
	return nil
}

// AsString returns the value as a string if it is a string, or ("", ErrWrongVariant) otherwise.
func (p ConfigAgentPlanPermissionBashUnion) AsString() (ConfigAgentPlanPermissionBashString, error) {
	if len(p.raw) == 0 || p.raw[0] != '"' {
		return "", wrongVariant("string variant", "non-string JSON value")
	}
	var s ConfigAgentPlanPermissionBashString
	if err := json.Unmarshal(p.raw, &s); err != nil {
		return "", err
	}
	return s, nil
}

// AsMap returns the value as a map if it is a map, or (nil, ErrWrongVariant) otherwise.
func (p ConfigAgentPlanPermissionBashUnion) AsMap() (ConfigAgentPlanPermissionBashMap, error) {
	if len(p.raw) == 0 || p.raw[0] != '{' {
		return nil, wrongVariant("map variant", "non-object JSON value")
	}
	var m ConfigAgentPlanPermissionBashMap
	if err := json.Unmarshal(p.raw, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (p ConfigAgentPlanPermissionBashUnion) MarshalJSON() ([]byte, error) {
	if p.raw == nil {
		return []byte("null"), nil
	}
	return p.raw, nil
}

type ConfigAgentPlanPermissionBashString string

const (
	ConfigAgentPlanPermissionBashStringAsk   ConfigAgentPlanPermissionBashString = "ask"
	ConfigAgentPlanPermissionBashStringAllow ConfigAgentPlanPermissionBashString = "allow"
	ConfigAgentPlanPermissionBashStringDeny  ConfigAgentPlanPermissionBashString = "deny"
)

func (r ConfigAgentPlanPermissionBashString) IsKnown() bool {
	switch r {
	case ConfigAgentPlanPermissionBashStringAsk, ConfigAgentPlanPermissionBashStringAllow, ConfigAgentPlanPermissionBashStringDeny:
		return true
	}
	return false
}

type ConfigAgentPlanPermissionBashMap map[string]ConfigAgentPlanPermissionBashMapItem

type ConfigAgentPlanPermissionBashMapItem string

const (
	ConfigAgentPlanPermissionBashMapAsk   ConfigAgentPlanPermissionBashMapItem = "ask"
	ConfigAgentPlanPermissionBashMapAllow ConfigAgentPlanPermissionBashMapItem = "allow"
	ConfigAgentPlanPermissionBashMapDeny  ConfigAgentPlanPermissionBashMapItem = "deny"
)

func (r ConfigAgentPlanPermissionBashMapItem) IsKnown() bool {
	switch r {
	case ConfigAgentPlanPermissionBashMapAsk, ConfigAgentPlanPermissionBashMapAllow, ConfigAgentPlanPermissionBashMapDeny:
		return true
	}
	return false
}

type ConfigAgentPlanPermissionEdit string

const (
	ConfigAgentPlanPermissionEditAsk   ConfigAgentPlanPermissionEdit = "ask"
	ConfigAgentPlanPermissionEditAllow ConfigAgentPlanPermissionEdit = "allow"
	ConfigAgentPlanPermissionEditDeny  ConfigAgentPlanPermissionEdit = "deny"
)

func (r ConfigAgentPlanPermissionEdit) IsKnown() bool {
	switch r {
	case ConfigAgentPlanPermissionEditAsk, ConfigAgentPlanPermissionEditAllow, ConfigAgentPlanPermissionEditDeny:
		return true
	}
	return false
}

type ConfigAgentPlanPermissionWebfetch string

const (
	ConfigAgentPlanPermissionWebfetchAsk   ConfigAgentPlanPermissionWebfetch = "ask"
	ConfigAgentPlanPermissionWebfetchAllow ConfigAgentPlanPermissionWebfetch = "allow"
	ConfigAgentPlanPermissionWebfetchDeny  ConfigAgentPlanPermissionWebfetch = "deny"
)

func (r ConfigAgentPlanPermissionWebfetch) IsKnown() bool {
	switch r {
	case ConfigAgentPlanPermissionWebfetchAsk, ConfigAgentPlanPermissionWebfetchAllow, ConfigAgentPlanPermissionWebfetchDeny:
		return true
	}
	return false
}

type ConfigCommand struct {
	Template    string `json:"template"`
	Agent       string `json:"agent"`
	Description string `json:"description"`
	Model       string `json:"model"`
	Subtask     bool   `json:"subtask"`
}

type ConfigExperimental struct {
	DisablePasteSummary bool                   `json:"disable_paste_summary"`
	Hook                ConfigExperimentalHook `json:"hook"`
}

type ConfigExperimentalHook struct {
	FileEdited       map[string][]ConfigExperimentalHookFileEdited `json:"file_edited"`
	SessionCompleted []ConfigExperimentalHookSessionCompleted      `json:"session_completed"`
}

type ConfigExperimentalHookFileEdited struct {
	Command     []string          `json:"command"`
	Environment map[string]string `json:"environment"`
}

type ConfigExperimentalHookSessionCompleted struct {
	Command     []string          `json:"command"`
	Environment map[string]string `json:"environment"`
}

type ConfigFormatter struct {
	Command     []string          `json:"command"`
	Disabled    bool              `json:"disabled"`
	Environment map[string]string `json:"environment"`
	Extensions  []string          `json:"extensions"`
}

// @deprecated Always uses stretch layout.
type ConfigLayout string

const (
	ConfigLayoutAuto    ConfigLayout = "auto"
	ConfigLayoutStretch ConfigLayout = "stretch"
)

func (r ConfigLayout) IsKnown() bool {
	switch r {
	case ConfigLayoutAuto, ConfigLayoutStretch:
		return true
	}
	return false
}

// ConfigLsp represents LSP (Language Server Protocol) configuration.
// It can be either ConfigLspDisabled or ConfigLspObject.
// ConfigLspObject is identified by the presence of the "command" field.
// ConfigLspDisabled is identified by the absence of "command" and disabled=true.
type ConfigLsp struct {
	// raw stores the full JSON for lazy unmarshaling
	raw json.RawMessage
}

func (r *ConfigLsp) UnmarshalJSON(data []byte) error {
	r.raw = append(json.RawMessage(nil), data...)
	return nil
}

// AsDisabled returns the config as ConfigLspDisabled if it has disabled=true without command field.
// Returns (nil, ErrWrongVariant) if it's not a disabled config.
func (r ConfigLsp) AsDisabled() (*ConfigLspDisabled, error) {
	if r.raw == nil {
		return nil, wrongVariant("disabled config", "nil raw")
	}

	// Use a combined struct to check both discriminator fields in one pass
	var probe struct {
		ConfigLspDisabled
		Command json.RawMessage `json:"command"`
	}
	if err := json.Unmarshal(r.raw, &probe); err != nil {
		return nil, err
	}
	if len(probe.Command) > 0 {
		return nil, wrongVariant("disabled config", "object config with command")
	}
	if !probe.Disabled {
		return nil, wrongVariant("disabled config", "config with disabled=false")
	}

	result := probe.ConfigLspDisabled
	return &result, nil
}

// AsObject returns the config as ConfigLspObject if it has a command field.
// Returns (nil, ErrWrongVariant) if it's not an object config.
func (r ConfigLsp) AsObject() (*ConfigLspObject, error) {
	if r.raw == nil {
		return nil, wrongVariant("object config", "nil raw")
	}
	var obj ConfigLspObject
	if err := json.Unmarshal(r.raw, &obj); err != nil {
		return nil, err
	}

	if len(obj.Command) == 0 {
		return nil, wrongVariant("object config", "config without command")
	}

	return &obj, nil
}

func (r ConfigLsp) MarshalJSON() ([]byte, error) {
	if r.raw == nil {
		return []byte("null"), nil
	}
	return r.raw, nil
}

type ConfigLspDisabled struct {
	Disabled ConfigLspDisabledDisabled `json:"disabled"`
}

type ConfigLspDisabledDisabled bool

const (
	ConfigLspDisabledDisabledTrue ConfigLspDisabledDisabled = true
)

func (r ConfigLspDisabledDisabled) IsKnown() bool {
	switch r {
	case ConfigLspDisabledDisabledTrue:
		return true
	}
	return false
}

type ConfigLspObject struct {
	Command        []string               `json:"command"`
	Disabled       bool                   `json:"disabled"`
	Env            map[string]string      `json:"env"`
	Extensions     []string               `json:"extensions"`
	Initialization map[string]interface{} `json:"initialization"`
}

// ConfigMcp represents MCP (Model Context Protocol) server configuration.
// It can be either McpLocalConfig or McpRemoteConfig, discriminated by the type field.
type ConfigMcp struct {
	// Type discriminator: "local" or "remote"
	Type ConfigMcpType `json:"type"`
	// raw stores the full JSON for lazy unmarshaling
	raw json.RawMessage
}

func (r *ConfigMcp) UnmarshalJSON(data []byte) error {
	// Peek at the discriminator field
	var peek struct {
		Type ConfigMcpType `json:"type"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return err
	}
	r.Type = peek.Type
	r.raw = append(json.RawMessage(nil), data...)
	return nil
}

// AsLocal returns the config as McpLocalConfig if type is "local".
// Returns (nil, ErrWrongVariant) if the type is not "local".
func (r ConfigMcp) AsLocal() (*McpLocalConfig, error) {
	if r.raw == nil {
		return nil, wrongVariant("local config", "nil raw")
	}
	if r.Type != ConfigMcpTypeLocal {
		return nil, wrongVariant("local", string(r.Type))
	}
	var local McpLocalConfig
	if err := json.Unmarshal(r.raw, &local); err != nil {
		return nil, fmt.Errorf("unmarshal %s Type: %w", r.Type, err)
	}
	return &local, nil
}

// AsRemote returns the config as McpRemoteConfig if type is "remote".
// Returns (nil, ErrWrongVariant) if the type is not "remote".
func (r ConfigMcp) AsRemote() (*McpRemoteConfig, error) {
	if r.raw == nil {
		return nil, wrongVariant("remote config", "nil raw")
	}
	if r.Type != ConfigMcpTypeRemote {
		return nil, wrongVariant("remote", string(r.Type))
	}
	var remote McpRemoteConfig
	if err := json.Unmarshal(r.raw, &remote); err != nil {
		return nil, fmt.Errorf("unmarshal %s Type: %w", r.Type, err)
	}
	return &remote, nil
}

func (r ConfigMcp) MarshalJSON() ([]byte, error) {
	if r.raw == nil {
		return []byte("null"), nil
	}
	return r.raw, nil
}

// Type of MCP server connection
type ConfigMcpType string

const (
	ConfigMcpTypeLocal  ConfigMcpType = "local"
	ConfigMcpTypeRemote ConfigMcpType = "remote"
)

func (r ConfigMcpType) IsKnown() bool {
	switch r {
	case ConfigMcpTypeLocal, ConfigMcpTypeRemote:
		return true
	}
	return false
}

// @deprecated Use `agent` field instead.
type ConfigMode struct {
	Build ConfigModeBuild `json:"build"`
	Plan  ConfigModePlan  `json:"plan"`
}

type ConfigModeBuild struct {
	// Description of when to use the agent
	Description string                    `json:"description"`
	Disable     bool                      `json:"disable"`
	Mode        ConfigModeBuildMode       `json:"mode"`
	Model       string                    `json:"model"`
	Permission  ConfigModeBuildPermission `json:"permission"`
	Prompt      string                    `json:"prompt"`
	Temperature float64                   `json:"temperature"`
	Tools       map[string]bool           `json:"tools"`
	TopP        float64                   `json:"top_p"`
}

type ConfigModeBuildMode string

const (
	ConfigModeBuildModeSubagent ConfigModeBuildMode = "subagent"
	ConfigModeBuildModePrimary  ConfigModeBuildMode = "primary"
	ConfigModeBuildModeAll      ConfigModeBuildMode = "all"
)

func (r ConfigModeBuildMode) IsKnown() bool {
	switch r {
	case ConfigModeBuildModeSubagent, ConfigModeBuildModePrimary, ConfigModeBuildModeAll:
		return true
	}
	return false
}

type ConfigModeBuildPermission struct {
	Bash     ConfigModeBuildPermissionBashUnion `json:"bash"`
	Edit     ConfigModeBuildPermissionEdit      `json:"edit"`
	Webfetch ConfigModeBuildPermissionWebfetch  `json:"webfetch"`
}

// ConfigModeBuildPermissionBashUnion can be either a string or a map.
// Use AsString() or AsMap() to access the value.
type ConfigModeBuildPermissionBashUnion struct {
	raw json.RawMessage
}

func (p *ConfigModeBuildPermissionBashUnion) UnmarshalJSON(data []byte) error {
	p.raw = append(json.RawMessage(nil), data...)
	return nil
}

// AsString returns the value as a string if it is a string, or ("", ErrWrongVariant) otherwise.
func (p ConfigModeBuildPermissionBashUnion) AsString() (ConfigModeBuildPermissionBashString, error) {
	if len(p.raw) == 0 || p.raw[0] != '"' {
		return "", wrongVariant("string variant", "non-string JSON value")
	}
	var s ConfigModeBuildPermissionBashString
	if err := json.Unmarshal(p.raw, &s); err != nil {
		return "", err
	}
	return s, nil
}

// AsMap returns the value as a map if it is a map, or (nil, ErrWrongVariant) otherwise.
func (p ConfigModeBuildPermissionBashUnion) AsMap() (ConfigModeBuildPermissionBashMap, error) {
	if len(p.raw) == 0 || p.raw[0] != '{' {
		return nil, wrongVariant("map variant", "non-object JSON value")
	}
	var m ConfigModeBuildPermissionBashMap
	if err := json.Unmarshal(p.raw, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (p ConfigModeBuildPermissionBashUnion) MarshalJSON() ([]byte, error) {
	if p.raw == nil {
		return []byte("null"), nil
	}
	return p.raw, nil
}

type ConfigModeBuildPermissionBashString string

const (
	ConfigModeBuildPermissionBashStringAsk   ConfigModeBuildPermissionBashString = "ask"
	ConfigModeBuildPermissionBashStringAllow ConfigModeBuildPermissionBashString = "allow"
	ConfigModeBuildPermissionBashStringDeny  ConfigModeBuildPermissionBashString = "deny"
)

func (r ConfigModeBuildPermissionBashString) IsKnown() bool {
	switch r {
	case ConfigModeBuildPermissionBashStringAsk, ConfigModeBuildPermissionBashStringAllow, ConfigModeBuildPermissionBashStringDeny:
		return true
	}
	return false
}

type ConfigModeBuildPermissionBashMap map[string]ConfigModeBuildPermissionBashMapItem

type ConfigModeBuildPermissionBashMapItem string

const (
	ConfigModeBuildPermissionBashMapAsk   ConfigModeBuildPermissionBashMapItem = "ask"
	ConfigModeBuildPermissionBashMapAllow ConfigModeBuildPermissionBashMapItem = "allow"
	ConfigModeBuildPermissionBashMapDeny  ConfigModeBuildPermissionBashMapItem = "deny"
)

func (r ConfigModeBuildPermissionBashMapItem) IsKnown() bool {
	switch r {
	case ConfigModeBuildPermissionBashMapAsk, ConfigModeBuildPermissionBashMapAllow, ConfigModeBuildPermissionBashMapDeny:
		return true
	}
	return false
}

type ConfigModeBuildPermissionEdit string

const (
	ConfigModeBuildPermissionEditAsk   ConfigModeBuildPermissionEdit = "ask"
	ConfigModeBuildPermissionEditAllow ConfigModeBuildPermissionEdit = "allow"
	ConfigModeBuildPermissionEditDeny  ConfigModeBuildPermissionEdit = "deny"
)

func (r ConfigModeBuildPermissionEdit) IsKnown() bool {
	switch r {
	case ConfigModeBuildPermissionEditAsk, ConfigModeBuildPermissionEditAllow, ConfigModeBuildPermissionEditDeny:
		return true
	}
	return false
}

type ConfigModeBuildPermissionWebfetch string

const (
	ConfigModeBuildPermissionWebfetchAsk   ConfigModeBuildPermissionWebfetch = "ask"
	ConfigModeBuildPermissionWebfetchAllow ConfigModeBuildPermissionWebfetch = "allow"
	ConfigModeBuildPermissionWebfetchDeny  ConfigModeBuildPermissionWebfetch = "deny"
)

func (r ConfigModeBuildPermissionWebfetch) IsKnown() bool {
	switch r {
	case ConfigModeBuildPermissionWebfetchAsk, ConfigModeBuildPermissionWebfetchAllow, ConfigModeBuildPermissionWebfetchDeny:
		return true
	}
	return false
}

type ConfigModePlan struct {
	// Description of when to use the agent
	Description string                   `json:"description"`
	Disable     bool                     `json:"disable"`
	Mode        ConfigModePlanMode       `json:"mode"`
	Model       string                   `json:"model"`
	Permission  ConfigModePlanPermission `json:"permission"`
	Prompt      string                   `json:"prompt"`
	Temperature float64                  `json:"temperature"`
	Tools       map[string]bool          `json:"tools"`
	TopP        float64                  `json:"top_p"`
}

type ConfigModePlanMode string

const (
	ConfigModePlanModeSubagent ConfigModePlanMode = "subagent"
	ConfigModePlanModePrimary  ConfigModePlanMode = "primary"
	ConfigModePlanModeAll      ConfigModePlanMode = "all"
)

func (r ConfigModePlanMode) IsKnown() bool {
	switch r {
	case ConfigModePlanModeSubagent, ConfigModePlanModePrimary, ConfigModePlanModeAll:
		return true
	}
	return false
}

type ConfigModePlanPermission struct {
	Bash     ConfigModePlanPermissionBashUnion `json:"bash"`
	Edit     ConfigModePlanPermissionEdit      `json:"edit"`
	Webfetch ConfigModePlanPermissionWebfetch  `json:"webfetch"`
}

// ConfigModePlanPermissionBashUnion can be either a string or a map.
// Use AsString() or AsMap() to access the value.
type ConfigModePlanPermissionBashUnion struct {
	raw json.RawMessage
}

func (p *ConfigModePlanPermissionBashUnion) UnmarshalJSON(data []byte) error {
	p.raw = append(json.RawMessage(nil), data...)
	return nil
}

// AsString returns the value as a string if it is a string, or ("", ErrWrongVariant) otherwise.
func (p ConfigModePlanPermissionBashUnion) AsString() (ConfigModePlanPermissionBashString, error) {
	if len(p.raw) == 0 || p.raw[0] != '"' {
		return "", wrongVariant("string variant", "non-string JSON value")
	}
	var s ConfigModePlanPermissionBashString
	if err := json.Unmarshal(p.raw, &s); err != nil {
		return "", err
	}
	return s, nil
}

// AsMap returns the value as a map if it is a map, or (nil, ErrWrongVariant) otherwise.
func (p ConfigModePlanPermissionBashUnion) AsMap() (ConfigModePlanPermissionBashMap, error) {
	if len(p.raw) == 0 || p.raw[0] != '{' {
		return nil, wrongVariant("map variant", "non-object JSON value")
	}
	var m ConfigModePlanPermissionBashMap
	if err := json.Unmarshal(p.raw, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (p ConfigModePlanPermissionBashUnion) MarshalJSON() ([]byte, error) {
	if p.raw == nil {
		return []byte("null"), nil
	}
	return p.raw, nil
}

type ConfigModePlanPermissionBashString string

const (
	ConfigModePlanPermissionBashStringAsk   ConfigModePlanPermissionBashString = "ask"
	ConfigModePlanPermissionBashStringAllow ConfigModePlanPermissionBashString = "allow"
	ConfigModePlanPermissionBashStringDeny  ConfigModePlanPermissionBashString = "deny"
)

func (r ConfigModePlanPermissionBashString) IsKnown() bool {
	switch r {
	case ConfigModePlanPermissionBashStringAsk, ConfigModePlanPermissionBashStringAllow, ConfigModePlanPermissionBashStringDeny:
		return true
	}
	return false
}

type ConfigModePlanPermissionBashMap map[string]ConfigModePlanPermissionBashMapItem

type ConfigModePlanPermissionBashMapItem string

const (
	ConfigModePlanPermissionBashMapAsk   ConfigModePlanPermissionBashMapItem = "ask"
	ConfigModePlanPermissionBashMapAllow ConfigModePlanPermissionBashMapItem = "allow"
	ConfigModePlanPermissionBashMapDeny  ConfigModePlanPermissionBashMapItem = "deny"
)

func (r ConfigModePlanPermissionBashMapItem) IsKnown() bool {
	switch r {
	case ConfigModePlanPermissionBashMapAsk, ConfigModePlanPermissionBashMapAllow, ConfigModePlanPermissionBashMapDeny:
		return true
	}
	return false
}

type ConfigModePlanPermissionEdit string

const (
	ConfigModePlanPermissionEditAsk   ConfigModePlanPermissionEdit = "ask"
	ConfigModePlanPermissionEditAllow ConfigModePlanPermissionEdit = "allow"
	ConfigModePlanPermissionEditDeny  ConfigModePlanPermissionEdit = "deny"
)

func (r ConfigModePlanPermissionEdit) IsKnown() bool {
	switch r {
	case ConfigModePlanPermissionEditAsk, ConfigModePlanPermissionEditAllow, ConfigModePlanPermissionEditDeny:
		return true
	}
	return false
}

type ConfigModePlanPermissionWebfetch string

const (
	ConfigModePlanPermissionWebfetchAsk   ConfigModePlanPermissionWebfetch = "ask"
	ConfigModePlanPermissionWebfetchAllow ConfigModePlanPermissionWebfetch = "allow"
	ConfigModePlanPermissionWebfetchDeny  ConfigModePlanPermissionWebfetch = "deny"
)

func (r ConfigModePlanPermissionWebfetch) IsKnown() bool {
	switch r {
	case ConfigModePlanPermissionWebfetchAsk, ConfigModePlanPermissionWebfetchAllow, ConfigModePlanPermissionWebfetchDeny:
		return true
	}
	return false
}

type ConfigPermission struct {
	Bash     ConfigPermissionBashUnion `json:"bash"`
	Edit     ConfigPermissionEdit      `json:"edit"`
	Webfetch ConfigPermissionWebfetch  `json:"webfetch"`
}

// ConfigPermissionBashUnion can be either a string or a map.
// Use AsString() or AsMap() to access the value.
type ConfigPermissionBashUnion struct {
	raw json.RawMessage
}

func (p *ConfigPermissionBashUnion) UnmarshalJSON(data []byte) error {
	p.raw = append(json.RawMessage(nil), data...)
	return nil
}

// AsString returns the value as a string if it is a string, or ("", ErrWrongVariant) otherwise.
func (p ConfigPermissionBashUnion) AsString() (ConfigPermissionBashString, error) {
	if len(p.raw) == 0 || p.raw[0] != '"' {
		return "", wrongVariant("string variant", "non-string JSON value")
	}
	var s ConfigPermissionBashString
	if err := json.Unmarshal(p.raw, &s); err != nil {
		return "", err
	}
	return s, nil
}

// AsMap returns the value as a map if it is a map, or (nil, ErrWrongVariant) otherwise.
func (p ConfigPermissionBashUnion) AsMap() (ConfigPermissionBashMap, error) {
	if len(p.raw) == 0 || p.raw[0] != '{' {
		return nil, wrongVariant("map variant", "non-object JSON value")
	}
	var m ConfigPermissionBashMap
	if err := json.Unmarshal(p.raw, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (p ConfigPermissionBashUnion) MarshalJSON() ([]byte, error) {
	if p.raw == nil {
		return []byte("null"), nil
	}
	return p.raw, nil
}

type ConfigPermissionBashString string

const (
	ConfigPermissionBashStringAsk   ConfigPermissionBashString = "ask"
	ConfigPermissionBashStringAllow ConfigPermissionBashString = "allow"
	ConfigPermissionBashStringDeny  ConfigPermissionBashString = "deny"
)

func (r ConfigPermissionBashString) IsKnown() bool {
	switch r {
	case ConfigPermissionBashStringAsk, ConfigPermissionBashStringAllow, ConfigPermissionBashStringDeny:
		return true
	}
	return false
}

type ConfigPermissionBashMap map[string]ConfigPermissionBashMapItem

type ConfigPermissionBashMapItem string

const (
	ConfigPermissionBashMapAsk   ConfigPermissionBashMapItem = "ask"
	ConfigPermissionBashMapAllow ConfigPermissionBashMapItem = "allow"
	ConfigPermissionBashMapDeny  ConfigPermissionBashMapItem = "deny"
)

func (r ConfigPermissionBashMapItem) IsKnown() bool {
	switch r {
	case ConfigPermissionBashMapAsk, ConfigPermissionBashMapAllow, ConfigPermissionBashMapDeny:
		return true
	}
	return false
}

type ConfigPermissionEdit string

const (
	ConfigPermissionEditAsk   ConfigPermissionEdit = "ask"
	ConfigPermissionEditAllow ConfigPermissionEdit = "allow"
	ConfigPermissionEditDeny  ConfigPermissionEdit = "deny"
)

func (r ConfigPermissionEdit) IsKnown() bool {
	switch r {
	case ConfigPermissionEditAsk, ConfigPermissionEditAllow, ConfigPermissionEditDeny:
		return true
	}
	return false
}

type ConfigPermissionWebfetch string

const (
	ConfigPermissionWebfetchAsk   ConfigPermissionWebfetch = "ask"
	ConfigPermissionWebfetchAllow ConfigPermissionWebfetch = "allow"
	ConfigPermissionWebfetchDeny  ConfigPermissionWebfetch = "deny"
)

func (r ConfigPermissionWebfetch) IsKnown() bool {
	switch r {
	case ConfigPermissionWebfetchAsk, ConfigPermissionWebfetchAllow, ConfigPermissionWebfetchDeny:
		return true
	}
	return false
}

type ConfigProvider struct {
	ID      string                         `json:"id"`
	API     string                         `json:"api"`
	Env     []string                       `json:"env"`
	Models  map[string]ConfigProviderModel `json:"models"`
	Name    string                         `json:"name"`
	Npm     string                         `json:"npm"`
	Options ConfigProviderOptions          `json:"options"`
}

type ConfigProviderModel struct {
	ID           string                         `json:"id"`
	Attachment   bool                           `json:"attachment"`
	Cost         ConfigProviderModelsCost       `json:"cost"`
	Experimental bool                           `json:"experimental"`
	Limit        ConfigProviderModelsLimit      `json:"limit"`
	Modalities   ConfigProviderModelsModalities `json:"modalities"`
	Name         string                         `json:"name"`
	Options      map[string]interface{}         `json:"options"`
	Provider     ConfigProviderModelsProvider   `json:"provider"`
	Reasoning    bool                           `json:"reasoning"`
	ReleaseDate  string                         `json:"release_date"`
	Status       ConfigProviderModelsStatus     `json:"status"`
	Temperature  bool                           `json:"temperature"`
	ToolCall     bool                           `json:"tool_call"`
}

type ConfigProviderModelsCost struct {
	Input      float64 `json:"input"`
	Output     float64 `json:"output"`
	CacheRead  float64 `json:"cache_read,omitempty"`
	CacheWrite float64 `json:"cache_write,omitempty"`
}

type ConfigProviderModelsLimit struct {
	Context float64 `json:"context"`
	Output  float64 `json:"output"`
}

type ConfigProviderModelsModalities struct {
	Input  []ConfigProviderModelsModalitiesInput  `json:"input"`
	Output []ConfigProviderModelsModalitiesOutput `json:"output"`
}

type ConfigProviderModelsModalitiesInput string

const (
	ConfigProviderModelsModalitiesInputText  ConfigProviderModelsModalitiesInput = "text"
	ConfigProviderModelsModalitiesInputAudio ConfigProviderModelsModalitiesInput = "audio"
	ConfigProviderModelsModalitiesInputImage ConfigProviderModelsModalitiesInput = "image"
	ConfigProviderModelsModalitiesInputVideo ConfigProviderModelsModalitiesInput = "video"
	ConfigProviderModelsModalitiesInputPdf   ConfigProviderModelsModalitiesInput = "pdf"
)

func (r ConfigProviderModelsModalitiesInput) IsKnown() bool {
	switch r {
	case ConfigProviderModelsModalitiesInputText, ConfigProviderModelsModalitiesInputAudio, ConfigProviderModelsModalitiesInputImage, ConfigProviderModelsModalitiesInputVideo, ConfigProviderModelsModalitiesInputPdf:
		return true
	}
	return false
}

type ConfigProviderModelsModalitiesOutput string

const (
	ConfigProviderModelsModalitiesOutputText  ConfigProviderModelsModalitiesOutput = "text"
	ConfigProviderModelsModalitiesOutputAudio ConfigProviderModelsModalitiesOutput = "audio"
	ConfigProviderModelsModalitiesOutputImage ConfigProviderModelsModalitiesOutput = "image"
	ConfigProviderModelsModalitiesOutputVideo ConfigProviderModelsModalitiesOutput = "video"
	ConfigProviderModelsModalitiesOutputPdf   ConfigProviderModelsModalitiesOutput = "pdf"
)

func (r ConfigProviderModelsModalitiesOutput) IsKnown() bool {
	switch r {
	case ConfigProviderModelsModalitiesOutputText, ConfigProviderModelsModalitiesOutputAudio, ConfigProviderModelsModalitiesOutputImage, ConfigProviderModelsModalitiesOutputVideo, ConfigProviderModelsModalitiesOutputPdf:
		return true
	}
	return false
}

type ConfigProviderModelsProvider struct {
	Npm string `json:"npm"`
}

type ConfigProviderModelsStatus string

const (
	ConfigProviderModelsStatusAlpha ConfigProviderModelsStatus = "alpha"
	ConfigProviderModelsStatusBeta  ConfigProviderModelsStatus = "beta"
)

func (r ConfigProviderModelsStatus) IsKnown() bool {
	switch r {
	case ConfigProviderModelsStatusAlpha, ConfigProviderModelsStatusBeta:
		return true
	}
	return false
}

type ConfigProviderOptions struct {
	APIKey  string `json:"apiKey"`
	BaseURL string `json:"baseURL"`
	// Timeout in milliseconds for requests to this provider. Default is 300000 (5
	// minutes). Set to false to disable timeout.
	Timeout ConfigProviderOptionsTimeoutUnion `json:"timeout"`
}

// ConfigProviderOptionsTimeoutUnion can be either an int64 or a bool.
// Use AsInt() or AsBool() to access the value.
type ConfigProviderOptionsTimeoutUnion struct {
	raw json.RawMessage
}

func (p *ConfigProviderOptionsTimeoutUnion) UnmarshalJSON(data []byte) error {
	p.raw = append(json.RawMessage(nil), data...)
	return nil
}

// AsInt returns the timeout as an int64 if it is a number, or (0, ErrWrongVariant) otherwise.
func (p ConfigProviderOptionsTimeoutUnion) AsInt() (int64, error) {
	if len(p.raw) == 0 {
		return 0, wrongVariant("int variant", "empty raw")
	}
	c := p.raw[0]
	if c != '-' && (c < '0' || c > '9') {
		return 0, wrongVariant("int variant", "non-numeric JSON value")
	}
	var i int64
	if err := json.Unmarshal(p.raw, &i); err != nil {
		return 0, err
	}
	return i, nil
}

// AsBool returns the timeout as a bool if it is a bool, or (false, ErrWrongVariant) otherwise.
func (p ConfigProviderOptionsTimeoutUnion) AsBool() (bool, error) {
	if len(p.raw) == 0 {
		return false, wrongVariant("bool variant", "empty raw")
	}
	c := p.raw[0]
	if c != 't' && c != 'f' {
		return false, wrongVariant("bool variant", "non-boolean JSON value")
	}
	var b bool
	if err := json.Unmarshal(p.raw, &b); err != nil {
		return false, err
	}
	return b, nil
}

func (p ConfigProviderOptionsTimeoutUnion) MarshalJSON() ([]byte, error) {
	if p.raw == nil {
		return []byte("null"), nil
	}
	return p.raw, nil
}

// Control sharing behavior:'manual' allows manual sharing via commands, 'auto'
// enables automatic sharing, 'disabled' disables all sharing
type ConfigShare string

const (
	ConfigShareManual   ConfigShare = "manual"
	ConfigShareAuto     ConfigShare = "auto"
	ConfigShareDisabled ConfigShare = "disabled"
)

func (r ConfigShare) IsKnown() bool {
	switch r {
	case ConfigShareManual, ConfigShareAuto, ConfigShareDisabled:
		return true
	}
	return false
}

// TUI specific settings
type ConfigTui struct {
	// TUI scroll speed
	ScrollSpeed float64 `json:"scroll_speed"`
}

type ConfigWatcher struct {
	Ignore []string `json:"ignore"`
}

// Custom keybind configurations
type KeybindsConfig struct {
	// Next agent
	AgentCycle string `json:"agent_cycle"`
	// Previous agent
	AgentCycleReverse string `json:"agent_cycle_reverse"`
	// List agents
	AgentList string `json:"agent_list"`
	// Exit the application
	AppExit string `json:"app_exit"`
	// Show help dialog
	AppHelp string `json:"app_help"`
	// Open external editor
	EditorOpen string `json:"editor_open"`
	// @deprecated Close file
	FileClose string `json:"file_close"`
	// @deprecated Split/unified diff
	FileDiffToggle string `json:"file_diff_toggle"`
	// @deprecated Currently not available. List files
	FileList string `json:"file_list"`
	// @deprecated Search file
	FileSearch string `json:"file_search"`
	// Clear input field
	InputClear string `json:"input_clear"`
	// Insert newline in input
	InputNewline string `json:"input_newline"`
	// Paste from clipboard
	InputPaste string `json:"input_paste"`
	// Submit input
	InputSubmit string `json:"input_submit"`
	// Leader key for keybind combinations
	Leader string `json:"leader"`
	// Copy message
	MessagesCopy string `json:"messages_copy"`
	// Navigate to first message
	MessagesFirst string `json:"messages_first"`
	// Scroll messages down by half page
	MessagesHalfPageDown string `json:"messages_half_page_down"`
	// Scroll messages up by half page
	MessagesHalfPageUp string `json:"messages_half_page_up"`
	// Navigate to last message
	MessagesLast string `json:"messages_last"`
	// @deprecated Toggle layout
	MessagesLayoutToggle string `json:"messages_layout_toggle"`
	// @deprecated Navigate to next message
	MessagesNext string `json:"messages_next"`
	// Scroll messages down by one page
	MessagesPageDown string `json:"messages_page_down"`
	// Scroll messages up by one page
	MessagesPageUp string `json:"messages_page_up"`
	// @deprecated Navigate to previous message
	MessagesPrevious string `json:"messages_previous"`
	// Redo message
	MessagesRedo string `json:"messages_redo"`
	// @deprecated use messages_undo. Revert message
	MessagesRevert string `json:"messages_revert"`
	// Undo message
	MessagesUndo string `json:"messages_undo"`
	// Next recent model
	ModelCycleRecent string `json:"model_cycle_recent"`
	// Previous recent model
	ModelCycleRecentReverse string `json:"model_cycle_recent_reverse"`
	// List available models
	ModelList string `json:"model_list"`
	// Create/update AGENTS.md
	ProjectInit string `json:"project_init"`
	// Cycle to next child session
	SessionChildCycle string `json:"session_child_cycle"`
	// Cycle to previous child session
	SessionChildCycleReverse string `json:"session_child_cycle_reverse"`
	// Compact the session
	SessionCompact string `json:"session_compact"`
	// Export session to editor
	SessionExport string `json:"session_export"`
	// Interrupt current session
	SessionInterrupt string `json:"session_interrupt"`
	// List all sessions
	SessionList string `json:"session_list"`
	// Create a new session
	SessionNew string `json:"session_new"`
	// Share current session
	SessionShare string `json:"session_share"`
	// Show session timeline
	SessionTimeline string `json:"session_timeline"`
	// Unshare current session
	SessionUnshare string `json:"session_unshare"`
	// @deprecated use agent_cycle. Next agent
	SwitchAgent string `json:"switch_agent"`
	// @deprecated use agent_cycle_reverse. Previous agent
	SwitchAgentReverse string `json:"switch_agent_reverse"`
	// @deprecated use agent_cycle. Next mode
	SwitchMode string `json:"switch_mode"`
	// @deprecated use agent_cycle_reverse. Previous mode
	SwitchModeReverse string `json:"switch_mode_reverse"`
	// List available themes
	ThemeList string `json:"theme_list"`
	// Toggle thinking blocks
	ThinkingBlocks string `json:"thinking_blocks"`
	// Toggle tool details
	ToolDetails string `json:"tool_details"`
}

type McpLocalConfig struct {
	// Command and arguments to run the MCP server
	Command []string `json:"command"`
	// Type of MCP server connection
	Type McpLocalConfigType `json:"type"`
	// Enable or disable the MCP server on startup
	Enabled bool `json:"enabled"`
	// Environment variables to set when running the MCP server
	Environment map[string]string `json:"environment"`
}

// Type of MCP server connection
type McpLocalConfigType string

const (
	McpLocalConfigTypeLocal McpLocalConfigType = "local"
)

func (r McpLocalConfigType) IsKnown() bool {
	switch r {
	case McpLocalConfigTypeLocal:
		return true
	}
	return false
}

type McpRemoteConfig struct {
	// Type of MCP server connection
	Type McpRemoteConfigType `json:"type"`
	// URL of the remote MCP server
	URL string `json:"url"`
	// Enable or disable the MCP server on startup
	Enabled bool `json:"enabled"`
	// Headers to send with the request
	Headers map[string]string `json:"headers"`
}

// Type of MCP server connection
type McpRemoteConfigType string

const (
	McpRemoteConfigTypeRemote McpRemoteConfigType = "remote"
)

func (r McpRemoteConfigType) IsKnown() bool {
	switch r {
	case McpRemoteConfigTypeRemote:
		return true
	}
	return false
}

// Auth represents authentication credentials. It can be one of OAuth, ApiAuth, or WellKnownAuth,
// discriminated by the type field.
type Auth struct {
	// Type discriminator: "oauth", "api", or "wellknown"
	Type AuthType `json:"type"`
	// raw stores the full JSON for lazy unmarshaling
	raw json.RawMessage
}

// AuthType is the discriminator for Auth union
type AuthType string

const (
	AuthTypeOAuth     AuthType = "oauth"
	AuthTypeAPI       AuthType = "api"
	AuthTypeWellKnown AuthType = "wellknown"
)

func (r AuthType) IsKnown() bool {
	switch r {
	case AuthTypeOAuth, AuthTypeAPI, AuthTypeWellKnown:
		return true
	}
	return false
}

// UnmarshalJSON implements json.Unmarshaler for Auth
func (a *Auth) UnmarshalJSON(data []byte) error {
	// Peek at discriminator
	var peek struct {
		Type AuthType `json:"type"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return err
	}
	a.Type = peek.Type
	a.raw = append(json.RawMessage(nil), data...)
	return nil
}

// AsOAuth returns the OAuth variant if the type is oauth
func (a Auth) AsOAuth() (*OAuth, error) {
	if a.Type != AuthTypeOAuth {
		return nil, wrongVariant("oauth", string(a.Type))
	}
	var oauth OAuth
	if err := json.Unmarshal(a.raw, &oauth); err != nil {
		return nil, fmt.Errorf("unmarshal %s Type: %w", a.Type, err)
	}
	return &oauth, nil
}

// AsAPI returns the ApiAuth variant if the type is api
func (a Auth) AsAPI() (*ApiAuth, error) {
	if a.Type != AuthTypeAPI {
		return nil, wrongVariant("api", string(a.Type))
	}
	var apiAuth ApiAuth
	if err := json.Unmarshal(a.raw, &apiAuth); err != nil {
		return nil, fmt.Errorf("unmarshal %s Type: %w", a.Type, err)
	}
	return &apiAuth, nil
}

// AsWellKnown returns the WellKnownAuth variant if the type is wellknown
func (a Auth) AsWellKnown() (*WellKnownAuth, error) {
	if a.Type != AuthTypeWellKnown {
		return nil, wrongVariant("wellknown", string(a.Type))
	}
	var wellKnown WellKnownAuth
	if err := json.Unmarshal(a.raw, &wellKnown); err != nil {
		return nil, fmt.Errorf("unmarshal %s Type: %w", a.Type, err)
	}
	return &wellKnown, nil
}

func (a Auth) MarshalJSON() ([]byte, error) {
	if a.raw == nil {
		return []byte("null"), nil
	}
	return a.raw, nil
}

// OAuth represents OAuth authentication credentials
type OAuth struct {
	Type    AuthType `json:"type"`
	Refresh string   `json:"refresh"`
	Access  string   `json:"access"`
	Expires float64  `json:"expires"`
}

func (OAuth) implementsAuthSetParamsAuthUnion() {}

// ApiAuth represents API key authentication
type ApiAuth struct {
	Type AuthType `json:"type"`
	Key  string   `json:"key"`
}

func (ApiAuth) implementsAuthSetParamsAuthUnion() {}

// WellKnownAuth represents well-known authentication
type WellKnownAuth struct {
	Type  AuthType `json:"type"`
	Key   string   `json:"key"`
	Token string   `json:"token"`
}

func (WellKnownAuth) implementsAuthSetParamsAuthUnion() {}

type ConfigGetParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [ConfigGetParams]'s query parameters as `url.Values`.
func (r ConfigGetParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type ConfigUpdateParams struct {
	// Config is the request body. The json:"-" tag prevents double-encoding
	// because MarshalJSON serializes this field as the top-level JSON object.
	Config    Config  `json:"-"`
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [ConfigUpdateParams]'s query parameters as `url.Values`.
func (r ConfigUpdateParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

// MarshalJSON marshals the Config field for the request body
func (r ConfigUpdateParams) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Config)
}

type ConfigProviderListParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [ConfigProviderListParams]'s query parameters as `url.Values`.
func (r ConfigProviderListParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type ConfigProviderListResponse struct {
	Default   map[string]string `json:"default"`
	Providers []ConfigProvider  `json:"providers"`
}
