package opencode

import (
	"encoding/json"
	"testing"
)

func TestConfigLsp_AsDisabled_ValidDisabled(t *testing.T) {
	jsonData := `{"disabled": true}`

	var lsp ConfigLsp
	if err := json.Unmarshal([]byte(jsonData), &lsp); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	disabled, ok := lsp.AsDisabled()
	if !ok {
		t.Fatal("AsDisabled() returned false, expected true")
	}

	if disabled == nil {
		t.Fatal("AsDisabled() returned nil")
	}

	if !bool(disabled.Disabled) {
		t.Errorf("Disabled = %v, want true", disabled.Disabled)
	}
}

func TestConfigLsp_AsObject_ValidObject(t *testing.T) {
	jsonData := `{
		"command": ["gopls", "-mode=stdio"],
		"disabled": false,
		"env": {"GOPATH": "/go"},
		"extensions": [".go", ".mod"],
		"initialization": {"usePlaceholders": true}
	}`

	var lsp ConfigLsp
	if err := json.Unmarshal([]byte(jsonData), &lsp); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	obj, ok := lsp.AsObject()
	if !ok {
		t.Fatal("AsObject() returned false, expected true")
	}

	if obj == nil {
		t.Fatal("AsObject() returned nil")
	}

	if len(obj.Command) != 2 {
		t.Errorf("Command length = %d, want 2", len(obj.Command))
	}

	if obj.Command[0] != "gopls" {
		t.Errorf("Command[0] = %s, want gopls", obj.Command[0])
	}

	if obj.Command[1] != "-mode=stdio" {
		t.Errorf("Command[1] = %s, want -mode=stdio", obj.Command[1])
	}

	if obj.Disabled {
		t.Errorf("Disabled = %v, want false", obj.Disabled)
	}

	if obj.Env["GOPATH"] != "/go" {
		t.Errorf("Env[GOPATH] = %s, want /go", obj.Env["GOPATH"])
	}

	if len(obj.Extensions) != 2 {
		t.Errorf("Extensions length = %d, want 2", len(obj.Extensions))
	}

	if obj.Extensions[0] != ".go" {
		t.Errorf("Extensions[0] = %s, want .go", obj.Extensions[0])
	}

	if val, ok := obj.Initialization["usePlaceholders"].(bool); !ok || !val {
		t.Errorf("Initialization[usePlaceholders] = %v, want true", obj.Initialization["usePlaceholders"])
	}
}

func TestConfigLsp_AsObject_MinimalCommand(t *testing.T) {
	jsonData := `{"command": ["gopls"]}`

	var lsp ConfigLsp
	if err := json.Unmarshal([]byte(jsonData), &lsp); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	obj, ok := lsp.AsObject()
	if !ok {
		t.Fatal("AsObject() returned false, expected true")
	}

	if obj == nil {
		t.Fatal("AsObject() returned nil")
	}

	if len(obj.Command) != 1 {
		t.Errorf("Command length = %d, want 1", len(obj.Command))
	}

	if obj.Command[0] != "gopls" {
		t.Errorf("Command[0] = %s, want gopls", obj.Command[0])
	}
}

func TestConfigLsp_WrongType_DisabledAsObject(t *testing.T) {
	jsonData := `{"disabled": true}`

	var lsp ConfigLsp
	if err := json.Unmarshal([]byte(jsonData), &lsp); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	obj, ok := lsp.AsObject()
	if ok {
		t.Errorf("AsObject() returned true, expected false for disabled config")
	}

	if obj != nil {
		t.Errorf("AsObject() returned non-nil, expected nil for disabled config")
	}
}

func TestConfigLsp_WrongType_ObjectAsDisabled(t *testing.T) {
	jsonData := `{"command": ["gopls"], "disabled": false}`

	var lsp ConfigLsp
	if err := json.Unmarshal([]byte(jsonData), &lsp); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	disabled, ok := lsp.AsDisabled()
	if ok {
		t.Errorf("AsDisabled() returned true, expected false for object config")
	}

	if disabled != nil {
		t.Errorf("AsDisabled() returned non-nil, expected nil for object config")
	}
}

func TestConfigLsp_InvalidJSON(t *testing.T) {
	jsonData := `{"disabled": "not a boolean"}`

	var lsp ConfigLsp
	// UnmarshalJSON for ConfigLsp just stores raw JSON, so it won't fail
	if err := json.Unmarshal([]byte(jsonData), &lsp); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// But AsDisabled should fail when trying to unmarshal the invalid JSON
	disabled, ok := lsp.AsDisabled()
	if ok {
		t.Errorf("AsDisabled() returned true for invalid JSON, expected false")
	}

	if disabled != nil {
		t.Errorf("AsDisabled() returned non-nil for invalid JSON, expected nil")
	}
}

func TestConfigLsp_EmptyJSON(t *testing.T) {
	jsonData := `{}`

	var lsp ConfigLsp
	if err := json.Unmarshal([]byte(jsonData), &lsp); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Empty JSON should not match either type
	obj, ok := lsp.AsObject()
	if ok {
		t.Errorf("AsObject() returned true for empty JSON, expected false")
	}
	if obj != nil {
		t.Errorf("AsObject() returned non-nil for empty JSON")
	}

	disabled, ok := lsp.AsDisabled()
	if ok {
		t.Errorf("AsDisabled() returned true for empty JSON, expected false")
	}
	if disabled != nil {
		t.Errorf("AsDisabled() returned non-nil for empty JSON")
	}
}

func TestConfigLsp_MalformedJSON(t *testing.T) {
	jsonData := `{"command": ["gopls", "disabled": true}`

	var lsp ConfigLsp
	// UnmarshalJSON stores raw, but it's malformed
	if err := json.Unmarshal([]byte(jsonData), &lsp); err == nil {
		t.Fatal("Expected unmarshal to fail for malformed JSON")
	}
}

func TestConfigLspDisabledDisabled_IsKnown(t *testing.T) {
	tests := []struct {
		name  string
		value ConfigLspDisabledDisabled
		want  bool
	}{
		{"True is known", ConfigLspDisabledDisabledTrue, true},
		{"False is unknown", ConfigLspDisabledDisabled(false), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsKnown(); got != tt.want {
				t.Errorf("IsKnown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigLsp_ObjectWithDisabledTrue(t *testing.T) {
	// Object config can have disabled=true along with command
	jsonData := `{"command": ["gopls"], "disabled": true}`

	var lsp ConfigLsp
	if err := json.Unmarshal([]byte(jsonData), &lsp); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Should be treated as Object, not Disabled, because command field exists
	obj, ok := lsp.AsObject()
	if !ok {
		t.Fatal("AsObject() returned false, expected true")
	}

	if obj == nil {
		t.Fatal("AsObject() returned nil")
	}

	if !obj.Disabled {
		t.Errorf("Disabled = %v, want true", obj.Disabled)
	}

	// Should not be treated as Disabled
	disabled, ok := lsp.AsDisabled()
	if ok {
		t.Errorf("AsDisabled() returned true, expected false when command field exists")
	}

	if disabled != nil {
		t.Errorf("AsDisabled() returned non-nil, expected nil when command field exists")
	}
}
