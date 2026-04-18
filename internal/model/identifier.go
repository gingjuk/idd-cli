// Package model defines the core data structures for IDD link validation.
package model

import (
	"fmt"
	"sort"
	"strings"
)

// IdentifierType represents the type of IDD identifier.
type IdentifierType string

const (
	TypeSpec     IdentifierType = "SPEC"
	TypeContract IdentifierType = "CONTRACT"
	TypeTest     IdentifierType = "TEST"
	TypeDesign   IdentifierType = "DESIGN"
)

// ParseIdentifierType converts a string to IdentifierType.
func ParseIdentifierType(s string) (IdentifierType, error) {
	switch strings.ToUpper(s) {
	case "SPEC":
		return TypeSpec, nil
	case "CONTRACT":
		return TypeContract, nil
	case "TEST":
		return TypeTest, nil
	case "DESIGN":
		return TypeDesign, nil
	default:
		return "", fmt.Errorf("unknown identifier type: %s", s)
	}
}

// Identifier represents a single IDD identifier found in docs or code.
type Identifier struct {
	// ID is the unique identifier (e.g., "SPEC-001")
	ID string
	// Type is the kind of identifier
	Type IdentifierType
	// Title is the human-readable title (extracted from document)
	Title string
	// Source is the file path or "inline" for code annotations
	Source string
	// Line is the line number where the identifier was found
	Line int
	// RawRef is the raw reference text as it appears in source
	RawRef string
	// Links are forward references from this identifier
	Links []string
}

// NewIdentifier creates a new identifier with the given fields.
func NewIdentifier(id string, idType IdentifierType, title, source string, line int) *Identifier {
	return &Identifier{
		ID:     id,
		Type:   idType,
		Title:  title,
		Source: source,
		Line:   line,
		RawRef: id,
		Links:  make([]string, 0),
	}
}

// AddLink adds a forward reference from this identifier.
func (i *Identifier) AddLink(ref string) {
	i.Links = append(i.Links, ref)
}

// IdentifierSet is a collection of all collected identifiers.
type IdentifierSet struct {
	Specs     []*Identifier
	Contracts []*Identifier
	Tests     []*Identifier
	Designs   []*Identifier

	byID map[string]*Identifier
}

// NewIdentifierSet creates a new empty identifier set.
func NewIdentifierSet() *IdentifierSet {
	return &IdentifierSet{
		Specs:     make([]*Identifier, 0),
		Contracts: make([]*Identifier, 0),
		Tests:     make([]*Identifier, 0),
		Designs:   make([]*Identifier, 0),
		byID:      make(map[string]*Identifier),
	}
}

// Add adds an identifier to the set.
func (s *IdentifierSet) Add(id *Identifier) {
	s.byID[id.ID] = id
	switch id.Type {
	case TypeSpec:
		s.Specs = append(s.Specs, id)
	case TypeContract:
		s.Contracts = append(s.Contracts, id)
	case TypeTest:
		s.Tests = append(s.Tests, id)
	case TypeDesign:
		s.Designs = append(s.Designs, id)
	}
}

// Get returns an identifier by ID.
func (s *IdentifierSet) Get(id string) (*Identifier, bool) {
	ident, ok := s.byID[id]
	return ident, ok
}

// Has returns true if the identifier exists in the set.
func (s *IdentifierSet) Has(id string) bool {
	_, ok := s.byID[id]
	return ok
}

// All returns all identifiers as a slice.
func (s *IdentifierSet) All() []*Identifier {
	result := make([]*Identifier, 0, len(s.byID))
	for _, id := range s.byID {
		result = append(result, id)
	}
	return result
}

// Count returns the total number of identifiers.
func (s *IdentifierSet) Count() int {
	return len(s.byID)
}

// Merge combines another identifier set into this one.
func (s *IdentifierSet) Merge(other *IdentifierSet) {
	for _, id := range other.All() {
		if !s.Has(id.ID) {
			s.Add(id)
		}
	}
}

// Annotation represents an annotation found in source code.
type Annotation struct {
	// Type is the annotation type (spec, contract, test, design)
	Type IdentifierType
	// Ref is the identifier reference (e.g., "SPEC-001")
	Ref string
	// Source is the file path
	Source string
	// Line is the line number
	Line int
	// Raw is the raw annotation text
	Raw string
	// Context is surrounding code context
	Context string
}

// NewAnnotation creates a new annotation.
func NewAnnotation(typ IdentifierType, ref, source, raw, context string, line int) *Annotation {
	return &Annotation{
		Type:    typ,
		Ref:     ref,
		Source:  source,
		Line:    line,
		Raw:     raw,
		Context: context,
	}
}

// ToIdentifier converts an annotation to an identifier.
func (a *Annotation) ToIdentifier() *Identifier {
	id := NewIdentifier(a.Ref, a.Type, "", a.Source, a.Line)
	id.RawRef = a.Raw
	return id
}

// ValidationError represents a single validation error.
type ValidationError struct {
	// Rule is the name of the validation rule that failed
	Rule string `json:"rule"`
	// Message is a human-readable error message
	Message string `json:"message"`
	// Source is the file path where the error was found
	Source string `json:"source,omitempty"`
	// Link is the identifier/link involved in the error
	Link string `json:"link,omitempty"`
	// Code is the specific code or line involved
	Code string `json:"code,omitempty"`
}

// Error implements error interface.
func (e ValidationError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Rule, e.Message)
}

// ValidationStats contains statistics about the validation run.
type ValidationStats struct {
	TotalIdentifiers  int `json:"total_identifiers"`
	TotalLinks        int `json:"total_links"`
	SpecsAnalyzed     int `json:"specs_analyzed"`
	TestsAnalyzed     int `json:"tests_analyzed"`
	ContractsAnalyzed int `json:"contracts_analyzed"`
	DesignsAnalyzed   int `json:"designs_analyzed"`
}

// ValidationResult contains the result of validation.
type ValidationResult struct {
	Valid    bool              `json:"valid"`
	Errors   []ValidationError `json:"errors,omitempty"`
	Warnings []ValidationError `json:"warnings,omitempty"`
	Stats    ValidationStats   `json:"stats"`
	Graph    *GraphSnapshot    `json:"graph,omitempty"`
}

// NewValidationResult creates a new validation result.
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		Errors:   make([]ValidationError, 0),
		Warnings: make([]ValidationError, 0),
	}
}

// AddError adds a validation error.
func (r *ValidationResult) AddError(rule, msg, source, link, code string) {
	r.Errors = append(r.Errors, ValidationError{
		Rule:    rule,
		Message: msg,
		Source:  source,
		Link:    link,
		Code:    code,
	})
	r.Valid = false
}

// AddWarning adds a validation warning.
func (r *ValidationResult) AddWarning(rule, msg, source, link, code string) {
	r.Warnings = append(r.Warnings, ValidationError{
		Rule:    rule,
		Message: msg,
		Source:  source,
		Link:    link,
		Code:    code,
	})
}

// Sort sorts errors and warnings by rule and message.
func (r *ValidationResult) Sort() {
	sort.Slice(r.Errors, func(i, j int) bool {
		if r.Errors[i].Rule != r.Errors[j].Rule {
			return r.Errors[i].Rule < r.Errors[j].Rule
		}
		return r.Errors[i].Message < r.Errors[j].Message
	})
	sort.Slice(r.Warnings, func(i, j int) bool {
		if r.Warnings[i].Rule != r.Warnings[j].Rule {
			return r.Warnings[i].Rule < r.Warnings[j].Rule
		}
		return r.Warnings[i].Message < r.Warnings[j].Message
	})
}

// Report is the final output report.
type Report struct {
	Tool      string           `json:"tool"`
	Version   string           `json:"version"`
	Timestamp string           `json:"timestamp"`
	Config    ConfigSummary    `json:"config"`
	Result    ValidationResult `json:"result"`
}

// ConfigSummary summarizes the config used for the run.
type ConfigSummary struct {
	DocPatterns  []string `json:"doc_patterns"`
	CodePatterns []string `json:"code_patterns"`
	Annotations  []string `json:"annotations"`
}

// GraphSnapshot is a summary of the linkage graph for the report.
type GraphSnapshot struct {
	Nodes []NodeSummary `json:"nodes"`
	Edges []EdgeSummary `json:"edges"`
}

// NodeSummary is a summary of a node for the report.
type NodeSummary struct {
	ID       string         `json:"id"`
	Type     IdentifierType `json:"type"`
	Outbound int            `json:"outbound"`
	Inbound  int            `json:"inbound"`
}

// EdgeSummary is a summary of an edge for the report.
type EdgeSummary struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Type     string `json:"type"`
	Verified bool   `json:"verified"`
	Source   string `json:"source,omitempty"`
	Line     int    `json:"line,omitempty"`
}
