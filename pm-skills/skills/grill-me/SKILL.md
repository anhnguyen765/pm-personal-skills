---
name: grill-me
description: A relentless interview to sharpen a plan or design. Stress-test assumptions before finalising a memo or PRD.
disable-model-invocation: true
reference: https://github.com/mattpocock/skills/blob/main/skills/productivity/grill-me/SKILL.md
---

# Grill Me

Run a `/grilling` session to relentlessly interview yourself about your plan, design, or proposal.

## When to Use

- **Before finalising a memo**: You've drafted the argument; grill the assumptions.
- **Before shipping a PRD**: You've drafted the design; stress-test it.
- **To validate a risky decision**: You're uncertain about a choice; get grilled on it.
- **Anytime you want a second opinion**: Use grilling to sharpen your thinking.

## What Happens

The `/grilling` skill will:

1. **Understand your proposal**: Ask you to briefly outline what you want grilled (a memo, a PRD section, a technical design).
2. **Map the decision tree**: Identify the core decisions and their dependencies.
3. **Ask one question at a time**: Wait for your answer before moving to the next.
4. **Provide recommendations**: For each question, suggest an answer based on ZaloPay Payment Core context.
5. **Keep pressing**: Don't accept hand-wavy answers. The goal is confidence, not just completion.

## Example Flow

You've drafted a memo arguing for adopting event sourcing. You type `/grill-me`.

> **Grilling**: What's your proposal?
> 
> **You**: Event sourcing for all Payment Core state machines.
> 
> **Grilling**: Q: If we switch to event sourcing, how do we handle the 10-year audit trail for legacy transactions that predate the system?
> 
> **My take**: We maintain a dual-write system: legacy transactions live in the old store; new transactions go to the event store. After 10 years, we can archive.
> 
> **Why I ask**: Audit compliance is a Payment Core requirement.
>
> **You**: Hmm, we hadn't thought about that. We'd need to run both systems...

And so on, until the memo is stress-tested.

## Integration with Memo and PRD

Both `/memo` and `/prd` workflows can trigger grilling at the end:

- **After `/memo`**: "Do you want me to grill your memo's assumptions?"
- **After `/prd`**: "Do you want me to grill the technical design?"

You can opt in or out — grilling is optional but recommended before finalising.

## Tips

- **Be honest**: If you don't know the answer, say so. Grilling exposes gaps.
- **Push back**: If you disagree with the recommended answer, defend your position. That's the point.
- **Take notes**: Jot down resolved decisions and revised assumptions. Update your memo or PRD.
- **Know when to stop**: If you're confident and all major branches are stress-tested, you're done.

## See Also

- `/grilling` — the core grilling skill (runs automatically; you don't invoke it directly).
- `/memo` — drafting decision memos (includes optional grilling step).
- `/prd` — drafting technical PRDs (includes optional grilling step).
