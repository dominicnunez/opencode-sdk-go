package opencode

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfigGet_SuccessWithUnionTypes(t *testing.T) {
	responseJSON := `{
		"model": "anthropic/claude-sonnet-4",
		"theme": "dark",
		"mcp": {
			"my-server": {
				"type": "local",
				"command": ["npx", "-y", "my-mcp-server"],
				"enabled": true,
				"environment": {"NODE_ENV": "production"}
			}
		},
		"agent": {
			"build": {
				"permission": {
					"bash": "allow",
					"edit": "ask",
					"webfetch": "deny"
				}
			},
			"general": {},
			"plan": {}
		},
		"provider": {
			"custom": {
				"id": "custom",
				"name": "Custom Provider",
				"api": "https://api.custom.com",
				"options": {
					"apiKey": "sk-test",
					"baseURL": "https://api.custom.com/v1",
					"timeout": 60000
				}
			}
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		if r.URL.Path != "/config" {
			t.Errorf("Expected path /config, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	result, err := client.Config.Get(context.Background(), &ConfigGetParams{})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if result.Model != "anthropic/claude-sonnet-4" {
		t.Errorf("Expected model anthropic/claude-sonnet-4, got %s", result.Model)
	}
	if result.Theme != "dark" {
		t.Errorf("Expected theme dark, got %s", result.Theme)
	}

	// Verify ConfigMcp union: AsLocal() on a local type entry
	mcpEntry, ok := result.Mcp["my-server"]
	if !ok {
		t.Fatal("Expected mcp entry 'my-server' to exist")
	}
	if mcpEntry.Type != ConfigMcpTypeLocal {
		t.Errorf("Expected mcp type 'local', got %s", mcpEntry.Type)
	}
	local, err := mcpEntry.AsLocal()
	if err != nil {
		t.Fatalf("AsLocal failed: %v", err)
	}
	if len(local.Command) != 3 || local.Command[0] != "npx" {
		t.Errorf("Expected command [npx -y my-mcp-server], got %v", local.Command)
	}
	if local.Environment["NODE_ENV"] != "production" {
		t.Errorf("Expected environment NODE_ENV=production, got %s", local.Environment["NODE_ENV"])
	}
	// AsRemote() on a local entry should fail
	_, err = mcpEntry.AsRemote()
	if err == nil {
		t.Error("Expected error calling AsRemote() on a local mcp entry")
	}

	// Verify ConfigAgentBuildPermissionBashUnion: AsString() on a string variant
	bashPerm := result.Agent.Build.Permission.Bash
	bashStr, err := bashPerm.AsString()
	if err != nil {
		t.Fatalf("AsString on bash permission failed: %v", err)
	}
	if bashStr != ConfigAgentBuildPermissionBashStringAllow {
		t.Errorf("Expected bash permission 'allow', got %s", bashStr)
	}
	// AsMap() on a string variant should fail
	_, err = bashPerm.AsMap()
	if err == nil {
		t.Error("Expected error calling AsMap() on a string bash permission")
	}

	// Verify ConfigProviderOptionsTimeoutUnion: AsInt() on a numeric variant
	provider, ok := result.Provider["custom"]
	if !ok {
		t.Fatal("Expected provider 'custom' to exist")
	}
	timeoutVal, err := provider.Options.Timeout.AsInt()
	if err != nil {
		t.Fatalf("AsInt on timeout failed: %v", err)
	}
	if timeoutVal != 60000 {
		t.Errorf("Expected timeout 60000, got %d", timeoutVal)
	}
	// AsBool() on a numeric variant should fail
	_, err = provider.Options.Timeout.AsBool()
	if err == nil {
		t.Error("Expected error calling AsBool() on a numeric timeout")
	}
}

func TestConfigGet_WithDirectoryQueryParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		if r.URL.Query().Get("directory") != "/workspace/project" {
			t.Errorf("Expected directory query param /workspace/project, got %s", r.URL.Query().Get("directory"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"model": "openai/gpt-4"}`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	result, err := client.Config.Get(context.Background(), &ConfigGetParams{
		Directory: ptrString("/workspace/project"),
	})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if result.Model != "openai/gpt-4" {
		t.Errorf("Expected model openai/gpt-4, got %s", result.Model)
	}
}

func TestConfigGet_NilParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"model": "anthropic/claude-sonnet-4"}`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	result, err := client.Config.Get(context.Background(), nil)
	if err != nil {
		t.Fatalf("Get with nil params failed: %v", err)
	}

	if result.Model != "anthropic/claude-sonnet-4" {
		t.Errorf("Expected model anthropic/claude-sonnet-4, got %s", result.Model)
	}
}
