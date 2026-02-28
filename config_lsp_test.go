package opencode

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestConfigLsp_AsDisabled_ValidDisabled(t *testing.T) {
	jsonData := `{"disabled": true}`

	var lsp ConfigLsp
	if err := json.Unmarshal([]byte(jsonData), &lsp); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	disabled, err := lsp.AsDisabled()
	if err != nil {
		t.Fatalf("AsDisabled() returned error: %v", err)
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

	obj, err := lsp.AsObject()
	if err != nil {
		t.Fatalf("AsObject() returned error: %v", err)
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

	obj, err := lsp.AsObject()
	if err != nil {
		t.Fatalf("AsObject() returned error: %v", err)
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

	obj, err := lsp.AsObject()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if obj != nil {
		t.Error("AsObject() should return nil for disabled config")
	}
}

func TestConfigLsp_WrongType_ObjectAsDisabled(t *testing.T) {
	jsonData := `{"command": ["gopls"], "disabled": false}`

	var lsp ConfigLsp
	if err := json.Unmarshal([]byte(jsonData), &lsp); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	disabled, err := lsp.AsDisabled()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if disabled != nil {
		t.Error("AsDisabled() should return nil for object config")
	}
}

func TestConfigLsp_InvalidJSON(t *testing.T) {
	jsonData := `{"disabled": "not a boolean"}`

	var lsp ConfigLsp
	if err := json.Unmarshal([]byte(jsonData), &lsp); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// AsDisabled should return error when trying to unmarshal the invalid JSON
	disabled, err := lsp.AsDisabled()
	if err == nil {
		t.Error("AsDisabled() should return error for invalid JSON")
	}
	if disabled != nil {
		t.Error("AsDisabled() should return nil for invalid JSON")
	}
}

func TestConfigLsp_EmptyJSON(t *testing.T) {
	jsonData := `{}`

	var lsp ConfigLsp
	if err := json.Unmarshal([]byte(jsonData), &lsp); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	obj, err := lsp.AsObject()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant for AsObject on empty JSON, got: %v", err)
	}
	if obj != nil {
		t.Error("AsObject() should return nil for empty JSON")
	}

	disabled, err := lsp.AsDisabled()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant for AsDisabled on empty JSON, got: %v", err)
	}
	if disabled != nil {
		t.Error("AsDisabled() should return nil for empty JSON")
	}
}

func TestConfigLsp_MalformedJSON(t *testing.T) {
	jsonData := `{"command": ["gopls", "disabled": true}`

	var lsp ConfigLsp
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
	jsonData := `{"command": ["gopls"], "disabled": true}`

	var lsp ConfigLsp
	if err := json.Unmarshal([]byte(jsonData), &lsp); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	obj, err := lsp.AsObject()
	if err != nil {
		t.Fatalf("AsObject() returned error: %v", err)
	}
	if obj == nil {
		t.Fatal("AsObject() returned nil")
	}
	if !obj.Disabled {
		t.Errorf("Disabled = %v, want true", obj.Disabled)
	}

	// Should not be treated as Disabled since command field exists
	disabled, err := lsp.AsDisabled()
	if !errors.Is(err, ErrWrongVariant) {
		t.Fatalf("expected ErrWrongVariant, got: %v", err)
	}
	if disabled != nil {
		t.Error("AsDisabled() should return nil when command field exists")
	}
}
