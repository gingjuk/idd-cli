package pattern

import (
	"fmt"
	"regexp"
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
	Prefix  string
	Regex   *regexp.Regexp
	Type    string
}

var AnnotationPatterns = []AnnotationPattern{
	{
		Prefix: "@spec",
		Regex:  regexp.MustCompile(`(?i)@spec\s+(\S+)`),
		Type:   "SPEC",
	},
	{
		Prefix: "@contract",
		Regex:  regexp.MustCompile(`(?i)@contract\s+(\S+)`),
		Type:   "CONTRACT",
	},
	{
		Prefix: "@test",
		Regex:  regexp.MustCompile(`(?i)@test\s+(\S+)`),
		Type:   "TEST",
	},
	{
		Prefix: "@design",
		Regex:  regexp.MustCompile(`(?i)@design\s+(\S+)`),
		Type:   "DESIGN",
	},
}

func ExtractIDDReferences(content string) []string {
	var refs []string
	for _, pat := range Patterns {
		matches := pat.Regex.FindAllStringSubmatch(content, -1)
		for _, m := range matches {
			if len(m) > 1 {
				refs = append(refs, m[1])
			}
		}
	}
	return refs
}

func ExtractAnnotations(content string) []string {
	var refs []string
	for _, pat := range AnnotationPatterns {
		matches := pat.Regex.FindAllStringSubmatch(content, -1)
		for _, m := range matches {
			if len(m) > 1 {
				refs = append(refs, m[1])
			}
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