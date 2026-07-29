// Package main tests the stable command and embedded-Skill boundaries.
package main

import (
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
