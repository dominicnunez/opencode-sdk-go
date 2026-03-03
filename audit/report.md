### [Bug] Session command rejects valid empty arguments
- **Severity**: Medium
- **File**: session.go:131
- **Details**: `SessionService.Command` rejects empty or whitespace-only `Arguments`, but the OpenAPI contract only requires the `arguments` field to be present and typed as string (`specs/openapi.yml:1428` and `specs/openapi.yml:1436`) with no `minLength` constraint.
- **Suggested fix**: Remove the `strings.TrimSpace(params.Arguments) == ""` validation and require only `params != nil` plus non-empty `Command`.

### [Configuration] Spec check script declares an unused dependency
- **Severity**: Low
- **File**: scripts/check-spec-update.sh:74
- **Details**: The script requires `grep` but never invokes it. This can fail in minimal environments even though the script logic does not need `grep`.
- **Suggested fix**: Remove `require_command grep` or add a real `grep` usage if that dependency is intentional.

### [Performance] Mock startup polling repeatedly rescans entire log
- **Severity**: Low
- **File**: scripts/mock:31
- **Details**: The daemon startup loop runs `grep -q` against `.prism.log` every 100ms. Each poll rescans the entire file, which scales poorly as logs grow.
- **Suggested fix**: Use incremental log reads (for example `tail -n` with bounded reads) or a direct health-check probe to avoid full-file rescans each poll.

### [Code Quality] Startup retry budget uses a hardcoded multiplier
- **Severity**: Low
- **File**: scripts/mock:10
- **Details**: `STARTUP_MAX_ATTEMPTS=$((STARTUP_TIMEOUT_SECONDS * 10))` hardcodes `10` instead of deriving attempts from `STARTUP_POLL_INTERVAL_SECONDS`, so changing the interval silently breaks timeout math.
- **Suggested fix**: Compute max attempts from timeout and poll interval (for example via `awk`/`bc`) and keep timeout semantics in one place.
