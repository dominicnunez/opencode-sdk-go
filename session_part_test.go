package opencode

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestPart_AsText_ValidTextPart(t *testing.T) {
	jsonData := `{
		"id": "part123",
		"messageID": "msg456",
		"sessionID": "sess789",
		"type": "text",
		"text": "Hello world",
		"synthetic": true,
		"metadata": {"key": "value"}
	}`

	var part Part
	if err := json.Unmarshal([]byte(jsonData), &part); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if part.Type != PartTypeText {
		t.Errorf("Expected type %s, got %s", PartTypeText, part.Type)
	}

	textPart, err := part.AsText()
	if err != nil {
		t.Fatal("AsText() should return true for type=text")
	}
	if textPart == nil {
		t.Fatal("AsText() should return non-nil TextPart")
	}
	if textPart.ID != "part123" {
		t.Errorf("Expected ID part123, got %s", textPart.ID)
	}
	if textPart.MessageID != "msg456" {
		t.Errorf("Expected MessageID msg456, got %s", textPart.MessageID)
	}
	if textPart.SessionID != "sess789" {
		t.Errorf("Expected SessionID sess789, got %s", textPart.SessionID)
	}
	if textPart.Text != "Hello world" {
		t.Errorf("Expected text 'Hello world', got %s", textPart.Text)
	}
	if !textPart.Synthetic {
		t.Error("Expected synthetic to be true")
	}
}

func TestPart_AsReasoning_ValidReasoningPart(t *testing.T) {
	jsonData := `{
		"id": "part123",
		"messageID": "msg456",
		"sessionID": "sess789",
		"type": "reasoning",
		"text": "Let me think...",
		"time": {"start": 100.5, "end": 200.5}
	}`

	var part Part
	if err := json.Unmarshal([]byte(jsonData), &part); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if part.Type != PartTypeReasoning {
		t.Errorf("Expected type %s, got %s", PartTypeReasoning, part.Type)
	}

	reasoningPart, err := part.AsReasoning()
	if err != nil {
		t.Fatal("AsReasoning() should return true for type=reasoning")
	}
	if reasoningPart == nil {
		t.Fatal("AsReasoning() should return non-nil ReasoningPart")
	}
	if reasoningPart.Text != "Let me think..." {
		t.Errorf("Expected text 'Let me think...', got %s", reasoningPart.Text)
	}
	if reasoningPart.Time.Start != 100.5 {
		t.Errorf("Expected time.start 100.5, got %f", reasoningPart.Time.Start)
	}
}

func TestPart_AsFile_ValidFilePart(t *testing.T) {
	jsonData := `{
		"id": "part123",
		"messageID": "msg456",
		"sessionID": "sess789",
		"type": "file",
		"mime": "text/plain",
		"url": "https://example.com/file.txt",
		"filename": "file.txt"
	}`

	var part Part
	if err := json.Unmarshal([]byte(jsonData), &part); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	filePart, err := part.AsFile()
	if err != nil {
		t.Fatal("AsFile() should return true for type=file")
	}
	if filePart.Mime != "text/plain" {
		t.Errorf("Expected mime text/plain, got %s", filePart.Mime)
	}
	if filePart.URL != "https://example.com/file.txt" {
		t.Errorf("Expected URL https://example.com/file.txt, got %s", filePart.URL)
	}
}

func TestPart_AsTool_ValidToolPart(t *testing.T) {
	jsonData := `{
		"id": "part123",
		"messageID": "msg456",
		"sessionID": "sess789",
		"type": "tool",
		"callID": "call999",
		"tool": "bash",
		"state": {"status": "completed"}
	}`

	var part Part
	if err := json.Unmarshal([]byte(jsonData), &part); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	toolPart, err := part.AsTool()
	if err != nil {
		t.Fatal("AsTool() should return true for type=tool")
	}
	if toolPart.CallID != "call999" {
		t.Errorf("Expected callID call999, got %s", toolPart.CallID)
	}
	if toolPart.Tool != "bash" {
		t.Errorf("Expected tool bash, got %s", toolPart.Tool)
	}
}

func TestPart_AsStepStart_ValidStepStartPart(t *testing.T) {
	jsonData := `{
		"id": "part123",
		"messageID": "msg456",
		"sessionID": "sess789",
		"type": "step-start",
		"snapshot": "snap123"
	}`

	var part Part
	if err := json.Unmarshal([]byte(jsonData), &part); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	stepStartPart, err := part.AsStepStart()
	if err != nil {
		t.Fatal("AsStepStart() should return true for type=step-start")
	}
	if stepStartPart.Snapshot != "snap123" {
		t.Errorf("Expected snapshot snap123, got %s", stepStartPart.Snapshot)
	}
}

func TestPart_AsStepFinish_ValidStepFinishPart(t *testing.T) {
	jsonData := `{
		"id": "part123",
		"messageID": "msg456",
		"sessionID": "sess789",
		"type": "step-finish",
		"cost": 0.05,
		"reason": "max_tokens",
		"tokens": {
			"cache": {"read": 100, "write": 50},
			"input": 200,
			"output": 150,
			"reasoning": 75
		}
	}`

	var part Part
	if err := json.Unmarshal([]byte(jsonData), &part); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	stepFinishPart, err := part.AsStepFinish()
	if err != nil {
		t.Fatal("AsStepFinish() should return true for type=step-finish")
	}
	if stepFinishPart.Cost != 0.05 {
		t.Errorf("Expected cost 0.05, got %f", stepFinishPart.Cost)
	}
	if stepFinishPart.Reason != "max_tokens" {
		t.Errorf("Expected reason max_tokens, got %s", stepFinishPart.Reason)
	}
	if stepFinishPart.Tokens.Input != 200 {
		t.Errorf("Expected tokens.input 200, got %d", stepFinishPart.Tokens.Input)
	}
}

func TestPart_AsSnapshot_ValidSnapshotPart(t *testing.T) {
	jsonData := `{
		"id": "part123",
		"messageID": "msg456",
		"sessionID": "sess789",
		"type": "snapshot",
		"snapshot": "snap789"
	}`

	var part Part
	if err := json.Unmarshal([]byte(jsonData), &part); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	snapshotPart, err := part.AsSnapshot()
	if err != nil {
		t.Fatal("AsSnapshot() should return true for type=snapshot")
	}
	if snapshotPart.Snapshot != "snap789" {
		t.Errorf("Expected snapshot snap789, got %s", snapshotPart.Snapshot)
	}
}

func TestPart_AsPatch_ValidPatchPart(t *testing.T) {
	jsonData := `{
		"id": "part123",
		"messageID": "msg456",
		"sessionID": "sess789",
		"type": "patch",
		"files": ["file1.go", "file2.go"],
		"hash": "abc123"
	}`

	var part Part
	if err := json.Unmarshal([]byte(jsonData), &part); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	patchPart, err := part.AsPatch()
	if err != nil {
		t.Fatal("AsPatch() should return true for type=patch")
	}
	if patchPart.Hash != "abc123" {
		t.Errorf("Expected hash abc123, got %s", patchPart.Hash)
	}
	if len(patchPart.Files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(patchPart.Files))
	}
}

func TestPart_AsAgent_ValidAgentPart(t *testing.T) {
	jsonData := `{
		"id": "part123",
		"messageID": "msg456",
		"sessionID": "sess789",
		"type": "agent",
		"name": "code-reviewer"
	}`

	var part Part
	if err := json.Unmarshal([]byte(jsonData), &part); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	agentPart, err := part.AsAgent()
	if err != nil {
		t.Fatal("AsAgent() should return true for type=agent")
	}
	if agentPart.Name != "code-reviewer" {
		t.Errorf("Expected name code-reviewer, got %s", agentPart.Name)
	}
}

func TestPart_AsRetry_ValidRetryPart(t *testing.T) {
	jsonData := `{
		"id": "part123",
		"messageID": "msg456",
		"sessionID": "sess789",
		"type": "retry",
		"attempt": 2,
		"error": {
			"name": "APIError",
			"data": {
				"isRetryable": true,
				"message": "Rate limit exceeded"
			}
		},
		"time": {"created": 123456789}
	}`

	var part Part
	if err := json.Unmarshal([]byte(jsonData), &part); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	retryPart, err := part.AsRetry()
	if err != nil {
		t.Fatal("AsRetry() should return true for type=retry")
	}
	if retryPart.Attempt != 2 {
		t.Errorf("Expected attempt 2, got %d", retryPart.Attempt)
	}
	if retryPart.Error.Name != PartRetryPartErrorNameAPIError {
		t.Errorf("Expected error name APIError, got %s", retryPart.Error.Name)
	}
	if retryPart.Error.Data.Message != "Rate limit exceeded" {
		t.Errorf("Expected error message 'Rate limit exceeded', got %s", retryPart.Error.Data.Message)
	}
}

func TestPart_WrongTypeReturnsNil(t *testing.T) {
	jsonData := `{
		"id": "part123",
		"messageID": "msg456",
		"sessionID": "sess789",
		"type": "text",
		"text": "Hello"
	}`

	var part Part
	if err := json.Unmarshal([]byte(jsonData), &part); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Try to get as reasoning when it's actually text
	reasoningPart, err := part.AsReasoning()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if reasoningPart != nil {
		t.Error("AsReasoning() should return nil for wrong type")
	}

	// Try to get as file when it's actually text
	filePart, err := part.AsFile()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if filePart != nil {
		t.Error("AsFile() should return nil for wrong type")
	}
}

func TestPart_InvalidJSON(t *testing.T) {
	jsonData := `{invalid json`

	var part Part
	err := json.Unmarshal([]byte(jsonData), &part)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestPart_MissingDiscriminator(t *testing.T) {
	jsonData := `{
		"id": "part123",
		"messageID": "msg456",
		"sessionID": "sess789"
	}`

	var part Part
	if err := json.Unmarshal([]byte(jsonData), &part); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Type should be empty
	if part.Type != "" {
		t.Errorf("Expected empty type, got %s", part.Type)
	}

	// All As* methods should return (nil, ErrWrongVariant) for wrong type
	if v, err := part.AsText(); !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	} else if v != nil {
		t.Error("AsText() should return nil for empty type")
	}
	if v, err := part.AsReasoning(); !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	} else if v != nil {
		t.Error("AsReasoning() should return nil for empty type")
	}
}

func TestPart_UnknownType(t *testing.T) {
	jsonData := `{
		"id": "part123",
		"messageID": "msg456",
		"sessionID": "sess789",
		"type": "unknown-type"
	}`

	var part Part
	if err := json.Unmarshal([]byte(jsonData), &part); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if part.Type != "unknown-type" {
		t.Errorf("Expected type unknown-type, got %s", part.Type)
	}

	if v, err := part.AsText(); !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	} else if v != nil {
		t.Error("AsText() should return nil for unknown type")
	}
	if v, err := part.AsTool(); !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	} else if v != nil {
		t.Error("AsTool() should return nil for unknown type")
	}
}
