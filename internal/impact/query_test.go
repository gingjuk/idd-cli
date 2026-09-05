package impact

import (
	"reflect"
	"strings"
	"testing"

	"github.com/jingxu9x/idd-cli/internal/collector"
	"github.com/jingxu9x/idd-cli/internal/graph"
	"github.com/jingxu9x/idd-cli/internal/model"
)

// @test TEST-CMD_IDD_CLI-008
func TestBuildReport(t *testing.T) {
	index := impactFixtureIndex()
	tests := []struct {
		name         string
		files        []FileChange
		wantIDs      []string
		wantWarnings int
	}{
		{
			name: "implementation change traverses complete review context",
			files: []FileChange{{
				Path: "internal/auth/service.go", Status: "modified",
				AnnotationIDs: []string{"SPEC-AUTH-001"},
			}},
			wantIDs:      []string{"SPEC-AUTH-001", "TEST-AUTH-001", "contract:auth#Authenticator", "component:auth#Service"},
			wantWarnings: 1,
		},
		{
			name: "canonical SPEC change suppresses unchanged warning",
			files: []FileChange{
				{Path: "docs/auth/spec.md", Status: "modified"},
				{Path: "internal/auth/service.go", Status: "modified"},
			},
			wantIDs:      []string{"SPEC-AUTH-001", "TEST-AUTH-001", "contract:auth#Authenticator", "component:auth#Service"},
			wantWarnings: 0,
		},
		{
			name: "unrelated prose remains only a changed file",
			files: []FileChange{{
				Path: "notes.md", Status: "untracked",
				AnnotationIDs: []string{"SPEC-EXAMPLE-999"},
			}},
			wantIDs:      []string{},
			wantWarnings: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			report, err := BuildReport(&ChangeSet{Schema: ChangeSetSchema, Base: "main", Files: test.files}, index)
			if err != nil {
				t.Fatalf("BuildReport() error = %v", err)
			}
			ids := make([]string, 0, len(report.Queue))
			for _, item := range report.Queue {
				ids = append(ids, item.ID)
			}
			if !reflect.DeepEqual(ids, test.wantIDs) {
				t.Errorf("BuildReport() IDs = %#v, want %#v", ids, test.wantIDs)
			}
			if len(report.Warnings) != test.wantWarnings {
				t.Errorf("BuildReport() warnings = %#v, want %d", report.Warnings, test.wantWarnings)
			}
			if test.name == "unrelated prose remains only a changed file" && len(report.Changes[0].AnnotationIDs) != 0 {
				t.Errorf("BuildReport() retained unindexed annotation examples: %#v", report.Changes[0].AnnotationIDs)
			}
		})
	}
}

// @test TEST-CMD_IDD_CLI-008
func TestBuildReportRetainsDeletedAnnotations(t *testing.T) {
	report, err := BuildReport(
		&ChangeSet{Schema: ChangeSetSchema, Base: "main", Files: []FileChange{{
			Path: "removed.go", Status: "deleted", AnnotationIDs: []string{"SPEC-REMOVED-001"},
		}}},
		impactFixtureIndex(),
	)
	if err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}
	if got := report.Changes[0].AnnotationIDs; !reflect.DeepEqual(got, []string{"SPEC-REMOVED-001"}) {
		t.Errorf("BuildReport() deleted annotation IDs = %#v", got)
	}
}

// @test TEST-CMD_IDD_CLI-008
func TestBuildReportUsesCollectedDeclarationRanges(t *testing.T) {
	index := impactFixtureIndex()
	analyses := []*collector.SourceAnalysis{{
		Path: "internal/auth/service.go",
		Declarations: []collector.SourceDeclaration{{
			Name: "service", Line: 20, EndLine: 30,
			Annotations: []collector.SourceAnnotation{{
				Kind: "implement", Refs: []string{"SPEC-AUTH-001"},
				Line: 19, Attached: true, Declaration: "service", DeclarationLine: 20,
			}},
		}},
	}}
	report, err := BuildReportWithAnalyses(
		&ChangeSet{Schema: ChangeSetSchema, Base: "main", Files: []FileChange{{
			Path: "internal/auth/service.go", Status: "modified",
			Hunks: []Hunk{{OldStart: 24, OldLines: 1, NewStart: 24, NewLines: 1}},
		}}},
		index,
		analyses,
	)
	if err != nil {
		t.Fatalf("BuildReportWithAnalyses() error = %v", err)
	}
	var spec *QueueItem
	for index := range report.Queue {
		if report.Queue[index].ID == "SPEC-AUTH-001" {
			spec = &report.Queue[index]
			break
		}
	}
	if spec == nil {
		t.Fatal("changed declaration did not seed SPEC-AUTH-001")
	}
	if len(spec.Reasons) == 0 || spec.Reasons[0].Kind != "changed_declaration" || spec.Reasons[0].Conservative {
		t.Errorf("SPEC reasons = %#v, want an exact changed_declaration", spec.Reasons)
	}
}

// @test TEST-CMD_IDD_CLI-008
func TestRender(t *testing.T) {
	report := &Report{
		Schema: ReportSchema, Base: "main", Depth: 2,
		Changes: []FileChange{{Path: "service.go", Status: "modified"}},
		Queue: []QueueItem{{
			ID: "SPEC-AUTH-001", Kind: "spec", Lifecycle: "active", Resolution: "resolved",
			Canonical: &Location{Path: "docs/auth/spec.md", Line: 8},
			Reasons:   []Reason{{Kind: "changed_annotation", Path: "service.go"}},
		}},
		Warnings: []Warning{{Rule: "related-spec-unchanged", Entity: "SPEC-AUTH-001", Paths: []string{"service.go"}, Message: "review alignment"}},
	}
	tests := []struct {
		name     string
		format   string
		contains string
		wantErr  bool
	}{
		{name: "JSON", format: "json", contains: `"schema": "idd.impacted.v1"`},
		{name: "Markdown", format: "llm-markdown", contains: "# IDD impacted review queue"},
		{name: "unsupported", format: "xml", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := Render(report, test.format)
			if (err != nil) != test.wantErr {
				t.Fatalf("Render() error = %v, wantErr %t", err, test.wantErr)
			}
			if !test.wantErr && !strings.Contains(string(output), test.contains) {
				t.Errorf("Render() output missing %q:\n%s", test.contains, output)
			}
		})
	}
}

func impactFixtureIndex() *graph.TraceIndex {
	identifiers := model.NewIdentifierSet()
	component := model.NewIdentifier("component:auth#Service", model.TypeDesign, "Service", "docs/auth/design.md", 8)
	component.Kind = "component"
	component.Derived = true
	contract := model.NewIdentifier("contract:auth#Authenticator", model.TypeContract, "Authenticator", "docs/auth/contract.md", 8)
	contract.Kind = "contract"
	contract.Derived = true
	spec := model.NewIdentifier("SPEC-AUTH-001", model.TypeSpec, "Authenticate", "docs/auth/spec.md", 8)
	spec.Status = "active"
	spec.AddTypedLinkAt(component.ID, model.LinkReferences, "docs/auth/spec.md", 10, "Components")
	spec.AddTypedLinkAt(contract.ID, model.LinkContract, "docs/auth/spec.md", 11, "Contracts")
	testRecord := model.NewIdentifier("TEST-AUTH-001", model.TypeTest, "Authentication evidence", "docs/auth/testing.md", 8)
	testRecord.AddTypedLinkAt(spec.ID, model.LinkImplements, "docs/auth/testing.md", 10, "Covers")
	implementation := model.NewIdentifier("SPEC-AUTH-001", model.TypeSpec, "service", "internal/auth/service.go", 20)
	implementation.SetOrigin(model.OriginCode)
	testDeclaration := model.NewIdentifier("TEST-AUTH-001", model.TypeTest, "TestService", "internal/auth/service_test.go", 20)
	testDeclaration.SetOrigin(model.OriginCode)
	for _, identifier := range []*model.Identifier{component, contract, spec, testRecord, implementation, testDeclaration} {
		identifiers.Add(identifier)
	}
	return graph.BuildTraceIndex(identifiers)
}
