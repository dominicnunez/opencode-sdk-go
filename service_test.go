package opencode_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dominicnunez/opencode-sdk-go"
)

// serviceTestHandler records the HTTP method, path, and query string from the
// incoming request, then writes responseBody as JSON.
type serviceTestHandler struct {
	t            *testing.T
	method       string
	path         string
	query        url.Values
	responseBody interface{}
}

func (h *serviceTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.method = r.Method
	h.path = r.URL.Path
	h.query = r.URL.Query()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(h.responseBody); err != nil {
		h.t.Fatalf("failed to encode response: %v", err)
	}
}

func TestAgentService_List(t *testing.T) {
	h := &serviceTestHandler{
		t: t,
		responseBody: []opencode.Agent{
			{Name: "coder", BuiltIn: true, Mode: opencode.AgentModePrimary},
		},
	}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	agents, err := client.Agent.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Agent.List failed: %v", err)
	}
	if h.method != http.MethodGet {
		t.Errorf("method: got %s, want GET", h.method)
	}
	if h.path != "/agent" {
		t.Errorf("path: got %s, want /agent", h.path)
	}
	if len(agents) != 1 || agents[0].Name != "coder" {
		t.Errorf("unexpected agents: %+v", agents)
	}
}

func TestAgentService_List_WithDirectory(t *testing.T) {
	h := &serviceTestHandler{t: t, responseBody: []opencode.Agent{}}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Agent.List(context.Background(), &opencode.AgentListParams{
		Directory: opencode.Ptr("/my/project"),
	})
	if err != nil {
		t.Fatalf("Agent.List failed: %v", err)
	}
	if got := h.query.Get("directory"); got != "/my/project" {
		t.Errorf("directory query param: got %q, want %q", got, "/my/project")
	}
}

func TestAppService_Log(t *testing.T) {
	h := &serviceTestHandler{t: t, responseBody: true}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.App.Log(context.Background(), &opencode.AppLogParams{
		Level:   opencode.LogLevelInfo,
		Message: "test message",
		Service: "test-svc",
	})
	if err != nil {
		t.Fatalf("App.Log failed: %v", err)
	}
	if h.method != http.MethodPost {
		t.Errorf("method: got %s, want POST", h.method)
	}
	if h.path != "/log" {
		t.Errorf("path: got %s, want /log", h.path)
	}
	if !result {
		t.Error("expected true result")
	}
}

func TestProjectService_List(t *testing.T) {
	h := &serviceTestHandler{
		t: t,
		responseBody: []opencode.Project{
			{ID: "proj_1", Worktree: "/home/user/code"},
		},
	}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	projects, err := client.Project.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Project.List failed: %v", err)
	}
	if h.method != http.MethodGet {
		t.Errorf("method: got %s, want GET", h.method)
	}
	if h.path != "/project" {
		t.Errorf("path: got %s, want /project", h.path)
	}
	if len(projects) != 1 || projects[0].ID != "proj_1" {
		t.Errorf("unexpected projects: %+v", projects)
	}
}

func TestProjectService_Current(t *testing.T) {
	h := &serviceTestHandler{
		t:            t,
		responseBody: opencode.Project{ID: "proj_cur", Worktree: "/code"},
	}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	project, err := client.Project.Current(context.Background(), nil)
	if err != nil {
		t.Fatalf("Project.Current failed: %v", err)
	}
	if h.method != http.MethodGet {
		t.Errorf("method: got %s, want GET", h.method)
	}
	if h.path != "/project/current" {
		t.Errorf("path: got %s, want /project/current", h.path)
	}
	if project.ID != "proj_cur" {
		t.Errorf("expected project ID proj_cur, got %s", project.ID)
	}
}

func TestCommandService_List(t *testing.T) {
	h := &serviceTestHandler{
		t: t,
		responseBody: []opencode.Command{
			{Name: "test-cmd", Template: "echo hi"},
		},
	}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	commands, err := client.Command.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Command.List failed: %v", err)
	}
	if h.method != http.MethodGet {
		t.Errorf("method: got %s, want GET", h.method)
	}
	if h.path != "/command" {
		t.Errorf("path: got %s, want /command", h.path)
	}
	if len(commands) != 1 || commands[0].Name != "test-cmd" {
		t.Errorf("unexpected commands: %+v", commands)
	}
}

func TestTuiService_AppendPrompt(t *testing.T) {
	h := &serviceTestHandler{t: t, responseBody: true}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Tui.AppendPrompt(context.Background(), &opencode.TuiAppendPromptParams{
		Text: "hello",
	})
	if err != nil {
		t.Fatalf("Tui.AppendPrompt failed: %v", err)
	}
	if h.method != http.MethodPost {
		t.Errorf("method: got %s, want POST", h.method)
	}
	if h.path != "/tui/append-prompt" {
		t.Errorf("path: got %s, want /tui/append-prompt", h.path)
	}
	if !result {
		t.Error("expected true result")
	}
}

func TestTuiService_ClearPrompt(t *testing.T) {
	h := &serviceTestHandler{t: t, responseBody: true}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Tui.ClearPrompt(context.Background(), nil)
	if err != nil {
		t.Fatalf("Tui.ClearPrompt failed: %v", err)
	}
	if h.method != http.MethodPost {
		t.Errorf("method: got %s, want POST", h.method)
	}
	if h.path != "/tui/clear-prompt" {
		t.Errorf("path: got %s, want /tui/clear-prompt", h.path)
	}
	if !result {
		t.Error("expected true result")
	}
}

func TestTuiService_ExecuteCommand(t *testing.T) {
	h := &serviceTestHandler{t: t, responseBody: true}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Tui.ExecuteCommand(context.Background(), &opencode.TuiExecuteCommandParams{
		Command: "/help",
	})
	if err != nil {
		t.Fatalf("Tui.ExecuteCommand failed: %v", err)
	}
	if h.method != http.MethodPost {
		t.Errorf("method: got %s, want POST", h.method)
	}
	if h.path != "/tui/execute-command" {
		t.Errorf("path: got %s, want /tui/execute-command", h.path)
	}
	if !result {
		t.Error("expected true result")
	}
}

func TestTuiService_OpenHelp(t *testing.T) {
	h := &serviceTestHandler{t: t, responseBody: true}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Tui.OpenHelp(context.Background(), nil)
	if err != nil {
		t.Fatalf("Tui.OpenHelp failed: %v", err)
	}
	if h.method != http.MethodPost {
		t.Errorf("method: got %s, want POST", h.method)
	}
	if h.path != "/tui/open-help" {
		t.Errorf("path: got %s, want /tui/open-help", h.path)
	}
	if !result {
		t.Error("expected true result")
	}
}

func TestTuiService_OpenModels(t *testing.T) {
	h := &serviceTestHandler{t: t, responseBody: true}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Tui.OpenModels(context.Background(), nil)
	if err != nil {
		t.Fatalf("Tui.OpenModels failed: %v", err)
	}
	if h.method != http.MethodPost {
		t.Errorf("method: got %s, want POST", h.method)
	}
	if h.path != "/tui/open-models" {
		t.Errorf("path: got %s, want /tui/open-models", h.path)
	}
	if !result {
		t.Error("expected true result")
	}
}

func TestTuiService_OpenSessions(t *testing.T) {
	h := &serviceTestHandler{t: t, responseBody: true}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Tui.OpenSessions(context.Background(), nil)
	if err != nil {
		t.Fatalf("Tui.OpenSessions failed: %v", err)
	}
	if h.method != http.MethodPost {
		t.Errorf("method: got %s, want POST", h.method)
	}
	if h.path != "/tui/open-sessions" {
		t.Errorf("path: got %s, want /tui/open-sessions", h.path)
	}
	if !result {
		t.Error("expected true result")
	}
}

func TestTuiService_OpenThemes(t *testing.T) {
	h := &serviceTestHandler{t: t, responseBody: true}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Tui.OpenThemes(context.Background(), nil)
	if err != nil {
		t.Fatalf("Tui.OpenThemes failed: %v", err)
	}
	if h.method != http.MethodPost {
		t.Errorf("method: got %s, want POST", h.method)
	}
	if h.path != "/tui/open-themes" {
		t.Errorf("path: got %s, want /tui/open-themes", h.path)
	}
	if !result {
		t.Error("expected true result")
	}
}

func TestTuiService_ShowToast(t *testing.T) {
	h := &serviceTestHandler{t: t, responseBody: true}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Tui.ShowToast(context.Background(), &opencode.TuiShowToastParams{
		Message: "done",
		Variant: opencode.ToastVariantSuccess,
	})
	if err != nil {
		t.Fatalf("Tui.ShowToast failed: %v", err)
	}
	if h.method != http.MethodPost {
		t.Errorf("method: got %s, want POST", h.method)
	}
	if h.path != "/tui/show-toast" {
		t.Errorf("path: got %s, want /tui/show-toast", h.path)
	}
	if !result {
		t.Error("expected true result")
	}
}

func TestTuiService_SubmitPrompt(t *testing.T) {
	h := &serviceTestHandler{t: t, responseBody: true}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Tui.SubmitPrompt(context.Background(), nil)
	if err != nil {
		t.Fatalf("Tui.SubmitPrompt failed: %v", err)
	}
	if h.method != http.MethodPost {
		t.Errorf("method: got %s, want POST", h.method)
	}
	if h.path != "/tui/submit-prompt" {
		t.Errorf("path: got %s, want /tui/submit-prompt", h.path)
	}
	if !result {
		t.Error("expected true result")
	}
}

func TestSessionService_Abort(t *testing.T) {
	h := &serviceTestHandler{t: t, responseBody: true}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Session.Abort(context.Background(), "ses_abc", nil)
	if err != nil {
		t.Fatalf("Session.Abort failed: %v", err)
	}
	if h.method != http.MethodPost {
		t.Errorf("method: got %s, want POST", h.method)
	}
	if h.path != "/session/ses_abc/abort" {
		t.Errorf("path: got %s, want /session/ses_abc/abort", h.path)
	}
	if !result {
		t.Error("expected true result")
	}
}

func TestSessionService_Children(t *testing.T) {
	h := &serviceTestHandler{
		t: t,
		responseBody: []opencode.Session{
			{ID: "child_1"},
		},
	}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	children, err := client.Session.Children(context.Background(), "ses_parent", nil)
	if err != nil {
		t.Fatalf("Session.Children failed: %v", err)
	}
	if h.method != http.MethodGet {
		t.Errorf("method: got %s, want GET", h.method)
	}
	if h.path != "/session/ses_parent/children" {
		t.Errorf("path: got %s, want /session/ses_parent/children", h.path)
	}
	if len(children) != 1 || children[0].ID != "child_1" {
		t.Errorf("unexpected children: %+v", children)
	}
}

func TestSessionService_Share(t *testing.T) {
	h := &serviceTestHandler{
		t:            t,
		responseBody: opencode.Session{ID: "ses_shared"},
	}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	session, err := client.Session.Share(context.Background(), "ses_123", nil)
	if err != nil {
		t.Fatalf("Session.Share failed: %v", err)
	}
	if h.method != http.MethodPost {
		t.Errorf("method: got %s, want POST", h.method)
	}
	if h.path != "/session/ses_123/share" {
		t.Errorf("path: got %s, want /session/ses_123/share", h.path)
	}
	if session.ID != "ses_shared" {
		t.Errorf("expected session ID ses_shared, got %s", session.ID)
	}
}

func TestPathService_Get(t *testing.T) {
	h := &serviceTestHandler{
		t: t,
		responseBody: opencode.Path{
			Config:    "/home/user/.config/opencode",
			Directory: "/home/user/project",
			State:     "/home/user/.local/state/opencode",
			Worktree:  "/home/user/project",
		},
	}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	path, err := client.Path.Get(context.Background(), nil)
	if err != nil {
		t.Fatalf("Path.Get failed: %v", err)
	}
	if h.method != http.MethodGet {
		t.Errorf("method: got %s, want GET", h.method)
	}
	if h.path != "/path" {
		t.Errorf("path: got %s, want /path", h.path)
	}
	if path.Config != "/home/user/.config/opencode" {
		t.Errorf("unexpected config path: %s", path.Config)
	}
}

func TestPathService_Get_WithDirectory(t *testing.T) {
	h := &serviceTestHandler{t: t, responseBody: opencode.Path{}}
	server := httptest.NewServer(h)
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Path.Get(context.Background(), &opencode.PathGetParams{
		Directory: opencode.Ptr("/my/project"),
	})
	if err != nil {
		t.Fatalf("Path.Get failed: %v", err)
	}
	if got := h.query.Get("directory"); got != "/my/project" {
		t.Errorf("directory query param: got %q, want %q", got, "/my/project")
	}
}
