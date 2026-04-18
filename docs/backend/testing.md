---
markers:
  - id: TEST-BE-001
    name: Core Validation Tests
---

# Test Cases (backend)

## TEST-BE-001: Core Validation Tests

**Status:** Done

**Purpose:**

Test cases for the idd-cli core validation logic.

### Model Tests (`internal/model/`)

| Test | Description |
|------|-------------|
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
|------|-------------|
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

**Spec Coverage:** SPEC-BE-001, SPEC-BE-002
