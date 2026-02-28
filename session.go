package opencode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dominicnunez/opencode-sdk-go/internal/queryparams"
	"github.com/dominicnunez/opencode-sdk-go/shared"
)

type SessionService struct {
	client      *Client
	Permissions *SessionPermissionService
}

func (s *SessionService) Create(ctx context.Context, params *SessionCreateParams) (*Session, error) {
	if params == nil {
		params = &SessionCreateParams{}
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

func (s *SessionService) Delete(ctx context.Context, id string, params *SessionDeleteParams) (bool, error) {
	if id == "" {
		return false, errors.New("missing required id parameter")
	}
	if params == nil {
		params = &SessionDeleteParams{}
	}
	var result bool
	err := s.client.do(ctx, http.MethodDelete, "session/"+id, params, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *SessionService) Abort(ctx context.Context, id string, params *SessionAbortParams) (bool, error) {
	if id == "" {
		return false, errors.New("missing required id parameter")
	}
	if params == nil {
		params = &SessionAbortParams{}
	}
	var result bool
	err := s.client.do(ctx, http.MethodPost, "session/"+id+"/abort", params, &result)
	if err != nil {
		return false, err
	}
	return result, nil
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

func (s *SessionService) Diff(ctx context.Context, id string, params *SessionDiffParams) ([]FileDiff, error) {
	if id == "" {
		return nil, errors.New("missing required id parameter")
	}
	if params == nil {
		params = &SessionDiffParams{}
	}
	var result []FileDiff
	err := s.client.do(ctx, http.MethodGet, "session/"+id+"/diff", params, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SessionService) Fork(ctx context.Context, id string, params *SessionForkParams) (*Session, error) {
	if id == "" {
		return nil, errors.New("missing required id parameter")
	}
	if params == nil {
		return nil, errors.New("params is required")
	}
	var result Session
	err := s.client.do(ctx, http.MethodPost, "session/"+id+"/fork", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SessionService) Shell(ctx context.Context, id string, params *SessionShellParams) (*AssistantMessage, error) {
	if id == "" {
		return nil, errors.New("missing required id parameter")
	}
	if params == nil {
		return nil, errors.New("params is required")
	}
	var result AssistantMessage
	err := s.client.do(ctx, http.MethodPost, "session/"+id+"/shell", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SessionService) Summarize(ctx context.Context, id string, params *SessionSummarizeParams) (bool, error) {
	if id == "" {
		return false, errors.New("missing required id parameter")
	}
	if params == nil {
		return false, errors.New("params is required")
	}
	var result bool
	err := s.client.do(ctx, http.MethodPost, "session/"+id+"/summarize", params, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *SessionService) Todo(ctx context.Context, id string, params *SessionTodoParams) ([]Todo, error) {
	if id == "" {
		return nil, errors.New("missing required id parameter")
	}
	if params == nil {
		params = &SessionTodoParams{}
	}
	var result []Todo
	err := s.client.do(ctx, http.MethodGet, "session/"+id+"/todo", params, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SessionService) Unrevert(ctx context.Context, id string, params *SessionUnrevertParams) (*Session, error) {
	if id == "" {
		return nil, errors.New("missing required id parameter")
	}
	if params == nil {
		params = &SessionUnrevertParams{}
	}
	var result Session
	err := s.client.do(ctx, http.MethodPost, "session/"+id+"/unrevert", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SessionService) Unshare(ctx context.Context, id string, params *SessionUnshareParams) (*Session, error) {
	if id == "" {
		return nil, errors.New("missing required id parameter")
	}
	if params == nil {
		params = &SessionUnshareParams{}
	}
	var result Session
	err := s.client.do(ctx, http.MethodDelete, "session/"+id+"/share", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type AgentPart struct {
	ID        string          `json:"id"`
	MessageID string          `json:"messageID"`
	Name      string          `json:"name"`
	SessionID string          `json:"sessionID"`
	Type      AgentPartType   `json:"type"`
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
	End   int64  `json:"end"`
	Start int64  `json:"start"`
	Value string `json:"value"`
}

type AgentPartInputParam struct {
	Name   string                     `json:"name"`
	Type   AgentPartInputType         `json:"type"`
	ID     *string                    `json:"id,omitempty"`
	Source *AgentPartInputSourceParam `json:"source,omitempty"`
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
	End   int64  `json:"end"`
	Start int64  `json:"start"`
	Value string `json:"value"`
}

type AssistantMessage struct {
	ID         string                 `json:"id"`
	Cost       float64                `json:"cost"`
	Mode       string                 `json:"mode"`
	ModelID    string                 `json:"modelID"`
	ParentID   string                 `json:"parentID"`
	Path       AssistantMessagePath   `json:"path"`
	ProviderID string                 `json:"providerID"`
	Role       AssistantMessageRole   `json:"role"`
	SessionID  string                 `json:"sessionID"`
	System     []string               `json:"system"`
	Time       AssistantMessageTime   `json:"time"`
	Tokens     AssistantMessageTokens `json:"tokens"`
	Error      AssistantMessageError  `json:"error"`
	Summary    bool                   `json:"summary"`
}

type AssistantMessagePath struct {
	Cwd  string `json:"cwd"`
	Root string `json:"root"`
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
	Created   float64 `json:"created"`
	Completed float64 `json:"completed"`
}

type AssistantMessageTokens struct {
	Cache     AssistantMessageTokensCache `json:"cache"`
	Input     int64                       `json:"input"`
	Output    int64                       `json:"output"`
	Reasoning int64                       `json:"reasoning"`
}

type AssistantMessageTokensCache struct {
	Read  int64 `json:"read"`
	Write int64 `json:"write"`
}

type AssistantMessageError struct {
	Name AssistantMessageErrorName `json:"name"`
	raw  json.RawMessage
}

func (r *AssistantMessageError) UnmarshalJSON(data []byte) error {
	// Peek at discriminator
	var peek struct {
		Name AssistantMessageErrorName `json:"name"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return err
	}
	r.Name = peek.Name
	r.raw = data
	return nil
}

func (r AssistantMessageError) AsProviderAuth() (*shared.ProviderAuthError, error) {
	if r.Name != AssistantMessageErrorNameProviderAuthError {
		return nil, ErrWrongVariant
	}
	var err shared.ProviderAuthError
	if err := json.Unmarshal(r.raw, &err); err != nil {
		return nil, fmt.Errorf("unmarshal %s Name: %w", r.Name, err)
	}
	return &err, nil
}

func (r AssistantMessageError) AsUnknown() (*shared.UnknownError, error) {
	if r.Name != AssistantMessageErrorNameUnknownError {
		return nil, ErrWrongVariant
	}
	var err shared.UnknownError
	if err := json.Unmarshal(r.raw, &err); err != nil {
		return nil, fmt.Errorf("unmarshal %s Name: %w", r.Name, err)
	}
	return &err, nil
}

func (r AssistantMessageError) AsOutputLength() (*AssistantMessageErrorMessageOutputLengthError, error) {
	if r.Name != AssistantMessageErrorNameMessageOutputLengthError {
		return nil, ErrWrongVariant
	}
	var err AssistantMessageErrorMessageOutputLengthError
	if err := json.Unmarshal(r.raw, &err); err != nil {
		return nil, fmt.Errorf("unmarshal %s Name: %w", r.Name, err)
	}
	return &err, nil
}

func (r AssistantMessageError) AsAborted() (*shared.MessageAbortedError, error) {
	if r.Name != AssistantMessageErrorNameMessageAbortedError {
		return nil, ErrWrongVariant
	}
	var err shared.MessageAbortedError
	if err := json.Unmarshal(r.raw, &err); err != nil {
		return nil, fmt.Errorf("unmarshal %s Name: %w", r.Name, err)
	}
	return &err, nil
}

func (r AssistantMessageError) AsAPI() (*AssistantMessageErrorAPIError, error) {
	if r.Name != AssistantMessageErrorNameAPIError {
		return nil, ErrWrongVariant
	}
	var err AssistantMessageErrorAPIError
	if err := json.Unmarshal(r.raw, &err); err != nil {
		return nil, fmt.Errorf("unmarshal %s Name: %w", r.Name, err)
	}
	return &err, nil
}

func (r AssistantMessageError) MarshalJSON() ([]byte, error) {
	if r.raw == nil {
		return []byte("null"), nil
	}
	return r.raw, nil
}

type AssistantMessageErrorMessageOutputLengthError struct {
	Data json.RawMessage                                   `json:"data"`
	Name AssistantMessageErrorMessageOutputLengthErrorName `json:"name"`
}

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
	Data AssistantMessageErrorAPIErrorData `json:"data"`
	Name AssistantMessageErrorAPIErrorName `json:"name"`
}

type AssistantMessageErrorAPIErrorData struct {
	IsRetryable     bool              `json:"isRetryable"`
	Message         string            `json:"message"`
	ResponseBody    *string           `json:"responseBody,omitempty"`
	ResponseHeaders map[string]string `json:"responseHeaders"`
	StatusCode      *int              `json:"statusCode,omitempty"`
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
	ID        string         `json:"id"`
	MessageID string         `json:"messageID"`
	Mime      string         `json:"mime"`
	SessionID string         `json:"sessionID"`
	Type      FilePartType   `json:"type"`
	URL       string         `json:"url"`
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

type FileDiff struct {
	File      string `json:"file"`
	Before    string `json:"before"`
	After     string `json:"after"`
	Additions int64  `json:"additions"`
	Deletions int64  `json:"deletions"`
}

type FilePartInputParam struct {
	Mime     string                    `json:"mime"`
	Type     FilePartInputType         `json:"type"`
	URL      string                    `json:"url"`
	ID       *string                   `json:"id,omitempty"`
	Filename *string                   `json:"filename,omitempty"`
	Source   *FilePartSourceUnionParam `json:"source,omitempty"`
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
	Type FilePartSourceType `json:"type"`
	// Embed raw JSON for lazy decode
	raw json.RawMessage `json:"-"`
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
func (r FilePartSource) AsFile() (*FileSource, error) {
	if r.Type != FilePartSourceTypeFile {
		return nil, ErrWrongVariant
	}
	var src FileSource
	if err := json.Unmarshal(r.raw, &src); err != nil {
		return nil, fmt.Errorf("unmarshal %s Type: %w", r.Type, err)
	}
	return &src, nil
}

// AsSymbol returns the source as a SymbolSource if Type is "symbol".
func (r FilePartSource) AsSymbol() (*SymbolSource, error) {
	if r.Type != FilePartSourceTypeSymbol {
		return nil, ErrWrongVariant
	}
	var src SymbolSource
	if err := json.Unmarshal(r.raw, &src); err != nil {
		return nil, fmt.Errorf("unmarshal %s Type: %w", r.Type, err)
	}
	return &src, nil
}

func (r FilePartSource) MarshalJSON() ([]byte, error) {
	if r.raw == nil {
		return []byte("null"), nil
	}
	return r.raw, nil
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
	Path  string                  `json:"path"`
	Text  FilePartSourceTextParam `json:"text"`
	Type  FilePartSourceType      `json:"type"`
	Kind  *int64                  `json:"kind,omitempty"`
	Name  *string                 `json:"name,omitempty"`
	Range any                     `json:"range,omitempty"`
}

func (r FilePartSourceParam) implementsFilePartSourceUnionParam() {}

// Satisfied by [FileSourceParam], [SymbolSourceParam], [FilePartSourceParam].
type FilePartSourceUnionParam interface {
	implementsFilePartSourceUnionParam()
}

type FilePartSourceText struct {
	End   int64  `json:"end"`
	Start int64  `json:"start"`
	Value string `json:"value"`
}

type FilePartSourceTextParam struct {
	End   int64  `json:"end"`
	Start int64  `json:"start"`
	Value string `json:"value"`
}

type FileSource struct {
	Path string             `json:"path"`
	Text FilePartSourceText `json:"text"`
	Type FileSourceType     `json:"type"`
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
	Path string                  `json:"path"`
	Text FilePartSourceTextParam `json:"text"`
	Type FileSourceType          `json:"type"`
}

func (r FileSourceParam) implementsFilePartSourceUnionParam() {}

// Message is either UserMessage or AssistantMessage, discriminated by Role.
type Message struct {
	ID        string      `json:"id"`
	Role      MessageRole `json:"role"`
	SessionID string      `json:"sessionID"`
	// Embed raw JSON for lazy decode
	raw json.RawMessage `json:"-"`
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
func (r Message) AsUser() (*UserMessage, error) {
	if r.Role != MessageRoleUser {
		return nil, ErrWrongVariant
	}
	var msg UserMessage
	if err := json.Unmarshal(r.raw, &msg); err != nil {
		return nil, fmt.Errorf("unmarshal %s Role: %w", r.Role, err)
	}
	return &msg, nil
}

// AsAssistant returns the AssistantMessage if the role is "assistant".
func (r Message) AsAssistant() (*AssistantMessage, error) {
	if r.Role != MessageRoleAssistant {
		return nil, ErrWrongVariant
	}
	var msg AssistantMessage
	if err := json.Unmarshal(r.raw, &msg); err != nil {
		return nil, fmt.Errorf("unmarshal %s Role: %w", r.Role, err)
	}
	return &msg, nil
}

func (r Message) MarshalJSON() ([]byte, error) {
	if r.raw == nil {
		return []byte("null"), nil
	}
	return r.raw, nil
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
	ID        string   `json:"id"`
	MessageID string   `json:"messageID"`
	SessionID string   `json:"sessionID"`
	Type      PartType `json:"type"`
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
func (r Part) AsText() (*TextPart, error) {
	if r.Type != PartTypeText {
		return nil, ErrWrongVariant
	}
	var part TextPart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, fmt.Errorf("unmarshal %s Type: %w", r.Type, err)
	}
	return &part, nil
}

// AsReasoning returns the part as a ReasoningPart if Type is "reasoning".
func (r Part) AsReasoning() (*ReasoningPart, error) {
	if r.Type != PartTypeReasoning {
		return nil, ErrWrongVariant
	}
	var part ReasoningPart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, fmt.Errorf("unmarshal %s Type: %w", r.Type, err)
	}
	return &part, nil
}

// AsFile returns the part as a FilePart if Type is "file".
func (r Part) AsFile() (*FilePart, error) {
	if r.Type != PartTypeFile {
		return nil, ErrWrongVariant
	}
	var part FilePart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, fmt.Errorf("unmarshal %s Type: %w", r.Type, err)
	}
	return &part, nil
}

// AsTool returns the part as a ToolPart if Type is "tool".
func (r Part) AsTool() (*ToolPart, error) {
	if r.Type != PartTypeTool {
		return nil, ErrWrongVariant
	}
	var part ToolPart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, fmt.Errorf("unmarshal %s Type: %w", r.Type, err)
	}
	return &part, nil
}

// AsStepStart returns the part as a StepStartPart if Type is "step-start".
func (r Part) AsStepStart() (*StepStartPart, error) {
	if r.Type != PartTypeStepStart {
		return nil, ErrWrongVariant
	}
	var part StepStartPart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, fmt.Errorf("unmarshal %s Type: %w", r.Type, err)
	}
	return &part, nil
}

// AsStepFinish returns the part as a StepFinishPart if Type is "step-finish".
func (r Part) AsStepFinish() (*StepFinishPart, error) {
	if r.Type != PartTypeStepFinish {
		return nil, ErrWrongVariant
	}
	var part StepFinishPart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, fmt.Errorf("unmarshal %s Type: %w", r.Type, err)
	}
	return &part, nil
}

// AsSnapshot returns the part as a SnapshotPart if Type is "snapshot".
func (r Part) AsSnapshot() (*SnapshotPart, error) {
	if r.Type != PartTypeSnapshot {
		return nil, ErrWrongVariant
	}
	var part SnapshotPart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, fmt.Errorf("unmarshal %s Type: %w", r.Type, err)
	}
	return &part, nil
}

// AsPatch returns the part as a PartPatchPart if Type is "patch".
func (r Part) AsPatch() (*PartPatchPart, error) {
	if r.Type != PartTypePatch {
		return nil, ErrWrongVariant
	}
	var part PartPatchPart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, fmt.Errorf("unmarshal %s Type: %w", r.Type, err)
	}
	return &part, nil
}

// AsAgent returns the part as an AgentPart if Type is "agent".
func (r Part) AsAgent() (*AgentPart, error) {
	if r.Type != PartTypeAgent {
		return nil, ErrWrongVariant
	}
	var part AgentPart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, fmt.Errorf("unmarshal %s Type: %w", r.Type, err)
	}
	return &part, nil
}

// AsRetry returns the part as a PartRetryPart if Type is "retry".
func (r Part) AsRetry() (*PartRetryPart, error) {
	if r.Type != PartTypeRetry {
		return nil, ErrWrongVariant
	}
	var part PartRetryPart
	if err := json.Unmarshal(r.raw, &part); err != nil {
		return nil, fmt.Errorf("unmarshal %s Type: %w", r.Type, err)
	}
	return &part, nil
}

func (r Part) MarshalJSON() ([]byte, error) {
	if r.raw == nil {
		return []byte("null"), nil
	}
	return r.raw, nil
}

type PartPatchPart struct {
	ID        string            `json:"id"`
	Files     []string          `json:"files"`
	Hash      string            `json:"hash"`
	MessageID string            `json:"messageID"`
	SessionID string            `json:"sessionID"`
	Type      PartPatchPartType `json:"type"`
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
	ID        string             `json:"id"`
	Attempt   int                `json:"attempt"`
	Error     PartRetryPartError `json:"error"`
	MessageID string             `json:"messageID"`
	SessionID string             `json:"sessionID"`
	Time      PartRetryPartTime  `json:"time"`
	Type      PartRetryPartType  `json:"type"`
}

type PartRetryPartError struct {
	Data PartRetryPartErrorData `json:"data"`
	Name PartRetryPartErrorName `json:"name"`
}

type PartRetryPartErrorData struct {
	IsRetryable     bool              `json:"isRetryable"`
	Message         string            `json:"message"`
	ResponseBody    *string           `json:"responseBody,omitempty"`
	ResponseHeaders map[string]string `json:"responseHeaders"`
	StatusCode      *int              `json:"statusCode,omitempty"`
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
	Created float64 `json:"created"`
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
	ID        string                 `json:"id"`
	MessageID string                 `json:"messageID"`
	SessionID string                 `json:"sessionID"`
	Text      string                 `json:"text"`
	Time      ReasoningPartTime      `json:"time"`
	Type      ReasoningPartType      `json:"type"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type ReasoningPartTime struct {
	Start float64  `json:"start"`
	End   *float64 `json:"end,omitempty"`
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
	ID        string          `json:"id"`
	Directory string          `json:"directory"`
	ProjectID string          `json:"projectID"`
	Time      SessionTime     `json:"time"`
	Title     string          `json:"title"`
	Version   string          `json:"version"`
	ParentID  *string         `json:"parentID,omitempty"`
	Revert    *SessionRevert  `json:"revert,omitempty"`
	Share     *SessionShare   `json:"share,omitempty"`
	Summary   *SessionSummary `json:"summary,omitempty"`
}

type SessionTime struct {
	Created    float64 `json:"created"`
	Updated    float64 `json:"updated"`
	Compacting float64 `json:"compacting,omitempty"`
}

type SessionRevert struct {
	MessageID string  `json:"messageID"`
	Diff      *string `json:"diff,omitempty"`
	PartID    *string `json:"partID,omitempty"`
	Snapshot  *string `json:"snapshot,omitempty"`
}

type SessionShare struct {
	URL string `json:"url"`
}

type SessionSummary struct {
	Diffs []SessionSummaryDiff `json:"diffs"`
}

type SessionSummaryDiff struct {
	Additions int64  `json:"additions"`
	After     string `json:"after"`
	Before    string `json:"before"`
	Deletions int64  `json:"deletions"`
	File      string `json:"file"`
}

type SnapshotPart struct {
	ID        string           `json:"id"`
	MessageID string           `json:"messageID"`
	SessionID string           `json:"sessionID"`
	Snapshot  string           `json:"snapshot"`
	Type      SnapshotPartType `json:"type"`
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
	ID        string               `json:"id"`
	Cost      float64              `json:"cost"`
	MessageID string               `json:"messageID"`
	Reason    string               `json:"reason"`
	SessionID string               `json:"sessionID"`
	Tokens    StepFinishPartTokens `json:"tokens"`
	Type      StepFinishPartType   `json:"type"`
	Snapshot  string               `json:"snapshot"`
}

type StepFinishPartTokens struct {
	Cache     StepFinishPartTokensCache `json:"cache"`
	Input     int64                     `json:"input"`
	Output    int64                     `json:"output"`
	Reasoning int64                     `json:"reasoning"`
}

type StepFinishPartTokensCache struct {
	Read  int64 `json:"read"`
	Write int64 `json:"write"`
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
	ID        string            `json:"id"`
	MessageID string            `json:"messageID"`
	SessionID string            `json:"sessionID"`
	Type      StepStartPartType `json:"type"`
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
	Kind  int64              `json:"kind"`
	Name  string             `json:"name"`
	Path  string             `json:"path"`
	Range SymbolSourceRange  `json:"range"`
	Text  FilePartSourceText `json:"text"`
	Type  SymbolSourceType   `json:"type"`
}

type SymbolSourceRange struct {
	End   SymbolSourceRangeEnd   `json:"end"`
	Start SymbolSourceRangeStart `json:"start"`
}

type SymbolSourceRangeEnd struct {
	Character int64 `json:"character"`
	Line      int64 `json:"line"`
}

type SymbolSourceRangeStart struct {
	Character int64 `json:"character"`
	Line      int64 `json:"line"`
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
	Kind  int64                   `json:"kind"`
	Name  string                  `json:"name"`
	Path  string                  `json:"path"`
	Range SymbolSourceRangeParam  `json:"range"`
	Text  FilePartSourceTextParam `json:"text"`
	Type  SymbolSourceType        `json:"type"`
}

func (r SymbolSourceParam) implementsFilePartSourceUnionParam() {}

type SymbolSourceRangeParam struct {
	End   SymbolSourceRangeEndParam   `json:"end"`
	Start SymbolSourceRangeStartParam `json:"start"`
}

type SymbolSourceRangeEndParam struct {
	Character int64 `json:"character"`
	Line      int64 `json:"line"`
}

type SymbolSourceRangeStartParam struct {
	Character int64 `json:"character"`
	Line      int64 `json:"line"`
}

type TextPart struct {
	ID        string                 `json:"id"`
	MessageID string                 `json:"messageID"`
	SessionID string                 `json:"sessionID"`
	Text      string                 `json:"text"`
	Type      TextPartType           `json:"type"`
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
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

type TextPartInputParam struct {
	Text      string                  `json:"text"`
	Type      TextPartInputType       `json:"type"`
	ID        *string                 `json:"id,omitempty"`
	Metadata  *map[string]interface{} `json:"metadata,omitempty"`
	Synthetic *bool                   `json:"synthetic,omitempty"`
	Time      *TextPartInputTimeParam `json:"time,omitempty"`
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
	Start float64  `json:"start"`
	End   *float64 `json:"end,omitempty"`
}

type ToolPart struct {
	ID        string                 `json:"id"`
	CallID    string                 `json:"callID"`
	MessageID string                 `json:"messageID"`
	SessionID string                 `json:"sessionID"`
	State     ToolPartState          `json:"state"`
	Tool      string                 `json:"tool"`
	Type      ToolPartType           `json:"type"`
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
func (r ToolPartState) AsPending() (*ToolStatePending, error) {
	if r.Status != ToolPartStateStatusPending {
		return nil, ErrWrongVariant
	}
	var state ToolStatePending
	if err := json.Unmarshal(r.raw, &state); err != nil {
		return nil, fmt.Errorf("unmarshal %s Status: %w", r.Status, err)
	}
	return &state, nil
}

// AsRunning returns the state as ToolStateRunning if Status is "running".
func (r ToolPartState) AsRunning() (*ToolStateRunning, error) {
	if r.Status != ToolPartStateStatusRunning {
		return nil, ErrWrongVariant
	}
	var state ToolStateRunning
	if err := json.Unmarshal(r.raw, &state); err != nil {
		return nil, fmt.Errorf("unmarshal %s Status: %w", r.Status, err)
	}
	return &state, nil
}

// AsCompleted returns the state as ToolStateCompleted if Status is "completed".
func (r ToolPartState) AsCompleted() (*ToolStateCompleted, error) {
	if r.Status != ToolPartStateStatusCompleted {
		return nil, ErrWrongVariant
	}
	var state ToolStateCompleted
	if err := json.Unmarshal(r.raw, &state); err != nil {
		return nil, fmt.Errorf("unmarshal %s Status: %w", r.Status, err)
	}
	return &state, nil
}

// AsError returns the state as ToolStateError if Status is "error".
func (r ToolPartState) AsError() (*ToolStateError, error) {
	if r.Status != ToolPartStateStatusError {
		return nil, ErrWrongVariant
	}
	var state ToolStateError
	if err := json.Unmarshal(r.raw, &state); err != nil {
		return nil, fmt.Errorf("unmarshal %s Status: %w", r.Status, err)
	}
	return &state, nil
}

func (r ToolPartState) MarshalJSON() ([]byte, error) {
	if r.raw == nil {
		return []byte("null"), nil
	}
	return r.raw, nil
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
	Input       map[string]interface{}   `json:"input"`
	Metadata    map[string]interface{}   `json:"metadata"`
	Output      string                   `json:"output"`
	Status      ToolStateCompletedStatus `json:"status"`
	Time        ToolStateCompletedTime   `json:"time"`
	Title       string                   `json:"title"`
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
	End       float64 `json:"end"`
	Start     float64 `json:"start"`
	Compacted float64 `json:"compacted"`
}

type ToolStateError struct {
	Error    string                 `json:"error"`
	Input    map[string]interface{} `json:"input"`
	Status   ToolStateErrorStatus   `json:"status"`
	Time     ToolStateErrorTime     `json:"time"`
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
	End   float64 `json:"end"`
	Start float64 `json:"start"`
}

type ToolStatePending struct {
	Status ToolStatePendingStatus `json:"status"`
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
	Input    interface{}            `json:"input"`
	Status   ToolStateRunningStatus `json:"status"`
	Time     ToolStateRunningTime   `json:"time"`
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
	Start float64 `json:"start"`
}

type UserMessage struct {
	ID        string             `json:"id"`
	Role      UserMessageRole    `json:"role"`
	SessionID string             `json:"sessionID"`
	Time      UserMessageTime    `json:"time"`
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
	Created float64 `json:"created"`
}

type UserMessageSummary struct {
	Diffs []FileDiff `json:"diffs"`
	Body  string     `json:"body"`
	Title string     `json:"title"`
}

type SessionCommandResponse struct {
	Info  AssistantMessage `json:"info"`
	Parts []Part           `json:"parts"`
}

type SessionMessageResponse struct {
	Info  Message `json:"info"`
	Parts []Part  `json:"parts"`
}

type SessionMessagesResponse struct {
	Info  Message `json:"info"`
	Parts []Part  `json:"parts"`
}

type SessionPromptResponse struct {
	Info  AssistantMessage `json:"info"`
	Parts []Part           `json:"parts"`
}

type SessionCreateParams struct {
	Directory *string `json:"-" query:"directory,omitempty"`
	ParentID  *string `json:"parentID,omitempty"`
	Title     *string `json:"title,omitempty"`
}

// URLQuery serializes [SessionCreateParams]'s query parameters as `url.Values`.
func (r SessionCreateParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionUpdateParams struct {
	Directory *string `json:"-" query:"directory,omitempty"`
	Title     *string `json:"title,omitempty"`
}

// URLQuery serializes [SessionUpdateParams]'s query parameters as `url.Values`.
func (r SessionUpdateParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionListParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionListParams]'s query parameters as `url.Values`.
func (r SessionListParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionDeleteParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionDeleteParams]'s query parameters as `url.Values`.
func (r SessionDeleteParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionAbortParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionAbortParams]'s query parameters as `url.Values`.
func (r SessionAbortParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionChildrenParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionChildrenParams]'s query parameters as `url.Values`.
func (r SessionChildrenParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionCommandParams struct {
	Arguments string  `json:"arguments"`
	Command   string  `json:"command"`
	Directory *string `json:"-" query:"directory,omitempty"`
	Agent     *string `json:"agent,omitempty"`
	MessageID *string `json:"messageID,omitempty"`
	Model     *string `json:"model,omitempty"`
}

// URLQuery serializes [SessionCommandParams]'s query parameters as `url.Values`.
func (r SessionCommandParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionDiffParams struct {
	Directory *string `query:"directory,omitempty"`
	MessageID *string `query:"messageID,omitempty"`
}

// URLQuery serializes [SessionDiffParams]'s query parameters as `url.Values`.
func (r SessionDiffParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionForkParams struct {
	MessageID *string `json:"messageID,omitempty"`
	Directory *string `json:"-" query:"directory,omitempty"`
}

// URLQuery serializes [SessionForkParams]'s query parameters as `url.Values`.
func (r SessionForkParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionGetParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionGetParams]'s query parameters as `url.Values`.
func (r SessionGetParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionInitParams struct {
	MessageID  string  `json:"messageID"`
	ModelID    string  `json:"modelID"`
	ProviderID string  `json:"providerID"`
	Directory  *string `json:"-" query:"directory,omitempty"`
}

// URLQuery serializes [SessionInitParams]'s query parameters as `url.Values`.
func (r SessionInitParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionMessageParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionMessageParams]'s query parameters as `url.Values`.
func (r SessionMessageParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionMessagesParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionMessagesParams]'s query parameters as `url.Values`.
func (r SessionMessagesParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionPromptParams struct {
	Parts     []SessionPromptParamsPartUnion `json:"parts"`
	Directory *string                        `json:"-" query:"directory,omitempty"`
	Agent     *string                        `json:"agent,omitempty"`
	MessageID *string                        `json:"messageID,omitempty"`
	Model     *SessionPromptParamsModel      `json:"model,omitempty"`
	NoReply   *bool                          `json:"noReply,omitempty"`
	System    *string                        `json:"system,omitempty"`
	Tools     *map[string]bool               `json:"tools,omitempty"`
}

// URLQuery serializes [SessionPromptParams]'s query parameters as `url.Values`.
func (r SessionPromptParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionPromptParamsPart struct {
	Type      SessionPromptParamsPartsType `json:"type"`
	ID        *string                      `json:"id,omitempty"`
	Filename  *string                      `json:"filename,omitempty"`
	Metadata  any                          `json:"metadata,omitempty"`
	Mime      *string                      `json:"mime,omitempty"`
	Name      *string                      `json:"name,omitempty"`
	Source    any                          `json:"source,omitempty"`
	Synthetic *bool                        `json:"synthetic,omitempty"`
	Text      *string                      `json:"text,omitempty"`
	Time      any                          `json:"time,omitempty"`
	URL       *string                      `json:"url,omitempty"`
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
	ModelID    string `json:"modelID"`
	ProviderID string `json:"providerID"`
}

type SessionRevertParams struct {
	MessageID string  `json:"messageID"`
	Directory *string `json:"-" query:"directory,omitempty"`
	PartID    *string `json:"partID,omitempty"`
}

// URLQuery serializes [SessionRevertParams]'s query parameters as `url.Values`.
func (r SessionRevertParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionShareParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionShareParams]'s query parameters as `url.Values`.
func (r SessionShareParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionShellParams struct {
	Agent     string  `json:"agent"`
	Command   string  `json:"command"`
	Directory *string `json:"-" query:"directory,omitempty"`
}

// URLQuery serializes [SessionShellParams]'s query parameters as `url.Values`.
func (r SessionShellParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionSummarizeParams struct {
	ModelID    string  `json:"modelID"`
	ProviderID string  `json:"providerID"`
	Directory  *string `json:"-" query:"directory,omitempty"`
}

// URLQuery serializes [SessionSummarizeParams]'s query parameters as `url.Values`.
func (r SessionSummarizeParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionTodoParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionTodoParams]'s query parameters as `url.Values`.
func (r SessionTodoParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionUnrevertParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionUnrevertParams]'s query parameters as `url.Values`.
func (r SessionUnrevertParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}

type SessionUnshareParams struct {
	Directory *string `query:"directory,omitempty"`
}

// URLQuery serializes [SessionUnshareParams]'s query parameters as `url.Values`.
func (r SessionUnshareParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}
