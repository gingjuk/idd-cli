---
related_files:
  spec: docs/internal/auth/spec.md
  contract: docs/internal/auth/contract.md
  design: docs/internal/auth/design.md
  testing: docs/internal/auth/testing.md
---

# Contracts (auth)

## Authentication Interface Contracts

**Status:** Done

**Overview:**

Contracts for authentication functionality in idd-cli.

### LoginRequest Contract

```go
type LoginRequest struct {
    Email    string  // User email address
    Password string  // User password (not stored in plaintext)
}
```

**Invariant:**

- `Email` must be a valid email format
- `Password` must be non-empty

### LoginResponse Contract

```go
type LoginResponse struct {
    Token string  // JWT or session token for authenticated user
}
```

**Invariant:**

- `Token` must be non-empty after successful authentication
- Token format must be parseable

### Error Handling

- Invalid credentials → return error with `ErrInvalidCredentials`
- Account locked → return error with `ErrAccountLocked`
- Network failure → return error with `ErrNetworkFailure`

**Related Specs:** `SPEC-INTERNAL_AUTH-001`
