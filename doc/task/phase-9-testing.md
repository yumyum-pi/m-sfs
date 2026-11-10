# Phase 9 — Testing

**Goal:** Establish automated test coverage for the security- and correctness-critical behaviors called out across the scope doc — especially tenant isolation, upload consistency, and validation robustness.

**Source:** Cross-cutting — validates behaviors from `doc/README.md` §5, §10, §11, §12, §14, §21, §28, and the risks in §32.

**Depends on:** Phases 1–6 (features under test must exist).

---

## Task 9.1 — Validator unit tests

Covers Phase 3 primitives:
- File size validation at/over/under the configured limit.
- Extension allow-list, including the "config replaces default" behavior (§8).
- MIME type detection against known-good and known-bad samples.
- Magic-byte validation: at least one renamed/mismatched-content case per supported type (§10.4's `malicious.exe` → `image.png` example).
- SVG sanitizer: at least one sample with embedded script/event-handler content (§11).

**Acceptance criteria**
- Each validator has both a passing and a failing test case; magic-byte and SVG tests specifically cover the adversarial cases named in §10.4/§11.

---

## Task 9.2 — Tenant isolation integration tests

Covers Phase 2's core security guarantee (§5, §17), exercised through the real HTTP APIs from Phases 4–6:
- A JWT for account A cannot see account B's assets via search (Phase 5.1).
- A JWT for account A cannot fetch account B's asset by ID (Phase 5.3) — must 404, not 403 (no existence leak).
- A JWT for account A cannot delete account B's asset (Phase 6.1) — must 404.
- A request with `?account_id=<other-tenant>` is ignored in favor of the JWT's `account_id` (Phase 2, Task 2.3).

**Acceptance criteria**
- All four scenarios above are automated tests, not just manual verification — this is the highest-risk area of the system per §5's explicit emphasis.

---

## Task 9.3 — Upload/delete failure-injection tests

Covers Phase 4's consistency guarantees (§12) and Phase 6's file-removal edge case:
- DB record creation fails after validation passes → no file left in `target_dir`, no orphaned temp file.
- Atomic move to `target_dir` fails after DB record is created → DB record is rolled back.
- Delete: DB soft-delete succeeds but physical file removal fails → behavior matches what was documented in Phase 6's "Note on file removal failure."

**Acceptance criteria**
- Tests use fault injection (e.g., a fake repository/filesystem that returns errors on demand) rather than relying on flaky real-world failures.
- No test run leaves orphaned files or DB rows behind (verified by inspecting temp/target dirs and DB state post-test).

---

## Task 9.4 — Pagination/sorting boundary tests

Covers Phase 5.2 (§21):
- `sort_by` outside the allow-list is rejected/handled per the documented behavior — not passed through to SQL.
- Pagination at boundaries: `page_size` larger than result set, `page` beyond the last page, default `page_size=20` when omitted.

**Acceptance criteria**
- A `sort_by` value crafted to look like a SQL injection attempt (e.g., `created_at; DROP TABLE assets`) is rejected by the allow-list, not passed to the query layer.

---

## Definition of Done for Phase 9
- Tenant isolation, upload/delete consistency, and adversarial file validation all have automated regression coverage — these are the three risk areas explicitly flagged in §32 (Key Risks).
