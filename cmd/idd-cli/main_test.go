// Package main tests the stable command and embedded-Skill boundaries.
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// @test TEST-CMD_IDD_CLI-001
func TestCommandSurfaceBehavior(t *testing.T) {
	for _, name := range []string{"run", "lint", "skills", "generate", "docs"} {
		if commandNamed(rootCmd, name) == nil {
			t.Errorf("root command missing %q", name)
		}
	}
	docs := commandNamed(rootCmd, "docs")
	for _, name := range []string{"init", "fix", "status", "review-context"} {
		if commandNamed(docs, name) == nil {
			t.Errorf("docs command missing %q", name)
		}
	}
}

// @test TEST-CMD_IDD_CLI-002
func TestRunUsesOneProjectRoot(t *testing.T) {
	originalDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDirectory); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	callerRoot := t.TempDir()
	targetRoot := t.TempDir()
	writeRunProjectFixture(t, callerRoot, "CALLER", true)
	writeRunProjectFixture(t, targetRoot, "TARGET", false)
	if err := os.Chdir(callerRoot); err != nil {
		t.Fatalf("Chdir(%s) error = %v", callerRoot, err)
	}

	savedConfigPath, savedOutputPath := cfgPath, outPath
	savedFormat, savedVerbose, savedNoConfig := format, verbose, noConfig
	t.Cleanup(func() {
		cfgPath, outPath = savedConfigPath, savedOutputPath
		format, verbose, noConfig = savedFormat, savedVerbose, savedNoConfig
	})
	cfgPath = ""
	outPath = filepath.Join(t.TempDir(), "report.json")
	format = "json"
	verbose = false
	noConfig = false

	err = run(nil, []string{targetRoot})
	if !errors.Is(err, errValidationFailed) {
		t.Fatalf("run(%s) error = %v, want validation failure", targetRoot, err)
	}
	currentDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() after run error = %v", err)
	}
	if currentDirectory != callerRoot {
		t.Errorf("working directory after run = %q, want %q", currentDirectory, callerRoot)
	}

	report, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("ReadFile(%s) error = %v", outPath, err)
	}
	text := string(report)
	if !strings.Contains(text, "SPEC-TARGET-001") ||
		!strings.Contains(text, "orphan-detection") {
		t.Errorf("target project evidence missing from report:\n%s", text)
	}
	if strings.Contains(text, "SPEC-CALLER-001") ||
		strings.Contains(text, "doc-code-correspondence") {
		t.Errorf("caller project leaked into target report:\n%s", text)
	}
}

// @test TEST-CMD_IDD_CLI-001
func TestResolveRunProjectRoot(t *testing.T) {
	directory := t.TempDir()
	file := filepath.Join(directory, "spec.md")
	if err := os.WriteFile(file, []byte("# test\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{name: "default current directory", want: mustAbsolutePath(t, ".")},
		{name: "explicit directory", args: []string{directory}, want: directory},
		{name: "file is rejected", args: []string{file}, wantErr: true},
		{name: "missing directory is rejected", args: []string{filepath.Join(directory, "missing")}, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := resolveRunProjectRoot(test.args)
			if (err != nil) != test.wantErr {
				t.Fatalf("resolveRunProjectRoot() error = %v, wantErr %t", err, test.wantErr)
			}
			if !test.wantErr && got != filepath.Clean(test.want) {
				t.Errorf("resolveRunProjectRoot() = %q, want %q", got, test.want)
			}
		})
	}
}

// @test TEST-CMD_IDD_CLI-001
func TestResolveInvocationOutputPath(t *testing.T) {
	absolute := filepath.Join(t.TempDir(), "report.json")
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "stdout default"},
		{name: "explicit stdout", path: "-", want: "-"},
		{name: "absolute path", path: absolute, want: absolute},
		{name: "relative path", path: "report.json", want: mustAbsolutePath(t, "report.json")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := resolveInvocationOutputPath(test.path)
			if err != nil {
				t.Fatalf("resolveInvocationOutputPath() error = %v", err)
			}
			if got != test.want {
				t.Errorf("resolveInvocationOutputPath() = %q, want %q", got, test.want)
			}
		})
	}
}

// @test TEST-CMD_IDD_CLI-002
func TestEmbeddedSkillIntegrationBehavior(t *testing.T) {
	paths := ListEmbeddedSkills()
	if len(paths) != 1 || paths[0] != "skills/SKILL.md" {
		t.Fatalf("ListEmbeddedSkills() = %#v", paths)
	}
	data, err := ReadEmbeddedSkill(paths[0])
	if err != nil {
		t.Fatalf("ReadEmbeddedSkill() error = %v", err)
	}
	if !strings.Contains(string(data), "# Intent-Driven Development (IDD)") {
		t.Error("embedded Skill body is incomplete")
	}
}

// @test-contract TEST-CMD_IDD_CLI-005
func TestCommandSurfaceContract(t *testing.T) {
	runCommand := commandNamed(rootCmd, "run")
	if runCommand == nil || runCommand.RunE == nil {
		t.Fatal("run command must expose an executable RunE boundary")
	}
	for _, flag := range []string{"config", "format", "output", "verbose", "no-config"} {
		if rootCmd.PersistentFlags().Lookup(flag) == nil {
			t.Errorf("persistent CLI contract missing --%s", flag)
		}
	}
	docs := commandNamed(rootCmd, "docs")
	if docs == nil {
		t.Fatal("document command group is absent from the CLI contract")
	}
	for _, name := range []string{"init", "fix", "status", "review-context"} {
		command := commandNamed(docs, name)
		if command == nil || command.RunE == nil {
			t.Errorf("document command %q is absent from the executable CLI contract", name)
		}
	}
	if review := commandNamed(docs, "review-context"); review != nil &&
		review.Flags().Lookup("docs-path") == nil {
		t.Error("review-context command is missing --docs-path")
	}
}

// @test-contract TEST-CMD_IDD_CLI-006
func TestEmbeddedSkillContract(t *testing.T) {
	data, err := ReadEmbeddedSkill("skills/SKILL.md")
	if err != nil {
		t.Fatalf("ReadEmbeddedSkill() error = %v", err)
	}
	info := parseSkillFrontmatter(string(data))
	if info.Name != "intent-driven-development" ||
		info.Audience != "agents" ||
		info.Workflow != "development" ||
		!info.Protected {
		t.Errorf("parsed SkillInfo = %#v", info)
	}
	if _, err := ReadEmbeddedSkill("skills/missing.md"); err == nil {
		t.Error("missing embedded Skill must return an error")
	}
}

func commandNamed(parent *cobra.Command, name string) *cobra.Command {
	if parent == nil {
		return nil
	}
	for _, command := range parent.Commands() {
		if command.Name() == name {
			return command
		}
	}
	return nil
}

func writeRunProjectFixture(t *testing.T, root, module string, allowOrphans bool) {
	t.Helper()
	docsRoot := filepath.Join(root, "docs")
	if err := os.MkdirAll(docsRoot, 0o755); err != nil {
		t.Fatalf("MkdirAll(%s) error = %v", docsRoot, err)
	}
	specID := "SPEC-" + module + "-001"
	document := fmt.Sprintf(`---
markers:
  - id: %s
    name: %s behavior
    describe: %s project evidence
---

## %s: %s behavior
`, specID, module, module, specID, module)
	if err := os.WriteFile(filepath.Join(docsRoot, "spec.md"), []byte(document), 0o644); err != nil {
		t.Fatalf("write project document: %v", err)
	}
	source := fmt.Sprintf(`package fixture

// @implement %s
func Run() {}
`, specID)
	if err := os.WriteFile(filepath.Join(root, "fixture.go"), []byte(source), 0o644); err != nil {
		t.Fatalf("write project source: %v", err)
	}
	configuration := fmt.Sprintf(`version: "1.0"
docs:
  patterns:
    - "docs/**/*.md"
code:
  patterns:
    - "**/*.go"
  annotations:
    spec: "@implement"
    test: "@test"
    test_contract: "@test-contract"
validation:
  allow_orphans: %t
  require_doc_code_correspondence: %t
`, allowOrphans, !allowOrphans)
	if err := os.WriteFile(filepath.Join(root, ".idd.yaml"), []byte(configuration), 0o644); err != nil {
		t.Fatalf("write project configuration: %v", err)
	}
}

func mustAbsolutePath(t *testing.T, path string) string {
	t.Helper()
	absolute, err := filepath.Abs(path)
	if err != nil {
		t.Fatalf("Abs(%s) error = %v", path, err)
	}
	return absolute
}
