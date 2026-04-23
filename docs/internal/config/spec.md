---
markers:
  - id: SPEC-INT_CFG-001
    name: Config Structure
  - id: SPEC-INT_CFG-002
    name: Docs Config Structure
  - id: SPEC-INT_CFG-003
    name: Identifier Patterns Structure
  - id: SPEC-INT_CFG-004
    name: Code Config Structure
  - id: SPEC-INT_CFG-005
    name: Validation Config Structure
  - id: SPEC-INT_CFG-006
    name: Output Config Structure
  - id: SPEC-INT_CFG-007
    name: Load Config Function
  - id: SPEC-INT_CFG-008
    name: Default Config Function
  - id: SPEC-INT_CFG-009
    name: Validate Config Function

related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Specification (config)

Configuration module for idd-cli.

## SPEC-INT_CFG-001: Config Structure

**Status:** Done

**Contract:** `Config`

**Design:** `ConfigModule`

**Requirement:**

Config is the root configuration structure that holds all settings for the IDD CLI validation tool.

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INT_CFG-001`

**Tests:** `TEST-INT_CFG-001`

---

## SPEC-INT_CFG-002: Docs Config Structure

**Status:** Done

**Contract:** `DocsConfig`

**Design:** `ConfigModule`

**Requirement:**

DocsConfig holds documentation-related configuration including patterns and ignore paths.

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INT_CFG-002`

**Tests:** `TEST-INT_CFG-002`

---

## SPEC-INT_CFG-003: Identifier Patterns Structure

**Status:** Done

**Contract:** `IdentifierPatterns`

**Design:** `ConfigModule`

**Requirement:**

IdentifierPatterns defines regex patterns for matching SPEC, TEST, and other IDD identifiers.

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INT_CFG-003`

**Tests:** `TEST-INT_CFG-003`

---

## SPEC-INT_CFG-004: Code Config Structure

**Status:** Done

**Contract:** `CodeConfig`

**Design:** `ConfigModule`

**Requirement:**

CodeConfig holds code-related configuration including patterns and annotations.

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INT_CFG-004`

**Tests:** `TEST-INT_CFG-003`

---

## SPEC-INT_CFG-005: Validation Config Structure

**Status:** Done

**Contract:** `ValidationConfig`

**Design:** `ConfigModule`

**Requirement:**

ValidationConfig holds validation rule settings for the IDD CLI.

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INT_CFG-005`

**Tests:** `TEST-INT_CFG-003`

---

## SPEC-INT_CFG-006: Output Config Structure

**Status:** Done

**Contract:** `OutputConfig`

**Design:** `ConfigModule`

**Requirement:**

OutputConfig holds output-related configuration settings.

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INT_CFG-006`

**Tests:** `TEST-INT_CFG-003`

---

## SPEC-INT_CFG-007: Load Config Function

**Status:** Done

**Contract:** `Load`

**Design:** `ConfigModule`

**Requirement:**

Load configuration from YAML files with validation.

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INT_CFG-007`

**Tests:** `TEST-INT_CFG-002`

---

## SPEC-INT_CFG-008: Default Config Function

**Status:** Done

**Contract:** `Default`

**Design:** `ConfigModule`

**Requirement:**

Return sensible default configuration.

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INT_CFG-008`

**Tests:** `TEST-INT_CFG-001`

---

## SPEC-INT_CFG-009: Validate Config Function

**Status:** Done

**Contract:** `Validate`

**Design:** `ConfigModule`

**Requirement:**

Validate configuration values are correct.

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INT_CFG-009`

**Tests:** `TEST-INT_CFG-003`
