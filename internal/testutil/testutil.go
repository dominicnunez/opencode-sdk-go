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
	resp, err := client.Get(url) //nolint:noctx // health check only, no context needed
	if err != nil {
		const envVar = "REQUIRE_MOCK_SERVER"
		if str, ok := os.LookupEnv(envVar); ok {
			require, parseErr := strconv.ParseBool(str)
			if parseErr != nil {
				t.Fatalf("invalid %s value %q: %s", envVar, str, parseErr)
			}
			if require {
				t.Fatalf("mock server not running and %s=true", envVar)
			}
		}
		t.Skip("mock server not running; set REQUIRE_MOCK_SERVER=true to fail instead of skip")
		return false
	}
	_ = resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("mock server returned non-2xx status: %d", resp.StatusCode)
		return false
	}
	return true
}
