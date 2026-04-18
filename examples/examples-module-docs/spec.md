---
markers:
  - id: SPEC-EXAMPLE-001
    name: User Authentication
---

# Specification Examples

This directory contains example IDD documents showing the expected format.

## SPEC-EXAMPLE-001: User Authentication

**Status:** Done

**Requirement:**

The system must support user authentication with email/password credentials.

**Implementation:** `internal/auth/auth.go`

**Acceptance Criteria:**
- Users can log in with email/password
- Sessions expire after 24 hours
- Passwords are hashed using bcrypt

**Tests:** [TEST-EXAMPLE-001](../testing.md#test-example-001-user-login-flow), [TEST-EXAMPLE-002](../testing.md#test-example-002-session-expiration)

**Related:** [CONTRACT-EXAMPLE-001](../contract.md#contract-example-001-password-hashing-interface), [DESIGN-EXAMPLE-001](../design.md#design-example-001-graph-first-architecture)
