# Opencode Go API Library

<a href="https://pkg.go.dev/github.com/anomalyco/opencode-sdk-go"><img src="https://pkg.go.dev/badge/github.com/anomalyco/opencode-sdk-go.svg" alt="Go Reference"></a>

The Opencode Go library provides convenient access to the [Opencode REST API](https://opencode.ai/docs)
from applications written in Go. This is a pure, idiomatic Go SDK built with Go's standard library.

## Installation

<!-- x-release-please-start-version -->

```go
import (
	"github.com/anomalyco/opencode-sdk-go" // imported as opencode
)
```

<!-- x-release-please-end -->

Or to pin the version:

<!-- x-release-please-start-version -->

```sh
go get -u 'github.com/anomalyco/opencode-sdk-go@v0.19.2'
```

<!-- x-release-please-end -->

## Requirements

This library requires Go 1.22+.

## Usage

The full API of this library can be found in [api.md](api.md).

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/anomalyco/opencode-sdk-go"
)

func main() {
	client, err := opencode.NewClient(
		opencode.WithBaseURL("http://localhost:8080"),
		opencode.WithTimeout(30 * time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	// List all sessions
	sessions, err := client.Session.List(context.TODO(), opencode.SessionListParams{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d sessions\n", len(sessions))

	// Create a new session
	session, err := client.Session.New(context.TODO(), opencode.SessionNewParams{
		Agent:    "general-purpose",
		Contents: []opencode.MessageUnionParam{
			opencode.UserMessageParam{
				Role:  opencode.MessageRoleUser,
				Parts: []opencode.PartInputParam{
					opencode.TextPartInputParam{
						Text: "Hello, how can you help me?",
					},
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created session: %s\n", session.ID)
}

```

### Request Parameters

Request parameters use idiomatic Go types:
- **Required fields** use direct types: `string`, `int64`, `bool`, `MyType`
- **Optional fields** use pointer types: `*string`, `*int64`, `*bool`, `*MyType`

Optional fields with `nil` values are automatically omitted from the request.

```go
params := opencode.SessionCommandParams{
	// Required field - direct type
	Command: "list-files",

	// Optional field - use pointer
	Directory: opencode.Ptr("./src"),

	// Optional field - nil means omit from request
	Agent: nil,
}
```

For convenience, we provide `Ptr[T](value T)` to create pointers:

```go
params := opencode.SessionUpdateParams{
	Title: opencode.Ptr("My Updated Title"),
}
```

### Response Objects

All response fields use standard Go types decoded with `encoding/json`:
- **String fields**: `string` (empty string if null/missing)
- **Number fields**: `int64`, `float64` (zero if null/missing)
- **Boolean fields**: `bool` (false if null/missing)
- **Object fields**: `SomeType` or `*SomeType` depending on nullability

```go
session, err := client.Session.Get(context.TODO(), "sess_xxx", nil)
if err != nil {
	log.Fatal(err)
}

// Access response fields directly
fmt.Printf("Session ID: %s\n", session.ID)
fmt.Printf("Title: %s\n", session.Title)
fmt.Printf("Created: %s\n", session.Time.Format(time.RFC3339))

// Check for empty values (zero values indicate null/missing)
if session.ParentID == "" {
	fmt.Println("This is a root session (no parent)")
}
```

### Client Configuration

Configure the client using functional options at initialization:

```go
client, err := opencode.NewClient(
	// Set the base URL (defaults to http://localhost:8080)
	opencode.WithBaseURL("https://api.opencode.ai"),

	// Configure timeout (defaults to 2 minutes)
	opencode.WithTimeout(60 * time.Second),

	// Configure max retries (defaults to 2)
	opencode.WithMaxRetries(5),

	// Use a custom HTTP client
	opencode.WithHTTPClient(&http.Client{
		Transport: myCustomTransport,
	}),
)
if err != nil {
	log.Fatal(err)
}
```

Available client options:
- `WithBaseURL(url string)` - Set the API base URL
- `WithTimeout(duration time.Duration)` - Set request timeout
- `WithMaxRetries(retries int)` - Set maximum retry attempts
- `WithHTTPClient(client *http.Client)` - Use a custom HTTP client

### Union Types

Some response types use discriminated unions. Use the `As*()` methods to safely extract the specific variant:

```go
// Messages can be UserMessage or AssistantMessage
messages, err := client.Session.Messages(context.TODO(), "sess_xxx", nil)
if err != nil {
	log.Fatal(err)
}

for _, msg := range messages {
	// Check which type of message
	if user, ok := msg.AsUser(); ok {
		fmt.Printf("User: %s\n", user.Parts[0].AsText().Text)
	} else if assistant, ok := msg.AsAssistant(); ok {
		fmt.Printf("Assistant: %s\n", assistant.Parts[0].AsText().Text)
	}
}

// Parts can be TextPart, FilePart, ToolPart, etc.
if textPart, ok := part.AsText(); ok {
	fmt.Printf("Text: %s\n", textPart.Text)
} else if filePart, ok := part.AsFile(); ok {
	fmt.Printf("File: %s\n", filePart.Source.AsFile().Path)
} else if toolPart, ok := part.AsTool(); ok {
	fmt.Printf("Tool: %s\n", toolPart.Name)
}
```

### Errors

When the API returns a non-success status code, we return an error with type
`*opencode.Error`. This contains the `StatusCode`, `*http.Request`, and
`*http.Response` values of the request, as well as the JSON of the error body
(much like other response objects in the SDK).

To handle errors, we recommend that you use the `errors.As` pattern:

```go
_, err := client.Session.List(context.TODO(), opencode.SessionListParams{})
if err != nil {
	var apierr *opencode.Error
	if errors.As(err, &apierr) {
		println(string(apierr.DumpRequest(true)))  // Prints the serialized HTTP request
		println(string(apierr.DumpResponse(true))) // Prints the serialized HTTP response
	}
	panic(err.Error()) // GET "/session": 400 Bad Request { ... }
}
```

When other errors occur, they are returned unwrapped; for example,
if HTTP transport fails, you might receive `*url.Error` wrapping `*net.OpError`.

### Timeouts

Configure client-level timeout when creating the client:

```go
client, err := opencode.NewClient(
	opencode.WithTimeout(60 * time.Second),
)
```

Or use context for per-request timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

sessions, err := client.Session.List(ctx, opencode.SessionListParams{})
if err != nil {
	log.Fatal(err)
}
```

Note: If a request is retried, the context timeout applies to the total time including all retries.

### Streaming Events

Subscribe to real-time events using Server-Sent Events (SSE):

```go
stream, err := client.Event.ListStreaming(context.TODO(), opencode.EventListStreamingParams{})
if err != nil {
	log.Fatal(err)
}
defer stream.Close()

for stream.Next() {
	event := stream.Current()

	// Handle specific event types
	if msgUpdated, ok := event.AsMessageUpdated(); ok {
		fmt.Printf("Message updated: %s\n", msgUpdated.Data.Info.ID)
	} else if sessionCreated, ok := event.AsSessionCreated(); ok {
		fmt.Printf("Session created: %s\n", sessionCreated.Data.Info.ID)
	}
}

if err := stream.Err(); err != nil {
	log.Fatal(err)
}
```

### Retries

The client automatically retries failed requests with exponential backoff (default: 2 retries).

Retries are triggered for:
- Connection errors
- 408 Request Timeout
- 429 Too Many Requests
- 5xx Server Errors

Configure retry behavior at client initialization:

```go
client, err := opencode.NewClient(
	// Disable retries
	opencode.WithMaxRetries(0),

	// Or increase retries
	opencode.WithMaxRetries(5),
)
```

The retry backoff schedule is: 500ms, 1s, 2s, 4s, 8s (capped at 8 seconds).

### Advanced Usage

#### Working with Authentication

Set authentication credentials for external services:

```go
// OAuth authentication
err := client.Auth.Set(context.TODO(), "provider-id", opencode.AuthSetParams{
	Auth: opencode.OAuth{
		Type:    opencode.AuthTypeOAuth,
		Refresh: "refresh_token_here",
		Access:  opencode.Ptr("access_token_here"),
		Expires: opencode.Ptr(int64(1234567890)),
	},
})

// API key authentication
err := client.Auth.Set(context.TODO(), "provider-id", opencode.AuthSetParams{
	Auth: opencode.ApiAuth{
		Type: opencode.AuthTypeAPI,
		Key:  "api_key_here",
	},
})
```

#### Forking Sessions

Create a new session from an existing one at a specific message:

```go
forked, err := client.Session.Fork(context.TODO(), "sess_xxx", opencode.SessionForkParams{
	MessageID: "msg_yyy",
	Directory: opencode.Ptr("/workspace/project"),
})
```

#### Running Shell Commands

Execute shell commands in a session:

```go
result, err := client.Session.Shell(context.TODO(), "sess_xxx", opencode.SessionShellParams{
	Agent:     "bash",
	Command:   "ls -la",
	Directory: opencode.Ptr("/workspace"),
})
```

### Working with Tools

List available tools and their schemas:

```go
// Get all tool IDs
toolIDs, err := client.Tool.IDs(context.TODO(), opencode.ToolIDsParams{})
if err != nil {
	log.Fatal(err)
}
fmt.Printf("Available tools: %v\n", *toolIDs)

// Get tool schemas for a specific provider/model
tools, err := client.Tool.List(context.TODO(), opencode.ToolListParams{
	Provider: "anthropic",
	Model:    "claude-sonnet-4-5",
})
if err != nil {
	log.Fatal(err)
}

for _, tool := range *tools {
	fmt.Printf("Tool: %s - %s\n", tool.ID, tool.Description)
	// tool.Parameters contains the JSON schema
}
```

### MCP Server Status

Check the status of MCP (Model Context Protocol) servers:

```go
status, err := client.Mcp.Status(context.TODO(), opencode.McpStatusParams{})
if err != nil {
	log.Fatal(err)
}

// status is map[string]interface{} containing dynamic server info
for serverID, serverInfo := range *status {
	fmt.Printf("Server %s: %+v\n", serverID, serverInfo)
}
```

### Custom HTTP Client

Use a custom HTTP client for advanced control over requests:

```go
customClient := &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	},
	Timeout: 30 * time.Second,
}

client, err := opencode.NewClient(
	opencode.WithHTTPClient(customClient),
)
```

For request logging or middleware, wrap the transport:

```go
type LoggingTransport struct {
	Base http.RoundTripper
}

func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	log.Printf("→ %s %s", req.Method, req.URL)

	resp, err := t.Base.RoundTrip(req)

	duration := time.Since(start)
	if err != nil {
		log.Printf("← Error: %v (took %s)", err, duration)
	} else {
		log.Printf("← %d %s (took %s)", resp.StatusCode, resp.Status, duration)
	}

	return resp, err
}

client, err := opencode.NewClient(
	opencode.WithHTTPClient(&http.Client{
		Transport: &LoggingTransport{Base: http.DefaultTransport},
	}),
)
```

## Semantic versioning

This package generally follows [SemVer](https://semver.org/spec/v2.0.0.html) conventions, though certain backwards-incompatible changes may be released as minor versions:

1. Changes to library internals which are technically public but not intended or documented for external use. _(Please open a GitHub issue to let us know if you are relying on such internals.)_
2. Changes that we do not expect to impact the vast majority of users in practice.

We take backwards-compatibility seriously and work hard to ensure you can rely on a smooth upgrade experience.

We are keen for your feedback; please open an [issue](https://www.github.com/anomalyco/opencode-sdk-go/issues) with questions, bugs, or suggestions.

## Contributing

See [the contributing documentation](./CONTRIBUTING.md).
