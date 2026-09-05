// Package collector provides documentation identifier collection functionality.
package collector

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/jingxu9x/idd-cli/internal/config"
	"github.com/jingxu9x/idd-cli/internal/graph"
	"github.com/jingxu9x/idd-cli/internal/model"
	"github.com/jingxu9x/idd-cli/pkg/pattern"
)

// DocCollector collects IDD identifiers from markdown documentation files,
// parsing frontmatter markers and extracting identifier references from content.
//
// @implement SPEC-INTERNAL_COLLECTOR-001
type DocCollector struct {
	cfg      *config.Config
	files    map[string][]byte
	mentions []graph.Occurrence
}

// NewDocCollector creates a new DocCollector with the given configuration.
//
// @implement SPEC-INTERNAL_COLLECTOR-002
func NewDocCollector(cfg *config.Config) *DocCollector {
	return &DocCollector{cfg: cfg, files: make(map[string][]byte)}
}

// Collect collects IDD identifiers from Markdown files at the target path.
// Self-describing document sets are discovered before legacy Markdown so one
// package always uses one deterministic parsing mode.
//
// @implement SPEC-INTERNAL_COLLECTOR-017
func (c *DocCollector) Collect(ctx context.Context, targetPath string) (*model.IdentifierSet, []*model.ValidationError, error) {
	set := model.NewIdentifierSet()
	c.files = make(map[string][]byte)
	c.mentions = nil
	var errors []*model.ValidationError
	if err := c.cfg.Validate(); err != nil {
		errors = append(errors, &model.ValidationError{Rule: "filesystem-scan", Message: err.Error(), Source: "configuration", Code: "invalid-pattern"})
		return set, errors, nil
	}

	resolvedTarget := c.cfg.ResolvePath(targetPath)
	info, err := os.Stat(resolvedTarget)
	if err != nil {
		errors = append(errors, filesystemScanFinding(c.cfg.DisplayPath(resolvedTarget), "stat", err))
		return set, errors, nil
	}

	var markdownPaths []string
	var centralCatalogPaths []string
	markdownSeen := make(map[string]bool)
	centralCatalogSeen := make(map[string]bool)
	addMarkdown := func(path string) {
		path = filepath.Clean(path)
		if markdownSeen[path] || c.shouldIgnore(path) {
			return
		}
		markdownSeen[path] = true
		markdownPaths = append(markdownPaths, path)
	}
	addCentralCatalog := func(path string) {
		path = filepath.Clean(path)
		if centralCatalogSeen[path] || c.shouldIgnore(path) {
			return
		}
		if packageFromDocumentPath(filepath.Join(filepath.Dir(path), "spec.md")) == "" {
			return
		}
		centralCatalogSeen[path] = true
		centralCatalogPaths = append(centralCatalogPaths, path)
	}

	if info.IsDir() {
		err := filepath.Walk(resolvedTarget, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				errors = append(errors, filesystemScanFinding(c.cfg.DisplayPath(path), "walk", err))
				return nil
			}
			if info.IsDir() {
				return nil
			}
			displayPath := c.cfg.DisplayPath(path)
			if c.shouldIgnore(displayPath) {
				return nil
			}
			if filepath.Base(displayPath) == "idd.yaml" {
				addCentralCatalog(displayPath)
			} else if filepath.Ext(displayPath) == ".md" {
				addMarkdown(displayPath)
			}
			return nil
		})
		if err != nil {
			return set, errors, err
		}
	} else {
		targetPath = c.cfg.DisplayPath(resolvedTarget)
		if filepath.Base(targetPath) == "idd.yaml" {
			addCentralCatalog(targetPath)
		} else if filepath.Ext(targetPath) == ".md" {
			addMarkdown(targetPath)
			if _, fixedDocument := iddDocumentRoles[filepath.Base(targetPath)]; fixedDocument {
				centralCatalogPath := filepath.Join(filepath.Dir(targetPath), "idd.yaml")
				if centralInfo, statErr := c.stat(centralCatalogPath); statErr == nil && !centralInfo.IsDir() {
					addCentralCatalog(centralCatalogPath)
				}
				for _, filename := range iddDocumentOrder {
					siblingPath := filepath.Join(filepath.Dir(targetPath), filename)
					if siblingInfo, statErr := c.stat(siblingPath); statErr == nil && !siblingInfo.IsDir() {
						addMarkdown(siblingPath)
					}
				}
			}
		}
	}

	sort.Strings(centralCatalogPaths)
	sort.Strings(markdownPaths)
	for _, path := range centralCatalogPaths {
		errors = append(errors, iddDocumentValidationError(
			"idd-document-migration",
			"package-local idd.yaml is not a marker source; move each declaration into its owning Markdown document",
			path,
			1,
			"",
			"central-catalog",
		))
	}

	invalidIDDPaths := make(map[string]bool)
	for _, path := range markdownPaths {
		filename := filepath.Base(path)
		if _, fixedDocument := iddDocumentRoles[filename]; fixedDocument {
			continue
		}
		if canonical, split := splitIDDDocumentFilename(filename); split {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-filename",
				fmt.Sprintf(
					"IDD role documents are not split; merge %s into %s",
					filename,
					canonical,
				),
				path,
				1,
				canonical,
				"split-role",
			))
			invalidIDDPaths[path] = true
			continue
		}
		data, readErr := c.readFile(path)
		if readErr != nil {
			errors = append(errors, filesystemScanFinding(path, "read", readErr))
			invalidIDDPaths[path] = true
		} else if hasIDDDocumentFrontmatter(data) {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-filename",
				"IDD frontmatter is allowed only in design.md, contract.md, spec.md, or testing.md",
				path,
				1,
				"",
				"noncanonical-role",
			))
			invalidIDDPaths[path] = true
		}
	}

	iddDirectories := make(map[string]bool)
	for _, path := range markdownPaths {
		if _, fixedDocument := iddDocumentRoles[filepath.Base(path)]; !fixedDocument {
			continue
		}
		data, readErr := c.readFile(path)
		if readErr != nil {
			if !invalidIDDPaths[path] {
				errors = append(errors, filesystemScanFinding(path, "read", readErr))
			}
			invalidIDDPaths[path] = true
		} else if hasIDDDocumentFrontmatter(data) {
			iddDirectories[filepath.Clean(filepath.Dir(path))] = true
		}
	}

	for _, directory := range sortedIDDDocumentPaths(iddDirectories) {
		errors = append(errors, c.collectIDDDocumentSet(directory, set)...)
	}
	for _, path := range markdownPaths {
		if invalidIDDPaths[path] {
			continue
		}
		if iddDirectories[filepath.Clean(filepath.Dir(path))] {
			if _, fixedDocument := iddDocumentRoles[filepath.Base(path)]; !fixedDocument {
				errors = append(errors, c.collectIDDPackageNarrative(path, set)...)
			}
			continue
		}
		errors = append(errors, c.collectFile(path, set)...)
	}
	c.collectMentions(markdownPaths, set)

	return set, errors, nil
}

// Mentions returns non-authoritative Markdown references found during the most
// recent collection. The content is read through the collector's scan cache,
// so trace does not perform a second filesystem scan.
func (c *DocCollector) Mentions() []graph.Occurrence {
	return append([]graph.Occurrence(nil), c.mentions...)
}

func (c *DocCollector) collectMentions(paths []string, set *model.IdentifierSet) {
	structured := graph.BuildTraceIndex(set)
	seen := make(map[string]bool)
	for _, path := range paths {
		content, ok := c.files[filepath.Clean(path)]
		if !ok {
			continue
		}
		role := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		for lineIndex, line := range strings.Split(string(content), "\n") {
			lineNumber := lineIndex + 1
			for _, ref := range pattern.ExtractIDDReferences(line) {
				if traceHasOccurrenceAt(structured.Occurrences(ref), path, lineNumber) {
					continue
				}
				key := strings.Join([]string{ref, filepath.Clean(path), fmt.Sprint(lineNumber)}, "\x00")
				if seen[key] {
					continue
				}
				seen[key] = true
				c.mentions = append(c.mentions, graph.Occurrence{
					EntityID: ref, Origin: "doc", DocumentRole: role,
					Kind: graph.OccurrenceMention, Path: path, Line: lineNumber,
				})
			}
		}
	}
	sort.Slice(c.mentions, func(i, j int) bool {
		if c.mentions[i].EntityID != c.mentions[j].EntityID {
			return c.mentions[i].EntityID < c.mentions[j].EntityID
		}
		if c.mentions[i].Path != c.mentions[j].Path {
			return c.mentions[i].Path < c.mentions[j].Path
		}
		return c.mentions[i].Line < c.mentions[j].Line
	})
}

func traceHasOccurrenceAt(occurrences []graph.Occurrence, path string, line int) bool {
	for _, occurrence := range occurrences {
		if filepath.Clean(occurrence.Path) == filepath.Clean(path) && occurrence.Line == line {
			return true
		}
	}
	return false
}

func hasLeadingFrontmatter(content string) bool {
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		return trimmed == "---"
	}
	return false
}

// collectFile parses a markdown file and extracts IDD identifiers from it.
// It handles frontmatter parsing, validation, and reference extraction.
func (c *DocCollector) collectFile(path string, set *model.IdentifierSet) []*model.ValidationError {
	var errors []*model.ValidationError

	content, err := c.readFile(path)
	if err != nil {
		return append(errors, filesystemScanFinding(path, "read", err))
	}

	fm, err := ParseFrontmatter(string(content))
	if err != nil {
		errors = append(errors, &model.ValidationError{
			Rule:    "frontmatter-parse",
			Message: err.Error(),
			Source:  path,
		})
	}

	if fm != nil {
		if fmErrors := ValidateFrontmatterMarkers(fm, string(content), path); len(fmErrors) > 0 {
			for _, e := range fmErrors {
				errors = append(errors, &model.ValidationError{
					Rule:    "frontmatter-mismatch",
					Message: e,
					Source:  path,
				})
			}
		}
	}

	lines := strings.Split(string(content), "\n")

	frontmatterIDs := make(map[string]bool)
	if fm != nil {
		for _, marker := range fm.Markers {
			frontmatterIDs[marker.ID] = true
		}
	}

	for i, line := range lines {
		refs := pattern.ExtractIDDReferences(line)
		for _, ref := range refs {
			idTypeStr := pattern.GetIdentifierType(ref)
			idType, _ := model.ParseIdentifierType(idTypeStr)
			if idType == "" {
				continue
			}

			if err := pattern.ValidateIdentifierFormat(ref); err != nil {
				errors = append(errors, &model.ValidationError{
					Rule:    "identifier-format",
					Message: err.Error(),
					Source:  path,
					Link:    ref,
				})
			}

			if frontmatterIDs[ref] {
				continue
			}
			if !set.Has(ref) {
				title := c.extractTitle(string(content), ref)
				id := model.NewIdentifier(ref, idType, title, path, i+1)
				id.RawRef = ref
				set.Add(id)
			}
		}
	}

	if fm != nil {
		for _, marker := range fm.Markers {
			idTypeStr := pattern.GetIdentifierType(marker.ID)
			idType, _ := model.ParseIdentifierType(idTypeStr)
			if idType == "" {
				continue
			}
			if err := ValidateDocumentStructure(path, string(idType)); err != nil {
				errors = append(errors, &model.ValidationError{
					Rule:    "document-structure",
					Message: err.Error(),
					Source:  path,
					Link:    marker.ID,
				})
			}
			existingID, exists := set.Get(marker.ID)
			if !exists {
				id := model.NewIdentifierWithDescribe(marker.ID, idType, marker.Name, marker.Describe, path, 0)
				id.RawRef = marker.ID
				set.Add(id)
			} else {
				existingID.Source = path
			}
		}
		// Add links based on explicit relationship fields
		// For TEST markers: parse **Spec Coverage:** field
		// For SPEC markers: parse **Tests:** field
		for _, marker := range fm.Markers {
			idTypeStr := pattern.GetIdentifierType(marker.ID)
			idType, _ := model.ParseIdentifierType(idTypeStr)
			if idType == "" {
				continue
			}

			id, ok := set.Get(marker.ID)
			if !ok {
				continue
			}

			sectionContent := extractSectionContent(string(content), marker.ID)
			if sectionContent == "" {
				continue
			}

			switch idType {
			case model.TypeTest:
				specs := extractSpecCoverage(sectionContent)
				id.Links = nil
				for _, spec := range specs {
					if spec != marker.ID {
						id.AddLink(spec)
					}
				}
			case model.TypeSpec:
				tests := extractTestsField(sectionContent)
				id.Links = nil
				for _, test := range tests {
					if test != marker.ID {
						id.AddLink(test)
					}
				}
				contracts := extractContractField(sectionContent)
				for _, contract := range contracts {
					if contract != marker.ID {
						id.AddLink(contract)
					}
				}
			}
		}
	}

	if strings.HasSuffix(path, "contract.md") && fm != nil {
		for _, marker := range fm.Markers {
			idTypeStr := pattern.GetIdentifierType(marker.ID)
			idType, _ := model.ParseIdentifierType(idTypeStr)
			if idType != model.TypeSpec {
				continue
			}
			id, ok := set.Get(marker.ID)
			if !ok {
				continue
			}
			sectionContent := extractSectionContent(string(content), marker.ID)
			if sectionContent == "" {
				continue
			}
			implements := extractImplementsField(sectionContent)
			for _, spec := range implements {
				if spec != marker.ID {
					id.AddLink(spec)
				}
			}
		}
	}

	return errors
}

// extractTitle extracts the title from a markdown heading that contains the given identifier.
func (c *DocCollector) extractTitle(content string, id string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, id) {
			if strings.Contains(line, "#") {
				return strings.TrimSpace(strings.SplitAfter(line, "#")[1])
			}
		}
	}
	return ""
}

// extractSectionContent extracts the markdown content under a section heading
// that starts with the given markerID. Handles H2 and H3 headings.
func extractSectionContent(content, markerID string) string {
	// Match ## or ### at start of line followed by markerID (allows H2 and H3 headings)
	pattern := regexp.MustCompile(`(?im)^#{2,3}\s+` + regexp.QuoteMeta(markerID) + `.*`)
	matches := pattern.FindStringSubmatchIndex(content)
	if len(matches) == 0 {
		return ""
	}
	result := content[matches[0]:matches[1]]
	remainingIdx := matches[1]
	for {
		nextMarkerIdx := -1
		nextPrefix := ""
		for _, prefix := range []string{"\n## ", "\n# ", "\n---"} {
			idx := strings.Index(content[remainingIdx:], prefix)
			if idx != -1 && (nextMarkerIdx == -1 || idx < nextMarkerIdx) {
				nextMarkerIdx = idx
				nextPrefix = prefix
			}
		}
		if nextMarkerIdx == -1 {
			result += content[remainingIdx:]
			break
		}
		// If we found ---, check what kind:
		// 1. Standalone --- line (section separator) - stop here
		// 2. --- inside a table row (| --- |) - skip and continue
		// 3. --- followed by ### (subsection separator within section) - skip and continue
		if nextPrefix == "\n---" {
			lineStart := remainingIdx + nextMarkerIdx + 1 // +1 to skip the newline
			lineEnd := lineStart
			for lineEnd < len(content) && content[lineEnd] != '\n' {
				lineEnd++
			}
			line := strings.TrimSpace(content[lineStart:lineEnd])
			// Check if this --- is part of a table row (contains |)
			if strings.Contains(line, "|") {
				// Table row separator, skip it
				remainingIdx = lineEnd
				continue
			}
			// Check if --- is followed by ### (subsection separator)
			afterDash := lineEnd + 1
			// Skip whitespace
			for afterDash < len(content) && (content[afterDash] == ' ' || content[afterDash] == '\t') {
				afterDash++
			}
			// Now find the start of the next line (skip blank lines)
			for afterDash < len(content) {
				// Check if we've reached a newline
				if content[afterDash] == '\n' {
					afterDash++
					continue
				}
				// Found non-newline content - check if it's a heading
				break
			}
			// Now check what the next non-blank line starts with
			// Only ## or # (possibly with leading whitespace before the #) at the start of content is a top-level heading
			// ### is a subsection, bold text (**), and other content should be skipped
			if afterDash < len(content) {
				// Look at what follows afterDash (after skipping any leading whitespace)
				skip := 0
				for afterDash+skip < len(content) && (content[afterDash+skip] == ' ' || content[afterDash+skip] == '\t') {
					skip++
				}
				next := content[afterDash+skip:]
				// If the next non-blank content starts with ## or # (H2/H1), it's a top-level heading (STOP)
				// If it starts with ###, it's a subsection heading - the --- is a subsection separator
				// within the same section, so DON'T STOP, continue extracting
				if strings.HasPrefix(next, "## ") || strings.HasPrefix(next, "##\n") ||
					strings.HasPrefix(next, "# ") || strings.HasPrefix(next, "#\n") {
					// Top-level heading - this is a section separator, stop
					result += content[remainingIdx : remainingIdx+nextMarkerIdx]
					break
				}
				// Not a top-level heading - continue past this --- to keep extracting
				// Append content up to and including the --- before continuing
				result += content[remainingIdx : remainingIdx+nextMarkerIdx+len(nextPrefix)]
				remainingIdx = afterDash
				continue
			}
			// Ran out of content - treat as section separator
			result += content[remainingIdx : remainingIdx+nextMarkerIdx]
			break
		}
		result += content[remainingIdx : remainingIdx+nextMarkerIdx]
		break
	}
	return result
}

// extractSpecCoverage extracts **Spec Coverage:** field value from section content.
func extractSpecCoverage(section string) []string {
	// Match **Spec Coverage:** followed by backtick-enclosed content
	re := regexp.MustCompile(`\*\*Spec Coverage:\*\*\s*(.+)`)
	match := re.FindStringSubmatch(section)
	if match == nil {
		return nil
	}
	// Split by comma and trim backticks and whitespace from each spec
	var specs []string
	for _, s := range strings.Split(match[1], ",") {
		s = strings.TrimSpace(s)
		s = strings.Trim(s, "`")
		if s != "" {
			specs = append(specs, s)
		}
	}
	return specs
}

// extractTestsField extracts **Tests:** field value from section content.
func extractTestsField(section string) []string {
	re := regexp.MustCompile(`\*\*Tests:\*\*\s*(.+)`)
	match := re.FindStringSubmatch(section)
	if match == nil {
		return nil
	}
	var tests []string
	for _, s := range strings.Split(match[1], ",") {
		s = strings.TrimSpace(s)
		s = strings.Trim(s, "`")
		if s != "" {
			tests = append(tests, s)
		}
	}
	return tests
}

// extractContractField extracts **Contract:** field value from section content.
func extractContractField(section string) []string {
	re := regexp.MustCompile(`\*\*Contract:\*\*\s*(.+)`)
	match := re.FindStringSubmatch(section)
	if match == nil {
		return nil
	}
	var contracts []string
	for _, s := range strings.Split(match[1], ",") {
		s = strings.TrimSpace(s)
		s = strings.Trim(s, "`")
		if s != "" {
			contracts = append(contracts, s)
		}
	}
	return contracts
}

// extractImplementsField extracts **Implements:** field value from section content.
func extractImplementsField(section string) []string {
	re := regexp.MustCompile(`\*\*Implements:\*\*\s*(.+)`)
	match := re.FindStringSubmatch(section)
	if match == nil {
		return nil
	}
	var specs []string
	for _, s := range strings.Split(match[1], ",") {
		s = strings.TrimSpace(s)
		s = strings.Trim(s, "`")
		if s != "" {
			specs = append(specs, s)
		}
	}
	return specs
}

func (c *DocCollector) readFile(path string) ([]byte, error) {
	path = filepath.Clean(path)
	if content, ok := c.files[path]; ok {
		return append([]byte(nil), content...), nil
	}
	content, err := os.ReadFile(c.cfg.ResolvePath(path))
	if err != nil {
		return nil, err
	}
	c.files[path] = append([]byte(nil), content...)
	return content, nil
}

func (c *DocCollector) stat(path string) (os.FileInfo, error) {
	return os.Stat(c.cfg.ResolvePath(path))
}

// shouldIgnore checks if a path should be ignored based on configured ignore patterns.
func (c *DocCollector) shouldIgnore(path string) bool {
	if len(c.cfg.Docs.IgnorePaths) == 0 {
		return false
	}
	for _, ignore := range c.cfg.Docs.IgnorePaths {
		if matchGlob(ignore, path) {
			return true
		}
	}
	return false
}
