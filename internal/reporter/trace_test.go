package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/jingxu9x/idd-cli/internal/collector"
)

// @test TEST-INTERNAL_REPORTER-016
func TestRenderTraceDossier(t *testing.T) {
	dossier := traceReporterFixture()
	tests := []struct {
		name      string
		format    string
		wantError string
		contains  []string
	}{
		{name: "json", format: "json", contains: []string{`"schema": "idd.trace.v1"`, `"occurrence_id": "occ-components"`, `"record_id": "SPEC-SAMPLE-001"`}},
		{name: "markdown", format: "markdown", contains: []string{"# IDD trace: SPEC-SAMPLE-001", "## Components and design context", "--`designed_by`-->", "field `Components`"}},
		{name: "llm markdown", format: "llm-markdown", contains: []string{"## Canonical declaration", "## Test evidence", "## Findings"}},
		{name: "unsupported", format: "yaml", wantError: "unsupported trace format"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			first, err := RenderTraceDossier(dossier, test.format)
			if test.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantError) {
					t.Fatalf("error = %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("RenderTraceDossier() error = %v", err)
			}
			second, err := RenderTraceDossier(dossier, test.format)
			if err != nil {
				t.Fatalf("second RenderTraceDossier() error = %v", err)
			}
			if !bytes.Equal(first, second) {
				t.Fatal("rendering is not byte stable")
			}
			for _, value := range test.contains {
				if !strings.Contains(string(first), value) {
					t.Errorf("output lacks %q", value)
				}
			}
			if test.format == "json" {
				if strings.Contains(string(first), `"mentions": null`) || strings.Contains(string(first), `"inbound_relations": null`) {
					t.Errorf("stable collection fields must be arrays, not null: %s", first)
				}
				var decoded collector.TraceDossier
				if err := json.Unmarshal(first, &decoded); err != nil {
					t.Fatalf("json.Unmarshal() error = %v", err)
				}
				if decoded.Occurrences[0].Path != "a.go" || decoded.Occurrences[1].Path != "z.go" {
					t.Errorf("occurrences not stably sorted: %#v", decoded.Occurrences)
				}
			}
		})
	}
}

// @test TEST-INTERNAL_REPORTER-016
func TestRenderTraceDossierUnresolvedAndNil(t *testing.T) {
	tests := []struct {
		name      string
		dossier   *collector.TraceDossier
		wantError string
		contains  []string
	}{
		{
			name:      "nil dossier is rejected",
			wantError: "trace dossier is nil",
		},
		{
			name: "unresolved dossier retains an explicit empty projection",
			dossier: &collector.TraceDossier{
				Schema:    collector.TraceSchema,
				Query:     collector.TraceQuery{ID: "SPEC-MISSING-001"},
				Status:    "unresolved",
				Truncated: true,
			},
			contains: []string{
				"No logical entity was resolved",
				"No unique canonical declaration was resolved",
				"No related Component was collected",
				"No related Contract was collected",
				"No related TEST evidence was collected",
				"Some authored records or declaration excerpts were truncated",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rendered, err := RenderTraceDossier(test.dossier, "markdown")
			if test.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantError) {
					t.Fatalf("error = %v, want containing %q", err, test.wantError)
				}
				return
			}
			if err != nil {
				t.Fatalf("RenderTraceDossier() error = %v", err)
			}
			for _, value := range test.contains {
				if !strings.Contains(string(rendered), value) {
					t.Errorf("output lacks %q", value)
				}
			}
		})
	}
}

func traceReporterFixture() *collector.TraceDossier {
	return &collector.TraceDossier{
		Schema:    collector.TraceSchema,
		Query:     collector.TraceQuery{ID: "SPEC-SAMPLE-001", Depth: 1, IncludeMentions: true},
		Status:    "resolved",
		Entity:    &collector.TraceEntity{ID: "SPEC-SAMPLE-001", Kind: "spec", Namespace: "SAMPLE", Lifecycle: "active", Resolution: "resolved"},
		Canonical: &collector.TraceCanonical{TraceLocation: collector.TraceLocation{Path: "docs/sample/spec.md", Line: 12}, Title: "Preserve input", Requirement: "Return the supplied value.", Acceptance: "The result equals the input."},
		Owners:    []collector.TraceOccurrence{{ID: "owner", EntityID: "SPEC-SAMPLE-001", Kind: "declaration", Origin: "doc", Path: "docs/sample/spec.md", Line: 12}},
		Occurrences: []collector.TraceOccurrence{
			{ID: "z", EntityID: "SPEC-SAMPLE-001", Kind: "implementation", Origin: "source", Path: "z.go", Line: 9},
			{ID: "a", EntityID: "SPEC-SAMPLE-001", Kind: "implementation", Origin: "source", Path: "a.go", Line: 3},
		},
		OutboundRelations: []collector.TraceRelation{{From: "SPEC-SAMPLE-001", To: "component:sample#Runner", Type: "designed_by", Depth: 1, Provenance: collector.TraceProvenance{OccurrenceID: "occ-components", Path: "docs/sample/spec.md", Line: 14, Field: "Components", RecordID: "SPEC-SAMPLE-001", OccurrenceKind: "reference"}}},
		InboundRelations:  []collector.TraceRelation{}, RelatedEntities: []collector.TraceRelatedEntity{},
		Components:      []collector.TraceComponent{{ID: "component:sample#Runner", Name: "Runner", Purpose: "Run behavior.", File: "docs/sample/design.md", Line: 10}},
		Contracts:       []collector.TraceContract{{ID: "contract:sample#Execution", Name: "Execution", Guarantees: "Preserve input.", File: "docs/sample/contract.md", Line: 10}},
		Implementations: []collector.ReviewContextDeclaration{},
		Tests:           []collector.TraceTest{{ID: "TEST-SAMPLE-001", Title: "Evidence", Kind: "test", Purpose: "Exercise behavior.", Oracle: "Result equals input.", Covers: []string{"SPEC-SAMPLE-001"}, Contracts: []string{}, Declarations: []collector.ReviewContextDeclaration{}}},
		Mentions:        []collector.TraceOccurrence{},
		Findings:        []collector.ReviewContextIssue{{Rule: "example", Message: "Review this evidence."}},
	}
}
