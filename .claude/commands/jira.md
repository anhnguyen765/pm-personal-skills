---
name: jira
description: Unified Jira skill for velocity analysis, stuck task resolution, and writing technical grooming-ready ticket descriptions for Payment Core.
---

# Jira Master

This unified skill combines advanced velocity analysis with technical ticket writing expertise for Payment Core.

## Core Capabilities

### 1. Weekly Activity Report
Generate an automated summary of all your Jira activities from the past week—tickets created, modified, or assigned to you.
- **Output**: Markdown table with ticket links, status, priority, domain, and descriptions.
- **Includes**: Summary statistics (count by status/priority) and automated key themes analysis.
- **Perfect for**: Weekly team syncs every Monday or asynchronous team updates.
- **Command**: `./pm-skills/skills/jira/scripts/weekly_activity_report.sh [email]`
- **Example**: `./pm-skills/skills/jira/scripts/weekly_activity_report.sh anhnhm3`
- **Output Location**: Direct to stdout—copy/paste into your weekly report doc or Confluence.

### 2. Velocity Management & Analysis
Maintain team momentum by identifying bottlenecks and providing actionable feedback.
- **Cross-Reference Jira**: Fetches real-time status and assignees for "Stuck" tasks.
- **Get issue**: `./pm-skills/tooling/jira/jira get <issue-key>`
- **Search issues**: `./pm-skills/tooling/jira/jira search "<jql>" [--limit N]`
- **Generate Commentary**: Provides suggested Jira comments to unblock progress.
- **Sprint Analysis**: Always use the prefix **"PCDPC"** (e.g., "PCDPC - Sprint 26.06.A"). Auto-prepend if missing.
- **Velocity**: `./pm-skills/skills/jira/scripts/velocity_analyzer` (for complex multi-sprint calculations)

### 3. Technical Ticket Writing (Payment Core)
Draft grooming-ready Jira tickets that bridge product vision and technical complexity.
- **Domain Focus**: Fund Movement Core (Order creation, acquiring domain) and Payment Engine (orchestration, transitions, ledger, ZaloPay Accounting System).
- **Standards**: Australian English (summarise, prioritise).
- **Template-Driven**: Follows a strict structure for Context, Requirements, and Acceptance Criteria.

## Workflow: Writing Jira Tickets

### 1) Context Gathering
Identify business objectives, user pain points, and affected flows (Pay, Refund, Settle, etc.).

### 2) Technical Integration
Incorporate technical constraints (AEv2, ZAS, idempotency) and Confluence-based requirements (Flow IDs, Accounting Codes).

### 3) Structured Output
Use Jira wiki style (`h2.`, `h3.`) with mandatory sections:
- **Context**: The "Why".
- **Requirements**: Product overview, impacted services, key logic changes, technical constraints, data/config needs.
- **Acceptance Criteria**: Functional, Technical (observability/integrity), Deployment (rollout strategy), and Documentation.

## Best Practices for Velocity
- **Highlight Blockers**: Probe technical blockers for "In Dev" tasks stuck > 3 days.
- **Clarify Priorities**: Validate necessity of "New" tasks already in the sprint.
- **Celebrate Progress**: Acknowledge when stuck tasks move to "Resolved" or "Done".

---

## Monday Weekly Report Workflow

**Every Monday morning**, run this command to generate your weekly activity summary:

```bash
./pm-skills/skills/jira/scripts/weekly_activity_report.sh anhnhm3
```

This outputs a markdown-formatted report with:
1. **Activity Table** — All tickets you created/modified last week with links, status, priority, domain
2. **Summary Statistics** — Counts by status (New, DONE, In Dev, Ready for Testing) and priority (P1, P2, P3)
3. **Key Themes** — Auto-detected from ticket domains; review and refine before sharing

**Use case**: Copy the output into your Monday standup, Confluence weekly post, or team Slack update.

## Resources
- `pm-skills/skills/jira/scripts/velocity_analyzer.go`: Core automation logic.
- `pm-skills/skills/jira/scripts/weekly_activity_report.sh`: Weekly activity report generator script.
