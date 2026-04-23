---
markers:
  - id: SPEC-INT_AUTH-001
    name: Authentication Module

related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Specification (auth)

## SPEC-INT_AUTH-001: Authentication Module

**Status:** Done

**Contract:** `LoginRequest`

**Design:** `AuthModule`

**Requirement:**

Authentication module provides user login and token generation functionality for idd-cli. It defines the core types for authentication requests and responses.

**Implementation:** `internal/auth/auth.go`

**Key Types:**

- `LoginRequest` — User login credentials (email, password)
- `LoginResponse` — Authentication result with token

**Acceptance Criteria:**

- [x] LoginRequest struct contains email and password fields
- [x] LoginResponse struct contains token field
- [x] Authentication types are annotated with @implement; tests use @test-contract for CONTRACT coverage
- [x] Test annotations link to `TEST-INT_AUTH-001`

**Tests:** `TEST-INT_AUTH-001`
