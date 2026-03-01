# Design

> Findings that describe behavior which is correct by design.
> Managed by sfk willie. Follow the entry format below.
>
> Entry format:
> ### Plain language description
> **Location:** `file/path:line` — optional context
> **Date:** YYYY-MM-DD
> **Reason:** Explanation (can be multiple lines)

### Event struct Data field maps to JSON key "properties"

**Location:** `event.go:326,349,377` — all EventXxx structs
**Date:** 2026-02-28

**Reason:** The Go field is named `Data` for call-site ergonomics (`evt.Data.Version` reads better than `evt.Properties.Version`). The JSON tag `json:"properties"` correctly matches the server's wire format per the OpenAPI spec. Renaming the Go field to `Properties` would make the SDK surface less intuitive, and renaming the JSON tag would break deserialization. The naming mismatch is an intentional tradeoff favoring Go-side readability.

### ListStreaming bypasses Client timeout and retry logic

**Location:** `event.go:49` — uses httpClient.Do directly
**Date:** 2026-02-28

**Reason:** SSE streams are long-lived connections that remain open indefinitely while events arrive. Applying the client's default 30s timeout would prematurely kill every SSE connection. Retries are also inappropriate for streaming — the caller should reconnect at the application level. Callers who need a deadline can set one via `context.WithTimeout` on the context they pass in. This is consistent with how other Go SSE/WebSocket libraries handle timeouts.

### ConfigProviderOptions exposes API key as a plain string field

**Location:** `config.go:1206` — APIKey field
**Date:** 2026-02-28

**Reason:** The struct faithfully reflects the OpenAPI spec schema. The server includes API keys in config responses, and the SDK must deserialize them. Adding `json:"-"` would silently drop data the server returns, breaking callers who need the value. Redaction is a presentation concern that belongs in the caller's logging/serialization layer, not in the SDK's data types.

### McpStatus is an untyped map

**Location:** `mcp.go:17`
**Date:** 2026-02-28

**Reason:** The OpenAPI spec defines the MCP status response with an empty schema (`"schema": {}`), meaning the response shape is intentionally unspecified. `map[string]interface{}` is the correct Go representation of an unconstrained JSON object. Creating a typed struct would impose structure the spec does not guarantee.

### APIError.Body and Message contain the same value

**Location:** `client.go:232-237`, `errors.go:17-22` — APIError construction
**Date:** 2026-02-26

**Reason:** `Body` preserves the raw HTTP response body for callers who need it (e.g. structured error parsing). `Message` is the human-readable error used by `Error()`. They happen to contain the same value today because the server returns plain-text errors, but they serve different semantic purposes. Collapsing them would prevent future differentiation without a breaking change.

### Stream is not safe for concurrent use across goroutines

**Location:** `packages/ssestream/ssestream.go:164-204` — Stream.Next() and Stream.Close()
**Date:** 2026-02-28

**Reason:** `Stream` follows the same contract as `bufio.Scanner`, `sql.Rows`, and other Go iterators: single-goroutine use. Adding `sync.Mutex` synchronization would add overhead to every `Next()` call for a usage pattern (calling `Close()` from a different goroutine) that is not the primary design intent. Callers who need cross-goroutine cancellation should use `context.WithCancel` on the context passed to the streaming method, which will unblock the underlying read.

### SSE stream error not returned directly

**Location:** `event.go:20-51`
**Date:** 2026-02-22

**Reason:** This is a standard pattern for streaming APIs in Go. The stream object must be returned so callers can iterate over events, and embedding the initial connection error in the stream allows a single return signature. The pattern is documented and callers are expected to check `stream.Err()` before iteration, similar to how database rows work.

### POST params serialized as both query string and JSON body

**Location:** `client.go:170-187` — doRaw query and body encoding
**Date:** 2026-02-28

**Reason:** When a params struct implements `URLQuery()` and is used with a POST method, `doRaw` encodes it as both query parameters and JSON body. This works correctly because `queryparams.Marshal` only encodes fields with `query:` tags, and body-only fields use `json:` tags (with `json:"-"` on query-only fields). The separation is enforced by struct tags, not by type splitting. Splitting every params struct into separate query and body types would be a large refactor with no behavioral benefit — the current contract is consistent across all endpoints.

### SessionPromptParamsPart uses `any` for optional fields

**Location:** `session.go:1792-1798` — Metadata, Source, Time fields
**Date:** 2026-02-28

**Reason:** `SessionPromptParamsPart` is the escape-hatch catch-all variant of `SessionPromptParamsPartUnion`. Callers who want type safety should use the typed variants (`TextPartInputParam`, `FilePartInputParam`, `AgentPartInputParam`). The `any` fields exist so callers can construct parts without importing every nested type. The typed variants already provide compile-time safety for callers who want it.

### FilePartSourceParam.Range typed as `any`

**Location:** `session.go:717` — Range field
**Date:** 2026-02-28

**Reason:** `FilePartSourceParam` is the catch-all variant of `FilePartSourceUnionParam` (alongside typed `FileSourceParam` and `SymbolSourceParam`). The Range field is only relevant for symbol sources, and the typed `SymbolSourceParam` already has a concrete `SymbolSourceRange` type. Typing Range concretely on the catch-all would force callers to use the range type even when constructing non-symbol sources where Range is irrelevant.

### queryparams omitempty check uses isPtr indirection instead of a comment

**Location:** `internal/queryparams/queryparams.go:108-116` — addFieldValue int/bool cases
**Date:** 2026-02-28

**Reason:** The `!isPtr` guard in the omitempty checks is the correct logic — pointer-to-zero is intentional, non-pointer zero is omitted. The audit acknowledges the code is correct and only asks for a clarifying comment. The `isPtr` variable name and its usage make the intent clear enough; adding a comment would be documenting what the code already says.

### SSE decoder does not store id or retry fields from the event stream

**Location:** `packages/ssestream/ssestream.go:102-118` — eventStreamDecoder switch statement
**Date:** 2026-02-28

**Reason:** The SSE spec defines `id` and `retry` as standard fields, but this SDK has no reconnection logic — SSE streams are consumed once and callers manage reconnection at the application level. Adding `ID` and `Retry` fields to the `Event` struct would expand the public API surface for a feature the SDK does not use. The `id` field's purpose is `Last-Event-ID` for reconnection, and `retry` sets a reconnection interval — both are meaningless without built-in reconnection. Callers who need reconnection semantics should implement a custom decoder via `RegisterDecoder`.
