//go:build integration

package opencode_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/dominicnunez/opencode-sdk-go"
	"github.com/dominicnunez/opencode-sdk-go/internal/testutil"
)

func TestPathGetWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client, err := opencode.NewClient(opencode.WithBaseURL(baseURL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Path.Get(context.TODO(), &opencode.PathGetParams{
		Directory: opencode.Ptr("directory"),
	})
	if err != nil {
		var apierr *opencode.APIError
		if errors.As(err, &apierr) {
			t.Log(apierr.Error())
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
