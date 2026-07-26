// Package engine tests syntax-tree-backed source validation.

// Spec: docs/internal/engine/spec.md
// Test: docs/internal/engine/testing.md
package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jingxu9x/idd-cli/internal/config"
	"github.com/jingxu9x/idd-cli/internal/model"
)

// @test TEST-INTERNAL_ENGINE-052
func TestSourcePublicAnnotationsAcrossLanguages(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		declaration string
	}{
		{"go", "service.go", "package service\n%sfunc Run() {}\n"},
		{"typescript", "service.ts", "%sexport function run() {}\n"},
		{"tsx", "view.tsx", "%sexport function View() { return <div/> }\n"},
		{"javascript", "service.js", "%sexport function run() {}\n"},
		{"cpp", "service.cpp", "%sint run() { return 0; }\n"},
		{"java", "Service.java", "%spublic class Service {}\n"},
		{"python", "service.py", "%sdef run():\n    pass\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, annotated := range []bool{true, false} {
				caseName := "missing"
				prefix := ""
				if annotated {
					caseName = "bound"
					prefix = sourceCommentPrefix(test.filename) + " @implement SPEC-LANGUAGE-001\n"
				}
				t.Run(caseName, func(t *testing.T) {
					dir := t.TempDir()
					path := filepath.Join(dir, test.filename)
					source := []byte(fmt.Sprintf(test.declaration, prefix))
					if err := os.WriteFile(path, source, 0o644); err != nil {
						t.Fatalf("WriteFile() error = %v", err)
					}
					cfg := &config.Config{
						Code: config.CodeConfig{Patterns: []string{path}},
						Validation: config.ValidationConfig{
							RequirePublicFuncAnnotation: true,
						},
					}
					result, err := New(cfg).Run(context.Background(), model.NewIdentifierSet())
					if err != nil {
						t.Fatalf("Run() error = %v", err)
					}
					hasMissing := resultHasRule(result, "public-func-annotation")
					if hasMissing == annotated {
						t.Errorf(
							"public-func-annotation present = %v, annotated = %v; errors = %#v",
							hasMissing,
							annotated,
							result.Errors,
						)
					}
				})
			}
		})
	}
}

// @test TEST-INTERNAL_ENGINE-052
func TestSourcePublicAnnotationsUseConfiguredPrefix(t *testing.T) {
	path := filepath.Join(t.TempDir(), "service.py")
	if err := os.WriteFile(
		path,
		[]byte("# @fulfills SPEC-LANGUAGE-001\ndef run():\n    pass\n"),
		0o644,
	); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	cfg := &config.Config{
		Code: config.CodeConfig{
			Patterns: []string{path},
			Annotations: map[string]string{
				"spec":          "@fulfills",
				"test":          "@checks",
				"test_contract": "@checks-contract",
			},
		},
		Validation: config.ValidationConfig{
			RequirePublicFuncAnnotation: true,
		},
	}
	result, err := New(cfg).Run(context.Background(), model.NewIdentifierSet())
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if resultHasRule(result, "public-func-annotation") ||
		resultHasRule(result, "annotation-placement") {
		t.Errorf("errors = %#v, want configured annotation to bind", result.Errors)
	}
}

// @test TEST-INTERNAL_ENGINE-053
func TestSourceTestAnnotationsAcrossLanguages(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		source   string
	}{
		{"go", "service_test.go", "package service\n// @test TEST-LANGUAGE-001\nfunc TestRun(t *testing.T) {}\n"},
		{"typescript", "service.test.ts", "// @test TEST-LANGUAGE-001\ntest('runs', () => {})\n"},
		{"tsx", "view.test.tsx", "// @test TEST-LANGUAGE-001\ntest('renders', () => <div/>)\n"},
		{"javascript", "service.test.js", "// @test TEST-LANGUAGE-001\ntest('runs', () => {})\n"},
		{"cpp", "service_test.cpp", "// @test TEST-LANGUAGE-001\nvoid TestRun() {}\n"},
		{"java", "ServiceTest.java", "class ServiceTest {\n// @test TEST-LANGUAGE-001\n@Test public void runs() {}\n}\n"},
		{"python", "test_service.py", "# @test TEST-LANGUAGE-001\ndef test_run():\n    pass\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), test.filename)
			if err := os.WriteFile(path, []byte(test.source), 0o644); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}
			cfg := &config.Config{
				Code: config.CodeConfig{Patterns: []string{path}},
				Validation: config.ValidationConfig{
					RequireTestAnnotation: true,
				},
			}
			ids := model.NewIdentifierSet()
			ids.Add(model.NewIdentifier(
				"TEST-LANGUAGE-001",
				model.TypeTest,
				"",
				"docs/language/testing.md",
				1,
			))
			result, err := New(cfg).Run(context.Background(), ids)
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			if resultHasRule(result, "test-annotation") ||
				resultHasRule(result, "annotation-placement") {
				t.Errorf("errors = %#v, want no source test annotation finding", result.Errors)
			}
		})
	}
}

// @test TEST-INTERNAL_ENGINE-053
func TestSourceContractTestRequiresContractAnnotation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "service_contract_test.cpp")
	if err := os.WriteFile(
		path,
		[]byte("// @test TEST-LANGUAGE-001\nvoid TestContract() {}\n"),
		0o644,
	); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	cfg := &config.Config{
		Code: config.CodeConfig{Patterns: []string{path}},
		Validation: config.ValidationConfig{
			RequireTestAnnotation: true,
		},
	}
	result, err := New(cfg).Run(context.Background(), model.NewIdentifierSet())
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !resultHasRule(result, "test-annotation") {
		t.Errorf("errors = %#v, want test-annotation", result.Errors)
	}
}

func sourceCommentPrefix(path string) string {
	if filepath.Ext(path) == ".py" {
		return "#"
	}
	return "//"
}

func resultHasRule(result *model.ValidationResult, rule string) bool {
	for _, validationError := range result.Errors {
		if validationError.Rule == rule {
			return true
		}
	}
	return false
}
