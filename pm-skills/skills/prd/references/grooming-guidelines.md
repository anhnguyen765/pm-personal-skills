# Sprint Grooming & Planning Guidelines

Use the PRD to drive technical alignment during grooming sessions.

## 1. Breakdown Strategy
- **Task Granularity**: Each 'Functional Requirement' should map to 1-3 Jira tasks.
- **Dependency Mapping**: Ensure 'Internal System' dependencies are identified and tickets created in respective backlogs.

## 2. Grooming Checklist
- [ ] **Idempotency**: Is the `request_id` or `partner_id` handling clear for all APIs?
- [ ] **State Machine**: Are all transition states (e.g., `PENDING` -> `SUCCESS`, `PENDING` -> `FAIL`) covered?
- [ ] **Timeout Handling**: What happens if the Partner Gateway doesn't respond in 15s?
- [ ] **Reconciliation**: Does the ledger entry match the gateway status?

## 3. Estimating
- Focus on the 'Risks & Mitigations' section to add buffer for high-risk technical integrations.
- Use 'Non-Functional Requirements' to define 'Definition of Done' (DoD) performance benchmarks.
