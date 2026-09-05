---
idd:
  version: "1.1"
  package: internal/auth
  namespace: INTERNAL_AUTH
---

# Specifications: internal/auth

## SPEC-INTERNAL_AUTH-001: Credential failure details

- **Components:** `AuthModule`
- **Contracts:** `Authenticator`

**Requirement:** Authenticate users with validated credentials without
disclosing which credential failed.

**Acceptance:** A valid credential fixture returns the complete expected
identity. Unknown identifiers and wrong secrets both return
`ErrInvalidCredentials`; malformed input avoids dependency work; cancellation
and operational failures remain distinguishable; and no failure returns a
partial identity or exposes a plaintext secret.

### Context and required behavior

Account existence is sensitive information. The authentication path must
therefore accept a syntactically valid identifier and non-empty secret, compare
them against the configured credential source, and expose the identity only
when both values are valid.

For every rejected credential attempt, the caller receives
`ErrInvalidCredentials` regardless of whether the identifier was unknown or the
secret was wrong. Internal telemetry may retain a private reason, provided it
does not include the plaintext secret or alter the public response.

### Failure and edge cases

- Invalid request shape is rejected before storage lookup.
- Cancellation and deadline expiry remain distinguishable from credential
  rejection.
- Store and verifier failures remain operational errors and must not be
  reported as invalid credentials.
- A failure never returns a partial or non-zero identity.
- Repeated calls do not create hidden session state or mutate credentials.

### Implementation boundary

This SPEC constrains credential verification and public failure semantics. It
does not define transport status codes, retry policy, rate limiting, password
hash algorithms, token issuance, session storage, authorization, or
multi-factor authentication. Those concerns belong to callers, injected
dependencies, or separate behavioral SPECs.

The implementation may use any credential store and password verifier that
satisfy the documented contract. It must not expose their concrete types
through the `Authenticator` boundary.

### Code association

The implementation declaration carries the SPEC identifier directly:

```go
// @implement SPEC-INTERNAL_AUTH-001
func (service *Service) Authenticate(
    ctx context.Context,
    credentials Credentials,
) (Identity, error) {
    // Implementation omitted from the documentation example.
}
```

idd-cli parses the declaration and joins this annotation to the SPEC record.
The source file does not repeat paths to `spec.md`, `design.md`, or
`contract.md`; document names and relationships remain owned by this
self-describing set.

### Rationale

Collapsing expected rejection reasons prevents callers from using response
differences for account discovery. Preserving cancellation and dependency
errors separately keeps operational failures diagnosable and avoids
misclassifying outages as user mistakes.

### Acceptance evidence

The acceptance paragraph is proved by observable cases for valid credentials,
unknown identifiers, incorrect secrets, malformed input, cancellation, and
dependency failure. Tests must prove both the returned identity/error and the
absence of secret leakage or partial success.
