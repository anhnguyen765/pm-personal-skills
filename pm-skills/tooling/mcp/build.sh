#!/usr/bin/env bash
# Build the pm-tools MCP server binary.
# Run from any directory — output lands next to this script.
set -e
DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$DIR"
go build -o pm-tools-mcp server.go
echo "Built: $DIR/pm-tools-mcp"
echo ""
echo "Claude Desktop config location:"
echo "  macOS:   ~/Library/Application Support/Claude/claude_desktop_config.json"
echo "  Windows: %APPDATA%\\Claude\\claude_desktop_config.json"
echo ""
echo "Merge the contents of claude_desktop_config.json into that file,"
echo "replacing the command path and credentials with your own values."
