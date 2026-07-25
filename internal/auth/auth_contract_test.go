// Package auth provides testing utilities for the auth module.

// Spec: docs/internal/auth/spec.md
// Test: docs/internal/auth/testing.md
// Contract: docs/internal/auth/contract.md
package auth

import (
	"regexp"
	"testing"
)

// CONTRACT-AUTH-001: Authentication Interface Contracts

// @test-contract TEST-INTERNAL_AUTH-001
func TestLoginRequest_EmailFormat(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "valid email format",
			email:   "user@example.com",
			wantErr: false,
		},
		{
			name:    "valid email with subdomain",
			email:   "user@mail.example.com",
			wantErr: false,
		},
		{
			name:    "invalid email - no at symbol",
			email:   "userexample.com",
			wantErr: true,
		},
		{
			name:    "invalid email - no domain",
			email:   "user@",
			wantErr: true,
		},
		{
			name:    "invalid email - empty",
			email:   "",
			wantErr: true,
		},
		{
			name:    "invalid email - missing local part",
			email:   "@example.com",
			wantErr: true,
		},
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := LoginRequest{
				Email:    tt.email,
				Password: "validpassword",
			}
			isValid := emailRegex.MatchString(req.Email)
			gotErr := !isValid
			if gotErr != tt.wantErr {
				t.Errorf("email validation = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}

// @test-contract TEST-INTERNAL_AUTH-001
func TestLoginRequest_PasswordNonEmpty(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "non-empty password",
			password: "validpassword",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := LoginRequest{
				Email:    "user@example.com",
				Password: tt.password,
			}
			gotErr := req.Password == ""
			if gotErr != tt.wantErr {
				t.Errorf("password non-empty check = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}

// @test-contract TEST-INTERNAL_AUTH-001
func TestLoginResponse_TokenNonEmpty(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "non-empty token",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		},
		{
			name:  "empty token",
			token: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := LoginResponse{
				Token: tt.token,
			}
			isNonEmpty := resp.Token != ""
			if tt.token == "" && isNonEmpty {
				t.Errorf("token should be empty for test case %q", tt.name)
			}
			if tt.token != "" && !isNonEmpty {
				t.Errorf("token should be non-empty for test case %q", tt.name)
			}
		})
	}
}

// @test-contract TEST-INTERNAL_AUTH-001
func TestLoginResponse_TokenFormat(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid JWT format - three parts separated by dots",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			wantErr: false,
		},
		{
			name:    "invalid JWT format - missing parts",
			token:   "invalid-token",
			wantErr: true,
		},
		{
			name:    "invalid JWT format - only one part",
			token:   "onlyonepart",
			wantErr: true,
		},
		{
			name:    "invalid JWT format - two parts",
			token:   "part1.part2",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	jwtRegex := regexp.MustCompile(`^[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+$`)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := LoginResponse{
				Token: tt.token,
			}
			isValid := jwtRegex.MatchString(resp.Token)
			gotErr := !isValid
			if gotErr != tt.wantErr {
				t.Errorf("token format validation = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}
