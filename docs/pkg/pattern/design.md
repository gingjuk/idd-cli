---
idd:
  version: "1.0"
  package: pkg/pattern
  document: design
---

# Design: pkg/pattern

## Component: SyntaxRecognizer

**Purpose:**

`SyntaxRecognizer` centralizes the lexical rules shared by documentation and
source collectors. It recognizes identifier-shaped text, extracts references
from explicit Markdown quoting, parses annotation lists, maps annotation
prefixes to identifier kinds, and performs strict identifier-format checks.

The component is deliberately lexical. It can say that text has a supported
shape, but it cannot establish that an identifier is declared, belongs to the
current package, uses the configured project pattern, or appears on an allowed
Go declaration. Collectors and the engine own those contextual decisions.

### Responsibilities

- compile the built-in identifier and annotation regular expressions once;
- distinguish broad recognition from strict format validation;
- require quotes or backticks before treating prose identifiers as document
  references;
- split comma-separated annotation targets without inventing relationships;
- expose the identifier kind associated with known syntax; and
- produce human-readable format errors for invalid three-part identifiers.

### Boundaries and known limitations

The `Patterns` registry is a Go map, so multi-type extraction order is not a
stable public ordering. Extraction also preserves duplicate occurrences.
`ExtractIDDReferences` uppercases the regex match before locating its quote in
the original content; quoted lowercase identifiers are therefore not
recognized reliably by the current implementation.

`ValidateIDPattern` and `GetIdentifierType` use unanchored regular-expression
matching and answer whether supported identifier text occurs in the input.
`ValidateIdentifierFormat` is the strict whole-string validator and should be
used when accepting an authored identifier. Even that strict validator
currently checks characters rather than explicitly rejecting an empty module
or number segment.

These details are documented as implementation boundaries, not recommended
authoring style. Tightening them requires regression tests because collectors
may currently depend on the permissive recognition path.

### Decisions and trade-offs

The package keeps built-in patterns separate from `.idd.yaml`. The config tells
collectors which annotations and documentation patterns to use; this package
defines the baseline grammar needed by internal validation and test fixtures.
Case-insensitive recognition is convenient when scanning text, while strict
module validation retains uppercase spelling to keep canonical identifiers
stable.

## Architecture

```text
Markdown text ── ExtractIDDReferences ──> quoted identifier references

source comments ── AnnotationPatterns
                └─ ExtractAnnotations ──> comma-split targets

candidate identifier
    ├─ GetIdentifierType / ValidateIDPattern  (broad recognition)
    └─ ValidateIdentifierFormat              (strict three-part grammar)
```

Recognition functions do not access files and retain no per-call state.
Compiled regular expressions are immutable package-level values and are safe
for concurrent reads.

## Package Layout

`idd.go` contains the registries and their parsing helpers in one place so the
grammar is visible as a coherent policy. File discovery remains in `pkg/walk`;
document structure parsing remains in `internal/collector`; relationship and
placement validation remain in `internal/engine`.

## Function Composition

`ExtractAnnotations` checks each annotation pattern in a deterministic slice
order and delegates list parsing to `SplitAnnotationRefs`.
`ExtractIDDReferences` checks every identifier pattern and delegates quote
detection to `isQuoted`. Strict validation splits the candidate into segments,
then dispatches to the section-marker or IDD-identifier validator.

## Dependencies

Only `fmt`, `regexp`, and `strings` are required. Unicode case normalization is
provided by the regular-expression engine, while module-character validation is
explicitly limited to ASCII uppercase letters, digits, and underscore.

## Testability Hooks

All behavior is pure and table-driven. Tests can supply text directly and
compare returned identifiers or error presence without filesystem fixtures.
The tables cover valid and invalid segment shapes, annotation kinds, reference
quoting, comma splitting, and multiple annotation forms.

The suite does not currently prove deterministic cross-type extraction order,
deduplication, quoted lowercase handling, or rejection of empty strict
segments. Those limitations must remain visible until dedicated behavior is
specified and tested.
