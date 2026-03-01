package opencode_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dominicnunez/opencode-sdk-go"
)

// injectionPayloads are path parameter values that attempt traversal, query
// injection, or encoded trickery. The SDK relies on server-generated UUIDs
// so these inputs are abnormal — tests document that Go's net/http
// neutralizes traversal via URL normalization.
var injectionPayloads = []struct {
	name             string
	id               string
	wantPathContains string
}{
	{"path traversal is neutralized", "../config", "/config"},
	{"slash in id", "foo/bar", "foo/bar"},
	{"query injection", "id?x=1", "id"},
	{"url-encoded traversal", "%2e%2e%2fconfig", "%2e%2e%2fconfig"},
}

// newInjectionServer returns a test server that captures the request path.
func newInjectionServer(receivedPath *string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.RawPath
		if p == "" {
			p = r.URL.Path
		}
		*receivedPath = p
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`"not found"`))
	}))
}

func TestPathParameterInjection_SessionGet(t *testing.T) {
	for _, tt := range injectionPayloads {
		t.Run(tt.name, func(t *testing.T) {
			var receivedPath string
			server := newInjectionServer(&receivedPath)
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

func TestPathParameterInjection_AuthSet(t *testing.T) {
	for _, tt := range injectionPayloads {
		t.Run(tt.name, func(t *testing.T) {
			var receivedPath string
			server := newInjectionServer(&receivedPath)
			defer server.Close()

			client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, _ = client.Auth.Set(context.Background(), tt.id, &opencode.AuthSetParams{
				Auth: opencode.ApiAuth{Key: "k"},
			})

			if receivedPath == "" {
				t.Fatal("server received no request")
			}
			if !strings.Contains(receivedPath, tt.wantPathContains) {
				t.Errorf("received path %q does not contain %q", receivedPath, tt.wantPathContains)
			}
		})
	}
}

func TestPathParameterInjection_SessionPermissionRespond(t *testing.T) {
	// Test injection in both the session id and permission id parameters.
	for _, tt := range injectionPayloads {
		t.Run("sessionID/"+tt.name, func(t *testing.T) {
			var receivedPath string
			server := newInjectionServer(&receivedPath)
			defer server.Close()

			client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, _ = client.Session.Permissions.Respond(
				context.Background(), tt.id, "safe-perm-id",
				&opencode.SessionPermissionRespondParams{Response: opencode.PermissionResponseOnce},
			)

			if receivedPath == "" {
				t.Fatal("server received no request")
			}
			if !strings.Contains(receivedPath, tt.wantPathContains) {
				t.Errorf("received path %q does not contain %q", receivedPath, tt.wantPathContains)
			}
		})

		t.Run("permissionID/"+tt.name, func(t *testing.T) {
			var receivedPath string
			server := newInjectionServer(&receivedPath)
			defer server.Close()

			client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, _ = client.Session.Permissions.Respond(
				context.Background(), "safe-session-id", tt.id,
				&opencode.SessionPermissionRespondParams{Response: opencode.PermissionResponseOnce},
			)

			if receivedPath == "" {
				t.Fatal("server received no request")
			}
			if !strings.Contains(receivedPath, tt.wantPathContains) {
				t.Errorf("received path %q does not contain %q", receivedPath, tt.wantPathContains)
			}
		})
	}
}

func TestPathParameterInjection_SessionMessage(t *testing.T) {
	// Session.Message has two dynamic segments: session/{id}/message/{messageID}.
	// Test injection in both parameters independently.
	for _, tt := range injectionPayloads {
		t.Run("sessionID/"+tt.name, func(t *testing.T) {
			var receivedPath string
			server := newInjectionServer(&receivedPath)
			defer server.Close()

			client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, _ = client.Session.Message(
				context.Background(), tt.id, "safe-message-id", nil,
			)

			if receivedPath == "" {
				t.Fatal("server received no request")
			}
			if !strings.Contains(receivedPath, tt.wantPathContains) {
				t.Errorf("received path %q does not contain %q", receivedPath, tt.wantPathContains)
			}
		})

		t.Run("messageID/"+tt.name, func(t *testing.T) {
			var receivedPath string
			server := newInjectionServer(&receivedPath)
			defer server.Close()

			client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, _ = client.Session.Message(
				context.Background(), "safe-session-id", tt.id, nil,
			)

			if receivedPath == "" {
				t.Fatal("server received no request")
			}
			if !strings.Contains(receivedPath, tt.wantPathContains) {
				t.Errorf("received path %q does not contain %q", receivedPath, tt.wantPathContains)
			}
		})
	}
}
