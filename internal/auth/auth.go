// Package auth provides authentication utilities.
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
// @implement SPEC-INTERNAL_AUTH-001
type LoginResponse struct {
	Token string
}
