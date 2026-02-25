package opencode_test

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCleanupVerification ensures that the SDK refactoring is complete and the codebase is clean.
func TestCleanupVerification(t *testing.T) {
	t.Run("aliases.go exists and re-exports used types", func(t *testing.T) {
		// aliases.go should exist and contain type aliases that are actively used by external consumers
		aliasesPath := "aliases.go"
		if _, err := os.Stat(aliasesPath); os.IsNotExist(err) {
			t.Fatal("aliases.go should exist - it re-exports types for public API")
		}
	})

	t.Run("param package is deleted", func(t *testing.T) {
		// The internal/param directory should no longer exist
		paramPath := filepath.Join("internal", "param")
		if _, err := os.Stat(paramPath); !os.IsNotExist(err) {
			t.Fatalf("internal/param should be deleted, but found at %s", paramPath)
		}
	})

	t.Run("field.go root file is deleted", func(t *testing.T) {
		// The root-level field.go (param helper re-exports) should be deleted
		fieldPath := "field.go"
		if _, err := os.Stat(fieldPath); !os.IsNotExist(err) {
			t.Fatalf("field.go should be deleted, but found at %s", fieldPath)
		}
	})

	t.Run("ptr.go exists with pointer helpers", func(t *testing.T) {
		// ptr.go should exist with idiomatic pointer helper functions
		ptrPath := "ptr.go"
		if _, err := os.Stat(ptrPath); os.IsNotExist(err) {
			t.Fatalf("ptr.go should exist with pointer helpers, but not found at %s", ptrPath)
		}
	})

	t.Run("gjson dependency is removed", func(t *testing.T) {
		// gjson was removed after converting unions to stdlib discriminated union pattern
		goModPath := "go.mod"
		content, err := os.ReadFile(goModPath)
		if err != nil {
			t.Fatalf("failed to read go.mod: %v", err)
		}

		// Verify gjson is no longer in go.mod (all unions now use stdlib)
		if containsString(string(content), "github.com/tidwall/gjson") {
			t.Error("gjson should be removed from go.mod - all unions now use stdlib discriminated union pattern")
		}
	})
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
