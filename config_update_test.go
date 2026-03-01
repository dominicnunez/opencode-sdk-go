package opencode

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfigUpdate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH method, got %s", r.Method)
		}
		if r.URL.Path != "/config" {
			t.Errorf("Expected path /config, got %s", r.URL.Path)
		}

		// Verify request body
		var received Config
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if received.Model != "anthropic/claude-sonnet-4" {
			t.Errorf("Expected model anthropic/claude-sonnet-4, got %s", received.Model)
		}
		if received.Theme != "dark" {
			t.Errorf("Expected theme dark, got %s", received.Theme)
		}

		// Return updated config
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Config{
			Model: "anthropic/claude-sonnet-4",
			Theme: "dark",
		})
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	params := &ConfigUpdateParams{
		Config: Config{
			Model: "anthropic/claude-sonnet-4",
			Theme: "dark",
		},
	}

	result, err := client.Config.Update(context.Background(), params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if result.Model != "anthropic/claude-sonnet-4" {
		t.Errorf("Expected model anthropic/claude-sonnet-4, got %s", result.Model)
	}
	if result.Theme != "dark" {
		t.Errorf("Expected theme dark, got %s", result.Theme)
	}
}

func TestConfigUpdate_WithDirectoryQueryParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH method, got %s", r.Method)
		}
		if r.URL.Path != "/config" {
			t.Errorf("Expected path /config, got %s", r.URL.Path)
		}

		// Verify query param
		if r.URL.Query().Get("directory") != "/workspace/project" {
			t.Errorf("Expected directory query param /workspace/project, got %s", r.URL.Query().Get("directory"))
		}

		// Verify request body
		var received Config
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Config{
			Model: "openai/gpt-4",
		})
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	dir := "/workspace/project"
	params := &ConfigUpdateParams{
		Config: Config{
			Model: "openai/gpt-4",
		},
		Directory: &dir,
	}

	result, err := client.Config.Update(context.Background(), params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if result.Model != "openai/gpt-4" {
		t.Errorf("Expected model openai/gpt-4, got %s", result.Model)
	}
}

func TestConfigUpdate_NilParams(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Config.Update(context.Background(), nil)
	if err == nil {
		t.Fatal("Expected error for nil params, got nil")
	}
	if err.Error() != "params is required" {
		t.Errorf("Expected error 'params is required', got %s", err.Error())
	}
}

func TestConfigUpdate_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	params := &ConfigUpdateParams{
		Config: Config{
			Model: "anthropic/claude-sonnet-4",
		},
	}

	_, err = client.Config.Update(context.Background(), params)
	if err == nil {
		t.Fatal("expected error for server error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, apiErr.StatusCode)
	}
}

func TestConfigUpdate_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{invalid json`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	params := &ConfigUpdateParams{
		Config: Config{
			Model: "anthropic/claude-sonnet-4",
		},
	}

	_, err = client.Config.Update(context.Background(), params)
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
}

func TestConfigUpdate_UnionTypeRoundTrip(t *testing.T) {
	mcpJSON := `{"type":"local","command":["/usr/bin/mcp-server"],"enabled":true,"environment":{"NODE_ENV":"production"}}`

	var mcp ConfigMcp
	if err := json.Unmarshal([]byte(mcpJSON), &mcp); err != nil {
		t.Fatalf("unmarshal ConfigMcp: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var received map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode request body: %v", err)
		}

		mcpRaw, ok := received["mcp"]
		if !ok {
			t.Fatal("expected mcp field in request body")
		}

		var mcpMap map[string]json.RawMessage
		if err := json.Unmarshal(mcpRaw, &mcpMap); err != nil {
			t.Fatalf("unmarshal mcp map: %v", err)
		}

		serverRaw, ok := mcpMap["test-server"]
		if !ok {
			t.Fatal("expected test-server key in mcp map")
		}

		var entry map[string]interface{}
		if err := json.Unmarshal(serverRaw, &entry); err != nil {
			t.Fatalf("unmarshal mcp entry: %v", err)
		}
		if entry["type"] != "local" {
			t.Errorf("expected type local, got %v", entry["type"])
		}
		cmd, ok := entry["command"].([]interface{})
		if !ok || len(cmd) != 1 || cmd[0] != "/usr/bin/mcp-server" {
			t.Errorf("unexpected command: %v", entry["command"])
		}

		// Echo back the config
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Config{})
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("create client: %v", err)
	}

	params := &ConfigUpdateParams{
		Config: Config{
			Mcp: map[string]ConfigMcp{
				"test-server": mcp,
			},
		},
	}

	_, err = client.Config.Update(context.Background(), params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
}

func TestConfigUpdateParams_MarshalJSON(t *testing.T) {
	params := ConfigUpdateParams{
		Config: Config{
			Model:      "anthropic/claude-sonnet-4",
			Theme:      "dark",
			Autoupdate: true,
		},
		Directory: nil,
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	var decoded Config
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Model != "anthropic/claude-sonnet-4" {
		t.Errorf("Expected model anthropic/claude-sonnet-4, got %s", decoded.Model)
	}
	if decoded.Theme != "dark" {
		t.Errorf("Expected theme dark, got %s", decoded.Theme)
	}
	if !decoded.Autoupdate {
		t.Error("Expected autoupdate true, got false")
	}
}

func TestConfigUpdateParams_MarshalJSON_OmitsZeroValues(t *testing.T) {
	params := ConfigUpdateParams{
		Config: Config{
			Theme: "dark",
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal into map failed: %v", err)
	}

	if raw["theme"] != "dark" {
		t.Errorf("Expected theme=dark, got %v", raw["theme"])
	}
	if _, ok := raw["autoshare"]; ok {
		t.Error("Zero-value bool 'autoshare' should be omitted from PATCH body")
	}
	if _, ok := raw["autoupdate"]; ok {
		t.Error("Zero-value bool 'autoupdate' should be omitted from PATCH body")
	}
	if _, ok := raw["snapshot"]; ok {
		t.Error("Zero-value bool 'snapshot' should be omitted from PATCH body")
	}
	if _, ok := raw["model"]; ok {
		t.Error("Zero-value string 'model' should be omitted from PATCH body")
	}
}

func TestConfigUpdateParams_URLQuery(t *testing.T) {
	tests := []struct {
		name      string
		params    ConfigUpdateParams
		expectDir string
		hasDir    bool
	}{
		{
			name: "with directory",
			params: ConfigUpdateParams{
				Config:    Config{Model: "test"},
				Directory: ptrString("/test/dir"),
			},
			expectDir: "/test/dir",
			hasDir:    true,
		},
		{
			name: "without directory",
			params: ConfigUpdateParams{
				Config:    Config{Model: "test"},
				Directory: nil,
			},
			hasDir: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vals, err := tt.params.URLQuery()
			if err != nil {
				t.Fatalf("URLQuery failed: %v", err)
			}

			if tt.hasDir {
				if vals.Get("directory") != tt.expectDir {
					t.Errorf("Expected directory %s, got %s", tt.expectDir, vals.Get("directory"))
				}
			} else {
				if vals.Has("directory") {
					t.Error("Expected no directory param, but it was present")
				}
			}
		})
	}
}
