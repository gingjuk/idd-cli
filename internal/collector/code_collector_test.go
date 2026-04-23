// Package collector provides testing utilities for the collector module.

// Spec: docs/internal/collector/spec.md
// Test: docs/internal/collector/testing.md
package collector

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jingxu9x/idd-cli/internal/config"
	"github.com/jingxu9x/idd-cli/internal/model"
)

// @test TEST-INT_COL-027
func TestExtractFunctionComment(t *testing.T) {
	// idd:ignore start
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
				"// @implement SPEC-BE-001",
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
				"// @implement SPEC-BE-001",
				"func Authenticate() {}",
			},
			annotationLine: 0,
			wantContains:   "[function: Authenticate]",
		},
		{
			name: "no comment",
			lines: []string{
				"// @implement SPEC-BE-001",
				"func Authenticate() {}",
			},
			annotationLine: 0,
			wantContains:   "[function: Authenticate]",
		},
		{
			name: "comment with func name",
			lines: []string{
				"// @implement SPEC-BE-001",
				"// Processes user login request",
				"func ProcessLogin() {}",
			},
			annotationLine: 0,
			wantContains:   "Processes user login request",
		},
	}
	// idd:ignore end

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractFunctionComment(tt.lines, tt.annotationLine)
			if !strings.Contains(result, tt.wantContains) {
				t.Errorf("extractFunctionComment got %q, want contains %q", result, tt.wantContains)
			}
		})
	}
}

// @test TEST-INT_COL-022
func TestExtractFunctionComment_OutOfBounds(t *testing.T) {
	lines := []string{"// @implement SPEC-BE-001"} // idd:ignore
	result := extractFunctionComment(lines, 10)
	if result != "" {
		t.Errorf("Out of bounds should return empty, got %q", result)
	}

	result = extractFunctionComment([]string{}, 0)
	if result != "" {
		t.Errorf("Empty lines should return empty, got %q", result)
	}
}

// @test TEST-INT_COL-023
func TestExtractFunctionComment_PointerReceiver(t *testing.T) {
	// idd:ignore start
	lines := []string{
		"// @implement SPEC-BE-001",
		"// Validates pointer receiver",
		"func (s *Service) Validate() {}",
	}
	// idd:ignore end
	result := extractFunctionComment(lines, 0)
	if !strings.Contains(result, "Validates pointer receiver") {
		t.Errorf("Expected comment with pointer receiver, got %q", result)
	}
}

// @test TEST-INT_COL-024
func TestCodeCollector_Collect(t *testing.T) {
	tmpDir := t.TempDir()

	// idd:ignore start
	code := `package main

// @implement SPEC-CMD_IDD-001
func Authenticate() error {
	return nil
}

// @implement SPEC-CMD_IDD-002
type Payment interface {
	Process(amount float64) error
}

// @test TEST-CMD_IDD-001
// Test authentication
func TestAuth(t *testing.T) {}
`
	// idd:ignore end
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

	if !set.Has("SPEC-CMD_IDD-001") {
		t.Error("Should have SPEC-CMD_IDD-001")
	}
	if !set.Has("SPEC-CMD_IDD-002") {
		t.Error("Should have SPEC-CMD_IDD-002")
	}
	if !set.Has("TEST-CMD_IDD-001") {
		t.Error("Should have TEST-CMD_IDD-001")
	}

	specID, ok := set.Get("SPEC-CMD_IDD-001")
	if !ok {
		t.Fatal("SPEC-CMD_IDD-001 not found")
	}
	if specID.Describe == "" {
		t.Error("SPEC-CMD_IDD-001 should have Describe extracted from function comment")
	}
	if specID.Origin != model.OriginCode {
		t.Errorf("Origin = %v, want %v", specID.Origin, model.OriginCode)
	}
}

// @test TEST-INT_COL-025
func TestCodeCollector_CollectGoFile(t *testing.T) {
	tmpDir := t.TempDir()

	// idd:ignore start
	code := `package main

// @implement SPEC-CMD_IDD-002
func Register() error {
	return nil
}
`
	// idd:ignore end
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

	if !set.Has("SPEC-CMD_IDD-002") {
		t.Error("Should have SPEC-CMD_IDD-002")
	}
}

// @test TEST-INT_COL-026
func TestCodeCollector_MultipleAnnotations(t *testing.T) {
	tmpDir := t.TempDir()

	// idd:ignore start
	code := `package main

// @implement SPEC-CMD_IDD-001, SPEC-CMD_IDD-002
func Multiple() error {
	return nil
}
`
	// idd:ignore end
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

	if !set.Has("SPEC-CMD_IDD-001") {
		t.Error("Should have SPEC-CMD_IDD-001")
	}
	if !set.Has("SPEC-CMD_IDD-002") {
		t.Error("Should have SPEC-CMD_IDD-002")
	}
}
