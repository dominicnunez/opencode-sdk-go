package opencode_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dominicnunez/opencode-sdk-go"
)

// allowedPathPayloads are inputs that should remain within a single path
// segment after escaping and decoding by net/http.
var allowedPathPayloads = []struct {
	name string
	id   string
}{
	{"simple id", "safe-id"},
	{"query injection is escaped", "id?x=1"},
	{"encoded traversal string stays literal", "%2e%2e%2fconfig"},
}

var rejectedPathPayloads = []string{".", "..", "../config", "foo/bar", "..\\config", "foo\\bar"}

type requestCapture struct {
	path    string
	rawPath string
	calls   int
}

// newInjectionServer returns a test server that captures the request path.
func newInjectionServer(capture *requestCapture) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capture.calls++
		decodedPath, err := url.PathUnescape(r.URL.Path)
		if err != nil {
			capture.path = r.URL.Path
		} else {
			capture.path = decodedPath
		}
		capture.rawPath = r.URL.RawPath
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`"not found"`))
	}))
}

func assertRequestPath(t *testing.T, capture requestCapture, want string) {
	t.Helper()
	if capture.calls != 1 {
		t.Fatalf("expected exactly one request, got %d", capture.calls)
	}
	if capture.path != want {
		t.Fatalf("expected path %q, got %q", want, capture.path)
	}
}

func assertDotSegmentRejected(t *testing.T, err error, capture requestCapture, id string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error for traversal/separator payload %q", id)
	}
	if capture.calls != 0 {
		t.Fatalf("expected no outbound request for traversal/separator payload %q, got %d", id, capture.calls)
	}
}

func TestPathParameterInjection_SessionGet(t *testing.T) {
	for _, tt := range allowedPathPayloads {
		t.Run(tt.name, func(t *testing.T) {
			capture := requestCapture{}
			server := newInjectionServer(&capture)
			defer server.Close()

			client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, _ = client.Session.Get(context.Background(), tt.id, nil)

			assertRequestPath(t, capture, "/session/"+tt.id)
		})
	}
}

func TestPathParameterInjection_SessionGet_RejectsTraversalOrSeparators(t *testing.T) {
	for _, id := range rejectedPathPayloads {
		t.Run(id, func(t *testing.T) {
			capture := requestCapture{}
			server := newInjectionServer(&capture)
			defer server.Close()

			client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, err = client.Session.Get(context.Background(), id, nil)
			assertDotSegmentRejected(t, err, capture, id)
		})
	}
}

func TestPathParameterInjection_AuthSet(t *testing.T) {
	for _, tt := range append(allowedPathPayloads, struct {
		name string
		id   string
	}{name: "simple id is unchanged", id: "anthropic"}) {
		t.Run(tt.name, func(t *testing.T) {
			capture := requestCapture{}
			server := newInjectionServer(&capture)
			defer server.Close()

			client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, _ = client.Auth.Set(context.Background(), tt.id, &opencode.AuthSetParams{
				Auth: opencode.ApiAuth{Key: "k"},
			})

			assertRequestPath(t, capture, "/auth/"+tt.id)
		})
	}
}

func TestPathParameterInjection_AuthSet_RejectsTraversalOrSeparators(t *testing.T) {
	for _, id := range rejectedPathPayloads {
		t.Run(id, func(t *testing.T) {
			capture := requestCapture{}
			server := newInjectionServer(&capture)
			defer server.Close()

			client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, err = client.Auth.Set(context.Background(), id, &opencode.AuthSetParams{
				Auth: opencode.ApiAuth{Key: "k"},
			})
			assertDotSegmentRejected(t, err, capture, id)
		})
	}
}

func TestPathParameterInjection_SessionPermissionRespond(t *testing.T) {
	for _, tt := range allowedPathPayloads {
		t.Run("sessionID/"+tt.name, func(t *testing.T) {
			capture := requestCapture{}
			server := newInjectionServer(&capture)
			defer server.Close()

			client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, _ = client.Session.Permissions.Respond(
				context.Background(), tt.id, "safe-perm-id",
				&opencode.SessionPermissionRespondParams{Response: opencode.PermissionResponseOnce},
			)

			assertRequestPath(
				t,
				capture,
				"/session/"+tt.id+"/permissions/safe-perm-id",
			)
		})

		t.Run("permissionID/"+tt.name, func(t *testing.T) {
			capture := requestCapture{}
			server := newInjectionServer(&capture)
			defer server.Close()

			client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, _ = client.Session.Permissions.Respond(
				context.Background(), "safe-session-id", tt.id,
				&opencode.SessionPermissionRespondParams{Response: opencode.PermissionResponseOnce},
			)

			assertRequestPath(
				t,
				capture,
				"/session/safe-session-id/permissions/"+tt.id,
			)
		})
	}
}

func TestPathParameterInjection_SessionPermissionRespond_RejectsTraversalOrSeparators(t *testing.T) {
	for _, id := range rejectedPathPayloads {
		t.Run("sessionID/"+id, func(t *testing.T) {
			capture := requestCapture{}
			server := newInjectionServer(&capture)
			defer server.Close()

			client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, err = client.Session.Permissions.Respond(
				context.Background(), id, "safe-perm-id",
				&opencode.SessionPermissionRespondParams{Response: opencode.PermissionResponseOnce},
			)
			assertDotSegmentRejected(t, err, capture, id)
		})

		t.Run("permissionID/"+id, func(t *testing.T) {
			capture := requestCapture{}
			server := newInjectionServer(&capture)
			defer server.Close()

			client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, err = client.Session.Permissions.Respond(
				context.Background(), "safe-session-id", id,
				&opencode.SessionPermissionRespondParams{Response: opencode.PermissionResponseOnce},
			)
			assertDotSegmentRejected(t, err, capture, id)
		})
	}
}

func TestPathParameterInjection_SessionMessage(t *testing.T) {
	for _, tt := range allowedPathPayloads {
		t.Run("sessionID/"+tt.name, func(t *testing.T) {
			capture := requestCapture{}
			server := newInjectionServer(&capture)
			defer server.Close()

			client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, _ = client.Session.Message(
				context.Background(), tt.id, "safe-message-id", nil,
			)

			assertRequestPath(t, capture, "/session/"+tt.id+"/message/safe-message-id")
		})

		t.Run("messageID/"+tt.name, func(t *testing.T) {
			capture := requestCapture{}
			server := newInjectionServer(&capture)
			defer server.Close()

			client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, _ = client.Session.Message(
				context.Background(), "safe-session-id", tt.id, nil,
			)

			assertRequestPath(t, capture, "/session/safe-session-id/message/"+tt.id)
		})
	}
}

func TestPathParameterInjection_SessionMessage_RejectsTraversalOrSeparators(t *testing.T) {
	for _, id := range rejectedPathPayloads {
		t.Run("sessionID/"+id, func(t *testing.T) {
			capture := requestCapture{}
			server := newInjectionServer(&capture)
			defer server.Close()

			client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, err = client.Session.Message(
				context.Background(), id, "safe-message-id", nil,
			)
			assertDotSegmentRejected(t, err, capture, id)
		})

		t.Run("messageID/"+id, func(t *testing.T) {
			capture := requestCapture{}
			server := newInjectionServer(&capture)
			defer server.Close()

			client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, err = client.Session.Message(
				context.Background(), "safe-session-id", id, nil,
			)
			assertDotSegmentRejected(t, err, capture, id)
		})
	}
}

func TestPathParameterInjection_SessionMethods(t *testing.T) {
	type callFunc func(client *opencode.Client, ctx context.Context, id string) error

	methods := []struct {
		name       string
		suffixPath string
		call       callFunc
	}{
		{"Update", "", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Update(ctx, id, nil)
			return err
		}},
		{"Delete", "", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Delete(ctx, id, nil)
			return err
		}},
		{"Abort", "/abort", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Abort(ctx, id, nil)
			return err
		}},
		{"Children", "/children", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Children(ctx, id, nil)
			return err
		}},
		{"Messages", "/message", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Messages(ctx, id, nil)
			return err
		}},
		{"Share", "/share", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Share(ctx, id, nil)
			return err
		}},
		{"Diff", "/diff", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Diff(ctx, id, nil)
			return err
		}},
		{"Fork", "/fork", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Fork(ctx, id, nil)
			return err
		}},
		{"Todo", "/todo", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Todo(ctx, id, nil)
			return err
		}},
		{"Unrevert", "/unrevert", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Unrevert(ctx, id, nil)
			return err
		}},
		{"Unshare", "/share", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Unshare(ctx, id, nil)
			return err
		}},
		{"Command", "/command", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Command(ctx, id, &opencode.SessionCommandParams{Command: "/help", Arguments: "about"})
			return err
		}},
		{"Init", "/init", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Init(ctx, id, &opencode.SessionInitParams{
				MessageID:  "msg_1",
				ModelID:    "model_1",
				ProviderID: "provider_1",
			})
			return err
		}},
		{"Prompt", "/message", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Prompt(ctx, id, &opencode.SessionPromptParams{
				Parts: []opencode.SessionPromptParamsPartUnion{
					opencode.TextPartInputParam{Type: opencode.TextPartInputTypeText, Text: "hello"},
				},
			})
			return err
		}},
		{"Revert", "/revert", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Revert(ctx, id, &opencode.SessionRevertParams{MessageID: "msg_1"})
			return err
		}},
		{"Shell", "/shell", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Shell(ctx, id, &opencode.SessionShellParams{
				Agent:   "bash",
				Command: "pwd",
			})
			return err
		}},
		{"Summarize", "/summarize", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Summarize(ctx, id, &opencode.SessionSummarizeParams{
				ModelID:    "model_1",
				ProviderID: "provider_1",
			})
			return err
		}},
	}

	for _, method := range methods {
		for _, tt := range allowedPathPayloads {
			t.Run(method.name+"/"+tt.name, func(t *testing.T) {
				capture := requestCapture{}
				server := newInjectionServer(&capture)
				defer server.Close()

				client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
				if err != nil {
					t.Fatalf("failed to create client: %v", err)
				}

				_ = method.call(client, context.Background(), tt.id)

				assertRequestPath(
					t,
					capture,
					"/session/"+tt.id+method.suffixPath,
				)
			})
		}
	}
}

func TestPathParameterInjection_SessionMethods_RejectsTraversalOrSeparators(t *testing.T) {
	type callFunc func(client *opencode.Client, ctx context.Context, id string) error

	methods := []struct {
		name string
		call callFunc
	}{
		{"Update", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Update(ctx, id, nil)
			return err
		}},
		{"Delete", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Delete(ctx, id, nil)
			return err
		}},
		{"Abort", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Abort(ctx, id, nil)
			return err
		}},
		{"Children", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Children(ctx, id, nil)
			return err
		}},
		{"Messages", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Messages(ctx, id, nil)
			return err
		}},
		{"Share", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Share(ctx, id, nil)
			return err
		}},
		{"Diff", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Diff(ctx, id, nil)
			return err
		}},
		{"Fork", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Fork(ctx, id, nil)
			return err
		}},
		{"Todo", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Todo(ctx, id, nil)
			return err
		}},
		{"Unrevert", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Unrevert(ctx, id, nil)
			return err
		}},
		{"Unshare", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Unshare(ctx, id, nil)
			return err
		}},
		{"Command", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Command(ctx, id, &opencode.SessionCommandParams{Command: "/help", Arguments: "about"})
			return err
		}},
		{"Init", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Init(ctx, id, &opencode.SessionInitParams{
				MessageID:  "msg_1",
				ModelID:    "model_1",
				ProviderID: "provider_1",
			})
			return err
		}},
		{"Prompt", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Prompt(ctx, id, &opencode.SessionPromptParams{
				Parts: []opencode.SessionPromptParamsPartUnion{
					opencode.TextPartInputParam{Type: opencode.TextPartInputTypeText, Text: "hello"},
				},
			})
			return err
		}},
		{"Revert", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Revert(ctx, id, &opencode.SessionRevertParams{MessageID: "msg_1"})
			return err
		}},
		{"Shell", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Shell(ctx, id, &opencode.SessionShellParams{
				Agent:   "bash",
				Command: "pwd",
			})
			return err
		}},
		{"Summarize", func(c *opencode.Client, ctx context.Context, id string) error {
			_, err := c.Session.Summarize(ctx, id, &opencode.SessionSummarizeParams{
				ModelID:    "model_1",
				ProviderID: "provider_1",
			})
			return err
		}},
	}

	for _, method := range methods {
		for _, id := range rejectedPathPayloads {
			t.Run(method.name+"/"+id, func(t *testing.T) {
				capture := requestCapture{}
				server := newInjectionServer(&capture)
				defer server.Close()

				client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
				if err != nil {
					t.Fatalf("failed to create client: %v", err)
				}

				err = method.call(client, context.Background(), id)
				assertDotSegmentRejected(t, err, capture, id)
			})
		}
	}
}
