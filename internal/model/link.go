// Package model provides link types for IDD identifier relationships.

// Spec: docs/internal/model/spec.md
// Contract: docs/internal/model/contract.md
package model

// LinkType represents the type of relationship between identifiers.
// @implement SPEC-INTERNAL_MODEL-012
type LinkType string

const (
	LinkImplements         LinkType = "implements"
	LinkTests              LinkType = "tests"
	LinkContractTests      LinkType = "contract_tests"
	LinkReferences         LinkType = "references"
	LinkAnnotates          LinkType = "annotates"
	LinkContract           LinkType = "contract"            // SPEC → contract (forward direction)
	LinkContractImplements LinkType = "contract_implements" // contract → SPEC (backward direction)
	LinkDependsOn          LinkType = "depends_on"
	LinkDependedBy         LinkType = "depended_by"
	LinkSupersedes         LinkType = "supersedes"
	LinkDeprecatedBy       LinkType = "deprecated_by"
)

// ReverseLinkType returns the reverse link type.
// @implement SPEC-INTERNAL_MODEL-012
func ReverseLinkType(lt LinkType) LinkType {
	switch lt {
	case LinkTests:
		return LinkImplements
	case LinkImplements:
		return LinkTests
	case LinkContractTests:
		return LinkContractTests
	case LinkReferences:
		return LinkReferences
	case LinkAnnotates:
		return LinkAnnotates
	case LinkContract:
		return LinkContractImplements
	case LinkContractImplements:
		return LinkContract
	case LinkDependsOn:
		return LinkDependedBy
	case LinkDependedBy:
		return LinkDependsOn
	case LinkSupersedes:
		return LinkDeprecatedBy
	case LinkDeprecatedBy:
		return LinkSupersedes
	default:
		return lt
	}
}

// IdentifierLink is an explicitly typed outgoing relationship owned by an
// Identifier. It is used when source and target identifier types cannot
// disambiguate relationship semantics.
// @implement SPEC-INTERNAL_MODEL-012
type IdentifierLink struct {
	Ref  string
	Type LinkType
}

// Link represents a directed relationship between two identifiers.
// @implement SPEC-INTERNAL_MODEL-013
type Link struct {
	From   string
	To     string
	Type   LinkType
	Source string
	Line   int
}

// NewLink creates a new link.
// @implement SPEC-INTERNAL_MODEL-015
func NewLink(from, to string, linkType LinkType, source string, line int) *Link {
	return &Link{
		From:   from,
		To:     to,
		Type:   linkType,
		Source: source,
		Line:   line,
	}
}
