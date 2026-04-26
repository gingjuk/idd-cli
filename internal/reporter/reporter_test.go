// Package reporter provides testing utilities for the reporter module.

// Spec: docs/internal/reporter/spec.md
// Test: docs/internal/reporter/testing.md
package reporter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jingxu9x/idd-cli/internal/config"
	"github.com/jingxu9x/idd-cli/internal/model"
)

// @test TEST-INTERNAL_REPORTER-009
func TestNew(t *testing.T) {
	cfg := config.Default()

	tests := []struct {
		format string
		want   string
	}{
		{"json", "json"},
		{"markdown", "markdown"},
		{"", "json"},
		{"JSON", "JSON"},
	}

	for _, tt := range tests {
		r := New(cfg, tt.format)
		if r.format != tt.want {
			t.Errorf("New(format=%q) format = %q, want %q", tt.format, r.format, tt.want)
		}
	}
}

// @test TEST-INTERNAL_REPORTER-001
func TestReporter_Generate(t *testing.T) {
	cfg := config.Default()
	r := New(cfg, "json")

	result := &model.ValidationResult{
		Valid:    true,
		Errors:   []model.ValidationError{},
		Warnings: []model.ValidationError{},
	}

	report, err := r.Generate(result)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if report.Tool != "idd-cli" {
		t.Errorf("Tool = %q, want %q", report.Tool, "idd-cli")
	}
	if report.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", report.Version, "1.0.0")
	}
	if report.Timestamp == "" {
		t.Error("Timestamp should not be empty")
	}
}

// @test TEST-INTERNAL_REPORTER-002
func TestReporter_Write_JSON(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "report.json")

	cfg := config.Default()
	r := New(cfg, "json")

	report := &model.Report{
		Tool:      "idd-cli",
		Version:   "1.0.0",
		Timestamp: "2024-01-01T00:00:00Z",
		Config:    model.ConfigSummary{},
		Result: model.ValidationResult{
			Valid: true,
		},
	}

	err := r.Write(report, outputPath)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}
}

// @test TEST-INTERNAL_REPORTER-003
func TestReporter_Write_Stdout(t *testing.T) {
	cfg := config.Default()
	r := New(cfg, "json")

	report := &model.Report{
		Tool:      "idd-cli",
		Version:   "1.0.0",
		Timestamp: "2024-01-01T00:00:00Z",
		Config:    model.ConfigSummary{},
		Result: model.ValidationResult{
			Valid: true,
		},
	}

	err := r.Write(report, "-")
	if err != nil {
		t.Fatalf("Write to stdout failed: %v", err)
	}
}

// @test TEST-INTERNAL_REPORTER-004
func TestReporter_Write_Markdown(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "report.md")

	cfg := config.Default()
	r := New(cfg, "markdown")

	report := &model.Report{
		Tool:      "idd-cli",
		Version:   "1.0.0",
		Timestamp: "2024-01-01T00:00:00Z",
		Config:    model.ConfigSummary{},
		Result: model.ValidationResult{
			Valid:  true,
			Errors: []model.ValidationError{},
			Warnings: []model.ValidationError{
				{Rule: "test-rule", Message: "test warning", Source: "test.go"},
			},
		},
	}

	err := r.Write(report, outputPath)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "# IDD Linkage Report") {
		t.Error("Markdown should contain report header")
	}
	if !strings.Contains(content, "test warning") {
		t.Error("Markdown should contain warning")
	}
}

// @test TEST-INTERNAL_REPORTER-005
func TestReporter_Write_UnsupportedFormat(t *testing.T) {
	cfg := config.Default()
	r := New(cfg, "xml")

	report := &model.Report{
		Tool:    "idd-cli",
		Version: "1.0.0",
	}

	err := r.Write(report, "/dev/null")
	if err == nil {
		t.Error("Expected error for unsupported format")
	}
}

// @test TEST-INTERNAL_REPORTER-006
func TestReporter_writeMarkdown_Errors(t *testing.T) {
	cfg := config.Default()
	r := New(cfg, "markdown")

	report := &model.Report{
		Tool:      "idd-cli",
		Version:   "1.0.0",
		Timestamp: "2024-01-01T00:00:00Z",
		Config:    model.ConfigSummary{},
		Result: model.ValidationResult{
			Valid: false,
			Errors: []model.ValidationError{
				{Rule: "error-rule", Message: "error message", Source: "error.go"},
			},
		},
	}

	var sb strings.Builder
	err := r.writeMarkdown(report, &sb)
	if err != nil {
		t.Fatalf("writeMarkdown failed: %v", err)
	}

	content := sb.String()
	if !strings.Contains(content, "error message") {
		t.Error("Markdown should contain error")
	}
	if !strings.Contains(content, "error-rule") {
		t.Error("Markdown should contain error rule")
	}
}

// @test TEST-INTERNAL_REPORTER-007
func TestReporter_writeMarkdown_Stats(t *testing.T) {
	cfg := config.Default()
	r := New(cfg, "markdown")

	report := &model.Report{
		Tool:      "idd-cli",
		Version:   "1.0.0",
		Timestamp: "2024-01-01T00:00:00Z",
		Config:    model.ConfigSummary{},
		Result: model.ValidationResult{
			Valid: true,
			Stats: model.ValidationStats{
				TotalIdentifiers:  10,
				TotalLinks:        5,
				SpecsAnalyzed:     3,
				TestsAnalyzed:     3,
				ContractsAnalyzed: 2,
				DesignsAnalyzed:   2,
			},
		},
	}

	var sb strings.Builder
	err := r.writeMarkdown(report, &sb)
	if err != nil {
		t.Fatalf("writeMarkdown failed: %v", err)
	}

	content := sb.String()
	if !strings.Contains(content, "**Total Identifiers:**") {
		t.Errorf("Markdown should contain stats, got: %s", content)
	}
}

// @test TEST-INTERNAL_REPORTER-008
func TestReporter_writeMarkdown_Graph(t *testing.T) {
	cfg := config.Default()
	r := New(cfg, "markdown")

	report := &model.Report{
		Tool:      "idd-cli",
		Version:   "1.0.0",
		Timestamp: "2024-01-01T00:00:00Z",
		Config:    model.ConfigSummary{},
		Result: model.ValidationResult{
			Valid: true,
			Graph: &model.GraphSnapshot{
				Nodes: []model.NodeSummary{
					{ID: "SPEC-001", Type: "SPEC", Inbound: 0, Outbound: 1},
				},
				Edges: []model.EdgeSummary{
					{From: "SPEC-001", To: "TEST-001", Type: "tests", Verified: true},
				},
			},
		},
	}

	var sb strings.Builder
	err := r.writeMarkdown(report, &sb)
	if err != nil {
		t.Fatalf("writeMarkdown failed: %v", err)
	}

	content := sb.String()
	if !strings.Contains(content, "SPEC-001") {
		t.Error("Markdown should contain graph nodes")
	}
}

// @test TEST-INTERNAL_REPORTER-010
func TestStatusIcon(t *testing.T) {
	if statusIcon(true) != "✅ PASS" {
		t.Errorf("statusIcon(true) = %q, want %q", statusIcon(true), "✅ PASS")
	}
	if statusIcon(false) != "❌ FAIL" {
		t.Errorf("statusIcon(false) = %q, want %q", statusIcon(false), "❌ FAIL")
	}
}
