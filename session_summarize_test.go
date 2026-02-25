package opencode

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSessionSummarize_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/session/sess_123/summarize" {
			t.Errorf("Expected path /session/sess_123/summarize, got %s", r.URL.Path)
		}

		// Verify request body
		var body SessionSummarizeParams
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if body.ProviderID != "anthropic" {
			t.Errorf("Expected providerID 'anthropic', got %s", body.ProviderID)
		}
		if body.ModelID != "claude-sonnet-4-5" {
			t.Errorf("Expected modelID 'claude-sonnet-4-5', got %s", body.ModelID)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("true"))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	result, err := client.Session.Summarize(context.Background(), "sess_123", &SessionSummarizeParams{
		ProviderID: "anthropic",
		ModelID:    "claude-sonnet-4-5",
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}
}

func TestSessionSummarize_WithDirectory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify directory query param
		if r.URL.Query().Get("directory") != "/home/user/project" {
			t.Errorf("Expected directory '/home/user/project', got %s", r.URL.Query().Get("directory"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("false"))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	directory := "/home/user/project"
	result, err := client.Session.Summarize(context.Background(), "sess_123", &SessionSummarizeParams{
		ProviderID: "anthropic",
		ModelID:    "claude-sonnet-4-5",
		Directory:  &directory,
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

func TestSessionSummarize_MissingID(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Session.Summarize(context.Background(), "", &SessionSummarizeParams{
		ProviderID: "anthropic",
		ModelID:    "claude-sonnet-4-5",
	})
	if err == nil {
		t.Fatal("Expected error for missing id, got nil")
	}
	if err.Error() != "missing required id parameter" {
		t.Errorf("Expected 'missing required id parameter', got %v", err)
	}
}

func TestSessionSummarize_MissingParams(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Session.Summarize(context.Background(), "sess_123", nil)
	if err == nil {
		t.Fatal("Expected error for missing params, got nil")
	}
	if err.Error() != "missing required params" {
		t.Errorf("Expected 'missing required params', got %v", err)
	}
}

func TestSessionSummarize_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Session.Summarize(context.Background(), "sess_123", &SessionSummarizeParams{
		ProviderID: "anthropic",
		ModelID:    "claude-sonnet-4-5",
	})
	if err == nil {
		t.Fatal("Expected error for server error, got nil")
	}
}

func TestSessionSummarize_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Session.Summarize(context.Background(), "sess_123", &SessionSummarizeParams{
		ProviderID: "anthropic",
		ModelID:    "claude-sonnet-4-5",
	})
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
}

func TestSessionSummarizeParams_Marshal(t *testing.T) {
	directory := "/home/user/project"
	params := SessionSummarizeParams{
		ProviderID: "anthropic",
		ModelID:    "claude-sonnet-4-5",
		Directory:  &directory,
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal params: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal params: %v", err)
	}

	if result["providerID"] != "anthropic" {
		t.Errorf("Expected providerID 'anthropic', got %v", result["providerID"])
	}
	if result["modelID"] != "claude-sonnet-4-5" {
		t.Errorf("Expected modelID 'claude-sonnet-4-5', got %v", result["modelID"])
	}
	// directory should not be in JSON body (it's a query param)
	if _, ok := result["directory"]; ok {
		t.Errorf("Expected directory to be omitted from JSON body, but it was present")
	}
}

func TestSessionSummarizeParams_URLQuery(t *testing.T) {
	directory := "/home/user/project"
	params := SessionSummarizeParams{
		ProviderID: "anthropic",
		ModelID:    "claude-sonnet-4-5",
		Directory:  &directory,
	}

	values, err := params.URLQuery()
	if err != nil {
		t.Fatalf("Failed to get URL query: %v", err)
	}

	if values.Get("directory") != "/home/user/project" {
		t.Errorf("Expected directory '/home/user/project', got %s", values.Get("directory"))
	}
}

func TestSessionSummarizeParams_URLQuery_NoDirectory(t *testing.T) {
	params := SessionSummarizeParams{
		ProviderID: "anthropic",
		ModelID:    "claude-sonnet-4-5",
	}

	values, err := params.URLQuery()
	if err != nil {
		t.Fatalf("Failed to get URL query: %v", err)
	}

	if values.Get("directory") != "" {
		t.Errorf("Expected empty directory, got %s", values.Get("directory"))
	}
}

func TestSessionSummarize_BothResponses(t *testing.T) {
	tests := []struct {
		name     string
		response string
		expected bool
	}{
		{
			name:     "true response",
			response: "true",
			expected: true,
		},
		{
			name:     "false response",
			response: "false",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client, err := NewClient(WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			result, err := client.Session.Summarize(context.Background(), "sess_123", &SessionSummarizeParams{
				ProviderID: "anthropic",
				ModelID:    "claude-sonnet-4-5",
			})
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSessionSummarize_URLEncoding(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify URL encoding of directory param
		decodedDir, err := url.QueryUnescape(r.URL.Query().Get("directory"))
		if err != nil {
			t.Fatalf("Failed to decode directory: %v", err)
		}
		if decodedDir != "/home/user/my project" {
			t.Errorf("Expected directory '/home/user/my project', got %s", decodedDir)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("true"))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	directory := "/home/user/my project"
	_, err = client.Session.Summarize(context.Background(), "sess_123", &SessionSummarizeParams{
		ProviderID: "anthropic",
		ModelID:    "claude-sonnet-4-5",
		Directory:  &directory,
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
