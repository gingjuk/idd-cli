// Package pattern provides IDD identifier and annotation pattern matching.

// Spec: docs/pkg/pattern/spec.md
// Contract: docs/pkg/pattern/contract.md
package pattern

import (
	"fmt"
	"regexp"
	"strings"
)

// IDDPattern defines a regex pattern for IDD identifier recognition.
// It holds the pattern type, the compiled regular expression, and an example identifier.
//
// Key Patterns:
//   - SPEC pattern: SPEC-[A-Z]+-[0-9]+
//   - TEST pattern: TEST-[A-Z]+-[0-9]+
//   - CONTRACT pattern: CONTRACT-[A-Z]+-[0-9]+
//   - DESIGN pattern: DESIGN-[A-Z]+-[0-9]+
//
// @implement SPEC-PKG_PATTERN-001
type IDDPattern struct {
	Type    string
	Regex   *regexp.Regexp
	Example string
}

var Patterns = map[string]*IDDPattern{
	"SPEC": {
		Type:    "SPEC",
		Regex:   regexp.MustCompile(`(?i)\b(SPEC-[A-Z0-9_]+-[0-9]+)\b`),
		Example: "SPEC-BE-001",
	},
	"CONTRACT": {
		Type:    "CONTRACT",
		Regex:   regexp.MustCompile(`(?i)\b(CONTRACT-[A-Z0-9_]+-[0-9]+)\b`),
		Example: "CONTRACT-BE-001",
	},
	"TEST": {
		Type:    "TEST",
		Regex:   regexp.MustCompile(`(?i)\b(TEST-[A-Z0-9_]+-[0-9]+)\b`),
		Example: "TEST-BE-001",
	},
	"DESIGN": {
		Type:    "DESIGN",
		Regex:   regexp.MustCompile(`(?i)\b(DESIGN-[A-Z0-9_]+-[0-9]+)\b`),
		Example: "DESIGN-BE-001",
	},
	"PATTERN": {
		Type:    "PATTERN",
		Regex:   regexp.MustCompile(`(?i)\b(PATTERN-[A-Z0-9_]+-[0-9]+)\b`),
		Example: "PATTERN-BE-001",
	},
	"WALK": {
		Type:    "WALK",
		Regex:   regexp.MustCompile(`(?i)\b(WALK-[A-Z0-9_]+-[0-9]+)\b`),
		Example: "WALK-BE-001",
	},
}

// AnnotationPattern defines a regex pattern for code annotations.
// It matches annotation prefixes like @implement, @test, and @test-contract,
// extracting the referenced IDD identifiers.
//
// Key Patterns:
//   - @implement → SPEC
//   - @test → TEST
//   - @test-contract → TEST
//
// @implement SPEC-PKG_PATTERN-002
type AnnotationPattern struct {
	Prefix string
	Regex  *regexp.Regexp
	Type   string
}

var AnnotationPatterns = []AnnotationPattern{
	{
		Prefix: "@implement",
		Regex:  regexp.MustCompile(`(?i)@implement\s+([A-Z][A-Z0-9_]*-[A-Z0-9_]+-[0-9]+(?:\s*,\s*[A-Z][A-Z0-9_]*-[A-Z0-9_]+-[0-9]+)*)`),
		Type:   "SPEC",
	},
	{
		Prefix: "@test",
		Regex:  regexp.MustCompile(`(?i)@test\s+([A-Z][A-Z0-9_]*-[A-Z0-9_]+-[0-9]+(?:\s*,\s*[A-Z][A-Z0-9_]*-[A-Z0-9_]+-[0-9]+)*)`),
		Type:   "TEST",
	},
	{
		Prefix: "@test-contract",
		Regex:  regexp.MustCompile(`(?i)@test-contract\s+([A-Z][A-Z0-9_]*-[A-Z0-9_]+-[0-9]+(?:\s*,\s*[A-Z][A-Z0-9_]*-[A-Z0-9_]+-[0-9]+)*)`),
		Type:   "TEST",
	},
}

// ExtractIDDReferences extracts all IDD identifier references from content.
// It searches for pattern matches in the provided text and returns a list of
// matching identifiers that are not quoted or backtick-wrapped.
//
// @implement SPEC-PKG_PATTERN-004
func ExtractIDDReferences(content string) []string {
	var refs []string
	for _, pat := range Patterns {
		matches := pat.Regex.FindAllStringSubmatch(content, -1)
		for _, m := range matches {
			if len(m) > 1 {
				id := strings.ToUpper(m[1])
				if isQuoted(content, id) {
					refs = append(refs, id)
				}
			}
		}
	}
	return refs
}

// isQuoted checks if an identifier in content is wrapped in quotes or backticks.
func isQuoted(content, id string) bool {
	idx := 0
	for {
		pos := strings.Index(content[idx:], id)
		if pos == -1 {
			return false
		}
		actualPos := idx + pos
		if actualPos > 0 && (content[actualPos-1] == '"' || content[actualPos-1] == '`') {
			return true
		}
		endPos := actualPos + len(id)
		if endPos < len(content) && (content[endPos] == '"' || content[endPos] == '`') {
			return true
		}
		idx = actualPos + 1
	}
}

// ExtractAnnotations filters out identifiers that are wrapped in backticks or quotes.
// It extracts IDD references from annotation comments like @implement, @test, and @test-contract.
//
// @implement SPEC-PKG_PATTERN-008
func ExtractAnnotations(content string) []string {
	var refs []string
	for _, pat := range AnnotationPatterns {
		matches := pat.Regex.FindAllStringSubmatch(content, -1)
		for _, m := range matches {
			if len(m) > 1 {
				ids := SplitAnnotationRefs(m[1])
				refs = append(refs, ids...)
			}
		}
	}
	return refs
}

// SplitAnnotationRefs splits comma-separated IDD references and trims whitespace.
// @implement SPEC-PKG_PATTERN-005
func SplitAnnotationRefs(s string) []string {
	var refs []string
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			refs = append(refs, part)
		}
	}
	return refs
}

// GetIdentifierType determines the identifier type from a reference string.
// It returns one of: "SPEC", "TEST", "CONTRACT", "DESIGN", "PATTERN", "WALK",
// or an empty string if the reference does not match any known pattern.
//
// @implement SPEC-PKG_PATTERN-006
func GetIdentifierType(ref string) string {
	for name, pat := range Patterns {
		if pat.Regex.MatchString(ref) {
			return name
		}
	}
	return ""
}

// GetAnnotationType maps annotation prefixes to identifier types.
// It returns the mapped type for known prefixes. @implement maps to SPEC,
// @test maps to TEST. Returns empty string for unknown prefixes.
//
// @implement SPEC-PKG_PATTERN-003
func GetAnnotationType(prefix string) string {
	for _, pat := range AnnotationPatterns {
		if pat.Prefix == prefix {
			return pat.Type
		}
	}
	return ""
}

// ValidateIDPattern validates an identifier against IDD pattern rules.
// It checks if the identifier matches any of the known IDD patterns
// (SPEC, CONTRACT, TEST, DESIGN, PATTERN, WALK) and returns an error if invalid.
//
// @implement SPEC-PKG_PATTERN-007
func ValidateIDPattern(id string) error {
	for _, pat := range Patterns {
		if pat.Regex.MatchString(id) {
			return nil
		}
	}
	return fmt.Errorf("invalid IDD identifier: %s", id)
}

// ValidateIdentifierFormat checks that an identifier follows the strict format TYPE-MODULE-NUMBER
// or is a section marker (TYPE-NUMBER).
// Returns an error if:
// - Identifier does not have 2 or 3 parts separated by hyphens
// - TYPE is not one of SPEC, CONTRACT, TEST, DESIGN, PATTERN, WALK
// - MODULE contains non-uppercase letters or non-alphanumeric characters
// - NUMBER is not purely digits
// @implement SPEC-PKG_PATTERN-009
func ValidateIdentifierFormat(id string) error {
	parts := strings.Split(id, "-")

	if len(parts) == 2 {
		return validateSectionMarker(id, parts)
	} else if len(parts) == 3 {
		return validateIDDIdentifier(id, parts)
	}

	return fmt.Errorf("identifier '%s' must have 2 or 3 parts, found %d parts", id, len(parts))
}

func validateSectionMarker(id string, parts []string) error {
	typePart := strings.ToUpper(parts[0])
	if typePart != "PATTERN" && typePart != "WALK" {
		return fmt.Errorf("identifier '%s' has invalid TYPE '%s', expected PATTERN or WALK for section markers", id, parts[0])
	}

	numberPart := parts[1]
	for _, c := range numberPart {
		if c < '0' || c > '9' {
			return fmt.Errorf("identifier '%s' has invalid NUMBER '%s', must be digits only", id, numberPart)
		}
	}
	return fmt.Errorf("identifier '%s' uses internal section marker TYPE '%s', not a valid IDD identifier (PATTERN-*, WALK-* are internal section markers)", id, typePart)
}

func validateIDDIdentifier(id string, parts []string) error {
	typePart := strings.ToUpper(parts[0])
	if typePart != "SPEC" && typePart != "CONTRACT" && typePart != "TEST" && typePart != "DESIGN" {
		return fmt.Errorf("identifier '%s' has invalid TYPE '%s', expected SPEC, CONTRACT, TEST, or DESIGN", id, parts[0])
	}

	modulePart := parts[1]
	for _, c := range modulePart {
		if (c < 'A' || c > 'Z') && (c < '0' || c > '9') && c != '_' {
			return fmt.Errorf("identifier '%s' has invalid MODULE '%s', must be uppercase letters, digits, or underscore only", id, modulePart)
		}
	}

	numberPart := parts[2]
	for _, c := range numberPart {
		if c < '0' || c > '9' {
			return fmt.Errorf("identifier '%s' has invalid NUMBER '%s', must be digits only", id, numberPart)
		}
	}
	return nil
}
