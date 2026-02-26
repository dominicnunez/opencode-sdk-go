package opencode

import (
	"encoding/json"
	"testing"
)

// TestParamMarshal verifies that param structs marshal correctly with stdlib encoding/json
func TestParamMarshal(t *testing.T) {
	t.Run("SessionCreateParams with fields", func(t *testing.T) {
		parentID := "parent-123"
		title := "Test Session"
		params := SessionCreateParams{
			ParentID: &parentID,
			Title:    &title,
		}

		data, err := json.Marshal(params)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		// Unmarshal to verify structure
		var decoded map[string]interface{}
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if decoded["parentID"] != "parent-123" {
			t.Errorf("Expected parentID=parent-123, got %v", decoded["parentID"])
		}
		if decoded["title"] != "Test Session" {
			t.Errorf("Expected title=Test Session, got %v", decoded["title"])
		}
	})

	t.Run("SessionCreateParams omits nil fields", func(t *testing.T) {
		params := SessionCreateParams{
			ParentID: nil, // Should be omitted
			Title:    nil, // Should be omitted
		}

		data, err := json.Marshal(params)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		var decoded map[string]interface{}
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if _, exists := decoded["parentID"]; exists {
			t.Errorf("Expected parentID field to be omitted, but it exists: %v", decoded["parentID"])
		}
		if _, exists := decoded["title"]; exists {
			t.Errorf("Expected title field to be omitted, but it exists: %v", decoded["title"])
		}
	})

	t.Run("TextPartInputParam marshals correctly", func(t *testing.T) {
		params := TextPartInputParam{
			Type: TextPartInputTypeText,
			Text: "Hello, world!",
		}

		data, err := json.Marshal(params)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		var decoded map[string]interface{}
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if decoded["type"] != "text" {
			t.Errorf("Expected type=text, got %v", decoded["type"])
		}
		if decoded["text"] != "Hello, world!" {
			t.Errorf("Expected text=Hello, world!, got %v", decoded["text"])
		}
	})

	t.Run("AgentPartInputParam with nested struct", func(t *testing.T) {
		agentID := "agent-456"
		params := AgentPartInputParam{
			Type: AgentPartInputTypeAgent,
			Name: "Test Agent",
			ID:   &agentID,
			Source: &AgentPartInputSourceParam{
				Start: 0,
				End:   100,
				Value: "test-value",
			},
		}

		data, err := json.Marshal(params)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		var decoded map[string]interface{}
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if decoded["type"] != "agent" {
			t.Errorf("Expected type=agent, got %v", decoded["type"])
		}
		if decoded["name"] != "Test Agent" {
			t.Errorf("Expected name=Test Agent, got %v", decoded["name"])
		}
		if decoded["id"] != "agent-456" {
			t.Errorf("Expected id=agent-456, got %v", decoded["id"])
		}

		source, ok := decoded["source"].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected source to be object, got %T", decoded["source"])
		}
		// JSON numbers are float64 by default
		if source["start"] != float64(0) {
			t.Errorf("Expected source.start=0, got %v", source["start"])
		}
		if source["end"] != float64(100) {
			t.Errorf("Expected source.end=100, got %v", source["end"])
		}
		if source["value"] != "test-value" {
			t.Errorf("Expected source.value=test-value, got %v", source["value"])
		}
	})

	t.Run("FileSourceParam marshals correctly", func(t *testing.T) {
		params := FileSourceParam{
			Type: FileSourceTypeFile,
			Path: "/path/to/file.go",
			Text: FilePartSourceTextParam{
				Start: 0,
				End:   100,
				Value: "file content",
			},
		}

		data, err := json.Marshal(params)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		var decoded map[string]interface{}
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if decoded["type"] != "file" {
			t.Errorf("Expected type=file, got %v", decoded["type"])
		}
		if decoded["path"] != "/path/to/file.go" {
			t.Errorf("Expected path=/path/to/file.go, got %v", decoded["path"])
		}

		text, ok := decoded["text"].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected text to be object, got %T", decoded["text"])
		}
		if text["value"] != "file content" {
			t.Errorf("Expected text.value=file content, got %v", text["value"])
		}
	})

	t.Run("AppLogParams marshals correctly", func(t *testing.T) {
		params := AppLogParams{
			Level:   LogLevelInfo,
			Message: "Test log message",
			Service: "test-service",
		}

		data, err := json.Marshal(params)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		var decoded map[string]interface{}
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if decoded["level"] != "info" {
			t.Errorf("Expected level=info, got %v", decoded["level"])
		}
		if decoded["message"] != "Test log message" {
			t.Errorf("Expected message=Test log message, got %v", decoded["message"])
		}
		if decoded["service"] != "test-service" {
			t.Errorf("Expected service=test-service, got %v", decoded["service"])
		}
	})
}

// TestRoundTrip verifies that params can be marshaled and unmarshaled correctly
func TestRoundTrip(t *testing.T) {
	t.Run("SessionPromptParams round trip", func(t *testing.T) {
		agent := "test-agent"
		original := SessionPromptParams{
			Parts: []SessionPromptParamsPartUnion{},
			Agent: &agent,
		}

		// Marshal
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		// Unmarshal
		var decoded SessionPromptParams
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		// Verify
		if decoded.Agent == nil {
			t.Fatal("Expected agent to be non-nil")
		}
		if *decoded.Agent != *original.Agent {
			t.Errorf("Expected agent=%s, got %s", *original.Agent, *decoded.Agent)
		}
	})

	t.Run("SessionCommandParams round trip with complex fields", func(t *testing.T) {
		model := "claude-opus-4"
		original := SessionCommandParams{
			Command:   "npm install",
			Arguments: "--save",
			Model:     &model,
		}

		// Marshal
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		// Unmarshal
		var decoded SessionCommandParams
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		// Verify
		if decoded.Command != original.Command {
			t.Errorf("Expected command=%s, got %s", original.Command, decoded.Command)
		}
		if decoded.Arguments != original.Arguments {
			t.Errorf("Expected arguments=%s, got %s", original.Arguments, decoded.Arguments)
		}
		if decoded.Model == nil {
			t.Fatal("Expected model to be non-nil")
		}
		if *decoded.Model != *original.Model {
			t.Errorf("Expected model=%s, got %s", *original.Model, *decoded.Model)
		}
	})
}
