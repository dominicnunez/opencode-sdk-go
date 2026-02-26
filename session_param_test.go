package opencode_test

import (
	"encoding/json"
	"testing"

	"github.com/dominicnunez/opencode-sdk-go"
)

// TestSessionCreateParams_DirectTypes verifies that required fields are direct types
// and optional fields are pointers with proper JSON marshaling
func TestSessionCreateParams_DirectTypes(t *testing.T) {
	params := &opencode.SessionCreateParams{
		Directory: opencode.Ptr("/tmp/test"),
		ParentID:  opencode.Ptr("parent123"),
		Title:     opencode.Ptr("Test Session"),
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	// Verify JSON structure
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	if result["parentID"] != "parent123" {
		t.Errorf("Expected parentID to be parent123, got %v", result["parentID"])
	}
	if result["title"] != "Test Session" {
		t.Errorf("Expected title to be 'Test Session', got %v", result["title"])
	}
}

// TestSessionCommandParams_RequiredFields verifies that required fields work as direct types
func TestSessionCommandParams_RequiredFields(t *testing.T) {
	params := &opencode.SessionCommandParams{
		Arguments: "test args",
		Command:   "test command",
		Directory: opencode.Ptr("/tmp"),
		Agent:     opencode.Ptr("agent1"),
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	// Required fields should always be present
	if result["arguments"] != "test args" {
		t.Errorf("Expected arguments to be 'test args', got %v", result["arguments"])
	}
	if result["command"] != "test command" {
		t.Errorf("Expected command to be 'test command', got %v", result["command"])
	}

	// Optional fields should be present when set
	if result["agent"] != "agent1" {
		t.Errorf("Expected agent to be 'agent1', got %v", result["agent"])
	}
}

// TestSessionPromptParams_ComplexTypes verifies nested structs and slices work correctly
func TestSessionPromptParams_ComplexTypes(t *testing.T) {
	// Use AgentPartInputParam which has simpler structure (no nested pointers with MarshalJSON)
	params := &opencode.SessionPromptParams{
		Parts: []opencode.SessionPromptParamsPartUnion{
			opencode.AgentPartInputParam{
				Name:   "test-agent",
				Type:   opencode.AgentPartInputTypeAgent,
				ID:     opencode.Ptr("agent1"),
				Source: &opencode.AgentPartInputSourceParam{
					End:   100,
					Start: 0,
					Value: "test",
				},
			},
		},
		Agent: opencode.Ptr("test-agent"),
		Model: &opencode.SessionPromptParamsModel{
			ModelID:    "gpt-4",
			ProviderID: "openai",
		},
		NoReply: opencode.Ptr(false),
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	// Verify parts array
	parts, ok := result["parts"].([]interface{})
	if !ok || len(parts) != 1 {
		t.Errorf("Expected parts to be array with 1 element, got %v", result["parts"])
	}

	// Verify optional fields
	if result["agent"] != "test-agent" {
		t.Errorf("Expected agent to be 'test-agent', got %v", result["agent"])
	}
	if result["noReply"] != false {
		t.Errorf("Expected noReply to be false, got %v", result["noReply"])
	}

	// Verify nested model struct
	model, ok := result["model"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected model to be object, got %T", result["model"])
	}
	if model["modelID"] != "gpt-4" {
		t.Errorf("Expected modelID to be 'gpt-4', got %v", model["modelID"])
	}
	if model["providerID"] != "openai" {
		t.Errorf("Expected providerID to be 'openai', got %v", model["providerID"])
	}
}

// TestSessionInitParams_MixedRequiredOptional verifies mixed required/optional fields
func TestSessionInitParams_MixedRequiredOptional(t *testing.T) {
	params := &opencode.SessionInitParams{
		MessageID:  "msg123",
		ModelID:    "gpt-4",
		ProviderID: "openai",
		Directory:  opencode.Ptr("/tmp/test"),
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	// Required fields should be present
	if result["messageID"] != "msg123" {
		t.Errorf("Expected messageID to be 'msg123', got %v", result["messageID"])
	}
	if result["modelID"] != "gpt-4" {
		t.Errorf("Expected modelID to be 'gpt-4', got %v", result["modelID"])
	}
	if result["providerID"] != "openai" {
		t.Errorf("Expected providerID to be 'openai', got %v", result["providerID"])
	}
}

// TestAgentPartInputSourceParam_DirectTypes verifies all required fields are direct types
func TestAgentPartInputSourceParam_DirectTypes(t *testing.T) {
	param := opencode.AgentPartInputSourceParam{
		End:   100,
		Start: 0,
		Value: "test value",
	}

	data, err := json.Marshal(param)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	// All fields are required - verify they're all present
	if result["end"] != float64(100) {
		t.Errorf("Expected end to be 100, got %v", result["end"])
	}
	if result["start"] != float64(0) {
		t.Errorf("Expected start to be 0, got %v", result["start"])
	}
	if result["value"] != "test value" {
		t.Errorf("Expected value to be 'test value', got %v", result["value"])
	}
}
