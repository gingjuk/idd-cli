// Package engine provides IDD validation engine.

// Spec: docs/internal/engine/spec.md
// Contract: docs/internal/engine/contract.md

package engine

import (
	"fmt"
	"regexp"
	"strings"
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
		for i, line := range lines {
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
					fmt.Sprintf("annotation should have space before identifier (file: %s, line: %d)", path, i+1),
					fmt.Sprintf("%s:%d", path, i+1),
					"",
					"",
				)
			}

			if hasImplement && !specIdRegex.MatchString(line) {
				e.result.AddError(
					"annotation-missing-identifier",
					fmt.Sprintf("@implement missing identifier (file: %s, line: %d)", path, i+1),
					fmt.Sprintf("%s:%d", path, i+1),
					"",
					"",
				)
			}
			if hasTestContract && !specIdRegex.MatchString(line) {
				e.result.AddError(
					"annotation-missing-identifier",
					fmt.Sprintf("@test-contract missing identifier (file: %s, line: %d)", path, i+1),
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
		for i := 0; i < len(lines)-1; i++ {
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
						fmt.Sprintf("%s annotations should not be split across lines (file: %s, lines: %d, %d)", currentAnnotation, path, i+1, j+1),
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
