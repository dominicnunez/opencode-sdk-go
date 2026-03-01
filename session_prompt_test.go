package opencode

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionPrompt_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/session/sess_1/message" {
			t.Errorf("expected path /session/sess_1/message, got %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var raw map[string]interface{}
		if err := json.Unmarshal(body, &raw); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		parts, ok := raw["parts"].([]interface{})
		if !ok || len(parts) == 0 {
			t.Fatal("expected non-empty parts array in body")
		}

		resp := SessionPromptResponse{
			Info: AssistantMessage{
				ID:        "msg_1",
				Role:      AssistantMessageRoleAssistant,
				SessionID: "sess_1",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Session.Prompt(context.Background(), "sess_1", &SessionPromptParams{
		Parts: []SessionPromptParamsPartUnion{
			TextPartInputParam{
				Text: "hello",
				Type: TextPartInputTypeText,
			},
		},
	})
	if err != nil {
		t.Fatalf("Prompt failed: %v", err)
	}
	if result.Info.ID != "msg_1" {
		t.Errorf("expected message ID msg_1, got %s", result.Info.ID)
	}
	if result.Info.SessionID != "sess_1" {
		t.Errorf("expected session ID sess_1, got %s", result.Info.SessionID)
	}
}

func TestSessionPrompt_MissingID(t *testing.T) {
	client, err := NewClient(WithBaseURL("http://localhost"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Prompt(context.Background(), "", &SessionPromptParams{})
	if err == nil {
		t.Fatal("expected error for missing ID, got nil")
	}
	if err.Error() != "missing required id parameter" {
		t.Errorf("expected 'missing required id parameter', got: %v", err)
	}
}

func TestSessionPrompt_NilParams(t *testing.T) {
	client, err := NewClient(WithBaseURL("http://localhost"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Prompt(context.Background(), "sess_1", nil)
	if err == nil {
		t.Fatal("expected error for nil params, got nil")
	}
	if err.Error() != "params is required" {
		t.Errorf("expected 'params is required', got: %v", err)
	}
}

func TestSessionPrompt_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Prompt(context.Background(), "sess_1", &SessionPromptParams{
		Parts: []SessionPromptParamsPartUnion{
			TextPartInputParam{Text: "hello", Type: TextPartInputTypeText},
		},
	})
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
