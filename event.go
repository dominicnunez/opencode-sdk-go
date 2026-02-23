// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package opencode

import (
	"context"
	"net/http"
	"net/url"
	"reflect"

	"github.com/dominicnunez/opencode-sdk-go/internal/apijson"
	"github.com/dominicnunez/opencode-sdk-go/internal/apiquery"
	"github.com/dominicnunez/opencode-sdk-go/internal/param"
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
	// This field can have the runtime type of
	// [EventInstallationUpdatedData],
	// [EventLspClientDiagnosticsData],
	// [EventMessageUpdatedData],
	// [EventMessageRemovedData],
	// [EventMessagePartUpdatedData],
	// [EventMessagePartRemovedData],
	// [EventSessionCompactedData], [Permission],
	// [EventPermissionRepliedData],
	// [EventFileEditedData],
	// [EventFileWatcherUpdatedData],
	// [EventTodoUpdatedData],
	// [EventSessionIdleData],
	// [EventSessionCreatedData],
	// [EventSessionUpdatedData],
	// [EventSessionDeletedData],
	// [EventSessionErrorData], [interface{}],
	// [EventIdeInstalledData].
	Data interface{}           `json:"properties,required"`
	Type       EventType `json:"type,required"`
	JSON       eventJSON `json:"-"`
	union      EventUnion
}

// eventJSON contains the JSON metadata for the struct
// [Event]
type eventJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r eventJSON) RawJSON() string {
	return r.raw
}

func (r *Event) UnmarshalJSON(data []byte) (err error) {
	*r = Event{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [EventUnion] interface which you can cast to the
// specific types for more type safety.
//
// Possible runtime types of the union are
// [EventInstallationUpdated],
// [EventLspClientDiagnostics],
// [EventMessageUpdated], [EventMessageRemoved],
// [EventMessagePartUpdated],
// [EventMessagePartRemoved],
// [EventSessionCompacted],
// [EventPermissionUpdated],
// [EventPermissionReplied], [EventFileEdited],
// [EventFileWatcherUpdated], [EventTodoUpdated],
// [EventSessionIdle], [EventSessionCreated],
// [EventSessionUpdated], [EventSessionDeleted],
// [EventSessionError], [EventServerConnected],
// [EventIdeInstalled].
func (r Event) AsUnion() EventUnion {
	return r.union
}

// Union satisfied by [EventInstallationUpdated],
// [EventLspClientDiagnostics],
// [EventMessageUpdated], [EventMessageRemoved],
// [EventMessagePartUpdated],
// [EventMessagePartRemoved],
// [EventSessionCompacted],
// [EventPermissionUpdated],
// [EventPermissionReplied], [EventFileEdited],
// [EventFileWatcherUpdated], [EventTodoUpdated],
// [EventSessionIdle], [EventSessionCreated],
// [EventSessionUpdated], [EventSessionDeleted],
// [EventSessionError], [EventServerConnected] or
// [EventIdeInstalled].
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
	Data EventInstallationUpdatedData `json:"properties,required"`
	Type       EventInstallationUpdatedType       `json:"type,required"`
	JSON       eventInstallationUpdatedJSON       `json:"-"`
}

// eventInstallationUpdatedJSON contains the JSON metadata for the
// struct [EventInstallationUpdated]
type eventInstallationUpdatedJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventInstallationUpdated) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventInstallationUpdatedJSON) RawJSON() string {
	return r.raw
}

func (r EventInstallationUpdated) implementsEvent() {}

type EventInstallationUpdatedData struct {
	Version string                                                  `json:"version,required"`
	JSON    eventInstallationUpdatedDataJSON `json:"-"`
}

// eventInstallationUpdatedDataJSON contains the JSON
// metadata for the struct [EventInstallationUpdatedData]
type eventInstallationUpdatedDataJSON struct {
	Version     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventInstallationUpdatedData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventInstallationUpdatedDataJSON) RawJSON() string {
	return r.raw
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
	Data EventLspClientDiagnosticsData `json:"properties,required"`
	Type       EventLspClientDiagnosticsType       `json:"type,required"`
	JSON       eventLspClientDiagnosticsJSON       `json:"-"`
}

// eventLspClientDiagnosticsJSON contains the JSON metadata for
// the struct [EventLspClientDiagnostics]
type eventLspClientDiagnosticsJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventLspClientDiagnostics) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventLspClientDiagnosticsJSON) RawJSON() string {
	return r.raw
}

func (r EventLspClientDiagnostics) implementsEvent() {}

type EventLspClientDiagnosticsData struct {
	Path     string                                                   `json:"path,required"`
	ServerID string                                                   `json:"serverID,required"`
	JSON     eventLspClientDiagnosticsDataJSON `json:"-"`
}

// eventLspClientDiagnosticsDataJSON contains the JSON
// metadata for the struct [EventLspClientDiagnosticsData]
type eventLspClientDiagnosticsDataJSON struct {
	Path        apijson.Field
	ServerID    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventLspClientDiagnosticsData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventLspClientDiagnosticsDataJSON) RawJSON() string {
	return r.raw
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
	Data EventMessageUpdatedData `json:"properties,required"`
	Type       EventMessageUpdatedType       `json:"type,required"`
	JSON       eventMessageUpdatedJSON       `json:"-"`
}

// eventMessageUpdatedJSON contains the JSON metadata for the
// struct [EventMessageUpdated]
type eventMessageUpdatedJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventMessageUpdated) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventMessageUpdatedJSON) RawJSON() string {
	return r.raw
}

func (r EventMessageUpdated) implementsEvent() {}

type EventMessageUpdatedData struct {
	Info Message                                            `json:"info,required"`
	JSON eventMessageUpdatedDataJSON `json:"-"`
}

// eventMessageUpdatedDataJSON contains the JSON metadata
// for the struct [EventMessageUpdatedData]
type eventMessageUpdatedDataJSON struct {
	Info        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventMessageUpdatedData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventMessageUpdatedDataJSON) RawJSON() string {
	return r.raw
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
	Data EventMessageRemovedData `json:"properties,required"`
	Type       EventMessageRemovedType       `json:"type,required"`
	JSON       eventMessageRemovedJSON       `json:"-"`
}

// eventMessageRemovedJSON contains the JSON metadata for the
// struct [EventMessageRemoved]
type eventMessageRemovedJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventMessageRemoved) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventMessageRemovedJSON) RawJSON() string {
	return r.raw
}

func (r EventMessageRemoved) implementsEvent() {}

type EventMessageRemovedData struct {
	MessageID string                                             `json:"messageID,required"`
	SessionID string                                             `json:"sessionID,required"`
	JSON      eventMessageRemovedDataJSON `json:"-"`
}

// eventMessageRemovedDataJSON contains the JSON metadata
// for the struct [EventMessageRemovedData]
type eventMessageRemovedDataJSON struct {
	MessageID   apijson.Field
	SessionID   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventMessageRemovedData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventMessageRemovedDataJSON) RawJSON() string {
	return r.raw
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
	Data EventMessagePartUpdatedData `json:"properties,required"`
	Type       EventMessagePartUpdatedType       `json:"type,required"`
	JSON       eventMessagePartUpdatedJSON       `json:"-"`
}

// eventMessagePartUpdatedJSON contains the JSON metadata for the
// struct [EventMessagePartUpdated]
type eventMessagePartUpdatedJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventMessagePartUpdated) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventMessagePartUpdatedJSON) RawJSON() string {
	return r.raw
}

func (r EventMessagePartUpdated) implementsEvent() {}

type EventMessagePartUpdatedData struct {
	Part  Part                                                   `json:"part,required"`
	Delta string                                                 `json:"delta"`
	JSON  eventMessagePartUpdatedDataJSON `json:"-"`
}

// eventMessagePartUpdatedDataJSON contains the JSON
// metadata for the struct [EventMessagePartUpdatedData]
type eventMessagePartUpdatedDataJSON struct {
	Part        apijson.Field
	Delta       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventMessagePartUpdatedData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventMessagePartUpdatedDataJSON) RawJSON() string {
	return r.raw
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
	Data EventMessagePartRemovedData `json:"properties,required"`
	Type       EventMessagePartRemovedType       `json:"type,required"`
	JSON       eventMessagePartRemovedJSON       `json:"-"`
}

// eventMessagePartRemovedJSON contains the JSON metadata for the
// struct [EventMessagePartRemoved]
type eventMessagePartRemovedJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventMessagePartRemoved) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventMessagePartRemovedJSON) RawJSON() string {
	return r.raw
}

func (r EventMessagePartRemoved) implementsEvent() {}

type EventMessagePartRemovedData struct {
	MessageID string                                                 `json:"messageID,required"`
	PartID    string                                                 `json:"partID,required"`
	SessionID string                                                 `json:"sessionID,required"`
	JSON      eventMessagePartRemovedDataJSON `json:"-"`
}

// eventMessagePartRemovedDataJSON contains the JSON
// metadata for the struct [EventMessagePartRemovedData]
type eventMessagePartRemovedDataJSON struct {
	MessageID   apijson.Field
	PartID      apijson.Field
	SessionID   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventMessagePartRemovedData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventMessagePartRemovedDataJSON) RawJSON() string {
	return r.raw
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
	Data EventSessionCompactedData `json:"properties,required"`
	Type       EventSessionCompactedType       `json:"type,required"`
	JSON       eventSessionCompactedJSON       `json:"-"`
}

// eventSessionCompactedJSON contains the JSON metadata for the
// struct [EventSessionCompacted]
type eventSessionCompactedJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventSessionCompacted) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventSessionCompactedJSON) RawJSON() string {
	return r.raw
}

func (r EventSessionCompacted) implementsEvent() {}

type EventSessionCompactedData struct {
	SessionID string                                               `json:"sessionID,required"`
	JSON      eventSessionCompactedDataJSON `json:"-"`
}

// eventSessionCompactedDataJSON contains the JSON metadata
// for the struct [EventSessionCompactedData]
type eventSessionCompactedDataJSON struct {
	SessionID   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventSessionCompactedData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventSessionCompactedDataJSON) RawJSON() string {
	return r.raw
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
	Data Permission                                  `json:"properties,required"`
	Type       EventPermissionUpdatedType `json:"type,required"`
	JSON       eventPermissionUpdatedJSON `json:"-"`
}

// eventPermissionUpdatedJSON contains the JSON metadata for the
// struct [EventPermissionUpdated]
type eventPermissionUpdatedJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventPermissionUpdated) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventPermissionUpdatedJSON) RawJSON() string {
	return r.raw
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
	Data EventPermissionRepliedData `json:"properties,required"`
	Type       EventPermissionRepliedType       `json:"type,required"`
	JSON       eventPermissionRepliedJSON       `json:"-"`
}

// eventPermissionRepliedJSON contains the JSON metadata for the
// struct [EventPermissionReplied]
type eventPermissionRepliedJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventPermissionReplied) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventPermissionRepliedJSON) RawJSON() string {
	return r.raw
}

func (r EventPermissionReplied) implementsEvent() {}

type EventPermissionRepliedData struct {
	PermissionID string                                                `json:"permissionID,required"`
	Response     string                                                `json:"response,required"`
	SessionID    string                                                `json:"sessionID,required"`
	JSON         eventPermissionRepliedDataJSON `json:"-"`
}

// eventPermissionRepliedDataJSON contains the JSON metadata
// for the struct [EventPermissionRepliedData]
type eventPermissionRepliedDataJSON struct {
	PermissionID apijson.Field
	Response     apijson.Field
	SessionID    apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *EventPermissionRepliedData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventPermissionRepliedDataJSON) RawJSON() string {
	return r.raw
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
	Data EventFileEditedData `json:"properties,required"`
	Type       EventFileEditedType       `json:"type,required"`
	JSON       eventFileEditedJSON       `json:"-"`
}

// eventFileEditedJSON contains the JSON metadata for the struct
// [EventFileEdited]
type eventFileEditedJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventFileEdited) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventFileEditedJSON) RawJSON() string {
	return r.raw
}

func (r EventFileEdited) implementsEvent() {}

type EventFileEditedData struct {
	File string                                         `json:"file,required"`
	JSON eventFileEditedDataJSON `json:"-"`
}

// eventFileEditedDataJSON contains the JSON metadata for
// the struct [EventFileEditedData]
type eventFileEditedDataJSON struct {
	File        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventFileEditedData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventFileEditedDataJSON) RawJSON() string {
	return r.raw
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
	Data EventFileWatcherUpdatedData `json:"properties,required"`
	Type       EventFileWatcherUpdatedType       `json:"type,required"`
	JSON       eventFileWatcherUpdatedJSON       `json:"-"`
}

// eventFileWatcherUpdatedJSON contains the JSON metadata for the
// struct [EventFileWatcherUpdated]
type eventFileWatcherUpdatedJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventFileWatcherUpdated) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventFileWatcherUpdatedJSON) RawJSON() string {
	return r.raw
}

func (r EventFileWatcherUpdated) implementsEvent() {}

type EventFileWatcherUpdatedData struct {
	Event EventFileWatcherUpdatedDataEvent `json:"event,required"`
	File  string                                                  `json:"file,required"`
	JSON  eventFileWatcherUpdatedDataJSON  `json:"-"`
}

// eventFileWatcherUpdatedDataJSON contains the JSON
// metadata for the struct [EventFileWatcherUpdatedData]
type eventFileWatcherUpdatedDataJSON struct {
	Event       apijson.Field
	File        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventFileWatcherUpdatedData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventFileWatcherUpdatedDataJSON) RawJSON() string {
	return r.raw
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
	Data EventTodoUpdatedData `json:"properties,required"`
	Type       EventTodoUpdatedType       `json:"type,required"`
	JSON       eventTodoUpdatedJSON       `json:"-"`
}

// eventTodoUpdatedJSON contains the JSON metadata for the struct
// [EventTodoUpdated]
type eventTodoUpdatedJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventTodoUpdated) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventTodoUpdatedJSON) RawJSON() string {
	return r.raw
}

func (r EventTodoUpdated) implementsEvent() {}

type EventTodoUpdatedData struct {
	SessionID string                                            `json:"sessionID,required"`
	Todos     []Todo `json:"todos,required"`
	JSON      eventTodoUpdatedDataJSON   `json:"-"`
}

// eventTodoUpdatedDataJSON contains the JSON metadata for
// the struct [EventTodoUpdatedData]
type eventTodoUpdatedDataJSON struct {
	SessionID   apijson.Field
	Todos       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventTodoUpdatedData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventTodoUpdatedDataJSON) RawJSON() string {
	return r.raw
}

type Todo struct {
	// Unique identifier for the todo item
	ID string `json:"id,required"`
	// Brief description of the task
	Content string `json:"content,required"`
	// Priority level of the task: high, medium, low
	Priority string `json:"priority,required"`
	// Current status of the task: pending, in_progress, completed, cancelled
	Status string                                              `json:"status,required"`
	JSON   todoJSON `json:"-"`
}

// todoJSON contains the JSON metadata
// for the struct [Todo]
type todoJSON struct {
	ID          apijson.Field
	Content     apijson.Field
	Priority    apijson.Field
	Status      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *Todo) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r todoJSON) RawJSON() string {
	return r.raw
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
	Data EventSessionIdleData `json:"properties,required"`
	Type       EventSessionIdleType       `json:"type,required"`
	JSON       eventSessionIdleJSON       `json:"-"`
}

// eventSessionIdleJSON contains the JSON metadata for the struct
// [EventSessionIdle]
type eventSessionIdleJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventSessionIdle) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventSessionIdleJSON) RawJSON() string {
	return r.raw
}

func (r EventSessionIdle) implementsEvent() {}

type EventSessionIdleData struct {
	SessionID string                                          `json:"sessionID,required"`
	JSON      eventSessionIdleDataJSON `json:"-"`
}

// eventSessionIdleDataJSON contains the JSON metadata for
// the struct [EventSessionIdleData]
type eventSessionIdleDataJSON struct {
	SessionID   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventSessionIdleData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventSessionIdleDataJSON) RawJSON() string {
	return r.raw
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
	Data EventSessionCreatedData `json:"properties,required"`
	Type       EventSessionCreatedType       `json:"type,required"`
	JSON       eventSessionCreatedJSON       `json:"-"`
}

// eventSessionCreatedJSON contains the JSON metadata for the
// struct [EventSessionCreated]
type eventSessionCreatedJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventSessionCreated) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventSessionCreatedJSON) RawJSON() string {
	return r.raw
}

func (r EventSessionCreated) implementsEvent() {}

type EventSessionCreatedData struct {
	Info Session                                            `json:"info,required"`
	JSON eventSessionCreatedDataJSON `json:"-"`
}

// eventSessionCreatedDataJSON contains the JSON metadata
// for the struct [EventSessionCreatedData]
type eventSessionCreatedDataJSON struct {
	Info        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventSessionCreatedData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventSessionCreatedDataJSON) RawJSON() string {
	return r.raw
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
	Data EventSessionUpdatedData `json:"properties,required"`
	Type       EventSessionUpdatedType       `json:"type,required"`
	JSON       eventSessionUpdatedJSON       `json:"-"`
}

// eventSessionUpdatedJSON contains the JSON metadata for the
// struct [EventSessionUpdated]
type eventSessionUpdatedJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventSessionUpdated) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventSessionUpdatedJSON) RawJSON() string {
	return r.raw
}

func (r EventSessionUpdated) implementsEvent() {}

type EventSessionUpdatedData struct {
	Info Session                                            `json:"info,required"`
	JSON eventSessionUpdatedDataJSON `json:"-"`
}

// eventSessionUpdatedDataJSON contains the JSON metadata
// for the struct [EventSessionUpdatedData]
type eventSessionUpdatedDataJSON struct {
	Info        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventSessionUpdatedData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventSessionUpdatedDataJSON) RawJSON() string {
	return r.raw
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
	Data EventSessionDeletedData `json:"properties,required"`
	Type       EventSessionDeletedType       `json:"type,required"`
	JSON       eventSessionDeletedJSON       `json:"-"`
}

// eventSessionDeletedJSON contains the JSON metadata for the
// struct [EventSessionDeleted]
type eventSessionDeletedJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventSessionDeleted) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventSessionDeletedJSON) RawJSON() string {
	return r.raw
}

func (r EventSessionDeleted) implementsEvent() {}

type EventSessionDeletedData struct {
	Info Session                                            `json:"info,required"`
	JSON eventSessionDeletedDataJSON `json:"-"`
}

// eventSessionDeletedDataJSON contains the JSON metadata
// for the struct [EventSessionDeletedData]
type eventSessionDeletedDataJSON struct {
	Info        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventSessionDeletedData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventSessionDeletedDataJSON) RawJSON() string {
	return r.raw
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
	Data EventSessionErrorData `json:"properties,required"`
	Type       EventSessionErrorType       `json:"type,required"`
	JSON       eventSessionErrorJSON       `json:"-"`
}

// eventSessionErrorJSON contains the JSON metadata for the struct
// [EventSessionError]
type eventSessionErrorJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventSessionError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventSessionErrorJSON) RawJSON() string {
	return r.raw
}

func (r EventSessionError) implementsEvent() {}

type EventSessionErrorData struct {
	Error     SessionError `json:"error"`
	SessionID string                                            `json:"sessionID"`
	JSON      eventSessionErrorDataJSON  `json:"-"`
}

// eventSessionErrorDataJSON contains the JSON metadata for
// the struct [EventSessionErrorData]
type eventSessionErrorDataJSON struct {
	Error       apijson.Field
	SessionID   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventSessionErrorData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventSessionErrorDataJSON) RawJSON() string {
	return r.raw
}

type SessionError struct {
	// This field can have the runtime type of [shared.ProviderAuthErrorData],
	// [shared.UnknownErrorData], [interface{}], [shared.MessageAbortedErrorData],
	// [SessionAPIErrorData].
	Data  interface{}                                           `json:"data,required"`
	Name  SessionErrorName `json:"name,required"`
	JSON  sessionErrorJSON `json:"-"`
	union SessionErrorUnion
}

// sessionErrorJSON contains the JSON metadata
// for the struct [SessionError]
type sessionErrorJSON struct {
	Data        apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r sessionErrorJSON) RawJSON() string {
	return r.raw
}

func (r *SessionError) UnmarshalJSON(data []byte) (err error) {
	*r = SessionError{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SessionErrorUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are [shared.ProviderAuthError],
// [shared.UnknownError],
// [MessageOutputLengthError],
// [shared.MessageAbortedError],
// [SessionAPIError].
func (r SessionError) AsUnion() SessionErrorUnion {
	return r.union
}

// Union satisfied by [shared.ProviderAuthError], [shared.UnknownError],
// [MessageOutputLengthError],
// [shared.MessageAbortedError] or
// [SessionAPIError].
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
	Data interface{}                                                                   `json:"data,required"`
	Name MessageOutputLengthErrorName `json:"name,required"`
	JSON messageOutputLengthErrorJSON `json:"-"`
}

// messageOutputLengthErrorJSON
// contains the JSON metadata for the struct
// [MessageOutputLengthError]
type messageOutputLengthErrorJSON struct {
	Data        apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *MessageOutputLengthError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r messageOutputLengthErrorJSON) RawJSON() string {
	return r.raw
}

func (r MessageOutputLengthError) ImplementsSessionError() {
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
	Data SessionAPIErrorData `json:"data,required"`
	Name SessionAPIErrorName `json:"name,required"`
	JSON sessionAPIErrorJSON `json:"-"`
}

// sessionAPIErrorJSON contains the JSON
// metadata for the struct
// [SessionAPIError]
type sessionAPIErrorJSON struct {
	Data        apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SessionAPIError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionAPIErrorJSON) RawJSON() string {
	return r.raw
}

func (r SessionAPIError) ImplementsSessionError() {
}

type SessionAPIErrorData struct {
	IsRetryable     bool                                                              `json:"isRetryable,required"`
	Message         string                                                            `json:"message,required"`
	ResponseBody    string                                                            `json:"responseBody"`
	ResponseHeaders map[string]string                                                 `json:"responseHeaders"`
	StatusCode      float64                                                           `json:"statusCode"`
	JSON            sessionAPIErrorDataJSON `json:"-"`
}

// sessionAPIErrorDataJSON contains the
// JSON metadata for the struct
// [SessionAPIErrorData]
type sessionAPIErrorDataJSON struct {
	IsRetryable     apijson.Field
	Message         apijson.Field
	ResponseBody    apijson.Field
	ResponseHeaders apijson.Field
	StatusCode      apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SessionAPIErrorData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionAPIErrorDataJSON) RawJSON() string {
	return r.raw
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
	Data interface{}                               `json:"properties,required"`
	Type       EventServerConnectedType `json:"type,required"`
	JSON       eventServerConnectedJSON `json:"-"`
}

// eventServerConnectedJSON contains the JSON metadata for the
// struct [EventServerConnected]
type eventServerConnectedJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventServerConnected) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventServerConnectedJSON) RawJSON() string {
	return r.raw
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
	Data EventIdeInstalledData `json:"properties,required"`
	Type       EventIdeInstalledType       `json:"type,required"`
	JSON       eventIdeInstalledJSON       `json:"-"`
}

// eventIdeInstalledJSON contains the JSON metadata for the struct
// [EventIdeInstalled]
type eventIdeInstalledJSON struct {
	Data  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventIdeInstalled) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventIdeInstalledJSON) RawJSON() string {
	return r.raw
}

func (r EventIdeInstalled) implementsEvent() {}

type EventIdeInstalledData struct {
	Ide  string                                           `json:"ide,required"`
	JSON eventIdeInstalledDataJSON `json:"-"`
}

// eventIdeInstalledDataJSON contains the JSON metadata for
// the struct [EventIdeInstalledData]
type eventIdeInstalledDataJSON struct {
	Ide         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EventIdeInstalledData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r eventIdeInstalledDataJSON) RawJSON() string {
	return r.raw
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
	Directory param.Field[string] `query:"directory"`
}

// URLQuery serializes [EventListParams]'s query parameters as `url.Values`.
func (r EventListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
