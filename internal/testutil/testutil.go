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
		const SKIP_MOCK_TESTS = "SKIP_MOCK_TESTS"
		if str, ok := os.LookupEnv(SKIP_MOCK_TESTS); ok {
			skip, err := strconv.ParseBool(str)
			if err != nil {
				t.Errorf("strconv.ParseBool(os.LookupEnv(%s)) failed: %s", SKIP_MOCK_TESTS, err)
				return false
			}
			if !skip {
				t.Errorf("The test will not run without a mock Prism server running against your OpenAPI spec. You can set the environment variable %s to true to skip running any tests that require the mock server", SKIP_MOCK_TESTS)
				return false
			}
		}
		t.Skip("The test will not run without a mock Prism server running against your OpenAPI spec")
		return false
	}
	return true
}
