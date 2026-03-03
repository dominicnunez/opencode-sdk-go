package opencode_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestScriptsLint_MissingGolangciLintShowsInstallMessage(t *testing.T) {
	repoRoot, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	fakeBinDir := t.TempDir()
	fakeGoPath := filepath.Join(fakeBinDir, "go")
	fakeGo := "#!/usr/bin/env bash\nset -euo pipefail\nexit 0\n"
	if err := os.WriteFile(fakeGoPath, []byte(fakeGo), 0o600); err != nil {
		t.Fatalf("write fake go: %v", err)
	}
	// #nosec G302 -- test helper must be executable to stub go via PATH.
	if err := os.Chmod(fakeGoPath, 0o755); err != nil {
		t.Fatalf("chmod fake go: %v", err)
	}

	fakeDirnamePath := filepath.Join(fakeBinDir, "dirname")
	fakeDirname := "#!/usr/bin/env bash\nset -euo pipefail\nprintf '%s\\n' \"${1%/*}\"\n"
	if err := os.WriteFile(fakeDirnamePath, []byte(fakeDirname), 0o600); err != nil {
		t.Fatalf("write fake dirname: %v", err)
	}
	// #nosec G302 -- test helper must be executable to provide dirname in PATH.
	if err := os.Chmod(fakeDirnamePath, 0o755); err != nil {
		t.Fatalf("chmod fake dirname: %v", err)
	}

	// #nosec G204 -- test runs a repository-local script from a controlled path.
	cmd := exec.Command("bash", filepath.Join(repoRoot, "scripts/lint"))
	cmd.Env = append(os.Environ(), "PATH="+fakeBinDir+":/usr/bin:/bin")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected scripts/lint to fail when golangci-lint is unavailable, output:\n%s", string(out))
	}

	output := string(out)
	if !strings.Contains(output, "golangci-lint is required") {
		t.Fatalf("expected install guidance for missing golangci-lint, got output:\n%s", output)
	}
}
