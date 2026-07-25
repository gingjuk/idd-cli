---
idd:
  version: "1.0"
  package: internal/auth
  document: spec
---

# Specifications: internal/auth

## SPEC-INTERNAL_AUTH-001: Authentication value boundaries

- **Design:** `AuthTypes`
- **Contract:** `AuthenticationValues`

**Requirement:** Define login credential and token response value types without
implementing authentication behavior.

`LoginRequest` carries an email address and password. `LoginResponse` carries a
token returned by a future authentication service. The package currently
defines these values only; it does not authenticate credentials, issue tokens,
or expose login/logout functions.

The source declarations implementing this behavior are
`internal/auth/auth.go:LoginRequest` and
`internal/auth/auth.go:LoginResponse`.
