package opencode

import (
	"encoding/json"
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

	var err AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &err); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	if err.Name != AssistantMessageErrorNameProviderAuthError {
		t.Errorf("Expected Name to be ProviderAuthError, got %s", err.Name)
	}

	providerAuthErr, ok := err.AsProviderAuth()
	if !ok {
		t.Fatal("AsProviderAuth() should return true for ProviderAuthError")
	}
	if providerAuthErr == nil {
		t.Fatal("AsProviderAuth() should return non-nil error")
	}
	if providerAuthErr.Data.Message != "Authentication failed" {
		t.Errorf("Expected message 'Authentication failed', got '%s'", providerAuthErr.Data.Message)
	}
	if providerAuthErr.Data.ProviderID != "provider-123" {
		t.Errorf("Expected providerID 'provider-123', got '%s'", providerAuthErr.Data.ProviderID)
	}

	// Wrong type should return false
	if _, ok := err.AsUnknown(); ok {
		t.Error("AsUnknown() should return false for ProviderAuthError")
	}
}

func TestAssistantMessageError_AsUnknown_ValidUnknownError(t *testing.T) {
	jsonData := []byte(`{
		"name": "UnknownError",
		"data": {
			"message": "Something went wrong"
		}
	}`)

	var err AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &err); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	if err.Name != AssistantMessageErrorNameUnknownError {
		t.Errorf("Expected Name to be UnknownError, got %s", err.Name)
	}

	unknownErr, ok := err.AsUnknown()
	if !ok {
		t.Fatal("AsUnknown() should return true for UnknownError")
	}
	if unknownErr == nil {
		t.Fatal("AsUnknown() should return non-nil error")
	}
	if unknownErr.Data.Message != "Something went wrong" {
		t.Errorf("Expected message 'Something went wrong', got '%s'", unknownErr.Data.Message)
	}

	// Wrong type should return false
	if _, ok := err.AsAborted(); ok {
		t.Error("AsAborted() should return false for UnknownError")
	}
}

func TestAssistantMessageError_AsOutputLength_ValidOutputLengthError(t *testing.T) {
	jsonData := []byte(`{
		"name": "MessageOutputLengthError",
		"data": {}
	}`)

	var err AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &err); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	if err.Name != AssistantMessageErrorNameMessageOutputLengthError {
		t.Errorf("Expected Name to be MessageOutputLengthError, got %s", err.Name)
	}

	outputLengthErr, ok := err.AsOutputLength()
	if !ok {
		t.Fatal("AsOutputLength() should return true for MessageOutputLengthError")
	}
	if outputLengthErr == nil {
		t.Fatal("AsOutputLength() should return non-nil error")
	}

	// Wrong type should return false
	if _, ok := err.AsAPI(); ok {
		t.Error("AsAPI() should return false for MessageOutputLengthError")
	}
}

func TestAssistantMessageError_AsAborted_ValidAbortedError(t *testing.T) {
	jsonData := []byte(`{
		"name": "MessageAbortedError",
		"data": {
			"message": "Request aborted by user"
		}
	}`)

	var err AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &err); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	if err.Name != AssistantMessageErrorNameMessageAbortedError {
		t.Errorf("Expected Name to be MessageAbortedError, got %s", err.Name)
	}

	abortedErr, ok := err.AsAborted()
	if !ok {
		t.Fatal("AsAborted() should return true for MessageAbortedError")
	}
	if abortedErr == nil {
		t.Fatal("AsAborted() should return non-nil error")
	}
	if abortedErr.Data.Message != "Request aborted by user" {
		t.Errorf("Expected message 'Request aborted by user', got '%s'", abortedErr.Data.Message)
	}

	// Wrong type should return false
	if _, ok := err.AsProviderAuth(); ok {
		t.Error("AsProviderAuth() should return false for MessageAbortedError")
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

	var err AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &err); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	if err.Name != AssistantMessageErrorNameAPIError {
		t.Errorf("Expected Name to be APIError, got %s", err.Name)
	}

	apiErr, ok := err.AsAPI()
	if !ok {
		t.Fatal("AsAPI() should return true for APIError")
	}
	if apiErr == nil {
		t.Fatal("AsAPI() should return non-nil error")
	}
	if !apiErr.Data.IsRetryable {
		t.Error("Expected IsRetryable to be true")
	}
	if apiErr.Data.Message != "Rate limit exceeded" {
		t.Errorf("Expected message 'Rate limit exceeded', got '%s'", apiErr.Data.Message)
	}
	if apiErr.Data.StatusCode != 429 {
		t.Errorf("Expected StatusCode 429, got %f", apiErr.Data.StatusCode)
	}
	if apiErr.Data.ResponseBody != "Too many requests" {
		t.Errorf("Expected ResponseBody 'Too many requests', got '%s'", apiErr.Data.ResponseBody)
	}
	if apiErr.Data.ResponseHeaders["retry-after"] != "60" {
		t.Errorf("Expected retry-after header '60', got '%s'", apiErr.Data.ResponseHeaders["retry-after"])
	}

	// Wrong type should return false
	if _, ok := err.AsOutputLength(); ok {
		t.Error("AsOutputLength() should return false for APIError")
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

	var err AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &err); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	// All other As* methods should return false
	if _, ok := err.AsUnknown(); ok {
		t.Error("AsUnknown() should return false for ProviderAuthError")
	}
	if _, ok := err.AsOutputLength(); ok {
		t.Error("AsOutputLength() should return false for ProviderAuthError")
	}
	if _, ok := err.AsAborted(); ok {
		t.Error("AsAborted() should return false for ProviderAuthError")
	}
	if _, ok := err.AsAPI(); ok {
		t.Error("AsAPI() should return false for ProviderAuthError")
	}
}

func TestAssistantMessageError_InvalidJSON(t *testing.T) {
	jsonData := []byte(`{invalid json}`)

	var err AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &err); unmarshalErr == nil {
		t.Error("Expected Unmarshal to fail for invalid JSON")
	}
}

func TestAssistantMessageError_MissingName(t *testing.T) {
	jsonData := []byte(`{
		"data": {
			"message": "Some error"
		}
	}`)

	var err AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &err); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	// Name should be empty string
	if err.Name != "" {
		t.Errorf("Expected Name to be empty, got %s", err.Name)
	}

	// All As* methods should return false for empty name
	if _, ok := err.AsProviderAuth(); ok {
		t.Error("AsProviderAuth() should return false for missing name")
	}
	if _, ok := err.AsUnknown(); ok {
		t.Error("AsUnknown() should return false for missing name")
	}
	if _, ok := err.AsOutputLength(); ok {
		t.Error("AsOutputLength() should return false for missing name")
	}
	if _, ok := err.AsAborted(); ok {
		t.Error("AsAborted() should return false for missing name")
	}
	if _, ok := err.AsAPI(); ok {
		t.Error("AsAPI() should return false for missing name")
	}
}

func TestAssistantMessageError_UnknownName(t *testing.T) {
	jsonData := []byte(`{
		"name": "UnknownErrorType",
		"data": {}
	}`)

	var err AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &err); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	if err.Name != "UnknownErrorType" {
		t.Errorf("Expected Name to be 'UnknownErrorType', got %s", err.Name)
	}

	// All As* methods should return false for unknown name
	if _, ok := err.AsProviderAuth(); ok {
		t.Error("AsProviderAuth() should return false for unknown name")
	}
	if _, ok := err.AsUnknown(); ok {
		t.Error("AsUnknown() should return false for unknown name")
	}
	if _, ok := err.AsOutputLength(); ok {
		t.Error("AsOutputLength() should return false for unknown name")
	}
	if _, ok := err.AsAborted(); ok {
		t.Error("AsAborted() should return false for unknown name")
	}
	if _, ok := err.AsAPI(); ok {
		t.Error("AsAPI() should return false for unknown name")
	}
}

func TestAssistantMessageError_MalformedData(t *testing.T) {
	// Valid name but data structure doesn't match expected type
	jsonData := []byte(`{
		"name": "ProviderAuthError",
		"data": "this should be an object not a string"
	}`)

	var err AssistantMessageError
	if unmarshalErr := json.Unmarshal(jsonData, &err); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	if err.Name != AssistantMessageErrorNameProviderAuthError {
		t.Errorf("Expected Name to be ProviderAuthError, got %s", err.Name)
	}

	// AsProviderAuth should fail because data structure is invalid
	providerAuthErr, ok := err.AsProviderAuth()
	if ok {
		t.Error("AsProviderAuth() should return false for malformed data")
	}
	if providerAuthErr != nil {
		t.Error("AsProviderAuth() should return nil for malformed data")
	}
}
