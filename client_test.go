package opencode_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dominicnunez/opencode-sdk-go"
	"github.com/dominicnunez/opencode-sdk-go/internal"
)

type closureTransport struct {
	fn func(req *http.Request) (*http.Response, error)
}

func (t *closureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.fn(req)
}

func TestUserAgentHeader(t *testing.T) {
	var userAgent string
	client, err := opencode.NewClient(
		opencode.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					userAgent = req.Header.Get("User-Agent")
					return &http.Response{
						StatusCode: http.StatusOK,
					}, nil
				},
			},
		}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Session.List(context.Background(), &opencode.SessionListParams{})
	if err == nil {
		t.Error("expected decode error from empty response body")
	}
	if userAgent != fmt.Sprintf("Opencode/Go %s", internal.PackageVersion) {
		t.Errorf("Expected User-Agent to be correct, but got: %#v", userAgent)
	}
}

func TestRetryOn429(t *testing.T) {
	attempts := 0
	client, err := opencode.NewClient(
		opencode.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					attempts++
					return &http.Response{
						StatusCode: http.StatusTooManyRequests,
						Status:     "429 Too Many Requests",
						Header:     http.Header{},
						Body:       io.NopCloser(strings.NewReader("rate limited")),
					}, nil
				},
			},
		}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Session.List(context.Background(), &opencode.SessionListParams{})
	if err == nil {
		t.Error("expected error after exhausting retries on 429")
	}
	if want := 3; attempts != want {
		t.Errorf("expected %d attempts, got %d", want, attempts)
	}
}

func TestRetryOn408(t *testing.T) {
	attempts := 0
	client, err := opencode.NewClient(
		opencode.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					attempts++
					return &http.Response{
						StatusCode: http.StatusRequestTimeout,
						Status:     "408 Request Timeout",
						Header:     http.Header{},
						Body:       io.NopCloser(strings.NewReader("request timeout")),
					}, nil
				},
			},
		}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Session.List(context.Background(), &opencode.SessionListParams{})
	if err == nil {
		t.Error("expected error after exhausting retries on 408")
	}
	if want := 3; attempts != want {
		t.Errorf("expected %d attempts, got %d", want, attempts)
	}
}

func TestContextCancel(t *testing.T) {
	client, err := opencode.NewClient(
		opencode.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					<-req.Context().Done()
					return nil, req.Context().Err()
				},
			},
		}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = client.Session.List(cancelCtx, &opencode.SessionListParams{})
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got: %v", err)
	}
}

func TestContextCancelDelay(t *testing.T) {
	client, err := opencode.NewClient(
		opencode.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					<-req.Context().Done()
					return nil, req.Context().Err()
				},
			},
		}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	cancelCtx, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)
	defer cancel()
	_, err = client.Session.List(cancelCtx, &opencode.SessionListParams{})
	if err == nil {
		t.Error("expected there to be a cancel error")
	}
}

func TestContextDeadline(t *testing.T) {
	testTimeout := time.After(3 * time.Second)
	testDone := make(chan struct{})

	deadline := time.Now().Add(100 * time.Millisecond)
	deadlineCtx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	go func() {
		client, err := opencode.NewClient(
			opencode.WithHTTPClient(&http.Client{
				Transport: &closureTransport{
					fn: func(req *http.Request) (*http.Response, error) {
						<-req.Context().Done()
						return nil, req.Context().Err()
					},
				},
			}),
		)
		if err != nil {
			t.Errorf("failed to create client: %v", err)
			close(testDone)
			return
		}
		_, err = client.Session.List(deadlineCtx, &opencode.SessionListParams{})
		if err == nil {
			t.Error("expected there to be a deadline error")
		}
		close(testDone)
	}()

	select {
	case <-testTimeout:
		t.Fatal("client didn't finish in time")
	case <-testDone:
		if diff := time.Since(deadline); diff < -30*time.Millisecond || 30*time.Millisecond < diff {
			t.Fatalf("client did not return within 30ms of context deadline, got %s", diff)
		}
	}
}

func TestContextDeadlineStreaming(t *testing.T) {
	testTimeout := time.After(3 * time.Second)
	testDone := make(chan struct{})

	deadline := time.Now().Add(100 * time.Millisecond)
	deadlineCtx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	go func() {
		client, err := opencode.NewClient(
			opencode.WithHTTPClient(&http.Client{
				Transport: &closureTransport{
					fn: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: 200,
							Status:     "200 OK",
							Body: io.NopCloser(
								io.Reader(readerFunc(func([]byte) (int, error) {
									<-req.Context().Done()
									return 0, req.Context().Err()
								})),
							),
						}, nil
					},
				},
			}),
		)
		if err != nil {
			t.Errorf("failed to create client: %v", err)
			close(testDone)
			return
		}
		stream := client.Event.ListStreaming(deadlineCtx, &opencode.EventListParams{})
		for stream.Next() {
			_ = stream.Current()
		}
		if stream.Err() == nil {
			t.Error("expected there to be a deadline error")
		}
		close(testDone)
	}()

	select {
	case <-testTimeout:
		t.Fatal("client didn't finish in time")
	case <-testDone:
		if diff := time.Since(deadline); diff < -30*time.Millisecond || 30*time.Millisecond < diff {
			t.Fatalf("client did not return within 30ms of context deadline, got %s", diff)
		}
	}
}

func TestListStreaming_BaseURLQueryParamsPreservedWithMethodParams(t *testing.T) {
	var receivedQuery string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		// Return empty body — stream will error, but we only care about the URL
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL+"?token=abc"),
		opencode.WithHTTPClient(server.Client()),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	stream := client.Event.ListStreaming(context.Background(), &opencode.EventListParams{
		Directory: opencode.Ptr("/test"),
	})
	// Drain the stream so the request is made
	for stream.Next() {
	}
	_ = stream.Close()

	if !strings.Contains(receivedQuery, "token=abc") {
		t.Errorf("expected query to contain token=abc, got %q", receivedQuery)
	}
	if !strings.Contains(receivedQuery, "directory=%2Ftest") {
		t.Errorf("expected query to contain directory=%%2Ftest, got %q", receivedQuery)
	}
}

func TestListStreaming_BaseURLQueryParamsPreservedWithNoMethodParams(t *testing.T) {
	var receivedQuery string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL+"?token=abc"),
		opencode.WithHTTPClient(server.Client()),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	stream := client.Event.ListStreaming(context.Background(), nil)
	for stream.Next() {
	}
	_ = stream.Close()

	if receivedQuery != "token=abc" {
		t.Errorf("expected query to be %q, got %q", "token=abc", receivedQuery)
	}
}

func TestListStreaming_EmptyBody502_ReturnsAPIErrorWithStatusText(t *testing.T) {
	client, err := opencode.NewClient(
		opencode.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusBadGateway,
						Status:     "502 Bad Gateway",
						Header:     http.Header{},
						Body:       io.NopCloser(strings.NewReader("")),
					}, nil
				},
			},
		}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	stream := client.Event.ListStreaming(context.Background(), nil)
	for stream.Next() {
	}

	var apiErr *opencode.APIError
	if !errors.As(stream.Err(), &apiErr) {
		t.Fatalf("expected *opencode.APIError, got %T: %v", stream.Err(), stream.Err())
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
	_ = stream.Close()
}

func TestBaseURL_WithPathComponent_ResolvesCorrectly(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL + "/api/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, _ = client.Session.List(context.Background(), nil)

	if receivedPath != "/api/v1/session" {
		t.Errorf("expected path /api/v1/session, got %s", receivedPath)
	}
}

func TestBaseURL_WithTrailingSlash_ResolvesCorrectly(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL + "/api/v1/"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, _ = client.Session.List(context.Background(), nil)

	if receivedPath != "/api/v1/session" {
		t.Errorf("expected path /api/v1/session, got %s", receivedPath)
	}
}

type readerFunc func([]byte) (int, error)

func (f readerFunc) Read(p []byte) (int, error) { return f(p) }
func (f readerFunc) Close() error               { return nil }
