---
idd:
  version: "1.0"
  package: internal/auth
  document: spec
---

# Specifications: internal/auth

## SPEC-INTERNAL_AUTH-001: Credential failure details

- **Design:** `AuthModule`
- **Contract:** `Authenticator`

**Requirement:** Authenticate users with validated credentials without
disclosing which credential failed.

Credential validation may record an internal failure reason for operations, but
the returned contract exposes one stable invalid-credentials error. Cancellation
and dependency failures remain distinguishable from rejected credentials.
