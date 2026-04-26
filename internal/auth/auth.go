// Package auth provides authentication utilities.

// Spec: docs/internal/auth/spec.md
// Contract: docs/internal/auth/contract.md
package auth

// LoginRequest represents user login credentials (email, password).
//
// @implement SPEC-INTERNAL_AUTH-001
type LoginRequest struct {
	Email    string
	Password string
}

// LoginResponse represents authentication result with token.
//
// @test TEST-INTERNAL_AUTH-001
// @implement SPEC-INTERNAL_AUTH-001
type LoginResponse struct {
	Token string
}
