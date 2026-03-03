package opencode_test

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckSpecUpdate_MissingGHShowsInstallError(t *testing.T) {
	scriptPath, _, fakeBin := setupCheckSpecUpdateWorkspace(t, "local-spec")
	writeExecutable(t, fakeBin, "dirname", "#!/bin/sh\nprintf '%s\\n' \"${1%/*}\"\n")

	cmd := exec.Command("bash", scriptPath)
	cmd.Env = append(os.Environ(), "PATH="+fakeBin)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected script to fail without gh, output:\n%s", string(out))
	}
	if !strings.Contains(string(out), "gh is required but not installed") {
		t.Fatalf("expected missing gh guidance, got output:\n%s", string(out))
	}
}

func TestCheckSpecUpdate_UpdateFailsWhenDownloadFails(t *testing.T) {
	const localSpec = "old-spec"
	const upstreamURL = "https://example.com/spec.yml"
	scriptPath, specPath, fakeBin := setupCheckSpecUpdateWorkspace(t, localSpec)

	stats := makeStatsYAML(sha256Hex("new-spec"), upstreamURL, "42")
	writeExecutable(t, fakeBin, "gh", "#!/usr/bin/env bash\nset -euo pipefail\nprintf '%s' \"$TEST_STATS_B64\"\n")
	writeExecutable(t, fakeBin, "curl", "#!/usr/bin/env bash\nset -euo pipefail\nexit 22\n")

	cmd := exec.Command("bash", scriptPath, "--update")
	cmd.Env = append(os.Environ(),
		"PATH="+fakeBin+string(os.PathListSeparator)+os.Getenv("PATH"),
		"TEST_STATS_B64="+base64.StdEncoding.EncodeToString([]byte(stats)),
	)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected update to fail when download fails, output:\n%s", string(out))
	}
	if !strings.Contains(string(out), "failed to download upstream spec") {
		t.Fatalf("expected download failure message, got output:\n%s", string(out))
	}

	unchanged, readErr := os.ReadFile(specPath)
	if readErr != nil {
		t.Fatalf("read spec after failed update: %v", readErr)
	}
	if string(unchanged) != localSpec {
		t.Fatalf("spec changed on failed download, got %q", string(unchanged))
	}
}

func TestCheckSpecUpdate_UpdateFailsOnHashMismatch(t *testing.T) {
	const localSpec = "old-spec"
	const upstreamURL = "https://example.com/spec.yml"
	const downloadedSpec = "downloaded-spec"
	const expectedSpec = "different-spec"
	scriptPath, specPath, fakeBin := setupCheckSpecUpdateWorkspace(t, localSpec)

	stats := makeStatsYAML(sha256Hex(expectedSpec), upstreamURL, "42")
	writeExecutable(t, fakeBin, "gh", "#!/usr/bin/env bash\nset -euo pipefail\nprintf '%s' \"$TEST_STATS_B64\"\n")
	writeExecutable(t, fakeBin, "curl", curlWriteBodyScript())

	cmd := exec.Command("bash", scriptPath, "--update")
	cmd.Env = append(os.Environ(),
		"PATH="+fakeBin+string(os.PathListSeparator)+os.Getenv("PATH"),
		"TEST_STATS_B64="+base64.StdEncoding.EncodeToString([]byte(stats)),
		"TEST_CURL_BODY="+downloadedSpec,
	)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected update to fail on hash mismatch, output:\n%s", string(out))
	}
	if !strings.Contains(string(out), "does not match upstream") {
		t.Fatalf("expected hash mismatch message, got output:\n%s", string(out))
	}

	unchanged, readErr := os.ReadFile(specPath)
	if readErr != nil {
		t.Fatalf("read spec after failed update: %v", readErr)
	}
	if string(unchanged) != localSpec {
		t.Fatalf("spec changed on hash mismatch, got %q", string(unchanged))
	}
}

func TestCheckSpecUpdate_UpdateReplacesSpecOnMatchingHash(t *testing.T) {
	const localSpec = "old-spec"
	const upstreamURL = "https://example.com/spec.yml"
	const downloadedSpec = "new-spec"
	scriptPath, specPath, fakeBin := setupCheckSpecUpdateWorkspace(t, localSpec)

	stats := makeStatsYAML(sha256Hex(downloadedSpec), upstreamURL, "42")
	writeExecutable(t, fakeBin, "gh", "#!/usr/bin/env bash\nset -euo pipefail\nprintf '%s' \"$TEST_STATS_B64\"\n")
	writeExecutable(t, fakeBin, "curl", curlWriteBodyScript())

	cmd := exec.Command("bash", scriptPath, "--update")
	cmd.Env = append(os.Environ(),
		"PATH="+fakeBin+string(os.PathListSeparator)+os.Getenv("PATH"),
		"TEST_STATS_B64="+base64.StdEncoding.EncodeToString([]byte(stats)),
		"TEST_CURL_BODY="+downloadedSpec,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected successful update, err=%v output:\n%s", err, string(out))
	}
	if !strings.Contains(string(out), "Spec updated") {
		t.Fatalf("expected success message, got output:\n%s", string(out))
	}

	updated, readErr := os.ReadFile(specPath)
	if readErr != nil {
		t.Fatalf("read updated spec: %v", readErr)
	}
	if string(updated) != downloadedSpec {
		t.Fatalf("spec was not replaced, got %q", string(updated))
	}
}

func TestCheckSpecUpdate_Base64DecodeFallsBackToBSDFlag(t *testing.T) {
	const localSpec = "stable-spec"
	const upstreamURL = "https://example.com/spec.yml"

	scriptPath, _, fakeBin := setupCheckSpecUpdateWorkspace(t, localSpec)

	stats := makeStatsYAML(sha256Hex(localSpec), upstreamURL, "42")
	writeExecutable(t, fakeBin, "gh", "#!/usr/bin/env bash\nset -euo pipefail\nprintf '%s' \"$TEST_STATS_B64\"\n")

	realBase64Path, err := exec.LookPath("base64")
	if err != nil {
		t.Fatalf("failed to locate real base64 binary: %v", err)
	}

	writeExecutable(t, fakeBin, "base64", "#!/usr/bin/env bash\nset -euo pipefail\nmode=\"${1:-}\"\ninput=\"$(cat)\"\ncase \"$mode\" in\n  --decode|-d)\n    exit 2\n    ;;\n  -D)\n    printf '%s' \"$input\" | \"$REAL_BASE64\" --decode\n    ;;\n  *)\n    exit 2\n    ;;\nesac\n")

	cmd := exec.Command("bash", scriptPath)
	cmd.Env = append(os.Environ(),
		"PATH="+fakeBin+string(os.PathListSeparator)+os.Getenv("PATH"),
		"TEST_STATS_B64="+base64.StdEncoding.EncodeToString([]byte(stats)),
		"REAL_BASE64="+realBase64Path,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected script to succeed using BSD base64 decode flag fallback, err=%v output:\n%s", err, string(out))
	}
	if !strings.Contains(string(out), "Spec is up to date") {
		t.Fatalf("expected up-to-date message, got output:\n%s", string(out))
	}
}

func setupCheckSpecUpdateWorkspace(t *testing.T, initialSpec string) (scriptPath, specPath, fakeBin string) {
	t.Helper()
	const testDirPerms = 0o750

	repoRoot, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	sourceScript := filepath.Join(repoRoot, "scripts/check-spec-update.sh")
	scriptBytes, err := os.ReadFile(sourceScript)
	if err != nil {
		t.Fatalf("read source script: %v", err)
	}

	projectDir := t.TempDir()
	scriptsDir := filepath.Join(projectDir, "scripts")
	specsDir := filepath.Join(projectDir, "specs")
	fakeBin = filepath.Join(projectDir, "fakebin")
	if err := os.MkdirAll(scriptsDir, testDirPerms); err != nil {
		t.Fatalf("mkdir scripts dir: %v", err)
	}
	if err := os.MkdirAll(specsDir, testDirPerms); err != nil {
		t.Fatalf("mkdir specs dir: %v", err)
	}
	if err := os.MkdirAll(fakeBin, testDirPerms); err != nil {
		t.Fatalf("mkdir fakebin: %v", err)
	}

	scriptPath = filepath.Join(scriptsDir, "check-spec-update.sh")
	specPath = filepath.Join(specsDir, "openapi.yml")
	if err := os.WriteFile(scriptPath, scriptBytes, 0o600); err != nil {
		t.Fatalf("write script copy: %v", err)
	}
	if err := os.WriteFile(specPath, []byte(initialSpec), 0o600); err != nil {
		t.Fatalf("write initial spec: %v", err)
	}

	return scriptPath, specPath, fakeBin
}

func writeExecutable(t *testing.T, dir, name, body string) {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write fake %s: %v", name, err)
	}
	// #nosec G302 -- test helper must be executable to stub external tools.
	if err := os.Chmod(path, 0o755); err != nil {
		t.Fatalf("chmod fake %s: %v", name, err)
	}
}

func sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func makeStatsYAML(hash, url, endpoints string) string {
	return "openapi_spec_hash: " + hash + "\n" +
		"openapi_spec_url: " + url + "\n" +
		"configured_endpoints: " + endpoints + "\n"
}

func curlWriteBodyScript() string {
	return "#!/usr/bin/env bash\n" +
		"set -euo pipefail\n" +
		"out=\"\"\n" +
		"while [ \"$#\" -gt 0 ]; do\n" +
		"  if [ \"$1\" = \"-o\" ]; then\n" +
		"    out=\"$2\"\n" +
		"    shift 2\n" +
		"    continue\n" +
		"  fi\n" +
		"  shift\n" +
		"done\n" +
		"if [ -z \"$out\" ]; then\n" +
		"  echo \"missing -o\" >&2\n" +
		"  exit 2\n" +
		"fi\n" +
		"printf '%s' \"$TEST_CURL_BODY\" > \"$out\"\n"
}
