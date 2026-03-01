package opencode_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	opencode "github.com/dominicnunez/opencode-sdk-go"
)

func TestConfigService_Providers_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/config/providers" {
			t.Errorf("Expected path /config/providers, got %s", r.URL.Path)
		}

		response := opencode.ConfigProviderListResponse{
			Default: map[string]string{
				"anthropic": "claude-3-5-sonnet-20241022",
			},
			Providers: []opencode.ConfigProvider{
				{
					ID:   "anthropic",
					Name: "Anthropic",
					Env:  []string{"ANTHROPIC_API_KEY"},
					Models: map[string]opencode.ConfigProviderModel{
						"claude-3-5-sonnet-20241022": {
							ID:          "claude-3-5-sonnet-20241022",
							Name:        "Claude 3.5 Sonnet",
							Attachment:  true,
							Reasoning:   false,
							Temperature: true,
							ToolCall:    true,
							Cost: opencode.ConfigProviderModelsCost{
								Input:  3.0,
								Output: 15.0,
							},
							Limit: opencode.ConfigProviderModelsLimit{
								Context: 200000,
								Output:  8192,
							},
						},
					},
				},
			},
		}

		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	result, err := client.Config.Providers(ctx, &opencode.ConfigProviderListParams{})
	if err != nil {
		t.Fatalf("Config.Providers failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if len(result.Default) == 0 {
		t.Error("Expected non-empty default map")
	}

	if result.Default["anthropic"] != "claude-3-5-sonnet-20241022" {
		t.Errorf("Expected default model 'claude-3-5-sonnet-20241022', got '%s'", result.Default["anthropic"])
	}

	if len(result.Providers) != 1 {
		t.Fatalf("Expected 1 provider, got %d", len(result.Providers))
	}

	provider := result.Providers[0]
	if provider.ID != "anthropic" {
		t.Errorf("Expected provider ID 'anthropic', got '%s'", provider.ID)
	}

	if provider.Name != "Anthropic" {
		t.Errorf("Expected provider name 'Anthropic', got '%s'", provider.Name)
	}

	if len(provider.Models) != 1 {
		t.Fatalf("Expected 1 model, got %d", len(provider.Models))
	}

	model := provider.Models["claude-3-5-sonnet-20241022"]
	if model.ID != "claude-3-5-sonnet-20241022" {
		t.Errorf("Expected model ID 'claude-3-5-sonnet-20241022', got '%s'", model.ID)
	}

	if model.Cost.Input != 3.0 {
		t.Errorf("Expected input cost 3.0, got %f", model.Cost.Input)
	}

	if model.Limit.Context != 200000 {
		t.Errorf("Expected context limit 200000, got %f", model.Limit.Context)
	}
}

func TestConfigService_Providers_WithDirectory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("directory") != "/test/dir" {
			t.Errorf("Expected directory query param '/test/dir', got '%s'", r.URL.Query().Get("directory"))
		}

		response := opencode.ConfigProviderListResponse{
			Default:   map[string]string{},
			Providers: []opencode.ConfigProvider{},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	dir := "/test/dir"
	_, err = client.Config.Providers(ctx, &opencode.ConfigProviderListParams{
		Directory: &dir,
	})
	if err != nil {
		t.Fatalf("Config.Providers failed: %v", err)
	}
}

func TestConfigService_Providers_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := opencode.ConfigProviderListResponse{
			Default:   map[string]string{},
			Providers: []opencode.ConfigProvider{},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	result, err := client.Config.Providers(ctx, &opencode.ConfigProviderListParams{})
	if err != nil {
		t.Fatalf("Config.Providers failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if len(result.Default) != 0 {
		t.Errorf("Expected empty default map, got %d entries", len(result.Default))
	}

	if len(result.Providers) != 0 {
		t.Errorf("Expected empty providers list, got %d entries", len(result.Providers))
	}
}

func TestConfigService_Providers_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal server error"))
	}))
	defer server.Close()

	client, err := opencode.NewClient(
		opencode.WithBaseURL(server.URL),
		opencode.WithMaxRetries(0),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	_, err = client.Config.Providers(ctx, &opencode.ConfigProviderListParams{})
	if err == nil {
		t.Fatal("expected error for server error, got nil")
	}
	var apiErr *opencode.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, apiErr.StatusCode)
	}
}

func TestConfigService_Providers_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client, err := opencode.NewClient(opencode.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	_, err = client.Config.Providers(ctx, &opencode.ConfigProviderListParams{})
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
}

func TestConfigProviderListParams_URLQuery(t *testing.T) {
	tests := []struct {
		name     string
		params   opencode.ConfigProviderListParams
		expected map[string]string
	}{
		{
			name: "with_directory",
			params: opencode.ConfigProviderListParams{
				Directory: ptrString("/test"),
			},
			expected: map[string]string{
				"directory": "/test",
			},
		},
		{
			name:     "without_directory",
			params:   opencode.ConfigProviderListParams{},
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, err := tt.params.URLQuery()
			if err != nil {
				t.Fatalf("URLQuery failed: %v", err)
			}

			for key, expectedValue := range tt.expected {
				if values.Get(key) != expectedValue {
					t.Errorf("Expected %s=%s, got %s", key, expectedValue, values.Get(key))
				}
			}

			if len(values) != len(tt.expected) {
				t.Errorf("Expected %d query params, got %d", len(tt.expected), len(values))
			}
		})
	}
}
