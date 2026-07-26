// Package collector provides testing utilities for the collector module.
package collector

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jingxu9x/idd-cli/internal/config"
	"github.com/jingxu9x/idd-cli/internal/model"
)

// @test-contract TEST-INTERNAL_COLLECTOR-024
func TestCodeCollector_Collect(t *testing.T) {
	tmpDir := t.TempDir()

	// idd:ignore start
	code := `package main

// @implement SPEC-CMD_IDD_CLI-001
func Authenticate() error {
	return nil
}

// @implement SPEC-CMD_IDD_CLI-002
type Payment interface {
	Process(amount float64) error
}

// @test TEST-CMD_IDD_CLI-001
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

	if !set.Has("SPEC-CMD_IDD_CLI-001") {
		t.Error("Should have SPEC-CMD_IDD_CLI-001")
	}
	if !set.Has("SPEC-CMD_IDD_CLI-002") {
		t.Error("Should have SPEC-CMD_IDD_CLI-002")
	}
	if !set.Has("TEST-CMD_IDD_CLI-001") {
		t.Error("Should have TEST-CMD_IDD_CLI-001")
	}

	specID, ok := set.Get("SPEC-CMD_IDD_CLI-001")
	if !ok {
		t.Fatal("SPEC-CMD_IDD_CLI-001 not found")
	}
	if specID.Describe == "" {
		t.Error("SPEC-CMD_IDD_CLI-001 should have Describe extracted from function comment")
	}
	if specID.Origin != model.OriginCode {
		t.Errorf("Origin = %v, want %v", specID.Origin, model.OriginCode)
	}
}

// @test TEST-INTERNAL_COLLECTOR-025
func TestCodeCollector_CollectGoFile(t *testing.T) {
	tmpDir := t.TempDir()

	// idd:ignore start
	code := `package main

// @implement SPEC-CMD_IDD_CLI-002
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

	if !set.Has("SPEC-CMD_IDD_CLI-002") {
		t.Error("Should have SPEC-CMD_IDD_CLI-002")
	}
}

// @test TEST-INTERNAL_COLLECTOR-026
func TestCodeCollector_MultipleAnnotations(t *testing.T) {
	tmpDir := t.TempDir()

	// idd:ignore start
	code := `package main

// @implement SPEC-CMD_IDD_CLI-001, SPEC-CMD_IDD_CLI-002
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

	if !set.Has("SPEC-CMD_IDD_CLI-001") {
		t.Error("Should have SPEC-CMD_IDD_CLI-001")
	}
	if !set.Has("SPEC-CMD_IDD_CLI-002") {
		t.Error("Should have SPEC-CMD_IDD_CLI-002")
	}
}

// @test TEST-INTERNAL_COLLECTOR-026
func TestCodeCollector_TestAnnotationKind(t *testing.T) {
	tests := []struct {
		name       string
		annotation string
		wantKind   string
	}{
		{name: "behavior test", annotation: "@test", wantKind: "test"},
		{name: "contract test", annotation: "@test-contract", wantKind: "contract"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// idd:ignore start
			code := "package sample\n\n// " + tt.annotation + " TEST-SAMPLE-001\nfunc TestBehavior() {}\n"
			// idd:ignore end
			if err := os.WriteFile(filepath.Join(tmpDir, "sample_test.go"), []byte(code), 0o644); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}

			set, err := NewCodeCollector(config.Default()).Collect(nil, tmpDir)
			if err != nil {
				t.Fatalf("Collect() error = %v", err)
			}
			identifier, ok := set.Get("TEST-SAMPLE-001")
			if !ok {
				t.Fatal("TEST-SAMPLE-001 not collected")
			}
			if identifier.Kind != tt.wantKind {
				t.Errorf("Kind = %q, want %q", identifier.Kind, tt.wantKind)
			}
		})
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-024
func TestCodeCollector_CollectsAttachedAnnotationsAcrossLanguages(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		source string
	}{
		{"go", "service.go", "package service\n// @implement SPEC-LANGUAGE-001\nfunc Run() {}\n"},
		{"typescript", "service.ts", "// @implement SPEC-LANGUAGE-001\nexport function run() {}\n"},
		{"tsx", "view.tsx", "// @implement SPEC-LANGUAGE-001\nexport function View() { return <div/> }\n"},
		{"javascript", "service.js", "// @implement SPEC-LANGUAGE-001\nexport function run() {}\n"},
		{"cpp", "service.cpp", "// @implement SPEC-LANGUAGE-001\nint run() { return 0; }\n"},
		{"java", "Service.java", "// @implement SPEC-LANGUAGE-001\npublic class Service {}\n"},
		{"python", "service.py", "# @implement SPEC-LANGUAGE-001\ndef run():\n    pass\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), test.path)
			if err := os.WriteFile(path, []byte(test.source), 0o644); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}
			set, findings, err := NewCodeCollector(config.Default()).CollectWithErrors(nil, path)
			if err != nil {
				t.Fatalf("CollectWithErrors() error = %v", err)
			}
			if len(findings) != 0 {
				t.Fatalf("findings = %#v, want none", findings)
			}
			if !set.HasOrigin("SPEC-LANGUAGE-001", model.OriginCode) {
				t.Errorf("identifiers = %#v, want SPEC-LANGUAGE-001", set.AllIdentifiers())
			}
		})
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-024
func TestCodeCollector_ReportsParseErrorsWithoutRegexFallback(t *testing.T) {
	path := filepath.Join(t.TempDir(), "broken.py")
	source := "# @implement SPEC-LANGUAGE-001\ndef broken(:\n"
	if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	set, findings, err := NewCodeCollector(config.Default()).CollectWithErrors(nil, path)
	if err != nil {
		t.Fatalf("CollectWithErrors() error = %v", err)
	}
	if set.Count() != 0 {
		t.Errorf("Identifier count = %d, want zero", set.Count())
	}
	if len(findings) != 1 || findings[0].Rule != "source-parse" {
		t.Fatalf("findings = %#v, want one source-parse finding", findings)
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-024
func TestCodeCollector_ConfiguredUnsupportedSourcesBecomeFindings(t *testing.T) {
	tests := []struct {
		name       string
		targetFile bool
	}{
		{name: "configured directory source"},
		{name: "explicit source file", targetFile: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			path := filepath.Join(tmpDir, "service.rs")
			if err := os.WriteFile(path, []byte("// @implement SPEC-LANGUAGE-001\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			cfg := config.Default()
			cfg.Code.Patterns = []string{"**/*.rs"}
			target := tmpDir
			if test.targetFile {
				target = path
			}

			set, findings, err := NewCodeCollector(cfg).CollectWithErrors(nil, target)
			if err != nil {
				t.Fatalf("CollectWithErrors() error = %v", err)
			}
			if set.Count() != 0 {
				t.Errorf("Identifier count = %d, want zero", set.Count())
			}
			if len(findings) != 1 ||
				findings[0].Rule != "source-parse" ||
				findings[0].Code != "unsupported-extension" {
				t.Fatalf("findings = %#v, want one unsupported source-parse finding", findings)
			}
		})
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-024
func TestCodeCollector_UsesConfiguredAnnotationPrefixes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "service.py")
	source := "# @fulfills SPEC-LANGUAGE-001\ndef run():\n    pass\n"
	if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	cfg := config.Default()
	cfg.Code.Annotations = map[string]string{
		"spec":          "@fulfills",
		"test":          "@checks",
		"test_contract": "@checks-contract",
	}

	set, findings, err := NewCodeCollector(cfg).CollectWithErrors(nil, path)
	if err != nil {
		t.Fatalf("CollectWithErrors() error = %v", err)
	}
	if len(findings) != 0 || !set.HasOrigin("SPEC-LANGUAGE-001", model.OriginCode) {
		t.Errorf("set = %#v, findings = %#v", set.AllIdentifiers(), findings)
	}
}
