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

**Status:** Done

**Contract:** `IDDPattern`
**Design:** `PatternModule`

**Requirement:**

Define regex patterns for IDD identifier recognition.

**Implementation:** `pkg/pattern/idd.go`

**Key Patterns:**

- SPEC pattern: `SPEC-[A-Z]+-[0-9]+`
- TEST pattern: `TEST-[A-Z]+-[0-9]+`
- CONTRACT pattern: `CONTRACT-[A-Z]+-[0-9]+`
- DESIGN pattern: `DESIGN-[A-Z]+-[0-9]+`
**Tests:** `TEST-PKG_PATTERN-001`

---

## SPEC-PKG_PATTERN-002: Code Annotation Pattern

**Status:** Done

**Contract:** `AnnotationPattern`
**Design:** `PatternModule`

**Requirement:**

Define regex patterns for code annotations (@implement, @test, @test-contract).

**Implementation:** `pkg/pattern/idd.go`

**Key Patterns:**

- Annotation prefix pattern
- Annotation value pattern
**Tests:** `TEST-PKG_PATTERN-002`

---

## SPEC-PKG_PATTERN-003: Annotation Type Mapping

**Status:** Done

**Contract:** `GetAnnotationType`
**Design:** `PatternModule`

**Requirement:**

Map annotation prefixes to identifier types.

**Implementation:** `pkg/pattern/idd.go`

**Key Functions:**

- `GetAnnotationType(prefix string) string`
- Maps: @implement → SPEC, @test → TEST, @test-contract → TEST
**Tests:** `TEST-PKG_PATTERN-003`

---

## SPEC-PKG_PATTERN-004: Reference Extraction

**Status:** Done

**Contract:** `ExtractIDDReferences`
**Design:** `PatternModule`

**Requirement:**

Extract IDD references from content.

**Implementation:** `pkg/pattern/idd.go`

**Key Functions:**

- `ExtractIDDReferences(content string) []string`
- Extracts all IDD identifier references from text
**Tests:** `TEST-PKG_PATTERN-004`

---

## SPEC-PKG_PATTERN-005: Reference Splitting

**Status:** Done

**Contract:** `SplitAnnotationRefs`
**Design:** `PatternModule`

**Requirement:**

Split comma-separated IDD references.

**Implementation:** `pkg/pattern/idd.go`

**Key Functions:**

- `SplitAnnotationRefs(s string) []string`
- Trims whitespace from each reference
**Tests:** `TEST-PKG_PATTERN-005`

---

## SPEC-PKG_PATTERN-006: Identifier Type Detection

**Status:** Done

**Contract:** `GetIdentifierType`
**Design:** `PatternModule`

**Requirement:**

Determine identifier type from reference string.

**Implementation:** `pkg/pattern/idd.go`

**Key Functions:**

- `GetIdentifierType(ref string) string`
- Returns "SPEC", "TEST", "CONTRACT", "DESIGN" or empty
**Tests:** `TEST-PKG_PATTERN-006`

---

## SPEC-PKG_PATTERN-007: Pattern Validation

**Status:** Done

**Contract:** `ValidateIDPattern`
**Design:** `PatternModule`

**Requirement:**

Validate identifier against IDD pattern.

**Implementation:** `pkg/pattern/idd.go`

**Key Functions:**

- `ValidateIDPattern(id string) error`
- Validates against known patterns
**Tests:** `TEST-PKG_PATTERN-007`

---

## SPEC-PKG_PATTERN-008: IDD Reference Filter

**Status:** Done

**Contract:** `ExtractAnnotations`
**Design:** `PatternModule`

**Requirement:**

Filter out quoted/backtick-wrapped identifiers.

**Implementation:** `pkg/pattern/idd.go`

**Key Functions:**

- `ExtractAnnotations(content string) []string`
- Filters out identifiers in backticks or quotes
**Tests:** `TEST-PKG_PATTERN-008`

---

## SPEC-PKG_PATTERN-009: Pattern Constants

**Status:** Done

**Contract:** `PatternConstants`
**Design:** `PatternModule`

**Requirement:**

Define pattern constants for reuse.

**Implementation:** `pkg/pattern/idd.go`

**Key Constants:**

- Pattern strings
- Annotation prefixes

**Tests:** `TEST-PKG_PATTERN-007`
