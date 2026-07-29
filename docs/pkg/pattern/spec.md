---
idd:
  version: "1.0"
  package: pkg/pattern
---

# Specifications: pkg/pattern

## SPEC-PKG_PATTERN-001: Built-in identifier recognition registry

- **Design:** `SyntaxRecognizer`
- **Contract:** `IdentifierSyntax`

**Requirement:**

The package must provide precompiled case-insensitive recognition patterns for
SPEC, CONTRACT, TEST, DESIGN, and the internal PATTERN and WALK forms. Project
module segments may contain uppercase letters, digits, and underscores.

### Boundary and rationale

The registry supports fast lexical scanning and examples; it is not the
project's configurable documentation source and is not a declaration catalog.
Map iteration order is intentionally not an output guarantee.

### Acceptance evidence

**Acceptance:**

Pattern-validation tests distinguish supported forms from unknown types, and
strict-format tests separately reject internal section markers as authored IDD
identifiers.

## SPEC-PKG_PATTERN-002: Source annotation recognition registry

- **Design:** `SyntaxRecognizer`
- **Contract:** `AnnotationSyntax`

**Requirement:**

The package must recognize `@implement`, `@test`, and `@test-contract`
followed by one or more comma-separated three-part identifiers and associate
the implementation prefix with SPEC and the two testing prefixes with TEST.

**Acceptance:** Table-driven prefix tests map the three exact supported
spellings to SPEC or TEST and reject differently cased, unknown, and empty
prefixes. Extraction cases cover one annotation, mixed kinds, comma-separated
targets, and content with no annotation.

### Edge cases

Whitespace around commas is accepted. A missing identifier is not matched by
the extraction regex and is diagnosed later by engine annotation validation.
The registry does not decide whether a recognized annotation is placed on a
valid declaration.

## SPEC-PKG_PATTERN-003: Annotation prefix type mapping

- **Design:** `SyntaxRecognizer`
- **Contract:** `AnnotationSyntax`

**Requirement:**

Exact known annotation prefixes must map to their identifier type; unknown or
differently cased prefixes must return no type.

### Acceptance evidence

**Acceptance:**

Table-driven tests include all three prefixes, an uppercase variation, an
unknown value, and empty input.

## SPEC-PKG_PATTERN-004: Explicit document reference extraction

- **Design:** `SyntaxRecognizer`
- **Contract:** `DocumentReferenceSyntax`

**Requirement:**

Documentation scanning must only return identifier-shaped prose references when
the occurrence is explicitly quoted or backtick-delimited, and returned
identifiers must be normalized to uppercase.

**Acceptance:** Table-driven extraction compares exact returned slices for
backtick-delimited, double-quoted, unquoted, empty, and partial-match inputs.
Plain identifier-shaped prose yields no reference, while accepted references
retain the current normalization and adjacency behavior.

### Failure and implementation boundary

Plain prose must not create graph edges accidentally. Duplicate occurrences
remain duplicates, balanced Markdown delimiters are not parsed, and ordering
between identifier types is not stable. Quoted lowercase references are a known
current limitation because quote detection searches with the normalized value.

## SPEC-PKG_PATTERN-005: Comma-separated annotation target splitting

- **Design:** `SyntaxRecognizer`
- **Contract:** `AnnotationSyntax`

**Requirement:**

Annotation target lists must split on commas, trim surrounding whitespace,
preserve target spelling, and omit empty elements.

**Acceptance:** One target, several targets with surrounding whitespace,
trailing commas, repeated spellings, and empty input produce the expected
slices. Returned values keep their original case and order.

### Non-goals

Splitting does not validate identifier grammar, remove duplicates, or accept
whitespace-only separation as multiple references. Those checks occur after
lexical extraction.

## SPEC-PKG_PATTERN-006: Broad identifier type detection

- **Design:** `SyntaxRecognizer`
- **Contract:** `IdentifierSyntax`

**Requirement:**

Given arbitrary text, the recognizer must return the first registry type whose
regular expression finds supported identifier text, or the empty string when
none is present.

**Acceptance:** Every registered identifier family, mixed-case type text, and
identifiers embedded in surrounding prose return their registered type.
Unknown text and empty input return the empty string.

### Boundary

The input is not required to consist solely of the identifier. Callers that
accept authored IDs must use strict validation instead.

## SPEC-PKG_PATTERN-007: Broad pattern validity check

- **Design:** `SyntaxRecognizer`
- **Contract:** `IdentifierSyntax`

**Requirement:**

The broad validity check must succeed when any built-in identifier pattern is
found and return an error naming the input otherwise.

**Acceptance:** Representative SPEC, CONTRACT, TEST, DESIGN, PATTERN, and WALK
text is accepted even inside surrounding content. Unknown and empty values
return a non-nil error containing the rejected input.

### Non-goals

This check does not enforce whole-string format, package module ownership, or
whether PATTERN and WALK forms are valid author-facing IDs.

## SPEC-PKG_PATTERN-008: Annotation target extraction

- **Design:** `SyntaxRecognizer`
- **Contract:** `AnnotationSyntax`

**Requirement:**

Annotation extraction must scan for every known prefix, split each captured
comma list, and return the targets without assigning graph relationships or
discarding duplicates.

### Acceptance evidence

**Acceptance:**

Tests cover single annotations, mixed annotation types, comma-separated
targets, and content containing no annotations.

## SPEC-PKG_PATTERN-009: Strict identifier format validation

- **Design:** `SyntaxRecognizer`
- **Contract:** `IdentifierSyntax`

**Requirement:**

Strict validation must reject candidates with the wrong segment count,
unsupported IDD types, lowercase or punctuation-bearing modules, non-numeric
number characters, and internal two-part section markers.

### Known boundary

The current implementation validates the characters found in module and number
segments but does not explicitly reject an empty segment. This document does
not claim a stronger guarantee than the code provides.

### Acceptance evidence

**Acceptance:**

The table-driven validator suite covers all supported IDD types, compound
underscore modules, invalid segment counts, unsupported types, lowercase
modules, punctuation, non-digit numbers, empty overall input, and
case-insensitive type spelling.
