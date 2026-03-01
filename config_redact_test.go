package opencode

import (
	"fmt"
	"strings"
	"testing"
)

func TestOAuth_String_RedactsCredentials(t *testing.T) {
	o := OAuth{
		Type:    AuthTypeOAuth,
		Access:  "secret_access_token",
		Refresh: "secret_refresh_token",
		Expires: 3600,
	}

	s := o.String()
	if strings.Contains(s, "secret_access_token") {
		t.Error("String() leaked Access credential")
	}
	if strings.Contains(s, "secret_refresh_token") {
		t.Error("String() leaked Refresh credential")
	}
	if !strings.Contains(s, "[REDACTED]") {
		t.Error("String() missing [REDACTED] placeholder")
	}
	if !strings.Contains(s, "3600") {
		t.Error("String() should include non-sensitive Expires field")
	}
}

func TestOAuth_GoString_RedactsCredentials(t *testing.T) {
	o := OAuth{Access: "secret", Refresh: "secret"}
	s := fmt.Sprintf("%#v", o)
	if strings.Contains(s, "secret") {
		t.Error("GoString() leaked credentials via fmt #v")
	}
}

func TestApiAuth_String_RedactsKey(t *testing.T) {
	a := ApiAuth{Type: AuthTypeAPI, Key: "sk-secret-key"}

	s := a.String()
	if strings.Contains(s, "sk-secret-key") {
		t.Error("String() leaked Key credential")
	}
	if !strings.Contains(s, "[REDACTED]") {
		t.Error("String() missing [REDACTED] placeholder")
	}
}

func TestApiAuth_GoString_RedactsKey(t *testing.T) {
	a := ApiAuth{Key: "sk-secret-key"}
	s := fmt.Sprintf("%#v", a)
	if strings.Contains(s, "sk-secret-key") {
		t.Error("GoString() leaked Key via fmt #v")
	}
}

func TestWellKnownAuth_String_RedactsCredentials(t *testing.T) {
	w := WellKnownAuth{
		Type:  AuthTypeWellKnown,
		Key:   "secret-key",
		Token: "secret-token",
	}

	s := w.String()
	if strings.Contains(s, "secret-key") {
		t.Error("String() leaked Key credential")
	}
	if strings.Contains(s, "secret-token") {
		t.Error("String() leaked Token credential")
	}
	if !strings.Contains(s, "[REDACTED]") {
		t.Error("String() missing [REDACTED] placeholder")
	}
}

func TestWellKnownAuth_GoString_RedactsCredentials(t *testing.T) {
	w := WellKnownAuth{Key: "secret-key", Token: "secret-token"}
	s := fmt.Sprintf("%#v", w)
	if strings.Contains(s, "secret-key") || strings.Contains(s, "secret-token") {
		t.Error("GoString() leaked credentials via fmt #v")
	}
}

func TestOAuth_FmtV_RedactsCredentials(t *testing.T) {
	o := OAuth{Access: "secret_access", Refresh: "secret_refresh"}
	s := fmt.Sprintf("%v", o)
	if strings.Contains(s, "secret_access") || strings.Contains(s, "secret_refresh") {
		t.Error("fmt v leaked credentials")
	}
}

func TestConfigProviderOptions_String_RedactsAPIKey(t *testing.T) {
	o := ConfigProviderOptions{
		APIKey:  "sk-super-secret",
		BaseURL: "https://api.example.com",
	}

	s := o.String()
	if strings.Contains(s, "sk-super-secret") {
		t.Error("String() leaked APIKey credential")
	}
	if !strings.Contains(s, "[REDACTED]") {
		t.Error("String() missing [REDACTED] placeholder")
	}
	if !strings.Contains(s, "https://api.example.com") {
		t.Error("String() should include non-sensitive BaseURL field")
	}
}

func TestConfigProviderOptions_GoString_RedactsAPIKey(t *testing.T) {
	o := ConfigProviderOptions{APIKey: "sk-super-secret"}
	s := fmt.Sprintf("%#v", o)
	if strings.Contains(s, "sk-super-secret") {
		t.Error("GoString() leaked APIKey via fmt #v")
	}
}

func TestMcpRemoteConfig_String_RedactsHeaders(t *testing.T) {
	r := McpRemoteConfig{
		Type:    McpRemoteConfigTypeRemote,
		URL:     "https://mcp.example.com",
		Enabled: true,
		Headers: map[string]string{
			"Authorization": "Bearer secret-token-123",
			"X-Custom":      "custom-value",
		},
	}

	s := r.String()
	if strings.Contains(s, "secret-token-123") {
		t.Error("String() leaked Authorization header value")
	}
	if strings.Contains(s, "custom-value") {
		t.Error("String() leaked custom header value")
	}
	if !strings.Contains(s, "2 redacted") {
		t.Error("String() should show header count")
	}
	if !strings.Contains(s, "https://mcp.example.com") {
		t.Error("String() should include non-sensitive URL field")
	}
}

func TestMcpRemoteConfig_GoString_RedactsHeaders(t *testing.T) {
	r := McpRemoteConfig{
		Headers: map[string]string{"Authorization": "Bearer secret"},
	}
	s := fmt.Sprintf("%#v", r)
	if strings.Contains(s, "secret") {
		t.Error("GoString() leaked header values via fmt #v")
	}
}

func TestMcpLocalConfig_String_RedactsEnvironment(t *testing.T) {
	c := McpLocalConfig{
		Type:    McpLocalConfigTypeLocal,
		Command: []string{"npx", "mcp-server"},
		Enabled: true,
		Environment: map[string]string{
			"API_KEY":    "sk-secret-key-12345",
			"AUTH_TOKEN": "secret-token-67890",
		},
	}

	s := c.String()
	if strings.Contains(s, "sk-secret-key-12345") {
		t.Error("String() leaked API_KEY value")
	}
	if strings.Contains(s, "secret-token-67890") {
		t.Error("String() leaked AUTH_TOKEN value")
	}
	if !strings.Contains(s, "2 redacted") {
		t.Error("String() should show environment variable count")
	}
	if !strings.Contains(s, "npx") {
		t.Error("String() should include non-sensitive Command field")
	}
}

func TestMcpLocalConfig_GoString_RedactsEnvironment(t *testing.T) {
	c := McpLocalConfig{
		Environment: map[string]string{"SECRET": "hunter2"},
	}
	s := fmt.Sprintf("%#v", c)
	if strings.Contains(s, "hunter2") {
		t.Error("GoString() leaked environment values via fmt #v")
	}
}

func TestConfigLspObject_String_RedactsEnv(t *testing.T) {
	c := ConfigLspObject{
		Command: []string{"gopls"},
		Env: map[string]string{
			"API_KEY":    "sk-secret-key-12345",
			"AUTH_TOKEN": "secret-token-67890",
		},
	}

	s := c.String()
	if strings.Contains(s, "sk-secret-key-12345") {
		t.Error("String() leaked API_KEY value")
	}
	if strings.Contains(s, "secret-token-67890") {
		t.Error("String() leaked AUTH_TOKEN value")
	}
	if !strings.Contains(s, "2 redacted") {
		t.Error("String() should show env variable count")
	}
	if !strings.Contains(s, "gopls") {
		t.Error("String() should include non-sensitive Command field")
	}
}

func TestConfigLspObject_GoString_RedactsEnv(t *testing.T) {
	c := ConfigLspObject{Env: map[string]string{"SECRET": "hunter2"}}
	s := fmt.Sprintf("%#v", c)
	if strings.Contains(s, "hunter2") {
		t.Error("GoString() leaked env values via fmt #v")
	}
}

func TestConfigFormatter_String_RedactsEnvironment(t *testing.T) {
	c := ConfigFormatter{
		Command: []string{"prettier", "--write"},
		Environment: map[string]string{
			"API_KEY": "sk-secret-key-12345",
		},
	}

	s := c.String()
	if strings.Contains(s, "sk-secret-key-12345") {
		t.Error("String() leaked API_KEY value")
	}
	if !strings.Contains(s, "1 redacted") {
		t.Error("String() should show environment variable count")
	}
	if !strings.Contains(s, "prettier") {
		t.Error("String() should include non-sensitive Command field")
	}
}

func TestConfigFormatter_GoString_RedactsEnvironment(t *testing.T) {
	c := ConfigFormatter{Environment: map[string]string{"SECRET": "hunter2"}}
	s := fmt.Sprintf("%#v", c)
	if strings.Contains(s, "hunter2") {
		t.Error("GoString() leaked environment values via fmt #v")
	}
}

func TestConfigExperimentalHookFileEdited_String_RedactsEnvironment(t *testing.T) {
	c := ConfigExperimentalHookFileEdited{
		Command:     []string{"notify", "--hook"},
		Environment: map[string]string{"SECRET_TOKEN": "abc123"},
	}

	s := c.String()
	if strings.Contains(s, "abc123") {
		t.Error("String() leaked SECRET_TOKEN value")
	}
	if !strings.Contains(s, "1 redacted") {
		t.Error("String() should show environment variable count")
	}
	if !strings.Contains(s, "notify") {
		t.Error("String() should include non-sensitive Command field")
	}
}

func TestConfigExperimentalHookFileEdited_GoString_RedactsEnvironment(t *testing.T) {
	c := ConfigExperimentalHookFileEdited{Environment: map[string]string{"SECRET": "hunter2"}}
	s := fmt.Sprintf("%#v", c)
	if strings.Contains(s, "hunter2") {
		t.Error("GoString() leaked environment values via fmt #v")
	}
}

func TestConfigExperimentalHookSessionCompleted_String_RedactsEnvironment(t *testing.T) {
	c := ConfigExperimentalHookSessionCompleted{
		Command:     []string{"cleanup", "--all"},
		Environment: map[string]string{"API_KEY": "secret-key", "TOKEN": "secret-token"},
	}

	s := c.String()
	if strings.Contains(s, "secret-key") {
		t.Error("String() leaked API_KEY value")
	}
	if strings.Contains(s, "secret-token") {
		t.Error("String() leaked TOKEN value")
	}
	if !strings.Contains(s, "2 redacted") {
		t.Error("String() should show environment variable count")
	}
	if !strings.Contains(s, "cleanup") {
		t.Error("String() should include non-sensitive Command field")
	}
}

func TestConfigExperimentalHookSessionCompleted_GoString_RedactsEnvironment(t *testing.T) {
	c := ConfigExperimentalHookSessionCompleted{Environment: map[string]string{"SECRET": "hunter2"}}
	s := fmt.Sprintf("%#v", c)
	if strings.Contains(s, "hunter2") {
		t.Error("GoString() leaked environment values via fmt #v")
	}
}
