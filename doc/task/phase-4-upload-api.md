# Phase 4 — Upload API (Write Path)

**Goal:** Implement `POST /api/v1/assets`, orchestrating Phase 2 (auth/tenancy) and Phase 3 (validation primitives) into the full upload workflow, with strict filesystem/database consistency.

**Source:** `doc/README.md` §12 (File Upload), §13 (Temporary Storage), §14.1 (Upload Asset API).

**Depends on:** Phase 1 (DB repository), Phase 2 (auth/tenancy middleware), Phase 3 (validators, sanitizer, naming).

---

## Task 4.1 — Implement `POST /api/v1/assets` workflow

Per §12, the exact required sequence is:

```text
Client Upload
   → Temporary Directory
   → Validate File Size
   → Validate Extension
   → Validate MIME Type
   → Validate File Signature
   → Sanitize SVG (if applicable)
   → Generate UUID
   → Generate Stored Filename
   → Create Database Record
       (failure → Rollback / Cleanup)
   → Atomically Move File
       (failure → Rollback DB Record)
   → Success
```

Per §14.1, the endpoint's responsibilities are: authenticate user, authorize user, extract `account_id` (from Phase 2), upload file, validate file (Phase 3 validators in order above), generate asset UUID, generate stored filename, store metadata, move file to target directory, create audit log.

**Acceptance criteria**
- The validation steps run in the order specified in §12 (size → extension → MIME → signature → SVG sanitize), short-circuiting on first failure.
- Successful upload returns metadata for the created asset in the standard success envelope (Phase 7 will formalize the envelope shape; use a reasonable shape now that Phase 7 can retrofit).

---

## Task 4.2 — Filesystem/database consistency & rollback

Per §12: *"No orphaned database records or unmanaged files should remain after failed operations."*

Two failure points must be handled explicitly:
1. **DB record creation fails** → clean up (delete) the temp file; nothing is written to `target_dir`.
2. **Atomic move to `target_dir` fails after the DB record was created** → roll back (delete) the DB record.

- "Atomically move" means the file must not be observable in `target_dir` in a partial/corrupt state — use a rename-based move within the same filesystem, or write-then-rename, rather than a streamed copy that could be interrupted.

**Acceptance criteria**
- Simulate a DB failure post-validation: no file appears in `target_dir`, no orphaned temp file remains.
- Simulate a move failure post-DB-write: the DB record is rolled back (not left as a dangling reference to a nonexistent file).
- (Full failure-injection test suite tracked in Phase 9, Task 9.3 — this task just needs the rollback logic implemented and manually verifiable.)

---

## Task 4.3 — Audit log entry on upload

Per §14.1 and §25: create an audit log entry (`action = Upload`) recording `asset_id`, `account_id`, `performed_by`, `timestamp`, and relevant metadata, using the audit repository from Phase 1.

**Acceptance criteria**
- Every successful upload produces exactly one audit log row.
- Failed uploads that are meaningful management-operation failures are considered for audit per §25 ("Failed management operations where relevant") — at minimum, log via operational logging (Phase 0); audit-log failed uploads if the team decides they're relevant during implementation.

---

## Task 4.4 — Temp file cleanup mechanism

Per §13: *"The application should include a cleanup mechanism for abandoned or failed temporary files."*

- Implement a cleanup routine (e.g., a periodic background sweep) that removes temp files older than a threshold that were never promoted to `target_dir` (covers crashes mid-upload, not just handled errors already cleaned up by Task 4.2).

**Acceptance criteria**
- A temp file artificially aged past the threshold is removed by the cleanup routine without affecting recently-written in-progress uploads.

---

## Definition of Done for Phase 4
- `POST /api/v1/assets` is fully functional end-to-end against a real DB and filesystem, including audit logging and rollback on both failure points.
- Temp directory never accumulates orphaned files under normal or failure operation.
