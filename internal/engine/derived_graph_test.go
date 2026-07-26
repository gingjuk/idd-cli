// Package engine tests derived document relationships in the linkage graph.

// Spec: docs/internal/engine/spec.md
// Test: docs/internal/engine/testing.md
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
		name      string
		addLink   bool
		wantError bool
	}{
		{"covered", true, false},
		{"uncovered", false, true},
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
			if test.addLink {
				testIdentifier.AddLink(contract.ID)
			}
			ids.Add(testIdentifier)

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
