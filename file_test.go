package opencode_test

import (
	"context"
	"strings"
	"testing"

	"github.com/dominicnunez/opencode-sdk-go"
)

// TestFileService_RequiredFieldValidation tests that required parameters are validated
func TestFileService_RequiredFieldValidation(t *testing.T) {
	tests := []struct {
		name           string
		method         func(context.Context, *opencode.Client) error
		expectedErrMsg string
	}{
		{
			name: "List with missing Path",
			method: func(ctx context.Context, client *opencode.Client) error {
				_, err := client.File.List(ctx, &opencode.FileListParams{
					Path: "",
				})
				return err
			},
			expectedErrMsg: "required path parameter",
		},
		{
			name: "Read with missing Path",
			method: func(ctx context.Context, client *opencode.Client) error {
				_, err := client.File.Read(ctx, &opencode.FileReadParams{
					Path: "",
				})
				return err
			},
			expectedErrMsg: "required path parameter",
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
