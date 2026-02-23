package opencode

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dominicnunez/opencode-sdk-go/internal/requestconfig"
	"github.com/dominicnunez/opencode-sdk-go/option"
)

const (
	DefaultBaseURL    = "http://localhost:54321"
	DefaultTimeout    = 30 * time.Second
	DefaultMaxRetries = 2
)

type Client struct {
	baseURL        string
	httpClient     *http.Client
	maxRetries     int
	timeout        time.Duration
	defaultOptions []option.RequestOption

	Session *SessionService
	Event   *EventService
	Agent   *AgentService
	App     *AppService
	Config  *ConfigService
	File    *FileService
	Find    *FindService
	Path    *PathService
	Project *ProjectService
	Command *CommandService
	Tui     *TuiService
}

type ClientOption func(*Client) error

func NewClient(opts ...ClientOption) (*Client, error) {
	c := &Client{
		baseURL:    os.Getenv("OPENCODE_BASE_URL"),
		httpClient: http.DefaultClient,
		maxRetries: DefaultMaxRetries,
		timeout:    DefaultTimeout,
	}

	if c.baseURL == "" {
		c.baseURL = DefaultBaseURL
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
	c.Config = &ConfigService{client: c}
	c.File = &FileService{client: c}
	c.Find = &FindService{client: c}
	c.Path = &PathService{client: c}
	c.Project = &ProjectService{client: c}
	c.Command = &CommandService{client: c}
	c.Tui = &TuiService{client: c}

	c.Session.Permissions = &SessionPermissionService{client: c}

	return c, nil
}

func WithBaseURL(rawURL string) ClientOption {
	return func(c *Client) error {
		u, err := url.Parse(rawURL)
		if err != nil {
			return fmt.Errorf("parse base URL: %w", err)
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			return fmt.Errorf("base URL must use http or https scheme, got %q", u.Scheme)
		}
		c.baseURL = rawURL
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
		c.timeout = d
		return nil
	}
}

func WithMaxRetries(n int) ClientOption {
	return func(c *Client) error {
		if n < 0 {
			return errors.New("max retries cannot be negative")
		}
		if n > 10 {
			return errors.New("max retries cannot exceed 10")
		}
		c.maxRetries = n
		return nil
	}
}

func WithRequestOption(opt option.RequestOption) ClientOption {
	return func(c *Client) error {
		c.defaultOptions = append(c.defaultOptions, opt)
		return nil
	}
}

func (c *Client) do(ctx context.Context, method, path string, params, result interface{}, opts ...option.RequestOption) error {
	allOpts := []option.RequestOption{
		option.WithBaseURL(c.baseURL),
		option.WithHTTPClient(c.httpClient),
		option.WithMaxRetries(c.maxRetries),
	}
	allOpts = append(allOpts, c.defaultOptions...)
	allOpts = append(allOpts, opts...)
	return requestconfig.ExecuteNewRequest(ctx, method, path, params, result, allOpts...)
}
