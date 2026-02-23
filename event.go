package opencode

import (
	"context"
	"net/http"
	"net/url"
	"reflect"

	"github.com/dominicnunez/opencode-sdk-go/internal/apijson"
	"github.com/dominicnunez/opencode-sdk-go/internal/apiquery"
	"github.com/dominicnunez/opencode-sdk-go/internal/requestconfig"
	"github.com/dominicnunez/opencode-sdk-go/option"
	"github.com/dominicnunez/opencode-sdk-go/packages/ssestream"
	"github.com/dominicnunez/opencode-sdk-go/shared"
	"github.com/tidwall/gjson"
)

type EventService struct {
	client *Client
}

func (s *EventService) ListStreaming(ctx context.Context, params *EventListParams, opts ...option.RequestOption) *ssestream.Stream[Event] {
	if params == nil {
		params = &EventListParams{}
	}
	var raw *http.Response
	allOpts := []option.RequestOption{
		option.WithHeader("Accept", "text/event-stream"),
		option.WithBaseURL(s.client.baseURL),
		option.WithHTTPClient(s.client.httpClient),
		option.WithMaxRetries(s.client.maxRetries),
	}
	allOpts = append(allOpts, s.client.defaultOptions...)
	allOpts = append(allOpts, opts...)
	err := requestconfig.ExecuteNewRequest(ctx, http.MethodGet, "event", params, &raw, allOpts...)
	return ssestream.NewStream[Event](ssestream.NewDecoder(raw), err)
}

type Event struct {
	Data  interface{} `json:"properties"`
	Type  EventType   `json:"type"`
	union EventUnion
}

func (r *Event) UnmarshalJSON(data []byte) error {
	*r = Event{}
	err := apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, r)
}

func (r Event) AsUnion() EventUnion {
	return r.union
}

type EventUnion interface {
	implementsEvent()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*EventUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(EventInstallationUpdated{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(EventLspClientDiagnostics{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(EventMessageUpdated{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(EventMessageRemoved{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(EventMessagePartUpdated{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(EventMessagePartRemoved{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(EventSessionCompacted{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(EventPermissionUpdated{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(EventPermissionReplied{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(EventFileEdited{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(EventFileWatcherUpdated{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(EventTodoUpdated{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(EventSessionIdle{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(EventSessionCreated{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(EventSessionUpdated{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(EventSessionDeleted{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(EventSessionError{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(EventServerConnected{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(EventIdeInstalled{}),
		},
	)
}

type EventInstallationUpdated struct {
	Data EventInstallationUpdatedData   `json:"properties"`
	Type EventInstallationUpdatedType   `json:"type"`
}

func (r EventInstallationUpdated) implementsEvent() {}

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

func (r EventLspClientDiagnostics) implementsEvent() {}

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

func (r EventMessageUpdated) implementsEvent() {}

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

func (r EventMessageRemoved) implementsEvent() {}

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

func (r EventMessagePartUpdated) implementsEvent() {}

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

func (r EventMessagePartRemoved) implementsEvent() {}

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

func (r EventSessionCompacted) implementsEvent() {}

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

func (r EventPermissionUpdated) implementsEvent() {}

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

func (r EventPermissionReplied) implementsEvent() {}

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

func (r EventFileEdited) implementsEvent() {}

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

func (r EventFileWatcherUpdated) implementsEvent() {}

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

func (r EventTodoUpdated) implementsEvent() {}

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

func (r EventSessionIdle) implementsEvent() {}

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

func (r EventSessionCreated) implementsEvent() {}

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

func (r EventSessionUpdated) implementsEvent() {}

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

func (r EventSessionDeleted) implementsEvent() {}

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

func (r EventSessionError) implementsEvent() {}

type EventSessionErrorData struct {
	Error     *SessionError `json:"error,omitempty"`
	SessionID *string       `json:"sessionID,omitempty"`
}

type SessionError struct {
	Data  interface{}        `json:"data"`
	Name  SessionErrorName   `json:"name"`
	union SessionErrorUnion
}

func (r *SessionError) UnmarshalJSON(data []byte) error {
	*r = SessionError{}
	err := apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, r)
}

func (r SessionError) AsUnion() SessionErrorUnion {
	return r.union
}

type SessionErrorUnion interface {
	ImplementsSessionError()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SessionErrorUnion)(nil)).Elem(),
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
			Type:       reflect.TypeOf(MessageOutputLengthError{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(shared.MessageAbortedError{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SessionAPIError{}),
		},
	)
}

type MessageOutputLengthError struct {
	Data interface{}                  `json:"data"`
	Name MessageOutputLengthErrorName `json:"name"`
}

func (r MessageOutputLengthError) ImplementsSessionError() {}

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

func (r SessionAPIError) ImplementsSessionError() {}

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

func (r EventServerConnected) implementsEvent() {}

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

func (r EventIdeInstalled) implementsEvent() {}

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
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
