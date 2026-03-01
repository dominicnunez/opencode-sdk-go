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

### apierror.Error is unused but exported as a public type alias

**Location:** `internal/apierror/apierror.go:12-17` — aliased as `opencode.Error` in `aliases.go:8`
**Date:** 2026-02-28

**Reason:** `apierror.Error` is never constructed anywhere in the SDK — it's a Stainless leftover. However, it's exposed as the public type `opencode.Error`. Removing it would be a breaking API change for any caller that references the type. The type is inert (never returned by any SDK method), so it causes no runtime harm.

### Six BashUnion types have identical method implementations

**Location:** `config.go:148-186, 291-329, 434-472, 807-844, 950-988, 1064-1102`
**Date:** 2026-02-28

**Reason:** Each BashUnion type maps to a distinct schema in the OpenAPI spec and is a separate public API type. Extracting a generic helper would require Go generics with type constraints, adding complexity for 6 small types (~40 lines each) that rarely change. The methods are generated-pattern code where a bug fix can be applied mechanically. The cost of the generic abstraction outweighs the risk of divergence.

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
