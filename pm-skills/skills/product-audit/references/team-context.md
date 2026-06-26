# Team Context: Payment Core (FMC & PE) & Cross-border

This document defines the specific domain context and standards for audits performed by the Payment Core team.

## 1. Product & Domain Definition
- **The Product**: **Payment Core (FMC & PE)**.
- **FMC (Fund Movement Core)**: The core orchestrator that handles **Acquiring Core (AC)** and **Refund Core**. It is responsible for order creation, refund logic, and managing product-level payment requests.
- **PE (Payment Engine)**: The technical execution engine beneath FMC, responsible for transaction orchestration, idempotency, and interfacing with Asset/Accounting layers.
- **Internal Users (Related Services/Layers)**:
    - **Product Layers**: Merchant Server, Zalopay App (Client), and Mini-apps.
    - **AE (Asset Exchange)**: The service responsible for fund movement and source-of-fund management.
    - **ZAS (ZaloPay Accounting System)**: The financial source of truth. Responsible for all **ledger entries**, accounting-related operations (**hạch toán**), and final financial reconciliation.
- **External Users**: Merchant platforms and third-party technical systems.

## 2. System-as-User Philosophy
Every audit must be conducted through the lens of the **System-as-User**. In this perspective:
- **Users** are the Product Layers (AC, Refund Core, Merchant Server) that 'consume' FMC/PE APIs.
- **UX Labels & Copy** translate to **API Parameters, Response Schemas, and Documentation**.
- **CTA (Call to Action)** translates to **Final API Execution** (e.g., the `Confirm` or `Commit` call).
- **Navigation** translates to **Service Orchestration & Sequence Flows**.
- **Accessibility** translates to **API Discoverability, Latency, and Availability**.

## 3. Technical Source of Truth (System-Heavy Context)
... (rest of the file remains same) ...


### Mandatory Technical References:
- **Zalopay Developer Centre (v2)**: [https://developers.zalopay.vn/v2/](https://developers.zalopay.vn/v2/) - Primary reference for Corev2 integrations.
- **Zalopay Developer Centre (v1)**: [https://developers.zalopay.vn/v1/](https://developers.zalopay.vn/v1/) - Reference for legacy and parallel v1 integrations.
- **Zalopay Technology Management (ZTM):** Confluence space for architecture blueprints and engineering standards.
- **Payment Engine Technical Space:** Confluence space for detailed system design and sequence diagrams.
- **Zalopay Merchant Technology:** Confluence space for detailed Cross-border Back-end logic.

### High-Signal Keywords:
- `FMC`, `PE`, `Fund Movement Core`, `Payment Engine`, `Idempotency`, `AEv2`, `Corev2`, `ACv2`, `Asset Exchange`, `ZAS`, `TDCore`.

## 4. Team-Specific Standards
- **Idempotency:** All API endpoints must implement idempotency at all layers using `requestId`.
- **Param Integrity:** Parameter naming must be consistent across domains (e.g., avoid mixing `sof_id` and `asset_token` without a justified mapping layer).
- **Orchestration Safety:** Every 'CTA' (Confirm/Execute) must be preceded by a freshness check (Stale-check) using the `Adjusted` (Hiệu chỉnh) timestamp.
- **Data Integrity:** Cross-service reconciliation is mandatory; FMC ledger entries must perfectly sync with Asset Exchange fund movements and ZAS postings.
- **Error Quality**: Responses must provide an `InstructionCode` or `ResolutionHint` for calling services to decide (Retry/Fail/Check History).

## 5. Detailed Technical Audit Considerations (System-as-User Lens)
These criteria must be cross-referenced for **every** audit to ensure system-level integrity.

| Dimension | System-as-User Interpretation |
| :--- | :--- |
| **Naming & Labels** | **Param Integrity:** Verify consistency of parameter names across ACv2, PE, and AE. Ensure naming matches domain standards to prevent integration friction. |
| **Error Quality** | **Communication Quality:** Review gRPC error codes. Responses must provide actionable guidance ("Hướng xử lý") for calling services. |
| **System Safeguards** | **Orchestration Safety:** Evaluate protection against: Duplicate requests (Fundloss), Peak traffic (System down), and Invalid operations (Double refund, over-capture). |
| **Security/Scam** | **Auth Integrity:** Distinguish between gRPC and public requests. Verify internal client identity and enforce strict authorization/permission controls. |
| **Feedback Mechanism** | **Callback Reliability:** Review proactive callback mechanisms (Webhooks). Ensure deduplication TTLs are aligned with partner retry intervals. |
| **Event Traceability** | **Universal ID:** Ensure a `Correlation ID` exists in every gRPC context to enable end-to-end querying across ACv2, PE, AE, and ZAS. |
| **Data Freshness** | **CTA Stale-Check:** Every 'Confirm' call must validate the **Adjusted** (Hiệu chỉnh) timestamp to ensure the system is not acting on stale data. |
| **Sync Reliability** | **Atomic Sync:** Audit the synchronization between Asset Exchange deductions and FMC ledger postings in ZAS. |
| **Error Recovery** | **Orchestration Resilience:** Map out flows involving Retry/Replay mechanisms. Define the logic, trigger conditions, and responsible actors for recovery paths. |
