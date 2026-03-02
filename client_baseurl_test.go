package opencode_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dominicnunez/opencode-sdk-go"
)

func TestNewClient_EnvBaseURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]opencode.Session{})
	}))
	defer server.Close()

	t.Setenv("OPENCODE_BASE_URL", server.URL)

	client, err := opencode.NewClient()
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = client.Session.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Session.List: %v", err)
	}
}

func TestNewClient_EnvBaseURL_InvalidURL(t *testing.T) {
	t.Setenv("OPENCODE_BASE_URL", "://bad")

	_, err := opencode.NewClient()
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestNewClient_WithBaseURLOverridesInvalidEnvBaseURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]opencode.Session{})
	}))
	defer server.Close()

	t.Setenv("OPENCODE_BASE_URL", "://bad")

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = client.Session.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Session.List: %v", err)
	}
}

func TestNewClient_EnvBaseURL_BadScheme(t *testing.T) {
	t.Setenv("OPENCODE_BASE_URL", "ftp://localhost:8080")

	_, err := opencode.NewClient()
	if err == nil {
		t.Fatal("expected error for non-http scheme")
	}
	if !strings.Contains(err.Error(), "http or https") {
		t.Errorf("expected scheme error, got: %v", err)
	}
}

func TestNewClient_EnvBaseURL_HTTPRemoteHostRejected(t *testing.T) {
	t.Setenv("OPENCODE_BASE_URL", "http://example.com:8080")

	_, err := opencode.NewClient()
	if err == nil {
		t.Fatal("expected error for insecure remote http base URL")
	}
	if !strings.Contains(err.Error(), "https for non-loopback hosts") {
		t.Errorf("expected insecure transport error, got: %v", err)
	}
}

func TestWithBaseURL_HTTPLoopbackHostsAllowed(t *testing.T) {
	tests := []string{
		"http://localhost:54321",
		"http://127.0.0.1:54321",
		"http://[::1]:54321",
	}

	for _, rawURL := range tests {
		t.Run(rawURL, func(t *testing.T) {
			_, err := opencode.NewClient(opencode.WithBaseURL(rawURL))
			if err != nil {
				t.Fatalf("expected loopback URL to be accepted, got: %v", err)
			}
		})
	}
}

func TestWithBaseURL_HTTPRemoteHostRejected(t *testing.T) {
	_, err := opencode.NewClient(opencode.WithBaseURL("http://example.com"))
	if err == nil {
		t.Fatal("expected error for insecure remote http base URL")
	}
	if !strings.Contains(err.Error(), "https for non-loopback hosts") {
		t.Errorf("expected insecure transport error, got: %v", err)
	}
}

func TestNewClient_EnvBaseURL_FallbackToDefault(t *testing.T) {
	t.Setenv("OPENCODE_BASE_URL", "")

	client, err := opencode.NewClient()
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	// Should have fallen back to DefaultBaseURL — just verify it was created
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestWithBaseURL_RejectsQueryParameters(t *testing.T) {
	tests := []string{
		"https://example.com?token=secret",
		"https://example.com?api_key=secret",
		"https://example.com?access_token=secret",
		"https://example.com?Access-Token=secret",
		"https://example.com?apikey=secret",
		"https://example.com?auth=secret",
		"https://example.com?key=secret",
		"https://example.com?bearer=secret",
		"https://example.com?authorization=secret",
		"https://example.com?workspace=abc",
		"https://example.com?directory=%2Ftmp",
	}

	for _, rawURL := range tests {
		t.Run(rawURL, func(t *testing.T) {
			_, err := opencode.NewClient(opencode.WithBaseURL(rawURL))
			if err == nil {
				t.Fatal("expected error for query parameter in base URL")
			}
			if !strings.Contains(err.Error(), "must not include query parameters") {
				t.Fatalf("expected query parameter error, got: %v", err)
			}
		})
	}
}

func TestWithBaseURL_RejectsUserInfoAndEmptyHost(t *testing.T) {
	tests := []struct {
		name   string
		rawURL string
		want   string
	}{
		{
			name:   "reject userinfo credentials",
			rawURL: "https://username@example.com",
			want:   "must not include user info",
		},
		{
			name:   "reject https URL with empty host",
			rawURL: "https:///api",
			want:   "must include a host",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := opencode.NewClient(opencode.WithBaseURL(tc.rawURL))
			if err == nil {
				t.Fatalf("expected error for invalid base URL: %s", tc.rawURL)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected error containing %q, got: %v", tc.want, err)
			}
		})
	}
}

func TestBuildURL_NoBaseQueryNoParams(t *testing.T) {
	var receivedQuery string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]opencode.Session{})
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = client.Session.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Session.List: %v", err)
	}

	if receivedQuery != "" {
		t.Errorf("expected empty query, got %q", receivedQuery)
	}
}
