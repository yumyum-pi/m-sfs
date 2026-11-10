# Phase 3 — File Validation & Storage Primitives

**Goal:** Build the reusable, standalone building blocks (path resolution, validators, sanitizer, naming) that the upload workflow (Phase 4) will orchestrate. Each of these should be independently unit-testable without needing an HTTP request.

**Source:** `doc/README.md` §6 (Asset Storage), §7 (File Naming Convention), §8 (Supported File Types), §9 (File Size Limits), §10 (File Validation), §11 (SVG Security), §13 (Temporary Storage).

**Depends on:** Phase 0 (config for dirs/limits/types).

---

## Task 3.1 — Storage path resolution

Per §6 and §13:
- `target_dir`: final storage location, flat structure, no tenant subdirectories. Example: `/var/static/assets/acme_550e8400.png`.
- `temp_dir`: staging location for uploads before validation completes, e.g. `/var/static/tmp`. Must **not** be exposed by Nginx.

Implement a small storage-path utility that resolves both from config (Phase 0) and exposes helpers for building temp and final paths.

**Acceptance criteria**
- Both directories are configurable and validated to exist/be writable at startup (fail fast rather than fail on first upload).

---

## Task 3.2 — File size validation

Per §9 and §10.1: reject files exceeding `max_file_size` (default 5 MB, configurable).

**Acceptance criteria**
- Files at/under the limit pass; files over the limit are rejected with a clear error before further processing.

---

## Task 3.3 — Extension validation

Per §8 and §10.2: validate against the configured `supported_files` list. Reminder: config **replaces** the default list, it does not extend it (already parsed in Phase 0, Task 0.2 — this task consumes that resolved list).

**Acceptance criteria**
- Unsupported extensions are rejected.
- Default list (SVG, PNG, JPEG, WebP, GIF, PDF) is used only when no override is configured.

---

## Task 3.4 — MIME type validation

Per §10.3: validate the detected MIME type of the uploaded content (not just the `Content-Type` header the client sends).

**Acceptance criteria**
- Detected MIME type must be consistent with an allowed type for the (validated) extension.

---

## Task 3.5 — File signature / magic-byte validation

Per §10.4, this is a hard security requirement: validation **must not rely solely on** file extension or client-provided MIME type. The actual file content must be inspected.

Example failure case from §10.4: a file renamed from `malicious.exe` to `image.png` must be rejected because its actual byte signature doesn't match a PNG.

**Acceptance criteria**
- A binary whose magic bytes don't match its extension/claimed MIME type is rejected, even if extension and client MIME type both look valid.
- Test with at least one renamed-executable-style case per supported type.

---

## Task 3.6 — SVG sanitization

Per §11: SVG files must be sanitized before permanent storage, and sanitization must happen **before** the file is moved to the final asset directory.

- Choose/implement a sanitizer that strips or rejects unsafe SVG content (e.g., embedded scripts, external entity references, event handler attributes).
- This must run as part of the temp-directory validation pipeline (Phase 4), never after the file reaches `target_dir`.

**Acceptance criteria**
- A crafted SVG containing `<script>`, `onload=`, or similar is either sanitized (unsafe content stripped) or rejected outright — pick and document one behavior, then enforce it consistently.
- Non-SVG files are unaffected by this step.

---

## Task 3.7 — UUID generation & stored filename convention

Per §7: every uploaded file gets a UUID, and the stored filename format is:

```text
{account_id}_{asset_uuid}.{extension}
```

Example: `acme_550e8400-e29b-41d4-a716-446655440000.png`. The **original** uploaded filename is preserved unchanged in metadata (`original_file_name`), only the stored file is renamed.

**Acceptance criteria**
- Generated filenames are collision-resistant (UUID-based) and match the exact format above.
- Duplicate uploads (same content/name) are allowed and produce distinct stored filenames (§20) — this is a natural consequence of UUID-based naming, not a separate dedup check.

---

## Definition of Done for Phase 3
- Each validator (size, extension, MIME, magic bytes), the SVG sanitizer, and the filename generator exist as independently callable, independently unit-tested functions/packages.
- Phase 4 can compose these in sequence without needing to know their internals.
