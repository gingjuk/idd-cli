// Package engine tests derived document relationships in the linkage graph.
package engine

import (
	"context"
	"testing"

	"github.com/jingxu9x/idd-cli/internal/config"
	"github.com/jingxu9x/idd-cli/internal/model"
)

// @test TEST-INTERNAL_ENGINE-054
func TestEngineContractCoverageUsesNamedContractRelationship(t *testing.T) {
	tests := []struct {
		name       string
		addLink    bool
		testStatus string
		addCode    bool
		wantError  bool
	}{
		{name: "active contract test with source evidence", addLink: true, addCode: true},
		{name: "uncovered", wantError: true},
		{name: "document relation without source evidence", addLink: true, wantError: true},
		{name: "planned contract test is not evidence", addLink: true, addCode: true, testStatus: "planned", wantError: true},
		{name: "deprecated contract test is not evidence", addLink: true, addCode: true, testStatus: "deprecated", wantError: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := &config.Config{
				Validation: config.ValidationConfig{
					RequireContractTestCoverage: true,
					AllowOrphans:                true,
				},
			}
			ids := model.NewIdentifierSet()
			contract := model.NewIdentifier(
				"contract:demo#Service",
				model.TypeContract,
				"Service",
				"docs/demo/contract.md",
				10,
			)
			contract.Derived = true
			ids.Add(contract)
			testIdentifier := model.NewIdentifier(
				"TEST-DEMO-001",
				model.TypeTest,
				"Service contract",
				"docs/demo/testing.md",
				10,
			)
			testIdentifier.Kind = "contract"
			testIdentifier.Status = test.testStatus
			if test.addLink {
				testIdentifier.AddLink(contract.ID)
			}
			ids.Add(testIdentifier)
			if test.addCode {
				code := model.NewIdentifier(testIdentifier.ID, model.TypeTest, "", "service_test.go", 20)
				code.SetOrigin(model.OriginCode)
				code.Kind = "contract"
				ids.Add(code)
			}

			result, err := New(cfg).Run(context.Background(), ids)
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			if got := resultHasRule(result, "contract-test-coverage"); got != test.wantError {
				t.Errorf("contract-test-coverage = %v, want %v; errors = %#v", got, test.wantError, result.Errors)
			}
		})
	}
}

// @test TEST-INTERNAL_ENGINE-054
func TestEngineSpecCoverageRequiresActiveBoundEvidence(t *testing.T) {
	tests := []struct {
		name       string
		specStatus string
		testStatus string
		addTest    bool
		addCode    bool
		wantError  bool
	}{
		{name: "active spec with active bound test", addTest: true, addCode: true},
		{name: "active spec without test", wantError: true},
		{name: "active spec with unbound test", addTest: true, wantError: true},
		{name: "planned test cannot satisfy active spec", addTest: true, addCode: true, testStatus: "planned", wantError: true},
		{name: "deprecated test cannot satisfy active spec", addTest: true, addCode: true, testStatus: "deprecated", wantError: true},
		{name: "planned spec may have no evidence", specStatus: "planned"},
		{name: "deprecated spec may have no evidence", specStatus: "deprecated"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := &config.Config{Validation: config.ValidationConfig{RequireSpecTestCoverage: true, AllowOrphans: true}}
			ids := model.NewIdentifierSet()
			spec := model.NewIdentifier("SPEC-DEMO-001", model.TypeSpec, "behavior", "docs/demo/spec.md", 10)
			spec.Status = test.specStatus
			ids.Add(spec)
			if test.addTest {
				evidence := model.NewIdentifier("TEST-DEMO-001", model.TypeTest, "evidence", "docs/demo/testing.md", 10)
				evidence.Status = test.testStatus
				evidence.AddTypedLink(spec.ID, model.LinkImplements)
				spec.AddTypedLink(evidence.ID, model.LinkTests)
				ids.Add(evidence)
				if test.addCode {
					code := model.NewIdentifier(evidence.ID, model.TypeTest, "", "demo_test.go", 20)
					code.SetOrigin(model.OriginCode)
					code.Kind = "test"
					ids.Add(code)
				}
			}
			result, err := New(cfg).Run(context.Background(), ids)
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			if got := resultHasRule(result, "spec-missing-tests"); got != test.wantError {
				t.Errorf("spec-missing-tests = %v, want %v; errors = %#v", got, test.wantError, result.Errors)
			}
		})
	}
}

// @test TEST-INTERNAL_ENGINE-054
func TestEngineComponentDependencyCycleUsesTypedEdgesOnly(t *testing.T) {
	tests := []struct {
		name      string
		linkType  model.LinkType
		wantCycle bool
	}{
		{"dependency cycle", model.LinkDependsOn, true},
		{"lifecycle relationship is not dependency", model.LinkSupersedes, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := &config.Config{Validation: config.ValidationConfig{AllowOrphans: true}}
			ids := model.NewIdentifierSet()
			first := model.NewIdentifier("component:demo#A", model.TypeDesign, "A", "design.md", 10)
			first.Derived = true
			first.AddTypedLink("component:demo#B", test.linkType)
			ids.Add(first)
			second := model.NewIdentifier("component:demo#B", model.TypeDesign, "B", "design.md", 20)
			second.Derived = true
			second.AddTypedLink("component:demo#A", test.linkType)
			ids.Add(second)

			result, err := New(cfg).Run(context.Background(), ids)
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			if got := resultHasRule(result, "component-dependency-cycle"); got != test.wantCycle {
				t.Errorf("component-dependency-cycle = %v, want %v; errors = %#v", got, test.wantCycle, result.Errors)
			}
		})
	}
}

// @test TEST-INTERNAL_ENGINE-054
func TestEngineGraphTargetValidationPreservesLegacyFreeFormLinks(t *testing.T) {
	tests := []struct {
		name      string
		addLink   func(*model.Identifier)
		wantError bool
	}{
		{
			name: "legacy free-form reference remains compatible",
			addLink: func(identifier *model.Identifier) {
				identifier.AddLink("LegacyContractName")
			},
		},
		{
			name: "missing derived Contract target is rejected",
			addLink: func(identifier *model.Identifier) {
				identifier.AddTypedLink("contract:other/package#Boundary", model.LinkContract)
			},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := &config.Config{
				Validation: config.ValidationConfig{AllowOrphans: true},
			}
			ids := model.NewIdentifierSet()
			spec := model.NewIdentifier(
				"SPEC-DEMO-001",
				model.TypeSpec,
				"Demo",
				"docs/demo/spec.md",
				10,
			)
			test.addLink(spec)
			ids.Add(spec)

			result, err := New(cfg).Run(context.Background(), ids)
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			if got := resultHasRule(result, "idd-document-reference"); got != test.wantError {
				t.Errorf(
					"idd-document-reference = %v, want %v; errors = %#v",
					got,
					test.wantError,
					result.Errors,
				)
			}
		})
	}
}
