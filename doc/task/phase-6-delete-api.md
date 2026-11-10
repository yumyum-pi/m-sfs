# Phase 6 — Delete API

**Goal:** Implement soft-deletion of assets such that public URLs stop resolving while metadata and audit history are retained.

**Source:** `doc/README.md` §14.4 (Delete Asset), §28 (Soft Delete Behaviour), §19 (Asset Metadata — `deleted_at`).

**Depends on:** Phase 1 (DB repository), Phase 2 (auth/tenancy), Phase 3 (target-dir path resolution).

---

## Task 6.1 — `DELETE /api/v1/assets/{id}`

Per §14.4 and §28, the required flow is:

```text
DELETE Request
   → Authorization Check
   → Mark Asset as Deleted (set deleted_at)
   → Remove from Search Results
   → Remove Public File (delete from target_dir)
   → Create Audit Log
```

This is the **explicit resolved interpretation** from §14.4: soft-delete the DB record **and** remove the physical file from `target_dir`, so the public URL returns 404 while metadata/audit history remain queryable. Do not implement "DB-only" soft delete that leaves the file publicly servable by Nginx.

- Authorization: only allow deletion within the caller's own tenant (Phase 2 tenant isolation applies here too — a client cannot delete another tenant's asset even if it somehow knows the ID).
- The frontend (out of scope for this project) is responsible for delete confirmation UX — the API itself does not need a confirmation step, just the delete operation.

**Acceptance criteria**
- After deletion: `deleted_at` is set, the asset no longer appears in `GET /api/v1/assets` search (Phase 5) or `GET /api/v1/assets/{id}` (Phase 5), the physical file no longer exists in `target_dir`, and an audit log entry exists.
- Deleting an already-deleted or nonexistent asset returns a sensible error (e.g., 404) rather than a silent success or a crash.
- Deleting an asset belonging to another tenant returns 404 (not 403), consistent with Phase 5's non-disclosure behavior.

---

## Task 6.2 — Audit log entry on delete

Per §25: create an audit log entry (`action = Delete`) with `asset_id`, `account_id`, `performed_by`, `timestamp`.

**Acceptance criteria**
- Audit row for the delete is created and remains queryable after the asset row is marked deleted (verifies Phase 1's no-cascade-delete design).

---

## Note on file removal failure

If the DB soft-delete succeeds but physical file removal fails (e.g., filesystem error), decide and document the recovery behavior — e.g., retry via the same temp-cleanup-style background sweep pattern used in Phase 4, Task 4.4, or surface an error and leave the asset queryable for a retry. This should not silently leave a "deleted" asset's file still publicly servable.

**Acceptance criteria**
- A simulated file-removal failure does not leave the public file servable indefinitely without at least being logged as an operational error (Phase 0 logging) for follow-up.

---

## Definition of Done for Phase 6
- Deleted assets are unreachable via both the read APIs and the public Nginx-served URL, while remaining in the audit trail.
