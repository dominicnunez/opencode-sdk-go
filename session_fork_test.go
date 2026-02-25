package opencode

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionFork_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/session/ses_123/fork" {
			t.Errorf("expected /session/ses_123/fork, got %s", r.URL.Path)
		}

		// Verify request body
		body, _ := io.ReadAll(r.Body)
		var params SessionForkParams
		if err := json.Unmarshal(body, &params); err != nil {
			t.Errorf("failed to unmarshal request body: %v", err)
		}
		if params.MessageID != "msg_456" {
			t.Errorf("expected messageID msg_456, got %s", params.MessageID)
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Session{
			ID:        "ses_789",
			Directory: "/test",
			ProjectID: "proj_1",
			ParentID:  "ses_123",
			Title:     "Forked Session",
			Version:   "1.0.0",
			Time: SessionTime{
				Created: 1234567890,
				Updated: 1234567890,
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	params := &SessionForkParams{
		MessageID: "msg_456",
	}
	session, err := client.Session.Fork(context.Background(), "ses_123", params)
	if err != nil {
		t.Fatalf("Fork() failed: %v", err)
	}

	if session.ID != "ses_789" {
		t.Errorf("expected ID ses_789, got %s", session.ID)
	}
	if session.ParentID != "ses_123" {
		t.Errorf("expected ParentID ses_123, got %s", session.ParentID)
	}
	if session.Title != "Forked Session" {
		t.Errorf("expected Title 'Forked Session', got %s", session.Title)
	}
}

func TestSessionFork_WithDirectory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameter
		if dir := r.URL.Query().Get("directory"); dir != "/custom/path" {
			t.Errorf("expected directory /custom/path, got %s", dir)
		}

		// Verify request body
		body, _ := io.ReadAll(r.Body)
		var params SessionForkParams
		if err := json.Unmarshal(body, &params); err != nil {
			t.Errorf("failed to unmarshal request body: %v", err)
		}
		if params.MessageID != "msg_789" {
			t.Errorf("expected messageID msg_789, got %s", params.MessageID)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Session{
			ID:        "ses_999",
			Directory: "/custom/path",
			ProjectID: "proj_2",
			Title:     "Custom Directory Fork",
			Version:   "1.0.0",
			Time: SessionTime{
				Created: 1234567890,
				Updated: 1234567890,
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dir := "/custom/path"
	params := &SessionForkParams{
		MessageID: "msg_789",
		Directory: &dir,
	}
	session, err := client.Session.Fork(context.Background(), "ses_456", params)
	if err != nil {
		t.Fatalf("Fork() failed: %v", err)
	}

	if session.Directory != "/custom/path" {
		t.Errorf("expected Directory /custom/path, got %s", session.Directory)
	}
}

func TestSessionFork_MissingID(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	params := &SessionForkParams{
		MessageID: "msg_123",
	}
	_, err = client.Session.Fork(context.Background(), "", params)
	if err == nil {
		t.Error("expected error for missing ID, got nil")
	}
	if err.Error() != "missing required id parameter" {
		t.Errorf("expected 'missing required id parameter', got %s", err.Error())
	}
}

func TestSessionFork_MissingParams(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Fork(context.Background(), "ses_123", nil)
	if err == nil {
		t.Error("expected error for missing params, got nil")
	}
	if err.Error() != "missing required params" {
		t.Errorf("expected 'missing required params', got %s", err.Error())
	}
}

func TestSessionFork_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "server error"}`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	params := &SessionForkParams{
		MessageID: "msg_123",
	}
	_, err = client.Session.Fork(context.Background(), "ses_123", params)
	if err == nil {
		t.Error("expected error for server error, got nil")
	}
}

func TestSessionFork_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "ses_123", "invalid json`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	params := &SessionForkParams{
		MessageID: "msg_123",
	}
	_, err = client.Session.Fork(context.Background(), "ses_123", params)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestSessionForkParams_Marshal(t *testing.T) {
	params := SessionForkParams{
		MessageID: "msg_123",
	}
	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("failed to marshal SessionForkParams: %v", err)
	}

	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if unmarshaled["messageID"] != "msg_123" {
		t.Errorf("expected messageID msg_123, got %v", unmarshaled["messageID"])
	}
}

func TestSessionForkParams_MarshalWithDirectory(t *testing.T) {
	dir := "/test/path"
	params := SessionForkParams{
		MessageID: "msg_123",
		Directory: &dir,
	}
	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("failed to marshal SessionForkParams: %v", err)
	}

	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if unmarshaled["messageID"] != "msg_123" {
		t.Errorf("expected messageID msg_123, got %v", unmarshaled["messageID"])
	}
	// Directory is a query param, not in JSON body
	if _, exists := unmarshaled["directory"]; exists {
		t.Error("directory should not be in JSON body")
	}
}

func TestSessionForkParams_URLQuery(t *testing.T) {
	dir := "/test/path"
	params := SessionForkParams{
		MessageID: "msg_123",
		Directory: &dir,
	}
	values, err := params.URLQuery()
	if err != nil {
		t.Fatalf("URLQuery() failed: %v", err)
	}

	if values.Get("directory") != "/test/path" {
		t.Errorf("expected directory /test/path, got %s", values.Get("directory"))
	}
}

func TestSessionForkParams_URLQueryWithoutDirectory(t *testing.T) {
	params := SessionForkParams{
		MessageID: "msg_123",
	}
	values, err := params.URLQuery()
	if err != nil {
		t.Fatalf("URLQuery() failed: %v", err)
	}

	if dir := values.Get("directory"); dir != "" {
		t.Errorf("expected empty directory, got %s", dir)
	}
}
