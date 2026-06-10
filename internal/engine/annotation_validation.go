// Package engine provides IDD validation engine.

// Spec: docs/internal/engine/spec.md
// Contract: docs/internal/engine/contract.md

package engine

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jingxu9x/idd-cli/internal/model"
)

func (e *Engine) validateAnnotationIdentifiers() {
	implementRegex := regexp.MustCompile(`(?i)@implement\b`)
	testRegex := regexp.MustCompile(`(?i)@test\b`)
	testContractRegex := regexp.MustCompile(`(?i)@test-contract\b`)
	specIdRegex := regexp.MustCompile(`[A-Z0-9_]+-[A-Z0-9_]+-[0-9]+`)
	// Check for space between annotation and identifier: "// @implement ID" not "// @implementID"
	// Supports annotations like @test-contract with hyphens
	annotationWithSpaceRegex := regexp.MustCompile(`(?i)^\s*//\s*@[\w-]+\s+[A-Z]`)
	// Only match lines where annotation appears at start of comment (after // and optional space)
	annotationAtStartRegex := regexp.MustCompile(`(?i)^\s*//\s*@(implement|test|test-contract)\b`)

	e.walkCodeFiles(func(path string, lines []string) {
		ignoreScope := false
		for i, line := range lines {
			// Honor // idd:ignore start/end scope markers so callers can
			// exclude fixture content (e.g. raw string fixtures in tests) from
			// the identifier check.
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

			trimmed := strings.TrimSpace(line)
			if !strings.HasPrefix(trimmed, "//") {
				continue
			}
			// Only check lines where annotation appears at start of comment
			if !annotationAtStartRegex.MatchString(line) {
				continue
			}
			hasImplement := implementRegex.MatchString(line)
			hasTest := testRegex.MatchString(line)
			hasTestContract := testContractRegex.MatchString(line)
			if !hasImplement && !hasTest && !hasTestContract {
				continue
			}

			// Check for missing space between annotation and identifier
			hasProperSpacing := annotationWithSpaceRegex.MatchString(line)
			hasIdentifier := specIdRegex.MatchString(line)
			if hasIdentifier && !hasProperSpacing {
				e.result.AddError(
					"annotation-format",
					"annotation should have a space before identifier (e.g. `// @implement SPEC-FOO-001`)", // idd:ignore
					fmt.Sprintf("%s:%d", path, i+1),
					"",
					"",
				)
			}

			if hasImplement && !specIdRegex.MatchString(line) {
				e.result.AddError(
					"annotation-missing-identifier",
					"@implement missing identifier (e.g. `// @implement SPEC-FOO-001`)", // idd:ignore
					fmt.Sprintf("%s:%d", path, i+1),
					"",
					"",
				)
			}
			if hasTestContract && !specIdRegex.MatchString(line) {
				e.result.AddError(
					"annotation-missing-identifier",
					"@test-contract missing identifier (e.g. `// @test-contract TEST-FOO-001`)", // idd:ignore
					fmt.Sprintf("%s:%d", path, i+1),
					"",
					"",
				)
			}
		}
	})
}

func (e *Engine) validateConsecutiveAnnotations() {
	annotationPrefixes := []string{"@implement", "@test", "@test-contract"}
	annotationRegex := regexp.MustCompile(`(?i)^\s*//\s*@(\w+)\s+`)

	e.walkCodeFiles(func(path string, lines []string) {
		ignoreScope := false
		for i := 0; i < len(lines)-1; i++ {
			// Honor // idd:ignore start/end scope markers.
			if strings.Contains(lines[i], "// idd:ignore start") || strings.Contains(lines[i], "//idd:ignore-start") {
				ignoreScope = true
				continue
			}
			if strings.Contains(lines[i], "// idd:ignore end") || strings.Contains(lines[i], "//idd:ignore-end") {
				ignoreScope = false
				continue
			}
			if ignoreScope {
				continue
			}

			currentLine := strings.TrimSpace(lines[i])
			if !strings.HasPrefix(currentLine, "//") {
				continue
			}

			currentAnnotation := ""
			if match := annotationRegex.FindStringSubmatch(currentLine); len(match) > 1 {
				currentAnnotation = "@" + strings.ToLower(match[1])
			}
			if currentAnnotation == "" {
				continue
			}

			prefixMatched := false
			for _, prefix := range annotationPrefixes {
				if strings.EqualFold(currentAnnotation, prefix) {
					prefixMatched = true
					break
				}
			}
			if !prefixMatched {
				continue
			}

			for j := i + 1; j < len(lines); j++ {
				nextLine := strings.TrimSpace(lines[j])
				if nextLine == "" {
					continue
				}
				if !strings.HasPrefix(nextLine, "//") {
					break
				}

				nextAnnotation := ""
				if match := annotationRegex.FindStringSubmatch(nextLine); len(match) > 1 {
					nextAnnotation = "@" + strings.ToLower(match[1])
				}
				if nextAnnotation == "" {
					continue
				}

				if strings.EqualFold(nextAnnotation, currentAnnotation) {
					e.result.AddError(
						"annotation-consecutive-line",
						fmt.Sprintf("%s annotations should not be split across lines (lines %d and %d); use comma-separated identifiers on one line", currentAnnotation, i+1, j+1),
						fmt.Sprintf("%s:%d", path, i+1),
						"",
						"",
					)
					break
				}
				break
			}
		}
	})
}

// validateDuplicateIDs checks that no identifier is defined in more than one
// pkg directory — both on the doc side and the code side.
func (e *Engine) validateDuplicateIDs(ids *model.IdentifierSet) {
	e.reportDuplicates(ids.DuplicateDocGroups(), "doc")
	e.reportDuplicates(ids.DuplicateCodeGroups(), "code")
}

func (e *Engine) reportDuplicates(groups [][]*model.Identifier, side string) {
	for _, group := range groups {
		dirs := make([]string, len(group))
		for i, id := range group {
			dirs[i] = filepath.Dir(id.Source)
		}
		for i, id := range group {
			others := make([]string, 0, len(dirs)-1)
			for j, d := range dirs {
				if j != i {
					others = append(others, d)
				}
			}
			suggestion := suggestRename(id.ID, id.Source)
			msg := fmt.Sprintf("%s (%s) also defined in: %s", id.ID, side, strings.Join(others, ", "))
			if suggestion != "" {
				msg += fmt.Sprintf("; consider renaming to %s", suggestion)
			}
			e.result.AddError("duplicate-id", msg, id.Source, id.ID, "")
		}
	}
}

// suggestRename derives a rename suggestion by replacing the MODULE segment of id
// with a name derived from the source file's directory path.
func suggestRename(id, sourcePath string) string {
	parts := strings.SplitN(id, "-", 3)
	if len(parts) != 3 {
		return ""
	}
	newModule := moduleFromPath(sourcePath)
	if newModule == "" || newModule == parts[1] {
		return ""
	}
	return parts[0] + "-" + newModule + "-" + parts[2]
}

// moduleFromPath derives an uppercase module name from a file path.
// For doc paths it uses directory components after "docs/".
// For code paths it uses the last 1-2 directory components.
func moduleFromPath(path string) string {
	dir := filepath.ToSlash(filepath.Dir(path))
	segs := strings.Split(dir, "/")
	for i, seg := range segs {
		if seg == "docs" && i+1 < len(segs) {
			first := pathSegToModule(segs[i+1])
			if i+2 < len(segs) && segs[i+2] != "" {
				return first + "_" + pathSegToModule(segs[i+2])
			}
			return first
		}
	}
	var tail []string
	for j := len(segs) - 1; j >= 0 && len(tail) < 2; j-- {
		if segs[j] != "" && segs[j] != "." {
			tail = append([]string{segs[j]}, tail...)
		}
	}
	result := make([]string, len(tail))
	for i, s := range tail {
		result[i] = pathSegToModule(s)
	}
	return strings.Join(result, "_")
}

func pathSegToModule(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, "-", "_"))
}
