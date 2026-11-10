# Phase 2 — Auth & Tenant Isolation

**Goal:** Every management API call is authenticated via JWT, authorized by role, and automatically scoped to the caller's tenant — with no way for a client to access another tenant's assets.

**Source:** `doc/README.md` §5 (Multi-Tenancy), §17 (Authorization).

**Depends on:** Phase 0 (config for JWT settings, logging for auth failures).

---

## Task 2.1 — JWT authentication middleware

Per §17, implement middleware that:
1. Extracts and validates the JWT (signature, expiry, issuer/audience as configured in Phase 0's `jwt` config block).
2. Extracts claims identifying: user, account/tenant (`account_id`), and role/authorization context.
3. Rejects unauthenticated or invalid-token requests before any handler logic runs.

**Acceptance criteria**
- Requests with missing, malformed, expired, or badly-signed tokens are rejected (401) before reaching business logic.
- Valid-token requests have parsed claims available to downstream handlers (e.g., via request context).

---

## Task 2.2 — Role-based authorization

Per §17, apply role checks after authentication. §17 notes exact role definitions should integrate with the organization's existing authorization model where applicable — implement this as a pluggable/extensible role-check layer rather than hardcoding roles that may need to change.

**Acceptance criteria**
- A request with a valid JWT but insufficient role/permission is rejected (403).
- Role requirements are declared per-endpoint (e.g., upload vs. delete may require different roles), not globally.

---

## Task 2.3 — Tenant isolation enforcement

Per §5, `account_id` is the tenant identifier and **must come from the JWT**, never trusted from client input.

Enforcement flow (§5, §17):

```text
Authenticated User
   → Extract account_id from JWT
   → Apply account filter
   → Execute operation
```

- Implement a helper that all data-access code paths (search, get, upload, delete) must use to obtain the effective `account_id` — always JWT-derived.
- **Explicit security requirement:** if a client passes `?account_id=other` (or any tenant-like filter) that differs from the JWT's `account_id`, the system must ignore/override it with the JWT-derived value, not merge or trust it. Example from §5:

```text
JWT account_id: acme
Request: GET /api/v1/assets?account_id=other
System: automatically enforces account_id from JWT (acme), ignoring "other"
```

**Acceptance criteria**
- A user-supplied `account_id` query parameter never changes which tenant's data is returned/affected.
- This behavior has a dedicated automated test (also tracked in Phase 9, Task 9.2) since it's the core multi-tenant security guarantee of the whole system.

---

## Definition of Done for Phase 2
- Middleware chain (auth → authz → tenant-scope) is composable and reusable by every handler built in later phases.
- No handler code written after this phase needs to re-implement JWT parsing or tenant filtering — it consumes the context set by this middleware.
