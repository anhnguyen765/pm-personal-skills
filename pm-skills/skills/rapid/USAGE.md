# Weekly Report Skill - Usage Guide

## Quick Start

To generate your weekly report, run the following command from the project root:

```bash
./pm-skills/skills/weekly-report/scripts/generator
```

## Workflow

The generator will guide you through 6 steps:

### Step 1: Fetch Latest Report
Automatically retrieves your most recent report from Confluence (Payment Platform space) to use as a reference.

### Step 2: Pull Jira Activities
Queries your Jira activities from the past 7 days (resolved, in-progress, and recently updated issues).

### Step 3: What I Did Last Week
Review the Jira-populated tasks and:
- Confirm or edit the task name
- Set the status (DONE, ONGOING, TBD, TBC)
- Add or edit notes (timeline, dependencies, concerns, links)
- Optionally add tasks not captured in Jira (meetings, grooming, stakeholder calls)

### Step 4: What I Will Do This Week
Enter tasks you plan to work on this week:
- Task name
- Status (TODO, ONGOING, TBD, TBC)
- Notes (deadlines, stakeholders, context)

### Step 5: Where I Need Support
Identify blockers and support needed:
- Blocker/task name
- Status (TBD, TBC)
- Details (who to ask, timeline, dependencies)

### Step 6: Generate & Save
The tool generates a markdown report and saves it to:
```
pm-skills/output/sandbox/weekly-report-[YYYY-MM-DD].md
```

Copy the content to your Confluence page at:
**Payment Platform > Meeting Notes/2026/PO Weekly Reports/Anhnhm3**

## Status Codes

- **DONE**: Completed in past week
- **ONGOING**: Started, still in progress
- **TODO**: Planned, not yet started
- **TBD**: To Be Determined (waiting for clarity/decision)
- **TBC**: To Be Confirmed (waiting for approval or confirmation)

## Notes Format

In the Notes column, you can include:
- Timeline and deadlines
- Stakeholder names and teams
- Email thread references (e.g., *Re: [Cross-border] Location Gating Rule Grooming*)
- Concerns or blockers
- Confluence or Jira links
- Vietnamese + English mixed naturally

Newlines will be automatically converted to `<br>` tags for Confluence compatibility.

## Sandbox Before Publishing

The report is always saved to `pm-skills/output/sandbox/` first. Review it before copying to Confluence:

```bash
cat pm-skills/output/sandbox/weekly-report-*.md
```

## Environment Setup

Ensure `.env` is configured with:
```
JIRA_URL=https://jira.zalopay.vn
JIRA_API_TOKEN=<your-token>
CONFLUENCE_URL=https://confluence.zalopay.vn
CONFLUENCE_EMAIL=<your-email>
CONFLUENCE_API_TOKEN=<your-token>
```

## Troubleshooting

**Jira issues not fetching?**
- Check `.env` configuration
- Verify your Jira API token is valid
- Ensure your assignment JQL is correct (uses `currentUser()`)

**Can't find latest Confluence report?**
- Verify the page URL: https://confluence.zalopay.vn/x/kSB6Eg
- Check that you have read access to the Payment Platform space

**Report formatting issues?**
- Pipe characters (`|`) in task names will be escaped automatically
- Newlines in Notes will convert to `<br>` for Confluence
- Review the generated markdown file in sandbox before publishing
