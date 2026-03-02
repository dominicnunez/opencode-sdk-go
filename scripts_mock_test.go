package opencode_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestScriptsMock_NoArgsUsesDefaultSpecPath(t *testing.T) {
	repoRoot, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	fakeBinDir := t.TempDir()
	recordedArgsPath := filepath.Join(t.TempDir(), "npm_args.txt")
	fakeNPMPath := filepath.Join(fakeBinDir, "npm")
	fakeNPM := "#!/usr/bin/env bash\nset -euo pipefail\nprintf '%s' \"$*\" > \"$TEST_NPM_ARGS_FILE\"\n"
	if err := os.WriteFile(fakeNPMPath, []byte(fakeNPM), 0o755); err != nil {
		t.Fatalf("write fake npm: %v", err)
	}

	cmd := exec.Command("bash", filepath.Join(repoRoot, "scripts/mock"))
	cmd.Env = append(os.Environ(),
		"PATH="+fakeBinDir+string(os.PathListSeparator)+os.Getenv("PATH"),
		"TEST_NPM_ARGS_FILE="+recordedArgsPath,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("scripts/mock without args failed: %v\noutput:\n%s", err, string(out))
	}

	output := string(out)
	if !strings.Contains(output, "Starting mock server with URL specs/openapi.yml") {
		t.Fatalf("unexpected script output: %q", output)
	}

	recordedArgs, err := os.ReadFile(recordedArgsPath)
	if err != nil {
		t.Fatalf("read recorded npm args: %v", err)
	}
	if !strings.Contains(string(recordedArgs), "prism mock specs/openapi.yml") {
		t.Fatalf("expected npm args to include default spec path, got %q", string(recordedArgs))
	}
}
