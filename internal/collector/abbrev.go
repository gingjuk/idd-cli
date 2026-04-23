// Package collector provides path abbreviation mappings.

// Spec: docs/internal/collector/spec.md
// Contract: docs/internal/collector/contract.md
package collector

import "strings"

var pathAbbrevs = map[string]string{
	"internal":   "INT",
	"config":     "CFG",
	"contract":   "CON",
	"similarity": "SIM",
	"pattern":    "PAT",
	"reporter":   "RPT",
	"identifier": "ID",
	"collector":  "COL",
	"engine":     "ENG",
	"graph":      "GRPH",
	"model":      "MOD",
	"cmd":        "CMD",
	"idd-cli":    "IDD",
	"backend":    "BACKEND",
}

func abbrevOrUpper(s string) string {
	if v, ok := pathAbbrevs[s]; ok {
		return v
	}
	return strings.ToUpper(s)
}
