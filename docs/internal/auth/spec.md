---
markers:
  - id: SPEC-INTERNAL_AUTH-001
    name: Authentication Module

related_files:
  spec: docs/internal/auth/spec.md
  contract: docs/internal/auth/contract.md
  design: docs/internal/auth/design.md
  testing: docs/internal/auth/testing.md

---

# Specification (auth)

## SPEC-INTERNAL_AUTH-001: Authentication Module

**Design:** `AuthModule`

**Contract:** `LoginRequest`

**Requirement:**

Authentication module provides user login and token generation functionality for idd-cli. It defines the core types for authentication requests and responses.

**Tests:** `TEST-INTERNAL_AUTH-001`
**Status:** Done

**Implementation:** `internal/auth/auth.go`

**Key Types:**

- `LoginRequest` — User login credentials (email, password)
- `LoginResponse` — Authentication result with token

**Acceptance Criteria:**

- [x] LoginRequest struct contains email and password fields
- [x] LoginResponse struct contains token field
- [x] Authentication types are annotated with @implement; tests use @test-contract for CONTRACT coverage
- [x] Test annotations link to `TEST-INTERNAL_AUTH-001`
