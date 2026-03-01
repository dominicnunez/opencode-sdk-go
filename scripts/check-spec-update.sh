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

# Fetch upstream .stats.yml
echo "Checking upstream spec..."
STATS=$(gh api "repos/$UPSTREAM_REPO/contents/.stats.yml" --jq '.content' | base64 -d)

UPSTREAM_HASH=$(echo "$STATS" | grep 'openapi_spec_hash' | awk '{print $2}')
UPSTREAM_URL=$(echo "$STATS" | grep 'openapi_spec_url' | awk '{print $2}')
UPSTREAM_ENDPOINTS=$(echo "$STATS" | grep 'configured_endpoints' | awk '{print $2}')

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
    echo "Downloading updated spec..."
    curl -sL "$UPSTREAM_URL" -o "$LOCAL_SPEC"
    NEW_HASH=$(shasum -a 256 "$LOCAL_SPEC" | awk '{print $1}')
    echo "New local hash: $NEW_HASH"
    if [ "$NEW_HASH" != "$UPSTREAM_HASH" ]; then
        echo "ERROR: Downloaded spec hash ($NEW_HASH) does not match upstream ($UPSTREAM_HASH)"
        rm "$LOCAL_SPEC"
        exit 1
    fi
    echo "✅ Spec updated. Review changes with: git diff specs/openapi.yml"
    echo ""
    echo "Next steps:"
    echo "  1. Diff the spec to see what changed"
    echo "  2. Update SDK methods/types to match"
    echo "  3. Run tests: go test -race ./..."
    echo "  4. Commit: git add specs/openapi.yml && git commit -m 'chore: update OpenAPI spec'"
else
    echo "Run with --update to download the new spec:"
    echo "  $0 --update"
fi

exit 1
