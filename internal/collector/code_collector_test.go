package collector

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jingxu9x/idd-link-validator/internal/config"
	"github.com/jingxu9x/idd-link-validator/internal/model"
)

func TestExtractFunctionComment(t *testing.T) {
	tests := []struct {
		name           string
		lines          []string
		annotationLine int
		wantContains   string
	}{
		{
			name: "comment after annotation",
			lines: []string{
				"package main",
				"// @spec SPEC-BE-001",
				"// Validates user credentials",
				"// and issues JWT tokens",
				"func Authenticate() {}",
			},
			annotationLine: 1,
			wantContains:   "Validates user credentials",
		},
		{
			name: "function comment with func",
			lines: []string{
				"// @spec SPEC-BE-001",
				"func Authenticate() {}",
			},
			annotationLine: 0,
			wantContains:   "[function: Authenticate]",
		},
		{
			name: "no comment",
			lines: []string{
				"// @spec SPEC-BE-001",
				"func Authenticate() {}",
			},
			annotationLine: 0,
			wantContains:   "[function: Authenticate]",
		},
		{
			name: "comment with func name",
			lines: []string{
				"// @spec SPEC-BE-001",
				"// Processes user login request",
				"func ProcessLogin() {}",
			},
			annotationLine: 0,
			wantContains:   "Processes user login request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractFunctionComment(tt.lines, tt.annotationLine)
			if !strings.Contains(result, tt.wantContains) {
				t.Errorf("extractFunctionComment got %q, want contains %q", result, tt.wantContains)
			}
		})
	}
}

func TestExtractFunctionComment_OutOfBounds(t *testing.T) {
	lines := []string{"// @spec SPEC-BE-001"}
	result := extractFunctionComment(lines, 10)
	if result != "" {
		t.Errorf("Out of bounds should return empty, got %q", result)
	}

	result = extractFunctionComment([]string{}, 0)
	if result != "" {
		t.Errorf("Empty lines should return empty, got %q", result)
	}
}

func TestExtractFunctionComment_PointerReceiver(t *testing.T) {
	lines := []string{
		"// @spec SPEC-BE-001",
		"// Validates pointer receiver",
		"func (s *Service) Validate() {}",
	}
	result := extractFunctionComment(lines, 0)
	if !strings.Contains(result, "Validates pointer receiver") {
		t.Errorf("Expected comment with pointer receiver, got %q", result)
	}
}

func TestCodeCollector_Collect(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

// @spec SPEC-BE-001
// Validates user credentials
func Authenticate() error {
	return nil
}

// @contract CONTRACT-BE-001
// Defines payment interface
type Payment interface {
	Process(amount float64) error
}

// @test TEST-BE-001
// Test authentication
func TestAuth(t *testing.T) {}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cfg := config.Default()
	coll := NewCodeCollector(cfg)
	set, err := coll.Collect(nil, tmpDir)
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if !set.Has("SPEC-BE-001") {
		t.Error("Should have SPEC-BE-001")
	}
	if !set.Has("CONTRACT-BE-001") {
		t.Error("Should have CONTRACT-BE-001")
	}
	if !set.Has("TEST-BE-001") {
		t.Error("Should have TEST-BE-001")
	}

	specID, ok := set.Get("SPEC-BE-001")
	if !ok {
		t.Fatal("SPEC-BE-001 not found")
	}
	if specID.Describe == "" {
		t.Error("SPEC-BE-001 should have Describe extracted from function comment")
	}
	if specID.Origin != model.OriginCode {
		t.Errorf("Origin = %v, want %v", specID.Origin, model.OriginCode)
	}
}

func TestCodeCollector_CollectGoFile(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

// @spec SPEC-BE-002
// Handles user registration
func Register() error {
	return nil
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(code), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cfg := config.Default()
	coll := NewCodeCollector(cfg)
	set, err := coll.Collect(nil, tmpDir)
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if !set.Has("SPEC-BE-002") {
		t.Error("Should have SPEC-BE-002")
	}
}

func TestCodeCollector_MultipleAnnotations(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

// @spec SPEC-BE-001, SPEC-BE-002
func Multiple() error {
	return nil
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cfg := config.Default()
	coll := NewCodeCollector(cfg)
	set, err := coll.Collect(nil, tmpDir)
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if !set.Has("SPEC-BE-001") {
		t.Error("Should have SPEC-BE-001")
	}
	if !set.Has("SPEC-BE-002") {
		t.Error("Should have SPEC-BE-002")
	}
}
