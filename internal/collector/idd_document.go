// Package collector provides documentation identifier collection functionality.

// Spec: docs/internal/collector/spec.md
// Contract: docs/internal/collector/contract.md
package collector

import (
	"bytes"
	stderrors "errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jingxu9x/idd-cli/internal/model"
	"github.com/jingxu9x/idd-cli/pkg/pattern"
	"github.com/yuin/goldmark"
	goldmarkast "github.com/yuin/goldmark/ast"
	goldmarktext "github.com/yuin/goldmark/text"
	"gopkg.in/yaml.v3"
)

const iddDocumentVersion = "1.0"

var (
	iddDocumentRoles = map[string]string{
		"design.md":   "design",
		"contract.md": "contract",
		"spec.md":     "spec",
		"testing.md":  "testing",
	}
	iddDocumentOrder = []string{"design.md", "contract.md", "spec.md", "testing.md"}
)

// IDDDocument holds the minimal identity from frontmatter plus semantic records
// parsed from the human-readable Markdown body. Semantic slices are never
// marshaled back into YAML.
//
// @implement SPEC-INTERNAL_COLLECTOR-001
type IDDDocument struct {
	Version    string            `yaml:"version"`
	Package    string            `yaml:"package"`
	Document   string            `yaml:"document"`
	Components []string          `yaml:"-"`
	Contracts  []string          `yaml:"-"`
	Specs      []IDDDocumentSpec `yaml:"-"`
	Tests      []IDDDocumentTest `yaml:"-"`
}

// IDDDocumentSpec is the single structured owner for a SPEC declaration.
//
// @implement SPEC-INTERNAL_COLLECTOR-001
type IDDDocumentSpec struct {
	ID          string `yaml:"id"`
	Title       string `yaml:"title"`
	Requirement string `yaml:"requirement"`
	Design      string `yaml:"design"`
	Contract    string `yaml:"contract"`
}

// IDDDocumentTest is the single structured owner for a TEST declaration and
// its SPEC coverage.
//
// @implement SPEC-INTERNAL_COLLECTOR-001
type IDDDocumentTest struct {
	ID      string   `yaml:"id"`
	Title   string   `yaml:"title"`
	Purpose string   `yaml:"purpose"`
	Kind    string   `yaml:"kind"`
	Covers  []string `yaml:"covers"`
}

type parsedIDDDocument struct {
	Path          string
	Metadata      *IDDDocument
	Body          []byte
	Raw           []byte
	Root          *yaml.Node
	Index         iddDocumentNodeIndex
	Issues        []iddMarkdownIssue
	BodyStartLine int
}

type iddMarkdownIssue struct {
	Rule    string
	Message string
	Line    int
	Link    string
	Code    string
}

type iddSemanticFrontmatterError struct {
	Field string
	Line  int
}

func (e *iddSemanticFrontmatterError) Error() string {
	return fmt.Sprintf(
		"frontmatter field %q stores semantic records; move each value into a human-readable Markdown record",
		e.Field,
	)
}

// ParseIDDDocument reads an IDD metadata block from leading Markdown
// frontmatter. A legacy document without an `idd` key returns a nil document
// and no error.
//
// @implement SPEC-INTERNAL_COLLECTOR-001
func ParseIDDDocument(data []byte) (*IDDDocument, []byte, error) {
	parsed, detected, err := parseIDDDocument("", data)
	if err != nil {
		return nil, nil, err
	}
	if !detected {
		return nil, data, nil
	}
	return parsed.Metadata, parsed.Body, nil
}

// MarshalIDDDocument emits canonical minimal identity frontmatter and preserves the
// supplied Markdown body exactly.
//
// @implement SPEC-INTERNAL_COLLECTOR-001
func MarshalIDDDocument(document *IDDDocument, body []byte) ([]byte, error) {
	if document == nil {
		return nil, fmt.Errorf("IDD document is nil")
	}
	return marshalIDDDocument(document, body, nil)
}

func parseIDDDocument(path string, data []byte) (*parsedIDDDocument, bool, error) {
	frontmatter, body, rootLineOffset, bodyStartLine, found, err := splitLeadingFrontmatter(data)
	if err != nil {
		return nil, frontmatterContainsIDDKey(frontmatter), err
	}
	if !found {
		return nil, false, nil
	}

	var yamlDocument yaml.Node
	if err := yaml.Unmarshal(frontmatter, &yamlDocument); err != nil {
		return nil,
			frontmatterContainsIDDKey(frontmatter),
			fmt.Errorf("decode IDD frontmatter at line %d: %w", yamlErrorLine(err)+rootLineOffset, err)
	}
	root := documentRoot(&yamlDocument)
	iddNode := mappingValue(root, "idd")
	if iddNode == nil {
		return nil, false, nil
	}

	for _, field := range []string{"components", "contracts", "specs", "tests"} {
		if value := mappingValue(iddNode, field); value != nil {
			line := value.Line + rootLineOffset
			if line <= 0 {
				line = 1
			}
			return nil, true, &iddSemanticFrontmatterError{Field: field, Line: line}
		}
	}

	metadata, err := decodeIDDMetadata(iddNode)
	if err != nil {
		metadataLineOffset := rootLineOffset
		if iddNode.Line > 0 {
			metadataLineOffset += iddNode.Line - 1
		}
		return nil, true, fmt.Errorf(
			"decode IDD metadata at line %d: %w",
			yamlErrorLine(err)+metadataLineOffset,
			err,
		)
	}
	parsed := &parsedIDDDocument{
		Path:     path,
		Metadata: metadata,
		Body:     body,
		Raw:      data,
		Root:     root,
		Index: iddDocumentNodeIndex{
			root:             iddNode,
			lineOffset:       rootLineOffset,
			recordLines:      make(map[string][]int),
			recordFieldLines: make(map[string][]map[string]int),
		},
		BodyStartLine: bodyStartLine,
	}
	role := metadata.Document
	if pathRole, ok := iddDocumentRoles[filepath.Base(path)]; ok {
		role = pathRole
	}
	parseIDDMarkdownRecords(parsed, role)
	return parsed, true, nil
}

func decodeIDDMetadata(node *yaml.Node) (*IDDDocument, error) {
	var encoded bytes.Buffer
	encoder := yaml.NewEncoder(&encoded)
	if err := encoder.Encode(node); err != nil {
		return nil, fmt.Errorf("encode IDD metadata for validation: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return nil, fmt.Errorf("close IDD metadata encoder: %w", err)
	}

	decoder := yaml.NewDecoder(bytes.NewReader(encoded.Bytes()))
	decoder.KnownFields(true)
	var metadata IDDDocument
	if err := decoder.Decode(&metadata); err != nil {
		return nil, fmt.Errorf("decode IDD metadata: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("decode IDD metadata: multiple YAML documents are not supported")
		}
		return nil, fmt.Errorf("decode IDD metadata: %w", err)
	}
	return &metadata, nil
}

func splitLeadingFrontmatter(data []byte) (
	frontmatter []byte,
	body []byte,
	rootLineOffset int,
	bodyStartLine int,
	found bool,
	err error,
) {
	lines := strings.SplitAfter(string(data), "\n")
	start := -1
	for index, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.TrimSpace(line) != "---" {
			return nil, data, 0, 1, false, nil
		}
		start = index
		break
	}
	if start == -1 {
		return nil, data, 0, 1, false, nil
	}

	end := -1
	for index := start + 1; index < len(lines); index++ {
		if strings.TrimSpace(lines[index]) == "---" {
			end = index
			break
		}
	}
	if end == -1 {
		partial := []byte(strings.Join(lines[start+1:], ""))
		return partial, nil, start + 1, len(lines), true, fmt.Errorf("decode IDD frontmatter: closing --- delimiter is missing")
	}

	return []byte(strings.Join(lines[start+1:end], "")),
		[]byte(strings.Join(lines[end+1:], "")),
		start + 1,
		end + 2,
		true,
		nil
}

func frontmatterContainsIDDKey(frontmatter []byte) bool {
	for _, line := range strings.Split(string(frontmatter), "\n") {
		if strings.TrimLeft(line, " \t") != line {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "idd:") {
			return true
		}
	}
	return false
}

func hasIDDDocumentFrontmatter(data []byte) bool {
	frontmatter, _, _, _, found, err := splitLeadingFrontmatter(data)
	if !found {
		return false
	}
	if err != nil {
		return frontmatterContainsIDDKey(frontmatter)
	}
	var document yaml.Node
	if yaml.Unmarshal(frontmatter, &document) != nil {
		return frontmatterContainsIDDKey(frontmatter)
	}
	return mappingValue(documentRoot(&document), "idd") != nil
}

func marshalIDDDocument(document *IDDDocument, body []byte, existingRoot *yaml.Node) ([]byte, error) {
	root := existingRoot
	if root == nil {
		root = mappingNode()
	}
	if root.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("frontmatter root must be a YAML mapping")
	}
	replaceMapping(root, "idd", iddMetadataNode(document))

	yamlDocument := &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{root}}
	var encoded bytes.Buffer
	encoder := yaml.NewEncoder(&encoded)
	encoder.SetIndent(2)
	if err := encoder.Encode(yamlDocument); err != nil {
		return nil, fmt.Errorf("encode IDD document frontmatter: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return nil, fmt.Errorf("close IDD document encoder: %w", err)
	}

	var output bytes.Buffer
	output.WriteString("---\n")
	output.Write(encoded.Bytes())
	output.WriteString("---\n")
	output.Write(body)
	return output.Bytes(), nil
}

func iddMetadataNode(document *IDDDocument) *yaml.Node {
	root := mappingNode()
	appendMapping(root, scalarNode("version"), quotedScalarNode(document.Version))
	appendMapping(root, scalarNode("package"), scalarNode(document.Package))
	appendMapping(root, scalarNode("document"), scalarNode(document.Document))
	return root
}

type iddMarkdownRecord struct {
	heading string
	line    int
	blocks  []goldmarkast.Node
}

func parseIDDMarkdownRecords(document *parsedIDDDocument, role string) {
	source := document.Body
	root := goldmark.DefaultParser().Parse(goldmarktext.NewReader(source))
	records := splitIDDMarkdownRecords(document, root)

	switch role {
	case "design":
		parseIDDNameRecords(document, records, "Component:", "components", "component")
	case "contract":
		parseIDDNameRecords(document, records, "Contract:", "contracts", "contract")
	case "spec":
		parseIDDSpecRecords(document, records)
		document.Issues = append(document.Issues, findIDDRecordTables(document, "SPEC-")...)
	case "testing":
		parseIDDTestRecords(document, records)
		document.Issues = append(document.Issues, findIDDRecordTables(document, "TEST-")...)
	}
}

func splitIDDMarkdownRecords(document *parsedIDDDocument, root goldmarkast.Node) []iddMarkdownRecord {
	var records []iddMarkdownRecord
	var current *iddMarkdownRecord
	for child := root.FirstChild(); child != nil; child = child.NextSibling() {
		heading, ok := child.(*goldmarkast.Heading)
		if ok && heading.Level == 2 {
			if current != nil {
				records = append(records, *current)
			}
			current = &iddMarkdownRecord{
				heading: markdownNodeText(heading, document.Body),
				line:    markdownNodeLine(document, heading),
			}
			continue
		}
		if current != nil {
			current.blocks = append(current.blocks, child)
		}
	}
	if current != nil {
		records = append(records, *current)
	}
	return records
}

func parseIDDNameRecords(
	document *parsedIDDDocument,
	records []iddMarkdownRecord,
	headingPrefix string,
	section string,
	code string,
) {
	for _, record := range records {
		if !strings.HasPrefix(record.heading, headingPrefix) {
			continue
		}
		name := strings.TrimSpace(strings.TrimPrefix(record.heading, headingPrefix))
		name = trimMarkdownScalar(name)
		switch section {
		case "components":
			document.Metadata.Components = append(document.Metadata.Components, name)
		case "contracts":
			document.Metadata.Contracts = append(document.Metadata.Contracts, name)
		}
		document.Index.addRecord(section, record.line, map[string]int{"name": record.line})
		if name == "" {
			document.Issues = append(document.Issues, iddMarkdownIssue{
				Rule:    "idd-document-schema",
				Message: strings.TrimSuffix(headingPrefix, ":") + " heading requires a concrete name",
				Line:    record.line,
				Code:    code,
			})
		}
		if !iddRecordHasNarrative(record, document.Body) {
			document.Issues = append(document.Issues, iddMarkdownIssue{
				Rule:    "idd-document-schema",
				Message: fmt.Sprintf("%s %q requires a human-readable narrative", strings.TrimSuffix(headingPrefix, ":"), name),
				Line:    record.line,
				Code:    code,
			})
		}
	}
}

func parseIDDSpecRecords(document *parsedIDDDocument, records []iddMarkdownRecord) {
	for _, record := range records {
		if !strings.HasPrefix(record.heading, "SPEC-") {
			continue
		}
		id, title := parseIDDRecordHeading(record.heading)
		fields, fieldLines, fieldIssues := parseIDDRecordFields(
			document,
			record,
			map[string]string{"design": "Design", "contract": "Contract"},
		)
		document.Issues = append(document.Issues, fieldIssues...)
		requirement, requirementLine, requirementIssues := parseIDDLabeledParagraph(document, record, "Requirement")
		document.Issues = append(document.Issues, requirementIssues...)

		document.Metadata.Specs = append(document.Metadata.Specs, IDDDocumentSpec{
			ID:          id,
			Title:       title,
			Requirement: requirement,
			Design:      fields["design"],
			Contract:    fields["contract"],
		})
		fieldLines["id"] = record.line
		fieldLines["title"] = record.line
		if requirementLine > 0 {
			fieldLines["requirement"] = requirementLine
		}
		document.Index.addRecord("specs", record.line, fieldLines)
	}
}

func parseIDDTestRecords(document *parsedIDDDocument, records []iddMarkdownRecord) {
	for _, record := range records {
		if !strings.HasPrefix(record.heading, "TEST-") {
			continue
		}
		id, title := parseIDDRecordHeading(record.heading)
		fields, fieldLines, fieldIssues := parseIDDRecordFields(
			document,
			record,
			map[string]string{"kind": "Kind", "covers": "Covers"},
		)
		document.Issues = append(document.Issues, fieldIssues...)
		purpose, purposeLine, purposeIssues := parseIDDLabeledParagraph(document, record, "Purpose")
		document.Issues = append(document.Issues, purposeIssues...)

		document.Metadata.Tests = append(document.Metadata.Tests, IDDDocumentTest{
			ID:      id,
			Title:   title,
			Purpose: purpose,
			Kind:    strings.ToLower(fields["kind"]),
			Covers:  extractIDDRefs(fields["covers"]),
		})
		fieldLines["id"] = record.line
		fieldLines["title"] = record.line
		if purposeLine > 0 {
			fieldLines["purpose"] = purposeLine
		}
		document.Index.addRecord("tests", record.line, fieldLines)
	}
}

func parseIDDRecordHeading(heading string) (string, string) {
	id, title, found := strings.Cut(heading, ":")
	if !found {
		return strings.TrimSpace(heading), ""
	}
	return strings.TrimSpace(id), strings.TrimSpace(title)
}

func parseIDDRecordFields(
	document *parsedIDDDocument,
	record iddMarkdownRecord,
	expected map[string]string,
) (map[string]string, map[string]int, []iddMarkdownIssue) {
	values := make(map[string]string, len(expected))
	lines := make(map[string]int, len(expected))
	var issues []iddMarkdownIssue

	if len(record.blocks) == 0 {
		return values, lines, issues
	}
	list, ok := record.blocks[0].(*goldmarkast.List)
	if !ok {
		return values, lines, issues
	}
	for child := list.FirstChild(); child != nil; child = child.NextSibling() {
		item, ok := child.(*goldmarkast.ListItem)
		if !ok {
			continue
		}
		text := markdownNodeText(item, document.Body)
		label, value, found := strings.Cut(text, ":")
		if !found {
			issues = append(issues, iddMarkdownIssue{
				Rule:    "idd-document-markdown",
				Message: "record metadata items must use '- **Field:** value'",
				Line:    markdownNodeLine(document, item),
				Link:    recordIdentifier(record.heading),
				Code:    "field-format",
			})
			continue
		}
		key := strings.ToLower(strings.TrimSpace(label))
		canonical, supported := expected[key]
		if !supported {
			issues = append(issues, iddMarkdownIssue{
				Rule:    "idd-document-schema",
				Message: fmt.Sprintf("record metadata field %q is not supported here", strings.TrimSpace(label)),
				Line:    markdownNodeLine(document, item),
				Link:    recordIdentifier(record.heading),
				Code:    key,
			})
			continue
		}
		if _, duplicate := values[key]; duplicate {
			issues = append(issues, iddMarkdownIssue{
				Rule:    "idd-document-schema",
				Message: fmt.Sprintf("record repeats %s metadata", canonical),
				Line:    markdownNodeLine(document, item),
				Link:    recordIdentifier(record.heading),
				Code:    key,
			})
			continue
		}
		values[key] = trimMarkdownScalar(value)
		lines[key] = markdownNodeLine(document, item)
	}
	return values, lines, issues
}

func parseIDDLabeledParagraph(
	document *parsedIDDDocument,
	record iddMarkdownRecord,
	label string,
) (string, int, []iddMarkdownIssue) {
	var value string
	var line int
	var issues []iddMarkdownIssue
	prefix := label + ":"
	for _, block := range record.blocks {
		if heading, ok := block.(*goldmarkast.Heading); ok && heading.Level >= 3 {
			break
		}
		paragraph, ok := block.(*goldmarkast.Paragraph)
		if !ok {
			continue
		}
		text := markdownNodeText(paragraph, document.Body)
		if !strings.HasPrefix(text, prefix) {
			continue
		}
		if value != "" {
			issues = append(issues, iddMarkdownIssue{
				Rule:    "idd-document-schema",
				Message: fmt.Sprintf("record repeats %s", label),
				Line:    markdownNodeLine(document, paragraph),
				Link:    recordIdentifier(record.heading),
				Code:    strings.ToLower(label),
			})
			continue
		}
		value = strings.TrimSpace(strings.TrimPrefix(text, prefix))
		line = markdownNodeLine(document, paragraph)
	}
	return value, line, issues
}

func iddRecordHasNarrative(record iddMarkdownRecord, source []byte) bool {
	for _, block := range record.blocks {
		if _, heading := block.(*goldmarkast.Heading); heading {
			continue
		}
		if markdownNodeText(block, source) != "" {
			return true
		}
		lines := block.Lines()
		if lines != nil && lines.Len() > 0 && strings.TrimSpace(string(lines.Value(source))) != "" {
			return true
		}
	}
	return false
}

func findIDDRecordTables(document *parsedIDDDocument, prefix string) []iddMarkdownIssue {
	var issues []iddMarkdownIssue
	inFence := false
	for index, line := range strings.Split(string(document.Body), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
			continue
		}
		if inFence || !strings.HasPrefix(trimmed, "|") || !strings.Contains(trimmed, prefix) {
			continue
		}
		refs := extractIDDRefs(trimmed)
		link := ""
		if len(refs) > 0 {
			link = refs[0]
		}
		issues = append(issues, iddMarkdownIssue{
			Rule:    "idd-document-markdown",
			Message: "IDD records must use human-readable H2 record blocks; authored registry tables are not canonical",
			Line:    document.BodyStartLine + index,
			Link:    link,
			Code:    "record-format",
		})
	}
	return issues
}

func recordIdentifier(heading string) string {
	identifier, _, _ := strings.Cut(heading, ":")
	return strings.TrimSpace(identifier)
}

func markdownNodeText(node goldmarkast.Node, source []byte) string {
	var builder strings.Builder
	_ = goldmarkast.Walk(node, func(child goldmarkast.Node, entering bool) (goldmarkast.WalkStatus, error) {
		if !entering {
			return goldmarkast.WalkContinue, nil
		}
		switch typed := child.(type) {
		case *goldmarkast.Text:
			builder.Write(typed.Value(source))
			if typed.SoftLineBreak() || typed.HardLineBreak() {
				builder.WriteByte(' ')
			}
		case *goldmarkast.String:
			builder.Write(typed.Value)
		}
		return goldmarkast.WalkContinue, nil
	})
	return strings.Join(strings.Fields(builder.String()), " ")
}

func markdownNodeLine(document *parsedIDDDocument, node goldmarkast.Node) int {
	offset := node.Pos()
	if lines := node.Lines(); lines != nil && lines.Len() > 0 {
		offset = lines.At(0).Start
	}
	if offset < 0 {
		offset = 0
	}
	if offset > len(document.Body) {
		offset = len(document.Body)
	}
	return document.BodyStartLine + bytes.Count(document.Body[:offset], []byte("\n"))
}

func trimMarkdownScalar(value string) string {
	value = strings.TrimSpace(value)
	for len(value) >= 2 && value[0] == '`' && value[len(value)-1] == '`' {
		value = strings.TrimSpace(value[1 : len(value)-1])
	}
	return value
}

func (c *DocCollector) collectIDDDocumentSet(directory string, set *model.IdentifierSet) []*model.ValidationError {
	documents := make(map[string]*parsedIDDDocument, len(iddDocumentOrder))
	var errors []*model.ValidationError

	for _, filename := range iddDocumentOrder {
		role := iddDocumentRoles[filename]
		path := filepath.Join(directory, filename)
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-set",
					fmt.Sprintf("self-describing document set is missing %s", filename),
					path,
					1,
					"",
					role,
				))
				continue
			}
			errors = append(errors, iddDocumentValidationError(
				"idd-document-parse",
				fmt.Sprintf("read %s: %v", filename, err),
				path,
				1,
				"",
				role,
			))
			continue
		}

		parsed, detected, parseErr := parseIDDDocument(path, data)
		if parseErr != nil {
			var semanticError *iddSemanticFrontmatterError
			if stderrors.As(parseErr, &semanticError) {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-migration",
					semanticError.Error(),
					path,
					semanticError.Line,
					"",
					"yaml-semantics",
				))
				continue
			}
			errors = append(errors, iddDocumentValidationError(
				"idd-document-parse",
				parseErr.Error(),
				path,
				yamlErrorLine(parseErr),
				"",
				role,
			))
			continue
		}
		if !detected {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-set",
				fmt.Sprintf("%s must contain an idd frontmatter block because another package document is self-describing", filename),
				path,
				1,
				"",
				role,
			))
			continue
		}
		documents[role] = parsed
		errors = append(errors, validateIDDDocument(parsed, role)...)
	}

	errors = append(errors, validateIDDDocumentReferences(documents)...)
	addIDDDocumentIdentifiers(documents, set)
	declared := declaredIDDIdentifiers(documents)
	for _, document := range documents {
		errors = append(errors, validateIDDDocumentMarkdown(document, declared)...)
	}
	return errors
}

func validateIDDDocument(document *parsedIDDDocument, expectedRole string) []*model.ValidationError {
	metadata := document.Metadata
	index := document.Index
	var errors []*model.ValidationError
	addError := func(rule, message, link, code string, line int) {
		errors = append(errors, iddDocumentValidationError(rule, message, document.Path, line, link, code))
	}

	if metadata.Version != iddDocumentVersion {
		addError("idd-document-identity", `version must be "1.0"`, "", "version", index.topFieldLine("version"))
	}
	expectedPackage := packageFromDocumentPath(document.Path)
	if expectedPackage == "" {
		addError(
			"idd-document-identity",
			"self-describing documents must be located at docs/<package>/<role>.md",
			"",
			"package",
			index.topFieldLine("package"),
		)
	} else if metadata.Package == "" {
		addError("idd-document-identity", "package is required", "", "package", index.topFieldLine("package"))
	} else if metadata.Package != expectedPackage {
		addError(
			"idd-document-identity",
			fmt.Sprintf("package %q does not match document path; expected %q", metadata.Package, expectedPackage),
			"",
			"package",
			index.topFieldLine("package"),
		)
	}
	if metadata.Document != expectedRole {
		addError(
			"idd-document-identity",
			fmt.Sprintf("document must be %q in %s", expectedRole, filepath.Base(document.Path)),
			"",
			"document",
			index.topFieldLine("document"),
		)
	}

	switch expectedRole {
	case "design":
		errors = append(errors, validateIDDNames(document, "components", metadata.Components)...)
	case "contract":
		errors = append(errors, validateIDDNames(document, "contracts", metadata.Contracts)...)
	case "spec":
		errors = append(errors, validateIDDSpecs(document, expectedPackage)...)
	case "testing":
		errors = append(errors, validateIDDTests(document, expectedPackage)...)
	}

	if mappingValue(document.Root, "markers") != nil || mappingValue(document.Root, "related_files") != nil {
		addError(
			"idd-document-markdown",
			"self-describing documents must not repeat legacy markers or related_files metadata",
			"",
			"duplicate-metadata",
			1,
		)
	}
	for _, issue := range document.Issues {
		errors = append(errors, iddDocumentValidationError(
			issue.Rule,
			issue.Message,
			document.Path,
			issue.Line,
			issue.Link,
			issue.Code,
		))
	}
	return errors
}

func validateIDDNames(document *parsedIDDDocument, field string, values []string) []*model.ValidationError {
	var errors []*model.ValidationError
	seen := make(map[string]bool, len(values))
	for valueIndex, value := range values {
		line := document.Index.recordLine(field, valueIndex)
		if strings.TrimSpace(value) == "" || isIDDPlaceholder(value, "", field) {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf("%s entries must contain concrete, meaningful names", field),
				document.Path,
				line,
				"",
				field,
			))
			continue
		}
		if seen[value] {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf("%s repeats %q", field, value),
				document.Path,
				line,
				"",
				field,
			))
			continue
		}
		seen[value] = true
		if !strings.Contains(string(document.Body), value) {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-reference",
				fmt.Sprintf("%s entry %q has no narrative definition in %s", field, value, filepath.Base(document.Path)),
				document.Path,
				line,
				"",
				field,
			))
		}
	}
	return errors
}

func validateIDDSpecs(document *parsedIDDDocument, packagePath string) []*model.ValidationError {
	var errors []*model.ValidationError
	module := moduleFromPackage(packagePath)
	seen := make(map[string]bool, len(document.Metadata.Specs))
	for recordIndex, spec := range document.Metadata.Specs {
		line := document.Index.recordLine("specs", recordIndex)
		validateIDDRecordField(&errors, document, "SPEC", spec.ID, "id", spec.ID, line)
		validateIDDRecordField(&errors, document, "SPEC", spec.ID, "title", spec.Title, document.Index.recordFieldLine("specs", recordIndex, "title"))
		validateIDDRecordField(&errors, document, "SPEC", spec.ID, "requirement", spec.Requirement, document.Index.recordFieldLine("specs", recordIndex, "requirement"))
		validateIDDRecordField(&errors, document, "SPEC", spec.ID, "design", spec.Design, document.Index.recordFieldLine("specs", recordIndex, "design"))
		validateIDDRecordField(&errors, document, "SPEC", spec.ID, "contract", spec.Contract, document.Index.recordFieldLine("specs", recordIndex, "contract"))

		if spec.ID == "" {
			continue
		}
		if err := pattern.ValidateIdentifierFormat(spec.ID); err != nil || !strings.HasPrefix(spec.ID, "SPEC-") {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf("SPEC record has invalid id %q; expected SPEC-<MODULE>-<NUMBER>", spec.ID),
				document.Path,
				document.Index.recordFieldLine("specs", recordIndex, "id"),
				spec.ID,
				"id",
			))
		} else if module != "" && !strings.HasPrefix(spec.ID, "SPEC-"+module+"-") {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf("SPEC id %s does not use package-derived module %s", spec.ID, module),
				document.Path,
				document.Index.recordFieldLine("specs", recordIndex, "id"),
				spec.ID,
				"id",
			))
		}
		if seen[spec.ID] {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf("duplicate SPEC id %s", spec.ID),
				document.Path,
				line,
				spec.ID,
				"id",
			))
		}
		seen[spec.ID] = true
	}
	return errors
}

func validateIDDTests(document *parsedIDDDocument, packagePath string) []*model.ValidationError {
	var errors []*model.ValidationError
	module := moduleFromPackage(packagePath)
	seen := make(map[string]bool, len(document.Metadata.Tests))
	for recordIndex, test := range document.Metadata.Tests {
		line := document.Index.recordLine("tests", recordIndex)
		validateIDDRecordField(&errors, document, "TEST", test.ID, "id", test.ID, line)
		validateIDDRecordField(&errors, document, "TEST", test.ID, "title", test.Title, document.Index.recordFieldLine("tests", recordIndex, "title"))
		validateIDDRecordField(&errors, document, "TEST", test.ID, "purpose", test.Purpose, document.Index.recordFieldLine("tests", recordIndex, "purpose"))
		validateIDDRecordField(&errors, document, "TEST", test.ID, "kind", test.Kind, document.Index.recordFieldLine("tests", recordIndex, "kind"))

		if test.ID != "" {
			if err := pattern.ValidateIdentifierFormat(test.ID); err != nil || !strings.HasPrefix(test.ID, "TEST-") {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-schema",
					fmt.Sprintf("TEST record has invalid id %q; expected TEST-<MODULE>-<NUMBER>", test.ID),
					document.Path,
					document.Index.recordFieldLine("tests", recordIndex, "id"),
					test.ID,
					"id",
				))
			} else if module != "" && !strings.HasPrefix(test.ID, "TEST-"+module+"-") {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-schema",
					fmt.Sprintf("TEST id %s does not use package-derived module %s", test.ID, module),
					document.Path,
					document.Index.recordFieldLine("tests", recordIndex, "id"),
					test.ID,
					"id",
				))
			}
			if seen[test.ID] {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-schema",
					fmt.Sprintf("duplicate TEST id %s", test.ID),
					document.Path,
					line,
					test.ID,
					"id",
				))
			}
			seen[test.ID] = true
		}

		if test.Kind != "" && test.Kind != "test" && test.Kind != "contract" {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf("TEST %s kind must be test or contract", test.ID),
				document.Path,
				document.Index.recordFieldLine("tests", recordIndex, "kind"),
				test.ID,
				"kind",
			))
		}
		if len(test.Covers) == 0 {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf("TEST %s covers must contain at least one SPEC id", test.ID),
				document.Path,
				document.Index.recordFieldLine("tests", recordIndex, "covers"),
				test.ID,
				"covers",
			))
		}
		covered := make(map[string]bool, len(test.Covers))
		for _, specID := range test.Covers {
			if covered[specID] {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-schema",
					fmt.Sprintf("TEST %s repeats covered SPEC %s", test.ID, specID),
					document.Path,
					document.Index.recordFieldLine("tests", recordIndex, "covers"),
					test.ID,
					"covers",
				))
			}
			covered[specID] = true
		}
	}
	return errors
}

func validateIDDRecordField(
	errors *[]*model.ValidationError,
	document *parsedIDDDocument,
	recordType string,
	recordID string,
	field string,
	value string,
	line int,
) {
	value = strings.TrimSpace(value)
	if value == "" {
		label := recordID
		if label == "" {
			label = recordType + " record"
		}
		*errors = append(*errors, iddDocumentValidationError(
			"idd-document-schema",
			fmt.Sprintf("%s field %s is required", label, field),
			document.Path,
			line,
			recordID,
			field,
		))
		return
	}
	if !isIDDPlaceholder(value, recordID, field) {
		return
	}
	*errors = append(*errors, iddDocumentValidationError(
		"idd-document-schema",
		fmt.Sprintf("%s field %s must contain concrete, meaningful content instead of %q", recordID, field, value),
		document.Path,
		line,
		recordID,
		field,
	))
}

func isIDDPlaceholder(value, recordID, field string) bool {
	normalized := strings.ToLower(strings.TrimSpace(value))
	return (field != "id" && recordID != "" && value == recordID) ||
		normalized == "tbd" ||
		normalized == "todo" ||
		normalized == "auto-generated" ||
		normalized == "autogenerated"
}

func validateIDDDocumentReferences(documents map[string]*parsedIDDDocument) []*model.ValidationError {
	var errors []*model.ValidationError
	designs := make(map[string]bool)
	contracts := make(map[string]bool)
	specs := make(map[string]bool)
	if document := documents["design"]; document != nil {
		designs = stringSet(document.Metadata.Components)
	}
	if document := documents["contract"]; document != nil {
		contracts = stringSet(document.Metadata.Contracts)
	}
	if document := documents["spec"]; document != nil {
		for _, spec := range document.Metadata.Specs {
			if spec.ID != "" {
				specs[spec.ID] = true
			}
		}
		for recordIndex, spec := range document.Metadata.Specs {
			if spec.Design != "" && !designs[spec.Design] {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-reference",
					fmt.Sprintf("SPEC %s references undeclared design %q", spec.ID, spec.Design),
					document.Path,
					document.Index.recordFieldLine("specs", recordIndex, "design"),
					spec.ID,
					"design",
				))
			}
			if spec.Contract != "" && !contracts[spec.Contract] {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-reference",
					fmt.Sprintf("SPEC %s references undeclared contract %q", spec.ID, spec.Contract),
					document.Path,
					document.Index.recordFieldLine("specs", recordIndex, "contract"),
					spec.ID,
					"contract",
				))
			}
		}
	}
	if document := documents["testing"]; document != nil {
		for recordIndex, test := range document.Metadata.Tests {
			for _, specID := range test.Covers {
				if specs[specID] {
					continue
				}
				errors = append(errors, iddDocumentValidationError(
					"idd-document-reference",
					fmt.Sprintf("TEST %s covers undefined SPEC %s", test.ID, specID),
					document.Path,
					document.Index.recordFieldLine("tests", recordIndex, "covers"),
					test.ID,
					"covers",
				))
			}
		}
	}
	return errors
}

func addIDDDocumentIdentifiers(documents map[string]*parsedIDDDocument, set *model.IdentifierSet) {
	specIdentifiers := make(map[string]*model.Identifier)
	if document := documents["spec"]; document != nil {
		for recordIndex, spec := range document.Metadata.Specs {
			if spec.ID == "" || !strings.HasPrefix(spec.ID, "SPEC-") {
				continue
			}
			identifier := model.NewIdentifierWithDescribe(
				spec.ID,
				model.TypeSpec,
				spec.Title,
				spec.Requirement,
				document.Path,
				document.Index.recordLine("specs", recordIndex),
			)
			identifier.RawRef = spec.ID
			set.Add(identifier)
			specIdentifiers[spec.ID] = identifier
		}
	}
	if document := documents["testing"]; document != nil {
		for recordIndex, test := range document.Metadata.Tests {
			if test.ID == "" || !strings.HasPrefix(test.ID, "TEST-") {
				continue
			}
			identifier := model.NewIdentifierWithDescribe(
				test.ID,
				model.TypeTest,
				test.Title,
				test.Purpose,
				document.Path,
				document.Index.recordLine("tests", recordIndex),
			)
			identifier.RawRef = test.ID
			identifier.Kind = test.Kind
			for _, specID := range test.Covers {
				identifier.AddLink(specID)
				if spec := specIdentifiers[specID]; spec != nil {
					spec.AddLink(test.ID)
				}
			}
			set.Add(identifier)
		}
	}
}

func declaredIDDIdentifiers(documents map[string]*parsedIDDDocument) map[string]bool {
	declared := make(map[string]bool)
	for _, document := range documents {
		for _, spec := range document.Metadata.Specs {
			if spec.ID != "" {
				declared[spec.ID] = true
			}
		}
		for _, test := range document.Metadata.Tests {
			if test.ID != "" {
				declared[test.ID] = true
			}
		}
	}
	return declared
}

func validateIDDDocumentMarkdown(document *parsedIDDDocument, declared map[string]bool) []*model.ValidationError {
	var errors []*model.ValidationError
	content := string(document.Raw)
	for _, formattingError := range ValidateMarkerFormatting(content, document.Path) {
		errors = append(errors, iddDocumentValidationError(
			"idd-document-markdown",
			formattingError,
			document.Path,
			1,
			"",
			"marker-format",
		))
	}

	inCodeBlock := false
	for lineIndex, line := range strings.Split(string(document.Body), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock || !strings.HasPrefix(trimmed, "#") {
			continue
		}
		for _, ref := range extractIDDRefs(trimmed) {
			lineNumber := document.BodyStartLine + lineIndex
			if !declared[ref] {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-reference",
					fmt.Sprintf("Markdown detail heading references %s, but this package document set does not declare it", ref),
					document.Path,
					lineNumber,
					ref,
					"heading",
				))
				continue
			}
			if err := ValidateHeadingFormat(ref, trimmed); err != nil {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-markdown",
					err.Error(),
					document.Path,
					lineNumber,
					ref,
					"heading-format",
				))
			}
		}
	}
	return errors
}

func (c *DocCollector) collectIDDPackageNarrative(path string, set *model.IdentifierSet) []*model.ValidationError {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var errors []*model.ValidationError
	for _, formattingError := range ValidateMarkerFormatting(string(content), path) {
		errors = append(errors, iddDocumentValidationError(
			"idd-document-markdown",
			formattingError,
			path,
			1,
			"",
			"marker-format",
		))
	}
	inCodeBlock := false
	for lineIndex, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock || !strings.HasPrefix(trimmed, "#") {
			continue
		}
		for _, ref := range extractIDDRefs(trimmed) {
			if set.Has(ref) {
				continue
			}
			errors = append(errors, iddDocumentValidationError(
				"idd-document-reference",
				fmt.Sprintf("Markdown detail heading references %s, but this package document set does not declare it", ref),
				path,
				lineIndex+1,
				ref,
				"heading",
			))
		}
	}
	return errors
}

func iddDocumentValidationError(rule, message, path string, line int, link, code string) *model.ValidationError {
	if line <= 0 {
		line = 1
	}
	return &model.ValidationError{
		Rule:    rule,
		Message: message,
		Source:  fmt.Sprintf("%s:%d", path, line),
		Link:    link,
		Code:    code,
	}
}

func packageFromDocumentPath(path string) string {
	parts := strings.Split(filepath.ToSlash(filepath.Clean(filepath.Dir(path))), "/")
	docsIndex := -1
	for index, part := range parts {
		if part == "docs" {
			docsIndex = index
		}
	}
	if docsIndex == -1 || docsIndex+1 >= len(parts) {
		return ""
	}
	return strings.Join(parts[docsIndex+1:], "/")
}

func moduleFromPackage(packagePath string) string {
	packagePath = filepath.ToSlash(filepath.Clean(packagePath))
	if packagePath == "." || packagePath == "" {
		return ""
	}
	parts := strings.Split(packagePath, "/")
	for index, part := range parts {
		parts[index] = strings.ToUpper(strings.ReplaceAll(part, "-", "_"))
	}
	return strings.Join(parts, "_")
}

func stringSet(values []string) map[string]bool {
	result := make(map[string]bool, len(values))
	for _, value := range values {
		if value != "" {
			result[value] = true
		}
	}
	return result
}

func yamlErrorLine(err error) int {
	message := err.Error()
	lineMarker := "line "
	index := strings.Index(message, lineMarker)
	if index == -1 {
		return 1
	}
	var line int
	if _, scanErr := fmt.Sscanf(message[index:], "line %d:", &line); scanErr != nil {
		return 1
	}
	return line
}

type iddDocumentNodeIndex struct {
	root             *yaml.Node
	lineOffset       int
	recordLines      map[string][]int
	recordFieldLines map[string][]map[string]int
}

func (i iddDocumentNodeIndex) line(node *yaml.Node) int {
	if node == nil || node.Line <= 0 {
		return 1
	}
	return node.Line + i.lineOffset
}

func (i iddDocumentNodeIndex) topFieldLine(field string) int {
	value := mappingValue(i.root, field)
	if value == nil {
		return i.line(i.root)
	}
	return i.line(value)
}

func (i iddDocumentNodeIndex) recordLine(section string, recordIndex int) int {
	if recordIndex >= 0 && recordIndex < len(i.recordLines[section]) {
		return i.recordLines[section][recordIndex]
	}
	return i.topFieldLine("document")
}

func (i iddDocumentNodeIndex) recordFieldLine(section string, recordIndex int, field string) int {
	if recordIndex >= 0 && recordIndex < len(i.recordFieldLines[section]) {
		if line := i.recordFieldLines[section][recordIndex][field]; line > 0 {
			return line
		}
	}
	return i.recordLine(section, recordIndex)
}

func (i *iddDocumentNodeIndex) addRecord(section string, line int, fields map[string]int) {
	if i.recordLines == nil {
		i.recordLines = make(map[string][]int)
	}
	if i.recordFieldLines == nil {
		i.recordFieldLines = make(map[string][]map[string]int)
	}
	i.recordLines[section] = append(i.recordLines[section], line)
	i.recordFieldLines[section] = append(i.recordFieldLines[section], fields)
}

func documentRoot(document *yaml.Node) *yaml.Node {
	if document == nil {
		return nil
	}
	if document.Kind == yaml.DocumentNode && len(document.Content) > 0 {
		return document.Content[0]
	}
	return document
}

func mappingValue(mapping *yaml.Node, key string) *yaml.Node {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return nil
	}
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Value == key {
			return mapping.Content[index+1]
		}
	}
	return nil
}

func replaceMapping(mapping *yaml.Node, key string, value *yaml.Node) {
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Value == key {
			mapping.Content[index+1] = value
			return
		}
	}
	appendMapping(mapping, scalarNode(key), value)
}

func mappingNode() *yaml.Node {
	return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
}

func scalarNode(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value}
}

func quotedScalarNode(value string) *yaml.Node {
	node := scalarNode(value)
	node.Style = yaml.DoubleQuotedStyle
	return node
}

func appendMapping(mapping, key, value *yaml.Node) {
	mapping.Content = append(mapping.Content, key, value)
}

func sortedIDDDocumentPaths(directories map[string]bool) []string {
	result := make([]string, 0, len(directories))
	for directory := range directories {
		result = append(result, directory)
	}
	sort.Strings(result)
	return result
}
