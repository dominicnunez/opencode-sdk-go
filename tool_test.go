package opencode

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestToolService_IDs_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/experimental/tool/ids" {
			t.Errorf("expected path /experimental/tool/ids, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]string{"tool1", "tool2", "tool3"})
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ids, err := client.Tool.IDs(context.Background(), nil)
	if err != nil {
		t.Fatalf("IDs failed: %v", err)
	}

	if ids == nil {
		t.Fatal("expected non-nil ToolIDs")
	}

	if len(*ids) != 3 {
		t.Errorf("expected 3 tool IDs, got %d", len(*ids))
	}

	expected := []string{"tool1", "tool2", "tool3"}
	for i, id := range *ids {
		if id != expected[i] {
			t.Errorf("expected tool ID %s at index %d, got %s", expected[i], i, id)
		}
	}
}

func TestToolService_IDs_WithDirectory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dir := r.URL.Query().Get("directory")
		if dir != "/test/path" {
			t.Errorf("expected directory=/test/path, got %s", dir)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]string{"tool1"})
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dir := "/test/path"
	ids, err := client.Tool.IDs(context.Background(), &ToolIDsParams{
		Directory: &dir,
	})
	if err != nil {
		t.Fatalf("IDs failed: %v", err)
	}

	if ids == nil || len(*ids) != 1 {
		t.Errorf("expected 1 tool ID, got %v", ids)
	}
}

func TestToolService_IDs_EmptyArray(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]string{})
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ids, err := client.Tool.IDs(context.Background(), nil)
	if err != nil {
		t.Fatalf("IDs failed: %v", err)
	}

	if ids == nil {
		t.Fatal("expected non-nil ToolIDs")
	}

	if len(*ids) != 0 {
		t.Errorf("expected 0 tool IDs, got %d", len(*ids))
	}
}

func TestToolService_IDs_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ids, err := client.Tool.IDs(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if ids != nil {
		t.Errorf("expected nil ToolIDs on error, got %v", ids)
	}
}

func TestToolService_IDs_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("not valid json"))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ids, err := client.Tool.IDs(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if ids != nil {
		t.Errorf("expected nil ToolIDs on error, got %v", ids)
	}
}

func TestToolService_List_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/experimental/tool" {
			t.Errorf("expected path /experimental/tool, got %s", r.URL.Path)
		}

		provider := r.URL.Query().Get("provider")
		model := r.URL.Query().Get("model")
		if provider != "anthropic" {
			t.Errorf("expected provider=anthropic, got %s", provider)
		}
		if model != "claude-3-opus" {
			t.Errorf("expected model=claude-3-opus, got %s", model)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]ToolListItem{
			{
				ID:          "Read",
				Description: "Read a file from the filesystem",
				Parameters:  json.RawMessage(`{"type":"object","properties":{"file_path":{"type":"string","description":"The path to the file"}},"required":["file_path"]}`),
			},
			{
				ID:          "Write",
				Description: "Write a file to the filesystem",
				Parameters:  json.RawMessage(`{"type":"object","properties":{"file_path":{"type":"string","description":"The path to the file"},"content":{"type":"string","description":"The content to write"}},"required":["file_path","content"]}`),
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tools, err := client.Tool.List(context.Background(), &ToolListParams{
		Provider: "anthropic",
		Model:    "claude-3-opus",
	})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if tools == nil {
		t.Fatal("expected non-nil ToolList")
	}

	if len(*tools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(*tools))
	}

	// Verify first tool
	tool1 := (*tools)[0]
	if tool1.ID != "Read" {
		t.Errorf("expected tool ID 'Read', got %s", tool1.ID)
	}
	if tool1.Description != "Read a file from the filesystem" {
		t.Errorf("expected description 'Read a file from the filesystem', got %s", tool1.Description)
	}
	if tool1.Parameters == nil {
		t.Error("expected non-nil Parameters")
	}

	// Verify parameters structure for Read tool
	var params map[string]interface{}
	if err := json.Unmarshal(tool1.Parameters, &params); err != nil {
		t.Errorf("failed to unmarshal Parameters: %v", err)
	} else {
		if params["type"] != "object" {
			t.Errorf("expected type=object, got %v", params["type"])
		}
	}

	// Verify second tool
	tool2 := (*tools)[1]
	if tool2.ID != "Write" {
		t.Errorf("expected tool ID 'Write', got %s", tool2.ID)
	}
}

func TestToolService_List_WithDirectory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dir := r.URL.Query().Get("directory")
		provider := r.URL.Query().Get("provider")
		model := r.URL.Query().Get("model")

		if dir != "/test/path" {
			t.Errorf("expected directory=/test/path, got %s", dir)
		}
		if provider != "openai" {
			t.Errorf("expected provider=openai, got %s", provider)
		}
		if model != "gpt-4" {
			t.Errorf("expected model=gpt-4, got %s", model)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]ToolListItem{})
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dir := "/test/path"
	tools, err := client.Tool.List(context.Background(), &ToolListParams{
		Directory: &dir,
		Provider:  "openai",
		Model:     "gpt-4",
	})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if tools == nil {
		t.Fatal("expected non-nil ToolList")
	}
}

func TestToolService_List_NilParams(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tools, err := client.Tool.List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error when params is nil, got nil")
	}

	if tools != nil {
		t.Errorf("expected nil ToolList on error, got %v", tools)
	}
}

func TestToolService_List_MissingProvider(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Tool.List(context.Background(), &ToolListParams{
		Model: "claude-3-opus",
	})
	if err == nil {
		t.Fatal("expected error when Provider is empty")
	}
	if !strings.Contains(err.Error(), "required query parameter") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestToolService_List_MissingModel(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Tool.List(context.Background(), &ToolListParams{
		Provider: "anthropic",
	})
	if err == nil {
		t.Fatal("expected error when Model is empty")
	}
	if !strings.Contains(err.Error(), "required query parameter") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestToolService_List_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tools, err := client.Tool.List(context.Background(), &ToolListParams{
		Provider: "anthropic",
		Model:    "claude-3-opus",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if tools != nil {
		t.Errorf("expected nil ToolList on error, got %v", tools)
	}
}

func TestToolService_List_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tools, err := client.Tool.List(context.Background(), &ToolListParams{
		Provider: "anthropic",
		Model:    "claude-3-opus",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if tools != nil {
		t.Errorf("expected nil ToolList on error, got %v", tools)
	}
}

func TestToolIDsParams_URLQuery(t *testing.T) {
	tests := []struct {
		name     string
		params   ToolIDsParams
		expected map[string]string
	}{
		{
			name:     "no directory",
			params:   ToolIDsParams{},
			expected: map[string]string{},
		},
		{
			name: "with directory",
			params: ToolIDsParams{
				Directory: ptrString("/test/dir"),
			},
			expected: map[string]string{
				"directory": "/test/dir",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := tt.params.URLQuery()
			if err != nil {
				t.Fatalf("URLQuery failed: %v", err)
			}

			for k, v := range tt.expected {
				vals := query[k]
				if len(vals) == 0 || vals[0] != v {
					t.Errorf("expected %s=%s, got %v", k, v, vals)
				}
			}

			// Verify no extra keys
			if len(query) != len(tt.expected) {
				t.Errorf("expected %d query params, got %d", len(tt.expected), len(query))
			}
		})
	}
}

func TestToolListParams_URLQuery(t *testing.T) {
	tests := []struct {
		name     string
		params   ToolListParams
		expected map[string]string
	}{
		{
			name: "required params only",
			params: ToolListParams{
				Provider: "anthropic",
				Model:    "claude-3-opus",
			},
			expected: map[string]string{
				"provider": "anthropic",
				"model":    "claude-3-opus",
			},
		},
		{
			name: "with directory",
			params: ToolListParams{
				Directory: ptrString("/project"),
				Provider:  "openai",
				Model:     "gpt-4",
			},
			expected: map[string]string{
				"directory": "/project",
				"provider":  "openai",
				"model":     "gpt-4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := tt.params.URLQuery()
			if err != nil {
				t.Fatalf("URLQuery failed: %v", err)
			}

			for k, v := range tt.expected {
				vals := query[k]
				if len(vals) == 0 || vals[0] != v {
					t.Errorf("expected %s=%s, got %v", k, v, vals)
				}
			}

			// Verify no extra keys
			if len(query) != len(tt.expected) {
				t.Errorf("expected %d query params, got %d", len(tt.expected), len(query))
			}
		})
	}
}

func TestToolListItem_Unmarshal(t *testing.T) {
	jsonData := `{
		"id": "Bash",
		"description": "Execute bash commands",
		"parameters": {
			"type": "object",
			"properties": {
				"command": {
					"type": "string",
					"description": "The command to execute"
				}
			},
			"required": ["command"]
		}
	}`

	var item ToolListItem
	err := json.Unmarshal([]byte(jsonData), &item)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if item.ID != "Bash" {
		t.Errorf("expected ID=Bash, got %s", item.ID)
	}

	if item.Description != "Execute bash commands" {
		t.Errorf("expected description='Execute bash commands', got %s", item.Description)
	}

	var params map[string]interface{}
	if err := json.Unmarshal(item.Parameters, &params); err != nil {
		t.Fatalf("failed to unmarshal Parameters: %v", err)
	}

	if params["type"] != "object" {
		t.Errorf("expected type=object, got %v", params["type"])
	}

	props, ok := params["properties"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected properties to be map[string]interface{}, got %T", params["properties"])
	}

	command, ok := props["command"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected command to be map[string]interface{}, got %T", props["command"])
	}

	if command["type"] != "string" {
		t.Errorf("expected command type=string, got %v", command["type"])
	}
}
