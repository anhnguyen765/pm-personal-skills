---
name: grilling
description: Interview the user relentlessly about a plan or design. Use when the user wants to stress-test a plan before building, to validate assumptions, or to strengthen a memo or PRD.
---

# Grilling

Interview the user relentlessly about every aspect of a plan, design, or proposal until we reach a shared understanding. Walk down each branch of the design tree, resolving dependencies between decisions one-by-one.

## Core Rules

1. **One question at a time** — Asking multiple questions at once is overwhelming. Wait for feedback on each before continuing.
2. **Recommended answers** — For each question, provide your best-guess answer. This anchors the conversation and speeds up rounds.
3. **Explore before asking** — If a question can be answered by exploring the codebase, Jira, or Confluence, do that instead of asking.
4. **Relentless** — Keep going until the plan is stress-tested or assumptions are validated. Don't accept hand-wavy answers.

## Workflow

### 1. Scope the Design
Understand what is being proposed: a feature, technical decision, operational change, memo argument, or PRD structure.

### 2. Identify the Decision Tree
Map the branches: what decisions drive the design? What depends on what?
- Core decisions (e.g., "Do we use event sourcing?").
- Derivative decisions (e.g., "How do we handle idempotency?").

### 3. Walk the Tree — One Question at a Time

**Question format**:
```
Q: [Question that exposes a decision or assumption]

My take: [Your recommended answer based on ZaloPay context]

Why I ask: [Why this decision matters to the overall design]
```

**Example**:
> **Q**: If we switch to event sourcing, how do we handle the 10-year audit trail for legacy transactions?
> 
> **My take**: Dual-write system: legacy transactions stay in the old store; new transactions go to the event store.
> 
> **Why I ask**: Audit compliance is a Payment Core requirement — we can't lose history.

### 4. Resolve Conflicts
When answers conflict with prior decisions or constraints, call it out:
> That conflicts with [prior decision]. Do we need to reconsider [decision X]?

### 5. Stop When
- All major branches are stress-tested.
- Assumptions are validated or revised.
- The user has clear confidence in the design.
- Or the user says "stop" — they own the pacing.

## Context: ZaloPay Payment Core

Use these anchors when grilling:

| Anchor | Detail |
|---|---|
| **Services** | AC (Account Core), AE (Authorization Engine), PA (Payment Aggregator), ZAS (Zalopay Auth Service), RC (Reconciliation Core) |
| **Constraints** | Idempotency, audit trails, cross-border compliance, ledger consistency |
| **High-stakes** | Event sourcing, eventual consistency, state machine design, failure modes |
| **Reference** | Confluence: "Zalopay Technology Management", "Payment Engine" |

## Example Questions

**Design question** (tech debt):
> **Q**: If we deprecate UPM for AEv2, how do we handle in-flight transactions that reference UPM?
> 
> **My take**: Run both systems in parallel for 6 months, with fallback to UPM if AEv2 fails. Then sunset UPM.
>
> **Why I ask**: Deprecation strategy is make-or-break for rollout safety.

**Assumption question** (business logic):
> **Q**: You're assuming all merchants accept the new idempotency key format. What if a legacy integration rejects it?
> 
> **My take**: Gate behind a feature flag; roll out to new merchants first.
>
> **Why I ask**: Compatibility is a hidden cost of the change.

**Edge case question** (failure mode):
> **Q**: What happens if the idempotency service is down during a payment request?
> 
> **My take**: Degrade gracefully: accept the request without deduplication and flag it for manual review.
>
> **Why I ask**: Availability beats perfect deduplication.

## Tips

- **Be relentless**: If the answer is vague ("it'll be fine"), keep pressing.
- **Maintain the tree**: Don't jump between unrelated topics — stay on branches.
- **Reference the docs**: Link to Confluence pages or Jira decisions when relevant.
- **Record outcomes**: Suggest the user update the memo or PRD with resolved decisions.

## When Grilling Triggers

- **During `/memo`**: After drafting, grill the core assumptions.
- **During `/prd`**: After drafting, grill the technical design and edge cases.
- **Explicit**: User types `/grill-me` or says "grill me" or "stress-test this".

## See Also

- `/grill-me` — Quick trigger for a grilling session.
- `/memo` — Decision memo drafting (includes optional grilling).
- `/prd` — PRD drafting (includes optional grilling).
