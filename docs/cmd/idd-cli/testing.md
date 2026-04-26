---
markers:
  - id: TEST-CMD_IDD_CLI-001
    name: Core Validation Tests
  - id: TEST-CMD_IDD_CLI-002
    name: Module Integration Tests

related_files:
  spec: docs/cmd/idd-cli/spec.md
  contract: docs/cmd/idd-cli/contract.md
  design: docs/cmd/idd-cli/design.md
  testing: docs/cmd/idd-cli/testing.md
---

# Test Cases (backend)

## TEST-CMD_IDD_CLI-001: Core Validation Tests

**Status:** Done

**Purpose:**

Test cases for the idd-cli core validation logic.

### Model Tests (`internal/model/`)

| Test | Description |
| --- | --- |
| `TestParseIdentifierType` | Validates identifier type parsing |
| `TestNewIdentifier` | Tests identifier creation |
| `TestIdentifierSet_Add_Get_Has` | Tests collection operations |
| `TestIdentifierSet_Count` | Tests counting identifiers |
| `TestIdentifierSet_All` | Tests retrieval of all identifiers |
| `TestIdentifierSet_Merge` | Tests merging identifier sets |
| `TestIdentifier_AddLink` | Tests adding forward references |
| `TestNewAnnotation` | Tests annotation creation |
| `TestAnnotation_ToIdentifier` | Tests annotation conversion |
| `TestValidationError_Error` | Tests error formatting |
| `TestNewValidationResult` | Tests result initialization |
| `TestValidationResult_AddError` | Tests error accumulation |
| `TestValidationResult_AddWarning` | Tests warning accumulation |
| `TestValidationResult_Sort` | Tests error sorting |

### Graph Tests (`internal/graph/`)

| Test | Description |
| --- | --- |
| `TestNewLinkageGraph` | Tests graph initialization |
| `TestLinkageGraph_AddNode` | Tests node addition |
| `TestLinkageGraph_AddEdge` | Tests edge addition |
| `TestLinkageGraph_GetNode` | Tests node lookup |
| `TestLinkageGraph_NodeInOutEdges` | Tests edge traversal |
| `TestLinkageGraph_GetOutboundByType` | Tests outbound filtering |
| `TestLinkageGraph_GetInboundByType` | Tests inbound filtering |
| `TestLinkageGraph_GetBacklinks` | Tests backlink lookup |
| `TestLinkageGraph_VerifyBidirectionalLinks` | Tests verification |
| `TestLinkageGraph_ToSnapshot` | Tests serialization |
| `TestLinkageGraph_Stats` | Tests statistics |
| `TestLinkageGraph_ValidateCompleteness` | Tests completeness |

**Spec Coverage:** `SPEC-CMD_IDD_CLI-001`, `SPEC-CMD_IDD_CLI-002`, `SPEC-CMD_IDD_CLI-003`, `SPEC-CMD_IDD_CLI-004`, `SPEC-CMD_IDD_CLI-005`, `SPEC-CMD_IDD_CLI-006`, `SPEC-CMD_IDD_CLI-007`, `SPEC-CMD_IDD_CLI-008`, `SPEC-CMD_IDD_CLI-009`, `SPEC-CMD_IDD_CLI-010`

---

## TEST-CMD_IDD_CLI-002: Module Integration Tests

**Status:** Done

**Purpose:**

Integration tests covering the interaction between modules including config loading, engine orchestration, and reporter output.

### Config Tests (`internal/config/`)

| Test | Description |
| --- | --- |
| `TestLoadConfig` | Tests loading config from YAML file |
| `TestDefaultConfig` | Tests default config values |
| `TestConfigFileSearch` | Tests config file discovery |

### Engine Tests (`internal/engine/`)

| Test | Description |
| --- | --- |
| `TestEngine_Run` | Tests end-to-end validation run |
| `TestEngine_Collect` | Tests collection orchestration |
| `TestEngine_BuildGraph` | Tests graph building |
| `TestEngine_Validate` | Tests validation rule execution |

### Reporter Tests (`internal/reporter/`)

| Test | Description |
| --- | --- |
| `TestJSONReporter` | Tests JSON output format |
| `TestMarkdownReporter` | Tests Markdown output format |

### Similarity Tests (`internal/similarity/`)

| Test | Description |
| --- | --- |
| `TestTFIDF_Vectorize` | Tests TF-IDF vectorization |
| `TestTFIDF_CosineSimilarity` | Tests cosine similarity computation |
| `TestTFIDF_FindDuplicates` | Tests duplicate detection |

**Spec Coverage:** `SPEC-CMD_IDD_CLI-001`, `SPEC-CMD_IDD_CLI-002`, `SPEC-CMD_IDD_CLI-003`, `SPEC-CMD_IDD_CLI-004`, `SPEC-CMD_IDD_CLI-005`, `SPEC-CMD_IDD_CLI-006`, `SPEC-CMD_IDD_CLI-007`, `SPEC-CMD_IDD_CLI-008`, `SPEC-CMD_IDD_CLI-009`, `SPEC-CMD_IDD_CLI-010`
