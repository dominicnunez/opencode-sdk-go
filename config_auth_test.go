package opencode

import (
	"encoding/json"
	"testing"
)

// Test OAuth variant with all required fields
func TestAuth_AsOAuth_Valid(t *testing.T) {
	jsonData := `{
		"type": "oauth",
		"refresh": "refresh_token_123",
		"access": "access_token_456",
		"expires": 1234567890
	}`

	var auth Auth
	err := json.Unmarshal([]byte(jsonData), &auth)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if auth.Type != AuthTypeOAuth {
		t.Errorf("Expected type %s, got %s", AuthTypeOAuth, auth.Type)
	}

	oauth, ok := auth.AsOAuth()
	if !ok {
		t.Fatal("AsOAuth() should return true for oauth type")
	}
	if oauth == nil {
		t.Fatal("AsOAuth() returned nil")
	}

	if oauth.Type != AuthTypeOAuth {
		t.Errorf("Expected OAuth.Type %s, got %s", AuthTypeOAuth, oauth.Type)
	}
	if oauth.Refresh != "refresh_token_123" {
		t.Errorf("Expected Refresh 'refresh_token_123', got %s", oauth.Refresh)
	}
	if oauth.Access != "access_token_456" {
		t.Errorf("Expected Access 'access_token_456', got %s", oauth.Access)
	}
	if oauth.Expires != 1234567890 {
		t.Errorf("Expected Expires 1234567890, got %d", oauth.Expires)
	}
}

// Test ApiAuth variant with all required fields
func TestAuth_AsAPI_Valid(t *testing.T) {
	jsonData := `{
		"type": "api",
		"key": "api_key_secret_789"
	}`

	var auth Auth
	err := json.Unmarshal([]byte(jsonData), &auth)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if auth.Type != AuthTypeAPI {
		t.Errorf("Expected type %s, got %s", AuthTypeAPI, auth.Type)
	}

	apiAuth, ok := auth.AsAPI()
	if !ok {
		t.Fatal("AsAPI() should return true for api type")
	}
	if apiAuth == nil {
		t.Fatal("AsAPI() returned nil")
	}

	if apiAuth.Type != AuthTypeAPI {
		t.Errorf("Expected ApiAuth.Type %s, got %s", AuthTypeAPI, apiAuth.Type)
	}
	if apiAuth.Key != "api_key_secret_789" {
		t.Errorf("Expected Key 'api_key_secret_789', got %s", apiAuth.Key)
	}
}

// Test WellKnownAuth variant with all required fields
func TestAuth_AsWellKnown_Valid(t *testing.T) {
	jsonData := `{
		"type": "wellknown",
		"key": "wellknown_key_abc",
		"token": "wellknown_token_xyz"
	}`

	var auth Auth
	err := json.Unmarshal([]byte(jsonData), &auth)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if auth.Type != AuthTypeWellKnown {
		t.Errorf("Expected type %s, got %s", AuthTypeWellKnown, auth.Type)
	}

	wellKnown, ok := auth.AsWellKnown()
	if !ok {
		t.Fatal("AsWellKnown() should return true for wellknown type")
	}
	if wellKnown == nil {
		t.Fatal("AsWellKnown() returned nil")
	}

	if wellKnown.Type != AuthTypeWellKnown {
		t.Errorf("Expected WellKnownAuth.Type %s, got %s", AuthTypeWellKnown, wellKnown.Type)
	}
	if wellKnown.Key != "wellknown_key_abc" {
		t.Errorf("Expected Key 'wellknown_key_abc', got %s", wellKnown.Key)
	}
	if wellKnown.Token != "wellknown_token_xyz" {
		t.Errorf("Expected Token 'wellknown_token_xyz', got %s", wellKnown.Token)
	}
}

// Test wrong type returns (nil, false) for OAuth when type is api
func TestAuth_AsOAuth_WrongType(t *testing.T) {
	jsonData := `{
		"type": "api",
		"key": "api_key_secret"
	}`

	var auth Auth
	err := json.Unmarshal([]byte(jsonData), &auth)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	oauth, ok := auth.AsOAuth()
	if ok {
		t.Error("AsOAuth() should return false for api type")
	}
	if oauth != nil {
		t.Error("AsOAuth() should return nil for wrong type")
	}
}

// Test wrong type returns (nil, false) for ApiAuth when type is oauth
func TestAuth_AsAPI_WrongType(t *testing.T) {
	jsonData := `{
		"type": "oauth",
		"refresh": "refresh_token",
		"access": "access_token",
		"expires": 1234567890
	}`

	var auth Auth
	err := json.Unmarshal([]byte(jsonData), &auth)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	apiAuth, ok := auth.AsAPI()
	if ok {
		t.Error("AsAPI() should return false for oauth type")
	}
	if apiAuth != nil {
		t.Error("AsAPI() should return nil for wrong type")
	}
}

// Test wrong type returns (nil, false) for WellKnownAuth when type is api
func TestAuth_AsWellKnown_WrongType(t *testing.T) {
	jsonData := `{
		"type": "api",
		"key": "api_key_secret"
	}`

	var auth Auth
	err := json.Unmarshal([]byte(jsonData), &auth)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	wellKnown, ok := auth.AsWellKnown()
	if ok {
		t.Error("AsWellKnown() should return false for api type")
	}
	if wellKnown != nil {
		t.Error("AsWellKnown() should return nil for wrong type")
	}
}

// Test invalid JSON returns error
func TestAuth_InvalidJSON(t *testing.T) {
	jsonData := `{invalid json`

	var auth Auth
	err := json.Unmarshal([]byte(jsonData), &auth)
	if err == nil {
		t.Error("Unmarshal should fail for invalid JSON")
	}
}

// Test missing type field
func TestAuth_MissingType(t *testing.T) {
	jsonData := `{
		"key": "some_key"
	}`

	var auth Auth
	err := json.Unmarshal([]byte(jsonData), &auth)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Type should be empty string
	if auth.Type != "" {
		t.Errorf("Expected empty type, got %s", auth.Type)
	}

	// All As* methods should return (nil, false)
	if oauth, ok := auth.AsOAuth(); ok || oauth != nil {
		t.Error("AsOAuth() should return (nil, false) for missing type")
	}
	if apiAuth, ok := auth.AsAPI(); ok || apiAuth != nil {
		t.Error("AsAPI() should return (nil, false) for missing type")
	}
	if wellKnown, ok := auth.AsWellKnown(); ok || wellKnown != nil {
		t.Error("AsWellKnown() should return (nil, false) for missing type")
	}
}

// Test unknown type
func TestAuth_UnknownType(t *testing.T) {
	jsonData := `{
		"type": "unknown_auth_type",
		"key": "some_key"
	}`

	var auth Auth
	err := json.Unmarshal([]byte(jsonData), &auth)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if auth.Type.IsKnown() {
		t.Error("IsKnown() should return false for unknown type")
	}

	// All As* methods should return (nil, false)
	if oauth, ok := auth.AsOAuth(); ok || oauth != nil {
		t.Error("AsOAuth() should return (nil, false) for unknown type")
	}
	if apiAuth, ok := auth.AsAPI(); ok || apiAuth != nil {
		t.Error("AsAPI() should return (nil, false) for unknown type")
	}
	if wellKnown, ok := auth.AsWellKnown(); ok || wellKnown != nil {
		t.Error("AsWellKnown() should return (nil, false) for unknown type")
	}
}

// Test malformed OAuth data (missing required field)
func TestAuth_MalformedOAuth(t *testing.T) {
	jsonData := `{
		"type": "oauth",
		"refresh": "refresh_token"
	}`

	var auth Auth
	err := json.Unmarshal([]byte(jsonData), &auth)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if auth.Type != AuthTypeOAuth {
		t.Errorf("Expected type %s, got %s", AuthTypeOAuth, auth.Type)
	}

	oauth, ok := auth.AsOAuth()
	if !ok {
		t.Fatal("AsOAuth() should return true for oauth type even if fields are missing")
	}
	if oauth == nil {
		t.Fatal("AsOAuth() returned nil")
	}

	// Check that missing fields are zero values
	if oauth.Refresh != "refresh_token" {
		t.Errorf("Expected Refresh 'refresh_token', got %s", oauth.Refresh)
	}
	if oauth.Access != "" {
		t.Errorf("Expected Access empty string, got %s", oauth.Access)
	}
	if oauth.Expires != 0 {
		t.Errorf("Expected Expires 0, got %d", oauth.Expires)
	}
}

// Test empty JSON
func TestAuth_EmptyJSON(t *testing.T) {
	jsonData := `{}`

	var auth Auth
	err := json.Unmarshal([]byte(jsonData), &auth)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Type should be empty
	if auth.Type != "" {
		t.Errorf("Expected empty type, got %s", auth.Type)
	}

	// All As* methods should return (nil, false)
	if oauth, ok := auth.AsOAuth(); ok || oauth != nil {
		t.Error("AsOAuth() should return (nil, false) for empty JSON")
	}
	if apiAuth, ok := auth.AsAPI(); ok || apiAuth != nil {
		t.Error("AsAPI() should return (nil, false) for empty JSON")
	}
	if wellKnown, ok := auth.AsWellKnown(); ok || wellKnown != nil {
		t.Error("AsWellKnown() should return (nil, false) for empty JSON")
	}
}

// Test AuthType.IsKnown() for all known types
func TestAuthType_IsKnown(t *testing.T) {
	tests := []struct {
		name     string
		authType AuthType
		expected bool
	}{
		{"oauth is known", AuthTypeOAuth, true},
		{"api is known", AuthTypeAPI, true},
		{"wellknown is known", AuthTypeWellKnown, true},
		{"unknown type", AuthType("unknown"), false},
		{"empty type", AuthType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.authType.IsKnown()
			if result != tt.expected {
				t.Errorf("IsKnown() = %v, want %v", result, tt.expected)
			}
		})
	}
}
