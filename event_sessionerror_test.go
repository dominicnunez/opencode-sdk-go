package opencode

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestSessionError_AsProviderAuth_ValidProviderAuthError(t *testing.T) {
	jsonData := `{"name":"ProviderAuthError","data":{"message":"auth failed","providerID":"test-provider"}}`
	var sessErr SessionError
	if e := json.Unmarshal([]byte(jsonData), &sessErr); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}
	if sessErr.Name != SessionErrorNameProviderAuthError {
		t.Errorf("Expected name ProviderAuthError, got %s", sessErr.Name)
	}

	providerAuthErr, err := sessErr.AsProviderAuth()
	if err != nil {
		t.Fatalf("AsProviderAuth error: %v", err)
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
	var sessErr SessionError
	if e := json.Unmarshal([]byte(jsonData), &sessErr); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}
	if sessErr.Name != SessionErrorNameUnknownError {
		t.Errorf("Expected name UnknownError, got %s", sessErr.Name)
	}

	unknownErr, err := sessErr.AsUnknown()
	if err != nil {
		t.Fatalf("AsUnknown error: %v", err)
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
	var sessErr SessionError
	if e := json.Unmarshal([]byte(jsonData), &sessErr); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}
	if sessErr.Name != SessionErrorNameMessageOutputLengthError {
		t.Errorf("Expected name MessageOutputLengthError, got %s", sessErr.Name)
	}

	outputLengthErr, err := sessErr.AsOutputLength()
	if err != nil {
		t.Fatalf("AsOutputLength error: %v", err)
	}
	if outputLengthErr == nil {
		t.Fatal("Expected non-nil MessageOutputLengthError")
	}
	if outputLengthErr.Data == nil {
		t.Error("Expected non-nil Data field")
	}
}

func TestSessionError_AsAborted_ValidMessageAbortedError(t *testing.T) {
	jsonData := `{"name":"MessageAbortedError","data":{"message":"user cancelled"}}`
	var sessErr SessionError
	if e := json.Unmarshal([]byte(jsonData), &sessErr); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}
	if sessErr.Name != SessionErrorNameMessageAbortedError {
		t.Errorf("Expected name MessageAbortedError, got %s", sessErr.Name)
	}

	abortedErr, err := sessErr.AsAborted()
	if err != nil {
		t.Fatalf("AsAborted error: %v", err)
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
	var sessErr SessionError
	if e := json.Unmarshal([]byte(jsonData), &sessErr); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}
	if sessErr.Name != SessionErrorNameAPIError {
		t.Errorf("Expected name APIError, got %s", sessErr.Name)
	}

	apiErr, err := sessErr.AsAPI()
	if err != nil {
		t.Fatalf("AsAPI error: %v", err)
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

func TestSessionError_WrongType_ReturnsErrWrongVariant(t *testing.T) {
	jsonData := `{"name":"ProviderAuthError","data":{"message":"auth failed","providerID":"test"}}`
	var sessErr SessionError
	if e := json.Unmarshal([]byte(jsonData), &sessErr); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}

	if v, err := sessErr.AsUnknown(); !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	} else if v != nil {
		t.Error("Expected AsUnknown to return nil for ProviderAuthError")
	}
	if v, err := sessErr.AsOutputLength(); !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	} else if v != nil {
		t.Error("Expected AsOutputLength to return nil for ProviderAuthError")
	}
	if v, err := sessErr.AsAborted(); !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	} else if v != nil {
		t.Error("Expected AsAborted to return nil for ProviderAuthError")
	}
	if v, err := sessErr.AsAPI(); !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	} else if v != nil {
		t.Error("Expected AsAPI to return nil for ProviderAuthError")
	}
}

func TestSessionError_InvalidJSON(t *testing.T) {
	invalidJSON := `{"name":"ProviderAuthError","data":{malformed}}`
	var sessErr SessionError
	if e := json.Unmarshal([]byte(invalidJSON), &sessErr); e == nil {
		t.Fatal("Expected Unmarshal to fail on malformed JSON")
	}
}

func TestSessionError_MalformedData_AsMethodFails(t *testing.T) {
	jsonData := `{"name":"ProviderAuthError","data":"not an object"}`
	var sessErr SessionError
	if e := json.Unmarshal([]byte(jsonData), &sessErr); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}

	_, err := sessErr.AsProviderAuth()
	if err == nil {
		t.Error("Expected AsProviderAuth to return error for wrong data structure")
	}
}

func TestSessionError_MissingName(t *testing.T) {
	jsonData := `{"data":{"message":"test"}}`
	var sessErr SessionError
	if e := json.Unmarshal([]byte(jsonData), &sessErr); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}

	if sessErr.Name != "" {
		t.Errorf("Expected empty name, got %s", sessErr.Name)
	}

	if v, err := sessErr.AsProviderAuth(); !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	} else if v != nil {
		t.Error("Expected AsProviderAuth to return nil for empty name")
	}
}

func TestSessionError_UnknownName(t *testing.T) {
	jsonData := `{"name":"UnknownErrorType","data":{"message":"test"}}`
	var sessErr SessionError
	if e := json.Unmarshal([]byte(jsonData), &sessErr); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}

	if sessErr.Name != "UnknownErrorType" {
		t.Errorf("Expected name UnknownErrorType, got %s", sessErr.Name)
	}
	if sessErr.Name.IsKnown() {
		t.Error("Expected IsKnown to return false for unknown error name")
	}

	if v, err := sessErr.AsProviderAuth(); !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	} else if v != nil {
		t.Error("Expected AsProviderAuth to return nil for unknown name")
	}
}

func TestSessionError_EmptyJSON(t *testing.T) {
	jsonData := `{}`
	var sessErr SessionError
	if e := json.Unmarshal([]byte(jsonData), &sessErr); e != nil {
		t.Fatalf("Unmarshal failed: %v", e)
	}

	if sessErr.Name != "" {
		t.Errorf("Expected empty name, got %s", sessErr.Name)
	}

	if v, err := sessErr.AsUnknown(); !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	} else if v != nil {
		t.Error("Expected AsUnknown to return nil for empty JSON")
	}
}
