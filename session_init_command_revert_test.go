package opencode

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestSessionInit_Success verifies Init sends POST to /session/{id}/init with correct body and decodes bool response
func TestSessionInit_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/session/sess_123/init" {
			t.Errorf("expected path /session/sess_123/init, got %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var parsed map[string]interface{}
		if err := json.Unmarshal(body, &parsed); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}
		if parsed["messageID"] != "msg_001" {
			t.Errorf("expected messageID msg_001, got %v", parsed["messageID"])
		}
		if parsed["modelID"] != "gpt-4" {
			t.Errorf("expected modelID gpt-4, got %v", parsed["modelID"])
		}
		if parsed["providerID"] != "openai" {
			t.Errorf("expected providerID openai, got %v", parsed["providerID"])
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(true)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Session.Init(context.Background(), "sess_123", &SessionInitParams{
		MessageID:  "msg_001",
		ModelID:    "gpt-4",
		ProviderID: "openai",
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if !result {
		t.Error("expected Init to return true, got false")
	}
}

// TestSessionCommand_Success verifies Command sends POST to /session/{id}/command with correct body and decodes response
func TestSessionCommand_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/session/sess_123/command" {
			t.Errorf("expected path /session/sess_123/command, got %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var parsed map[string]interface{}
		if err := json.Unmarshal(body, &parsed); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}
		if parsed["command"] != "/ask" {
			t.Errorf("expected command /ask, got %v", parsed["command"])
		}
		if parsed["arguments"] != "what is this project" {
			t.Errorf("expected arguments 'what is this project', got %v", parsed["arguments"])
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"info": map[string]interface{}{
				"id":         "msg_001",
				"sessionID":  "sess_123",
				"role":       "assistant",
				"cost":       0.0,
				"mode":       "",
				"modelID":    "",
				"parentID":   "",
				"path":       map[string]interface{}{"cwd": "", "root": ""},
				"providerID": "",
				"system":     []string{},
				"time":       map[string]interface{}{"created": 0.0, "completed": 0.0},
				"tokens":     map[string]interface{}{"input": 0, "output": 0, "reasoning": 0, "cache": map[string]interface{}{"read": 0, "write": 0}},
			},
			"parts": []interface{}{},
		})
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.Session.Command(context.Background(), "sess_123", &SessionCommandParams{
		Command:   "/ask",
		Arguments: "what is this project",
	})
	if err != nil {
		t.Fatalf("Command failed: %v", err)
	}

	if resp.Info.ID != "msg_001" {
		t.Errorf("expected info ID msg_001, got %s", resp.Info.ID)
	}
	if resp.Info.SessionID != "sess_123" {
		t.Errorf("expected info sessionID sess_123, got %s", resp.Info.SessionID)
	}
	if resp.Info.Role != AssistantMessageRoleAssistant {
		t.Errorf("expected info role assistant, got %s", resp.Info.Role)
	}
}

// TestSessionRevert_Success verifies Revert sends POST to /session/{id}/revert with correct body and decodes session response
func TestSessionRevert_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/session/sess_123/revert" {
			t.Errorf("expected path /session/sess_123/revert, got %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var parsed map[string]interface{}
		if err := json.Unmarshal(body, &parsed); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}
		if parsed["messageID"] != "msg_to_revert" {
			t.Errorf("expected messageID msg_to_revert, got %v", parsed["messageID"])
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":        "sess_123",
			"directory": "/test/path",
			"projectID": "proj_456",
			"title":     "Reverted Session",
			"version":   "1.0.0",
			"time": map[string]interface{}{
				"created": 1234567890.0,
				"updated": 1234567900.0,
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	session, err := client.Session.Revert(context.Background(), "sess_123", &SessionRevertParams{
		MessageID: "msg_to_revert",
	})
	if err != nil {
		t.Fatalf("Revert failed: %v", err)
	}

	if session.ID != "sess_123" {
		t.Errorf("expected session ID sess_123, got %s", session.ID)
	}
	if session.Title != "Reverted Session" {
		t.Errorf("expected title 'Reverted Session', got %s", session.Title)
	}
	if session.ProjectID != "proj_456" {
		t.Errorf("expected projectID proj_456, got %s", session.ProjectID)
	}
}

func TestSessionCommand_RequiresCommand(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Command(context.Background(), "sess_123", &SessionCommandParams{})
	if err == nil {
		t.Fatal("expected error for missing command")
	}
	if err.Error() != "missing required command parameter" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSessionInit_RequiresFields(t *testing.T) {
	tests := []struct {
		name    string
		params  *SessionInitParams
		wantErr string
	}{
		{
			name:    "missing message id",
			params:  &SessionInitParams{ModelID: "gpt-4", ProviderID: "openai"},
			wantErr: "missing required messageID parameter",
		},
		{
			name:    "missing model id",
			params:  &SessionInitParams{MessageID: "msg_001", ProviderID: "openai"},
			wantErr: "missing required modelID parameter",
		},
		{
			name:    "missing provider id",
			params:  &SessionInitParams{MessageID: "msg_001", ModelID: "gpt-4"},
			wantErr: "missing required providerID parameter",
		},
	}

	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Session.Init(context.Background(), "sess_123", tt.params)
			if err == nil {
				t.Fatal("expected validation error")
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("expected %q, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestSessionRevert_RequiresMessageID(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Revert(context.Background(), "sess_123", &SessionRevertParams{})
	if err == nil {
		t.Fatal("expected error for missing messageID")
	}
	if err.Error() != "missing required messageID parameter" {
		t.Fatalf("unexpected error: %v", err)
	}
}
