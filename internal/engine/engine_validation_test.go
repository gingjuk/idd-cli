// Package engine tests document, package, and annotation validation rules.

// Spec: docs/internal/engine/spec.md
// Test: docs/internal/engine/testing.md
package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jingxu9x/idd-cli/internal/config"
	"github.com/jingxu9x/idd-cli/internal/model"
)

// @test TEST-INTERNAL_ENGINE-003
func TestEngine_validateDesignSections_SelfDescribingContent(t *testing.T) {
	complete := `# Design

## Architecture
Concrete architecture.

## Package Layout
Concrete package layout.

## Function Composition
Concrete call flow.

## Dependencies
No runtime dependencies.

## Testability Hooks
Temporary filesystem fixtures.
`
	tests := []struct {
		name      string
		content   string
		wantError bool
	}{
		{
			name:    "all self-describing design sections have content",
			content: complete,
		},
		{
			name:      "empty self-describing design section is incomplete",
			content:   strings.Replace(complete, "## Architecture\nConcrete architecture.", "## Architecture", 1),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			docsDir := filepath.Join(t.TempDir(), "docs", "internal", "example")
			if err := os.MkdirAll(docsDir, 0o755); err != nil {
				t.Fatalf("MkdirAll() error = %v", err)
			}
			content := `---
idd:
  version: "1.0"
  package: internal/example
  document: design
  components: []
---

` + tt.content
			if err := os.WriteFile(filepath.Join(docsDir, "design.md"), []byte(content), 0o644); err != nil {
				t.Fatalf("WriteFile(design.md) error = %v", err)
			}

			cfg := config.Default()
			cfg.Docs.Patterns = []string{filepath.Join(docsDir, "**/*.md")}
			eng := New(cfg)
			eng.validateDesignSections()

			hasError := false
			for _, validationErr := range eng.result.Errors {
				if validationErr.Rule == "design-sections" {
					hasError = true
				}
			}
			if hasError != tt.wantError {
				t.Errorf("design-sections error = %v, want %v: %v", hasError, tt.wantError, eng.result.Errors)
			}
		})
	}
}

// @test TEST-INTERNAL_ENGINE-003
func TestEngine_validateDocumentTestKinds(t *testing.T) {
	tests := []struct {
		name      string
		docKind   string
		codeKinds []string
		wantError bool
	}{
		{name: "behavior kind uses test annotation", docKind: "test", codeKinds: []string{"test"}},
		{name: "contract kind uses contract annotation", docKind: "contract", codeKinds: []string{"contract"}},
		{name: "behavior kind rejects contract annotation", docKind: "test", codeKinds: []string{"contract"}, wantError: true},
		{name: "contract kind rejects test annotation", docKind: "contract", codeKinds: []string{"test"}, wantError: true},
		{name: "document kind is exclusive", docKind: "test", codeKinds: []string{"test", "contract"}, wantError: true},
		{name: "legacy test has no kind contract", codeKinds: []string{"contract"}},
		{name: "missing code remains doc code concern", docKind: "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ids := model.NewIdentifierSet()
			documented := model.NewIdentifier("TEST-INTERNAL_SAMPLE-001", model.TypeTest, "Sample test", "docs/internal/sample/testing.md", 8)
			documented.Kind = tt.docKind
			ids.Add(documented)
			for index, kind := range tt.codeKinds {
				code := model.NewIdentifier(
					"TEST-INTERNAL_SAMPLE-001",
					model.TypeTest,
					"",
					fmt.Sprintf("internal/sample/sample_%d_test.go", index),
					10,
				)
				code.SetOrigin(model.OriginCode)
				code.Kind = kind
				ids.Add(code)
			}

			eng := New(config.Default())
			eng.buildGraph(ids)
			eng.validateDocumentTestKinds()

			hasError := false
			for _, validationErr := range eng.result.Errors {
				if validationErr.Rule == "idd-document-test-kind" {
					hasError = true
				}
			}
			if hasError != tt.wantError {
				t.Errorf("idd-document-test-kind error = %v, want %v: %v", hasError, tt.wantError, eng.result.Errors)
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
	tests := []struct {
		name                  string
		spec                  string
		rootDoc               bool
		selfDescribingSibling bool
	}{
		{
			name: "legacy frontmatter declares related files",
			spec: `---
related_files:
  spec: spec.md
  contract: contract.md
---
# Spec
Content
`,
		},
		{
			name: "self-describing document does not repeat related files",
			spec: `---
idd:
  version: "1.0"
  package: test
  document: spec
  specs: []
---
# Specifications
`,
		},
		{
			name:                  "one self-describing sibling selects mode for the package",
			spec:                  "# Specifications\n",
			selfDescribingSibling: true,
		},
		{
			name:    "root documentation is not a package document",
			spec:    "# IDD introduction\n",
			rootDoc: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			docsDir := filepath.Join(t.TempDir(), "docs", "test")
			filename := "spec.md"
			if tt.rootDoc {
				docsDir = filepath.Dir(docsDir)
				filename = "idd-intro.md"
			}
			if err := os.MkdirAll(docsDir, 0o755); err != nil {
				t.Fatalf("MkdirAll() error = %v", err)
			}
			if err := os.WriteFile(filepath.Join(docsDir, filename), []byte(tt.spec), 0o644); err != nil {
				t.Fatalf("WriteFile(%s) error = %v", filename, err)
			}
			if tt.selfDescribingSibling {
				contract := "---\nidd: {version: \"1.0\", package: test, document: contract, contracts: []}\n---\n# Contracts\n"
				if err := os.WriteFile(filepath.Join(docsDir, "contract.md"), []byte(contract), 0o644); err != nil {
					t.Fatalf("WriteFile(contract.md) error = %v", err)
				}
			}
			cfg := config.Default()
			cfg.Validation.AllowOrphans = true
			cfg.Docs.Patterns = []string{filepath.Join(docsDir, "**/*.md")}
			eng := New(cfg)
			eng.validateRelatedFiles()

			if len(eng.result.Errors) != 0 {
				t.Errorf("validateRelatedFiles() errors = %v, want none", eng.result.Errors)
			}
		})
	}
}

// @test TEST-INTERNAL_ENGINE-029
func TestHasIDDDocumentMetadata(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  bool
	}{
		{
			name:  "top-level idd block",
			lines: []string{"---", "idd:", `  version: "1.0"`, "---"},
			want:  true,
		},
		{
			name:  "flow-style top-level idd block",
			lines: []string{"---", `idd: {version: "1.0"}`, "---"},
			want:  true,
		},
		{
			name:  "nested idd key is generic metadata",
			lines: []string{"---", "site:", "  idd:", "    enabled: true", "---"},
		},
		{
			name:  "idd text outside frontmatter",
			lines: []string{"# Notes", "idd:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasIDDDocumentMetadata(tt.lines); got != tt.want {
				t.Errorf("hasIDDDocumentMetadata() = %v, want %v", got, tt.want)
			}
		})
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
// the syntax-tree public declaration rule.
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

// TestEngine_AnnotationIdentifier_HonorsIgnoreScope verifies that syntax-tree
// annotation validation honors idd-ignore scope markers so fixtures embedded
// in the test source are not reported as real violations.
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

// TestEngine_AnnotationPlacement_HonorsIgnoreScope verifies that syntax-tree
// placement validation skips content inside an // idd:ignore start/end block.
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
// syntax-tree consecutive-annotation validation skips content inside an
// // idd:ignore start/end block.
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
