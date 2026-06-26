# PM Skills Pack — Setup Guide

Welcome to the **PM Skills Pack**, a comprehensive suite of Product Management tools and skills for teams working with Jira, Confluence, and document generation workflows.

## Table of Contents
- [Quick Start](#quick-start)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Skills Overview](#skills-overview)
- [Tooling](#tooling)
- [Troubleshooting](#troubleshooting)

---

## Quick Start

```bash
# 1. Clone or extract this repository
cd pm-skills

# 2. Create your .env file with Jira and Confluence credentials
cp ../.env.example ../.env
# Edit ../.env with your actual credentials

# 3. Test a skill
/jira                    # Jira velocity analysis
/prd                     # Generate a PRD
/ascii                   # Create ASCII diagrams
```

---

## Prerequisites

Before using the PM Skills Pack, ensure you have:

### Required
- **macOS, Linux, or WSL2** (Windows)
- **Jira account** with API token access
- **Confluence account** with API token access
- **Claude Code** (installed and configured)

### Optional (depending on skills)
- **Excel** or a spreadsheet application (for Product Audit reports)
- **Go 1.16+** (if compiling tooling from source)
- **Bash 4.0+** for script execution

---

## Installation

### 1. Clone the Repository
```bash
git clone <repository-url> pm-skills
cd pm-skills
```

### 2. Set File Permissions
Make scripts executable:
```bash
chmod +x pm-skills/tooling/jira/*
chmod +x pm-skills/tooling/confluence/*
chmod +x pm-skills/tooling/web/*
chmod +x pm-skills/scripts/*.sh
chmod +x pm-skills/skills/*/scripts/*.sh
```

### 3. Verify Installation
```bash
ls -la pm-skills/tooling/jira/
# You should see: jira, daily_report, velocity_history (or similar binaries)
```

---

## Configuration

### Step 1: Create Your `.env` File

In the **parent directory** of `pm-skills/` (typically your project root), create a `.env` file:

```bash
cd ..  # Move to parent directory
cp pm-skills/CONFIG.env.example .env
```

### Step 2: Fill in Your Credentials

Edit `.env` with your actual credentials:

```env
# Jira Configuration
JIRA_URL=https://your-instance.atlassian.net
JIRA_EMAIL=your-email@company.com
JIRA_API_TOKEN=your-jira-api-token

# Confluence Configuration
CONFLUENCE_URL=https://your-instance.atlassian.net/wiki
CONFLUENCE_EMAIL=your-email@company.com
CONFLUENCE_API_TOKEN=your-confluence-api-token

# Product Audit (optional)
PRODUCT_AUDIT_TEMPLATE_PATH=/path/to/product_audit_template.xlsx
PRODUCT_AUDIT_OUTPUT_DIR=./pm-skills/output/sandbox/
```

### Step 3: Obtain Credentials

#### Jira API Token
1. Go to https://id.atlassian.com/manage-profile/security/api-tokens
2. Create an API token
3. Copy the token to your `.env` file

#### Confluence API Token
1. Go to https://id.atlassian.com/manage-profile/security/api-tokens
2. Create an API token (same as Jira)
3. Copy the token to your `.env` file

---

## Skills Overview

Each skill is a self-contained Claude Code module. Invoke them with `/` in Claude Code:

### Core PM Skills

| Skill | Command | Purpose |
|-------|---------|---------|
| Jira | `/jira` | Velocity analysis, sprint reports, issue grooming |
| PRD | `/prd` | Generate system-oriented Product Requirements Documents |
| Product Audit | `/product-audit` | Quality and compliance auditing with Excel export |
| Advisor | `/advisor` | Synthesise Confluence documentation for architectural guidance |
| Memo | `/memo` | Socratic debate + executive decision memos |
| Grilling | `/grilling` | Stress-test assumptions before building |
| ASCII | `/ascii` | Professional ASCII diagrams for system flows |
| Translator | `/translator` | Translate features to technical, product, and business views |
| Data Analysis | `/data-analysis` | Synthesise KPIs and identify actionable insights |
| Rapid | `/rapid` | Weekly activity reports and team summaries |

### Skill Structure

Each skill in `skills/<skill-name>/` contains:

```
skills/jira/
├── SKILL.md              # Skill definition and prompts
├── README.md             # Skill documentation
├── USAGE.md              # Usage examples (if applicable)
├── scripts/              # Helper scripts (optional)
│   ├── velocity_analyzer
│   └── weekly_activity_report.sh
├── assets/               # Templates, examples, references
└── references/           # Documentation, guidelines
```

---

## Tooling

The `tooling/` directory provides Go binaries for direct Jira, Confluence, and web access. These power the skills from the CLI.

### Jira Tooling

```bash
# Get a full issue
./pm-skills/tooling/jira/jira get PCDPC-1234

# Search with JQL
./pm-skills/tooling/jira/jira search "project = PCDPC AND status != Done" --limit 20

# Generate daily standup report
./pm-skills/tooling/jira/daily_report

# Generate velocity history
./pm-skills/tooling/jira/velocity_history
```

**Environment variables used:**
- `JIRA_URL`: Your Jira instance URL
- `JIRA_EMAIL`: Your Jira email
- `JIRA_API_TOKEN`: Your Jira API token

### Confluence Tooling

```bash
# Search Confluence
./pm-skills/tooling/confluence/confluence search "FX Management" --limit 10

# Fetch page content by ID
./pm-skills/tooling/confluence/confluence fetch 12345678

# Fetch full page with formatting
./pm-skills/tooling/confluence/fetch_full 12345678
```

**Environment variables used:**
- `CONFLUENCE_URL`: Your Confluence instance URL
- `CONFLUENCE_EMAIL`: Your Confluence email
- `CONFLUENCE_API_TOKEN`: Your Confluence API token

### Web Tooling

```bash
# Fetch and clean web content (plain text)
./pm-skills/tooling/web/web fetch "https://example.com/api"

# Fetch as HTML
./pm-skills/tooling/web/web fetch "https://example.com" --format html

# Fetch as JSON
./pm-skills/tooling/web/web fetch "https://api.example.com/data" --format json

# Custom timeout
./pm-skills/tooling/web/web fetch "https://example.com" --timeout 60
```

**Environment variables (optional):**
- `WEB_USER_AGENT`: Custom user agent
- `WEB_TIMEOUT`: Default timeout in seconds (default: 30)
- `WEB_CACHE_DIR`: Cache directory for responses

See [`tooling/WEB_FETCHER.md`](./tooling/WEB_FETCHER.md) for detailed examples.

---

## Common Workflows

### Generate a Weekly Activity Report
```bash
./pm-skills/skills/jira/scripts/weekly_activity_report.sh your-jira-user
```
See [`skills/jira/WEEKLY_REPORT_GUIDE.md`](./skills/jira/WEEKLY_REPORT_GUIDE.md) for details.

### Analyze Sprint Velocity
```bash
./pm-skills/skills/jira/scripts/velocity_analyzer
```

### Fetch Confluence Documentation
```bash
./pm-skills/tooling/confluence/fetch_full 12345678
```

### Run a Product Audit
```bash
# Generate audit JSON
go run pm-skills/skills/product-audit/scripts/excel_helper.go \
  -mode write \
  -json-file pm-skills/skills/product-audit/scripts/my_audit.json
```

---

## Output & Artifacts

Generated artifacts are saved to:
- **Sandbox (staging)**: `pm-skills/output/sandbox/`
- **Final reports**: `pm-skills/output/reports/`

Before publishing artifacts externally (to Jira, OneDrive, etc.), stage them in the sandbox for review.

---

## Project Structure

```
pm-skills/
├── README.md                  # Overview and quick reference
├── SETUP.md                   # This file — setup instructions
├── skills/                    # Individual PM skills
│   ├── jira/                  # Jira automation and reporting
│   ├── prd/                   # PRD generation
│   ├── product-audit/         # Quality auditing
│   ├── advisor/               # Confluence-based architecture guidance
│   ├── memo/                  # Decision memo drafting
│   ├── ascii/                 # ASCII diagrams
│   ├── translator/            # Stakeholder communication
│   ├── data-analysis/         # KPI synthesis
│   ├── grilling/              # Interview/stress-testing
│   └── ... (other skills)
├── tooling/                   # CLI utilities
│   ├── jira/                  # Jira binaries and utilities
│   ├── confluence/            # Confluence binaries and utilities
│   ├── web/                   # Web fetcher and utilities
│   └── README.md              # Tooling reference
├── bundles/                   # Pre-packaged skill artifacts
├── output/                    # Generated reports and artifacts
│   ├── sandbox/               # Staging area for review
│   └── reports/               # Final published reports
└── scripts/                   # Build and packaging scripts
```

---

## Troubleshooting

### Error: `command not found: jira`
**Solution**: Ensure the tooling is executable:
```bash
chmod +x pm-skills/tooling/jira/*
```

### Error: `JIRA_URL not set`
**Solution**: Create and populate your `.env` file (see [Configuration](#configuration)).

### Error: `401 Unauthorized`
**Solution**: Verify your Jira/Confluence credentials:
1. Check that `JIRA_API_TOKEN` and `CONFLUENCE_API_TOKEN` are correct
2. Regenerate tokens at https://id.atlassian.com/manage-profile/security/api-tokens
3. Ensure tokens are placed in the `.env` file without extra spaces

### Error: `Connection refused`
**Solution**: Verify your Jira/Confluence URLs are correct:
```bash
# Test connectivity
curl -u your-email:your-api-token https://your-instance.atlassian.net/rest/api/3/myself
```

### Skill not responding
**Solution**:
1. Ensure Claude Code is running the latest version
2. Check the `.env` file for missing or malformed credentials
3. Run a tooling command directly to isolate the issue:
   ```bash
   ./pm-skills/tooling/jira/jira get PCDPC-1 --verbose
   ```

---

## Support & Contributing

If you encounter issues or have suggestions:
1. Check the [Troubleshooting](#troubleshooting) section above
2. Review skill-specific documentation in `skills/<skill-name>/README.md`
3. Open an issue on the project repository

---

## License & Attribution

This PM Skills Pack is provided as-is for Product Management teams. Ensure you comply with your organization's policies when sharing or distributing artifacts.

---

**Last updated**: 26 June 2026
