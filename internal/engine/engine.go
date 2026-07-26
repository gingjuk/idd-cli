// Package engine provides the core validation engine.
package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/jingxu9x/idd-cli/internal/collector"
	"github.com/jingxu9x/idd-cli/internal/config"
	"github.com/jingxu9x/idd-cli/internal/graph"
	"github.com/jingxu9x/idd-cli/internal/model"
	"github.com/jingxu9x/idd-cli/internal/similarity"
	"github.com/jingxu9x/idd-cli/pkg/pattern"
)

// Engine is the core validation engine that orchestrates the validation pipeline.
// It coordinates collection, graph building, and validation rules to process
// identifiers through the complete validation workflow.
// The Engine holds the configuration, a linkage graph for tracking relationships
// between identifiers, and accumulates validation results.
// @implement SPEC-CMD_IDD_CLI-001, SPEC-CMD_IDD_CLI-002, SPEC-CMD_IDD_CLI-003, SPEC-CMD_IDD_CLI-004, SPEC-CMD_IDD_CLI-005, SPEC-CMD_IDD_CLI-006, SPEC-CMD_IDD_CLI-007, SPEC-INTERNAL_ENGINE-001
type Engine struct {
	cfg               *config.Config
	graph             *graph.LinkageGraph
	result            *model.ValidationResult
	sourceAnalyses    []*collector.SourceAnalysis
	sourceAnalysesSet bool
}

// New creates a new Engine instance with the given configuration.
// It initializes the engine with an empty linkage graph and validation result,
// preparing it to run validation against identifiers.
// @implement SPEC-INTERNAL_ENGINE-002
func New(cfg *config.Config) *Engine {
	return &Engine{
		cfg:    cfg,
		graph:  graph.NewLinkageGraph(),
		result: model.NewValidationResult(),
	}
}

// Run executes the validation pipeline for the given identifiers.
// It builds the linkage graph from the identifiers, runs all validation rules
// based on configuration, and returns the accumulated validation result.
// @implement SPEC-INTERNAL_ENGINE-004
func (e *Engine) Run(ctx context.Context, ids *model.IdentifierSet) (*model.ValidationResult, error) {
	e.validateDuplicateIDs(ids)
	e.buildGraph(ids)
	e.result.Stats = e.graph.Stats()
	e.validate()
	if e.cfg.Output.IncludeGraph {
		e.result.Graph = e.graph.ToSnapshot()
	}
	return e.result, nil
}

// AddStructuralErrors appends pre-built validation errors to the result.
// This allows callers to inject structural errors detected outside the
// normal validation pipeline.
// @implement SPEC-INTERNAL_ENGINE-001
func (e *Engine) AddStructuralErrors(errors []*model.ValidationError) {
	for _, err := range errors {
		if err == nil {
			continue
		}
		e.result.AddError(err.Rule, err.Message, err.Source, err.Link, err.Code)
	}
}

// SetSourceAnalyses provides the normalized syntax-tree results collected for
// the current run. Engine tests and library callers that omit this call are
// handled by a lazy syntax-tree collection from configured code patterns.
// @implement SPEC-INTERNAL_ENGINE-001
func (e *Engine) SetSourceAnalyses(analyses []*collector.SourceAnalysis) {
	e.sourceAnalyses = append([]*collector.SourceAnalysis(nil), analyses...)
	e.sourceAnalysesSet = true
}

// buildGraph constructs the linkage graph from the identifier set.
// It adds all identifiers as nodes, then creates edges based on link references
// between identifiers, inferring link types from the identifier types.
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
		if id.Kind != "" {
			if id.Origin == model.OriginDoc {
				node.Metadata["kind_doc"] = id.Kind
			} else {
				node.Metadata["kind_code_"+id.Kind] = true
			}
		}
		if id.Derived {
			node.Metadata["derived"] = true
			node.Metadata["display_name"] = id.Title
		}
		srcKey := "source_" + string(id.Origin)
		if _, exists := node.Metadata[srcKey]; !exists && id.Source != "" {
			node.Metadata[srcKey] = id.Source
		}
		if _, exists := node.Metadata["source_file"]; !exists && id.Source != "" {
			node.Metadata["source_file"] = id.Source
		}
	}

	for _, id := range ids.AllIdentifiers() {
		for _, linkRef := range id.Links {
			linkType := e.inferLinkType(id.Type, linkRef)
			e.graph.AddEdge(id.ID, linkRef, linkType, id.Source, id.Line)
		}
		for _, typedLink := range id.TypedLinks {
			e.graph.AddEdge(id.ID, typedLink.Ref, typedLink.Type, id.Source, id.Line)
		}
	}

	e.graph.VerifyBidirectionalLinks()
}

// inferLinkType determines the appropriate link type based on the source
// identifier type and target reference. It maps SPEC->TEST as LinkTests,
// SPEC->CONTRACT as LinkContract, TEST->SPEC as LinkImplements, etc.
func (e *Engine) inferLinkType(fromType model.IdentifierType, toRef string) model.LinkType {
	toType := model.TypeSpec
	if node, ok := e.graph.GetNode(toRef); ok {
		toType = node.Type
	} else if refType := pattern.GetIdentifierType(toRef); refType != "" {
		if t, err := model.ParseIdentifierType(refType); err == nil {
			toType = t
		}
	}

	switch fromType {
	case model.TypeSpec:
		if toType == model.TypeTest {
			return model.LinkTests
		}
		if toType == model.TypeContract {
			return model.LinkContract
		}
		return model.LinkReferences
	case model.TypeTest:
		if toType == model.TypeContract {
			return model.LinkContractTests
		}
		return model.LinkImplements
	case model.TypeContract:
		return model.LinkContractImplements
	case model.TypeDesign:
		return model.LinkReferences
	default:
		return model.LinkReferences
	}
}

// validate runs all enabled validation rules based on the engine configuration.
// It checks completeness, contract coverage, design sections, doc link consistency,
// orphan detection, doc-code correspondence, and various annotation requirements.
func (e *Engine) validate() {
	if e.cfg.Validation.RequireSpecTestCoverage {
		for _, err := range e.graph.ValidateCompleteness() {
			source := err.Source
			if source == "" && err.Link != "" {
				if node, ok := e.graph.GetNode(err.Link); ok {
					source = nodeSourceByOrigin(node, model.OriginDoc)
				}
			}
			e.result.AddError(err.Rule, err.Message, source, err.Link, err.Code)
		}
	}

	if e.cfg.Validation.RequireContractTestCoverage {
		e.validateContractTestCoverage()
	}

	if e.cfg.Validation.RequireDesignSections {
		e.validateDesignSections()
	}

	if e.cfg.Validation.RequireSpecFields {
		e.validateSpecRequiredFields()
	}
	e.validateDocumentTestKinds()
	e.validateGraphTargets()
	e.validateComponentDependencyCycles()
	e.validateContractDesignMarkers()
	e.validateDocPathExists()

	if e.cfg.Validation.RequireDocLinkConsistency {
		e.validateDocLinkConsistency()
		e.validateDuplicateHeadingIdentifiers()
	}

	if !e.cfg.Validation.AllowOrphans {
		e.validateNoOrphans()
	}

	if e.cfg.Validation.RequireDocCodeCorrespondence {
		e.validateDocCodeCorrespondence()
	}

	if e.cfg.Validation.RequirePublicFuncAnnotation {
		e.validateSourcePublicAnnotations()
	}

	if e.cfg.Validation.RequireRelatedFiles {
		e.validateRelatedFiles()
	}

	if e.cfg.Validation.RequirePkgDocFiles {
		e.validatePkgDocFiles()
	}

	if e.cfg.Validation.RequireTestAnnotation {
		e.validateSourceTestAnnotations()
	}

	if e.cfg.Validation.RequireAnnotationIdentifier {
		e.validateSourceAnnotationIdentifiers()
	}

	if e.cfg.Validation.RequireAnnotationOnSameLine {
		e.validateSourceConsecutiveAnnotations()
	}

	if e.cfg.Validation.ConsistencyCheck.Enabled {
		e.validateConsistency()
	}

	if len(e.result.Errors) == 0 {
		e.result.Valid = true
	}

	e.result.Sort()
}

// validateNoOrphans detects identifiers with no connections in the graph.
func (e *Engine) validateNoOrphans() {
	for _, node := range e.graph.Nodes() {
		if len(node.InEdges()) == 0 && len(node.OutEdges()) == 0 {
			if node.Type == model.TypeContract || node.Type == model.TypeDesign {
				e.result.AddWarning(
					"orphan-detection",
					fmt.Sprintf("%s is not referenced by any identifier", node.ID),
					nodeSource(node),
					"",
					"",
				)
			} else {
				e.result.AddError(
					"orphan-detection",
					fmt.Sprintf("%s has no connections", node.ID),
					nodeSource(node),
					"",
					"",
				)
			}
		}
	}
}

// validateContractTestCoverage checks that every CONTRACT identifier has at least
// one TEST identifier with @test-contract annotation referencing it.
func (e *Engine) validateContractTestCoverage() {
	for _, node := range e.graph.Nodes() {
		if node.Type != model.TypeContract {
			continue
		}
		if len(e.graph.GetInboundByType(node.ID, model.LinkContractTests)) == 0 {
			name := node.ID
			if displayName, ok := node.Metadata["display_name"].(string); ok && displayName != "" {
				name = displayName
			}
			e.result.AddError(
				"contract-test-coverage",
				fmt.Sprintf("Contract %s has no kind contract TEST relationship", name),
				nodeSource(node),
				node.ID,
				"contracts",
			)
		}
	}
}

func (e *Engine) validateGraphTargets() {
	for _, edge := range e.graph.Edges() {
		if !strings.HasPrefix(edge.To, "component:") &&
			!strings.HasPrefix(edge.To, "contract:") {
			continue
		}
		if _, exists := e.graph.GetNode(edge.To); exists {
			continue
		}
		e.result.AddError(
			"idd-document-reference",
			fmt.Sprintf("%s references missing graph target %s", edge.From, edge.To),
			fmt.Sprintf("%s:%d", edge.Source, edge.Line),
			edge.From,
			"reference",
		)
	}
}

func (e *Engine) validateComponentDependencyCycles() {
	state := make(map[string]uint8)
	var stack []string
	reported := make(map[string]bool)
	var visit func(string)
	visit = func(nodeID string) {
		state[nodeID] = 1
		stack = append(stack, nodeID)
		node, _ := e.graph.GetNode(nodeID)
		if node != nil {
			for _, edge := range node.OutEdges() {
				if edge.Type != model.LinkDependsOn {
					continue
				}
				target, ok := e.graph.GetNode(edge.To)
				if !ok || target.Type != model.TypeDesign {
					continue
				}
				switch state[target.ID] {
				case 0:
					visit(target.ID)
				case 1:
					if !reported[target.ID] {
						reported[target.ID] = true
						e.result.AddError(
							"component-dependency-cycle",
							fmt.Sprintf("Component dependency cycle detected: %s -> %s", strings.Join(stack, " -> "), target.ID),
							nodeSource(node),
							nodeID,
							"depends-on",
						)
					}
				}
			}
		}
		stack = stack[:len(stack)-1]
		state[nodeID] = 2
	}
	for nodeID, node := range e.graph.Nodes() {
		if node.Type == model.TypeDesign && state[nodeID] == 0 {
			visit(nodeID)
		}
	}
}

// validateDesignSections checks that design.md files contain all required
// sections: architecture, package layout, function composition, testability
// hooks, and dependencies.
func (e *Engine) validateDesignSections() {
	requiredSections := []string{
		"architecture",
		"package layout",
		"function composition",
		"testability hooks",
		"dependencies",
	}

	designRegex := regexp.MustCompile(`(?i)^#{2,}\s*(architecture|package layout|function composition|testability hooks|dependencies)\s*$`)

	e.walkDocFiles(func(path string, lines []string) {
		if !strings.HasSuffix(path, "design.md") {
			return
		}

		foundSections := make(map[string]int)
		for lineIndex, line := range lines {
			if matches := designRegex.FindStringSubmatch(strings.TrimSpace(line)); len(matches) > 1 {
				foundSections[strings.ToLower(matches[1])] = lineIndex
			}
		}

		for _, section := range requiredSections {
			lineIndex, found := foundSections[section]
			if !found {
				e.result.AddError(
					"design-sections",
					fmt.Sprintf("design.md missing required section: %s", section),
					path,
					"",
					section,
				)
				continue
			}

			if !hasIDDDocumentMetadata(lines) && !directoryHasIDDDocumentMetadata(filepath.Dir(path)) {
				continue
			}
			hasContent := false
			for nextIndex := lineIndex + 1; nextIndex < len(lines); nextIndex++ {
				nextLine := strings.TrimSpace(lines[nextIndex])
				if strings.HasPrefix(nextLine, "## ") {
					break
				}
				if nextLine != "" {
					hasContent = true
					break
				}
			}
			if !hasContent {
				e.result.AddError(
					"design-sections",
					fmt.Sprintf("design.md required section is empty: %s", section),
					fmt.Sprintf("%s:%d", path, lineIndex+1),
					"",
					section,
				)
			}
		}
	})
}

// validateDocLinkConsistency checks that docs use the correct link fields:
// SPEC docs should use **Tests:** not **Spec Coverage:**, and TEST docs
// should use **Spec Coverage:** not **Tests:**. It also verifies that all
// graph edges have corresponding backlinks.
func (e *Engine) validateDocLinkConsistency() {
	specCoverageRegex := regexp.MustCompile(`\*\*Spec Coverage:\*\*`)
	testsRegex := regexp.MustCompile(`\*\*Tests:\*\*`)

	e.walkDocFiles(func(path string, lines []string) {
		isSpec := strings.HasSuffix(path, "spec.md")
		isTest := strings.HasSuffix(path, "testing.md")

		for i, line := range lines {
			if isSpec && specCoverageRegex.MatchString(line) {
				e.result.AddError(
					"doc-link-consistency",
					"spec.md should use **Tests:** not **Spec Coverage:**",
					fmt.Sprintf("%s:%d", path, i+1),
					"",
					"",
				)
			}
			if isTest && testsRegex.MatchString(line) {
				e.result.AddError(
					"doc-link-consistency",
					"testing.md should use **Spec Coverage:** not **Tests:**",
					fmt.Sprintf("%s:%d", path, i+1),
					"",
					"",
				)
			}
		}
	})

	e.validateContractInterfaceConsistency()
}

// validateContractInterfaceConsistency checks that the interface name referenced
// in SPEC's **Contract:** field exists in the corresponding contract.md file.
// If the interface is not found, it emits a warning (not an error).
func (e *Engine) validateContractInterfaceConsistency() {
	contractFieldRegex := regexp.MustCompile(`(?i)\*\*Contract:\*\*\s*(.+)`)
	implementsInterfaceRegex := regexp.MustCompile(`implements\s+interface\s+([A-Za-z0-9_]+)`)

	contractInterfaces := make(map[string]map[string]bool)

	e.walkDocFiles(func(path string, lines []string) {
		if !strings.HasSuffix(path, "contract.md") {
			return
		}
		interfaces := make(map[string]bool)
		contentStr := strings.Join(lines, "\n")

		// Match ## Interface: <Name> or ### <Name> headings
		if match := regexp.MustCompile(`(?i)^#{2,3}\s+interface:\s*([A-Za-z][A-Za-z0-9_]*)`).FindStringSubmatch(contentStr); len(match) > 1 {
			interfaces[match[1]] = true
		}
		// Match backtick-quoted identifiers that look like interface names
		for _, line := range lines {
			if match := regexp.MustCompile("`([A-Z][a-zA-Z0-9_]*)`").FindStringSubmatch(line); len(match) > 1 {
				name := match[1]
				if name != "contract.md" && name != "design.md" && name != "spec.md" && name != "testing.md" && !strings.HasSuffix(name, ".md") {
					interfaces[name] = true
				}
			}
		}
		// Match func declarations in Go code blocks
		funcRegex := regexp.MustCompile(`(?m)^func\s+(?:\([^)]+\)\s+)?([A-Z][a-zA-Z0-9_]*)\s*\(`)
		for _, match := range funcRegex.FindAllStringSubmatch(contentStr, -1) {
			if len(match) > 1 {
				interfaces[match[1]] = true
			}
		}
		if len(interfaces) > 0 {
			contractInterfaces[path] = interfaces
		}
	})

	e.walkDocFiles(func(path string, lines []string) {
		if !strings.HasSuffix(path, "spec.md") {
			return
		}
		if hasIDDDocumentMetadata(lines) || directoryHasIDDDocumentMetadata(filepath.Dir(path)) {
			return
		}

		dir := filepath.Dir(path)
		contractPath := filepath.Join(dir, "contract.md")

		for i, line := range lines {
			if matches := contractFieldRegex.FindStringSubmatch(line); len(matches) > 1 {
				fieldValue := strings.TrimSpace(matches[1])

				if fieldValue == "" || fieldValue == "`contract.md`" || fieldValue == "contract.md" {
					continue
				}

				var interfaceName string
				if implMatches := implementsInterfaceRegex.FindStringSubmatch(fieldValue); len(implMatches) > 1 {
					interfaceName = implMatches[1]
				} else {
					backtickRegex := regexp.MustCompile("`([^`]+)`")
					if btMatches := backtickRegex.FindAllStringSubmatch(fieldValue, -1); len(btMatches) > 0 {
						for _, m := range btMatches {
							content := m[1]
							if content != "contract.md" && !strings.Contains(content, ".md") {
								interfaceName = content
							}
						}
					}
				}

				if interfaceName == "" {
					continue
				}

				if contractIdents, ok := contractInterfaces[contractPath]; ok {
					if !contractIdents[interfaceName] {
						// Fallback: check if interfaceName appears anywhere as capitalized identifier in contract.md
						found := false
						if contractContent, err := os.ReadFile(contractPath); err == nil {
							fallbackRegex := regexp.MustCompile(`\b([A-Z][a-zA-Z0-9_]*)\b`)
							for _, m := range fallbackRegex.FindAllStringSubmatch(string(contractContent), -1) {
								if len(m) > 1 && m[1] == interfaceName {
									found = true
									break
								}
							}
						}
						if !found {
							e.result.AddWarning(
								"contract-interface-consistency",
								fmt.Sprintf("SPEC references interface '%s' in **Contract:** but it was not found in %s", interfaceName, contractPath),
								fmt.Sprintf("%s:%d", path, i+1),
								interfaceName,
								"",
							)
						}
					}
				}
			}
		}
	})
}

type specRequiredField struct {
	name       string
	definition string
	regex      *regexp.Regexp
}

var specRequiredFields = []specRequiredField{
	{
		name:       "Design",
		definition: "**Design:** implements architecture `<ComponentName>`",
		regex:      regexp.MustCompile(`(?i)^\s*\*\*Design:\*\*`),
	},
	{
		name:       "Contract",
		definition: "**Contract:** implements interface `<InterfaceName>`",
		regex:      regexp.MustCompile(`(?i)^\s*\*\*Contract:\*\*`),
	},
	{
		name:       "Requirement",
		definition: "**Requirement:** [What this spec describes]",
		regex:      regexp.MustCompile(`(?i)^\s*\*\*Requirement:\*\*`),
	},
	{
		name:       "Tests",
		definition: "**Tests:** `TEST-<MODULE>-001`",
		regex:      regexp.MustCompile(`(?i)^\s*\*\*Tests:\*\*`),
	},
}

// validateSpecRequiredFields checks that every SPEC section in spec.md contains
// the fields required by the IDD spec template.
func (e *Engine) validateSpecRequiredFields() {
	headingRegex := regexp.MustCompile(`(?i)^#{2}\s+(SPEC-[A-Z0-9_]+-[0-9]+)(?::\s*.*)?$`)
	nextSectionRegex := regexp.MustCompile(`^#{1,2}\s+`)
	fieldRegex := regexp.MustCompile(`^\s*\*\*([^*:]+):\*\*`)
	expectedOrder := specRequiredFieldOrder()

	e.walkDocFiles(func(path string, lines []string) {
		if !strings.HasSuffix(path, "spec.md") {
			return
		}
		if hasIDDDocumentMetadata(lines) || directoryHasIDDDocumentMetadata(filepath.Dir(path)) {
			return
		}

		markerIDs := specFrontmatterMarkerIDs(lines)
		for i := 0; i < len(lines); i++ {
			trimmed := strings.TrimSpace(lines[i])
			matches := headingRegex.FindStringSubmatch(trimmed)
			if len(matches) <= 1 {
				continue
			}

			specID := matches[1]
			if !markerIDs[specID] {
				continue
			}
			found := make(map[string]bool, len(specRequiredFields))
			foundOrder := make([]string, 0, len(specRequiredFields))
			optionalFieldOrderError := false

			for j := i + 1; j < len(lines); j++ {
				nextTrimmed := strings.TrimSpace(lines[j])
				if nextSectionRegex.MatchString(nextTrimmed) {
					break
				}
				requiredField := false
				for _, field := range specRequiredFields {
					if field.regex.MatchString(nextTrimmed) {
						if !found[field.name] {
							foundOrder = append(foundOrder, field.name)
							found[field.name] = true
						}
						requiredField = true
						break
					}
				}
				if requiredField || len(foundOrder) == len(specRequiredFields) || optionalFieldOrderError {
					continue
				}
				fieldMatches := fieldRegex.FindStringSubmatch(nextTrimmed)
				if len(fieldMatches) <= 1 {
					continue
				}
				optionalFieldOrderError = true
				e.result.AddError(
					"spec-required-fields",
					fmt.Sprintf("SPEC %s optional field **%s:** appears before required fields are complete; optional fields must come after required fields: %s", specID, fieldMatches[1], expectedOrder),
					fmt.Sprintf("%s:%d", path, j+1),
					specID,
					"optional-field-order",
				)
			}

			missingField := false
			for _, field := range specRequiredFields {
				if found[field.name] {
					continue
				}
				missingField = true
				e.result.AddError(
					"spec-required-fields",
					fmt.Sprintf("SPEC %s missing required field **%s:**; define it as: %s", specID, field.name, field.definition),
					fmt.Sprintf("%s:%d", path, i+1),
					specID,
					field.name,
				)
			}
			if missingField || specRequiredFieldsInOrder(foundOrder) {
				continue
			}
			e.result.AddError(
				"spec-required-fields",
				fmt.Sprintf("SPEC %s required fields are out of order; expected order: %s", specID, expectedOrder),
				fmt.Sprintf("%s:%d", path, i+1),
				specID,
				"field-order",
			)
		}
	})
}

func specRequiredFieldsInOrder(foundOrder []string) bool {
	if len(foundOrder) != len(specRequiredFields) {
		return false
	}
	for i, field := range specRequiredFields {
		if foundOrder[i] != field.name {
			return false
		}
	}
	return true
}

// validateDocumentTestKinds checks that a self-described TEST's kind agrees
// with every source annotation for that identifier. Legacy TESTs have no kind
// metadata and remain on the existing validation path.
func (e *Engine) validateDocumentTestKinds() {
	for _, node := range e.graph.Nodes() {
		if node.Type != model.TypeTest {
			continue
		}
		documentedKind, ok := node.Metadata["kind_doc"].(string)
		if !ok || documentedKind == "" || node.Metadata[string(model.OriginCode)] != true {
			continue
		}

		hasTestAnnotation := node.Metadata["kind_code_test"] == true
		hasContractAnnotation := node.Metadata["kind_code_contract"] == true
		matches := documentedKind == "test" && hasTestAnnotation && !hasContractAnnotation ||
			documentedKind == "contract" && hasContractAnnotation && !hasTestAnnotation
		if matches {
			continue
		}

		expectedAnnotation := "@test"
		if documentedKind == "contract" {
			expectedAnnotation = "@test-contract"
		}
		e.result.AddError(
			"idd-document-test-kind",
			fmt.Sprintf(
				"self-described TEST %s has kind %q but its code annotations do not exclusively use %s",
				node.ID,
				documentedKind,
				expectedAnnotation,
			),
			nodeSourceByOrigin(node, model.OriginCode),
			node.ID,
			expectedAnnotation,
		)
	}
}

func specRequiredFieldOrder() string {
	parts := make([]string, 0, len(specRequiredFields))
	for _, field := range specRequiredFields {
		parts = append(parts, fmt.Sprintf("**%s:**", field.name))
	}
	return strings.Join(parts, ", ")
}

func specFrontmatterMarkerIDs(lines []string) map[string]bool {
	markerIDs := make(map[string]bool)
	idRegex := regexp.MustCompile(`^\s*-\s+id:\s*(SPEC-[A-Z0-9_]+-[0-9]+)\s*$`)
	inFrontmatter := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			break
		}
		if !inFrontmatter {
			continue
		}
		if matches := idRegex.FindStringSubmatch(line); len(matches) > 1 {
			markerIDs[matches[1]] = true
		}
	}

	return markerIDs
}

// validateDuplicateHeadingIdentifiers checks that H2 headings within the same
// file do not have duplicate identifier names. This catches cases where multiple
// sections (e.g., "## Test" and "## Contract Test") reuse the same TEST identifiers.
func (e *Engine) validateDuplicateHeadingIdentifiers() {
	// Match H2 headings with identifiers like ## TEST-FOO-001: Title
	headingRegex := regexp.MustCompile(`(?i)^#{2}\s+(TEST-[A-Z0-9_]+-[0-9]+|SPEC-[A-Z0-9_]+-[0-9]+|CONTRACT-[A-Z0-9_]+-[0-9]+):\s*`)

	e.walkDocFiles(func(path string, lines []string) {
		seen := make(map[string][]int) // identifier -> line numbers
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if matches := headingRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
				id := matches[1]
				seen[id] = append(seen[id], i+1)
			}
		}

		for id, lines := range seen {
			if len(lines) > 1 {
				e.result.AddError(
					"duplicate-heading-identifier",
					fmt.Sprintf("%s appears %d times in %s (lines: %v)", id, len(lines), path, lines),
					fmt.Sprintf("%s:%d", path, lines[0]),
					"",
					"",
				)
			}
		}
	})
}

// validateContractDesignMarkers checks that contract.md and design.md files
// do NOT have markers in their frontmatter - they should only have related_files.
func (e *Engine) validateContractDesignMarkers() {
	e.walkDocFiles(func(path string, lines []string) {
		if !strings.HasSuffix(path, "contract.md") && !strings.HasSuffix(path, "design.md") {
			return
		}

		// Parse frontmatter markers
		startIdx, endIdx := -1, -1
		inFrontmatter := false
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "---" {
				if !inFrontmatter {
					startIdx = i
					inFrontmatter = true
				} else {
					endIdx = i
					break
				}
			}
		}

		if startIdx == -1 || endIdx == -1 {
			return // No frontmatter, that's fine
		}

		// Check if markers field exists in frontmatter
		for i := startIdx; i < endIdx; i++ {
			trimmed := strings.TrimSpace(lines[i])
			if strings.HasPrefix(trimmed, "markers:") {
				e.result.AddError(
					"frontmatter-markers",
					fmt.Sprintf("%s should not have a markers field in frontmatter", filepath.Base(path)),
					path,
					"",
					"",
				)
				return
			}
		}
	})
}

// validateDocPathExists checks that subdirectories under docs/ correspond to
// real packages in the codebase. If a docs subdirectory doesn't match any
// package path, emit a warning.
func (e *Engine) validateDocPathExists() {
	docsRoot := "docs"

	entries, err := os.ReadDir(docsRoot)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Skip special directories like backend, core, contract, etc.
		docDir := filepath.Join(docsRoot, entry.Name())
		pkgPath := entry.Name()

		// Skip root-level doc directories that don't correspond to packages
		// (backend, core, contract are reference/example docs)
		if entry.Name() == "backend" || entry.Name() == "core" || entry.Name() == "contract" {
			continue
		}

		// Check if this doc directory corresponds to a real package
		// e.g., docs/internal/auth -> internal/auth
		fullPkgPath := pkgPath
		if entry.Name() == "internal" || entry.Name() == "pkg" {
			// For internal/ and pkg/, check subdirectories
			subEntries, _ := os.ReadDir(docDir)
			for _, subEntry := range subEntries {
				if !subEntry.IsDir() {
					continue
				}
				if entry.Name() == "internal" {
					fullPkgPath = filepath.Join("internal", subEntry.Name())
				} else if entry.Name() == "pkg" {
					fullPkgPath = filepath.Join("pkg", subEntry.Name())
				}

				// Check if package exists
				if !e.packageExists(fullPkgPath) && !e.isIgnoredDocPath(docDir) {
					e.result.AddWarning(
						"doc-path-missing",
						fmt.Sprintf("docs directory '%s' has no corresponding package '%s'. Add to ignore_paths if this is user documentation.", docDir, fullPkgPath),
						docDir,
						"",
						"",
					)
				}
			}
		} else {
			// Check if package exists at docs/<name> -> <name>
			if !e.packageExists(pkgPath) && !e.isIgnoredDocPath(docDir) {
				e.result.AddWarning(
					"doc-path-missing",
					fmt.Sprintf("docs directory '%s' has no corresponding package '%s'. Add to ignore_paths if this is user documentation.", docDir, pkgPath),
					docDir,
					"",
					"",
				)
			}
		}
	}
}

// packageExists checks whether a package directory exists at the given path.
func (e *Engine) packageExists(pkgPath string) bool {
	// Check if the package path exists
	info, err := os.Stat(pkgPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// isIgnoredDocPath checks whether a doc path should be ignored based on
// configured ignore patterns.
func (e *Engine) isIgnoredDocPath(docPath string) bool {
	for _, ignore := range e.cfg.Docs.IgnorePaths {
		if matched, _ := filepath.Match(ignore, docPath); matched {
			return true
		}
		// Handle ** glob pattern: docs/api/** matches docs/api and docs/api/...
		if strings.Contains(ignore, "**") {
			prefix := strings.TrimSuffix(ignore, "/**")
			prefix = strings.TrimSuffix(prefix, "**")
			if strings.HasPrefix(docPath, prefix+"/") || docPath == prefix {
				return true
			}
			continue
		}
		if strings.Contains(docPath, ignore) {
			return true
		}
	}
	return false
}

// validateDocCodeCorrespondence checks that identifiers have both documentation
// and code annotations. It reports errors when an identifier is documented
// but has no code annotation, or vice versa.
func (e *Engine) validateDocCodeCorrespondence() {
	for _, node := range e.graph.Nodes() {
		if node.Metadata["derived"] == true {
			continue
		}
		hasDoc, hasCode := false, false
		if node.Metadata != nil {
			hasDoc, _ = node.Metadata[string(model.OriginDoc)].(bool)
			hasCode, _ = node.Metadata[string(model.OriginCode)].(bool)
		}

		if hasDoc && !hasCode {
			e.result.AddError(
				"doc-code-correspondence",
				fmt.Sprintf("%s is documented but has no matching source annotation", node.ID),
				nodeSourceByOrigin(node, model.OriginDoc),
				"",
				"",
			)
		}

		if hasCode && !hasDoc {
			e.result.AddError(
				"doc-code-correspondence",
				fmt.Sprintf("%s has a source annotation but no matching documentation record", node.ID),
				nodeSourceByOrigin(node, model.OriginCode),
				"",
				"",
			)
		}
	}
}

// validateConsistency checks that the describe fields in documentation and code
// annotations have sufficient similarity, using configurable threshold. Low
// similarity indicates the doc and code descriptions may be out of sync.
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
		if docDescribe == "" || codeDescribe == "" || isFunctionLocator(codeDescribe) {
			continue
		}
		score := similarity.Score(docDescribe, codeDescribe)
		if score < threshold {
			e.result.AddWarning(
				"consistency-check",
				fmt.Sprintf("%s has low similarity between doc and code (score: %.2f < threshold: %.2f)", node.ID, score, threshold),
				nodeSource(node),
				fmt.Sprintf("doc: %s | code: %s", truncate(docDescribe, 50), truncate(codeDescribe, 50)),
				fmt.Sprintf("%.3f", score),
			)
		}
	}
}

func isFunctionLocator(description string) bool {
	description = strings.TrimSpace(description)
	return strings.HasPrefix(description, "[function: ") && strings.HasSuffix(description, "]")
}

// nodeSource returns the source file path for a graph node from metadata,
// falling back to the node ID when no source was recorded.
func nodeSource(node *graph.Node) string {
	if node.Metadata != nil {
		if src, ok := node.Metadata["source_file"].(string); ok && src != "" {
			return src
		}
	}
	return node.ID
}

// nodeSourceByOrigin returns the source file for the specific origin (doc/code),
// falling back to nodeSource when not available.
func nodeSourceByOrigin(node *graph.Node, origin model.Origin) string {
	if node.Metadata != nil {
		key := "source_" + string(origin)
		if src, ok := node.Metadata[key].(string); ok && src != "" {
			return src
		}
	}
	return nodeSource(node)
}

// truncate truncates a string to the specified maximum length, appending
// "..." if the string was shortened.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// validateRelatedFiles checks that documentation files have the required
// related_files field in their frontmatter.
func (e *Engine) validateRelatedFiles() {
	frontmatterRegex := regexp.MustCompile(`^related_files:`)

	e.walkDocFiles(func(path string, lines []string) {
		if filepath.Base(filepath.Dir(filepath.Clean(path))) == "docs" {
			return
		}
		if hasIDDDocumentMetadata(lines) || directoryHasIDDDocumentMetadata(filepath.Dir(path)) {
			return
		}

		startIdx, endIdx := -1, -1
		inFrontmatter := false
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "---" {
				if !inFrontmatter {
					startIdx = i
					inFrontmatter = true
				} else {
					endIdx = i
					break
				}
			}
		}

		if startIdx == -1 || endIdx == -1 {
			e.result.AddError(
				"related-files",
				"document missing frontmatter",
				path, "", "",
			)
			return
		}

		hasRelatedFiles := false
		for i := startIdx; i < endIdx; i++ {
			if frontmatterRegex.MatchString(strings.TrimSpace(lines[i])) {
				hasRelatedFiles = true
				break
			}
		}

		if !hasRelatedFiles {
			e.result.AddError(
				"related-files",
				"frontmatter missing related_files",
				path, "", "",
			)
		}
	})
}

func hasIDDDocumentMetadata(lines []string) bool {
	inFrontmatter := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" && !inFrontmatter {
			continue
		}
		if trimmed == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			return false
		}
		if !inFrontmatter {
			return false
		}
		if strings.TrimLeft(line, " \t") == line && strings.HasPrefix(trimmed, "idd:") {
			return true
		}
	}
	return false
}

func directoryHasIDDDocumentMetadata(directory string) bool {
	for _, filename := range []string{"design.md", "contract.md", "spec.md", "testing.md"} {
		data, err := os.ReadFile(filepath.Join(directory, filename))
		if err != nil {
			continue
		}
		if hasIDDDocumentMetadata(strings.Split(string(data), "\n")) {
			return true
		}
	}
	return false
}

// validatePkgDocFiles maps every scanned source package directory, including
// nested packages, to docs/<package>/ and requires the four canonical files.
func (e *Engine) validatePkgDocFiles() {
	required := []string{"design.md", "contract.md", "spec.md", "testing.md"}
	docsRoots := configuredDocsRoots(e.cfg.Docs.Patterns)
	docsRootSet := make(map[string]bool, len(docsRoots))
	for _, root := range docsRoots {
		docsRootSet[filepath.Clean(root)] = true
	}

	dirFiles := make(map[string]map[string]bool)
	e.walkDocFiles(func(path string, lines []string) {
		if e.isIgnoredDocPath(path) {
			return
		}
		dir := filepath.Clean(filepath.Dir(path))
		if docsRootSet[dir] {
			return // skip files sitting directly in the docs root
		}
		if dirFiles[dir] == nil {
			dirFiles[dir] = make(map[string]bool)
		}
		dirFiles[dir][filepath.Base(path)] = true
	})

	for _, analysis := range e.ensureSourceAnalyses() {
		docsDir, ok := docsDirectoryForSource(analysis.Path, docsRoots)
		if !ok || e.isIgnoredDocPath(docsDir) {
			continue
		}
		if dirFiles[docsDir] == nil {
			dirFiles[docsDir] = make(map[string]bool)
		}
	}

	directories := make([]string, 0, len(dirFiles))
	for dir := range dirFiles {
		directories = append(directories, dir)
	}
	sort.Strings(directories)
	for _, dir := range directories {
		files := dirFiles[dir]
		for _, req := range required {
			if !files[req] {
				e.result.AddError(
					"pkg-doc-files",
					fmt.Sprintf("pkg doc directory '%s' missing required file: %s", dir, req),
					filepath.Join(dir, req),
					"",
					req,
				)
			}
		}
	}
}

func configuredDocsRoots(patterns []string) []string {
	roots := make(map[string]bool)
	for _, pattern := range patterns {
		literal := pattern
		if wildcard := strings.IndexAny(literal, "*?["); wildcard >= 0 {
			literal = literal[:wildcard]
		}
		literal = strings.TrimRight(literal, `/\`)
		if literal == "" {
			continue
		}
		if filepath.Ext(literal) != "" {
			literal = filepath.Dir(literal)
		}
		for current := filepath.Clean(literal); ; current = filepath.Dir(current) {
			if filepath.Base(current) == "docs" {
				roots[current] = true
				break
			}
			parent := filepath.Dir(current)
			if parent == current || current == "." {
				break
			}
		}
	}
	result := make([]string, 0, len(roots))
	for root := range roots {
		result = append(result, root)
	}
	sort.Strings(result)
	return result
}

func docsDirectoryForSource(sourcePath string, docsRoots []string) (string, bool) {
	sourceAbsolute, err := filepath.Abs(sourcePath)
	if err != nil {
		return "", false
	}

	var selected string
	selectedRootLength := -1
	for _, docsRoot := range docsRoots {
		docsAbsolute, absErr := filepath.Abs(docsRoot)
		if absErr != nil {
			continue
		}
		projectRoot := filepath.Dir(docsAbsolute)
		relativeSource, relErr := filepath.Rel(projectRoot, sourceAbsolute)
		if relErr != nil ||
			relativeSource == ".." ||
			strings.HasPrefix(relativeSource, ".."+string(filepath.Separator)) {
			continue
		}
		packagePath := filepath.Dir(relativeSource)
		if packagePath == "." || packagePath == "" {
			continue
		}
		if len(projectRoot) <= selectedRootLength {
			continue
		}
		selected = filepath.Clean(filepath.Join(docsRoot, packagePath))
		selectedRootLength = len(projectRoot)
	}
	return selected, selected != ""
}
