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
