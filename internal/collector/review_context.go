package collector

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/jingxu9x/idd-cli/internal/config"
	"github.com/jingxu9x/idd-cli/internal/graph"
	"github.com/jingxu9x/idd-cli/internal/model"
	"github.com/jingxu9x/idd-cli/pkg/pattern"
	"github.com/yuin/goldmark"
	goldmarkast "github.com/yuin/goldmark/ast"
	goldmarktext "github.com/yuin/goldmark/text"
)

const (
	reviewContextSchema            = "idd.spec_review_context.v2"
	reviewContextBatchSchema       = "idd.spec_review_context_batch.v2"
	reviewContextMaxSpecs          = 10
	reviewRecordCharacters         = 20000
	reviewExcerptCharacters        = 8000
	reviewContextExcerptCharacters = 64000
)

// SpecReviewContextBatch is an ordered collection of independently reviewable
// SPEC evidence bundles produced from one documentation and source scan.
// @implement SPEC-INTERNAL_COLLECTOR-026
type SpecReviewContextBatch struct {
	Schema           string               `json:"schema"`
	RequestedSpecIDs []string             `json:"requested_spec_ids"`
	Contexts         []*SpecReviewContext `json:"contexts"`
}

// SpecReviewContext is an evidence-only bundle for reviewing one SPEC.
// @implement SPEC-INTERNAL_COLLECTOR-026
type SpecReviewContext struct {
	Schema       string                     `json:"schema"`
	Spec         ReviewContextSpec          `json:"spec"`
	Components   []ReviewContextComponent   `json:"components"`
	Contracts    []ReviewContextContract    `json:"contracts"`
	Tests        []ReviewContextTest        `json:"tests"`
	Declarations []ReviewContextDeclaration `json:"declarations"`
	Issues       []ReviewContextIssue       `json:"issues,omitempty"`
	Truncated    bool                       `json:"truncated"`
}

// ReviewContextSpec preserves the selected authored SPEC and its location.
// @implement SPEC-INTERNAL_COLLECTOR-026
type ReviewContextSpec struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Package     string   `json:"package"`
	File        string   `json:"file"`
	Line        int      `json:"line"`
	Components  []string `json:"components"`
	Contracts   []string `json:"contracts"`
	Requirement string   `json:"requirement"`
	Acceptance  string   `json:"acceptance"`
	Markdown    string   `json:"markdown"`
	Truncated   bool     `json:"truncated"`
}

// ReviewContextComponent preserves the Component and its authored design
// boundary selected through the shared trace index.
type ReviewContextComponent struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Purpose   string `json:"purpose"`
	Ownership string `json:"ownership,omitempty"`
	Boundary  string `json:"boundary,omitempty"`
	Decisions string `json:"decisions,omitempty"`
	File      string `json:"file"`
	Line      int    `json:"line"`
	Markdown  string `json:"markdown"`
	Truncated bool   `json:"truncated"`
}

// ReviewContextContract preserves the guarantees referenced by the SPEC.
// @implement SPEC-INTERNAL_COLLECTOR-026
type ReviewContextContract struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Guarantees string `json:"guarantees"`
	File       string `json:"file"`
	Line       int    `json:"line"`
	Markdown   string `json:"markdown"`
	Truncated  bool   `json:"truncated"`
}

// ReviewContextTest preserves one authored TEST that covers the selected SPEC.
// @implement SPEC-INTERNAL_COLLECTOR-026
type ReviewContextTest struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Kind      string   `json:"kind"`
	Purpose   string   `json:"purpose"`
	Oracle    string   `json:"oracle"`
	Contracts []string `json:"contracts,omitempty"`
	File      string   `json:"file"`
	Line      int      `json:"line"`
	Markdown  string   `json:"markdown"`
	Truncated bool     `json:"truncated"`
}

// ReviewContextDeclaration preserves one AST-bound implementation declaration.
// @implement SPEC-INTERNAL_COLLECTOR-026
type ReviewContextDeclaration struct {
	Role       string   `json:"role"`
	References []string `json:"references"`
	Language   string   `json:"language"`
	Kind       string   `json:"kind"`
	Name       string   `json:"name"`
	File       string   `json:"file"`
	Line       int      `json:"line"`
	EndLine    int      `json:"end_line"`
	Annotation string   `json:"annotation"`
	Excerpt    string   `json:"excerpt"`
	Truncated  bool     `json:"truncated"`
}

// ReviewContextIssue identifies collection evidence that may be incomplete.
// @implement SPEC-INTERNAL_COLLECTOR-026
type ReviewContextIssue struct {
	Rule    string `json:"rule"`
	Message string `json:"message"`
	Source  string `json:"source,omitempty"`
}

// BuildSpecReviewContext assembles bounded evidence for an external semantic
// review. It does not score or approve the selected documentation.
// @implement SPEC-INTERNAL_COLLECTOR-026, SPEC-INTERNAL_COLLECTOR-028
func BuildSpecReviewContext(
	ctx context.Context,
	cfg *config.Config,
	docsPath string,
	sourcePath string,
	specID string,
) (*SpecReviewContext, error) {
	batch, err := BuildSpecReviewContexts(ctx, cfg, docsPath, sourcePath, []string{specID})
	if err != nil {
		return nil, err
	}
	return batch.Contexts[0], nil
}

// BuildSpecReviewContexts assembles an ordered bounded batch while sharing one
// documentation collection and one source collection across all requested IDs.
// @implement SPEC-INTERNAL_COLLECTOR-026, SPEC-INTERNAL_COLLECTOR-028
func BuildSpecReviewContexts(
	ctx context.Context,
	cfg *config.Config,
	docsPath string,
	sourcePath string,
	specIDs []string,
) (*SpecReviewContextBatch, error) {
	requested, err := normalizeReviewSpecIDs(specIDs)
	if err != nil {
		return nil, err
	}
	project, err := BuildTraceProject(ctx, cfg, docsPath, sourcePath)
	if err != nil {
		return nil, err
	}
	batch := &SpecReviewContextBatch{
		Schema:           reviewContextBatchSchema,
		RequestedSpecIDs: requested,
		Contexts:         make([]*SpecReviewContext, 0, len(requested)),
	}
	for _, specID := range requested {
		dossier, traceErr := project.Trace(specID, TraceOptions{Depth: 2})
		if traceErr != nil {
			return nil, traceErr
		}
		if dossier.Entity == nil || dossier.Entity.Resolution == string(graph.ResolutionUnresolved) {
			return nil, fmt.Errorf("SPEC %s was not found in a canonical spec.md record", specID)
		}
		if dossier.Entity.Resolution == string(graph.ResolutionAmbiguous) {
			locations := make([]string, 0, len(dossier.Owners))
			for _, owner := range dossier.Owners {
				locations = append(locations, fmt.Sprintf("%s:%d", owner.Path, owner.Line))
			}
			return nil, fmt.Errorf("SPEC %s has multiple document owners: %s", specID, strings.Join(locations, ", "))
		}
		contextValue := projectTraceReviewContext(dossier)
		sortReviewContext(contextValue)
		batch.Contexts = append(batch.Contexts, contextValue)
	}
	return batch, nil
}

func projectTraceReviewContext(dossier *TraceDossier) *SpecReviewContext {
	canonical := dossier.Canonical
	result := &SpecReviewContext{
		Schema:       reviewContextSchema,
		Components:   make([]ReviewContextComponent, 0, len(dossier.Components)),
		Contracts:    make([]ReviewContextContract, 0, len(dossier.Contracts)),
		Tests:        make([]ReviewContextTest, 0, len(dossier.Tests)),
		Declarations: make([]ReviewContextDeclaration, 0, len(dossier.Implementations)),
		Issues:       append([]ReviewContextIssue(nil), dossier.Findings...),
		Truncated:    dossier.Truncated,
	}
	if canonical != nil {
		result.Spec = ReviewContextSpec{
			ID: dossier.Query.ID, Title: canonical.Title, Package: canonical.Package,
			File: canonical.Path, Line: canonical.Line, Requirement: canonical.Requirement,
			Acceptance: canonical.Acceptance, Markdown: canonical.Markdown, Truncated: canonical.Truncated,
		}
	}
	for _, declaration := range dossier.Implementations {
		if containsString(declaration.References, dossier.Query.ID) {
			result.Declarations = append(result.Declarations, declaration)
		}
	}
	for _, relation := range dossier.OutboundRelations {
		switch relation.Type {
		case string(model.LinkDesignedBy), string(model.LinkReferences):
			if strings.HasPrefix(relation.To, "component:") {
				result.Spec.Components = append(result.Spec.Components, relation.To)
			}
		case string(model.LinkConstrainedBy), string(model.LinkContract):
			if strings.HasPrefix(relation.To, "contract:") {
				result.Spec.Contracts = append(result.Spec.Contracts, relation.To)
			}
		}
	}
	result.Spec.Components = uniqueSortedStrings(result.Spec.Components)
	result.Spec.Contracts = uniqueSortedStrings(result.Spec.Contracts)
	for _, component := range dossier.Components {
		result.Components = append(result.Components, ReviewContextComponent(component))
	}
	for _, contract := range dossier.Contracts {
		result.Contracts = append(result.Contracts, ReviewContextContract(contract))
	}
	for _, test := range dossier.Tests {
		if !containsString(test.Covers, dossier.Query.ID) {
			continue
		}
		result.Tests = append(result.Tests, ReviewContextTest{
			ID: test.ID, Title: test.Title, Kind: test.Kind, Purpose: test.Purpose,
			Oracle: test.Oracle, Contracts: append([]string(nil), test.Contracts...),
			File: test.File, Line: test.Line, Markdown: test.Markdown, Truncated: test.Truncated,
		})
		result.Declarations = append(result.Declarations, test.Declarations...)
	}
	return result
}

func uniqueSortedStrings(values []string) []string {
	seen := make(map[string]bool, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value != "" && !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	sort.Strings(result)
	return result
}

func normalizeReviewSpecIDs(specIDs []string) ([]string, error) {
	if len(specIDs) == 0 {
		return nil, fmt.Errorf("at least one SPEC identifier is required")
	}
	seen := make(map[string]bool, len(specIDs))
	requested := make([]string, 0, len(specIDs))
	for _, specID := range specIDs {
		if !strings.HasPrefix(specID, "SPEC-") || pattern.ValidateIdentifierFormat(specID) != nil {
			return nil, fmt.Errorf("%q is not a valid SPEC identifier", specID)
		}
		if seen[specID] {
			continue
		}
		seen[specID] = true
		requested = append(requested, specID)
	}
	if len(requested) > reviewContextMaxSpecs {
		return nil, fmt.Errorf(
			"review-context accepts at most %d unique SPEC identifiers; split this request into smaller focused batches",
			reviewContextMaxSpecs,
		)
	}
	return requested, nil
}

func addReviewDeclarationEvidence(
	result *SpecReviewContext,
	analyses []*SourceAnalysis,
	sourceLines map[string][]string,
	specID string,
	resolvePath func(string) string,
) error {
	remaining := reviewContextExcerptCharacters
	testIDs := make(map[string]bool, len(result.Tests))
	for _, test := range result.Tests {
		testIDs[test.ID] = true
	}
	for _, analysis := range analyses {
		lines, found := sourceLines[analysis.Path]
		if !found {
			path := analysis.Path
			if resolvePath != nil {
				path = resolvePath(path)
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("read declaration source %s: %w", analysis.Path, err)
			}
			lines = strings.Split(string(content), "\n")
			sourceLines[analysis.Path] = lines
		}
		declarations := make(map[string]SourceDeclaration)
		for _, declaration := range analysis.Declarations {
			declarations[sourceDeclarationKey(declaration.Name, declaration.Line)] = declaration
		}
		for _, annotation := range analysis.Annotations {
			role, references := reviewDeclarationRole(annotation, specID, testIDs)
			if role == "" {
				continue
			}
			declaration := declarations[sourceDeclarationKey(annotation.Declaration, annotation.DeclarationLine)]
			excerpt, truncated := boundedDeclarationExcerpt(lines, declaration, remaining)
			remaining -= len(excerpt)
			if remaining < 0 {
				remaining = 0
			}
			result.Truncated = result.Truncated || truncated
			result.Declarations = append(result.Declarations, ReviewContextDeclaration{
				Role:       role,
				References: references,
				Language:   analysis.Language,
				Kind:       declaration.Kind,
				Name:       declaration.Name,
				File:       analysis.Path,
				Line:       declaration.Line,
				EndLine:    declaration.EndLine,
				Annotation: annotation.Raw,
				Excerpt:    excerpt,
				Truncated:  truncated,
			})
		}
	}
	return nil
}

func reviewDeclarationRole(
	annotation SourceAnnotation,
	specID string,
	testIDs map[string]bool,
) (string, []string) {
	if annotation.Ignored || !annotation.Attached {
		return "", nil
	}
	if annotation.Kind == "implement" && containsString(annotation.Refs, specID) {
		return "implementation", []string{specID}
	}
	if annotation.Kind != "test" && annotation.Kind != "test-contract" {
		return "", nil
	}
	var references []string
	for _, reference := range annotation.Refs {
		if testIDs[reference] {
			references = append(references, reference)
		}
	}
	if len(references) == 0 {
		return "", nil
	}
	return "test", references
}

func boundedDeclarationExcerpt(
	lines []string,
	declaration SourceDeclaration,
	remaining int,
) (string, bool) {
	start := declaration.Line - 1
	if start < 0 {
		start = 0
	}
	end := declaration.EndLine
	if end > len(lines) {
		end = len(lines)
	}
	if end < start {
		end = start
	}
	excerpt := strings.Join(lines[start:end], "\n")
	limit := reviewExcerptCharacters
	if remaining < limit {
		limit = remaining
	}
	if limit < 0 {
		limit = 0
	}
	if len(excerpt) <= limit {
		return excerpt, false
	}
	return truncateReviewText(excerpt, limit), true
}

func boundedReviewRecordMarkdown(
	document *parsedIDDDocument,
	section string,
	recordIndex int,
) (string, bool) {
	lines := strings.Split(string(document.Raw), "\n")
	start := document.Index.recordLine(section, recordIndex) - 1
	if start < 0 {
		start = 0
	}
	end := len(lines)
	root := goldmark.DefaultParser().Parse(goldmarktext.NewReader(document.Body))
	recordLine := document.Index.recordLine(section, recordIndex)
	for child := root.FirstChild(); child != nil; child = child.NextSibling() {
		heading, ok := child.(*goldmarkast.Heading)
		if !ok || heading.Level != 2 {
			continue
		}
		line := markdownNodeLine(document, heading)
		if line > recordLine {
			end = line - 1
			break
		}
	}
	if end < start {
		end = start
	}
	markdown := strings.TrimSpace(strings.Join(lines[start:end], "\n"))
	if len(markdown) <= reviewRecordCharacters {
		return markdown, false
	}
	return truncateReviewText(markdown, reviewRecordCharacters), true
}

func truncateReviewText(value string, limit int) string {
	if limit >= len(value) {
		return value
	}
	if limit <= 0 {
		return ""
	}
	for limit > 0 && !utf8.RuneStart(value[limit]) {
		limit--
	}
	return value[:limit]
}

func validationErrorSourcePath(source string) string {
	index := strings.LastIndex(source, ":")
	if index == -1 {
		return source
	}
	if _, err := strconv.Atoi(source[index+1:]); err != nil {
		return source
	}
	return source[:index]
}

func sortReviewContext(result *SpecReviewContext) {
	sort.Slice(result.Components, func(i, j int) bool {
		if result.Components[i].ID != result.Components[j].ID {
			return result.Components[i].ID < result.Components[j].ID
		}
		if result.Components[i].File != result.Components[j].File {
			return result.Components[i].File < result.Components[j].File
		}
		return result.Components[i].Line < result.Components[j].Line
	})
	sort.Slice(result.Contracts, func(i, j int) bool {
		if result.Contracts[i].ID != result.Contracts[j].ID {
			return result.Contracts[i].ID < result.Contracts[j].ID
		}
		if result.Contracts[i].File != result.Contracts[j].File {
			return result.Contracts[i].File < result.Contracts[j].File
		}
		return result.Contracts[i].Line < result.Contracts[j].Line
	})
	sort.Slice(result.Tests, func(i, j int) bool {
		if result.Tests[i].ID != result.Tests[j].ID {
			return result.Tests[i].ID < result.Tests[j].ID
		}
		return result.Tests[i].File < result.Tests[j].File
	})
	sort.Slice(result.Declarations, func(i, j int) bool {
		if result.Declarations[i].File != result.Declarations[j].File {
			return result.Declarations[i].File < result.Declarations[j].File
		}
		if result.Declarations[i].Line != result.Declarations[j].Line {
			return result.Declarations[i].Line < result.Declarations[j].Line
		}
		return result.Declarations[i].Name < result.Declarations[j].Name
	})
	sort.Slice(result.Issues, func(i, j int) bool {
		if result.Issues[i].Source != result.Issues[j].Source {
			return result.Issues[i].Source < result.Issues[j].Source
		}
		if result.Issues[i].Rule != result.Issues[j].Rule {
			return result.Issues[i].Rule < result.Issues[j].Rule
		}
		return result.Issues[i].Message < result.Issues[j].Message
	})
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
