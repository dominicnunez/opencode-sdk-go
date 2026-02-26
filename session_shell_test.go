package opencode

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionService_Shell_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/session/sess123/shell" {
			t.Errorf("expected /session/sess123/shell, got %s", r.URL.Path)
		}

		// Verify request body
		var body SessionShellParams
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.Agent != "bash" {
			t.Errorf("expected agent=bash, got %s", body.Agent)
		}
		if body.Command != "ls -la" {
			t.Errorf("expected command='ls -la', got %s", body.Command)
		}

		// Return AssistantMessage
		resp := AssistantMessage{
			ID:        "msg456",
			Role:      AssistantMessageRoleAssistant,
			SessionID: "sess123",
			Cost:      0.01,
			Mode:      "auto",
			ModelID:   "claude-sonnet-4.5",
			ParentID:  "msg123",
			ProviderID: "anthropic",
			System:    []string{},
			Path: AssistantMessagePath{
				Cwd:  "/home/user",
				Root: "/home/user",
			},
			Time: AssistantMessageTime{
				Created: 1234567890,
			},
			Tokens: AssistantMessageTokens{
				Input:  100,
				Output: 200,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	params := &SessionShellParams{
		Agent:   "bash",
		Command: "ls -la",
	}

	result, err := client.Session.Shell(context.Background(), "sess123", params)
	if err != nil {
		t.Fatalf("Shell error: %v", err)
	}

	if result.ID != "msg456" {
		t.Errorf("expected ID=msg456, got %s", result.ID)
	}
	if result.SessionID != "sess123" {
		t.Errorf("expected SessionID=sess123, got %s", result.SessionID)
	}
	if result.Role != AssistantMessageRoleAssistant {
		t.Errorf("expected role=assistant, got %s", result.Role)
	}
}

func TestSessionService_Shell_WithDirectory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query param
		if dir := r.URL.Query().Get("directory"); dir != "/tmp" {
			t.Errorf("expected directory=/tmp, got %s", dir)
		}

		// Return minimal AssistantMessage
		resp := AssistantMessage{
			ID:        "msg789",
			Role:      AssistantMessageRoleAssistant,
			SessionID: "sess123",
			Cost:      0.01,
			Mode:      "auto",
			ModelID:   "claude-sonnet-4.5",
			ParentID:  "msg123",
			ProviderID: "anthropic",
			System:    []string{},
			Path: AssistantMessagePath{
				Cwd:  "/tmp",
				Root: "/tmp",
			},
			Time: AssistantMessageTime{
				Created: 1234567890,
			},
			Tokens: AssistantMessageTokens{
				Input:  50,
				Output: 100,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	directory := "/tmp"
	params := &SessionShellParams{
		Agent:     "bash",
		Command:   "pwd",
		Directory: &directory,
	}

	result, err := client.Session.Shell(context.Background(), "sess123", params)
	if err != nil {
		t.Fatalf("Shell error: %v", err)
	}

	if result.Path.Cwd != "/tmp" {
		t.Errorf("expected cwd=/tmp, got %s", result.Path.Cwd)
	}
}

func TestSessionService_Shell_MissingID(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	params := &SessionShellParams{
		Agent:   "bash",
		Command: "ls",
	}

	_, err = client.Session.Shell(context.Background(), "", params)
	if err == nil {
		t.Fatal("expected error for missing ID, got nil")
	}
	if err.Error() != "missing required id parameter" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestSessionService_Shell_MissingParams(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.Session.Shell(context.Background(), "sess123", nil)
	if err == nil {
		t.Fatal("expected error for missing params, got nil")
	}
	if err.Error() != "params is required" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestSessionService_Shell_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"internal error"}`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	params := &SessionShellParams{
		Agent:   "bash",
		Command: "ls",
	}

	_, err = client.Session.Shell(context.Background(), "sess123", params)
	if err == nil {
		t.Fatal("expected error for server error, got nil")
	}
}

func TestSessionService_Shell_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{invalid json`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	params := &SessionShellParams{
		Agent:   "bash",
		Command: "ls",
	}

	_, err = client.Session.Shell(context.Background(), "sess123", params)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestSessionShellParams_Marshal(t *testing.T) {
	params := SessionShellParams{
		Agent:   "bash",
		Command: "echo hello",
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if decoded["agent"] != "bash" {
		t.Errorf("expected agent=bash, got %v", decoded["agent"])
	}
	if decoded["command"] != "echo hello" {
		t.Errorf("expected command='echo hello', got %v", decoded["command"])
	}
}

func TestSessionShellParams_URLQuery_WithDirectory(t *testing.T) {
	directory := "/home/user"
	params := SessionShellParams{
		Agent:     "bash",
		Command:   "ls",
		Directory: &directory,
	}

	values, err := params.URLQuery()
	if err != nil {
		t.Fatalf("URLQuery error: %v", err)
	}

	if values.Get("directory") != "/home/user" {
		t.Errorf("expected directory=/home/user, got %s", values.Get("directory"))
	}
}

func TestSessionShellParams_URLQuery_WithoutDirectory(t *testing.T) {
	params := SessionShellParams{
		Agent:   "bash",
		Command: "ls",
	}

	values, err := params.URLQuery()
	if err != nil {
		t.Fatalf("URLQuery error: %v", err)
	}

	if values.Get("directory") != "" {
		t.Errorf("expected empty directory, got %s", values.Get("directory"))
	}
}
