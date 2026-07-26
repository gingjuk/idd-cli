// Package engine provides syntax-tree-backed source validation.
package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jingxu9x/idd-cli/internal/collector"
	"github.com/jingxu9x/idd-cli/pkg/pattern"
)

func (e *Engine) ensureSourceAnalyses() []*collector.SourceAnalysis {
	if e.sourceAnalysesSet {
		return e.sourceAnalyses
	}

	seen := make(map[string]bool)
	e.walkCodeFiles(func(path string, _ []string) {
		if seen[path] || !collector.SupportedSourcePath(path) {
			return
		}
		seen[path] = true
		content, err := os.ReadFile(path)
		if err != nil {
			e.result.AddError("source-parse", err.Error(), path, "", "")
			return
		}
		analysis, err := collector.AnalyzeSourceWithAnnotations(
			path,
			content,
			e.cfg.Code.Annotations,
		)
		if err != nil {
			e.result.AddError("source-parse", err.Error(), path, "", "")
			return
		}
		e.sourceAnalyses = append(e.sourceAnalyses, analysis)
		for _, parseError := range analysis.ParseErrors {
			e.result.AddError(
				"source-parse",
				fmt.Sprintf("%s source cannot be bound safely: %s", analysis.Language, parseError.Message),
				fmt.Sprintf("%s:%d", path, parseError.Line),
				"",
				analysis.Language,
			)
		}
	})
	sort.Slice(e.sourceAnalyses, func(i, j int) bool {
		return e.sourceAnalyses[i].Path < e.sourceAnalyses[j].Path
	})
	e.sourceAnalysesSet = true
	return e.sourceAnalyses
}

func (e *Engine) validateSourcePublicAnnotations() {
	for _, analysis := range e.ensureSourceAnalyses() {
		if len(analysis.ParseErrors) > 0 {
			continue
		}
		for _, annotation := range analysis.Annotations {
			if annotation.Ignored || annotation.Kind != "implement" || annotation.Attached {
				continue
			}
			e.result.AddError(
				"annotation-placement",
				fmt.Sprintf("%s is not attached to a declaration", sourceAnnotationPrefix(annotation)),
				fmt.Sprintf("%s:%d", analysis.Path, annotation.Line),
				"",
				"",
			)
		}
		for _, declaration := range analysis.Declarations {
			if declaration.Ignored {
				continue
			}
			implementAnnotations := sourceAnnotationsByKind(declaration, "implement")
			if len(implementAnnotations) > 1 {
				prefix := sourceAnnotationPrefix(implementAnnotations[1])
				e.result.AddError(
					"duplicate-annotation",
					fmt.Sprintf(
						"%s %s has multiple %s annotations; use comma-separated identifiers on one annotation",
						declaration.Kind,
						declaration.Name,
						prefix,
					),
					fmt.Sprintf("%s:%d", analysis.Path, implementAnnotations[1].Line),
					"",
					"",
				)
			}
			if sourceValidationTestPath(analysis.Path) ||
				!declaration.Public || declaration.TestKind != "" ||
				declaration.Name == "main" || declaration.Name == "init" {
				continue
			}
			if len(implementAnnotations) == 0 {
				prefix := e.configuredSourceAnnotationPrefix("spec", "@implement")
				e.result.AddError(
					"public-func-annotation",
					fmt.Sprintf(
						"public %s %s missing %s",
						declaration.Kind,
						declaration.Name,
						prefix,
					),
					fmt.Sprintf("%s:%d", analysis.Path, declaration.Line),
					"",
					"",
				)
			}
		}
	}
}

func sourceValidationTestPath(path string) bool {
	name := strings.ToLower(filepath.Base(path))
	extension := filepath.Ext(name)
	base := strings.TrimSuffix(name, extension)
	return strings.HasSuffix(base, "_test") ||
		strings.Contains(name, ".test.") ||
		strings.Contains(name, ".spec.") ||
		strings.HasPrefix(name, "test_") ||
		strings.HasSuffix(base, "test")
}

func (e *Engine) validateSourceTestAnnotations() {
	for _, analysis := range e.ensureSourceAnalyses() {
		if len(analysis.ParseErrors) > 0 {
			continue
		}
		for _, annotation := range analysis.Annotations {
			if annotation.Ignored ||
				annotation.Kind != "test" && annotation.Kind != "test-contract" ||
				annotation.Attached {
				continue
			}
			e.result.AddError(
				"annotation-placement",
				fmt.Sprintf("%s is not attached to a test declaration", sourceAnnotationPrefix(annotation)),
				fmt.Sprintf("%s:%d", analysis.Path, annotation.Line),
				"",
				"",
			)
		}
		for _, declaration := range analysis.Declarations {
			if declaration.Ignored {
				continue
			}
			testAnnotations := append(
				sourceAnnotationsByKind(declaration, "test"),
				sourceAnnotationsByKind(declaration, "test-contract")...,
			)
			if declaration.TestKind == "" {
				for _, annotation := range testAnnotations {
					e.result.AddError(
						"test-annotation",
						fmt.Sprintf(
							"%s is attached to non-test %s %s",
							sourceAnnotationPrefix(annotation),
							declaration.Kind,
							declaration.Name,
						),
						fmt.Sprintf("%s:%d", analysis.Path, annotation.Line),
						"",
						"",
					)
				}
				continue
			}

			behaviorAnnotations := sourceAnnotationsByKind(declaration, "test")
			contractAnnotations := sourceAnnotationsByKind(declaration, "test-contract")
			expected := "test"
			expectedRole := "test"
			matching := behaviorAnnotations
			if declaration.TestKind == "contract" || len(contractAnnotations) > 0 {
				expected = "test-contract"
				expectedRole = "test_contract"
				matching = contractAnnotations
			}
			if len(matching) == 0 {
				expectedPrefix := e.configuredSourceAnnotationPrefix(
					expectedRole,
					"@"+expected,
				)
				e.result.AddError(
					"test-annotation",
					fmt.Sprintf(
						"test %s missing %s annotation",
						declaration.Name,
						expectedPrefix,
					),
					fmt.Sprintf("%s:%d", analysis.Path, declaration.Line),
					"",
					"",
				)
				continue
			}
			if len(behaviorAnnotations) > 0 && len(contractAnnotations) > 0 {
				e.result.AddError(
					"test-annotation",
					fmt.Sprintf(
						"test %s mixes %s and %s annotations",
						declaration.Name,
						sourceAnnotationPrefix(behaviorAnnotations[0]),
						sourceAnnotationPrefix(contractAnnotations[0]),
					),
					fmt.Sprintf("%s:%d", analysis.Path, declaration.Line),
					"",
					"",
				)
				continue
			}
			if len(matching) > 1 {
				prefix := sourceAnnotationPrefix(matching[1])
				e.result.AddError(
					"duplicate-annotation",
					fmt.Sprintf(
						"test %s has multiple %s annotations; use comma-separated identifiers on one annotation",
						declaration.Name,
						prefix,
					),
					fmt.Sprintf("%s:%d", analysis.Path, matching[1].Line),
					"",
					"",
				)
			}
			for _, annotation := range matching {
				for _, ref := range annotation.Refs {
					if _, ok := e.graph.GetNode(ref); ok {
						continue
					}
					e.result.AddError(
						"test-annotation",
						fmt.Sprintf("test %s references undefined %s", declaration.Name, ref),
						fmt.Sprintf("%s:%d", analysis.Path, annotation.Line),
						ref,
						"",
					)
				}
			}
		}
	}
}

func (e *Engine) validateSourceAnnotationIdentifiers() {
	for _, analysis := range e.ensureSourceAnalyses() {
		if len(analysis.ParseErrors) > 0 {
			continue
		}
		for _, annotation := range analysis.Annotations {
			if annotation.Ignored {
				continue
			}
			source := fmt.Sprintf("%s:%d", analysis.Path, annotation.Line)
			prefix := sourceAnnotationPrefix(annotation)
			if !annotation.HasSeparator && len(annotation.Refs) > 0 {
				e.result.AddError(
					"annotation-format",
					fmt.Sprintf("%s must be followed by a space before its identifier", prefix),
					source,
					"",
					"",
				)
			}
			if len(annotation.Refs) == 0 {
				e.result.AddError(
					"annotation-missing-identifier",
					fmt.Sprintf("%s missing identifier", prefix),
					source,
					"",
					"",
				)
				continue
			}
			expectedPrefix := "SPEC-"
			if annotation.Kind == "test" || annotation.Kind == "test-contract" {
				expectedPrefix = "TEST-"
			}
			for _, ref := range annotation.Refs {
				if strings.HasPrefix(ref, expectedPrefix) &&
					pattern.ValidateIdentifierFormat(ref) == nil {
					continue
				}
				e.result.AddError(
					"annotation-invalid-identifier",
					fmt.Sprintf(
						"%s identifier %q must be a valid %s identifier",
						prefix,
						ref,
						strings.TrimSuffix(expectedPrefix, "-"),
					),
					source,
					ref,
					"",
				)
			}
		}
	}
}

func (e *Engine) validateSourceConsecutiveAnnotations() {
	for _, analysis := range e.ensureSourceAnalyses() {
		if len(analysis.ParseErrors) > 0 {
			continue
		}
		for _, declaration := range analysis.Declarations {
			if declaration.Ignored {
				continue
			}
			for _, kind := range []string{"implement", "test", "test-contract"} {
				annotations := sourceAnnotationsByKind(declaration, kind)
				if len(annotations) < 2 {
					continue
				}
				sort.Slice(annotations, func(i, j int) bool {
					return annotations[i].Line < annotations[j].Line
				})
				e.result.AddError(
					"annotation-consecutive-line",
					fmt.Sprintf(
						"%s annotations for %s %s are split across lines %d and %d; use comma-separated identifiers on one line",
						sourceAnnotationPrefix(annotations[0]),
						declaration.Kind,
						declaration.Name,
						annotations[0].Line,
						annotations[1].Line,
					),
					fmt.Sprintf("%s:%d", analysis.Path, annotations[0].Line),
					"",
					"",
				)
			}
		}
	}
}

func sourceAnnotationsByKind(
	declaration collector.SourceDeclaration,
	kind string,
) []collector.SourceAnnotation {
	annotations := make([]collector.SourceAnnotation, 0)
	for _, annotation := range declaration.Annotations {
		if !annotation.Ignored && annotation.Kind == kind {
			annotations = append(annotations, annotation)
		}
	}
	return annotations
}

func sourceAnnotationPrefix(annotation collector.SourceAnnotation) string {
	if annotation.Prefix != "" {
		return annotation.Prefix
	}
	return "@" + annotation.Kind
}

func (e *Engine) configuredSourceAnnotationPrefix(role, fallback string) string {
	if e.cfg != nil && e.cfg.Code.Annotations != nil {
		if prefix := strings.TrimSpace(e.cfg.Code.Annotations[role]); prefix != "" {
			return prefix
		}
	}
	return fallback
}
