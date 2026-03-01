# Misreads

> Findings where the audit misread the code or described behavior that doesn't occur.
> Managed by sfk willie. Follow the entry format below.
>
> Entry format:
> ### Plain language description
> **Location:** `file/path:line` — optional context
> **Date:** YYYY-MM-DD
> **Reason:** Explanation (can be multiple lines)

### check-spec-update.sh silently treats empty hash or URL as valid

**Location:** `scripts/check-spec-update.sh:19-30` — hash/URL parsing from .stats.yml
**Date:** 2026-02-28

**Reason:** The script uses `set -euo pipefail` (line 6). If `grep 'openapi_spec_hash'` finds no match, the pipeline exits non-zero, and `pipefail` propagates that failure to the command substitution, which causes the script to abort under `set -e`. The same applies to the URL parsing on line 20. The scenario where `UPSTREAM_HASH` or `UPSTREAM_URL` is set to empty string without the script aborting cannot occur — `grep` returning no matches is a pipeline failure, not a successful empty result. The audit correctly noted `set -euo pipefail` is in effect but then incorrectly concluded the empty-variable scenario was still reachable.

### RegisterDecoder race between parse and write

**Location:** `packages/ssestream/ssestream.go:59-67` — RegisterDecoder locking
**Date:** 2026-02-28

**Reason:** The audit claims a race exists because `mime.ParseMediaType` and `strings.ToLower` execute before the write lock is acquired. These are pure functions operating on the function's input parameter — they access no shared state. The shared state (`decoderTypes` map) is properly protected by the mutex at line 64. The "last writer wins" behavior described by the audit is identical regardless of whether the parse runs inside or outside the lock, since both goroutines would compute the same `mediaType` from the same input. There is no data race.

### Reusing bytes.Buffer across retry iterations is fragile

**Location:** `client.go:192-284` — retry loop body encoding
**Date:** 2026-02-28

**Reason:** The audit describes the buffer reuse pattern as "fragile" but acknowledges the code is correct today: the re-encode block at line 278-284 creates a fresh buffer on every retry, guarded by the same method check as the initial encode at line 187. The described "bug" is purely hypothetical — "any future refactor that adjusts the re-encode guard independently would silently send empty bodies." The current code has no bug; the two guards are structurally identical and correct. This is speculative fragility, not a real defect.

### McpRemoteConfig.Headers exposes auth tokens in plain struct fields

**Location:** `config.go:1491` — Headers field
**Date:** 2026-02-28

**Reason:** This is a duplicate of the existing intentional design decision "ConfigProviderOptions exposes API key as a plain string field." The audit itself acknowledges "the same rationale from the APIKey exception applies." The struct reflects the OpenAPI spec schema, and redaction is a caller concern. Already categorized.

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

### RegisterDecoder silently accepts empty content type, creating unreachable decoder

**Location:** `packages/ssestream/ssestream.go:55-58` — RegisterDecoder with empty string
**Date:** 2026-02-28

**Reason:** The audit claims a decoder registered under the empty-string key "can never match because `mime.ParseMediaType("")` returns an error, making `mediaType` empty." This is factually wrong. `mime.ParseMediaType("")` does return an error, but the error is discarded at line 36 (`mediaType, _, _ := mime.ParseMediaType(...)`), and the returned `mediaType` is the empty string `""`. This means `decoderTypes[""]` WOULD match when a response has an empty or missing Content-Type header. The decoder is reachable, not unreachable. The finding's premise — that this is dead/unreachable registration — is incorrect.

### Timer not stopped on normal completion in retry backoff

**Location:** `client.go:279-285` — backoff timer select
**Date:** 2026-02-28

**Reason:** The finding acknowledges the code is correct: after a timer fires via `<-timer.C`, calling `Stop()` is a documented no-op in Go. The timer's internal goroutine exits upon firing. The finding's own suggested fix is "No action needed." The inconsistency with the `ctx.Done()` branch is stylistic, not a bug.

### Non-pointer string fields always omitted when empty in queryparams

**Location:** `internal/queryparams/queryparams.go:124-131` — addFieldValue string case
**Date:** 2026-02-28

**Reason:** The behavior is explicitly documented in the `Marshal` doc comment (lines 18-21): "Empty non-pointer strings are always omitted from the output, regardless of whether 'omitempty' is set." The `omitempty` tag being redundant on string fields is intentional — the finding describes documented, working-as-designed behavior and concludes "No code change needed."

### Session.Delete returns bool but may fail on empty response body

**Location:** `session.go:74-80` — Delete method
**Date:** 2026-02-28

**Reason:** The OpenAPI spec (`specs/openapi.yml:449-455`) defines the DELETE session endpoint as returning HTTP 200 with `application/json` body of type `boolean`. The server returns JSON `true`/`false`, not 204 No Content. The `json.Decoder.Decode` call correctly deserializes this into a `bool`. The finding's concern about 204 No Content is based on a general assumption about DELETE endpoints, not the actual spec.

### Redundant SessionPromptParamsPart catch-all type duplicates union interface

**Location:** `session.go:1788-1813` — SessionPromptParamsPart
**Date:** 2026-02-28

**Reason:** The finding itself concludes "This is intentional per the exceptions doc. No action needed." The catch-all type is a documented escape hatch for SDK consumers who need to construct parts without importing every nested type. The typed variants provide compile-time safety for callers who want it.

### queryparams.Marshal accepts non-struct values after pointer dereference

**Location:** `internal/queryparams/queryparams.go:28-34` — Marshal type check
**Date:** 2026-02-28

**Reason:** The finding describes behavior with inputs that no call site uses (`**SomeStruct`, `*string`). All call sites pass structs or `*struct`. The error message for unsupported types is correct ("expected struct, got X"). The finding itself concludes "No action needed — all call sites pass structs or `*struct`. The error message is adequate." This is a theoretical concern with no practical impact.

### CI fork filter described as inverted but uses the standard anti-duplication pattern

**Location:** `.github/workflows/ci.yml:11,27` — job-level `if` condition
**Date:** 2026-02-28

**Reason:** The condition `github.event_name == 'push' || github.event.pull_request.head.repo.fork` is the standard
pattern to avoid duplicate CI runs. Same-repo PRs already receive CI from the `push` event (which fires for
every branch push). Fork PRs don't trigger `push` events, so the fork check ensures they still get CI via
the `pull_request` event. The audit's suggested "fix" (`!github.event.pull_request.head.repo.fork`) would
actually cause same-repo PRs to run CI twice — once from push, once from pull_request.

### doRaw has no timeout guard, allowing future non-streaming callers to lack timeout protection

**Location:** `client.go:149,169` — `do` vs `doRaw` timeout inconsistency
**Date:** 2026-02-28

**Reason:** The finding claims `ListStreaming` calls `doRaw` directly, creating an inconsistency where `doRaw` callers have no timeout. This is factually wrong — `ListStreaming` (event.go:44) uses `s.client.httpClient.Do(req)` directly, not `doRaw`. `doRaw` is only ever called by `do` (client.go:152), which always wraps the context with `WithTimeout`. The concern about future callers of `doRaw` lacking timeout is speculative and already covered by the "ListStreaming bypasses Client timeout and retry logic" intentional design entry.

### apierror.Error has overlapping StatusCode field that is never read

**Location:** `internal/apierror/apierror.go:13-16` — redundant StatusCode field
**Date:** 2026-02-28

**Reason:** The finding correctly identifies that `Error()` reads `Response.StatusCode` instead of the `StatusCode` field, but the finding itself concludes "no action needed" since the type is a Stainless leftover. This is already tracked in the Won't Fix section as "apierror.Error is unused but exported as a public type alias" — the type is never constructed anywhere in the SDK, making the field overlap entirely theoretical.
