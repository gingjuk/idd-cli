# TEST-BE-001 User Login Flow

## Test Case

Tests the complete login flow including validation.

## Specs Covered

- SPEC-BE-001 (User Authentication)

## Code Coverage

```go
// @test SPEC-BE-001
func TestLogin(t *testing.T) {
    // test implementation
}
```

## Expected Results

1. Valid credentials return success
2. Invalid credentials return error
3. Session is created on success