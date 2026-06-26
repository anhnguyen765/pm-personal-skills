#!/usr/bin/env bash
set -euo pipefail

# Bundles skills into .skill (ZIP) files and deploys to active .gemini/skills.
# Uses a sandbox-first approach for verification.

REPO="$(cd "$(dirname "$0")/../.." && pwd)"
SKILLS_DIR="$REPO/pm-skills/skills"
SANDBOX_DIR="$REPO/pm-skills/output/sandbox"
BUNDLES_DIR="$REPO/pm-skills/bundles"
ACTIVE_SKILLS_DIR="$REPO/.gemini/skills"

mkdir -p "$SANDBOX_DIR"
mkdir -p "$BUNDLES_DIR"

# 1. Compile Go Scripts
echo "Compiling Go scripts..."
find "$SKILLS_DIR" -name "*.go" | while read -r go_file; do
  script_dir=$(dirname "$go_file")
  binary_name=$(basename "$go_file" .go)
  echo "  -> Compiling $binary_name in $script_dir"
  (cd "$script_dir" && go build -o "$binary_name" "$(basename "$go_file")")
done

# 2. Package to Sandbox
for skill_path in "$SKILLS_DIR"/*; do
  if [ -d "$skill_path" ]; then
    name=$(basename "$skill_path")
    echo "Packaging $name..."
    
    sandbox_bundle="$SANDBOX_DIR/$name.skill"
    (cd "$SKILLS_DIR" && zip -rq "$sandbox_bundle" "$name")
    echo "  -> Sandbox bundle created: $sandbox_bundle"
  fi
done

echo ""
echo "Verification & Deployment:"
echo "1. Verify bundles in: $SANDBOX_DIR"
echo "2. To deploy to active CLI, run: cp -r $SKILLS_DIR/* $ACTIVE_SKILLS_DIR/"
echo "3. To publish bundles, run: cp $SANDBOX_DIR/*.skill $BUNDLES_DIR/"
