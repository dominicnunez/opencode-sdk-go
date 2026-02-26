package opencode

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestToolPartState_AsPending_ValidPendingState(t *testing.T) {
	jsonData := `{"status": "pending"}`
	var state ToolPartState
	if err := json.Unmarshal([]byte(jsonData), &state); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if state.Status != ToolPartStateStatusPending {
		t.Errorf("expected status pending, got %s", state.Status)
	}

	pending, err := state.AsPending()
	if err != nil {
		t.Fatal("AsPending() should return true for pending status")
	}
	if pending == nil {
		t.Fatal("AsPending() should return non-nil pointer")
	}
	if pending.Status != ToolStatePendingStatusPending {
		t.Errorf("expected pending status, got %s", pending.Status)
	}
}

func TestToolPartState_AsRunning_ValidRunningState(t *testing.T) {
	jsonData := `{
		"status": "running",
		"input": {"command": "test"},
		"time": {"start": 1234567890.5},
		"metadata": {"key": "value"},
		"title": "Running Test"
	}`
	var state ToolPartState
	if err := json.Unmarshal([]byte(jsonData), &state); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if state.Status != ToolPartStateStatusRunning {
		t.Errorf("expected status running, got %s", state.Status)
	}

	running, err := state.AsRunning()
	if err != nil {
		t.Fatal("AsRunning() should return true for running status")
	}
	if running == nil {
		t.Fatal("AsRunning() should return non-nil pointer")
	}
	if running.Status != ToolStateRunningStatusRunning {
		t.Errorf("expected running status, got %s", running.Status)
	}
	if running.Title != "Running Test" {
		t.Errorf("expected title 'Running Test', got %s", running.Title)
	}
	if running.Time.Start != 1234567890.5 {
		t.Errorf("expected start time 1234567890.5, got %f", running.Time.Start)
	}
}

func TestToolPartState_AsCompleted_ValidCompletedState(t *testing.T) {
	jsonData := `{
		"status": "completed",
		"input": {"command": "test"},
		"metadata": {"key": "value"},
		"output": "Test completed successfully",
		"time": {"start": 1234567890.5, "end": 1234567900.5, "compacted": 10.0},
		"title": "Completed Test",
		"attachments": []
	}`
	var state ToolPartState
	if err := json.Unmarshal([]byte(jsonData), &state); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if state.Status != ToolPartStateStatusCompleted {
		t.Errorf("expected status completed, got %s", state.Status)
	}

	completed, err := state.AsCompleted()
	if err != nil {
		t.Fatal("AsCompleted() should return true for completed status")
	}
	if completed == nil {
		t.Fatal("AsCompleted() should return non-nil pointer")
	}
	if completed.Status != ToolStateCompletedStatusCompleted {
		t.Errorf("expected completed status, got %s", completed.Status)
	}
	if completed.Title != "Completed Test" {
		t.Errorf("expected title 'Completed Test', got %s", completed.Title)
	}
	if completed.Output != "Test completed successfully" {
		t.Errorf("expected output 'Test completed successfully', got %s", completed.Output)
	}
	if completed.Time.Start != 1234567890.5 {
		t.Errorf("expected start time 1234567890.5, got %f", completed.Time.Start)
	}
	if completed.Time.End != 1234567900.5 {
		t.Errorf("expected end time 1234567900.5, got %f", completed.Time.End)
	}
}

func TestToolPartState_AsError_ValidErrorState(t *testing.T) {
	jsonData := `{
		"status": "error",
		"error": "Command execution failed",
		"input": {"command": "test"},
		"time": {"start": 1234567890.5, "end": 1234567900.5},
		"metadata": {"errorCode": 500}
	}`
	var state ToolPartState
	if err := json.Unmarshal([]byte(jsonData), &state); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if state.Status != ToolPartStateStatusError {
		t.Errorf("expected status error, got %s", state.Status)
	}

	errState, err := state.AsError()
	if err != nil {
		t.Fatal("AsError() should return true for error status")
	}
	if errState == nil {
		t.Fatal("AsError() should return non-nil pointer")
	}
	if errState.Status != ToolStateErrorStatusError {
		t.Errorf("expected error status, got %s", errState.Status)
	}
	if errState.Error != "Command execution failed" {
		t.Errorf("expected error 'Command execution failed', got %s", errState.Error)
	}
	if errState.Time.Start != 1234567890.5 {
		t.Errorf("expected start time 1234567890.5, got %f", errState.Time.Start)
	}
	if errState.Time.End != 1234567900.5 {
		t.Errorf("expected end time 1234567900.5, got %f", errState.Time.End)
	}
}

func TestToolPartState_WrongTypeReturnsNilFalse(t *testing.T) {
	// Test that asking for the wrong type returns (nil, false)
	jsonData := `{"status": "pending"}`
	var state ToolPartState
	if err := json.Unmarshal([]byte(jsonData), &state); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// Should succeed for pending
	pending, err := state.AsPending()
	if err != nil {
		t.Fatalf("AsPending error: %v", err)
	}
	if pending == nil {
		t.Error("AsPending() should succeed for pending status")
	}

	// Should return ErrWrongVariant for other types
	if running, err := state.AsRunning(); !errors.Is(err, ErrWrongVariant) || running != nil {
		t.Error("AsRunning() should return (nil, ErrWrongVariant) for pending status")
	}
	if completed, err := state.AsCompleted(); !errors.Is(err, ErrWrongVariant) || completed != nil {
		t.Error("AsCompleted() should return (nil, ErrWrongVariant) for pending status")
	}
	if errState, err := state.AsError(); !errors.Is(err, ErrWrongVariant) || errState != nil {
		t.Error("AsError() should return (nil, ErrWrongVariant) for pending status")
	}
}

func TestToolPartState_InvalidJSON(t *testing.T) {
	jsonData := `{"status": "running", "time": "invalid"}`
	var state ToolPartState
	// Should unmarshal the discriminator successfully
	if err := json.Unmarshal([]byte(jsonData), &state); err != nil {
		t.Fatalf("unmarshal discriminator failed: %v", err)
	}

	// AsRunning should return error due to invalid time field
	running, err := state.AsRunning()
	if err == nil {
		t.Error("AsRunning() should return error for invalid JSON")
	}
	if running != nil {
		t.Error("AsRunning() should return nil for invalid JSON")
	}
}

func TestToolPartState_UnknownStatus(t *testing.T) {
	jsonData := `{"status": "unknown"}`
	var state ToolPartState
	if err := json.Unmarshal([]byte(jsonData), &state); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// All As* methods should return (nil, ErrWrongVariant)
	if pending, err := state.AsPending(); !errors.Is(err, ErrWrongVariant) || pending != nil {
		t.Error("AsPending() should return (nil, ErrWrongVariant) for unknown status")
	}
	if running, err := state.AsRunning(); !errors.Is(err, ErrWrongVariant) || running != nil {
		t.Error("AsRunning() should return (nil, ErrWrongVariant) for unknown status")
	}
	if completed, err := state.AsCompleted(); !errors.Is(err, ErrWrongVariant) || completed != nil {
		t.Error("AsCompleted() should return (nil, ErrWrongVariant) for unknown status")
	}
	if errState, err := state.AsError(); !errors.Is(err, ErrWrongVariant) || errState != nil {
		t.Error("AsError() should return (nil, ErrWrongVariant) for unknown status")
	}
}

func TestToolPartState_MissingStatus(t *testing.T) {
	jsonData := `{"input": {"command": "test"}}`
	var state ToolPartState
	if err := json.Unmarshal([]byte(jsonData), &state); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// Status should be empty string (zero value)
	if state.Status != "" {
		t.Errorf("expected empty status, got %s", state.Status)
	}

	// All As* methods should return (nil, ErrWrongVariant)
	if pending, err := state.AsPending(); !errors.Is(err, ErrWrongVariant) || pending != nil {
		t.Error("AsPending() should return (nil, ErrWrongVariant) for missing status")
	}
	if running, err := state.AsRunning(); !errors.Is(err, ErrWrongVariant) || running != nil {
		t.Error("AsRunning() should return (nil, ErrWrongVariant) for missing status")
	}
	if completed, err := state.AsCompleted(); !errors.Is(err, ErrWrongVariant) || completed != nil {
		t.Error("AsCompleted() should return (nil, ErrWrongVariant) for missing status")
	}
	if errState, err := state.AsError(); !errors.Is(err, ErrWrongVariant) || errState != nil {
		t.Error("AsError() should return (nil, ErrWrongVariant) for missing status")
	}
}

func TestToolPartState_EmptyJSON(t *testing.T) {
	jsonData := `{}`
	var state ToolPartState
	if err := json.Unmarshal([]byte(jsonData), &state); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// Status should be empty string (zero value)
	if state.Status != "" {
		t.Errorf("expected empty status, got %s", state.Status)
	}

	// All As* methods should return (nil, ErrWrongVariant)
	if pending, err := state.AsPending(); !errors.Is(err, ErrWrongVariant) || pending != nil {
		t.Error("AsPending() should return (nil, ErrWrongVariant) for empty JSON")
	}
}
