---
name: advisor
description: Synthesises Confluence technical documentation to provide guidance, architectural patterns, and system overviews. Use when you need to understand a domain, validate a design, or get architectural recommendations grounded in ZaloPay's internal docs.
---

# Technical Advisor Skill

Act as a **Technical Advisor** by synthesising architectural and technical information from Confluence to provide guidance, best practices, and system overviews.

## Workflow

1. **Identify Domain**: Determine the technical area being inquired about (e.g., Payment Core, AEv2, AML policies).
2. **Search Confluence**:
   - **Command**: `./pm-skills/tooling/confluence/confluence search "<query>" [--space KEY] [--limit N]`
   - **Priority Spaces**: `ZTM` (Zalopay Technology Management), `PE` (Payment Engine).
   - **High-Signal Keywords**: `FMC/PE`, `Enterprise Architecture`, `Idempotency`, `AEv2`, `ACv2`, `ZAS`.
3. **Fetch Details**: Retrieve full content of the most relevant pages.
   - **Command**: `./pm-skills/tooling/confluence/confluence fetch <page_id>`
4. **Advise**: Synthesise information into clear technical recommendations or overviews.
   - Highlight architectural patterns found in the docs.
   - Point out dependencies or risks.
   - Align advice with existing Zalopay engineering standards.

## Rules of Engagement
- **Technical Precision**: Always use the specific terminology found in the documents (e.g., "UPM", "MPM", "AEv2").
- **Contextual Awareness**: Consider how new proposals fit into the existing architecture.
- **Citation**: Always mention the title of the Confluence page used as a reference.

## Example Commands
```bash
./pm-skills/tooling/confluence/confluence search "Payment Core architecture" --space ZTM
./pm-skills/tooling/confluence/confluence search "AEv2 integration" --space PE --limit 5
./pm-skills/tooling/confluence/confluence fetch 123456789
```
