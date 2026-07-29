// Package collector provides code identifier collection functionality.
package collector

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jingxu9x/idd-cli/internal/config"
	"github.com/jingxu9x/idd-cli/internal/model"
	"github.com/jingxu9x/idd-cli/pkg/pattern"
)

// CodeCollector collects IDD annotations from supported source files, extracting
// @implement, @test, and @test-contract annotations bound to declarations.
//
// @implement SPEC-INTERNAL_COLLECTOR-003
type CodeCollector struct {
	cfg      *config.Config
	analyses []*SourceAnalysis
}

// NewCodeCollector creates a new CodeCollector with the given configuration.
//
// @implement SPEC-INTERNAL_COLLECTOR-004
func NewCodeCollector(cfg *config.Config) *CodeCollector {
	return &CodeCollector{cfg: cfg}
}

// Collect collects IDD identifiers from code annotations in source files at the
// target path. If targetPath is a directory, it recursively walks configured
// source files.
//
// @implement SPEC-INTERNAL_COLLECTOR-024
func (c *CodeCollector) Collect(ctx context.Context, targetPath string) (*model.IdentifierSet, error) {
	set, _, err := c.CollectWithErrors(ctx, targetPath)
	return set, err
}

// CollectWithErrors collects attached annotations and returns source-parse
// findings separately so callers can add them to the normal validation report.
// @implement SPEC-INTERNAL_COLLECTOR-024
func (c *CodeCollector) CollectWithErrors(
	ctx context.Context,
	targetPath string,
) (*model.IdentifierSet, []*model.ValidationError, error) {
	set := model.NewIdentifierSet()
	c.analyses = nil
	var validationErrors []*model.ValidationError

	resolvedTarget := c.cfg.ResolvePath(targetPath)
	info, err := os.Stat(resolvedTarget)
	if err != nil {
		return set, validationErrors, nil
	}

	if info.IsDir() {
		err := filepath.Walk(resolvedTarget, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			displayPath := c.cfg.DisplayPath(path)
			if c.shouldIgnore(displayPath) {
				return nil
			}
			if !c.matchesConfiguredSource(displayPath) {
				return nil
			}
			if !SupportedSourcePath(displayPath) {
				validationErrors = append(validationErrors, unsupportedSourceFinding(displayPath))
				return nil
			}
			fileErrors, collectErr := c.collectFile(displayPath, set)
			if collectErr != nil {
				return collectErr
			}
			validationErrors = append(validationErrors, fileErrors...)
			return nil
		})
		if err != nil {
			return set, validationErrors, err
		}
	} else {
		displayPath := c.cfg.DisplayPath(resolvedTarget)
		if !SupportedSourcePath(displayPath) {
			validationErrors = append(validationErrors, unsupportedSourceFinding(displayPath))
			return set, validationErrors, nil
		}
		fileErrors, collectErr := c.collectFile(displayPath, set)
		if collectErr != nil {
			return set, validationErrors, collectErr
		}
		validationErrors = append(validationErrors, fileErrors...)
	}

	return set, validationErrors, nil
}

func unsupportedSourceFinding(path string) *model.ValidationError {
	extension := filepath.Ext(path)
	if extension == "" {
		extension = "<none>"
	}
	return &model.ValidationError{
		Rule: "source-parse",
		Message: fmt.Sprintf(
			"source extension %q has no configured AST grammar; remove it from code.patterns or use a supported language",
			extension,
		),
		Source: path,
		Code:   "unsupported-extension",
	}
}

// Analyses returns the syntax-tree source models from the most recent
// collection.
// @implement SPEC-INTERNAL_COLLECTOR-024
func (c *CodeCollector) Analyses() []*SourceAnalysis {
	return append([]*SourceAnalysis(nil), c.analyses...)
}

// collectFile parses one source file and extracts annotations attached to
// declarations.
func (c *CodeCollector) collectFile(
	path string,
	set *model.IdentifierSet,
) ([]*model.ValidationError, error) {
	content, err := os.ReadFile(c.cfg.ResolvePath(path))
	if err != nil {
		return nil, err
	}
	analysis, err := AnalyzeSourceWithAnnotations(path, content, c.cfg.Code.Annotations)
	if err != nil {
		return nil, err
	}
	c.analyses = append(c.analyses, analysis)
	var validationErrors []*model.ValidationError
	for _, parseError := range analysis.ParseErrors {
		validationErrors = append(validationErrors, &model.ValidationError{
			Rule:    "source-parse",
			Message: fmt.Sprintf("%s source cannot be bound safely: %s", analysis.Language, parseError.Message),
			Source:  fmt.Sprintf("%s:%d", path, parseError.Line),
			Code:    analysis.Language,
		})
	}
	if len(analysis.ParseErrors) > 0 {
		return validationErrors, nil
	}

	declarations := make(map[string]SourceDeclaration)
	for _, declaration := range analysis.Declarations {
		declarations[sourceDeclarationKey(declaration.Name, declaration.Line)] = declaration
	}
	lines := strings.Split(string(content), "\n")
	for _, annotation := range analysis.Annotations {
		if annotation.Ignored || !annotation.Attached {
			continue
		}
		idType := model.TypeSpec
		if annotation.Kind == "test" || annotation.Kind == "test-contract" {
			idType = model.TypeTest
		}
		declaration := declarations[sourceDeclarationKey(annotation.Declaration, annotation.DeclarationLine)]
		contextStart := annotation.Line - 3
		if contextStart < 0 {
			contextStart = 0
		}
		contextEnd := declaration.EndLine + 1
		if contextEnd > len(lines) {
			contextEnd = len(lines)
		}
		contextText := strings.Join(lines[contextStart:contextEnd], "\n")
		declarationDescription := fmt.Sprintf("[%s: %s]", declaration.Kind, declaration.Name)
		for _, ref := range annotation.Refs {
			if pattern.ValidateIdentifierFormat(ref) != nil {
				continue
			}
			if (idType == model.TypeSpec && !strings.HasPrefix(ref, "SPEC-")) ||
				(idType == model.TypeTest && !strings.HasPrefix(ref, "TEST-")) {
				continue
			}
			ann := model.NewAnnotationWithComment(
				idType,
				ref,
				path,
				annotation.Raw,
				contextText,
				declarationDescription,
				annotation.Line,
			)
			id := ann.ToIdentifier()
			id.SetOrigin(model.OriginCode)
			id.Kind = annotation.Kind
			if id.Kind == "test-contract" {
				id.Kind = "contract"
			}
			set.Add(id)
		}
	}
	return validationErrors, nil
}

func sourceDeclarationKey(name string, line int) string {
	return fmt.Sprintf("%s:%d", name, line)
}

func (c *CodeCollector) matchesConfiguredSource(path string) bool {
	if len(c.cfg.Code.Patterns) == 0 {
		return true
	}
	slashPath := filepath.ToSlash(path)
	for _, configuredPattern := range c.cfg.Code.Patterns {
		if matchGlob(configuredPattern, slashPath) {
			return true
		}
		if strings.HasPrefix(configuredPattern, "**/*") &&
			strings.HasSuffix(slashPath, strings.TrimPrefix(configuredPattern, "**/*")) {
			return true
		}
	}
	return false
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
