# Phase 8 — Operations

**Goal:** Make the service operable — health/readiness signals for orchestration, and a documented filesystem permission model shared with Nginx/DevOps.

**Source:** `doc/README.md` §26 (Health Checks), §27 (File Permissions).

**Depends on:** Phase 1 (DB connectivity to check in readiness), Phase 0 (filesystem config).

---

## Task 8.1 — Liveness endpoint

Per §26: `GET /health` indicates whether the application process is running.

**Acceptance criteria**
- Returns a fast, dependency-free 200 response whenever the process is up (no DB/filesystem checks here — that's readiness).

---

## Task 8.2 — Readiness endpoint

Per §26: `GET /ready` indicates whether the application is ready to process requests, and may validate:
- Database connectivity.
- Required filesystem availability (i.e., `target_dir` / `temp_dir` from Phase 0/3 are present and writable).

**Acceptance criteria**
- Returns non-200 (e.g., 503) if the DB is unreachable or the storage directories are missing/unwritable.
- Returns 200 only when both checks pass.

---

## Task 8.3 — Filesystem permission documentation

Per §27, the deployment (not the application code) must enforce:

```text
Golang Application → Read, Write, Delete on the asset directory
Nginx               → Read Only on the asset directory
```

- Nginx must not be able to modify, delete, or upload files.
- The temporary upload directory must not be publicly accessible via Nginx.

This is a **deployment/ops document**, produced collaboratively with DevOps (§32 lists "Shared filesystem permissions" and "Nginx configuration" as DevOps-owned dependencies) — not application code.

**Acceptance criteria**
- A short doc (e.g., `doc/deployment.md`, produced in Phase 10) specifies the exact ownership/permission bits expected on `target_dir` and `temp_dir`, and explicitly states `temp_dir` must be excluded from Nginx's served paths.

---

## Definition of Done for Phase 8
- `/health` and `/ready` are wired into the router and reflect real dependency state.
- Filesystem permission expectations are written down for DevOps handoff (finalized in Phase 10).
