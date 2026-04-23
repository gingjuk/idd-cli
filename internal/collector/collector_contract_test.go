// Package collector provides testing utilities for the collector module.

// Spec: docs/internal/collector/spec.md
// Test: docs/internal/collector/testing.md
// Contract: docs/internal/collector/contract.md
package collector

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jingxu9x/idd-cli/internal/config"
	"github.com/jingxu9x/idd-cli/internal/model"
)

// =============================================================================
// CONTRACT-COLLECTOR-001: Collector Interface Contracts
// =============================================================================

// -----------------------------------------------------------------------------
// Frontmatter Parsing Contract Tests
// -----------------------------------------------------------------------------

func TestContract_Frontmatter_ValidYAML(t *testing.T) {
	content := `---
markers:
  - id: SPEC-BE-001
    name: Test Specification
    describe: A test spec
---
# SPEC-BE-001

This document describes SPEC-BE-001.
`

	fm, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("ParseFrontmatter returned error: %v", err)
	}
	if fm == nil {
		t.Fatal("ParseFrontmatter returned nil frontmatter")
	}
	if len(fm.Markers) != 1 {
		t.Fatalf("len(fm.Markers) = %d, want 1", len(fm.Markers))
	}
	if fm.Markers[0].ID != "SPEC-BE-001" {
		t.Errorf("fm.Markers[0].ID = %q, want %q", fm.Markers[0].ID, "SPEC-BE-001")
	}
	if fm.Markers[0].Name != "Test Specification" {
		t.Errorf("fm.Markers[0].Name = %q, want %q", fm.Markers[0].Name, "Test Specification")
	}
	if fm.Markers[0].Describe != "A test spec" {
		t.Errorf("fm.Markers[0].Describe = %q, want %q", fm.Markers[0].Describe, "A test spec")
	}
}

func TestContract_Frontmatter_NoFrontmatter(t *testing.T) {
	content := `# Document

This document has no frontmatter.
`

	fm, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("ParseFrontmatter returned error: %v", err)
	}
	if fm != nil {
		t.Errorf("ParseFrontmatter returned non-nil frontmatter for content without frontmatter")
	}
}

func TestContract_Frontmatter_EmptyMarkers(t *testing.T) {
	content := `---
markers: []
---

# Document
`

	fm, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("ParseFrontmatter returned error: %v", err)
	}
	if fm == nil {
		t.Fatal("ParseFrontmatter returned nil frontmatter")
	}
	if len(fm.Markers) != 0 {
		t.Errorf("len(fm.Markers) = %d, want 0", len(fm.Markers))
	}
}

func TestContract_Frontmatter_MultipleMarkers(t *testing.T) {
	content := `---
markers:
  - id: SPEC-BE-001
    name: First Spec
    describe: Description one
  - id: TEST-BE-001
    name: Test Case
    describe: Test description
  - id: CONTRACT-BE-001
    name: Contract
    describe: Contract description
---

# SPEC-BE-001

Content referencing TEST-BE-001.
`

	fm, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("ParseFrontmatter returned error: %v", err)
	}
	if fm == nil {
		t.Fatal("ParseFrontmatter returned nil frontmatter")
	}
	if len(fm.Markers) != 3 {
		t.Fatalf("len(fm.Markers) = %d, want 3", len(fm.Markers))
	}
}

func TestContract_Frontmatter_InvalidYAML(t *testing.T) {
	content := `---
markers:
  - id: SPEC-BE-001
    name: Test
    describe: Missing colon
  invalid yaml here
---

# Document
`

	_, err := ParseFrontmatter(content)
	if err == nil {
		t.Error("ParseFrontmatter expected error for invalid YAML, got nil")
	}
}

func TestContract_ValidateMarkers_ValidMarkers(t *testing.T) {
	content := `---
markers:
  - id: SPEC-BE-001
    name: Test Spec
    describe: A test
---

# SPEC-BE-001: Test Spec Description

This document describes ` + "`SPEC-BE-001`" + `.
`

	fm, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("ParseFrontmatter returned error: %v", err)
	}

	errors := ValidateFrontmatterMarkers(fm, content, "test.md")
	if len(errors) != 0 {
		t.Errorf("ValidateFrontmatterMarkers returned %d errors, want 0: %v", len(errors), errors)
	}
}

func TestContract_ValidateMarkers_MarkerNotInContent(t *testing.T) {
	content := `---
markers:
  - id: SPEC-BE-001
    name: Test Spec
    describe: A test
---

# Introduction

This document does NOT contain SPEC-BE-001 in content.
`

	fm, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("ParseFrontmatter returned error: %v", err)
	}

	errors := ValidateFrontmatterMarkers(fm, content, "test.md")
	if len(errors) == 0 {
		t.Error("ValidateFrontmatterMarkers should return error for marker not in content")
	}
}

func TestContract_ValidateMarkers_BareMarkerNoBackticks(t *testing.T) {
	content := `---
markers:
  - id: SPEC-BE-001
    name: Test Spec
    describe: A test
---

# SPEC-BE-001

The identifier SPEC-BE-001 should be wrapped in backticks.
`

	fm, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("ParseFrontmatter returned error: %v", err)
	}

	errors := ValidateFrontmatterMarkers(fm, content, "test.md")
	if len(errors) == 0 {
		t.Error("ValidateFrontmatterMarkers should return error for bare marker without backticks")
	}
}

func TestContract_ValidateMarkers_NilFrontmatter(t *testing.T) {
	errors := ValidateFrontmatterMarkers(nil, "some content", "test.md")
	if len(errors) != 0 {
		t.Errorf("ValidateFrontmatterMarkers should return no errors for nil frontmatter, got %d", len(errors))
	}
}

// -----------------------------------------------------------------------------
// DocCollector Interface Contract Tests
// -----------------------------------------------------------------------------

func TestContract_NewDocCollector_ReturnsNonNil(t *testing.T) {
	cfg := config.Default()
	collector := NewDocCollector(cfg)
	if collector == nil {
		t.Fatal("NewDocCollector returned nil")
	}
}

func TestContract_NewDocCollector_SetsConfig(t *testing.T) {
	cfg := config.Default()
	collector := NewDocCollector(cfg)
	if collector.cfg != cfg {
		t.Error("NewDocCollector did not set config")
	}
}

func TestContract_DocCollector_Collect_Signature(t *testing.T) {
	cfg := config.Default()
	collector := NewDocCollector(cfg)

	var set *model.IdentifierSet
	var errors []*model.ValidationError
	var err error

	_ = set
	_ = errors
	_ = err
	_ = collector.Collect
}

func TestContract_DocCollector_Collect_FileNotFound(t *testing.T) {
	cfg := config.Default()
	collector := NewDocCollector(cfg)

	set, errors, err := collector.Collect(context.Background(), "/nonexistent/path/to/file.md")
	if err != nil {
		t.Fatalf("Collect should not return error for nonexistent file, got: %v", err)
	}
	if set == nil {
		t.Error("Collect returned nil IdentifierSet")
	}
	if set.Count() != 0 {
		t.Errorf("Collect should return empty set for nonexistent file, got count %d", set.Count())
	}
	if errors != nil {
		t.Error("Collect should return nil errors for nonexistent file")
	}
}

func TestContract_DocCollector_Collect_ValidMarkdownFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "spec.md")
	content := `---
markers:
  - id: SPEC-BE-001
    name: Test Specification
    describe: A test spec
---

# SPEC-BE-001

This document describes ` + "`SPEC-BE-001`" + `.
`

	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := config.Default()
	collector := NewDocCollector(cfg)

	set, _, err := collector.Collect(context.Background(), tmpFile)
	if err != nil {
		t.Fatalf("Collect returned error: %v", err)
	}
	if set == nil {
		t.Fatal("Collect returned nil IdentifierSet")
	}
	if set.Count() == 0 {
		t.Error("Collect should return at least one identifier")
	}
}

func TestContract_DocCollector_Collect_Directory(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "docs")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	content := `---
markers:
  - id: SPEC-BE-001
    name: Test
    describe: Test desc
---

# SPEC-BE-001

Content for ` + "`SPEC-BE-001`" + `.
`

	if err := os.WriteFile(filepath.Join(subDir, "spec.md"), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := config.Default()
	collector := NewDocCollector(cfg)

	set, _, err := collector.Collect(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("Collect returned error: %v", err)
	}
	if set == nil {
		t.Fatal("Collect returned nil IdentifierSet")
	}
	if set.Count() == 0 {
		t.Error("Collect should find identifiers in directory")
	}
}

func TestContract_DocCollector_Collect_IgnoresNonMarkdownFiles(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.go")

	// idd:ignore start
	code := `// @implement SPEC-BE-001
package main
`
	// idd:ignore end
	if err := os.WriteFile(tmpFile, []byte(code), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := config.Default()
	collector := NewDocCollector(cfg)

	set, _, err := collector.Collect(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("Collect returned error: %v", err)
	}
	if set.Has("SPEC-BE-001") {
		t.Error("DocCollector should not collect code annotations from .go files")
	}
}

func TestContract_DocCollector_Collect_ExtractsReferencesInContent(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "spec.md")

	content := "---\nmarkers:\n  - id: SPEC-BE-001\n    name: Test Spec\n    describe: Test description\n---\n\n# SPEC-BE-001\n\nThis spec references `TEST-BE-001` and `CONTRACT-BE-001` in backticks.\n"

	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := config.Default()
	collector := NewDocCollector(cfg)

	set, _, err := collector.Collect(context.Background(), tmpFile)
	if err != nil {
		t.Fatalf("Collect returned error: %v", err)
	}

	if !set.Has("SPEC-BE-001") {
		t.Error("Collect should find SPEC-BE-001 from frontmatter")
	}

	if !set.Has("TEST-BE-001") {
		t.Error("Collect should find TEST-BE-001 referenced in content")
	}
	if !set.Has("CONTRACT-BE-001") {
		t.Error("Collect should find CONTRACT-BE-001 referenced in content")
	}
}

// -----------------------------------------------------------------------------
// CodeCollector Interface Contract Tests
// -----------------------------------------------------------------------------

func TestContract_NewCodeCollector_ReturnsNonNil(t *testing.T) {
	cfg := config.Default()
	collector := NewCodeCollector(cfg)
	if collector == nil {
		t.Fatal("NewCodeCollector returned nil")
	}
}

func TestContract_NewCodeCollector_SetsConfig(t *testing.T) {
	cfg := config.Default()
	collector := NewCodeCollector(cfg)
	if collector.cfg != cfg {
		t.Error("NewCodeCollector did not set config")
	}
}

func TestContract_CodeCollector_Collect_Signature(t *testing.T) {
	cfg := config.Default()
	collector := NewCodeCollector(cfg)

	var set *model.IdentifierSet
	var err error

	_ = set
	_ = err
	_ = collector.Collect
}

func TestContract_CodeCollector_Collect_FileNotFound(t *testing.T) {
	cfg := config.Default()
	collector := NewCodeCollector(cfg)

	set, err := collector.Collect(context.Background(), "/nonexistent/path/to/file.go")
	if err != nil {
		t.Fatalf("Collect should not return error for nonexistent file, got: %v", err)
	}
	if set == nil {
		t.Error("Collect returned nil IdentifierSet")
	}
	if set.Count() != 0 {
		t.Errorf("Collect should return empty set for nonexistent file, got count %d", set.Count())
	}
}

func TestContract_CodeCollector_Collect_GoSourceFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.go")

	// idd:ignore start
	code := `package main

// @implement SPEC-BE-001
func Authenticate() {
}
`
	// idd:ignore end
	if err := os.WriteFile(tmpFile, []byte(code), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := config.Default()
	collector := NewCodeCollector(cfg)

	set, err := collector.Collect(context.Background(), tmpFile)
	if err != nil {
		t.Fatalf("Collect returned error: %v", err)
	}
	if set == nil {
		t.Fatal("Collect returned nil IdentifierSet")
	}
	if !set.Has("SPEC-BE-001") {
		t.Error("Collect should find SPEC-BE-001 annotation")
	}
}

func TestContract_CodeCollector_Collect_TypeScriptFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.ts")

	// idd:ignore start
	code := `// @test-contract TEST-BE-001
function test() {
}
`
	// idd:ignore end
	if err := os.WriteFile(tmpFile, []byte(code), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := config.Default()
	collector := NewCodeCollector(cfg)

	set, err := collector.Collect(context.Background(), tmpFile)
	if err != nil {
		t.Fatalf("Collect returned error: %v", err)
	}
	if set == nil {
		t.Fatal("Collect returned nil IdentifierSet")
	}
	if !set.Has("TEST-BE-001") {
		t.Error("Collect should find TEST-BE-001 annotation in .ts file")
	}
}

func TestContract_CodeCollector_Collect_JavaScriptFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.js")

	// idd:ignore start
	code := `// @implement SPEC-BE-001
function implement() {
}
`
	// idd:ignore end
	if err := os.WriteFile(tmpFile, []byte(code), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := config.Default()
	collector := NewCodeCollector(cfg)

	set, err := collector.Collect(context.Background(), tmpFile)
	if err != nil {
		t.Fatalf("Collect returned error: %v", err)
	}
	if set == nil {
		t.Fatal("Collect returned nil IdentifierSet")
	}
	if !set.Has("SPEC-BE-001") {
		t.Error("Collect should find SPEC-BE-001 annotation in .js file")
	}
}

func TestContract_CodeCollector_Collect_Directory(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "src")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// idd:ignore start
	code := `// @implement SPEC-BE-001
package main
`
	// idd:ignore end
	if err := os.WriteFile(filepath.Join(subDir, "main.go"), []byte(code), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := config.Default()
	collector := NewCodeCollector(cfg)

	set, err := collector.Collect(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("Collect returned error: %v", err)
	}
	if set == nil {
		t.Fatal("Collect returned nil IdentifierSet")
	}
	if !set.Has("SPEC-BE-001") {
		t.Error("Collect should find SPEC-BE-001 in directory")
	}
}

func TestContract_CodeCollector_Collect_IgnoresMarkdownFiles(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "readme.md")

	content := `---
markers:
  - id: SPEC-BE-001
    name: Test
    describe: Test desc
---

# SPEC-BE-001

Content.
`

	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := config.Default()
	collector := NewCodeCollector(cfg)

	set, err := collector.Collect(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("Collect returned error: %v", err)
	}
	if set.Has("SPEC-BE-001") {
		t.Error("CodeCollector should not collect from .md files")
	}
}

func TestContract_CodeCollector_Collect_MultipleAnnotations(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.go")

	// idd:ignore start
	code := `package main

// @implement SPEC-TEST-001, SPEC-TEST-002
func Authenticate() {
}

// @test-contract TEST-TEST-001
func TestAuthenticate() {
}
`
	// idd:ignore end
	if err := os.WriteFile(tmpFile, []byte(code), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := config.Default()
	collector := NewCodeCollector(cfg)

	set, err := collector.Collect(context.Background(), tmpFile)
	if err != nil {
		t.Fatalf("Collect returned error: %v", err)
	}

	if set.Count() != 3 {
		t.Errorf("Expected 3 identifiers, got %d", set.Count())
	}
}

func TestContract_CodeCollector_Collect_OriginSet(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.go")

	// idd:ignore start
	code := `// @implement SPEC-BE-001
package main
`
	// idd:ignore end
	if err := os.WriteFile(tmpFile, []byte(code), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := config.Default()
	collector := NewCodeCollector(cfg)

	set, err := collector.Collect(context.Background(), tmpFile)
	if err != nil {
		t.Fatalf("Collect returned error: %v", err)
	}

	if !set.HasOrigin("SPEC-BE-001", model.OriginCode) {
		t.Error("Identifier from CodeCollector should have OriginCode")
	}
}

// -----------------------------------------------------------------------------
// Error Handling Rules Contract Tests
// -----------------------------------------------------------------------------

func TestContract_ErrorHandling_FileNotFound_ReturnsEmptySet(t *testing.T) {
	cfg := config.Default()

	docCollector := NewDocCollector(cfg)
	set, errors, err := docCollector.Collect(context.Background(), "/nonexistent/file.md")
	if err != nil {
		t.Errorf("DocCollector Collect: error = %v, want nil", err)
	}
	if set.Count() != 0 {
		t.Errorf("DocCollector Collect: count = %d, want 0", set.Count())
	}
	if errors != nil {
		t.Errorf("DocCollector Collect: errors = %v, want nil", errors)
	}

	codeCollector := NewCodeCollector(cfg)
	codeSet, codeErr := codeCollector.Collect(context.Background(), "/nonexistent/file.go")
	if codeErr != nil {
		t.Errorf("CodeCollector Collect: error = %v, want nil", codeErr)
	}
	if codeSet.Count() != 0 {
		t.Errorf("CodeCollector Collect: count = %d, want 0", codeSet.Count())
	}
}

func TestContract_ErrorHandling_InvalidFrontmatter_AddsValidationError(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "spec.md")

	content := `---
markers:
  - id: SPEC-BE-001
    name: Test
    invalid yaml here
---

# SPEC-BE-001
`

	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := config.Default()
	collector := NewDocCollector(cfg)

	_, errors, err := collector.Collect(context.Background(), tmpFile)
	if err != nil {
		t.Fatalf("Collect returned unexpected error: %v", err)
	}

	if len(errors) == 0 {
		t.Error("Collect should add ValidationError for invalid frontmatter")
	}
}

func TestContract_ErrorHandling_BareMarker_AddsValidationError(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "spec.md")

	content := `---
markers:
  - id: SPEC-BE-001
    name: Test
    describe: Test
---

# SPEC-BE-001

SPEC-BE-001 should be wrapped in backticks here.
`

	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := config.Default()
	collector := NewDocCollector(cfg)

	_, errors, _ := collector.Collect(context.Background(), tmpFile)

	found := false
	for _, e := range errors {
		if e.Rule == "frontmatter-mismatch" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Collect should add frontmatter-mismatch error for bare marker")
	}
}

func TestContract_ErrorHandling_ModulePrefixMismatch_AddsValidationError(t *testing.T) {
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs", "BE")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	tmpFile := filepath.Join(docsDir, "spec.md")

	content := `---
markers:
  - id: SPEC-FE-001
    name: Frontend Spec
    describe: A frontend spec
---

# SPEC-FE-001

Content.
`

	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	cfg := config.Default()
	collector := NewDocCollector(cfg)

	_, errors, _ := collector.Collect(context.Background(), tmpFile)

	found := false
	for _, e := range errors {
		if e.Rule == "module-prefix-mismatch" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Collect should add module-prefix-mismatch error")
	}
}
