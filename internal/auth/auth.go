package auth

// @spec SPEC-001
// @contract CONTRACT-001
type LoginRequest struct {
	Email    string
	Password string
}

// @spec SPEC-001
// @test TEST-001
type LoginResponse struct {
	Token string
}
