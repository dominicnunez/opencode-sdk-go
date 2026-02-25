package requestconfig

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestNewRequestConfig verifies basic request config creation
func TestNewRequestConfig(t *testing.T) {
	ctx := context.Background()

	// Test with nil body
	cfg, err := NewRequestConfig(ctx, "GET", "http://example.com", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Request.Method != "GET" {
		t.Errorf("expected method GET, got %s", cfg.Request.Method)
	}

	// Test with JSON body
	type TestBody struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	body := TestBody{Name: "test", Count: 42}
	cfg, err = NewRequestConfig(ctx, "POST", "http://example.com", body, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Request.Header.Get(HeaderContentType) != ContentTypeJSON {
		t.Errorf("expected content type %s, got %s", ContentTypeJSON, cfg.Request.Header.Get(HeaderContentType))
	}
}

// TestRequestConfigExecute verifies the request execution flow
func TestRequestConfigExecute(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result": "success"}`))
	}))
	defer server.Close()

	ctx := context.Background()
	type Response struct {
		Result string `json:"result"`
	}
	var response Response

	cfg, err := NewRequestConfig(ctx, "GET", "/", nil, &response)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Set base URL
	cfg.DefaultBaseURL, _ = url.Parse(server.URL)
	err = cfg.Execute()
	if err != nil {
		t.Fatalf("unexpected error during execute: %v", err)
	}

	if response.Result != "success" {
		t.Errorf("expected result 'success', got '%s'", response.Result)
	}
}

// TestRequestConfigRetry verifies retry logic
func TestRequestConfigRetry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result": "ok"}`))
	}))
	defer server.Close()

	ctx := context.Background()
	type Response struct {
		Result string `json:"result"`
	}
	var response Response

	cfg, err := NewRequestConfig(ctx, "GET", "/", nil, &response)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Set base URL
	cfg.DefaultBaseURL, _ = url.Parse(server.URL)
	err = cfg.Execute()
	if err != nil {
		t.Fatalf("unexpected error during execute: %v", err)
	}

	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
	if response.Result != "ok" {
		t.Errorf("expected result 'ok', got '%s'", response.Result)
	}
}

// TestShouldRetry verifies retry decision logic
func TestShouldRetry(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		header     map[string]string
		want       bool
	}{
		{"success", http.StatusOK, nil, false},
		{"bad request", http.StatusBadRequest, nil, false},
		{"timeout", http.StatusRequestTimeout, nil, true},
		{"conflict", http.StatusConflict, nil, true},
		{"too many requests", http.StatusTooManyRequests, nil, true},
		{"internal error", http.StatusInternalServerError, nil, true},
		{"bad gateway", http.StatusBadGateway, nil, true},
		{"explicit retry true", http.StatusOK, map[string]string{"x-should-retry": "true"}, true},
		{"explicit retry false", http.StatusInternalServerError, map[string]string{"x-should-retry": "false"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			res := &http.Response{
				StatusCode: tt.statusCode,
				Header:     http.Header{},
			}
			for k, v := range tt.header {
				res.Header.Set(k, v)
			}

			got := shouldRetry(req, res)
			if got != tt.want {
				t.Errorf("shouldRetry() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRequestWithByteBody verifies []byte body handling
func TestRequestWithByteBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if string(body) != "test data" {
			t.Errorf("expected body 'test data', got '%s'", string(body))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx := context.Background()
	body := []byte("test data")

	cfg, err := NewRequestConfig(ctx, "POST", "/", body, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Set base URL
	cfg.DefaultBaseURL, _ = url.Parse(server.URL)
	err = cfg.Execute()
	if err != nil {
		t.Fatalf("unexpected error during execute: %v", err)
	}
}

// TestRequestWithReaderBody verifies io.Reader body handling
func TestRequestWithReaderBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if string(body) != "reader data" {
			t.Errorf("expected body 'reader data', got '%s'", string(body))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx := context.Background()
	body := strings.NewReader("reader data")

	cfg, err := NewRequestConfig(ctx, "POST", "/", body, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Set base URL
	cfg.DefaultBaseURL, _ = url.Parse(server.URL)
	err = cfg.Execute()
	if err != nil {
		t.Fatalf("unexpected error during execute: %v", err)
	}
}
