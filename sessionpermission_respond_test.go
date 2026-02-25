package opencode

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSessionPermissionService_Respond_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Verify path
		expectedPath := "/session/sess-123/permissions/perm-456"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Verify request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		var params struct {
			Response PermissionResponse `json:"response"`
		}
		if err := json.Unmarshal(body, &params); err != nil {
			t.Fatalf("Failed to unmarshal request body: %v", err)
		}
		if params.Response != PermissionResponseOnce {
			t.Errorf("Expected response 'once', got %s", params.Response)
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(true)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	result, err := client.Session.Permissions.Respond(
		context.Background(),
		"sess-123",
		"perm-456",
		&SessionPermissionRespondParams{
			Response: PermissionResponseOnce,
		},
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result {
		t.Errorf("Expected result true, got false")
	}
}

func TestSessionPermissionService_Respond_WithDirectoryParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Verify path
		expectedPath := "/session/sess-123/permissions/perm-456"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Verify query params
		if r.URL.Query().Get("directory") != "/test/dir" {
			t.Errorf("Expected directory query param '/test/dir', got %s", r.URL.Query().Get("directory"))
		}

		// Verify request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		var params struct {
			Response PermissionResponse `json:"response"`
		}
		if err := json.Unmarshal(body, &params); err != nil {
			t.Fatalf("Failed to unmarshal request body: %v", err)
		}
		if params.Response != PermissionResponseAlways {
			t.Errorf("Expected response 'always', got %s", params.Response)
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(true)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	directory := "/test/dir"
	result, err := client.Session.Permissions.Respond(
		context.Background(),
		"sess-123",
		"perm-456",
		&SessionPermissionRespondParams{
			Response:  PermissionResponseAlways,
			Directory: &directory,
		},
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result {
		t.Errorf("Expected result true, got false")
	}
}

func TestSessionPermissionService_Respond_RejectResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		var params struct {
			Response PermissionResponse `json:"response"`
		}
		if err := json.Unmarshal(body, &params); err != nil {
			t.Fatalf("Failed to unmarshal request body: %v", err)
		}
		if params.Response != PermissionResponseReject {
			t.Errorf("Expected response 'reject', got %s", params.Response)
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(true)
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	result, err := client.Session.Permissions.Respond(
		context.Background(),
		"sess-123",
		"perm-456",
		&SessionPermissionRespondParams{
			Response: PermissionResponseReject,
		},
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result {
		t.Errorf("Expected result true, got false")
	}
}

func TestSessionPermissionService_Respond_MissingSessionID(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Session.Permissions.Respond(
		context.Background(),
		"",
		"perm-456",
		&SessionPermissionRespondParams{
			Response: PermissionResponseOnce,
		},
	)
	if err == nil {
		t.Error("Expected error for missing session ID, got nil")
	}
	if err.Error() != "missing required id parameter" {
		t.Errorf("Expected error message 'missing required id parameter', got %s", err.Error())
	}
}

func TestSessionPermissionService_Respond_MissingPermissionID(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Session.Permissions.Respond(
		context.Background(),
		"sess-123",
		"",
		&SessionPermissionRespondParams{
			Response: PermissionResponseOnce,
		},
	)
	if err == nil {
		t.Error("Expected error for missing permission ID, got nil")
	}
	if err.Error() != "missing required permissionID parameter" {
		t.Errorf("Expected error message 'missing required permissionID parameter', got %s", err.Error())
	}
}

func TestSessionPermissionService_Respond_NilParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request body has empty response field (zero value)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		var params struct {
			Response PermissionResponse `json:"response"`
		}
		if err := json.Unmarshal(body, &params); err != nil {
			t.Fatalf("Failed to unmarshal request body: %v", err)
		}
		// Zero value for PermissionResponse is empty string
		if params.Response != "" {
			t.Errorf("Expected empty response, got %s", params.Response)
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

	result, err := client.Session.Permissions.Respond(
		context.Background(),
		"sess-123",
		"perm-456",
		nil,
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result {
		t.Errorf("Expected result true, got false")
	}
}

func TestSessionPermissionService_Respond_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal server error"))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Session.Permissions.Respond(
		context.Background(),
		"sess-123",
		"perm-456",
		&SessionPermissionRespondParams{
			Response: PermissionResponseOnce,
		},
	)
	if err == nil {
		t.Error("Expected error from server, got nil")
	}
}

func TestSessionPermissionService_Respond_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Session.Permissions.Respond(
		context.Background(),
		"sess-123",
		"perm-456",
		&SessionPermissionRespondParams{
			Response: PermissionResponseOnce,
		},
	)
	if err == nil {
		t.Error("Expected error from invalid JSON, got nil")
	}
}

func TestSessionPermissionRespondParams_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		params   SessionPermissionRespondParams
		expected string
	}{
		{
			name: "once response",
			params: SessionPermissionRespondParams{
				Response: PermissionResponseOnce,
			},
			expected: `{"response":"once"}`,
		},
		{
			name: "always response",
			params: SessionPermissionRespondParams{
				Response: PermissionResponseAlways,
			},
			expected: `{"response":"always"}`,
		},
		{
			name: "reject response",
			params: SessionPermissionRespondParams{
				Response: PermissionResponseReject,
			},
			expected: `{"response":"reject"}`,
		},
		{
			name: "with directory (should not be in JSON)",
			params: SessionPermissionRespondParams{
				Response:  PermissionResponseOnce,
				Directory: ptrString("/test/dir"),
			},
			expected: `{"response":"once"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.params)
			if err != nil {
				t.Fatalf("Failed to marshal params: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("Expected JSON %s, got %s", tt.expected, string(data))
			}
		})
	}
}

func TestSessionPermissionRespondParams_URLQuery(t *testing.T) {
	tests := []struct {
		name     string
		params   SessionPermissionRespondParams
		expected url.Values
	}{
		{
			name: "no directory",
			params: SessionPermissionRespondParams{
				Response: PermissionResponseOnce,
			},
			expected: url.Values{},
		},
		{
			name: "with directory",
			params: SessionPermissionRespondParams{
				Response:  PermissionResponseOnce,
				Directory: ptrString("/test/dir"),
			},
			expected: url.Values{"directory": []string{"/test/dir"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, err := tt.params.URLQuery()
			if err != nil {
				t.Fatalf("Failed to get URL query: %v", err)
			}

			// Compare the values
			if len(values) != len(tt.expected) {
				t.Errorf("Expected %d query params, got %d", len(tt.expected), len(values))
			}

			for key, expectedVals := range tt.expected {
				actualVals, ok := values[key]
				if !ok {
					t.Errorf("Expected query param %s not found", key)
					continue
				}
				if len(actualVals) != len(expectedVals) {
					t.Errorf("Expected %d values for %s, got %d", len(expectedVals), key, len(actualVals))
					continue
				}
				for i, expectedVal := range expectedVals {
					if actualVals[i] != expectedVal {
						t.Errorf("Expected value %s for %s[%d], got %s", expectedVal, key, i, actualVals[i])
					}
				}
			}
		})
	}
}
