package pattern

import (
	"fmt"
	"regexp"
	"strings"
)

type IDDPattern struct {
	Type    string
	Regex   *regexp.Regexp
	Example string
}

var Patterns = map[string]*IDDPattern{
	"SPEC": {
		Type:    "SPEC",
		Regex:   regexp.MustCompile(`(?i)\b(SPEC-[A-Z]+-[0-9]+)\b`),
		Example: "SPEC-BE-001",
	},
	"CONTRACT": {
		Type:    "CONTRACT",
		Regex:   regexp.MustCompile(`(?i)\b(CONTRACT-[A-Z]+-[0-9]+)\b`),
		Example: "CONTRACT-BE-001",
	},
	"TEST": {
		Type:    "TEST",
		Regex:   regexp.MustCompile(`(?i)\b(TEST-[A-Z]+-[0-9]+)\b`),
		Example: "TEST-BE-001",
	},
	"DESIGN": {
		Type:    "DESIGN",
		Regex:   regexp.MustCompile(`(?i)\b(DESIGN-[A-Z]+-[0-9]+)\b`),
		Example: "DESIGN-BE-001",
	},
}

type AnnotationPattern struct {
	Prefix string
	Regex  *regexp.Regexp
	Type   string
}

var AnnotationPatterns = []AnnotationPattern{
	{
		Prefix: "@spec",
		Regex:  regexp.MustCompile(`(?i)@spec\s+([A-Z]+-[A-Z]+-[0-9]+(?:\s*,\s*[A-Z]+-[A-Z]+-[0-9]+)*)`),
		Type:   "SPEC",
	},
	{
		Prefix: "@contract",
		Regex:  regexp.MustCompile(`(?i)@contract\s+([A-Z]+-[A-Z]+-[0-9]+(?:\s*,\s*[A-Z]+-[A-Z]+-[0-9]+)*)`),
		Type:   "CONTRACT",
	},
	{
		Prefix: "@test",
		Regex:  regexp.MustCompile(`(?i)@test\s+([A-Z]+-[A-Z]+-[0-9]+(?:\s*,\s*[A-Z]+-[A-Z]+-[0-9]+)*)`),
		Type:   "TEST",
	},
	{
		Prefix: "@design",
		Regex:  regexp.MustCompile(`(?i)@design\s+([A-Z]+-[A-Z]+-[0-9]+(?:\s*,\s*[A-Z]+-[A-Z]+-[0-9]+)*)`),
		Type:   "DESIGN",
	},
}

func ExtractIDDReferences(content string) []string {
	var refs []string
	for _, pat := range Patterns {
		matches := pat.Regex.FindAllStringSubmatch(content, -1)
		for _, m := range matches {
			if len(m) > 1 {
				id := strings.ToUpper(m[1])
				if !isQuoted(content, m[1]) {
					refs = append(refs, id)
				}
			}
		}
	}
	return refs
}

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

func GetIdentifierType(ref string) string {
	for name, pat := range Patterns {
		if pat.Regex.MatchString(ref) {
			return name
		}
	}
	return ""
}

func GetAnnotationType(prefix string) string {
	for _, pat := range AnnotationPatterns {
		if pat.Prefix == prefix {
			return pat.Type
		}
	}
	return ""
}

func ValidateIDPattern(id string) error {
	for _, pat := range Patterns {
		if pat.Regex.MatchString(id) {
			return nil
		}
	}
	return fmt.Errorf("invalid IDD identifier: %s", id)
}
