package opencode

import (
	"encoding/json"
	"testing"
)

func TestAuth_AsOAuth_Valid(t *testing.T) {
	jsonData := `{
		"type": "oauth",
		"refresh": "refresh_token_123",
		"access": "access_token_456",
		"expires": 1234567890
	}`

	var auth Auth
	if err := json.Unmarshal([]byte(jsonData), &auth); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if auth.Type != AuthTypeOAuth {
		t.Errorf("Expected type %s, got %s", AuthTypeOAuth, auth.Type)
	}

	oauth, err := auth.AsOAuth()
	if err != nil {
		t.Fatalf("AsOAuth() error: %v", err)
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
		t.Errorf("Expected Expires 1234567890, got %v", oauth.Expires)
	}
}

func TestAuth_AsAPI_Valid(t *testing.T) {
	jsonData := `{
		"type": "api",
		"key": "api_key_secret_789"
	}`

	var auth Auth
	if err := json.Unmarshal([]byte(jsonData), &auth); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if auth.Type != AuthTypeAPI {
		t.Errorf("Expected type %s, got %s", AuthTypeAPI, auth.Type)
	}

	apiAuth, err := auth.AsAPI()
	if err != nil {
		t.Fatalf("AsAPI() error: %v", err)
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

func TestAuth_AsWellKnown_Valid(t *testing.T) {
	jsonData := `{
		"type": "wellknown",
		"key": "wellknown_key_abc",
		"token": "wellknown_token_xyz"
	}`

	var auth Auth
	if err := json.Unmarshal([]byte(jsonData), &auth); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if auth.Type != AuthTypeWellKnown {
		t.Errorf("Expected type %s, got %s", AuthTypeWellKnown, auth.Type)
	}

	wellKnown, err := auth.AsWellKnown()
	if err != nil {
		t.Fatalf("AsWellKnown() error: %v", err)
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

func TestAuth_AsOAuth_WrongType(t *testing.T) {
	jsonData := `{"type": "api", "key": "api_key_secret"}`

	var auth Auth
	if err := json.Unmarshal([]byte(jsonData), &auth); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	oauth, err := auth.AsOAuth()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if oauth != nil {
		t.Error("AsOAuth() should return nil for api type")
	}
}

func TestAuth_AsAPI_WrongType(t *testing.T) {
	jsonData := `{"type": "oauth", "refresh": "r", "access": "a", "expires": 1234567890}`

	var auth Auth
	if err := json.Unmarshal([]byte(jsonData), &auth); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	apiAuth, err := auth.AsAPI()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if apiAuth != nil {
		t.Error("AsAPI() should return nil for oauth type")
	}
}

func TestAuth_AsWellKnown_WrongType(t *testing.T) {
	jsonData := `{"type": "api", "key": "api_key_secret"}`

	var auth Auth
	if err := json.Unmarshal([]byte(jsonData), &auth); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	wellKnown, err := auth.AsWellKnown()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if wellKnown != nil {
		t.Error("AsWellKnown() should return nil for api type")
	}
}

func TestAuth_InvalidJSON(t *testing.T) {
	jsonData := `{invalid json`

	var auth Auth
	if err := json.Unmarshal([]byte(jsonData), &auth); err == nil {
		t.Error("Unmarshal should fail for invalid JSON")
	}
}

func TestAuth_MissingType(t *testing.T) {
	jsonData := `{"key": "some_key"}`

	var auth Auth
	if err := json.Unmarshal([]byte(jsonData), &auth); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if auth.Type != "" {
		t.Errorf("Expected empty type, got %s", auth.Type)
	}

	if oauth, err := auth.AsOAuth(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if oauth != nil {
		t.Error("AsOAuth() should return nil for missing type")
	}
	if apiAuth, err := auth.AsAPI(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if apiAuth != nil {
		t.Error("AsAPI() should return nil for missing type")
	}
	if wellKnown, err := auth.AsWellKnown(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if wellKnown != nil {
		t.Error("AsWellKnown() should return nil for missing type")
	}
}

func TestAuth_UnknownType(t *testing.T) {
	jsonData := `{"type": "unknown_auth_type", "key": "some_key"}`

	var auth Auth
	if err := json.Unmarshal([]byte(jsonData), &auth); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if auth.Type.IsKnown() {
		t.Error("IsKnown() should return false for unknown type")
	}

	if oauth, err := auth.AsOAuth(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if oauth != nil {
		t.Error("AsOAuth() should return nil for unknown type")
	}
	if apiAuth, err := auth.AsAPI(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if apiAuth != nil {
		t.Error("AsAPI() should return nil for unknown type")
	}
	if wellKnown, err := auth.AsWellKnown(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if wellKnown != nil {
		t.Error("AsWellKnown() should return nil for unknown type")
	}
}

func TestAuth_MalformedOAuth(t *testing.T) {
	jsonData := `{"type": "oauth", "refresh": "refresh_token"}`

	var auth Auth
	if err := json.Unmarshal([]byte(jsonData), &auth); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	oauth, err := auth.AsOAuth()
	if err != nil {
		t.Fatalf("AsOAuth() error: %v", err)
	}
	if oauth == nil {
		t.Fatal("AsOAuth() returned nil")
	}
	if oauth.Refresh != "refresh_token" {
		t.Errorf("Expected Refresh 'refresh_token', got %s", oauth.Refresh)
	}
	if oauth.Access != "" {
		t.Errorf("Expected Access empty string, got %s", oauth.Access)
	}
}

func TestAuth_EmptyJSON(t *testing.T) {
	jsonData := `{}`

	var auth Auth
	if err := json.Unmarshal([]byte(jsonData), &auth); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if auth.Type != "" {
		t.Errorf("Expected empty type, got %s", auth.Type)
	}

	if oauth, err := auth.AsOAuth(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if oauth != nil {
		t.Error("AsOAuth() should return nil for empty JSON")
	}
}

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
