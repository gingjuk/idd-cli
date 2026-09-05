package reporter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/jingxu9x/idd-cli/internal/collector"
)

// RenderTraceDossier renders one stable ID-centered dossier.
// @implement SPEC-INTERNAL_REPORTER-013
func RenderTraceDossier(dossier *collector.TraceDossier, format string) ([]byte, error) {
	if dossier == nil {
		return nil, fmt.Errorf("trace dossier is nil")
	}
	stable := stableTraceDossier(dossier)
	switch format {
	case "json":
		data, err := json.MarshalIndent(stable, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal trace dossier: %w", err)
		}
		return append(data, '\n'), nil
	case "markdown", "llm-markdown":
		return []byte(renderTraceMarkdown(stable)), nil
	default:
		return nil, fmt.Errorf(
			"unsupported trace format %q; expected json, markdown, or llm-markdown",
			format,
		)
	}
}

func stableTraceDossier(source *collector.TraceDossier) *collector.TraceDossier {
	result := *source
	result.Owners = append([]collector.TraceOccurrence{}, source.Owners...)
	result.Occurrences = append([]collector.TraceOccurrence{}, source.Occurrences...)
	result.OutboundRelations = append([]collector.TraceRelation{}, source.OutboundRelations...)
	result.InboundRelations = append([]collector.TraceRelation{}, source.InboundRelations...)
	result.RelatedEntities = append([]collector.TraceRelatedEntity{}, source.RelatedEntities...)
	result.Components = append([]collector.TraceComponent{}, source.Components...)
	result.Contracts = append([]collector.TraceContract{}, source.Contracts...)
	result.Implementations = append([]collector.ReviewContextDeclaration{}, source.Implementations...)
	result.Tests = append([]collector.TraceTest{}, source.Tests...)
	result.Mentions = append([]collector.TraceOccurrence{}, source.Mentions...)
	result.Findings = append([]collector.ReviewContextIssue{}, source.Findings...)

	sort.Slice(result.Owners, func(i, j int) bool { return lessTraceOccurrence(result.Owners[i], result.Owners[j]) })
	sort.Slice(result.Occurrences, func(i, j int) bool { return lessTraceOccurrence(result.Occurrences[i], result.Occurrences[j]) })
	sort.Slice(result.Mentions, func(i, j int) bool { return lessTraceOccurrence(result.Mentions[i], result.Mentions[j]) })
	sort.Slice(result.OutboundRelations, func(i, j int) bool {
		return lessTraceRelation(result.OutboundRelations[i], result.OutboundRelations[j])
	})
	sort.Slice(result.InboundRelations, func(i, j int) bool { return lessTraceRelation(result.InboundRelations[i], result.InboundRelations[j]) })
	sort.Slice(result.RelatedEntities, func(i, j int) bool {
		if result.RelatedEntities[i].Depth != result.RelatedEntities[j].Depth {
			return result.RelatedEntities[i].Depth < result.RelatedEntities[j].Depth
		}
		return result.RelatedEntities[i].ID < result.RelatedEntities[j].ID
	})
	sort.Slice(result.Components, func(i, j int) bool {
		if result.Components[i].ID != result.Components[j].ID {
			return result.Components[i].ID < result.Components[j].ID
		}
		return lessLocation(result.Components[i].File, result.Components[i].Line, result.Components[j].File, result.Components[j].Line)
	})
	sort.Slice(result.Contracts, func(i, j int) bool {
		if result.Contracts[i].ID != result.Contracts[j].ID {
			return result.Contracts[i].ID < result.Contracts[j].ID
		}
		return lessLocation(result.Contracts[i].File, result.Contracts[i].Line, result.Contracts[j].File, result.Contracts[j].Line)
	})
	sort.Slice(result.Implementations, func(i, j int) bool { return lessDeclaration(result.Implementations[i], result.Implementations[j]) })
	for index := range result.Tests {
		result.Tests[index].Covers = sortedStrings(result.Tests[index].Covers)
		result.Tests[index].Contracts = sortedStrings(result.Tests[index].Contracts)
		result.Tests[index].Declarations = append([]collector.ReviewContextDeclaration{}, result.Tests[index].Declarations...)
		sort.Slice(result.Tests[index].Declarations, func(i, j int) bool {
			return lessDeclaration(result.Tests[index].Declarations[i], result.Tests[index].Declarations[j])
		})
	}
	sort.Slice(result.Tests, func(i, j int) bool {
		if result.Tests[i].ID != result.Tests[j].ID {
			return result.Tests[i].ID < result.Tests[j].ID
		}
		return lessLocation(result.Tests[i].File, result.Tests[i].Line, result.Tests[j].File, result.Tests[j].Line)
	})
	sort.Slice(result.Findings, func(i, j int) bool {
		if result.Findings[i].Rule != result.Findings[j].Rule {
			return result.Findings[i].Rule < result.Findings[j].Rule
		}
		if result.Findings[i].Source != result.Findings[j].Source {
			return result.Findings[i].Source < result.Findings[j].Source
		}
		return result.Findings[i].Message < result.Findings[j].Message
	})
	return &result
}

func lessTraceOccurrence(left, right collector.TraceOccurrence) bool {
	if left.Kind != right.Kind {
		return left.Kind < right.Kind
	}
	if left.EntityID != right.EntityID {
		return left.EntityID < right.EntityID
	}
	if left.Path != right.Path {
		return left.Path < right.Path
	}
	if left.Line != right.Line {
		return left.Line < right.Line
	}
	if left.Field != right.Field {
		return left.Field < right.Field
	}
	if left.RecordID != right.RecordID {
		return left.RecordID < right.RecordID
	}
	return left.ID < right.ID
}

func lessTraceRelation(left, right collector.TraceRelation) bool {
	if left.Type != right.Type {
		return left.Type < right.Type
	}
	if left.From != right.From {
		return left.From < right.From
	}
	if left.To != right.To {
		return left.To < right.To
	}
	if left.Provenance.Path != right.Provenance.Path {
		return left.Provenance.Path < right.Provenance.Path
	}
	if left.Provenance.Line != right.Provenance.Line {
		return left.Provenance.Line < right.Provenance.Line
	}
	if left.Provenance.Field != right.Provenance.Field {
		return left.Provenance.Field < right.Provenance.Field
	}
	if left.Provenance.RecordID != right.Provenance.RecordID {
		return left.Provenance.RecordID < right.Provenance.RecordID
	}
	return left.Provenance.OccurrenceID < right.Provenance.OccurrenceID
}

func lessDeclaration(left, right collector.ReviewContextDeclaration) bool {
	if left.File != right.File {
		return left.File < right.File
	}
	if left.Line != right.Line {
		return left.Line < right.Line
	}
	return left.Name < right.Name
}

func lessLocation(leftPath string, leftLine int, rightPath string, rightLine int) bool {
	if leftPath != rightPath {
		return leftPath < rightPath
	}
	return leftLine < rightLine
}

func sortedStrings(values []string) []string {
	result := append([]string{}, values...)
	sort.Strings(result)
	return result
}

func renderTraceMarkdown(dossier *collector.TraceDossier) string {
	var output strings.Builder
	writeMarkdownHeading(&output, 1, "IDD trace: "+dossier.Query.ID)
	fmt.Fprintf(&output, "Schema: `%s`\n\n", dossier.Schema)
	fmt.Fprintf(&output, "Status: `%s`\n\n", dossier.Status)
	fmt.Fprintf(&output, "Query depth: `%d`; mentions included: `%t`\n\n", dossier.Query.Depth, dossier.Query.IncludeMentions)

	writeMarkdownHeading(&output, 2, "Entity")
	if dossier.Entity == nil {
		output.WriteString("No logical entity was resolved. Typed occurrences and findings are retained below.\n\n")
	} else {
		fmt.Fprintf(&output, "- ID: `%s`\n- Kind: `%s`\n- Resolution: `%s`\n", dossier.Entity.ID, dossier.Entity.Kind, dossier.Entity.Resolution)
		if dossier.Entity.Namespace != "" {
			fmt.Fprintf(&output, "- Namespace: `%s`\n", dossier.Entity.Namespace)
		}
		if dossier.Entity.Lifecycle != "" {
			fmt.Fprintf(&output, "- Lifecycle: `%s`\n", dossier.Entity.Lifecycle)
		}
		output.WriteString("\n")
	}

	writeMarkdownHeading(&output, 2, "Canonical declaration")
	if dossier.Canonical == nil {
		output.WriteString("No unique canonical declaration was resolved.\n\n")
	} else {
		fmt.Fprintf(&output, "Location: `%s:%d`\n\n", dossier.Canonical.Path, dossier.Canonical.Line)
		writeTraceCanonicalFields(&output, dossier.Canonical)
		if dossier.Canonical.Markdown != "" {
			writeFencedBlock(&output, "markdown", dossier.Canonical.Markdown)
		}
		fmt.Fprintf(&output, "Record truncated: `%t`\n\n", dossier.Canonical.Truncated)
	}

	writeTraceOccurrences(&output, "Canonical owners", dossier.Owners)
	writeTraceComponents(&output, dossier.Components)
	writeTraceContracts(&output, dossier.Contracts)
	writeTraceDeclarations(&output, "Implementation declarations", dossier.Implementations)
	writeTraceTests(&output, dossier.Tests)
	writeTraceRelations(&output, "Outbound relations", dossier.OutboundRelations)
	writeTraceRelations(&output, "Inbound relations", dossier.InboundRelations)
	writeTraceOccurrences(&output, "Typed occurrences", dossier.Occurrences)
	if dossier.Query.IncludeMentions {
		writeTraceOccurrences(&output, "Mentions", dossier.Mentions)
	}
	if len(dossier.RelatedEntities) > 0 {
		writeMarkdownHeading(&output, 2, "Related entities")
		for _, entity := range dossier.RelatedEntities {
			fmt.Fprintf(&output, "- `%s` (%s, %s, depth %d)\n", entity.ID, entity.Kind, entity.Resolution, entity.Depth)
		}
		output.WriteString("\n")
	}
	if len(dossier.Findings) > 0 {
		writeMarkdownHeading(&output, 2, "Findings")
		for _, finding := range dossier.Findings {
			fmt.Fprintf(&output, "- `%s`", finding.Rule)
			if finding.Source != "" {
				fmt.Fprintf(&output, " at `%s`", finding.Source)
			}
			fmt.Fprintf(&output, ": %s\n", finding.Message)
		}
		output.WriteString("\n")
	}
	if dossier.Truncated {
		output.WriteString("Some authored records or declaration excerpts were truncated by the dossier budget.\n")
	}
	return output.String()
}

func writeTraceCanonicalFields(output *strings.Builder, canonical *collector.TraceCanonical) {
	fields := []struct{ label, value string }{
		{"Title", canonical.Title}, {"Package", canonical.Package}, {"Purpose", canonical.Purpose},
		{"Ownership", canonical.Ownership}, {"Boundary", canonical.Boundary}, {"Decisions", canonical.Decisions},
		{"Guarantees", canonical.Guarantees}, {"Requirement", canonical.Requirement}, {"Acceptance", canonical.Acceptance},
		{"Kind", canonical.Kind}, {"Oracle", canonical.Oracle},
	}
	for _, field := range fields {
		if field.value != "" {
			fmt.Fprintf(output, "- %s: %s\n", field.label, field.value)
		}
	}
	output.WriteString("\n")
}

func writeTraceOccurrences(output *strings.Builder, title string, occurrences []collector.TraceOccurrence) {
	writeMarkdownHeading(output, 2, title)
	if len(occurrences) == 0 {
		output.WriteString("None.\n\n")
		return
	}
	for _, occurrence := range occurrences {
		fmt.Fprintf(output, "- `%s:%d` — `%s` / `%s`", occurrence.Path, occurrence.Line, occurrence.Origin, occurrence.Kind)
		if occurrence.Field != "" {
			fmt.Fprintf(output, ", field `%s`", occurrence.Field)
		}
		if occurrence.RecordID != "" {
			fmt.Fprintf(output, ", record `%s`", occurrence.RecordID)
		}
		output.WriteString("\n")
	}
	output.WriteString("\n")
}

func writeTraceComponents(output *strings.Builder, components []collector.TraceComponent) {
	writeMarkdownHeading(output, 2, "Components and design context")
	if len(components) == 0 {
		output.WriteString("No related Component was collected.\n\n")
		return
	}
	for _, component := range components {
		writeMarkdownHeading(output, 3, component.Name)
		fmt.Fprintf(output, "- ID: `%s`\n- Location: `%s:%d`\n\n", component.ID, component.File, component.Line)
		for _, field := range []struct{ label, value string }{{"Purpose", component.Purpose}, {"Ownership", component.Ownership}, {"Boundary", component.Boundary}, {"Decisions", component.Decisions}} {
			if field.value != "" {
				fmt.Fprintf(output, "%s: %s\n\n", field.label, field.value)
			}
		}
		writeFencedBlock(output, "markdown", component.Markdown)
	}
}

func writeTraceContracts(output *strings.Builder, contracts []collector.TraceContract) {
	writeMarkdownHeading(output, 2, "Contracts")
	if len(contracts) == 0 {
		output.WriteString("No related Contract was collected.\n\n")
		return
	}
	for _, contract := range contracts {
		writeMarkdownHeading(output, 3, contract.Name)
		fmt.Fprintf(output, "- ID: `%s`\n- Location: `%s:%d`\n\nGuarantees: %s\n\n", contract.ID, contract.File, contract.Line, contract.Guarantees)
		writeFencedBlock(output, "markdown", contract.Markdown)
	}
}

func writeTraceDeclarations(output *strings.Builder, title string, declarations []collector.ReviewContextDeclaration) {
	writeMarkdownHeading(output, 2, title)
	if len(declarations) == 0 {
		output.WriteString("None.\n\n")
		return
	}
	for _, declaration := range declarations {
		writeMarkdownHeading(output, 3, declaration.Kind+" "+declaration.Name)
		fmt.Fprintf(output, "- Role: `%s`\n- Language: `%s`\n- Location: `%s:%d-%d`\n- Annotation: `%s`\n\n", declaration.Role, declaration.Language, declaration.File, declaration.Line, declaration.EndLine, strings.TrimSpace(declaration.Annotation))
		writeFencedBlock(output, markdownLanguage(declaration.Language), declaration.Excerpt)
	}
}

func writeTraceTests(output *strings.Builder, tests []collector.TraceTest) {
	writeMarkdownHeading(output, 2, "Test evidence")
	if len(tests) == 0 {
		output.WriteString("No related TEST evidence was collected.\n\n")
		return
	}
	for _, test := range tests {
		writeMarkdownHeading(output, 3, test.ID+": "+test.Title)
		fmt.Fprintf(output, "- Kind: `%s`\n- Location: `%s:%d`\n- Covers: `%s`\n", test.Kind, test.File, test.Line, strings.Join(test.Covers, "`, `"))
		if len(test.Contracts) > 0 {
			fmt.Fprintf(output, "- Contracts: `%s`\n", strings.Join(test.Contracts, "`, `"))
		}
		fmt.Fprintf(output, "\nPurpose: %s\n\nOracle: %s\n\n", test.Purpose, test.Oracle)
		writeFencedBlock(output, "markdown", test.Markdown)
		if len(test.Declarations) > 0 {
			writeTraceDeclarations(output, "Executable declarations for "+test.ID, test.Declarations)
		}
	}
}

func writeTraceRelations(output *strings.Builder, title string, relations []collector.TraceRelation) {
	writeMarkdownHeading(output, 2, title)
	if len(relations) == 0 {
		output.WriteString("None.\n\n")
		return
	}
	for _, relation := range relations {
		fmt.Fprintf(output, "- `%s` --`%s`--> `%s` (depth %d; `%s:%d`", relation.From, relation.Type, relation.To, relation.Depth, relation.Provenance.Path, relation.Provenance.Line)
		if relation.Provenance.Field != "" {
			fmt.Fprintf(output, ", field `%s`", relation.Provenance.Field)
		}
		output.WriteString(")\n")
	}
	output.WriteString("\n")
}
