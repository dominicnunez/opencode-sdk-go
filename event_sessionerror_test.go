package opencode

import (
	"encoding/json"
	"testing"
)

func TestSessionError_AsProviderAuth_ValidProviderAuthError(t *testing.T) {
	jsonData := `{"name":"ProviderAuthError","data":{"message":"auth failed","providerID":"test-provider"}}`
	var err SessionError
	if e := json.Unmarshal([]byte(jsonData), &err); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}
	if err.Name != SessionErrorNameProviderAuthError {
		t.Errorf("Expected name ProviderAuthError, got %s", err.Name)
	}

	providerAuthErr, ok := err.AsProviderAuth()
	if !ok {
		t.Fatal("Expected AsProviderAuth to return true")
	}
	if providerAuthErr == nil {
		t.Fatal("Expected non-nil ProviderAuthError")
	}
	if providerAuthErr.Data.Message != "auth failed" {
		t.Errorf("Expected message 'auth failed', got %s", providerAuthErr.Data.Message)
	}
	if providerAuthErr.Data.ProviderID != "test-provider" {
		t.Errorf("Expected providerID 'test-provider', got %s", providerAuthErr.Data.ProviderID)
	}
}

func TestSessionError_AsUnknown_ValidUnknownError(t *testing.T) {
	jsonData := `{"name":"UnknownError","data":{"message":"something went wrong"}}`
	var err SessionError
	if e := json.Unmarshal([]byte(jsonData), &err); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}
	if err.Name != SessionErrorNameUnknownError {
		t.Errorf("Expected name UnknownError, got %s", err.Name)
	}

	unknownErr, ok := err.AsUnknown()
	if !ok {
		t.Fatal("Expected AsUnknown to return true")
	}
	if unknownErr == nil {
		t.Fatal("Expected non-nil UnknownError")
	}
	if unknownErr.Data.Message != "something went wrong" {
		t.Errorf("Expected message 'something went wrong', got %s", unknownErr.Data.Message)
	}
}

func TestSessionError_AsOutputLength_ValidMessageOutputLengthError(t *testing.T) {
	jsonData := `{"name":"MessageOutputLengthError","data":{"limit":1000,"used":1500}}`
	var err SessionError
	if e := json.Unmarshal([]byte(jsonData), &err); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}
	if err.Name != SessionErrorNameMessageOutputLengthError {
		t.Errorf("Expected name MessageOutputLengthError, got %s", err.Name)
	}

	outputLengthErr, ok := err.AsOutputLength()
	if !ok {
		t.Fatal("Expected AsOutputLength to return true")
	}
	if outputLengthErr == nil {
		t.Fatal("Expected non-nil MessageOutputLengthError")
	}
	// Data is interface{} so we can check it's not nil
	if outputLengthErr.Data == nil {
		t.Error("Expected non-nil Data field")
	}
}

func TestSessionError_AsAborted_ValidMessageAbortedError(t *testing.T) {
	jsonData := `{"name":"MessageAbortedError","data":{"message":"user cancelled"}}`
	var err SessionError
	if e := json.Unmarshal([]byte(jsonData), &err); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}
	if err.Name != SessionErrorNameMessageAbortedError {
		t.Errorf("Expected name MessageAbortedError, got %s", err.Name)
	}

	abortedErr, ok := err.AsAborted()
	if !ok {
		t.Fatal("Expected AsAborted to return true")
	}
	if abortedErr == nil {
		t.Fatal("Expected non-nil MessageAbortedError")
	}
	if abortedErr.Data.Message != "user cancelled" {
		t.Errorf("Expected message 'user cancelled', got %s", abortedErr.Data.Message)
	}
}

func TestSessionError_AsAPI_ValidSessionAPIError(t *testing.T) {
	jsonData := `{"name":"APIError","data":{"isRetryable":true,"message":"timeout","statusCode":504,"responseBody":"Gateway Timeout","responseHeaders":{"content-type":"text/plain"}}}`
	var err SessionError
	if e := json.Unmarshal([]byte(jsonData), &err); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}
	if err.Name != SessionErrorNameAPIError {
		t.Errorf("Expected name APIError, got %s", err.Name)
	}

	apiErr, ok := err.AsAPI()
	if !ok {
		t.Fatal("Expected AsAPI to return true")
	}
	if apiErr == nil {
		t.Fatal("Expected non-nil SessionAPIError")
	}
	if !apiErr.Data.IsRetryable {
		t.Error("Expected IsRetryable to be true")
	}
	if apiErr.Data.Message != "timeout" {
		t.Errorf("Expected message 'timeout', got %s", apiErr.Data.Message)
	}
	if apiErr.Data.StatusCode == nil || *apiErr.Data.StatusCode != 504 {
		t.Errorf("Expected statusCode 504, got %v", apiErr.Data.StatusCode)
	}
	if apiErr.Data.ResponseBody == nil || *apiErr.Data.ResponseBody != "Gateway Timeout" {
		t.Errorf("Expected responseBody 'Gateway Timeout', got %v", apiErr.Data.ResponseBody)
	}
	if apiErr.Data.ResponseHeaders["content-type"] != "text/plain" {
		t.Errorf("Expected header content-type=text/plain, got %s", apiErr.Data.ResponseHeaders["content-type"])
	}
}

func TestSessionError_WrongType_ReturnsNilFalse(t *testing.T) {
	jsonData := `{"name":"ProviderAuthError","data":{"message":"auth failed","providerID":"test"}}`
	var err SessionError
	if e := json.Unmarshal([]byte(jsonData), &err); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}

	// Try wrong As* methods
	if unknownErr, ok := err.AsUnknown(); ok || unknownErr != nil {
		t.Error("Expected AsUnknown to return (nil, false) for ProviderAuthError")
	}
	if outputLengthErr, ok := err.AsOutputLength(); ok || outputLengthErr != nil {
		t.Error("Expected AsOutputLength to return (nil, false) for ProviderAuthError")
	}
	if abortedErr, ok := err.AsAborted(); ok || abortedErr != nil {
		t.Error("Expected AsAborted to return (nil, false) for ProviderAuthError")
	}
	if apiErr, ok := err.AsAPI(); ok || apiErr != nil {
		t.Error("Expected AsAPI to return (nil, false) for ProviderAuthError")
	}
}

func TestSessionError_InvalidJSON(t *testing.T) {
	invalidJSON := `{"name":"ProviderAuthError","data":{malformed}}`
	var err SessionError
	// Unmarshal should fail on malformed JSON
	if e := json.Unmarshal([]byte(invalidJSON), &err); e == nil {
		t.Fatal("Expected Unmarshal to fail on malformed JSON")
	}
}

func TestSessionError_MalformedData_AsMethodFails(t *testing.T) {
	// Valid JSON but with wrong data structure
	jsonData := `{"name":"ProviderAuthError","data":"not an object"}`
	var err SessionError
	if e := json.Unmarshal([]byte(jsonData), &err); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}

	// AsProviderAuth should fail when trying to unmarshal the full data
	if providerAuthErr, ok := err.AsProviderAuth(); ok || providerAuthErr != nil {
		t.Error("Expected AsProviderAuth to fail on wrong data structure")
	}
}

func TestSessionError_MissingName(t *testing.T) {
	jsonData := `{"data":{"message":"test"}}`
	var err SessionError
	if e := json.Unmarshal([]byte(jsonData), &err); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}

	// Name should be empty
	if err.Name != "" {
		t.Errorf("Expected empty name, got %s", err.Name)
	}

	// All As* methods should return (nil, false)
	if _, ok := err.AsProviderAuth(); ok {
		t.Error("Expected AsProviderAuth to return false for empty name")
	}
}

func TestSessionError_UnknownName(t *testing.T) {
	jsonData := `{"name":"UnknownErrorType","data":{"message":"test"}}`
	var err SessionError
	if e := json.Unmarshal([]byte(jsonData), &err); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}

	// Name should be set even if unknown
	if err.Name != "UnknownErrorType" {
		t.Errorf("Expected name UnknownErrorType, got %s", err.Name)
	}

	// IsKnown should return false
	if err.Name.IsKnown() {
		t.Error("Expected IsKnown to return false for unknown error name")
	}

	// All As* methods should return (nil, false)
	if _, ok := err.AsProviderAuth(); ok {
		t.Error("Expected AsProviderAuth to return false for unknown name")
	}
}

func TestSessionError_EmptyJSON(t *testing.T) {
	jsonData := `{}`
	var err SessionError
	if e := json.Unmarshal([]byte(jsonData), &err); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}

	// Name should be empty
	if err.Name != "" {
		t.Errorf("Expected empty name, got %s", err.Name)
	}

	// All As* methods should return (nil, false)
	if _, ok := err.AsUnknown(); ok {
		t.Error("Expected AsUnknown to return false for empty JSON")
	}
}
