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

**Design:** `ConfigModule`

**Contract:** `Config`

**Requirement:**

Config is the root configuration structure that holds all settings for the IDD CLI validation tool.

**Tests:** `TEST-INTERNAL_CONFIG-001`

**Status:** Done

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INTERNAL_CONFIG-001`

---

## SPEC-INTERNAL_CONFIG-002: Docs Config Structure

**Design:** `ConfigModule`

**Contract:** `DocsConfig`

**Requirement:**

DocsConfig holds documentation-related configuration including patterns and ignore paths.

**Tests:** `TEST-INTERNAL_CONFIG-002`

**Status:** Done

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INTERNAL_CONFIG-002`

---

## SPEC-INTERNAL_CONFIG-003: Identifier Patterns Structure

**Design:** `ConfigModule`

**Contract:** `IdentifierPatterns`

**Requirement:**

IdentifierPatterns defines regex patterns for matching SPEC, TEST, and other IDD identifiers.

**Tests:** `TEST-INTERNAL_CONFIG-003`

**Status:** Done

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INTERNAL_CONFIG-003`

---

## SPEC-INTERNAL_CONFIG-004: Code Config Structure

**Design:** `ConfigModule`

**Contract:** `CodeConfig`

**Requirement:**

CodeConfig holds code-related configuration including patterns and annotations.

**Tests:** `TEST-INTERNAL_CONFIG-003`

**Status:** Done

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INTERNAL_CONFIG-004`

---

## SPEC-INTERNAL_CONFIG-005: Validation Config Structure

**Design:** `ConfigModule`

**Contract:** `ValidationConfig`

**Requirement:**

ValidationConfig holds validation rule settings for the IDD CLI, including the `require_spec_fields` toggle that requires each frontmatter-declared SPEC section to include `Design`, `Contract`, `Requirement`, and `Tests` in that order before any optional fields.

**Tests:** `TEST-INTERNAL_CONFIG-003`

**Status:** Done

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INTERNAL_CONFIG-005`

---

## SPEC-INTERNAL_CONFIG-006: Output Config Structure

**Design:** `ConfigModule`

**Contract:** `OutputConfig`

**Requirement:**

OutputConfig holds output-related configuration settings.

**Tests:** `TEST-INTERNAL_CONFIG-003`

**Status:** Done

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INTERNAL_CONFIG-006`

---

## SPEC-INTERNAL_CONFIG-007: Load Config Function

**Design:** `ConfigModule`

**Contract:** `Load`

**Requirement:**

Load configuration from YAML files with validation.

**Tests:** `TEST-INTERNAL_CONFIG-002`

**Status:** Done

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INTERNAL_CONFIG-007`

---

## SPEC-INTERNAL_CONFIG-008: Default Config Function

**Design:** `ConfigModule`

**Contract:** `Default`

**Requirement:**

Return sensible default configuration.

**Tests:** `TEST-INTERNAL_CONFIG-001`

**Status:** Done

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INTERNAL_CONFIG-008`

---

## SPEC-INTERNAL_CONFIG-009: Validate Config Function

**Design:** `ConfigModule`

**Contract:** `Validate`

**Requirement:**

Validate configuration values are correct.

**Tests:** `TEST-INTERNAL_CONFIG-003`
**Status:** Done

**Implementation:** `internal/config/config.go`

**Code Annotation:** `@implement SPEC-INTERNAL_CONFIG-009`
