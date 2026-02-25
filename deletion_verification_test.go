package opencode

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDeletionVerification verifies that internal/requestconfig and option packages were deleted
func TestDeletionVerification(t *testing.T) {
	t.Run("requestconfig_directory_deleted", func(t *testing.T) {
		path := filepath.Join("internal", "requestconfig")
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("internal/requestconfig directory still exists at %s", path)
		}
	})

	t.Run("option_directory_deleted", func(t *testing.T) {
		path := "option"
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("option directory still exists at %s", path)
		}
	})

	t.Run("no_imports_of_deleted_packages", func(t *testing.T) {
		// Walk through all .go files to ensure no imports of deleted packages remain
		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip vendor, .git, and test directories
			if info.IsDir() && (info.Name() == "vendor" || info.Name() == ".git" || info.Name() == "node_modules") {
				return filepath.SkipDir
			}

			// Only check .go files, excluding this test file itself
			if !info.IsDir() && strings.HasSuffix(path, ".go") && path != "deletion_verification_test.go" {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				// Check for imports of deleted packages
				if strings.Contains(string(content), "github.com/dominicnunez/opencode-sdk-go/internal/requestconfig") {
					t.Errorf("File %s still imports internal/requestconfig", path)
				}
				if strings.Contains(string(content), "github.com/dominicnunez/opencode-sdk-go/option") {
					t.Errorf("File %s still imports option package", path)
				}
			}

			return nil
		})

		if err != nil {
			t.Fatalf("Error walking directory: %v", err)
		}
	})

	t.Run("aliases_no_longer_export_option_types", func(t *testing.T) {
		// Read aliases.go to ensure it doesn't reference deleted types
		content, err := os.ReadFile("aliases.go")
		if os.IsNotExist(err) {
			// aliases.go may have been deleted - that's fine
			return
		}
		if err != nil {
			t.Fatalf("Error reading aliases.go: %v", err)
		}

		// Check that RequestOption alias is removed
		if strings.Contains(string(content), "RequestOption") {
			t.Error("aliases.go still references RequestOption type")
		}
	})
}
