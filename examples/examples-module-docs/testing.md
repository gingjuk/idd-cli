---
markers:
  - id: TEST-BE-001
    name: User Login Flow
  - id: TEST-BE-002
    name: Session Expiration
---

# Test Examples

This directory contains example IDD test documents showing the expected format.

## TEST-BE-001: User Login Flow

**Status:** Done

**Purpose:**

Test cases for the user authentication login flow.

**Spec Coverage:** [SPEC-BE-001](../spec.md#spec-be-001-user-authentication)

### Test Cases

| Test | Description |
|------|-------------|
| `TestLogin_Success` | Valid credentials return success with session token |
| `TestLogin_InvalidPassword` | Invalid password returns error |
| `TestLogin_UserNotFound` | Non-existent user returns error |

**Code Example:**
```go
// TEST-BE-001: User Login Flow
// Covers: SPEC-BE-001
func TestLogin_Success(t *testing.T) {
    req := LoginRequest{Email: "user@test.com", Password: "valid"}
    resp := Login(req)
    if resp.Token == "" {
        t.Error("expected non-empty token")
    }
}
```

---

## TEST-BE-002: Session Expiration

**Status:** Done

**Purpose:**

Test cases for session expiration behavior.

**Spec Coverage:** [SPEC-BE-001](../spec.md#spec-be-001-user-authentication)

### Test Cases

| Test | Description |
|------|-------------|
| `TestSession_ExpiresAfter24Hours` | Session becomes invalid after 24 hours |
| `TestSession_Refresh` | Session can be refreshed before expiration |

**Code Example:**
```go
// TEST-BE-002: Session Expiration
// Covers: SPEC-BE-001
func TestSession_ExpiresAfter24Hours(t *testing.T) {
    session := CreateSession(time.Now().Add(-25 * time.Hour))
    if session.IsValid() {
        t.Error("expected session to be expired")
    }
}
```
