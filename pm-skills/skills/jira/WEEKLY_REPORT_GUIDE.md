# Weekly Activity Report — Quick Start Guide

## 🎯 Purpose
Generate an automated summary of all your Jira activities from the past week—perfect for Monday team standups, Confluence posts, or asynchronous team updates.

## 📋 What You Get
1. **Activity Table** — All tickets you created/modified last week with direct Jira links
2. **Summary Statistics** — Breakdown by status (New, DONE, In Dev, Ready for Testing) and priority (P1, P2, P3)
3. **Key Themes** — Pre-filled themes based on your ticket domains; review and refine
4. **Copy-Paste Ready** — Markdown-formatted output ready for team sharing

## 🚀 How to Use (Every Monday Morning)

### Option 1: Command Line
```bash
./pm-skills/skills/jira/scripts/weekly_activity_report.sh anhnhm3
```

### Option 2: With Custom Email
```bash
./pm-skills/skills/jira/scripts/weekly_activity_report.sh your.email@example.com
```

### Option 3: From Claude Code (via `/jira` slash command)
The skill now includes this capability—ask Claude to "Generate my weekly activity report for Monday."

## 📤 What to Do With the Output

1. **Copy the entire output** (everything from "# Your Jira Activities" to the end)
2. **Paste into**:
   - Monday standup doc (Confluence/Docs)
   - Team Slack channel
   - Email update to your manager
   - Weekly team sync notes

3. **Customise**:
   - Review the **Key Themes** section and refine to match your week's focus
   - Add 1-2 sentence commentary if desired
   - Highlight blockers or wins

## 📊 Example Output Structure

```
# Your Jira Activities — Last Week

| # | Ticket | Status | Priority | Domain | Description |
...

## Summary Statistics

| Metric | Count |
| **Total Activities** | 16 |
...

## Key Themes (Review & Update)

- **Cross-border & Payments**: ...
- **EMVCO & QR Global**: ...
...

✨ **Ready to share**: Copy this entire output into your Monday standup...
```

## 🎨 Tips & Tricks

- **Filter by week**: The script automatically searches `updated >= -7d` (last 7 days)
- **All projects**: Searches PC (Payment Core), ZPBO (Cross-border), and PMS (PM Sync)
- **Only your work**: Shows tickets where you're assignee OR reporter
- **Reusable template**: Save the output as a template and modify as needed
- **Team visibility**: Run for team members too by passing their email: `./weekly_activity_report.sh colleague@example.com`

## ⚡ Troubleshooting

| Issue | Solution |
|-------|----------|
| Script not found | Ensure you're in the repo root: `cd /path/to/jora` |
| Permission denied | Run: `chmod +x ./pm-skills/skills/jira/scripts/weekly_activity_report.sh` |
| No results | Check your email spelling and ensure you have recent Jira activity |
| Jira tool error | Verify `.env` credentials are correct (see `.env.example`) |

## 📚 Related Commands

- **Velocity analysis**: `./pm-skills/skills/jira/scripts/velocity_analyzer`
- **Jira search**: `./pm-skills/tooling/jira/jira search "<jql>" --limit N`
- **Get issue**: `./pm-skills/tooling/jira/jira get PC-1234`

---

**Updated**: 2026-06-15 | **Skill Version**: Jira Master v1.1
