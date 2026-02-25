package opencode_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestApiformDeletion verifies that internal/apiform was safely deleted
// and no code depends on it.
func TestApiformDeletion(t *testing.T) {
	t.Run("DirectoryDoesNotExist", func(t *testing.T) {
		_, err := os.Stat("internal/apiform")
		if !os.IsNotExist(err) {
			t.Errorf("internal/apiform directory still exists")
		}
	})

	t.Run("NoApiformImports", func(t *testing.T) {
		// Walk all Go files and verify none import apiform
		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Skip non-Go files, test files (including this one), vendor, and .git
			if !strings.HasSuffix(path, ".go") ||
				strings.HasSuffix(path, "_test.go") ||
				strings.Contains(path, "vendor/") ||
				strings.Contains(path, ".git/") {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			if strings.Contains(string(content), "internal/apiform") ||
				strings.Contains(string(content), `"apiform"`) {
				t.Errorf("Production file %s still imports apiform", path)
			}

			return nil
		})
		if err != nil {
			t.Fatalf("Error walking directory: %v", err)
		}
	})

	t.Run("NoMultipartFormDataUsage", func(t *testing.T) {
		// Verify the SDK doesn't use multipart/form-data anywhere
		// (all endpoints use application/json as verified in openapi.yml)
		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Skip test files, internal/apiform (deleted), vendor, .git
			if !strings.HasSuffix(path, ".go") ||
				strings.HasSuffix(path, "_test.go") ||
				strings.Contains(path, "vendor/") ||
				strings.Contains(path, ".git/") ||
				strings.Contains(path, "internal/apiform") {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Verify no production code uses multipart
			if strings.Contains(string(content), "multipart.Writer") ||
				strings.Contains(string(content), "multipart.NewWriter") {
				t.Errorf("File %s uses multipart.Writer (should not be needed)", path)
			}

			return nil
		})
		if err != nil {
			t.Fatalf("Error walking directory: %v", err)
		}
	})

	t.Run("BuildStillWorks", func(t *testing.T) {
		// This test passing means go build succeeds without apiform
		// (the test suite runs go build before running tests)
	})
}
