---
idd:
  version: "1.0"
  package: internal/auth
  document: testing
---

# Testing: internal/auth

## TEST-INTERNAL_AUTH-001: Authentication behavior

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_AUTH-001`

**Purpose:** Verify accepted and rejected credential outcomes.

## TEST-INTERNAL_AUTH-002: Authenticator contract

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_AUTH-001`

**Purpose:** Verify every `Authenticator` implementation preserves the public
boundary.

## Strategy

Behavior tests cover valid credentials, unknown users, invalid secrets,
cancelled contexts, and credential-store failures. The contract suite reuses
the public cases for every `Authenticator` implementation and verifies that
unknown-user and invalid-secret failures are indistinguishable.
