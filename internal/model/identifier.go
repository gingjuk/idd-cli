// Package model defines the core data structures for IDD link validation.
// @spec SPEC-BE-005
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

// Origin indicates where an identifier was found.
type Origin string

const (
	OriginDoc  Origin = "doc"
	OriginCode Origin = "code"
)

// Identifier represents a single IDD identifier found in docs or code.
type Identifier struct {
	ID       string
	Type     IdentifierType
	Title    string
	Describe string
	Source   string
	Line     int
	RawRef   string
	Links    []string
	Origin   Origin
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
		Origin: OriginDoc,
	}
}

// NewIdentifierWithDescribe creates a new identifier with describe field.
func NewIdentifierWithDescribe(id string, idType IdentifierType, title, describe, source string, line int) *Identifier {
	return &Identifier{
		ID:       id,
		Type:     idType,
		Title:    title,
		Describe: describe,
		Source:   source,
		Line:     line,
		RawRef:   id,
		Links:    make([]string, 0),
		Origin:   OriginDoc,
	}
}

// AddLink adds a forward reference from this identifier.
func (i *Identifier) AddLink(ref string) {
	i.Links = append(i.Links, ref)
}

// SetOrigin sets the origin of this identifier.
func (i *Identifier) SetOrigin(origin Origin) {
	i.Origin = origin
}

// IdentifierSet is a collection of all collected identifiers.
type IdentifierSet struct {
	Specs     []*Identifier
	Contracts []*Identifier
	Tests     []*Identifier
	Designs   []*Identifier

	byID map[string][]*Identifier
}

// NewIdentifierSet creates a new empty identifier set.
func NewIdentifierSet() *IdentifierSet {
	return &IdentifierSet{
		Specs:     make([]*Identifier, 0),
		Contracts: make([]*Identifier, 0),
		Tests:     make([]*Identifier, 0),
		Designs:   make([]*Identifier, 0),
		byID:      make(map[string][]*Identifier),
	}
}

// Add adds an identifier to the set.
func (s *IdentifierSet) Add(id *Identifier) {
	s.byID[id.ID] = append(s.byID[id.ID], id)
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

// Get returns the first identifier by ID (for backward compatibility).
func (s *IdentifierSet) Get(id string) (*Identifier, bool) {
	ids, ok := s.byID[id]
	if !ok || len(ids) == 0 {
		return nil, false
	}
	return ids[0], true
}

// GetAll returns all identifiers with the given ID.
func (s *IdentifierSet) GetAll(id string) []*Identifier {
	return s.byID[id]
}

// Has returns true if the identifier exists in the set.
func (s *IdentifierSet) Has(id string) bool {
	ids, ok := s.byID[id]
	return ok && len(ids) > 0
}

// All returns all unique identifiers as a slice (one per ID).
func (s *IdentifierSet) All() []*Identifier {
	result := make([]*Identifier, 0, len(s.byID))
	for _, ids := range s.byID {
		if len(ids) > 0 {
			result = append(result, ids[0])
		}
	}
	return result
}

// AllIdentifiers returns all identifiers including duplicates (multiple origins).
func (s *IdentifierSet) AllIdentifiers() []*Identifier {
	result := make([]*Identifier, 0)
	for _, ids := range s.byID {
		result = append(result, ids...)
	}
	return result
}

// ByOrigin returns all identifiers with the specified origin.
func (s *IdentifierSet) ByOrigin(origin Origin) []*Identifier {
	var result []*Identifier
	for _, ids := range s.byID {
		for _, id := range ids {
			if id.Origin == origin {
				result = append(result, id)
			}
		}
	}
	return result
}

// HasOrigin returns true if at least one identifier with the given ID has the specified origin.
func (s *IdentifierSet) HasOrigin(id string, origin Origin) bool {
	ids, ok := s.byID[id]
	if !ok {
		return false
	}
	for _, id := range ids {
		if id.Origin == origin {
			return true
		}
	}
	return false
}

// Count returns the total number of identifiers.
func (s *IdentifierSet) Count() int {
	return len(s.byID)
}

// Merge combines another identifier set into this one.
func (s *IdentifierSet) Merge(other *IdentifierSet) {
	for _, id := range other.AllIdentifiers() {
		s.Add(id)
	}
}

// Annotation represents an annotation found in source code.
type Annotation struct {
	Type            IdentifierType
	Ref             string
	Source          string
	Line            int
	Raw             string
	Context         string
	FunctionComment string
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

// NewAnnotationWithComment creates a new annotation with function comment.
func NewAnnotationWithComment(typ IdentifierType, ref, source, raw, context, funcComment string, line int) *Annotation {
	return &Annotation{
		Type:            typ,
		Ref:             ref,
		Source:          source,
		Line:            line,
		Raw:             raw,
		Context:         context,
		FunctionComment: funcComment,
	}
}

// ToIdentifier converts an annotation to an identifier.
func (a *Annotation) ToIdentifier() *Identifier {
	id := NewIdentifier(a.Ref, a.Type, "", a.Source, a.Line)
	id.RawRef = a.Raw
	if a.FunctionComment != "" {
		id.Describe = a.FunctionComment
	}
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
