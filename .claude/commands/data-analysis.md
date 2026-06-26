---
name: data-analysis
description: Synthesises platform metrics, cross-border KPIs, and market data to identify next actionable items for Payment Core and Cross-border PO. Use when analyzing provided data, identifying trends, and surfacing data-backed next steps.
---

# Data Analysis Skill

This skill empowers the Payment Core Platform PO and Cross-border PO with data-driven insights, trend analysis, and prioritised next-step recommendations grounded in provided metrics and data.

## Core Capabilities

### 1. Payment Core Platform Metrics Analysis
Analyse performance indicators for Payment Core (Fund Movement Core and Payment Engine, including Zalopay Accounting System).
- **Metrics Tracked**: Transaction volume, success rates, latency, error rates, user adoption, settlement velocity, ledger consistency.
- **Data Source**: Analyst-provided datasets, dashboards, and reports.
- **Trend Detection**: Identify spikes, regressions, seasonal patterns, and inflection points.
- **Segmentation**: Analyse by payment method, acquiring domain, settlement flow, or ledger type.

### 2. Cross-Border Domain Metrics
Monitor outbound and inbound cross-border payment flows and user experiences.
- **Core Focus**: Cross-border payments for Zalopay users and partner users in Vietnam.
- **KPIs**: Outbound/inbound transaction volume, settlement corridors, success rates, user conversion, partner adoption.
- **User Experience**: Payment flow completion rates, error patterns, friction points by corridor.
- **Stakeholder Alignment**: Data analysis grounded in cross-functional collaboration with all squads (led by Merchant Solutions).

### 3. Next Actionable Items (Prioritised)
Synthesise data to surface high-impact next steps grounded in provided metrics.
- **Quick Wins**: Low-effort, high-impact improvements (e.g., reduce manual settlement reviews by 20%, improve UX friction).
- **Growth Levers**: Tactical moves to accelerate cross-border corridor expansion, user conversion, or settlement velocity.
- **Risk Mitigation**: Fraud patterns, compliance risks, or operational bottlenecks identified in data.
- **Feature Prioritisation**: Data-backed business cases for new Payment Core or cross-border features.

### 4. Root Cause Analysis
Investigate data anomalies and provide diagnostic insight.
- **Trend Correlation**: Link metric shifts to business events (campaign launches, feature releases, market changes).
- **Pattern Clustering**: Group failures or drop-offs by corridor, user segment, payment method, or flow stage.
- **Cross-Domain Context**: Consider interactions between Payment Core infrastructure and cross-border experiences.

---

## Workflow: Data Analysis for Next Actions

### 1) Define the Question
Clarify the metric scope and analytical goal (e.g., "Why did cross-border settlement velocity drop 15% last week?" or "What are the top friction points in the outbound payment flow?").

### 2) Ingest & Validate Data
- Accept datasets, dashboards, and reports provided by the analyst or team.
- Clarify the time period, segmentation, and any data quality notes.
- Identify relevant dimensions: corridor, payment method, user segment, flow stage, or ledger type.

### 3) Analyse Trends
- Normalise metrics over time (YoY, QoQ, WoW).
- Identify inflection points and correlate with known business events or changes.
- Segment by relevant dimensions (corridor, payment method, user cohort, flow step).
- Detect anomalies and outliers.

### 4) Synthesise Insights
Produce a structured output:
- **Current State**: What the data shows right now.
- **Trends**: Recent trajectory and momentum (↑/↓/→).
- **Root Causes**: Hypotheses grounded in the data, with supporting evidence.
- **Next 3 Actions**: Prioritised by impact × effort, with data-backed rationale.

---

## Rules of Engagement
- **Precision**: Always cite the exact metric, time period, and source (Jira dashboard, Confluence page, or Grafana board).
- **Contextual Awareness**: Consider platform constraints (AEv2 limitations, settlement SLAs, compliance windows).
- **Actionability**: Every insight must map to a next step (code review, feature flag, campaign, investigation ticket).
- **Australian English**: Summarise, prioritise, centralised.

## Example Scenarios

| Scenario | Analysis Approach |
|---|---|
| **Weekly cross-border performance review** | Ingest provided KPI data (volume, success rate, settlement time by corridor); identify trends and top friction points; recommend priority actions. |
| **Root cause: Payment Core latency spike** | Analyse provided latency metrics by acquiring domain and payment engine service; correlate with transaction volume or ledger load; surface infrastructure improvements. |
| **Cross-border corridor expansion** | Compare metrics across existing corridors (success rates, user conversion, partner adoption); identify best-performing patterns; recommend next corridor entry strategy. |
| **User experience friction analysis** | Segment provided flow completion data by step and user cohort; identify drop-off patterns; prioritise UX improvements. |
| **Feature ROI analysis** | Compare pre/post metrics for a Payment Core or cross-border feature; quantify business impact; support go/no-go decisions. |

---

## Key Metrics by Domain

### Payment Core Platform (FMC/PE)
**Fund Movement Core (FMC)**
- Order creation success rate
- Acquiring domain latency
- Settlement velocity (D+1, D+2)

**Payment Engine (PE)**
- Transaction orchestration latency (p50, p95, p99)
- Payment state transition duration
- Ledger consistency checks (passed/failed)

**Zalopay Accounting System (ZAS, part of PE)**
- Settlement reconciliation time
- Ledger integrity and audit compliance
- Accounting flow latency

### Cross-Border Domain
**Outbound Payments**
- Corridor-specific transaction volume and growth
- Payment success rates by corridor
- Settlement time and velocity
- User conversion rate from domestic to cross-border

**Inbound Payments**
- Receiving flow volume and partner adoption
- Payout success rates by method
- User experience metrics (flow completion, drop-off points)

**Partner Experience**
- Partner onboarding metrics
- Settlement SLA compliance
- Partner adoption and transaction volume growth

---

## Integration Points
- **Data Ingestion**: Accept datasets from analytics team, BI dashboards, and provided reports.
- **Confluence**: KPI definitions, roadmap context, architectural constraints (AEv2, idempotency, settlement compliance).
- **Cross-Functional Collaboration**: Align insights with Payment Core squads (FX, AML, Settlement) and Merchant Solutions (cross-border).
