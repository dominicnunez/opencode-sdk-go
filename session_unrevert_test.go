package opencode

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSessionUnrevert_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/session/sess_123/unrevert" {
			t.Errorf("expected /session/sess_123/unrevert, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":        "sess_123",
			"directory": "/test/path",
			"projectID": "proj_456",
			"title":     "Test Session",
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

	session, err := client.Session.Unrevert(context.Background(), "sess_123", &SessionUnrevertParams{})
	if err != nil {
		t.Fatalf("Unrevert failed: %v", err)
	}

	if session.ID != "sess_123" {
		t.Errorf("expected session ID sess_123, got %s", session.ID)
	}
	if session.Directory != "/test/path" {
		t.Errorf("expected directory /test/path, got %s", session.Directory)
	}
	if session.ProjectID != "proj_456" {
		t.Errorf("expected project ID proj_456, got %s", session.ProjectID)
	}
}

func TestSessionUnrevert_WithDirectoryQueryParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		query := r.URL.Query()
		if query.Get("directory") != "/custom/dir" {
			t.Errorf("expected directory query param /custom/dir, got %s", query.Get("directory"))
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":        "sess_123",
			"directory": "/custom/dir",
			"projectID": "proj_456",
			"title":     "Test Session",
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

	dir := "/custom/dir"
	session, err := client.Session.Unrevert(context.Background(), "sess_123", &SessionUnrevertParams{
		Directory: &dir,
	})
	if err != nil {
		t.Fatalf("Unrevert failed: %v", err)
	}

	if session.ID != "sess_123" {
		t.Errorf("expected session ID sess_123, got %s", session.ID)
	}
	if session.Directory != "/custom/dir" {
		t.Errorf("expected directory /custom/dir, got %s", session.Directory)
	}
}

func TestSessionUnrevert_NilParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id": "sess_123",
		})
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	session, err := client.Session.Unrevert(context.Background(), "sess_123", nil)
	if err != nil {
		t.Fatalf("Unrevert failed: %v", err)
	}

	if session.ID != "sess_123" {
		t.Errorf("expected session ID sess_123, got %s", session.ID)
	}
}

func TestSessionUnrevert_MissingID(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Unrevert(context.Background(), "", &SessionUnrevertParams{})
	if err == nil {
		t.Fatal("expected error for missing ID, got nil")
	}
	if err.Error() != "missing required id parameter" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestSessionUnrevert_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Unrevert(context.Background(), "sess_123", &SessionUnrevertParams{})
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

func TestSessionUnrevert_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id": "sess_123", invalid json`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Unrevert(context.Background(), "sess_123", &SessionUnrevertParams{})
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestSessionUnrevertParams_URLQuery(t *testing.T) {
	tests := []struct {
		name   string
		params SessionUnrevertParams
		want   url.Values
	}{
		{
			name:   "no directory",
			params: SessionUnrevertParams{},
			want:   url.Values{},
		},
		{
			name: "with directory",
			params: SessionUnrevertParams{
				Directory: ptrString("/custom/dir"),
			},
			want: url.Values{"directory": []string{"/custom/dir"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.params.URLQuery()
			if err != nil {
				t.Fatalf("URLQuery failed: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Errorf("URLQuery() length = %d, want %d", len(got), len(tt.want))
			}
			for key, wantValues := range tt.want {
				gotValues := got[key]
				if len(gotValues) != len(wantValues) {
					t.Errorf("URLQuery()[%s] length = %d, want %d", key, len(gotValues), len(wantValues))
				}
				for i, wantValue := range wantValues {
					if i < len(gotValues) && gotValues[i] != wantValue {
						t.Errorf("URLQuery()[%s][%d] = %s, want %s", key, i, gotValues[i], wantValue)
					}
				}
			}
		})
	}
}

func ptrString(s string) *string {
	return &s
}
