package opencode_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dominicnunez/opencode-sdk-go"
)

func TestSessionUnshare_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE method, got %s", r.Method)
		}
		if r.URL.Path != "/session/ses_123/share" {
			t.Errorf("expected path /session/ses_123/share, got %s", r.URL.Path)
		}

		response := opencode.Session{
			ID:        "ses_123",
			Directory: "/test",
			ProjectID: "prj_123",
			Time: opencode.SessionTime{
				Created: 1234567890,
				Updated: 1234567890,
			},
			Title:   "Test Session",
			Version: "1",
			// Share is omitted when unshared
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	session, err := client.Session.Unshare(context.Background(), "ses_123", nil)
	if err != nil {
		t.Fatalf("Unshare failed: %v", err)
	}

	if session.ID != "ses_123" {
		t.Errorf("expected session ID ses_123, got %s", session.ID)
	}
	// Verify Share is nil after unsharing
	if session.Share != nil {
		t.Errorf("expected Share to be nil after unsharing, got %+v", session.Share)
	}
}

func TestSessionUnshare_WithDirectory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE method, got %s", r.Method)
		}
		if r.URL.Path != "/session/ses_456/share" {
			t.Errorf("expected path /session/ses_456/share, got %s", r.URL.Path)
		}

		// Verify directory query param
		query := r.URL.Query()
		if query.Get("directory") != "/home/test" {
			t.Errorf("expected directory query param /home/test, got %s", query.Get("directory"))
		}

		response := opencode.Session{
			ID:        "ses_456",
			Directory: "/home/test",
			ProjectID: "prj_456",
			Time: opencode.SessionTime{
				Created: 1234567890,
				Updated: 1234567890,
			},
			Title:   "Test Session",
			Version: "1",
			// Share is omitted when unshared
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	dir := "/home/test"
	params := &opencode.SessionUnshareParams{
		Directory: &dir,
	}

	session, err := client.Session.Unshare(context.Background(), "ses_456", params)
	if err != nil {
		t.Fatalf("Unshare failed: %v", err)
	}

	if session.ID != "ses_456" {
		t.Errorf("expected session ID ses_456, got %s", session.ID)
	}
	if session.Directory != "/home/test" {
		t.Errorf("expected directory /home/test, got %s", session.Directory)
	}
}

func TestSessionUnshare_MissingID(t *testing.T) {
	client, err := opencode.NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Unshare(context.Background(), "", nil)
	if err == nil {
		t.Fatal("expected error for missing ID, got nil")
	}
	if err.Error() != "missing required id parameter" {
		t.Errorf("expected 'missing required id parameter' error, got %s", err.Error())
	}
}

func TestSessionUnshare_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "session not found"}`))
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Unshare(context.Background(), "ses_999", nil)
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

func TestSessionUnshare_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{invalid json`))
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Unshare(context.Background(), "ses_123", nil)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestSessionUnshareParams_URLQuery(t *testing.T) {
	tests := []struct {
		name     string
		params   opencode.SessionUnshareParams
		expected url.Values
	}{
		{
			name: "with directory",
			params: opencode.SessionUnshareParams{
				Directory: ptr("/home/user/project"),
			},
			expected: url.Values{"directory": []string{"/home/user/project"}},
		},
		{
			name:     "without directory",
			params:   opencode.SessionUnshareParams{},
			expected: url.Values{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, err := tt.params.URLQuery()
			if err != nil {
				t.Fatalf("URLQuery failed: %v", err)
			}

			if len(values) != len(tt.expected) {
				t.Errorf("expected %d query params, got %d", len(tt.expected), len(values))
			}

			for key, expectedVals := range tt.expected {
				gotVals := values[key]
				if len(gotVals) != len(expectedVals) {
					t.Errorf("for key %s: expected %d values, got %d", key, len(expectedVals), len(gotVals))
					continue
				}
				for i := range expectedVals {
					if gotVals[i] != expectedVals[i] {
						t.Errorf("for key %s[%d]: expected %s, got %s", key, i, expectedVals[i], gotVals[i])
					}
				}
			}
		})
	}
}

// Helper function for pointer to string
func ptr(s string) *string {
	return &s
}
