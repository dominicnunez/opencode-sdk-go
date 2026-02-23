# SDK Redesign Plan

**Module path:** `github.com/dominicnunez/opencode-sdk-go`

**Goal:** Transform from generated SDK to proper hand-crafted Go SDK optimized for developer iteration speed.

---

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Auth | No changes | Local server handles auth, no API key needed |
| Streaming | Keep iterator pattern | Matches stdlib (database/sql), clean error handling |
| Client options | Functional options `WithX()` | Idiomatic, extensible, composable |
| Errors | Typed errors + sentinel | `errors.Is(err, ErrNotFound)` for instant debugging |
| Params | Simple structs with `omitempty` | Remove `param.Field` complexity |
| Module path | `github.com/dominicnunez/opencode-sdk-go` | Matches Springfield namespace |

---

## File Structure (Target)

```
opencode-sdk-go/
├── client.go           # Client + options
├── errors.go           # Typed errors
├── session.go          # Session service + types
├── event.go            # Event service + types (includes streaming)
├── agent.go            # Agent service + types
├── config.go           # Config service + types
├── file.go             # File service + types
├── find.go             # Find service + types
├── path.go             # Path service + types
├── project.go          # Project service + types
├── command.go          # Command service + types
├── tui.go              # TUI service + types
├── aliases.go          # Type aliases (keep)
├── field.go            # Field type (simplify)
├── shared/             # Shared types (keep)
├── internal/
│   ├── apijson/        # Keep - JSON encoding
│   ├── apiquery/       # Keep - query encoding  
│   ├── apiform/        # Keep - form encoding
│   ├── apierror/       # Keep - error types
│   ├── requestconfig/  # Simplify
│   └── param/          # DELETE after params simplified
└── packages/
    └── ssestream/      # Keep - SSE streaming
```

---

## Phase 1: Module Path Update

### Files to update
- `go.mod` - change module path
- All `.go` files - update import paths from `github.com/anomalyco/opencode-sdk-go` to `github.com/dominicnunez/opencode-sdk-go`

### Commands
```bash
# Update go.mod
# Update all import paths with sed/find-replace
```

---

## Phase 2: Typed Errors

### Create `errors.go`

```go
package opencode

import (
    "errors"
    "fmt"
)

var (
    ErrNotFound      = errors.New("resource not found")
    ErrUnauthorized  = errors.New("unauthorized")
    ErrRateLimited   = errors.New("rate limited")
    ErrInvalidRequest = errors.New("invalid request")
    ErrInternal      = errors.New("internal server error")
)

type APIError struct {
    StatusCode int
    Message    string
    RequestID  string
    Body       string
}

func (e *APIError) Error() string {
    if e.RequestID != "" {
        return fmt.Sprintf("%s (status %d, request %s)", e.Message, e.StatusCode, e.RequestID)
    }
    return fmt.Sprintf("%s (status %d)", e.Message, e.StatusCode)
}

func (e *APIError) Is(target error) bool {
    switch {
    case e.StatusCode == 404:
        return errors.Is(target, ErrNotFound)
    case e.StatusCode == 401 || e.StatusCode == 403:
        return errors.Is(target, ErrUnauthorized)
    case e.StatusCode == 429:
        return errors.Is(target, ErrRateLimited)
    case e.StatusCode >= 400 && e.StatusCode < 500:
        return errors.Is(target, ErrInvalidRequest)
    case e.StatusCode >= 500:
        return errors.Is(target, ErrInternal)
    }
    return false
}

func IsNotFoundError(err error) bool {
    var apiErr *APIError
    if errors.As(err, &apiErr) {
        return apiErr.StatusCode == 404
    }
    return errors.Is(err, ErrNotFound)
}
```

---

## Phase 3: Simplified Client

### Rewrite `client.go`

```go
package opencode

import (
    "context"
    "net/http"
    "os"
    "time"

    "github.com/dominicnunez/opencode-sdk-go/internal/requestconfig"
)

const (
    DefaultBaseURL    = "http://localhost:54321"
    DefaultTimeout    = 30 * time.Second
    DefaultMaxRetries = 2
)

type Client struct {
    baseURL    string
    httpClient *http.Client
    maxRetries int
    timeout    time.Duration
    
    // Services
    Session  *SessionService
    Event    *EventService
    Agent    *AgentService
    App      *AppService
    Config   *ConfigService
    File     *FileService
    Find     *FindService
    Path     *PathService
    Project  *ProjectService
    Command  *CommandService
    Tui      *TuiService
}

type ClientOption func(*Client) error

func NewClient(opts ...ClientOption) (*Client, error) {
    c := &Client{
        baseURL:    os.Getenv("OPENCODE_BASE_URL"),
        httpClient: http.DefaultClient,
        maxRetries: DefaultMaxRetries,
        timeout:    DefaultTimeout,
    }
    
    // Apply defaults
    if c.baseURL == "" {
        c.baseURL = DefaultBaseURL
    }
    
    // Apply options
    for _, opt := range opts {
        if err := opt(c); err != nil {
            return nil, err
        }
    }
    
    // Initialize services
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
    
    return c, nil
}

func WithBaseURL(url string) ClientOption {
    return func(c *Client) error {
        c.baseURL = url
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
        c.maxRetries = n
        return nil
    }
}

// Execute makes a raw HTTP request with the client's configuration.
func (c *Client) Execute(ctx context.Context, method, path string, params, result interface{}) error {
    // Use internal requestconfig
}
```

---

## Phase 4: Service Pattern

### Rewrite each service (starting with session.go)

**Pattern:**

```go
package opencode

import (
    "context"
    "net/http"
)

type SessionService struct {
    client *Client
}

// Create a new session
func (s *SessionService) Create(ctx context.Context, params *SessionCreateParams) (*Session, error) {
    if params == nil {
        params = &SessionCreateParams{}
    }
    var result Session
    err := s.client.Execute(ctx, http.MethodPost, "session", params, &result)
    if err != nil {
        return nil, err
    }
    return &result, nil
}

// Get a session by ID
func (s *SessionService) Get(ctx context.Context, id string) (*Session, error) {
    var result Session
    err := s.client.Execute(ctx, http.MethodGet, "session/"+id, nil, &result)
    if err != nil {
        return nil, err
    }
    return &result, nil
}

// List all sessions
func (s *SessionService) List(ctx context.Context) ([]Session, error) {
    var result []Session
    err := s.client.Execute(ctx, http.MethodGet, "session", nil, &result)
    if err != nil {
        return nil, err
    }
    return result, nil
}

// Delete a session
func (s *SessionService) Delete(ctx context.Context, id string) error {
    return s.client.Execute(ctx, http.MethodDelete, "session/"+id, nil, nil)
}

// Update a session
func (s *SessionService) Update(ctx context.Context, id string, params *SessionUpdateParams) (*Session, error) {
    if params == nil {
        params = &SessionUpdateParams{}
    }
    var result Session
    err := s.client.Execute(ctx, http.MethodPatch, "session/"+id, params, &result)
    if err != nil {
        return nil, err
    }
    return &result, nil
}
```

---

## Phase 5: Simplified Params

### Current (complex):
```go
type SessionNewParams struct {
    ParentID param.Field[string] `json:"parentId"`
}
```

### New (simple):
```go
type SessionCreateParams struct {
    ParentID string `json:"parentId,omitempty"`
}

type SessionUpdateParams struct {
    // fields
}

type SessionListParams struct {
    // fields
}
```

**Action:** Replace all `param.Field[T]` with direct `T` types, use `omitempty` for optional fields.

---

## Phase 6: Simplified Response Types

### Current (complex):
```go
type Session struct {
    ID        string      `json:"id,required"`
    JSON      sessionJSON `json:"-"`
}
type sessionJSON struct {
    raw string
    // ...
}
```

### New (simple):
```go
type Session struct {
    ID        string    `json:"id"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
    // Direct fields, no JSON metadata
}
```

**Action:** Remove `JSON` metadata fields from response types. Keep raw JSON access via optional `.RawJSON()` method if needed.

---

## Phase 7: Simplify internal/requestconfig

1. Remove unused code paths
2. Simplify error handling to use typed errors
3. Remove `param.Field` handling
4. Keep core HTTP execution logic

---

## Phase 8: Delete internal/param

After all params are simplified to use direct types, the `internal/param` package becomes unused. Delete it.

---

## Phase 9: Update Tests

1. Update import paths
2. Update test cases to use new params pattern
3. Add tests for typed errors
4. Run all tests

---

## Phase 10: Final Cleanup

1. Run `go mod tidy`
2. Run `golangci-lint run ./...`
3. Run `go test ./...`
4. Verify build succeeds

---

## Execution Checklist

- [ ] Phase 1: Module path update
- [ ] Phase 2: Create errors.go
- [ ] Phase 3: Rewrite client.go
- [ ] Phase 4: Rewrite session.go
- [ ] Phase 4: Rewrite event.go
- [ ] Phase 4: Rewrite agent.go
- [ ] Phase 4: Rewrite app.go
- [ ] Phase 4: Rewrite config.go
- [ ] Phase 4: Rewrite file.go
- [ ] Phase 4: Rewrite find.go
- [ ] Phase 4: Rewrite path.go
- [ ] Phase 4: Rewrite project.go
- [ ] Phase 4: Rewrite command.go
- [ ] Phase 4: Rewrite tui.go
- [ ] Phase 4: Rewrite sessionpermission.go
- [ ] Phase 5: Simplify params (remove param.Field)
- [ ] Phase 6: Simplify response types
- [ ] Phase 7: Simplify requestconfig
- [ ] Phase 8: Delete internal/param
- [ ] Phase 9: Update tests
- [ ] Phase 10: Final cleanup

---

## Notes

- Keep streaming pattern (`ssestream.Stream[T]`) - already good
- Keep encoding internals (`apijson`, `apiquery`, `apiform`) - they work
- Keep `shared/` package - contains shared types
- Service method naming: `Create`, `Get`, `List`, `Update`, `Delete` (standard CRUD)
