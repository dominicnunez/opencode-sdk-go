package opencode

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionMessages_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/session/sess_1/message" {
			t.Errorf("expected path /session/sess_1/message, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"info":{"id":"msg_1","role":"user","sessionID":"sess_1"},"parts":[]},
			{"info":{"id":"msg_2","role":"assistant","sessionID":"sess_1"},"parts":[]}
		]`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	messages, err := client.Session.Messages(context.Background(), "sess_1", nil)
	if err != nil {
		t.Fatalf("Messages failed: %v", err)
	}
	if len(messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(messages))
	}
	if messages[0].Info.ID != "msg_1" {
		t.Errorf("expected first message ID msg_1, got %s", messages[0].Info.ID)
	}
	if messages[1].Info.ID != "msg_2" {
		t.Errorf("expected second message ID msg_2, got %s", messages[1].Info.ID)
	}
}

func TestSessionMessages_MissingID(t *testing.T) {
	client, err := NewClient(WithBaseURL("http://localhost"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Messages(context.Background(), "", nil)
	if err == nil {
		t.Fatal("expected error for missing ID, got nil")
	}
	if err.Error() != "missing required id parameter" {
		t.Errorf("expected 'missing required id parameter', got: %v", err)
	}
}

func TestSessionMessages_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Messages(context.Background(), "sess_1", nil)
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
