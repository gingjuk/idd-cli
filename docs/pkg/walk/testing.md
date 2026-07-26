---
idd:
  version: "1.0"
  package: pkg/walk
---

# Testing: pkg/walk

## TEST-PKG_WALK-001: No-match traversal is a successful no-op

- **Kind:** `contract`
- **Covers:** `SPEC-PKG_WALK-001`, `SPEC-PKG_WALK-002`
- **Contracts:** `FileTraversal`

**Purpose:**

Prove that a valid pattern with no matches returns nil and never invokes the
visitor. The callback itself fails the test if called, making absence of a
visit the oracle.

**Oracle:** The test passes only when its assertions confirm a valid
pattern with no matches returns nil and never invokes the visitor. The callback itself
fails the test if called, making absence of a visit the oracle.

## TEST-PKG_WALK-002: Direct visitor errors stop traversal

- **Kind:** `contract`
- **Covers:** `SPEC-PKG_WALK-001`
- **Contracts:** `ExtensionMatch`

**Purpose:**

Prove that an error returned while visiting a directly matched file is returned
unchanged. A temporary file creates one deterministic match and a sentinel
`os.ErrPermission` distinguishes propagation from wrapping or suppression.

**Oracle:** The test passes only when its assertions confirm an
error returned while visiting a directly matched file is returned unchanged. A temporary
file creates one deterministic match and a sentinel `os.ErrPermission` distinguishes
propagation from wrapping or suppression.

## TEST-PKG_WALK-003: Matched directories expose descendants

- **Kind:** `test`
- **Covers:** `SPEC-PKG_WALK-002`

**Purpose:**

Prove that matching a directory visits a file below it. The test records paths
from an isolated temporary tree and uses the descendant basename as its
observable oracle.

**Oracle:** The test passes only when its assertions confirm
matching a directory visits a file below it. The test records paths from an isolated
temporary tree and uses the descendant basename as its observable oracle.

## TEST-PKG_WALK-004: Overlapping patterns do not duplicate callbacks

- **Kind:** `test`
- **Covers:** `SPEC-PKG_WALK-002`

**Purpose:**

Prove that the same exact file matched by both a wildcard and a literal pattern
is delivered once. A callback counter detects duplicate traversal independent
of ordering.

**Oracle:** The test passes only when its assertions confirm the
same exact file matched by both a wildcard and a literal pattern is delivered once. A
callback counter detects duplicate traversal independent of ordering.

## TEST-PKG_WALK-005: Nil patterns are accepted

- **Kind:** `test`
- **Covers:** `SPEC-PKG_WALK-002`

**Purpose:**

Prove that a nil pattern slice returns nil with zero callbacks, preserving the
no-op contract for optional configuration.

**Oracle:** The test passes only when its assertions confirm a nil
pattern slice returns nil with zero callbacks, preserving the no-op contract for
optional configuration.

## TEST-PKG_WALK-006: Empty patterns are accepted

- **Kind:** `test`
- **Covers:** `SPEC-PKG_WALK-002`

**Purpose:**

Prove that an allocated but empty pattern slice has the same zero-visit,
zero-error behavior as nil input.

**Oracle:** The test passes only when its assertions confirm an
allocated but empty pattern slice has the same zero-visit, zero-error behavior as nil
input.

## TEST-PKG_WALK-007: Direct file patterns visit their matches

- **Kind:** `test`
- **Covers:** `SPEC-PKG_WALK-002`

**Purpose:**

Prove that a wildcard selecting Go files visits the matching direct files while
leaving an unmatched text file and an unmatched nested Go file outside the
result. The test uses a temporary directory so host filesystem contents cannot
affect the count.

**Oracle:** The test passes only when its assertions confirm a
wildcard selecting Go files visits the matching direct files while leaving an unmatched
text file and an unmatched nested Go file outside the result. The test uses a temporary
directory so host filesystem contents cannot affect the count.

## TEST-PKG_WALK-008: Extension matching uses the final exact suffix

- **Kind:** `test`
- **Covers:** `SPEC-PKG_WALK-003`

**Purpose:**

Prove exact final-extension semantics across allowed suffixes, unsupported
suffixes, compound filenames, missing suffixes, and empty allow-lists. The
table's expected boolean is the complete oracle and no filesystem fixture is
required.

**Oracle:** The test passes only when its assertions confirm exact
final-extension semantics across allowed suffixes, unsupported suffixes, compound
filenames, missing suffixes, and empty allow-lists. The table's expected boolean is the
complete oracle and no filesystem fixture is required.

## Strategy

Traversal tests use `t.TempDir` and create only the paths needed by each
scenario. They assert observable callback counts, paths, and errors rather than
private traversal state. Tests are synchronous and contain no timing or
platform-external dependencies.

The suite deliberately does not claim coverage of invalid glob syntax,
permission-denied `Stat`, symlink behavior, platform-specific ordering, or
recursive descendant visitor-error propagation. Those cases need explicit
fixtures before the contract can be strengthened or changed.
