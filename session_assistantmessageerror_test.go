package opencode

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestAssistantMessageError_AsProviderAuth_ValidProviderAuthError(t *testing.T) {
	jsonData := []byte(`{
		"name": "ProviderAuthError",
		"data": {
			"message": "Authentication failed",
			"providerID": "provider-123"
		}
	}`)

	var ame AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &ame); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	if ame.Name != AssistantMessageErrorNameProviderAuthError {
		t.Errorf("Expected Name to be ProviderAuthError, got %s", ame.Name)
	}

	providerAuthErr, err := ame.AsProviderAuth()
	if err != nil {
		t.Fatalf("AsProviderAuth() error: %v", err)
	}
	if providerAuthErr == nil {
		t.Fatal("AsProviderAuth() should return non-nil")
	}
	if providerAuthErr.Data.Message != "Authentication failed" {
		t.Errorf("Expected message 'Authentication failed', got '%s'", providerAuthErr.Data.Message)
	}
	if providerAuthErr.Data.ProviderID != "provider-123" {
		t.Errorf("Expected providerID 'provider-123', got '%s'", providerAuthErr.Data.ProviderID)
	}

	// Wrong type should return (nil, ErrWrongVariant)
	v, err := ame.AsUnknown()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if v != nil {
		t.Error("AsUnknown() should return nil for ProviderAuthError")
	}
}

func TestAssistantMessageError_AsUnknown_ValidUnknownError(t *testing.T) {
	jsonData := []byte(`{
		"name": "UnknownError",
		"data": {
			"message": "Something went wrong"
		}
	}`)

	var ame AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &ame); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	if ame.Name != AssistantMessageErrorNameUnknownError {
		t.Errorf("Expected Name to be UnknownError, got %s", ame.Name)
	}

	unknownErr, err := ame.AsUnknown()
	if err != nil {
		t.Fatalf("AsUnknown() error: %v", err)
	}
	if unknownErr == nil {
		t.Fatal("AsUnknown() should return non-nil")
	}
	if unknownErr.Data.Message != "Something went wrong" {
		t.Errorf("Expected message 'Something went wrong', got '%s'", unknownErr.Data.Message)
	}

	// Wrong type should return (nil, ErrWrongVariant)
	v, err := ame.AsAborted()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if v != nil {
		t.Error("AsAborted() should return nil for UnknownError")
	}
}

func TestAssistantMessageError_AsOutputLength_ValidOutputLengthError(t *testing.T) {
	jsonData := []byte(`{
		"name": "MessageOutputLengthError",
		"data": {}
	}`)

	var ame AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &ame); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	if ame.Name != AssistantMessageErrorNameMessageOutputLengthError {
		t.Errorf("Expected Name to be MessageOutputLengthError, got %s", ame.Name)
	}

	outputLengthErr, err := ame.AsOutputLength()
	if err != nil {
		t.Fatalf("AsOutputLength() error: %v", err)
	}
	if outputLengthErr == nil {
		t.Fatal("AsOutputLength() should return non-nil")
	}

	// Wrong type should return (nil, ErrWrongVariant)
	v, err := ame.AsAPI()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if v != nil {
		t.Error("AsAPI() should return nil for MessageOutputLengthError")
	}
}

func TestAssistantMessageError_AsAborted_ValidAbortedError(t *testing.T) {
	jsonData := []byte(`{
		"name": "MessageAbortedError",
		"data": {
			"message": "Request aborted by user"
		}
	}`)

	var ame AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &ame); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	if ame.Name != AssistantMessageErrorNameMessageAbortedError {
		t.Errorf("Expected Name to be MessageAbortedError, got %s", ame.Name)
	}

	abortedErr, err := ame.AsAborted()
	if err != nil {
		t.Fatalf("AsAborted() error: %v", err)
	}
	if abortedErr == nil {
		t.Fatal("AsAborted() should return non-nil")
	}
	if abortedErr.Data.Message != "Request aborted by user" {
		t.Errorf("Expected message 'Request aborted by user', got '%s'", abortedErr.Data.Message)
	}

	// Wrong type should return (nil, ErrWrongVariant)
	v, err := ame.AsProviderAuth()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if v != nil {
		t.Error("AsProviderAuth() should return nil for MessageAbortedError")
	}
}

func TestAssistantMessageError_AsAPI_ValidAPIError(t *testing.T) {
	jsonData := []byte(`{
		"name": "APIError",
		"data": {
			"isRetryable": true,
			"message": "Rate limit exceeded",
			"statusCode": 429,
			"responseBody": "Too many requests",
			"responseHeaders": {
				"retry-after": "60"
			}
		}
	}`)

	var ame AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &ame); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	if ame.Name != AssistantMessageErrorNameAPIError {
		t.Errorf("Expected Name to be APIError, got %s", ame.Name)
	}

	apiErr, err := ame.AsAPI()
	if err != nil {
		t.Fatalf("AsAPI() error: %v", err)
	}
	if apiErr == nil {
		t.Fatal("AsAPI() should return non-nil")
	}
	if !apiErr.Data.IsRetryable {
		t.Error("Expected IsRetryable to be true")
	}
	if apiErr.Data.Message != "Rate limit exceeded" {
		t.Errorf("Expected message 'Rate limit exceeded', got '%s'", apiErr.Data.Message)
	}
	if apiErr.Data.StatusCode != 429 {
		t.Errorf("Expected StatusCode 429, got %v", apiErr.Data.StatusCode)
	}
	if apiErr.Data.ResponseBody != "Too many requests" {
		t.Errorf("Expected ResponseBody 'Too many requests', got '%s'", apiErr.Data.ResponseBody)
	}
	if apiErr.Data.ResponseHeaders["retry-after"] != "60" {
		t.Errorf("Expected retry-after header '60', got '%s'", apiErr.Data.ResponseHeaders["retry-after"])
	}

	// Wrong type should return (nil, ErrWrongVariant)
	v, err := ame.AsOutputLength()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if v != nil {
		t.Error("AsOutputLength() should return nil for APIError")
	}
}

func TestAssistantMessageError_WrongType(t *testing.T) {
	jsonData := []byte(`{
		"name": "ProviderAuthError",
		"data": {
			"message": "Auth failed",
			"providerID": "provider-123"
		}
	}`)

	var ame AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &ame); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	if v, err := ame.AsUnknown(); !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	} else if v != nil {
		t.Error("AsUnknown() should return nil for ProviderAuthError")
	}
	if v, err := ame.AsOutputLength(); !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	} else if v != nil {
		t.Error("AsOutputLength() should return nil for ProviderAuthError")
	}
	if v, err := ame.AsAborted(); !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	} else if v != nil {
		t.Error("AsAborted() should return nil for ProviderAuthError")
	}
	if v, err := ame.AsAPI(); !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	} else if v != nil {
		t.Error("AsAPI() should return nil for ProviderAuthError")
	}
}

func TestAssistantMessageError_InvalidJSON(t *testing.T) {
	jsonData := []byte(`{invalid json}`)

	var ame AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &ame); unmarshalErr == nil {
		t.Error("Expected Unmarshal to fail for invalid JSON")
	}
}

func TestAssistantMessageError_MissingName(t *testing.T) {
	jsonData := []byte(`{
		"data": {
			"message": "Some error"
		}
	}`)

	var ame AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &ame); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	if ame.Name != "" {
		t.Errorf("Expected Name to be empty, got %s", ame.Name)
	}

	if v, err := ame.AsProviderAuth(); !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	} else if v != nil {
		t.Error("AsProviderAuth() should return nil for missing name")
	}
	if v, err := ame.AsUnknown(); !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	} else if v != nil {
		t.Error("AsUnknown() should return nil for missing name")
	}
}

func TestAssistantMessageError_UnknownName(t *testing.T) {
	jsonData := []byte(`{
		"name": "UnknownErrorType",
		"data": {}
	}`)

	var ame AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &ame); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	if ame.Name != "UnknownErrorType" {
		t.Errorf("Expected Name to be 'UnknownErrorType', got %s", ame.Name)
	}

	if v, err := ame.AsProviderAuth(); !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	} else if v != nil {
		t.Error("AsProviderAuth() should return nil for unknown name")
	}
}

func TestAssistantMessageError_MalformedData(t *testing.T) {
	jsonData := []byte(`{
		"name": "ProviderAuthError",
		"data": "this should be an object not a string"
	}`)

	var ame AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &ame); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	if ame.Name != AssistantMessageErrorNameProviderAuthError {
		t.Errorf("Expected Name to be ProviderAuthError, got %s", ame.Name)
	}

	// AsProviderAuth should return error because data structure is invalid
	_, err := ame.AsProviderAuth()
	if err == nil {
		t.Error("AsProviderAuth() should return error for malformed data")
	}
}
