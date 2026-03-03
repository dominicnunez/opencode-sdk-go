package opencode

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestClientDoRaw_HeadDoesNotSendBodyOrContentType(t *testing.T) {
	var requestBody []byte
	var contentType string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodHead {
			t.Errorf("Expected HEAD, got %s", r.Method)
		}
		if r.URL.Path != "/probe" {
			t.Errorf("Expected path /probe, got %s", r.URL.Path)
		}

		contentType = r.Header.Get("Content-Type")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}
		requestBody = body

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("create client: %v", err)
	}

	resp, err := client.doRaw(context.Background(), http.MethodHead, "probe", map[string]string{"x": "y"})
	if err != nil {
		t.Fatalf("doRaw HEAD failed: %v", err)
	}
	_ = resp.Body.Close()

	if len(requestBody) != 0 {
		t.Fatalf("Expected no request body for HEAD, got %q", string(requestBody))
	}
	if contentType != "" {
		t.Fatalf("Expected no Content-Type header for HEAD, got %q", contentType)
	}
}

func TestClientDoRaw_RetryDrainCapsBodyRead(t *testing.T) {
	var attempts int32
	var bytesRead int64
	largeBodySize := int64(maxRetryBodyDrainSize * 8)

	client, err := NewClient(
		WithBaseURL("http://127.0.0.1"),
		WithMaxRetries(1),
		WithHTTPClient(&http.Client{
			Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				attempt := atomic.AddInt32(&attempts, 1)
				if attempt == 1 {
					return &http.Response{
						StatusCode: http.StatusTooManyRequests,
						Header:     http.Header{},
						Body: &countingReadCloser{
							remaining: largeBodySize,
							bytesRead: &bytesRead,
						},
					}, nil
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{},
					Body:       io.NopCloser(strings.NewReader("{}")),
				}, nil
			}),
		}),
	)
	if err != nil {
		t.Fatalf("create client: %v", err)
	}

	resp, err := client.doRaw(context.Background(), http.MethodGet, "session", nil)
	if err != nil {
		t.Fatalf("doRaw failed: %v", err)
	}
	_ = resp.Body.Close()

	if got := atomic.LoadInt64(&bytesRead); got > int64(maxRetryBodyDrainSize) {
		t.Fatalf("expected drained bytes <= %d, got %d", maxRetryBodyDrainSize, got)
	}
	if gotAttempts := atomic.LoadInt32(&attempts); gotAttempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", gotAttempts)
	}
}

func TestClientDoRaw_ContextErrorsIncludeOperationContext(t *testing.T) {
	client, err := NewClient(
		WithBaseURL("http://127.0.0.1"),
		WithMaxRetries(0),
		WithHTTPClient(&http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				<-req.Context().Done()
				return nil, req.Context().Err()
			}),
		}),
	)
	if err != nil {
		t.Fatalf("create client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	resp, err := client.doRaw(ctx, http.MethodGet, "session", nil)
	if resp != nil {
		_ = resp.Body.Close()
	}
	if err == nil {
		t.Fatal("expected error from canceled context")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if !strings.Contains(err.Error(), "GET session") {
		t.Fatalf("expected operation context in error, got %v", err)
	}
}

type countingReadCloser struct {
	remaining int64
	bytesRead *int64
	closed    bool
}

func (r *countingReadCloser) Read(p []byte) (int, error) {
	if r.closed {
		return 0, io.ErrClosedPipe
	}
	if r.remaining <= 0 {
		return 0, io.EOF
	}
	n := int64(len(p))
	if n > r.remaining {
		n = r.remaining
	}
	for i := int64(0); i < n; i++ {
		p[i] = 'x'
	}
	r.remaining -= n
	atomic.AddInt64(r.bytesRead, n)
	return int(n), nil
}

func (r *countingReadCloser) Close() error {
	r.closed = true
	return nil
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
