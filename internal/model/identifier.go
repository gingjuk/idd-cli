// Package model defines the core data structures for IDD link validation.

// Spec: docs/internal/model/spec.md
// Contract: docs/internal/model/contract.md
package model

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// IdentifierType represents the type of IDD identifier.
// @implement SPEC-INTERNAL_MODEL-001
type IdentifierType string

const (
	TypeSpec     IdentifierType = "SPEC"
	TypeContract IdentifierType = "CONTRACT"
	TypeTest     IdentifierType = "TEST"
	TypeDesign   IdentifierType = "DESIGN"
)

// ParseIdentifierType parses identifier type from string representation.
// @implement SPEC-INTERNAL_MODEL-001
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

// Origin indicates where an identifier was found (doc or code).
// @implement SPEC-INTERNAL_MODEL-014
type Origin string

const (
	OriginDoc  Origin = "doc"
	OriginCode Origin = "code"
)

// Identifier represents a single IDD identifier found in docs or code.
// @implement SPEC-INTERNAL_MODEL-018
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
	// Kind distinguishes behavior and contract TEST definitions/annotations.
	Kind string
}

// NewIdentifier creates a new identifier with given fields.
// @implement SPEC-INTERNAL_MODEL-002
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
// @implement SPEC-INTERNAL_MODEL-002
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
// @implement SPEC-INTERNAL_MODEL-003
func (i *Identifier) AddLink(ref string) {
	i.Links = append(i.Links, ref)
}

// SetOrigin sets the origin of this identifier.
// @implement SPEC-INTERNAL_MODEL-017
func (i *Identifier) SetOrigin(origin Origin) {
	i.Origin = origin
}

// IdentifierSet stores identifiers by type and ID.
// @implement SPEC-INTERNAL_MODEL-004
type IdentifierSet struct {
	Specs     []*Identifier
	Contracts []*Identifier
	Tests     []*Identifier
	Designs   []*Identifier

	byID map[string][]*Identifier
}

// NewIdentifierSet creates a new empty identifier set.
// @implement SPEC-INTERNAL_MODEL-004
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
// @implement SPEC-INTERNAL_MODEL-004
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

// Get returns the first identifier by ID.
// @implement SPEC-INTERNAL_MODEL-005
func (s *IdentifierSet) Get(id string) (*Identifier, bool) {
	ids, ok := s.byID[id]
	if !ok || len(ids) == 0 {
		return nil, false
	}
	return ids[0], true
}

// GetAll returns all identifiers with the given ID.
// @implement SPEC-INTERNAL_MODEL-005
func (s *IdentifierSet) GetAll(id string) []*Identifier {
	return s.byID[id]
}

// Has returns true if the identifier exists in the set.
// @implement SPEC-INTERNAL_MODEL-005
func (s *IdentifierSet) Has(id string) bool {
	ids, ok := s.byID[id]
	return ok && len(ids) > 0
}

// All returns all unique identifiers as a slice (one per ID), sorted by ID.
// @implement SPEC-INTERNAL_MODEL-006
func (s *IdentifierSet) All() []*Identifier {
	result := make([]*Identifier, 0, len(s.byID))
	for _, ids := range s.byID {
		if len(ids) > 0 {
			result = append(result, ids[0])
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}

// AllIdentifiers returns all identifiers including duplicates (multiple origins), sorted by ID.
// @implement SPEC-INTERNAL_MODEL-006
func (s *IdentifierSet) AllIdentifiers() []*Identifier {
	result := make([]*Identifier, 0)
	for _, ids := range s.byID {
		result = append(result, ids...)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}

// ByOrigin returns all identifiers with the specified origin, sorted by ID.
// @implement SPEC-INTERNAL_MODEL-004
func (s *IdentifierSet) ByOrigin(origin Origin) []*Identifier {
	result := make([]*Identifier, 0)
	for _, ids := range s.byID {
		for _, id := range ids {
			if id.Origin == origin {
				result = append(result, id)
			}
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}

// HasOrigin returns true if at least one identifier with the given ID has the specified origin.
// @implement SPEC-INTERNAL_MODEL-004
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
// @implement SPEC-INTERNAL_MODEL-006
func (s *IdentifierSet) Count() int {
	return len(s.byID)
}

// DuplicateDocGroups returns groups of doc-origin identifiers that share the same ID
// across multiple source directories (packages), indicating a naming conflict.
// Multiple files in the same directory referencing the same ID are not a conflict.
func (s *IdentifierSet) DuplicateDocGroups() [][]*Identifier {
	return s.duplicatesAcrossDirectories(OriginDoc)
}

// DuplicateCodeGroups returns groups of code-origin identifiers that share the same ID
// across multiple source packages, indicating a naming conflict.
func (s *IdentifierSet) DuplicateCodeGroups() [][]*Identifier {
	return s.duplicatesAcrossDirectories(OriginCode)
}

func (s *IdentifierSet) duplicatesAcrossDirectories(origin Origin) [][]*Identifier {
	// byDir maps idString -> directory -> first Identifier from that directory
	byDir := make(map[string]map[string]*Identifier)
	for _, ids := range s.byID {
		for _, id := range ids {
			if id.Origin != origin {
				continue
			}
			dir := filepath.Dir(id.Source)
			if byDir[id.ID] == nil {
				byDir[id.ID] = make(map[string]*Identifier)
			}
			if _, exists := byDir[id.ID][dir]; !exists {
				byDir[id.ID][dir] = id
			}
		}
	}

	var groups [][]*Identifier
	for _, dirMap := range byDir {
		if len(dirMap) <= 1 {
			continue
		}
		group := make([]*Identifier, 0, len(dirMap))
		for _, id := range dirMap {
			group = append(group, id)
		}
		sort.Slice(group, func(i, j int) bool { return group[i].Source < group[j].Source })
		groups = append(groups, group)
	}
	return groups
}

// Merge combines another identifier set into this one.
// @implement SPEC-INTERNAL_MODEL-007
func (s *IdentifierSet) Merge(other *IdentifierSet) {
	for _, id := range other.AllIdentifiers() {
		s.Add(id)
	}
}

// Annotation represents an annotation found in source code.
// @implement SPEC-INTERNAL_MODEL-008
type Annotation struct {
	Type            IdentifierType
	Ref             string
	Source          string
	Line            int
	Raw             string
	Context         string
	FunctionComment string
}

// NewAnnotation creates a new annotation from code.
// @implement SPEC-INTERNAL_MODEL-008
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
// @implement SPEC-INTERNAL_MODEL-008
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
// @implement SPEC-INTERNAL_MODEL-009
func (a *Annotation) ToIdentifier() *Identifier {
	id := NewIdentifier(a.Ref, a.Type, "", a.Source, a.Line)
	id.RawRef = a.Raw
	if a.FunctionComment != "" {
		id.Describe = a.FunctionComment
	}
	return id
}

// ValidationError represents a single validation error.
// @implement SPEC-INTERNAL_MODEL-010
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
// @implement SPEC-INTERNAL_MODEL-010
type ValidationStats struct {
	TotalIdentifiers  int `json:"total_identifiers"`
	TotalLinks        int `json:"total_links"`
	SpecsAnalyzed     int `json:"specs_analyzed"`
	TestsAnalyzed     int `json:"tests_analyzed"`
	ContractsAnalyzed int `json:"contracts_analyzed"`
	DesignsAnalyzed   int `json:"designs_analyzed"`
}

// ValidationResult contains the result of validation.
// @implement SPEC-INTERNAL_MODEL-010
type ValidationResult struct {
	Valid    bool              `json:"valid"`
	Errors   []ValidationError `json:"errors,omitempty"`
	Warnings []ValidationError `json:"warnings,omitempty"`
	Stats    ValidationStats   `json:"stats"`
	Graph    *GraphSnapshot    `json:"graph,omitempty"`
}

// NewValidationResult creates a new validation result.
// @implement SPEC-INTERNAL_MODEL-010
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		Errors:   make([]ValidationError, 0),
		Warnings: make([]ValidationError, 0),
	}
}

// AddError adds a validation error.
// @implement SPEC-INTERNAL_MODEL-011
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
// @implement SPEC-INTERNAL_MODEL-011
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
// @implement SPEC-INTERNAL_MODEL-011
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
// @implement SPEC-INTERNAL_MODEL-010
type Report struct {
	Tool      string           `json:"tool"`
	Version   string           `json:"version"`
	Timestamp string           `json:"timestamp"`
	Config    ConfigSummary    `json:"config"`
	Result    ValidationResult `json:"result"`
}

// ConfigSummary summarizes the config used for the run.
// @implement SPEC-INTERNAL_MODEL-010
type ConfigSummary struct {
	DocPatterns  []string          `json:"doc_patterns"`
	CodePatterns []string          `json:"code_patterns"`
	Annotations  map[string]string `json:"annotations"`
}

// GraphSnapshot is a summary of the linkage graph for the report.
// @implement SPEC-INTERNAL_MODEL-010
type GraphSnapshot struct {
	Nodes []NodeSummary `json:"nodes"`
	Edges []EdgeSummary `json:"edges"`
}

// NodeSummary is a summary of a node for the report.
// @implement SPEC-INTERNAL_MODEL-010
type NodeSummary struct {
	ID       string         `json:"id"`
	Type     IdentifierType `json:"type"`
	Outbound int            `json:"outbound"`
	Inbound  int            `json:"inbound"`
}

// EdgeSummary is a summary of an edge for the report.
// @implement SPEC-INTERNAL_MODEL-010
type EdgeSummary struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Type     string `json:"type"`
	Verified bool   `json:"verified"`
	Source   string `json:"source,omitempty"`
	Line     int    `json:"line,omitempty"`
}

// LLMReport is a finding-centered report shape for LLM analysis and repair.
// @implement SPEC-INTERNAL_MODEL-037
type LLMReport struct {
	Schema   string       `json:"schema"`
	Status   string       `json:"status"`
	Summary  LLMSummary   `json:"summary"`
	Findings []LLMFinding `json:"findings"`
}

// LLMSummary summarizes validation failures for the LLM report.
// @implement SPEC-INTERNAL_MODEL-037
type LLMSummary struct {
	Errors     int               `json:"errors"`
	Warnings   int               `json:"warnings"`
	TopRules   []string          `json:"top_rules,omitempty"`
	RuleGroups []LLMFindingGroup `json:"rule_groups,omitempty"`
}

// LLMFindingGroup summarizes repeated findings that share the same rule and severity.
// @implement SPEC-INTERNAL_MODEL-037
type LLMFindingGroup struct {
	Severity       string   `json:"severity"`
	Rule           string   `json:"rule"`
	Title          string   `json:"title"`
	Count          int      `json:"count"`
	Files          []string `json:"files,omitempty"`
	Identifiers    []string `json:"identifiers,omitempty"`
	FindingIndexes []int    `json:"finding_indexes"`
	SuggestedFix   string   `json:"suggested_fix"`
}

// LLMFinding is a self-contained validation issue for LLM consumption.
// @implement SPEC-INTERNAL_MODEL-037
type LLMFinding struct {
	Severity           string                 `json:"severity"`
	Rule               string                 `json:"rule"`
	Title              string                 `json:"title"`
	Location           LLMLocation            `json:"location,omitempty"`
	Identifier         string                 `json:"identifier,omitempty"`
	Problem            string                 `json:"problem"`
	Expected           string                 `json:"expected,omitempty"`
	Actual             string                 `json:"actual,omitempty"`
	SuggestedFix       string                 `json:"suggested_fix"`
	RelatedIdentifiers []LLMRelatedIdentifier `json:"related_identifiers,omitempty"`
}

// LLMLocation identifies the file and line associated with a finding.
// @implement SPEC-INTERNAL_MODEL-037
type LLMLocation struct {
	File string `json:"file,omitempty"`
	Line int    `json:"line,omitempty"`
}

// LLMRelatedIdentifier describes an identifier related to a finding.
// @implement SPEC-INTERNAL_MODEL-037
type LLMRelatedIdentifier struct {
	ID       string `json:"id"`
	Relation string `json:"relation"`
}
