---
markers:
  - id: SPEC-INTERNAL_CONFIG-001
    name: Config Structure
  - id: SPEC-INTERNAL_CONFIG-002
    name: Docs Config Structure
  - id: SPEC-INTERNAL_CONFIG-003
    name: Identifier Patterns Structure
  - id: SPEC-INTERNAL_CONFIG-004
    name: Code Config Structure
  - id: SPEC-INTERNAL_CONFIG-005
    name: Validation Config Structure
  - id: SPEC-INTERNAL_CONFIG-006
    name: Output Config Structure
  - id: SPEC-INTERNAL_CONFIG-007
    name: Load Config Function
  - id: SPEC-INTERNAL_CONFIG-008
    name: Default Config Function
  - id: SPEC-INTERNAL_CONFIG-009
    name: Validate Config Function

related_files:
  spec: docs/internal/config/spec.md
  testing: docs/internal/config/testing.md
---

# Specification (config)

Configuration module for idd-cli.

## SPEC-INTERNAL_CONFIG-001: Config Structure

**Status:** Done

**Contract:** `Config`

**Design:** `ConfigModule`

**Requirement:**

Config is the root configuration structure that holds all settings for the IDD CLI validation tool.

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INTERNAL_CONFIG-001`

**Tests:** `TEST-INTERNAL_CONFIG-001`

---

## SPEC-INTERNAL_CONFIG-002: Docs Config Structure

**Status:** Done

**Contract:** `DocsConfig`

**Design:** `ConfigModule`

**Requirement:**

DocsConfig holds documentation-related configuration including patterns and ignore paths.

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INTERNAL_CONFIG-002`

**Tests:** `TEST-INTERNAL_CONFIG-002`

---

## SPEC-INTERNAL_CONFIG-003: Identifier Patterns Structure

**Status:** Done

**Contract:** `IdentifierPatterns`

**Design:** `ConfigModule`

**Requirement:**

IdentifierPatterns defines regex patterns for matching SPEC, TEST, and other IDD identifiers.

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INTERNAL_CONFIG-003`

**Tests:** `TEST-INTERNAL_CONFIG-003`

---

## SPEC-INTERNAL_CONFIG-004: Code Config Structure

**Status:** Done

**Contract:** `CodeConfig`

**Design:** `ConfigModule`

**Requirement:**

CodeConfig holds code-related configuration including patterns and annotations.

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INTERNAL_CONFIG-004`

**Tests:** `TEST-INTERNAL_CONFIG-003`

---

## SPEC-INTERNAL_CONFIG-005: Validation Config Structure

**Status:** Done

**Contract:** `ValidationConfig`

**Design:** `ConfigModule`

**Requirement:**

ValidationConfig holds validation rule settings for the IDD CLI.

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INTERNAL_CONFIG-005`

**Tests:** `TEST-INTERNAL_CONFIG-003`

---

## SPEC-INTERNAL_CONFIG-006: Output Config Structure

**Status:** Done

**Contract:** `OutputConfig`

**Design:** `ConfigModule`

**Requirement:**

OutputConfig holds output-related configuration settings.

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INTERNAL_CONFIG-006`

**Tests:** `TEST-INTERNAL_CONFIG-003`

---

## SPEC-INTERNAL_CONFIG-007: Load Config Function

**Status:** Done

**Contract:** `Load`

**Design:** `ConfigModule`

**Requirement:**

Load configuration from YAML files with validation.

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INTERNAL_CONFIG-007`

**Tests:** `TEST-INTERNAL_CONFIG-002`

---

## SPEC-INTERNAL_CONFIG-008: Default Config Function

**Status:** Done

**Contract:** `Default`

**Design:** `ConfigModule`

**Requirement:**

Return sensible default configuration.

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INTERNAL_CONFIG-008`

**Tests:** `TEST-INTERNAL_CONFIG-001`

---

## SPEC-INTERNAL_CONFIG-009: Validate Config Function

**Status:** Done

**Contract:** `Validate`

**Design:** `ConfigModule`

**Requirement:**

Validate configuration values are correct.

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INTERNAL_CONFIG-009`

**Tests:** `TEST-INTERNAL_CONFIG-003`
