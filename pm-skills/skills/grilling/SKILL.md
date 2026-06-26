---
name: grilling
description: Interview the user relentlessly about a plan or design. Use when the user wants to stress-test a plan before building, to validate assumptions, or to strengthen a memo or PRD.
reference: https://github.com/mattpocock/skills/blob/main/skills/productivity/grilling/SKILL.md
---

# Grilling Skill

Interview the user relentlessly about every aspect of a plan, design, or proposal until we reach a shared understanding. Walk down each branch of the design tree, resolving dependencies between decisions one-by-one.

## Core Principles

- **One question at a time**: Asking multiple questions at once is overwhelming. Wait for feedback on each before continuing.
- **Recommended answers**: For each question, provide your best-guess answer. This anchors the conversation and speeds up rounds.
- **Explore before asking**: If a question can be answered by exploring the codebase, Jira, or Confluence, do that instead of asking.
- **Dependency-first**: Questions should expose hidden assumptions and decision dependencies, not just gather facts.
- **Relentless**: Keep going until the plan is stress-tested or assumptions are validated. Don't accept hand-wavy answers.

## Workflow

### 1. **Scope the Design**
   Understand what is being proposed: a feature, a technical decision, an operational change, a memo argument, a PRD structure.

### 2. **Identify the Decision Tree**
   Map the branches: what decisions drive the design? What depends on what?
   - Core decisions (e.g., "Do we use event sourcing?").
   - Derivative decisions (e.g., "How do we handle idempotency?" depends on event sourcing).

### 3. **Walk the Tree**
   For each node, ask one question. Wait for the answer. Ask the next.
   
   **Question format**:
   > **Q**: [Question that exposes a decision or assumption]
   > 
   > **My take**: [Your recommended answer based on ZaloPay context, Payment Core constraints, or engineering standards]
   >
   > **Why I ask**: [Why this decision matters to the overall design]

### 4. **Resolve Conflicts**
   When answers conflict with prior decisions or constraints, call it out:
   > That conflicts with [prior decision]. Do we need to reconsider [decision X]?

### 5. **Stop When**
   - All major branches are stress-tested.
   - Assumptions are validated or revised.
   - The user has clear confidence in the design.
   - Or the user says "stop" — they own the pacing.

## Context: ZaloPay Payment Core

Use these anchors when grilling:

- **Payment Core services**: AC (Account Core), AE (Authorization Engine), PA (Payment Aggregator), ZAS (Zalopay Auth Service), RC (Reconciliation Core).
- **Key constraints**: Idempotency, audit trails, cross-border compliance, ledger consistency.
- **High-stakes decisions**: Event sourcing, eventual consistency, state machine design, failure modes.
- **Reference**: Confluence spaces "Zalopay Technology Management" and "Payment Engine".

## When to Trigger

- **During `/memo`**: After the user drafts a decision memo, grill the core assumptions.
- **During `/prd`**: After the user drafts a PRD, grill the technical design and edge cases.
- **Explicit request**: User types `/grill-me` or says "grill me" or "stress-test this".

## Example Questions

**Design question** (tech debt):
> **Q**: If we deprecate UPM for AEv2, how do we handle in-flight transactions that reference UPM?
> 
> **My take**: We run both systems in parallel for 6 months, with a fallback to UPM if AEv2 fails. Then sunset UPM.
>
> **Why I ask**: Deprecation strategy is make-or-break for rollout safety.

**Assumption question** (business logic):
> **Q**: You're assuming all merchants accept the new idempotency key format. What if a legacy integration rejects it?
> 
> **My take**: We gate behind a feature flag and only roll out to new merchants first.
>
> **Why I ask**: Compatibility is a hidden cost of the change.

**Edge case question** (failure mode):
> **Q**: What happens if the idempotency service is down during a payment request?
> 
> **My take**: We degrade gracefully: accept the request without deduplication and flag it for manual review.
>
> **Why I ask**: Availability beats perfect deduplication; manual review catches real duplicates.

## Failure Modes

- **Accepting hand-wavy answers**: "It'll be fine" is not an answer. Keep pressing.
- **Losing the tree**: Jumping between unrelated topics. Maintain a mental map of dependencies.
- **Running out of time**: If the user is tired, summarize and suggest a follow-up session.
- **Assuming you know the answer**: Ask; don't tell. The user's answer may surprise you.

## Notes

- **Reference the docs**: If a question points to a Confluence page or Jira decision, link it.
- **Escalate if needed**: If a question reveals a blocker (e.g., "legal hasn't reviewed this"), note it and suggest a follow-up.
- **Record outcomes**: Suggest the user capture resolved decisions in the original memo or PRD.
