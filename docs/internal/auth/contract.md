---
idd:
  version: "1.0"
  package: internal/auth
  document: contract
---

# Contracts: internal/auth

## Contract: AuthenticationValues

**Guarantees:**

`AuthenticationValues` is a data-shape contract for callers that will exchange
login credentials and token responses. It describes the meaning and handling
expectations of the fields; the current structs do not enforce those
expectations at runtime.

### LoginRequest

```go
type LoginRequest struct {
    Email    string  // User email address
    Password string  // User password (not stored in plaintext)
}
```

#### Inputs and validity

`Email` is expected to contain a syntactically valid address and `Password` is
expected to be non-empty. The contract does not prescribe normalization,
internationalized-address handling, password strength, or credential
verification. A future service must define those policies before treating the
request as authenticated input.

#### Ownership and security

The caller owns both strings and their lifecycle. Consumers may read the
password only for the operation that needs it and must not log or persist it in
plaintext. The struct offers no automatic redaction, zeroization, or defensive
copying.

#### Failure behavior

Constructing `LoginRequest` cannot fail, and an invalid value remains
representable. Any component that accepts the value must return its own defined
validation error rather than assuming the struct proves validity.

### LoginResponse

```go
type LoginResponse struct {
    Token string  // JWT or session token for authenticated user
}
```

#### Output and validity

`Token` is expected to be non-empty after successful authentication. Current
contract evidence treats a three-segment, dot-separated token as valid, but the
type itself neither fixes a JWT standard nor validates claims, signature,
expiry, audience, or issuer.

#### Ownership and lifecycle

The recipient owns storage and disclosure decisions for the token. The value
does not provide expiry, refresh, revocation, or secure-storage behavior.

#### Invariants and compatibility

A successful future authentication result should not expose an empty token.
Changing the field type or defining a permanent token format is a compatibility
decision that requires a separate behavioral contract. Until then,
`LoginResponse` is a transport-neutral container rather than a token service.

### Package-level boundary

This contract intentionally excludes authentication algorithms, credential
stores, token generation, session management, authorization, transports, and
error mapping. Those capabilities cannot be inferred from the presence of the
request and response values.
