package opencode_test

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

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
	if got := err.Error(); !strings.Contains(got, "GET event:") {
		t.Fatalf("expected operation context in error string, got %q", got)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream := client.Event.ListStreaming(ctx, nil)
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
			if got := err.Error(); !strings.Contains(got, "GET event:") {
				t.Fatalf("expected operation context in error string, got %q", got)
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

func TestListStreaming_MissingContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("event: message\ndata: {\"type\":\"message.updated\"}\n\n"))
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	stream := client.Event.ListStreaming(context.Background(), nil)
	defer func() { _ = stream.Close() }()
	if stream.Next() {
		t.Fatal("expected Next() to return false when Content-Type is missing")
	}

	err = stream.Err()
	if err == nil {
		t.Fatal("expected non-nil error for missing content type")
	}
	if !strings.Contains(err.Error(), "unexpected content type") {
		t.Errorf("expected error about unexpected content type, got: %v", err)
	}
}

func TestListStreaming_ExplicitDeadlineStaysOpenPastClientTimeout(t *testing.T) {
	const clientTimeout = 50 * time.Millisecond
	sendEvent := make(chan struct{})

	client, err := opencode.NewClient(
		opencode.WithTimeout(clientTimeout),
		opencode.WithHTTPClient(&http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"text/event-stream"}},
					Body:       newBlockedSSEBody(sendEvent, []byte("event: message\ndata: {\"type\":\"message.updated\"}\n\n")),
				}, nil
			}),
		}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream := client.Event.ListStreaming(ctx, nil)
	defer func() { _ = stream.Close() }()

	time.Sleep(clientTimeout + 20*time.Millisecond)
	close(sendEvent)

	if !stream.Next() {
		t.Fatalf("expected event after waiting longer than client timeout, got err: %v", stream.Err())
	}
}

func TestListStreaming_NoDeadlineIgnoresHTTPClientTimeoutAfterConnect(t *testing.T) {
	const httpClientTimeout = 50 * time.Millisecond
	const eventDelay = httpClientTimeout + 50*time.Millisecond

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("response writer does not support flushing")
		}
		flusher.Flush()

		time.Sleep(eventDelay)
		_, _ = w.Write([]byte("event: message\ndata: {\"type\":\"message.updated\"}\n\n"))
		flusher.Flush()
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL),
		opencode.WithHTTPClient(&http.Client{Timeout: httpClientTimeout}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	stream := client.Event.ListStreaming(context.Background(), nil)
	defer func() { _ = stream.Close() }()

	if !stream.Next() {
		t.Fatalf("expected event after waiting longer than http client timeout, got err: %v", stream.Err())
	}
}

func TestListStreaming_NoDeadlineUsesClientTimeoutDuringConnect(t *testing.T) {
	const clientTimeout = 50 * time.Millisecond
	const waitForResult = 250 * time.Millisecond

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("event: message\ndata: {\"type\":\"message.updated\"}\n\n"))
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL),
		opencode.WithTimeout(clientTimeout),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	done := make(chan string, 1)
	go func() {
		stream := client.Event.ListStreaming(context.Background(), nil)
		defer func() { _ = stream.Close() }()

		if stream.Next() {
			done <- "expected Next() to return false when connect deadline is exceeded"
			return
		}

		streamErr := stream.Err()
		if streamErr == nil {
			done <- "expected non-nil error when connect deadline is exceeded"
			return
		}
		if !strings.Contains(streamErr.Error(), "timeout") {
			done <- "expected timeout-related error when waiting for initial response headers"
			return
		}
		done <- ""
	}()

	select {
	case result := <-done:
		if result != "" {
			t.Fatal(result)
		}
	case <-time.After(waitForResult):
		t.Fatalf("expected ListStreaming to apply client timeout during connect and return within %s", waitForResult)
	}
}

func TestListStreaming_NoDeadlineUsesClientTimeoutDuringTLSHandshake(t *testing.T) {
	const clientTimeout = 50 * time.Millisecond
	const waitForResult = 300 * time.Millisecond

	baseURL, cleanup := startStalledTLSHandshakeEndpoint(t)
	defer cleanup()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(baseURL),
		opencode.WithTimeout(clientTimeout),
		opencode.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // controlled test endpoint
			},
		}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	done := make(chan string, 1)
	go func() {
		stream := client.Event.ListStreaming(context.Background(), nil)
		defer func() { _ = stream.Close() }()

		if stream.Next() {
			done <- "expected Next() to return false when TLS handshake exceeds client timeout"
			return
		}

		streamErr := stream.Err()
		if streamErr == nil {
			done <- "expected non-nil error when TLS handshake exceeds client timeout"
			return
		}
		if !strings.Contains(streamErr.Error(), "timeout") {
			done <- "expected timeout-related error when TLS handshake exceeds client timeout"
			return
		}
		done <- ""
	}()

	select {
	case result := <-done:
		if result != "" {
			t.Fatal(result)
		}
	case <-time.After(waitForResult):
		t.Fatalf("expected ListStreaming to apply client timeout during TLS handshake and return within %s", waitForResult)
	}
}

func TestListStreaming_CustomTransportWithExplicitDeadline(t *testing.T) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	stream := client.Event.ListStreaming(ctx, nil)
	defer func() { _ = stream.Close() }()

	if stream.Next() {
		t.Fatal("expected Next() to return false when connect deadline is exceeded")
	}
	if !errors.Is(stream.Err(), context.DeadlineExceeded) {
		t.Fatalf("expected context deadline exceeded, got: %v", stream.Err())
	}
}

type blockedSSEBody struct {
	ready <-chan struct{}
	data  []byte
	sent  bool
}

func newBlockedSSEBody(ready <-chan struct{}, data []byte) *blockedSSEBody {
	return &blockedSSEBody{ready: ready, data: data}
}

func (b *blockedSSEBody) Read(p []byte) (int, error) {
	if b.sent {
		return 0, io.EOF
	}
	<-b.ready
	b.sent = true
	n := copy(p, b.data)
	return n, nil
}

func (b *blockedSSEBody) Close() error { return nil }

func TestListStreaming_NoDeadlineWithCustomTransportUsesClientTimeoutDuringConnect(t *testing.T) {
	const clientTimeout = 50 * time.Millisecond
	const waitForResult = 250 * time.Millisecond

	var calls int32
	client, err := opencode.NewClient(
		opencode.WithTimeout(clientTimeout),
		opencode.WithHTTPClient(&http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				atomic.AddInt32(&calls, 1)
				<-req.Context().Done()
				return nil, req.Context().Err()
			}),
		}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	done := make(chan string, 1)
	go func() {
		stream := client.Event.ListStreaming(context.Background(), nil)
		defer func() { _ = stream.Close() }()

		if stream.Next() {
			done <- "expected Next() to return false when connect deadline is exceeded"
			return
		}

		streamErr := stream.Err()
		if streamErr == nil {
			done <- "expected non-nil error when connect deadline is exceeded"
			return
		}
		if !strings.Contains(streamErr.Error(), "timeout") && !errors.Is(streamErr, context.Canceled) {
			done <- "expected timeout-related error when waiting for custom transport connect"
			return
		}
		done <- ""
	}()

	select {
	case result := <-done:
		if result != "" {
			t.Fatal(result)
		}
	case <-time.After(waitForResult):
		t.Fatalf("expected ListStreaming to apply client timeout during custom transport connect and return within %s", waitForResult)
	}

	if atomic.LoadInt32(&calls) != 1 {
		t.Fatalf("expected transport to be called once, got %d calls", calls)
	}
}

func TestListStreaming_CustomTransportTimeoutRaceReturnsTimeoutAndClosesBody(t *testing.T) {
	const clientTimeout = 30 * time.Millisecond
	const connectDelay = 75 * time.Millisecond

	var closedCount int32
	body := &closeTrackedBody{
		reader: strings.NewReader("event: message\ndata: {\"type\":\"message.updated\"}\n\n"),
		closed: &closedCount,
	}

	client, err := opencode.NewClient(
		opencode.WithTimeout(clientTimeout),
		opencode.WithHTTPClient(&http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				time.Sleep(connectDelay)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"text/event-stream"}},
					Body:       body,
					Request:    req,
				}, nil
			}),
		}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	stream := client.Event.ListStreaming(context.Background(), nil)
	if stream.Next() {
		_ = stream.Close()
		t.Fatal("expected stream connect timeout when delayed response exceeds client timeout")
	}
	streamErr := stream.Err()
	if streamErr == nil {
		t.Fatal("expected timeout error when delayed response exceeds client timeout")
	}
	if !errors.Is(streamErr, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded error, got: %v", streamErr)
	}
	if got := atomic.LoadInt32(&closedCount); got != 1 {
		t.Fatalf("expected response body to be closed exactly once, got %d", got)
	}
}

func TestListStreaming_CustomTransportConnectBeforeTimeoutSucceeds(t *testing.T) {
	const clientTimeout = 40 * time.Millisecond
	const connectDelay = 35 * time.Millisecond

	var closedCount int32
	body := &closeTrackedBody{
		reader: strings.NewReader("event: message\ndata: {\"type\":\"message.updated\"}\n\n"),
		closed: &closedCount,
	}

	client, err := opencode.NewClient(
		opencode.WithTimeout(clientTimeout),
		opencode.WithHTTPClient(&http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				time.Sleep(connectDelay)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"text/event-stream"}},
					Body:       body,
					Request:    req,
				}, nil
			}),
		}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	stream := client.Event.ListStreaming(context.Background(), nil)
	if !stream.Next() {
		_ = stream.Close()
		t.Fatalf("expected first event after connect before timeout, got err: %v", stream.Err())
	}

	event := stream.Current()
	if event.Type != opencode.EventTypeMessageUpdated {
		_ = stream.Close()
		t.Fatalf("expected message.updated event type, got: %s", event.Type)
	}
	if got := atomic.LoadInt32(&closedCount); got != 0 {
		_ = stream.Close()
		t.Fatalf("expected response body to remain open for active stream, got %d closes", got)
	}

	if closeErr := stream.Close(); closeErr != nil {
		t.Fatalf("failed to close stream: %v", closeErr)
	}
	if got := atomic.LoadInt32(&closedCount); got != 1 {
		t.Fatalf("expected response body to be closed once by stream.Close, got %d closes", got)
	}
}

func startStalledTLSHandshakeEndpoint(t *testing.T) (string, func()) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}

	done := make(chan struct{})
	var wg sync.WaitGroup
	var mu sync.Mutex
	var connections []net.Conn

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				select {
				case <-done:
					return
				default:
					return
				}
			}

			mu.Lock()
			connections = append(connections, conn)
			mu.Unlock()

			wg.Add(1)
			go func(c net.Conn) {
				defer wg.Done()
				<-done
				_ = c.Close()
			}(conn)
		}
	}()

	cleanup := func() {
		close(done)
		_ = listener.Close()

		mu.Lock()
		for _, c := range connections {
			_ = c.Close()
		}
		mu.Unlock()

		wg.Wait()
	}

	return "https://" + listener.Addr().String(), cleanup
}

type closeTrackedBody struct {
	reader io.Reader
	closed *int32
}

func (b *closeTrackedBody) Read(p []byte) (int, error) {
	return b.reader.Read(p)
}

func (b *closeTrackedBody) Close() error {
	atomic.AddInt32(b.closed, 1)
	return nil
}
