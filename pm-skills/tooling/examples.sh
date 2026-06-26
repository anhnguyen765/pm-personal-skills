#!/bin/bash
# Example usage of tooling binaries

set -e

TOOLING_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=== Web Fetcher Examples ==="
echo ""
echo "Fetching example.com as plain text:"
"$TOOLING_DIR/web/web" fetch "https://example.com"
echo ""
echo "---"
echo ""
echo "Fetching example.com as HTML:"
"$TOOLING_DIR/web/web" fetch "https://example.com" --format html | head -20
echo ""
echo "---"
echo ""
echo "Fetching with custom timeout:"
"$TOOLING_DIR/web/web" fetch "https://example.com" --timeout 15
