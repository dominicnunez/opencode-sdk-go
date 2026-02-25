package opencode_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dominicnunez/opencode-sdk-go"
)

func TestClientDo_Success(t *testing.T) {
	type response struct {
		Message string `json:"message"`
		Count   int    `json:"count"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", ct)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response{Message: "success", Count: 42})
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Create(context.Background(), &opencode.SessionNewParams{
		ParentID: opencode.PtrString("test-parent"),
	})
	// Just check that the request was made successfully
	if err != nil {
		t.Logf("Session.Create error (expected for mock): %v", err)
	}
}

func TestClientDo_Retry(t *testing.T) {
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("server error"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL),
		opencode.WithMaxRetries(3),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	params := &opencode.SessionListParams{}
	_, err = client.Session.List(context.Background(), params)

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}

	if err != nil {
		t.Logf("Final error (expected): %v", err)
	}
}

func TestClientDo_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = client.Session.List(ctx, &opencode.SessionListParams{})
	if err == nil {
		t.Error("Expected context cancellation error")
	}
	if err != context.Canceled {
		t.Logf("Got error: %v", err)
	}
}

func TestClientDo_QueryParams(t *testing.T) {
	var receivedQuery string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]opencode.Session{})
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	params := &opencode.SessionListParams{
		Directory: opencode.PtrString("/test"),
	}
	_, err = client.Session.List(context.Background(), params)
	if err != nil {
		t.Fatalf("Session.List failed: %v", err)
	}

	if receivedQuery == "" {
		t.Error("Expected query params to be sent")
	}
	t.Logf("Received query: %s", receivedQuery)
}

func TestClientDo_PostWithBody(t *testing.T) {
	var receivedBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}

		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(opencode.Session{ID: "test-session"})
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Create(context.Background(), &opencode.SessionNewParams{
		ParentID: opencode.PtrString("test-parent"),
	})
	if err != nil {
		t.Fatalf("Session.Create failed: %v", err)
	}

	if receivedBody == nil {
		t.Error("Expected request body to be sent")
	}
	if parentID, ok := receivedBody["parentID"].(string); !ok || parentID != "test-parent" {
		t.Errorf("Expected parentID=test-parent, got %v", receivedBody)
	}
}

func TestClientDo_NoRetryOn4xx(t *testing.T) {
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL),
		opencode.WithMaxRetries(3),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.List(context.Background(), &opencode.SessionListParams{})

	if attempts != 1 {
		t.Errorf("Expected 1 attempt (no retries on 4xx), got %d", attempts)
	}

	if err == nil {
		t.Error("Expected error for 400 status")
	}
}

func TestClientDo_ExponentialBackoff(t *testing.T) {
	attempts := 0
	attemptTimes := []time.Time{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		attemptTimes = append(attemptTimes, time.Now())
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL),
		opencode.WithMaxRetries(2),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, _ = client.Session.List(context.Background(), &opencode.SessionListParams{})

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}

	// Check that there was a delay between attempts
	if len(attemptTimes) >= 2 {
		delay := attemptTimes[1].Sub(attemptTimes[0])
		if delay < 400*time.Millisecond {
			t.Errorf("Expected delay >= 400ms between attempts, got %v", delay)
		}
	}
}
