package reporter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jingxu9x/idd-cli/internal/collector"
)

// RenderSpecReviewContext renders an evidence-only single-SPEC bundle.
// @implement SPEC-INTERNAL_REPORTER-013
func RenderSpecReviewContext(context *collector.SpecReviewContext, format string) ([]byte, error) {
	switch format {
	case "json":
		data, err := json.MarshalIndent(context, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal SPEC review context: %w", err)
		}
		return append(data, '\n'), nil
	case "markdown", "llm-markdown":
		return []byte(renderSpecReviewContextMarkdown(context, 0)), nil
	default:
		return nil, fmt.Errorf(
			"unsupported review-context format %q; expected json, markdown, or llm-markdown",
			format,
		)
	}
}

// RenderSpecReviewContexts renders one context with its stable single schema or
// several ordered contexts with the batch schema.
// @implement SPEC-INTERNAL_REPORTER-013
func RenderSpecReviewContexts(
	batch *collector.SpecReviewContextBatch,
	format string,
) ([]byte, error) {
	if batch == nil || len(batch.Contexts) == 0 {
		return nil, fmt.Errorf("SPEC review context batch is empty")
	}
	if len(batch.Contexts) == 1 {
		return RenderSpecReviewContext(batch.Contexts[0], format)
	}
	switch format {
	case "json":
		data, err := json.MarshalIndent(batch, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal SPEC review context batch: %w", err)
		}
		return append(data, '\n'), nil
	case "markdown", "llm-markdown":
		var output strings.Builder
		output.WriteString("# SPEC review context batch\n\n")
		fmt.Fprintf(&output, "Schema: `%s`\n\n", batch.Schema)
		fmt.Fprintf(&output, "Requested SPECs: `%s`\n\n", strings.Join(batch.RequestedSpecIDs, "`, `"))
		output.WriteString("Each section is independent evidence for human or LLM review. The batch contains no semantic score or approval.\n\n")
		for index, contextValue := range batch.Contexts {
			if index > 0 {
				output.WriteString("---\n\n")
			}
			output.WriteString(renderSpecReviewContextMarkdown(contextValue, 1))
		}
		return []byte(output.String()), nil
	default:
		return nil, fmt.Errorf(
			"unsupported review-context format %q; expected json, markdown, or llm-markdown",
			format,
		)
	}
}

func renderSpecReviewContextMarkdown(context *collector.SpecReviewContext, headingOffset int) string {
	var output strings.Builder
	writeMarkdownHeading(&output, 1+headingOffset, "Review context: "+context.Spec.ID)
	fmt.Fprintf(&output, "Schema: `%s`\n\n", context.Schema)
	output.WriteString("This is evidence for human or LLM review. It contains no semantic score or approval.\n\n")

	writeMarkdownHeading(&output, 2+headingOffset, "SPEC")
	fmt.Fprintf(&output, "- Title: %s\n", context.Spec.Title)
	fmt.Fprintf(&output, "- Package: `%s`\n", context.Spec.Package)
	fmt.Fprintf(&output, "- Location: `%s:%d`\n", context.Spec.File, context.Spec.Line)
	fmt.Fprintf(&output, "- Components: `%s`\n", strings.Join(context.Spec.Components, "`, `"))
	if len(context.Spec.Contracts) > 0 {
		fmt.Fprintf(&output, "- Contracts: `%s`\n", strings.Join(context.Spec.Contracts, "`, `"))
	}
	output.WriteString("\n")
	writeMarkdownHeading(&output, 3+headingOffset, "Requirement")
	fmt.Fprintf(&output, "%s\n\n", context.Spec.Requirement)
	writeMarkdownHeading(&output, 3+headingOffset, "Acceptance")
	fmt.Fprintf(&output, "%s\n\n", context.Spec.Acceptance)
	writeMarkdownHeading(&output, 3+headingOffset, "Complete authored SPEC record")
	writeFencedBlock(&output, "markdown", context.Spec.Markdown)
	fmt.Fprintf(&output, "Record truncated: `%t`\n\n", context.Spec.Truncated)

	writeMarkdownHeading(&output, 2+headingOffset, "Components and design context")
	if len(context.Components) == 0 {
		output.WriteString("No matching Component record was collected.\n\n")
	}
	for _, component := range context.Components {
		writeMarkdownHeading(&output, 3+headingOffset, component.Name)
		fmt.Fprintf(&output, "- ID: `%s`\n- Location: `%s:%d`\n\n", component.ID, component.File, component.Line)
		for _, field := range []struct{ label, value string }{{"Purpose", component.Purpose}, {"Ownership", component.Ownership}, {"Boundary", component.Boundary}, {"Decisions", component.Decisions}} {
			if field.value != "" {
				fmt.Fprintf(&output, "%s: %s\n\n", field.label, field.value)
			}
		}
		writeFencedBlock(&output, "markdown", component.Markdown)
		fmt.Fprintf(&output, "Record truncated: `%t`\n\n", component.Truncated)
	}

	writeMarkdownHeading(&output, 2+headingOffset, "Contracts")
	if len(context.Contracts) == 0 {
		output.WriteString("No matching Contract record was collected.\n\n")
	}
	for _, contract := range context.Contracts {
		writeMarkdownHeading(&output, 3+headingOffset, contract.Name)
		fmt.Fprintf(&output, "- ID: `%s`\n- Location: `%s:%d`\n\n", contract.ID, contract.File, contract.Line)
		fmt.Fprintf(&output, "%s\n\n", contract.Guarantees)
		writeFencedBlock(&output, "markdown", contract.Markdown)
		fmt.Fprintf(&output, "Record truncated: `%t`\n\n", contract.Truncated)
	}

	writeMarkdownHeading(&output, 2+headingOffset, "Covering TEST records")
	if len(context.Tests) == 0 {
		output.WriteString("No covering TEST record was collected.\n\n")
	}
	for _, test := range context.Tests {
		writeMarkdownHeading(&output, 3+headingOffset, test.ID+": "+test.Title)
		fmt.Fprintf(&output, "- Kind: `%s`\n", test.Kind)
		fmt.Fprintf(&output, "- Location: `%s:%d`\n", test.File, test.Line)
		if len(test.Contracts) > 0 {
			fmt.Fprintf(&output, "- Contracts: `%s`\n", strings.Join(test.Contracts, "`, `"))
		}
		fmt.Fprintf(&output, "\nPurpose: %s\n\nOracle: %s\n\n", test.Purpose, test.Oracle)
		writeFencedBlock(&output, "markdown", test.Markdown)
		fmt.Fprintf(&output, "Record truncated: `%t`\n\n", test.Truncated)
	}

	writeMarkdownHeading(&output, 2+headingOffset, "Annotated declarations")
	if len(context.Declarations) == 0 {
		output.WriteString("No matching source declaration was collected.\n\n")
	}
	for _, declaration := range context.Declarations {
		writeMarkdownHeading(&output, 3+headingOffset, declaration.Kind+" "+declaration.Name)
		fmt.Fprintf(
			&output,
			"- Role: `%s`\n- References: `%s`\n- Language: `%s`\n- Location: `%s:%d-%d`\n- Annotation: `%s`\n- Excerpt truncated: `%t`\n\n",
			declaration.Role,
			strings.Join(declaration.References, "`, `"),
			declaration.Language,
			declaration.File,
			declaration.Line,
			declaration.EndLine,
			strings.TrimSpace(declaration.Annotation),
			declaration.Truncated,
		)
		writeFencedBlock(&output, markdownLanguage(declaration.Language), declaration.Excerpt)
	}

	if len(context.Issues) > 0 {
		writeMarkdownHeading(&output, 2+headingOffset, "Collection issues")
		for _, issue := range context.Issues {
			fmt.Fprintf(&output, "- `%s`", issue.Rule)
			if issue.Source != "" {
				fmt.Fprintf(&output, " at `%s`", issue.Source)
			}
			fmt.Fprintf(&output, ": %s\n", issue.Message)
		}
		output.WriteString("\n")
	}
	if context.Truncated {
		output.WriteString("Some record or declaration evidence was truncated by the context budget.\n\n")
	}

	writeMarkdownHeading(&output, 2+headingOffset, "Review questions")
	output.WriteString("- Does the Requirement state one observable behavior and its boundary?\n")
	output.WriteString("- Does Acceptance provide concrete evidence for that behavior?\n")
	output.WriteString("- Do the Component boundaries, Contract guarantees, TEST oracles, and implementation agree with the SPEC?\n")
	output.WriteString("- Are important failures, side effects, or edge cases absent from the authored records?\n")
	return output.String()
}

func writeMarkdownHeading(output *strings.Builder, level int, title string) {
	fmt.Fprintf(output, "%s %s\n\n", strings.Repeat("#", level), title)
}

func markdownLanguage(language string) string {
	switch language {
	case "typescript":
		return "ts"
	case "javascript":
		return "js"
	case "cpp":
		return "cpp"
	default:
		return language
	}
}

func writeFencedBlock(output *strings.Builder, language string, content string) {
	fenceLength := 3
	for _, run := range strings.FieldsFunc(content, func(r rune) bool { return r != '`' }) {
		if len(run) >= fenceLength {
			fenceLength = len(run) + 1
		}
	}
	fence := strings.Repeat("`", fenceLength)
	fmt.Fprintf(output, "%s%s\n%s\n%s\n\n", fence, language, content, fence)
}
