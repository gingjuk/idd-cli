// Package collector provides frontmatter parsing and validation functionality.
package collector

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// @implement SPEC-INTERNAL_COLLECTOR-005
type Marker struct {
	ID       string `yaml:"id"`
	Name     string `yaml:"name"`
	Describe string `yaml:"describe"`
}

// @implement SPEC-INTERNAL_COLLECTOR-006
type RelatedFiles struct {
	Spec     string `yaml:"spec,omitempty"`
	Contract string `yaml:"contract,omitempty"`
	Design   string `yaml:"design,omitempty"`
	Testing  string `yaml:"testing,omitempty"`
}

// @implement SPEC-INTERNAL_COLLECTOR-007
type Frontmatter struct {
	Markers      []Marker      `yaml:"markers"`
	RelatedFiles *RelatedFiles `yaml:"related_files,omitempty"`
}

// @implement SPEC-INTERNAL_COLLECTOR-008
func ParseFrontmatter(content string) (*Frontmatter, error) {
	lines := strings.Split(content, "\n")
	startIdx, endIdx := -1, -1
	codeFence := ""

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		var transition bool
		codeFence, transition = markdownFenceTransition(codeFence, trimmed)
		if transition {
			continue
		}
		if codeFence != "" {
			continue
		}
		if trimmed == "---" {
			if startIdx == -1 {
				startIdx = i
			} else {
				endIdx = i
				break
			}
		}
	}

	if startIdx == -1 || endIdx == -1 || endIdx <= startIdx+1 {
		return nil, nil
	}

	yamlContent := strings.Join(lines[startIdx+1:endIdx], "\n")
	var fm Frontmatter
	if err := yaml.Unmarshal([]byte(yamlContent), &fm); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}
	return &fm, nil
}

// @implement SPEC-INTERNAL_COLLECTOR-009
func ValidateFrontmatterMarkers(fm *Frontmatter, content string, filePath string) []string {
	var errors []string
	if fm == nil {
		return errors
	}

	definedMarkers := extractDefinedMarkers(content)
	referencedMarkers := extractReferencedMarkers(content)
	headingLines := extractHeadingLines(content)

	frontmatterIDs := make(map[string]bool)
	for _, marker := range fm.Markers {
		frontmatterIDs[marker.ID] = true
	}

	for _, marker := range fm.Markers {
		if !definedMarkers[marker.ID] {
			if referencedMarkers[marker.ID] {
				errors = append(errors, fmt.Sprintf(
					"frontmatter marker '%s' is only referenced in content but not defined as heading (file: %s)",
					marker.ID, filePath,
				))
			} else {
				errors = append(errors, fmt.Sprintf(
					"frontmatter marker '%s' not found in content (file: %s)",
					marker.ID, filePath,
				))
			}
		}
	}

	for id, line := range headingLines {
		if err := ValidateHeadingFormat(id, line); err != nil {
			errors = append(errors, fmt.Sprintf("%s (file: %s)", err.Error(), filePath))
		}
	}

	if errs := ValidateMarkerFormatting(content, filePath); len(errs) > 0 {
		errors = append(errors, errs...)
	}

	return errors
}

// @implement SPEC-INTERNAL_COLLECTOR-010
func ValidateMarkerFormatting(content string, filePath string) []string {
	var errors []string
	lines := strings.Split(content, "\n")
	codeFence := ""
	inFrontmatter := false

	for lineNum, line := range lines {
		trimmed := strings.TrimSpace(line)
		var transition bool
		codeFence, transition = markdownFenceTransition(codeFence, trimmed)
		if transition {
			continue
		}
		if codeFence != "" {
			continue
		}
		if trimmed == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			} else {
				inFrontmatter = false
				continue
			}
		}
		if inFrontmatter {
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			continue
		}

		bareMarkers := extractBareMarkers(line)
		for _, marker := range bareMarkers {
			errors = append(errors, fmt.Sprintf(
				"doc marker '%s' should be wrapped in backticks (file: %s, line: %d)",
				marker, filePath, lineNum+1,
			))
		}
	}
	return errors
}

func extractBareMarkers(line string) []string {
	var bare []string
	// First, remove double-backtick examples (e.g., `` `SPEC-BE-001` ``)
	backtickPattern := regexp.MustCompile("``[^`]*`[^`]+`[^`]*``")
	lineWithoutExamples := backtickPattern.ReplaceAllString(line, "")

	// Remove single-backtick wrapped identifiers (these are properly formatted)
	// Replace with spaces to preserve word boundaries
	singleBacktickPattern := regexp.MustCompile("`[^`]+`")
	lineWithoutBackticks := singleBacktickPattern.ReplaceAllString(lineWithoutExamples, " ")

	// Now find any remaining bare identifiers
	matches := idPattern.FindAllStringSubmatch(lineWithoutBackticks, -1)
	for _, m := range matches {
		if len(m) > 1 {
			bare = append(bare, m[1])
		}
	}
	return bare
}

func extractDefinedMarkers(content string) map[string]bool {
	markers := make(map[string]bool)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			for _, ref := range extractIDDRefs(trimmed) {
				markers[ref] = true
			}
		}
	}
	return markers
}

func extractReferencedMarkers(content string) map[string]bool {
	markers := make(map[string]bool)
	for _, line := range strings.Split(content, "\n") {
		for _, ref := range extractIDDRefs(line) {
			markers[ref] = true
		}
	}
	return markers
}

func extractHeadingLines(content string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(content, "\n")
	codeFence := ""
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		var transition bool
		codeFence, transition = markdownFenceTransition(codeFence, trimmed)
		if transition {
			continue
		}
		if codeFence != "" {
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			for _, ref := range extractIDDRefs(trimmed) {
				result[ref] = trimmed
			}
		}
	}
	return result
}

func markdownFenceTransition(current, line string) (string, bool) {
	if current != "" {
		run := markdownFenceRun(line, current[0])
		if run >= len(current) && strings.TrimSpace(line[run:]) == "" {
			return "", true
		}
		return current, false
	}
	for _, marker := range []byte{'`', '~'} {
		if run := markdownFenceRun(line, marker); run >= 3 {
			return line[:run], true
		}
	}
	return "", false
}

func markdownFenceRun(line string, marker byte) int {
	run := 0
	for run < len(line) && line[run] == marker {
		run++
	}
	return run
}

// @implement SPEC-INTERNAL_COLLECTOR-011
func ValidateHeadingFormat(id, heading string) error {
	colonIdx := strings.Index(heading, ":")
	if colonIdx == -1 {
		return fmt.Errorf("heading for '%s' missing ':' separator, expected format: '%s: <description>'", id, id)
	}
	headingId := strings.TrimSpace(strings.TrimLeft(heading[:colonIdx], "# "))
	if headingId != id {
		return fmt.Errorf("heading identifier '%s' does not match marker '%s'", headingId, id)
	}
	desc := strings.TrimSpace(heading[colonIdx+1:])
	if desc == "" {
		return fmt.Errorf("heading for '%s' missing description after ':'", id)
	}

	parts := strings.Split(id, "-")

	if len(parts) == 2 {
		if err := validateSectionMarkerFormat(id, parts); err != nil {
			return err
		}
	} else if len(parts) == 3 {
		if err := validateIDDIdentifierFormat(id, parts); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("identifier '%s' must have 2 or 3 parts, found %d parts", id, len(parts))
	}

	return nil
}

func validateSectionMarkerFormat(id string, parts []string) error {
	typePart := strings.ToUpper(parts[0])
	if typePart != "PATTERN" && typePart != "WALK" {
		return fmt.Errorf("identifier '%s' has invalid TYPE '%s' for section marker, expected PATTERN or WALK", id, parts[0])
	}

	numberPart := parts[1]
	for _, c := range numberPart {
		if c < '0' || c > '9' {
			return fmt.Errorf("identifier '%s' has invalid NUMBER '%s', must be digits only", id, numberPart)
		}
	}
	return fmt.Errorf("identifier '%s' uses internal section marker TYPE '%s', not a valid IDD identifier (PATTERN-*, WALK-* are internal section markers)", id, typePart)
}

func validateIDDIdentifierFormat(id string, parts []string) error {
	validIDDTypes := map[string]bool{
		"SPEC":     true,
		"CONTRACT": true,
		"TEST":     true,
		"DESIGN":   true,
	}

	typePart := strings.ToUpper(parts[0])
	if !validIDDTypes[typePart] {
		return fmt.Errorf("identifier '%s' has invalid TYPE '%s', expected SPEC, CONTRACT, TEST, or DESIGN", id, parts[0])
	}

	modulePart := parts[1]
	for _, c := range modulePart {
		if (c < 'A' || c > 'Z') && (c < '0' || c > '9') && c != '_' {
			return fmt.Errorf("identifier '%s' has invalid MODULE '%s', must be uppercase letters, digits, or underscore only", id, modulePart)
		}
	}

	numberPart := parts[2]
	for _, c := range numberPart {
		if c < '0' || c > '9' {
			return fmt.Errorf("identifier '%s' has invalid NUMBER '%s', must be digits only", id, numberPart)
		}
	}
	return nil
}

var idPattern = regexp.MustCompile(`\b(SPEC-[A-Z0-9_]+-[0-9]+|CONTRACT-[A-Z0-9_]+-[0-9]+|TEST-[A-Z0-9_]+-[0-9]+|DESIGN-[A-Z0-9_]+-[0-9]+|PATTERN-[A-Z0-9_]+-[0-9]+|WALK-[A-Z0-9_]+-[0-9]+)\b`)

func extractIDDRefs(line string) []string {
	var refs []string
	matches := idPattern.FindAllStringSubmatch(line, -1)
	for _, m := range matches {
		if len(m) > 1 {
			refs = append(refs, m[1])
		}
	}
	return refs
}

// @implement SPEC-INTERNAL_COLLECTOR-012
func GetExpectedFilename(idType string) string {
	switch strings.ToUpper(idType) {
	case "SPEC":
		return "spec.md"
	case "CONTRACT":
		return "contract.md"
	case "TEST":
		return "testing.md"
	case "DESIGN":
		return "design.md"
	}
	return ""
}

// @implement SPEC-INTERNAL_COLLECTOR-013
func ValidateDocumentStructure(filePath string, idType string) error {
	filename := strings.ToLower(filepath.Base(filePath))

	if isRootDocFile(filePath) {
		return nil
	}

	expected := GetExpectedFilename(idType)
	if expected != "" && filename != expected {
		return fmt.Errorf("file named '%s' but should be '%s' for type %s (file: %s)",
			filename, expected, idType, filePath)
	}
	return nil
}

func isRootDocFile(filePath string) bool {
	dir := filepath.Dir(filePath)
	filename := strings.ToLower(filepath.Base(filePath))
	if dir == "docs" || dir == "docs/" {
		return true
	}
	if strings.HasSuffix(dir, "/docs") && !strings.Contains(dir, "/docs/") {
		return true
	}
	_ = filename
	return false
}
