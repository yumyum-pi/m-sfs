# Phase 0 — Project Bootstrap

**Goal:** Stand up a runnable Go application skeleton with configuration loading and logging, so every later phase has a place to plug into.

**Source:** `doc/README.md` §23 (Application Configuration), §24 (Application Logging).

**Depends on:** Nothing (first phase).

---

## Task 0.1 — Initialize Go module and project layout

Create the base repository structure.

- `go.mod` with an appropriate module path.
- Suggested layout:
  - `cmd/server/main.go` — entrypoint.
  - `internal/config/` — config loading.
  - `internal/logging/` — logger setup.
  - `internal/db/` — database layer (Phase 1).
  - `internal/auth/` — JWT/authz (Phase 2).
  - `internal/asset/` — validation, storage, business logic (Phase 3+).
  - `internal/api/` — HTTP handlers/router (Phase 4+).
  - `migrations/` — SQL migrations (Phase 1).

**Acceptance criteria**
- `go build ./...` succeeds with an empty `main` that starts and exits cleanly.

---

## Task 0.2 — Implement `config.json` loading

Per §23, config is loaded once at startup from `config.json`. No hot reload — the process must be restarted for config changes to take effect.

Implement a config struct covering all four areas from the conceptual example in §23:

```json
{
  "files": {
    "target_dir": "/var/static/assets",
    "temp_dir": "/var/static/tmp",
    "max_file_size": "5MB",
    "supported_files": ["svg", "png", "jpeg", "webp", "gif", "pdf"]
  },
  "jwt": {},
  "database": {},
  "logging": { "level": "INFO" }
}
```

- `files.target_dir` / `files.temp_dir`: string paths (§6, §13).
- `files.max_file_size`: human-readable size string, default `5MB` (§9). Parse to bytes.
- `files.supported_files`: list of extensions. **When present, replaces the default list entirely** — it does not extend it (§8). Default list when omitted: SVG, PNG, JPEG, WebP, GIF, PDF.
- `jwt`: fields needed for token validation (signing method, key/secret or JWKS URL, expected issuer/audience — finalize exact fields with Phase 2).
- `database`: connection settings sufficient to support both MySQL (prod) and SQLite (dev) per §18.
- `logging.level`: one of `DEBUG`, `INFO`, `WARN`, `ERROR` (§24).

**Acceptance criteria**
- Missing `config.json` or invalid JSON fails startup with a clear error.
- `supported_files` override behavior is unit-tested: providing `["png", "webp"]` results in exactly PNG/WebP being accepted, not PNG/WebP/SVG/etc.
- `max_file_size` parsing handles `"5MB"`-style strings.

---

## Task 0.3 — Application logging

Implement structured operational logging supporting DEBUG/INFO/WARN/ERROR levels, configurable via `logging.level` (§24).

- Log at minimum: startup/shutdown, and provide the logger as a dependency other phases will use for upload success/failure, validation failures, DB errors, file operation failures, authorization failures, and unexpected errors.
- **Hard requirement:** JWT tokens and other sensitive information must never be written to logs (§24). Add a redaction convention/helper now so later phases follow it by default.

**Acceptance criteria**
- Log level is configurable and filters output correctly (e.g., `WARN` suppresses DEBUG/INFO).
- A log-scrubbing test or code-review checklist item exists to prevent tokens/secrets from being logged.

---

## Definition of Done for Phase 0
- App builds and starts, reading `config.json` and initializing the logger.
- Config and logging are exposed as reusable internal packages for later phases.
