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

type SessionService struct {
	client      *Client
	Permissions *SessionPermissionService
}

func (s *SessionService) Create(ctx context.Context, params *SessionNewParams) (*Session, error) {
	if params == nil {
		params = &SessionNewParams{}
	}
	var result Session
	err := s.client.do(ctx, http.MethodPost, "session", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SessionService) Update(ctx context.Context, id string, params *SessionUpdateParams) (*Session, error) {
	if id == "" {
		return nil, errors.New("missing required id parameter")
	}
	if params == nil {
		params = &SessionUpdateParams{}
	}
	var result Session
	err := s.client.do(ctx, http.MethodPatch, "session/"+id, params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SessionService) List(ctx context.Context, params *SessionListParams) ([]Session, error) {
	if params == nil {
		params = &SessionListParams{}
	}
	var result []Session
	err := s.client.do(ctx, http.MethodGet, "session", params, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SessionService) Get(ctx context.Context, id string, params *SessionGetParams) (*Session, error) {
	if id == "" {
		return nil, errors.New("missing required id parameter")
	}
	if params == nil {
		params = &SessionGetParams{}
	}
	var result Session
	err := s.client.do(ctx, http.MethodGet, "session/"+id, params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SessionService) Delete(ctx context.Context, id string, params *SessionDeleteParams) error {
	if id == "" {
		return errors.New("missing required id parameter")
	}
	if params == nil {
		params = &SessionDeleteParams{}
	}
	return s.client.do(ctx, http.MethodDelete, "session/"+id, params, nil)
}

func (s *SessionService) Abort(ctx context.Context, id string, params *SessionAbortParams) error {
	if id == "" {
		return errors.New("missing required id parameter")
	}
	if params == nil {
		params = &SessionAbortParams{}
	}
	return s.client.do(ctx, http.MethodPost, "session/"+id+"/abort", params, nil)
}

func (s *SessionService) Children(ctx context.Context, id string, params *SessionChildrenParams) ([]Session, error) {
	if id == "" {
		return nil, errors.New("missing required id parameter")
	}
	if params == nil {
		params = &SessionChildrenParams{}
	}
	var result []Session
	err := s.client.do(ctx, http.MethodGet, "session/"+id+"/children", params, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SessionService) Command(ctx context.Context, id string, params *SessionCommandParams) (*SessionCommandResponse, error) {
	if id == "" {
		return nil, errors.New("missing required id parameter")
	}
	if params == nil {
		params = &SessionCommandParams{}
	}
	var result SessionCommandResponse
	err := s.client.do(ctx, http.MethodPost, "session/"+id+"/command", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SessionService) Init(ctx context.Context, id string, params *SessionInitParams) (bool, error) {
	if id == "" {
		return false, errors.New("missing required id parameter")
	}
	if params == nil {
		params = &SessionInitParams{}
	}
	var result bool
	err := s.client.do(ctx, http.MethodPost, "session/"+id+"/init", params, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *SessionService) Message(ctx context.Context, id string, messageID string, params *SessionMessageParams) (*SessionMessageResponse, error) {
	if id == "" {
		return nil, errors.New("missing required id parameter")
	}
	if messageID == "" {
		return nil, errors.New("missing required messageID parameter")
	}
	if params == nil {
		params = &SessionMessageParams{}
	}
	var result SessionMessageResponse
	err := s.client.do(ctx, http.MethodGet, fmt.Sprintf("session/%s/message/%s", id, messageID), params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SessionService) Messages(ctx context.Context, id string, params *SessionMessagesParams) ([]SessionMessagesResponse, error) {
	if id == "" {
		return nil, errors.New("missing required id parameter")
	}
	if params == nil {
		params = &SessionMessagesParams{}
	}
	var result []SessionMessagesResponse
	err := s.client.do(ctx, http.MethodGet, "session/"+id+"/message", params, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SessionService) Prompt(ctx context.Context, id string, params *SessionPromptParams) (*SessionPromptResponse, error) {
	if id == "" {
		return nil, errors.New("missing required id parameter")
	}
	if params == nil {
		params = &SessionPromptParams{}
	}
	var result SessionPromptResponse
	err := s.client.do(ctx, http.MethodPost, "session/"+id+"/message", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SessionService) Revert(ctx context.Context, id string, params *SessionRevertParams) (*Session, error) {
	if id == "" {
		return nil, errors.New("missing required id parameter")
	}
	if params == nil {
		params = &SessionRevertParams{}
	}
	var result Session
	err := s.client.do(ctx, http.MethodPost, "session/"+id+"/revert", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SessionService) Share(ctx context.Context, id string, params *SessionShareParams) (*Session, error) {
	if id == "" {
		return nil, errors.New("missing required id parameter")
	}
	if params == nil {
		params = &SessionShareParams{}
	}
	var result Session
	err := s.client.do(ctx, http.MethodPost, "session/"+id+"/share", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type AgentPart struct {
	ID        string          `json:"id,required"`
	MessageID string          `json:"messageID,required"`
	Name      string          `json:"name,required"`
	SessionID string          `json:"sessionID,required"`
	Type      AgentPartType   `json:"type,required"`
	Source    AgentPartSource `json:"source"`
}

type AgentPartType string

const (
	AgentPartTypeAgent AgentPartType = "agent"
)

func (r AgentPartType) IsKnown() bool {
	switch r {
	case AgentPartTypeAgent:
		return true
	}
	return false
}

type AgentPartSource struct {
	End   int64               `json:"end,required"`
	Start int64               `json:"start,required"`
	Value string              `json:"value,required"`
}

type AgentPartInputParam struct {
	Name string `json:"name,required"`
	Type AgentPartInputType `json:"type,required"`
	ID *string `json:"id,omitempty"`
	Source *AgentPartInputSourceParam `json:"source,omitempty"`
}

func (r AgentPartInputParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r AgentPartInputParam) implementsSessionPromptParamsPartUnion() {}

type AgentPartInputType string

const (
	AgentPartInputTypeAgent AgentPartInputType = "agent"
)

func (r AgentPartInputType) IsKnown() bool {
	switch r {
	case AgentPartInputTypeAgent:
		return true
	}
	return false
}

type AgentPartInputSourceParam struct {
	End int64 `json:"end,required"`
	Start int64 `json:"start,required"`
	Value string `json:"value,required"`
}

func (r AgentPartInputSourceParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type AssistantMessage struct {
	ID         string                 `json:"id,required"`
	Cost       float64                `json:"cost,required"`
	Mode       string                 `json:"mode,required"`
	ModelID    string                 `json:"modelID,required"`
	ParentID   string                 `json:"parentID,required"`
	Path       AssistantMessagePath   `json:"path,required"`
	ProviderID string                 `json:"providerID,required"`
	Role       AssistantMessageRole   `json:"role,required"`
	SessionID  string                 `json:"sessionID,required"`
	System     []string               `json:"system,required"`
	Time       AssistantMessageTime   `json:"time,required"`
	Tokens     AssistantMessageTokens `json:"tokens,required"`
	Error      AssistantMessageError  `json:"error"`
	Summary    bool                   `json:"summary"`
}

type AssistantMessagePath struct {
	Cwd  string                   `json:"cwd,required"`
	Root string                   `json:"root,required"`
}

type AssistantMessageRole string

const (
	AssistantMessageRoleAssistant AssistantMessageRole = "assistant"
)

func (r AssistantMessageRole) IsKnown() bool {
	switch r {
	case AssistantMessageRoleAssistant:
		return true
	}
	return false
}

type AssistantMessageTime struct {
	Created   float64                  `json:"created,required"`
	Completed float64                  `json:"completed"`
}

type AssistantMessageTokens struct {
	Cache     AssistantMessageTokensCache `json:"cache,required"`
	Input     float64                     `json:"input,required"`
	Output    float64                     `json:"output,required"`
	Reasoning float64                     `json:"reasoning,required"`
}

type AssistantMessageTokensCache struct {
	Read  float64                         `json:"read,required"`
	Write float64                         `json:"write,required"`
}

type AssistantMessageError struct {
	// This field can have the runtime type of [shared.ProviderAuthErrorData],
	// [shared.UnknownErrorData], [interface{}], [shared.MessageAbortedErrorData],
	// [AssistantMessageErrorAPIErrorData].
	Data  interface{}               `json:"data,required"`
	Name  AssistantMessageErrorName `json:"name,required"`
	union AssistantMessageErrorUnion
}

func (r *AssistantMessageError) UnmarshalJSON(data []byte) (err error) {
	*r = AssistantMessageError{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [AssistantMessageErrorUnion] interface which you can cast to
// the specific types for more type safety.
//
// Possible runtime types of the union are [shared.ProviderAuthError],
// [shared.UnknownError], [AssistantMessageErrorMessageOutputLengthError],
// [shared.MessageAbortedError], [AssistantMessageErrorAPIError].
func (r AssistantMessageError) AsUnion() AssistantMessageErrorUnion {
	return r.union
}

// Union satisfied by [shared.ProviderAuthError], [shared.UnknownError],
// [AssistantMessageErrorMessageOutputLengthError], [shared.MessageAbortedError] or
// [AssistantMessageErrorAPIError].
type AssistantMessageErrorUnion interface {
	ImplementsAssistantMessageError()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*AssistantMessageErrorUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(shared.ProviderAuthError{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(shared.UnknownError{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(AssistantMessageErrorMessageOutputLengthError{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(shared.MessageAbortedError{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(AssistantMessageErrorAPIError{}),
		},
	)
}

type AssistantMessageErrorMessageOutputLengthError struct {
	Data interface{}                                       `json:"data,required"`
	Name AssistantMessageErrorMessageOutputLengthErrorName `json:"name,required"`
}

func (r AssistantMessageErrorMessageOutputLengthError) ImplementsAssistantMessageError() {}

type AssistantMessageErrorMessageOutputLengthErrorName string

const (
	AssistantMessageErrorMessageOutputLengthErrorNameMessageOutputLengthError AssistantMessageErrorMessageOutputLengthErrorName = "MessageOutputLengthError"
)

func (r AssistantMessageErrorMessageOutputLengthErrorName) IsKnown() bool {
	switch r {
	case AssistantMessageErrorMessageOutputLengthErrorNameMessageOutputLengthError:
		return true
	}
	return false
}

type AssistantMessageErrorAPIError struct {
	Data AssistantMessageErrorAPIErrorData `json:"data,required"`
	Name AssistantMessageErrorAPIErrorName `json:"name,required"`
}

func (r AssistantMessageErrorAPIError) ImplementsAssistantMessageError() {}

type AssistantMessageErrorAPIErrorData struct {
	IsRetryable     bool                                  `json:"isRetryable,required"`
	Message         string                                `json:"message,required"`
	ResponseBody    string                                `json:"responseBody"`
	ResponseHeaders map[string]string                     `json:"responseHeaders"`
	StatusCode      float64                               `json:"statusCode"`
}

type AssistantMessageErrorAPIErrorName string

const (
	AssistantMessageErrorAPIErrorNameAPIError AssistantMessageErrorAPIErrorName = "APIError"
)

func (r AssistantMessageErrorAPIErrorName) IsKnown() bool {
	switch r {
	case AssistantMessageErrorAPIErrorNameAPIError:
		return true
	}
	return false
}

type AssistantMessageErrorName string

const (
	AssistantMessageErrorNameProviderAuthError        AssistantMessageErrorName = "ProviderAuthError"
	AssistantMessageErrorNameUnknownError             AssistantMessageErrorName = "UnknownError"
	AssistantMessageErrorNameMessageOutputLengthError AssistantMessageErrorName = "MessageOutputLengthError"
	AssistantMessageErrorNameMessageAbortedError      AssistantMessageErrorName = "MessageAbortedError"
	AssistantMessageErrorNameAPIError                 AssistantMessageErrorName = "APIError"
)

func (r AssistantMessageErrorName) IsKnown() bool {
	switch r {
	case AssistantMessageErrorNameProviderAuthError, AssistantMessageErrorNameUnknownError, AssistantMessageErrorNameMessageOutputLengthError, AssistantMessageErrorNameMessageAbortedError, AssistantMessageErrorNameAPIError:
		return true
	}
	return false
}

type FilePart struct {
	ID        string         `json:"id,required"`
	MessageID string         `json:"messageID,required"`
	Mime      string         `json:"mime,required"`
	SessionID string         `json:"sessionID,required"`
	Type      FilePartType   `json:"type,required"`
	URL       string         `json:"url,required"`
	Filename  string         `json:"filename"`
	Source    FilePartSource `json:"source"`
}

type FilePartType string

const (
	FilePartTypeFile FilePartType = "file"
)

func (r FilePartType) IsKnown() bool {
	switch r {
	case FilePartTypeFile:
		return true
	}
	return false
}

type FilePartInputParam struct {
	Mime string `json:"mime,required"`
	Type FilePartInputType `json:"type,required"`
	URL string `json:"url,required"`
	ID *string `json:"id,omitempty"`
	Filename *string `json:"filename,omitempty"`
	Source *FilePartSourceUnionParam `json:"source,omitempty"`
}

func (r FilePartInputParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r FilePartInputParam) implementsSessionPromptParamsPartUnion() {}

type FilePartInputType string

const (
	FilePartInputTypeFile FilePartInputType = "file"
)

func (r FilePartInputType) IsKnown() bool {
	switch r {
	case FilePartInputTypeFile:
		return true
	}
	return false
}

// FilePartSource is either FileSource or SymbolSource, discriminated by Type.
type FilePartSource struct {
	Type FilePartSourceType `json:"type,required"`
	// Embed raw JSON for lazy decode
	raw []byte `json:"-"`
}

func (r *FilePartSource) UnmarshalJSON(data []byte) error {
	// Peek at discriminator
	var peek struct {
		Type FilePartSourceType `json:"type"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return err
	}
	r.Type = peek.Type
	r.raw = data
	return nil
}

// AsFile returns the source as a FileSource if Type is "file".
func (r FilePartSource) AsFile() (*FileSource, bool) {
	if r.Type != FilePartSourceTypeFile {
		return nil, false
	}
	var src FileSource
	if err := json.Unmarshal(r.raw, &src); err != nil {
		return nil, false
	}
	return &src, true
}

// AsSymbol returns the source as a SymbolSource if Type is "symbol".
func (r FilePartSource) AsSymbol() (*SymbolSource, bool) {
	if r.Type != FilePartSourceTypeSymbol {
		return nil, false
	}
	var src SymbolSource
	if err := json.Unmarshal(r.raw, &src); err != nil {
		return nil, false
	}
	return &src, true
}

type FilePartSourceType string

const (
	FilePartSourceTypeFile   FilePartSourceType = "file"
	FilePartSourceTypeSymbol FilePartSourceType = "symbol"
)

func (r FilePartSourceType) IsKnown() bool {
	switch r {
	case FilePartSourceTypeFile, FilePartSourceTypeSymbol:
		return true
	}
	return false
}

type FilePartSourceParam struct {
	Path string `json:"path,required"`
	Text FilePartSourceTextParam `json:"text,required"`
	Type FilePartSourceType `json:"type,required"`
	Kind *int64 `json:"kind,omitempty"`
	Name *string `json:"name,omitempty"`
	Range any `json:"range,omitempty"`
}

func (r FilePartSourceParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r FilePartSourceParam) implementsFilePartSourceUnionParam() {}

// Satisfied by [FileSourceParam], [SymbolSourceParam], [FilePartSourceParam].
type FilePartSourceUnionParam interface {
	implementsFilePartSourceUnionParam()
}

type FilePartSourceText struct {
	End   int64                  `json:"end,required"`
	Start int64                  `json:"start,required"`
	Value string                 `json:"value,required"`
}

type FilePartSourceTextParam struct {
	End int64 `json:"end,required"`
	Start int64 `json:"start,required"`
	Value string `json:"value,required"`
}

func (r FilePartSourceTextParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type FileSource struct {
	Path string             `json:"path,required"`
	Text FilePartSourceText `json:"text,required"`
	Type FileSourceType     `json:"type,required"`
}

type FileSourceType string

const (
	FileSourceTypeFile FileSourceType = "file"
)

func (r FileSourceType) IsKnown() bool {
	switch r {
	case FileSourceTypeFile:
		return true
	}
	return false
}

type FileSourceParam struct {
	Path string `json:"path,required"`
	Text FilePartSourceTextParam `json:"text,required"`
	Type FileSourceType `json:"type,required"`
}

func (r FileSourceParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r FileSourceParam) implementsFilePartSourceUnionParam() {}

// Message is either UserMessage or AssistantMessage, discriminated by Role.
type Message struct {
	ID        string      `json:"id,required"`
	Role      MessageRole `json:"role,required"`
	SessionID string      `json:"sessionID,required"`
	// Embed raw JSON for lazy decode
	raw []byte `json:"-"`
}

func (r *Message) UnmarshalJSON(data []byte) error {
	// Peek at common fields including discriminator
	var peek struct {
		ID        string      `json:"id"`
		Role      MessageRole `json:"role"`
		SessionID string      `json:"sessionID"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return err
	}
	r.ID = peek.ID
	r.Role = peek.Role
	r.SessionID = peek.SessionID
	r.raw = data
	return nil
}

// AsUser returns the UserMessage if the role is "user".
func (r Message) AsUser() (*UserMessage, bool) {
	if r.Role != MessageRoleUser {
		return nil, false
	}
	var msg UserMessage
	if err := json.Unmarshal(r.raw, &msg); err != nil {
		return nil, false
	}
	return &msg, true
}

// AsAssistant returns the AssistantMessage if the role is "assistant".
func (r Message) AsAssistant() (*AssistantMessage, bool) {
	if r.Role != MessageRoleAssistant {
		return nil, false
	}
	var msg AssistantMessage
	if err := json.Unmarshal(r.raw, &msg); err != nil {
		return nil, false
	}
	return &msg, true
}

type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
)

func (r MessageRole) IsKnown() bool {
	switch r {
	case MessageRoleUser, MessageRoleAssistant:
		return true
	}
	return false
}

// Part is a discriminated union type representing different kinds of message parts.
// Use the As* methods to access the specific part type based on the Type field.
type Part struct {
	ID        string   `json:"id,required"`
	MessageID string   `json:"messageID,required"`
	SessionID string   `json:"sessionID,required"`
	Type      PartType `json:"type,required"`
	raw       json.RawMessage
}

func (r *Part) UnmarshalJSON(data []byte) error {
	// Peek at discriminator to determine the type
	var peek struct {
		ID        string   `json:"id"`
		MessageID string   `json:"messageID"`
		SessionID string   `json:"sessionID"`
		Type      PartType `json:"type"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return err
	}
	r.ID = peek.ID
	r.MessageID = peek.MessageID
	r.SessionID = peek.SessionID
	r.Type = peek.Type
	r.raw = data
	return nil
}

// AsText returns the part as a TextPart if Type is "text".
func (r Part) AsText() (*TextPart, bool) {
	if r.Type != PartTypeText {
		return nil, false
	}
	var part TextPart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, false
	}
	return &part, true
}

// AsReasoning returns the part as a ReasoningPart if Type is "reasoning".
func (r Part) AsReasoning() (*ReasoningPart, bool) {
	if r.Type != PartTypeReasoning {
		return nil, false
	}
	var part ReasoningPart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, false
	}
	return &part, true
}

// AsFile returns the part as a FilePart if Type is "file".
func (r Part) AsFile() (*FilePart, bool) {
	if r.Type != PartTypeFile {
		return nil, false
	}
	var part FilePart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, false
	}
	return &part, true
}

// AsTool returns the part as a ToolPart if Type is "tool".
func (r Part) AsTool() (*ToolPart, bool) {
	if r.Type != PartTypeTool {
		return nil, false
	}
	var part ToolPart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, false
	}
	return &part, true
}

// AsStepStart returns the part as a StepStartPart if Type is "step-start".
func (r Part) AsStepStart() (*StepStartPart, bool) {
	if r.Type != PartTypeStepStart {
		return nil, false
	}
	var part StepStartPart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, false
	}
	return &part, true
}

// AsStepFinish returns the part as a StepFinishPart if Type is "step-finish".
func (r Part) AsStepFinish() (*StepFinishPart, bool) {
	if r.Type != PartTypeStepFinish {
		return nil, false
	}
	var part StepFinishPart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, false
	}
	return &part, true
}

// AsSnapshot returns the part as a SnapshotPart if Type is "snapshot".
func (r Part) AsSnapshot() (*SnapshotPart, bool) {
	if r.Type != PartTypeSnapshot {
		return nil, false
	}
	var part SnapshotPart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, false
	}
	return &part, true
}

// AsPatch returns the part as a PartPatchPart if Type is "patch".
func (r Part) AsPatch() (*PartPatchPart, bool) {
	if r.Type != PartTypePatch {
		return nil, false
	}
	var part PartPatchPart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, false
	}
	return &part, true
}

// AsAgent returns the part as an AgentPart if Type is "agent".
func (r Part) AsAgent() (*AgentPart, bool) {
	if r.Type != PartTypeAgent {
		return nil, false
	}
	var part AgentPart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, false
	}
	return &part, true
}

// AsRetry returns the part as a PartRetryPart if Type is "retry".
func (r Part) AsRetry() (*PartRetryPart, bool) {
	if r.Type != PartTypeRetry {
		return nil, false
	}
	var part PartRetryPart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, false
	}
	return &part, true
}

type PartPatchPart struct {
	ID        string            `json:"id,required"`
	Files     []string          `json:"files,required"`
	Hash      string            `json:"hash,required"`
	MessageID string            `json:"messageID,required"`
	SessionID string            `json:"sessionID,required"`
	Type      PartPatchPartType `json:"type,required"`
}

type PartPatchPartType string

const (
	PartPatchPartTypePatch PartPatchPartType = "patch"
)

func (r PartPatchPartType) IsKnown() bool {
	switch r {
	case PartPatchPartTypePatch:
		return true
	}
	return false
}

type PartRetryPart struct {
	ID        string             `json:"id,required"`
	Attempt   float64            `json:"attempt,required"`
	Error     PartRetryPartError `json:"error,required"`
	MessageID string             `json:"messageID,required"`
	SessionID string             `json:"sessionID,required"`
	Time      PartRetryPartTime  `json:"time,required"`
	Type      PartRetryPartType  `json:"type,required"`
}

type PartRetryPartError struct {
	Data PartRetryPartErrorData `json:"data,required"`
	Name PartRetryPartErrorName `json:"name,required"`
}

type PartRetryPartErrorData struct {
	IsRetryable     bool                       `json:"isRetryable,required"`
	Message         string                     `json:"message,required"`
	ResponseBody    string                     `json:"responseBody"`
	ResponseHeaders map[string]string          `json:"responseHeaders"`
	StatusCode      float64                    `json:"statusCode"`
}

type PartRetryPartErrorName string

const (
	PartRetryPartErrorNameAPIError PartRetryPartErrorName = "APIError"
)

func (r PartRetryPartErrorName) IsKnown() bool {
	switch r {
	case PartRetryPartErrorNameAPIError:
		return true
	}
	return false
}

type PartRetryPartTime struct {
	Created float64               `json:"created,required"`
}

type PartRetryPartType string

const (
	PartRetryPartTypeRetry PartRetryPartType = "retry"
)

func (r PartRetryPartType) IsKnown() bool {
	switch r {
	case PartRetryPartTypeRetry:
		return true
	}
	return false
}

type PartType string

const (
	PartTypeText       PartType = "text"
	PartTypeReasoning  PartType = "reasoning"
	PartTypeFile       PartType = "file"
	PartTypeTool       PartType = "tool"
	PartTypeStepStart  PartType = "step-start"
	PartTypeStepFinish PartType = "step-finish"
	PartTypeSnapshot   PartType = "snapshot"
	PartTypePatch      PartType = "patch"
	PartTypeAgent      PartType = "agent"
	PartTypeRetry      PartType = "retry"
)

func (r PartType) IsKnown() bool {
	switch r {
	case PartTypeText, PartTypeReasoning, PartTypeFile, PartTypeTool, PartTypeStepStart, PartTypeStepFinish, PartTypeSnapshot, PartTypePatch, PartTypeAgent, PartTypeRetry:
		return true
	}
	return false
}

type ReasoningPart struct {
	ID        string                 `json:"id,required"`
	MessageID string                 `json:"messageID,required"`
	SessionID string                 `json:"sessionID,required"`
	Text      string                 `json:"text,required"`
	Time      ReasoningPartTime      `json:"time,required"`
	Type      ReasoningPartType      `json:"type,required"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type ReasoningPartTime struct {
	Start float64               `json:"start,required"`
	End   float64               `json:"end"`
}

type ReasoningPartType string

const (
	ReasoningPartTypeReasoning ReasoningPartType = "reasoning"
)

func (r ReasoningPartType) IsKnown() bool {
	switch r {
	case ReasoningPartTypeReasoning:
		return true
	}
	return false
}

type Session struct {
	ID        string         `json:"id,required"`
	Directory string         `json:"directory,required"`
	ProjectID string         `json:"projectID,required"`
	Time      SessionTime    `json:"time,required"`
	Title     string         `json:"title,required"`
	Version   string         `json:"version,required"`
	ParentID  string         `json:"parentID"`
	Revert    SessionRevert  `json:"revert"`
	Share     SessionShare   `json:"share"`
	Summary   SessionSummary `json:"summary"`
}

type SessionTime struct {
	Created    float64         `json:"created,required"`
	Updated    float64         `json:"updated,required"`
	Compacting float64         `json:"compacting"`
}

type SessionRevert struct {
	MessageID string            `json:"messageID,required"`
	Diff      string            `json:"diff"`
	PartID    string            `json:"partID"`
	Snapshot  string            `json:"snapshot"`
}

type SessionShare struct {
	URL  string           `json:"url,required"`
}

type SessionSummary struct {
	Diffs []SessionSummaryDiff `json:"diffs,required"`
}

type SessionSummaryDiff struct {
	Additions float64                `json:"additions,required"`
	After     string                 `json:"after,required"`
	Before    string                 `json:"before,required"`
	Deletions float64                `json:"deletions,required"`
	File      string                 `json:"file,required"`
}

type SnapshotPart struct {
	ID        string           `json:"id,required"`
	MessageID string           `json:"messageID,required"`
	SessionID string           `json:"sessionID,required"`
	Snapshot  string           `json:"snapshot,required"`
	Type      SnapshotPartType `json:"type,required"`
}

type SnapshotPartType string

const (
	SnapshotPartTypeSnapshot SnapshotPartType = "snapshot"
)

func (r SnapshotPartType) IsKnown() bool {
	switch r {
	case SnapshotPartTypeSnapshot:
		return true
	}
	return false
}

type StepFinishPart struct {
	ID        string               `json:"id,required"`
	Cost      float64              `json:"cost,required"`
	MessageID string               `json:"messageID,required"`
	Reason    string               `json:"reason,required"`
	SessionID string               `json:"sessionID,required"`
	Tokens    StepFinishPartTokens `json:"tokens,required"`
	Type      StepFinishPartType   `json:"type,required"`
	Snapshot  string               `json:"snapshot"`
}

type StepFinishPartTokens struct {
	Cache     StepFinishPartTokensCache `json:"cache,required"`
	Input     float64                   `json:"input,required"`
	Output    float64                   `json:"output,required"`
	Reasoning float64                   `json:"reasoning,required"`
}

type StepFinishPartTokensCache struct {
	Read  float64                       `json:"read,required"`
	Write float64                       `json:"write,required"`
}

type StepFinishPartType string

const (
	StepFinishPartTypeStepFinish StepFinishPartType = "step-finish"
)

func (r StepFinishPartType) IsKnown() bool {
	switch r {
	case StepFinishPartTypeStepFinish:
		return true
	}
	return false
}

type StepStartPart struct {
	ID        string            `json:"id,required"`
	MessageID string            `json:"messageID,required"`
	SessionID string            `json:"sessionID,required"`
	Type      StepStartPartType `json:"type,required"`
	Snapshot  string            `json:"snapshot"`
}

type StepStartPartType string

const (
	StepStartPartTypeStepStart StepStartPartType = "step-start"
)

func (r StepStartPartType) IsKnown() bool {
	switch r {
	case StepStartPartTypeStepStart:
		return true
	}
	return false
}

type SymbolSource struct {
	Kind  int64              `json:"kind,required"`
	Name  string             `json:"name,required"`
	Path  string             `json:"path,required"`
	Range SymbolSourceRange  `json:"range,required"`
	Text  FilePartSourceText `json:"text,required"`
	Type  SymbolSourceType   `json:"type,required"`
}

type SymbolSourceRange struct {
	End   SymbolSourceRangeEnd   `json:"end,required"`
	Start SymbolSourceRangeStart `json:"start,required"`
}

type SymbolSourceRangeEnd struct {
	Character float64                  `json:"character,required"`
	Line      float64                  `json:"line,required"`
}

type SymbolSourceRangeStart struct {
	Character float64                    `json:"character,required"`
	Line      float64                    `json:"line,required"`
}

type SymbolSourceType string

const (
	SymbolSourceTypeSymbol SymbolSourceType = "symbol"
)

func (r SymbolSourceType) IsKnown() bool {
	switch r {
	case SymbolSourceTypeSymbol:
		return true
	}
	return false
}

type SymbolSourceParam struct {
	Kind int64 `json:"kind,required"`
	Name string `json:"name,required"`
	Path string `json:"path,required"`
	Range SymbolSourceRangeParam `json:"range,required"`
	Text FilePartSourceTextParam `json:"text,required"`
	Type SymbolSourceType `json:"type,required"`
}

func (r SymbolSourceParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SymbolSourceParam) implementsFilePartSourceUnionParam() {}

type SymbolSourceRangeParam struct {
	End SymbolSourceRangeEndParam `json:"end,required"`
	Start SymbolSourceRangeStartParam `json:"start,required"`
}

func (r SymbolSourceRangeParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SymbolSourceRangeEndParam struct {
	Character float64 `json:"character,required"`
	Line float64 `json:"line,required"`
}

func (r SymbolSourceRangeEndParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SymbolSourceRangeStartParam struct {
	Character float64 `json:"character,required"`
	Line float64 `json:"line,required"`
}

func (r SymbolSourceRangeStartParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type TextPart struct {
	ID        string                 `json:"id,required"`
	MessageID string                 `json:"messageID,required"`
	SessionID string                 `json:"sessionID,required"`
	Text      string                 `json:"text,required"`
	Type      TextPartType           `json:"type,required"`
	Metadata  map[string]interface{} `json:"metadata"`
	Synthetic bool                   `json:"synthetic"`
	Time      TextPartTime           `json:"time"`
}

type TextPartType string

const (
	TextPartTypeText TextPartType = "text"
)

func (r TextPartType) IsKnown() bool {
	switch r {
	case TextPartTypeText:
		return true
	}
	return false
}

type TextPartTime struct {
	Start float64          `json:"start,required"`
	End   float64          `json:"end"`
}

type TextPartInputParam struct {
	Text string `json:"text,required"`
	Type TextPartInputType `json:"type,required"`
	ID *string `json:"id,omitempty"`
	Metadata *map[string]interface{} `json:"metadata,omitempty"`
	Synthetic *bool `json:"synthetic,omitempty"`
	Time *TextPartInputTimeParam `json:"time,omitempty"`
}

func (r TextPartInputParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r TextPartInputParam) implementsSessionPromptParamsPartUnion() {}

type TextPartInputType string

const (
	TextPartInputTypeText TextPartInputType = "text"
)

func (r TextPartInputType) IsKnown() bool {
	switch r {
	case TextPartInputTypeText:
		return true
	}
	return false
}

type TextPartInputTimeParam struct {
	Start float64 `json:"start,required"`
	End *float64 `json:"end,omitempty"`
}

func (r TextPartInputTimeParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ToolPart struct {
	ID        string                 `json:"id,required"`
	CallID    string                 `json:"callID,required"`
	MessageID string                 `json:"messageID,required"`
	SessionID string                 `json:"sessionID,required"`
	State     ToolPartState          `json:"state,required"`
	Tool      string                 `json:"tool,required"`
	Type      ToolPartType           `json:"type,required"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// ToolPartState is a discriminated union type representing the state of a tool.
// Use the As* methods to access the specific state type based on the Status discriminator.
type ToolPartState struct {
	Status ToolPartStateStatus `json:"status"`
	raw    json.RawMessage
}

func (r *ToolPartState) UnmarshalJSON(data []byte) error {
	// Peek at the discriminator (status field)
	var peek struct {
		Status ToolPartStateStatus `json:"status"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return err
	}
	r.Status = peek.Status
	r.raw = data
	return nil
}

// AsPending returns the state as ToolStatePending if Status is "pending".
func (r ToolPartState) AsPending() (*ToolStatePending, bool) {
	if r.Status != ToolPartStateStatusPending {
		return nil, false
	}
	var state ToolStatePending
	if err := json.Unmarshal(r.raw, &state); err != nil {
		return nil, false
	}
	return &state, true
}

// AsRunning returns the state as ToolStateRunning if Status is "running".
func (r ToolPartState) AsRunning() (*ToolStateRunning, bool) {
	if r.Status != ToolPartStateStatusRunning {
		return nil, false
	}
	var state ToolStateRunning
	if err := json.Unmarshal(r.raw, &state); err != nil {
		return nil, false
	}
	return &state, true
}

// AsCompleted returns the state as ToolStateCompleted if Status is "completed".
func (r ToolPartState) AsCompleted() (*ToolStateCompleted, bool) {
	if r.Status != ToolPartStateStatusCompleted {
		return nil, false
	}
	var state ToolStateCompleted
	if err := json.Unmarshal(r.raw, &state); err != nil {
		return nil, false
	}
	return &state, true
}

// AsError returns the state as ToolStateError if Status is "error".
func (r ToolPartState) AsError() (*ToolStateError, bool) {
	if r.Status != ToolPartStateStatusError {
		return nil, false
	}
	var state ToolStateError
	if err := json.Unmarshal(r.raw, &state); err != nil {
		return nil, false
	}
	return &state, true
}

type ToolPartStateStatus string

const (
	ToolPartStateStatusPending   ToolPartStateStatus = "pending"
	ToolPartStateStatusRunning   ToolPartStateStatus = "running"
	ToolPartStateStatusCompleted ToolPartStateStatus = "completed"
	ToolPartStateStatusError     ToolPartStateStatus = "error"
)

func (r ToolPartStateStatus) IsKnown() bool {
	switch r {
	case ToolPartStateStatusPending, ToolPartStateStatusRunning, ToolPartStateStatusCompleted, ToolPartStateStatusError:
		return true
	}
	return false
}

type ToolPartType string

const (
	ToolPartTypeTool ToolPartType = "tool"
)

func (r ToolPartType) IsKnown() bool {
	switch r {
	case ToolPartTypeTool:
		return true
	}
	return false
}

type ToolStateCompleted struct {
	Input       map[string]interface{}   `json:"input,required"`
	Metadata    map[string]interface{}   `json:"metadata,required"`
	Output      string                   `json:"output,required"`
	Status      ToolStateCompletedStatus `json:"status,required"`
	Time        ToolStateCompletedTime   `json:"time,required"`
	Title       string                   `json:"title,required"`
	Attachments []FilePart               `json:"attachments"`
}

type ToolStateCompletedStatus string

const (
	ToolStateCompletedStatusCompleted ToolStateCompletedStatus = "completed"
)

func (r ToolStateCompletedStatus) IsKnown() bool {
	switch r {
	case ToolStateCompletedStatusCompleted:
		return true
	}
	return false
}

type ToolStateCompletedTime struct {
	End       float64                    `json:"end,required"`
	Start     float64                    `json:"start,required"`
	Compacted float64                    `json:"compacted"`
}

type ToolStateError struct {
	Error    string                 `json:"error,required"`
	Input    map[string]interface{} `json:"input,required"`
	Status   ToolStateErrorStatus   `json:"status,required"`
	Time     ToolStateErrorTime     `json:"time,required"`
	Metadata map[string]interface{} `json:"metadata"`
}

type ToolStateErrorStatus string

const (
	ToolStateErrorStatusError ToolStateErrorStatus = "error"
)

func (r ToolStateErrorStatus) IsKnown() bool {
	switch r {
	case ToolStateErrorStatusError:
		return true
	}
	return false
}

type ToolStateErrorTime struct {
	End   float64                `json:"end,required"`
	Start float64                `json:"start,required"`
}

type ToolStatePending struct {
	Status ToolStatePendingStatus `json:"status,required"`
}

type ToolStatePendingStatus string

const (
	ToolStatePendingStatusPending ToolStatePendingStatus = "pending"
)

func (r ToolStatePendingStatus) IsKnown() bool {
	switch r {
	case ToolStatePendingStatusPending:
		return true
	}
	return false
}

type ToolStateRunning struct {
	Input    interface{}            `json:"input,required"`
	Status   ToolStateRunningStatus `json:"status,required"`
	Time     ToolStateRunningTime   `json:"time,required"`
	Metadata map[string]interface{} `json:"metadata"`
	Title    string                 `json:"title"`
}

type ToolStateRunningStatus string

const (
	ToolStateRunningStatusRunning ToolStateRunningStatus = "running"
)

func (r ToolStateRunningStatus) IsKnown() bool {
	switch r {
	case ToolStateRunningStatusRunning:
		return true
	}
	return false
}

type ToolStateRunningTime struct {
	Start float64                  `json:"start,required"`
}

type UserMessage struct {
	ID        string             `json:"id,required"`
	Role      UserMessageRole    `json:"role,required"`
	SessionID string             `json:"sessionID,required"`
	Time      UserMessageTime    `json:"time,required"`
	Summary   UserMessageSummary `json:"summary"`
}

type UserMessageRole string

const (
	UserMessageRoleUser UserMessageRole = "user"
)

func (r UserMessageRole) IsKnown() bool {
	switch r {
	case UserMessageRoleUser:
		return true
	}
	return false
}

type UserMessageTime struct {
	Created float64             `json:"created,required"`
}

type UserMessageSummary struct {
	Diffs []UserMessageSummaryDiff `json:"diffs,required"`
	Body  string                   `json:"body"`
	Title string                   `json:"title"`
}

type UserMessageSummaryDiff struct {
	Additions float64                    `json:"additions,required"`
	After     string                     `json:"after,required"`
	Before    string                     `json:"before,required"`
	Deletions float64                    `json:"deletions,required"`
	File      string                     `json:"file,required"`
}

type SessionCommandResponse struct {
	Info  AssistantMessage           `json:"info,required"`
	Parts []Part                     `json:"parts,required"`
}

type SessionMessageResponse struct {
	Info  Message                    `json:"info,required"`
	Parts []Part                     `json:"parts,required"`
}

type SessionMessagesResponse struct {
	Info  Message                     `json:"info,required"`
	Parts []Part                      `json:"parts,required"`
}

type SessionPromptResponse struct {
	Info  AssistantMessage          `json:"info,required"`
	Parts []Part                    `json:"parts,required"`
}

type SessionNewParams struct {
	Directory *string `query:"directory,omitempty"`
	ParentID *string `json:"parentID,omitempty"`
	Title *string `json:"title,omitempty"`
}

func (r SessionNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// URLQuery serializes [SessionNewParams]'s query parameters as `url.Values`.
func (r SessionNewParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionUpdateParams struct {
	Directory *string `query:"directory,omitempty"`
	Title *string `json:"title,omitempty"`
}

func (r SessionUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// URLQuery serializes [SessionUpdateParams]'s query parameters as `url.Values`.
func (r SessionUpdateParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionListParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionListParams]'s query parameters as `url.Values`.
func (r SessionListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionDeleteParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionDeleteParams]'s query parameters as `url.Values`.
func (r SessionDeleteParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionAbortParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionAbortParams]'s query parameters as `url.Values`.
func (r SessionAbortParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionChildrenParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionChildrenParams]'s query parameters as `url.Values`.
func (r SessionChildrenParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionCommandParams struct {
	Arguments string `json:"arguments,required"`
	Command string `json:"command,required"`
	Directory *string `query:"directory,omitempty"`
	Agent *string `json:"agent,omitempty"`
	MessageID *string `json:"messageID,omitempty"`
	Model *string `json:"model,omitempty"`
}

func (r SessionCommandParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// URLQuery serializes [SessionCommandParams]'s query parameters as `url.Values`.
func (r SessionCommandParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionGetParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionGetParams]'s query parameters as `url.Values`.
func (r SessionGetParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionInitParams struct {
	MessageID string `json:"messageID,required"`
	ModelID string `json:"modelID,required"`
	ProviderID string `json:"providerID,required"`
	Directory *string `query:"directory,omitempty"`
}

func (r SessionInitParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// URLQuery serializes [SessionInitParams]'s query parameters as `url.Values`.
func (r SessionInitParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionMessageParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionMessageParams]'s query parameters as `url.Values`.
func (r SessionMessageParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionMessagesParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionMessagesParams]'s query parameters as `url.Values`.
func (r SessionMessagesParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionPromptParams struct {
	Parts []SessionPromptParamsPartUnion `json:"parts,required"`
	Directory *string `query:"directory,omitempty"`
	Agent *string `json:"agent,omitempty"`
	MessageID *string `json:"messageID,omitempty"`
	Model *SessionPromptParamsModel `json:"model,omitempty"`
	NoReply *bool `json:"noReply,omitempty"`
	System *string `json:"system,omitempty"`
	Tools *map[string]bool `json:"tools,omitempty"`
}

func (r SessionPromptParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// URLQuery serializes [SessionPromptParams]'s query parameters as `url.Values`.
func (r SessionPromptParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionPromptParamsPart struct {
	Type SessionPromptParamsPartsType `json:"type,required"`
	ID *string `json:"id,omitempty"`
	Filename *string `json:"filename,omitempty"`
	Metadata any `json:"metadata,omitempty"`
	Mime *string `json:"mime,omitempty"`
	Name *string `json:"name,omitempty"`
	Source any `json:"source,omitempty"`
	Synthetic *bool `json:"synthetic,omitempty"`
	Text *string `json:"text,omitempty"`
	Time any `json:"time,omitempty"`
	URL *string `json:"url,omitempty"`
}

func (r SessionPromptParamsPart) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SessionPromptParamsPart) implementsSessionPromptParamsPartUnion() {}

// Satisfied by [TextPartInputParam], [FilePartInputParam], [AgentPartInputParam],
// [SessionPromptParamsPart].
type SessionPromptParamsPartUnion interface {
	implementsSessionPromptParamsPartUnion()
}

type SessionPromptParamsPartsType string

const (
	SessionPromptParamsPartsTypeText  SessionPromptParamsPartsType = "text"
	SessionPromptParamsPartsTypeFile  SessionPromptParamsPartsType = "file"
	SessionPromptParamsPartsTypeAgent SessionPromptParamsPartsType = "agent"
)

func (r SessionPromptParamsPartsType) IsKnown() bool {
	switch r {
	case SessionPromptParamsPartsTypeText, SessionPromptParamsPartsTypeFile, SessionPromptParamsPartsTypeAgent:
		return true
	}
	return false
}

type SessionPromptParamsModel struct {
	ModelID string `json:"modelID,required"`
	ProviderID string `json:"providerID,required"`
}

func (r SessionPromptParamsModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SessionRevertParams struct {
	MessageID string `json:"messageID,required"`
	Directory *string `query:"directory,omitempty"`
	PartID *string `json:"partID,omitempty"`
}

func (r SessionRevertParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// URLQuery serializes [SessionRevertParams]'s query parameters as `url.Values`.
func (r SessionRevertParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionShareParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionShareParams]'s query parameters as `url.Values`.
func (r SessionShareParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionShellParams struct {
	Agent string `json:"agent,required"`
	Command string `json:"command,required"`
	Directory *string `query:"directory,omitempty"`
}

func (r SessionShellParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// URLQuery serializes [SessionShellParams]'s query parameters as `url.Values`.
func (r SessionShellParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionSummarizeParams struct {
	ModelID string `json:"modelID,required"`
	ProviderID string `json:"providerID,required"`
	Directory *string `query:"directory,omitempty"`
}

func (r SessionSummarizeParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// URLQuery serializes [SessionSummarizeParams]'s query parameters as `url.Values`.
func (r SessionSummarizeParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionUnrevertParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionUnrevertParams]'s query parameters as `url.Values`.
func (r SessionUnrevertParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionUnshareParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionUnshareParams]'s query parameters as `url.Values`.
func (r SessionUnshareParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
