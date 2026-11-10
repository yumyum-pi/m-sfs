# Phase 5 — Read APIs

**Goal:** Implement asset search/listing and single-asset metadata retrieval, both tenant-scoped, paginated, and sortable within a safe allow-list.

**Source:** `doc/README.md` §14.2 (Search Assets), §14.3 (Get Asset Metadata), §21 (Search and Pagination).

**Depends on:** Phase 1 (DB repository), Phase 2 (tenant isolation).

---

## Task 5.1 — `GET /api/v1/assets` search

Per §14.2, support:
- **Search** over `original_file_name` and `stored_file_name`.
- **Filters**: file type (extension), `created_by`, `created_date`, `updated_date`.
- Tenant/account isolation automatically enforced via JWT (Phase 2, Task 2.3) — a client cannot search another tenant's assets regardless of query parameters.
- Soft-deleted assets must be excluded from normal search results (§14.4, §28).

**Acceptance criteria**
- Search matches on filename fields (exact or partial match — decide and document one behavior, e.g., case-insensitive substring match, during implementation).
- Each filter can be applied independently or combined.
- Results never include soft-deleted assets or assets belonging to another tenant, even under adversarial query parameters.

---

## Task 5.2 — Pagination and sorting

Per §21:
- Default page size: **20**.
- Query shape: `GET /api/v1/assets?page=1&page_size=20`.
- Sorting: `GET /api/v1/assets?page=1&page_size=20&sort_by=created_at&sort_order=desc`.
- **Hard requirement:** sorting fields must go through an **explicit allow-list** — do not pass arbitrary user-provided field names into SQL `ORDER BY` (SQL injection / information-disclosure risk). Recommended initial allow-list: `created_at`, `updated_at`, `original_file_name`, `file_size`.

**Acceptance criteria**
- Requesting a `sort_by` value outside the allow-list is rejected (400) or falls back to a safe default — pick one behavior and document it; rejecting with a clear error is recommended for API clarity.
- Pagination response includes enough metadata for clients to page through results (total count, current page, page size — finalize exact shape with Phase 7's response standard).

---

## Task 5.3 — `GET /api/v1/assets/{id}` metadata retrieval

Per §14.3: returns metadata only, **never the file content**. This is confirmed in-scope per the README's note that it's a necessary companion to Search.

**Acceptance criteria**
- Returns the full metadata record (same fields as §19) for an asset owned by the caller's tenant.
- Returns not-found (404) both when the asset doesn't exist and when it belongs to another tenant — do not leak existence of other tenants' assets via a distinct error code/message.
- Returns not-found for soft-deleted assets (consistent with search exclusion), unless a future admin/audit-only endpoint is explicitly added — out of scope for this task.

---

## Definition of Done for Phase 5
- Both read endpoints are tenant-safe, exclude soft-deleted assets, and paginate/sort safely.
- Allow-list enforcement for sorting is unit-tested (cross-referenced in Phase 9, Task 9.4).
