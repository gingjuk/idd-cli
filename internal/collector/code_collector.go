package collector

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jingxu9x/idd-link-validator/internal/config"
	"github.com/jingxu9x/idd-link-validator/internal/model"
	"github.com/jingxu9x/idd-link-validator/pkg/pattern"
)

type CodeCollector struct {
	cfg *config.Config
}

func NewCodeCollector(cfg *config.Config) *CodeCollector {
	return &CodeCollector{cfg: cfg}
}

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

func (c *CodeCollector) collectFile(path string, set *model.IdentifierSet) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		for _, pat := range pattern.AnnotationPatterns {
			matches := pat.Regex.FindAllStringSubmatch(line, -1)
			for _, m := range matches {
				if len(m) > 1 {
					refs := pattern.SplitAnnotationRefs(m[1])
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
						set.Add(id)
					}
				}
			}
		}
	}

	return nil
}

func matchGlob(pattern, fullPath string) bool {
	if pattern == "**" {
		return true
	}
	if strings.Contains(fullPath, pattern) {
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

var funcDeclRegex = regexp.MustCompile(`^func\s+(?:\([^)]+\)\s+)?(\w+)\s*\(`)

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
