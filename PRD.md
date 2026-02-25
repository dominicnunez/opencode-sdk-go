# PRD: Complete Idiomatic Go SDK — Full Spec Coverage, Zero Stainless

Reference spec: `openapi.yml` (OpenAPI 3.1.1, 51 endpoints, committed in repo root)

This PRD has two goals:
1. Implement all 51 spec endpoints as idiomatic Go service methods
2. Remove every remaining Stainless internal — `apijson`, `apiquery`, `apiform`, `requestconfig`, `option`, `gjson`, `reflect`-based union registration — and replace with clean stdlib Go

**Rule: the OpenAPI spec defines WHAT the SDK does. This PRD defines HOW — and the how is pure Go.**

---

## Phase 1: Replace Stainless HTTP internals with stdlib client

Replace `internal/requestconfig` + `option/` with a single `client.go` `do()` method that uses `net/http` + `encoding/json` directly.

- [x] `client.go` already has `NewClient` with functional options (`WithBaseURL`, `WithHTTPClient`, `WithTimeout`, `WithMaxRetries`) — keep this
- [x] Rewrite `Client.do()` to NOT delegate to `requestconfig.ExecuteNewRequest`. Instead: build `*http.Request` directly, marshal JSON body with `encoding/json`, set headers, execute with retry loop, unmarshal response with `encoding/json`. Handle query params via `url.Values` from param structs' `URLQuery()` methods
- [x] Add `Client.doRaw()` variant that returns `*http.Response` for SSE streaming (used by `EventService.ListStreaming`)
- [x] Delete `internal/requestconfig/` entirely
- [x] Delete `option/requestoption.go` and `option/middleware.go` — replace any needed options with `ClientOption` functional options on the client itself
- [x] Update all service methods to use the new `do()` / `doRaw()` — no more `opts ...option.RequestOption` on individual methods

---

## Phase 2: Replace Stainless JSON machinery with encoding/json

Remove `internal/apijson` (2,185 lines), `internal/apiquery`, `internal/apiform` and replace with stdlib.

- [x] Delete `internal/apijson/` — all of it (decoder, encoder, field, port, registry, tag)
- [x] Delete `internal/apiform/` — all of it
- [x] Replace `internal/apiquery/` with a simple `queryParams(v interface{}) url.Values` helper that reads `query:"name"` struct tags (or just use the existing `URLQuery()` methods on param structs and delete apiquery too)
- [x] Remove `github.com/tidwall/gjson` dependency
- [x] Run `go mod tidy` to clean deps

---

## Phase 3: Replace union types with idiomatic Go discriminated unions

The spec has 6 union types. Stainless handles them with `reflect` + `gjson` + `apijson.RegisterUnion` + `init()` blocks. Replace with type-switch on a discriminator field.

**Pattern for all unions:**

```go
// Message is either UserMessage or AssistantMessage, discriminated by Role.
type Message struct {
	Role MessageRole `json:"role"`
	// Embed raw JSON for lazy decode
	raw json.RawMessage
}

func (m *Message) UnmarshalJSON(data []byte) error {
	// Peek at discriminator
	var peek struct {
		Role MessageRole `json:"role"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return err
	}
	m.Role = peek.Role
	m.raw = data
	return nil
}

func (m Message) AsUser() (*UserMessage, bool) {
	if m.Role != MessageRoleUser { return nil, false }
	var msg UserMessage
	if err := json.Unmarshal(m.raw, &msg); err != nil { return nil, false }
	return &msg, true
}

func (m Message) AsAssistant() (*AssistantMessage, bool) {
	if m.Role != MessageRoleAssistant { return nil, false }
	var msg AssistantMessage
	if err := json.Unmarshal(m.raw, &msg); err != nil { return nil, false }
	return &msg, true
}
```

Apply this pattern to each union:

- [x] **Message** (discriminator: `role`) → `UserMessage`, `AssistantMessage`. Remove `MessageUnion` interface, `apijson.RegisterUnion` init block, `apijson.Port` call. Add `AsUser()`, `AsAssistant()` methods.

- [x] **Part** (discriminator: `type`) → `TextPart`, `ReasoningPart`, `FilePart`, `ToolPart`, `StepStartPart`, `StepFinishPart`, `SnapshotPart`, `PatchPart`, `AgentPart`, `RetryPart`. Remove `PartUnion` interface and init block. Add `AsText()`, `AsReasoning()`, `AsFile()`, `AsTool()`, `AsStepStart()`, `AsStepFinish()`, `AsSnapshot()`, `AsPatch()`, `AsAgent()`, `AsRetry()` methods.

- [x] **ToolState** (discriminator: `status`) → `ToolStatePending`, `ToolStateRunning`, `ToolStateCompleted`, `ToolStateError`. Remove `ToolStateUnion` interface and init block. Add `AsPending()`, `AsRunning()`, `AsCompleted()`, `AsError()` methods.

- [x] **FilePartSource** (discriminator: `type`) → `FileSource`, `SymbolSource`. Remove `FilePartSourceUnion` interface and init block. Add `AsFile()`, `AsSymbol()` methods.

- [x] **Event** (discriminator: `type`) → 19 event types. Remove `EventUnion` interface and init block. Add `AsMessageUpdated()`, `AsSessionCreated()`, etc. for each event type.

- [x] **AssistantMessageError** (discriminator: `name`) → `MessageAbortedError`, `MessageOutputLengthError`, `APIError`, `ProviderAuthError`, `UnknownError`. Remove `AssistantMessageErrorUnion` interface and init block. Add `AsAborted()`, `AsOutputLength()`, `AsAPI()`, `AsProviderAuth()`, `AsUnknown()` methods.

- [x] **Auth** (discriminator: `type`) → `OAuth`, `ApiAuth`, `WellKnownAuth`. Remove `AuthUnion` interface and init block. Add `AsOAuth()`, `AsAPI()`, `AsWellKnown()` methods.

- [x] **ConfigMcp** (discriminator: `type`) → `McpLocalConfig`, `McpRemoteConfig`. Remove `ConfigMcpUnion` interface and init block. Add `AsLocal()`, `AsRemote()` methods.

- [x] **ConfigLsp** (discriminator: presence of `command` field) → `ConfigLspDisabled`, `ConfigLspObject`. Remove `ConfigLspUnion` interface and init block. Add `AsDisabled()`, `AsObject()` methods.

- [x] **SessionError** (discriminator: `name`) → `ProviderAuthError`, `UnknownError`, `MessageOutputLengthError`, `MessageAbortedError`, `SessionAPIError`. Remove `SessionErrorUnion` interface and init block. Add `AsProviderAuth()`, `AsUnknown()`, `AsOutputLength()`, `AsAborted()`, `AsAPI()` methods.

- [x] **PermissionPattern** (type-based: `string | array`) → convert to discriminated union with `AsString()` and `AsArray()` methods. Remove init block from sessionpermission.go.

- [x] Delete all `func init()` blocks that call `apijson.RegisterUnion` (7 remaining: all in config.go)
- [x] Remove all `reflect` and `gjson` imports

---

## Phase 4: Add missing service methods (18 endpoints)

All follow the existing pattern in the codebase. Use the spec for request/response types.

### Session service (`session.go`) — add 10 missing methods:

- [x] `Diff(ctx, id, params) → ([]FileDiff, error)` — `GET /session/{id}/diff`
- [x] `Fork(ctx, id, params) → (*Session, error)` — `POST /session/{id}/fork`. Params: `messageID string` (required)
- [x] `Shell(ctx, id, params) → (*AssistantMessage, error)` — `POST /session/{id}/shell`. Params: `agent string` (required), `command string` (required), `directory *string` (optional)
- [x] `Summarize(ctx, id, params) → (bool, error)` — `POST /session/{id}/summarize`. Params: `providerID string` (required), `modelID string` (required), `directory *string` (optional query param)
- [x] `Todo(ctx, id, params) → ([]Todo, error)` — `GET /session/{id}/todo`
- [x] `Unrevert(ctx, id, params) → (*Session, error)` — `POST /session/{id}/unrevert`
- [x] `Unshare(ctx, id, params) → (*Session, error)` — `DELETE /session/{id}/share`

Verify these already exist (my earlier scan found them but coverage script missed due to path concat):
- [x] Verify `Delete` works — `DELETE /session/{id}`
- [x] Verify `Get` works — `GET /session/{id}`
- [x] Verify `Update` works — `PATCH /session/{id}`

### Config service (`config.go`) — add 1 missing method:

- [ ] `Update(ctx, params) → (*Config, error)` — `PATCH /config`. Params: full Config struct

### Auth service — create new `auth.go`:

- [ ] Create `AuthService` struct on Client
- [ ] `Set(ctx, id, params) → error` — `PUT /auth/{id}`. Params: Auth union (OAuth | ApiAuth | WellKnownAuth)
- [ ] Wire into `Client` in `NewClient()`

### MCP service — create new `mcp.go`:

- [ ] Create `McpService` struct on Client
- [ ] `Status(ctx, params) → (*McpStatus, error)` — `GET /mcp`
- [ ] Wire into `Client` in `NewClient()`

### Tool service — create new `tool.go`:

- [ ] Create `ToolService` struct on Client
- [ ] `IDs(ctx, params) → (*ToolIDs, error)` — `GET /experimental/tool/ids`
- [ ] `List(ctx, params) → (*ToolList, error)` — `GET /experimental/tool`
- [ ] Wire into `Client` in `NewClient()`

### Session Permissions — verify coverage:

- [ ] Verify `SessionPermissionService.Reply(ctx, id, permissionID, params)` exists and matches `POST /session/{id}/permissions/{permissionID}`

---

## Phase 5: Add missing response/param types from spec

Check each schema in `openapi.yml` against existing Go types. Add any missing ones.

- [ ] `Todo` struct — `content string`, `status string`, `priority string`, `id string`
- [x] `FileDiff` struct — `file string`, `before string`, `after string`, `additions int64`, `deletions int64`  
- [ ] `SessionShellResponse` — check spec for response schema
- [ ] `SessionForkParams` — `messageID string` (required), `directory *string`
- [ ] `McpStatus` — check spec for response schema of `GET /mcp`
- [ ] `ToolIDs` — check spec
- [ ] `ToolList` / `ToolListItem` — `id string`, `description string`, `parameters interface{}`
- [ ] Verify all existing types match spec field names and types. Fix any drift.

---

## Phase 6: Clean up event.go streaming

- [ ] Keep `ssestream` package — it's clean and works
- [ ] Update `EventService.ListStreaming` to use `Client.doRaw()` instead of `requestconfig.ExecuteNewRequest`
- [ ] Ensure Event union type uses the Phase 3 discriminated union pattern (switch on `type` field)

---

## Phase 7: Update and expand tests

- [ ] Update all existing tests to work without `option.RequestOption` params
- [ ] Add tests for each new service method (auth, mcp, tool, missing session methods)
- [ ] Add tests for each union type's `As*()` methods — verify correct and incorrect discriminator values
- [ ] Add tests for `Client.do()` — mock HTTP server, verify request construction, headers, query params, JSON body
- [ ] Remove `internal/sessiontest/` if it only tested Stainless patterns
- [ ] Remove `internal/testutil/` if unused after cleanup

---

## Phase 8: Final cleanup

- [ ] Delete `internal/apijson/` (all files)
- [ ] Delete `internal/apiform/` (all files)
- [ ] Delete `internal/apiquery/` (all files — or keep minimal query helper if needed)
- [ ] Delete `internal/requestconfig/` (all files)
- [ ] Delete `option/` package (all files)
- [ ] Delete `internal/timeformat/` if unused
- [ ] Delete `aliases.go` if it only re-exports removed types
- [ ] Delete `ptr.go` if it only has `String()`, `Int()`, `Float()`, `Bool()` helpers that are no longer needed (or keep if tests use them for pointer construction)
- [ ] Remove `github.com/tidwall/gjson` from go.mod
- [ ] Run `go mod tidy`
- [ ] Run `go vet ./...`
- [ ] Run `go test -race ./...`
- [ ] Run `go build ./...`
- [ ] Run `golangci-lint run ./...` (add `.golangci.yml` if not present)
- [ ] Verify 0 imports of deleted packages remain
- [ ] Update README.md with new usage examples showing idiomatic patterns

---

## What NOT to change

- **Keep `ssestream` package** — already clean, handles SSE decoding well
- **Keep all spec-defined types** — every struct that maps to an OpenAPI schema stays
- **Keep functional options on Client** — `WithBaseURL`, `WithHTTPClient`, `WithTimeout`, `WithMaxRetries`
- **Keep `shared/` package** — contains error types used across services
- **Keep service pattern** — `ServiceName` struct with pointer to `Client`, methods take `(ctx, params)` return `(result, error)`

---

## Endpoint checklist (51 total)

When complete, every box should be checked:

### Session (16)
- [ ] `POST   /session` — session.create
- [ ] `GET    /session` — session.list
- [ ] `GET    /session/{id}` — session.get
- [ ] `PATCH  /session/{id}` — session.update
- [ ] `DELETE /session/{id}` — session.delete
- [ ] `POST   /session/{id}/abort` — session.abort
- [ ] `GET    /session/{id}/children` — session.children
- [ ] `POST   /session/{id}/command` — session.command
- [ ] `GET    /session/{id}/diff` — session.diff
- [ ] `POST   /session/{id}/fork` — session.fork
- [ ] `POST   /session/{id}/init` — session.init
- [ ] `GET    /session/{id}/message` — session.messages
- [ ] `GET    /session/{id}/message/{messageID}` — session.message
- [ ] `POST   /session/{id}/message` — session.prompt
- [ ] `POST   /session/{id}/revert` — session.revert
- [ ] `POST   /session/{id}/share` — session.share
- [x] `POST   /session/{id}/shell` — session.shell
- [ ] `POST   /session/{id}/summarize` — session.summarize
- [x] `GET    /session/{id}/todo` — session.todo
- [x] `POST   /session/{id}/unrevert` — session.unrevert
- [ ] `DELETE /session/{id}/share` — session.unshare
- [ ] `POST   /session/{id}/permissions/{permissionID}` — session.permission.reply

### Config (3)
- [ ] `GET    /config` — config.get
- [ ] `PATCH  /config` — config.update
- [ ] `GET    /config/providers` — config.providers

### Project (2)
- [ ] `GET    /project` — project.list
- [ ] `GET    /project/current` — project.current

### File (3)
- [ ] `GET    /file` — file.list
- [ ] `GET    /file/content` — file.read
- [ ] `GET    /file/status` — file.status

### Find (3)
- [ ] `GET    /find` — find.text
- [ ] `GET    /find/file` — find.files
- [ ] `GET    /find/symbol` — find.symbols

### Event (1)
- [ ] `GET    /event` — event.subscribe (SSE stream)

### TUI (8)
- [ ] `POST   /tui/append-prompt` — tui.appendPrompt
- [ ] `POST   /tui/clear-prompt` — tui.clearPrompt
- [ ] `POST   /tui/execute-command` — tui.executeCommand
- [ ] `POST   /tui/open-help` — tui.openHelp
- [ ] `POST   /tui/open-models` — tui.openModels
- [ ] `POST   /tui/open-sessions` — tui.openSessions
- [ ] `POST   /tui/open-themes` — tui.openThemes
- [ ] `POST   /tui/show-toast` — tui.showToast
- [ ] `POST   /tui/submit-prompt` — tui.submitPrompt

### App (2)
- [ ] `GET    /agent` — app.agents
- [ ] `POST   /log` — app.log

### Auth (1)
- [ ] `PUT    /auth/{id}` — auth.set

### Command (1)
- [ ] `GET    /command` — command.list

### Path (1)
- [ ] `GET    /path` — path.get

### MCP (1)
- [ ] `GET    /mcp` — mcp.status

### Tool (2)
- [ ] `GET    /experimental/tool/ids` — tool.ids
- [ ] `GET    /experimental/tool` — tool.list
