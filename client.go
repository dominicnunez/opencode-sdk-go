package opencode

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dominicnunez/opencode-sdk-go/internal"
)

const (
	DefaultBaseURL    = "http://localhost:54321"
	DefaultTimeout    = 30 * time.Second
	DefaultMaxRetries = 2

	maxRetryCap      = 10
	initialBackoff   = 500 * time.Millisecond
	maxBackoff       = 8 * time.Second
	maxErrorBodySize = 1 << 20 // 1 MB
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	maxRetries int
	timeout    time.Duration
	userAgent  string

	Session *SessionService
	Event   *EventService
	Agent   *AgentService
	App     *AppService
	Auth    *AuthService
	Config  *ConfigService
	File    *FileService
	Find    *FindService
	Mcp     *McpService
	Path    *PathService
	Project *ProjectService
	Command *CommandService
	Tui     *TuiService
	Tool    *ToolService
}

type ClientOption func(*Client) error

func parseBaseURL(rawURL string) (*url.URL, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("parse base URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("base URL must use http or https scheme, got %q", parsed.Scheme)
	}
	if !strings.HasSuffix(parsed.Path, "/") {
		parsed.Path += "/"
	}
	return parsed, nil
}

func NewClient(opts ...ClientOption) (*Client, error) {
	rawURL := os.Getenv("OPENCODE_BASE_URL")
	if rawURL == "" {
		rawURL = DefaultBaseURL
	}
	parsed, err := parseBaseURL(rawURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		baseURL:    parsed,
		httpClient: &http.Client{},
		maxRetries: DefaultMaxRetries,
		timeout:    DefaultTimeout,
		userAgent:  "Opencode/Go " + internal.PackageVersion,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	c.Session = &SessionService{client: c}
	c.Event = &EventService{client: c}
	c.Agent = &AgentService{client: c}
	c.App = &AppService{client: c}
	c.Auth = &AuthService{client: c}
	c.Config = &ConfigService{client: c}
	c.File = &FileService{client: c}
	c.Find = &FindService{client: c}
	c.Mcp = &McpService{client: c}
	c.Path = &PathService{client: c}
	c.Project = &ProjectService{client: c}
	c.Command = &CommandService{client: c}
	c.Tui = &TuiService{client: c}
	c.Tool = &ToolService{client: c}

	c.Session.Permissions = &SessionPermissionService{client: c}

	return c, nil
}

func WithBaseURL(rawURL string) ClientOption {
	return func(c *Client) error {
		u, err := parseBaseURL(rawURL)
		if err != nil {
			return err
		}
		c.baseURL = u
		return nil
	}
}

func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) error {
		if hc == nil {
			return errors.New("http client cannot be nil")
		}
		c.httpClient = hc
		return nil
	}
}

func WithTimeout(d time.Duration) ClientOption {
	return func(c *Client) error {
		if d <= 0 {
			return errors.New("timeout must be positive")
		}
		c.timeout = d
		return nil
	}
}

func WithMaxRetries(n int) ClientOption {
	return func(c *Client) error {
		if n < 0 {
			return errors.New("max retries cannot be negative")
		}
		if n > maxRetryCap {
			return fmt.Errorf("max retries cannot exceed %d", maxRetryCap)
		}
		c.maxRetries = n
		return nil
	}
}

func (c *Client) do(ctx context.Context, method, path string, params, result interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.doRaw(ctx, method, path, params)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if result == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decode %s %s response: %w", method, path, err)
	}
	return nil
}

// buildURL resolves path against the base URL and merges query parameters from
// the base URL and the params struct (if it implements URLQuery).
func (c *Client) buildURL(path string, params interface{}) (*url.URL, error) {
	fullURL := c.baseURL.ResolveReference(&url.URL{Path: path})
	mergedQuery := fullURL.Query()
	for k, vs := range c.baseURL.Query() {
		mergedQuery[k] = vs
	}

	if params != nil {
		if queryer, ok := params.(interface{ URLQuery() (url.Values, error) }); ok {
			query, err := queryer.URLQuery()
			if err != nil {
				return nil, fmt.Errorf("encode query params: %w", err)
			}
			for k, vs := range query {
				mergedQuery[k] = vs
			}
		}
	}
	fullURL.RawQuery = mergedQuery.Encode()
	return fullURL, nil
}

func (c *Client) doRaw(ctx context.Context, method, path string, params interface{}) (*http.Response, error) {
	fullURL, err := c.buildURL(path, params)
	if err != nil {
		return nil, err
	}

	var body io.Reader

	// Handle JSON body for POST/PATCH/PUT
	if params != nil && method != http.MethodGet && method != http.MethodDelete {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(params); err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		body = &buf
	}

	// Build request with retry loop
	var resp *http.Response
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, method, fullURL.String(), body)
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}

		// Set headers
		if method != http.MethodGet && method != http.MethodDelete {
			req.Header.Set("Content-Type", "application/json")
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", c.userAgent)

		// Execute request
		resp, lastErr = c.httpClient.Do(req)

		// Close body from transport errors that still return a response
		// (e.g., custom HTTP clients that don't follow stdlib's contract)
		if lastErr != nil && resp != nil {
			_ = resp.Body.Close()
		}

		// Check context cancellation
		if ctx.Err() != nil {
			if resp != nil {
				_ = resp.Body.Close()
			}
			return nil, ctx.Err()
		}

		// Success - return response
		if lastErr == nil && resp.StatusCode < 400 {
			return resp, nil
		}

		// Error response - don't retry client errors (4xx except specific cases)
		if lastErr == nil && resp.StatusCode >= 400 {
			if !isRetryableStatus(resp.StatusCode) || attempt >= c.maxRetries {
				return nil, readAPIError(resp, maxErrorBodySize)
			}

			// Drain and close body before retry to enable connection reuse
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}

		// Wait before retry (exponential backoff)
		if attempt < c.maxRetries {
			delay := initialBackoff * (1 << attempt)
			if delay <= 0 || delay > maxBackoff {
				delay = maxBackoff
			}
			timer := time.NewTimer(delay)
			select {
			case <-timer.C:
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			}

			// Reset body for retry if needed
			if params != nil && method != http.MethodGet && method != http.MethodDelete {
				var buf bytes.Buffer
				if err := json.NewEncoder(&buf).Encode(params); err != nil {
					return nil, fmt.Errorf("marshal request body for retry: %w", err)
				}
				body = &buf
			}
		}
	}

	// All retries exhausted — only reachable via transport errors (lastErr != nil).
	// HTTP errors return structured *APIError inside the loop.
	return nil, fmt.Errorf("request failed after %d retries: %w", c.maxRetries, lastErr)
}
