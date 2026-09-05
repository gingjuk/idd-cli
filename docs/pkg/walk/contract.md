---
idd:
  version: "1.1"
  package: pkg/walk
  namespace: PKG_WALK
---

# Contracts: pkg/walk

## Contract: FileTraversal

**Guarantees:**

`FileTraversal` provides synchronous visitation of entries selected by one or
more platform-native glob patterns.

```go
type FileVisitor func(path string, info os.FileInfo) error

func Walk(patterns []string, visitor FileVisitor) error
```

### Inputs and ownership

`patterns` is read in caller-provided order and is not retained. Each element
uses `filepath.Glob` syntax on the current platform. A nil or empty slice is a
valid no-op. The visitor must be non-nil when any match can occur.

The callback receives the matched path string and the `os.FileInfo` obtained
for that path. Those values are borrowed for the duration of the call. The
visitor may inspect the filesystem, but mutation during traversal can make
later metadata stale or paths disappear.

### Ordering and duplicate behavior

Direct matches are visited in pattern order and in the order returned by
`filepath.Glob`. A directly matched directory is visited before its
descendants. The same exact path string is visited at most once per `Walk`
invocation, including overlaps between direct patterns and recursive descent.

### Errors and side effects

The function itself performs reads only. It returns a visitor error raised for
a direct glob match without wrapping it. Invalid glob patterns, failed direct
`Stat` operations, and filesystem walk errors are skipped.

### Non-guarantees

The contract does not promise canonical-file uniqueness across symlinks,
fail-fast handling for every filesystem error, callback safety for a nil
visitor, or a stable ordering beyond the platform-native glob and walk order
described above.

### Known limitations

If a descendant visitor returns an error during recursive descent, that
directory walk stops, but the outer pattern loop currently suppresses the
returned walk error. The package also does not guard a nil callback. These are
current implementation limitations, not compatibility promises; callers that
require fail-fast traversal or callback validation must provide that boundary
themselves until the implementation changes.

### Compatibility commitments

Changing glob syntax, canonicalizing deduplication keys, or parallelizing
callback execution would change stable observable behavior and requires
explicit specification and test updates. Propagating recursive visitor errors
or returning a defined error for a nil visitor may be introduced as a defect
fix without treating the present limitation as a promise; that change still
requires focused tests and release communication.

## Contract: ExtensionMatch

**Guarantees:**

```go
func MatchAnyExtensions(path string, extensions []string) bool
```

The function compares `filepath.Ext(path)` against each supplied extension
using exact, case-sensitive string equality. Extensions therefore include the
leading dot, such as `.go`. An empty list, a path without an extension, or a
multi-suffix path whose final suffix is not listed returns false. The function
performs no filesystem access and retains no input.
