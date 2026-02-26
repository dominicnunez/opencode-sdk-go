package opencode

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestSessionGet_Success verifies Session.Get returns session data
func TestSessionGet_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/session/sess_123" {
			t.Errorf("Expected path /session/sess_123, got %s", r.URL.Path)
		}

		response := Session{
			ID:        "sess_123",
			Title:     "Test Session",
			ProjectID: "proj_456",
			Directory: "/test/dir",
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	session, err := client.Session.Get(context.Background(), "sess_123", nil)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if session.ID != "sess_123" {
		t.Errorf("Expected session ID sess_123, got %s", session.ID)
	}
	if session.Title != "Test Session" {
		t.Errorf("Expected title 'Test Session', got %s", session.Title)
	}
	if session.ProjectID != "proj_456" {
		t.Errorf("Expected project ID proj_456, got %s", session.ProjectID)
	}
	if session.Directory != "/test/dir" {
		t.Errorf("Expected directory /test/dir, got %s", session.Directory)
	}
}

// TestSessionGet_WithDirectoryParam verifies query params are passed correctly
func TestSessionGet_WithDirectoryParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryDir := r.URL.Query().Get("directory")
		if queryDir != "/custom/dir" {
			t.Errorf("Expected directory query param /custom/dir, got %s", queryDir)
		}

		response := Session{
			ID:        "sess_123",
			Directory: queryDir,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	dir := "/custom/dir"
	session, err := client.Session.Get(context.Background(), "sess_123", &SessionGetParams{
		Directory: &dir,
	})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if session.Directory != "/custom/dir" {
		t.Errorf("Expected directory /custom/dir, got %s", session.Directory)
	}
}

// TestSessionGet_MissingID verifies error is returned when ID is missing
func TestSessionGet_MissingID(t *testing.T) {
	client, err := NewClient(WithBaseURL("http://localhost"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Session.Get(context.Background(), "", nil)
	if err == nil {
		t.Fatal("Expected error for missing ID, got nil")
	}
	if err.Error() != "missing required id parameter" {
		t.Errorf("Expected 'missing required id parameter' error, got: %v", err)
	}
}

// TestSessionGet_ServerError verifies error handling for server errors
func TestSessionGet_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Session.Get(context.Background(), "sess_123", nil)
	if err == nil {
		t.Fatal("Expected error for server error, got nil")
	}
}

// TestSessionUpdate_Success verifies Session.Update updates session and returns updated data
func TestSessionUpdate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}
		if r.URL.Path != "/session/sess_123" {
			t.Errorf("Expected path /session/sess_123, got %s", r.URL.Path)
		}

		// Verify request body
		body, _ := io.ReadAll(r.Body)
		var params SessionUpdateParams
		if err := json.Unmarshal(body, &params); err != nil {
			t.Fatalf("Failed to unmarshal request body: %v", err)
		}

		if params.Title == nil || *params.Title != "Updated Title" {
			t.Errorf("Expected title 'Updated Title', got %v", params.Title)
		}

		response := Session{
			ID:    "sess_123",
			Title: "Updated Title",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	newTitle := "Updated Title"
	session, err := client.Session.Update(context.Background(), "sess_123", &SessionUpdateParams{
		Title: &newTitle,
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if session.ID != "sess_123" {
		t.Errorf("Expected session ID sess_123, got %s", session.ID)
	}
	if session.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got %s", session.Title)
	}
}

// TestSessionUpdate_WithDirectoryParam verifies query params and body params work together
func TestSessionUpdate_WithDirectoryParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryDir := r.URL.Query().Get("directory")
		if queryDir != "/custom/dir" {
			t.Errorf("Expected directory query param /custom/dir, got %s", queryDir)
		}

		body, _ := io.ReadAll(r.Body)
		var params SessionUpdateParams
		_ = json.Unmarshal(body, &params)

		if params.Title == nil || *params.Title != "New Title" {
			t.Errorf("Expected title 'New Title', got %v", params.Title)
		}

		response := Session{
			ID:        "sess_123",
			Title:     *params.Title,
			Directory: queryDir,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	dir := "/custom/dir"
	title := "New Title"
	session, err := client.Session.Update(context.Background(), "sess_123", &SessionUpdateParams{
		Directory: &dir,
		Title:     &title,
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if session.Title != "New Title" {
		t.Errorf("Expected title 'New Title', got %s", session.Title)
	}
	if session.Directory != "/custom/dir" {
		t.Errorf("Expected directory /custom/dir, got %s", session.Directory)
	}
}

// TestSessionUpdate_MissingID verifies error is returned when ID is missing
func TestSessionUpdate_MissingID(t *testing.T) {
	client, err := NewClient(WithBaseURL("http://localhost"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Session.Update(context.Background(), "", &SessionUpdateParams{})
	if err == nil {
		t.Fatal("Expected error for missing ID, got nil")
	}
	if err.Error() != "missing required id parameter" {
		t.Errorf("Expected 'missing required id parameter' error, got: %v", err)
	}
}

// TestSessionUpdate_NilParams verifies nil params are handled gracefully
func TestSessionUpdate_NilParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := Session{
			ID:    "sess_123",
			Title: "Unchanged",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	session, err := client.Session.Update(context.Background(), "sess_123", nil)
	if err != nil {
		t.Fatalf("Update with nil params failed: %v", err)
	}

	if session.ID != "sess_123" {
		t.Errorf("Expected session ID sess_123, got %s", session.ID)
	}
}

// TestSessionDelete_Success verifies Session.Delete deletes session successfully
func TestSessionDelete_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/session/sess_123" {
			t.Errorf("Expected path /session/sess_123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.Session.Delete(context.Background(), "sess_123", nil)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

// TestSessionDelete_WithDirectoryParam verifies query params are passed correctly
func TestSessionDelete_WithDirectoryParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryDir := r.URL.Query().Get("directory")
		if queryDir != "/custom/dir" {
			t.Errorf("Expected directory query param /custom/dir, got %s", queryDir)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	dir := "/custom/dir"
	err = client.Session.Delete(context.Background(), "sess_123", &SessionDeleteParams{
		Directory: &dir,
	})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

// TestSessionDelete_MissingID verifies error is returned when ID is missing
func TestSessionDelete_MissingID(t *testing.T) {
	client, err := NewClient(WithBaseURL("http://localhost"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.Session.Delete(context.Background(), "", nil)
	if err == nil {
		t.Fatal("Expected error for missing ID, got nil")
	}
	if err.Error() != "missing required id parameter" {
		t.Errorf("Expected 'missing required id parameter' error, got: %v", err)
	}
}

// TestSessionDelete_ServerError verifies error handling for server errors
func TestSessionDelete_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.Session.Delete(context.Background(), "sess_123", nil)
	if err == nil {
		t.Fatal("Expected error for server error, got nil")
	}
}

// TestSessionUpdateParams_URLQuery verifies query param serialization
func TestSessionUpdateParams_URLQuery(t *testing.T) {
	tests := []struct {
		name   string
		params SessionUpdateParams
		expect string
	}{
		{
			name:   "No directory",
			params: SessionUpdateParams{},
			expect: "",
		},
		{
			name: "With directory",
			params: SessionUpdateParams{
				Directory: Ptr("/test/dir"),
			},
			expect: "directory=%2Ftest%2Fdir",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := tt.params.URLQuery()
			if err != nil {
				t.Fatalf("URLQuery failed: %v", err)
			}

			encoded := query.Encode()
			if encoded != tt.expect {
				t.Errorf("Expected query %q, got %q", tt.expect, encoded)
			}
		})
	}
}

// TestSessionDeleteParams_URLQuery verifies query param serialization
func TestSessionDeleteParams_URLQuery(t *testing.T) {
	tests := []struct {
		name   string
		params SessionDeleteParams
		expect string
	}{
		{
			name:   "No directory",
			params: SessionDeleteParams{},
			expect: "",
		},
		{
			name: "With directory",
			params: SessionDeleteParams{
				Directory: Ptr("/test/dir"),
			},
			expect: "directory=%2Ftest%2Fdir",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := tt.params.URLQuery()
			if err != nil {
				t.Fatalf("URLQuery failed: %v", err)
			}

			encoded := query.Encode()
			if encoded != tt.expect {
				t.Errorf("Expected query %q, got %q", tt.expect, encoded)
			}
		})
	}
}

// TestSessionGetParams_URLQuery verifies query param serialization
func TestSessionGetParams_URLQuery(t *testing.T) {
	tests := []struct {
		name   string
		params SessionGetParams
		expect string
	}{
		{
			name:   "No directory",
			params: SessionGetParams{},
			expect: "",
		},
		{
			name: "With directory",
			params: SessionGetParams{
				Directory: Ptr("/test/dir"),
			},
			expect: "directory=%2Ftest%2Fdir",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := tt.params.URLQuery()
			if err != nil {
				t.Fatalf("URLQuery failed: %v", err)
			}

			encoded := query.Encode()
			if encoded != tt.expect {
				t.Errorf("Expected query %q, got %q", tt.expect, encoded)
			}
		})
	}
}
