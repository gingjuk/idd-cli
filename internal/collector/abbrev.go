// Package collector provides path abbreviation mappings.

// Spec: docs/internal/collector/spec.md
// Contract: docs/internal/collector/contract.md
package collector

import "strings"

// abbrevOrUpper converts a path component to its module prefix form:
// uppercase the component and replace hyphens with underscores.
func abbrevOrUpper(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, "-", "_"))
}
