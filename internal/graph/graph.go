// Package graph provides graph data structures for IDD link validation.

// Spec: docs/internal/graph/spec.md
// Contract: docs/internal/graph/contract.md
package graph

import (
	"fmt"

	"github.com/jingxu9x/idd-cli/internal/model"
)

// Node represents an IDD identifier within the linkage graph with incoming and outgoing edges.
//
// @implement SPEC-INTERNAL_GRAPH-001
type Node struct {
	ID       string
	Type     model.IdentifierType
	Metadata map[string]interface{}
	inEdges  []*Edge
	outEdges []*Edge
}

// Edge represents a directed relationship between two nodes in the linkage graph.
//
// @implement SPEC-INTERNAL_GRAPH-002
type Edge struct {
	From     string
	To       string
	Type     model.LinkType
	Source   string
	Line     int
	Verified bool
}

// LinkageGraph manages nodes and edges for IDD identifier validation.
//
// @implement SPEC-INTERNAL_GRAPH-003
type LinkageGraph struct {
	nodes map[string]*Node
	edges []*Edge
	idx   *Index
}

// Index provides fast lookup structures for nodes by ID, type, and backlinks.
//
// @implement SPEC-INTERNAL_GRAPH-004
type Index struct {
	byID      map[string]*Node
	byType    map[model.IdentifierType][]*Node
	backlinks map[string][]string
}

// NewLinkageGraph creates a new empty linkage graph with initialized maps.
// @implement SPEC-INTERNAL_GRAPH-005, SPEC-INTERNAL_GRAPH-006, SPEC-INTERNAL_GRAPH-009
func NewLinkageGraph() *LinkageGraph {
	return &LinkageGraph{
		nodes: make(map[string]*Node),
		edges: make([]*Edge, 0),
		idx: &Index{
			byID:      make(map[string]*Node),
			byType:    make(map[model.IdentifierType][]*Node),
			backlinks: make(map[string][]string),
		},
	}
}

// AddNode adds a node to the graph if it doesn't exist.
//
// @implement SPEC-INTERNAL_GRAPH-007
func (g *LinkageGraph) AddNode(id string, idType model.IdentifierType) *Node {
	if n, ok := g.nodes[id]; ok {
		return n
	}
	n := &Node{
		ID:       id,
		Type:     idType,
		Metadata: make(map[string]interface{}),
		inEdges:  make([]*Edge, 0),
		outEdges: make([]*Edge, 0),
	}
	g.nodes[id] = n
	g.idx.byID[id] = n
	g.idx.byType[idType] = append(g.idx.byType[idType], n)
	return n
}

// AddEdge adds a directed edge between two nodes.
//
// @implement SPEC-INTERNAL_GRAPH-008
func (g *LinkageGraph) AddEdge(from, to string, edgeType model.LinkType, source string, line int) {
	edge := &Edge{
		From:   from,
		To:     to,
		Type:   edgeType,
		Source: source,
		Line:   line,
	}
	g.edges = append(g.edges, edge)

	if fromNode, ok := g.nodes[from]; ok {
		fromNode.outEdges = append(fromNode.outEdges, edge)
	}
	if toNode, ok := g.nodes[to]; ok {
		toNode.inEdges = append(toNode.inEdges, edge)
	}

	g.idx.backlinks[to] = append(g.idx.backlinks[to], from)
}

func (g *LinkageGraph) Nodes() map[string]*Node {
	return g.nodes
}

func (g *LinkageGraph) Edges() []*Edge {
	return g.edges
}

func (g *LinkageGraph) GetNode(id string) (*Node, bool) {
	n, ok := g.nodes[id]
	return n, ok
}

func (n *Node) InEdges() []*Edge {
	return n.inEdges
}

func (n *Node) OutEdges() []*Edge {
	return n.outEdges
}

func (g *LinkageGraph) GetOutboundByType(nodeID string, linkType model.LinkType) []*Edge {
	node, ok := g.nodes[nodeID]
	if !ok {
		return nil
	}
	var result []*Edge
	for _, e := range node.outEdges {
		if e.Type == linkType {
			result = append(result, e)
		}
	}
	return result
}

func (g *LinkageGraph) GetInboundByType(nodeID string, linkType model.LinkType) []*Edge {
	node, ok := g.nodes[nodeID]
	if !ok {
		return nil
	}
	var result []*Edge
	for _, e := range node.inEdges {
		if e.Type == linkType {
			result = append(result, e)
		}
	}
	return result
}

func (g *LinkageGraph) GetBacklinks(nodeID string) []string {
	if links, ok := g.idx.backlinks[nodeID]; ok {
		return links
	}
	return []string{}
}

func (g *LinkageGraph) NodeCount() int {
	return len(g.nodes)
}

func (g *LinkageGraph) EdgeCount() int {
	return len(g.edges)
}

func (g *LinkageGraph) VerifyBidirectionalLinks() {
	for _, edge := range g.edges {
		reverseType := model.ReverseLinkType(edge.Type)
		reverseEdges := g.GetInboundByType(edge.From, reverseType)
		edge.Verified = false
		for _, re := range reverseEdges {
			if re.From == edge.To {
				edge.Verified = true
				break
			}
		}
	}
}

func (g *LinkageGraph) ToSnapshot() *model.GraphSnapshot {
	nodes := make([]model.NodeSummary, 0, len(g.nodes))
	for _, n := range g.nodes {
		nodes = append(nodes, model.NodeSummary{
			ID:       n.ID,
			Type:     n.Type,
			Outbound: len(n.outEdges),
			Inbound:  len(n.inEdges),
		})
	}
	edges := make([]model.EdgeSummary, 0, len(g.edges))
	for _, e := range g.edges {
		edges = append(edges, model.EdgeSummary{
			From:     e.From,
			To:       e.To,
			Type:     string(e.Type),
			Verified: e.Verified,
			Source:   e.Source,
			Line:     e.Line,
		})
	}
	return &model.GraphSnapshot{Nodes: nodes, Edges: edges}
}

func (g *LinkageGraph) Stats() model.ValidationStats {
	stats := model.ValidationStats{
		TotalIdentifiers: g.NodeCount(),
		TotalLinks:       g.EdgeCount(),
	}
	for _, n := range g.nodes {
		switch n.Type {
		case model.TypeSpec:
			stats.SpecsAnalyzed++
		case model.TypeTest:
			stats.TestsAnalyzed++
		case model.TypeContract:
			stats.ContractsAnalyzed++
		case model.TypeDesign:
			stats.DesignsAnalyzed++
		}
	}
	return stats
}

func (g *LinkageGraph) ValidateCompleteness() []model.ValidationError {
	var errors []model.ValidationError
	for _, node := range g.nodes {
		if node.Type == model.TypeSpec {
			testLinks := g.GetOutboundByType(node.ID, model.LinkTests)
			if len(testLinks) == 0 {
				errors = append(errors, model.ValidationError{
					Rule:    "spec-missing-tests",
					Message: fmt.Sprintf("SPEC %s has no **Tests:** field", node.ID),
					Link:    node.ID,
				})
			}
		}
		if node.Type == model.TypeTest {
			specLinks := g.GetOutboundByType(node.ID, model.LinkImplements)
			if len(specLinks) == 0 {
				errors = append(errors, model.ValidationError{
					Rule:    "test-missing-coverage",
					Message: fmt.Sprintf("TEST %s has no **Spec Coverage:** field", node.ID),
					Link:    node.ID,
				})
			}
		}
	}
	return errors
}
