package engine

// @spec SPEC-BE-004

import (
	"context"
	"fmt"
	"time"

	"github.com/jingxu9x/idd-link-validator/internal/config"
	"github.com/jingxu9x/idd-link-validator/internal/graph"
	"github.com/jingxu9x/idd-link-validator/internal/model"
	"github.com/jingxu9x/idd-link-validator/internal/similarity"
	"github.com/jingxu9x/idd-link-validator/pkg/pattern"
)

type Engine struct {
	cfg    *config.Config
	graph  *graph.LinkageGraph
	result *model.ValidationResult
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
	for _, id := range ids.AllIdentifiers() {
		node := e.graph.AddNode(id.ID, id.Type)
		if node.Metadata == nil {
			node.Metadata = make(map[string]interface{})
		}
		if _, exists := node.Metadata[string(id.Origin)]; !exists {
			node.Metadata[string(id.Origin)] = true
		}
		if id.Describe != "" {
			if _, exists := node.Metadata["describe_"+string(id.Origin)]; !exists {
				node.Metadata["describe_"+string(id.Origin)] = id.Describe
			}
		}
	}

	for _, id := range ids.AllIdentifiers() {
		for _, linkRef := range id.Links {
			linkType := e.inferLinkType(id.Type, linkRef)
			e.graph.AddEdge(id.ID, linkRef, linkType, id.Source, id.Line)
		}
	}

	e.graph.VerifyBidirectionalLinks()
}

func (e *Engine) inferLinkType(fromType model.IdentifierType, toRef string) model.LinkType {
	toType := model.TypeSpec
	if refType := pattern.GetIdentifierType(toRef); refType != "" {
		if t, err := model.ParseIdentifierType(refType); err == nil {
			toType = t
		}
	}

	switch fromType {
	case model.TypeSpec:
		if toType == model.TypeTest {
			return model.LinkTests
		}
		return model.LinkReferences
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

	if e.cfg.Validation.RequireDocCodeCorrespondence {
		e.validateDocCodeCorrespondence()
	}

	if e.cfg.Validation.ConsistencyCheck.Enabled {
		e.validateConsistency()
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

func (e *Engine) validateDocCodeCorrespondence() {
	for _, node := range e.graph.Nodes() {
		hasDoc, hasCode := false, false
		if node.Metadata != nil {
			hasDoc, _ = node.Metadata[string(model.OriginDoc)].(bool)
			hasCode, _ = node.Metadata[string(model.OriginCode)].(bool)
		}

		if hasDoc && !hasCode {
			e.result.AddError(
				"doc-code-correspondence",
				fmt.Sprintf("%s is documented but has no code annotation", node.ID),
				node.ID,
				"",
				"",
			)
		}

		if hasCode && !hasDoc {
			e.result.AddError(
				"doc-code-correspondence",
				fmt.Sprintf("%s has code annotation but no documentation", node.ID),
				node.ID,
				"",
				"",
			)
		}
	}
}

func (e *Engine) validateConsistency() {
	threshold := e.cfg.Validation.ConsistencyCheck.Threshold
	for _, node := range e.graph.Nodes() {
		if node.Metadata == nil {
			continue
		}
		hasDoc, _ := node.Metadata[string(model.OriginDoc)].(bool)
		hasCode, _ := node.Metadata[string(model.OriginCode)].(bool)
		if !hasDoc || !hasCode {
			continue
		}
		docDescribe, _ := node.Metadata["describe_"+string(model.OriginDoc)].(string)
		codeDescribe, _ := node.Metadata["describe_"+string(model.OriginCode)].(string)
		if docDescribe == "" || codeDescribe == "" {
			continue
		}
		score := similarity.Score(docDescribe, codeDescribe)
		if score < threshold {
			e.result.AddWarning(
				"consistency-check",
				fmt.Sprintf("%s has low similarity between doc and code (score: %.2f < threshold: %.2f)", node.ID, score, threshold),
				node.ID,
				fmt.Sprintf("doc: %s | code: %s", truncate(docDescribe, 50), truncate(codeDescribe, 50)),
				fmt.Sprintf("%.3f", score),
			)
		}
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func (e *Engine) BuildReport() *model.Report {
	return &model.Report{
		Tool:      "idd-cli",
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
