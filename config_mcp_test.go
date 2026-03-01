package opencode

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestConfigMcp_AsLocal_ValidLocalConfig(t *testing.T) {
	jsonData := `{
		"type": "local",
		"command": ["/usr/bin/mcp-server", "--port", "8080"],
		"enabled": true,
		"environment": {
			"NODE_ENV": "production",
			"PORT": "8080"
		}
	}`

	var mcp ConfigMcp
	if err := json.Unmarshal([]byte(jsonData), &mcp); err != nil {
		t.Fatalf("Failed to unmarshal ConfigMcp: %v", err)
	}

	if mcp.Type != ConfigMcpTypeLocal {
		t.Errorf("Expected type %q, got %q", ConfigMcpTypeLocal, mcp.Type)
	}

	local, err := mcp.AsLocal()
	if err != nil {
		t.Fatal("AsLocal() returned false for local type")
	}
	if local == nil {
		t.Fatal("AsLocal() returned nil for local type")
	}

	// Verify all fields
	if len(local.Command) != 3 {
		t.Errorf("Expected 3 command elements, got %d", len(local.Command))
	}
	if local.Command[0] != "/usr/bin/mcp-server" {
		t.Errorf("Expected command[0] %q, got %q", "/usr/bin/mcp-server", local.Command[0])
	}
	if local.Command[1] != "--port" {
		t.Errorf("Expected command[1] %q, got %q", "--port", local.Command[1])
	}
	if local.Command[2] != "8080" {
		t.Errorf("Expected command[2] %q, got %q", "8080", local.Command[2])
	}
	if !local.Enabled {
		t.Error("Expected Enabled to be true")
	}
	if len(local.Environment) != 2 {
		t.Errorf("Expected 2 environment variables, got %d", len(local.Environment))
	}
	if local.Environment["NODE_ENV"] != "production" {
		t.Errorf("Expected NODE_ENV %q, got %q", "production", local.Environment["NODE_ENV"])
	}
	if local.Environment["PORT"] != "8080" {
		t.Errorf("Expected PORT %q, got %q", "8080", local.Environment["PORT"])
	}
}

func TestConfigMcp_AsRemote_ValidRemoteConfig(t *testing.T) {
	jsonData := `{
		"type": "remote",
		"url": "https://mcp.example.com/api",
		"enabled": false,
		"headers": {
			"Authorization": "Bearer token123",
			"X-API-Key": "key456"
		}
	}`

	var mcp ConfigMcp
	if err := json.Unmarshal([]byte(jsonData), &mcp); err != nil {
		t.Fatalf("Failed to unmarshal ConfigMcp: %v", err)
	}

	if mcp.Type != ConfigMcpTypeRemote {
		t.Errorf("Expected type %q, got %q", ConfigMcpTypeRemote, mcp.Type)
	}

	remote, err := mcp.AsRemote()
	if err != nil {
		t.Fatal("AsRemote() returned false for remote type")
	}
	if remote == nil {
		t.Fatal("AsRemote() returned nil for remote type")
	}

	// Verify all fields
	if remote.URL != "https://mcp.example.com/api" {
		t.Errorf("Expected URL %q, got %q", "https://mcp.example.com/api", remote.URL)
	}
	if remote.Enabled {
		t.Error("Expected Enabled to be false")
	}
	if len(remote.Headers) != 2 {
		t.Errorf("Expected 2 headers, got %d", len(remote.Headers))
	}
	if remote.Headers["Authorization"] != "Bearer token123" {
		t.Errorf("Expected Authorization %q, got %q", "Bearer token123", remote.Headers["Authorization"])
	}
	if remote.Headers["X-API-Key"] != "key456" {
		t.Errorf("Expected X-API-Key %q, got %q", "key456", remote.Headers["X-API-Key"])
	}
}

func TestConfigMcp_AsLocal_WrongType(t *testing.T) {
	jsonData := `{
		"type": "remote",
		"url": "https://mcp.example.com/api"
	}`

	var mcp ConfigMcp
	if err := json.Unmarshal([]byte(jsonData), &mcp); err != nil {
		t.Fatalf("Failed to unmarshal ConfigMcp: %v", err)
	}

	local, err := mcp.AsLocal()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if local != nil {
		t.Error("AsLocal() should return nil for remote type")
	}
}

func TestConfigMcp_AsRemote_WrongType(t *testing.T) {
	jsonData := `{
		"type": "local",
		"command": ["/usr/bin/mcp-server"]
	}`

	var mcp ConfigMcp
	if err := json.Unmarshal([]byte(jsonData), &mcp); err != nil {
		t.Fatalf("Failed to unmarshal ConfigMcp: %v", err)
	}

	remote, err := mcp.AsRemote()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if remote != nil {
		t.Error("AsRemote() should return nil for local type")
	}
}

func TestConfigMcp_InvalidJSON(t *testing.T) {
	jsonData := `{
		"type": "local",
		"command": "not-an-array"
	}`

	var mcp ConfigMcp
	if err := json.Unmarshal([]byte(jsonData), &mcp); err != nil {
		t.Fatalf("Failed to unmarshal ConfigMcp: %v", err)
	}

	// Should succeed in unmarshaling the type, but fail when trying to get as local
	if mcp.Type != ConfigMcpTypeLocal {
		t.Errorf("Expected type %q, got %q", ConfigMcpTypeLocal, mcp.Type)
	}

	local, err := mcp.AsLocal()
	if err == nil {
		t.Error("AsLocal() should return error for malformed JSON")
	}
	if local != nil {
		t.Error("AsLocal() should return nil for malformed JSON")
	}
}

func TestConfigMcp_MissingType(t *testing.T) {
	jsonData := `{
		"command": ["/usr/bin/mcp-server"]
	}`

	var mcp ConfigMcp
	if err := json.Unmarshal([]byte(jsonData), &mcp); err != nil {
		t.Fatalf("Failed to unmarshal ConfigMcp: %v", err)
	}

	if mcp.Type != "" {
		t.Errorf("Expected empty type, got %q", mcp.Type)
	}

	local, err := mcp.AsLocal()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if local != nil {
		t.Error("AsLocal() should return nil for missing type")
	}

	remote, err := mcp.AsRemote()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if remote != nil {
		t.Error("AsRemote() should return nil for missing type")
	}
}

func TestConfigMcp_UnknownType(t *testing.T) {
	jsonData := `{
		"type": "unknown"
	}`

	var mcp ConfigMcp
	if err := json.Unmarshal([]byte(jsonData), &mcp); err != nil {
		t.Fatalf("Failed to unmarshal ConfigMcp: %v", err)
	}

	if mcp.Type == ConfigMcpTypeLocal || mcp.Type == ConfigMcpTypeRemote {
		t.Errorf("Expected unknown type, got %q", mcp.Type)
	}

	local, err := mcp.AsLocal()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if local != nil {
		t.Error("AsLocal() should return nil for unknown type")
	}

	remote, err := mcp.AsRemote()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if remote != nil {
		t.Error("AsRemote() should return nil for unknown type")
	}
}

func TestConfigMcp_EmptyJSON(t *testing.T) {
	jsonData := `{}`

	var mcp ConfigMcp
	if err := json.Unmarshal([]byte(jsonData), &mcp); err != nil {
		t.Fatalf("Failed to unmarshal ConfigMcp: %v", err)
	}

	if mcp.Type != "" {
		t.Errorf("Expected empty type, got %q", mcp.Type)
	}

	local, err := mcp.AsLocal()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if local != nil {
		t.Error("AsLocal() should return nil for empty JSON")
	}

	remote, err := mcp.AsRemote()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if remote != nil {
		t.Error("AsRemote() should return nil for empty JSON")
	}
}

func TestConfigMcp_MalformedJSON(t *testing.T) {
	jsonData := `{
		"type": "local",
		"command": ["/usr/bin/mcp-server"],
		"environment": "not-a-map"
	}`

	var mcp ConfigMcp
	if err := json.Unmarshal([]byte(jsonData), &mcp); err != nil {
		t.Fatalf("Failed to unmarshal ConfigMcp: %v", err)
	}

	// Type should be set correctly
	if mcp.Type != ConfigMcpTypeLocal {
		t.Errorf("Expected type %q, got %q", ConfigMcpTypeLocal, mcp.Type)
	}

	// AsLocal should return error due to malformed environment field
	local, err := mcp.AsLocal()
	if err == nil {
		t.Error("AsLocal() should return error for malformed environment")
	}
	if local != nil {
		t.Error("AsLocal() should return nil for malformed environment")
	}
}

func TestConfigMcpType_IsKnown(t *testing.T) {
	tests := []struct {
		name     string
		mcpType  ConfigMcpType
		expected bool
	}{
		{"local type", ConfigMcpTypeLocal, true},
		{"remote type", ConfigMcpTypeRemote, true},
		{"unknown type", ConfigMcpType("unknown"), false},
		{"empty type", ConfigMcpType(""), false},
		{"random string", ConfigMcpType("random"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mcpType.IsKnown(); got != tt.expected {
				t.Errorf("IsKnown() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestConfigMcp_LocalMinimalFields(t *testing.T) {
	// Test with only required fields for local config
	jsonData := `{
		"type": "local",
		"command": ["/usr/bin/mcp-server"]
	}`

	var mcp ConfigMcp
	if err := json.Unmarshal([]byte(jsonData), &mcp); err != nil {
		t.Fatalf("Failed to unmarshal ConfigMcp: %v", err)
	}

	local, err := mcp.AsLocal()
	if err != nil {
		t.Fatal("AsLocal() returned false for valid minimal local config")
	}
	if local == nil {
		t.Fatal("AsLocal() returned nil for valid minimal local config")
	}

	if len(local.Command) != 1 {
		t.Errorf("Expected 1 command element, got %d", len(local.Command))
	}
	if local.Command[0] != "/usr/bin/mcp-server" {
		t.Errorf("Expected command[0] %q, got %q", "/usr/bin/mcp-server", local.Command[0])
	}
	if local.Enabled {
		t.Error("Expected Enabled to be false (zero value)")
	}
	if len(local.Environment) != 0 {
		t.Errorf("Expected nil or empty environment, got %d entries", len(local.Environment))
	}
}

func TestConfigMcp_RemoteMinimalFields(t *testing.T) {
	// Test with only required fields for remote config
	jsonData := `{
		"type": "remote",
		"url": "https://mcp.example.com/api"
	}`

	var mcp ConfigMcp
	if err := json.Unmarshal([]byte(jsonData), &mcp); err != nil {
		t.Fatalf("Failed to unmarshal ConfigMcp: %v", err)
	}

	remote, err := mcp.AsRemote()
	if err != nil {
		t.Fatal("AsRemote() returned false for valid minimal remote config")
	}
	if remote == nil {
		t.Fatal("AsRemote() returned nil for valid minimal remote config")
	}

	if remote.URL != "https://mcp.example.com/api" {
		t.Errorf("Expected URL %q, got %q", "https://mcp.example.com/api", remote.URL)
	}
	if remote.Enabled {
		t.Error("Expected Enabled to be false (zero value)")
	}
	if len(remote.Headers) != 0 {
		t.Errorf("Expected nil or empty headers, got %d entries", len(remote.Headers))
	}
}
