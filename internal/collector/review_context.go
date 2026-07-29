package collector

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/jingxu9x/idd-cli/internal/config"
	"github.com/jingxu9x/idd-cli/internal/model"
	"github.com/jingxu9x/idd-cli/pkg/pattern"
	"github.com/yuin/goldmark"
	goldmarkast "github.com/yuin/goldmark/ast"
	goldmarktext "github.com/yuin/goldmark/text"
)

const (
	reviewContextSchema            = "idd.spec_review_context.v1"
	reviewContextBatchSchema       = "idd.spec_review_context_batch.v1"
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
	Contract     *ReviewContextContract     `json:"contract,omitempty"`
	Tests        []ReviewContextTest        `json:"tests"`
	Declarations []ReviewContextDeclaration `json:"declarations"`
	Issues       []ReviewContextIssue       `json:"issues,omitempty"`
	Truncated    bool                       `json:"truncated"`
}

// ReviewContextSpec preserves the selected authored SPEC and its location.
// @implement SPEC-INTERNAL_COLLECTOR-026
type ReviewContextSpec struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Package     string `json:"package"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Design      string `json:"design"`
	Contract    string `json:"contract"`
	Requirement string `json:"requirement"`
	Acceptance  string `json:"acceptance"`
	Markdown    string `json:"markdown"`
	Truncated   bool   `json:"truncated"`
}

// ReviewContextContract preserves the guarantees referenced by the SPEC.
// @implement SPEC-INTERNAL_COLLECTOR-026
type ReviewContextContract struct {
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
// @implement SPEC-INTERNAL_COLLECTOR-026
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
// @implement SPEC-INTERNAL_COLLECTOR-026
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
	docCollector := NewDocCollector(cfg)
	docSet, docErrors, err := docCollector.Collect(ctx, docsPath)
	if err != nil {
		return nil, fmt.Errorf("collect documentation: %w", err)
	}

	batch := &SpecReviewContextBatch{
		Schema:           reviewContextBatchSchema,
		RequestedSpecIDs: requested,
		Contexts:         make([]*SpecReviewContext, 0, len(requested)),
	}
	for _, specID := range requested {
		contextValue, buildErr := buildReviewDocumentContext(docSet, docErrors, specID)
		if buildErr != nil {
			return nil, buildErr
		}
		batch.Contexts = append(batch.Contexts, contextValue)
	}

	codeCollector := NewCodeCollector(cfg)
	_, codeErrors, err := codeCollector.CollectWithErrors(ctx, sourcePath)
	if err != nil {
		return nil, fmt.Errorf("collect source: %w", err)
	}
	analyses := codeCollector.Analyses()
	sort.Slice(analyses, func(i, j int) bool { return analyses[i].Path < analyses[j].Path })
	sourceLines := make(map[string][]string)
	for _, contextValue := range batch.Contexts {
		specID := contextValue.Spec.ID
		contextValue.Issues = append(
			contextValue.Issues,
			reviewIssuesForSpecSource(codeErrors, specID)...,
		)
		if err := addReviewDeclarationEvidence(
			contextValue,
			analyses,
			sourceLines,
			specID,
		); err != nil {
			return nil, err
		}
		sortReviewContext(contextValue)
	}
	return batch, nil
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

func buildReviewDocumentContext(
	docSet *model.IdentifierSet,
	docErrors []*model.ValidationError,
	specID string,
) (*SpecReviewContext, error) {
	specDocument, specRecord, specIndex, err := resolveReviewSpec(docSet, specID)
	if err != nil {
		return nil, err
	}

	specMarkdown, specTruncated := boundedReviewRecordMarkdown(specDocument, "specs", specIndex)
	result := &SpecReviewContext{
		Schema: reviewContextSchema,
		Spec: ReviewContextSpec{
			ID:          specRecord.ID,
			Title:       specRecord.Title,
			Package:     specDocument.Metadata.Package,
			File:        specDocument.Path,
			Line:        specDocument.Index.recordLine("specs", specIndex),
			Design:      specRecord.Design,
			Contract:    specRecord.Contract,
			Requirement: specRecord.Requirement,
			Acceptance:  specRecord.Acceptance,
			Markdown:    specMarkdown,
			Truncated:   specTruncated,
		},
		Tests:        make([]ReviewContextTest, 0),
		Declarations: make([]ReviewContextDeclaration, 0),
		Issues:       reviewIssuesForDocumentPackage(docErrors, filepath.Dir(specDocument.Path)),
		Truncated:    specTruncated,
	}

	if err := addReviewDocumentEvidence(result, docSet, specDocument, specRecord); err != nil {
		return nil, err
	}
	return result, nil
}

func resolveReviewSpec(
	set *model.IdentifierSet,
	specID string,
) (*parsedIDDDocument, IDDDocumentSpec, int, error) {
	var owners []*model.Identifier
	for _, identifier := range set.GetAll(specID) {
		if identifier.Origin == model.OriginDoc && filepath.Base(identifier.Source) == "spec.md" {
			owners = append(owners, identifier)
		}
	}
	if len(owners) == 0 {
		return nil, IDDDocumentSpec{}, 0, fmt.Errorf("SPEC %s was not found in canonical spec.md records", specID)
	}
	if len(owners) > 1 {
		locations := make([]string, 0, len(owners))
		for _, owner := range owners {
			locations = append(locations, fmt.Sprintf("%s:%d", owner.Source, owner.Line))
		}
		sort.Strings(locations)
		return nil, IDDDocumentSpec{}, 0, fmt.Errorf(
			"SPEC %s has multiple document owners: %s",
			specID,
			strings.Join(locations, ", "),
		)
	}

	path := owners[0].Source
	parsed, err := readParsedIDDDocument(path)
	if err != nil {
		return nil, IDDDocumentSpec{}, 0, err
	}
	for index, spec := range parsed.Metadata.Specs {
		if spec.ID == specID {
			return parsed, spec, index, nil
		}
	}
	return nil, IDDDocumentSpec{}, 0, fmt.Errorf("SPEC %s was not found in %s", specID, path)
}

func addReviewDocumentEvidence(
	result *SpecReviewContext,
	docSet *model.IdentifierSet,
	specDocument *parsedIDDDocument,
	spec IDDDocumentSpec,
) error {
	if err := addReviewContractEvidence(result, docSet, specDocument, spec); err != nil {
		return err
	}
	return addReviewTestEvidence(result, docSet, spec.ID)
}

func addReviewContractEvidence(
	result *SpecReviewContext,
	docSet *model.IdentifierSet,
	specDocument *parsedIDDDocument,
	spec IDDDocumentSpec,
) error {
	targetPackage, targetName := splitIDDScopedName(
		specDocument.Metadata.Package,
		spec.Contract,
	)
	owners := append(
		[]*model.Identifier(nil),
		docSet.GetAll(derivedContractID(targetPackage, targetName))...,
	)
	sort.Slice(owners, func(i, j int) bool {
		if owners[i].Source != owners[j].Source {
			return owners[i].Source < owners[j].Source
		}
		return owners[i].Line < owners[j].Line
	})
	for _, owner := range owners {
		if owner.Origin != model.OriginDoc || owner.Kind != "contract" {
			continue
		}
		document, err := readParsedIDDDocument(owner.Source)
		if err != nil {
			return err
		}
		for index, contract := range document.Metadata.ContractRecords {
			line := document.Index.recordLine("contracts", index)
			if contract.Name != targetName || line != owner.Line {
				continue
			}
			markdown, truncated := boundedReviewRecordMarkdown(document, "contracts", index)
			result.Contract = &ReviewContextContract{
				Name:       contract.Name,
				Guarantees: contract.Guarantees,
				File:       document.Path,
				Line:       line,
				Markdown:   markdown,
				Truncated:  truncated,
			}
			result.Truncated = result.Truncated || truncated
			return nil
		}
	}
	return nil
}

func addReviewTestEvidence(
	result *SpecReviewContext,
	docSet *model.IdentifierSet,
	specID string,
) error {
	documents := make(map[string]*parsedIDDDocument)
	seen := make(map[string]bool)
	for _, identifier := range docSet.Tests {
		if identifier.Origin != model.OriginDoc ||
			!identifierHasTypedLink(identifier, specID, model.LinkImplements) {
			continue
		}
		document := documents[identifier.Source]
		if document == nil {
			var err error
			document, err = readParsedIDDDocument(identifier.Source)
			if err != nil {
				return err
			}
			documents[identifier.Source] = document
		}
		for index, test := range document.Metadata.Tests {
			line := document.Index.recordLine("tests", index)
			key := strings.Join(
				[]string{document.Path, strconv.Itoa(line), test.ID},
				"\x00",
			)
			if seen[key] || test.ID != identifier.ID || line != identifier.Line {
				continue
			}
			seen[key] = true
			markdown, truncated := boundedReviewRecordMarkdown(document, "tests", index)
			result.Tests = append(result.Tests, ReviewContextTest{
				ID:        test.ID,
				Title:     test.Title,
				Kind:      test.Kind,
				Purpose:   test.Purpose,
				Oracle:    test.Oracle,
				Contracts: append([]string(nil), test.Contracts...),
				File:      document.Path,
				Line:      line,
				Markdown:  markdown,
				Truncated: truncated,
			})
			result.Truncated = result.Truncated || truncated
		}
	}
	return nil
}

func identifierHasTypedLink(
	identifier *model.Identifier,
	target string,
	linkType model.LinkType,
) bool {
	for _, link := range identifier.TypedLinks {
		if link.Ref == target && link.Type == linkType {
			return true
		}
	}
	return false
}

func addReviewDeclarationEvidence(
	result *SpecReviewContext,
	analyses []*SourceAnalysis,
	sourceLines map[string][]string,
	specID string,
) error {
	remaining := reviewContextExcerptCharacters
	testIDs := make(map[string]bool, len(result.Tests))
	for _, test := range result.Tests {
		testIDs[test.ID] = true
	}
	for _, analysis := range analyses {
		lines, found := sourceLines[analysis.Path]
		if !found {
			content, err := os.ReadFile(analysis.Path)
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

func readParsedIDDDocument(path string) (*parsedIDDDocument, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read review document %s: %w", path, err)
	}
	parsed, detected, err := parseIDDDocument(path, data)
	if err != nil {
		return nil, fmt.Errorf("parse review document %s: %w", path, err)
	}
	if !detected {
		return nil, fmt.Errorf("%s is not a self-describing IDD document", path)
	}
	return parsed, nil
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

func reviewIssuesForDocumentPackage(
	errors []*model.ValidationError,
	directory string,
) []ReviewContextIssue {
	issues := make([]ReviewContextIssue, 0, len(errors))
	for _, validationError := range errors {
		sourcePath := validationErrorSourcePath(validationError.Source)
		if sourcePath != directory && filepath.Dir(sourcePath) != directory {
			continue
		}
		issues = append(issues, ReviewContextIssue{
			Rule:    validationError.Rule,
			Message: validationError.Message,
			Source:  validationError.Source,
		})
	}
	return issues
}

func reviewIssuesForSpecSource(
	errors []*model.ValidationError,
	specID string,
) []ReviewContextIssue {
	issues := make([]ReviewContextIssue, 0, len(errors))
	for _, validationError := range errors {
		sourcePath := validationErrorSourcePath(validationError.Source)
		content, err := os.ReadFile(sourcePath)
		if err != nil || !strings.Contains(string(content), specID) {
			continue
		}
		issues = append(issues, ReviewContextIssue{
			Rule:    validationError.Rule,
			Message: validationError.Message,
			Source:  validationError.Source,
		})
	}
	return issues
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
