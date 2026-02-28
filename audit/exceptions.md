# Audit Exceptions

> Items validated as false positives or accepted as won't-fix.
> Managed by willie audit loop. Do not edit format manually.
>
> Entry format:
> ### Plain language description
> **Location:** `file/path:line` — optional context
> **Date:** YYYY-MM-DD
> **Reason:** Explanation (can be multiple lines)

## False Positives

<!-- Findings where the audit misread the code or described behavior that doesn't occur -->

### Retry loop body is not re-readable after first attempt

**Location:** `client.go:180-244` — retry loop body re-encoding
**Date:** 2026-02-25

**Reason:** The audit claims the final retry attempt (`attempt == c.maxRetries`) reuses a stale body because the re-encode block is guarded by `attempt < c.maxRetries`. This is incorrect. For HTTP error responses (status >= 400), when `attempt >= c.maxRetries` the code returns immediately at line 213-217 and never loops back. For transport errors, when `attempt == c.maxRetries` the re-encode block is skipped, but the loop counter increments to `attempt = c.maxRetries + 1` which fails the loop condition `attempt <= c.maxRetries`, so no further request is made with a stale body. In both cases, the body is never reused without re-encoding.

### Retry on transport error reuses exhausted body reader

**Location:** `client.go:180-244` — retry loop transport error path
**Date:** 2026-02-25

**Reason:** Same root misunderstanding as the previous finding. The audit claims the final attempt uses a stale body reader, but tracing the control flow shows: on the last iteration where `attempt == c.maxRetries`, the `attempt < c.maxRetries` guard prevents re-encoding, but also prevents delay. The loop then increments `attempt` past `c.maxRetries` and exits. No request is ever issued with an exhausted body.

### ConfigProviderOptionsTimeoutUnion.AsInt accepts boolean JSON values

**Location:** `config.go:1236-1241` — AsInt union accessor
**Date:** 2026-02-25

**Reason:** The audit claims `json.Unmarshal` into `int64` succeeds for JSON booleans (`true` → `1`, `false` → `0`). This is factually wrong. Go's `encoding/json` returns an error: "cannot unmarshal bool into Go value of type int64". Verified empirically. `AsInt()` correctly returns `(0, false)` for boolean input, and `AsBool()` correctly returns `(false, false)` for numeric input. The union discriminates types correctly.

### Inconsistent error message format for missing required parameters

**Location:** `session.go:54-55,74-75,86-87,98-99` and other service files
**Date:** 2026-02-22

**Reason:** The audit claims error messages are "inconsistent with the pattern used elsewhere." However, all parameter validation messages in the codebase follow the same format: `missing required X parameter`. The audit provides no evidence of actual inconsistency and cannot cite any examples of different phrasing because none exist.

### SSE buffer size integer overflow claim

**Location:** `packages/ssestream/ssestream.go:45`
**Date:** 2026-02-22

**Reason:** The audit claims `bufio.MaxScanTokenSize<<sseBufferMultiplier` (64KB << 9 = ~32MB) "could theoretically overflow on 32-bit systems." This is mathematically incorrect. The result is 33,554,432 bytes (~32MB), which is well under the 32-bit signed int maximum of 2,147,483,647 (~2.1GB). No overflow is possible.

### Backoff bit-shift overflow on high retry counts

**Location:** `client.go:246` — exponential backoff calculation
**Date:** 2026-02-26

**Reason:** The audit flags `initialBackoffMs*(1<<attempt)` as fragile if maxRetries were raised above ~22. However, `WithMaxRetries` hard-caps at 10, and the `maxBackoff` clamp catches any large value regardless. The overflow scenario requires violating an enforced invariant. The report itself concludes "No action needed."

### APIError.Is() does not match other APIError instances via errors.Is()

**Location:** `errors.go:31-45` — APIError.Is implementation
**Date:** 2026-02-28

**Reason:** The `Is()` implementation is correct and idiomatic. Its purpose is to map HTTP status codes to sentinel errors (`ErrNotFound`, `ErrRateLimited`, etc.) so callers can write `errors.Is(err, ErrNotFound)`. The audit's claim that `Is` should also match other `*APIError` values by status code is not standard Go practice — that would make two distinct error instances with the same status code semantically "equal", which is confusing. Callers who need type matching use `errors.As`. Go's `errors.Is` already handles pointer equality before calling the custom `Is` method, so identical pointers work. The implementation correctly supplements the default behavior.

### Event.UnmarshalJSON silently accepts unknown event types

**Location:** `event.go:77-88` — UnmarshalJSON discriminator handling
**Date:** 2026-02-28

**Reason:** The audit itself concludes "No code change needed." Accepting unknown event types is intentional forward-compatibility — the SDK must handle server-side additions gracefully without breaking existing callers. Callers can check `e.Type.IsKnown()` or use a default case in type switches. All other union types in the codebase behave identically. This is standard practice for SDKs that consume versioned APIs.

### Ptr helper function appears unused

**Location:** `ptr.go:4-6` — generic Ptr[T] function
**Date:** 2026-02-28

**Reason:** The audit claims `Ptr[T]` is unused, but it's exercised in `readme_test.go` (lines 312, 318, 324). More importantly, it's a public API convenience helper intended for SDK consumers who need pointer values for optional fields (e.g., `opencode.Ptr("value")` for `*string` params). Internal non-usage is expected — the SDK itself doesn't need the helper because it constructs structs directly.

### SSE response body not closed when NewDecoder returns nil

**Location:** `event.go:51-70` — ListStreaming success path
**Date:** 2026-02-26

**Reason:** The audit claims that if `NewDecoder` returns nil (because `res` or `res.Body` is nil at ssestream.go:31), the response body leaks. This scenario cannot occur in the `ListStreaming` code path. At event.go:51, `httpClient.Do(req)` returns successfully (no error), which guarantees both `resp` and `resp.Body` are non-nil per Go's `net/http` contract. `NewDecoder` will always receive a valid response and return a non-nil decoder. The general concern about callers needing to call `stream.Close()` is standard Go resource management (like `os.File.Close()`), not a code bug.

### ToolStateRunning.Input typed as interface{} instead of map[string]interface{}

**Location:** `session.go:1530` — ToolStateRunning.Input field type
**Date:** 2026-02-28

**Reason:** The audit claims this is an inconsistency with ToolStateCompleted.Input and ToolStateError.Input which use `map[string]interface{}`. However, the OpenAPI spec defines ToolStateRunning's `input` with an empty schema (`"input": {}`), while ToolStateCompleted and ToolStateError define it as `"input": {"type": "object", "propertyNames": {"type": "string"}, "additionalProperties": {}}`. The code correctly reflects this spec difference: unconstrained schema maps to `interface{}`, typed object schema maps to `map[string]interface{}`. The types are intentionally different per the spec.

### Session.Update error message format is inconsistent with other methods

**Location:** `session.go:34` — missing required id parameter error
**Date:** 2026-02-28

**Reason:** The audit compares `Session.Update`'s error format against a `missing required parameter 'X' (received empty string)` format documented in exceptions.md. However, that format doesn't exist anywhere in the codebase. All 30+ parameter validation messages across session.go, sessionpermission.go, and auth.go use the identical format: `missing required id parameter` (or `missing required permissionID parameter`, `missing required messageID parameter`). The code is internally consistent. The exceptions.md entry describing the `'X' (received empty string)` format is itself stale.

### apierror.Error stores live http.Request and http.Response references

**Location:** `internal/apierror/apierror.go:12-17` — Error struct fields
**Date:** 2026-02-28

**Reason:** The finding itself acknowledges this is already documented as won't-fix in exceptions.md — `apierror.Error` is never constructed anywhere in the SDK (it's a Stainless leftover exposed as `opencode.Error`). Since the type is inert, the references can never pin memory in practice. This is a duplicate of the existing "apierror.Error is unused but exported as a public type alias" won't-fix entry.

### ConfigProviderListResponse described as a map type alias but is actually a struct

**Location:** `config.go:1657-1660` — ConfigProviderListResponse type definition
**Date:** 2026-02-28

**Reason:** The audit claims `ConfigProviderListResponse` is "a `map[string]Provider`" and flags a naming mismatch between the type name suggesting a wrapper and it being "a type alias for a map." This is factually wrong. `ConfigProviderListResponse` is defined as a struct with `Default map[string]string` and `Providers []ConfigProvider` fields — it is a proper response wrapper struct, not a map alias. The finding's premise is entirely incorrect.

### ConfigProviderModelsLimit uses float64 instead of int64 for token limits

**Location:** `config.go:1213-1216` — Context and Output fields
**Date:** 2026-02-28

**Reason:** The OpenAPI spec defines both `limit.context` and `limit.output` as `"type": "number"`, not `"type": "integer"`. `float64` is the correct Go mapping for JSON `number`. The `ModelLimit` type in `app.go` using `int64` is the one with the looser spec mapping, not `ConfigProviderModelsLimit`. The two types represent different spec schemas (one under `Config.Provider.Models`, one under `Model`) and their field types correctly reflect the spec definitions.

### CI workflow claims scripts are missing but they exist

**Location:** `.github/workflows/ci.yml:31,46,49` — references to `./scripts/lint`, `./scripts/bootstrap`, `./scripts/test`
**Date:** 2026-02-28

**Reason:** The audit claims these scripts don't exist, stating "only `scripts/check-spec-update.sh` exists in the repository." This is factually wrong. All three scripts (`scripts/lint`, `scripts/bootstrap`, `scripts/test`) exist in the repository alongside `scripts/check-spec-update.sh`, `scripts/format`, and `scripts/mock`.

### Duplicate error struct types for APIError data across session types

**Location:** `session.go:542-553, session.go:1028-1039` — APIError data structs
**Date:** 2026-02-28

**Reason:** The audit claims three structurally identical types exist: `AssistantMessageErrorAPIErrorData`, `SessionAPIErrorData`, and `PartRetryPartErrorData`. However, `SessionAPIErrorData` does not exist anywhere in the codebase (line 830-841 contains `MessageRole` and the `Part` union type, not an error data struct). The audit also claims an `omitempty` tag difference between the types, but the two types that do exist (`AssistantMessageErrorAPIErrorData` at line 547 and `PartRetryPartErrorData` at line 1033) have identical struct tags. The finding's cited locations, type count, and tag variance claim are all factually wrong.

### Retries-exhausted on HTTP error returns untyped string instead of structured error

**Location:** `client.go:283` — retry loop exit path
**Date:** 2026-02-28

**Reason:** The audit claims that when all retries are exhausted on retryable HTTP errors (408, 429, 5xx),
the code falls through to line 283 and returns `errors.New("request failed: retries exhausted")`. This is
factually wrong. Tracing the control flow: on the final iteration (`attempt == maxRetries`) with a retryable
HTTP status, the condition `!shouldRetry || attempt >= c.maxRetries` at line 233 evaluates to true, so the
code enters the block and returns a structured `*APIError` at line 242-247 with status code, message, request
ID, and body. Line 283 is dead code — it can only be reached if the loop completes without either a transport
error (`lastErr != nil` returns at line 280) or an HTTP error (returns at line 242). In practice, every
iteration produces one of those two outcomes.

### CI fork filter described as inverted but uses the standard anti-duplication pattern

**Location:** `.github/workflows/ci.yml:11,27` — job-level `if` condition
**Date:** 2026-02-28

**Reason:** The condition `github.event_name == 'push' || github.event.pull_request.head.repo.fork` is the standard
pattern to avoid duplicate CI runs. Same-repo PRs already receive CI from the `push` event (which fires for
every branch push). Fork PRs don't trigger `push` events, so the fork check ensures they still get CI via
the `pull_request` event. The audit's suggested "fix" (`!github.event.pull_request.head.repo.fork`) would
actually cause same-repo PRs to run CI twice — once from push, once from pull_request.

## Won't Fix

<!-- Real findings not worth fixing — architectural cost, external constraints, etc. -->

### apierror.Error is unused but exported as a public type alias

**Location:** `internal/apierror/apierror.go:12-17` — aliased as `opencode.Error` in `aliases.go:8`
**Date:** 2026-02-28

**Reason:** `apierror.Error` is never constructed anywhere in the SDK — it's a Stainless leftover. However, it's exposed as the public type `opencode.Error`. Removing it would be a breaking API change for any caller that references the type. The type is inert (never returned by any SDK method), so it causes no runtime harm.


### Path parameters not URL-encoded in service methods

**Location:** `session.go`, `sessionpermission.go` — string concatenation for path segments
**Date:** 2026-02-25

**Reason:** Path parameters (session IDs, permission IDs, message IDs) are constructed via string concatenation without `url.PathEscape()`. The IDs in this SDK are server-generated UUIDs that do not contain special characters, so path injection is not a practical concern for normal usage. Adding escaping to every path construction would add noise for no real-world benefit.

### httputil dump errors ignored in debugging methods

**Location:** `internal/apierror/apierror.go:44,46,51`
**Date:** 2026-02-22

**Reason:** The DumpRequest and DumpResponse methods are debugging utilities that return []byte. Adding error return values would be a breaking API change. For debugging purposes, returning empty output when the dump fails is acceptable behavior—the caller can inspect the returned bytes to determine if useful information was captured. Adding logging in library code is not idiomatic Go.

### Deprecated config fields still parsed

**Location:** `config.go:53-54,68,73-74`
**Date:** 2026-02-22

**Reason:** The Config struct fields (autoshare, mode, layout) reflect the upstream OpenAPI spec. The spec defines these fields as deprecated. Removing them would break deserialization of API responses that still include them. The deprecation comments are accurate and guide users to migrate.

### Default base URL uses plaintext HTTP

**Location:** `client.go:19` — DefaultBaseURL constant
**Date:** 2026-02-26

**Reason:** The SDK targets a local dev server (`localhost:54321`). The `WithBaseURL` validator intentionally allows `http://` for localhost usage. Adding hostname validation to reject non-localhost HTTP would be disproportionate: callers who set a remote `OPENCODE_BASE_URL` are explicitly overriding the default and responsible for their transport security. The SDK doesn't handle auth credentials itself — `AuthService.Set` is a passthrough to the server API.

### Bytes buffer allocation in SSE hot path

**Location:** `packages/ssestream/ssestream.go:81`
**Date:** 2026-02-22

**Reason:** The `bytes.NewBuffer(nil)` call per event is a minor allocation in a streaming context. For typical usage patterns, the GC overhead is negligible. Using `sync.Pool` would add complexity for an optimization that would only benefit extremely high-throughput scenarios. No performance issue has been reported or measured.

### POST methods with query-only params send an empty JSON body

**Location:** `tui.go` — TuiClearPromptParams, TuiOpenHelpParams, etc.
**Date:** 2026-02-28

**Reason:** When a params struct has only `query:` tagged fields and no `json:` fields, `doRaw` serializes it as `{}`. The server accepts empty bodies, and the overhead is a few bytes per request. Eliminating this would require splitting the query/body encoding path in `doRaw` or special-casing these methods, which adds complexity for no behavioral benefit.

## Intentional Design Decisions

<!-- Findings that describe behavior which is correct by design -->

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
