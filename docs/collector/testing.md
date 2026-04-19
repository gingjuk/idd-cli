---
markers:
  - id: TEST-COLLECTOR-001
    name: Document Collector Tests
  - id: TEST-COLLECTOR-002
    name: Code Collector Tests
---

# Test Cases (collector)

## TEST-COLLECTOR-001: Document Collector Tests

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

### Frontmatter Tests

| Test | Description |
| --- | ----------- |
| `TestParseFrontmatter` | Tests frontmatter YAML parsing |
| `TestParseFrontmatter/valid_frontmatter` | Valid frontmatter with markers |
| `TestParseFrontmatter/no_frontmatter` | No frontmatter returns nil |
| `TestParseFrontmatter/empty_content` | Empty content handled gracefully |
| `TestParseFrontmatter/code_block_with_dashes` | Code blocks don't interfere |
| `TestParseFrontmatter/frontmatter_with_describe` | Marker with describe field |

### Validation Tests

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

**Spec Coverage:** `SPEC-COLLECTOR-001`, `SPEC-BE-001`

---

## TEST-COLLECTOR-002: Code Collector Tests

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

**Spec Coverage:** `SPEC-COLLECTOR-002`, `SPEC-BE-001`
