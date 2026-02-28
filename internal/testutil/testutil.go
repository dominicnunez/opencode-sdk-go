package testutil

import (
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

const testServerTimeout = 5 * time.Second

func CheckTestServer(t *testing.T, url string) bool {
	client := &http.Client{Timeout: testServerTimeout}
	if _, err := client.Get(url); err != nil {
		const skipEnvVar = "SKIP_MOCK_TESTS"
		if str, ok := os.LookupEnv(skipEnvVar); ok {
			skip, err := strconv.ParseBool(str)
			if err != nil {
				t.Fatalf("invalid %s value %q: %s", skipEnvVar, str, err)
			}
			if !skip {
				t.Fatalf("mock server not running and %s=false; start the server or set %s=true to skip", skipEnvVar, skipEnvVar)
			}
		}
		t.Skip("mock server not running; set SKIP_MOCK_TESTS=true to skip")
		return false
	}
	return true
}
