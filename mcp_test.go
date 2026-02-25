package opencode

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMcpService_Status_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if r.URL.Path != "/mcp" {
			t.Errorf("expected path /mcp, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"enabled": true,
			"servers": map[string]interface{}{
				"local-server": map[string]interface{}{
					"status":  "running",
					"command": []string{"/bin/mcp-server"},
				},
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	status, err := client.Mcp.Status(context.Background(), nil)
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}

	if status == nil {
		t.Fatal("expected non-nil status")
	}

	enabled, ok := (*status)["enabled"].(bool)
	if !ok || !enabled {
		t.Errorf("expected enabled=true, got %v", (*status)["enabled"])
	}
}

func TestMcpService_Status_WithDirectory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		directory := r.URL.Query().Get("directory")
		if directory != "/test/path" {
			t.Errorf("expected directory=/test/path, got %s", directory)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"enabled": false,
			"servers": map[string]interface{}{},
		})
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dir := "/test/path"
	status, err := client.Mcp.Status(context.Background(), &McpStatusParams{
		Directory: &dir,
	})
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}

	if status == nil {
		t.Fatal("expected non-nil status")
	}

	enabled, ok := (*status)["enabled"].(bool)
	if !ok || enabled {
		t.Errorf("expected enabled=false, got %v", (*status)["enabled"])
	}
}

func TestMcpService_Status_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	status, err := client.Mcp.Status(context.Background(), nil)
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}

	if status == nil {
		t.Fatal("expected non-nil status")
	}

	if len(*status) != 0 {
		t.Errorf("expected empty status, got %v", *status)
	}
}

func TestMcpService_Status_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	status, err := client.Mcp.Status(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if status != nil {
		t.Errorf("expected nil status on error, got %v", status)
	}
}

func TestMcpService_Status_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{invalid json"))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	status, err := client.Mcp.Status(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
	if status != nil {
		t.Errorf("expected nil status on error, got %v", status)
	}
}

func TestMcpService_Status_ComplexNestedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"enabled": true,
			"servers": map[string]interface{}{
				"server1": map[string]interface{}{
					"status": "running",
					"tools": []interface{}{
						map[string]interface{}{
							"name":        "tool1",
							"description": "A test tool",
							"parameters": map[string]interface{}{
								"param1": "string",
								"param2": 123,
							},
						},
					},
				},
			},
			"metadata": map[string]interface{}{
				"version": "1.0.0",
				"count":   5,
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	status, err := client.Mcp.Status(context.Background(), nil)
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}

	if status == nil {
		t.Fatal("expected non-nil status")
	}

	// Verify nested structure is preserved
	servers, ok := (*status)["servers"].(map[string]interface{})
	if !ok {
		t.Fatal("expected servers map")
	}

	server1, ok := servers["server1"].(map[string]interface{})
	if !ok {
		t.Fatal("expected server1 map")
	}

	serverStatus, ok := server1["status"].(string)
	if !ok || serverStatus != "running" {
		t.Errorf("expected server1 status=running, got %v", server1["status"])
	}

	metadata, ok := (*status)["metadata"].(map[string]interface{})
	if !ok {
		t.Fatal("expected metadata map")
	}

	version, ok := metadata["version"].(string)
	if !ok || version != "1.0.0" {
		t.Errorf("expected metadata version=1.0.0, got %v", metadata["version"])
	}
}

func TestMcpStatusParams_URLQuery(t *testing.T) {
	tests := []struct {
		name   string
		params McpStatusParams
		want   map[string]string
	}{
		{
			name:   "with directory",
			params: McpStatusParams{Directory: ptrString("/test/path")},
			want:   map[string]string{"directory": "/test/path"},
		},
		{
			name:   "without directory",
			params: McpStatusParams{},
			want:   map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, err := tt.params.URLQuery()
			if err != nil {
				t.Fatalf("URLQuery() error = %v", err)
			}

			for key, expectedValue := range tt.want {
				if got := values.Get(key); got != expectedValue {
					t.Errorf("URLQuery() key %s = %v, want %v", key, got, expectedValue)
				}
			}

			// Verify no extra keys
			if len(values) != len(tt.want) {
				t.Errorf("URLQuery() returned %d keys, want %d", len(values), len(tt.want))
			}
		})
	}
}
