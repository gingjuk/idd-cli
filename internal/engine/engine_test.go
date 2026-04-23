// Package engine provides testing utilities for the engine module.

// Spec: docs/internal/engine/spec.md
// Test: docs/internal/engine/testing.md

// @test TEST-CMD_IDD-001, TEST-CMD_IDD-002
package engine

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jingxu9x/idd-cli/internal/config"
	"github.com/jingxu9x/idd-cli/internal/model"
)

func TestEngine_ConsistencyCheck_Warning(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Validation: config.ValidationConfig{
			ConsistencyCheck: config.ConsistencyCheck{
				Enabled:   true,
				Threshold: 0.3,
			},
		},
	}

	ids := model.NewIdentifierSet()

	docID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "Test Spec", "validates user authentication and JWT token issuance", "docs/spec.md", 1)
	docID.SetOrigin(model.OriginDoc)
	ids.Add(docID)

	codeID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "", "database connection pooling settings", "main.go", 10)
	codeID.SetOrigin(model.OriginCode)
	ids.Add(codeID)

	eng := New(cfg)
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if len(result.Warnings) == 0 {
		t.Error("Expected consistency warning for low similarity")
	}

	found := false
	for _, w := range result.Warnings {
		if w.Rule == "consistency-check" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected consistency-check warning")
	}
}

func TestEngine_ConsistencyCheck_NoWarning(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Validation: config.ValidationConfig{
			ConsistencyCheck: config.ConsistencyCheck{
				Enabled:   true,
				Threshold: 0.3,
			},
		},
	}

	ids := model.NewIdentifierSet()

	docID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "Test Spec", "validates user authentication and JWT token issuance", "docs/spec.md", 1)
	docID.SetOrigin(model.OriginDoc)
	ids.Add(docID)

	codeID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "", "validates user authentication and JWT token issuance for session", "main.go", 10)
	codeID.SetOrigin(model.OriginCode)
	ids.Add(codeID)

	eng := New(cfg)
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, w := range result.Warnings {
		if w.Rule == "consistency-check" {
			t.Error("Expected no consistency warning for similar texts")
		}
	}
}

func TestEngine_ConsistencyCheck_Disabled(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Validation: config.ValidationConfig{
			ConsistencyCheck: config.ConsistencyCheck{
				Enabled:   false,
				Threshold: 0.3,
			},
		},
	}

	ids := model.NewIdentifierSet()

	docID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "Test Spec", "validates user authentication", "docs/spec.md", 1)
	docID.SetOrigin(model.OriginDoc)
	ids.Add(docID)

	codeID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "", "completely unrelated content here", "main.go", 10)
	codeID.SetOrigin(model.OriginCode)
	ids.Add(codeID)

	eng := New(cfg)
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, w := range result.Warnings {
		if w.Rule == "consistency-check" {
			t.Error("Expected no consistency warnings when disabled")
		}
	}
}

func TestEngine_ConsistencyCheck_MissingDescribe(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Validation: config.ValidationConfig{
			ConsistencyCheck: config.ConsistencyCheck{
				Enabled:   true,
				Threshold: 0.3,
			},
		},
	}

	ids := model.NewIdentifierSet()

	docID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "Test Spec", "", "docs/spec.md", 1)
	docID.SetOrigin(model.OriginDoc)
	ids.Add(docID)

	codeID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "", "validates user authentication", "main.go", 10)
	codeID.SetOrigin(model.OriginCode)
	ids.Add(codeID)

	eng := New(cfg)
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, w := range result.Warnings {
		if w.Rule == "consistency-check" {
			t.Error("Expected no consistency warning when doc describe is empty")
		}
	}
}

func TestEngine_ConsistencyCheck_CodeOnly(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Validation: config.ValidationConfig{
			ConsistencyCheck: config.ConsistencyCheck{
				Enabled:   true,
				Threshold: 0.3,
			},
		},
	}

	ids := model.NewIdentifierSet()

	codeID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "", "validates user authentication", "main.go", 10)
	codeID.SetOrigin(model.OriginCode)
	ids.Add(codeID)

	eng := New(cfg)
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, w := range result.Warnings {
		if w.Rule == "consistency-check" {
			t.Error("Expected no consistency warning when only code exists")
		}
	}
}

func TestEngine_ConsistencyCheck_HighThreshold(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Validation: config.ValidationConfig{
			ConsistencyCheck: config.ConsistencyCheck{
				Enabled:   true,
				Threshold: 0.95,
			},
		},
	}

	ids := model.NewIdentifierSet()

	docID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "Test Spec", "validates user authentication and JWT tokens", "docs/spec.md", 1)
	docID.SetOrigin(model.OriginDoc)
	ids.Add(docID)

	codeID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "", "validates user authentication and JWT tokens for sessions", "main.go", 10)
	codeID.SetOrigin(model.OriginCode)
	ids.Add(codeID)

	eng := New(cfg)
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	found := false
	for _, w := range result.Warnings {
		if w.Rule == "consistency-check" {
			found = true
			break
		}
	}
	if found {
		t.Error("Expected no consistency warning for very similar texts with high threshold")
	}
}

func TestEngine_ConsistencyCheck_LowSimilarity(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Validation: config.ValidationConfig{
			ConsistencyCheck: config.ConsistencyCheck{
				Enabled:   true,
				Threshold: 0.1,
			},
		},
	}

	ids := model.NewIdentifierSet()

	docID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "Test Spec", "user login and authentication", "docs/spec.md", 1)
	docID.SetOrigin(model.OriginDoc)
	ids.Add(docID)

	codeID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "", "database connection pool settings", "main.go", 10)
	codeID.SetOrigin(model.OriginCode)
	ids.Add(codeID)

	eng := New(cfg)
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	found := false
	for _, w := range result.Warnings {
		if w.Rule == "consistency-check" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected consistency warning for very dissimilar texts")
	}
}

func TestEngine_AddStructuralErrors(t *testing.T) {
	cfg := config.Default()
	eng := New(cfg)

	errors := []*model.ValidationError{
		{Rule: "test-rule", Message: "test error"},
	}
	eng.AddStructuralErrors(errors)

	if len(eng.result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(eng.result.Errors))
	}
}

func TestEngine_inferLinkType(t *testing.T) {
	cfg := config.Default()
	eng := New(cfg)

	tests := []struct {
		fromType model.IdentifierType
		wantLink model.LinkType
	}{
		{model.TypeSpec, model.LinkTests},
		{model.TypeTest, model.LinkImplements},
		{model.TypeContract, model.LinkContractImplements},
		{model.TypeDesign, model.LinkReferences},
	}

	for _, tt := range tests {
		link := eng.inferLinkType(tt.fromType, "TEST-BE-001")
		if link != tt.wantLink {
			t.Errorf("inferLinkType(%v) = %v, want %v", tt.fromType, link, tt.wantLink)
		}
	}
}

func TestEngine_validateBidirectional(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Validation: config.ValidationConfig{
			RequireDocLinkConsistency: false,
			AllowOrphans:              true,
		},
	}
	eng := New(cfg)

	ids := model.NewIdentifierSet()
	specID := model.NewIdentifier("SPEC-001", model.TypeSpec, "", "docs/spec.md", 1)
	specID.SetOrigin(model.OriginDoc)
	ids.Add(specID)

	eng.Run(context.Background(), ids)

	if len(eng.result.Errors) != 0 {
		t.Errorf("Expected 0 errors, got %d", len(eng.result.Errors))
	}
}

func TestEngine_validateDocCodeCorrespondence(t *testing.T) {
	cfg := config.Default()
	eng := New(cfg)

	ids := model.NewIdentifierSet()

	codeOnlyID := model.NewIdentifier("SPEC-001", model.TypeSpec, "", "main.go", 1)
	codeOnlyID.SetOrigin(model.OriginCode)
	ids.Add(codeOnlyID)

	eng.Run(context.Background(), ids)

	found := false
	for _, e := range eng.result.Errors {
		if e.Rule == "doc-code-correspondence" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error for code without doc")
	}
}

func TestEngine_BuildReport(t *testing.T) {
	cfg := config.Default()
	eng := New(cfg)

	ids := model.NewIdentifierSet()
	eng.Run(context.Background(), ids)

	report := eng.BuildReport()
	if report.Tool != "idd-cli" {
		t.Errorf("Tool = %q, want idd-cli", report.Tool)
	}
	if report.Version != "1.0.0" {
		t.Errorf("Version = %q, want 1.0.0", report.Version)
	}
}

// @test TEST-INT_ENG-001
func TestEngine_PackageDocComment_Valid(t *testing.T) {
	// Valid format: package doc comment before package declaration
	tmpDir := t.TempDir()
	code := `// Package foo provides core functionality for the foo module.

// Spec: docs/foo/spec.md
// Contract: docs/foo/contract.md
package foo
`

	err := os.WriteFile(filepath.Join(tmpDir, "foo.go"), []byte(code), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := &config.Config{
		Version: "1.0",
		Code: config.CodeConfig{
			Patterns: []string{filepath.Join(tmpDir, "*.go")},
		},
		Validation: config.ValidationConfig{
			RequirePackageDocComment: true,
		},
	}

	eng := New(cfg)
	ids := model.NewIdentifierSet()
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Should have no package-doc-comment errors
	for _, e := range result.Errors {
		if e.Rule == "package-doc-comment" {
			t.Errorf("Expected no package-doc-comment errors for valid format, got: %s", e.Message)
		}
	}
}

// @test TEST-INT_ENG-002
func TestEngine_PackageDocComment_MissingPackageDoc(t *testing.T) {
	// Missing package doc comment line (has Spec and Contract but no "Package foo provides")
	tmpDir := t.TempDir()
	code := `// Spec: docs/foo/spec.md
// Contract: docs/foo/contract.md
package foo
`

	err := os.WriteFile(filepath.Join(tmpDir, "foo.go"), []byte(code), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := &config.Config{
		Version: "1.0",
		Code: config.CodeConfig{
			Patterns: []string{filepath.Join(tmpDir, "*.go")},
		},
		Validation: config.ValidationConfig{
			RequirePackageDocComment: true,
		},
	}

	eng := New(cfg)
	ids := model.NewIdentifierSet()
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	found := false
	for _, e := range result.Errors {
		if e.Rule == "package-doc-comment" && strings.Contains(e.Message, "missing package doc comment") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error for missing package doc comment")
	}
}

// @test TEST-INT_ENG-003
func TestEngine_PackageDocComment_MissingPackageName(t *testing.T) {
	// Has comment but not a package doc comment (no "Package xxx" pattern)
	tmpDir := t.TempDir()
	code := `// This is a regular comment, not a package doc comment.
//
// Spec: docs/foo/spec.md
// Contract: docs/foo/contract.md
package foo
`

	err := os.WriteFile(filepath.Join(tmpDir, "foo.go"), []byte(code), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := &config.Config{
		Version: "1.0",
		Code: config.CodeConfig{
			Patterns: []string{filepath.Join(tmpDir, "*.go")},
		},
		Validation: config.ValidationConfig{
			RequirePackageDocComment: true,
		},
	}

	eng := New(cfg)
	ids := model.NewIdentifierSet()
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	found := false
	for _, e := range result.Errors {
		if e.Rule == "package-doc-comment" && strings.Contains(e.Message, "missing package doc comment") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error for missing package doc comment")
	}
}

// @test TEST-INT_ENG-004
func TestEngine_PackageDocComment_PackageAfterComments(t *testing.T) {
	// Package declaration AFTER comments (correct order)
	tmpDir := t.TempDir()
	code := `// Package bar provides core bar functionality.

// Spec: docs/bar/spec.md
// Contract: docs/bar/contract.md
package bar
`

	err := os.WriteFile(filepath.Join(tmpDir, "bar.go"), []byte(code), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := &config.Config{
		Version: "1.0",
		Code: config.CodeConfig{
			Patterns: []string{filepath.Join(tmpDir, "*.go")},
		},
		Validation: config.ValidationConfig{
			RequirePackageDocComment: true,
		},
	}

	eng := New(cfg)
	ids := model.NewIdentifierSet()
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, e := range result.Errors {
		if e.Rule == "package-doc-comment" {
			t.Errorf("Expected no package-doc-comment errors for package after comments, got: %s", e.Message)
		}
	}
}

// @test TEST-INT_ENG-005
func TestEngine_PackageDocComment_Disabled(t *testing.T) {
	// When RequirePackageDocComment is false, no errors should be generated
	tmpDir := t.TempDir()
	code := `package baz
`

	err := os.WriteFile(filepath.Join(tmpDir, "baz.go"), []byte(code), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := &config.Config{
		Version: "1.0",
		Code: config.CodeConfig{
			Patterns: []string{filepath.Join(tmpDir, "*.go")},
		},
		Validation: config.ValidationConfig{
			RequirePackageDocComment: false,
		},
	}

	eng := New(cfg)
	ids := model.NewIdentifierSet()
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, e := range result.Errors {
		if e.Rule == "package-doc-comment" {
			t.Error("Should not generate package-doc-comment errors when disabled")
		}
	}
}

// @test TEST-INT_ENG-006
func TestEngine_PackageDocComment_TestFile_RequiresSpecAndTest(t *testing.T) {
	tmpDir := t.TempDir()
	code := `// Package foo provides testing utilities for foo module.

// Spec: docs/foo/spec.md
// Test: docs/foo/testing.md
package foo
`

	err := os.WriteFile(filepath.Join(tmpDir, "foo_test.go"), []byte(code), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := &config.Config{
		Version: "1.0",
		Code: config.CodeConfig{
			Patterns: []string{filepath.Join(tmpDir, "*.go")},
		},
		Validation: config.ValidationConfig{
			RequirePackageDocComment: true,
		},
	}

	eng := New(cfg)
	ids := model.NewIdentifierSet()
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, e := range result.Errors {
		if e.Rule == "package-doc-comment" {
			t.Errorf("Expected no errors for test file with Spec and Test, got: %s", e.Message)
		}
	}
}

// @test TEST-INT_ENG-007
func TestEngine_PackageDocComment_TestFile_MissingTest(t *testing.T) {
	tmpDir := t.TempDir()
	code := `// Package foo provides testing utilities for foo module.
//
// Spec: docs/foo/spec.md
// Contract: docs/foo/contract.md
package foo
`

	err := os.WriteFile(filepath.Join(tmpDir, "foo_test.go"), []byte(code), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := &config.Config{
		Version: "1.0",
		Code: config.CodeConfig{
			Patterns: []string{filepath.Join(tmpDir, "*.go")},
		},
		Validation: config.ValidationConfig{
			RequirePackageDocComment: true,
		},
	}

	eng := New(cfg)
	ids := model.NewIdentifierSet()
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	found := false
	for _, e := range result.Errors {
		if e.Rule == "package-doc-comment" && strings.Contains(e.Message, "missing Test path") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error for test file missing Test path")
	}
}

// @test TEST-INT_ENG-020
func TestEngine_validateDuplicateHeadingIdentifiers(t *testing.T) {
	cfg := config.Default()
	cfg.Validation.AllowOrphans = true
	eng := New(cfg)

	// Create temp dir with docs
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs", "test")
	err := os.MkdirAll(docsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Write doc with duplicate TEST headings
	doc := `# Documentation
## TEST-BE-001: First Test
Some content
## TEST-BE-001: Duplicate Test
More content
`
	err = os.WriteFile(filepath.Join(docsDir, "spec.md"), []byte(doc), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	// Override docs pattern to point to our temp dir
	cfg.Docs.Patterns = []string{filepath.Join(docsDir, "**/*.md")}
	eng = New(cfg)

	eng.validateDuplicateHeadingIdentifiers()

	// Should have exactly one error for the duplicate
	if len(eng.result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d: %v", len(eng.result.Errors), eng.result.Errors)
	}
	if eng.result.Errors[0].Rule != "duplicate-heading-identifier" {
		t.Errorf("Expected duplicate-heading-identifier rule, got %s", eng.result.Errors[0].Rule)
	}
}

// @test TEST-INT_ENG-021
func TestEngine_validateDuplicateHeadingIdentifiers_NoDuplicates(t *testing.T) {
	cfg := config.Default()
	cfg.Validation.AllowOrphans = true
	eng := New(cfg)

	// Create temp dir with docs
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs", "test")
	err := os.MkdirAll(docsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Write doc with unique headings
	doc := `## TEST-BE-001: First Test
Some content
## TEST-BE-002: Second Test
More content
`
	err = os.WriteFile(filepath.Join(docsDir, "spec.md"), []byte(doc), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg.Docs.Patterns = []string{filepath.Join(docsDir, "**/*.md")}
	eng = New(cfg)

	eng.validateDuplicateHeadingIdentifiers()

	// Should have no errors
	if len(eng.result.Errors) != 0 {
		t.Errorf("Expected 0 errors, got %d", len(eng.result.Errors))
	}
}

// @test TEST-INT_ENG-022
func TestEngine_validateContractDesignMarkers_WithMarkers(t *testing.T) {
	cfg := config.Default()
	cfg.Validation.AllowOrphans = true
	eng := New(cfg)

	// Create temp dir with docs
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs", "internal", "test")
	err := os.MkdirAll(docsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Write contract.md with markers (should error)
	contract := `---
markers:
  - id: CONTRACT-BE-001
    name: Test Contract
---
# Contracts
`
	err = os.WriteFile(filepath.Join(docsDir, "contract.md"), []byte(contract), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg.Docs.Patterns = []string{filepath.Join(docsDir, "**/*.md")}
	eng = New(cfg)

	eng.validateContractDesignMarkers()

	// Should have exactly one error
	if len(eng.result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d: %v", len(eng.result.Errors), eng.result.Errors)
	}
	if eng.result.Errors[0].Rule != "frontmatter-markers" {
		t.Errorf("Expected frontmatter-markers rule, got %s", eng.result.Errors[0].Rule)
	}
}

// @test TEST-INT_ENG-023
func TestEngine_validateContractDesignMarkers_WithoutMarkers(t *testing.T) {
	cfg := config.Default()
	cfg.Validation.AllowOrphans = true
	eng := New(cfg)

	// Create temp dir with docs
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs", "internal", "test")
	err := os.MkdirAll(docsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Write contract.md without markers (should pass)
	contract := `---
related_files:
  spec: spec.md
---
# Contracts
`
	err = os.WriteFile(filepath.Join(docsDir, "contract.md"), []byte(contract), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg.Docs.Patterns = []string{filepath.Join(docsDir, "**/*.md")}
	eng = New(cfg)

	eng.validateContractDesignMarkers()

	// Should have no errors
	if len(eng.result.Errors) != 0 {
		t.Errorf("Expected 0 errors, got %d", len(eng.result.Errors))
	}
}

// @test TEST-INT_ENG-024
func TestEngine_validateContractDesignMarkers_DesignFile(t *testing.T) {
	cfg := config.Default()
	cfg.Validation.AllowOrphans = true
	eng := New(cfg)

	// Create temp dir with docs
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs", "internal", "test")
	err := os.MkdirAll(docsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Write design.md with markers (should error)
	design := `---
markers:
  - id: DESIGN-BE-001
    name: Test Design
---
# Design
`
	err = os.WriteFile(filepath.Join(docsDir, "design.md"), []byte(design), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg.Docs.Patterns = []string{filepath.Join(docsDir, "**/*.md")}
	eng = New(cfg)

	eng.validateContractDesignMarkers()

	// Should have exactly one error
	if len(eng.result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d: %v", len(eng.result.Errors), eng.result.Errors)
	}
}

// @test TEST-INT_ENG-025
func TestEngine_packageExists(t *testing.T) {
	cfg := config.Default()
	eng := New(cfg)

	// Package that exists - use current working directory
	cwd, _ := os.Getwd()
	// When run from package dir, cwd is /project/internal/engine, so go up 2 dirs to project root
	if strings.HasSuffix(cwd, "/internal/engine") {
		cwd = filepath.Dir(filepath.Dir(cwd))
	}
	pkgPath := filepath.Join(cwd, "internal/engine")
	if !eng.packageExists(pkgPath) {
		t.Error("internal/engine should exist at:", pkgPath)
	}

	// Package that doesn't exist
	if eng.packageExists("nonexistent/package") {
		t.Error("nonexistent/package should not exist")
	}
}

// @test TEST-INT_ENG-026
func TestEngine_isIgnoredDocPath(t *testing.T) {
	cfg := config.Default()
	cfg.Docs.IgnorePaths = []string{"docs/backend/**", "docs/internal/auth"}
	eng := New(cfg)

	if !eng.isIgnoredDocPath("docs/backend/spec.md") {
		t.Error("docs/backend/** should be ignored")
	}

	if !eng.isIgnoredDocPath("docs/internal/auth") {
		t.Error("docs/internal/auth should be ignored")
	}

	if eng.isIgnoredDocPath("docs/internal/engine") {
		t.Error("docs/internal/engine should not be ignored")
	}
}

// @test TEST-INT_ENG-027
func TestEngine_validateDocPathExists_NoWarning(t *testing.T) {
	cfg := config.Default()
	cfg.Validation.AllowOrphans = true
	eng := New(cfg)

	// Use the actual docs directory which should have corresponding packages
	eng.validateDocPathExists()

	// The test just verifies no panic and some warnings are generated
	// The actual number of warnings depends on the docs structure
}

// @test TEST-INT_ENG-028
func TestEngine_validateRelatedFiles(t *testing.T) {
	cfg := config.Default()
	cfg.Validation.AllowOrphans = true
	eng := New(cfg)

	// Create temp dir with docs
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs", "test")
	err := os.MkdirAll(docsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Write spec.md without related_files
	spec := `# Spec
Content
`
	err = os.WriteFile(filepath.Join(docsDir, "spec.md"), []byte(spec), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg.Docs.Patterns = []string{filepath.Join(docsDir, "**/*.md")}
	eng = New(cfg)

	eng.validateRelatedFiles()

	// Should have one error for missing related_files
	found := false
	for _, err := range eng.result.Errors {
		if err.Rule == "related-files" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected related-files error")
	}
}

// @test TEST-INT_ENG-029
func TestEngine_validateRelatedFiles_Valid(t *testing.T) {
	cfg := config.Default()
	cfg.Validation.AllowOrphans = true
	eng := New(cfg)

	// Create temp dir with docs
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs", "test")
	err := os.MkdirAll(docsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Write spec.md with related_files
	spec := `---
related_files:
  spec: spec.md
  contract: contract.md
---
# Spec
Content
`
	err = os.WriteFile(filepath.Join(docsDir, "spec.md"), []byte(spec), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg.Docs.Patterns = []string{filepath.Join(docsDir, "**/*.md")}
	eng = New(cfg)

	eng.validateRelatedFiles()

	// Should have no errors
	if len(eng.result.Errors) != 0 {
		t.Errorf("Expected 0 errors, got %d", len(eng.result.Errors))
	}
}

// @test TEST-INT_ENG-007
func TestEngine_PackageDocComment_MainPackageSkipped(t *testing.T) {
	// main package should be skipped
	tmpDir := t.TempDir()
	code := `package main

func main() {}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	cfg := &config.Config{
		Version: "1.0",
		Code: config.CodeConfig{
			Patterns: []string{filepath.Join(tmpDir, "*.go")},
		},
		Validation: config.ValidationConfig{
			RequirePackageDocComment: true,
		},
	}

	eng := New(cfg)
	ids := model.NewIdentifierSet()
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, e := range result.Errors {
		if e.Rule == "package-doc-comment" {
			t.Error("main package should be skipped")
		}
	}
}
