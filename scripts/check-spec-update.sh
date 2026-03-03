#!/usr/bin/env bash
# Check if the upstream OpenAPI spec has changed since our last download.
# Compares the spec hash from anomalyco/opencode-sdk-go .stats.yml
# against our local specs/openapi.yml.

set -euo pipefail

UPSTREAM_REPO="anomalyco/opencode-sdk-go"
LOCAL_SPEC="specs/openapi.yml"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

require_command() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "ERROR: $1 is required but not installed" >&2
        exit 1
    fi
}

extract_stats_value() {
    local key="$1"
    local value
    value=$(echo "$STATS" | awk -F': ' -v key="$key" '$1 == key {print $2; exit}')
    value=${value#\"}
    value=${value%\"}
    if [ -z "$value" ]; then
        echo "ERROR: missing $key in upstream .stats.yml" >&2
        exit 1
    fi
    echo "$value"
}

decode_base64() {
    local input="$1"
    local decoded

    if decoded=$(printf '%s' "$input" | base64 --decode 2>/dev/null); then
        printf '%s' "$decoded"
        return 0
    fi
    if decoded=$(printf '%s' "$input" | base64 -d 2>/dev/null); then
        printf '%s' "$decoded"
        return 0
    fi
    if decoded=$(printf '%s' "$input" | base64 -D 2>/dev/null); then
        printf '%s' "$decoded"
        return 0
    fi

    return 1
}

require_command gh
require_command base64
require_command shasum
require_command awk
require_command grep

# Fetch upstream .stats.yml
echo "Checking upstream spec..."
if ! STATS_B64=$(gh api "repos/$UPSTREAM_REPO/contents/.stats.yml" --jq '.content'); then
    echo "ERROR: failed to fetch upstream .stats.yml from GitHub API" >&2
    exit 1
fi
if ! STATS=$(decode_base64 "$STATS_B64"); then
    echo "ERROR: failed to decode upstream .stats.yml content" >&2
    exit 1
fi

UPSTREAM_HASH=$(extract_stats_value "openapi_spec_hash")
UPSTREAM_URL=$(extract_stats_value "openapi_spec_url")
UPSTREAM_ENDPOINTS=$(extract_stats_value "configured_endpoints")

# Hash our local spec
LOCAL_HASH=$(shasum -a 256 "$LOCAL_SPEC" | awk '{print $1}')

echo "Upstream hash:    $UPSTREAM_HASH"
echo "Local hash:       $LOCAL_HASH"
echo "Upstream endpoints: $UPSTREAM_ENDPOINTS"

if [ "$UPSTREAM_HASH" = "$LOCAL_HASH" ]; then
    echo "✅ Spec is up to date."
    exit 0
fi

echo ""
echo "⚠️  Spec has changed!"
echo "Upstream URL: $UPSTREAM_URL"
echo ""

if [ "${1:-}" = "--update" ]; then
    require_command curl
    echo "Downloading updated spec..."
    TMPFILE=$(mktemp)
    trap 'rm -f "$TMPFILE"' EXIT
    if ! curl -fsSL "$UPSTREAM_URL" -o "$TMPFILE"; then
        echo "ERROR: failed to download upstream spec from $UPSTREAM_URL" >&2
        exit 1
    fi
    NEW_HASH=$(shasum -a 256 "$TMPFILE" | awk '{print $1}')
    echo "New local hash: $NEW_HASH"
    if [ "$NEW_HASH" != "$UPSTREAM_HASH" ]; then
        echo "ERROR: Downloaded spec hash ($NEW_HASH) does not match upstream ($UPSTREAM_HASH)"
        exit 1
    fi
    mv "$TMPFILE" "$LOCAL_SPEC"
    echo "✅ Spec updated. Review changes with: git diff specs/openapi.yml"
    echo ""
    echo "Next steps:"
    echo "  1. Diff the spec to see what changed"
    echo "  2. Update SDK methods/types to match"
    echo "  3. Run tests: go test -race ./..."
    echo "  4. Commit: git add specs/openapi.yml && git commit -m 'chore: update OpenAPI spec'"
    exit 0
else
    echo "Run with --update to download the new spec:"
    echo "  $0 --update"
fi

exit 1
