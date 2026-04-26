---
markers:
  - id: TEST-INTERNAL_COLLECTOR-001
    name: Document Collector Tests
  - id: TEST-INTERNAL_COLLECTOR-002
    name: Code Collector Tests
  - id: TEST-INTERNAL_COLLECTOR-003
    name: Frontmatter Parsing Tests
  - id: TEST-INTERNAL_COLLECTOR-004
    name: Frontmatter Validation Tests
  - id: TEST-INTERNAL_COLLECTOR-005
    name: File Path Validation Tests
  - id: TEST-INTERNAL_COLLECTOR-006
    name: Title Extraction Tests
  - id: TEST-INTERNAL_COLLECTOR-007
    name: Module Prefix Tests
  - id: TEST-INTERNAL_COLLECTOR-008
    name: Document Structure Tests
  - id: TEST-INTERNAL_COLLECTOR-009
    name: Annotation Extraction Tests
  - id: TEST-INTERNAL_COLLECTOR-010
    name: Function Context Tests
  - id: TEST-INTERNAL_COLLECTOR-011
    name: Code Origin Tests
  - id: TEST-INTERNAL_COLLECTOR-012
    name: Multi-Annotation Tests
  - id: TEST-INTERNAL_COLLECTOR-013
    name: Path Ignore Tests
  - id: TEST-INTERNAL_COLLECTOR-019
    name: Document Type Matching Tests
  - id: TEST-INTERNAL_COLLECTOR-014
    name: Code File Discovery Tests
  - id: TEST-INTERNAL_COLLECTOR-015
    name: Language Support Tests
  - id: TEST-INTERNAL_COLLECTOR-016
    name: Collector Integration Tests
  - id: TEST-INTERNAL_COLLECTOR-017
    name: Edge Case Tests
  - id: TEST-INTERNAL_COLLECTOR-018
    name: Error Handling Tests
  - id: TEST-INTERNAL_COLLECTOR-020
    name: Marker Extraction Tests
  - id: TEST-INTERNAL_COLLECTOR-021
    name: IDD Reference Tests
  - id: TEST-INTERNAL_COLLECTOR-022
    name: Code Collector Basic Collection
  - id: TEST-INTERNAL_COLLECTOR-023
    name: Code Collector Go File
  - id: TEST-INTERNAL_COLLECTOR-024
    name: Code Collector Multiple Annotations
  - id: TEST-INTERNAL_COLLECTOR-025
    name: Code Collector Function Context
  - id: TEST-INTERNAL_COLLECTOR-026
    name: Code Collector Multi-File
  - id: TEST-INTERNAL_COLLECTOR-027
    name: Code Collector Integration
  - id: TEST-INTERNAL_COLLECTOR-028
    name: Doc Collector Section Title Extraction
  - id: TEST-INTERNAL_COLLECTOR-029
    name: Doc Collector Multi-Section Handling

related_files:
  spec: docs/internal/collector/spec.md
  contract: docs/internal/collector/contract.md
  design: docs/internal/collector/design.md
  testing: docs/internal/collector/testing.md
---

# Test Cases

## Test

## TEST-INTERNAL_COLLECTOR-001: Document Collector Tests

**Status:** Done

**Purpose:**

Test the `DocCollector` for collecting identifiers from markdown documentation.

### Test Cases

| Test | Description |
| --- | ----------- |
| `TestDocCollector_Collect` | Tests collecting from directory |
| `TestDocCollector_Collect_WithFrontmatter` | Tests frontmatter parsing and marker validation |
| `TestDocCollector_Collect_FileNotFound` | Tests handling of missing file |
| `TestDocCollector_Collect_NonMarkdownFile` | Tests that non-.md files are skipped |
| `TestDocCollector_extractTitle` | Tests title extraction from heading |
| `TestDocCollector_extractTitle_NotFound` | Tests title extraction when no heading found |
| `TestDocCollector_Collect_DirectoryWithNoSpec` | Tests directory with no spec files |
| `TestNewDocCollector` | Tests collector initialization |

#### Frontmatter Tests

| Test | Description |
| --- | ----------- |
| `TestParseFrontmatter` | Tests frontmatter YAML parsing |
| `TestParseFrontmatter/valid_frontmatter` | Valid frontmatter with markers |
| `TestParseFrontmatter/no_frontmatter` | No frontmatter returns nil |
| `TestParseFrontmatter/empty_content` | Empty content handled gracefully |
| `TestParseFrontmatter/code_block_with_dashes` | Code blocks don't interfere |
| `TestParseFrontmatter/frontmatter_with_describe` | Marker with describe field |

#### Validation Tests

| Test | Description |
| --- | ----------- |
| `TestValidateFrontmatterMarkers` | Tests marker-to-heading matching |
| `TestExtractDefinedMarkers` | Tests extraction of defined markers |
| `TestExtractReferencedMarkers` | Tests extraction of referenced markers |
| `TestExtractIDDRefs` | Tests IDD reference regex |
| `TestGetExpectedFilename` | Tests filename from type mapping |
| `TestValidateDocumentStructure` | Tests doc structure validation |
| `TestIsRootDocFile` | Tests root doc detection |
| `TestExtractModuleName` | Tests module name extraction |
| `TestIsKnownAbbreviation` | Tests abbreviation validation |
| `TestValidateModulePrefix` | Tests prefix matching |

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-002`

---

## TEST-INTERNAL_COLLECTOR-002: Code Collector Tests

**Status:** Done

**Purpose:**

Test the `CodeCollector` for collecting annotations from source code.

### Test Cases

| Test | Description |
| --- | ----------- |
| `TestCodeCollector_Collect` | Tests collecting from directory |
| `TestCodeCollector_CollectGoFile` | Tests Go file annotation extraction |
| `TestCodeCollector_MultipleAnnotations` | Tests multiple annotations on same line |
| `TestExtractFunctionComment` | Tests function comment extraction |
| `TestExtractFunctionComment/comment_after_annotation` | Comment lines after annotation |
| `TestExtractFunctionComment/function_comment_with_func` | Function comment with func keyword |
| `TestExtractFunctionComment/no_comment` | No comment found |
| `TestExtractFunctionComment/comment_with_func_name` | Comment containing func name |
| `TestExtractFunctionComment_OutOfBounds` | Edge case for line bounds |
| `TestExtractFunctionComment_PointerReceiver` | Tests pointer receiver methods |

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`

---

## TEST-INTERNAL_COLLECTOR-003: ParseFrontmatter Tests

**Status:** Done

**Purpose:**

Test frontmatter YAML parsing with various valid and invalid inputs.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-003`

---

## TEST-INTERNAL_COLLECTOR-004: Frontmatter Edge Cases

**Status:** Done

**Purpose:**

Test frontmatter parsing edge cases like code blocks and empty content.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-003`

---

## TEST-INTERNAL_COLLECTOR-005: Marker Validation

**Status:** Done

**Purpose:**

Test validation that frontmatter markers match actual headings.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-004`

---

## TEST-INTERNAL_COLLECTOR-006: Validation Errors

**Status:** Done

**Purpose:**

Test error reporting for malformed markers.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-004`

---

## TEST-INTERNAL_COLLECTOR-007: File Path Structure

**Status:** Done

**Purpose:**

Test document structure validation by file path.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-005`

---

## TEST-INTERNAL_COLLECTOR-008: Expected Filename

**Status:** Done

**Purpose:**

Test expected filename generation by identifier type.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-005`

---

## TEST-INTERNAL_COLLECTOR-009: Title From Heading

**Status:** Done

**Purpose:**

Test extracting title from markdown heading.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-006`

---

## TEST-INTERNAL_COLLECTOR-010: Title Not Found

**Status:** Done

**Purpose:**

Test handling when no heading found.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-006`

---

## TEST-INTERNAL_COLLECTOR-011: Module Prefix Validation

**Status:** Done

**Purpose:**

Test module prefix matching against directory structure.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-007`

---

## TEST-INTERNAL_COLLECTOR-012: Known Abbreviations

**Status:** Done

**Purpose:**

Test recognition of known module abbreviations.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-007`

---

## TEST-INTERNAL_COLLECTOR-019: Document Type Matching

**Status:** Done

**Purpose:**

Test validation that document type matches filename.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-008`

---

## TEST-INTERNAL_COLLECTOR-014: Root Document Detection

**Status:** Done

**Purpose:**

Test detection of root-level vs module-level documents.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-008`

---

## TEST-INTERNAL_COLLECTOR-015: Extract Annotations

**Status:** Done

**Purpose:**

Test extraction of @implement, @test, @test-contract annotations.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-009`

---

## TEST-INTERNAL_COLLECTOR-016: Function Context

**Status:** Done

**Purpose:**

Test extraction of function name and preceding comments.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-010`

---

## TEST-INTERNAL_COLLECTOR-017: Origin Code Tracking

**Status:** Done

**Purpose:**

Test setting origin to OriginCode for code annotations.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-011`

---

## TEST-INTERNAL_COLLECTOR-018: Multiple Annotations

**Status:** Done

**Purpose:**

Test handling multiple annotations on same line.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-012`

---

## TEST-INTERNAL_COLLECTOR-013: Path Ignore Tests

**Status:** Done

**Purpose:**

Test respecting ignore_paths configuration.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-013`

---

## TEST-INTERNAL_COLLECTOR-020: IDD Reference Extraction

**Status:** Done

**Purpose:**

Test IDD reference regex from markdown content.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-006`

---

## TEST-INTERNAL_COLLECTOR-021: Collector Integration

**Status:** Done

**Purpose:**

Test end-to-end collection from both docs and code.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-014`, `SPEC-INTERNAL_COLLECTOR-015`

---

## TEST-INTERNAL_COLLECTOR-022: Code Collector Basic Collection

**Status:** Done

**Purpose:**

Test basic code annotation collection.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`

---

## TEST-INTERNAL_COLLECTOR-023: Code Collector Go File

**Status:** Done

**Purpose:**

Test Go file annotation extraction.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`

---

## TEST-INTERNAL_COLLECTOR-024: Code Collector Multiple Annotations

**Status:** Done

**Purpose:**

Test handling multiple annotations.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-009`, `SPEC-INTERNAL_COLLECTOR-012`

---

## TEST-INTERNAL_COLLECTOR-025: Code Collector Function Context

**Status:** Done

**Purpose:**

Test function context extraction.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-010`

---

## TEST-INTERNAL_COLLECTOR-026: Code Collector Multi-File

**Status:** Done

**Purpose:**

Test multi-file annotation collection.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`

---

## TEST-INTERNAL_COLLECTOR-027: Code Collector Integration

**Status:** Done

**Purpose:**

Integration test for code collector.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`

---

## TEST-INTERNAL_COLLECTOR-028: Doc Collector Section Title Extraction

**Status:** Done

**Purpose:**

Test that DocCollector correctly extracts section titles from documentation.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-003`

---

## TEST-INTERNAL_COLLECTOR-029: Doc Collector Multi-Section Handling

**Status:** Done

**Purpose:**

Test that DocCollector handles documents with multiple sections.

**Spec Coverage:** `SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-004`

---

## Contract Test

There are no contract tests for this module. The collector module focuses on identifier collection without contract definitions.
