package opencode

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestGolangciLint verifies that the codebase passes golangci-lint checks.
// This test ensures code quality standards are maintained.
func TestGolangciLint(t *testing.T) {
	// Try to find golangci-lint in GOPATH/bin or PATH
	var lintPath string

	// Get GOPATH using go env command
	cmd := exec.Command("go", "env", "GOPATH")
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		gopath := string(output)
		// Remove trailing newline
		if len(gopath) > 0 && gopath[len(gopath)-1] == '\n' {
			gopath = gopath[:len(gopath)-1]
		}
		gopathLint := filepath.Join(gopath, "bin", "golangci-lint")
		if _, err := os.Stat(gopathLint); err == nil {
			lintPath = gopathLint
		}
	}

	if lintPath == "" {
		// Try PATH
		path, err := exec.LookPath("golangci-lint")
		if err != nil {
			t.Skip("golangci-lint not found in PATH or GOPATH/bin, skipping")
			return
		}
		lintPath = path
	}

	cmd = exec.Command(lintPath, "run", "./...")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("golangci-lint failed:\n%s\nError: %v", string(output), err)
	}

	if len(output) > 0 {
		t.Logf("golangci-lint output:\n%s", string(output))
	}
}
