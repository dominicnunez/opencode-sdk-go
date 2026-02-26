package opencode_test

import (
	"encoding/json"
	"testing"

	"github.com/dominicnunez/opencode-sdk-go"
)

func TestMessage_AsUser_ValidUserMessage(t *testing.T) {
	jsonData := `{
		"id": "msg123",
		"role": "user",
		"sessionID": "ses456",
		"time": {
			"created": 1234567890.5
		},
		"summary": {
			"diffs": [],
			"body": "test body",
			"title": "test title"
		}
	}`

	var msg opencode.Message
	if err := json.Unmarshal([]byte(jsonData), &msg); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify discriminator fields are populated
	if msg.ID != "msg123" {
		t.Errorf("expected ID msg123, got %s", msg.ID)
	}
	if msg.Role != opencode.MessageRoleUser {
		t.Errorf("expected role user, got %s", msg.Role)
	}
	if msg.SessionID != "ses456" {
		t.Errorf("expected sessionID ses456, got %s", msg.SessionID)
	}

	// Test AsUser - should succeed
	userMsg, err := msg.AsUser()
	if err != nil {
		t.Fatal("AsUser should return true for user role")
	}
	if userMsg == nil {
		t.Fatal("userMsg should not be nil")
	}
	if userMsg.ID != "msg123" {
		t.Errorf("expected ID msg123, got %s", userMsg.ID)
	}
	if userMsg.Time.Created != 1234567890.5 {
		t.Errorf("expected created time 1234567890.5, got %f", userMsg.Time.Created)
	}
	if userMsg.Summary.Body != "test body" {
		t.Errorf("expected body 'test body', got %s", userMsg.Summary.Body)
	}

	// Test AsAssistant - should fail
	assistantMsg, err := msg.AsAssistant()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if assistantMsg != nil {
		t.Error("AsAssistant should return false for user role")
	}
	if assistantMsg != nil {
		t.Error("assistantMsg should be nil for user role")
	}
}

func TestMessage_AsAssistant_ValidAssistantMessage(t *testing.T) {
	jsonData := `{
		"id": "msg789",
		"role": "assistant",
		"sessionID": "ses456",
		"cost": 0.05,
		"mode": "code",
		"modelID": "claude-3",
		"parentID": "msg123",
		"path": {
			"cwd": "/home/user",
			"root": "/home/user/project"
		},
		"providerID": "anthropic",
		"system": ["sys1", "sys2"],
		"time": {
			"created": 1234567890.5,
			"completed": 1234567900.5
		},
		"tokens": {
			"cache": {
				"read": 100,
				"write": 50
			},
			"input": 1000,
			"output": 500,
			"reasoning": 250
		},
		"summary": false
	}`

	var msg opencode.Message
	if err := json.Unmarshal([]byte(jsonData), &msg); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify discriminator fields
	if msg.ID != "msg789" {
		t.Errorf("expected ID msg789, got %s", msg.ID)
	}
	if msg.Role != opencode.MessageRoleAssistant {
		t.Errorf("expected role assistant, got %s", msg.Role)
	}

	// Test AsAssistant - should succeed
	assistantMsg, err := msg.AsAssistant()
	if err != nil {
		t.Fatal("AsAssistant should return true for assistant role")
	}
	if assistantMsg == nil {
		t.Fatal("assistantMsg should not be nil")
	}
	if assistantMsg.ID != "msg789" {
		t.Errorf("expected ID msg789, got %s", assistantMsg.ID)
	}
	if assistantMsg.Cost != 0.05 {
		t.Errorf("expected cost 0.05, got %f", assistantMsg.Cost)
	}
	if assistantMsg.ModelID != "claude-3" {
		t.Errorf("expected modelID claude-3, got %s", assistantMsg.ModelID)
	}
	if len(assistantMsg.System) != 2 {
		t.Errorf("expected 2 system items, got %d", len(assistantMsg.System))
	}
	if assistantMsg.Tokens.Input != 1000 {
		t.Errorf("expected input tokens 1000, got %d", assistantMsg.Tokens.Input)
	}

	// Test AsUser - should fail
	userMsg, err := msg.AsUser()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userMsg != nil {
		t.Error("AsUser should return false for assistant role")
	}
	if userMsg != nil {
		t.Error("userMsg should be nil for assistant role")
	}
}

func TestMessage_InvalidJSON(t *testing.T) {
	invalidJSON := `{"id": "msg123", "role": "invalid"}`

	var msg opencode.Message
	if err := json.Unmarshal([]byte(invalidJSON), &msg); err != nil {
		t.Fatalf("unmarshal should not fail on unknown role: %v", err)
	}

	// Both As* methods should return false for invalid role
	userMsg, err := msg.AsUser()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userMsg != nil {
		t.Error("AsUser should return false for invalid role")
	}
	if userMsg != nil {
		t.Error("userMsg should be nil for invalid role")
	}

	assistantMsg, err := msg.AsAssistant()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if assistantMsg != nil {
		t.Error("AsAssistant should return false for invalid role")
	}
	if assistantMsg != nil {
		t.Error("assistantMsg should be nil for invalid role")
	}
}

func TestMessage_MalformedJSON(t *testing.T) {
	malformedJSON := `{"id": "msg123", "role": "user", "time": "not-a-time-object"}`

	var msg opencode.Message
	// UnmarshalJSON should succeed for Message (only peeks at discriminator)
	if err := json.Unmarshal([]byte(malformedJSON), &msg); err != nil {
		t.Fatalf("unmarshal should succeed for discriminator peek: %v", err)
	}

	// AsUser should return error when trying to unmarshal the full UserMessage
	userMsg, err := msg.AsUser()
	if err == nil {
		t.Fatal("AsUser should return error when full unmarshal fails")
	}
	if userMsg != nil {
		t.Error("userMsg should be nil when unmarshal fails")
	}
}

func TestMessage_EmptyJSON(t *testing.T) {
	emptyJSON := `{}`

	var msg opencode.Message
	if err := json.Unmarshal([]byte(emptyJSON), &msg); err != nil {
		t.Fatalf("unmarshal should not fail on empty object: %v", err)
	}

	// Both As* methods should return false for empty role
	userMsg, err := msg.AsUser()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userMsg != nil {
		t.Error("AsUser should return false for empty role")
	}
	if userMsg != nil {
		t.Error("userMsg should be nil for empty role")
	}

	assistantMsg, err := msg.AsAssistant()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if assistantMsg != nil {
		t.Error("AsAssistant should return false for empty role")
	}
	if assistantMsg != nil {
		t.Error("assistantMsg should be nil for empty role")
	}
}
