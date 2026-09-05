package reporter

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/jingxu9x/idd-cli/internal/collector"
)

// @test-contract TEST-INTERNAL_REPORTER-014
func TestRenderSpecReviewContext(t *testing.T) {
	contextValue := &collector.SpecReviewContext{
		Schema: "idd.spec_review_context.v2",
		Spec: collector.ReviewContextSpec{
			ID:          "SPEC-SAMPLE-001",
			Title:       "Preserve input",
			Package:     "sample",
			File:        "docs/sample/spec.md",
			Line:        9,
			Components:  []string{"component:sample#Runner"},
			Contracts:   []string{"contract:sample#Execution"},
			Requirement: "Return the supplied value.",
			Acceptance:  "The result equals the input.",
		},
		Components: []collector.ReviewContextComponent{{
			ID:       "component:sample#Runner",
			Name:     "Runner",
			Purpose:  "Execute the sample behavior.",
			Boundary: "Accept and return a value.",
			File:     "docs/sample/design.md",
			Line:     9,
		}},
		Contracts: []collector.ReviewContextContract{{
			ID:         "contract:sample#Execution",
			Name:       "Execution",
			Guarantees: "The input is preserved.",
			File:       "docs/sample/contract.md",
			Line:       9,
		}},
		Tests: []collector.ReviewContextTest{{
			ID:      "TEST-SAMPLE-001",
			Title:   "Preservation evidence",
			Kind:    "test",
			Purpose: "Exercise preservation.",
			Oracle:  "Result equals input.",
			File:    "docs/sample/testing.md",
			Line:    9,
		}},
		Declarations: []collector.ReviewContextDeclaration{{
			Role:       "implementation",
			References: []string{"SPEC-SAMPLE-001"},
			Language:   "go",
			Kind:       "function",
			Name:       "execute",
			File:       "sample.go",
			Line:       4,
			EndLine:    6,
			Annotation: "// @implement SPEC-SAMPLE-001",
			Excerpt:    "func execute(value string) string { return value }",
		}},
		Truncated: true,
	}
	tests := []struct {
		name       string
		format     string
		wantError  bool
		wantValues []string
	}{
		{
			name:       "json",
			format:     "json",
			wantValues: []string{`"schema": "idd.spec_review_context.v2"`, `"id": "SPEC-SAMPLE-001"`},
		},
		{
			name:   "markdown",
			format: "markdown",
			wantValues: []string{
				"# Review context: SPEC-SAMPLE-001",
				"## Review questions",
				"func execute",
				"record or declaration evidence was truncated",
			},
		},
		{
			name:       "llm markdown",
			format:     "llm-markdown",
			wantValues: []string{"# Review context: SPEC-SAMPLE-001", "TEST-SAMPLE-001", "Execution"},
		},
		{name: "unsupported", format: "xml", wantError: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rendered, err := RenderSpecReviewContext(contextValue, test.format)
			if test.wantError {
				if err == nil {
					t.Fatal("RenderSpecReviewContext() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("RenderSpecReviewContext() error = %v", err)
			}
			for _, value := range test.wantValues {
				if !strings.Contains(string(rendered), value) {
					t.Errorf("rendered output does not contain %q", value)
				}
			}
			if test.format == "json" {
				var decoded collector.SpecReviewContext
				if err := json.Unmarshal(rendered, &decoded); err != nil {
					t.Fatalf("json.Unmarshal() error = %v", err)
				}
			}
			if strings.Contains(string(rendered), "semantic_score") ||
				strings.Contains(string(rendered), "review_passed") {
				t.Errorf("rendered output added a semantic verdict: %s", rendered)
			}
		})
	}
}

// @test-contract TEST-INTERNAL_REPORTER-014
func TestRenderSpecReviewContexts(t *testing.T) {
	contexts := []*collector.SpecReviewContext{
		{
			Schema: "idd.spec_review_context.v2",
			Spec: collector.ReviewContextSpec{
				ID:      "SPEC-SAMPLE-002",
				Title:   "Second",
				Package: "sample",
			},
		},
		{
			Schema: "idd.spec_review_context.v2",
			Spec: collector.ReviewContextSpec{
				ID:      "SPEC-SAMPLE-001",
				Title:   "First",
				Package: "sample",
			},
		},
	}
	tests := []struct {
		name       string
		batch      *collector.SpecReviewContextBatch
		format     string
		wantError  bool
		wantValues []string
	}{
		{
			name: "single retains single schema",
			batch: &collector.SpecReviewContextBatch{
				Schema:           "idd.spec_review_context_batch.v2",
				RequestedSpecIDs: []string{"SPEC-SAMPLE-002"},
				Contexts:         contexts[:1],
			},
			format: "json",
			wantValues: []string{
				`"schema": "idd.spec_review_context.v2"`,
				`"id": "SPEC-SAMPLE-002"`,
			},
		},
		{
			name: "batch json preserves order",
			batch: &collector.SpecReviewContextBatch{
				Schema:           "idd.spec_review_context_batch.v2",
				RequestedSpecIDs: []string{"SPEC-SAMPLE-002", "SPEC-SAMPLE-001"},
				Contexts:         contexts,
			},
			format: "json",
			wantValues: []string{
				`"schema": "idd.spec_review_context_batch.v2"`,
				`"requested_spec_ids"`,
				`"id": "SPEC-SAMPLE-002"`,
				`"id": "SPEC-SAMPLE-001"`,
			},
		},
		{
			name: "batch markdown has independent sections",
			batch: &collector.SpecReviewContextBatch{
				Schema:           "idd.spec_review_context_batch.v2",
				RequestedSpecIDs: []string{"SPEC-SAMPLE-002", "SPEC-SAMPLE-001"},
				Contexts:         contexts,
			},
			format: "llm-markdown",
			wantValues: []string{
				"# SPEC review context batch",
				"## Review context: SPEC-SAMPLE-002",
				"## Review context: SPEC-SAMPLE-001",
			},
		},
		{
			name:      "empty batch",
			batch:     &collector.SpecReviewContextBatch{},
			format:    "json",
			wantError: true,
		},
		{
			name: "unsupported batch format",
			batch: &collector.SpecReviewContextBatch{
				Contexts: contexts,
			},
			format:    "xml",
			wantError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rendered, err := RenderSpecReviewContexts(test.batch, test.format)
			if test.wantError {
				if err == nil {
					t.Fatal("RenderSpecReviewContexts() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("RenderSpecReviewContexts() error = %v", err)
			}
			text := string(rendered)
			for _, value := range test.wantValues {
				if !strings.Contains(text, value) {
					t.Errorf("rendered output does not contain %q", value)
				}
			}
			if test.name == "batch json preserves order" &&
				strings.Index(text, "SPEC-SAMPLE-002") > strings.LastIndex(text, "SPEC-SAMPLE-001") {
				t.Errorf("batch output did not preserve request order: %s", text)
			}
		})
	}
}
