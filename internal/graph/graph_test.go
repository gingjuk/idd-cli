// Package graph provides testing utilities for the graph module.

// Spec: docs/internal/graph/spec.md
// Test: docs/internal/graph/testing.md
package graph

import (
	"testing"

	"github.com/jingxu9x/idd-cli/internal/model"
)

// @test TEST-INT_GRPH-015
func TestNewLinkageGraph(t *testing.T) {
	g := NewLinkageGraph()

	if g.NodeCount() != 0 {
		t.Errorf("NodeCount() = %d, want 0", g.NodeCount())
	}
	if g.EdgeCount() != 0 {
		t.Errorf("EdgeCount() = %d, want 0", g.EdgeCount())
	}
}

// @test TEST-INT_GRPH-001
func TestLinkageGraph_AddNode(t *testing.T) {
	g := NewLinkageGraph()

	n := g.AddNode("SPEC-001", model.TypeSpec)
	if n == nil {
		t.Fatal("AddNode returned nil")
	}
	if n.ID != "SPEC-001" {
		t.Errorf("n.ID = %q, want %q", n.ID, "SPEC-001")
	}
	if g.NodeCount() != 1 {
		t.Errorf("NodeCount() = %d, want 1", g.NodeCount())
	}

	n2 := g.AddNode("SPEC-001", model.TypeSpec)
	if n2 != n {
		t.Error("AddNode should return existing node for duplicate")
	}
	if g.NodeCount() != 1 {
		t.Errorf("NodeCount() = %d, want 1 (no new node)", g.NodeCount())
	}
}

// @test TEST-INT_GRPH-002
func TestLinkageGraph_AddEdge(t *testing.T) {
	g := NewLinkageGraph()
	g.AddNode("SPEC-001", model.TypeSpec)
	g.AddNode("TEST-001", model.TypeTest)

	g.AddEdge("SPEC-001", "TEST-001", model.LinkTests, "test.go", 10)

	if g.EdgeCount() != 1 {
		t.Errorf("EdgeCount() = %d, want 1", g.EdgeCount())
	}

	edges := g.Edges()
	if len(edges) != 1 {
		t.Errorf("len(Edges()) = %d, want 1", len(edges))
	}
	if edges[0].From != "SPEC-001" {
		t.Errorf("edges[0].From = %q, want %q", edges[0].From, "SPEC-001")
	}
}

// @test TEST-INT_GRPH-003
func TestLinkageGraph_GetNode(t *testing.T) {
	g := NewLinkageGraph()
	g.AddNode("SPEC-001", model.TypeSpec)

	n, ok := g.GetNode("SPEC-001")
	if !ok {
		t.Error("GetNode(SPEC-001) ok = false, want true")
	}
	if n.ID != "SPEC-001" {
		t.Errorf("n.ID = %q, want %q", n.ID, "SPEC-001")
	}

	_, ok = g.GetNode("NONEXISTENT")
	if ok {
		t.Error("GetNode(NONEXISTENT) ok = true, want false")
	}
}

// @test TEST-INT_GRPH-004
func TestLinkageGraph_NodeInOutEdges(t *testing.T) {
	g := NewLinkageGraph()
	g.AddNode("SPEC-001", model.TypeSpec)
	g.AddNode("TEST-001", model.TypeTest)
	g.AddNode("TEST-002", model.TypeTest)

	g.AddEdge("SPEC-001", "TEST-001", model.LinkTests, "test.go", 10)
	g.AddEdge("SPEC-001", "TEST-002", model.LinkTests, "test.go", 11)

	n, _ := g.GetNode("SPEC-001")
	outEdges := n.OutEdges()
	if len(outEdges) != 2 {
		t.Errorf("len(OutEdges()) = %d, want 2", len(outEdges))
	}

	n2, _ := g.GetNode("TEST-001")
	inEdges := n2.InEdges()
	if len(inEdges) != 1 {
		t.Errorf("len(InEdges()) for TEST-001 = %d, want 1", len(inEdges))
	}
}

// @test TEST-INT_GRPH-005
func TestLinkageGraph_GetOutboundByType(t *testing.T) {
	g := NewLinkageGraph()
	g.AddNode("SPEC-001", model.TypeSpec)
	g.AddNode("TEST-001", model.TypeTest)
	g.AddNode("CONTRACT-001", model.TypeContract)

	g.AddEdge("SPEC-001", "TEST-001", model.LinkTests, "test.go", 10)
	g.AddEdge("SPEC-001", "CONTRACT-001", model.LinkImplements, "test.go", 11)

	edges := g.GetOutboundByType("SPEC-001", model.LinkTests)
	if len(edges) != 1 {
		t.Errorf("len(GetOutboundByType(LinkTests)) = %d, want 1", len(edges))
	}

	edges = g.GetOutboundByType("SPEC-001", model.LinkImplements)
	if len(edges) != 1 {
		t.Errorf("len(GetOutboundByType(LinkImplements)) = %d, want 1", len(edges))
	}

	edges = g.GetOutboundByType("NONEXISTENT", model.LinkTests)
	if edges != nil {
		t.Errorf("GetOutboundByType for nonexistent should return nil")
	}
}

// @test TEST-INT_GRPH-006
func TestLinkageGraph_GetInboundByType(t *testing.T) {
	g := NewLinkageGraph()
	g.AddNode("SPEC-001", model.TypeSpec)
	g.AddNode("SPEC-002", model.TypeSpec)
	g.AddNode("TEST-001", model.TypeTest)

	g.AddEdge("SPEC-001", "TEST-001", model.LinkTests, "test.go", 10)
	g.AddEdge("SPEC-002", "TEST-001", model.LinkTests, "test.go", 11)

	edges := g.GetInboundByType("TEST-001", model.LinkTests)
	if len(edges) != 2 {
		t.Errorf("len(GetInboundByType) = %d, want 2", len(edges))
	}
}

// @test TEST-INT_GRPH-007
func TestLinkageGraph_GetBacklinks(t *testing.T) {
	g := NewLinkageGraph()
	g.AddNode("SPEC-001", model.TypeSpec)
	g.AddNode("TEST-001", model.TypeTest)

	g.AddEdge("SPEC-001", "TEST-001", model.LinkTests, "test.go", 10)

	bl := g.GetBacklinks("TEST-001")
	if len(bl) != 1 {
		t.Errorf("len(GetBacklinks) = %d, want 1", len(bl))
	}
	if bl[0] != "SPEC-001" {
		t.Errorf("GetBacklinks[0] = %q, want %q", bl[0], "SPEC-001")
	}
}

// @test TEST-INT_GRPH-008
func TestLinkageGraph_VerifyBidirectionalLinks(t *testing.T) {
	g := NewLinkageGraph()
	g.AddNode("SPEC-001", model.TypeSpec)
	g.AddNode("TEST-001", model.TypeTest)

	g.AddEdge("SPEC-001", "TEST-001", model.LinkTests, "test.go", 10)
	g.AddEdge("TEST-001", "SPEC-001", model.LinkImplements, "test.go", 11)

	g.VerifyBidirectionalLinks()

	edges := g.Edges()
	for _, e := range edges {
		if !e.Verified {
			t.Errorf("Edge %s->%s should be verified", e.From, e.To)
		}
	}
}

// @test TEST-INT_GRPH-009
func TestLinkageGraph_VerifyBidirectionalLinks_Unverified(t *testing.T) {
	g := NewLinkageGraph()
	g.AddNode("SPEC-001", model.TypeSpec)
	g.AddNode("TEST-001", model.TypeTest)

	// Only one direction
	g.AddEdge("SPEC-001", "TEST-001", model.LinkTests, "test.go", 10)

	g.VerifyBidirectionalLinks()

	edges := g.Edges()
	if edges[0].Verified {
		t.Error("One-way link should not be verified")
	}
}

// @test TEST-INT_GRPH-010
func TestLinkageGraph_ToSnapshot(t *testing.T) {
	g := NewLinkageGraph()
	g.AddNode("SPEC-001", model.TypeSpec)
	g.AddNode("TEST-001", model.TypeTest)
	g.AddEdge("SPEC-001", "TEST-001", model.LinkTests, "test.go", 10)

	snap := g.ToSnapshot()

	if len(snap.Nodes) != 2 {
		t.Errorf("len(Nodes) = %d, want 2", len(snap.Nodes))
	}
	if len(snap.Edges) != 1 {
		t.Errorf("len(Edges) = %d, want 1", len(snap.Edges))
	}
}

// @test TEST-INT_GRPH-011
func TestLinkageGraph_Stats(t *testing.T) {
	g := NewLinkageGraph()
	g.AddNode("SPEC-001", model.TypeSpec)
	g.AddNode("SPEC-002", model.TypeSpec)
	g.AddNode("TEST-001", model.TypeTest)
	g.AddNode("CONTRACT-001", model.TypeContract)

	stats := g.Stats()

	if stats.TotalIdentifiers != 4 {
		t.Errorf("TotalIdentifiers = %d, want 4", stats.TotalIdentifiers)
	}
	if stats.SpecsAnalyzed != 2 {
		t.Errorf("SpecsAnalyzed = %d, want 2", stats.SpecsAnalyzed)
	}
	if stats.TestsAnalyzed != 1 {
		t.Errorf("TestsAnalyzed = %d, want 1", stats.TestsAnalyzed)
	}
	if stats.ContractsAnalyzed != 1 {
		t.Errorf("ContractsAnalyzed = %d, want 1", stats.ContractsAnalyzed)
	}
}

// @test TEST-INT_GRPH-012
func TestLinkageGraph_ValidateCompleteness(t *testing.T) {
	g := NewLinkageGraph()
	g.AddNode("SPEC-001", model.TypeSpec)
	g.AddNode("TEST-001", model.TypeTest)

	g.AddEdge("SPEC-001", "TEST-001", model.LinkTests, "test.go", 10)
	g.AddEdge("TEST-001", "SPEC-001", model.LinkImplements, "test.go", 11)

	errs := g.ValidateCompleteness()
	if len(errs) != 0 {
		t.Errorf("ValidateCompleteness() returned %d errors, want 0", len(errs))
	}
}

// @test TEST-INT_GRPH-013
func TestLinkageGraph_ValidateCompleteness_NoLinks(t *testing.T) {
	g := NewLinkageGraph()
	g.AddNode("SPEC-001", model.TypeSpec)
	g.AddNode("TEST-001", model.TypeTest)

	errs := g.ValidateCompleteness()
	if len(errs) != 2 {
		t.Errorf("ValidateCompleteness() returned %d errors, want 2", len(errs))
	}
}

// @test TEST-INT_GRPH-014
func TestLinkageGraph_Nodes_Edges(t *testing.T) {
	g := NewLinkageGraph()
	g.AddNode("SPEC-001", model.TypeSpec)
	g.AddEdge("SPEC-001", "TEST-001", model.LinkTests, "test.go", 10)

	nodes := g.Nodes()
	if len(nodes) != 1 {
		t.Errorf("len(Nodes()) = %d, want 1", len(nodes))
	}

	edges := g.Edges()
	if len(edges) != 1 {
		t.Errorf("len(Edges()) = %d, want 1", len(edges))
	}
}
