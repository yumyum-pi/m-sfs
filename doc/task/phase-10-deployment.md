# Phase 10 — Deployment Readiness

**Goal:** Give DevOps and the (separately built) frontend team what they need to deploy and integrate with this service correctly, including explicit scope boundaries so no one builds against unsupported behavior.

**Source:** `doc/README.md` §16 (Public Static File Access), §23 (Application Configuration), §27 (File Permissions), §29 (In Scope), §30 (Out of Scope), §32 (Key Risks and Dependencies).

**Depends on:** All prior phases (documents the finished system).

---

## Task 10.1 — Reference `config.json`

Produce a real, complete example `config.json` (not the conceptual snippet in §23) reflecting the actual fields implemented in Phase 0, including realistic values for `files`, `jwt`, `database` (both a MySQL example and a SQLite example), and `logging`.

**Acceptance criteria**
- The example config loads successfully against the Phase 0 config loader with no modification.

---

## Task 10.2 — Deployment notes

Document, for DevOps handoff:
- **Nginx setup expectations** (§16): static files served directly by Nginx from `target_dir`, on a separate domain/URL from the management API (example from §16: `static.example.com` vs `asset-manager.example.com`), with no JWT and no pass-through to the Go app.
- **Filesystem permissions** (§27, finalized from Phase 8, Task 8.3): Go app gets read/write/delete on `target_dir`/`temp_dir`; Nginx gets read-only on `target_dir` only — `temp_dir` must not be in Nginx's served path at all.
- **Dependencies owned by other teams** (§32): Nginx configuration, shared filesystem permissions, JWT claim structure, role definitions, MySQL availability, and `account_id` presence in JWT are all owned outside this project — call these out explicitly as prerequisites, not implementation tasks for this codebase.

**Acceptance criteria**
- A DevOps engineer unfamiliar with the project can provision the deployment (directories, permissions, Nginx config skeleton) from this doc alone, without reading the full `doc/README.md`.

---

## Task 10.3 — Explicit scope-boundary summary for integrators

Per §30, the following are explicitly **not** supported by this service, and should be restated plainly for whoever builds the frontend or integrates with the API, so they don't design against unsupported capabilities:
- No asset replace/update (§15, §30) — clients must upload-new + delete-old; asset URLs are immutable.
- No static file serving by the Go app — only Nginx serves files publicly.
- No in-app image processing (resize/crop/thumbnail/convert/compress).
- No CDN integration.
- No object storage (S3/Azure Blob/GCS) — local filesystem only.
- No hot config reload — restart required after config changes.
- No duplicate-file rejection — duplicates are allowed by design (§20).

**Acceptance criteria**
- This list is published somewhere the frontend/integrator team will actually see it (e.g., linked from the API's top-level docs or an integration README), not buried only in this task doc.

---

## Definition of Done for Phase 10
- A DevOps engineer and a frontend/integrator engineer each have a self-contained document answering "what do I need to know to deploy against / build against this service," without needing to read the full scope document.
