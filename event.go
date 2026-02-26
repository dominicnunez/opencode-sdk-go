package opencode

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/dominicnunez/opencode-sdk-go/internal"
	"github.com/dominicnunez/opencode-sdk-go/internal/queryparams"
	"github.com/dominicnunez/opencode-sdk-go/packages/ssestream"
	"github.com/dominicnunez/opencode-sdk-go/shared"
)

type EventService struct {
	client *Client
}

func (s *EventService) ListStreaming(ctx context.Context, params *EventListParams) *ssestream.Stream[Event] {
	if params == nil {
		params = &EventListParams{}
	}

	// Build URL with query params
	u, err := url.Parse(s.client.baseURL)
	if err != nil {
		return ssestream.NewStream[Event](nil, err)
	}
	fullURL := u.ResolveReference(&url.URL{Path: "event"})

	if params != nil {
		query, err := params.URLQuery()
		if err != nil {
			return ssestream.NewStream[Event](nil, err)
		}
		fullURL.RawQuery = query.Encode()
	}

	// Create request with SSE headers
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return ssestream.NewStream[Event](nil, err)
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("User-Agent", fmt.Sprintf("Opencode/Go %s", internal.PackageVersion))

	// Execute request
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return ssestream.NewStream[Event](nil, err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return ssestream.NewStream[Event](nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body)))
	}

	return ssestream.NewStream[Event](ssestream.NewDecoder(resp), nil)
}

type Event struct {
	Type EventType `json:"type"`
	raw  json.RawMessage
}

func (e *Event) UnmarshalJSON(data []byte) error {
	// Peek at discriminator
	var peek struct {
		Type EventType `json:"type"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return err
	}
	e.Type = peek.Type
	e.raw = data
	return nil
}

// AsInstallationUpdated returns the event as EventInstallationUpdated if Type is "installation.updated".
func (e Event) AsInstallationUpdated() (*EventInstallationUpdated, bool) {
	if e.Type != EventTypeInstallationUpdated {
		return nil, false
	}
	var evt EventInstallationUpdated
	if err := json.Unmarshal(e.raw, &evt); err != nil {
		return nil, false
	}
	return &evt, true
}

// AsLspClientDiagnostics returns the event as EventLspClientDiagnostics if Type is "lsp.client.diagnostics".
func (e Event) AsLspClientDiagnostics() (*EventLspClientDiagnostics, bool) {
	if e.Type != EventTypeLspClientDiagnostics {
		return nil, false
	}
	var evt EventLspClientDiagnostics
	if err := json.Unmarshal(e.raw, &evt); err != nil {
		return nil, false
	}
	return &evt, true
}

// AsMessageUpdated returns the event as EventMessageUpdated if Type is "message.updated".
func (e Event) AsMessageUpdated() (*EventMessageUpdated, bool) {
	if e.Type != EventTypeMessageUpdated {
		return nil, false
	}
	var evt EventMessageUpdated
	if err := json.Unmarshal(e.raw, &evt); err != nil {
		return nil, false
	}
	return &evt, true
}

// AsMessageRemoved returns the event as EventMessageRemoved if Type is "message.removed".
func (e Event) AsMessageRemoved() (*EventMessageRemoved, bool) {
	if e.Type != EventTypeMessageRemoved {
		return nil, false
	}
	var evt EventMessageRemoved
	if err := json.Unmarshal(e.raw, &evt); err != nil {
		return nil, false
	}
	return &evt, true
}

// AsMessagePartUpdated returns the event as EventMessagePartUpdated if Type is "message.part.updated".
func (e Event) AsMessagePartUpdated() (*EventMessagePartUpdated, bool) {
	if e.Type != EventTypeMessagePartUpdated {
		return nil, false
	}
	var evt EventMessagePartUpdated
	if err := json.Unmarshal(e.raw, &evt); err != nil {
		return nil, false
	}
	return &evt, true
}

// AsMessagePartRemoved returns the event as EventMessagePartRemoved if Type is "message.part.removed".
func (e Event) AsMessagePartRemoved() (*EventMessagePartRemoved, bool) {
	if e.Type != EventTypeMessagePartRemoved {
		return nil, false
	}
	var evt EventMessagePartRemoved
	if err := json.Unmarshal(e.raw, &evt); err != nil {
		return nil, false
	}
	return &evt, true
}

// AsSessionCompacted returns the event as EventSessionCompacted if Type is "session.compacted".
func (e Event) AsSessionCompacted() (*EventSessionCompacted, bool) {
	if e.Type != EventTypeSessionCompacted {
		return nil, false
	}
	var evt EventSessionCompacted
	if err := json.Unmarshal(e.raw, &evt); err != nil {
		return nil, false
	}
	return &evt, true
}

// AsPermissionUpdated returns the event as EventPermissionUpdated if Type is "permission.updated".
func (e Event) AsPermissionUpdated() (*EventPermissionUpdated, bool) {
	if e.Type != EventTypePermissionUpdated {
		return nil, false
	}
	var evt EventPermissionUpdated
	if err := json.Unmarshal(e.raw, &evt); err != nil {
		return nil, false
	}
	return &evt, true
}

// AsPermissionReplied returns the event as EventPermissionReplied if Type is "permission.replied".
func (e Event) AsPermissionReplied() (*EventPermissionReplied, bool) {
	if e.Type != EventTypePermissionReplied {
		return nil, false
	}
	var evt EventPermissionReplied
	if err := json.Unmarshal(e.raw, &evt); err != nil {
		return nil, false
	}
	return &evt, true
}

// AsFileEdited returns the event as EventFileEdited if Type is "file.edited".
func (e Event) AsFileEdited() (*EventFileEdited, bool) {
	if e.Type != EventTypeFileEdited {
		return nil, false
	}
	var evt EventFileEdited
	if err := json.Unmarshal(e.raw, &evt); err != nil {
		return nil, false
	}
	return &evt, true
}

// AsFileWatcherUpdated returns the event as EventFileWatcherUpdated if Type is "file.watcher.updated".
func (e Event) AsFileWatcherUpdated() (*EventFileWatcherUpdated, bool) {
	if e.Type != EventTypeFileWatcherUpdated {
		return nil, false
	}
	var evt EventFileWatcherUpdated
	if err := json.Unmarshal(e.raw, &evt); err != nil {
		return nil, false
	}
	return &evt, true
}

// AsTodoUpdated returns the event as EventTodoUpdated if Type is "todo.updated".
func (e Event) AsTodoUpdated() (*EventTodoUpdated, bool) {
	if e.Type != EventTypeTodoUpdated {
		return nil, false
	}
	var evt EventTodoUpdated
	if err := json.Unmarshal(e.raw, &evt); err != nil {
		return nil, false
	}
	return &evt, true
}

// AsSessionIdle returns the event as EventSessionIdle if Type is "session.idle".
func (e Event) AsSessionIdle() (*EventSessionIdle, bool) {
	if e.Type != EventTypeSessionIdle {
		return nil, false
	}
	var evt EventSessionIdle
	if err := json.Unmarshal(e.raw, &evt); err != nil {
		return nil, false
	}
	return &evt, true
}

// AsSessionCreated returns the event as EventSessionCreated if Type is "session.created".
func (e Event) AsSessionCreated() (*EventSessionCreated, bool) {
	if e.Type != EventTypeSessionCreated {
		return nil, false
	}
	var evt EventSessionCreated
	if err := json.Unmarshal(e.raw, &evt); err != nil {
		return nil, false
	}
	return &evt, true
}

// AsSessionUpdated returns the event as EventSessionUpdated if Type is "session.updated".
func (e Event) AsSessionUpdated() (*EventSessionUpdated, bool) {
	if e.Type != EventTypeSessionUpdated {
		return nil, false
	}
	var evt EventSessionUpdated
	if err := json.Unmarshal(e.raw, &evt); err != nil {
		return nil, false
	}
	return &evt, true
}

// AsSessionDeleted returns the event as EventSessionDeleted if Type is "session.deleted".
func (e Event) AsSessionDeleted() (*EventSessionDeleted, bool) {
	if e.Type != EventTypeSessionDeleted {
		return nil, false
	}
	var evt EventSessionDeleted
	if err := json.Unmarshal(e.raw, &evt); err != nil {
		return nil, false
	}
	return &evt, true
}

// AsSessionError returns the event as EventSessionError if Type is "session.error".
func (e Event) AsSessionError() (*EventSessionError, bool) {
	if e.Type != EventTypeSessionError {
		return nil, false
	}
	var evt EventSessionError
	if err := json.Unmarshal(e.raw, &evt); err != nil {
		return nil, false
	}
	return &evt, true
}

// AsServerConnected returns the event as EventServerConnected if Type is "server.connected".
func (e Event) AsServerConnected() (*EventServerConnected, bool) {
	if e.Type != EventTypeServerConnected {
		return nil, false
	}
	var evt EventServerConnected
	if err := json.Unmarshal(e.raw, &evt); err != nil {
		return nil, false
	}
	return &evt, true
}

// AsIdeInstalled returns the event as EventIdeInstalled if Type is "ide.installed".
func (e Event) AsIdeInstalled() (*EventIdeInstalled, bool) {
	if e.Type != EventTypeIdeInstalled {
		return nil, false
	}
	var evt EventIdeInstalled
	if err := json.Unmarshal(e.raw, &evt); err != nil {
		return nil, false
	}
	return &evt, true
}

type EventInstallationUpdated struct {
	Data EventInstallationUpdatedData   `json:"properties"`
	Type EventInstallationUpdatedType   `json:"type"`
}


type EventInstallationUpdatedData struct {
	Version string `json:"version"`
}

type EventInstallationUpdatedType string

const (
	EventInstallationUpdatedTypeInstallationUpdated EventInstallationUpdatedType = "installation.updated"
)

func (r EventInstallationUpdatedType) IsKnown() bool {
	switch r {
	case EventInstallationUpdatedTypeInstallationUpdated:
		return true
	}
	return false
}

type EventLspClientDiagnostics struct {
	Data EventLspClientDiagnosticsData  `json:"properties"`
	Type EventLspClientDiagnosticsType  `json:"type"`
}


type EventLspClientDiagnosticsData struct {
	Path     string `json:"path"`
	ServerID string `json:"serverID"`
}

type EventLspClientDiagnosticsType string

const (
	EventLspClientDiagnosticsTypeLspClientDiagnostics EventLspClientDiagnosticsType = "lsp.client.diagnostics"
)

func (r EventLspClientDiagnosticsType) IsKnown() bool {
	switch r {
	case EventLspClientDiagnosticsTypeLspClientDiagnostics:
		return true
	}
	return false
}

type EventMessageUpdated struct {
	Data EventMessageUpdatedData `json:"properties"`
	Type EventMessageUpdatedType `json:"type"`
}


type EventMessageUpdatedData struct {
	Info Message `json:"info"`
}

type EventMessageUpdatedType string

const (
	EventMessageUpdatedTypeMessageUpdated EventMessageUpdatedType = "message.updated"
)

func (r EventMessageUpdatedType) IsKnown() bool {
	switch r {
	case EventMessageUpdatedTypeMessageUpdated:
		return true
	}
	return false
}

type EventMessageRemoved struct {
	Data EventMessageRemovedData `json:"properties"`
	Type EventMessageRemovedType `json:"type"`
}


type EventMessageRemovedData struct {
	MessageID string `json:"messageID"`
	SessionID string `json:"sessionID"`
}

type EventMessageRemovedType string

const (
	EventMessageRemovedTypeMessageRemoved EventMessageRemovedType = "message.removed"
)

func (r EventMessageRemovedType) IsKnown() bool {
	switch r {
	case EventMessageRemovedTypeMessageRemoved:
		return true
	}
	return false
}

type EventMessagePartUpdated struct {
	Data EventMessagePartUpdatedData `json:"properties"`
	Type EventMessagePartUpdatedType `json:"type"`
}


type EventMessagePartUpdatedData struct {
	Part  Part    `json:"part"`
	Delta *string `json:"delta,omitempty"`
}

type EventMessagePartUpdatedType string

const (
	EventMessagePartUpdatedTypeMessagePartUpdated EventMessagePartUpdatedType = "message.part.updated"
)

func (r EventMessagePartUpdatedType) IsKnown() bool {
	switch r {
	case EventMessagePartUpdatedTypeMessagePartUpdated:
		return true
	}
	return false
}

type EventMessagePartRemoved struct {
	Data EventMessagePartRemovedData `json:"properties"`
	Type EventMessagePartRemovedType `json:"type"`
}


type EventMessagePartRemovedData struct {
	MessageID string `json:"messageID"`
	PartID    string `json:"partID"`
	SessionID string `json:"sessionID"`
}

type EventMessagePartRemovedType string

const (
	EventMessagePartRemovedTypeMessagePartRemoved EventMessagePartRemovedType = "message.part.removed"
)

func (r EventMessagePartRemovedType) IsKnown() bool {
	switch r {
	case EventMessagePartRemovedTypeMessagePartRemoved:
		return true
	}
	return false
}

type EventSessionCompacted struct {
	Data EventSessionCompactedData `json:"properties"`
	Type EventSessionCompactedType `json:"type"`
}


type EventSessionCompactedData struct {
	SessionID string `json:"sessionID"`
}

type EventSessionCompactedType string

const (
	EventSessionCompactedTypeSessionCompacted EventSessionCompactedType = "session.compacted"
)

func (r EventSessionCompactedType) IsKnown() bool {
	switch r {
	case EventSessionCompactedTypeSessionCompacted:
		return true
	}
	return false
}

type EventPermissionUpdated struct {
	Data Permission                `json:"properties"`
	Type EventPermissionUpdatedType `json:"type"`
}


type EventPermissionUpdatedType string

const (
	EventPermissionUpdatedTypePermissionUpdated EventPermissionUpdatedType = "permission.updated"
)

func (r EventPermissionUpdatedType) IsKnown() bool {
	switch r {
	case EventPermissionUpdatedTypePermissionUpdated:
		return true
	}
	return false
}

type EventPermissionReplied struct {
	Data EventPermissionRepliedData `json:"properties"`
	Type EventPermissionRepliedType `json:"type"`
}


type EventPermissionRepliedData struct {
	PermissionID string `json:"permissionID"`
	Response     string `json:"response"`
	SessionID    string `json:"sessionID"`
}

type EventPermissionRepliedType string

const (
	EventPermissionRepliedTypePermissionReplied EventPermissionRepliedType = "permission.replied"
)

func (r EventPermissionRepliedType) IsKnown() bool {
	switch r {
	case EventPermissionRepliedTypePermissionReplied:
		return true
	}
	return false
}

type EventFileEdited struct {
	Data EventFileEditedData `json:"properties"`
	Type EventFileEditedType `json:"type"`
}


type EventFileEditedData struct {
	File string `json:"file"`
}

type EventFileEditedType string

const (
	EventFileEditedTypeFileEdited EventFileEditedType = "file.edited"
)

func (r EventFileEditedType) IsKnown() bool {
	switch r {
	case EventFileEditedTypeFileEdited:
		return true
	}
	return false
}

type EventFileWatcherUpdated struct {
	Data EventFileWatcherUpdatedData `json:"properties"`
	Type EventFileWatcherUpdatedType `json:"type"`
}


type EventFileWatcherUpdatedData struct {
	Event EventFileWatcherUpdatedDataEvent `json:"event"`
	File  string                           `json:"file"`
}

type EventFileWatcherUpdatedDataEvent string

const (
	EventFileWatcherUpdatedDataEventAdd    EventFileWatcherUpdatedDataEvent = "add"
	EventFileWatcherUpdatedDataEventChange EventFileWatcherUpdatedDataEvent = "change"
	EventFileWatcherUpdatedDataEventUnlink EventFileWatcherUpdatedDataEvent = "unlink"
)

func (r EventFileWatcherUpdatedDataEvent) IsKnown() bool {
	switch r {
	case EventFileWatcherUpdatedDataEventAdd, EventFileWatcherUpdatedDataEventChange, EventFileWatcherUpdatedDataEventUnlink:
		return true
	}
	return false
}

type EventFileWatcherUpdatedType string

const (
	EventFileWatcherUpdatedTypeFileWatcherUpdated EventFileWatcherUpdatedType = "file.watcher.updated"
)

func (r EventFileWatcherUpdatedType) IsKnown() bool {
	switch r {
	case EventFileWatcherUpdatedTypeFileWatcherUpdated:
		return true
	}
	return false
}

type EventTodoUpdated struct {
	Data EventTodoUpdatedData `json:"properties"`
	Type EventTodoUpdatedType `json:"type"`
}


type EventTodoUpdatedData struct {
	SessionID string `json:"sessionID"`
	Todos     []Todo `json:"todos"`
}

type Todo struct {
	ID       string `json:"id"`
	Content  string `json:"content"`
	Priority string `json:"priority"`
	Status   string `json:"status"`
}

type EventTodoUpdatedType string

const (
	EventTodoUpdatedTypeTodoUpdated EventTodoUpdatedType = "todo.updated"
)

func (r EventTodoUpdatedType) IsKnown() bool {
	switch r {
	case EventTodoUpdatedTypeTodoUpdated:
		return true
	}
	return false
}

type EventSessionIdle struct {
	Data EventSessionIdleData `json:"properties"`
	Type EventSessionIdleType `json:"type"`
}


type EventSessionIdleData struct {
	SessionID string `json:"sessionID"`
}

type EventSessionIdleType string

const (
	EventSessionIdleTypeSessionIdle EventSessionIdleType = "session.idle"
)

func (r EventSessionIdleType) IsKnown() bool {
	switch r {
	case EventSessionIdleTypeSessionIdle:
		return true
	}
	return false
}

type EventSessionCreated struct {
	Data EventSessionCreatedData `json:"properties"`
	Type EventSessionCreatedType `json:"type"`
}


type EventSessionCreatedData struct {
	Info Session `json:"info"`
}

type EventSessionCreatedType string

const (
	EventSessionCreatedTypeSessionCreated EventSessionCreatedType = "session.created"
)

func (r EventSessionCreatedType) IsKnown() bool {
	switch r {
	case EventSessionCreatedTypeSessionCreated:
		return true
	}
	return false
}

type EventSessionUpdated struct {
	Data EventSessionUpdatedData `json:"properties"`
	Type EventSessionUpdatedType `json:"type"`
}


type EventSessionUpdatedData struct {
	Info Session `json:"info"`
}

type EventSessionUpdatedType string

const (
	EventSessionUpdatedTypeSessionUpdated EventSessionUpdatedType = "session.updated"
)

func (r EventSessionUpdatedType) IsKnown() bool {
	switch r {
	case EventSessionUpdatedTypeSessionUpdated:
		return true
	}
	return false
}

type EventSessionDeleted struct {
	Data EventSessionDeletedData `json:"properties"`
	Type EventSessionDeletedType `json:"type"`
}


type EventSessionDeletedData struct {
	Info Session `json:"info"`
}

type EventSessionDeletedType string

const (
	EventSessionDeletedTypeSessionDeleted EventSessionDeletedType = "session.deleted"
)

func (r EventSessionDeletedType) IsKnown() bool {
	switch r {
	case EventSessionDeletedTypeSessionDeleted:
		return true
	}
	return false
}

type EventSessionError struct {
	Data EventSessionErrorData `json:"properties"`
	Type EventSessionErrorType `json:"type"`
}


type EventSessionErrorData struct {
	Error     *SessionError `json:"error,omitempty"`
	SessionID *string       `json:"sessionID,omitempty"`
}

type SessionError struct {
	Name SessionErrorName `json:"name"`
	raw  json.RawMessage
}

func (r *SessionError) UnmarshalJSON(data []byte) error {
	// Peek at discriminator
	var peek struct {
		Name SessionErrorName `json:"name"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return err
	}
	r.Name = peek.Name
	r.raw = data
	return nil
}

func (r SessionError) AsProviderAuth() (*shared.ProviderAuthError, bool) {
	if r.Name != SessionErrorNameProviderAuthError {
		return nil, false
	}
	var err shared.ProviderAuthError
	if e := json.Unmarshal(r.raw, &err); e != nil {
		return nil, false
	}
	return &err, true
}

func (r SessionError) AsUnknown() (*shared.UnknownError, bool) {
	if r.Name != SessionErrorNameUnknownError {
		return nil, false
	}
	var err shared.UnknownError
	if e := json.Unmarshal(r.raw, &err); e != nil {
		return nil, false
	}
	return &err, true
}

func (r SessionError) AsOutputLength() (*MessageOutputLengthError, bool) {
	if r.Name != SessionErrorNameMessageOutputLengthError {
		return nil, false
	}
	var err MessageOutputLengthError
	if e := json.Unmarshal(r.raw, &err); e != nil {
		return nil, false
	}
	return &err, true
}

func (r SessionError) AsAborted() (*shared.MessageAbortedError, bool) {
	if r.Name != SessionErrorNameMessageAbortedError {
		return nil, false
	}
	var err shared.MessageAbortedError
	if e := json.Unmarshal(r.raw, &err); e != nil {
		return nil, false
	}
	return &err, true
}

func (r SessionError) AsAPI() (*SessionAPIError, bool) {
	if r.Name != SessionErrorNameAPIError {
		return nil, false
	}
	var err SessionAPIError
	if e := json.Unmarshal(r.raw, &err); e != nil {
		return nil, false
	}
	return &err, true
}

type MessageOutputLengthError struct {
	Data interface{}                  `json:"data"`
	Name MessageOutputLengthErrorName `json:"name"`
}

type MessageOutputLengthErrorName string

const (
	MessageOutputLengthErrorNameMessageOutputLengthError MessageOutputLengthErrorName = "MessageOutputLengthError"
)

func (r MessageOutputLengthErrorName) IsKnown() bool {
	switch r {
	case MessageOutputLengthErrorNameMessageOutputLengthError:
		return true
	}
	return false
}

type SessionAPIError struct {
	Data SessionAPIErrorData `json:"data"`
	Name SessionAPIErrorName `json:"name"`
}

type SessionAPIErrorData struct {
	IsRetryable     bool              `json:"isRetryable"`
	Message         string            `json:"message"`
	ResponseBody    *string           `json:"responseBody,omitempty"`
	ResponseHeaders map[string]string `json:"responseHeaders,omitempty"`
	StatusCode      *float64          `json:"statusCode,omitempty"`
}

type SessionAPIErrorName string

const (
	SessionAPIErrorNameAPIError SessionAPIErrorName = "APIError"
)

func (r SessionAPIErrorName) IsKnown() bool {
	switch r {
	case SessionAPIErrorNameAPIError:
		return true
	}
	return false
}

type SessionErrorName string

const (
	SessionErrorNameProviderAuthError        SessionErrorName = "ProviderAuthError"
	SessionErrorNameUnknownError             SessionErrorName = "UnknownError"
	SessionErrorNameMessageOutputLengthError SessionErrorName = "MessageOutputLengthError"
	SessionErrorNameMessageAbortedError      SessionErrorName = "MessageAbortedError"
	SessionErrorNameAPIError                 SessionErrorName = "APIError"
)

func (r SessionErrorName) IsKnown() bool {
	switch r {
	case SessionErrorNameProviderAuthError, SessionErrorNameUnknownError, SessionErrorNameMessageOutputLengthError, SessionErrorNameMessageAbortedError, SessionErrorNameAPIError:
		return true
	}
	return false
}

type EventSessionErrorType string

const (
	EventSessionErrorTypeSessionError EventSessionErrorType = "session.error"
)

func (r EventSessionErrorType) IsKnown() bool {
	switch r {
	case EventSessionErrorTypeSessionError:
		return true
	}
	return false
}

type EventServerConnected struct {
	Data interface{}              `json:"properties"`
	Type EventServerConnectedType `json:"type"`
}


type EventServerConnectedType string

const (
	EventServerConnectedTypeServerConnected EventServerConnectedType = "server.connected"
)

func (r EventServerConnectedType) IsKnown() bool {
	switch r {
	case EventServerConnectedTypeServerConnected:
		return true
	}
	return false
}

type EventIdeInstalled struct {
	Data EventIdeInstalledData `json:"properties"`
	Type EventIdeInstalledType `json:"type"`
}


type EventIdeInstalledData struct {
	Ide string `json:"ide"`
}

type EventIdeInstalledType string

const (
	EventIdeInstalledTypeIdeInstalled EventIdeInstalledType = "ide.installed"
)

func (r EventIdeInstalledType) IsKnown() bool {
	switch r {
	case EventIdeInstalledTypeIdeInstalled:
		return true
	}
	return false
}

type EventType string

const (
	EventTypeInstallationUpdated  EventType = "installation.updated"
	EventTypeLspClientDiagnostics EventType = "lsp.client.diagnostics"
	EventTypeMessageUpdated       EventType = "message.updated"
	EventTypeMessageRemoved       EventType = "message.removed"
	EventTypeMessagePartUpdated   EventType = "message.part.updated"
	EventTypeMessagePartRemoved   EventType = "message.part.removed"
	EventTypeSessionCompacted     EventType = "session.compacted"
	EventTypePermissionUpdated    EventType = "permission.updated"
	EventTypePermissionReplied    EventType = "permission.replied"
	EventTypeFileEdited           EventType = "file.edited"
	EventTypeFileWatcherUpdated   EventType = "file.watcher.updated"
	EventTypeTodoUpdated          EventType = "todo.updated"
	EventTypeSessionIdle          EventType = "session.idle"
	EventTypeSessionCreated       EventType = "session.created"
	EventTypeSessionUpdated       EventType = "session.updated"
	EventTypeSessionDeleted       EventType = "session.deleted"
	EventTypeSessionError         EventType = "session.error"
	EventTypeServerConnected      EventType = "server.connected"
	EventTypeIdeInstalled         EventType = "ide.installed"
)

func (r EventType) IsKnown() bool {
	switch r {
	case EventTypeInstallationUpdated, EventTypeLspClientDiagnostics, EventTypeMessageUpdated, EventTypeMessageRemoved, EventTypeMessagePartUpdated, EventTypeMessagePartRemoved, EventTypeSessionCompacted, EventTypePermissionUpdated, EventTypePermissionReplied, EventTypeFileEdited, EventTypeFileWatcherUpdated, EventTypeTodoUpdated, EventTypeSessionIdle, EventTypeSessionCreated, EventTypeSessionUpdated, EventTypeSessionDeleted, EventTypeSessionError, EventTypeServerConnected, EventTypeIdeInstalled:
		return true
	}
	return false
}

type EventListParams struct {
	Directory *string `query:"directory,omitempty"`
}

func (r EventListParams) URLQuery() (url.Values, error) {
	return queryparams.Marshal(r)
}
