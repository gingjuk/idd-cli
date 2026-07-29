---
idd:
  version: "1.0"
  package: pkg/walk
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
valid no-op. The visitor must be non-nil when any match can occur; the package
does not guard a nil callback.

The callback receives the matched path string and the `os.FileInfo` obtained
for that path. Those values are borrowed for the duration of the call. The
visitor may inspect the filesystem, but mutation during traversal can make
later metadata stale or paths disappear.

### Ordering and duplicate behavior

Direct matches are visited in pattern order and in the order returned by
`filepath.Glob`. A directly matched directory is visited before its
descendants. The same exact path string is visited at most once per `Walk`
invocation, including overlaps between direct patterns and recursive descent.
No promise is made for canonical-file uniqueness across symlinks or lexically
different paths.

### Errors and side effects

The function itself performs reads only. It returns a visitor error raised for
a direct glob match without wrapping it. Invalid glob patterns, failed direct
`Stat` operations, and filesystem walk errors are skipped.

If a descendant visitor returns an error during recursive descent, that
directory walk stops, but the outer pattern loop currently suppresses the
returned walk error. Callers requiring fail-fast behavior for every descendant
must not assume it from this contract.

### Compatibility

Changing glob syntax, canonicalizing deduplication keys, parallelizing callback
execution, or propagating recursive visitor failures would change observable
ordering or failure behavior. Such changes require explicit specification and
test updates.

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
