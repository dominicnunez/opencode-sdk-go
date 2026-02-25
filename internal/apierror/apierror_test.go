package apierror

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		contains []string
	}{
		{
			name: "complete error with all fields",
			err: &Error{
				StatusCode: 404,
				Request: &http.Request{
					Method: "GET",
					URL: &url.URL{
						Scheme: "https",
						Host:   "api.opencode.ai",
						Path:   "/v1/session",
					},
				},
				Response: &http.Response{
					StatusCode: 404,
				},
				Body: `{"error": "not found"}`,
			},
			contains: []string{
				"GET",
				"https://api.opencode.ai/v1/session",
				"404",
				"Not Found",
				`{"error": "not found"}`,
			},
		},
		{
			name: "error with empty body",
			err: &Error{
				StatusCode: 500,
				Request: &http.Request{
					Method: "POST",
					URL: &url.URL{
						Scheme: "https",
						Host:   "api.opencode.ai",
						Path:   "/v1/command",
					},
				},
				Response: &http.Response{
					StatusCode: 500,
				},
				Body: "",
			},
			contains: []string{
				"POST",
				"https://api.opencode.ai/v1/command",
				"500",
				"Internal Server Error",
				"(no response body)",
			},
		},
		{
			name: "error with nil request",
			err: &Error{
				StatusCode: 403,
				Request:    nil,
				Response: &http.Response{
					StatusCode: 403,
				},
				Body: `{"error": "forbidden"}`,
			},
			contains: []string{
				"403",
				"Forbidden",
				`{"error": "forbidden"}`,
			},
		},
		{
			name: "error with nil response",
			err: &Error{
				StatusCode: 400,
				Request: &http.Request{
					Method: "PUT",
					URL: &url.URL{
						Scheme: "https",
						Host:   "api.opencode.ai",
						Path:   "/v1/update",
					},
				},
				Response: nil,
				Body:     `{"error": "bad request"}`,
			},
			contains: []string{
				"PUT",
				"https://api.opencode.ai/v1/update",
				`{"error": "bad request"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Error() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestError_DumpRequest(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		body bool
		want string
	}{
		{
			name: "nil request returns nil",
			err: &Error{
				Request: nil,
			},
			body: true,
			want: "",
		},
		{
			name: "dumps request without body",
			err: &Error{
				Request: &http.Request{
					Method: "GET",
					URL: &url.URL{
						Scheme: "https",
						Host:   "api.opencode.ai",
						Path:   "/v1/session",
					},
					Proto:      "HTTP/1.1",
					ProtoMajor: 1,
					ProtoMinor: 1,
					Header:     http.Header{},
				},
			},
			body: false,
			want: "GET /v1/session HTTP/1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.DumpRequest(tt.body)
			if tt.want == "" {
				if result != nil {
					t.Errorf("DumpRequest() = %q, want nil", result)
				}
			} else {
				if !strings.Contains(string(result), tt.want) {
					t.Errorf("DumpRequest() = %q, should contain %q", result, tt.want)
				}
			}
		})
	}
}

func TestError_DumpResponse(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		body bool
		want string
	}{
		{
			name: "nil response returns nil",
			err: &Error{
				Response: nil,
			},
			body: true,
			want: "",
		},
		{
			name: "dumps response without body",
			err: &Error{
				Response: &http.Response{
					Status:     "404 Not Found",
					StatusCode: 404,
					Proto:      "HTTP/1.1",
					ProtoMajor: 1,
					ProtoMinor: 1,
					Header:     http.Header{},
				},
			},
			body: false,
			want: "HTTP/1.1 404 Not Found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.DumpResponse(tt.body)
			if tt.want == "" {
				if result != nil {
					t.Errorf("DumpResponse() = %q, want nil", result)
				}
			} else {
				if !strings.Contains(string(result), tt.want) {
					t.Errorf("DumpResponse() = %q, should contain %q", result, tt.want)
				}
			}
		})
	}
}
