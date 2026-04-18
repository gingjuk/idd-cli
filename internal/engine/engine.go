package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/yourorg/idd-link-validator/internal/config"
	"github.com/yourorg/idd-link-validator/internal/graph"
	"github.com/yourorg/idd-link-validator/internal/model"
	"github.com/yourorg/idd-link-validator/pkg/pattern"
)

type Engine struct {
	cfg      *config.Config
	graph    *graph.LinkageGraph
	result   *model.ValidationResult
}

func New(cfg *config.Config) *Engine {
	return &Engine{
		cfg:    cfg,
		graph:  graph.NewLinkageGraph(),
		result: model.NewValidationResult(),
	}
}

func (e *Engine) Run(ctx context.Context, ids *model.IdentifierSet) (*model.ValidationResult, error) {
	e.buildGraph(ids)
	e.result.Stats = e.graph.Stats()
	e.validate()
	if e.cfg.Output.IncludeGraph {
		e.result.Graph = e.graph.ToSnapshot()
	}
	return e.result, nil
}

func (e *Engine) AddStructuralErrors(errors []*model.ValidationError) {
	for _, err := range errors {
		e.result.AddError(err.Rule, err.Message, err.Source, err.Link, err.Code)
	}
}

func (e *Engine) buildGraph(ids *model.IdentifierSet) {
	for _, id := range ids.All() {
		e.graph.AddNode(id.ID, id.Type)
	}

	for _, id := range ids.All() {
		for _, linkRef := range id.Links {
			linkType := e.inferLinkType(id.Type, linkRef)
			e.graph.AddEdge(id.ID, linkRef, linkType, id.Source, id.Line)
		}
	}

	e.graph.VerifyBidirectionalLinks()
}

func (e *Engine) inferLinkType(fromType model.IdentifierType, toRef string) model.LinkType {
	switch fromType {
	case model.TypeSpec:
		return model.LinkTests
	case model.TypeTest:
		return model.LinkImplements
	case model.TypeContract:
		return model.LinkReferences
	case model.TypeDesign:
		return model.LinkReferences
	default:
		return model.LinkReferences
	}
}

func (e *Engine) validate() {
	e.result.Valid = true

	if e.cfg.Validation.RequireSpecTestCoverage {
		for _, err := range e.graph.ValidateCompleteness() {
			e.result.AddError(err.Rule, err.Message, err.Source, err.Link, err.Code)
		}
	}

	if e.cfg.Validation.RequireBidirectional {
		e.validateBidirectional()
	}

	if !e.cfg.Validation.AllowOrphans {
		e.validateNoOrphans()
	}

	e.result.Sort()
}

func (e *Engine) validateBidirectional() {
	for _, edge := range e.graph.Edges() {
		if !edge.Verified {
			e.result.AddError(
				"bidirectional-linkage",
				fmt.Sprintf("Link %s → %s has no backlink", edge.From, edge.To),
				fmt.Sprintf("%s:%d", edge.Source, edge.Line),
				fmt.Sprintf("%s→%s", edge.From, edge.To),
				"",
			)
		}
	}
}

func (e *Engine) validateNoOrphans() {
	for _, node := range e.graph.Nodes() {
		if len(node.InEdges()) == 0 && len(node.OutEdges()) == 0 {
			idType := pattern.GetIdentifierType(node.ID)
			if idType == "CONTRACT" || idType == "DESIGN" {
				e.result.AddWarning(
					"orphan-detection",
					fmt.Sprintf("%s is not referenced by any identifier", node.ID),
					node.ID,
					"",
					"",
				)
			} else {
				e.result.AddError(
					"orphan-detection",
					fmt.Sprintf("%s has no connections", node.ID),
					node.ID,
					"",
					"",
				)
			}
		}
	}
}

func (e *Engine) BuildReport() *model.Report {
	return &model.Report{
		Tool:      "idd-link-validator",
		Version:   "1.0.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Config: model.ConfigSummary{
			DocPatterns:  e.cfg.Docs.Patterns,
			CodePatterns: e.cfg.Code.Patterns,
			Annotations:  e.cfg.Code.Annotations,
		},
		Result: *e.result,
	}
}