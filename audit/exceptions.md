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

**Reason:** The audit claims error messages are "inconsistent with the pattern used elsewhere" and suggests "some use 'received empty string' while others might use different phrasing." However, all 17 error messages for missing required parameters in the codebase follow the exact same format: `missing required parameter 'X' (received empty string)`. The audit provides no evidence of actual inconsistency and cannot cite any examples of different phrasing because none exist.

### SSE buffer size integer overflow claim

**Location:** `packages/ssestream/ssestream.go:45`
**Date:** 2026-02-22

**Reason:** The audit claims `bufio.MaxScanTokenSize<<sseBufferMultiplier` (64KB << 9 = ~32MB) "could theoretically overflow on 32-bit systems." This is mathematically incorrect. The result is 33,554,432 bytes (~32MB), which is well under the 32-bit signed int maximum of 2,147,483,647 (~2.1GB). No overflow is possible.

## Won't Fix

<!-- Real findings not worth fixing — architectural cost, external constraints, etc. -->

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

### Bytes buffer allocation in SSE hot path

**Location:** `packages/ssestream/ssestream.go:81`
**Date:** 2026-02-22

**Reason:** The `bytes.NewBuffer(nil)` call per event is a minor allocation in a streaming context. For typical usage patterns, the GC overhead is negligible. Using `sync.Pool` would add complexity for an optimization that would only benefit extremely high-throughput scenarios. No performance issue has been reported or measured.

## Intentional Design Decisions

<!-- Findings that describe behavior which is correct by design -->

### SSE stream error not returned directly

**Location:** `event.go:20-51`
**Date:** 2026-02-22

**Reason:** This is a standard pattern for streaming APIs in Go. The stream object must be returned so callers can iterate over events, and embedding the initial connection error in the stream allows a single return signature. The pattern is documented and callers are expected to check `stream.Err()` before iteration, similar to how database rows work.
