---
idd:
  version: "1.0"
  package: pkg/pattern
  document: contract
---

# Contracts: pkg/pattern

## Contract: IdentifierSyntax

**Guarantees:**

The identifier contract exposes two deliberately different levels of
recognition.

Broad recognition through `Patterns`, `GetIdentifierType`, and
`ValidateIDPattern` locates supported identifier text case-insensitively. It is
suitable for scanning content but is not a whole-string acceptance check.

Strict validation through `ValidateIdentifierFormat` expects either a
three-part `TYPE-MODULE-NUMBER` candidate or a two-part internal section marker.
Only `SPEC`, `CONTRACT`, `TEST`, and `DESIGN` are valid IDD types. Module
characters are uppercase ASCII letters, digits, or underscore; number
characters are decimal digits. `PATTERN-*` and `WALK-*` section markers are
recognized only to return an explicit “internal marker” error.

### Outputs and errors

Recognition returns a type name or nil error when any supported pattern is
found. Strict validation returns descriptive errors containing the rejected
identifier and the violated segment rule. No function verifies declaration
existence, package-derived module ownership, or configured regex policy.

### Compatibility

Consumers must choose broad or strict behavior intentionally. Anchoring the
broad regexes, changing accepted types, or changing case rules would alter
collector behavior and requires coordinated parser and validation tests.

## Contract: DocumentReferenceSyntax

**Guarantees:**

`ExtractIDDReferences` returns identifier occurrences that appear adjacent to a
double quote or backtick. Returned values are uppercased, duplicates are
preserved, and ordering across identifier types is unspecified because the
registry is map-backed.

The function is a lexical filter, not a Markdown renderer. It does not prove
balanced delimiters and currently relies on an uppercase lookup when checking
the original text, which limits lowercase quoted input.

## Contract: AnnotationSyntax

**Guarantees:**

Known source annotations are `@implement`, `@test`, and `@test-contract`.
`ExtractAnnotations` accepts a comma-separated list of three-part identifiers
after one of those prefixes and returns the split references in annotation
pattern order and source occurrence order. `SplitAnnotationRefs` trims
whitespace, drops empty comma elements, preserves spelling, and does not
validate each part.

`GetAnnotationType` is an exact, case-sensitive prefix lookup:
`@implement` maps to `SPEC`; both test annotations map to `TEST`. Unknown
prefixes return the empty string.

Annotation placement, declaration kind, TEST kind, configured prefix values,
and doc/code correspondence are outside this contract and belong to collectors
and the engine.
