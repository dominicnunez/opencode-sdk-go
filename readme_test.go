package opencode

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestREADMEExamples verifies that the examples in README.md are correct and compile
func TestREADMEExamples(t *testing.T) {
	t.Run("BasicUsage", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/session" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode([]Session{
					{ID: "sess_1", Title: "Test Session"},
				})
			}
		}))
		defer server.Close()

		client, err := NewClient(
			WithBaseURL(server.URL),
			WithTimeout(30*time.Second),
		)
		if err != nil {
			t.Fatal(err)
		}

		sessions, err := client.Session.List(context.TODO(), &SessionListParams{})
		if err != nil {
			t.Fatal(err)
		}

		if len(sessions) != 1 {
			t.Errorf("expected 1 session, got %d", len(sessions))
		}
	})

	t.Run("UnionTypes", func(t *testing.T) {
		msgJSON := `{"id":"msg_1","sessionID":"sess_1","role":"user","parts":[]}`
		var msg Message
		if err := json.Unmarshal([]byte(msgJSON), &msg); err != nil {
			t.Fatal(err)
		}

		user, err := msg.AsUser()
		if err != nil {
			t.Fatalf("AsUser error: %v", err)
		}
		if user == nil {
			t.Fatal("expected AsUser to succeed")
		}
		if user.Role != UserMessageRoleUser {
			t.Error("expected user message role")
		}

		partJSON := `{"id":"part_1","messageID":"msg_1","sessionID":"sess_1","type":"text","text":"Hello"}`
		var part Part
		if err := json.Unmarshal([]byte(partJSON), &part); err != nil {
			t.Fatal(err)
		}

		textPart, err := part.AsText()
		if err != nil {
			t.Fatalf("AsText error: %v", err)
		}
		if textPart == nil {
			t.Fatal("expected AsText to succeed")
		}
		if textPart.Text != "Hello" {
			t.Errorf("expected 'Hello', got %s", textPart.Text)
		}
	})

	t.Run("Authentication", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/auth/provider-id" && r.Method == http.MethodPut {
				var body map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if body["type"] != string(AuthTypeOAuth) {
					t.Fatalf("expected oauth type, got %v", body["type"])
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte("true"))
			}
		}))
		defer server.Close()

		client, err := NewClient(WithBaseURL(server.URL))
		if err != nil {
			t.Fatal(err)
		}

		success, err := client.Auth.Set(context.TODO(), "provider-id", &AuthSetParams{
			Auth: OAuth{
				Refresh: "refresh_token",
				Access:  "access_token",
				Expires: 1234567890,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		if !success {
			t.Error("expected success=true")
		}
	})

	t.Run("StreamingEvents", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/event" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "text/event-stream")
				w.Header().Set("Cache-Control", "no-cache")
				w.Header().Set("Connection", "keep-alive")

				_, _ = w.Write([]byte("data: {\"type\":\"message.updated\",\"properties\":{\"info\":{\"id\":\"msg_1\",\"sessionID\":\"sess_1\",\"role\":\"user\",\"parts\":[]}}}\n\n"))

				flusher, ok := w.(http.Flusher)
				if ok {
					flusher.Flush()
				}
			}
		}))
		defer server.Close()

		client, err := NewClient(WithBaseURL(server.URL))
		if err != nil {
			t.Fatal(err)
		}

		stream := client.Event.ListStreaming(context.TODO(), &EventListParams{})
		defer func() { _ = stream.Close() }()

		if !stream.Next() {
			t.Fatalf("expected at least one event, got err: %v", stream.Err())
		}
		evt := stream.Current()
		if evt.Type != EventTypeMessageUpdated {
			t.Errorf("expected event type %q, got %q", EventTypeMessageUpdated, evt.Type)
		}

		updated, err := evt.AsMessageUpdated()
		if err != nil {
			t.Fatalf("AsMessageUpdated: %v", err)
		}
		if updated.Data.Info.ID != "msg_1" {
			t.Errorf("expected message ID %q, got %q", "msg_1", updated.Data.Info.ID)
		}
	})

	t.Run("ToolsAPI", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/experimental/tool/ids" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode([]string{"tool_1", "tool_2"})
			} else if r.URL.Path == "/experimental/tool" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode([]ToolListItem{
					{
						ID:          "tool_1",
						Description: "Test tool",
						Parameters:  json.RawMessage(`{"type":"object"}`),
					},
				})
			}
		}))
		defer server.Close()

		client, err := NewClient(WithBaseURL(server.URL))
		if err != nil {
			t.Fatal(err)
		}

		toolIDs, err := client.Tool.IDs(context.TODO(), &ToolIDsParams{})
		if err != nil {
			t.Fatal(err)
		}
		if len(*toolIDs) != 2 {
			t.Errorf("expected 2 tool IDs, got %d", len(*toolIDs))
		}

		tools, err := client.Tool.List(context.TODO(), &ToolListParams{
			Provider: "anthropic",
			Model:    "claude-sonnet-4-5",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(*tools) != 1 {
			t.Errorf("expected 1 tool, got %d", len(*tools))
		}
	})

	t.Run("MCPStatus", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/mcp" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"server1": map[string]interface{}{"status": "running"},
				})
			}
		}))
		defer server.Close()

		client, err := NewClient(WithBaseURL(server.URL))
		if err != nil {
			t.Fatal(err)
		}

		status, err := client.Mcp.Status(context.TODO(), &McpStatusParams{})
		if err != nil {
			t.Fatal(err)
		}

		if len(*status) != 1 {
			t.Errorf("expected 1 server, got %d", len(*status))
		}
	})

	t.Run("CustomHTTPClient", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]Session{})
		}))
		defer server.Close()

		customClient := &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
			Timeout: 30 * time.Second,
		}

		client, err := NewClient(
			WithBaseURL(server.URL),
			WithHTTPClient(customClient),
		)
		if err != nil {
			t.Fatal(err)
		}

		sessions, err := client.Session.List(context.TODO(), nil)
		if err != nil {
			t.Fatalf("expected custom HTTP client to work, got: %v", err)
		}
		if sessions == nil {
			t.Error("expected non-nil response")
		}
	})
}
