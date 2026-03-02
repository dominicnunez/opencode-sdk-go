package opencode

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
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
