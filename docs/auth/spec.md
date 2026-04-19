---
markers:
  - id: SPEC-AUTH-001
    name: Authentication Module
  - id: SPEC-AUTH-002
    name: Authorization Module
---

# Specification (auth)

## SPEC-AUTH-001: Authentication Module

**Status:** Done

**Requirement:**

Authentication module provides user login and token generation functionality for idd-cli. It defines the core types for authentication requests and responses.

**Implementation:** `internal/auth/auth.go`

**Key Types:**

- `LoginRequest` — User login credentials (email, password)
- `LoginResponse` — Authentication result with token

**Acceptance Criteria:**

- [x] LoginRequest struct contains email and password fields
- [x] LoginResponse struct contains token field
- [x] Authentication types are annotated with @spec and @contract
- [x] Test annotations link to `TEST-AUTH-001`

**Tests:** `TEST-AUTH-001`

**Related:** `CONTRACT-AUTH-001`

---

## SPEC-AUTH-002: Authorization Module

**Status:** Planned

**Requirement:**

Authorization module will handle permission checking and access control for idd-cli operations.

**Implementation:** Not yet implemented

**Key Features:**

- Role-based access control (RBAC)
- Permission verification before operations
- Access control lists (ACL)

**Acceptance Criteria:**

- [ ] Define permission types and roles
- [ ] Implement permission checking logic
- [ ] Add authorization middleware
- [ ] Integrate with CLI commands

**Related:** `SPEC-AUTH-001`
