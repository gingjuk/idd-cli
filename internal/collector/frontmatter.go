package collector

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Marker struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

type Frontmatter struct {
	Markers []Marker `yaml:"markers"`
}

func ParseFrontmatter(content string) (*Frontmatter, error) {
	lines := strings.Split(content, "\n")
	startIdx, endIdx := -1, -1

	for i, line := range lines {
		if strings.TrimSpace(line) == "---" {
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

	contentMarkers := make(map[string]bool)
	for _, line := range strings.Split(content, "\n") {
		for _, ref := range extractIDDRefs(line) {
			contentMarkers[ref] = true
		}
	}

	for _, marker := range fm.Markers {
		if !contentMarkers[marker.ID] {
			errors = append(errors, fmt.Sprintf(
				"frontmatter marker '%s' not found in content (file: %s)",
				marker.ID, filePath,
			))
		}
	}
	return errors
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
	expected := GetExpectedFilename(idType)
	if expected != "" && filename != expected {
		return fmt.Errorf("file named '%s' but should be '%s' for type %s (file: %s)",
			filename, expected, idType, filePath)
	}
	return nil
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
		return fmt.Errorf(
			"identifier module '%s' does not match directory '%s' (id: %s, file: %s)",
			idModule, moduleFromPath, id, filePath,
		)
	}
	return nil
}
