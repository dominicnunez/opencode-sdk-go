package opencode_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dominicnunez/opencode-sdk-go"
)

// TestPathParameterInjection verifies that path parameters containing special
// characters (traversal sequences, query strings, encoded chars) are sent
// to the server as-is. The SDK relies on server-generated UUIDs for IDs, so
// these inputs are abnormal — the test documents current behavior.
func TestPathParameterInjection(t *testing.T) {
	tests := []struct {
		name             string
		id               string
		wantPathContains string
	}{
		// Go's net/http resolves "../config" to "/config" — traversal
		// is neutralized by URL normalization, not by the SDK.
		{"path traversal is neutralized", "../config", "/config"},
		{"slash in id", "foo/bar", "foo/bar"},
		{"query injection", "id?x=1", "id"},
		{"url-encoded traversal", "%2e%2e%2fconfig", "%2e%2e%2fconfig"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var receivedPath string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedPath = r.URL.RawPath
				if receivedPath == "" {
					receivedPath = r.URL.Path
				}
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`"not found"`))
			}))
			defer server.Close()

			client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, _ = client.Session.Get(context.Background(), tt.id, nil)

			if receivedPath == "" {
				t.Fatal("server received no request")
			}
			if !strings.Contains(receivedPath, tt.wantPathContains) {
				t.Errorf("received path %q does not contain %q", receivedPath, tt.wantPathContains)
			}
		})
	}
}
