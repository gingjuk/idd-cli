// Package reporter provides output formatting for validation results.

// Spec: docs/internal/reporter/spec.md
// Contract: docs/internal/reporter/contract.md
package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jingxu9x/idd-cli/internal/config"
	"github.com/jingxu9x/idd-cli/internal/model"
)

// Reporter generates validation reports in LLM JSON and Markdown formats.
// @implement SPEC-INTERNAL_REPORTER-001
type Reporter struct {
	cfg    *config.Config
	format string
}

type ruleInfo struct {
	Severity    string
	Title       string
	Explanation string
	FixHint     string
	Expected    string
}

var (
	identifierPattern = regexp.MustCompile(`\b(?:SPEC|TEST|CONTRACT|DESIGN)-[A-Z0-9_]+-[0-9]+\b`)

	ruleInfoByName = map[string]ruleInfo{
		"doc-link-consistency": {
			Severity:    "error",
			Title:       "Document link field is inconsistent",
			Explanation: "IDD documents must use the backlink field that matches the document type.",
			FixHint:     "Rename the incorrect relationship field to the expected field and keep the same identifier references.",
		},
		"spec-missing-tests": {
			Severity:    "error",
			Title:       "SPEC is missing test coverage",
			Explanation: "Every SPEC must be covered by at least one TEST relationship.",
			FixHint:     "In self-describing mode, add the SPEC to the TEST record's **Covers:** field in testing.md. In legacy mode, add matching **Tests:** and **Spec Coverage:** fields.",
			Expected:    "**Covers:** `SPEC-...`, or **Tests:** `TEST-...`",
		},
		"test-missing-coverage": {
			Severity:    "error",
			Title:       "TEST is missing spec coverage",
			Explanation: "Every TEST must reference at least one SPEC it covers.",
			FixHint:     "In self-describing mode, populate the testing.md TEST record's **Covers:** field. In legacy mode, add a **Spec Coverage:** field.",
			Expected:    "**Covers:** `SPEC-...`, or **Spec Coverage:** `SPEC-...`",
		},
		"orphan-detection": {
			Severity:    "error",
			Title:       "Identifier is not linked",
			Explanation: "IDD identifiers should be connected to related specs, tests, contracts, designs, or code annotations.",
			FixHint:     "Add the missing relationship field or remove the unused identifier if it is obsolete.",
		},
		"duplicate-id": {
			Severity:    "error",
			Title:       "Identifier is duplicated",
			Explanation: "The same identifier appears in multiple incompatible locations.",
			FixHint:     "Rename one identifier group or merge the duplicate declarations so each module identifier is unique.",
		},
		"annotation-format": {
			Severity:    "error",
			Title:       "Annotation format is invalid",
			Explanation: "Code annotations must include a supported annotation keyword and at least one valid IDD identifier.",
			FixHint:     "Rewrite the annotation with the expected keyword and identifier list.",
		},
		"annotation-identifier": {
			Severity:    "error",
			Title:       "Annotation identifier is missing or invalid",
			Explanation: "Annotations must reference concrete IDD identifiers.",
			FixHint:     "Add the missing identifier or replace the invalid value with an existing SPEC or TEST identifier.",
		},
		"annotation-placement": {
			Severity:    "error",
			Title:       "Annotation placement is invalid",
			Explanation: "IDD annotations should be placed on the declaration comment, not inside function bodies.",
			FixHint:     "Move the annotation to the comment immediately preceding the declaration it documents.",
		},
		"public-func-annotation": {
			Severity:    "error",
			Title:       "Public declaration is missing an implementation annotation",
			Explanation: "Public API declarations must be traceable to SPEC identifiers.",
			FixHint:     "Add an @implement annotation referencing the SPEC implemented by this declaration.",
		},
		"package-doc-comment": {
			Severity:    "error",
			Title:       "Package documentation comment is incomplete",
			Explanation: "Package comments must include the package summary and related IDD document paths.",
			FixHint:     "Update the package comment with the required Spec, Contract, Design, or Test paths.",
		},
		"design-sections": {
			Severity:    "error",
			Title:       "Design document is missing a required section",
			Explanation: "Design documents must contain the required architecture and implementation-planning sections.",
			FixHint:     "Add the missing section to design.md with concrete module details.",
		},
		"spec-required-fields": {
			Severity:    "error",
			Title:       "SPEC is missing a required field",
			Explanation: "SPEC entries must include required traceability fields such as Contract, Design, Requirement, and Tests.",
			FixHint:     "Add the missing required field to the SPEC entry.",
		},
		"idd-document-parse": {
			Severity:    "error",
			Title:       "IDD document frontmatter cannot be parsed",
			Explanation: "The document's idd block must be valid YAML containing only version, package, and document identity.",
			FixHint:     "Correct the YAML syntax or misspelled field at the reported Markdown line, then rerun validation.",
		},
		"idd-document-identity": {
			Severity:    "error",
			Title:       "IDD document identity is inconsistent",
			Explanation: "The package and document role must agree with docs/<package>/<role>.md.",
			FixHint:     "Run idd-cli docs fix on the reported document or package directory to repair structural identity values.",
		},
		"idd-document-set": {
			Severity:    "error",
			Title:       "Self-describing IDD document set is incomplete",
			Explanation: "A self-describing package has design.md, contract.md, spec.md, and testing.md, each with its own idd block.",
			FixHint:     "Run idd-cli docs fix on the package documentation directory to create missing structural skeletons.",
		},
		"idd-document-schema": {
			Severity:    "error",
			Title:       "IDD document record is incomplete",
			Explanation: "Each document role owns a small set of level-two Markdown records and required fixed fields.",
			FixHint:     "Fill the reported Markdown field with concrete package-specific content; do not use TBD or generated placeholder text.",
		},
		"idd-document-reference": {
			Severity:    "error",
			Title:       "IDD document reference is unresolved",
			Explanation: "SPEC and TEST relationships must resolve across the four self-describing package documents.",
			FixHint:     "Correct the reference or add the declaration to the owning document with meaningful content.",
		},
		"idd-document-markdown": {
			Severity:    "error",
			Title:       "Self-describing Markdown record is inconsistent",
			Explanation: "Role-owned declarations use bounded level-two Markdown records; the idd block contains identity only.",
			FixHint:     "Use the canonical record heading and fixed fields, and remove duplicate legacy markers, related_files metadata, or declaration tables.",
		},
		"idd-document-migration": {
			Severity:    "error",
			Title:       "IDD semantic catalog must be migrated",
			Explanation: "Semantic records do not belong in a package-local idd.yaml or in document YAML frontmatter.",
			FixHint:     "Move components to design.md Component sections, contracts to Contract sections, SPECs to spec.md records, and TESTs with Covers fields to testing.md.",
		},
		"idd-document-test-kind": {
			Severity:    "error",
			Title:       "Documented TEST kind and code annotation disagree",
			Explanation: "A testing.md TEST with kind test uses @test, while kind contract uses @test-contract.",
			FixHint:     "Correct the TEST kind or replace the mismatched source annotation so one identifier uses exactly one annotation kind.",
		},
	}
)

// New creates a new Reporter with the given configuration and output format.
// Defaults to JSON if format is empty. JSON uses the LLM report schema.
// @implement SPEC-INTERNAL_REPORTER-003
func New(cfg *config.Config, format string) *Reporter {
	if format == "" {
		format = "json"
	}
	return &Reporter{cfg: cfg, format: format}
}

// Generate creates a complete report from a validation result, including tool metadata,
// config summary, and the validation result.
// @implement SPEC-INTERNAL_REPORTER-004
func (r *Reporter) Generate(result *model.ValidationResult) (*model.Report, error) {
	report := &model.Report{
		Tool:      "idd-cli",
		Version:   "1.0.0",
		Timestamp: time.Now().Format(time.RFC3339),
		Config: model.ConfigSummary{
			DocPatterns:  r.cfg.Docs.Patterns,
			CodePatterns: r.cfg.Code.Patterns,
			Annotations:  r.cfg.Code.Annotations,
		},
		Result: *result,
	}
	return report, nil
}

// Write outputs the report to stdout or file based on output path.
// @implement SPEC-INTERNAL_REPORTER-005
func (r *Reporter) Write(report *model.Report, output string) error {
	var writer io.Writer
	if output == "" || output == "-" {
		writer = os.Stdout
	} else {
		file, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer func() { _ = file.Close() }()
		writer = file
	}

	switch strings.ToLower(r.format) {
	case "json":
		enc := json.NewEncoder(writer)
		enc.SetIndent("", "  ")
		return enc.Encode(r.buildLLMReport(report))
	case "markdown", "md":
		return r.writeMarkdown(report, writer)
	case "llm-markdown", "llm-md":
		return r.writeLLMMarkdown(report, writer)
	default:
		return fmt.Errorf("unsupported format: %s", r.format)
	}
}

func (r *Reporter) writeMarkdown(report *model.Report, w io.Writer) error {
	var sb strings.Builder

	sb.WriteString("# IDD Linkage Report\n\n")
	fmt.Fprintf(&sb, "**Tool:** %s v%s  \n", report.Tool, report.Version)
	fmt.Fprintf(&sb, "**Timestamp:** %s  \n", report.Timestamp)
	fmt.Fprintf(&sb, "**Status:** %s  \n\n", statusIcon(report.Result.Valid))

	if len(report.Result.Errors) > 0 {
		sb.WriteString("## Errors\n\n")
		for _, err := range report.Result.Errors {
			fmt.Fprintf(&sb, "- [%s] %s", err.Rule, err.Message)
			if err.Source != "" {
				fmt.Fprintf(&sb, " (%s)", err.Source)
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	if len(report.Result.Warnings) > 0 {
		sb.WriteString("## Warnings\n\n")
		for _, warn := range report.Result.Warnings {
			fmt.Fprintf(&sb, "- [%s] %s", warn.Rule, warn.Message)
			if warn.Source != "" {
				fmt.Fprintf(&sb, " (%s)", warn.Source)
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Stats\n\n")
	fmt.Fprintf(&sb, "- **Total Identifiers:** %d\n", report.Result.Stats.TotalIdentifiers)
	fmt.Fprintf(&sb, "- **Total Links:** %d\n", report.Result.Stats.TotalLinks)
	fmt.Fprintf(&sb, "- **SPECs:** %d\n", report.Result.Stats.SpecsAnalyzed)
	fmt.Fprintf(&sb, "- **TESTs:** %d\n", report.Result.Stats.TestsAnalyzed)
	fmt.Fprintf(&sb, "- **CONTRACTs:** %d\n", report.Result.Stats.ContractsAnalyzed)
	fmt.Fprintf(&sb, "- **DESIGNs:** %d\n\n", report.Result.Stats.DesignsAnalyzed)

	if report.Result.Graph != nil && len(report.Result.Graph.Nodes) > 0 {
		sb.WriteString("## Linkage Graph\n\n")
		sb.WriteString("### Nodes\n\n")
		for _, node := range report.Result.Graph.Nodes {
			fmt.Fprintf(&sb, "- `%s` (%s) — in: %d, out: %d\n",
				node.ID, node.Type, node.Inbound, node.Outbound)
		}

		if len(report.Result.Graph.Edges) > 0 {
			sb.WriteString("\n### Edges\n\n")
			for _, edge := range report.Result.Graph.Edges {
				verified := "✓"
				if !edge.Verified {
					verified = "✗"
				}
				fmt.Fprintf(&sb, "- %s `%s` → `%s` %s\n",
					verified, edge.From, edge.To, edge.Type)
			}
		}
	}

	_, err := w.Write([]byte(sb.String()))
	return err
}

func statusIcon(valid bool) string {
	if valid {
		return "✅ PASS"
	}
	return "❌ FAIL"
}

// buildLLMReport converts the complete report into a finding-centered LLM report.
// @implement SPEC-INTERNAL_REPORTER-011, SPEC-INTERNAL_REPORTER-012
func (r *Reporter) buildLLMReport(report *model.Report) *model.LLMReport {
	status := "fail"
	if report.Result.Valid {
		status = "pass"
	}

	llmReport := &model.LLMReport{
		Schema: "idd.llm_report.v1",
		Status: status,
		Summary: model.LLMSummary{
			Errors:   len(report.Result.Errors),
			Warnings: len(report.Result.Warnings),
			TopRules: topRules(report.Result.Errors, report.Result.Warnings),
		},
		Findings: make([]model.LLMFinding, 0, len(report.Result.Errors)+len(report.Result.Warnings)),
	}

	for _, validationErr := range report.Result.Errors {
		llmReport.Findings = append(llmReport.Findings, r.buildFinding(validationErr, "error", report.Result.Graph))
	}
	for _, validationErr := range report.Result.Warnings {
		llmReport.Findings = append(llmReport.Findings, r.buildFinding(validationErr, "warning", report.Result.Graph))
	}
	llmReport.Summary.RuleGroups = groupFindings(llmReport.Findings)

	return llmReport
}

func (r *Reporter) buildFinding(validationErr model.ValidationError, severity string, graphSnapshot *model.GraphSnapshot) model.LLMFinding {
	info := lookupRuleInfo(validationErr.Rule, severity)
	identifier := primaryIdentifier(validationErr)
	expected, actual := expectedActual(validationErr, info)

	finding := model.LLMFinding{
		Severity:           severity,
		Rule:               validationErr.Rule,
		Title:              info.Title,
		Location:           parseLocation(validationErr.Source),
		Identifier:         identifier,
		Problem:            problemText(validationErr, info),
		Expected:           expected,
		Actual:             actual,
		SuggestedFix:       info.FixHint,
		RelatedIdentifiers: relatedIdentifiers(validationErr, identifier, graphSnapshot),
	}

	if finding.Severity == "" && info.Severity != "" {
		finding.Severity = info.Severity
	}
	if finding.Severity == "" {
		finding.Severity = severity
	}
	if finding.Title == "" {
		finding.Title = humanizeRule(validationErr.Rule)
	}
	if finding.SuggestedFix == "" {
		finding.SuggestedFix = "Inspect the referenced file and update the IDD documentation or annotation so it satisfies the validation rule."
	}

	return finding
}

func lookupRuleInfo(rule, defaultSeverity string) ruleInfo {
	if info, ok := ruleInfoByName[rule]; ok {
		if info.Severity == "" {
			info.Severity = defaultSeverity
		}
		return info
	}
	return ruleInfo{
		Severity:    defaultSeverity,
		Title:       humanizeRule(rule),
		Explanation: "The validation rule reported this issue.",
		FixHint:     "Inspect the finding location and update the IDD documentation or code annotation to satisfy the rule.",
	}
}

func parseLocation(source string) model.LLMLocation {
	if source == "" {
		return model.LLMLocation{}
	}

	file := source
	line := 0
	if idx := strings.LastIndex(source, ":"); idx > 0 && idx < len(source)-1 {
		if parsedLine, err := strconv.Atoi(source[idx+1:]); err == nil {
			file = source[:idx]
			line = parsedLine
		}
	}

	return model.LLMLocation{File: file, Line: line}
}

func primaryIdentifier(validationErr model.ValidationError) string {
	if identifierPattern.MatchString(validationErr.Link) {
		return identifierPattern.FindString(validationErr.Link)
	}
	for _, candidate := range []string{validationErr.Message, validationErr.Code, validationErr.Source} {
		if match := identifierPattern.FindString(candidate); match != "" {
			return match
		}
	}
	return ""
}

func expectedActual(validationErr model.ValidationError, info ruleInfo) (string, string) {
	expected := info.Expected
	actual := validationErr.Code

	if validationErr.Rule == "doc-link-consistency" {
		switch {
		case strings.Contains(validationErr.Message, "spec.md should use **Tests:**"):
			expected = "**Tests:** `TEST-...`"
			if actual == "" {
				actual = "**Spec Coverage:**"
			}
		case strings.Contains(validationErr.Message, "testing.md should use **Spec Coverage:**"):
			expected = "**Spec Coverage:** `SPEC-...`"
			if actual == "" {
				actual = "**Tests:**"
			}
		}
	}

	return expected, actual
}

func problemText(validationErr model.ValidationError, info ruleInfo) string {
	if info.Explanation == "" {
		return validationErr.Message
	}
	return fmt.Sprintf("%s %s", validationErr.Message, info.Explanation)
}

func relatedIdentifiers(validationErr model.ValidationError, primary string, graphSnapshot *model.GraphSnapshot) []model.LLMRelatedIdentifier {
	seen := make(map[string]bool)
	related := make([]model.LLMRelatedIdentifier, 0)

	add := func(id, relation string) {
		if id == "" || id == primary || seen[id] {
			return
		}
		seen[id] = true
		related = append(related, model.LLMRelatedIdentifier{ID: id, Relation: relation})
	}

	for _, match := range identifierPattern.FindAllString(validationErr.Message+" "+validationErr.Code, -1) {
		add(match, "mentioned")
	}

	if graphSnapshot != nil && primary != "" {
		for _, edge := range graphSnapshot.Edges {
			switch {
			case edge.From == primary:
				add(edge.To, edge.Type)
			case edge.To == primary:
				add(edge.From, edge.Type)
			}
		}
	}

	return related
}

func topRules(errors, warnings []model.ValidationError) []string {
	counts := make(map[string]int)
	for _, err := range errors {
		counts[err.Rule]++
	}
	for _, warn := range warnings {
		counts[warn.Rule]++
	}

	rules := make([]string, 0, len(counts))
	for rule := range counts {
		rules = append(rules, rule)
	}
	sort.Slice(rules, func(i, j int) bool {
		if counts[rules[i]] == counts[rules[j]] {
			return rules[i] < rules[j]
		}
		return counts[rules[i]] > counts[rules[j]]
	})
	if len(rules) > 5 {
		rules = rules[:5]
	}
	return rules
}

func groupFindings(findings []model.LLMFinding) []model.LLMFindingGroup {
	type groupAccumulator struct {
		group      model.LLMFindingGroup
		fileSet    map[string]bool
		idSet      map[string]bool
		sortKey    string
		firstIndex int
	}

	groupsByKey := make(map[string]*groupAccumulator)
	for i, finding := range findings {
		key := finding.Severity + "\x00" + finding.Rule
		acc, ok := groupsByKey[key]
		if !ok {
			acc = &groupAccumulator{
				group: model.LLMFindingGroup{
					Severity:     finding.Severity,
					Rule:         finding.Rule,
					Title:        finding.Title,
					SuggestedFix: finding.SuggestedFix,
				},
				fileSet:    make(map[string]bool),
				idSet:      make(map[string]bool),
				sortKey:    key,
				firstIndex: i,
			}
			groupsByKey[key] = acc
		}

		acc.group.Count++
		acc.group.FindingIndexes = append(acc.group.FindingIndexes, i+1)
		if finding.Location.File != "" && !acc.fileSet[finding.Location.File] {
			acc.fileSet[finding.Location.File] = true
			acc.group.Files = append(acc.group.Files, finding.Location.File)
		}
		if finding.Identifier != "" && !acc.idSet[finding.Identifier] {
			acc.idSet[finding.Identifier] = true
			acc.group.Identifiers = append(acc.group.Identifiers, finding.Identifier)
		}
		for _, related := range finding.RelatedIdentifiers {
			if related.ID != "" && !acc.idSet[related.ID] {
				acc.idSet[related.ID] = true
				acc.group.Identifiers = append(acc.group.Identifiers, related.ID)
			}
		}
	}

	accumulators := make([]*groupAccumulator, 0, len(groupsByKey))
	for _, acc := range groupsByKey {
		sort.Strings(acc.group.Files)
		sort.Strings(acc.group.Identifiers)
		accumulators = append(accumulators, acc)
	}
	sort.Slice(accumulators, func(i, j int) bool {
		if accumulators[i].group.Count != accumulators[j].group.Count {
			return accumulators[i].group.Count > accumulators[j].group.Count
		}
		if accumulators[i].group.Severity != accumulators[j].group.Severity {
			return accumulators[i].group.Severity < accumulators[j].group.Severity
		}
		if accumulators[i].group.Rule != accumulators[j].group.Rule {
			return accumulators[i].group.Rule < accumulators[j].group.Rule
		}
		return accumulators[i].firstIndex < accumulators[j].firstIndex
	})

	groups := make([]model.LLMFindingGroup, 0, len(accumulators))
	for _, acc := range accumulators {
		groups = append(groups, acc.group)
	}
	return groups
}

func humanizeRule(rule string) string {
	if rule == "" {
		return "Validation issue"
	}
	parts := strings.Fields(strings.ReplaceAll(rule, "-", " "))
	for i := range parts {
		parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
	}
	return strings.Join(parts, " ")
}

func (r *Reporter) writeLLMMarkdown(report *model.Report, w io.Writer) error {
	llmReport := r.buildLLMReport(report)
	var sb strings.Builder

	sb.WriteString("# IDD Validation Report\n\n")
	fmt.Fprintf(&sb, "Status: %s\n\n", strings.ToUpper(llmReport.Status))

	sb.WriteString("## Summary\n\n")
	fmt.Fprintf(&sb, "- Errors: %d\n", llmReport.Summary.Errors)
	fmt.Fprintf(&sb, "- Warnings: %d\n", llmReport.Summary.Warnings)
	if len(llmReport.Summary.TopRules) > 0 {
		fmt.Fprintf(&sb, "- Top Rules: %s\n", strings.Join(llmReport.Summary.TopRules, ", "))
	}
	sb.WriteString("\n")

	if len(llmReport.Findings) == 0 {
		sb.WriteString("## Highest priority fixes\n\nNo findings.\n")
		_, err := w.Write([]byte(sb.String()))
		return err
	}

	if len(llmReport.Summary.RuleGroups) > 0 {
		sb.WriteString("## Finding groups\n\n")
		for _, group := range llmReport.Summary.RuleGroups {
			fmt.Fprintf(&sb, "- %s %s: %d finding(s)", strings.ToUpper(group.Severity), group.Rule, group.Count)
			if len(group.Files) > 0 {
				fmt.Fprintf(&sb, " across %d file(s)", len(group.Files))
			}
			fmt.Fprintf(&sb, "; findings: %s\n", joinIndexes(group.FindingIndexes))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Highest priority fixes\n\n")
	for i, finding := range llmReport.Findings {
		location := finding.Location.File
		if location == "" {
			location = "unknown"
		}
		if finding.Location.Line > 0 {
			location = fmt.Sprintf("%s:%d", location, finding.Location.Line)
		}

		fmt.Fprintf(&sb, "### %d. %s\n", i+1, location)
		fmt.Fprintf(&sb, "Rule: %s\n", finding.Rule)
		if finding.Identifier != "" {
			fmt.Fprintf(&sb, "Identifier: %s\n", finding.Identifier)
		}
		fmt.Fprintf(&sb, "Severity: %s\n\n", strings.ToUpper(finding.Severity))
		sb.WriteString("Problem:\n")
		fmt.Fprintf(&sb, "%s\n\n", finding.Problem)
		if finding.Expected != "" {
			sb.WriteString("Expected:\n")
			fmt.Fprintf(&sb, "%s\n\n", finding.Expected)
		}
		if finding.Actual != "" {
			sb.WriteString("Actual:\n")
			fmt.Fprintf(&sb, "%s\n\n", finding.Actual)
		}
		sb.WriteString("Fix:\n")
		fmt.Fprintf(&sb, "%s\n", finding.SuggestedFix)
		if len(finding.RelatedIdentifiers) > 0 {
			sb.WriteString("\nRelated:\n")
			for _, related := range finding.RelatedIdentifiers {
				fmt.Fprintf(&sb, "- %s (%s)\n", related.ID, related.Relation)
			}
		}
		sb.WriteString("\n")
	}

	_, err := w.Write([]byte(sb.String()))
	return err
}

func joinIndexes(indexes []int) string {
	parts := make([]string, 0, len(indexes))
	for _, index := range indexes {
		parts = append(parts, strconv.Itoa(index))
	}
	return strings.Join(parts, ", ")
}
