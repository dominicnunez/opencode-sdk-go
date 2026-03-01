package opencode_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/dominicnunez/opencode-sdk-go"
)

func TestPermissionPattern_AsString(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		want     string
		wantErr  bool
	}{
		{
			name:     "valid string pattern",
			jsonData: `"*.go"`,
			want:     "*.go",
			wantErr:  false,
		},
		{
			name:     "empty string pattern",
			jsonData: `""`,
			want:     "",
			wantErr:  false,
		},
		{
			name:     "array pattern returns error",
			jsonData: `["*.go", "*.ts"]`,
			want:     "",
			wantErr:  true,
		},
		{
			name:     "object pattern returns error",
			jsonData: `{"foo": "bar"}`,
			want:     "",
			wantErr:  true,
		},
		{
			name:     "number returns ErrWrongVariant",
			jsonData: `123`,
			want:     "",
			wantErr:  true,
		},
		{
			name:     "boolean returns ErrWrongVariant",
			jsonData: `true`,
			want:     "",
			wantErr:  true,
		},
		{
			name:     "null returns ErrWrongVariant",
			jsonData: `null`,
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pattern opencode.PermissionPattern
			if err := json.Unmarshal([]byte(tt.jsonData), &pattern); err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}

			got, err := pattern.AsString()
			if (err != nil) != tt.wantErr {
				t.Errorf("AsString() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, opencode.ErrWrongVariant) {
				t.Errorf("AsString() error = %v, want ErrWrongVariant", err)
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
		wantErr  bool
	}{
		{
			name:     "valid array pattern",
			jsonData: `["*.go", "*.ts"]`,
			want:     []string{"*.go", "*.ts"},
			wantErr:  false,
		},
		{
			name:     "single element array",
			jsonData: `["*.go"]`,
			want:     []string{"*.go"},
			wantErr:  false,
		},
		{
			name:     "empty array",
			jsonData: `[]`,
			want:     []string{},
			wantErr:  false,
		},
		{
			name:     "string pattern returns error",
			jsonData: `"*.go"`,
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "object pattern returns error",
			jsonData: `{"foo": "bar"}`,
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "number returns ErrWrongVariant",
			jsonData: `123`,
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "boolean returns ErrWrongVariant",
			jsonData: `false`,
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "null returns ErrWrongVariant",
			jsonData: `null`,
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pattern opencode.PermissionPattern
			if err := json.Unmarshal([]byte(tt.jsonData), &pattern); err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}

			got, err := pattern.AsArray()
			if (err != nil) != tt.wantErr {
				t.Errorf("AsArray() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, opencode.ErrWrongVariant) {
				t.Errorf("AsArray() error = %v, want ErrWrongVariant", err)
			}
			if err == nil && !stringSlicesEqual(got, tt.want) {
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
			if err := json.Unmarshal([]byte(tt.jsonData), &pattern); err == nil {
				// Try to access the value - should fail
				if _, err := pattern.AsString(); err == nil {
					t.Error("AsString() succeeded on invalid JSON, expected error")
				}
				if _, err := pattern.AsArray(); err == nil {
					t.Error("AsArray() succeeded on invalid JSON, expected error")
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

	pattern, err := perm.Pattern.AsString()
	if err != nil {
		t.Fatalf("Pattern.AsString() returned error: %v", err)
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

	patterns, err := perm.Pattern.AsArray()
	if err != nil {
		t.Fatalf("Pattern.AsArray() returned error: %v", err)
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

func TestPermissionResponse_IsKnown(t *testing.T) {
	tests := []struct {
		value opencode.PermissionResponse
		known bool
	}{
		{opencode.PermissionResponseOnce, true},
		{opencode.PermissionResponseAlways, true},
		{opencode.PermissionResponseReject, true},
		{opencode.PermissionResponse(""), false},
		{opencode.PermissionResponse("invalid"), false},
		{opencode.PermissionResponse("ONCE"), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.value), func(t *testing.T) {
			if got := tt.value.IsKnown(); got != tt.known {
				t.Errorf("PermissionResponse(%q).IsKnown() = %v, want %v", tt.value, got, tt.known)
			}
		})
	}
}

func TestPermissionPattern_MarshalJSON_Roundtrip(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
	}{
		{"string variant", `"*.go"`},
		{"array variant", `["*.go","*.ts"]`},
		{"empty array", `[]`},
		{"empty string", `""`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p1 opencode.PermissionPattern
			if err := json.Unmarshal([]byte(tt.jsonData), &p1); err != nil {
				t.Fatalf("Unmarshal: %v", err)
			}

			data, err := p1.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON: %v", err)
			}

			var p2 opencode.PermissionPattern
			if err := json.Unmarshal(data, &p2); err != nil {
				t.Fatalf("Unmarshal roundtrip: %v", err)
			}

			data2, err := p2.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON roundtrip: %v", err)
			}

			if string(data) != string(data2) {
				t.Errorf("roundtrip mismatch: %s != %s", data, data2)
			}
		})
	}
}

func TestPermissionPattern_MarshalJSON_NilRaw(t *testing.T) {
	var p opencode.PermissionPattern
	data, err := p.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}
	if string(data) != "null" {
		t.Errorf("MarshalJSON(nil raw) = %s, want null", data)
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
