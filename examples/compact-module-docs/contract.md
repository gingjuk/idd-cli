---
idd:
  version: "1.0"
  package: internal/auth
  document: contract
---

# Contracts: internal/auth

## Contract: Authenticator

`Authenticate` returns an identity only when all credentials are valid. Invalid
usernames and invalid secrets return the same public error so callers cannot
use the boundary for account discovery.

```go
type Authenticator interface {
    Authenticate(ctx context.Context, credentials Credentials) (Identity, error)
}
```
