package shared

import (
	"encoding/json"
	"testing"
)

func TestMessageAbortedErrorUnmarshal(t *testing.T) {
	jsonData := `{
		"name": "MessageAbortedError",
		"data": {
			"message": "Message was aborted"
		}
	}`

	var err MessageAbortedError
	if unmarshalErr := json.Unmarshal([]byte(jsonData), &err); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	if err.Name != MessageAbortedErrorNameMessageAbortedError {
		t.Errorf("Expected Name to be %q, got %q", MessageAbortedErrorNameMessageAbortedError, err.Name)
	}

	if err.Data.Message != "Message was aborted" {
		t.Errorf("Expected Data.Message to be %q, got %q", "Message was aborted", err.Data.Message)
	}
}

func TestProviderAuthErrorUnmarshal(t *testing.T) {
	jsonData := `{
		"name": "ProviderAuthError",
		"data": {
			"message": "Authentication failed",
			"providerID": "provider123"
		}
	}`

	var err ProviderAuthError
	if unmarshalErr := json.Unmarshal([]byte(jsonData), &err); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	if err.Name != ProviderAuthErrorNameProviderAuthError {
		t.Errorf("Expected Name to be %q, got %q", ProviderAuthErrorNameProviderAuthError, err.Name)
	}

	if err.Data.Message != "Authentication failed" {
		t.Errorf("Expected Data.Message to be %q, got %q", "Authentication failed", err.Data.Message)
	}

	if err.Data.ProviderID != "provider123" {
		t.Errorf("Expected Data.ProviderID to be %q, got %q", "provider123", err.Data.ProviderID)
	}
}

func TestUnknownErrorUnmarshal(t *testing.T) {
	jsonData := `{
		"name": "UnknownError",
		"data": {
			"message": "An unknown error occurred"
		}
	}`

	var err UnknownError
	if unmarshalErr := json.Unmarshal([]byte(jsonData), &err); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	if err.Name != UnknownErrorNameUnknownError {
		t.Errorf("Expected Name to be %q, got %q", UnknownErrorNameUnknownError, err.Name)
	}

	if err.Data.Message != "An unknown error occurred" {
		t.Errorf("Expected Data.Message to be %q, got %q", "An unknown error occurred", err.Data.Message)
	}
}

func TestMessageAbortedErrorNameIsKnown(t *testing.T) {
	if !MessageAbortedErrorNameMessageAbortedError.IsKnown() {
		t.Error("Expected MessageAbortedErrorNameMessageAbortedError to be known")
	}

	unknownName := MessageAbortedErrorName("InvalidName")
	if unknownName.IsKnown() {
		t.Error("Expected unknown name to return false for IsKnown()")
	}
}

func TestProviderAuthErrorNameIsKnown(t *testing.T) {
	if !ProviderAuthErrorNameProviderAuthError.IsKnown() {
		t.Error("Expected ProviderAuthErrorNameProviderAuthError to be known")
	}

	unknownName := ProviderAuthErrorName("InvalidName")
	if unknownName.IsKnown() {
		t.Error("Expected unknown name to return false for IsKnown()")
	}
}

func TestUnknownErrorNameIsKnown(t *testing.T) {
	if !UnknownErrorNameUnknownError.IsKnown() {
		t.Error("Expected UnknownErrorNameUnknownError to be known")
	}

	unknownName := UnknownErrorName("InvalidName")
	if unknownName.IsKnown() {
		t.Error("Expected unknown name to return false for IsKnown()")
	}
}
