# Task Documents — Static Asset Management Service

This directory breaks the project scope defined in [`../README.md`](../README.md) into ordered, phase-based implementation documents. Each phase builds only on what earlier phases deliver, so the system stays runnable and testable throughout development.

## Phase Index

| Phase | Document | Summary |
| ----- | -------- | ------- |
| 0 | [phase-0-bootstrap.md](phase-0-bootstrap.md) | Go module, project layout, config loading, logging |
| 1 | [phase-1-database.md](phase-1-database.md) | Asset & audit schema, MySQL/SQLite abstraction, migrations |
| 2 | [phase-2-auth-tenancy.md](phase-2-auth-tenancy.md) | JWT auth, RBAC, tenant isolation |
| 3 | [phase-3-file-validation.md](phase-3-file-validation.md) | Storage paths, size/extension/MIME/magic-byte checks, SVG sanitization, naming |
| 4 | [phase-4-upload-api.md](phase-4-upload-api.md) | `POST /api/v1/assets` end-to-end upload workflow |
| 5 | [phase-5-read-apis.md](phase-5-read-apis.md) | Search, pagination, sorting, get-by-id |
| 6 | [phase-6-delete-api.md](phase-6-delete-api.md) | Soft delete + physical file removal |
| 7 | [phase-7-cross-cutting-api.md](phase-7-cross-cutting-api.md) | Standard response envelope, error codes, audit log integration |
| 8 | [phase-8-operations.md](phase-8-operations.md) | Health/readiness endpoints, filesystem permissions |
| 9 | [phase-9-testing.md](phase-9-testing.md) | Unit, integration, tenant-isolation, and failure-injection tests |
| 10 | [phase-10-deployment.md](phase-10-deployment.md) | Reference config, deployment notes, scope boundaries |

## Ordering Principle

```text
infra/config → data layer → auth/tenancy → validation primitives
   → write path (upload) → read paths (search/get) → delete path
   → cross-cutting (responses/audit) → operations → tests → deployment docs
```

Every task references the source section(s) of `doc/README.md` it implements, so requirements can always be traced back to the approved scope document. Nothing here introduces new scope — these documents only decompose existing scope into buildable units.
