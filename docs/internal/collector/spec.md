---
markers:
  - id: SPEC-INTERNAL_COLLECTOR-001
    name: Document Collector
  - id: SPEC-INTERNAL_COLLECTOR-002
    name: NewDocCollector
  - id: SPEC-INTERNAL_COLLECTOR-003
    name: Code Collector
  - id: SPEC-INTERNAL_COLLECTOR-004
    name: NewCodeCollector
  - id: SPEC-INTERNAL_COLLECTOR-005
    name: Frontmatter Parsing
  - id: SPEC-INTERNAL_COLLECTOR-006
    name: Frontmatter Validation
  - id: SPEC-INTERNAL_COLLECTOR-007
    name: File Path Validation
  - id: SPEC-INTERNAL_COLLECTOR-008
    name: Title Extraction
  - id: SPEC-INTERNAL_COLLECTOR-009
    name: Module Prefix Validation
  - id: SPEC-INTERNAL_COLLECTOR-010
    name: Document Structure Validation
  - id: SPEC-INTERNAL_COLLECTOR-011
    name: Annotation Extraction
  - id: SPEC-INTERNAL_COLLECTOR-012
    name: Function Context Extraction
  - id: SPEC-INTERNAL_COLLECTOR-013
    name: Code Origin Tracking

  - id: SPEC-INTERNAL_COLLECTOR-017
    name: DocCollector.Collect
  - id: SPEC-INTERNAL_COLLECTOR-024
    name: CodeCollector.Collect
  - id: SPEC-INTERNAL_COLLECTOR-025
    name: SplitAnnotationRefs

related_files:
  spec: docs/internal/collector/spec.md
  contract: docs/internal/collector/contract.md
  design: docs/internal/collector/design.md
  testing: docs/internal/collector/testing.md
---

# Specification (collector)

## SPEC-INTERNAL_COLLECTOR-001: Document Collector

**Contract:** `DocCollector`

**Design:** `CollectorModule`

**Status:** Done

**Requirement:**

The Document Collector (`DocCollector`) must collect IDD identifiers from markdown documentation files, parsing frontmatter markers and extracting identifier references from content.

**Implementation:** `internal/collector/doc_collector.go`

**Key Types:**

- `DocCollector` — Main collector struct with config reference
- `Frontmatter` — YAML frontmatter with marker definitions
- `Marker` — Individual marker with ID, name, and describe fields

**Key Functionality:**

- Recursively walk directories to find `.md` files
- Parse YAML frontmatter to extract defined markers
- Extract IDD references from markdown content using regex
- Validate frontmatter markers match actual content headings
- Validate module prefix matches directory structure
- Report errors for malformed markers (not wrapped in backticks)

**Tests:** `TEST-INTERNAL_COLLECTOR-001`

**Public Functions:**

## SPEC-INTERNAL_COLLECTOR-016: DocCollector.NewDocCollector

**Function Signature:**
`func NewDocCollector(cfg *config.Config) *DocCollector`

**Purpose:** Creates a new DocCollector with the given configuration.

**Parameters:**

- `cfg`: Configuration pointer

**Returns:** A new DocCollector instance

---

## SPEC-INTERNAL_COLLECTOR-017: DocCollector.Collect

**Function Signature:**
`func (c *DocCollector) Collect(ctx context.Context, targetPath string) (*model.IdentifierSet, []*model.ValidationError, error)`

**Purpose:** Collects IDD identifiers from markdown files at the target path. If targetPath is a directory, recursively walks to find all `.md` files.

**Parameters:**

- `ctx`: Context for cancellation
- `targetPath`: Path to markdown file or directory

**Returns:** IdentifierSet with collected identifiers, validation errors, and any error encountered

**Tests:** `TEST-INTERNAL_COLLECTOR-017`

---

## SPEC-INTERNAL_COLLECTOR-018: DocCollector.ParseFrontmatter

**Function Signature:**
`func ParseFrontmatter(content string) (*Frontmatter, error)`

**Purpose:** Parses YAML frontmatter from markdown content. Looks for content between opening `---` and closing `---` markers, excluding code blocks.

**Parameters:**

- `content`: Raw markdown content

**Returns:** Parsed Frontmatter pointer or nil if no frontmatter found, error if parsing fails

---

## SPEC-INTERNAL_COLLECTOR-019: DocCollector.ValidateFrontmatterMarkers

**Function Signature:**
`func ValidateFrontmatterMarkers(fm *Frontmatter, content string, filePath string) []string`

**Purpose:** Validates that frontmatter markers match actual content headings and are properly formatted (wrapped in backticks).

**Parameters:**

- `fm`: Parsed frontmatter
- `content`: Full markdown content
- `filePath`: Path to the file for error messages

**Returns:** List of validation error messages (empty if valid)

---

## SPEC-INTERNAL_COLLECTOR-020: DocCollector.ValidateDocumentStructure

**Function Signature:**
`func ValidateDocumentStructure(filePath string, idType string) error`

**Purpose:** Validates that a document's filename matches its identifier type (e.g., SPEC should be in `spec.md`, TEST in `testing.md`).

**Parameters:**

- `filePath`: Path to the document
- `idType`: Expected identifier type

**Returns:** Error if filename doesn't match expected pattern

---

## SPEC-INTERNAL_COLLECTOR-021: DocCollector.ValidateModulePrefix

**Function Signature:**
`func ValidateModulePrefix(id string, filePath string) error`

**Purpose:** Validates that the module part of an identifier matches the directory structure (e.g., an identifier with module "BE" should be in docs/backend/).

**Parameters:**

- `id`: The identifier to validate
- `filePath`: Path to the file containing the identifier

**Returns:** Error if module prefix doesn't match directory

---

## SPEC-INTERNAL_COLLECTOR-022: DocCollector.GetExpectedFilename

**Function Signature:**
`func GetExpectedFilename(idType string) string`

**Purpose:** Returns the expected filename for a given identifier type.

**Parameters:**

- `idType`: Identifier type (SPEC, CONTRACT, TEST, DESIGN)

**Returns:** Expected filename (spec.md, contract.md, testing.md, design.md) or empty string if unknown

---

**Acceptance Criteria:**

- [x] Collects identifiers from markdown files in directory trees
- [x] Parses frontmatter markers correctly
- [x] Detects defined markers vs referenced markers
- [x] Extracts title from heading containing identifier
- [x] Validates module prefix consistency
- [x] Ignores paths configured in `ignore_paths`
- [x] Reports errors for bare markers (not wrapped in backticks)

**Tests:** `TEST-INTERNAL_COLLECTOR-001`

**Related:** `CON-INTERNAL_COLLECTOR-001`

---

## SPEC-INTERNAL_COLLECTOR-002: Code Collector

**Contract:** `CodeCollector`

**Design:** `CollectorModule`

**Status:** Done

**Requirement:**

The Code Collector (`CodeCollector`) must collect IDD annotations from source code files (Go, TypeScript, JavaScript), extracting `@implement`, `@test`, and `@test-contract` annotations.

**Implementation:** `internal/collector/code_collector.go`

**Key Types:**

- `CodeCollector` — Main collector struct with config reference

**Key Functionality:**

- Walk directory trees to find source files (`.go`, `.ts`, `.tsx`, `.js`)
- Extract annotations using configured patterns
- Extract function context (function name, preceding comments)
- Build identifiers from annotations with code location
- Set origin to `model.OriginCode` for code-based identifiers

**Tests:** `TEST-INTERNAL_COLLECTOR-002`

**Public Functions:**

## SPEC-INTERNAL_COLLECTOR-023: CodeCollector.NewCodeCollector

**Function Signature:**
`func NewCodeCollector(cfg *config.Config) *CodeCollector`

**Purpose:** Creates a new CodeCollector with the given configuration.

**Parameters:**

- `cfg`: Configuration pointer

**Returns:** A new CodeCollector instance

---

## SPEC-INTERNAL_COLLECTOR-024: CodeCollector.Collect

**Function Signature:**
`func (c *CodeCollector) Collect(ctx context.Context, targetPath string) (*model.IdentifierSet, error)`

**Purpose:** Collects IDD identifiers from code annotations in source files at the target path. If targetPath is a directory, recursively walks to find all `.go`, `.ts`, `.tsx`, `.js` files.

**Parameters:**

- `ctx`: Context for cancellation
- `targetPath`: Path to source file or directory

**Returns:** IdentifierSet with collected identifiers from code, and any error encountered

**Tests:** `TEST-INTERNAL_COLLECTOR-024`

---

## SPEC-INTERNAL_COLLECTOR-025: CodeCollector.SplitAnnotationRefs

**Function Signature:**
`func SplitAnnotationRefs(s string) []string`

**Purpose:** Splits comma-separated IDD references from an annotation and trims whitespace. Used to handle multiple references in a single annotation like `@implement` `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-002`.

**Parameters:**

- `s`: Comma-separated identifier references

**Returns:** Slice of individual trimmed identifier references

---

**Acceptance Criteria:**

- [x] Supports Go, TypeScript, and JavaScript files
- [x] Extracts `@implement`, `@test`, `@test-contract` annotations
- [x] Captures function name and preceding comments as context
- [x] Reports file path and line number for each annotation
- [x] Ignores paths configured in `ignore_paths`
- [x] Handles multiple annotations on same line

**Tests:** `TEST-INTERNAL_COLLECTOR-001`, `TEST-INTERNAL_COLLECTOR-002`, `TEST-INTERNAL_COLLECTOR-003`, `TEST-INTERNAL_COLLECTOR-004`, `TEST-INTERNAL_COLLECTOR-005`, `TEST-INTERNAL_COLLECTOR-006`, `TEST-INTERNAL_COLLECTOR-007`, `TEST-INTERNAL_COLLECTOR-008`, `TEST-INTERNAL_COLLECTOR-009`, `TEST-INTERNAL_COLLECTOR-010`, `TEST-INTERNAL_COLLECTOR-011`, `TEST-INTERNAL_COLLECTOR-012`, `TEST-INTERNAL_COLLECTOR-013`, `TEST-INTERNAL_COLLECTOR-014`, `TEST-INTERNAL_COLLECTOR-015`, `TEST-INTERNAL_COLLECTOR-016`, `TEST-INTERNAL_COLLECTOR-017`, `TEST-INTERNAL_COLLECTOR-021`, `TEST-INTERNAL_COLLECTOR-022`, `TEST-INTERNAL_COLLECTOR-023`, `TEST-INTERNAL_COLLECTOR-024`, `TEST-INTERNAL_COLLECTOR-025`, `TEST-INTERNAL_COLLECTOR-026`, `TEST-INTERNAL_COLLECTOR-027`

---

## SPEC-INTERNAL_COLLECTOR-003: Frontmatter Parsing

**Contract:** `ParseFrontmatter`

**Design:** `CollectorModule`

**Status:** Done

**Requirement:**

Frontmatter parsing must handle various YAML structures including markers with id, name, and describe fields.

**Implementation:** `internal/collector/frontmatter.go`

**Tests:** `TEST-INTERNAL_COLLECTOR-003`, `TEST-INTERNAL_COLLECTOR-004`

---

## SPEC-INTERNAL_COLLECTOR-004: Frontmatter Validation

**Contract:** `ValidateFrontmatterMarkers`

**Design:** `CollectorModule`

**Status:** Done

**Requirement:**

Frontmatter validation must check that markers in YAML match the actual document headings.

**Implementation:** `internal/collector/frontmatter.go`

**Tests:** `TEST-INTERNAL_COLLECTOR-005`, `TEST-INTERNAL_COLLECTOR-006`

---

## SPEC-INTERNAL_COLLECTOR-005: File Path Validation

**Contract:** `ValidateDocumentStructure`

**Design:** `CollectorModule`

**Status:** Done

**Requirement:**

File path validation ensures proper document structure and naming conventions.

**Implementation:** `internal/collector/doc_collector.go`

**Tests:** `TEST-INTERNAL_COLLECTOR-007`, `TEST-INTERNAL_COLLECTOR-008`

---

## SPEC-INTERNAL_COLLECTOR-006: Title Extraction

**Contract:** `ExtractTitle`

**Design:** `CollectorModule`

**Status:** Done

**Requirement:**

Extract document titles from markdown headings containing identifiers.

**Implementation:** `internal/collector/doc_collector.go`

**Tests:** `TEST-INTERNAL_COLLECTOR-009`, `TEST-INTERNAL_COLLECTOR-010`

---

## SPEC-INTERNAL_COLLECTOR-007: Module Prefix Validation

**Contract:** `ValidateModulePrefix`

**Design:** `CollectorModule`

**Status:** Done

**Requirement:**

Validate that identifier module prefixes match directory structure.

**Implementation:** `internal/collector/doc_collector.go`

**Tests:** `TEST-INTERNAL_COLLECTOR-011`, `TEST-INTERNAL_COLLECTOR-012`

---

## SPEC-INTERNAL_COLLECTOR-008: Document Structure Validation

**Contract:** `GetExpectedFilename`

**Design:** `CollectorModule`

**Status:** Done

**Requirement:**

Validate document structure including expected filenames by identifier type.

**Implementation:** `internal/collector/doc_collector.go`

**Tests:** `TEST-INTERNAL_COLLECTOR-019`, `TEST-INTERNAL_COLLECTOR-014`

---

## SPEC-INTERNAL_COLLECTOR-009: Annotation Extraction

**Contract:** `ExtractAnnotations`

**Design:** `CollectorModule`

**Status:** Done

**Requirement:**

Extract IDD annotations from code including @implement, @test, @test-contract.

**Implementation:** `internal/collector/code_collector.go`

**Tests:** `TEST-INTERNAL_COLLECTOR-015`, `TEST-INTERNAL_COLLECTOR-024`

---

## SPEC-INTERNAL_COLLECTOR-010: Function Context Extraction

**Contract:** `ExtractFunctionContext`

**Design:** `CollectorModule`

**Status:** Done

**Requirement:**

Extract function name and preceding comments as context for code annotations.

**Implementation:** `internal/collector/code_collector.go`

**Tests:** `TEST-INTERNAL_COLLECTOR-016`, `TEST-INTERNAL_COLLECTOR-025`

---

## SPEC-INTERNAL_COLLECTOR-011: Code Origin Tracking

**Contract:** `SetOrigin`

**Design:** `CollectorModule`

**Status:** Done

**Requirement:**

Set origin to OriginCode for code-based identifiers to distinguish from doc origins.

**Implementation:** `internal/collector/code_collector.go`

**Tests:** `TEST-INTERNAL_COLLECTOR-017`

---

## SPEC-INTERNAL_COLLECTOR-012: Multi-Annotation Handling

**Contract:** `SplitAnnotationRefs`

**Design:** `CollectorModule`

**Status:** Done

**Requirement:**

Handle multiple annotations on the same line or multiple references in single annotation.

**Implementation:** `internal/collector/code_collector.go`

**Tests:** `TEST-INTERNAL_COLLECTOR-024`

---

## SPEC-INTERNAL_COLLECTOR-013: Path Ignore Patterns

**Contract:** `ShouldIgnore`

**Design:** `CollectorModule`

**Status:** Done

**Requirement:**

Respect ignore_paths configuration when collecting from directories.

**Implementation:** `internal/collector/doc_collector.go`, `internal/collector/code_collector.go`

**Tests:** `TEST-INTERNAL_COLLECTOR-013`

---

