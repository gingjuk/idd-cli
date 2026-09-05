package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/jingxu9x/idd-cli/internal/collector"
)

// @test TEST-CMD_IDD_CLI-007
func TestTraceIdentifierExitContract(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T, root string)
		id         string
		wantStatus string
		wantError  error
	}{
		{name: "resolved succeeds", setup: func(t *testing.T, root string) { writeTraceCLIDocs(t, root, "sample", "active") }, id: "SPEC-SAMPLE-001", wantStatus: "resolved"},
		{name: "planned succeeds", setup: func(t *testing.T, root string) { writeTraceCLIDocs(t, root, "sample", "planned") }, id: "SPEC-SAMPLE-001", wantStatus: "planned"},
		{name: "unresolved emits then fails", setup: func(t *testing.T, root string) {}, id: "SPEC-SAMPLE-999", wantStatus: "unresolved", wantError: errTraceResolution},
		{name: "ambiguous emits then fails", setup: func(t *testing.T, root string) {
			writeTraceCLIDocs(t, root, "sample", "active")
			writeTraceCLIDocs(t, root, "other", "active")
		}, id: "SPEC-SAMPLE-001", wantStatus: "ambiguous", wantError: errTraceResolution},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			test.setup(t, root)
			output := filepath.Join(root, "trace.json")
			restore := setTraceTestGlobals(root, output)
			defer restore()
			err := traceIdentifier(nil, []string{test.id})
			if !errors.Is(err, test.wantError) {
				t.Fatalf("traceIdentifier() error = %v, want %v", err, test.wantError)
			}
			data, readErr := os.ReadFile(output)
			if readErr != nil {
				t.Fatalf("ReadFile(output) error = %v", readErr)
			}
			var dossier collector.TraceDossier
			if err := json.Unmarshal(data, &dossier); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if dossier.Status != test.wantStatus {
				t.Fatalf("status = %q, want %q", dossier.Status, test.wantStatus)
			}
			if test.wantStatus == "ambiguous" && (dossier.Canonical != nil || len(dossier.Owners) != 2) {
				t.Fatalf("ambiguous dossier selected an owner: %#v", dossier)
			}
		})
	}
}

// @test TEST-CMD_IDD_CLI-007
func TestParseTraceIncludes(t *testing.T) {
	tests := []struct {
		name      string
		values    []string
		want      bool
		wantError bool
	}{
		{name: "default"},
		{name: "mentions", values: []string{"mentions"}, want: true},
		{name: "case insensitive duplicate", values: []string{"MENTIONS", "mentions"}, want: true},
		{name: "unsupported", values: []string{"source"}, wantError: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseTraceIncludes(test.values)
			if (err != nil) != test.wantError || got != test.want {
				t.Fatalf("parseTraceIncludes() = %t, %v", got, err)
			}
		})
	}
}

func setTraceTestGlobals(root, output string) func() {
	oldProjectRoot, oldOutput, oldFormat := traceProjectRoot, outPath, format
	oldDepth, oldIncludes, oldNoConfig, oldConfig := traceDepth, traceIncludes, noConfig, cfgPath
	traceProjectRoot, outPath, format = root, output, "json"
	traceDepth, traceIncludes, noConfig, cfgPath = 1, nil, true, ""
	return func() {
		traceProjectRoot, outPath, format = oldProjectRoot, oldOutput, oldFormat
		traceDepth, traceIncludes, noConfig, cfgPath = oldDepth, oldIncludes, oldNoConfig, oldConfig
	}
}

func writeTraceCLIDocs(t *testing.T, root, packageName, status string) {
	t.Helper()
	directory := filepath.Join(root, "docs", packageName)
	namespace := "SAMPLE"
	frontmatter := "---\nidd:\n  version: \"1.1\"\n  package: " + packageName + "\n  namespace: " + namespace + "\n---\n\n"
	files := map[string]string{
		"design.md":   frontmatter + "# Design\n\n## Component: Runner\n\n**Purpose:** Run the behavior.\n\n**Ownership:** Own execution.\n\n**Boundary:** Accept and return values.\n\n**Decisions:** Keep execution deterministic.\n",
		"contract.md": frontmatter + "# Contracts\n\nNo stable external contract is declared.\n",
		"spec.md":     frontmatter + "# Specifications\n\n## SPEC-SAMPLE-001: Run behavior\n\n- **Status:** `" + status + "`\n- **Components:** `Runner`\n\n**Requirement:** Run the selected behavior.\n\n**Acceptance:** The behavior returns its value.\n",
		"testing.md":  frontmatter + "# Testing\n\nNo executable evidence is registered for this fixture.\n",
	}
	if err := os.MkdirAll(directory, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(directory, name), []byte(content), 0644); err != nil {
			t.Fatalf("WriteFile(%s) error = %v", name, err)
		}
	}
}
