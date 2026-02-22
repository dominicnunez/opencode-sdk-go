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

## Won't Fix

<!-- Real findings not worth fixing — architectural cost, external constraints, etc. -->

### Panic in library code can crash applications

**Location:** `internal/apijson/decoder.go:211`, `internal/apijson/field.go:30`, `internal/apiquery/encoder.go:190,235,246`, `option/requestoption.go:101`
**Date:** 2026-02-22

**Reason:** This SDK is auto-generated from an OpenAPI spec by Stainless. The panics in generated code occur in edge cases that should never happen in normal usage (union type not in registry, invalid map keys, unsupported array formats). Replacing these panics with error returns would require modifying the Stainless generator, which is an external tool. The `WithMaxRetries` panic is documented in the function's docstring.

### Deprecated config fields still parsed

**Location:** `config.go:53-54,68,73-74,1058`
**Date:** 2026-02-22

**Reason:** The Config struct and its deprecated fields (autoshare, mode, layout) are generated from the OpenAPI spec. We cannot control what fields the spec defines. The deprecation comments are accurate and users should migrate, but removing the fields would break the generated code contract with the API.

### Global sync.Map for encoder/decoder caching grows unbounded

**Location:** `internal/apijson/decoder.go:18`, `internal/apijson/encoder.go:19`, `internal/apiquery/encoder.go:15`
**Date:** 2026-02-22

**Reason:** This is a standard caching pattern for reflection-based serialization. The cache is bounded by the number of distinct types used, which in practice is limited and stable for a given application. Memory profiling would be needed to demonstrate an actual problem before adding complexity like LRU eviction.

## Intentional Design Decisions

<!-- Findings that describe behavior which is correct by design -->

### Hardcoded default base URL uses localhost

**Location:** `option/requestoption.go:266`
**Date:** 2026-02-22

**Reason:** This SDK is designed for local development against the opencode CLI server which runs on localhost:54321 by default. Users targeting a production API should explicitly set the base URL using `WithBaseURL()`. The function name `WithEnvironmentProduction` is misleading (it should perhaps be `WithEnvironmentLocal`), but changing it now would be a breaking change.

### Debug middleware logs sensitive data

**Location:** `option/middleware.go:23-33`
**Date:** 2026-02-22

**Reason:** The `WithDebugLog` function is explicitly documented as "for debugging and development purposes only" and "should not be used in production." Users must explicitly opt-in by passing this middleware. Adding header redaction would add complexity for a debugging tool and could hide issues that debugging is meant to reveal.

### SSE stream error not returned directly

**Location:** `event.go:42-51`
**Date:** 2026-02-22

**Reason:** This is a standard pattern for streaming APIs in Go. The stream object must be returned so callers can iterate over events, and embedding the initial connection error in the stream allows a single return signature. The pattern is documented and callers are expected to check `stream.Err()` before iteration, similar to how database rows work.
