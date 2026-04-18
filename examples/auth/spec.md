---
markers:
  - id: SPEC-BE-001
    name: User Authentication
---

# Specification Examples

This directory contains example IDD documents showing the expected format.

## SPEC-BE-001: User Authentication

**Status:** Done

**Requirement:**

The system must support user authentication with email/password credentials.

**Implementation:** `internal/auth/auth.go`

**Acceptance Criteria:**
- Users can log in with email/password
- Sessions expire after 24 hours
- Passwords are hashed using bcrypt

**Tests:** [TEST-BE-001](../testing.md#test-be-001-user-login-flow), [TEST-BE-002](../testing.md#test-be-002-session-expiration)

**Related:** [CONTRACT-BE-001](../contract.md#contract-be-001-password-hashing-interface), [DESIGN-BE-001](../design.md#design-be-001-graph-first-architecture)
