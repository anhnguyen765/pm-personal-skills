---
name: jira
description: Unified Jira skill for velocity analysis, stuck task resolution, and writing technical grooming-ready ticket descriptions for Payment Core.
---

# Jira Master

This unified skill combines advanced velocity analysis with technical ticket writing expertise for Payment Core.

## Core Capabilities

### 1. Velocity Management & Analysis
Maintain team momentum by identifying bottlenecks and providing actionable feedback.
- **Cross-Reference Jira**: Fetches real-time status and assignees for "Stuck" tasks.
- **Generate Commentary**: Provides suggested Jira comments to unblock progress.
- **Sprint Analysis**: When performing any "sprint analysis" or querying a sprint, always use the prefix **"PCDPC"** (e.g., "PCDPC - Sprint 26.06.A"). Automatically prepend this prefix if it is missing from the user's request.
- **Command**: `./scripts/velocity_analyzer` 

### 2. Technical Ticket Writing (Payment Core)
Draft grooming-ready Jira tickets that bridge product vision and technical complexity.
- **Domain Focus**: Fund Movement Core (Order creation, acquiring domain) and Payment Engine (orchestration, transitions, ledger, Zalopay Accounting System).
- **Standards**: Australian English (summarise, prioritise).
- **Template-Driven**: Follows a strict structure for Context, Requirements, and Acceptance Criteria.

---

## Workflow: Writing Jira Tickets

### 1) Context Gathering
Identify business objectives, user pain points, and affected flows (Pay, Refund, Settle, etc.).

### 2) Technical Integration
Incorporate technical constraints (AEv2, ZAS, idempotency) and Confluence-based requirements (Flow IDs, Accounting Codes).

### 3) Structured Output
Use Jira wiki style (`h2.`, `h3.`) with the following mandatory sections:
- **Context**: The "Why".
- **Requirements**: Product overview, impacted services, key logic changes, technical constraints, and data/config needs.
- **Acceptance Criteria**: Functional, Technical (observability/integrity), Deployment (rollout strategy), and Documentation.

---

## Best Practices for Velocity
- **Highlight Blockers**: Probe technical blockers for "In Dev" tasks stuck for > 3 days.
- **Clarify Priorities**: Validate the necessity of "New" tasks already in the sprint.
- **Celebrate Progress**: Acknowledge when stuck tasks move to "Resolved" or "Done".

## Resources
- `scripts/velocity_analyzer.go`: Core automation logic for velocity checks.
