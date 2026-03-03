# Risks

> Real findings consciously accepted — architectural cost, external constraints, disproportionate effort.
> Managed by sfk willie. Follow the entry format below.
>
> Entry format:
> ### Plain language description
> **Location:** `file/path:line` — optional context
> **Date:** YYYY-MM-DD
> **Reason:** Explanation (can be multiple lines)

### Union types cannot be constructed programmatically for serialization

**Location:** `config.go`, `event.go`, `session.go` — all union types using `raw json.RawMessage`
**Date:** 2026-02-28

**Reason:** Adding constructor functions (e.g. `NewConfigMcpLocal`) for every union type would expand the public API surface significantly. The current approach stores `json.RawMessage` internally and provides `As*()` accessors for deserialization. Callers who need to construct unions can JSON-roundtrip through the variant type. Response-only unions (Event, Part) don't need constructors at all. The few request unions (ConfigMcp, ConfigUpdateParams) are low-traffic enough that a JSON roundtrip is acceptable.

### Six BashUnion types have identical method implementations

**Location:** `config.go:148-186, 291-329, 434-472, 807-844, 950-988, 1064-1102`
**Date:** 2026-02-28

**Reason:** Each BashUnion type maps to a distinct schema in the OpenAPI spec and is a separate public API type. Extracting a generic helper would require Go generics with type constraints, adding complexity for 6 small types (~40 lines each) that rarely change. The methods are generated-pattern code where a bug fix can be applied mechanically. The cost of the generic abstraction outweighs the risk of divergence.

### Deprecated config fields still parsed

**Location:** `config.go:53-54,68,73-74`
**Date:** 2026-02-22

**Reason:** The Config struct fields (autoshare, mode, layout) reflect the upstream OpenAPI spec. The spec defines these fields as deprecated. Removing them would break deserialization of API responses that still include them. The deprecation comments are accurate and guide users to migrate.

### Bytes buffer allocation in SSE hot path

**Location:** `packages/ssestream/ssestream.go:81`
**Date:** 2026-02-22

**Reason:** The `bytes.NewBuffer(nil)` call per event is a minor allocation in a streaming context. For typical usage patterns, the GC overhead is negligible. Using `sync.Pool` would add complexity for an optimization that would only benefit extremely high-throughput scenarios. No performance issue has been reported or measured.

### ListStreaming response body not closed if readAPIError panics on custom transport

**Location:** `event.go:56-57` — non-2xx status path
**Date:** 2026-02-28

**Reason:** If a custom transport's `io.ReadCloser` panics during `io.ReadAll` inside `readAPIError`, the response body is never closed. This requires a custom transport that violates Go's `io.Reader` contract (panicking instead of returning an error). The `doRaw` path avoids this via a defer in `do`, but `ListStreaming` intentionally bypasses `do` for SSE semantics. Adding a defer-close before `readAPIError` would close the body twice on the normal path (since `readAPIError` already closes it). The risk is theoretical — no well-behaved transport panics from `Read`.

### SSE stream data accumulation not configurable by callers

**Location:** `packages/ssestream/ssestream.go:87,141-148` — maxDataBytes field and limit check
**Date:** 2026-03-01

**Reason:** The `eventStreamDecoder` has an unexported `maxDataBytes` field that defaults to `maxSSETokenSize` (32MB). Many small `data:` lines can accumulate up to this limit before the check triggers. The 32MB default is generous but matches the scanner token limit for consistency. Exposing `maxDataBytes` via a `DecoderOption` would expand the public API surface of the `ssestream` package. The current limit is adequate for the SDK's use case (OpenCode server events are small JSON objects). Callers who need tighter bounds can implement a custom decoder via `RegisterDecoder`.
