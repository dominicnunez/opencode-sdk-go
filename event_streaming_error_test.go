package opencode_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dominicnunez/opencode-sdk-go"
)

func TestListStreaming_JSONErrorBody(t *testing.T) {
	const jsonBody = `{"error": "rate limited", "retryAfter": 30}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Request-Id", "req-abc-123")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(jsonBody))
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	stream := client.Event.ListStreaming(context.Background(), nil)
	if stream.Next() {
		t.Fatal("expected Next() to return false on error status")
	}

	err = stream.Err()
	if err == nil {
		t.Fatal("expected non-nil error from stream")
	}

	var apiErr *opencode.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *opencode.APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusTooManyRequests {
		t.Errorf("status: got %d, want %d", apiErr.StatusCode, http.StatusTooManyRequests)
	}
	if apiErr.Body != jsonBody {
		t.Errorf("body: got %q, want %q", apiErr.Body, jsonBody)
	}
	if apiErr.RequestID != "req-abc-123" {
		t.Errorf("request ID: got %q, want %q", apiErr.RequestID, "req-abc-123")
	}
}

func TestListStreaming_ErrorStatus(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantStatus int
	}{
		{"401 unauthorized", http.StatusUnauthorized, "unauthorized", http.StatusUnauthorized},
		{"403 forbidden", http.StatusForbidden, "forbidden", http.StatusForbidden},
		{"404 not found", http.StatusNotFound, "not found", http.StatusNotFound},
		{"500 internal", http.StatusInternalServerError, "internal error", http.StatusInternalServerError},
		{"502 bad gateway", http.StatusBadGateway, "", http.StatusBadGateway},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.body != "" {
					_, _ = w.Write([]byte(tt.body))
				}
			}))
			defer server.Close()

			client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			stream := client.Event.ListStreaming(context.Background(), nil)
			if stream.Next() {
				t.Fatal("expected Next() to return false on error status")
			}

			err = stream.Err()
			if err == nil {
				t.Fatal("expected non-nil error from stream")
			}

			var apiErr *opencode.APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("expected *opencode.APIError, got %T: %v", err, err)
			}
			if apiErr.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, apiErr.StatusCode)
			}
		})
	}
}
