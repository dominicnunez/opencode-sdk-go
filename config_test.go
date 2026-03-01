package opencode_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/dominicnunez/opencode-sdk-go"
	"github.com/dominicnunez/opencode-sdk-go/internal/testutil"
)

func TestConfigGetWithOptionalParams(t *testing.T) {
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
	_, err = client.Config.Get(context.TODO(), &opencode.ConfigGetParams{
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

func TestConfigUnmarshal(t *testing.T) {
	jsonData := `{
		"$schema": "https://opencode.ai/config-schema.json",
		"model": "anthropic/claude-3-5-sonnet",
		"autoupdate": true,
		"autoshare": false,
		"snapshot": true
	}`

	var config opencode.Config
	err := json.Unmarshal([]byte(jsonData), &config)
	if err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}

	if config.Schema != "https://opencode.ai/config-schema.json" {
		t.Errorf("expected schema to be 'https://opencode.ai/config-schema.json', got %s", config.Schema)
	}

	if config.Model != "anthropic/claude-3-5-sonnet" {
		t.Errorf("expected model to be 'anthropic/claude-3-5-sonnet', got %s", config.Model)
	}

	if !config.Autoupdate {
		t.Error("expected autoupdate to be true")
	}

	if config.Autoshare {
		t.Error("expected autoshare to be false")
	}

	if !config.Snapshot {
		t.Error("expected snapshot to be true")
	}
}

func TestConfigAgentUnmarshal(t *testing.T) {
	jsonData := `{
		"description": "General purpose agent",
		"model": "anthropic/claude-3-opus",
		"temperature": 0.8,
		"tools": {
			"bash": true,
			"edit": false
		}
	}`

	var agent opencode.ConfigAgentGeneral
	err := json.Unmarshal([]byte(jsonData), &agent)
	if err != nil {
		t.Fatalf("failed to unmarshal agent config: %v", err)
	}

	if agent.Description != "General purpose agent" {
		t.Errorf("expected description to be 'General purpose agent', got %s", agent.Description)
	}

	if agent.Model != "anthropic/claude-3-opus" {
		t.Errorf("expected model to be 'anthropic/claude-3-opus', got %s", agent.Model)
	}

	if agent.Temperature != 0.8 {
		t.Errorf("expected temperature to be 0.8, got %f", agent.Temperature)
	}

	if len(agent.Tools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(agent.Tools))
	}

	if !agent.Tools["bash"] {
		t.Error("expected bash tool to be enabled")
	}

	if agent.Tools["edit"] {
		t.Error("expected edit tool to be disabled")
	}
}

func TestConfigCommandUnmarshal(t *testing.T) {
	jsonData := `{
		"template": "Run the tests",
		"agent": "test-runner",
		"description": "Runs the test suite"
	}`

	var cmd opencode.ConfigCommand
	err := json.Unmarshal([]byte(jsonData), &cmd)
	if err != nil {
		t.Fatalf("failed to unmarshal command config: %v", err)
	}

	if cmd.Template != "Run the tests" {
		t.Errorf("expected template to be 'Run the tests', got %s", cmd.Template)
	}

	if cmd.Agent != "test-runner" {
		t.Errorf("expected agent to be 'test-runner', got %s", cmd.Agent)
	}

	if cmd.Description != "Runs the test suite" {
		t.Errorf("expected description to be 'Runs the test suite', got %s", cmd.Description)
	}
}

func TestConfigProviderUnmarshal(t *testing.T) {
	jsonData := `{
		"id": "custom-provider",
		"name": "Custom Provider",
		"api": "openai",
		"options": {
			"apiKey": "test-key",
			"baseURL": "https://api.custom.com"
		}
	}`

	var provider opencode.ConfigProvider
	err := json.Unmarshal([]byte(jsonData), &provider)
	if err != nil {
		t.Fatalf("failed to unmarshal provider config: %v", err)
	}

	if provider.ID != "custom-provider" {
		t.Errorf("expected id to be 'custom-provider', got %s", provider.ID)
	}

	if provider.Name != "Custom Provider" {
		t.Errorf("expected name to be 'Custom Provider', got %s", provider.Name)
	}

	if provider.API != "openai" {
		t.Errorf("expected api to be 'openai', got %s", provider.API)
	}

	if provider.Options.APIKey != "test-key" {
		t.Errorf("expected apiKey to be 'test-key', got %s", provider.Options.APIKey)
	}

	if provider.Options.BaseURL != "https://api.custom.com" {
		t.Errorf("expected baseURL to be 'https://api.custom.com', got %s", provider.Options.BaseURL)
	}
}
