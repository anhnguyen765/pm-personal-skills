---
name: advisor
description: Synthesises Confluence technical documentation to provide guidance, architectural patterns, and system overviews.
---

# Technical Advisor Skill

This skill acts as a Technical Advisor by synthesising architectural and technical information from Confluence to provide guidance, best practices, and system overviews.

## Workflow

1.  **Identify Domain**: Determine the technical area the user is inquiring about (e.g., Payment Core, AEv2, AML policies).
2.  **Search Confluence**: Use the advisor script to find relevant technical documentation.
    -   `./scripts/advisor search "topic name"`
    -   **Priority Spaces**: "Zalopay Technology Management", "Payment Engine".
    -   **High-Signal Keywords**: `FMC/PE`, `Enterprise Architecture`, `Idempotency`, `AEv2`.
3.  **Fetch Details**: Retrieve the full content of the most relevant pages.
    -   `./scripts/advisor fetch <page_id>`
4.  **Advise**: Synthesise the information to provide a clear technical recommendation or overview.
    -   Highlight architectural patterns found in the docs.
    -   Point out dependencies or risks.
    -   Align advice with existing Zalopay engineering standards.

## Rules of Engagement

-   **Technical Precision**: Always use the specific terminology found in the documents (e.g., "UPM", "MPM", "AEv2").
-   **Contextual Awareness**: Consider how new proposals fit into the existing architecture described in Confluence.
-   **Citation**: Always mention the title of the Confluence page used as a reference.

## Example Commands

-   `go run scripts/advisor.go search "Payment Core architecture"`
-   `go run scripts/advisor.go fetch 12345678`
