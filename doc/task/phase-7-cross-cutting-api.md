# Phase 7 — Cross-Cutting API Concerns

**Goal:** Standardize how every endpoint from Phases 4–6 responds, and finalize the audit logging contract across all operations.

**Source:** `doc/README.md` §22 (API Response Standard), §25 (Audit Logging).

**Depends on:** Phases 4, 5, 6 (endpoints to retrofit/standardize).

---

## Task 7.1 — Standard response envelope

Per §22, all APIs must follow a common format.

Success:

```json
{
  "success": true,
  "data": {},
  "message": "Asset uploaded successfully"
}
```

Error:

```json
{
  "success": false,
  "error": {
    "code": "INVALID_FILE_TYPE",
    "message": "The uploaded file type is not supported"
  }
}
```

Per §22, the final spec should define: error codes, HTTP status codes, validation errors, authorization errors, and pagination response structure. Concretely for this project:

- Define a fixed error-code enum covering at minimum: invalid file type, file too large, invalid file signature, unauthorized (401), forbidden (403), not found (404), validation error (400), internal error (500).
- Map each error code to an HTTP status.
- Define the `data` shape for paginated list responses (e.g., `{ items: [...], page, page_size, total }`) and retrofit Phase 5's search endpoint to it.

**Acceptance criteria**
- Every endpoint built in Phases 4–6 is retrofitted to return this envelope shape, including error paths (validation failures, auth failures, not-found).
- A single shared response-writing helper/middleware is used by all handlers, so the shape can't drift per-endpoint.

---

## Task 7.2 — Audit logging integration review

Per §25: audit events cover asset uploaded, asset deleted, and "failed management operations where relevant." This task is a review pass across Phases 4 and 6 to confirm:

- Upload (Task 4.3) and Delete (Task 6.2) audit entries are consistently shaped (same metadata field usage).
- Decide which failure cases (if any) beyond successful upload/delete warrant an audit entry vs. just an operational log entry (Phase 0) — document the decision so it's applied consistently rather than ad hoc per endpoint.
- Confirm audit logs are queryable independently of operational logs (separate store/table, per §25 "separately from operational logs").

**Acceptance criteria**
- A short written note (can live in this doc or a follow-up ADR) states which operations are audited and why, so future endpoints follow the same rule.

---

## Definition of Done for Phase 7
- Every API response in the system — success or error — conforms to the §22 envelope.
- Audit logging behavior is consistent and documented across upload and delete.
