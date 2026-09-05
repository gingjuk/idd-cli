---
idd:
  version: "1.1"
  package: pkg/walk
  namespace: PKG_WALK
---

# Specifications: pkg/walk

## SPEC-PKG_WALK-001: Synchronous visitor boundary

- **Components:** `FilesystemWalker`
- **Contracts:** `FileTraversal`

**Requirement:**

Filesystem traversal must deliver each selected entry to a caller-supplied
synchronous visitor together with its file metadata. A direct-match visitor
error must stop processing and be returned unchanged so callers can implement
fail-fast collection.

### Boundary and failure behavior

The visitor owns domain parsing and may return any error value. The traversal
layer neither wraps that direct error nor converts it into a validation
finding. The callback must be non-nil when matches are possible. Recursive
descendant errors have the narrower behavior documented by
`SPEC-PKG_WALK-002`.

### Acceptance evidence

**Acceptance:**

A sentinel error returned for a directly matched file is the exact error
returned by `Walk`; a no-match invocation never calls the visitor.

## SPEC-PKG_WALK-002: Pattern expansion, recursion, and deduplication

- **Components:** `FilesystemWalker`
- **Contracts:** `FileTraversal`

**Requirement:**

Traversal must process standard-library glob patterns deterministically,
recurse when a direct match is a directory, and invoke the visitor no more than
once for the same exact path string within one call.

### Required behavior

- Nil and empty pattern collections succeed without visits.
- A pattern with no matches is a successful no-op.
- Invalid patterns and entries that cannot be stated are skipped.
- A directly matched directory is visited and then traversed depth-first using
  `filepath.Walk`.
- Overlapping patterns share one visited set.
- Filesystem or visitor errors returned from recursive descent stop that
  directory walk but are currently suppressed by the outer loop.

### Implementation boundary

The specification does not require doublestar expansion, path
canonicalization, hidden-file filtering, symlink following, parallelism, or
cross-call caching. Selection inside a visited directory belongs to the
visitor.

### Acceptance evidence

**Acceptance:**

Temporary directory trees demonstrate direct file visits, descendant visits,
empty inputs, and exact-path duplicate suppression. A later behavior change to
recursive error propagation requires a new regression scenario because the
current suite only proves direct-match error propagation.

## SPEC-PKG_WALK-003: Exact final-extension matching

- **Components:** `FilesystemWalker`
- **Contracts:** `ExtensionMatch`

**Requirement:**

Extension filtering must compare the final suffix returned by
`filepath.Ext` with a caller-provided allow-list using exact case-sensitive
equality.

### Edge cases and non-goals

The allow-list entries include their leading dot. A path without a suffix,
an empty allow-list, or `file.go.txt` tested against `.go` must not match.
The helper does not normalize case, infer MIME type, inspect file contents, or
treat compound suffixes specially.

### Acceptance evidence

**Acceptance:**

Table-driven cases cover supported and unsupported suffixes, missing
extensions, compound filenames, and empty allow-lists.
