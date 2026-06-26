---
name: new-skill
description: Reference for writing and editing PM skills well — design, structure, and principles for predictable automation.
disable-model-invocation: true
reference: https://github.com/mattpocock/skills/tree/main/skills/productivity/writing-great-skills
---

# Writing PM Skills

A PM skill exists to wrangle **determinism** out of stochastic agents. **Predictability** — the agent taking the same _process_ every run, not producing the same output — is the root virtue.

This guide covers the anatomy of a skill, design principles, and the workflow for creating one in the Jora PM suite.

## Anatomy of a Skill

Every skill lives in `pm-skills/skills/<skill-name>/` and comprises:

### Core Files

- **`SKILL.md`** — the primary artifact. Combines:
  - **Frontmatter**: `name`, `description`, `disable-model-invocation` (optional).
  - **Body**: Steps, reference, and rules the agent follows.
- **Supporting files** (optional):
  - `GLOSSARY.md` — definitions and terminology specific to the skill.
  - `scripts/` — executables (Go binaries, shell scripts) the skill calls.
  - `templates/` — default document templates or output shapes.
  - `examples/` — sample runs or expected artifacts.

### Invocation Modes

**User-invoked skills** (`disable-model-invocation: true`):
- Only fire when the user types the skill name.
- Zero context load — description is hidden from the agent.
- Cognitive load on _you_ — you must remember it exists.
- Use when: the skill is exploratory, manual, or fired only on explicit request.

**Model-invoked skills** (omit `disable-model-invocation`):
- Agent can autonomously fire the skill based on description.
- Adds to context load — description sits in the window every turn.
- Use when: the agent must reach the skill independently or another skill must call it.

**Default for Jora**: User-invoked. PM skills are typically driver-class tools you summon on demand.

## Information Hierarchy

A skill's content lives on a ladder; each rung trades **context load** against **legwork** (digging the agent does):

1. **In-skill step** (primary tier) — ordered actions in `SKILL.md`.
   - Each step ends on a **completion criterion**: the condition that marks it done.
   - Checkable: can the agent tell done from not-done?
   - Exhaustive: "every step completed", not "make progress".

2. **In-skill reference** (secondary tier) — definitions, rules, facts consulted on demand.
   - Often flat peers: every rule of a process on one rung.
   - Flat peers are fine; they are not a smell.

3. **External reference** (tertiary tier) — pushed out to separate files via a **context pointer**.
   - Load only when the pointer fires.
   - Example: `GLOSSARY.md` for domain terminology.

**Progressive disclosure**: Keep the top legible by pushing specialist material down. Inline only what every branch needs.

## Design Principles

### Predictability

**Leading words** anchor behaviour. A single token (e.g., _relentless_, _tracer_, _red_) carries pre-trained intuition and appears consistently:

- **In the description**: guides invocation.
- **In the body**: guides execution.
- **In code comments**: guides maintenance.

Hunt for opportunities to collapse duplication into a leading word. "fast, deterministic, low-overhead" → _tight_. "a loop you can trust" → _red_.

### Pruning

Keep each meaning in a **single source of truth** — one place, so a change is one edit.

**Relevance test**: does this line still bear on what the skill does? Delete lines that fail, not rewriting.

**No-op test** (per sentence): does it change behaviour vs. the agent's default?
- "be thorough" (when the agent is already thorough) is a no-op.
- "be _relentless_" (a stronger word) is not.

### Granularity (When to Split)

Split off a new skill only when the cut earns it:

- **By invocation**: a distinct **leading word** that should trigger independently.
  - Cost: context load for the new description.
  - Worth it when: the agent must reach it on its own.
- **By sequence**: when later steps tempt the agent to rush the current one (**premature completion**).
  - Cost: cognitive load on you.
  - Worth it when: hiding post-steps encourages deeper **legwork**.

## Creating a New Skill: Workflow

### 1. Define the Scope

Answer:
- **Leading word**: What is the one-token core idea? (e.g., _audit_, _velocity_, _synthesis_).
- **Invocation**: User-invoked or model-invoked?
- **Branches**: What distinct triggering scenarios does it handle?
- **Audiences**: Who drives the skill? (PM, Eng, Data).

### 2. Draft the Body

Start with **steps** if it is a process; start with **reference** if it is a toolkit.

- **Write for predictability**: each step should be the same every run.
- **Completion criterion**: how does the agent know the step is done? Make it checkable and exhaustive.
- **External scripts**: if calling a Go binary or shell script, use `./scripts/name`.
  - Keep scripts focused on data retrieval and transformation, not reasoning.
  - Document the script signature in the skill body.

### 3. Extract Reference

Push **specialist reference** to a separate file (e.g., `GLOSSARY.md`) using a **context pointer**.

- **Context pointer**: a named link that signals when the agent should load the file.
- Example: "If you need Jira terminology, consult the [[Jira-Glossary]]" → `/GLOSSARY.md#jira-terminology`.

### 4. Test Against Failure Modes

Ask:

- **Premature completion**: Does the completion criterion invite rushing? Sharpen it.
- **Duplication**: Is a meaning stated more than once? Collapse it.
- **Sediment**: Is old material still live? Prune it.
- **Sprawl**: Is the skill too long? Disclose reference or split by branch.
- **No-op**: Does every sentence change behaviour? Delete no-ops.

### 5. Mirror to `.claude/commands/`

Copy the skill's path structure:

```
pm-skills/skills/<name>/SKILL.md
    ↓
.claude/commands/<name>.md
```

The command file is symlink or copy; keep them synchronized.

### 6. Bundle and Test

Run the bundle script:

```bash
./pm-skills/scripts/bundle_skills.sh
```

Output lands in `pm-skills/bundles/` as a `.skill` file (ZIP archive with versioning).

### 7. Stage Before Publishing

Place outputs in `pm-skills/output/sandbox/` for review, then move to `pm-skills/output/` for final publication.

## Structure: Example

For a skill named `example-skill`:

```
pm-skills/skills/example-skill/
├── SKILL.md                 # Primary artifact
├── GLOSSARY.md              # (optional) Domain terms
├── scripts/
│   ├── example.go           # Go binary
│   ├── helper.sh            # Shell helper
│   └── README.md            # Script docs
├── templates/
│   └── output_template.md   # (optional) Default template
└── examples/
    └── sample_run.md        # (optional) Expected output
```

## Zola Pay PM Suite Conventions

- **Sprint naming**: Always use **"PCDPC"** prefix for Payment Core (e.g., "PCDPC - Sprint 26.06.A").
- **Spaces**: Priority Confluence spaces: "Zalopay Technology Management", "Payment Engine".
- **Keywords**: High-signal search terms: `FMC/PE`, `AEv2`, `ZAS`, `Idempotency`, `ACv2`.
- **Tone**: Professional, direct, concise. Australian English (summarise, prioritise, centralised).
- **Output**: Always stage in `pm-skills/output/sandbox/` before final publication.

## Failure Modes Checklist

- [ ] **Premature completion**: Completion criterion is vague or invites rushing. Sharpen it.
- [ ] **Duplication**: Same meaning stated twice. Collapse into a leading word or single source.
- [ ] **Sediment**: Old material no longer relevant. Prune it.
- [ ] **Sprawl**: Skill too long, even if every line is live. Disclose reference or split.
- [ ] **No-op**: Sentences that don't change agent behaviour. Delete them.

## References

- **Original skill-writing guide**: [mattpocock/skills](https://github.com/mattpocock/skills/tree/main/skills/productivity/writing-great-skills) — the foundation for this skill.
- **Jora PM Skills**: `pm-skills/README.md` — project conventions and tooling.
- **Claude Code**: [claude-code-guide skill](/) — agent SDK and automation patterns.
