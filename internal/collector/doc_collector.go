// Package collector provides documentation identifier collection functionality.

// Spec: docs/internal/collector/spec.md
// Contract: docs/internal/collector/contract.md
package collector

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jingxu9x/idd-cli/internal/config"
	"github.com/jingxu9x/idd-cli/internal/model"
	"github.com/jingxu9x/idd-cli/pkg/pattern"
)

// DocCollector collects IDD identifiers from markdown documentation files,
// parsing frontmatter markers and extracting identifier references from content.
//
// @implement SPEC-INT_COL-001
type DocCollector struct {
	cfg *config.Config
}

// NewDocCollector creates a new DocCollector with the given configuration.
//
// @implement SPEC-INT_COL-002
func NewDocCollector(cfg *config.Config) *DocCollector {
	return &DocCollector{cfg: cfg}
}

// Collect collects IDD identifiers from markdown files at the target path.
// If targetPath is a directory, recursively walks to find all .md files.
//
// @implement SPEC-INT_COL-017
func (c *DocCollector) Collect(ctx context.Context, targetPath string) (*model.IdentifierSet, []*model.ValidationError, error) {
	set := model.NewIdentifierSet()
	var errors []*model.ValidationError

	info, err := os.Stat(targetPath)
	if err != nil {
		return set, errors, nil
	}

	if info.IsDir() {
		err := filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			if c.shouldIgnore(path) {
				return nil
			}
			if filepath.Ext(path) == ".md" {
				fileErrors := c.collectFile(path, set)
				errors = append(errors, fileErrors...)
			}
			return nil
		})
		if err != nil {
			return set, errors, err
		}
	} else {
		fileErrors := c.collectFile(targetPath, set)
		errors = append(errors, fileErrors...)
	}

	return set, errors, nil
}

// collectFile parses a markdown file and extracts IDD identifiers from it.
// It handles frontmatter parsing, validation, and reference extraction.
func (c *DocCollector) collectFile(path string, set *model.IdentifierSet) []*model.ValidationError {
	var errors []*model.ValidationError

	content, err := os.ReadFile(path)
	if err != nil {
		return errors
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

			if err := ValidateModulePrefix(ref, path); err != nil {
				errors = append(errors, &model.ValidationError{
					Rule:    "module-prefix-mismatch",
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
			if err := ValidateModulePrefix(marker.ID, path); err != nil {
				errors = append(errors, &model.ValidationError{
					Rule:    "module-prefix-mismatch",
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
