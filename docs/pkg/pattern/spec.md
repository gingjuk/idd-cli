---
markers:
  - id: SPEC-PKG_PATTERN-001
    name: IDD Pattern Regex
  - id: SPEC-PKG_PATTERN-002
    name: Code Annotation Pattern
  - id: SPEC-PKG_PATTERN-003
    name: Annotation Type Mapping
  - id: SPEC-PKG_PATTERN-004
    name: Reference Extraction
  - id: SPEC-PKG_PATTERN-005
    name: Reference Splitting
  - id: SPEC-PKG_PATTERN-006
    name: Identifier Type Detection
  - id: SPEC-PKG_PATTERN-007
    name: Pattern Validation
  - id: SPEC-PKG_PATTERN-008
    name: IDD Reference Filter
  - id: SPEC-PKG_PATTERN-009
    name: Pattern Constants

related_files:
  spec: docs/pkg/pattern/spec.md
  contract: docs/pkg/pattern/contract.md
  design: docs/pkg/pattern/design.md
  testing: docs/pkg/pattern/testing.md

---

# Specification (pattern)

## SPEC-PKG_PATTERN-001: IDD Pattern Regex

**Design:** `PatternModule`

**Contract:** `IDDPattern`

**Requirement:**

Define regex patterns for IDD identifier recognition.

**Tests:** `TEST-PKG_PATTERN-001`

**Status:** Done

**Implementation:** `pkg/pattern/idd.go`

**Key Patterns:**

- SPEC pattern: `SPEC-[A-Z]+-[0-9]+`
- TEST pattern: `TEST-[A-Z]+-[0-9]+`
- CONTRACT pattern: `CONTRACT-[A-Z]+-[0-9]+`
- DESIGN pattern: `DESIGN-[A-Z]+-[0-9]+`

---

## SPEC-PKG_PATTERN-002: Code Annotation Pattern

**Design:** `PatternModule`

**Contract:** `AnnotationPattern`

**Requirement:**

Define regex patterns for code annotations (@implement, @test, @test-contract).

**Tests:** `TEST-PKG_PATTERN-002`

**Status:** Done

**Implementation:** `pkg/pattern/idd.go`

**Key Patterns:**

- Annotation prefix pattern
- Annotation value pattern

---

## SPEC-PKG_PATTERN-003: Annotation Type Mapping

**Design:** `PatternModule`

**Contract:** `GetAnnotationType`

**Requirement:**

Map annotation prefixes to identifier types.

**Tests:** `TEST-PKG_PATTERN-003`

**Status:** Done

**Implementation:** `pkg/pattern/idd.go`

**Key Functions:**

- `GetAnnotationType(prefix string) string`
- Maps: @implement → SPEC, @test → TEST, @test-contract → TEST

---

## SPEC-PKG_PATTERN-004: Reference Extraction

**Design:** `PatternModule`

**Contract:** `ExtractIDDReferences`

**Requirement:**

Extract IDD references from content.

**Tests:** `TEST-PKG_PATTERN-004`

**Status:** Done

**Implementation:** `pkg/pattern/idd.go`

**Key Functions:**

- `ExtractIDDReferences(content string) []string`
- Extracts all IDD identifier references from text

---

## SPEC-PKG_PATTERN-005: Reference Splitting

**Design:** `PatternModule`

**Contract:** `SplitAnnotationRefs`

**Requirement:**

Split comma-separated IDD references.

**Tests:** `TEST-PKG_PATTERN-005`

**Status:** Done

**Implementation:** `pkg/pattern/idd.go`

**Key Functions:**

- `SplitAnnotationRefs(s string) []string`
- Trims whitespace from each reference

---

## SPEC-PKG_PATTERN-006: Identifier Type Detection

**Design:** `PatternModule`

**Contract:** `GetIdentifierType`

**Requirement:**

Determine identifier type from reference string.

**Tests:** `TEST-PKG_PATTERN-006`

**Status:** Done

**Implementation:** `pkg/pattern/idd.go`

**Key Functions:**

- `GetIdentifierType(ref string) string`
- Returns "SPEC", "TEST", "CONTRACT", "DESIGN" or empty

---

## SPEC-PKG_PATTERN-007: Pattern Validation

**Design:** `PatternModule`

**Contract:** `ValidateIDPattern`

**Requirement:**

Validate identifier against IDD pattern.

**Tests:** `TEST-PKG_PATTERN-007`

**Status:** Done

**Implementation:** `pkg/pattern/idd.go`

**Key Functions:**

- `ValidateIDPattern(id string) error`
- Validates against known patterns

---

## SPEC-PKG_PATTERN-008: IDD Reference Filter

**Design:** `PatternModule`

**Contract:** `ExtractAnnotations`

**Requirement:**

Filter out quoted/backtick-wrapped identifiers.

**Tests:** `TEST-PKG_PATTERN-008`

**Status:** Done

**Implementation:** `pkg/pattern/idd.go`

**Key Functions:**

- `ExtractAnnotations(content string) []string`
- Filters out identifiers in backticks or quotes

---

## SPEC-PKG_PATTERN-009: Pattern Constants

**Design:** `PatternModule`

**Contract:** `PatternConstants`

**Requirement:**

Define pattern constants for reuse.

**Tests:** `TEST-PKG_PATTERN-007`
**Status:** Done

**Implementation:** `pkg/pattern/idd.go`

**Key Constants:**

- Pattern strings
- Annotation prefixes
