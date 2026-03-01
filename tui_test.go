package opencode_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/dominicnunez/opencode-sdk-go"
	"github.com/dominicnunez/opencode-sdk-go/internal/testutil"
)

func TestTuiAppendPromptWithOptionalParams(t *testing.T) {
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
	_, err = client.Tui.AppendPrompt(context.TODO(), &opencode.TuiAppendPromptParams{
		Text:      "text",
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

func TestTuiClearPromptWithOptionalParams(t *testing.T) {
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
	_, err = client.Tui.ClearPrompt(context.TODO(), &opencode.TuiClearPromptParams{
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

func TestTuiExecuteCommandWithOptionalParams(t *testing.T) {
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
	_, err = client.Tui.ExecuteCommand(context.TODO(), &opencode.TuiExecuteCommandParams{
		Command:   "command",
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

func TestTuiOpenHelpWithOptionalParams(t *testing.T) {
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
	_, err = client.Tui.OpenHelp(context.TODO(), &opencode.TuiOpenHelpParams{
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

func TestTuiOpenModelsWithOptionalParams(t *testing.T) {
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
	_, err = client.Tui.OpenModels(context.TODO(), &opencode.TuiOpenModelsParams{
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

func TestTuiOpenSessionsWithOptionalParams(t *testing.T) {
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
	_, err = client.Tui.OpenSessions(context.TODO(), &opencode.TuiOpenSessionsParams{
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

func TestTuiOpenThemesWithOptionalParams(t *testing.T) {
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
	_, err = client.Tui.OpenThemes(context.TODO(), &opencode.TuiOpenThemesParams{
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

func TestTuiShowToastWithOptionalParams(t *testing.T) {
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
	_, err = client.Tui.ShowToast(context.TODO(), &opencode.TuiShowToastParams{
		Message:   "message",
		Variant:   opencode.ToastVariantInfo,
		Directory: opencode.Ptr("directory"),
		Title:     "title",
	})
	if err != nil {
		var apierr *opencode.APIError
		if errors.As(err, &apierr) {
			t.Log(apierr.Error())
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestTuiSubmitPromptWithOptionalParams(t *testing.T) {
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
	_, err = client.Tui.SubmitPrompt(context.TODO(), &opencode.TuiSubmitPromptParams{
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
