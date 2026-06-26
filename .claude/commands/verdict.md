---
name: verdict
description: Executive product brief — synthesise input signals and metrics into an execution overview and strategic direction assessment for PO and leadership
disable-model-invocation: true
---

# Verdict: Execution & Strategic Review Brief

**Leading word**: _Verdict_ — distil your input (narrative + metrics) into a balanced executive brief covering execution health and strategic direction.

Use this skill when:
- **Weekly/sprint leadership sync**: Present product execution status + forward strategy.
- **Monthly/quarterly review**: Reflect on performance + replan next priority wave.
- **Executive handoff**: Brief PO and leadership on what's working, what's breaking, and where we're heading.

Output: Concise, structured brief (2–3 pages) in PM voice.

---

## Process

### 1. Collect Your Input
**Done when**: You have provided free-form narrative + key metrics for this review period.

Provide (in any format—bullets, paragraphs, mixed):
- **Execution signals**: Velocity trend, cycle time, quality metrics (incident count, escape rate), deployment frequency.
- **User/market signals**: Support volume, NPS trend, customer feedback themes, competitive movement.
- **Roadmap/strategy signals**: What shipped vs. plan, blockers, capacity constraints, strategic bets in flight.
- **Team signals**: Capacity changes, knowledge gaps, process friction, morale notes.

_Example input you might provide_:
```
- Velocity down 12% last sprint; 2 devs out for onboarding
- Critical incidents ↑2 in reconciliation; root cause = stale caching strategy
- Payment validation cycle time ↓8% after parallelisation work
- Customer support: +15% volume, 60% billing-related (strategy shift impact?)
- Q3 roadmap: 70% shipped, 2 features pushed to Q4 (vendor integration scope creep)
- Team: New AC lead starting mid-sprint; PA team at capacity ceiling
```

### 2. Parse Into Execution + Strategy Layers
**Done when**: You've separated _how we're executing_ from _where we're heading_.

**Execution layer** (health of current delivery):
- Velocity, quality, incident trends, deployment cadence.
- Capacity constraints (people, infra, process).
- Technical debt or friction blocking speed.

**Strategy layer** (progress toward goals):
- Roadmap completion vs. plan.
- Market signals (customer demand, competitive moves).
- Strategic bets: are they unblocking velocity or creating it?

_Example parse_:
```
EXECUTION: Velocity down, but driven by onboarding (temporary). Quality concern in Reconciliation (structural).
STRATEGY: 70% Q3 roadmap shipped; vendor integration slip impacts market timing. Payment validation gains compound quarterly goal.
```

### 3. Render Health & Direction
**Done when**: You can summarise in 1–2 sentences: execution status + strategic bearing.

Execution summary (up/stable/down + primary reason):
```
Execution is _stable with localized friction_: core velocity sustained despite onboarding load; 
quality spike in Reconciliation is addressable within current sprint.
```

Strategy summary (on-track/at-risk + why):
```
Strategic: Q3 roadmap 70% complete; vendor integration slip moderates Q4 market entry but 
payment-validation gains position us ahead on core metric.
```

### 4. Structure for PO + Leadership
**Done when**: You've mapped your input to: signals, why it matters, and what's next.

Think of:
- **Signals**: What's working (green), what's breaking (red), what's unclear (yellow). Adapt format to review type (bullets, narrative, combo).
- **Why**: Root causes, second-order effects, capacity constraints. Market impact. Why should leadership care _now_?
- **Next**: 3–5 actions ranked by unblock + strategy progress. Owner, criterion, target date.

### 5. Write the Brief
**Done when**: You have a 2–3 page document in PM voice, ready for PO + leadership circulation.

**Guardrails** (not a template—structure as needed, but include these):

1. **Open with a clear judgment** (1–2 sentences): Execution status + strategic bearing. This is the headline.
2. **Surface signals clearly**: What's working (green), what's breaking (red), what's unclear (yellow). Use bullets or narrative; adapt to your review type.
3. **Explain why it matters** (1 paragraph): Root causes, second-order effects, capacity constraints, market impact. Why should they care?
4. **State next priorities** (3–5 actions, ranked): Owner, success criterion, target date. Rank by unblock + strategy progress.
5. **Close with confidence note**: Data quality and known gaps (e.g., "Based on sprint metrics; NPS data pending").

**Length**: 2–3 pages. If context requires more, move specialist detail to appendix or Confluence link.

**Tone**: Direct PM voice. Lead with what matters to PO + leadership (market timing, capacity, risk). Balance candour with optimism.

---

## Output

Verdict artifact lands in `pm-skills/output/verdict_[sprint]_[date].md`.

Example:
```
pm-skills/output/verdict_PCDPC-26.06.A_2026-06-26.md
```

---

## Notes

- **Audience**: PO, leadership team. Assume tech fluency but no sprint-level detail.
- **Frequency**: Weekly (during sprints), monthly (strategic review), or on-demand (crisis/opportunity).
