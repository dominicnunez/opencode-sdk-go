package opencode

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
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
	backoffJitterDiv = 2
	maxErrorBodySize = 1 << 20 // 1 MB

	defaultMaxSuccessBodySize int64 = 8 << 20 // 8 MB
)

const (
	dotPathSegment       = "."
	doubleDotPathSegment = ".."
	maxRetryAfterSeconds = int64(1<<63-1) / int64(time.Second)
)

var retryBackoffRandInt63n = rand.Int63n

func blockRedirects(*http.Request, []*http.Request) error {
	return http.ErrUseLastResponse
}

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	maxRetries int
	timeout    time.Duration
	userAgent  string
	baseURLSet bool
	// maxSuccessBodySize limits successful JSON response bodies.
	// A value of 0 disables the limit.
	maxSuccessBodySize int64

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
	host := parsed.Hostname()
	if host == "" {
		return nil, errors.New("base URL must include a host")
	}
	if parsed.User != nil {
		return nil, errors.New("base URL must not include user info; configure authentication explicitly")
	}
	if parsed.Scheme == "http" && !isLoopbackHost(host) {
		return nil, fmt.Errorf("base URL must use https for non-loopback hosts, got %q", parsed.Hostname())
	}
	if err := validateBaseURLQuery(parsed.Query()); err != nil {
		return nil, err
	}
	if !strings.HasSuffix(parsed.Path, "/") {
		parsed.Path += "/"
	}
	return parsed, nil
}

func validateBaseURLQuery(query url.Values) error {
	for key := range query {
		return fmt.Errorf("base URL must not include query parameters (found %q)", key)
	}
	return nil
}

func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func NewClient(opts ...ClientOption) (*Client, error) {
	defaultBaseURL, err := parseBaseURL(DefaultBaseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		baseURL:            defaultBaseURL,
		httpClient:         &http.Client{CheckRedirect: blockRedirects},
		maxRetries:         DefaultMaxRetries,
		timeout:            DefaultTimeout,
		userAgent:          "Opencode/Go " + internal.PackageVersion,
		maxSuccessBodySize: defaultMaxSuccessBodySize,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	if !c.baseURLSet {
		if rawURL := os.Getenv("OPENCODE_BASE_URL"); rawURL != "" {
			parsed, err := parseBaseURL(rawURL)
			if err != nil {
				return nil, err
			}
			c.baseURL = parsed
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
		c.baseURLSet = true
		return nil
	}
}

func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) error {
		if hc == nil {
			return errors.New("http client cannot be nil")
		}
		clone := *hc
		clone.CheckRedirect = blockRedirects
		c.httpClient = &clone
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

func WithMaxSuccessBodySize(n int64) ClientOption {
	return func(c *Client) error {
		if n < 0 {
			return errors.New("max success body size cannot be negative")
		}
		c.maxSuccessBodySize = n
		return nil
	}
}

type countingReader struct {
	reader io.Reader
	count  int64
}

func (r *countingReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	r.count += int64(n)
	return n, err
}

func decodeSuccessBodyLimitExceeded(reader *countingReader, limit int64) bool {
	return reader != nil && limit > 0 && reader.count > limit
}

func isMethodRetryable(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodPut, http.MethodDelete:
		return true
	default:
		return false
	}
}

func (c *Client) do(ctx context.Context, method, path string, params, result interface{}) error {
	if ctx == nil {
		return ErrContextRequired
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	resp, err := c.doRaw(ctx, method, path, params)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if result == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}

	successReader := io.Reader(resp.Body)
	var responseCounter *countingReader
	if c.maxSuccessBodySize > 0 {
		responseCounter = &countingReader{
			reader: io.LimitReader(resp.Body, c.maxSuccessBodySize+1),
		}
		successReader = responseCounter
	}

	dec := json.NewDecoder(successReader)
	if err := dec.Decode(result); err != nil {
		if decodeSuccessBodyLimitExceeded(responseCounter, c.maxSuccessBodySize) {
			return fmt.Errorf("decode %s %s response: response body exceeds %d bytes limit", method, path, c.maxSuccessBodySize)
		}
		return fmt.Errorf("decode %s %s response: %w", method, path, err)
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("decode %s %s response: unexpected trailing JSON value", method, path)
		}
		if decodeSuccessBodyLimitExceeded(responseCounter, c.maxSuccessBodySize) {
			return fmt.Errorf("decode %s %s response: response body exceeds %d bytes limit", method, path, c.maxSuccessBodySize)
		}
		return fmt.Errorf("decode %s %s response: %w", method, path, err)
	}
	if decodeSuccessBodyLimitExceeded(responseCounter, c.maxSuccessBodySize) {
		return fmt.Errorf("decode %s %s response: response body exceeds %d bytes limit", method, path, c.maxSuccessBodySize)
	}
	return nil
}

// buildURL joins a validated endpoint path to the base URL and applies query
// parameters from the params struct (if it implements URLQuery).
func (c *Client) buildURL(path string, params interface{}) (*url.URL, error) {
	if err := validateEndpointPath(path); err != nil {
		return nil, err
	}

	fullURL := *c.baseURL
	fullURL.Path = joinEndpointPath(c.baseURL.Path, path)
	fullURL.RawPath = ""
	mergedQuery := fullURL.Query()

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
	return &fullURL, nil
}

func validateEndpointPath(endpointPath string) error {
	for _, segment := range strings.Split(endpointPath, "/") {
		if segment == "" {
			continue
		}
		if isDotPathSegment(segment) {
			return fmt.Errorf("endpoint path contains forbidden dot-segment %q", segment)
		}

		decodedSegment, err := url.PathUnescape(segment)
		if err != nil {
			continue
		}
		if isDotPathSegment(decodedSegment) {
			return fmt.Errorf("endpoint path contains forbidden dot-segment %q", decodedSegment)
		}
		if strings.Contains(decodedSegment, "/") || strings.Contains(decodedSegment, "\\") {
			return fmt.Errorf("endpoint path contains forbidden separator in segment %q", decodedSegment)
		}
	}

	return nil
}

func isDotPathSegment(segment string) bool {
	return segment == dotPathSegment || segment == doubleDotPathSegment
}

func joinEndpointPath(basePath, endpointPath string) string {
	trimmedBase := strings.TrimSuffix(basePath, "/")
	trimmedEndpoint := strings.TrimPrefix(endpointPath, "/")

	if trimmedBase == "" {
		if trimmedEndpoint == "" {
			return "/"
		}
		return "/" + trimmedEndpoint
	}
	if trimmedEndpoint == "" {
		return trimmedBase + "/"
	}
	return trimmedBase + "/" + trimmedEndpoint
}

func (c *Client) doRaw(ctx context.Context, method, path string, params interface{}) (*http.Response, error) {
	if ctx == nil {
		return nil, ErrContextRequired
	}

	fullURL, err := c.buildURL(path, params)
	if err != nil {
		return nil, err
	}

	var bodyBytes []byte

	// Handle JSON body for POST/PATCH/PUT
	if params != nil && method != http.MethodGet && method != http.MethodDelete {
		var err error
		bodyBytes, err = json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
	}

	// Build request with retry loop
	var resp *http.Response
	var lastErr error
	maxRequestRetries := c.maxRetries
	if !isMethodRetryable(method) {
		maxRequestRetries = 0
	}

	for attempt := 0; attempt <= maxRequestRetries; attempt++ {
		var body io.Reader
		if len(bodyBytes) > 0 {
			body = bytes.NewReader(bodyBytes)
		}

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
			resp = nil
		}

		// Check context cancellation
		if ctx.Err() != nil {
			if resp != nil {
				_ = resp.Body.Close()
			}
			return nil, ctx.Err()
		}

		// Success — only 2xx responses are valid JSON API results.
		if lastErr == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp, nil
		}

		// Any non-2xx HTTP response is surfaced as an API error.
		// Retry only retryable statuses (408, 429, 5xx).
		if lastErr == nil {
			if !isRetryableStatus(resp.StatusCode) || attempt >= maxRequestRetries {
				return nil, readAPIError(resp, maxErrorBodySize)
			}

			retryDelay := retryDelayWithServerGuidance(attempt, resp, ctx, time.Now())

			// Drain and close body before retry to enable connection reuse
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()

			timer := time.NewTimer(retryDelay)
			select {
			case <-timer.C:
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			}
			continue
		}

		// No more retries — let the loop condition handle exit
		if attempt >= maxRequestRetries {
			continue
		}

		// Backoff before retry. Transport errors (lastErr != nil) skip the
		// delay on the penultimate attempt because the final retry is
		// best-effort — sleeping up to maxBackoff for a likely-unreachable
		// host wastes wall-clock time without improving success odds.
		skipDelay := lastErr != nil && attempt == maxRequestRetries-1
		if !skipDelay {
			delay := retryBackoffDelay(attempt)
			timer := time.NewTimer(delay)
			select {
			case <-timer.C:
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			}
		}
	}

	// All retries exhausted. HTTP errors return structured *APIError inside
	// the loop, so this path should only occur for transport failures.
	if lastErr != nil {
		return nil, fmt.Errorf("%s %s request failed after %d retries: %w", method, path, maxRequestRetries, lastErr)
	}
	return nil, fmt.Errorf("%s %s request failed after %d retries", method, path, maxRequestRetries)
}

func retryBackoffDelay(attempt int) time.Duration {
	delay := retryBackoffBaseDelay(attempt)
	jitterSpan := delay / backoffJitterDiv
	if jitterSpan <= 0 {
		return delay
	}
	return delay - jitterSpan + time.Duration(retryBackoffRandInt63n(int64(jitterSpan)+1))
}

func retryBackoffBaseDelay(attempt int) time.Duration {
	delay := initialBackoff
	for i := 0; i < attempt; i++ {
		if delay >= maxBackoff || delay > maxBackoff/2 {
			return maxBackoff
		}
		delay *= 2
	}
	if delay <= 0 || delay > maxBackoff {
		return maxBackoff
	}
	return delay
}

func retryDelayWithServerGuidance(attempt int, resp *http.Response, ctx context.Context, now time.Time) time.Duration {
	delay := retryBackoffDelay(attempt)
	fromServer := false

	if resp != nil {
		if retryAfterDelay, ok := parseRetryAfterDelay(resp.Header.Get("Retry-After"), now); ok {
			delay = retryAfterDelay
			fromServer = true
		}
	}
	if !fromServer && delay > maxBackoff {
		delay = maxBackoff
	}
	if delay < 0 {
		delay = 0
	}

	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return 0
		}
		if delay > remaining {
			delay = remaining
		}
	}

	return delay
}

func parseRetryAfterDelay(headerValue string, now time.Time) (time.Duration, bool) {
	value := strings.TrimSpace(headerValue)
	if value == "" {
		return 0, false
	}

	seconds, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		if seconds < 0 || seconds > maxRetryAfterSeconds {
			return 0, false
		}
		return time.Duration(seconds) * time.Second, true
	}

	retryAt, err := http.ParseTime(value)
	if err != nil {
		return 0, false
	}
	delay := retryAt.Sub(now)
	if delay < 0 {
		return 0, true
	}
	return delay, true
}
