package collector

import (
	"context"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/jingxu9x/idd-cli/internal/config"
)

// @test TEST-INTERNAL_COLLECTOR-033
func TestTraceProjectQueries(t *testing.T) {
	root := t.TempDir()
	writeReviewFixture(t, root, "sample", "SPEC-SAMPLE-001")
	writeReviewFile(t, filepath.Join(root, "source", "sample.go"), `package sample

// @implement SPEC-SAMPLE-001
func execute(value string) string { return value }

// @test TEST-SAMPLE-001
func TestExecute() {}
`)
	designPath := filepath.Join(root, "docs", "sample", "design.md")
	design := readTestFile(t, designPath)
	design = strings.Replace(design, "## Component: Runner\n\n**Purpose:**", "## Component: Runner\n\n- **Depends on:** `Helper`\n\n**Purpose:**", 1)
	design += `

## Component: Helper

**Purpose:** Support sample execution.

**Ownership:** Own internal value helpers.

**Boundary:** Accept values only from Runner.

**Decisions:** Remain an internal implementation boundary.

The behavior is tracked by ` + "`SPEC-SAMPLE-001`" + `.
`
	writeTestFile(t, designPath, design)

	cfg := traceTestConfig(t, root)
	project, err := BuildTraceProject(context.Background(), cfg, ".", ".")
	if err != nil {
		t.Fatalf("BuildTraceProject() error = %v", err)
	}

	tests := []struct {
		name          string
		id            string
		options       TraceOptions
		wantKind      string
		wantStatus    string
		wantCanonical bool
		check         func(t *testing.T, dossier *TraceDossier)
	}{
		{
			name: "SPEC dossier includes semantics evidence and provenance", id: "SPEC-SAMPLE-001",
			options: TraceOptions{Depth: 1}, wantKind: "spec", wantStatus: "resolved", wantCanonical: true,
			check: func(t *testing.T, dossier *TraceDossier) {
				if dossier.Canonical.Requirement != "Return the supplied value without mutation." {
					t.Errorf("canonical = %#v", dossier.Canonical)
				}
				if len(dossier.Components) != 1 || dossier.Components[0].Name != "Runner" {
					t.Errorf("components = %#v", dossier.Components)
				}
				if len(dossier.Contracts) != 1 || dossier.Contracts[0].Name != "Execution" {
					t.Errorf("contracts = %#v", dossier.Contracts)
				}
				if len(dossier.Tests) != 1 || len(dossier.Tests[0].Declarations) != 1 {
					t.Errorf("tests = %#v", dossier.Tests)
				}
				if len(dossier.Implementations) != 1 || dossier.Implementations[0].Name != "execute" {
					t.Errorf("implementations = %#v", dossier.Implementations)
				}
				if !relationHasField(dossier.OutboundRelations, "Components") || !relationHasField(dossier.OutboundRelations, "Covers") {
					t.Errorf("outbound relations lack field provenance: %#v", dossier.OutboundRelations)
				}
				if len(dossier.Mentions) != 0 || relationHasType(dossier.InboundRelations, "mentions") {
					t.Errorf("mentions leaked into default dossier: %#v %#v", dossier.Mentions, dossier.InboundRelations)
				}
			},
		},
		{name: "TEST", id: "TEST-SAMPLE-001", options: TraceOptions{Depth: 1}, wantKind: "test", wantStatus: "resolved", wantCanonical: true},
		{name: "Contract", id: "contract:sample#Execution", options: TraceOptions{Depth: 1}, wantKind: "contract", wantStatus: "resolved", wantCanonical: true},
		{name: "Component", id: "component:sample#Runner", options: TraceOptions{Depth: 1}, wantKind: "component", wantStatus: "resolved", wantCanonical: true},
		{
			name: "depth two reaches dependency", id: "SPEC-SAMPLE-001", options: TraceOptions{Depth: 2},
			wantKind: "spec", wantStatus: "resolved", wantCanonical: true,
			check: func(t *testing.T, dossier *TraceDossier) {
				if !relatedAtDepth(dossier.RelatedEntities, "component:sample#Helper", 2) {
					t.Errorf("related entities = %#v", dossier.RelatedEntities)
				}
			},
		},
		{
			name: "mentions are opt in", id: "SPEC-SAMPLE-001", options: TraceOptions{Depth: 1, IncludeMentions: true},
			wantKind: "spec", wantStatus: "resolved", wantCanonical: true,
			check: func(t *testing.T, dossier *TraceDossier) {
				if len(dossier.Mentions) != 1 || dossier.Mentions[0].Kind != "mention" {
					t.Errorf("mentions = %#v", dossier.Mentions)
				}
			},
		},
		{name: "unknown remains queryable", id: "SPEC-SAMPLE-999", options: TraceOptions{Depth: 1}, wantStatus: "unresolved"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dossier, err := project.Trace(test.id, test.options)
			if err != nil {
				t.Fatalf("Trace() error = %v", err)
			}
			if dossier.Status != test.wantStatus {
				t.Errorf("status = %q, want %q", dossier.Status, test.wantStatus)
			}
			if dossier.Entity != nil && dossier.Entity.Kind != test.wantKind {
				t.Errorf("kind = %q, want %q", dossier.Entity.Kind, test.wantKind)
			}
			if (dossier.Canonical != nil) != test.wantCanonical {
				t.Errorf("canonical = %#v, want present %t", dossier.Canonical, test.wantCanonical)
			}
			if test.check != nil {
				test.check(t, dossier)
			}
		})
	}
}

// @test TEST-INTERNAL_COLLECTOR-033
func TestTraceProjectResolutionAndLifecycle(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T, root string)
		id         string
		wantStatus string
		wantOwners int
	}{
		{
			name: "ambiguous preserves every owner", id: "SPEC-SAMPLE-001", wantStatus: "ambiguous", wantOwners: 2,
			setup: func(t *testing.T, root string) {
				writeReviewFixture(t, root, "sample", "SPEC-SAMPLE-001")
				writeReviewFixture(t, root, "other", "SPEC-SAMPLE-001")
			},
		},
		{
			name: "planned permits empty evidence", id: "SPEC-SAMPLE-001", wantStatus: "planned", wantOwners: 1,
			setup: func(t *testing.T, root string) {
				writeReviewFixture(t, root, "sample", "SPEC-SAMPLE-001")
				specPath := filepath.Join(root, "docs", "sample", "spec.md")
				spec := readTestFile(t, specPath)
				spec = strings.Replace(spec, "- **Components:**", "- **Status:** `planned`\n- **Components:**", 1)
				writeTestFile(t, specPath, spec)
				testingPath := filepath.Join(root, "docs", "sample", "testing.md")
				writeTestFile(t, testingPath, `---
idd:
  version: "1.1"
  package: sample
  namespace: SAMPLE
---

# Testing

No executable evidence exists while the approved behavior is planned.
`)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			test.setup(t, root)
			project, err := BuildTraceProject(context.Background(), traceTestConfig(t, root), ".", ".")
			if err != nil {
				t.Fatalf("BuildTraceProject() error = %v", err)
			}
			dossier, err := project.Trace(test.id, TraceOptions{Depth: 1})
			if err != nil {
				t.Fatalf("Trace() error = %v", err)
			}
			if dossier.Status != test.wantStatus || len(dossier.Owners) != test.wantOwners {
				t.Fatalf("status/owners = %q/%d, want %q/%d", dossier.Status, len(dossier.Owners), test.wantStatus, test.wantOwners)
			}
			if test.wantStatus == "planned" && (len(dossier.Implementations) != 0 || len(dossier.Tests) != 0) {
				t.Errorf("planned evidence = implementations %#v tests %#v", dossier.Implementations, dossier.Tests)
			}
		})
	}
}

func traceTestConfig(t *testing.T, root string) *config.Config {
	t.Helper()
	cfg := config.Default()
	if err := cfg.SetWorkdir(root); err != nil {
		t.Fatalf("SetWorkdir() error = %v", err)
	}
	cfg.Docs.Patterns = []string{"docs/**/*.md"}
	cfg.Docs.IgnorePaths = nil
	cfg.Code.Patterns = []string{"source/**/*.go"}
	cfg.Code.IgnorePaths = nil
	return cfg
}

func relationHasField(relations []TraceRelation, field string) bool {
	for _, relation := range relations {
		if relation.Provenance.Field == field {
			return true
		}
	}
	return false
}

func relationHasType(relations []TraceRelation, relationType string) bool {
	for _, relation := range relations {
		if relation.Type == relationType {
			return true
		}
	}
	return false
}

func relatedAtDepth(entities []TraceRelatedEntity, id string, depth int) bool {
	for _, entity := range entities {
		if entity.ID == id && entity.Depth == depth {
			return true
		}
	}
	return false
}

// @test TEST-INTERNAL_COLLECTOR-033
func TestTraceProjectExternalWorkdirIsStable(t *testing.T) {
	root := t.TempDir()
	writeReviewFixture(t, root, "sample", "SPEC-SAMPLE-001")
	project, err := BuildTraceProject(context.Background(), traceTestConfig(t, root), ".", ".")
	if err != nil {
		t.Fatalf("BuildTraceProject() error = %v", err)
	}
	dossier, err := project.Trace("SPEC-SAMPLE-001", TraceOptions{Depth: 1})
	if err != nil {
		t.Fatalf("Trace() error = %v", err)
	}
	if dossier.Canonical == nil || filepath.IsAbs(dossier.Canonical.Path) || dossier.Canonical.Path != filepath.ToSlash(filepath.Join("docs", "sample", "spec.md")) {
		t.Fatalf("canonical location = %#v", dossier.Canonical)
	}
	got := []string{dossier.Components[0].Name, dossier.Contracts[0].Name}
	if !reflect.DeepEqual(got, []string{"Runner", "Execution"}) {
		t.Fatalf("context = %#v", got)
	}
}

// @test TEST-INTERNAL_COLLECTOR-033
func TestValidateTraceQuery(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		options TraceOptions
		wantErr bool
	}{
		{name: "SPEC", id: "SPEC-SAMPLE-001", options: TraceOptions{Depth: 1}},
		{name: "scoped component", id: "component:sample#Runner", options: TraceOptions{Depth: 5}},
		{name: "empty", wantErr: true},
		{name: "wrong shape", id: "sample", wantErr: true},
		{name: "negative depth", id: "SPEC-SAMPLE-001", options: TraceOptions{Depth: -1}, wantErr: true},
		{name: "excessive depth", id: "SPEC-SAMPLE-001", options: TraceOptions{Depth: 6}, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := validateTraceQuery(test.id, test.options); (err != nil) != test.wantErr {
				t.Fatalf("validateTraceQuery() error = %v, wantErr %t", err, test.wantErr)
			}
		})
	}
}
