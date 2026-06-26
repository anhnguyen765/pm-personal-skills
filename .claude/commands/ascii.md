---
name: ascii
description: Create professional ASCII diagrams for system architecture and component interactions. Use when the user requests a high-level overview of system flows, API sequences, or architectural component interactions in plain-text format.
---

# SA ASCII Diagrams

Provides a standardised way to draw high-signal, professional ASCII diagrams that mimic Solution Architecture (SA) interaction diagrams.

## Core Principles

1. **High Signal, Low Noise**: Focus on critical components and data flows. Avoid unnecessary decorative elements.
2. **Standardised Component Shapes**:
   - `[ Component ]`: Standard service or application.
   - `( Database )`: Persistence layers.
   - `<  Broker  >`: Message queues or event brokers.
   - `[[ External ]]`: Third-party services or external systems.
3. **Directional Flow**: Use clear arrows (`--->`, `<---`, `<==>`) to indicate direction of data or control flow.
4. **Numbered Sequences**: Use `1. Action`, `2. Result` for sequential interactions.

## Diagram Patterns

### 1. Component Interaction (Box & Arrow)
Best for showing static relationships and high-level data flow.

```text
+----------------+       +-----------------+       +-----------------+
|   Client UI    | ----> |  API Gateway    | ----> |  Payment Core   |
|   (Mobile/Web) | <---- |   (Nginx/Kong)  | <---- |     (FMC)       |
+----------------+       +--------+--------+       +--------+--------+
                                  |                         |
                                  |                 +-------+-------+
                                  |                 |      PE       |
                                  +---------------> | (Idempotency) |
                                                    +-------+-------+
```

### 2. Sequence Interaction (Timeline)
Best for showing step-by-step logic and time-ordered events.

```text
User           FMC (Core)          PE (Engine)           ZAS (Accounting)
 |                 |                   |                      |
 | 1. Confirm Pay  |                   |                      |
 |---------------->|                   |                      |
 |                 | 2. Check Dup      |                      |
 |                 |------------------>|                      |
 |                 | 3. Record Journal |                      |
 |                 |----------------------------------------->|
 |                 | <----------------------------------------|
 |                 | 4. Success        |                      |
 |<----------------|                   |                      |
```

## Workflow Guidance

1. **Identify Actors**: List all internal services, external APIs, and databases involved.
2. **Define the Trigger**: Start with the user or system event that initiates the flow.
3. **Trace the Path**: Follow the primary "happy path" first, then add critical error branches.
4. **Annotate**: Label every arrow with the action or data being transferred.
