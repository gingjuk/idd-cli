// Package reporter provides testing utilities for the reporter module.
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

// @test-contract TEST-INTERNAL_REPORTER-001
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

	var parsed model.LLMReport
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}
	if parsed.Schema != "idd.llm_report.v1" {
		t.Fatalf("Schema = %q, want idd.llm_report.v1", parsed.Schema)
	}
	if parsed.Status != "pass" {
		t.Fatalf("Status = %q, want pass", parsed.Status)
	}
}

// @test-contract TEST-INTERNAL_REPORTER-011
func TestReporter_Write_DefaultJSONUsesLLMSchema(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "report.json")

	cfg := config.Default()
	r := New(cfg, "")

	report := &model.Report{
		Tool:      "idd-cli",
		Version:   "1.0.0",
		Timestamp: "2024-01-01T00:00:00Z",
		Config:    model.ConfigSummary{},
		Result: model.ValidationResult{
			Valid: true,
		},
	}

	if err := r.Write(report, outputPath); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var parsed model.LLMReport
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}
	if parsed.Schema != "idd.llm_report.v1" {
		t.Fatalf("Schema = %q, want idd.llm_report.v1", parsed.Schema)
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

// @test-contract TEST-INTERNAL_REPORTER-005
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

// @test-contract TEST-INTERNAL_REPORTER-011
func TestReporter_Write_JSONFindingReport(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "report.json")

	cfg := config.Default()
	r := New(cfg, "json")

	report := sampleLLMSourceReport()
	if err := r.Write(report, outputPath); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var parsed model.LLMReport
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}
	if parsed.Schema != "idd.llm_report.v1" {
		t.Fatalf("Schema = %q, want idd.llm_report.v1", parsed.Schema)
	}
	if parsed.Status != "fail" {
		t.Fatalf("Status = %q, want fail", parsed.Status)
	}
	if parsed.Summary.Errors != 1 {
		t.Fatalf("Summary.Errors = %d, want 1", parsed.Summary.Errors)
	}
	if len(parsed.Findings) != 1 {
		t.Fatalf("len(Findings) = %d, want 1", len(parsed.Findings))
	}
	finding := parsed.Findings[0]
	if finding.Location.File != "docs/internal/foo/spec.md" {
		t.Fatalf("Location.File = %q, want docs/internal/foo/spec.md", finding.Location.File)
	}
	if finding.Location.Line != 27 {
		t.Fatalf("Location.Line = %d, want 27", finding.Location.Line)
	}
	if finding.Identifier != "SPEC-INTERNAL_FOO-001" {
		t.Fatalf("Identifier = %q, want SPEC-INTERNAL_FOO-001", finding.Identifier)
	}
	if finding.SuggestedFix == "" {
		t.Fatal("SuggestedFix should not be empty")
	}
}

// @test TEST-INTERNAL_REPORTER-012
func TestReporter_Write_LLMMarkdown(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "llm-report.md")

	cfg := config.Default()
	r := New(cfg, "llm-markdown")

	report := sampleLLMSourceReport()
	if err := r.Write(report, outputPath); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	content := string(data)
	for _, want := range []string{
		"# IDD Validation Report",
		"Status: FAIL",
		"## Finding groups",
		"docs/internal/foo/spec.md:27",
		"Rule: doc-link-consistency",
		"Problem:",
		"Fix:",
		"TEST-INTERNAL_FOO-001",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("LLM Markdown missing %q:\n%s", want, content)
		}
	}
}

// @test TEST-INTERNAL_REPORTER-013
func TestReporter_buildLLMReport_EnrichesFindings(t *testing.T) {
	cfg := config.Default()
	r := New(cfg, "json")

	llmReport := r.buildLLMReport(sampleLLMSourceReport())
	if len(llmReport.Summary.TopRules) != 1 || llmReport.Summary.TopRules[0] != "doc-link-consistency" {
		t.Fatalf("TopRules = %v, want [doc-link-consistency]", llmReport.Summary.TopRules)
	}
	if len(llmReport.Summary.RuleGroups) != 1 {
		t.Fatalf("len(RuleGroups) = %d, want 1", len(llmReport.Summary.RuleGroups))
	}
	if llmReport.Summary.RuleGroups[0].Count != 1 {
		t.Fatalf("RuleGroups[0].Count = %d, want 1", llmReport.Summary.RuleGroups[0].Count)
	}
	if len(llmReport.Summary.RuleGroups[0].FindingIndexes) != 1 || llmReport.Summary.RuleGroups[0].FindingIndexes[0] != 1 {
		t.Fatalf("RuleGroups[0].FindingIndexes = %v, want [1]", llmReport.Summary.RuleGroups[0].FindingIndexes)
	}
	if len(llmReport.Findings) != 1 {
		t.Fatalf("len(Findings) = %d, want 1", len(llmReport.Findings))
	}

	finding := llmReport.Findings[0]
	if finding.Severity != "error" {
		t.Fatalf("Severity = %q, want error", finding.Severity)
	}
	if finding.Title == "" {
		t.Fatal("Title should not be empty")
	}
	if finding.Expected == "" {
		t.Fatal("Expected should not be empty")
	}
	if finding.Actual == "" {
		t.Fatal("Actual should not be empty")
	}
	if len(finding.RelatedIdentifiers) != 1 {
		t.Fatalf("len(RelatedIdentifiers) = %d, want 1", len(finding.RelatedIdentifiers))
	}
	if finding.RelatedIdentifiers[0].ID != "TEST-INTERNAL_FOO-001" {
		t.Fatalf("RelatedIdentifiers[0].ID = %q, want TEST-INTERNAL_FOO-001", finding.RelatedIdentifiers[0].ID)
	}
}

// @test TEST-INTERNAL_REPORTER-013
func TestReporter_buildLLMReport_PreservesWarningSeverity(t *testing.T) {
	cfg := config.Default()
	r := New(cfg, "json")

	report := &model.Report{
		Result: model.ValidationResult{
			Valid:  true,
			Errors: []model.ValidationError{},
			Warnings: []model.ValidationError{
				{
					Rule:    "orphan-detection",
					Message: "DESIGN-INTERNAL_FOO-001 is not referenced by any identifier",
					Source:  "docs/internal/foo/design.md",
					Link:    "DESIGN-INTERNAL_FOO-001",
				},
			},
		},
	}

	llmReport := r.buildLLMReport(report)
	if len(llmReport.Findings) != 1 {
		t.Fatalf("len(Findings) = %d, want 1", len(llmReport.Findings))
	}
	if llmReport.Findings[0].Severity != "warning" {
		t.Fatalf("Severity = %q, want warning", llmReport.Findings[0].Severity)
	}
}

// @test TEST-INTERNAL_REPORTER-013
func TestReporter_buildLLMReport_GroupsRepeatedFindings(t *testing.T) {
	cfg := config.Default()
	r := New(cfg, "json")

	report := &model.Report{
		Result: model.ValidationResult{
			Valid: false,
			Errors: []model.ValidationError{
				{
					Rule:    "doc-link-consistency",
					Message: "spec.md should use **Tests:** not **Spec Coverage:**",
					Source:  "docs/internal/foo/spec.md:12",
					Link:    "SPEC-INTERNAL_FOO-001",
				},
				{
					Rule:    "doc-link-consistency",
					Message: "spec.md should use **Tests:** not **Spec Coverage:**",
					Source:  "docs/internal/bar/spec.md:21",
					Link:    "SPEC-INTERNAL_BAR-001",
				},
			},
		},
	}

	llmReport := r.buildLLMReport(report)
	if len(llmReport.Summary.RuleGroups) != 1 {
		t.Fatalf("len(RuleGroups) = %d, want 1", len(llmReport.Summary.RuleGroups))
	}

	group := llmReport.Summary.RuleGroups[0]
	if group.Rule != "doc-link-consistency" {
		t.Fatalf("Rule = %q, want doc-link-consistency", group.Rule)
	}
	if group.Count != 2 {
		t.Fatalf("Count = %d, want 2", group.Count)
	}
	if got, want := group.FindingIndexes, []int{1, 2}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("FindingIndexes = %v, want %v", got, want)
	}
	if len(group.Files) != 2 {
		t.Fatalf("len(Files) = %d, want 2", len(group.Files))
	}
	if len(group.Identifiers) != 2 {
		t.Fatalf("len(Identifiers) = %d, want 2", len(group.Identifiers))
	}
}

// @test TEST-INTERNAL_REPORTER-013
func TestLookupRuleInfo_IDDDocumentFindings(t *testing.T) {
	tests := []struct {
		rule        string
		wantFixText string
	}{
		{rule: "idd-document-parse", wantFixText: "YAML"},
		{rule: "idd-document-identity", wantFixText: "docs fix"},
		{rule: "idd-document-set", wantFixText: "docs fix"},
		{rule: "idd-document-schema", wantFixText: "field"},
		{rule: "idd-document-incomplete", wantFixText: "docs status"},
		{rule: "idd-document-reference", wantFixText: "declaration"},
		{rule: "idd-document-markdown", wantFixText: "legacy"},
		{rule: "idd-document-migration", wantFixText: "testing.md"},
		{rule: "idd-document-test-kind", wantFixText: "annotation"},
		{rule: "source-parse", wantFixText: "supported source extension"},
		{rule: "doc-code-correspondence", wantFixText: "@implement"},
	}

	for _, tt := range tests {
		t.Run(tt.rule, func(t *testing.T) {
			info := lookupRuleInfo(tt.rule, "error")
			if info.Title == "" || info.Explanation == "" {
				t.Fatalf("lookupRuleInfo(%q) returned incomplete metadata: %#v", tt.rule, info)
			}
			if !strings.Contains(info.FixHint, tt.wantFixText) {
				t.Errorf("lookupRuleInfo(%q).FixHint = %q, want text %q", tt.rule, info.FixHint, tt.wantFixText)
			}
		})
	}
}

func sampleLLMSourceReport() *model.Report {
	return &model.Report{
		Tool:      "idd-cli",
		Version:   "1.0.0",
		Timestamp: "2024-01-01T00:00:00Z",
		Config:    model.ConfigSummary{},
		Result: model.ValidationResult{
			Valid: false,
			Errors: []model.ValidationError{
				{
					Rule:    "doc-link-consistency",
					Message: "spec.md should use **Tests:** not **Spec Coverage:**",
					Source:  "docs/internal/foo/spec.md:27",
					Link:    "SPEC-INTERNAL_FOO-001",
					Code:    "**Spec Coverage:** `TEST-INTERNAL_FOO-001`",
				},
			},
			Warnings: []model.ValidationError{},
			Stats: model.ValidationStats{
				TotalIdentifiers: 2,
				TotalLinks:       1,
				SpecsAnalyzed:    1,
				TestsAnalyzed:    1,
			},
			Graph: &model.GraphSnapshot{
				Nodes: []model.NodeSummary{
					{ID: "SPEC-INTERNAL_FOO-001", Type: model.TypeSpec, Outbound: 1},
					{ID: "TEST-INTERNAL_FOO-001", Type: model.TypeTest, Inbound: 1},
				},
				Edges: []model.EdgeSummary{
					{
						From:     "SPEC-INTERNAL_FOO-001",
						To:       "TEST-INTERNAL_FOO-001",
						Type:     "tests",
						Verified: false,
						Source:   "docs/internal/foo/spec.md",
						Line:     27,
					},
				},
			},
		},
	}
}
