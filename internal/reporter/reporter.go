package reporter

// @spec SPEC-BE-006

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/jingxu9x/idd-link-validator/internal/config"
	"github.com/jingxu9x/idd-link-validator/internal/model"
)

type Reporter struct {
	cfg    *config.Config
	format string
}

func New(cfg *config.Config, format string) *Reporter {
	if format == "" {
		format = "json"
	}
	return &Reporter{cfg: cfg, format: format}
}

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
		return enc.Encode(report)
	case "markdown", "md":
		return r.writeMarkdown(report, writer)
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
