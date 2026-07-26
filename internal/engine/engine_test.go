// Package engine provides testing utilities for the engine module.

// Spec: docs/internal/engine/spec.md
// Test: docs/internal/engine/testing.md

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

// @test-contract TEST-INTERNAL_ENGINE-007
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

// @test-contract TEST-INTERNAL_ENGINE-007
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

// @test-contract TEST-INTERNAL_ENGINE-007
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

// @test-contract TEST-INTERNAL_ENGINE-007
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

// @test-contract TEST-INTERNAL_ENGINE-007
func TestEngine_ConsistencyCheck_FunctionLocatorOnly(t *testing.T) {
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
	docID := model.NewIdentifierWithDescribe(
		"SPEC-INTERNAL_SAMPLE-001",
		model.TypeSpec,
		"Sample behavior",
		"validate self-describing package document relationships",
		"docs/internal/sample/spec.md",
		1,
	)
	ids.Add(docID)
	codeID := model.NewIdentifierWithDescribe(
		"SPEC-INTERNAL_SAMPLE-001",
		model.TypeSpec,
		"",
		"[function: UnrelatedName]",
		"internal/sample/sample.go",
		10,
	)
	codeID.SetOrigin(model.OriginCode)
	ids.Add(codeID)

	result, err := New(cfg).Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	for _, warning := range result.Warnings {
		if warning.Rule == "consistency-check" {
			t.Fatalf("Run() warning = %v, want function locator ignored", warning)
		}
	}
}

// @test-contract TEST-INTERNAL_ENGINE-007
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

// @test-contract TEST-INTERNAL_ENGINE-007
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

// @test-contract TEST-INTERNAL_ENGINE-007
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

// @test TEST-INTERNAL_ENGINE-004
func TestEngine_ValidateCompleteness_UsesDocSourceForMissingDocFields(t *testing.T) {
	tests := []struct {
		name       string
		docID      *model.Identifier
		codeID     *model.Identifier
		wantRule   string
		wantSource string
	}{
		{
			name:       "test missing spec coverage",
			docID:      model.NewIdentifier("TEST-BE-001", model.TypeTest, "Kanban", "docs/backend/testing.md", 12),
			codeID:     model.NewIdentifier("TEST-BE-001", model.TypeTest, "", "internal/service/flow_workspace_test.go", 34),
			wantRule:   "test-missing-coverage",
			wantSource: "docs/backend/testing.md",
		},
		{
			name:       "spec missing tests",
			docID:      model.NewIdentifier("SPEC-BE-001", model.TypeSpec, "Kanban", "docs/backend/spec.md", 12),
			codeID:     model.NewIdentifier("SPEC-BE-001", model.TypeSpec, "", "internal/service/flow_workspace.go", 34),
			wantRule:   "spec-missing-tests",
			wantSource: "docs/backend/spec.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Version: "1.0",
				Validation: config.ValidationConfig{
					RequireSpecTestCoverage: true,
				},
			}
			eng := New(cfg)
			ids := model.NewIdentifierSet()

			tt.codeID.SetOrigin(model.OriginCode)
			ids.Add(tt.codeID)
			tt.docID.SetOrigin(model.OriginDoc)
			ids.Add(tt.docID)

			result, err := eng.Run(context.Background(), ids)
			if err != nil {
				t.Fatalf("Run failed: %v", err)
			}

			for _, err := range result.Errors {
				if err.Rule != tt.wantRule {
					continue
				}
				if err.Source != tt.wantSource {
					t.Fatalf("Source = %q, want %q", err.Source, tt.wantSource)
				}
				return
			}
			t.Fatalf("Expected %s error, got %v", tt.wantRule, result.Errors)
		})
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

// @test-contract TEST-INTERNAL_ENGINE-007
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

// @test TEST-INTERNAL_ENGINE-003
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
