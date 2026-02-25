package opencode

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"

	"github.com/dominicnunez/opencode-sdk-go/internal/apijson"
	"github.com/dominicnunez/opencode-sdk-go/internal/apiquery"
	"github.com/dominicnunez/opencode-sdk-go/shared"
	"github.com/tidwall/gjson"
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

type Config struct {
	// JSON schema reference for configuration validation
	Schema string `json:"$schema"`
	// Agent configuration, see https://opencode.ai/docs/agent
	Agent ConfigAgent `json:"agent"`
	// @deprecated Use 'share' field instead. Share newly created sessions
	// automatically
	Autoshare bool `json:"autoshare"`
	// Automatically update to the latest version
	Autoupdate bool `json:"autoupdate"`
	// Command configuration, see https://opencode.ai/docs/commands
	Command map[string]ConfigCommand `json:"command"`
	// Disable providers that are loaded automatically
	DisabledProviders []string                   `json:"disabled_providers"`
	Experimental      ConfigExperimental         `json:"experimental"`
	Formatter         map[string]ConfigFormatter `json:"formatter"`
	// Additional instruction files or patterns to include
	Instructions []string `json:"instructions"`
	// Custom keybind configurations
	Keybinds KeybindsConfig `json:"keybinds"`
	// @deprecated Always uses stretch layout.
	Layout ConfigLayout         `json:"layout"`
	Lsp    map[string]ConfigLsp `json:"lsp"`
	// MCP (Model Context Protocol) server configurations
	Mcp map[string]ConfigMcp `json:"mcp"`
	// @deprecated Use `agent` field instead.
	Mode ConfigMode `json:"mode"`
	// Model to use in the format of provider/model, eg anthropic/claude-2
	Model      string           `json:"model"`
	Permission ConfigPermission `json:"permission"`
	Plugin     []string         `json:"plugin"`
	// Custom provider configurations and model overrides
	Provider map[string]ConfigProvider `json:"provider"`
	// Control sharing behavior:'manual' allows manual sharing via commands, 'auto'
	// enables automatic sharing, 'disabled' disables all sharing
	Share ConfigShare `json:"share"`
	// Small model to use for tasks like title generation in the format of
	// provider/model
	SmallModel string `json:"small_model"`
	Snapshot   bool   `json:"snapshot"`
	// Theme name to use for the interface
	Theme string          `json:"theme"`
	Tools map[string]bool `json:"tools"`
	// TUI specific settings
	Tui ConfigTui `json:"tui"`
	// Custom username to display in conversations instead of system username
	Username string        `json:"username"`
	Watcher  ConfigWatcher `json:"watcher"`
}

// Agent configuration, see https://opencode.ai/docs/agent
type ConfigAgent struct {
	Build       ConfigAgentBuild       `json:"build"`
	General     ConfigAgentGeneral     `json:"general"`
	Plan        ConfigAgentPlan        `json:"plan"`
	ExtraFields map[string]ConfigAgent `json:"-,extras"`
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
	ExtraFields map[string]interface{}     `json:"-,extras"`
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

// Union satisfied by [ConfigAgentBuildPermissionBashString] or
// [ConfigAgentBuildPermissionBashMap].
type ConfigAgentBuildPermissionBashUnion interface {
	implementsConfigAgentBuildPermissionBashUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ConfigAgentBuildPermissionBashUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(ConfigAgentBuildPermissionBashString("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ConfigAgentBuildPermissionBashMap{}),
		},
	)
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

func (r ConfigAgentBuildPermissionBashString) implementsConfigAgentBuildPermissionBashUnion() {}

type ConfigAgentBuildPermissionBashMap map[string]ConfigAgentBuildPermissionBashMapItem

func (r ConfigAgentBuildPermissionBashMap) implementsConfigAgentBuildPermissionBashUnion() {}

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
	ExtraFields map[string]interface{}       `json:"-,extras"`
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

// Union satisfied by [ConfigAgentGeneralPermissionBashString] or
// [ConfigAgentGeneralPermissionBashMap].
type ConfigAgentGeneralPermissionBashUnion interface {
	implementsConfigAgentGeneralPermissionBashUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ConfigAgentGeneralPermissionBashUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(ConfigAgentGeneralPermissionBashString("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ConfigAgentGeneralPermissionBashMap{}),
		},
	)
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

func (r ConfigAgentGeneralPermissionBashString) implementsConfigAgentGeneralPermissionBashUnion() {}

type ConfigAgentGeneralPermissionBashMap map[string]ConfigAgentGeneralPermissionBashMapItem

func (r ConfigAgentGeneralPermissionBashMap) implementsConfigAgentGeneralPermissionBashUnion() {}

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
	ExtraFields map[string]interface{}    `json:"-,extras"`
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

// Union satisfied by [ConfigAgentPlanPermissionBashString] or
// [ConfigAgentPlanPermissionBashMap].
type ConfigAgentPlanPermissionBashUnion interface {
	implementsConfigAgentPlanPermissionBashUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ConfigAgentPlanPermissionBashUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(ConfigAgentPlanPermissionBashString("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ConfigAgentPlanPermissionBashMap{}),
		},
	)
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

func (r ConfigAgentPlanPermissionBashString) implementsConfigAgentPlanPermissionBashUnion() {}

type ConfigAgentPlanPermissionBashMap map[string]ConfigAgentPlanPermissionBashMapItem

func (r ConfigAgentPlanPermissionBashMap) implementsConfigAgentPlanPermissionBashUnion() {}

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
	Template    string            `json:"template,required"`
	Agent       string            `json:"agent"`
	Description string            `json:"description"`
	Model       string            `json:"model"`
	Subtask     bool              `json:"subtask"`
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
	Command     []string                             `json:"command,required"`
	Environment map[string]string                    `json:"environment"`
}

type ConfigExperimentalHookSessionCompleted struct {
	Command     []string                                   `json:"command,required"`
	Environment map[string]string                          `json:"environment"`
}

type ConfigFormatter struct {
	Command     []string            `json:"command"`
	Disabled    bool                `json:"disabled"`
	Environment map[string]string   `json:"environment"`
	Extensions  []string            `json:"extensions"`
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

type ConfigLsp struct {
	// This field can have the runtime type of [[]string].
	Command  interface{} `json:"command"`
	Disabled bool        `json:"disabled"`
	// This field can have the runtime type of [map[string]string].
	Env interface{} `json:"env"`
	// This field can have the runtime type of [[]string].
	Extensions interface{} `json:"extensions"`
	// This field can have the runtime type of [map[string]interface{}].
	Initialization interface{}   `json:"initialization"`
	union          ConfigLspUnion
}

func (r *ConfigLsp) UnmarshalJSON(data []byte) (err error) {
	*r = ConfigLsp{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [ConfigLspUnion] interface which you can cast to the specific
// types for more type safety.
//
// Possible runtime types of the union are [ConfigLspDisabled], [ConfigLspObject].
func (r ConfigLsp) AsUnion() ConfigLspUnion {
	return r.union
}

// Union satisfied by [ConfigLspDisabled] or [ConfigLspObject].
type ConfigLspUnion interface {
	implementsConfigLsp()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ConfigLspUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ConfigLspDisabled{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ConfigLspObject{}),
		},
	)
}

type ConfigLspDisabled struct {
	Disabled ConfigLspDisabledDisabled `json:"disabled,required"`
}

func (r ConfigLspDisabled) implementsConfigLsp() {}

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
	Command        []string               `json:"command,required"`
	Disabled       bool                   `json:"disabled"`
	Env            map[string]string      `json:"env"`
	Extensions     []string               `json:"extensions"`
	Initialization map[string]interface{} `json:"initialization"`
}

func (r ConfigLspObject) implementsConfigLsp() {}

type ConfigMcp struct {
	// Type of MCP server connection
	Type ConfigMcpType `json:"type,required"`
	// This field can have the runtime type of [[]string].
	Command interface{} `json:"command"`
	// Enable or disable the MCP server on startup
	Enabled bool `json:"enabled"`
	// This field can have the runtime type of [map[string]string].
	Environment interface{} `json:"environment"`
	// This field can have the runtime type of [map[string]string].
	Headers interface{} `json:"headers"`
	// URL of the remote MCP server
	URL   string        `json:"url"`
	union ConfigMcpUnion
}



func (r *ConfigMcp) UnmarshalJSON(data []byte) (err error) {
	*r = ConfigMcp{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [ConfigMcpUnion] interface which you can cast to the specific
// types for more type safety.
//
// Possible runtime types of the union are [McpLocalConfig], [McpRemoteConfig].
func (r ConfigMcp) AsUnion() ConfigMcpUnion {
	return r.union
}

// Union satisfied by [McpLocalConfig] or [McpRemoteConfig].
type ConfigMcpUnion interface {
	implementsConfigMcp()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ConfigMcpUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(McpLocalConfig{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(McpRemoteConfig{}),
		},
	)
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
	Build       ConfigModeBuild       `json:"build"`
	Plan        ConfigModePlan        `json:"plan"`
	ExtraFields map[string]ConfigMode `json:"-,extras"`
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
	ExtraFields map[string]interface{}    `json:"-,extras"`
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




// Union satisfied by [ConfigModeBuildPermissionBashString] or
// [ConfigModeBuildPermissionBashMap].
type ConfigModeBuildPermissionBashUnion interface {
	implementsConfigModeBuildPermissionBashUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ConfigModeBuildPermissionBashUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(ConfigModeBuildPermissionBashString("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ConfigModeBuildPermissionBashMap{}),
		},
	)
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

func (r ConfigModeBuildPermissionBashString) implementsConfigModeBuildPermissionBashUnion() {}

type ConfigModeBuildPermissionBashMap map[string]ConfigModeBuildPermissionBashMapItem

func (r ConfigModeBuildPermissionBashMap) implementsConfigModeBuildPermissionBashUnion() {}

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
	ExtraFields map[string]interface{}   `json:"-,extras"`
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




// Union satisfied by [ConfigModePlanPermissionBashString] or
// [ConfigModePlanPermissionBashMap].
type ConfigModePlanPermissionBashUnion interface {
	implementsConfigModePlanPermissionBashUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ConfigModePlanPermissionBashUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(ConfigModePlanPermissionBashString("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ConfigModePlanPermissionBashMap{}),
		},
	)
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

func (r ConfigModePlanPermissionBashString) implementsConfigModePlanPermissionBashUnion() {}

type ConfigModePlanPermissionBashMap map[string]ConfigModePlanPermissionBashMapItem

func (r ConfigModePlanPermissionBashMap) implementsConfigModePlanPermissionBashUnion() {}

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




// Union satisfied by [ConfigPermissionBashString] or [ConfigPermissionBashMap].
type ConfigPermissionBashUnion interface {
	implementsConfigPermissionBashUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ConfigPermissionBashUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(ConfigPermissionBashString("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ConfigPermissionBashMap{}),
		},
	)
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

func (r ConfigPermissionBashString) implementsConfigPermissionBashUnion() {}

type ConfigPermissionBashMap map[string]ConfigPermissionBashMapItem

func (r ConfigPermissionBashMap) implementsConfigPermissionBashUnion() {}

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
	Input      float64                      `json:"input,required"`
	Output     float64                      `json:"output,required"`
	CacheRead  float64                      `json:"cache_read"`
	CacheWrite float64                      `json:"cache_write"`
}




type ConfigProviderModelsLimit struct {
	Context float64                       `json:"context,required"`
	Output  float64                       `json:"output,required"`
}




type ConfigProviderModelsModalities struct {
	Input  []ConfigProviderModelsModalitiesInput  `json:"input,required"`
	Output []ConfigProviderModelsModalitiesOutput `json:"output,required"`
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
	Npm  string                           `json:"npm,required"`
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
	Timeout     ConfigProviderOptionsTimeoutUnion `json:"timeout"`
	ExtraFields map[string]interface{}            `json:"-,extras"`
}




// Timeout in milliseconds for requests to this provider. Default is 300000 (5
// minutes). Set to false to disable timeout.
//
// Union satisfied by [shared.UnionInt] or [shared.UnionBool].
type ConfigProviderOptionsTimeoutUnion interface {
	ImplementsConfigProviderOptionsTimeoutUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ConfigProviderOptionsTimeoutUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionInt(0)),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.True,
			Type:       reflect.TypeOf(shared.UnionBool(false)),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.False,
			Type:       reflect.TypeOf(shared.UnionBool(false)),
		},
	)
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
	ScrollSpeed float64       `json:"scroll_speed"`
}




type ConfigWatcher struct {
	Ignore []string          `json:"ignore"`
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
	ToolDetails string             `json:"tool_details"`
}




type McpLocalConfig struct {
	// Command and arguments to run the MCP server
	Command []string `json:"command,required"`
	// Type of MCP server connection
	Type McpLocalConfigType `json:"type,required"`
	// Enable or disable the MCP server on startup
	Enabled bool `json:"enabled"`
	// Environment variables to set when running the MCP server
	Environment map[string]string  `json:"environment"`
}




func (r McpLocalConfig) implementsConfigMcp() {}

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
	Type McpRemoteConfigType `json:"type,required"`
	// URL of the remote MCP server
	URL string `json:"url,required"`
	// Enable or disable the MCP server on startup
	Enabled bool `json:"enabled"`
	// Headers to send with the request
	Headers map[string]string   `json:"headers"`
}




func (r McpRemoteConfig) implementsConfigMcp() {}

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
	AuthTypeOAuth      AuthType = "oauth"
	AuthTypeAPI        AuthType = "api"
	AuthTypeWellKnown  AuthType = "wellknown"
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
	a.raw = data
	return nil
}

// AsOAuth returns the OAuth variant if the type is oauth
func (a Auth) AsOAuth() (*OAuth, bool) {
	if a.Type != AuthTypeOAuth {
		return nil, false
	}
	var oauth OAuth
	if err := json.Unmarshal(a.raw, &oauth); err != nil {
		return nil, false
	}
	return &oauth, true
}

// AsAPI returns the ApiAuth variant if the type is api
func (a Auth) AsAPI() (*ApiAuth, bool) {
	if a.Type != AuthTypeAPI {
		return nil, false
	}
	var apiAuth ApiAuth
	if err := json.Unmarshal(a.raw, &apiAuth); err != nil {
		return nil, false
	}
	return &apiAuth, true
}

// AsWellKnown returns the WellKnownAuth variant if the type is wellknown
func (a Auth) AsWellKnown() (*WellKnownAuth, bool) {
	if a.Type != AuthTypeWellKnown {
		return nil, false
	}
	var wellKnown WellKnownAuth
	if err := json.Unmarshal(a.raw, &wellKnown); err != nil {
		return nil, false
	}
	return &wellKnown, true
}

// OAuth represents OAuth authentication credentials
type OAuth struct {
	Type    AuthType `json:"type"`
	Refresh string   `json:"refresh"`
	Access  string   `json:"access"`
	Expires int64    `json:"expires"`
}

// ApiAuth represents API key authentication
type ApiAuth struct {
	Type AuthType `json:"type"`
	Key  string   `json:"key"`
}

// WellKnownAuth represents well-known authentication
type WellKnownAuth struct {
	Type  AuthType `json:"type"`
	Key   string   `json:"key"`
	Token string   `json:"token"`
}

type ConfigGetParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [ConfigGetParams]'s query parameters as `url.Values`.
func (r ConfigGetParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
