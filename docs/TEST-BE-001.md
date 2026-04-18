# Testing: TEST-BE-001

## Core Validation Tests

**Version:** 1.0.0
**Last Updated:** 2026-04-18

## Overview

Test cases for the IDD Link Validator core validation logic.

## Test Coverage

### Model Tests (`internal/model/identifier_test.go`)

| Test | Description |
|------|-------------|
| `TestParseIdentifierType` | Validates identifier type parsing for SPEC, TEST, CONTRACT, DESIGN |
| `TestNewIdentifier` | Tests identifier creation with all fields |
| `TestIdentifierSet_Add_Get_Has` | Tests collection operations |
| `TestIdentifierSet_Count` | Tests counting multiple identifiers |
| `TestIdentifierSet_All` | Tests retrieval of all identifiers |
| `TestIdentifierSet_Merge` | Tests merging two identifier sets |
| `TestIdentifier_AddLink` | Tests adding forward references |
| `TestNewAnnotation` | Tests annotation creation |
| `TestAnnotation_ToIdentifier` | Tests annotation to identifier conversion |
| `TestValidationError_Error` | Tests error message formatting |
| `TestNewValidationResult` | Tests result initialization |
| `TestValidationResult_AddError` | Tests error accumulation |
| `TestValidationResult_AddWarning` | Tests warning accumulation |
| `TestValidationResult_Sort` | Tests error/warning sorting |

### Model Link Tests (`internal/model/link_test.go`)

| Test | Description |
|------|-------------|
| `TestLinkType_Constants` | Validates link type constants |
| `TestReverseLinkType` | Tests link type reversal logic |
| `TestNewLink` | Tests link creation |

### Graph Tests (`internal/graph/graph_test.go`)

| Test | Description |
|------|-------------|
| `TestNewLinkageGraph` | Tests graph initialization |
| `TestLinkageGraph_AddNode` | Tests node addition and deduplication |
| `TestLinkageGraph_AddEdge` | Tests edge addition |
| `TestLinkageGraph_GetNode` | Tests node lookup |
| `TestLinkageGraph_NodeInOutEdges` | Tests edge traversal |
| `TestLinkageGraph_GetOutboundByType` | Tests outbound edge filtering |
| `TestLinkageGraph_GetInboundByType` | Tests inbound edge filtering |
| `TestLinkageGraph_GetBacklinks` | Tests backlink lookup |
| `TestLinkageGraph_VerifyBidirectionalLinks` | Tests bidirectional verification |
| `TestLinkageGraph_VerifyBidirectionalLinks_Unverified` | Tests unidirectional edge detection |
| `TestLinkageGraph_ToSnapshot` | Tests graph serialization |
| `TestLinkageGraph_Stats` | Tests statistics generation |
| `TestLinkageGraph_ValidateCompleteness` | Tests completeness validation |
| `TestLinkageGraph_Nodes_Edges` | Tests node/edge accessors |

## Related Documents

- **SPEC-BE-001** — IDD Link Validator Overview
- **CONTRACT-BE-001** — Collector Interface Contracts
