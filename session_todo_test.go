package opencode_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dominicnunez/opencode-sdk-go"
)

func TestSessionService_Todo_Success(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if r.URL.Path != "/session/ses_123/todo" {
			t.Errorf("expected path /session/ses_123/todo, got %s", r.URL.Path)
		}

		// Return sample Todo array
		todos := []opencode.Todo{
			{
				ID:       "todo_1",
				Content:  "Complete feature X",
				Priority: "high",
				Status:   "in_progress",
			},
			{
				ID:       "todo_2",
				Content:  "Write tests",
				Priority: "medium",
				Status:   "pending",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(todos)
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	todos, err := client.Session.Todo(context.Background(), "ses_123", &opencode.SessionTodoParams{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(todos) != 2 {
		t.Fatalf("expected 2 todos, got %d", len(todos))
	}
	if todos[0].ID != "todo_1" {
		t.Errorf("expected first todo ID to be 'todo_1', got %s", todos[0].ID)
	}
	if todos[0].Content != "Complete feature X" {
		t.Errorf("expected first todo content to be 'Complete feature X', got %s", todos[0].Content)
	}
	if todos[0].Priority != "high" {
		t.Errorf("expected first todo priority to be 'high', got %s", todos[0].Priority)
	}
	if todos[0].Status != "in_progress" {
		t.Errorf("expected first todo status to be 'in_progress', got %s", todos[0].Status)
	}
	if todos[1].ID != "todo_2" {
		t.Errorf("expected second todo ID to be 'todo_2', got %s", todos[1].ID)
	}
}

func TestSessionService_Todo_WithDirectory(t *testing.T) {
	// Mock server that verifies directory query param
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("directory") != "/path/to/dir" {
			t.Errorf("expected directory query param '/path/to/dir', got %s", r.URL.Query().Get("directory"))
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]opencode.Todo{})
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dir := "/path/to/dir"
	_, err = client.Session.Todo(context.Background(), "ses_123", &opencode.SessionTodoParams{
		Directory: &dir,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSessionService_Todo_EmptyArray(t *testing.T) {
	// Mock server that returns empty array
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	todos, err := client.Session.Todo(context.Background(), "ses_123", &opencode.SessionTodoParams{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(todos) != 0 {
		t.Errorf("expected 0 todos, got %d", len(todos))
	}
}

func TestSessionService_Todo_MissingID(t *testing.T) {
	client, err := opencode.NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Todo(context.Background(), "", &opencode.SessionTodoParams{})
	if err == nil {
		t.Fatal("expected error for missing id, got nil")
	}
	if err.Error() != "missing required id parameter" {
		t.Errorf("expected 'missing required id parameter' error, got %v", err)
	}
}

func TestSessionService_Todo_NilParams(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]opencode.Todo{})
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// nil params should be handled gracefully
	_, err = client.Session.Todo(context.Background(), "ses_123", nil)
	if err != nil {
		t.Fatalf("expected no error with nil params, got %v", err)
	}
}

func TestSessionService_Todo_ServerError(t *testing.T) {
	// Mock server that returns 500 error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Todo(context.Background(), "ses_123", &opencode.SessionTodoParams{})
	if err == nil {
		t.Fatal("expected error for server error, got nil")
	}
	var apiErr *opencode.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, apiErr.StatusCode)
	}
}

func TestSessionService_Todo_InvalidJSON(t *testing.T) {
	// Mock server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{invalid json`))
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Todo(context.Background(), "ses_123", &opencode.SessionTodoParams{})
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestTodo_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"id": "todo_abc",
		"content": "Test task",
		"priority": "low",
		"status": "completed"
	}`

	var todo opencode.Todo
	err := json.Unmarshal([]byte(jsonData), &todo)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if todo.ID != "todo_abc" {
		t.Errorf("expected ID 'todo_abc', got %s", todo.ID)
	}
	if todo.Content != "Test task" {
		t.Errorf("expected content 'Test task', got %s", todo.Content)
	}
	if todo.Priority != "low" {
		t.Errorf("expected priority 'low', got %s", todo.Priority)
	}
	if todo.Status != "completed" {
		t.Errorf("expected status 'completed', got %s", todo.Status)
	}
}

func TestSessionTodoParams_URLQuery(t *testing.T) {
	tests := []struct {
		name     string
		params   opencode.SessionTodoParams
		expected map[string]string
	}{
		{
			name:     "no params",
			params:   opencode.SessionTodoParams{},
			expected: map[string]string{},
		},
		{
			name: "with directory",
			params: opencode.SessionTodoParams{
				Directory: ptrString("/my/dir"),
			},
			expected: map[string]string{
				"directory": "/my/dir",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := tt.params.URLQuery()
			if err != nil {
				t.Fatalf("URLQuery failed: %v", err)
			}

			for key, expectedValue := range tt.expected {
				actualValue := query.Get(key)
				if actualValue != expectedValue {
					t.Errorf("expected %s=%s, got %s", key, expectedValue, actualValue)
				}
			}

			// Verify no extra keys
			if len(query) != len(tt.expected) {
				t.Errorf("expected %d query params, got %d", len(tt.expected), len(query))
			}
		})
	}
}

func ptrString(s string) *string {
	return &s
}
