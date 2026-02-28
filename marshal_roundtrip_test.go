package opencode

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestMarshalJSON_RoundTrip_PermissionPattern(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{"string pattern", `"*.go"`},
		{"array pattern", `["*.go","*.ts"]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p PermissionPattern
			if err := json.Unmarshal([]byte(tt.json), &p); err != nil {
				t.Fatalf("Unmarshal: %v", err)
			}

			got, err := json.Marshal(p)
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}

			// Compare compacted JSON to ignore whitespace differences
			var wantBuf, gotBuf bytes.Buffer
			if err := json.Compact(&wantBuf, []byte(tt.json)); err != nil {
				t.Fatalf("compact want: %v", err)
			}
			if err := json.Compact(&gotBuf, got); err != nil {
				t.Fatalf("compact got: %v", err)
			}

			if wantBuf.String() != gotBuf.String() {
				t.Errorf("round-trip mismatch:\n  want: %s\n  got:  %s", wantBuf.String(), gotBuf.String())
			}
		})
	}
}

func TestMarshalJSON_RoundTrip_Event(t *testing.T) {
	original := `{"type":"installation.updated","properties":{"version":"1.2.3"}}`

	var event Event
	if err := json.Unmarshal([]byte(original), &event); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	got, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var wantBuf, gotBuf bytes.Buffer
	if err := json.Compact(&wantBuf, []byte(original)); err != nil {
		t.Fatalf("compact want: %v", err)
	}
	if err := json.Compact(&gotBuf, got); err != nil {
		t.Fatalf("compact got: %v", err)
	}

	if wantBuf.String() != gotBuf.String() {
		t.Errorf("round-trip mismatch:\n  want: %s\n  got:  %s", wantBuf.String(), gotBuf.String())
	}
}

func TestMarshalJSON_RoundTrip_Part(t *testing.T) {
	original := `{"type":"text","id":"p1","messageID":"m1","sessionID":"s1","text":"hello"}`

	var part Part
	if err := json.Unmarshal([]byte(original), &part); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	got, err := json.Marshal(part)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var wantBuf, gotBuf bytes.Buffer
	if err := json.Compact(&wantBuf, []byte(original)); err != nil {
		t.Fatalf("compact want: %v", err)
	}
	if err := json.Compact(&gotBuf, got); err != nil {
		t.Fatalf("compact got: %v", err)
	}

	if wantBuf.String() != gotBuf.String() {
		t.Errorf("round-trip mismatch:\n  want: %s\n  got:  %s", wantBuf.String(), gotBuf.String())
	}
}

func TestMarshalJSON_RoundTrip_ConfigMcp(t *testing.T) {
	original := `{"type":"local","command":["/usr/bin/mcp"],"enabled":true,"environment":{"KEY":"val"}}`

	var mcp ConfigMcp
	if err := json.Unmarshal([]byte(original), &mcp); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	got, err := json.Marshal(mcp)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var wantBuf, gotBuf bytes.Buffer
	if err := json.Compact(&wantBuf, []byte(original)); err != nil {
		t.Fatalf("compact want: %v", err)
	}
	if err := json.Compact(&gotBuf, got); err != nil {
		t.Fatalf("compact got: %v", err)
	}

	if wantBuf.String() != gotBuf.String() {
		t.Errorf("round-trip mismatch:\n  want: %s\n  got:  %s", wantBuf.String(), gotBuf.String())
	}
}

func TestMarshalJSON_ZeroValue_ReturnsNull(t *testing.T) {
	tests := []struct {
		name string
		val  interface{}
	}{
		{"Event", Event{}},
		{"Part", Part{}},
		{"PermissionPattern", PermissionPattern{}},
		{"ConfigMcp", ConfigMcp{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.val)
			if err != nil {
				t.Fatalf("Marshal zero-value %s: %v", tt.name, err)
			}
			if string(got) != "null" {
				t.Errorf("expected null for zero-value %s, got %s", tt.name, string(got))
			}
		})
	}
}
