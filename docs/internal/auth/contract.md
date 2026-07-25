---
idd:
  version: "1.0"
  package: internal/auth
  document: contract
---

# Contracts: internal/auth

## Contract: AuthenticationValues

### LoginRequest

```go
type LoginRequest struct {
    Email    string  // User email address
    Password string  // User password (not stored in plaintext)
}
```

- `Email` is expected to contain a syntactically valid email address.
- `Password` is expected to be non-empty.
- The value type stores inputs but does not enforce either constraint.

### LoginResponse

```go
type LoginResponse struct {
    Token string  // JWT or session token for authenticated user
}
```

- `Token` is expected to be non-empty after successful authentication.
- Token consumers currently expect a three-segment, dot-separated value.
- The value type stores the token but does not generate or validate it.
