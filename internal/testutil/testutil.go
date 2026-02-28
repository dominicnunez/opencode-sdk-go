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
	return true
}
