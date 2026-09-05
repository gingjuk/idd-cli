---
idd:
  version: "1.1"
  package: pkg/walk
  namespace: PKG_WALK
---

# Design: pkg/walk

## Component: FilesystemWalker

**Purpose:**

`FilesystemWalker` is a small synchronous boundary between configured glob
patterns and collectors that need to inspect matching filesystem entries. It
owns pattern expansion, directory descent, and duplicate suppression. The
visitor owns every domain-specific decision, including whether a matched entry
is a document, source file, directory, or irrelevant path.

Keeping traversal independent from IDD parsing lets collectors reuse one
ordering and deduplication policy without coupling this package to identifiers,
configuration structures, or validation results.

**Ownership:**

The component owns synchronous glob expansion, recursive descent, duplicate
path suppression, callback sequencing, and exact extension matching for
callers that need a small general filesystem walker.

**Boundary:**

It does not interpret IDD files, apply ignore or documentation-unit policy,
follow semantic relationships, or decide whether an I/O failure is tolerable.
Validation-facing scans enforce fail-closed behavior at their owning boundary.

**Decisions:**

Standard-library path operations and a per-call visited set keep this utility
dependency-free and deterministic for returned path strings. Synchronous
callbacks expose ordering and failure behavior without background state.

### Responsibilities

- expand each pattern with the standard library's `filepath.Glob`;
- visit direct matches in pattern order and glob-result order;
- recursively visit an entry when the direct match is a directory;
- suppress repeated visits to the same path string across overlapping
  patterns and recursive traversal;
- expose file metadata to a synchronous callback; and
- provide an exact extension predicate for callers that already have a path.

The component does not canonicalize paths, interpret doublestar syntax, filter
hidden files, follow symbolic-link directories, or decide which failures should
be reported as IDD findings.

### Boundaries and failure containment

Glob syntax errors, failed `os.Stat` calls, and filesystem errors reported by
`filepath.Walk` are skipped so one inaccessible path does not abort collection.
A visitor error on a direct glob match is returned immediately. For a visitor
error raised below a recursively matched directory, `filepath.Walk` stops that
walk, but the outer loop currently continues to the next pattern and returns no
error for that recursive failure.

That asymmetry is part of the present implementation and is documented so
callers do not assume stronger error propagation. Normalizing it would be an
observable compatibility change rather than a documentation repair.

### Decisions and trade-offs

The package uses exact path strings as deduplication keys. This is inexpensive
and deterministic for the paths returned by a single glob operation, but two
lexically different paths to the same file can still be visited twice. It also
uses the standard library rather than a doublestar dependency; recursive
collection is obtained by matching a directory and walking it, not by assigning
special meaning to `**` inside this package.

## Architecture

```text
configured patterns
       |
       v
filepath.Glob -- invalid pattern --> skip pattern
       |
       v
direct matches -- stat failure --> skip match
       |
       +--> visitor(match)
       |
       +--> if directory: filepath.Walk(match)
                              |
                              +--> visitor(descendant)
```

One `visited` map lives for the duration of a `Walk` call. A path is marked
before the callback runs, so an overlapping later pattern cannot cause a second
callback even if the first visit returned successfully.

## Package Layout

`files.go` contains both the traversal abstraction and the extension helper.
They remain together because both are filesystem-name operations with no IDD
domain dependency. Adding collector-specific filters here would invert that
dependency and make this utility responsible for validation policy.

## Function Composition

`Walk` expands one pattern at a time, stats each new direct match, invokes the
visitor, and then descends if the match is a directory. Recursive traversal
shares the same visited map as direct matching. `MatchAnyExtensions` is a
separate pure predicate and is not called implicitly by `Walk`.

There is no background work, retained process state, cancellation mechanism, or
concurrency. The visitor completes before traversal advances.

## Dependencies

The package depends only on `os` and `path/filepath`. Its behavior therefore
inherits platform path syntax, glob ordering, extension parsing, and
`filepath.Walk` symbolic-link behavior.

## Testability Hooks

Tests create isolated temporary trees and observe callback paths and returned
errors. The synchronous visitor is the primary seam: it can count visits,
record order, or inject a sentinel error without mocking the filesystem API.
Extension matching is tested as a pure table-driven predicate.

Recursive visitor-error propagation is not currently covered by a dedicated
test; the documented asymmetry comes directly from the outer error-handling
branch and should receive an explicit regression test before that behavior is
intentionally changed.
