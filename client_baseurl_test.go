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

func TestBuildURL_BaseURLQueryMergedWithParams(t *testing.T) {
	var receivedQuery string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]opencode.Session{})
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL + "?token=abc"))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = client.Session.List(context.Background(), &opencode.SessionListParams{
		Directory: opencode.Ptr("/mydir"),
	})
	if err != nil {
		t.Fatalf("Session.List: %v", err)
	}

	if !strings.Contains(receivedQuery, "token=abc") {
		t.Errorf("expected base URL query param preserved, got %q", receivedQuery)
	}
	if !strings.Contains(receivedQuery, "directory=") {
		t.Errorf("expected params query param present, got %q", receivedQuery)
	}
}

func TestBuildURL_ParamsOverrideBaseURLQueryKey(t *testing.T) {
	var receivedQuery string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]opencode.Session{})
	}))
	defer server.Close()

	// Base URL has directory=old, params will set directory=/new
	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL + "?directory=old"))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = client.Session.List(context.Background(), &opencode.SessionListParams{
		Directory: opencode.Ptr("/new"),
	})
	if err != nil {
		t.Fatalf("Session.List: %v", err)
	}

	// Params struct should override the base URL's directory key
	if strings.Contains(receivedQuery, "directory=old") {
		t.Errorf("expected params to override base URL query key, got %q", receivedQuery)
	}
	if !strings.Contains(receivedQuery, "directory=") {
		t.Errorf("expected directory param present, got %q", receivedQuery)
	}
}

func TestBuildURL_NilParams(t *testing.T) {
	var receivedQuery string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]opencode.Session{})
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL + "?token=abc"))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = client.Session.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Session.List: %v", err)
	}

	if receivedQuery != "token=abc" {
		t.Errorf("expected only base URL query params, got %q", receivedQuery)
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
