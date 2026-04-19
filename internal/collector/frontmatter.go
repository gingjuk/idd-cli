package collector

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Marker struct {
	ID       string `yaml:"id"`
	Name     string `yaml:"name"`
	Describe string `yaml:"describe"`
}

type Frontmatter struct {
	Markers []Marker `yaml:"markers"`
}

func ParseFrontmatter(content string) (*Frontmatter, error) {
	lines := strings.Split(content, "\n")
	startIdx, endIdx := -1, -1
	inCodeBlock := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
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

func ValidateFrontmatterMarkers(fm *Frontmatter, content string, filePath string) []string {
	var errors []string
	if fm == nil {
		return errors
	}

	definedMarkers := extractDefinedMarkers(content)
	referencedMarkers := extractReferencedMarkers(content)

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

	if errs := ValidateMarkerFormatting(content, filePath); len(errs) > 0 {
		errors = append(errors, errs...)
	}

	return errors
}

func ValidateMarkerFormatting(content string, filePath string) []string {
	var errors []string
	lines := strings.Split(content, "\n")
	inCodeBlock := false
	inFrontmatter := false

	for lineNum, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
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
	backtickPattern := regexp.MustCompile("`([^`]+)`")
	lineWithoutBackticks := backtickPattern.ReplaceAllString(line, "")

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

var idPattern = regexp.MustCompile(`\b(SPEC-[A-Z]+-[0-9]+|CONTRACT-[A-Z]+-[0-9]+|TEST-[A-Z]+-[0-9]+|DESIGN-[A-Z]+-[0-9]+)\b`)

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

func ExtractModuleName(filePath string) string {
	parts := strings.Split(filepath.Dir(filePath), string(filepath.Separator))
	for i, part := range parts {
		if part == "docs" && i+1 < len(parts) {
			module := parts[i+1]
			if module != "" && module != "." {
				return strings.ToUpper(module)
			}
		}
	}
	return ""
}

var knownAbbrevs = map[string]string{
	"BE":  "BACKEND",
	"FE":  "FRONTEND",
	"E2E": "E2E",
}

func isKnownAbbreviation(short, long string) bool {
	for abbrev, full := range knownAbbrevs {
		if strings.ToUpper(short) == abbrev && strings.ToUpper(long) == full {
			return true
		}
	}
	return false
}

func ValidateModulePrefix(id string, filePath string) error {
	moduleFromPath := ExtractModuleName(filePath)
	if moduleFromPath == "" {
		return nil
	}

	parts := strings.Split(id, "-")
	if len(parts) < 2 {
		return nil
	}
	idModule := parts[1]

	if idModule != moduleFromPath {
		if !isKnownAbbreviation(idModule, moduleFromPath) {
			return fmt.Errorf(
				"identifier module '%s' does not match directory '%s' (id: %s, file: %s)",
				idModule, moduleFromPath, id, filePath,
			)
		}
	}
	return nil
}
