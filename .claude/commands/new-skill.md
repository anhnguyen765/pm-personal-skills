---
name: new-skill
description: Reference for writing and editing PM skills well — design, structure, and principles for predictable automation.
---

# Writing PM Skills

A PM skill is a deterministic automation tool that does the same _process_ every run. This guide covers the anatomy, design principles, and workflow for creating one in the Jora PM suite.

## Skill Anatomy

Every skill lives in `pm-skills/skills/<skill-name>/` with:

- **`SKILL.md`** — the primary artifact (frontmatter + body with steps or reference).
- **`GLOSSARY.md`** (optional) — definitions and domain terminology.
- **`scripts/`** (optional) — Go binaries or shell helpers the skill calls.
- **`templates/`** (optional) — default output templates or shapes.
- **`examples/`** (optional) — sample runs or expected artifacts.

## Invocation Modes

| Mode | Use When | Cost |
|---|---|---|
| **User-invoked** (`disable-model-invocation: true`) | Skill fires only on explicit `/skill-name` command. Exploratory or manual skills. | Zero context load; cognitive load on you (must remember it exists). |
| **Model-invoked** (default) | Agent can autonomously fire based on description. Another skill can reach it. | Context load (description in window every turn); agent reaches it reliably. |

**Default for Jora**: User-invoked. PM skills are typically driver-class tools you summon on demand.

## Information Hierarchy

Keep the top legible by moving specialist material down:

1. **In-skill step** (primary) — ordered actions in `SKILL.md` with checkable completion criterion.
2. **In-skill reference** (secondary) — definitions and rules consulted on demand.
3. **External reference** (tertiary) — pushed to `GLOSSARY.md` or linked files via context pointers.

## Design Principles

### Predictability → Leading Words

A single strong token (e.g., _relentless_, _audit_, _synthesis_) anchors behaviour:

- In the description: guides invocation.
- In the body: guides execution.
- In code: guides maintenance.

Collapse duplication into a leading word: "fast, deterministic, low-overhead" → _tight_.

### Pruning

**Relevance test**: Delete lines that don't bear on what the skill does.

**No-op test** (sentence-by-sentence): Does it change behaviour vs. the agent's default?
- "be thorough" is often a no-op.
- "be _relentless_" is not.

### Granularity (When to Split)

Split a new skill only when:

- **By invocation**: A distinct leading word that should trigger independently (costs context load for description).
- **By sequence**: Later steps tempt rushing the current one; hiding them encourages deeper work (costs cognitive load on you).

## Creating a New Skill: Workflow

### 1. Define Scope

- **Leading word**: One-token core idea (e.g., _audit_, _velocity_, _synthesis_).
- **Invocation**: User-invoked or model-invoked?
- **Branches**: Distinct triggering scenarios?
- **Audience**: PM, Eng, Data?

### 2. Draft the Body

**For a process**: Start with **steps**, each with a checkable completion criterion.

**For a toolkit**: Start with **reference** (rules, definitions, facts).

- Write for predictability: same process, every run.
- Completion criterion: how does the agent know it's done? Make it exhaustive.
- External scripts: Use `./scripts/name`; keep them focused on data, not reasoning.

### 3. Extract Reference

Push specialist material to `GLOSSARY.md` or a linked file using **context pointers**:

```markdown
For Jira terminology, see [[Jira-Glossary]] → GLOSSARY.md#jira-terminology
```

### 4. Test Against Failure Modes

- **Premature completion**: Does the criterion invite rushing? Sharpen it.
- **Duplication**: Same meaning stated twice? Collapse it.
- **Sediment**: Old material still live? Prune it.
- **Sprawl**: Skill too long? Disclose reference or split by branch.
- **No-op**: Does every sentence change behaviour? Delete it.

### 5. Mirror to `.claude/commands/`

Synchronise the command file with the skill:

```
pm-skills/skills/<name>/SKILL.md → .claude/commands/<name>.md
```

Update the CLAUDE.md instructions table if needed.

### 6. Bundle and Test

```bash
./pm-skills/scripts/bundle_skills.sh
```

Output lands in `pm-skills/bundles/` (`.skill` files = ZIP archives).

### 7. Stage Before Publishing

- Review in `pm-skills/output/sandbox/`.
- Move to `pm-skills/output/` for publication.

## Jora PM Suite Conventions

- **Sprint naming**: Use **"PCDPC"** prefix (e.g., "PCDPC - Sprint 26.06.A").
- **Confluence spaces**: Priority: "Zalopay Technology Management", "Payment Engine".
- **Keywords**: `FMC/PE`, `AEv2`, `ZAS`, `Idempotency`, `ACv2`.
- **Tone**: Professional, direct, concise. Australian English (summarise, prioritise, centralised).
- **Output**: Stage in `pm-skills/output/sandbox/` before final publication.

## Example Skill Structure

```
pm-skills/skills/example-skill/
├── SKILL.md                 # Primary artifact
├── GLOSSARY.md              # Domain terms
├── scripts/
│   ├── example.go           # Go binary
│   ├── helper.sh            # Shell helper
│   └── README.md            # Script docs
├── templates/
│   └── output_template.md   # Default template
└── examples/
    └── sample_run.md        # Expected output
```

## Quick Checklist

- [ ] Leading word chosen and used consistently.
- [ ] Completion criteria are checkable and exhaustive.
- [ ] Specialist reference extracted to `GLOSSARY.md` or linked file.
- [ ] Failure modes tested (premature completion, duplication, sediment, sprawl, no-op).
- [ ] Command file mirrored to `.claude/commands/`.
- [ ] Bundled and tested in `pm-skills/output/sandbox/`.
- [ ] Output documentation updated in `pm-skills/README.md`.

## References

- **Original skill-writing guide**: [mattpocock/skills — Writing Great Skills](https://github.com/mattpocock/skills/tree/main/skills/productivity/writing-great-skills)
- **Jora PM Skills**: `pm-skills/README.md`
- **Full skill reference**: `pm-skills/skills/new-skill/SKILL.md`
