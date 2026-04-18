package auth

// @spec SPEC-AUTH-001
// @contract CONTRACT-AUTH-001
type LoginRequest struct {
	Email    string
	Password string
}

// @spec SPEC-AUTH-001
// @test TEST-AUTH-001
type LoginResponse struct {
	Token string
}
