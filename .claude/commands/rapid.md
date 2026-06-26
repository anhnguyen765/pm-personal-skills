---
name: rapid
description: Generate weekly reports for team meetings by synthesising Jira activities, Confluence benchmarks, and user input on progress, new tasks, and blockers.
---

# RAPID Weekly Report Skill

This skill automates weekly progress reporting by pulling your Jira activities, fetching the previous week's report from Confluence as a benchmark, and guiding you through a structured interview to update task progress, identify new work, and surface blockers for team support.

## Core Capabilities

### 1. Confluence Benchmark Retrieval
Fetch your latest weekly report from the **Payment Platform** space under **Meeting Notes/2026/PO Weekly Reports/Anhnhm3** page tree to use as a reference and starting point.
- Identifies the most recent report
- Extracts existing task structure and status
- Uses completed items as a baseline for this week's update

### 2. Jira Activity Integration
Pull your Jira activities from the past week (resolves, in-progress, comments) to automatically populate the "What I Did Last Week" section.
- Filters by assignee (your account)
- Groups by epic or issue key
- Captures resolved and progressed items
- Identifies blocked or in-review work

### 3. Structured Report Generation
Guide you through an interactive workflow to gather:
- **Updates to last week's tasks**: Review Jira-populated items, edit status and notes
- **New tasks for this week**: Add any work not captured in Jira
- **Blockers and support requests**: Identify what needs team help to unblock

### 4. Formatted Output
Generate a markdown report matching your team's standard format:
- Header with date, owner, status
- Three main sections (Last Week / This Week / Support Needed)
- Markdown tables with Task, Status, Notes columns
- Ready for publication to Confluence

---

## Workflow: Weekly Report Generation

### 1) Fetch Latest Report
Retrieve your most recent report from Confluence (Payment Platform space) as a reference benchmark.
- If found: Use as baseline; highlight completed vs. carry-over tasks
- If not found: Start fresh with empty sections

### 2) Pull Jira Activities
Query your Jira activities from the past 7 days:
- Resolved issues
- In-progress work
- Recent comments/updates
- Blocked or review-waiting items

### 3) Interview: What I Did Last Week
Present Jira-populated tasks and ask:
- Confirm/edit task description
- Update status (DONE, ONGOING, TBD, TBC)
- Add/edit notes (timeline, dependencies, concerns, links)
- Add any tasks not in Jira (meetings, grooming, stakeholder alignment)

### 4) Interview: What I Will Do This Week
Prompt for:
- Carry-over tasks from last week still in progress
- New sprint/planned work
- Grooming, planning, or review activities
- Status for each (TODO, ONGOING, TBD, TBC)
- Detailed notes (deadlines, dependencies, stakeholder context)

### 5) Interview: Where I Need Support
Ask:
- Which tasks have blockers?
- What decisions or approvals are pending?
- What cross-functional help is needed?
- Who should be looped in? (include names/teams in notes)

### 6) Generate & Stage Report
- Render markdown report in your standard format
- Save to `pm-skills/output/sandbox/weekly-report-[date].md`
- Ready for copy-paste to Confluence or further editing

---

## Rules of Engagement

- **Jira Authority**: Resolved and completed items from Jira are the source of truth; notes and status can be edited interactively.
- **Confluence Baseline**: Previous week's report is used only as a reference—carry-over tasks will be re-evaluated.
- **Clarity in Notes**: Include timeline, stakeholder names, email chains (brief reference), and blockers clearly.
- **Vietnamese + English**: Mix languages naturally as in your reports (no forced translation).
- **Markdown Tables**: Tables must be clean and aligned for Confluence publishing.
- **Sandbox First**: Always stage to `pm-skills/output/sandbox/` before publishing to Confluence.

---

## Example Workflow

**Input:** User runs `/weekly-report` on 2026-06-22.

**Step 1 - Fetch:** Retrieves 2026-06-15 report from Confluence.
```
✓ Found latest report: 2026-06-15
Previous week had 10 completed, 5 ongoing tasks.
```

**Step 2 - Pull Jira:** Finds 8 resolved issues in past 7 days under Payment Core epics.

**Step 3 - Interview:**
```
Populate "What I Did Last Week" from Jira:
1. Cross-border Gating Rules Grooming (Jira: PC-1234) — Status? [User: TBD]
   Notes? [User: Cần đánh giá lại solution. Mail: Re: [Cross-border]...]
2. Payment Core Quarter Retro (Jira: PC-5678) — Status? [User: DONE]
   Notes? [User: (leave blank)]
...
```

**Step 4 & 5 - Gather This Week & Support:**
```
New tasks for this week?
[User: Sprint 7A Planning & Grooming, Finalise Product Planning for Q3]

Blockers or support needed?
[User: Partners Routing Logic (waiting HuyDBQ), Cross-border Gating Rules (QA review)]
```

**Step 6 - Generate:** Outputs markdown report to `pm-skills/output/sandbox/weekly-report-2026-06-22.md`.

---

## Integration Points

- **Confluence**: Fetch latest report from Payment Platform > Meeting Notes/2026/PO Weekly Reports/Anhnhm3
- **Jira**: Pull user's resolved and in-progress issues (past 7 days)
- **Output**: Markdown tables, ready for Confluence publishing
- **Sandbox**: `pm-skills/output/sandbox/weekly-report-[date].md`

---

## Key Fields in Report

### Header
- **Title**: `# [Name] - Weekly Report - [YYYY/MM/DD]`
- **Owner**: Your full name (from Jira user profile)
- **Status**: DONE (when report is complete)

### Section 1: What I Did Last Week
| # | Task | Status | Notes |
|---|------|--------|--------|
| (auto-numbered) | Task name or Jira key | DONE / ONGOING / TBD / TBC | Timeline, dependencies, links, concerns |

### Section 2: What I Will Do This Week
| # | Task | Status | Notes |
|---|------|--------|--------|
| (auto-numbered) | Task name | TODO / ONGOING / TBD / TBC | Deadlines, stakeholders, context |

### Section 3: Where I Need Support
| # | Task | Status | Notes |
|---|------|--------|--------|
| (auto-numbered) | Task/Blocker | TBD / TBC | Blocker details, who to ask, timeline |

---

## Status Codes
- **DONE**: Completed in past week
- **ONGOING**: Started, still in progress
- **TODO**: Planned, not yet started
- **TBD**: To Be Determined (waiting for clarity/decision)
- **TBC**: To Be Confirmed (waiting for approval or confirmation)

