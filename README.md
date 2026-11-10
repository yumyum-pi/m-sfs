# m-sfs — Static Asset Management Service

A multi-tenant Golang REST API for centrally managing public-facing static assets (icons, banners, email images, PDFs, etc.). The service owns upload, validation, metadata, and tenant isolation; **Nginx remains responsible for serving the files publicly**.

Full project scope: [doc/README.md](doc/README.md). Phased implementation plan: [doc/task/README.md](doc/task/README.md).

## Status

🚧 **Phase 0 (bootstrap) complete** — config loading and logging are in place; the HTTP server, database, auth, and asset APIs are not yet implemented. See the [task phases](doc/task/README.md) for what's next.

## Architecture

- **Management plane** (this app): authentication, authorization, upload, validation, metadata, search, soft delete, audit logging.
- **Serving plane** (Nginx): serves files directly from the shared storage directory; never touches the Go app.

Tenancy is derived from `account_id` in the JWT and enforced at the API layer — a caller can only ever see their own account's assets.

## Project Layout

```
cmd/server/       entrypoint (main.go)
internal/config/   config.json loading, validation, defaults
internal/logging/  structured logger
internal/db/       database abstraction (MySQL prod / SQLite dev) — Phase 1
internal/auth/     JWT auth + RBAC — Phase 2
internal/asset/    file validation, storage, asset business logic — Phase 3+
internal/api/      HTTP handlers and routing — Phase 4+
migrations/        SQL migrations — Phase 1
doc/               project scope and phase-by-phase task docs
```

## Requirements

- Go 1.26.5+

## Configuration

Copy the example config and adjust for your environment:

```bash
cp config.example.json config.json
```

```json
{
  "files": {
    "target_dir": "/var/static/assets",
    "temp_dir": "/var/static/tmp",
    "max_file_size": "5MB",
    "supported_files": ["svg", "png", "jpeg", "webp", "gif", "pdf"]
  },
  "jwt": {
    "signing_method": "HS256",
    "secret": "change-me",
    "issuer": "example-auth",
    "audience": "asset-manager"
  },
  "database": {
    "driver": "sqlite",
    "path": "./data/dev.db"
  },
  "logging": {
    "level": "INFO"
  }
}
```

Notes:
- `files.supported_files`, when set, **replaces** the default extension list rather than extending it.
- `files.max_file_size` accepts human-readable sizes (`"5MB"`, `"512KB"`, etc.).
- `logging.level` must be one of `DEBUG`, `INFO`, `WARN`, `ERROR`.
- `database.driver` is `sqlite` for local development or `mysql` for production.
- Configuration is loaded once at startup — there is no hot reload; restart the process after changes.
- `config.json` is git-ignored; never commit real secrets.

## Running

```bash
go run ./cmd/server -config config.json
```

## Testing

```bash
go test ./...
```
