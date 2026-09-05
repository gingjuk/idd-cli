package collector

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jingxu9x/idd-cli/internal/config"
	"github.com/jingxu9x/idd-cli/internal/graph"
	"github.com/jingxu9x/idd-cli/internal/model"
	"github.com/jingxu9x/idd-cli/pkg/pattern"
)

const maxTraceDepth = 5

// TraceProject owns one shared documentation/source collection and its
// post-collection index. Multiple queries and review projections reuse it.
// @implement SPEC-INTERNAL_COLLECTOR-028
type TraceProject struct {
	cfg      *config.Config
	index    *graph.TraceIndex
	analyses []*SourceAnalysis
	findings []model.ValidationError
}

// BuildTraceProject collects the project once and builds the same Engine-owned
// index used by validation. It deliberately does not implement an independent
// Markdown or source linkage scanner.
// @implement SPEC-INTERNAL_COLLECTOR-028
func BuildTraceProject(
	ctx context.Context,
	cfg *config.Config,
	docsPath string,
	sourcePath string,
) (*TraceProject, error) {
	docCollector := NewDocCollector(cfg)
	docSet, docErrors, err := docCollector.Collect(ctx, docsPath)
	if err != nil {
		return nil, fmt.Errorf("collect trace documentation: %w", err)
	}
	codeCollector := NewCodeCollector(cfg)
	codeSet, codeErrors, err := codeCollector.CollectWithErrors(ctx, sourcePath)
	if err != nil {
		return nil, fmt.Errorf("collect trace source: %w", err)
	}
	docSet.Merge(codeSet)

	findings := make([]model.ValidationError, 0, len(docErrors)+len(codeErrors))
	for _, finding := range append(docErrors, codeErrors...) {
		if finding != nil {
			findings = append(findings, *finding)
		}
	}
	return &TraceProject{
		cfg:      cfg,
		index:    graph.BuildTraceIndexWithOccurrences(docSet, docCollector.Mentions()),
		analyses: codeCollector.Analyses(),
		findings: findings,
	}, nil
}

// Index exposes the stable read-only project index to other deterministic
// projections such as impacted queries.
func (p *TraceProject) Index() *graph.TraceIndex {
	return p.index
}

// Analyses returns the already-collected syntax-tree models. Change-aware
// projections use declaration ranges without rescanning source files.
func (p *TraceProject) Analyses() []*SourceAnalysis {
	return append([]*SourceAnalysis(nil), p.analyses...)
}

// Trace builds a bounded dossier without suppressing broken or incomplete
// state. Callers decide whether unresolved/ambiguous status is an exit error.
// @implement SPEC-INTERNAL_COLLECTOR-028
func (p *TraceProject) Trace(id string, options TraceOptions) (*TraceDossier, error) {
	if err := validateTraceQuery(id, options); err != nil {
		return nil, err
	}
	dossier := &TraceDossier{
		Schema:            TraceSchema,
		Query:             TraceQuery{ID: id, Depth: options.Depth, IncludeMentions: options.IncludeMentions},
		Status:            string(graph.ResolutionUnresolved),
		Owners:            []TraceOccurrence{},
		Occurrences:       []TraceOccurrence{},
		OutboundRelations: []TraceRelation{},
		InboundRelations:  []TraceRelation{},
		RelatedEntities:   []TraceRelatedEntity{},
		Components:        []TraceComponent{},
		Contracts:         []TraceContract{},
		Implementations:   []ReviewContextDeclaration{},
		Tests:             []TraceTest{},
		Mentions:          []TraceOccurrence{},
		Findings:          []ReviewContextIssue{},
	}
	entity, found := p.index.Entity(id)
	if !found {
		dossier.Findings = append(dossier.Findings, ReviewContextIssue{
			Rule: "unknown-id", Message: fmt.Sprintf("%s has no declaration, reference, evidence, or mention", id),
		})
		return dossier, nil
	}
	dossier.Entity = traceEntity(entity)
	dossier.Status = string(entity.Resolution)
	if entity.Resolution == graph.ResolutionResolved && entity.Lifecycle == "planned" {
		dossier.Status = "planned"
	}
	for _, owner := range entity.Owners {
		dossier.Owners = append(dossier.Owners, traceOccurrence(owner))
	}
	for _, occurrence := range p.index.Occurrences(id) {
		converted := traceOccurrence(occurrence)
		if occurrence.Kind == graph.OccurrenceMention {
			if options.IncludeMentions {
				dossier.Mentions = append(dossier.Mentions, converted)
			}
			continue
		}
		dossier.Occurrences = append(dossier.Occurrences, converted)
	}
	if entity.Owner != nil {
		canonical, err := p.traceCanonical(*entity.Owner, entity.Kind)
		if err != nil {
			return nil, err
		}
		dossier.Canonical = canonical
		dossier.Truncated = canonical.Truncated
	}

	reachable := p.addTraceRelations(dossier, id, options.Depth, options.IncludeMentions)
	if err := p.addTraceSemanticContext(dossier, id, reachable); err != nil {
		return nil, err
	}
	p.addTraceFindings(dossier, reachable)
	return dossier, nil
}

func validateTraceQuery(id string, options TraceOptions) error {
	if strings.TrimSpace(id) == "" || id != strings.TrimSpace(id) {
		return fmt.Errorf("trace requires one non-empty ID without surrounding whitespace")
	}
	if !strings.HasPrefix(id, "component:") && !strings.HasPrefix(id, "contract:") && pattern.ValidateIdentifierFormat(id) != nil {
		return fmt.Errorf("%q is not a valid IDD identifier or scoped Component/Contract ID", id)
	}
	if options.Depth < 0 || options.Depth > maxTraceDepth {
		return fmt.Errorf("trace depth must be between 0 and %d", maxTraceDepth)
	}
	return nil
}

func traceEntity(entity graph.Entity) *TraceEntity {
	return &TraceEntity{ID: entity.ID, Kind: entity.Kind, Namespace: entity.Namespace, Lifecycle: entity.Lifecycle, Resolution: string(entity.Resolution)}
}

func traceOccurrence(occurrence graph.Occurrence) TraceOccurrence {
	return TraceOccurrence{
		ID: occurrence.ID, EntityID: occurrence.EntityID, Kind: string(occurrence.Kind), Origin: occurrence.Origin,
		DocumentRole: occurrence.DocumentRole, Path: occurrence.Path, Line: occurrence.Line,
		Field: occurrence.Field, RecordID: occurrence.RecordID,
	}
}

func traceRelation(relation graph.Relation, depth int) TraceRelation {
	return TraceRelation{
		From: relation.From, To: relation.To, Type: string(relation.Type), Depth: depth,
		Provenance: TraceProvenance{
			OccurrenceID: relation.Provenance.OccurrenceID, Path: relation.Provenance.Path, Line: relation.Provenance.Line,
			Field: relation.Provenance.Field, OccurrenceKind: string(relation.Provenance.Kind),
			RecordID: relation.Provenance.RecordID,
		},
	}
}

func (p *TraceProject) addTraceRelations(dossier *TraceDossier, root string, maxDepth int, includeMentions bool) map[string]int {
	reachable := map[string]int{root: 0}
	frontier := []string{root}
	seenRelations := make(map[string]bool)
	for depth := 1; depth <= maxDepth && len(frontier) > 0; depth++ {
		next := make([]string, 0)
		for _, current := range frontier {
			for _, direction := range []struct {
				relations []graph.Relation
				outbound  bool
			}{{p.index.Outbound(current), true}, {p.index.Inbound(current), false}} {
				for _, relation := range direction.relations {
					if relation.Type == model.LinkMentions && !includeMentions {
						continue
					}
					key := traceRelationKey(relation)
					if !seenRelations[key] {
						seenRelations[key] = true
						converted := traceRelation(relation, depth)
						if direction.outbound {
							dossier.OutboundRelations = append(dossier.OutboundRelations, converted)
						} else {
							dossier.InboundRelations = append(dossier.InboundRelations, converted)
						}
					}
					neighbor := relation.From
					if direction.outbound {
						neighbor = relation.To
					}
					if _, exists := reachable[neighbor]; exists {
						continue
					}
					neighborEntity, exists := p.index.Entity(neighbor)
					if !exists {
						continue
					}
					reachable[neighbor] = depth
					next = append(next, neighbor)
					dossier.RelatedEntities = append(dossier.RelatedEntities, TraceRelatedEntity{
						ID: neighborEntity.ID, Kind: neighborEntity.Kind, Lifecycle: neighborEntity.Lifecycle,
						Resolution: string(neighborEntity.Resolution), Depth: depth,
					})
				}
			}
		}
		sort.Strings(next)
		frontier = next
	}
	return reachable
}

func traceRelationKey(relation graph.Relation) string {
	return strings.Join([]string{relation.From, relation.To, string(relation.Type), relation.Provenance.Path, fmt.Sprint(relation.Provenance.Line), relation.Provenance.Field, relation.Provenance.RecordID}, "\x00")
}

func (p *TraceProject) traceCanonical(owner graph.Occurrence, kind string) (*TraceCanonical, error) {
	document, err := p.readParsedIDDDocument(owner.Path)
	if err != nil {
		return nil, err
	}
	canonical := &TraceCanonical{TraceLocation: TraceLocation{Path: owner.Path, Line: owner.Line}, Package: document.Metadata.Package}
	switch kind {
	case "spec":
		for index, record := range document.Metadata.Specs {
			if record.ID != owner.EntityID || document.Index.recordLine("specs", index) != owner.Line {
				continue
			}
			canonical.Title, canonical.Requirement, canonical.Acceptance = record.Title, record.Requirement, record.Acceptance
			canonical.Markdown, canonical.Truncated = boundedReviewRecordMarkdown(document, "specs", index)
			return canonical, nil
		}
	case "test":
		for index, record := range document.Metadata.Tests {
			if record.ID != owner.EntityID || document.Index.recordLine("tests", index) != owner.Line {
				continue
			}
			canonical.Title, canonical.Purpose, canonical.Oracle, canonical.Kind = record.Title, record.Purpose, record.Oracle, record.Kind
			canonical.Markdown, canonical.Truncated = boundedReviewRecordMarkdown(document, "tests", index)
			return canonical, nil
		}
	case "component":
		for index, record := range document.Metadata.ComponentRecords {
			if record.Name != owner.Title || document.Index.recordLine("components", index) != owner.Line {
				continue
			}
			canonical.Title, canonical.Purpose = record.Name, record.Purpose
			canonical.Ownership, canonical.Boundary, canonical.Decisions = record.Ownership, record.Boundary, record.Decisions
			canonical.Markdown, canonical.Truncated = boundedReviewRecordMarkdown(document, "components", index)
			return canonical, nil
		}
	case "contract":
		for index, record := range document.Metadata.ContractRecords {
			if record.Name != owner.Title || document.Index.recordLine("contracts", index) != owner.Line {
				continue
			}
			canonical.Title, canonical.Guarantees = record.Name, record.Guarantees
			canonical.Markdown, canonical.Truncated = boundedReviewRecordMarkdown(document, "contracts", index)
			return canonical, nil
		}
	}
	return canonical, nil
}

func (p *TraceProject) addTraceSemanticContext(dossier *TraceDossier, root string, reachable map[string]int) error {
	ids := make([]string, 0, len(reachable))
	for id := range reachable {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		entity, found := p.index.Entity(id)
		if !found || entity.Owner == nil {
			continue
		}
		switch entity.Kind {
		case "component":
			canonical, err := p.traceCanonical(*entity.Owner, entity.Kind)
			if err != nil {
				return err
			}
			dossier.Components = append(dossier.Components, TraceComponent{ID: id, Name: canonical.Title, Purpose: canonical.Purpose, Ownership: canonical.Ownership, Boundary: canonical.Boundary, Decisions: canonical.Decisions, File: canonical.Path, Line: canonical.Line, Markdown: canonical.Markdown, Truncated: canonical.Truncated})
			dossier.Truncated = dossier.Truncated || canonical.Truncated
		case "contract":
			canonical, err := p.traceCanonical(*entity.Owner, entity.Kind)
			if err != nil {
				return err
			}
			dossier.Contracts = append(dossier.Contracts, TraceContract{ID: id, Name: canonical.Title, Guarantees: canonical.Guarantees, File: canonical.Path, Line: canonical.Line, Markdown: canonical.Markdown, Truncated: canonical.Truncated})
			dossier.Truncated = dossier.Truncated || canonical.Truncated
		case "test":
			test, err := p.traceTest(entity)
			if err != nil {
				return err
			}
			if test != nil {
				dossier.Tests = append(dossier.Tests, *test)
				dossier.Truncated = dossier.Truncated || test.Truncated
			}
		}
	}
	return p.addTraceDeclarations(dossier, root, reachable)
}

func (p *TraceProject) traceTest(entity graph.Entity) (*TraceTest, error) {
	canonical, err := p.traceCanonical(*entity.Owner, entity.Kind)
	if err != nil {
		return nil, err
	}
	document, err := p.readParsedIDDDocument(entity.Owner.Path)
	if err != nil {
		return nil, err
	}
	for _, record := range document.Metadata.Tests {
		if record.ID != entity.ID {
			continue
		}
		return &TraceTest{ID: entity.ID, Title: canonical.Title, Kind: record.Kind, Purpose: record.Purpose, Oracle: record.Oracle, Covers: append([]string(nil), record.Covers...), Contracts: append([]string(nil), record.Contracts...), File: canonical.Path, Line: canonical.Line, Markdown: canonical.Markdown, Declarations: []ReviewContextDeclaration{}, Truncated: canonical.Truncated}, nil
	}
	return nil, nil
}

func (p *TraceProject) readParsedIDDDocument(path string) (*parsedIDDDocument, error) {
	data, err := os.ReadFile(p.cfg.ResolvePath(path))
	if err != nil {
		return nil, fmt.Errorf("read trace document %s: %w", path, err)
	}
	parsed, detected, err := parseIDDDocument(path, data)
	if err != nil {
		return nil, fmt.Errorf("parse trace document %s: %w", path, err)
	}
	if !detected {
		return nil, fmt.Errorf("%s is not a self-describing IDD document", path)
	}
	return parsed, nil
}

func (p *TraceProject) addTraceDeclarations(dossier *TraceDossier, root string, reachable map[string]int) error {
	specIDs := make([]string, 0)
	for id := range reachable {
		entity, found := p.index.Entity(id)
		if found && entity.Kind == "spec" {
			specIDs = append(specIDs, id)
		}
	}
	if len(specIDs) == 0 {
		specIDs = append(specIDs, root)
	}
	sort.Strings(specIDs)
	sourceLines := make(map[string][]string)
	seen := make(map[string]bool)
	for _, specID := range specIDs {
		temporary := &SpecReviewContext{Spec: ReviewContextSpec{ID: specID}, Tests: []ReviewContextTest{}, Declarations: []ReviewContextDeclaration{}}
		for _, test := range dossier.Tests {
			temporary.Tests = append(temporary.Tests, ReviewContextTest{ID: test.ID})
		}
		if err := addReviewDeclarationEvidence(temporary, p.analyses, sourceLines, specID, p.cfg.ResolvePath); err != nil {
			return err
		}
		for _, declaration := range temporary.Declarations {
			key := strings.Join([]string{declaration.Role, declaration.File, fmt.Sprint(declaration.Line), declaration.Name, strings.Join(declaration.References, ",")}, "\x00")
			if seen[key] {
				continue
			}
			seen[key] = true
			if declaration.Role == "implementation" {
				dossier.Implementations = append(dossier.Implementations, declaration)
				continue
			}
			for index := range dossier.Tests {
				if intersects(declaration.References, []string{dossier.Tests[index].ID}) {
					dossier.Tests[index].Declarations = append(dossier.Tests[index].Declarations, declaration)
				}
			}
		}
	}
	return nil
}

func intersects(left, right []string) bool {
	for _, a := range left {
		for _, b := range right {
			if a == b {
				return true
			}
		}
	}
	return false
}

func (p *TraceProject) addTraceFindings(dossier *TraceDossier, reachable map[string]int) {
	paths := make(map[string]bool)
	for _, occurrence := range dossier.Occurrences {
		paths[occurrence.Path] = true
	}
	for _, finding := range p.findings {
		_, related := reachable[finding.Link]
		if !related && finding.Link != dossier.Query.ID && !paths[validationErrorSourcePath(finding.Source)] {
			continue
		}
		dossier.Findings = append(dossier.Findings, ReviewContextIssue{Rule: finding.Rule, Message: finding.Message, Source: finding.Source})
	}
	if dossier.Entity != nil {
		switch dossier.Entity.Resolution {
		case string(graph.ResolutionUnresolved):
			dossier.Findings = append(dossier.Findings, ReviewContextIssue{Rule: "undefined-reference", Message: fmt.Sprintf("%s has occurrences but no canonical declaration", dossier.Query.ID)})
		case string(graph.ResolutionAmbiguous):
			locations := make([]string, 0, len(dossier.Owners))
			for _, owner := range dossier.Owners {
				locations = append(locations, fmt.Sprintf("%s:%d", owner.Path, owner.Line))
			}
			dossier.Findings = append(dossier.Findings, ReviewContextIssue{Rule: "duplicate-owner", Message: fmt.Sprintf("%s has multiple canonical declarations: %s", dossier.Query.ID, strings.Join(locations, ", "))})
		}
	}
}

// SourceOccurrences returns indexed occurrences in one path and optional line
// interval. It is the stable bridge for change-aware projections.
func (p *TraceProject) SourceOccurrences(path string, firstLine, lastLine int) []graph.Occurrence {
	return p.index.OccurrencesIn(filepath.ToSlash(filepath.Clean(path)), firstLine, lastLine)
}
