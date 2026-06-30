// Package engine provides the core validation engine.

// Spec: docs/internal/engine/spec.md
// Contract: docs/internal/engine/contract.md

package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

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
	cfg    *config.Config
	graph  *graph.LinkageGraph
	result *model.ValidationResult
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
func (e *Engine) AddStructuralErrors(errors []*model.ValidationError) {
	for _, err := range errors {
		if err == nil {
			continue
		}
		e.result.AddError(err.Rule, err.Message, err.Source, err.Link, err.Code)
	}
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
	}

	e.graph.VerifyBidirectionalLinks()
}

// inferLinkType determines the appropriate link type based on the source
// identifier type and target reference. It maps SPEC->TEST as LinkTests,
// SPEC->CONTRACT as LinkContract, TEST->SPEC as LinkImplements, etc.
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
		if toType == model.TypeContract {
			return model.LinkContract
		}
		return model.LinkReferences
	case model.TypeTest:
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
		e.validatePublicFuncAnnotations()
	}

	if e.cfg.Validation.RequirePackageDocComment {
		e.validatePackageDocComment()
		e.validateDocPathMatchesPackagePath()
	}

	if e.cfg.Validation.RequireRelatedFiles {
		e.validateRelatedFiles()
	}

	if e.cfg.Validation.RequirePkgDocFiles {
		e.validatePkgDocFiles()
	}

	if e.cfg.Validation.RequireTestAnnotation {
		e.validateTestAnnotations()
	}

	if e.cfg.Validation.RequireAnnotationIdentifier {
		e.validateAnnotationIdentifiers()
	}

	if e.cfg.Validation.RequireAnnotationOnSameLine {
		e.validateConsecutiveAnnotations()
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
			idType := pattern.GetIdentifierType(node.ID)
			if idType == "CONTRACT" || idType == "DESIGN" {
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
	contractTestMap := make(map[string]bool)

	testContractRegex := regexp.MustCompile(`@test-contract\s+(TEST-[A-Z0-9_]+-[0-9]+)`)

	e.walkCodeFiles(func(path string, lines []string) {
		for _, line := range lines {
			if matches := testContractRegex.FindStringSubmatch(line); len(matches) > 1 {
				contractID := matches[1]
				contractTestMap[contractID] = true
			}
		}
	})

	for _, node := range e.graph.Nodes() {
		idType := pattern.GetIdentifierType(node.ID)
		if idType == "CONTRACT" {
			if !contractTestMap[node.ID] {
				e.result.AddError(
					"contract-test-coverage",
					fmt.Sprintf("%s has no @test-contract TEST annotation", node.ID),
					nodeSource(node),
					"",
					"",
				)
			}
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

		foundSections := make(map[string]bool)
		for _, line := range lines {
			if matches := designRegex.FindStringSubmatch(strings.TrimSpace(line)); len(matches) > 1 {
				foundSections[strings.ToLower(matches[1])] = true
			}
		}

		for _, section := range requiredSections {
			if !foundSections[section] {
				e.result.AddError(
					"design-sections",
					fmt.Sprintf("design.md missing required section: %s", section),
					path,
					"",
					"",
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
		hasDoc, hasCode := false, false
		if node.Metadata != nil {
			hasDoc, _ = node.Metadata[string(model.OriginDoc)].(bool)
			hasCode, _ = node.Metadata[string(model.OriginCode)].(bool)
		}

		if hasDoc && !hasCode {
			e.result.AddError(
				"doc-code-correspondence",
				fmt.Sprintf("%s is documented but missing @implement annotation in code", node.ID),
				nodeSourceByOrigin(node, model.OriginDoc),
				"",
				"",
			)
		}

		if hasCode && !hasDoc {
			e.result.AddError(
				"doc-code-correspondence",
				fmt.Sprintf("%s has @implement annotation in code but no documentation", node.ID),
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
		if docDescribe == "" || codeDescribe == "" {
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

// validatePublicFuncAnnotations enforces the relaxed @implement placement rules:
//  1. Public functions/types MUST carry an @implement annotation.
//  2. Private functions/types MAY carry an @implement annotation; when they do,
//     it must still be placed directly above the declaration.
//
// Doc/code correspondence for any annotated identifier is enforced separately
// in validateDocCodeCorrespondence — every @implement (public or private) must
// have a matching doc entry.
func (e *Engine) validatePublicFuncAnnotations() {
	publicFuncRegex := regexp.MustCompile(`^func\s+([A-Z][a-zA-Z0-9]*)\s*\(`)
	publicTypeRegex := regexp.MustCompile(`^type\s+([A-Z][a-zA-Z0-9]*)\s*`)
	// "any" regexes match both exported and unexported declarations. Used for
	// placement validation, which accepts @implement above private targets too.
	anyFuncRegex := regexp.MustCompile(`^func\s+([a-zA-Z][a-zA-Z0-9]*)\s*\(`)
	anyMethodRegex := regexp.MustCompile(`^func\s+\([^)]+\)\s*([a-zA-Z][a-zA-Z0-9]*)\s*\(`)
	anyTypeRegex := regexp.MustCompile(`^type\s+([a-zA-Z][a-zA-Z0-9]*)\s*`)
	// Only match lines where @implement appears at start of comment (after // and optional space)
	implementAtStartRegex := regexp.MustCompile(`^\s*//\s*@implement\b`)

	e.walkCodeFiles(func(path string, lines []string) {
		annotated := make(map[string]bool)
		annotatedLine := make(map[string]int)

		// First pass: find all @implement annotations and their target
		// functions/types — both public and private. Private targets won't be
		// required to carry @implement, but recognising them here keeps the
		// placement check from flagging valid private placements as errors.
		type anno struct {
			line     int
			funcName string
			isMethod bool
		}
		var annotations []anno

		ignoreScope := false
		for i, line := range lines {
			// Check for idd:ignore scope markers
			if strings.Contains(line, "// idd:ignore start") || strings.Contains(line, "//idd:ignore-start") {
				ignoreScope = true
				continue
			}
			if strings.Contains(line, "// idd:ignore end") || strings.Contains(line, "//idd:ignore-end") {
				ignoreScope = false
				continue
			}
			if ignoreScope {
				continue
			}

			if !implementAtStartRegex.MatchString(line) {
				continue
			}
			// Find the next function/type after this @implement line (public or private).
			for j := i + 1; j < len(lines); j++ {
				nextLine := strings.TrimSpace(lines[j])
				if strings.HasPrefix(nextLine, "//") || nextLine == "" {
					continue // skip comment lines and blank lines
				}
				if match := anyFuncRegex.FindStringSubmatch(nextLine); len(match) > 1 {
					annotations = append(annotations, anno{i, match[1], false})
					break
				}
				if match := anyMethodRegex.FindStringSubmatch(nextLine); len(match) > 1 {
					annotations = append(annotations, anno{i, match[1], true})
					break
				}
				if match := anyTypeRegex.FindStringSubmatch(nextLine); len(match) > 1 {
					annotations = append(annotations, anno{i, match[1], false})
					break
				}
				// Not a function/type line — stop searching
				break
			}
		}

		// Check for duplicate annotations on the same function/type (public or
		// private). Use "method:" prefix for methods to distinguish from
		// standalone functions with the same name.
		for _, a := range annotations {
			key := a.funcName
			if a.isMethod {
				key = "method:" + key
			}
			if annotated[key] {
				e.result.AddError(
					"duplicate-annotation",
					fmt.Sprintf("function/type %s has multiple @implement annotations; use comma-separated identifiers on a single line", a.funcName),
					fmt.Sprintf("%s:%d", path, a.line+1),
					"",
					"",
				)
			}
			annotated[key] = true
			annotatedLine[a.funcName] = a.line + 1
		}

		// Second pass: any @implement that isn't placed directly above a
		// function/type (public or private) is a placement error.
		ignoreScope = false
		for i, line := range lines {
			// Honor // idd:ignore start/end scope markers here too — otherwise
			// fixtures in test source files would still trip placement checks
			// even when the annotation lookup is suppressed above.
			if strings.Contains(line, "// idd:ignore start") || strings.Contains(line, "//idd:ignore-start") {
				ignoreScope = true
				continue
			}
			if strings.Contains(line, "// idd:ignore end") || strings.Contains(line, "//idd:ignore-end") {
				ignoreScope = false
				continue
			}
			if ignoreScope {
				continue
			}
			if !implementAtStartRegex.MatchString(line) {
				continue
			}
			if i+1 < len(lines) {
				nextLine := strings.TrimSpace(lines[i+1])
				if anyFuncRegex.MatchString(nextLine) || anyMethodRegex.MatchString(nextLine) || anyTypeRegex.MatchString(nextLine) {
					continue // valid placement, handled in first pass
				}
			}
			// Only report placement error if this @implement didn't find a valid target
			found := false
			for _, a := range annotations {
				if a.line == i {
					found = true
					break
				}
			}
			if !found {
				e.result.AddError(
					"annotation-placement",
					"@implement not directly above function/type",
					fmt.Sprintf("%s:%d", path, i+1),
					"",
					"",
				)
			}
		}

		ignoreScope = false
		for i, line := range lines {
			// Check for idd:ignore scope markers
			if strings.Contains(line, "// idd:ignore start") || strings.Contains(line, "//idd:ignore-start") {
				ignoreScope = true
				continue
			}
			if strings.Contains(line, "// idd:ignore end") || strings.Contains(line, "//idd:ignore-end") {
				ignoreScope = false
				continue
			}
			if ignoreScope {
				continue
			}

			funcName := ""
			if match := publicFuncRegex.FindStringSubmatch(line); len(match) > 1 {
				funcName = match[1]
			} else if match := publicTypeRegex.FindStringSubmatch(line); len(match) > 1 {
				funcName = match[1]
			}

			if funcName != "" && funcName != "main" && funcName != "init" && !strings.HasPrefix(funcName, "Test") && !strings.HasSuffix(funcName, "Test") {
				// Check both funcName and "method:"+funcName since annotated map uses method prefix
				if !annotated[funcName] && !annotated["method:"+funcName] {
					e.result.AddError(
						"public-func-annotation",
						fmt.Sprintf("public function/type %s missing @implement", funcName),
						fmt.Sprintf("%s:%d", path, i+1),
						"",
						"",
					)
				}
			}
		}
	})
}

// validatePackageDocComment checks that each Go package file has a proper
// package doc comment with Package description, Spec path, and Contract path.
// It also verifies there's a blank line between Package describe and Spec.
func (e *Engine) validatePackageDocComment() {
	specPathRegex := regexp.MustCompile(`(?i)^//\s*Spec:\s*docs/`)
	contractPathRegex := regexp.MustCompile(`(?i)^//\s*Contract:\s*docs/`)
	testPathRegex := regexp.MustCompile(`(?i)^//\s*Test:\s*docs/`)
	packageDocRegex := regexp.MustCompile(`^//\s*Package\s+\w+`)

	e.walkCodeFiles(func(path string, lines []string) {
		if filepath.Ext(path) != ".go" {
			return
		}

		dir := filepath.Dir(path)
		isTestFile := strings.HasSuffix(path, "_test.go")
		isContractTestFile := strings.HasSuffix(path, "_contract_test.go")

		var packageName string
		var packageLine int
		for i, line := range lines {
			if strings.HasPrefix(line, "package ") {
				packageName = strings.TrimPrefix(strings.TrimSpace(line), "package ")
				packageLine = i
				break
			}
		}

		if packageName == "" || packageName == "main" {
			return
		}

		commentStart := packageLine - 10
		if commentStart < 0 {
			commentStart = 0
		}

		var packageDocLine, specLine, testOrContractLine, contractLine = -1, -1, -1, -1
		for i := commentStart; i < packageLine; i++ {
			line := strings.TrimSpace(lines[i])
			if line == "" {
				continue
			}
			if packageDocRegex.MatchString(line) && packageDocLine == -1 {
				packageDocLine = i
			}
			if specPathRegex.MatchString(line) && specLine == -1 {
				specLine = i
			}
			if !isTestFile || isContractTestFile {
				if contractPathRegex.MatchString(line) && contractLine == -1 {
					contractLine = i
				}
			}
			if testPathRegex.MatchString(line) && testOrContractLine == -1 {
				testOrContractLine = i
			}
		}

		if packageDocLine == -1 {
			e.result.AddError(
				"package-doc-comment",
				fmt.Sprintf(`package %s missing package doc comment; add "// Package %s ..." before package declaration`, packageName, packageName),
				path, "", "",
			)
			return
		}
		if specLine == -1 {
			e.result.AddError(
				"package-doc-comment",
				fmt.Sprintf(`package %s missing Spec path; add "// Spec: docs/%s/spec.md" after package doc comment`, packageName, dir),
				path, "", "",
			)
			return
		}

		if isTestFile {
			if testOrContractLine == -1 {
				e.result.AddError(
					"package-doc-comment",
					fmt.Sprintf(`package %s missing Test path; add "// Test: docs/%s/testing.md" after package doc comment`, packageName, dir),
					path, "", "",
				)
				return
			}
			if isContractTestFile && contractLine == -1 {
				e.result.AddError(
					"package-doc-comment",
					fmt.Sprintf(`package %s missing Contract path; add "// Contract: docs/%s/contract.md" after package doc comment`, packageName, dir),
					path, "", "",
				)
				return
			}
		} else {
			if contractLine == -1 {
				e.result.AddError(
					"package-doc-comment",
					fmt.Sprintf(`package %s missing Contract path; add "// Contract: docs/%s/contract.md" after package doc comment`, packageName, dir),
					path, "", "",
				)
				return
			}
		}

		blankLineBetweenPackageAndSpec := false
		for i := packageDocLine + 1; i < specLine; i++ {
			if strings.TrimSpace(lines[i]) == "" {
				blankLineBetweenPackageAndSpec = true
				break
			}
		}
		if !blankLineBetweenPackageAndSpec {
			e.result.AddError(
				"package-doc-comment",
				fmt.Sprintf("package %s missing blank line between package doc comment and Spec path", packageName),
				path, "", "",
			)
		}
	})
}

// validateDocPathMatchesPackagePath checks that docs paths in package comments
// follow the format docs/<package_path>/xxx.md where <package_path> matches
// the actual package directory structure (e.g., docs/internal/auth/spec.md
// for package at internal/auth).
func (e *Engine) validateDocPathMatchesPackagePath() {
	// Match docs paths like "docs/auth/spec.md" or "docs/internal/engine/spec.md"
	docPathRegex := regexp.MustCompile(`(?i)^//\s*(Spec|Contract|Test):\s*(docs/.*\.md)`)

	e.walkCodeFiles(func(path string, lines []string) {
		if filepath.Ext(path) != ".go" {
			return
		}

		// Skip test files for this validation - they use Test: not path matching
		if strings.HasSuffix(path, "_test.go") {
			return
		}

		// Get the package directory from the file path
		// e.g., "internal/auth/auth.go" -> "internal/auth"
		dir := filepath.Dir(path)

		var packageLine int
		for i, line := range lines {
			if strings.HasPrefix(line, "package ") {
				packageLine = i
				break
			}
		}

		if packageLine == 0 {
			return
		}

		commentStart := packageLine - 10
		if commentStart < 0 {
			commentStart = 0
		}

		// Check each docs path reference
		for i := commentStart; i < packageLine; i++ {
			line := strings.TrimSpace(lines[i])
			if m := docPathRegex.FindStringSubmatch(line); len(m) > 2 {
				docsPath := m[2] // e.g., "docs/auth/spec.md"

				// Extract the path after "docs/"
				// e.g., "docs/auth/spec.md" -> "auth"
				// e.g., "docs/internal/auth/spec.md" -> "internal/auth"
				rest := strings.TrimPrefix(docsPath, "docs/")
				rest = strings.TrimSuffix(rest, filepath.Base(docsPath))
				rest = strings.TrimSuffix(rest, "/")

				// Check if the docs path matches the package directory
				if rest != dir {
					e.result.AddError(
						"package-doc-path",
						fmt.Sprintf("docs path %s does not match package path %s (file: %s, line: %d)", docsPath, dir, path, i+1),
						fmt.Sprintf("%s:%d", path, i+1),
						"",
						"",
					)
				}
			}
		}
	})
}

// validateTestAnnotations checks that test functions have proper annotations
// (like @test or @test-contract), and that referenced identifiers exist in the graph.
func (e *Engine) validateTestAnnotations() {
	testFuncRegex := regexp.MustCompile(`^func\s+(Test[A-Z][a-zA-Z0-9]*|[A-Z][a-zA-Z0-9]*Test[A-Z][a-zA-Z0-9]*)\s*\(`)
	testAnnotationRegex := regexp.MustCompile(`@test\s+(TEST-[A-Z0-9_]+-[0-9]+)`)
	testContractAnnotationRegex := regexp.MustCompile(`@test-contract\s+(TEST-[A-Z0-9_]+-[0-9]+)`)

	e.walkCodeFiles(func(path string, lines []string) {
		if !strings.HasSuffix(path, "_test.go") {
			return
		}

		isContractTestFile := strings.HasSuffix(path, "_contract_test.go")

		ignoreScope := false
		for i, line := range lines {
			// Check for idd:ignore scope markers
			if strings.Contains(line, "// idd:ignore start") || strings.Contains(line, "//idd:ignore-start") {
				ignoreScope = true
				continue
			}
			if strings.Contains(line, "// idd:ignore end") || strings.Contains(line, "//idd:ignore-end") {
				ignoreScope = false
				continue
			}

			// Skip if inside ignore scope
			if ignoreScope {
				continue
			}

			if match := testFuncRegex.FindStringSubmatch(line); len(match) > 1 {
				funcName := match[1]
				found := false
				var testID string
				for j := i - 1; j >= 0 && j >= i-5; j-- {
					if strings.TrimSpace(lines[j]) == "" {
						continue
					}
					if strings.HasPrefix(lines[j], "func ") {
						break
					}
					if strings.HasPrefix(lines[j], "//") {
						if isContractTestFile {
							if m := testContractAnnotationRegex.FindStringSubmatch(lines[j]); len(m) > 1 {
								found = true
								testID = m[1]
								break
							}
						} else {
							if m := testAnnotationRegex.FindStringSubmatch(lines[j]); len(m) > 1 {
								found = true
								testID = m[1]
								break
							}
						}
					}
					if strings.HasPrefix(lines[j], "/*") {
						break
					}
				}
				if !found {
					if isContractTestFile {
						e.result.AddError(
							"test-annotation",
							fmt.Sprintf("test function %s missing @test-contract annotation", funcName),
							fmt.Sprintf("%s:%d", path, i+1),
							"",
							"",
						)
					} else {
						e.result.AddError(
							"test-annotation",
							fmt.Sprintf("test function %s missing @test annotation", funcName),
							fmt.Sprintf("%s:%d", path, i+1),
							"",
							"",
						)
					}
				} else if testID != "" {
					if node, ok := e.graph.GetNode(testID); !ok || node == nil {
						e.result.AddError(
							"test-annotation",
							fmt.Sprintf("test function %s references undefined %s", funcName, testID),
							fmt.Sprintf("%s:%d", path, i+1),
							"",
							"",
						)
					}
				}
			}
		}
	})
}

// validateRelatedFiles checks that documentation files have the required
// related_files field in their frontmatter.
func (e *Engine) validateRelatedFiles() {
	frontmatterRegex := regexp.MustCompile(`^related_files:`)

	e.walkDocFiles(func(path string, lines []string) {
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

// validatePkgDocFiles checks that every pkg docs directory contains all four
// required documentation files: spec.md, contract.md, testing.md, design.md.
func (e *Engine) validatePkgDocFiles() {
	required := []string{"spec.md", "contract.md", "testing.md", "design.md"}

	// Derive the docs root directories from configured patterns so the check
	// works with both relative ("docs/**/*.md") and absolute paths (tests).
	docsRoots := make(map[string]bool)
	for _, pat := range e.cfg.Docs.Patterns {
		slash := filepath.ToSlash(pat)
		base := strings.SplitN(slash, "**", 2)[0]
		docsRoots[filepath.Clean(base)] = true
	}

	dirFiles := make(map[string]map[string]bool)
	e.walkDocFiles(func(path string, lines []string) {
		if e.isIgnoredDocPath(path) {
			return
		}
		dir := filepath.Clean(filepath.Dir(path))
		if docsRoots[dir] {
			return // skip files sitting directly in the docs root
		}
		if dirFiles[dir] == nil {
			dirFiles[dir] = make(map[string]bool)
		}
		dirFiles[dir][filepath.Base(path)] = true
	})

	for dir, files := range dirFiles {
		for _, req := range required {
			if !files[req] {
				e.result.AddError(
					"pkg-doc-files",
					fmt.Sprintf("pkg doc directory '%s' missing required file: %s", dir, req),
					dir, "", "",
				)
			}
		}
	}
}

// walkCodeFiles iterates over all code files matching the configured patterns
// and invokes the visitor function for each file with its path and line content.
func (e *Engine) walkCodeFiles(visitor func(path string, lines []string)) {
	for _, globPattern := range e.cfg.Code.Patterns {
		e.walkGlob(globPattern, visitor)
	}
}

// walkDocFiles iterates over all documentation files matching the configured
// patterns and invokes the visitor function for each file with its path and lines.
func (e *Engine) walkDocFiles(visitor func(path string, lines []string)) {
	for _, globPattern := range e.cfg.Docs.Patterns {
		e.walkGlob(globPattern, visitor)
	}
}

// walkGlob walks the directory tree matching the given pattern and invokes
// the visitor function for each matching file. Handles non-recursive patterns.
func (e *Engine) walkGlob(pattern string, visitor func(path string, lines []string)) {
	hasRecursive := strings.Contains(pattern, "**")
	if hasRecursive {
		e.walkGlobRecursive(pattern, visitor)
		return
	}
	dir, file := filepath.Split(pattern)
	if dir == "" {
		dir = "."
	}
	if file == "" {
		return
	}
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		matched, err := filepath.Match(file, d.Name())
		if err != nil || !matched {
			return nil
		}
		fullPath := filepath.Join(dir, d.Name())
		if e.shouldIgnorePath(fullPath) {
			return nil
		}
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return nil
		}
		lines := strings.Split(string(content), "\n")
		visitor(fullPath, lines)
		return nil
	})
}

// walkGlobRecursive walks the directory tree matching the given recursive
// pattern (containing **/) and invokes the visitor function for each matching file.
func (e *Engine) walkGlobRecursive(pattern string, visitor func(path string, lines []string)) {
	dir, file := filepath.Split(pattern)
	if dir == "" {
		dir = "."
	}
	if file == "" {
		return
	}
	dir = strings.TrimSuffix(dir, "**/")
	if dir == "" {
		dir = "."
	}
	file = strings.TrimPrefix(file, "**/")

	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		matched, err := filepath.Match(file, d.Name())
		if err != nil || !matched {
			return nil
		}
		fullPath := path
		if e.shouldIgnorePath(fullPath) {
			return nil
		}
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return nil
		}
		lines := strings.Split(string(content), "\n")
		visitor(fullPath, lines)
		return nil
	})
}

// shouldIgnorePath checks whether a file path matches any of the configured
// ignore patterns for code files.
func (e *Engine) shouldIgnorePath(path string) bool {
	for _, ignore := range e.cfg.Code.IgnorePaths {
		if matched, _ := filepath.Match(ignore, filepath.Base(path)); matched {
			return true
		}
		if strings.Contains(path, ignore) {
			return true
		}
	}
	return false
}

// BuildReport generates a complete validation report with tool information,
// configuration summary, and the accumulated validation result.
// @implement SPEC-CMD_IDD_CLI-004
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
