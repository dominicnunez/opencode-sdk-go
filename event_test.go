package opencode

import (
	"encoding/json"
	"errors"
	"testing"
)

// Test Event discriminated union - AsInstallationUpdated
func TestEvent_AsInstallationUpdated(t *testing.T) {
	jsonData := `{"type":"installation.updated","properties":{"version":"1.2.3"}}`
	var event Event
	if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	evt, err := event.AsInstallationUpdated()
	if err != nil {
		t.Fatal("Expected AsInstallationUpdated to return true")
	}
	if evt.Data.Version != "1.2.3" {
		t.Errorf("Expected version 1.2.3, got %s", evt.Data.Version)
	}
	if evt.Type != EventInstallationUpdatedTypeInstallationUpdated {
		t.Errorf("Expected type installation.updated, got %s", evt.Type)
	}
}

// Test Event discriminated union - AsMessageUpdated
func TestEvent_AsMessageUpdated(t *testing.T) {
	jsonData := `{"type":"message.updated","properties":{"info":{"id":"msg123","sessionID":"sess456","role":"user"}}}`
	var event Event
	if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	evt, err := event.AsMessageUpdated()
	if err != nil {
		t.Fatal("Expected AsMessageUpdated to return true")
	}
	if evt.Data.Info.ID != "msg123" {
		t.Errorf("Expected message ID msg123, got %s", evt.Data.Info.ID)
	}
	if evt.Data.Info.SessionID != "sess456" {
		t.Errorf("Expected session ID sess456, got %s", evt.Data.Info.SessionID)
	}
}

// Test Event discriminated union - AsSessionCreated
func TestEvent_AsSessionCreated(t *testing.T) {
	jsonData := `{"type":"session.created","properties":{"info":{"id":"sess789","directory":"/test","projectID":"proj1","time":{"created":1234567890,"updated":1234567890},"title":"Test Session","version":"1.0"}}}`
	var event Event
	if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	evt, err := event.AsSessionCreated()
	if err != nil {
		t.Fatal("Expected AsSessionCreated to return true")
	}
	if evt.Data.Info.ID != "sess789" {
		t.Errorf("Expected session ID sess789, got %s", evt.Data.Info.ID)
	}
	if evt.Data.Info.Title != "Test Session" {
		t.Errorf("Expected title 'Test Session', got %s", evt.Data.Info.Title)
	}
}

// Test Event discriminated union - AsTodoUpdated
func TestEvent_AsTodoUpdated(t *testing.T) {
	jsonData := `{"type":"todo.updated","properties":{"sessionID":"sess123","todos":[{"id":"todo1","content":"Fix bug","priority":"high","status":"pending"}]}}`
	var event Event
	if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	evt, err := event.AsTodoUpdated()
	if err != nil {
		t.Fatal("Expected AsTodoUpdated to return true")
	}
	if evt.Data.SessionID != "sess123" {
		t.Errorf("Expected session ID sess123, got %s", evt.Data.SessionID)
	}
	if len(evt.Data.Todos) != 1 {
		t.Fatalf("Expected 1 todo, got %d", len(evt.Data.Todos))
	}
	if evt.Data.Todos[0].ID != "todo1" {
		t.Errorf("Expected todo ID todo1, got %s", evt.Data.Todos[0].ID)
	}
	if evt.Data.Todos[0].Content != "Fix bug" {
		t.Errorf("Expected content 'Fix bug', got %s", evt.Data.Todos[0].Content)
	}
}

// Test Event discriminated union - AsFileEdited
func TestEvent_AsFileEdited(t *testing.T) {
	jsonData := `{"type":"file.edited","properties":{"file":"main.go"}}`
	var event Event
	if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	evt, err := event.AsFileEdited()
	if err != nil {
		t.Fatal("Expected AsFileEdited to return true")
	}
	if evt.Data.File != "main.go" {
		t.Errorf("Expected file main.go, got %s", evt.Data.File)
	}
}

// Test Event discriminated union - wrong type returns false
func TestEvent_WrongType(t *testing.T) {
	jsonData := `{"type":"session.created","properties":{"info":{"id":"sess789","status":"idle"}}}`
	var event Event
	if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Try to get as wrong type - should return (nil, ErrWrongVariant)
	evt, err := event.AsInstallationUpdated()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if evt != nil {
		t.Error("Expected nil event when type doesn't match")
	}
}

// Test Event discriminated union - invalid JSON
func TestEvent_InvalidJSON(t *testing.T) {
	jsonData := `{"type":"installation.updated","properties":{"version":123}}`
	var event Event
	if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Type is correct but data is malformed
	evt, err := event.AsInstallationUpdated()
	if err != nil {
		t.Logf("AsInstallationUpdated returned error (expected for malformed data): %v", err)
	} else if evt != nil {
		if evt.Data.Version != "" {
			t.Logf("Version parsed as: %s", evt.Data.Version)
		}
	}
}

// Test Event discriminated union - missing type
func TestEvent_MissingType(t *testing.T) {
	jsonData := `{"properties":{"version":"1.0.0"}}`
	var event Event
	// Should unmarshal but with empty Type
	if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if event.Type != "" {
		t.Errorf("Expected empty type, got %s", event.Type)
	}

	// Trying to get any specific type should return (nil, ErrWrongVariant)
	val, err := event.AsInstallationUpdated()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if val != nil {
		t.Error("Expected AsInstallationUpdated to return nil when type is missing")
	}
}

// Test Event discriminated union - AsSessionError
func TestEvent_AsSessionError(t *testing.T) {
	jsonData := `{"type":"session.error","properties":{"sessionID":"sess123","error":{"name":"APIError","data":{"message":"Failed","isRetryable":true}}}}`
	var event Event
	if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	evt, err := event.AsSessionError()
	if err != nil {
		t.Fatal("Expected AsSessionError to return true")
	}
	if evt.Data.SessionID == nil || *evt.Data.SessionID != "sess123" {
		t.Errorf("Expected session ID sess123")
	}
	if evt.Data.Error == nil {
		t.Fatal("Expected error to be non-nil")
	}
}

// Test all 19 event types can be discriminated
func TestEvent_AllTypes(t *testing.T) {
	testCases := []struct {
		name      string
		eventType EventType
		jsonData  string
		checkFunc func(Event) bool
	}{
		{"InstallationUpdated", EventTypeInstallationUpdated, `{"type":"installation.updated","properties":{"version":"1.0"}}`, func(e Event) bool { v, err := e.AsInstallationUpdated(); return err == nil && v != nil }},
		{"LspClientDiagnostics", EventTypeLspClientDiagnostics, `{"type":"lsp.client.diagnostics","properties":{"path":"/test","serverID":"srv1"}}`, func(e Event) bool { v, err := e.AsLspClientDiagnostics(); return err == nil && v != nil }},
		{"MessageUpdated", EventTypeMessageUpdated, `{"type":"message.updated","properties":{"info":{"id":"m1","sessionID":"s1","role":"user"}}}`, func(e Event) bool { v, err := e.AsMessageUpdated(); return err == nil && v != nil }},
		{"MessageRemoved", EventTypeMessageRemoved, `{"type":"message.removed","properties":{"messageID":"m1","sessionID":"s1"}}`, func(e Event) bool { v, err := e.AsMessageRemoved(); return err == nil && v != nil }},
		{"MessagePartUpdated", EventTypeMessagePartUpdated, `{"type":"message.part.updated","properties":{"part":{"id":"p1","messageID":"m1","sessionID":"s1","type":"text"}}}`, func(e Event) bool { v, err := e.AsMessagePartUpdated(); return err == nil && v != nil }},
		{"MessagePartRemoved", EventTypeMessagePartRemoved, `{"type":"message.part.removed","properties":{"messageID":"m1","partID":"p1","sessionID":"s1"}}`, func(e Event) bool { v, err := e.AsMessagePartRemoved(); return err == nil && v != nil }},
		{"SessionCompacted", EventTypeSessionCompacted, `{"type":"session.compacted","properties":{"sessionID":"s1"}}`, func(e Event) bool { v, err := e.AsSessionCompacted(); return err == nil && v != nil }},
		{"PermissionUpdated", EventTypePermissionUpdated, `{"type":"permission.updated","properties":{"id":"perm1","sessionID":"s1","status":"pending"}}`, func(e Event) bool { v, err := e.AsPermissionUpdated(); return err == nil && v != nil }},
		{"PermissionReplied", EventTypePermissionReplied, `{"type":"permission.replied","properties":{"permissionID":"p1","response":"allow","sessionID":"s1"}}`, func(e Event) bool { v, err := e.AsPermissionReplied(); return err == nil && v != nil }},
		{"FileEdited", EventTypeFileEdited, `{"type":"file.edited","properties":{"file":"test.go"}}`, func(e Event) bool { v, err := e.AsFileEdited(); return err == nil && v != nil }},
		{"FileWatcherUpdated", EventTypeFileWatcherUpdated, `{"type":"file.watcher.updated","properties":{"event":"change","file":"test.go"}}`, func(e Event) bool { v, err := e.AsFileWatcherUpdated(); return err == nil && v != nil }},
		{"TodoUpdated", EventTypeTodoUpdated, `{"type":"todo.updated","properties":{"sessionID":"s1","todos":[]}}`, func(e Event) bool { v, err := e.AsTodoUpdated(); return err == nil && v != nil }},
		{"SessionIdle", EventTypeSessionIdle, `{"type":"session.idle","properties":{"sessionID":"s1"}}`, func(e Event) bool { v, err := e.AsSessionIdle(); return err == nil && v != nil }},
		{"SessionCreated", EventTypeSessionCreated, `{"type":"session.created","properties":{"info":{"id":"s1","directory":"/","projectID":"p1","time":{"created":0,"updated":0},"title":"","version":""}}}`, func(e Event) bool { v, err := e.AsSessionCreated(); return err == nil && v != nil }},
		{"SessionUpdated", EventTypeSessionUpdated, `{"type":"session.updated","properties":{"info":{"id":"s1","directory":"/","projectID":"p1","time":{"created":0,"updated":0},"title":"","version":""}}}`, func(e Event) bool { v, err := e.AsSessionUpdated(); return err == nil && v != nil }},
		{"SessionDeleted", EventTypeSessionDeleted, `{"type":"session.deleted","properties":{"info":{"id":"s1","directory":"/","projectID":"p1","time":{"created":0,"updated":0},"title":"","version":""}}}`, func(e Event) bool { v, err := e.AsSessionDeleted(); return err == nil && v != nil }},
		{"SessionError", EventTypeSessionError, `{"type":"session.error","properties":{"sessionID":"s1"}}`, func(e Event) bool { v, err := e.AsSessionError(); return err == nil && v != nil }},
		{"ServerConnected", EventTypeServerConnected, `{"type":"server.connected","properties":{}}`, func(e Event) bool { v, err := e.AsServerConnected(); return err == nil && v != nil }},
		{"IdeInstalled", EventTypeIdeInstalled, `{"type":"ide.installed","properties":{"ide":"vscode"}}`, func(e Event) bool { v, err := e.AsIdeInstalled(); return err == nil && v != nil }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var event Event
			if err := json.Unmarshal([]byte(tc.jsonData), &event); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			if event.Type != tc.eventType {
				t.Errorf("Expected type %s, got %s", tc.eventType, event.Type)
			}

			if !tc.checkFunc(event) {
				t.Errorf("Failed to convert to %s", tc.name)
			}
		})
	}
}
