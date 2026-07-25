// Package engine provides testing utilities for the engine module.

// Spec: docs/internal/engine/spec.md
// Test: docs/internal/engine/testing.md
// Contract: docs/internal/engine/contract.md

package engine

import (
	"context"
	"testing"

	"github.com/jingxu9x/idd-cli/internal/config"
	"github.com/jingxu9x/idd-cli/internal/graph"
	"github.com/jingxu9x/idd-cli/internal/model"
)

// CONTRACT-BE-002: Engine Validation Contracts
//
// The validation engine orchestrates the collection, linking, and validation process.
// This contract defines the engine's responsibilities and expected behaviors.
//
// Engine Interface:
//   - Engine struct with cfg, collector, graph, rules fields
//   - Rule interface with Name() and Validate(g *graph.LinkageGraph) []model.ValidationError
//
// Validation Process:
//   1. Collect — Gather identifiers from docs and code
//   2. Link — Build graph with edges from forward links
//   3. Verify — Check doc-link-consistency
//   4. Report — Aggregate errors and generate output

// @test-contract TEST-INTERNAL_ENGINE-001
func TestEngineStruct(t *testing.T) {
	cfg := config.Default()
	e := New(cfg)

	if e == nil {
		t.Fatal("New(cfg) returned nil")
	}

	if e.cfg == nil {
		t.Error("Engine.cfg is nil")
	}

	if e.graph == nil {
		t.Error("Engine.graph is nil")
	}

	if e.result == nil {
		t.Error("Engine.result is nil")
	}
}

// @test-contract TEST-INTERNAL_ENGINE-002
func TestEngineNewConstructor(t *testing.T) {
	cfg := config.Default()
	e := New(cfg)

	if e == nil {
		t.Fatal("New(cfg) should not return nil")
	}

	// Verify initial state - Valid is false initially (zero value)
	// Run() calls validate() which sets Valid=true before validation checks
	if e.result.Valid {
		t.Error("New engine should have Valid=false initially (before Run)")
	}

	if len(e.result.Errors) != 0 {
		t.Errorf("New engine should have no errors, got %d", len(e.result.Errors))
	}
}

// @test-contract TEST-INTERNAL_ENGINE-003
func TestEngineRunPerformsValidationProcess(t *testing.T) {
	cfg := config.Default()
	cfg.Validation.AllowOrphans = true // Allow orphans for this test

	e := New(cfg)

	// Create identifiers with doc-link-consistency
	ids := model.NewIdentifierSet()

	spec1 := model.NewIdentifier("SPEC-BE-001", model.TypeSpec, "Test Spec", "test.md", 1)
	spec1.AddLink("TEST-BE-001") // SPEC points to TEST
	ids.Add(spec1)

	test1 := model.NewIdentifier("TEST-BE-001", model.TypeTest, "Test Case", "test.go", 10)
	test1.AddLink("SPEC-BE-001") // TEST points back to SPEC (doc-link-consistency)
	ids.Add(test1)

	ctx := context.Background()
	result, err := e.Run(ctx, ids)

	if err != nil {
		t.Fatalf("Engine.Run() returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Engine.Run() returned nil result")
	}

	// Verify stats are populated (proves graph was built)
	if result.Stats.TotalIdentifiers != 2 {
		t.Errorf("Expected 2 identifiers, got %d", result.Stats.TotalIdentifiers)
	}

	if result.Stats.TotalLinks != 2 {
		t.Errorf("Expected 2 links, got %d", result.Stats.TotalLinks)
	}

	// Verify graph was built correctly
	if result.Stats.TotalLinks != 2 {
		t.Errorf("Expected 2 links, got %d", result.Stats.TotalLinks)
	}
}

// @test-contract TEST-INTERNAL_ENGINE-004
func TestEngineRunWithUnidirectionalLink(t *testing.T) {
	// NOTE: Bidirectional link checking was removed per user request.
	// Unidirectional links are now allowed without warnings.
	// This test verifies the engine still processes unidirectional links correctly.

	cfg := config.Default()
	cfg.Validation.AllowOrphans = true

	e := New(cfg)

	// Create identifiers with ONLY unidirectional link
	ids := model.NewIdentifierSet()

	spec1 := model.NewIdentifier("SPEC-BE-001", model.TypeSpec, "Test Spec", "test.md", 1)
	spec1.AddLink("TEST-BE-001") // SPEC points to TEST
	ids.Add(spec1)

	test1 := model.NewIdentifier("TEST-BE-001", model.TypeTest, "Test Case", "test.go", 10)
	// TEST does NOT link back to SPEC (unidirectional)
	ids.Add(test1)

	ctx := context.Background()
	result, err := e.Run(ctx, ids)

	if err != nil {
		t.Fatalf("Engine.Run() returned error: %v", err)
	}

	// Verify the graph was built with the link
	if result.Stats.TotalLinks != 1 {
		t.Errorf("Expected 1 link, got %d", result.Stats.TotalLinks)
	}

	// No missing-backlink warnings should be present (feature was removed)
	for _, warn := range result.Warnings {
		if warn.Rule == "missing-backlink" {
			t.Errorf("Unexpected missing-backlink warning: %s", warn.Message)
		}
	}
}

// @test-contract TEST-INTERNAL_ENGINE-005
func TestEngineValidationProcessSteps(t *testing.T) {
	cfg := config.Default()
	cfg.Validation.AllowOrphans = false // Orphans should generate errors

	e := New(cfg)

	// Create an orphan identifier (no links)
	ids := model.NewIdentifierSet()
	orphan := model.NewIdentifier("SPEC-BE-999", model.TypeSpec, "Orphan Spec", "orphan.md", 1)
	ids.Add(orphan)

	ctx := context.Background()
	result, err := e.Run(ctx, ids)

	if err != nil {
		t.Fatalf("Engine.Run() returned error: %v", err)
	}

	// Orphan should be detected
	foundOrphanError := false
	for _, err := range result.Errors {
		if err.Rule == "orphan-detection" {
			foundOrphanError = true
			break
		}
	}

	if !foundOrphanError {
		t.Error("Expected orphan-detection error for orphan identifier")
	}
}

// @test-contract TEST-INTERNAL_ENGINE-006
func TestRuleInterface(t *testing.T) {
	// CONTRACT-BE-002 defines Rule interface:
	// type Rule interface {
	//     Name() string
	//     Validate(g *graph.LinkageGraph) []model.ValidationError
	// }

	// Create a concrete implementation to test the interface
	rule := &testRule{name: "test-rule", shouldError: false}

	// Verify Name() method
	if rule.Name() != "test-rule" {
		t.Errorf("Rule.Name() = %q, want %q", rule.Name(), "test-rule")
	}

	// Verify Validate() method
	g := graph.NewLinkageGraph()
	errors := rule.Validate(g)

	if len(errors) != 0 {
		t.Errorf("Rule.Validate() returned %d errors, want 0", len(errors))
	}
}

// @test-contract TEST-INTERNAL_ENGINE-007
func TestRuleInterfaceWithErrors(t *testing.T) {
	rule := &testRule{name: "failing-rule", shouldError: true}

	g := graph.NewLinkageGraph()
	errors := rule.Validate(g)

	if len(errors) != 1 {
		t.Errorf("Rule.Validate() returned %d errors, want 1", len(errors))
	}

	if errors[0].Rule != "failing-rule" {
		t.Errorf("Error.Rule = %q, want %q", errors[0].Rule, "failing-rule")
	}
}

// @test-contract TEST-INTERNAL_ENGINE-008
func TestBidirectionalLinkRuleContract(t *testing.T) {
	// NOTE: Bidirectional link checking was removed per user request.
	// This test verifies the engine still processes links correctly.

	cfg := config.Default()
	cfg.Validation.AllowOrphans = true

	e := New(cfg)

	// Create a graph with bidirectional links
	ids := model.NewIdentifierSet()
	spec1 := model.NewIdentifier("SPEC-BE-001", model.TypeSpec, "Spec", "s.md", 1)
	spec1.AddLink("TEST-BE-001")
	ids.Add(spec1)

	test1 := model.NewIdentifier("TEST-BE-001", model.TypeTest, "Test", "t.go", 5)
	test1.AddLink("SPEC-BE-001")
	ids.Add(test1)

	ctx := context.Background()
	result, _ := e.Run(ctx, ids)

	// Verify the graph was built correctly
	if result.Stats.TotalLinks != 2 {
		t.Errorf("Expected 2 links, got %d", result.Stats.TotalLinks)
	}
}

// @test-contract TEST-INTERNAL_ENGINE-009
func TestOrphanRuleContract(t *testing.T) {
	// According to CONTRACT-BE-002, OrphanRule should detect unreferenced identifiers

	cfg := config.Default()
	cfg.Validation.AllowOrphans = false

	e := New(cfg)

	// Create orphan identifier
	ids := model.NewIdentifierSet()
	orphan := model.NewIdentifier("SPEC-BE-999", model.TypeSpec, "Orphan", "orphan.md", 1)
	ids.Add(orphan)

	ctx := context.Background()
	result, _ := e.Run(ctx, ids)

	// Orphan should be detected
	foundOrphan := false
	for _, err := range result.Errors {
		if err.Rule == "orphan-detection" {
			foundOrphan = true
			break
		}
	}

	if !foundOrphan {
		t.Error("Orphan rule should detect unreferenced identifier")
	}
}

// @test-contract TEST-INTERNAL_ENGINE-010
func TestEngineValidateMethod(t *testing.T) {
	cfg := config.Default()
	cfg.Validation.AllowOrphans = true
	cfg.Validation.RequireSpecTestCoverage = false
	cfg.Validation.RequireDocCodeCorrespondence = false

	e := New(cfg)

	// Add structural errors - AddStructuralErrors sets Valid=false via AddError
	e.AddStructuralErrors([]*model.ValidationError{
		{
			Rule:    "test-rule",
			Message: "Test error",
			Source:  "test.go",
			Link:    "TEST-001",
			Code:    "",
		},
	})

	// Verify AddStructuralErrors worked
	if e.result.Valid {
		t.Error("Result should not be valid after adding errors")
	}

	if len(e.result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(e.result.Errors))
	}

	// Run validation
	e.validate()

	// Since structural errors exist, Valid should remain false even after validate()
	if e.result.Valid {
		t.Error("Result should not be valid after validate() with existing errors")
	}
}

// @test-contract TEST-INTERNAL_ENGINE-011
func TestEngineBuildReport(t *testing.T) {
	cfg := config.Default()
	e := New(cfg)

	ids := model.NewIdentifierSet()
	spec1 := model.NewIdentifier("SPEC-BE-001", model.TypeSpec, "Test", "test.md", 1)
	ids.Add(spec1)

	ctx := context.Background()
	e.Run(ctx, ids)

	report := e.BuildReport()

	if report == nil {
		t.Fatal("BuildReport() returned nil")
	}

	if report.Tool != "idd-cli" {
		t.Errorf("Report.Tool = %q, want %q", report.Tool, "idd-cli")
	}

	if report.Version != "1.0.0" {
		t.Errorf("Report.Version = %q, want %q", report.Version, "1.0.0")
	}

	if report.Timestamp == "" {
		t.Error("Report.Timestamp should not be empty")
	}
}

// @test-contract TEST-INTERNAL_ENGINE-012
func TestEngineInferLinkType(t *testing.T) {
	cfg := config.Default()
	e := New(cfg)

	tests := []struct {
		fromType model.IdentifierType
		toRef    string
		want     model.LinkType
	}{
		{model.TypeSpec, "TEST-BE-001", model.LinkTests},
		{model.TypeSpec, "SPEC-BE-002", model.LinkReferences},
		{model.TypeSpec, "CONTRACT-BE-001", model.LinkContract},
		{model.TypeTest, "SPEC-BE-001", model.LinkImplements},
		{model.TypeContract, "SPEC-BE-001", model.LinkContractImplements},
		{model.TypeDesign, "SPEC-BE-001", model.LinkReferences},
	}

	for _, tt := range tests {
		got := e.inferLinkType(tt.fromType, tt.toRef)
		if got != tt.want {
			t.Errorf("inferLinkType(%v, %q) = %v, want %v", tt.fromType, tt.toRef, got, tt.want)
		}
	}
}

// testRule is a test implementation of the Rule interface per CONTRACT-BE-002
type testRule struct {
	name        string
	shouldError bool
}

func (r *testRule) Name() string {
	return r.name
}

func (r *testRule) Validate(g *graph.LinkageGraph) []model.ValidationError {
	if r.shouldError {
		return []model.ValidationError{
			{
				Rule:    r.name,
				Message: "Test error from rule",
				Source:  "test",
				Link:    "",
				Code:    "",
			},
		}
	}
	return nil
}

// Compile-time check that testRule implements Rule interface
var _ Rule = (*testRule)(nil)

// Rule interface as defined in CONTRACT-BE-002
// @implement SPEC-CMD_IDD_CLI-010
type Rule interface {
	Name() string
	Validate(g *graph.LinkageGraph) []model.ValidationError
}
