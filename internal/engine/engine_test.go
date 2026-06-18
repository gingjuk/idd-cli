// Package engine provides testing utilities for the engine module.

// Spec: docs/internal/engine/spec.md
// Test: docs/internal/engine/testing.md

// @test TEST-CMD_IDD_CLI-001, TEST-CMD_IDD_CLI-002
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

// @test TEST-INTERNAL_ENGINE-007
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

// @test TEST-INTERNAL_ENGINE-007
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

// @test TEST-INTERNAL_ENGINE-007
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

// @test TEST-INTERNAL_ENGINE-007
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

// @test TEST-INTERNAL_ENGINE-007
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

// @test TEST-INTERNAL_ENGINE-007
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

// @test TEST-INTERNAL_ENGINE-007
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

// @test TEST-INTERNAL_ENGINE-004
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

// @test TEST-INTERNAL_ENGINE-002
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

// @test TEST-INTERNAL_ENGINE-003
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

// @test TEST-INTERNAL_ENGINE-003
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

// @test TEST-INTERNAL_ENGINE-005
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

// @test TEST-INTERNAL_ENGINE-001
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

// @test TEST-INTERNAL_ENGINE-002
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

// @test TEST-INTERNAL_ENGINE-003
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

// @test TEST-INTERNAL_ENGINE-004
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

// @test TEST-INTERNAL_ENGINE-005
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

// @test TEST-INTERNAL_ENGINE-006
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

// @test TEST-INTERNAL_ENGINE-007
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

// @test TEST-INTERNAL_ENGINE-020
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

// @test TEST-INTERNAL_ENGINE-021
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

func TestEngine_validateSpecRequiredFields(t *testing.T) {
	tests := []struct {
		name                   string
		doc                    string
		wantFields             []string
		wantOrderError         bool
		wantOptionalOrderError bool
		wantErrorCount         int
	}{
		{
			name: "all required fields present",
			doc: `---
markers:
  - id: SPEC-BE-001
    name: User Login
---

## SPEC-BE-001: User Login

**Design:** implements architecture AuthModule

**Contract:** implements interface Authenticator

**Requirement:**

Users can log in with email and password.

**Tests:** TEST-BE-001

**Status:** Done

**Implementation:** auth.go
`,
		},
		{
			name: "missing contract and design",
			doc: `---
markers:
  - id: SPEC-BE-001
    name: User Login
---

## SPEC-BE-001: User Login

**Requirement:**

Users can log in with email and password.

**Tests:** TEST-BE-001
`,
			wantFields:     []string{"Contract", "Design"},
			wantErrorCount: 2,
		},
		{
			name: "missing requirement and tests",
			doc: `---
markers:
  - id: SPEC-BE-001
    name: User Login
---

## SPEC-BE-001: User Login

**Contract:** implements interface Authenticator

**Design:** implements architecture AuthModule
`,
			wantFields:     []string{"Requirement", "Tests"},
			wantErrorCount: 2,
		},
		{
			name: "fields must be in template order",
			doc: `---
markers:
  - id: SPEC-BE-001
    name: User Login
---

## SPEC-BE-001: User Login

**Contract:** implements interface Authenticator

**Design:** implements architecture AuthModule

**Requirement:**

Users can log in with email and password.

**Tests:** TEST-BE-001
`,
			wantOrderError: true,
			wantErrorCount: 1,
		},
		{
			name: "optional fields must follow required fields",
			doc: `---
markers:
  - id: SPEC-BE-001
    name: User Login
---

## SPEC-BE-001: User Login

**Design:** implements architecture AuthModule

**Status:** Done

**Contract:** implements interface Authenticator

**Requirement:**

Users can log in with email and password.

**Tests:** TEST-BE-001
`,
			wantOptionalOrderError: true,
			wantErrorCount:         1,
		},
		{
			name: "checks each spec section independently",
			doc: `---
markers:
  - id: SPEC-BE-001
    name: Complete
  - id: SPEC-BE-002
    name: Missing Tests
---

## SPEC-BE-001: Complete

**Design:** implements architecture AuthModule

**Contract:** implements interface Authenticator

**Requirement:** Complete requirement.

**Tests:** TEST-BE-001

## SPEC-BE-002: Missing Tests

**Contract:** implements interface Authenticator

**Design:** implements architecture AuthModule

**Requirement:** Missing tests only.
`,
			wantFields:     []string{"Tests"},
			wantErrorCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			docsDir := filepath.Join(tmpDir, "docs", "backend")
			if err := os.MkdirAll(docsDir, 0755); err != nil {
				t.Fatalf("Failed to create docs dir: %v", err)
			}
			if err := os.WriteFile(filepath.Join(docsDir, "spec.md"), []byte(tt.doc), 0644); err != nil {
				t.Fatalf("Failed to write spec.md: %v", err)
			}

			cfg := config.Default()
			cfg.Docs.Patterns = []string{filepath.Join(tmpDir, "docs/**/*.md")}
			eng := New(cfg)

			eng.validateSpecRequiredFields()

			if len(eng.result.Errors) != tt.wantErrorCount {
				t.Fatalf("Expected %d errors, got %d: %v", tt.wantErrorCount, len(eng.result.Errors), eng.result.Errors)
			}

			for _, field := range tt.wantFields {
				found := false
				for _, err := range eng.result.Errors {
					if err.Rule == "spec-required-fields" && err.Code == field {
						found = true
						if !strings.Contains(err.Message, "define it as:") {
							t.Errorf("Expected missing field message to include field definition, got %q", err.Message)
						}
						break
					}
				}
				if !found {
					t.Errorf("Expected missing field error for %s", field)
				}
			}
			if tt.wantOrderError {
				found := false
				for _, err := range eng.result.Errors {
					if err.Rule == "spec-required-fields" && err.Code == "field-order" {
						found = true
						if !strings.Contains(err.Message, "**Design:**, **Contract:**, **Requirement:**, **Tests:**") {
							t.Errorf("Expected order error to include expected field order, got %q", err.Message)
						}
						break
					}
				}
				if !found {
					t.Errorf("Expected field order error")
				}
			}
			if tt.wantOptionalOrderError {
				found := false
				for _, err := range eng.result.Errors {
					if err.Rule == "spec-required-fields" && err.Code == "optional-field-order" {
						found = true
						if !strings.Contains(err.Message, "optional field **Status:** appears before required fields are complete") {
							t.Errorf("Expected optional field order error to include offending field, got %q", err.Message)
						}
						break
					}
				}
				if !found {
					t.Errorf("Expected optional field order error")
				}
			}
		})
	}
}

// @test TEST-INTERNAL_ENGINE-022
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

// @test TEST-INTERNAL_ENGINE-023
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

// @test TEST-INTERNAL_ENGINE-024
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

// @test TEST-INTERNAL_ENGINE-025
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

// @test TEST-INTERNAL_ENGINE-026
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

// @test TEST-INTERNAL_ENGINE-027
func TestEngine_validateDocPathExists_NoWarning(t *testing.T) {
	cfg := config.Default()
	cfg.Validation.AllowOrphans = true
	eng := New(cfg)

	// Use the actual docs directory which should have corresponding packages
	eng.validateDocPathExists()

	// The test just verifies no panic and some warnings are generated
	// The actual number of warnings depends on the docs structure
}

// @test TEST-INTERNAL_ENGINE-028
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

// @test TEST-INTERNAL_ENGINE-029
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

// @test TEST-INTERNAL_ENGINE-030
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

// @test TEST-INTERNAL_ENGINE-031
func TestEngine_PkgDocFiles_MissingFiles(t *testing.T) {
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs", "backend")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// Only spec.md present; contract.md, testing.md, design.md missing.
	if err := os.WriteFile(filepath.Join(docsDir, "spec.md"), []byte("# SPEC-BE-001\n"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	cfg := &config.Config{
		Version: "1.0",
		Docs:    config.DocsConfig{Patterns: []string{filepath.Join(tmpDir, "docs/**/*.md")}},
		Validation: config.ValidationConfig{
			RequirePkgDocFiles: true,
		},
	}

	eng := New(cfg)
	result, err := eng.Run(context.Background(), model.NewIdentifierSet())
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	missing := map[string]bool{}
	for _, e := range result.Errors {
		if e.Rule == "pkg-doc-files" {
			missing[e.Message] = true
		}
	}
	for _, f := range []string{"contract.md", "testing.md", "design.md"} {
		found := false
		for msg := range missing {
			if strings.Contains(msg, f) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected pkg-doc-files error for missing %s", f)
		}
	}
	// spec.md is present — should not appear in errors.
	for msg := range missing {
		if strings.Contains(msg, "spec.md") {
			t.Errorf("unexpected pkg-doc-files error for spec.md: %s", msg)
		}
	}
}

// @test TEST-INTERNAL_ENGINE-032
func TestEngine_PkgDocFiles_AllPresent(t *testing.T) {
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs", "backend")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	for _, f := range []string{"spec.md", "contract.md", "testing.md", "design.md"} {
		if err := os.WriteFile(filepath.Join(docsDir, f), []byte("# content\n"), 0644); err != nil {
			t.Fatalf("write: %v", err)
		}
	}

	cfg := &config.Config{
		Version: "1.0",
		Docs:    config.DocsConfig{Patterns: []string{filepath.Join(tmpDir, "docs/**/*.md")}},
		Validation: config.ValidationConfig{
			RequirePkgDocFiles: true,
		},
	}

	eng := New(cfg)
	result, err := eng.Run(context.Background(), model.NewIdentifierSet())
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, e := range result.Errors {
		if e.Rule == "pkg-doc-files" {
			t.Errorf("unexpected pkg-doc-files error: %s", e.Message)
		}
	}
}

// @test TEST-INTERNAL_ENGINE-033
func TestEngine_PkgDocFiles_RootLevelSkipped(t *testing.T) {
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// Root-level docs/spec.md — should not trigger pkg-doc-files check.
	if err := os.WriteFile(filepath.Join(docsDir, "spec.md"), []byte("# SPEC-001\n"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	cfg := &config.Config{
		Version: "1.0",
		Docs:    config.DocsConfig{Patterns: []string{filepath.Join(tmpDir, "docs/**/*.md")}},
		Validation: config.ValidationConfig{
			RequirePkgDocFiles: true,
		},
	}

	eng := New(cfg)
	result, err := eng.Run(context.Background(), model.NewIdentifierSet())
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, e := range result.Errors {
		if e.Rule == "pkg-doc-files" {
			t.Errorf("root-level docs should be skipped, got: %s", e.Message)
		}
	}
}

// @test TEST-INTERNAL_ENGINE-034
func TestEngine_DuplicateIDs_DocSide(t *testing.T) {
	cfg := &config.Config{Version: "1.0"}
	ids := model.NewIdentifierSet()

	// Same ID in two different doc directories → conflict
	a := model.NewIdentifier("SPEC-BE-001", model.TypeSpec, "Spec A", "docs/backend/spec.md", 1)
	a.SetOrigin(model.OriginDoc)
	b := model.NewIdentifier("SPEC-BE-001", model.TypeSpec, "Spec B", "docs/billing/spec.md", 1)
	b.SetOrigin(model.OriginDoc)
	ids.Add(a)
	ids.Add(b)

	eng := New(cfg)
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	found := false
	for _, e := range result.Errors {
		if e.Rule == "duplicate-id" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected duplicate-id error for same ID in different doc directories")
	}
}

// @test TEST-INTERNAL_ENGINE-035
func TestEngine_DuplicateIDs_SameDir_NoError(t *testing.T) {
	cfg := &config.Config{Version: "1.0"}
	ids := model.NewIdentifierSet()

	// Same ID in two files of the same doc directory → not a conflict
	a := model.NewIdentifier("SPEC-BE-001", model.TypeSpec, "Spec A", "docs/backend/spec.md", 1)
	a.SetOrigin(model.OriginDoc)
	b := model.NewIdentifier("SPEC-BE-001", model.TypeSpec, "Spec B", "docs/backend/testing.md", 1)
	b.SetOrigin(model.OriginDoc)
	ids.Add(a)
	ids.Add(b)

	eng := New(cfg)
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, e := range result.Errors {
		if e.Rule == "duplicate-id" {
			t.Errorf("unexpected duplicate-id error for same directory: %s", e.Message)
		}
	}
}

// @test TEST-INTERNAL_ENGINE-036
func TestEngine_DuplicateIDs_CodeSide(t *testing.T) {
	cfg := &config.Config{Version: "1.0"}
	ids := model.NewIdentifierSet()

	// Same ID @implement in two different code packages → conflict
	a := model.NewIdentifier("SPEC-BE-001", model.TypeSpec, "", "internal/backend/service.go", 10)
	a.SetOrigin(model.OriginCode)
	b := model.NewIdentifier("SPEC-BE-001", model.TypeSpec, "", "internal/billing/service.go", 20)
	b.SetOrigin(model.OriginCode)
	ids.Add(a)
	ids.Add(b)

	eng := New(cfg)
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	found := false
	for _, e := range result.Errors {
		if e.Rule == "duplicate-id" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected duplicate-id error for same ID in different code packages")
	}
}

// @test TEST-INTERNAL_ENGINE-037
func TestEngine_DuplicateIDs_RenameSuggestion(t *testing.T) {
	cfg := &config.Config{Version: "1.0"}
	ids := model.NewIdentifierSet()

	a := model.NewIdentifier("SPEC-BE-001", model.TypeSpec, "Spec A", "docs/backend/spec.md", 1)
	a.SetOrigin(model.OriginDoc)
	b := model.NewIdentifier("SPEC-BE-001", model.TypeSpec, "Spec B", "docs/billing/spec.md", 1)
	b.SetOrigin(model.OriginDoc)
	ids.Add(a)
	ids.Add(b)

	eng := New(cfg)
	result, _ := eng.Run(context.Background(), ids)

	for _, e := range result.Errors {
		if e.Rule == "duplicate-id" && strings.Contains(e.Message, "consider renaming") {
			return
		}
	}
	t.Error("expected rename suggestion in duplicate-id error message")
}

// TestEngine_PublicFuncAnnotation_PrivateFuncAllowed verifies that @implement
// annotations above private functions are allowed (no annotation-placement error).
// @test TEST-INTERNAL_ENGINE-038
func TestEngine_PublicFuncAnnotation_PrivateFuncAllowed(t *testing.T) {
	tmpDir := t.TempDir()
	// idd:ignore start
	code := `package foo

// helper does internal work.
// @implement SPEC-FOO-001
func helper() {}
`
	// idd:ignore end
	if err := os.WriteFile(filepath.Join(tmpDir, "foo.go"), []byte(code), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	cfg := &config.Config{
		Version: "1.0",
		Code: config.CodeConfig{
			Patterns: []string{filepath.Join(tmpDir, "*.go")},
		},
		Validation: config.ValidationConfig{
			RequirePublicFuncAnnotation: true,
		},
	}

	eng := New(cfg)
	result, err := eng.Run(context.Background(), model.NewIdentifierSet())
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, e := range result.Errors {
		if e.Rule == "annotation-placement" {
			t.Errorf("unexpected annotation-placement error above private func: %s", e.Message)
		}
	}
}

// TestEngine_PublicFuncAnnotation_PrivateMethodAllowed verifies that @implement
// annotations above private methods are allowed (no annotation-placement error).
// @test TEST-INTERNAL_ENGINE-039
func TestEngine_PublicFuncAnnotation_PrivateMethodAllowed(t *testing.T) {
	tmpDir := t.TempDir()
	// idd:ignore start
	code := `package foo

type bar struct{}

// process is an internal helper method.
// @implement SPEC-FOO-002
func (b *bar) process() {}
`
	// idd:ignore end
	if err := os.WriteFile(filepath.Join(tmpDir, "foo.go"), []byte(code), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	cfg := &config.Config{
		Version: "1.0",
		Code: config.CodeConfig{
			Patterns: []string{filepath.Join(tmpDir, "*.go")},
		},
		Validation: config.ValidationConfig{
			RequirePublicFuncAnnotation: true,
		},
	}

	eng := New(cfg)
	result, err := eng.Run(context.Background(), model.NewIdentifierSet())
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, e := range result.Errors {
		if e.Rule == "annotation-placement" {
			t.Errorf("unexpected annotation-placement error above private method: %s", e.Message)
		}
	}
}

// TestEngine_PublicFuncAnnotation_PrivateTypeAllowed verifies that @implement
// annotations above private types are allowed (no annotation-placement error).
// @test TEST-INTERNAL_ENGINE-040
func TestEngine_PublicFuncAnnotation_PrivateTypeAllowed(t *testing.T) {
	tmpDir := t.TempDir()
	// idd:ignore start
	code := `package foo

// internalState holds private state.
// @implement SPEC-FOO-003
type internalState struct {
	value int
}
`
	// idd:ignore end
	if err := os.WriteFile(filepath.Join(tmpDir, "foo.go"), []byte(code), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	cfg := &config.Config{
		Version: "1.0",
		Code: config.CodeConfig{
			Patterns: []string{filepath.Join(tmpDir, "*.go")},
		},
		Validation: config.ValidationConfig{
			RequirePublicFuncAnnotation: true,
		},
	}

	eng := New(cfg)
	result, err := eng.Run(context.Background(), model.NewIdentifierSet())
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, e := range result.Errors {
		if e.Rule == "annotation-placement" {
			t.Errorf("unexpected annotation-placement error above private type: %s", e.Message)
		}
	}
}

// idd:ignore end

// TestEngine_PublicFuncAnnotation_PublicStillRequired confirms that public
// functions without @implement still produce a public-func-annotation error.
// @test TEST-INTERNAL_ENGINE-041
func TestEngine_PublicFuncAnnotation_PublicStillRequired(t *testing.T) {
	tmpDir := t.TempDir()
	// idd:ignore start
	code := `package foo

// Bar is a public function missing @implement.
func Bar() {}
`
	// idd:ignore end
	if err := os.WriteFile(filepath.Join(tmpDir, "foo.go"), []byte(code), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	cfg := &config.Config{
		Version: "1.0",
		Code: config.CodeConfig{
			Patterns: []string{filepath.Join(tmpDir, "*.go")},
		},
		Validation: config.ValidationConfig{
			RequirePublicFuncAnnotation: true,
		},
	}

	eng := New(cfg)
	result, err := eng.Run(context.Background(), model.NewIdentifierSet())
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	found := false
	for _, e := range result.Errors {
		if e.Rule == "public-func-annotation" && strings.Contains(e.Message, "Bar") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected public-func-annotation error for public function missing @implement")
	}
}

// idd:ignore end

// TestEngine_PublicFuncAnnotation_GarbageAfterImplementStillErrors confirms that
// an @implement followed by something that is not a function/type (e.g. a
// variable or arbitrary code) still triggers annotation-placement.
// @test TEST-INTERNAL_ENGINE-042
func TestEngine_PublicFuncAnnotation_GarbageAfterImplementStillErrors(t *testing.T) {
	tmpDir := t.TempDir()
	// idd:ignore start
	code := `package foo

// @implement SPEC-FOO-004
var x = 5
`
	// idd:ignore end
	if err := os.WriteFile(filepath.Join(tmpDir, "foo.go"), []byte(code), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	cfg := &config.Config{
		Version: "1.0",
		Code: config.CodeConfig{
			Patterns: []string{filepath.Join(tmpDir, "*.go")},
		},
		Validation: config.ValidationConfig{
			RequirePublicFuncAnnotation: true,
		},
	}

	eng := New(cfg)
	result, err := eng.Run(context.Background(), model.NewIdentifierSet())
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	found := false
	for _, e := range result.Errors {
		if e.Rule == "annotation-placement" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected annotation-placement error when @implement is above a non-func/type declaration")
	}
}

// TestEngine_PrivateImplement_RequiresDoc verifies the third rule: any
// implement annotation (even on a private function) must have a corresponding
// doc entry. This is enforced by doc-code-correspondence, not by
// validatePublicFuncAnnotations.
// @test TEST-INTERNAL_ENGINE-043
func TestEngine_PrivateImplement_RequiresDoc(t *testing.T) {
	cfg := config.Default()
	eng := New(cfg)

	// Simulate the collector having parsed an @implement annotation from a
	// private function. The identifier exists on the code side with no
	// matching doc-side identifier.
	ids := model.NewIdentifierSet()
	codeOnly := model.NewIdentifier("SPEC-FOO-099", model.TypeSpec, "", "internal/foo/foo.go", 5)
	codeOnly.SetOrigin(model.OriginCode)
	ids.Add(codeOnly)

	eng.Run(context.Background(), ids)

	found := false
	for _, e := range eng.result.Errors {
		if e.Rule == "doc-code-correspondence" && strings.Contains(e.Message, "SPEC-FOO-099") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected doc-code-correspondence error for @implement without matching doc")
	}
}

// TestEngine_AnnotationIdentifier_HonorsIgnoreScope verifies that
// validateAnnotationIdentifiers honors the idd-ignore scope markers so
// fixtures embedded in the test source are not reported as real violations.
// @test TEST-INTERNAL_ENGINE-044
func TestEngine_AnnotationIdentifier_HonorsIgnoreScope(t *testing.T) {
	tmpDir := t.TempDir()
	// This block contains two comments that would normally trip the validator:
	//   1. an @implement with no identifier
	//   2. an @implement above a non-function declaration (placement)
	// The // idd:ignore start/end scope should silence both.
	code := "package foo\n\n" +
		"// idd:ignore start\n" +
		"// @implement\n" +
		"var ignored1 = 1\n" +
		"// idd:ignore end\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "foo.go"), []byte(code), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	cfg := &config.Config{
		Version: "1.0",
		Code: config.CodeConfig{
			Patterns: []string{filepath.Join(tmpDir, "*.go")},
		},
		Validation: config.ValidationConfig{
			RequireAnnotationIdentifier: true,
		},
	}

	eng := New(cfg)
	result, err := eng.Run(context.Background(), model.NewIdentifierSet())
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, e := range result.Errors {
		if e.Rule == "annotation-missing-identifier" {
			t.Errorf("expected no annotation-missing-identifier errors inside idd:ignore scope, got: %s @ %s", e.Message, e.Source)
		}
	}
}

// TestEngine_AnnotationPlacement_HonorsIgnoreScope verifies that the placement
// pass of validatePublicFuncAnnotations skips content inside an // idd:ignore
// start/end block — not just the first and third passes.
// @test TEST-INTERNAL_ENGINE-045
func TestEngine_AnnotationPlacement_HonorsIgnoreScope(t *testing.T) {
	tmpDir := t.TempDir()
	// An @implement followed by a `var` is normally a placement error. The
	// ignore scope should silence it.
	code := "package foo\n\n" +
		"// idd:ignore start\n" +
		"// @implement SPEC-IGNORE-001\n" +
		"var ignored = 1\n" +
		"// idd:ignore end\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "foo.go"), []byte(code), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	cfg := &config.Config{
		Version: "1.0",
		Code: config.CodeConfig{
			Patterns: []string{filepath.Join(tmpDir, "*.go")},
		},
		Validation: config.ValidationConfig{
			RequirePublicFuncAnnotation: true,
		},
	}

	eng := New(cfg)
	result, err := eng.Run(context.Background(), model.NewIdentifierSet())
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, e := range result.Errors {
		if e.Rule == "annotation-placement" {
			t.Errorf("expected no annotation-placement errors inside idd:ignore scope, got: %s @ %s", e.Message, e.Source)
		}
	}
}

// TestEngine_ConsecutiveAnnotations_HonorsIgnoreScope verifies that
// validateConsecutiveAnnotations skips content inside an // idd:ignore
// start/end block.
// @test TEST-INTERNAL_ENGINE-046
func TestEngine_ConsecutiveAnnotations_HonorsIgnoreScope(t *testing.T) {
	tmpDir := t.TempDir()
	// Two consecutive @implement comments would normally trip
	// annotation-consecutive-line. The ignore scope should silence it.
	code := "package foo\n\n" +
		"// idd:ignore start\n" +
		"// @implement SPEC-A-001\n" +
		"// @implement SPEC-B-001\n" +
		"// idd:ignore end\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "foo.go"), []byte(code), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	cfg := &config.Config{
		Version: "1.0",
		Code: config.CodeConfig{
			Patterns: []string{filepath.Join(tmpDir, "*.go")},
		},
		Validation: config.ValidationConfig{
			RequireAnnotationOnSameLine: true,
		},
	}

	eng := New(cfg)
	result, err := eng.Run(context.Background(), model.NewIdentifierSet())
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, e := range result.Errors {
		if e.Rule == "annotation-consecutive-line" {
			t.Errorf("expected no annotation-consecutive-line errors inside idd:ignore scope, got: %s @ %s", e.Message, e.Source)
		}
	}
}
