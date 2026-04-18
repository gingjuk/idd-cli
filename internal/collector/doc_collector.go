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
	var fileRefs []string

	for i, line := range lines {
		refs := pattern.ExtractIDDReferences(line)
		for _, ref := range refs {
			idTypeStr := pattern.GetIdentifierType(ref)
			idType, _ := model.ParseIdentifierType(idTypeStr)
			if idType == "" {
				continue
			}

			if err := ValidateModulePrefix(ref, path); err != nil {
				errors = append(errors, &model.ValidationError{
					Rule:    "module-prefix-mismatch",
					Message: err.Error(),
					Source:  path,
					Link:    ref,
				})
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
		}
	}

	return errors
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
