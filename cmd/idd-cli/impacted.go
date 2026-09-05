package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jingxu9x/idd-cli/internal/collector"
	"github.com/jingxu9x/idd-cli/internal/impact"
)

var impactBase string

var docsImpactedCmd = &cobra.Command{
	Use:   "impacted [project-root]",
	Short: "Build an IDD review queue from Git changes",
	Long: `Compare a Git base revision with the complete current worktree and map
changed declarations and annotations through the shared ID-centered trace
index. The result is a deterministic review queue of related SPEC, TEST,
Contract, and Component entities with relationship provenance.

Tracked staged and unstaged changes and untracked files are included. A warning
is emitted when implementation or test evidence changed while a related
canonical SPEC file did not. The warning requests semantic review; it does not
claim that the specification is stale or score its meaning.

Example:
  idd-cli docs impacted --base main
  idd-cli docs impacted /path/to/project --base origin/main --format llm-markdown`,
	Args: cobra.MaximumNArgs(1),
	RunE: impactedDocs,
}

func init() {
	docsCmd.AddCommand(docsImpactedCmd)
	docsImpactedCmd.Flags().StringVar(
		&impactBase,
		"base",
		"HEAD",
		"Git revision to compare with the complete current worktree",
	)
}

// impactedDocs builds one shared trace project, then projects Git changes
// through its index and already-collected AST declaration spans.
func impactedDocs(cmd *cobra.Command, args []string) error {
	projectRoot, err := resolveRunProjectRoot(args)
	if err != nil {
		return err
	}
	cfg, err := loadCLIConfig(projectRoot)
	if err != nil {
		return err
	}
	writeConfigDeprecationWarnings(cfg)
	traceProject, err := collector.BuildTraceProject(context.Background(), cfg, ".", ".")
	if err != nil {
		return fmt.Errorf("build impacted trace index: %w", err)
	}
	changes, err := impact.CollectGitChanges(context.Background(), projectRoot, impactBase)
	if err != nil {
		return err
	}
	report, err := impact.BuildReportWithAnalyses(
		changes,
		traceProject.Index(),
		traceProject.Analyses(),
	)
	if err != nil {
		return err
	}
	rendered, err := impact.Render(report, format)
	if err != nil {
		return err
	}
	if outPath != "" && outPath != "-" {
		if err := os.WriteFile(outPath, rendered, 0o644); err != nil {
			return fmt.Errorf("write impacted report: %w", err)
		}
		return nil
	}
	_, err = os.Stdout.Write(rendered)
	return err
}
