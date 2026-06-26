---
name: product-audit
description: Systematically evaluate product quality, operational readiness, and compliance using the internal Product Audit Excel template. Use when auditing a feature, evaluating a product launch, or checking for UX/logic gaps.
---

# Product Audit: Setup & Usage Guide

Systematic product auditing through the **System-as-User** lens. Interfaces with the company-standard Product Audit Excel template, anchored in Confluence technical documentation.

## System-as-User Perspective
Treat Payment Core (PE/FMC) as the **Product** and internal services as **Users**:
- **Product**: Payment Core (PE & FMC).
- **Users**: ACv2, AE, ZAS, Refund Core, and other technical services.
- **Naming Consistency** = Param Integrity.
- **Confirm/Execute** = CTA (Call to Action).
- **Service Orchestration** = Navigation.

## Core Workflow

### 1. Research: Technical Source of Truth
- **Developer Centre**: Consult [Zalopay Developers v2](https://developers.zalopay.vn/v2/) and [v1](https://developers.zalopay.vn/v1/) for public API contracts.
- **Internal Specs**: Search "Zalopay Technology Management" (ZTM) and "Payment Engine" Confluence spaces:
  ```bash
  ./pm-skills/tooling/confluence/confluence search "<keyword>" --space ZTM
  ./pm-skills/tooling/confluence/confluence search "<keyword>" --space PE
  ```
- **High-Signal Keywords**: `FMC/PE`, `Asset Exchange`, `ZAS`, `AEv2`, `Corev2`, `ACv2`.

### 2. The Technical Advisor Module
Fetch and synthesise documentation before making audit judgments:
- **Search**: `./pm-skills/tooling/confluence/confluence search "<topic>" [--space ZTM|PE] [--limit N]`
- **Fetch**: `./pm-skills/tooling/confluence/confluence fetch <page_id>`

### 3. System-as-User Audit
Evaluate the flow against EVERY checklist item in the Excel template.
- **Surgical Technical Detail**: Evidence MUST include parameter names, gRPC context details, and orchestration steps.
- **Cross-Reference Standards**: Use Table 5 in `pm-skills/skills/product-audit/references/team-context.md`.

### 4. Technical Validation (Mandatory)
Before recording any finding, cross-verify observed behaviour against the official architectural blueprint.
- Confirm `AC -> AE` and `AE -> ZAS` interactions align with documented sequence diagrams.
- Every 'Pass', 'Flag', or 'Fail' must cite a specific Confluence page.
- Divergence from Source of Truth = 'High' severity 'Technical Misalignment' risk.

### 5. Output Generation
Update the Excel file:
```bash
go run pm-skills/skills/product-audit/scripts/excel_helper.go \
  -mode write -json-file <audit-json-file>
```
Generate a text-based **Audit Summary Report** for architectural risks.

## Audit Summary Structure
1. **AUDIT OVERVIEW**: Squad/Product, Auditor, Date, Source of Truth.
2. **AUDIT SCOPE**: Flows and user journeys reviewed.
3. **KEY RISKS IDENTIFIED (Flag/Fail)**: Top findings with dimension, criterion, and gap.
4. **HIGH SEVERITY ITEMS (🔴)**: Requires immediate action.
5. **MEDIUM SEVERITY ITEMS (🟡)**: Address this quarter.
6. **LOW SEVERITY / PASS HIGHLIGHTS**: Strengths and minor polish.
7. **PROPOSED NEXT STEPS**: Owner, action, target date.
8. **OPEN QUESTIONS / DEPENDENCIES**: Blockers needing other teams.
9. **AUDIT METRICS**: Final counts of Pass, Flag, Fail, N/A.

## Resources
- `pm-skills/skills/product-audit/references/audit-guidelines.md`
- `pm-skills/skills/product-audit/references/team-context.md`
