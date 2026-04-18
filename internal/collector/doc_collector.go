package collector

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/yourorg/idd-link-validator/internal/config"
	"github.com/yourorg/idd-link-validator/internal/model"
	"github.com/yourorg/idd-link-validator/pkg/pattern"
)

type DocCollector struct {
	cfg *config.Config
}

func NewDocCollector(cfg *config.Config) *DocCollector {
	return &DocCollector{cfg: cfg}
}

func (c *DocCollector) Collect(ctx context.Context, targetPath string) (*model.IdentifierSet, error) {
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
			if filepath.Ext(path) == ".md" {
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

func (c *DocCollector) collectFile(path string, set *model.IdentifierSet) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var fileRefs []string

	for i, line := range lines {
		refs := pattern.ExtractIDDReferences(line)
		for _, ref := range refs {
			idType, _ := model.ParseIdentifierType(pattern.GetIdentifierType(ref))
			if idType == "" {
				continue
			}

			if !set.Has(ref) {
				title := c.extractTitle(string(content), ref)
				id := model.NewIdentifier(ref, idType, title, path, i+1)
				id.RawRef = ref
				set.Add(id)
			}
			fileRefs = append(fileRefs, ref)
		}
	}

	return nil
}

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
