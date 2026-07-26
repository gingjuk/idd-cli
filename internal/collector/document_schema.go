// Package collector defines the shared schema for generated IDD documents.

// Spec: docs/internal/collector/spec.md
// Contract: docs/internal/collector/contract.md
package collector

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	goldmarkast "github.com/yuin/goldmark/ast"
	goldmarktext "github.com/yuin/goldmark/text"
)

const documentCompletionSchema = "idd.document_status.v1"

var scaffoldMarkerPattern = regexp.MustCompile(
	`^\s*<!--\s*idd:scaffold\s+slot="([a-z0-9.-]+)"\s*-->\s*$`,
)
var completionSlotPartPattern = regexp.MustCompile(`[^a-z0-9]+`)

// IncompleteSlot identifies one generated or structurally required document
// location that still needs human-authored content.
// @implement SPEC-INTERNAL_COLLECTOR-001
type IncompleteSlot struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Role   string `json:"role"`
	Slot   string `json:"slot"`
	Reason string `json:"reason"`
}

// DocumentCompletionStatus is the stable machine-readable work list returned
// by docs init and docs status.
// @implement SPEC-INTERNAL_COLLECTOR-001
type DocumentCompletionStatus struct {
	Schema          string           `json:"schema"`
	Status          string           `json:"status"`
	IncompleteSlots []IncompleteSlot `json:"incomplete_slots"`
}

type iddDocumentSlotSchema struct {
	Name         string
	Heading      string
	RecordPrefix string
	Guidance     []string
}

type iddDocumentRecordFieldSchema struct {
	Name             string
	Label            string
	List             bool
	Heading          bool
	RequiredWhenKind string
}

type iddDocumentRecordSchema struct {
	Section string
	Type    string
	Fields  []iddDocumentRecordFieldSchema
}

type iddDocumentRoleSchema struct {
	Role       string
	Title      string
	Collection iddDocumentSlotSchema
	Sections   []iddDocumentSlotSchema
	Record     iddDocumentRecordSchema
}

var iddDocumentSchemas = map[string]iddDocumentRoleSchema{
	"design": {
		Role:  "design",
		Title: "Design",
		Collection: iddDocumentSlotSchema{
			Name:         "design.components",
			RecordPrefix: "Component:",
			Guidance: []string{
				"Add one `## Component: <name>` section for each named component.",
				"Explain why it exists, its responsibilities and state ownership,",
				"what remains outside its boundary, how it collaborates, and the",
				"decisions and trade-offs behind that shape.",
			},
		},
		Record: iddDocumentRecordSchema{
			Section: "components",
			Type:    "Component",
			Fields: []iddDocumentRecordFieldSchema{
				{Name: "name", Label: "Component name", Heading: true},
				{Name: "purpose", Label: "Purpose"},
			},
		},
		Sections: []iddDocumentSlotSchema{
			{
				Name:    "design.architecture",
				Heading: "Architecture",
				Guidance: []string{
					"Describe data and control flow, state ownership, lifecycle,",
					"concurrency when relevant, and failure containment.",
				},
			},
			{
				Name:    "design.package-layout",
				Heading: "Package Layout",
				Guidance: []string{
					"Explain why responsibilities live in their selected packages,",
					"not only the directory tree.",
				},
			},
			{
				Name:    "design.function-composition",
				Heading: "Function Composition",
				Guidance: []string{
					"Describe meaningful execution paths and initialization or",
					"lifecycle order.",
				},
			},
			{
				Name:    "design.dependencies",
				Heading: "Dependencies",
				Guidance: []string{
					"Explain what each dependency contributes and which assumptions",
					"must remain true.",
				},
			},
			{
				Name:    "design.testability-hooks",
				Heading: "Testability Hooks",
				Guidance: []string{
					"Explain seams, fakes, observable outcomes, and difficult",
					"failure paths.",
				},
			},
		},
	},
	"contract": {
		Role:  "contract",
		Title: "Contracts",
		Collection: iddDocumentSlotSchema{
			Name:         "contract.records",
			RecordPrefix: "Contract:",
			Guidance: []string{
				"Add one `## Contract: <name>` section for each observable boundary.",
				"Describe the caller-visible capability,",
				"inputs and outputs, validity and ownership rules,",
				"errors, side effects, invariants, and compatibility expectations.",
				"A signature may clarify",
				"the boundary but does not replace its behavioral explanation.",
			},
		},
		Record: iddDocumentRecordSchema{
			Section: "contracts",
			Type:    "Contract",
			Fields: []iddDocumentRecordFieldSchema{
				{Name: "name", Label: "Contract name", Heading: true},
				{Name: "guarantees", Label: "Guarantees"},
			},
		},
	},
	"spec": {
		Role:  "spec",
		Title: "Specifications",
		Collection: iddDocumentSlotSchema{
			Name:         "spec.records",
			RecordPrefix: "SPEC-",
			Guidance: []string{
				"Add each cohesive behavior as `## SPEC-<MODULE>-<NUMBER>: <title>`,",
				"then fill the required fields listed below. Continue with required behavior,",
				"edge and failure cases, implementation boundary and non-goals,",
				"rationale, examples, and concrete acceptance evidence.",
				"The record is behavioral documentation, not an API listing.",
			},
		},
		Record: iddDocumentRecordSchema{
			Section: "specs",
			Type:    "SPEC",
			Fields: []iddDocumentRecordFieldSchema{
				{Name: "id", Label: "SPEC identifier", Heading: true},
				{Name: "title", Label: "title", Heading: true},
				{Name: "design", Label: "Design"},
				{Name: "contract", Label: "Contract"},
				{Name: "requirement", Label: "Requirement"},
				{Name: "acceptance", Label: "Acceptance"},
			},
		},
	},
	"testing": {
		Role:  "testing",
		Title: "Testing",
		Collection: iddDocumentSlotSchema{
			Name:         "testing.records",
			RecordPrefix: "TEST-",
			Guidance: []string{
				"Add each evidence record as `## TEST-<MODULE>-<NUMBER>: <title>`,",
				"then fill the required fields listed below. Explain",
				"positive/negative/boundary/failure scenarios,",
				"fixtures and isolation, oracles, nondeterminism controls, and",
				"meaningful exclusions.",
			},
		},
		Record: iddDocumentRecordSchema{
			Section: "tests",
			Type:    "TEST",
			Fields: []iddDocumentRecordFieldSchema{
				{Name: "id", Label: "TEST identifier", Heading: true},
				{Name: "title", Label: "title", Heading: true},
				{Name: "kind", Label: "Kind"},
				{Name: "covers", Label: "Covers", List: true},
				{Name: "contracts", Label: "Contracts", List: true, RequiredWhenKind: "contract"},
				{Name: "purpose", Label: "Purpose"},
				{Name: "oracle", Label: "Oracle"},
			},
		},
	},
}

func renderIDDDocumentTemplate(packagePath, role string) string {
	schema, ok := iddDocumentSchemas[role]
	if !ok {
		return "\n"
	}

	var builder strings.Builder
	fmt.Fprintf(&builder, "\n# %s: %s\n\n", schema.Title, packagePath)
	writeScaffoldSlot(&builder, schema.Collection, schema.Record)
	for _, section := range schema.Sections {
		fmt.Fprintf(&builder, "\n## %s\n\n", section.Heading)
		writeScaffoldSlot(&builder, section, iddDocumentRecordSchema{})
	}
	return builder.String()
}

func writeScaffoldSlot(
	builder *strings.Builder,
	slot iddDocumentSlotSchema,
	record iddDocumentRecordSchema,
) {
	fmt.Fprintf(builder, "<!-- idd:scaffold slot=%q -->\n\n", slot.Name)
	for _, line := range slot.Guidance {
		fmt.Fprintf(builder, "> %s\n", line)
	}
	required := make([]string, 0, len(record.Fields))
	conditional := make([]string, 0)
	for _, field := range record.Fields {
		if field.Heading {
			continue
		}
		if field.RequiredWhenKind != "" {
			conditional = append(
				conditional,
				fmt.Sprintf("%s when Kind is `%s`", field.Label, field.RequiredWhenKind),
			)
			continue
		}
		required = append(required, field.Label)
	}
	if len(required) > 0 {
		fmt.Fprintf(builder, "> Required fields: %s.\n", strings.Join(required, ", "))
	}
	if len(conditional) > 0 {
		fmt.Fprintf(builder, "> Conditionally required: %s.\n", strings.Join(conditional, ", "))
	}
}

// InspectDocumentCompletion inspects one IDD document or a directory containing
// one or more package document sets.
// @implement SPEC-INTERNAL_COLLECTOR-001
func InspectDocumentCompletion(path string) (*DocumentCompletionStatus, error) {
	targets, missing, err := documentCompletionTargets(path)
	if err != nil {
		return nil, err
	}

	slots := append([]IncompleteSlot(nil), missing...)
	for _, target := range targets {
		data, readErr := os.ReadFile(target)
		if readErr != nil {
			return nil, fmt.Errorf("read document completion target %s: %w", target, readErr)
		}
		parsed, detected, parseErr := parseIDDDocument(target, data)
		if parseErr != nil {
			return nil, fmt.Errorf("inspect document completion for %s: %w", target, parseErr)
		}
		if !detected {
			return nil, fmt.Errorf("%s is not a self-describing IDD document", target)
		}
		role := parsed.Metadata.Document
		if pathRole, ok := iddDocumentRoles[filepath.Base(target)]; ok {
			role = pathRole
		}
		slots = append(slots, inspectParsedDocumentCompletion(parsed, role, true)...)
	}
	sortIncompleteSlots(slots)

	status := "complete"
	if len(slots) > 0 {
		status = "incomplete"
	}
	if slots == nil {
		slots = make([]IncompleteSlot, 0)
	}
	return &DocumentCompletionStatus{
		Schema:          documentCompletionSchema,
		Status:          status,
		IncompleteSlots: slots,
	}, nil
}

func documentCompletionTargets(path string) ([]string, []IncompleteSlot, error) {
	cleanPath := filepath.Clean(path)
	info, err := os.Stat(cleanPath)
	if err != nil {
		return nil, nil, fmt.Errorf("inspect document status path: %w", err)
	}
	if !info.IsDir() {
		if _, ok := iddDocumentRoles[filepath.Base(cleanPath)]; !ok {
			return nil, nil, fmt.Errorf("document status path must name an IDD role Markdown file or directory")
		}
		return []string{cleanPath}, nil, nil
	}

	packageDirs := make(map[string]bool)
	err = filepath.WalkDir(cleanPath, func(current string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		if _, ok := iddDocumentRoles[entry.Name()]; !ok {
			return nil
		}
		data, readErr := os.ReadFile(current)
		if readErr != nil {
			return readErr
		}
		if hasIDDDocumentFrontmatter(data) {
			packageDirs[filepath.Dir(current)] = true
		}
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("walk document status path: %w", err)
	}
	if len(packageDirs) == 0 {
		return nil, nil, fmt.Errorf("no IDD role documents found under %s", cleanPath)
	}

	var targets []string
	var missing []IncompleteSlot
	for directory := range packageDirs {
		for _, filename := range iddDocumentOrder {
			target := filepath.Join(directory, filename)
			if _, statErr := os.Stat(target); statErr == nil {
				targets = append(targets, target)
				continue
			} else if !os.IsNotExist(statErr) {
				return nil, nil, fmt.Errorf("inspect document status target %s: %w", target, statErr)
			}
			role := iddDocumentRoles[filename]
			missing = append(missing, IncompleteSlot{
				File:   target,
				Line:   1,
				Role:   role,
				Slot:   role + ".document",
				Reason: "required document is missing",
			})
		}
	}
	sort.Strings(targets)
	return targets, missing, nil
}

func inspectParsedDocumentCompletion(
	document *parsedIDDDocument,
	role string,
	includeRecordFields bool,
) []IncompleteSlot {
	schema, ok := iddDocumentSchemas[role]
	if !ok {
		return nil
	}

	incomplete := make(map[string]IncompleteSlot)
	root := goldmark.DefaultParser().Parse(goldmarktext.NewReader(document.Body))
	for _, marker := range scaffoldMarkersInDocument(document, root) {
		slot := marker.Slot
		incomplete[slot] = IncompleteSlot{
			File:   document.Path,
			Line:   marker.Line,
			Role:   role,
			Slot:   slot,
			Reason: "generated scaffold marker remains",
		}
	}

	records := splitIDDMarkdownRecords(document, root)
	if !hasRecordPrefix(records, schema.Collection.RecordPrefix) {
		addMissingCompletionSlot(incomplete, IncompleteSlot{
			File:   document.Path,
			Line:   completionFallbackLine(document, schema.Collection.Name),
			Role:   role,
			Slot:   schema.Collection.Name,
			Reason: "required record has not been authored",
		})
	}
	for _, section := range schema.Sections {
		record, found := findCompletionSection(records, section.Heading)
		if !found {
			addMissingCompletionSlot(incomplete, IncompleteSlot{
				File:   document.Path,
				Line:   completionFallbackLine(document, section.Name),
				Role:   role,
				Slot:   section.Name,
				Reason: fmt.Sprintf("required section %q is missing", section.Heading),
			})
			continue
		}
		if !completionRecordHasEffectiveContent(record, document.Body) {
			addMissingCompletionSlot(incomplete, IncompleteSlot{
				File:   document.Path,
				Line:   record.line,
				Role:   role,
				Slot:   section.Name,
				Reason: fmt.Sprintf("required section %q has no authored content", section.Heading),
			})
		}
	}
	if includeRecordFields {
		for _, record := range iddDocumentRecordValuesForRole(document, role) {
			for _, field := range schema.Record.Fields {
				if !iddDocumentRecordFieldRequired(field, record.Kind) {
					continue
				}
				value := record.scalarValue(field.Name)
				if field.List {
					value = strings.Join(record.Lists[field.Name], ", ")
				}
				placeholderRecordID := record.Identifier
				if field.Name == "name" {
					placeholderRecordID = ""
				}
				if iddDocumentRecordFieldComplete(value, placeholderRecordID, field.Name) {
					continue
				}
				reason := fmt.Sprintf(
					"%s %s field %s requires concrete authored content",
					schema.Record.Type,
					record.displayName(),
					field.Label,
				)
				addMissingCompletionSlot(incomplete, IncompleteSlot{
					File: document.Path,
					Line: iddDocumentRecordFieldLine(
						document,
						schema.Record.Section,
						record.Index,
						field.Name,
					),
					Role:   role,
					Slot:   iddDocumentRecordFieldSlot(role, record, field.Name),
					Reason: reason,
				})
			}
		}
	}

	slots := make([]IncompleteSlot, 0, len(incomplete))
	for _, slot := range incomplete {
		slots = append(slots, slot)
	}
	sortIncompleteSlots(slots)
	return slots
}

type iddScaffoldMarker struct {
	Slot string
	Line int
}

func scaffoldMarkersInDocument(
	document *parsedIDDDocument,
	root goldmarkast.Node,
) []iddScaffoldMarker {
	var markers []iddScaffoldMarker
	_ = goldmarkast.Walk(root, func(node goldmarkast.Node, entering bool) (goldmarkast.WalkStatus, error) {
		if !entering {
			return goldmarkast.WalkContinue, nil
		}

		var raw []byte
		switch typed := node.(type) {
		case *goldmarkast.HTMLBlock:
			raw = append(raw, typed.Lines().Value(document.Body)...)
			if typed.HasClosure() {
				raw = append(raw, typed.ClosureLine.Value(document.Body)...)
			}
		case *goldmarkast.RawHTML:
			raw = typed.Segments.Value(document.Body)
		default:
			return goldmarkast.WalkContinue, nil
		}

		baseLine := markdownNodeLine(document, node)
		for lineIndex, line := range strings.Split(string(raw), "\n") {
			match := scaffoldMarkerPattern.FindStringSubmatch(line)
			if len(match) != 2 {
				continue
			}
			markers = append(markers, iddScaffoldMarker{
				Slot: match[1],
				Line: baseLine + lineIndex,
			})
		}
		return goldmarkast.WalkContinue, nil
	})
	return markers
}

type iddDocumentRecordValues struct {
	Index      int
	Identifier string
	Kind       string
	Scalars    map[string]string
	Lists      map[string][]string
}

func (record iddDocumentRecordValues) scalarValue(field string) string {
	if record.Scalars == nil {
		return ""
	}
	return record.Scalars[field]
}

func (record iddDocumentRecordValues) displayName() string {
	if strings.TrimSpace(record.Identifier) != "" {
		return record.Identifier
	}
	return fmt.Sprintf("record %d", record.Index+1)
}

func iddDocumentRecordValuesForRole(
	document *parsedIDDDocument,
	role string,
) []iddDocumentRecordValues {
	if document == nil || document.Metadata == nil {
		return nil
	}
	var records []iddDocumentRecordValues
	switch role {
	case "design":
		for index, component := range document.Metadata.ComponentRecords {
			records = append(records, iddDocumentRecordValues{
				Index:      index,
				Identifier: component.Name,
				Scalars: map[string]string{
					"name":    component.Name,
					"purpose": component.Purpose,
				},
			})
		}
	case "contract":
		for index, contract := range document.Metadata.ContractRecords {
			records = append(records, iddDocumentRecordValues{
				Index:      index,
				Identifier: contract.Name,
				Scalars: map[string]string{
					"name":       contract.Name,
					"guarantees": contract.Guarantees,
				},
			})
		}
	case "spec":
		for index, spec := range document.Metadata.Specs {
			records = append(records, iddDocumentRecordValues{
				Index:      index,
				Identifier: spec.ID,
				Scalars: map[string]string{
					"id":          spec.ID,
					"title":       spec.Title,
					"design":      spec.Design,
					"contract":    spec.Contract,
					"requirement": spec.Requirement,
					"acceptance":  spec.Acceptance,
				},
			})
		}
	case "testing":
		for index, test := range document.Metadata.Tests {
			records = append(records, iddDocumentRecordValues{
				Index:      index,
				Identifier: test.ID,
				Kind:       test.Kind,
				Scalars: map[string]string{
					"id":      test.ID,
					"title":   test.Title,
					"kind":    test.Kind,
					"purpose": test.Purpose,
					"oracle":  test.Oracle,
				},
				Lists: map[string][]string{
					"covers":    test.Covers,
					"contracts": test.Contracts,
				},
			})
		}
	}
	return records
}

func iddDocumentRecordFieldRequired(
	field iddDocumentRecordFieldSchema,
	recordKind string,
) bool {
	return field.RequiredWhenKind == "" || field.RequiredWhenKind == recordKind
}

func iddDocumentRecordFieldComplete(value, recordID, field string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	for _, item := range strings.Split(value, ",") {
		if isIDDPlaceholder(strings.TrimSpace(item), recordID, field) {
			return false
		}
	}
	return true
}

func iddDocumentRecordFieldLine(
	document *parsedIDDDocument,
	section string,
	recordIndex int,
	field string,
) int {
	line := document.Index.recordFieldLine(section, recordIndex, field)
	if line <= 0 {
		line = document.Index.recordLine(section, recordIndex)
	}
	if line <= 0 {
		line = document.BodyStartLine
	}
	if line <= 0 {
		return 1
	}
	return line
}

func iddDocumentRecordFieldSlot(
	role string,
	record iddDocumentRecordValues,
	field string,
) string {
	recordKey := strings.ToLower(record.Identifier)
	recordKey = completionSlotPartPattern.ReplaceAllString(recordKey, "-")
	recordKey = strings.Trim(recordKey, "-")
	if recordKey == "" {
		recordKey = fmt.Sprintf("record-%d", record.Index+1)
	}
	return role + "." + recordKey + "." + field
}

func hasRecordPrefix(records []iddMarkdownRecord, prefix string) bool {
	for _, record := range records {
		if strings.HasPrefix(record.heading, prefix) &&
			strings.TrimSpace(strings.TrimPrefix(record.heading, prefix)) != "" {
			return true
		}
	}
	return false
}

func findCompletionSection(records []iddMarkdownRecord, heading string) (iddMarkdownRecord, bool) {
	for _, record := range records {
		if strings.EqualFold(strings.TrimSpace(record.heading), heading) {
			return record, true
		}
	}
	return iddMarkdownRecord{}, false
}

func completionRecordHasEffectiveContent(record iddMarkdownRecord, source []byte) bool {
	for _, block := range record.blocks {
		switch block.(type) {
		case *goldmarkast.Heading, *goldmarkast.Blockquote, *goldmarkast.HTMLBlock, *goldmarkast.RawHTML:
			continue
		}
		text := markdownNodeText(block, source)
		if text == "" {
			if lines := block.Lines(); lines != nil && lines.Len() > 0 {
				text = strings.TrimSpace(string(lines.Value(source)))
			}
		}
		if text != "" && !isIDDPlaceholder(text, "", "") {
			return true
		}
	}
	return false
}

func completionFallbackLine(document *parsedIDDDocument, slot string) int {
	root := goldmark.DefaultParser().Parse(goldmarktext.NewReader(document.Body))
	for _, marker := range scaffoldMarkersInDocument(document, root) {
		if marker.Slot == slot {
			return marker.Line
		}
	}
	if document.BodyStartLine > 0 {
		return document.BodyStartLine
	}
	return 1
}

func addMissingCompletionSlot(slots map[string]IncompleteSlot, slot IncompleteSlot) {
	if _, exists := slots[slot.Slot]; !exists {
		slots[slot.Slot] = slot
	}
}

func sortIncompleteSlots(slots []IncompleteSlot) {
	sort.Slice(slots, func(i, j int) bool {
		if slots[i].File != slots[j].File {
			return slots[i].File < slots[j].File
		}
		if slots[i].Line != slots[j].Line {
			return slots[i].Line < slots[j].Line
		}
		return slots[i].Slot < slots[j].Slot
	})
}
