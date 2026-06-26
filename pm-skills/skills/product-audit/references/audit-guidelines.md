# Product Audit Guidelines

## 1. Purpose
The Product Audit skill supports systematic evaluation of product quality, operational readiness, and compliance with internal standards.

## 2. Audit Execution Workflow

### Step 1: Initialize
- Locate the template: `Add local storage path`.
- Copy it to the workspace for processing if needed.
- Use `scripts/excel_helper.go` to read the instructions and checklist.

### Step 2: Evaluation (for each row)
- **Understand the Rule**: Read Description, What to Check, and Failure Example.
- **Collect Evidence**: Review documentation, screenshots, or live behavior.
- **Technical Validation**: **(MANDATORY)** Cross-reference collected evidence with Confluence using the `technical-advisor` skill to ensure alignment with architectural designs.
- **Assign Status**: `Pass`, `Flag`, `Fail`, or `N/A`.
- **Document Findings**: Record evidence, findings, and severity in the specified columns. Any misalignment between design and implementation must be flagged.

### Step 3: Reporting
- **Populate Excel**: Use `excel_helper.go` in `write` mode to fill the results into the spreadsheet.
- **Generate Summary**: Create a structured report following the "Summary" tab template.

## 3. Checklist Definition (Columns F-H)
- **Description**: The concept or rule.
- **What to Check**: Specific verification steps.
- **Failure Example**: Patterns that indicate a fail status.

## 4. Result Fields (Columns I-K)
- **Status (I)**: The outcome (`Pass`, `Flag`, `Fail`, `N/A`).
- **Evidence & Findings (J)**: Combined technical evidence and findings discovered during the audit.
- **Severity (K)**: Impact level (`Low`, `Medium`, `High`).
