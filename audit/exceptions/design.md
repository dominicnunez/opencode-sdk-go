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

### readAPIError relies on net/http non-nil Header guarantee

**Location:** `errors.go:120` — resp.Header.Get without nil check
**Date:** 2026-03-01

**Reason:** `readAPIError` accesses `resp.Header.Get("X-Request-Id")` without a nil guard. Go's `net/http` guarantees non-nil headers on all responses returned by `http.Client.Do`. The function is only called from `Client.do` and `ListStreaming`, both of which receive responses from `http.Client.Do`. Custom transports that return nil headers violate the `net/http.RoundTripper` contract. Adding a defensive nil check would guard against a contract violation that indicates a broken transport, not an SDK bug.

### SSE decoder does not store id or retry fields from the event stream

**Location:** `packages/ssestream/ssestream.go:102-118` — eventStreamDecoder switch statement
**Date:** 2026-02-28

**Reason:** The SSE spec defines `id` and `retry` as standard fields, but this SDK has no reconnection logic — SSE streams are consumed once and callers manage reconnection at the application level. Adding `ID` and `Retry` fields to the `Event` struct would expand the public API surface for a feature the SDK does not use. The `id` field's purpose is `Last-Event-ID` for reconnection, and `retry` sets a reconnection interval — both are meaningless without built-in reconnection. Callers who need reconnection semantics should implement a custom decoder via `RegisterDecoder`.

### ErrInvalidRequest is a catch-all for 4xx without a dedicated sentinel

**Location:** `errors.go:56-73` — APIError.Is switch
**Date:** 2026-02-28

**Reason:** `Is()` evaluates top-down: 401, 403, 404, 408, and 429 match their dedicated sentinels first; the 400-499 range then catches everything else as `ErrInvalidRequest`. A caller doing `errors.Is(err, ErrInvalidRequest)` on a 429 gets `false` — this is correct because 429 has a more specific sentinel (`ErrRateLimited`). Renaming to `ErrClientError` would be a breaking public API change for clearer naming alone. The current sentinel names, combined with the `Is*Error()` helpers and thorough test coverage in `errors_test.go`, make the semantics unambiguous.

### AssistantMessageErrorAPIErrorData and SessionAPIErrorData are structurally identical

**Location:** `session.go:547-553`, `event.go:826-832` — two API error data structs
**Date:** 2026-03-01

**Reason:** Both types map to distinct OpenAPI spec schemas (`AssistantMessage.error.data` vs `Session.error.APIError.data`). They happen to be byte-for-byte identical today, but extracting a shared type would couple two independent spec schemas. If either schema adds or removes a field, the shared type would break the other. Keeping them separate preserves 1:1 spec fidelity at the cost of ~6 lines of duplication.

### ListStreaming uses stricter status check than doRaw

**Location:** `event.go:67` — `resp.StatusCode < 200 || resp.StatusCode >= 300`
**Date:** 2026-02-28

**Reason:** `doRaw` accepts any status < 400 as success because JSON API responses can legitimately use 2xx and 3xx codes. SSE streams require a 200 OK with a streaming body — a 204 No Content or 3xx redirect would produce an empty or missing event stream with no error signal. The stricter check ensures SSE callers get an explicit `*APIError` for any non-2xx response rather than a silent empty stream. The `do` path handles different HTTP semantics and the two checks are intentionally distinct.

### Event type-specific Type fields duplicate the parent Event discriminator

**Location:** `event.go:318,341,365` — all EventXxx structs
**Date:** 2026-02-28

**Reason:** Each `EventXxx` struct has its own `Type` field with a dedicated string type and `IsKnown()` method, despite the parent `Event.Type` already carrying this information. This is spec-driven — the OpenAPI spec defines a `type` property on each event schema. The per-struct types faithfully reflect the spec and ensure round-trip fidelity. Removing them would diverge from the spec and break callers who access `event.Type` on the concrete struct after calling `AsXxx()`.

### RegisterDecoder uses global mutable state without unregister

**Location:** `packages/ssestream/ssestream.go:52-64` — global decoder registry
**Date:** 2026-03-01

**Reason:** `RegisterDecoder` follows the same global-registration pattern as `sql.Register`, `image.RegisterFormat`, and `encoding.RegisterCodec` in the Go stdlib. Registrations are process-lifetime by design — decoders are registered at init time and never removed. Moving the registry to the `Client` struct would force callers to configure decoders per-client, which is unnecessary since content-type decoders are application-wide. The mutex prevents races, and the test helper `saveAndRestoreDecoders` adequately isolates test state.

### ListStreaming returns error via stream object on buildURL failure

**Location:** `event.go:39-41` — buildURL error path
**Date:** 2026-03-01

**Reason:** When `buildURL` fails, `ListStreaming` wraps the error into the stream via `ssestream.NewStream[Event](nil, err)`. The caller must check `stream.Err()` after iteration, which is documented in the method's godoc with a full usage example. This matches Go's iterator contract (`bufio.Scanner`, `sql.Rows`) where errors are deferred to an `Err()` method. Returning `(*Stream, error)` would break the single-return-value streaming API and force callers to handle two error paths instead of one.

### ListStreaming buildURL and request-creation error paths are not unit-tested

**Location:** `event.go:38-47` — two early-return error paths
**Date:** 2026-03-01

**Reason:** The `buildURL` failure path (line 39-41) is unreachable through the public API because `EventListParams.URLQuery()` delegates to `queryparams.Marshal` which cannot error for the field types in `EventListParams` (a single `*string`). The `http.NewRequestWithContext` failure path (line 45-47) requires a nil context or invalid method, neither of which can occur from normal usage — the method is hardcoded as `GET` and context comes from the caller. Both paths exist as defensive coding against future changes to the params struct or internal invariants. Testing them would require either bypassing the type system or injecting programming errors, providing no regression value.

### Backoff overflow guard is unreachable with current constants

**Location:** `client.go:283` — `delay <= 0` check in retry backoff
**Date:** 2026-03-01

**Reason:** With `maxRetryCap = 10` and `initialBackoff = 500ms`, the maximum product is `500ms * 1024 = 512s`, well within `time.Duration` (int64 nanoseconds) range, so the `delay <= 0` overflow check is unreachable. However, it's a zero-cost guard that protects against future changes to `maxRetryCap` or `initialBackoff`. The `maxBackoff` cap alone would mask an overflow (capping a negative duration to `maxBackoff` silently), so the explicit overflow check adds a meaningful safety net. Removing it saves nothing and introduces a latent risk.

### APIError.Is returns false for 3xx status codes

**Location:** `errors.go:66` — `Is()` method switch statement
**Date:** 2026-03-01

**Reason:** 3xx responses are not standard API error classes — they represent redirects, not client or server errors. The `Is()` method intentionally maps only 4xx and 5xx to sentinels. 3xx errors are created as `*APIError` by `doRaw` (allowing `errors.As` + `StatusCode` inspection), but have no sentinel because there is no actionable error category for redirects that callers would match against. Adding `ErrRedirect` would expand the public API for a status class that HTTP clients typically handle transparently. Callers who need to detect 3xx can use `errors.As` and check `StatusCode`, which is the correct pattern for uncommon status classes.

### PermissionTime.Created uses float64 for Unix timestamp

**Location:** `sessionpermission.go:64-66` — also `ProjectTime.Created`, `SessionTime`
**Date:** 2026-03-01

**Reason:** The OpenAPI spec defines these timestamps as JSON `number` (IEEE 754 double). Go's `float64` is the correct mapping. Sub-millisecond precision loss is inherent to the wire format, not the SDK. Changing the Go type to `int64` or `time.Time` would diverge from the spec and break round-trip fidelity for callers who re-serialize these values.

### SSE decoder silently drops field names with leading/trailing whitespace

**Location:** `packages/ssestream/ssestream.go:122` — switch on field name
**Date:** 2026-03-01

**Reason:** Per the SSE spec (W3C Server-Sent Events §9.2), field names are matched literally — the spec does not define trimming. A line like `data :value` produces field name `"data "` which does not match `"data"`, so the field is ignored. This is correct spec-compliant behavior. The SSE spec explicitly states that unknown field names must be ignored, and whitespace is significant in field names.

### ConfigProviderOptionsTimeoutUnion.AsInt accepts negative timeout values

**Location:** `config.go:1288-1289` — allows `-` as first character
**Date:** 2026-03-01

**Reason:** This is a deserialization type that faithfully represents wire data. A negative number is a valid JSON number, and the SDK should not reject it at the unmarshal layer. Validation of timeout semantics (positive, within range) belongs at the application layer, not in the SDK's type system. Silently clamping or rejecting values would violate the principle of faithful wire representation.

### doRaw discards response body silently when result parameter is nil

**Location:** `client.go:167-169` — `do()` drains body to `io.Discard`
**Date:** 2026-03-01

**Reason:** When `result == nil`, `do()` drains the response body and returns nil error. This is the correct behavior for methods that return only a success/failure signal (e.g. `Session.Delete`). The caller opted out of body parsing by passing nil. Logging or returning a warning for non-empty bodies would add noise for a scenario that is not an error — the server is free to return metadata in 200 responses that the caller doesn't need. Callers who want the body should pass a result pointer.

### Config union types serialize as null when zero-valued

**Location:** `config.go` — all union types with `MarshalJSON` returning `[]byte("null")` for `raw == nil`
**Date:** 2026-03-01

**Reason:** Union types store `json.RawMessage` internally and return `"null"` when never populated via `UnmarshalJSON`. For `ConfigUpdateParams`, this means zero-value union fields inside nested structs appear as `null` in the PATCH body. However, the nested structs themselves (e.g. `ConfigAgent`, `ConfigPermission`) also serialize their zero-value non-union fields (empty strings, false bools, zero ints). The `null` from unions is consistent with this behavior — Go's `omitempty` does not omit zero-value structs. Fixing only unions while leaving other zero-valued fields would create an inconsistency. The correct fix would require the server to use JSON Merge Patch (RFC 7396) semantics or the SDK to implement a sparse serializer that tracks which fields were explicitly set — both are disproportionate to the risk, which only materializes if the server interprets `null` as "unset" rather than "unchanged".

### Auth credential types expose sensitive fields via json.Marshal

**Location:** `config.go:1609-1661` — OAuth, ApiAuth, WellKnownAuth types
**Date:** 2026-03-01

**Reason:** These types must serialize credential fields (Access, Refresh, Key, Token) in the HTTP request body sent to the server — that is their purpose. `String()` and `GoString()` redact these fields to prevent accidental exposure via `fmt` or log output, which covers the common leak vector. Implementing `MarshalJSON` to redact would break the API since `Client.do()` uses `json.Marshal` to build request bodies. A context-aware `MarshalJSON` that switches between redacted and non-redacted modes would add complexity for a marginal benefit — callers who marshal these types to JSON are explicitly opting into serialization and should be aware of the contents. The type-level godoc already documents which fields are sensitive and that `String()` is the safe output method.

### queryparams emits non-pointer zero int/bool without omitempty

**Location:** `internal/queryparams/queryparams.go:147-170` — zero-value emission
**Date:** 2026-03-01

**Reason:** Non-pointer `int`/`bool` fields without `omitempty` are emitted at zero value. This is documented as intentional at line 149 and matches `encoding/json` semantics: zero without `omitempty` means "explicitly set to zero". No current SDK param struct triggers this — all int/bool query fields use either pointers or `omitempty`. Changing the default would break the semantic contract for any future struct that intentionally sends zero.

### ListStreaming ignores mime.ParseMediaType error

**Location:** `event.go:63` — also `ssestream.go:47`
**Date:** 2026-03-01

**Reason:** When `mime.ParseMediaType` errors on a malformed Content-Type, `mediaType` is empty, causing the code to fall through to the default SSE decoder. This is the correct fallback for streaming endpoints — SSE is the expected wire format and rejecting the stream for a malformed header would be overly strict. The same pattern in `ssestream.NewDecoder` ensures consistent behavior. Treating parse errors as "no explicit content type" and defaulting to SSE is a deliberate defensive choice.

### ConfigProviderOptions APIKey field exposed via json.Marshal

**Location:** `config.go:1264-1278` — ConfigProviderOptions type
**Date:** 2026-03-01

**Reason:** `ConfigProviderOptions` faithfully reflects the OpenAPI spec schema. The server includes API keys in config responses, and the SDK must deserialize them. Adding `json:"-"` would silently drop data the server returns, breaking callers who need the value. `String()` and `GoString()` already redact the APIKey for safe logging. The serialization concern belongs in the caller's logging/caching layer, not in the SDK's data types. Implementing a `Redact()` method would expand the public API surface for a concern that is the caller's responsibility.
