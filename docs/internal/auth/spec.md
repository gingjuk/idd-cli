---
idd:
  version: "1.1"
  package: internal/auth
  namespace: INTERNAL_AUTH
---

# Specifications: internal/auth

## SPEC-INTERNAL_AUTH-001: Authentication value boundaries

- **Components:** `AuthTypes`
- **Contracts:** `AuthenticationValues`

**Requirement:** Define login credential and token response value types without
implementing authentication behavior.

### Context and required behavior

`LoginRequest` carries an email address and password. `LoginResponse` carries a
token returned by a future authentication service. The package currently
defines these values only; it does not authenticate credentials, issue tokens,
or expose login/logout functions.

Both values must remain directly constructible Go structs so future callers can
exchange them without importing a transport, persistence adapter, or
cryptographic implementation. Their names and fields provide traceable domain
vocabulary; they do not confer runtime validity.

### Constraints and sensitive data

- Email and password validity are caller/service concerns until a validating
  component is explicitly introduced.
- Plaintext passwords must be treated as request-scoped sensitive values even
  though the string type cannot enforce that lifecycle.
- A successful future response is expected to contain a non-empty token, but
  token syntax and claim validation remain outside this package.
- IDD annotations on the declarations are documentation evidence only and must
  not be interpreted as executable validation.

### Implementation boundary and non-goals

This SPEC constrains the existence, names, field meanings, and value-only
nature of `LoginRequest` and `LoginResponse`. It does not require login
services, validation methods, token parsing, session persistence, middleware,
logout behavior, authorization, or external dependencies.

If any of those capabilities are added, they require their own Component,
Contract, behavioral SPEC, TEST evidence, and implementation annotation rather
than being folded into this record as an undocumented assumption.

### Acceptance evidence

**Acceptance:**

The source must expose both annotated value types without runtime dependencies.
Contract tests document representative email, password, and token expectations
while explicitly demonstrating that the tests—not methods on the structs—apply
those checks.

The source declarations implementing this behavior are
`internal/auth/auth.go:LoginRequest` and
`internal/auth/auth.go:LoginResponse`.
