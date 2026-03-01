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
// injection, or encoded trickery. All path parameters are escaped with
// url.PathEscape, so special characters are percent-encoded.
var injectionPayloads = []struct {
	name             string
	id               string
	wantPathContains string
}{
	{"path traversal is escaped", "../config", "..%2Fconfig"},
	{"slash in id is escaped", "foo/bar", "foo%2Fbar"},
	{"query injection is escaped", "id?x=1", "id%3Fx=1"},
	{"url-encoded traversal is double-escaped", "%2e%2e%2fconfig", "%252e%252e%252fconfig"},
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

// authSetPayloads has expected paths after url.PathEscape is applied to the
// provider ID. Unlike server-generated UUIDs, provider IDs are user-defined,
// so Auth.Set escapes them explicitly.
var authSetPayloads = []struct {
	name             string
	id               string
	wantPathContains string
}{
	{"path traversal is escaped", "../config", "..%2Fconfig"},
	{"slash in id is escaped", "foo/bar", "foo%2Fbar"},
	{"query injection is escaped", "id?x=1", "id%3Fx=1"},
	{"url-encoded traversal is double-escaped", "%2e%2e%2fconfig", "%252e%252e%252fconfig"},
	{"simple id is unchanged", "anthropic", "/auth/anthropic"},
}

func TestPathParameterInjection_AuthSet(t *testing.T) {
	for _, tt := range authSetPayloads {
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

func TestPathParameterInjection_SessionMethods(t *testing.T) {
	type callFunc func(client *opencode.Client, ctx context.Context, id string)

	methods := []struct {
		name string
		call callFunc
	}{
		{"Update", func(c *opencode.Client, ctx context.Context, id string) {
			_, _ = c.Session.Update(ctx, id, nil)
		}},
		{"Delete", func(c *opencode.Client, ctx context.Context, id string) {
			_, _ = c.Session.Delete(ctx, id, nil)
		}},
		{"Abort", func(c *opencode.Client, ctx context.Context, id string) {
			_, _ = c.Session.Abort(ctx, id, nil)
		}},
		{"Children", func(c *opencode.Client, ctx context.Context, id string) {
			_, _ = c.Session.Children(ctx, id, nil)
		}},
		{"Messages", func(c *opencode.Client, ctx context.Context, id string) {
			_, _ = c.Session.Messages(ctx, id, nil)
		}},
		{"Share", func(c *opencode.Client, ctx context.Context, id string) {
			_, _ = c.Session.Share(ctx, id, nil)
		}},
		{"Diff", func(c *opencode.Client, ctx context.Context, id string) {
			_, _ = c.Session.Diff(ctx, id, nil)
		}},
		{"Fork", func(c *opencode.Client, ctx context.Context, id string) {
			_, _ = c.Session.Fork(ctx, id, nil)
		}},
		{"Todo", func(c *opencode.Client, ctx context.Context, id string) {
			_, _ = c.Session.Todo(ctx, id, nil)
		}},
		{"Unrevert", func(c *opencode.Client, ctx context.Context, id string) {
			_, _ = c.Session.Unrevert(ctx, id, nil)
		}},
		{"Unshare", func(c *opencode.Client, ctx context.Context, id string) {
			_, _ = c.Session.Unshare(ctx, id, nil)
		}},
		{"Command", func(c *opencode.Client, ctx context.Context, id string) {
			_, _ = c.Session.Command(ctx, id, &opencode.SessionCommandParams{})
		}},
		{"Init", func(c *opencode.Client, ctx context.Context, id string) {
			_, _ = c.Session.Init(ctx, id, &opencode.SessionInitParams{})
		}},
		{"Prompt", func(c *opencode.Client, ctx context.Context, id string) {
			_, _ = c.Session.Prompt(ctx, id, &opencode.SessionPromptParams{})
		}},
		{"Revert", func(c *opencode.Client, ctx context.Context, id string) {
			_, _ = c.Session.Revert(ctx, id, &opencode.SessionRevertParams{})
		}},
		{"Shell", func(c *opencode.Client, ctx context.Context, id string) {
			_, _ = c.Session.Shell(ctx, id, &opencode.SessionShellParams{})
		}},
		{"Summarize", func(c *opencode.Client, ctx context.Context, id string) {
			_, _ = c.Session.Summarize(ctx, id, &opencode.SessionSummarizeParams{})
		}},
	}

	for _, method := range methods {
		for _, tt := range injectionPayloads {
			t.Run(method.name+"/"+tt.name, func(t *testing.T) {
				var receivedPath string
				server := newInjectionServer(&receivedPath)
				defer server.Close()

				client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
				if err != nil {
					t.Fatalf("failed to create client: %v", err)
				}

				method.call(client, context.Background(), tt.id)

				if receivedPath == "" {
					t.Fatal("server received no request")
				}
				if !strings.Contains(receivedPath, tt.wantPathContains) {
					t.Errorf("received path %q does not contain %q", receivedPath, tt.wantPathContains)
				}
			})
		}
	}
}
