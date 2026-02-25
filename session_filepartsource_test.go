package opencode

import (
	"encoding/json"
	"testing"
)

func TestFilePartSource_AsFile_ValidFileSource(t *testing.T) {
	jsonData := `{
		"type": "file",
		"path": "/home/user/project/main.go",
		"text": {
			"start": 0,
			"end": 100,
			"value": "package main\n\nfunc main() {\n\tprintln(\"Hello\")\n}"
		}
	}`

	var source FilePartSource
	if err := json.Unmarshal([]byte(jsonData), &source); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if source.Type != FilePartSourceTypeFile {
		t.Errorf("Expected type %q, got %q", FilePartSourceTypeFile, source.Type)
	}

	fileSource, ok := source.AsFile()
	if !ok {
		t.Fatal("AsFile() should succeed for file type")
	}
	if fileSource == nil {
		t.Fatal("AsFile() returned nil fileSource")
	}

	if fileSource.Path != "/home/user/project/main.go" {
		t.Errorf("Expected path %q, got %q", "/home/user/project/main.go", fileSource.Path)
	}
	if fileSource.Text.Start != 0 {
		t.Errorf("Expected text.start 0, got %d", fileSource.Text.Start)
	}
	if fileSource.Text.End != 100 {
		t.Errorf("Expected text.end 100, got %d", fileSource.Text.End)
	}
	expectedValue := "package main\n\nfunc main() {\n\tprintln(\"Hello\")\n}"
	if fileSource.Text.Value != expectedValue {
		t.Errorf("Expected text value %q, got %q", expectedValue, fileSource.Text.Value)
	}
}

func TestFilePartSource_AsSymbol_ValidSymbolSource(t *testing.T) {
	jsonData := `{
		"type": "symbol",
		"kind": 12,
		"name": "main",
		"path": "/home/user/project/main.go",
		"range": {
			"start": {
				"line": 3.0,
				"character": 0.0
			},
			"end": {
				"line": 5.0,
				"character": 1.0
			}
		},
		"text": {
			"start": 50,
			"end": 150,
			"value": "func main() {\n\tprintln(\"Hello\")\n}"
		}
	}`

	var source FilePartSource
	if err := json.Unmarshal([]byte(jsonData), &source); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if source.Type != FilePartSourceTypeSymbol {
		t.Errorf("Expected type %q, got %q", FilePartSourceTypeSymbol, source.Type)
	}

	symbolSource, ok := source.AsSymbol()
	if !ok {
		t.Fatal("AsSymbol() should succeed for symbol type")
	}
	if symbolSource == nil {
		t.Fatal("AsSymbol() returned nil symbolSource")
	}

	if symbolSource.Kind != 12 {
		t.Errorf("Expected kind 12, got %d", symbolSource.Kind)
	}
	if symbolSource.Name != "main" {
		t.Errorf("Expected name %q, got %q", "main", symbolSource.Name)
	}
	if symbolSource.Path != "/home/user/project/main.go" {
		t.Errorf("Expected path %q, got %q", "/home/user/project/main.go", symbolSource.Path)
	}
	if symbolSource.Range.Start.Line != 3.0 {
		t.Errorf("Expected range.start.line 3.0, got %f", symbolSource.Range.Start.Line)
	}
	if symbolSource.Range.End.Character != 1.0 {
		t.Errorf("Expected range.end.character 1.0, got %f", symbolSource.Range.End.Character)
	}
	if symbolSource.Text.Start != 50 {
		t.Errorf("Expected text.start 50, got %d", symbolSource.Text.Start)
	}
}

func TestFilePartSource_AsFile_WrongType(t *testing.T) {
	jsonData := `{
		"type": "symbol",
		"kind": 12,
		"name": "main",
		"path": "/home/user/project/main.go",
		"range": {
			"start": {"line": 3.0, "character": 0.0},
			"end": {"line": 5.0, "character": 1.0}
		},
		"text": {
			"start": 50,
			"end": 150,
			"value": "func main() {}"
		}
	}`

	var source FilePartSource
	if err := json.Unmarshal([]byte(jsonData), &source); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	fileSource, ok := source.AsFile()
	if ok {
		t.Error("AsFile() should fail for symbol type")
	}
	if fileSource != nil {
		t.Error("AsFile() should return nil for symbol type")
	}
}

func TestFilePartSource_AsSymbol_WrongType(t *testing.T) {
	jsonData := `{
		"type": "file",
		"path": "/home/user/project/main.go",
		"text": {
			"start": 0,
			"end": 100,
			"value": "package main"
		}
	}`

	var source FilePartSource
	if err := json.Unmarshal([]byte(jsonData), &source); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	symbolSource, ok := source.AsSymbol()
	if ok {
		t.Error("AsSymbol() should fail for file type")
	}
	if symbolSource != nil {
		t.Error("AsSymbol() should return nil for file type")
	}
}

func TestFilePartSource_InvalidJSON(t *testing.T) {
	jsonData := `{
		"type": "invalid",
		"path": "/home/user/project/main.go"
	}`

	var source FilePartSource
	if err := json.Unmarshal([]byte(jsonData), &source); err != nil {
		t.Fatalf("Unmarshal should succeed even with unknown type: %v", err)
	}

	// Type should be parsed even if unknown
	if source.Type != "invalid" {
		t.Errorf("Expected type %q, got %q", "invalid", source.Type)
	}

	// AsFile and AsSymbol should both fail for unknown type
	if fileSource, ok := source.AsFile(); ok || fileSource != nil {
		t.Error("AsFile() should fail for unknown type")
	}
	if symbolSource, ok := source.AsSymbol(); ok || symbolSource != nil {
		t.Error("AsSymbol() should fail for unknown type")
	}
}

func TestFilePartSource_MalformedJSON(t *testing.T) {
	jsonData := `{type": "file", "path": "/test"}`

	var source FilePartSource
	err := json.Unmarshal([]byte(jsonData), &source)
	if err == nil {
		t.Error("Expected unmarshal error for malformed JSON")
	}
}

func TestFilePartSource_EmptyJSON(t *testing.T) {
	jsonData := `{}`

	var source FilePartSource
	if err := json.Unmarshal([]byte(jsonData), &source); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Type should be empty string
	if source.Type != "" {
		t.Errorf("Expected empty type, got %q", source.Type)
	}

	// Both methods should fail
	if fileSource, ok := source.AsFile(); ok || fileSource != nil {
		t.Error("AsFile() should fail for empty type")
	}
	if symbolSource, ok := source.AsSymbol(); ok || symbolSource != nil {
		t.Error("AsSymbol() should fail for empty type")
	}
}

func TestFilePartSource_MissingType(t *testing.T) {
	jsonData := `{
		"path": "/home/user/project/main.go",
		"text": {
			"start": 0,
			"end": 100,
			"value": "package main"
		}
	}`

	var source FilePartSource
	if err := json.Unmarshal([]byte(jsonData), &source); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Type should be empty/zero value
	if source.Type != "" {
		t.Errorf("Expected empty type, got %q", source.Type)
	}

	// Both methods should fail
	if fileSource, ok := source.AsFile(); ok || fileSource != nil {
		t.Error("AsFile() should fail for missing type")
	}
	if symbolSource, ok := source.AsSymbol(); ok || symbolSource != nil {
		t.Error("AsSymbol() should fail for missing type")
	}
}

func TestFilePartSource_MalformedNestedJSON(t *testing.T) {
	// Valid outer structure but invalid inner data
	jsonData := `{
		"type": "file",
		"path": "/home/user/project/main.go",
		"text": "this should be an object, not a string"
	}`

	var source FilePartSource
	if err := json.Unmarshal([]byte(jsonData), &source); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Type should be parsed
	if source.Type != FilePartSourceTypeFile {
		t.Errorf("Expected type %q, got %q", FilePartSourceTypeFile, source.Type)
	}

	// AsFile should fail because text field is malformed
	fileSource, ok := source.AsFile()
	if ok {
		t.Error("AsFile() should fail for malformed nested JSON")
	}
	if fileSource != nil {
		t.Error("AsFile() should return nil for malformed nested JSON")
	}
}
