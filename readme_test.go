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
		// Mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/session" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode([]Session{
					{ID: "sess_1", Title: "Test Session"},
				})
			}
		}))
		defer server.Close()

		// Create client with mock server
		client, err := NewClient(
			WithBaseURL(server.URL),
			WithTimeout(30*time.Second),
		)
		if err != nil {
			t.Fatal(err)
		}

		// List sessions
		sessions, err := client.Session.List(context.TODO(), &SessionListParams{})
		if err != nil {
			t.Fatal(err)
		}

		if len(sessions) != 1 {
			t.Errorf("expected 1 session, got %d", len(sessions))
		}
	})

	t.Run("RequestParameters", func(t *testing.T) {
		// Verify direct types for required fields
		params := SessionCommandParams{
			Command: "list-files",
		}
		if params.Command != "list-files" {
			t.Errorf("expected 'list-files', got %s", params.Command)
		}

		// Verify pointer types for optional fields
		params.Directory = Ptr("./src")
		if params.Directory == nil || *params.Directory != "./src" {
			t.Error("expected directory to be './src'")
		}

		// Verify nil optional fields
		params.Agent = nil
		if params.Agent != nil {
			t.Error("expected agent to be nil")
		}
	})

	t.Run("UnionTypes", func(t *testing.T) {
		// Test Message union (UserMessage)
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

		// Test Part union (TextPart)
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

	t.Run("ClientConfiguration", func(t *testing.T) {
		customClient := &http.Client{Timeout: 10 * time.Second}

		client, err := NewClient(
			WithBaseURL("https://api.example.com"),
			WithTimeout(60*time.Second),
			WithMaxRetries(5),
			WithHTTPClient(customClient),
		)
		if err != nil {
			t.Fatal(err)
		}

		if client == nil {
			t.Error("expected non-nil client")
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

		// OAuth authentication
		success, err := client.Auth.Set(context.TODO(), "provider-id", &AuthSetParams{
			Auth: OAuth{
				Type:    AuthTypeOAuth,
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
		// Test that EventListStreaming returns a stream
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/event" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "text/event-stream")
				w.Header().Set("Cache-Control", "no-cache")
				w.Header().Set("Connection", "keep-alive")

				// Send a test event
				_, _ = w.Write([]byte("data: {\"type\":\"message.updated\",\"data\":{\"info\":{\"id\":\"msg_1\",\"sessionID\":\"sess_1\",\"role\":\"user\",\"parts\":[]}}}\n\n"))

				// Close connection
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

		if stream == nil {
			t.Error("expected non-nil stream")
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

		// Get tool IDs
		toolIDs, err := client.Tool.IDs(context.TODO(), &ToolIDsParams{})
		if err != nil {
			t.Fatal(err)
		}
		if len(*toolIDs) != 2 {
			t.Errorf("expected 2 tool IDs, got %d", len(*toolIDs))
		}

		// Get tool schemas
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
		customClient := &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
			Timeout: 30 * time.Second,
		}

		client, err := NewClient(
			WithHTTPClient(customClient),
		)
		if err != nil {
			t.Fatal(err)
		}

		if client == nil {
			t.Error("expected non-nil client")
		}
	})
}

// TestREADMELoggingTransport verifies the logging transport example compiles
func TestREADMELoggingTransport(t *testing.T) {
	type LoggingTransport struct {
		Base http.RoundTripper
	}

	roundTripFunc := func(t *LoggingTransport, req *http.Request) (*http.Response, error) {
		// Simplified version without actual logging for test
		return t.Base.RoundTrip(req)
	}

	transport := &LoggingTransport{Base: http.DefaultTransport}

	// Verify RoundTrip compiles
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	_, _ = roundTripFunc(transport, req)
}

// TestREADMEPtrHelper verifies the Ptr helper works as documented
func TestREADMEPtrHelper(t *testing.T) {
	// Test string pointer
	strPtr := Ptr("test")
	if strPtr == nil || *strPtr != "test" {
		t.Error("Ptr[string] failed")
	}

	// Test int64 pointer
	intPtr := Ptr(int64(42))
	if intPtr == nil || *intPtr != 42 {
		t.Error("Ptr[int64] failed")
	}

	// Test bool pointer
	boolPtr := Ptr(true)
	if boolPtr == nil || *boolPtr != true {
		t.Error("Ptr[bool] failed")
	}
}
