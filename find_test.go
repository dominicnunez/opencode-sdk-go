package opencode_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/dominicnunez/opencode-sdk-go"
	"github.com/dominicnunez/opencode-sdk-go/internal/testutil"
)

func TestFindFilesWithOptionalParams(t *testing.T) {
	t.Skip("Prism tests are disabled")
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
	_, err = client.Find.Files(context.TODO(), &opencode.FindFilesParams{
		Query:     "query",
		Directory: opencode.Ptr("directory"),
	})
	if err != nil {
		var apierr *opencode.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestFindSymbolsWithOptionalParams(t *testing.T) {
	t.Skip("Prism tests are disabled")
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
	_, err = client.Find.Symbols(context.TODO(), &opencode.FindSymbolsParams{
		Query:     "query",
		Directory: opencode.Ptr("directory"),
	})
	if err != nil {
		var apierr *opencode.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestFindTextWithOptionalParams(t *testing.T) {
	t.Skip("Prism tests are disabled")
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
	_, err = client.Find.Text(context.TODO(), &opencode.FindTextParams{
		Pattern:   "pattern",
		Directory: opencode.Ptr("directory"),
	})
	if err != nil {
		var apierr *opencode.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

// TestFindService_RequiredFieldValidation tests that required parameters are validated
func TestFindService_RequiredFieldValidation(t *testing.T) {
	tests := []struct {
		name           string
		method         func(context.Context, *opencode.Client) error
		expectedErrMsg string
	}{
		{
			name: "Files with missing Query",
			method: func(ctx context.Context, client *opencode.Client) error {
				_, err := client.Find.Files(ctx, &opencode.FindFilesParams{
					Query: "",
				})
				return err
			},
			expectedErrMsg: "required query parameter",
		},
		{
			name: "Symbols with missing Query",
			method: func(ctx context.Context, client *opencode.Client) error {
				_, err := client.Find.Symbols(ctx, &opencode.FindSymbolsParams{
					Query: "",
				})
				return err
			},
			expectedErrMsg: "required query parameter",
		},
		{
			name: "Text with missing Pattern",
			method: func(ctx context.Context, client *opencode.Client) error {
				_, err := client.Find.Text(ctx, &opencode.FindTextParams{
					Pattern: "",
				})
				return err
			},
			expectedErrMsg: "required query parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := opencode.NewClient(opencode.WithBaseURL("http://localhost:1"))
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			err = tt.method(context.Background(), client)
			if err == nil {
				t.Fatalf("Expected error for %s, got nil", tt.name)
			}
			if !strings.Contains(err.Error(), tt.expectedErrMsg) {
				t.Errorf("Expected error containing %q, got %q", tt.expectedErrMsg, err.Error())
			}
		})
	}
}
