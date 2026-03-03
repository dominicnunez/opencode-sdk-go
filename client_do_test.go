package opencode_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
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

func TestClientDo_BaseURLQueryParamsRejected(t *testing.T) {
	_, err := opencode.NewClient(opencode.WithBaseURL("https://example.com?workspace=xyz"))
	if err == nil {
		t.Fatal("expected base URL query parameters to be rejected")
	}
	if !strings.Contains(err.Error(), "must not include query parameters") {
		t.Fatalf("expected query parameter validation error, got: %v", err)
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

func TestClientDo_ErrorHandlingContract_AsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", "req_123")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte("slow down"))
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
		t.Fatal("expected error for 429 response")
	}

	var apiErr *opencode.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected errors.As(err, *APIError) to match, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected status %d, got %d", http.StatusTooManyRequests, apiErr.StatusCode)
	}
	if apiErr.RequestID != "req_123" {
		t.Fatalf("expected request ID %q, got %q", "req_123", apiErr.RequestID)
	}
	if apiErr.Body != "slow down" {
		t.Fatalf("expected body %q, got %q", "slow down", apiErr.Body)
	}
}

func TestClientDo_RetryAfterHeaderDelay(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte("slow down"))
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]opencode.Session{{
			ID:    "sess_1",
			Title: "Recovered",
			Time:  opencode.SessionTime{Created: 1, Updated: 1},
		}})
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL),
		opencode.WithMaxRetries(1),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	start := time.Now()
	_, err = client.Session.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Session.List failed: %v", err)
	}

	elapsed := time.Since(start)
	if elapsed < 900*time.Millisecond {
		t.Fatalf("expected retry delay to honor Retry-After header, elapsed %v", elapsed)
	}
	if attempts != 2 {
		t.Fatalf("expected two attempts, got %d", attempts)
	}
}

func TestClientDo_PostDoesNotRetryOnRetryableHTTPStatus(t *testing.T) {
	var bodies []string
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		body, _ := io.ReadAll(r.Body)
		bodies = append(bodies, string(body))

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("server error"))
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
	if err == nil {
		t.Fatal("expected Session.Create to fail without retrying POST")
	}

	var apiErr *opencode.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *opencode.APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, apiErr.StatusCode)
	}

	if attempts != 1 {
		t.Fatalf("expected 1 attempt for POST, got %d", attempts)
	}

	if len(bodies) != 1 {
		t.Fatalf("expected 1 request body, got %d", len(bodies))
	}

	if bodies[0] == "" {
		t.Fatal("expected non-empty request body")
	}

	if !strings.Contains(bodies[0], "test-parent") {
		t.Fatalf("expected body to contain %q, got %q", "test-parent", bodies[0])
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

	if !strings.Contains(err.Error(), "GET session") {
		t.Errorf("expected error to include request identity, got: %v", err)
	}
	if !strings.Contains(err.Error(), "2 retries") {
		t.Errorf("expected error to mention retry count, got: %v", err)
	}
}

func TestClientDo_PostTransportErrorDoesNotRetry(t *testing.T) {
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

	_, err = client.Session.Create(context.Background(), &opencode.SessionCreateParams{
		ParentID: opencode.Ptr("test-parent"),
	})
	if err == nil {
		t.Fatal("expected transport error for POST")
	}

	if attempts != 1 {
		t.Fatalf("expected 1 attempt for POST transport error, got %d", attempts)
	}
	if !errors.Is(err, transportErr) {
		t.Fatalf("expected error to wrap transport error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "0 retries") {
		t.Fatalf("expected error to mention zero retries, got: %v", err)
	}
}

func TestClientDo_DefaultClientTreats3xxAsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/other")
		w.WriteHeader(http.StatusMovedPermanently)
		_, _ = w.Write([]byte("moved"))
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
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

func TestClientDo_DefaultClientDoesNotForwardBodyOn307(t *testing.T) {
	redirectedRequestBodies := make(chan string, 1)
	redirectTarget := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		redirectedRequestBodies <- string(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer redirectTarget.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", redirectTarget.URL+"/capture")
		w.WriteHeader(http.StatusTemporaryRedirect)
		_, _ = w.Write([]byte("redirected"))
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Auth.Set(context.Background(), "provider-id", &opencode.AuthSetParams{
		Auth: opencode.ApiAuth{Key: "super-secret"},
	})
	if err == nil {
		t.Fatal("expected error for 307 response")
	}
	var apiErr *opencode.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *opencode.APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusTemporaryRedirect {
		t.Fatalf("expected status %d, got %d", http.StatusTemporaryRedirect, apiErr.StatusCode)
	}

	select {
	case body := <-redirectedRequestBodies:
		t.Fatalf("unexpected redirected request body: %q", body)
	case <-time.After(150 * time.Millisecond):
	}
}

func TestClientDo_WithHTTPClientStillDoesNotForwardBodyOn307(t *testing.T) {
	redirectedRequestBodies := make(chan string, 1)
	redirectTarget := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		redirectedRequestBodies <- string(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer redirectTarget.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", redirectTarget.URL+"/capture")
		w.WriteHeader(http.StatusTemporaryRedirect)
		_, _ = w.Write([]byte("redirected"))
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL),
		opencode.WithHTTPClient(&http.Client{}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Auth.Set(context.Background(), "provider-id", &opencode.AuthSetParams{
		Auth: opencode.ApiAuth{Key: "super-secret"},
	})
	if err == nil {
		t.Fatal("expected error for 307 response")
	}
	var apiErr *opencode.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *opencode.APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusTemporaryRedirect {
		t.Fatalf("expected status %d, got %d", http.StatusTemporaryRedirect, apiErr.StatusCode)
	}

	select {
	case body := <-redirectedRequestBodies:
		t.Fatalf("unexpected redirected request body: %q", body)
	case <-time.After(150 * time.Millisecond):
	}
}

func TestClientDo_ContextCancelledDuringBackoffDelay(t *testing.T) {
	attempts := 0
	firstAttemptDone := make(chan struct{}, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("server error"))
		if attempts == 1 {
			firstAttemptDone <- struct{}{}
		}
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
	go func() {
		<-firstAttemptDone
		cancel()
	}()

	_, err = client.Session.List(ctx, &opencode.SessionListParams{})

	if err == nil {
		t.Fatal("expected error after context cancellation")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got: %v", err)
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt before cancellation, got %d", attempts)
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

func TestClientDo_JSONDecodeRejectsTrailingBytes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"sess_1","title":"ok","time":{"created":1,"updated":1}}TRAILING`))
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Get(context.Background(), "sess_1", nil)
	if err == nil {
		t.Fatal("expected error for JSON response with trailing bytes")
	}
	if !strings.Contains(err.Error(), "decode") {
		t.Fatalf("expected decode error, got: %v", err)
	}
}

func TestClientDo_SuccessResponseExceedsDefaultSizeLimit(t *testing.T) {
	const defaultLimitExceededBodySize = 16 << 20
	oversizedBody := strings.Repeat("a", defaultLimitExceededBodySize)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"sess_1","title":"` + oversizedBody + `","time":{"created":1,"updated":1}}`))
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Get(context.Background(), "sess_1", nil)
	if err == nil {
		t.Fatal("expected size limit error for oversized successful response body")
	}
	if !strings.Contains(err.Error(), "response body exceeds") {
		t.Fatalf("expected body limit error, got: %v", err)
	}
}

func TestClientDo_SuccessResponseExceedsConfiguredSizeLimit(t *testing.T) {
	const maxSuccessBodySize = int64(1 << 20)
	oversizedBody := strings.Repeat("a", int(maxSuccessBodySize)+1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"sess_1","title":"` + oversizedBody + `","time":{"created":1,"updated":1}}`))
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL),
		opencode.WithMaxSuccessBodySize(maxSuccessBodySize),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.Get(context.Background(), "sess_1", nil)
	if err == nil {
		t.Fatal("expected size limit error for oversized successful response body")
	}
	if !strings.Contains(err.Error(), "response body exceeds") {
		t.Fatalf("expected body limit error, got: %v", err)
	}
}

func TestWithMaxSuccessBodySize_RejectsMaxInt64(t *testing.T) {
	_, err := opencode.NewClient(opencode.WithMaxSuccessBodySize(math.MaxInt64))
	if err == nil {
		t.Fatal("expected error when max success body size cannot be represented safely")
	}
	if !strings.Contains(err.Error(), "must be at most") {
		t.Fatalf("expected max size validation error, got: %v", err)
	}
}

func TestClientDo_SuccessResponseWithNearMaxConfiguredSize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("[]"))
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL),
		opencode.WithMaxSuccessBodySize(math.MaxInt64-1),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	sessions, err := client.Session.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected response decode to succeed near int64 boundary, got: %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("expected no sessions, got %d", len(sessions))
	}
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type trackedReadCloser struct {
	reader io.Reader
	closed bool
}

func (t *trackedReadCloser) Read(p []byte) (int, error) {
	return t.reader.Read(p)
}

func (t *trackedReadCloser) Close() error {
	t.closed = true
	return nil
}

func TestClientDo_Unexpected1xxReturnsAPIErrorAndClosesBody(t *testing.T) {
	trackedBody := &trackedReadCloser{reader: strings.NewReader("switching protocols")}

	client, err := opencode.NewClient(
		opencode.WithBaseURL("http://localhost"),
		opencode.WithMaxRetries(0),
		opencode.WithHTTPClient(&http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusSwitchingProtocols,
					Header:     make(http.Header),
					Body:       trackedBody,
					Request:    req,
				}, nil
			}),
		}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Session.List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for 1xx status")
	}

	var apiErr *opencode.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *opencode.APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("expected status %d, got %d", http.StatusSwitchingProtocols, apiErr.StatusCode)
	}
	if strings.Contains(err.Error(), "%!w(<nil>)") {
		t.Fatalf("unexpected nil-wrap formatting in error: %v", err)
	}
	if !trackedBody.closed {
		t.Fatal("expected 1xx response body to be closed")
	}
}
