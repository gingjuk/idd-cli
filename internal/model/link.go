package model

// LinkType represents the type of relationship between identifiers.
type LinkType string

const (
	LinkImplements LinkType = "implements"
	LinkTests     LinkType = "tests"
	LinkReferences LinkType = "references"
	LinkAnnotates  LinkType = "annotates"
)

// ReverseLinkType returns the reverse link type.
func ReverseLinkType(lt LinkType) LinkType {
	switch lt {
	case LinkTests:
		return LinkImplements
	case LinkImplements:
		return LinkTests
	case LinkReferences:
		return LinkReferences
	case LinkAnnotates:
		return LinkAnnotates
	default:
		return lt
	}
}

// Link represents a directed relationship between two identifiers.
type Link struct {
	From   string
	To     string
	Type   LinkType
	Source string
	Line   int
}

// NewLink creates a new link.
func NewLink(from, to string, linkType LinkType, source string, line int) *Link {
	return &Link{
		From:   from,
		To:     to,
		Type:   linkType,
		Source: source,
		Line:   line,
	}
}