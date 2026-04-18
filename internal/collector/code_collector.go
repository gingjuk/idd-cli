package collector

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/yourorg/idd-cli/internal/config"
	"github.com/yourorg/idd-cli/internal/model"
	"github.com/yourorg/idd-cli/pkg/pattern"
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
					ref := m[1]
					idType, _ := model.ParseIdentifierType(pat.Type)

					context := ""
					start := i - 2
					if start < 0 {
						start = 0
					}
					end := i + 3
					if end > len(lines) {
						end = len(lines)
					}
					context = strings.Join(lines[start:end], "\n")

					ann := model.NewAnnotation(idType, ref, path, strings.TrimSpace(line), context, i+1)
					id := ann.ToIdentifier()
					set.Add(id)
				}
			}
		}
	}

	return nil
}
