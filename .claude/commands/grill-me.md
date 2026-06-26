---
name: grill-me
description: A relentless interview to sharpen a plan or design. Stress-test assumptions before finalising a memo or PRD.
---

# Grill Me

Run a `/grilling` session to relentlessly interview yourself about your plan, design, or proposal.

## When to Use

| Scenario | Use Grilling |
|---|---|
| **Before finalising a memo** | You've drafted the argument; stress-test the assumptions. |
| **Before shipping a PRD** | You've drafted the design; validate the technical approach. |
| **To validate a risky decision** | You're uncertain about a choice; get grilled on it. |
| **Anytime you want a second opinion** | Use grilling to sharpen your thinking before committing. |

## What Happens

The grilling session will:

1. **Understand your proposal**: Briefly outline what you want grilled.
2. **Map the decision tree**: Identify core decisions and their dependencies.
3. **Ask one question at a time**: Wait for your answer before moving to the next.
4. **Provide recommendations**: For each question, suggest an answer based on ZaloPay Payment Core context.
5. **Keep pressing**: Don't accept hand-wavy answers. Goal is confidence, not just completion.

## Example: Grilling a Memo

You've drafted a memo arguing for adopting event sourcing. You type `/grill-me`.

```
Grilling: What's your proposal?

You: Event sourcing for all Payment Core state machines.

Grilling: Q: If we switch to event sourcing, how do we handle the 10-year 
audit trail for legacy transactions that predate the system?

My take: Dual-write system. Legacy transactions stay in the old store; 
new transactions go to the event store. After 10 years, archive.

Why I ask: Audit compliance is a Payment Core requirement.

You: Hmm, we hadn't thought about that. We'd need to run both systems...

Grilling: Exactly. That's a 6-month parallel run, at minimum. How does 
that affect the launch timeline you proposed?
```

And so on, until the memo is stress-tested.

## Integration with `/memo` and `/prd`

Both workflows can trigger grilling at the end:

- **After `/memo`**: "Do you want me to grill your memo's assumptions?"
- **After `/prd`**: "Do you want me to grill the technical design?"

You can opt in or out — grilling is optional but recommended before finalising.

## Tips

- **Be honest**: If you don't know the answer, say so. That's the point of grilling.
- **Push back**: If you disagree with the recommended answer, defend your position.
- **Take notes**: Jot down resolved decisions and revised assumptions. Update your memo or PRD.
- **Know when to stop**: Once you're confident and major branches are stress-tested, you're done.

## Flow

1. Invoke `/grill-me` (or the grilling session happens automatically in `/memo` or `/prd`).
2. Describe your proposal in 1–2 sentences.
3. Grilling asks one question. Answer.
4. Grilling may follow up or move to the next branch.
5. Repeat until done.
6. Update your memo or PRD with the insights you gained.

## See Also

- `/grilling` — The core grilling skill; runs automatically during grilling sessions.
- `/memo` — Decision memo drafting (includes optional grilling integration).
- `/prd` — PRD drafting (includes optional grilling integration).
