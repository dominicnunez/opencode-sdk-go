package opencode_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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
	defer func() { _ = stream.Close() }()
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

func TestListStreaming_ContextCancelDuringConnect(t *testing.T) {
	client, err := opencode.NewClient(
		opencode.WithHTTPClient(&http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				<-req.Context().Done()
				return nil, req.Context().Err()
			}),
		}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	stream := client.Event.ListStreaming(ctx, nil)
	defer func() { _ = stream.Close() }()
	if stream.Next() {
		t.Fatal("expected Next() to return false on cancelled context")
	}

	err = stream.Err()
	if err == nil {
		t.Fatal("expected non-nil error from stream")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected error to wrap context.Canceled, got: %v", err)
	}
	if got := err.Error(); !strings.Contains(got, "event stream request") {
		t.Errorf("expected error to contain 'event stream request', got: %q", got)
	}
}

func TestListStreaming_TransportErrorWrapped(t *testing.T) {
	transportErr := errors.New("connection refused")
	client, err := opencode.NewClient(
		opencode.WithHTTPClient(&http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return nil, transportErr
			}),
		}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	stream := client.Event.ListStreaming(context.Background(), nil)
	defer func() { _ = stream.Close() }()
	if stream.Next() {
		t.Fatal("expected Next() to return false on transport error")
	}

	err = stream.Err()
	if err == nil {
		t.Fatal("expected non-nil error from stream")
	}
	if !errors.Is(err, transportErr) {
		t.Fatalf("expected error to wrap transport error, got: %v", err)
	}
	if got := err.Error(); !strings.Contains(got, "event stream request") {
		t.Errorf("expected error to contain 'event stream request', got: %q", got)
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
			defer func() { _ = stream.Close() }()
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

func TestListStreaming_UnexpectedContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"events": []}`))
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	stream := client.Event.ListStreaming(context.Background(), nil)
	defer func() { _ = stream.Close() }()
	if stream.Next() {
		t.Fatal("expected Next() to return false on wrong content type")
	}

	err = stream.Err()
	if err == nil {
		t.Fatal("expected non-nil error for unexpected content type")
	}
	if !strings.Contains(err.Error(), "unexpected content type") {
		t.Errorf("expected error about unexpected content type, got: %v", err)
	}
	if !strings.Contains(err.Error(), "application/json") {
		t.Errorf("expected error to mention actual content type, got: %v", err)
	}
}
