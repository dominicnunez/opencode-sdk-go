package opencode

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestAuthService_Set_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}
		if r.URL.Path != "/auth/provider-id" {
			t.Errorf("Expected path /auth/provider-id, got %s", r.URL.Path)
		}

		// Verify request body contains Auth union
		var body OAuth
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if body.Type != AuthTypeOAuth {
			t.Errorf("Expected auth type oauth, got %s", body.Type)
		}
		if body.Refresh != "refresh_token" {
			t.Errorf("Expected refresh token 'refresh_token', got %s", body.Refresh)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(true)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	result, err := client.Auth.Set(context.Background(), "provider-id", &AuthSetParams{
		Auth: OAuth{
			Type:    AuthTypeOAuth,
			Refresh: "refresh_token",
			Access:  "access_token",
			Expires: 3600,
		},
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !result {
		t.Errorf("Expected result to be true, got false")
	}
}

func TestAuthService_Set_WithDirectory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify directory query parameter
		if r.URL.Query().Get("directory") != "/custom/dir" {
			t.Errorf("Expected directory query param '/custom/dir', got '%s'", r.URL.Query().Get("directory"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(true)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	dir := "/custom/dir"
	_, err = client.Auth.Set(context.Background(), "provider-id", &AuthSetParams{
		Auth: ApiAuth{
			Type: AuthTypeAPI,
			Key:  "api_key",
		},
		Directory: &dir,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestAuthService_Set_MissingID(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Auth.Set(context.Background(), "", &AuthSetParams{
		Auth: ApiAuth{Type: AuthTypeAPI, Key: "test"},
	})

	if err == nil {
		t.Fatal("Expected error for missing id, got nil")
	}
	if err.Error() != "id is required" {
		t.Errorf("Expected 'id is required' error, got %v", err)
	}
}

func TestAuthService_Set_MissingParams(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Auth.Set(context.Background(), "provider-id", nil)

	if err == nil {
		t.Fatal("Expected error for nil params, got nil")
	}
	if err.Error() != "params is required" {
		t.Errorf("Expected 'params is required' error, got %v", err)
	}
}

func TestAuthService_Set_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Auth.Set(context.Background(), "provider-id", &AuthSetParams{
		Auth: ApiAuth{Type: AuthTypeAPI, Key: "key"},
	})

	if err == nil {
		t.Fatal("Expected error from server, got nil")
	}
}

func TestAuthService_Set_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("not json"))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Auth.Set(context.Background(), "provider-id", &AuthSetParams{
		Auth: ApiAuth{Type: AuthTypeAPI, Key: "key"},
	})

	if err == nil {
		t.Fatal("Expected JSON decode error, got nil")
	}
}

func TestAuthSetParams_MarshalJSON_OAuth(t *testing.T) {
	params := AuthSetParams{
		Auth: OAuth{
			Type:    AuthTypeOAuth,
			Refresh: "ref",
			Access:  "acc",
			Expires: 3600,
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal params: %v", err)
	}

	// Verify that only Auth field is marshaled
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if _, hasDirectory := result["directory"]; hasDirectory {
		t.Error("Expected directory field to not be in marshaled JSON")
	}
	if result["type"] != "oauth" {
		t.Errorf("Expected type 'oauth', got %v", result["type"])
	}
}

func TestAuthSetParams_MarshalJSON_ApiAuth(t *testing.T) {
	params := AuthSetParams{
		Auth: ApiAuth{
			Type: AuthTypeAPI,
			Key:  "my-api-key",
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal params: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if result["type"] != "api" {
		t.Errorf("Expected type 'api', got %v", result["type"])
	}
	if result["key"] != "my-api-key" {
		t.Errorf("Expected key 'my-api-key', got %v", result["key"])
	}
}

func TestAuthSetParams_MarshalJSON_WellKnownAuth(t *testing.T) {
	params := AuthSetParams{
		Auth: WellKnownAuth{
			Type:  AuthTypeWellKnown,
			Key:   "wk-key",
			Token: "wk-token",
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal params: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if result["type"] != "wellknown" {
		t.Errorf("Expected type 'wellknown', got %v", result["type"])
	}
	if result["key"] != "wk-key" {
		t.Errorf("Expected key 'wk-key', got %v", result["key"])
	}
	if result["token"] != "wk-token" {
		t.Errorf("Expected token 'wk-token', got %v", result["token"])
	}
}

func TestAuthSetParams_URLQuery(t *testing.T) {
	tests := []struct {
		name     string
		params   AuthSetParams
		expected url.Values
	}{
		{
			name: "with directory",
			params: AuthSetParams{
				Auth:      ApiAuth{Type: AuthTypeAPI, Key: "k"},
				Directory: Ptr("/test/dir"),
			},
			expected: url.Values{"directory": []string{"/test/dir"}},
		},
		{
			name: "without directory",
			params: AuthSetParams{
				Auth: ApiAuth{Type: AuthTypeAPI, Key: "k"},
			},
			expected: url.Values{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := tt.params.URLQuery()
			if err != nil {
				t.Fatalf("URLQuery failed: %v", err)
			}
			if len(query) != len(tt.expected) {
				t.Errorf("Expected %d query params, got %d", len(tt.expected), len(query))
			}
			for key, expectedVals := range tt.expected {
				actualVals := query[key]
				if len(actualVals) != len(expectedVals) {
					t.Errorf("Expected %d values for key %s, got %d", len(expectedVals), key, len(actualVals))
					continue
				}
				for i, expectedVal := range expectedVals {
					if actualVals[i] != expectedVal {
						t.Errorf("Expected value %s for key %s, got %s", expectedVal, key, actualVals[i])
					}
				}
			}
		})
	}
}
