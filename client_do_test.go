package opencode_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dominicnunez/opencode-sdk-go"
)

func TestClientDo_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", ct)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(opencode.Session{
			ID:    "sess_1",
			Title: "Test Session",
			Time:  opencode.SessionTime{Created: 1, Updated: 1},
		})
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	session, err := client.Session.Create(context.Background(), &opencode.SessionCreateParams{
		ParentID: opencode.Ptr("test-parent"),
	})
	if err != nil {
		t.Fatalf("Session.Create failed: %v", err)
	}
	if session.ID != "sess_1" {
		t.Errorf("expected session ID %q, got %q", "sess_1", session.ID)
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
		_ = json.NewEncoder(w).Encode([]opencode.Session{
			{ID: "sess_1", Title: "Recovered", Time: opencode.SessionTime{Created: 1, Updated: 1}},
		})
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL),
		opencode.WithMaxRetries(3),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	sessions, err := client.Session.List(context.Background(), &opencode.SessionListParams{})
	if err != nil {
		t.Fatalf("Session.List failed after retries: %v", err)
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
	if len(sessions) != 1 || sessions[0].ID != "sess_1" {
		t.Errorf("expected 1 session with ID sess_1, got %v", sessions)
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
		t.Fatal("expected context cancellation error")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got: %v", err)
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
		Directory: opencode.Ptr("/test"),
	}
	_, err = client.Session.List(context.Background(), params)
	if err != nil {
		t.Fatalf("Session.List failed: %v", err)
	}

	if receivedQuery == "" {
		t.Fatal("Expected query params to be sent")
	}
	if !strings.Contains(receivedQuery, "directory=%2Ftest") {
		t.Errorf("expected query to contain directory=%%2Ftest, got %q", receivedQuery)
	}
}

func TestClientDo_BaseURLQueryParamsPreserved(t *testing.T) {
	var receivedQuery string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]opencode.Session{})
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL + "?token=xyz"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// List with no query params of its own — base URL token should survive
	_, err = client.Session.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Session.List failed: %v", err)
	}

	if receivedQuery != "token=xyz" {
		t.Errorf("Expected base URL query params preserved, got %q", receivedQuery)
	}
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

	_, err = client.Session.Create(context.Background(), &opencode.SessionCreateParams{
		ParentID: opencode.Ptr("test-parent"),
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
		t.Fatalf("Expected 3 attempts, got %d", attempts)
	}

	// With 3 attempts we have 2 inter-attempt delays.
	// Backoff formula: initialBackoff * (1 << attempt) → 500ms, 1000ms.
	// Verify each delay is at least 80% of the expected value (timing tolerance)
	// and that the second delay is meaningfully larger than the first.
	if len(attemptTimes) < 3 {
		t.Fatal("not enough attempt timestamps recorded")
	}

	delay1 := attemptTimes[1].Sub(attemptTimes[0])
	delay2 := attemptTimes[2].Sub(attemptTimes[1])

	if delay1 < 400*time.Millisecond {
		t.Errorf("first delay should be ~500ms, got %v", delay1)
	}
	if delay2 < 800*time.Millisecond {
		t.Errorf("second delay should be ~1000ms, got %v", delay2)
	}
	if delay2 <= delay1 {
		t.Errorf("delays should increase exponentially: delay1=%v, delay2=%v", delay1, delay2)
	}
}

func TestClientDo_EmptyBody502_ReturnsAPIErrorWithStatusText(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL),
		opencode.WithMaxRetries(0),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for 502 response")
	}

	var apiErr *opencode.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *opencode.APIError, got %T: %v", err, err)
	}
	if apiErr.Message == "" {
		t.Error("expected non-empty Message on APIError for empty body 502")
	}
	if apiErr.Message != "Bad Gateway" {
		t.Errorf("expected Message to be %q, got %q", "Bad Gateway", apiErr.Message)
	}
	if apiErr.StatusCode != http.StatusBadGateway {
		t.Errorf("expected StatusCode %d, got %d", http.StatusBadGateway, apiErr.StatusCode)
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt with maxRetries=0, got %d", attempts)
	}
}

func TestClientDo_MaxRetriesZero_ExactlyOneAttempt(t *testing.T) {
	// WithMaxRetries(0) means exactly one attempt with no retries.
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("server error"))
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL),
		opencode.WithMaxRetries(0),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.List(context.Background(), &opencode.SessionListParams{})
	if err == nil {
		t.Fatal("expected error for 500 response")
	}

	if attempts != 1 {
		t.Errorf("expected exactly 1 attempt with maxRetries=0, got %d", attempts)
	}
}

func TestClientDo_PostBodyReencodedOnRetry(t *testing.T) {
	var bodies []string
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		body, _ := io.ReadAll(r.Body)
		bodies = append(bodies, string(body))

		if attempts < 2 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("server error"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(opencode.Session{
			ID:   "sess_1",
			Time: opencode.SessionTime{Created: 1, Updated: 1},
		})
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL),
		opencode.WithMaxRetries(2),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Create(context.Background(), &opencode.SessionCreateParams{
		ParentID: opencode.Ptr("test-parent"),
	})
	if err != nil {
		t.Fatalf("Session.Create failed: %v", err)
	}

	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}

	if len(bodies) != 2 {
		t.Fatalf("expected 2 request bodies, got %d", len(bodies))
	}

	// Both attempts should receive the same non-empty body
	for i, body := range bodies {
		if body == "" {
			t.Errorf("attempt %d: expected non-empty body", i+1)
		}
		if !strings.Contains(body, "test-parent") {
			t.Errorf("attempt %d: expected body to contain 'test-parent', got %q", i+1, body)
		}
	}

	if bodies[0] != bodies[1] {
		t.Errorf("body mismatch between attempts:\n  attempt 1: %s\n  attempt 2: %s", bodies[0], bodies[1])
	}
}

func TestClientDo_TransportErrorRetryExhaustion(t *testing.T) {
	transportErr := errors.New("connection refused")
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not reach server")
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL),
		opencode.WithMaxRetries(2),
		opencode.WithHTTPClient(&http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				attempts++
				return nil, transportErr
			}),
		}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.List(context.Background(), &opencode.SessionListParams{})
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}

	if attempts != 3 {
		t.Errorf("expected 3 attempts (1 initial + 2 retries), got %d", attempts)
	}

	if !errors.Is(err, transportErr) {
		t.Errorf("expected error to wrap transport error, got: %v", err)
	}

	if !strings.Contains(err.Error(), "2 retries") {
		t.Errorf("expected error to mention retry count, got: %v", err)
	}
}

func TestClientDo_3xxRedirectIsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/other")
		w.WriteHeader(http.StatusMovedPermanently)
		_, _ = w.Write([]byte("moved"))
	}))
	defer server.Close()

	// Disable redirect following so the client sees the 301 directly
	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL),
		opencode.WithHTTPClient(&http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.List(context.Background(), &opencode.SessionListParams{})
	if err == nil {
		t.Fatal("expected error for 3xx response, got nil")
	}

	var apiErr *opencode.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *opencode.APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusMovedPermanently {
		t.Errorf("expected status %d, got %d", http.StatusMovedPermanently, apiErr.StatusCode)
	}
}

func TestClientDo_ContextCancelledDuringBackoffDelay(t *testing.T) {
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("server error"))
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL),
		opencode.WithMaxRetries(3),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context after the first attempt completes but during the backoff
	// delay (initialBackoff = 500ms). 100ms gives time for the request to
	// complete and enter the timer select, but is well under the 500ms delay.
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	_, err = client.Session.List(ctx, &opencode.SessionListParams{})
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected error after context cancellation")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got: %v", err)
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt before cancellation, got %d", attempts)
	}
	// Should return well before the full backoff (500ms) plus remaining
	// retries would take. 400ms gives generous timing tolerance.
	if elapsed > 400*time.Millisecond {
		t.Errorf("expected prompt return after cancellation, took %v", elapsed)
	}
}

func TestClientDo_ContextCancelledDuringInFlightRequest(t *testing.T) {
	// Verifies that cancelling a context while the HTTP request is in-flight
	// propagates context.Canceled through the SDK's error wrapping.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Hold the connection open until the client cancels.
		<-r.Context().Done()
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL),
		opencode.WithMaxRetries(0),
		opencode.WithTimeout(5*time.Second),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel after a short delay so the request is in-flight.
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	_, err = client.Session.List(ctx, &opencode.SessionListParams{})
	if err == nil {
		t.Fatal("expected error after context cancellation during in-flight request")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got: %v", err)
	}
}

func TestClientDo_JSONDecodeFailureOnSuccessResponse(t *testing.T) {
	// Verifies that HTTP 200 with invalid JSON body produces a decode error
	// containing method and path context.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("not json"))
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.List(context.Background(), &opencode.SessionListParams{})
	if err == nil {
		t.Fatal("expected error for invalid JSON response body")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "decode") {
		t.Errorf("expected error to contain %q, got: %v", "decode", errMsg)
	}
	if !strings.Contains(errMsg, "GET") {
		t.Errorf("expected error to contain HTTP method, got: %v", errMsg)
	}
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
