# Phase 1 — Database Layer

**Goal:** Provide persistent storage for asset metadata and audit events, working identically against MySQL (production) and SQLite (local development).

**Source:** `doc/README.md` §18 (Database), §19 (Asset Metadata), §25 (Audit Logging).

**Depends on:** Phase 0 (config for DB connection, logging for errors).

---

## Task 1.1 — Design the asset metadata schema

Per §19, the `assets` table must store:

| Field | Type notes |
| ----- | ---------- |
| `id` | UUID of the asset (primary key) |
| `account_id` | Tenant identifier — indexed, required |
| `original_file_name` | As uploaded, unchanged |
| `stored_file_name` | Generated `{account_id}_{asset_uuid}.{extension}` (§7) |
| `file_path` | Physical storage path |
| `file_size` | Bytes |
| `extension` | File extension |
| `mime_type` | Validated MIME type |
| `created_at` | Timestamp |
| `created_by` | User identifier |
| `updated_at` | Timestamp |
| `deleted_at` | Nullable — soft delete marker |

Note: §19 explicitly states the checksum field is **out of scope** — do not add it.

**Acceptance criteria**
- Schema/DDL documented and reviewed against the table above field-for-field.
- Index on `(account_id, deleted_at)` at minimum, to support tenant-scoped search excluding deleted assets efficiently.

---

## Task 1.2 — Design the audit log schema

Per §25, suggested audit metadata:

| Field | Description |
| ----- | ----------- |
| `id` | Audit event identifier |
| `asset_id` | Related asset |
| `account_id` | Related tenant |
| `action` | Upload/Delete |
| `performed_by` | User identifier |
| `timestamp` | Action timestamp |
| `metadata` | Additional contextual info (JSON/text) |

**Hard requirement:** audit logs must remain queryable/available even after the related asset is soft-deleted (§25, §14.4) — do not cascade-delete or hide audit rows when an asset is soft-deleted.

**Acceptance criteria**
- Audit table is independent of asset row lifecycle (no FK cascade delete).

---

## Task 1.3 — Database abstraction layer (MySQL + SQLite)

Per §18, the app must abstract DB operations enough to support MySQL in production and SQLite in local dev, selected via `database` config from Phase 0.

- Define a repository interface (e.g., `AssetRepository`, `AuditRepository`) with methods needed by later phases: create, get by ID, soft-delete, search/paginate/sort, create audit entry.
- Implement against both drivers, or a single driver-agnostic query layer (e.g., `database/sql` + a query builder that works for both dialects). Avoid MySQL-only or SQLite-only SQL features in shared code paths.

**Acceptance criteria**
- Same repository interface passes the same test suite against both a SQLite test DB and a MySQL test instance (can be deferred to Phase 9, but the interface must be dialect-agnostic now).

---

## Task 1.4 — Migrations

Write migrations (using whatever migration tool is chosen, e.g. `golang-migrate`) for:
- `assets` table (Task 1.1).
- `audit_logs` table (Task 1.2).

**Acceptance criteria**
- Migrations run cleanly against a fresh SQLite file and a fresh MySQL database.
- Migrations are idempotent / re-runnable in CI.

---

## Definition of Done for Phase 1
- Repository layer compiles and is unit-testable independent of HTTP/auth layers.
- Both `assets` and `audit_logs` tables exist and match the field lists in §19 and §25 exactly.
