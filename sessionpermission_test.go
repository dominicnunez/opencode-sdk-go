package opencode_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/dominicnunez/opencode-sdk-go"
	"github.com/dominicnunez/opencode-sdk-go/internal/testutil"
)

func TestSessionPermissionRespondWithOptionalParams(t *testing.T) {
	t.Skip("Prism tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client, err := opencode.NewClient(opencode.WithBaseURL(baseURL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Session.Permissions.Respond(
		context.TODO(),
		"id",
		"permissionID",
		&opencode.SessionPermissionRespondParams{
			Response:  opencode.PermissionResponseOnce,
			Directory: opencode.Ptr("directory"),
		},
	)
	if err != nil {
		var apierr *opencode.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestPermissionPattern_AsString(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		want     string
		wantOk   bool
	}{
		{
			name:     "valid string pattern",
			jsonData: `"*.go"`,
			want:     "*.go",
			wantOk:   true,
		},
		{
			name:     "empty string pattern",
			jsonData: `""`,
			want:     "",
			wantOk:   true,
		},
		{
			name:     "array pattern returns false",
			jsonData: `["*.go", "*.ts"]`,
			want:     "",
			wantOk:   false,
		},
		{
			name:     "object pattern returns false",
			jsonData: `{"foo": "bar"}`,
			want:     "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pattern opencode.PermissionPattern
			if err := json.Unmarshal([]byte(tt.jsonData), &pattern); err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}

			got, ok := pattern.AsString()
			if ok != tt.wantOk {
				t.Errorf("AsString() ok = %v, want %v", ok, tt.wantOk)
			}
			if got != tt.want {
				t.Errorf("AsString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPermissionPattern_AsArray(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		want     []string
		wantOk   bool
	}{
		{
			name:     "valid array pattern",
			jsonData: `["*.go", "*.ts"]`,
			want:     []string{"*.go", "*.ts"},
			wantOk:   true,
		},
		{
			name:     "single element array",
			jsonData: `["*.go"]`,
			want:     []string{"*.go"},
			wantOk:   true,
		},
		{
			name:     "empty array",
			jsonData: `[]`,
			want:     []string{},
			wantOk:   true,
		},
		{
			name:     "string pattern returns false",
			jsonData: `"*.go"`,
			want:     nil,
			wantOk:   false,
		},
		{
			name:     "object pattern returns false",
			jsonData: `{"foo": "bar"}`,
			want:     nil,
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pattern opencode.PermissionPattern
			if err := json.Unmarshal([]byte(tt.jsonData), &pattern); err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}

			got, ok := pattern.AsArray()
			if ok != tt.wantOk {
				t.Errorf("AsArray() ok = %v, want %v", ok, tt.wantOk)
			}
			if ok && !stringSlicesEqual(got, tt.want) {
				t.Errorf("AsArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermissionPattern_InvalidJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
	}{
		{
			name:     "invalid JSON",
			jsonData: `{invalid}`,
		},
		{
			name:     "truncated JSON",
			jsonData: `["*.go"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pattern opencode.PermissionPattern
			// UnmarshalJSON just stores raw bytes, so it won't error
			// The error will occur when calling AsString() or AsArray()
			if err := json.Unmarshal([]byte(tt.jsonData), &pattern); err == nil {
				// Try to access the value - should fail
				if _, ok := pattern.AsString(); ok {
					t.Error("AsString() succeeded on invalid JSON, expected false")
				}
				if _, ok := pattern.AsArray(); ok {
					t.Error("AsArray() succeeded on invalid JSON, expected false")
				}
			}
		})
	}
}

func TestPermission_UnmarshalWithStringPattern(t *testing.T) {
	jsonData := `{
		"id": "perm-123",
		"messageID": "msg-456",
		"metadata": {"key": "value"},
		"sessionID": "sess-789",
		"time": {"created": 1234567890.5},
		"title": "File Access",
		"type": "file",
		"callID": "call-abc",
		"pattern": "*.go"
	}`

	var perm opencode.Permission
	if err := json.Unmarshal([]byte(jsonData), &perm); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if perm.ID != "perm-123" {
		t.Errorf("ID = %q, want %q", perm.ID, "perm-123")
	}
	if perm.Pattern == nil {
		t.Fatal("Pattern is nil")
	}

	pattern, ok := perm.Pattern.AsString()
	if !ok {
		t.Fatal("Pattern.AsString() returned false")
	}
	if pattern != "*.go" {
		t.Errorf("Pattern = %q, want %q", pattern, "*.go")
	}
}

func TestPermission_UnmarshalWithArrayPattern(t *testing.T) {
	jsonData := `{
		"id": "perm-123",
		"messageID": "msg-456",
		"metadata": {"key": "value"},
		"sessionID": "sess-789",
		"time": {"created": 1234567890.5},
		"title": "File Access",
		"type": "file",
		"callID": "call-abc",
		"pattern": ["*.go", "*.ts", "*.js"]
	}`

	var perm opencode.Permission
	if err := json.Unmarshal([]byte(jsonData), &perm); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if perm.ID != "perm-123" {
		t.Errorf("ID = %q, want %q", perm.ID, "perm-123")
	}
	if perm.Pattern == nil {
		t.Fatal("Pattern is nil")
	}

	patterns, ok := perm.Pattern.AsArray()
	if !ok {
		t.Fatal("Pattern.AsArray() returned false")
	}
	want := []string{"*.go", "*.ts", "*.js"}
	if !stringSlicesEqual(patterns, want) {
		t.Errorf("Pattern = %v, want %v", patterns, want)
	}
}

func TestPermission_UnmarshalWithoutPattern(t *testing.T) {
	jsonData := `{
		"id": "perm-123",
		"messageID": "msg-456",
		"metadata": {},
		"sessionID": "sess-789",
		"time": {"created": 1234567890.5},
		"title": "File Access",
		"type": "file"
	}`

	var perm opencode.Permission
	if err := json.Unmarshal([]byte(jsonData), &perm); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if perm.Pattern != nil {
		t.Errorf("Pattern = %v, want nil", perm.Pattern)
	}
}

// Helper function to compare string slices
func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
