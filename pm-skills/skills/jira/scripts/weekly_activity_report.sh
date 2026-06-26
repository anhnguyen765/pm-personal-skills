#!/bin/bash

# Weekly Activity Report Generator
# Generates a markdown table with all Jira activities from the past week
# Usage: ./weekly_activity_report.sh [email]
# Default: uses anhnhm3 (from .env)

# Get user email from argument or environment
USER_EMAIL="${1:-anhnhm3}"

# Get the jira tool path
JIRA_TOOL="./pm-skills/tooling/jira/jira"

# Check if jira tool exists
if [ ! -f "$JIRA_TOOL" ]; then
    echo "Error: Jira tool not found at $JIRA_TOOL"
    exit 1
fi

# Create temporary files for results
TEMP_FILE=$(mktemp)
TEMP_TABLE=$(mktemp)
trap "rm -f $TEMP_FILE $TEMP_TABLE" EXIT

# Fetch all issues created or modified by the user in the past week
# Search across PC, ZPBO, and PMS projects
echo "Fetching your Jira activities from last week..." >&2

JQL_QUERY="(project = PC OR project = ZPBO OR project = PMS) AND updated >= -7d AND (assignee = $USER_EMAIL OR reporter = $USER_EMAIL)"

$JIRA_TOOL search "$JQL_QUERY" --limit 100 > "$TEMP_FILE" 2>&1

# Parse the results and create a formatted table
echo ""
echo "# Your Jira Activities — Last Week"
echo ""
echo "| # | Ticket | Status | Priority | Domain | Description |"
echo "|---|--------|--------|----------|--------|-------------|"

# Extract tickets and fetch details
line_num=0
while IFS='|' read -r _ key type status priority assignee updated summary _; do
    key=$(echo "$key" | xargs)

    # Skip header and empty lines
    if [[ "$key" == "Key" || -z "$key" ]]; then
        continue
    fi

    line_num=$((line_num + 1))

    # Clean up values
    type=$(echo "$type" | xargs)
    status=$(echo "$status" | xargs)
    priority=$(echo "$priority" | xargs)
    summary=$(echo "$summary" | xargs | cut -c1-75)

    # Extract domain from summary (text in brackets like [Domain])
    # Use sed for macOS compatibility instead of grep -oP
    domain=$(echo "$summary" | sed -E 's/^.*\[([^]]+)\].*/\1/' | head -1)
    if [ "$domain" = "$summary" ]; then
        domain="General"
    fi

    # Build table row with Jira link
    echo "| $line_num | [$key](https://jira.zalopay.vn/browse/$key) | $status | $priority | $domain | $summary |"

done < "$TEMP_FILE"

echo ""
echo "---"
echo ""

# Generate summary statistics
echo "## Summary Statistics"
echo ""

# Count by status
TOTAL=$(tail -n +2 "$TEMP_FILE" | wc -l | xargs)

NEW_COUNT=$(tail -n +2 "$TEMP_FILE" | grep -c " New " || echo "0")
DONE_COUNT=$(tail -n +2 "$TEMP_FILE" | grep -c " DONE " || echo "0")
IN_DEV=$(tail -n +2 "$TEMP_FILE" | grep -c " In Dev " || echo "0")
READY_TEST=$(tail -n +2 "$TEMP_FILE" | grep -c " Ready " || echo "0")

P1_COUNT=$(tail -n +2 "$TEMP_FILE" | grep -c " P1 " || echo "0")
P2_COUNT=$(tail -n +2 "$TEMP_FILE" | grep -c " P2 " || echo "0")
P3_COUNT=$(tail -n +2 "$TEMP_FILE" | grep -c " P3 " || echo "0")

echo "| Metric | Count |"
echo "|--------|-------|"
echo "| **Total Activities** | $TOTAL |"
echo "| **Status: New** | $NEW_COUNT |"
echo "| **Status: DONE** | $DONE_COUNT |"
echo "| **Status: In Dev** | $IN_DEV |"
echo "| **Status: Ready for Testing** | $READY_TEST |"
echo "| **Priority P1** | $P1_COUNT |"
echo "| **Priority P2** | $P2_COUNT |"
echo "| **Priority P3** | $P3_COUNT |"
echo ""

echo "---"
echo ""
echo "## Key Themes (Review & Update)"
echo ""
echo "Based on ticket domains, consider these themes for your report:"
echo ""
echo "- **Cross-border & Payments**: A+, Tenpay, Multi-market expansion"
echo "- **EMVCO & QR Global**: International payment flows, new asset types (pmcid 58)"
echo "- **Configuration & UX**: Error messages, SoF limits, fund flow routing"
echo "- **TapToPay & Lending**: Partial reversal, debt APIs, product features"
echo "- **Operations & Support**: UM checks, liquidation workflows"
echo ""
echo "---"
echo ""
echo "**Report generated**: $(date +'%Y-%m-%d %H:%M')"
echo "**Email**: $USER_EMAIL"
echo ""
echo "✨ **Ready to share**: Copy this entire output into your Monday standup, Confluence, or Slack update."
