// Package collector provides code identifier collection functionality.

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

// CodeCollector collects IDD annotations from source code files (Go, TypeScript,
// JavaScript), extracting @implement, @test, and @test-contract annotations.
//
// @implement SPEC-INTERNAL_COLLECTOR-003
type CodeCollector struct {
	cfg *config.Config
}

// NewCodeCollector creates a new CodeCollector with the given configuration.
//
// @implement SPEC-INTERNAL_COLLECTOR-004
func NewCodeCollector(cfg *config.Config) *CodeCollector {
	return &CodeCollector{cfg: cfg}
}

// Collect collects IDD identifiers from code annotations in source files at the
// target path. If targetPath is a directory, recursively walks to find all
// .go, .ts, .tsx, .js files.
//
// @implement SPEC-INTERNAL_COLLECTOR-024
func (c *CodeCollector) Collect(ctx context.Context, targetPath string) (*model.IdentifierSet, error) {
	set := model.NewIdentifierSet()

	info, err := os.Stat(targetPath)
	if err != nil {
		return set, nil
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
			ext := filepath.Ext(path)
			if ext == ".go" || ext == ".ts" || ext == ".tsx" || ext == ".js" {
				_ = c.collectFile(path, set)
			}
			return nil
		})
		if err != nil {
			return set, err
		}
	} else {
		_ = c.collectFile(targetPath, set)
	}

	return set, nil
}

// collectFile scans a source file for IDD annotations and extracts identifiers.
func (c *CodeCollector) collectFile(path string, set *model.IdentifierSet) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")

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

		// Check for single-line idd:ignore
		if strings.Contains(line, "// idd:ignore") || strings.Contains(line, "//idd:ignore") {
			continue
		}

		// Skip if inside ignore scope
		if ignoreScope {
			continue
		}

		for _, pat := range pattern.AnnotationPatterns {
			matches := pat.Regex.FindAllStringSubmatch(line, -1)
			for _, m := range matches {
				if len(m) > 1 {
					refs := SplitAnnotationRefs(m[1])
					idType, _ := model.ParseIdentifierType(pat.Type)

					ctx := ""
					start := i - 2
					if start < 0 {
						start = 0
					}
					end := i + 3
					if end > len(lines) {
						end = len(lines)
					}
					ctx = strings.Join(lines[start:end], "\n")

					funcComment := extractFunctionComment(lines, i)

					for _, ref := range refs {
						ann := model.NewAnnotationWithComment(idType, ref, path, strings.TrimSpace(line), ctx, funcComment, i+1)
						id := ann.ToIdentifier()
						id.SetOrigin(model.OriginCode)
						id.Kind = strings.TrimPrefix(pat.Prefix, "@")
						if id.Kind == "test-contract" {
							id.Kind = "contract"
						}
						set.Add(id)
					}
				}
			}
		}
	}

	return nil
}

// matchGlob matches a glob pattern against a full file path.
// Supports ** for matching any number of directories.
func matchGlob(pattern, fullPath string) bool {
	if pattern == "**" {
		return true
	}
	parts := strings.Split(pattern, "/**")
	if len(parts) == 2 {
		prefix := parts[0]
		if prefix != "" && (strings.HasPrefix(fullPath, prefix+"/") || strings.Contains(fullPath, "/"+prefix+"/")) {
			return true
		}
	}
	matched, _ := filepath.Match(pattern, fullPath)
	return matched
}

// shouldIgnore checks if a path should be ignored based on configured ignore patterns.
func (c *CodeCollector) shouldIgnore(path string) bool {
	if len(c.cfg.Code.IgnorePaths) == 0 {
		return false
	}
	for _, ignore := range c.cfg.Code.IgnorePaths {
		if matchGlob(ignore, path) {
			return true
		}
	}
	return false
}

// funcDeclRegex matches function declarations in Go code.
var funcDeclRegex = regexp.MustCompile(`^func\s+(?:\([^)]+\)\s+)?(\w+)\s*\(`)

// extractFunctionComment extracts function name and preceding comments as context
// for code annotations.
//
// extractFunctionComment extracts function name and preceding comments as context.
func extractFunctionComment(lines []string, annotationLine int) string {
	if annotationLine >= len(lines) {
		return ""
	}
	searchStart := annotationLine + 1
	if searchStart >= len(lines) {
		return ""
	}
	endSearch := annotationLine + 10
	if endSearch > len(lines) {
		endSearch = len(lines)
	}
	var commentLines []string
	for i := searchStart; i < endSearch; i++ {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "//") {
			commentText := strings.TrimPrefix(line, "//")
			commentText = strings.TrimSpace(commentText)
			if commentText != "" {
				commentLines = append(commentLines, commentText)
			}
		} else if strings.HasPrefix(line, "/*") {
			continue
		} else if line == "" {
			continue
		} else {
			break
		}
	}
	for i := searchStart; i < endSearch; i++ {
		line := strings.TrimSpace(lines[i])
		if match := funcDeclRegex.FindStringSubmatch(line); len(match) > 1 {
			funcName := match[1]
			if len(commentLines) > 0 {
				return strings.Join(commentLines, " ") + " [function: " + funcName + "]"
			}
			return "[function: " + funcName + "]"
		}
	}
	if len(commentLines) > 0 {
		return strings.Join(commentLines, " ")
	}
	return ""
}

// SplitAnnotationRefs splits comma-separated IDD references from an annotation
// and trims whitespace. Used to handle multiple references in a single
// annotation like `@implement` `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-002`.
//
// @implement SPEC-INTERNAL_COLLECTOR-025
func SplitAnnotationRefs(s string) []string {
	refs := strings.Split(s, ",")
	for i, ref := range refs {
		refs[i] = strings.TrimSpace(ref)
	}
	return refs
}
