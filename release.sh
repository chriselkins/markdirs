#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

if [[ $# -lt 2 ]]; then
  echo "Usage: $0 v1.0.2 \"Release notes or changelog\""
  exit 1
fi

VERSION="$1"
NOTES="$2"

if [ -z "$VERSION" ]; then
  echo "Usage: $0 v1.0.2 \"Release notes or changelog\""
  exit 1
fi

if ! command -v gh >/dev/null 2>&1; then
  echo "GitHub CLI (gh) is required. Install it: https://cli.github.com/"
  exit 2
fi

go mod tidy

# Ensure clean git state
git diff-index --quiet HEAD -- || { echo "Uncommitted changes!"; exit 1; }

echo "Tagging release ${VERSION}..."
git tag "${VERSION}"
git push
git push --tags

echo "Building release artifacts..."
make release

echo "Creating GitHub release and uploading artifacts..."
gh release create "${VERSION}" dist/*.zip dist/SHA256SUMS \
  --title "${VERSION}" \
  --notes "${NOTES:-"Release $VERSION"}"

echo "Release ${VERSION} published!"

# Prompt pkg.go.dev to fetch the new tag
echo "Requesting pkg.go.dev to fetch ${VERSION}..."
go list -m "github.com/chriselkins/markdirs@${VERSION}" || true
curl -sSf "https://proxy.golang.org/github.com/chriselkins/markdirs/@v/${VERSION}.info" > /dev/null || true
echo "To make sure pkg.go.dev indexes your new version, visit:"
echo "  https://pkg.go.dev/github.com/chriselkins/markdirs@${VERSION}"
echo "and click the 'Request' button if available."