---
markers:
  - id: CONTRACT-EXAMPLE-001
    name: Password Hashing Interface
  - id: CONTRACT-EXAMPLE-002
    name: Session Management Interface
---

# Contract Examples

This directory contains example IDD contract documents showing the expected format.

## CONTRACT-EXAMPLE-001: Password Hashing Interface

**Status:** Done

**Overview:**

Contract for secure password hashing and verification.

### Interface

```go
type PasswordHasher interface {
    Hash(password string) (string, error)
    Verify(password, hash string) error
}
```

### Error Handling

- **Empty password** — Return error
- **Hash failure** — Return error
- **Verification mismatch** — Return error

**Related Specs:** [SPEC-EXAMPLE-001](../spec.md#spec-example-001-user-authentication)

---

## CONTRACT-EXAMPLE-002: Session Management Interface

**Status:** Done

### Interface

```go
type SessionManager interface {
    Create(userID string) (*Session, error)
    Validate(token string) (*Session, error)
    Refresh(token string) (*Session, error)
    Revoke(token string) error
}
```

### Session Lifecycle

1. **Create** — Generate new session with 24-hour expiry
2. **Validate** — Check session is still valid
3. **Refresh** — Extend session lifetime
4. **Revoke** — Immediately invalidate session

**Related Specs:** [SPEC-EXAMPLE-001](../spec.md#spec-example-001-user-authentication)
