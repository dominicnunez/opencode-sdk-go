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

### Base URL query merge in buildURL is redundant because ResolveReference preserves query params

**Location:** `client.go:181-183` — buildURL base URL query loop
**Date:** 2026-02-28

**Reason:** The audit claims `url.URL.ResolveReference` preserves the base URL's query string when the reference has no `RawQuery`, making the loop at lines 181-183 redundant. This is factually wrong. `ResolveReference(&url.URL{Path: path})` produces a resolved URL with an empty `RawQuery` — the base URL's query params are dropped because the reference has a non-empty path. The loop is necessary to re-merge those params. Verified empirically: resolving `http://localhost:54321/?foo=bar` with `&url.URL{Path: "sessions"}` produces `http://localhost:54321/sessions` with no query string.

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

### Retry loop reuses consumed body reader on first retry

**Location:** `client.go:206-222` — retry loop body encoding
**Date:** 2026-02-28

**Reason:** The finding's own analysis contradicts its title. It states "The current code happens to work because re-encode always runs before the next request" and traces the correct control flow: (1) make request with body, (2) check retry, (3) backoff, (4) re-encode body, (5) loop back to make request. The body is always freshly encoded at step 4 before it is used at step 5. The "consumed-buffer window" described — where the `body` variable temporarily holds an exhausted buffer — is not a bug because no code reads `body` during that window. The variable is reassigned before its next use. This is the fourth variant of the same retry-body misread (see existing entries for lines 192-284, 180-244).

### Nil-params tests only verify error is non-nil, not message content

**Location:** `nilparams_test.go:127-128` — nil-params test assertions
**Date:** 2026-02-28

**Reason:** The finding claims "several other services with required params use the same `params is required` message but their nil-params tests only verify the error is non-nil, not the message content." This is factually wrong. `nilparams_test.go` covers all the named services (App.Log, File.List, File.Read, Find.Files, Find.Symbols, Find.Text, Tui.AppendPrompt, Tui.ExecuteCommand, Tui.ShowToast) and asserts `strings.Contains(err.Error(), "params is required")` at line 127. Additional dedicated tests (`config_update_test.go:126`, `auth_test.go:122`, `session_shell_test.go:174`, `session_summarize_test.go:114`, `sessionpermission_respond_test.go:240`) also check the error message. The error message format is consistent across all services and tested everywhere.

### ListStreaming bypasses client timeout on SSE connections

**Location:** `event.go:51` — uses httpClient.Do directly
**Date:** 2026-02-28

**Reason:** SSE streams are long-lived connections that remain open indefinitely while events arrive. Applying the client's default 30s timeout would prematurely kill every SSE connection. Callers who need a deadline can set one via `context.WithTimeout` on the context they pass in. The audit correctly describes the code path but misclassifies it as a bug — this is intentional design consistent with how other Go SSE/WebSocket libraries handle timeouts.

### ListStreaming bypasses client retry logic

**Location:** `event.go:51` — single HTTP request via httpClient.Do
**Date:** 2026-02-28

**Reason:** SSE streams are consumed once; retrying a stream is complex and callers should manage reconnection at the application level. The audit correctly describes the behavior but misclassifies it as a bug. `ListStreaming` intentionally bypasses `Client.doRaw()` for SSE semantics. The `WithMaxRetries` configuration governs JSON API calls, not streaming connections.

### apierror.Error unused Stainless leftover combines already-excepted sub-issues

**Location:** `internal/apierror/apierror.go:12-17`, `aliases.go:8`
**Date:** 2026-02-28

**Reason:** The finding rolls up multiple concerns (memory pinning from stored `*http.Request`/`*http.Response`, dead `StatusCode` field, `DumpRequest` body mutation) that are each already classified as exceptions: "apierror.Error stores live http.Request and http.Response references", "apierror.Error has overlapping StatusCode field that is never read", "httputil dump errors ignored in debugging methods", and "apierror.Error is unused but exported as a public type alias". The type is never constructed anywhere in the SDK, making all sub-issues theoretical. No new observation beyond existing exceptions.

### TestClientDo_Success described as silently swallowing decode errors

**Location:** `client_do_test.go:27-48` — Session.Create success test
**Date:** 2026-03-01

**Reason:** The audit claims the server returns `{"message": "success", "count": 42}` and the decode error is only logged with `t.Logf`. This is factually wrong. The server encodes a valid `Session` struct (line 27-31), the test asserts `err != nil` via `t.Fatalf` (line 43-44), and checks `session.ID` (line 46-48). The test properly validates the full request/response cycle.

### TestClientDo_Retry described as silently swallowing decode errors

**Location:** `client_do_test.go:62-86` — Session.List retry test
**Date:** 2026-03-01

**Reason:** The audit claims the server returns `{"status": "ok"}` which is not valid `[]Session` and the decode error is logged with `t.Logf`. This is factually wrong. The server encodes a valid `[]Session` (line 62-64), the test asserts `err` via `t.Fatalf` (line 77-78), and verifies `sessions[0].ID` (line 84-85). The response is correctly decoded and validated.

### TestRetryAfterMs described as testing a non-existent Retry-After-Ms feature

**Location:** `client_test.go:52-79` — 429 retry test
**Date:** 2026-03-01

**Reason:** The audit claims the test is named `TestRetryAfterMs` and sends a `Retry-After-Ms: 100` header. This is factually wrong. The test is named `TestRetryOn429` (line 52), sends no `Retry-After-Ms` header, and correctly tests that 429 responses are retried. The audit described a test that doesn't exist in the codebase.

### Refactoring archaeology tests described as serving no ongoing purpose

**Location:** `cleanup_verification_test.go`, `deletion_verification_test.go`, `apiform_deletion_test.go`
**Date:** 2026-03-01

**Reason:** All three files have already been deleted. The git status shows them with `D` prefix (staged deletions). The finding describes files that no longer exist in the working tree.

### apiform_deletion_test.go described as containing an empty test body

**Location:** `apiform_deletion_test.go:85-88` — BuildStillWorks subtest
**Date:** 2026-03-01

**Reason:** The file has been deleted. It no longer exists in the working tree (git status shows `D apiform_deletion_test.go`).

### TestClientDo_QueryParams described as only asserting query presence, not content

**Location:** `client_do_test.go:135-140` — query params test
**Date:** 2026-03-01

**Reason:** The audit claims the test "only checks `receivedQuery == ""`" and that "any non-empty query passes." This is factually wrong. Line 138 asserts `strings.Contains(receivedQuery, "directory=%2Ftest")`, which is a specific value assertion on the encoded query parameter content.

### No test covers retry with POST body re-encoding

**Location:** `client.go:284-289` — retry loop body re-encoding
**Date:** 2026-03-01

**Reason:** The audit claims no test exercises the POST body re-encoding path on retry. This is factually wrong. `TestClientDo_PostBodyReencodedOnRetry` (client_do_test.go:303-361) POSTs to a server that returns 500 on the first attempt and 200 on the second, then verifies both attempts received identical non-empty request bodies containing `"test-parent"`.

### No test covers transport-error retry exhaustion

**Location:** `client.go:296` — retry exhaustion error path
**Date:** 2026-03-01

**Reason:** The audit claims no test exercises the transport-error retry exhaustion path. This is factually wrong. `TestClientDo_TransportErrorRetryExhaustion` (client_do_test.go:363-402) uses a custom transport that always returns `connection refused`, verifies 3 attempts are made, and asserts the error wraps the transport error and mentions the retry count.

### SSE decoder error propagation described as untested

**Location:** `packages/ssestream/ssestream.go:216` — Stream.Next error propagation
**Date:** 2026-03-01

**Reason:** The audit claims the test suite's `mockDecoder` always returns `nil` from `Err()` and no test covers decoder errors mid-stream. This is factually wrong. `TestStream_DecoderErrorPropagation` (ssestream_test.go:252-265) uses an `errorDecoder` that returns `false` from `Next()` with a non-nil `Err()` ("connection reset by peer"), and verifies `stream.Err()` surfaces the decoder error via `errors.Is`.

### RegisterDecoder lookup path described as untested

**Location:** `packages/ssestream/ssestream.go:41-51` — NewDecoder content-type lookup
**Date:** 2026-03-01

**Reason:** The audit claims tests cover concurrent registration but never verify lookup or fallback. This is factually wrong. `TestRegisterDecoder_LookupByContentType` (ssestream_test.go:277-299) registers a custom decoder, creates a response with matching content-type, and verifies the custom factory was called. `TestNewDecoder_UnknownContentType_FallsBackToSSE` (ssestream_test.go:301-321) verifies an unknown content-type falls back to the SSE decoder and correctly parses events.

### readme_test.go described as containing vacuous assertions

**Location:** `readme_test.go:162, readme_test.go:266, readme_test.go:289-305`
**Date:** 2026-03-01

**Reason:** The audit's claims about the test code are factually wrong at every cited location. (1) `StreamingEvents` (lines 162-197) calls `stream.Next()`, retrieves the event, and asserts `evt.Type == EventTypeMessageUpdated` — not just `stream != nil`. (2) `CustomHTTPClient` (lines 270-301) makes an actual HTTP request to a mock server and asserts the response is valid — not just `client != nil`. (3) `TestREADMELoggingTransport` (lines 305-325) asserts both `err != nil` and `resp == nil` after calling `RoundTrip`, then closes the body — not "no assertion." The audit described test behavior that doesn't match the actual code.

### Path parameter injection via unsanitized user input in URL construction

**Location:** `session.go:40,67,82,97`, `sessionpermission.go:29`, `auth.go:31` — path segment interpolation
**Date:** 2026-03-01

**Reason:** Duplicate of existing known exception "Path parameters not URL-encoded in service methods." The IDs are server-generated UUIDs that do not contain special characters. The exception already documents this as a conscious design tradeoff — adding escaping to every path construction would add noise for no real-world benefit.

### SSE stream response body leak on non-2xx without Close

**Location:** `event.go:56-58` — non-2xx status path
**Date:** 2026-03-01

**Reason:** Duplicate of existing known exception "ListStreaming response body not closed if readAPIError panics on custom transport." The `readAPIError` function reads and closes the body (errors.go:96). A leak only occurs if a custom transport panics during `io.ReadAll`, which violates Go's `io.Reader` contract. The normal code path has no leak.

### Retry loop sends empty body on retries for requests with io.Reader body

**Location:** `client.go:223-292` — retry loop body re-encoding
**Date:** 2026-03-01

**Reason:** The finding's own text admits "In practice this is harmless because context cancellation returns early." This is the fifth variant of the same retry-body concern, already excepted four times: "Reusing bytes.Buffer across retry iterations is fragile" (client.go:192-284), "Retry loop body is not re-readable after first attempt" (client.go:180-244), "Retry on transport error reuses exhausted body reader" (client.go:180-244), and "Retry loop reuses consumed body reader on first retry" (client.go:206-222). The code re-encodes the body at lines 284-289 before every retry iteration. No request is ever made with a stale body.

### buildURL duplicates base URL query parameters

**Location:** `client.go:180-200` — buildURL base URL query loop
**Date:** 2026-03-01

**Reason:** The finding claims `ResolveReference` preserves the base URL's query string, making the loop at lines 183-185 redundant. This is factually wrong — already proven in existing exception "Base URL query merge in buildURL is redundant because ResolveReference preserves query params." `ResolveReference(&url.URL{Path: path})` drops the base URL's `RawQuery` because the reference has a non-empty path. The loop is necessary to re-merge base URL query parameters.

### apierror.Error retains full http.Request and http.Response

**Location:** `internal/apierror/apierror.go:12-17` — Error struct fields
**Date:** 2026-03-01

**Reason:** Duplicate of existing known exceptions: "apierror.Error stores live http.Request and http.Response references" and "apierror.Error unused Stainless leftover combines already-excepted sub-issues." The type is never constructed anywhere in the SDK — it's a Stainless leftover exposed as `opencode.Error`. Since the type is inert, the references can never pin memory in practice.

### apierror.Error appears to be dead code (Stainless artifact)

**Location:** `internal/apierror/apierror.go:1-60`, `aliases.go:8`
**Date:** 2026-03-01

**Reason:** Duplicate of existing known exception "apierror.Error is unused but exported as a public type alias." Removing it would be a breaking API change for any caller that references `opencode.Error`. The type is inert (never returned by any SDK method) so it causes no runtime harm.

### McpStatus typed as map[string]interface{} loses all type safety

**Location:** `mcp.go:17`
**Date:** 2026-03-01

**Reason:** Duplicate of existing known exception "McpStatus is an untyped map." The OpenAPI spec defines the MCP status response with an empty schema (`"schema": {}`), meaning the response shape is intentionally unspecified. `map[string]interface{}` is the correct Go representation of an unconstrained JSON object.

### FilePartSourceParam.Range uses any type

**Location:** `session.go:717` — Range field
**Date:** 2026-03-01

**Reason:** Duplicate of existing known exception "FilePartSourceParam.Range typed as `any`." `FilePartSourceParam` is the catch-all variant of `FilePartSourceUnionParam`. The typed `SymbolSourceParam` already has a concrete `SymbolSourceRange` type for callers who want type safety.

### Event.ListStreaming does not apply client timeout

**Location:** `event.go:31-61` — uses httpClient.Do directly
**Date:** 2026-03-01

**Reason:** Duplicate of existing known exceptions "ListStreaming bypasses Client timeout and retry logic" and "ListStreaming bypasses client timeout on SSE connections." SSE streams are long-lived connections; applying a 30s timeout would kill every connection. Callers use `context.WithTimeout` for deadlines.

### EventService.ListStreaming does not use retry logic

**Location:** `event.go:31-61` — single HTTP request via httpClient.Do
**Date:** 2026-03-01

**Reason:** Duplicate of existing known exception "ListStreaming bypasses client retry logic." SSE streams are consumed once; retrying is complex and callers should manage reconnection at the application level.

### apierror.Error credential leakage via DumpRequest

**Location:** `internal/apierror/apierror.go:12-17`, `aliases.go:8`
**Date:** 2026-03-01

**Reason:** Duplicate of multiple existing known exceptions covering the apierror.Error Stainless leftover. The type is never constructed anywhere in the SDK — no SDK method returns it. The credential leakage scenario requires constructing the type with a live `*http.Request` containing auth headers, which only a consumer could do (and they already have the request). The type is inert in practice.

### SSE decoder silently discards id and retry fields

**Location:** `packages/ssestream/ssestream.go:113-130` — switch statement
**Date:** 2026-03-01

**Reason:** Duplicate of existing known exception "SSE decoder does not store id or retry fields from the event stream." The SDK has no reconnection logic — SSE streams are consumed once. The `id` field's purpose is `Last-Event-ID` for reconnection, and `retry` sets a reconnection interval — both are meaningless without built-in reconnection.

### apierror.Error StatusCode field shadowed by Response.StatusCode in Error()

**Location:** `internal/apierror/apierror.go:13, 30-32`
**Date:** 2026-03-01

**Reason:** Duplicate of existing known exception "apierror.Error has overlapping StatusCode field that is never read." The type is a Stainless leftover never constructed by the SDK, making the field overlap entirely theoretical.

### apierror.Error dead code consolidation duplicates existing exceptions

**Location:** `internal/apierror/apierror.go:12-60`, `aliases.go:8`
**Date:** 2026-03-01

**Reason:** The finding rolls up sub-issues (memory pinning, dead StatusCode field, DumpRequest mutation, dead code) that are each already classified in known exceptions: "apierror.Error is unused but exported as a public type alias", "apierror.Error stores live http.Request and http.Response references", "apierror.Error has overlapping StatusCode field that is never read", "httputil dump errors ignored in debugging methods", and "apierror.Error unused Stainless leftover combines already-excepted sub-issues." The finding itself concludes "No action needed — already tracked." No new observation beyond existing exceptions.

### ConfigService.Update described as having no valid-payload test coverage

**Location:** `config_update_test.go` — ConfigService.Update test suite
**Date:** 2026-03-01

**Reason:** The audit claims "config_update_test.go only tests the nil-params validation error." This is factually wrong. `TestConfigUpdate_Success` (line 12) sends a `ConfigUpdateParams` with model and theme fields to a mock server, verifies the HTTP method is PATCH, decodes the request body and asserts field values, and checks the response. `TestConfigUpdate_WithDirectoryQueryParam` (line 67) additionally verifies query parameter encoding. `TestConfigUpdate_ServerError` (line 132) tests error handling. `TestConfigUpdate_InvalidJSON` (line 163) tests malformed responses. `TestConfigUpdateParams_MarshalJSON` (line 187) tests serialization. The test file has comprehensive coverage — the audit apparently read a stale or different version of the file.

### ListStreaming body close suggestion duplicates known exception for readAPIError panic path

**Location:** `event.go:58-59` — non-2xx status path
**Date:** 2026-03-01

**Reason:** The finding suggests adding `defer resp.Body.Close()` before the status check to guard against `readAPIError` panicking. This is already tracked as a known exception ("ListStreaming response body not closed if readAPIError panics on custom transport") which documents that such a panic requires a custom transport violating Go's `io.Reader` contract, and that adding a defer-close would double-close the body on the normal path. The finding acknowledges the existing exception tracking.

### No Retry-After header parsing described as a testing gap

**Location:** `client.go:258-267` — retry backoff for 429 responses
**Date:** 2026-03-01

**Reason:** The SDK does not implement `Retry-After` header support — 429 responses use the same exponential backoff as 5xx. `TestRetryOn429` correctly tests the actual behavior (retry count). The finding describes a missing *feature* (Retry-After support) as a testing gap, but there is nothing to test when the feature doesn't exist. The finding itself says "While not a bug."

### No SSE reconnection test described as a testing gap

**Location:** `event.go:33-63` — ListStreaming single-request design
**Date:** 2026-03-01

**Reason:** The SDK intentionally does not implement SSE reconnection — this is documented in known exceptions ("ListStreaming bypasses client retry logic"). Suggesting an example test for caller-side reconnection is a documentation/example request, not a test gap. The code under test has no reconnection logic to exercise.

### Auth types defined in config.go described as not visible from auth.go

**Location:** `auth.go:40-83` — MarshalJSON referencing OAuth, ApiAuth, WellKnownAuth
**Date:** 2026-03-01

**Reason:** The finding itself says "No code change needed — this is a readability observation." The types (`OAuth`, `ApiAuth`, `WellKnownAuth`) and their constants (`AuthTypeOAuth`, `AuthTypeAPI`, `AuthTypeWellKnown`) are defined in `config.go` and are findable via standard IDE navigation or grep. Cross-file type references are normal in Go packages. This is an observation about code organization, not a defect.

### Backoff delay overflow for high attempt values

**Location:** `client.go:280` — exponential backoff calculation
**Date:** 2026-03-01

**Reason:** Duplicate of existing exception "Backoff bit-shift overflow on high retry counts." `WithMaxRetries` hard-caps at 10, and the `delay <= 0 || delay > maxBackoff` guard at line 281 catches any overflow. The maximum intermediate value is `500ms * 1024 = 512s`, well within `int64` range. The finding is speculative — it describes a concern that requires violating the enforced `maxRetryCap` invariant.

### queryparams non-pointer zero-value int/bool emitted even without omitempty

**Location:** `internal/queryparams/queryparams.go:145-163` — addFieldValue int/bool cases
**Date:** 2026-03-01

**Reason:** The behavior is explicitly documented in a comment at lines 141-144: "unlike strings (always omitted when empty), non-pointer int/bool zero values are emitted unless omitempty is set." The finding itself says "No code change needed." This describes documented, working-as-designed behavior. No current params struct uses a non-pointer int/bool without `omitempty`, so the "test documenting this behavior" suggestion addresses a purely hypothetical future scenario.

### No test for readAPIError with truncated large response body

**Location:** `errors.go:97-125` — readAPIError truncation
**Date:** 2026-03-01

**Reason:** The audit claims no test verifies truncation behavior. This is factually wrong. `TestReadAPIError_BodyTruncation` (errors_test.go:590-631) has three subtests: "body within limit is not marked truncated" (asserts `Truncated == false`), "body exceeding limit is marked truncated" (asserts `Truncated == true` and `len(Body) == maxErrorBodySize`), and "body exactly one byte over limit is marked truncated." The exact scenarios described in the suggested fix are already tested.

### No test for readAPIError with ReadErr (partial body read failure)

**Location:** `errors.go:99` — readAPIError ReadErr field
**Date:** 2026-03-01

**Reason:** The audit claims no test verifies that a mid-stream read failure populates `ReadErr`. This is factually wrong. `TestReadAPIError_PartialReadError` (errors_test.go:520) has subtests including "read error stored in ReadErr field" which uses a custom `io.Reader` that returns partial data then an error, and asserts both `Body` contains partial content and `ReadErr` is non-nil with the expected error via `errors.Is`.

### SSE maxSSETokenSize allows 32MB per token with no backpressure

**Location:** `packages/ssestream/ssestream.go:18` — maxSSETokenSize constant
**Date:** 2026-03-01

**Reason:** The finding itself says "This is likely acceptable for the intended use case (local dev server)" and only suggests adding godoc. The constant is already documented with a comment at lines 15-17 explaining why 32MB is needed. This is a documentation suggestion for a conscious design choice targeting a local dev server, not a code defect.

### apierror.Error type alias has no integration test

**Location:** `internal/apierror/apierror_test.go` — Error type alias
**Date:** 2026-03-01

**Reason:** Duplicate of existing exception "apierror.Error is unused but exported as a public type alias." The type is a Stainless leftover never constructed by the SDK. The finding itself says "low priority" and acknowledges the type is already documented as never-constructed. Adding an integration test for a type that is never returned by any SDK method has no practical value.

### No test for retry behavior on 429 status codes

**Location:** `client_do_test.go` — missing 429 retry integration test
**Date:** 2026-03-01

**Reason:** The audit claims neither 408 nor 429 has a dedicated retry test. The 429 claim is factually wrong. `TestRetryOn429` (client_test.go:52-79) uses a custom transport returning 429 on every request, creates a client, calls `Session.List`, asserts an error after exhausting retries, and verifies exactly 3 attempts were made. This is a full integration retry test, not just a unit check of `isRetryableStatus`.
