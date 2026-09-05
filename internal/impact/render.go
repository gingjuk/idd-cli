package impact

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Render emits the stable machine schema or its evidence-only Markdown view.
// @implement SPEC-CMD_IDD_CLI-012
func Render(report *Report, format string) ([]byte, error) {
	if report == nil {
		return nil, fmt.Errorf("impacted report is nil")
	}
	switch strings.ToLower(format) {
	case "json":
		data, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal impacted report: %w", err)
		}
		return append(data, '\n'), nil
	case "markdown", "llm-markdown":
		return []byte(renderMarkdown(report)), nil
	default:
		return nil, fmt.Errorf("unsupported impacted format %q; expected json, markdown, or llm-markdown", format)
	}
}

func renderMarkdown(report *Report) string {
	var output strings.Builder
	output.WriteString("# IDD impacted review queue\n\n")
	fmt.Fprintf(&output, "Schema: `%s`  \nBase: `%s`  \nTraversal depth: `%d`\n\n", report.Schema, report.Base, report.Depth)
	output.WriteString("This queue is derived from Git changes and typed trace-index relationships. Warnings request semantic review; they do not claim that documentation is stale.\n\n")
	output.WriteString("## Changed files\n\n")
	if len(report.Changes) == 0 {
		output.WriteString("No changed files.\n\n")
	}
	for _, change := range report.Changes {
		fmt.Fprintf(&output, "- `%s` (%s)", change.Path, change.Status)
		if change.OldPath != "" {
			fmt.Fprintf(&output, ", from `%s`", change.OldPath)
		}
		if len(change.AnnotationIDs) > 0 {
			fmt.Fprintf(&output, "; changed annotations: `%s`", strings.Join(change.AnnotationIDs, "`, `"))
		}
		output.WriteString("\n")
	}
	output.WriteString("\n## Review queue\n\n")
	if len(report.Queue) == 0 {
		output.WriteString("No indexed IDD entities were impacted.\n\n")
	}
	for _, item := range report.Queue {
		fmt.Fprintf(&output, "### %s\n\n", item.ID)
		fmt.Fprintf(&output, "- Kind: `%s`\n- Lifecycle: `%s`\n- Resolution: `%s`\n", item.Kind, item.Lifecycle, item.Resolution)
		if item.Canonical != nil {
			fmt.Fprintf(&output, "- Canonical: `%s:%d`\n", item.Canonical.Path, item.Canonical.Line)
		}
		output.WriteString("- Reasons:\n")
		for _, reason := range item.Reasons {
			fmt.Fprintf(&output, "  - `%s`", reason.Kind)
			if reason.SourceID != "" {
				fmt.Fprintf(&output, " from `%s`", reason.SourceID)
			}
			if reason.Relation != "" {
				fmt.Fprintf(&output, " via `%s`", reason.Relation)
			}
			if reason.Path != "" {
				fmt.Fprintf(&output, " at `%s:%d`", reason.Path, reason.Line)
			}
			if reason.Conservative {
				output.WriteString(" (conservative changed-file match)")
			}
			output.WriteString("\n")
		}
		output.WriteString("\n")
	}
	output.WriteString("## Warnings\n\n")
	if len(report.Warnings) == 0 {
		output.WriteString("No review warnings.\n")
	}
	for _, warning := range report.Warnings {
		fmt.Fprintf(&output, "- `%s` for `%s`: %s (changes: `%s`)\n", warning.Rule, warning.Entity, warning.Message, strings.Join(warning.Paths, "`, `"))
	}
	return output.String()
}
