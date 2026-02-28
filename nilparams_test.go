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
		{
			name: "SessionService.Command",
			call: func() error {
				_, err := client.Session.Command(context.Background(), "ses_123", nil)
				return err
			},
		},
		{
			name: "SessionService.Init",
			call: func() error {
				_, err := client.Session.Init(context.Background(), "ses_123", nil)
				return err
			},
		},
		{
			name: "SessionService.Prompt",
			call: func() error {
				_, err := client.Session.Prompt(context.Background(), "ses_123", nil)
				return err
			},
		},
		{
			name: "SessionService.Revert",
			call: func() error {
				_, err := client.Session.Revert(context.Background(), "ses_123", nil)
				return err
			},
		},
		{
			name: "SessionPermissionService.Respond",
			call: func() error {
				_, err := client.Session.Permissions.Respond(context.Background(), "ses_123", "perm_456", nil)
				return err
			},
		},
		{
			name: "TuiService.AppendPrompt",
			call: func() error {
				_, err := client.Tui.AppendPrompt(context.Background(), nil)
				return err
			},
		},
		{
			name: "TuiService.ExecuteCommand",
			call: func() error {
				_, err := client.Tui.ExecuteCommand(context.Background(), nil)
				return err
			},
		},
		{
			name: "TuiService.ShowToast",
			call: func() error {
				_, err := client.Tui.ShowToast(context.Background(), nil)
				return err
			},
		},
		{
			name: "AppService.Log",
			call: func() error {
				_, err := client.App.Log(context.Background(), nil)
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
