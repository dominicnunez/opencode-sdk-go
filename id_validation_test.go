package opencode

import (
	"context"
	"errors"
	"testing"
)

func TestPathIdentifiers_RejectWhitespaceOnlyValues(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tests := []struct {
		name      string
		call      func() error
		wantParam string
	}{
		{
			name: "Auth.Set whitespace id",
			call: func() error {
				_, err := client.Auth.Set(context.Background(), " \t ", &AuthSetParams{
					Auth: ApiAuth{Key: "key"},
				})
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Get whitespace id",
			call: func() error {
				_, err := client.Session.Get(context.Background(), " \t ", nil)
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Update whitespace id",
			call: func() error {
				_, err := client.Session.Update(context.Background(), " \t ", nil)
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Delete whitespace id",
			call: func() error {
				_, err := client.Session.Delete(context.Background(), " \t ", nil)
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Abort whitespace id",
			call: func() error {
				_, err := client.Session.Abort(context.Background(), " \t ", nil)
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Children whitespace id",
			call: func() error {
				_, err := client.Session.Children(context.Background(), " \t ", nil)
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Command whitespace id",
			call: func() error {
				_, err := client.Session.Command(context.Background(), " \t ", &SessionCommandParams{
					Command: "echo hi",
				})
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Init whitespace id",
			call: func() error {
				_, err := client.Session.Init(context.Background(), " \t ", &SessionInitParams{
					MessageID:  "msg_123",
					ModelID:    "model_1",
					ProviderID: "provider_1",
				})
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Message whitespace id",
			call: func() error {
				_, err := client.Session.Message(context.Background(), " \t ", "msg_123", nil)
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Message whitespace message id",
			call: func() error {
				_, err := client.Session.Message(context.Background(), "sess_123", " \t ", nil)
				return err
			},
			wantParam: "messageID",
		},
		{
			name: "Session.Messages whitespace id",
			call: func() error {
				_, err := client.Session.Messages(context.Background(), " \t ", nil)
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Prompt whitespace id",
			call: func() error {
				_, err := client.Session.Prompt(context.Background(), " \t ", &SessionPromptParams{})
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Revert whitespace id",
			call: func() error {
				_, err := client.Session.Revert(context.Background(), " \t ", &SessionRevertParams{
					MessageID: "msg_123",
				})
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Share whitespace id",
			call: func() error {
				_, err := client.Session.Share(context.Background(), " \t ", nil)
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Diff whitespace id",
			call: func() error {
				_, err := client.Session.Diff(context.Background(), " \t ", nil)
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Fork whitespace id",
			call: func() error {
				_, err := client.Session.Fork(context.Background(), " \t ", nil)
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Shell whitespace id",
			call: func() error {
				_, err := client.Session.Shell(context.Background(), " \t ", &SessionShellParams{
					Agent:   "coder",
					Command: "pwd",
				})
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Summarize whitespace id",
			call: func() error {
				_, err := client.Session.Summarize(context.Background(), " \t ", &SessionSummarizeParams{
					ModelID:    "model_1",
					ProviderID: "provider_1",
				})
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Todo whitespace id",
			call: func() error {
				_, err := client.Session.Todo(context.Background(), " \t ", nil)
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Unrevert whitespace id",
			call: func() error {
				_, err := client.Session.Unrevert(context.Background(), " \t ", nil)
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Unshare whitespace id",
			call: func() error {
				_, err := client.Session.Unshare(context.Background(), " \t ", nil)
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Permissions.Respond whitespace id",
			call: func() error {
				_, err := client.Session.Permissions.Respond(context.Background(), " \t ", "perm_123", &SessionPermissionRespondParams{
					Response: PermissionResponseOnce,
				})
				return err
			},
			wantParam: "id",
		},
		{
			name: "Session.Permissions.Respond whitespace permission id",
			call: func() error {
				_, err := client.Session.Permissions.Respond(context.Background(), "sess_123", " \t ", &SessionPermissionRespondParams{
					Response: PermissionResponseOnce,
				})
				return err
			},
			wantParam: "permissionID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !errors.Is(err, &MissingRequiredParameterError{Parameter: tt.wantParam}) {
				t.Fatalf("expected missing required %s parameter, got %v", tt.wantParam, err)
			}
		})
	}
}
