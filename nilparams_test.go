package opencode_test

import (
	"context"
	"strings"
	"testing"

	"github.com/dominicnunez/opencode-sdk-go"
)

func TestNilParams_ReturnsError(t *testing.T) {
	client, err := opencode.NewClient(opencode.WithBaseURL("http://localhost:0"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "FindService.Files",
			call: func() error {
				_, err := client.Find.Files(context.Background(), nil)
				return err
			},
		},
		{
			name: "FindService.Symbols",
			call: func() error {
				_, err := client.Find.Symbols(context.Background(), nil)
				return err
			},
		},
		{
			name: "FindService.Text",
			call: func() error {
				_, err := client.Find.Text(context.Background(), nil)
				return err
			},
		},
		{
			name: "FileService.List",
			call: func() error {
				_, err := client.File.List(context.Background(), nil)
				return err
			},
		},
		{
			name: "FileService.Read",
			call: func() error {
				_, err := client.File.Read(context.Background(), nil)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("expected error for nil params, got nil")
			}
			if !strings.Contains(err.Error(), "params is required") {
				t.Errorf("expected error containing %q, got %q", "params is required", err.Error())
			}
		})
	}
}
