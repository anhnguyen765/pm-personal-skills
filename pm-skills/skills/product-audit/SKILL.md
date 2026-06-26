---
name: product-audit
description: Systematically evaluate product quality, operational readiness, and compliance using the internal Product Audit Excel template. Use this when a user needs to audit a feature, evaluate a product launch, or check for UX/logic gaps.
---

# Product Audit: Setup & Usage Guide

This skill enables systematic product auditing through the **System-as-User** lens. It interfaces with the company-standard Product Audit Excel template, anchored in technical documentation from Confluence.

## 🚀 Getting Started: Setup Guide

Follow these steps to configure the skill for the Payment Core environment.

### 1. Configure Environment Variables
Ensure your `.env` file in the project root contains the following keys.

```env
# Path to the master Excel template (usually in this skill's assets folder)
PRODUCT_AUDIT_TEMPLATE_PATH=/path/to/product-audit/assets/product_audit_template.xlsx

# Directory where your completed audits will be saved
PRODUCT_AUDIT_OUTPUT_DIR=/path/to/output/reports/

# Confluence API Configuration
CONFLUENCE_URL=https://your-company.atlassian.net/wiki
CONFLUENCE_EMAIL=your-email@company.com
CONFLUENCE_API_TOKEN=your_token_here
```

### 2. System-as-User Perspective
To ensure technical integrity, every audit must treat Payment Core (PE/FMC) as the **Product** and internal services as **Users**:
- **Product**: Payment Core (PE & FMC).
- **Users**: ACv2, AE, ZAS, Refund Core, and other technical systems.
- **Naming Consistency** = Param Integrity.
- **Confirm/Execute** = CTA (Call to Action).
- **Service Orchestration** = Navigation.

## 🛠 Core Workflow & Modules

This skill integrates a **Technical Advisor Module** to ensure every audit is grounded in verified documentation.

### 1. Research: Technical Source of Truth
-   **Prioritize Architecture & Integration**: 
    -   **Developer Centre**: Consult [Zalopay Developers v2](https://developers.zalopay.vn/v2/) and [v1](https://developers.zalopay.vn/v1/) for public API contracts.
    -   **Internal Specs**: Use the internal **Advisor Module** to search "Zalopay Technology Management" (ZTM) and "Payment Engine" Confluence spaces.
    -   **High-Signal Keywords**: Use `FMC/PE`, `Asset Exchange`, `ZAS`, `AEv2`, `Corev2`, `ACv2`.

### 2. The Technical Advisor Module
Use this internal capability to fetch and synthesize documentation before making audit judgments:
- **Search**: `go run scripts/advisor.go search "topic"`
- **Fetch**: `go run scripts/advisor.go fetch <page_id>`
- **Role**: Synthesize architectural patterns, point out dependencies, and align findings with Zalopay engineering standards.

### 3. The System-as-User Audit
Systematically evaluate the flow against **EVERY** checklist item in the Excel template.
-   **Surgical Technical Detail**: Evidence MUST include parameter names, gRPC context details, and orchestration steps (e.g., "Verified `UniversalID` propagation in gRPC metadata").
-   **Cross-Reference Standards**: Use Table 5 in `references/team-context.md` to audit param integrity, error quality, and orchestration safety.

### 4. Technical Validation & Cross-Verification
**Mandatory Step.** Before recording any finding or generating the Excel output, you MUST use the **Advisor Module** to cross-verify the observed system behavior against the official architectural blueprint.
-   **System-as-Actor Verification**: Confirm that the interaction between `AC -> AE` and `AE -> ZAS` aligns with the documented sequence diagrams.
-   **Grounding**: Every 'Pass', 'Flag', or 'Fail' result must cite a specific Confluence page or technical spec retrieved via the Advisor.
-   **Perspective Check**: Ensure 'Users' in the checklist are mapped to technical services (e.g., 'User Feedback' = 'gRPC Error Metadata').
-   **Resolve Discrepancies**: If current implementation diverges from the 'Source of Truth' found by the Advisor, mark it as a 'High' severity 'Technical Misalignment' risk.

### 5. Documentation & Evidence Collection
Record findings for **ALL** 75 criteria. 
-   **Evidence Quality**: Include gRPC method names, parameter keys, and specific state transition logic (e.g., "Verified AE's `Deduct` gRPC call includes the `adjusted_timestamp` for stale-checks").
-   **Automated N/A Handling**: Criteria not applicable to the backend orchestration layer will be populated as 'N/A' with a technical justification.

### 6. Output Generation
-   Update the Excel file using `scripts/excel_helper.go`.
-   Generate a text-based **Audit Summary Report** focused on architectural risks.

## 📖 Detailed Guidelines

- [Audit Guidelines](references/audit-guidelines.md): Methodology for system-level evaluation.
- [Team Context](references/team-context.md): Domain-specific standards for Payment Core.

## 🧰 Tooling Reference

- **Excel Helper**: `/product-audit/scripts/excel_helper.go`
- **Technical Advisor**: `/product-audit/scripts/advisor.go` (Internal capability for Confluence verification)
- **Confluence Fetcher**: `/product-audit/scripts/fetch_confluence.go` (Legacy)


## 📊 Audit Summary Structure

The generated text-based Audit Summary Report MUST follow this structure to match the Excel template's executive view:

1.  **AUDIT OVERVIEW**: 
    -   **Squad / Product**: (e.g., Payment Core / Normal Payment)
    -   **Auditor**: (Name and Team)
    -   **Date**: (DD MMM YYYY)
    -   **Source of Truth**: (Link/Title of the Confluence Doc)
2.  **AUDIT SCOPE**: Which product flows and user journeys were reviewed?
3.  **KEY RISKS IDENTIFIED (Flag / Fail)**: List top findings including dimension, criterion, and specific gap.
4.  **HIGH SEVERITY ITEMS (🔴)**: Items requiring immediate action.
5.  **MEDIUM SEVERITY ITEMS (🟡)**: Items to address this quarter.
6.  **LOW SEVERITY / PASS HIGHLIGHTS**: Notable strengths or minor polish items.
7.  **PROPOSED NEXT STEPS**: Owner, action, and target date.
8.  **OPEN QUESTIONS / DEPENDENCIES**: Blockers or items needing input from other teams.
9.  **AUDIT METRICS**: Final counts of Pass, Flag, Fail, and N/A.
