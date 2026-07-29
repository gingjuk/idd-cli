package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/jingxu9x/idd-cli/internal/collector"
	"github.com/jingxu9x/idd-cli/internal/config"
	"github.com/jingxu9x/idd-cli/internal/engine"
	"github.com/jingxu9x/idd-cli/internal/reporter"
)

var (
	cfgPath        string
	outPath        string
	format         string
	verbose        bool
	noConfig       bool
	reviewDocsPath string
	version        = "1.0.0"
)

var rootCmd = &cobra.Command{
	Use:   "idd-cli",
	Short: "IDD workflow guardrail for skills, documents, and code traceability",
	Long: `idd-cli embeds the IDD authoring skill and enforces the document model
that skill describes. The skill authors semantic content; idd-cli exports the
paired workflow, maintains safe document structure, and validates the complete
documentation/code graph.

Example usage:
  idd-cli generate skill -o idd-skill.md
  idd-cli docs init internal/auth
  idd-cli docs review-context SPEC-INTERNAL_AUTH-001
  idd-cli run . --format llm-markdown
  idd-cli run . --format json`,
	Version:      version,
	SilenceUsage: true,
}

var runCmd = &cobra.Command{
	Use:   "run [path]",
	Short: "Run IDD linkage validation",
	Long: `Run IDD linkage validation and emit the complete finding report.

Documentation is collected from [path], while source annotations are collected
from the current project working tree. Run from project root with "." for the
authoritative project gate. An invalid graph writes its report, then exits
non-zero. Passing proves enabled structural and traceability rules; the paired
IDD skill and human review still judge whether the authored explanation is
complete and correct.

Example:
  idd-cli run . --format llm-markdown
  idd-cli run . --config .idd.yaml --format json`,
	Args: cobra.MaximumNArgs(1),
	RunE: run,
}

var lintCmd = &cobra.Command{
	Use:   "lint [path]",
	Short: "Lint IDD linkage (alias for 'run')",
	Long: `Lint IDD linkage validation. This is an alias for 'run' and uses the
same project-root validation scope.

Example:
  idd-cli lint . --format llm-markdown
  idd-cli lint . --format json -o report.json`,
	Args: cobra.MaximumNArgs(1),
	RunE: run,
}

var skillsCmd = &cobra.Command{
	Use:   "skills",
	Short: "List IDD skills defined in this project",
	Long: `List IDD skills from the local skills/ directory, falling back to the
workflow embedded in this binary when no local Markdown skill is present.

Parses frontmatter from skill markdown files and outputs skill definitions.

Example:
  idd-cli skills
  idd-cli skills --format json`,
	RunE: listSkills,
}

var generateCmd = &cobra.Command{
	Use:   "generate [skill|skill --output file]",
	Short: "Export the IDD skill paired with this binary",
	Long: `Export the IDD workflow embedded in this binary.

Install the result through the agent's skill mechanism. Regenerate it after
upgrading idd-cli so authoring instructions and validator behavior stay
aligned.

Example:
  idd-cli generate skill --output idd-skill.md

Then install idd-skill.md as SKILL.md through the agent's normal skill
mechanism.`,
	RunE: generateSkill,
}

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Create and repair self-describing IDD package documents",
	Long: `Create and repair package-local self-describing IDD documents.

Each fixed Markdown file has minimal identity frontmatter and owns its
human-readable role-specific Markdown records. A SPEC update changes spec.md;
a TEST update or coverage change changes testing.md. The paired IDD skill
authors complete semantic content; these commands maintain safe structure only.
They do not shorten, summarize, or judge the adequacy of human-authored prose.`,
}

var docsInitCmd = &cobra.Command{
	Use:   "init <package>...",
	Short: "Create four self-describing IDD document skeletons",
	Long: `Create design.md, contract.md, spec.md, and testing.md under
docs/<package>/. Each file contains minimal identity frontmatter and
human-readable Markdown guidance for its role. Existing plain narrative bodies
are preserved. Existing IDD or legacy marker metadata is never overwritten.
All requested packages pass preflight before the first expected write. Use
this once for new packages, then let the paired IDD skill author complete
records with behavior, rationale, boundaries, failures, and evidence. The
generated guidance is a scaffold, not finished documentation.

Example:
  idd-cli docs init internal/auth
  idd-cli docs init internal/auth internal/config`,
	Args: cobra.MinimumNArgs(1),
	RunE: initPackageDocs,
}

var docsFixCmd = &cobra.Command{
	Use:   "fix <docs-package-or-document>...",
	Short: "Normalize safe document metadata without inventing semantics",
	Long: `Normalize version and package identity in self-describing IDD
frontmatter. The exact filename remains the sole document-role authority. A
directory target repairs the four-file set and creates missing skeletons. A
Markdown file target writes only that file. Markdown bodies are preserved
byte-for-byte. All targets are planned before writes, and overlaps are
deduplicated. Semantic records are neither rewritten nor invented. Use the
paired IDD skill to resolve semantic findings reported by run.

Example:
  idd-cli docs fix docs/internal/auth
  idd-cli docs fix docs/internal/auth/testing.md
  idd-cli docs fix docs/internal/auth docs/internal/config`,
	Args: cobra.MinimumNArgs(1),
	RunE: repairPackageDocs,
}

var docsStatusCmd = &cobra.Command{
	Use:   "status <docs-package-or-document>...",
	Short: "List IDD document slots that still require authored content",
	Long: `Inspect one or more self-describing IDD documents, package document
directories, or docs trees. The command uses the same role schema and
completion rules as docs init and run. Generated scaffold markers are
incomplete, and removing a marker does not pass unless the bounded section or
record contains effective authored content. Once a record exists, every missing
or placeholder required record field remains an incomplete slot.

JSON output contains normalized targets and a deterministic incomplete_slots
work list with file, line, role, slot, and reason. Repeated or overlapping
targets do not duplicate work items.

Example:
  idd-cli docs status docs/internal/auth --format json
  idd-cli docs status docs/internal/auth docs/internal/config --format json`,
	Args: cobra.MinimumNArgs(1),
	RunE: statusPackageDocs,
}

var docsReviewContextCmd = &cobra.Command{
	Use:   "review-context <SPEC-ID>...",
	Short: "Assemble focused evidence for reviewing one or more SPECs",
	Long: `Assemble bounded, read-only evidence bundles for one or more
canonical SPECs with one documentation scan and one source scan.

The --docs-path flag selects the documentation search root and defaults to
".". Source declarations are collected from the current project working tree,
matching the scope used by run. Repeated SPEC IDs are deduplicated in
first-request order, and one batch accepts at most ten unique IDs.

One SPEC retains JSON schema idd.spec_review_context.v1. Multiple SPECs use
idd.spec_review_context_batch.v1 with an ordered context per SPEC. Markdown
formats preserve the same independent Requirement, Acceptance, Contract,
covering TEST, implementation, and test evidence with neutral review questions.

This command does not invoke an LLM, score prose, approve documentation, or
participate in validation.

Example:
  idd-cli docs review-context SPEC-INTERNAL_AUTH-001
  idd-cli docs review-context SPEC-INTERNAL_AUTH-001 SPEC-INTERNAL_AUTH-002
  idd-cli docs review-context SPEC-INTERNAL_AUTH-001 --docs-path docs/internal/auth --format llm-markdown`,
	Args: cobra.MinimumNArgs(1),
	RunE: reviewSpecContext,
}

// @implement SPEC-CMD_IDD_CLI-009
type SkillInfo struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	License       string `json:"license,omitempty"`
	Compatibility string `json:"compatibility,omitempty"`
	Audience      string `json:"audience,omitempty"`
	Workflow      string `json:"workflow,omitempty"`
	Protected     bool   `json:"protected,omitempty"`
	Module        string `json:"module,omitempty"`
	Path          string `json:"path"`
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "", "Path to .idd.yaml config file (default: ./.idd.yaml)")
	rootCmd.PersistentFlags().StringVarP(&outPath, "output", "o", "", "Output file path (default: stdout)")
	rootCmd.PersistentFlags().StringVar(&format, "format", "json", "Output format (json, markdown, llm-markdown)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVar(&noConfig, "no-config", false, "Disable config file loading")

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(lintCmd)
	rootCmd.AddCommand(skillsCmd)
	rootCmd.AddCommand(generateCmd)
	docsCmd.AddCommand(docsInitCmd)
	docsCmd.AddCommand(docsFixCmd)
	docsCmd.AddCommand(docsStatusCmd)
	docsCmd.AddCommand(docsReviewContextCmd)
	docsReviewContextCmd.Flags().StringVar(
		&reviewDocsPath,
		"docs-path",
		".",
		"Documentation search root",
	)
	rootCmd.AddCommand(docsCmd)

	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	_ = viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format"))
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("no-config", rootCmd.PersistentFlags().Lookup("no-config"))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	targetPath := "."
	if len(args) > 0 {
		targetPath = args[0]
	}

	cfg, err := loadCLIConfig()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "idd-cli v%s\n", version)
		fmt.Fprintf(os.Stderr, "Validating: %s\n", targetPath)
		if !noConfig {
			fmt.Fprintf(os.Stderr, "Config: %s\n", cfgPath)
		}
		fmt.Fprintln(os.Stderr)
	}

	start := time.Now()

	docColl := collector.NewDocCollector(cfg)
	docSet, docErrors, err := docColl.Collect(ctx, targetPath)
	if err != nil {
		return fmt.Errorf("failed to collect doc identifiers: %w", err)
	}

	codeColl := collector.NewCodeCollector(cfg)
	codeSet, codeErrors, err := codeColl.CollectWithErrors(ctx, ".")
	if err != nil {
		return fmt.Errorf("failed to collect code identifiers: %w", err)
	}

	docSet.Merge(codeSet)

	if verbose {
		fmt.Fprintf(os.Stderr, "Collected %d identifiers in %v\n", docSet.Count(), time.Since(start))
		fmt.Fprintf(os.Stderr, "  SPECs: %d\n", len(docSet.Specs))
		fmt.Fprintf(os.Stderr, "  TESTs: %d\n", len(docSet.Tests))
		fmt.Fprintf(os.Stderr, "  CONTRACTS: %d\n", len(docSet.Contracts))
		fmt.Fprintf(os.Stderr, "  DESIGNs: %d\n", len(docSet.Designs))
		fmt.Fprintln(os.Stderr)
	}

	eng := engine.New(cfg)
	eng.SetSourceAnalyses(codeColl.Analyses())
	eng.AddStructuralErrors(append(docErrors, codeErrors...))
	result, err := eng.Run(ctx, docSet)
	if err != nil {
		return fmt.Errorf("failed to run validation: %w", err)
	}

	rep := reporter.New(cfg, format)
	report, err := rep.Generate(result)
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	if err := rep.Write(report, outPath); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Validation completed in %v\n", time.Since(start))
		if result.Valid {
			fmt.Fprintln(os.Stderr, "✓ Validation passed")
		} else {
			fmt.Fprintf(os.Stderr, "✗ Validation failed with %d errors\n", len(result.Errors))
		}
	}

	if !result.Valid {
		os.Exit(1)
	}

	return nil
}

func listSkills(cmd *cobra.Command, args []string) error {
	skillsPath := "skills"
	if len(args) > 0 {
		skillsPath = args[0]
	}

	skills, err := loadSkills(skillsPath)
	if err != nil {
		return fmt.Errorf("failed to load skills: %w", err)
	}

	if len(skills) == 0 {
		if verbose {
			fmt.Fprintln(os.Stderr, "No skills found")
		}
		return nil
	}

	if format == "json" {
		data, err := json.MarshalIndent(skills, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal skills: %w", err)
		}
		fmt.Println(string(data))
	} else {
		for _, s := range skills {
			fmt.Printf("## %s\n", s.Name)
			fmt.Printf("%s\n", s.Description)
			if s.Protected {
				fmt.Printf("🔒 Protected skill\n")
			}
			fmt.Printf("Module: %s\n", s.Module)
			fmt.Printf("Path: %s\n", s.Path)
			fmt.Println()
		}
	}

	return nil
}

func loadSkills(skillsPath string) ([]SkillInfo, error) {
	var skills []SkillInfo

	entries, err := os.ReadDir(skillsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return loadEmbeddedSkills()
		}
		return nil, err
	}

	hasMD := false
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			hasMD = true
			break
		}
	}

	if !hasMD {
		return loadEmbeddedSkills()
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		path := filepath.Join(skillsPath, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		skill := parseSkillFrontmatter(string(data))
		skill.Path = path
		skill.Module = extractModuleName(entry.Name())
		skills = append(skills, skill)
	}

	return skills, nil
}

func loadEmbeddedSkills() ([]SkillInfo, error) {
	var skills []SkillInfo

	paths := ListEmbeddedSkills()
	if len(paths) == 0 {
		return skills, nil
	}

	for _, p := range paths {
		data, err := ReadEmbeddedSkill(p)
		if err != nil {
			continue
		}
		skill := parseSkillFrontmatter(string(data))
		skill.Path = p
		skill.Module = extractModuleName(filepath.Base(p))
		skills = append(skills, skill)
	}

	return skills, nil
}

type frontmatter struct {
	Name          string `yaml:"name"`
	Description   string `yaml:"description"`
	License       string `yaml:"license"`
	Compatibility string `yaml:"compatibility"`
	Metadata      struct {
		Audience  string `yaml:"audience"`
		Workflow  string `yaml:"workflow"`
		Protected bool   `yaml:"protected"`
	} `yaml:"metadata"`
}

func parseSkillFrontmatter(content string) SkillInfo {
	var fm frontmatter

	lines := strings.Split(content, "\n")
	var yamlLines []string
	inFrontmatter := false

	for _, line := range lines {
		if strings.TrimSpace(line) == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			break
		}
		if inFrontmatter {
			yamlLines = append(yamlLines, line)
		}
	}

	if len(yamlLines) > 0 {
		yamlContent := strings.Join(yamlLines, "\n")
		_ = yaml.Unmarshal([]byte(yamlContent), &fm)
	}

	return SkillInfo{
		Name:          fm.Name,
		Description:   fm.Description,
		License:       fm.License,
		Compatibility: fm.Compatibility,
		Audience:      fm.Metadata.Audience,
		Workflow:      fm.Metadata.Workflow,
		Protected:     fm.Metadata.Protected,
	}
}

func extractModuleName(filename string) string {
	name := strings.TrimSuffix(filename, ".md")
	if strings.HasPrefix(name, "SKILL-") {
		return strings.TrimPrefix(name, "SKILL-")
	}
	return ""
}

func generateSkill(cmd *cobra.Command, args []string) error {
	if len(args) > 0 && args[0] != "skill" {
		return fmt.Errorf("unknown generate target: %s", args[0])
	}

	skillPath := "skills/SKILL.md"
	data, err := ReadEmbeddedSkill(skillPath)
	if err != nil {
		return fmt.Errorf("failed to read embedded skill: %w", err)
	}

	if outPath == "" || outPath == "-" {
		fmt.Print(string(data))
	} else {
		if err := os.WriteFile(outPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write skill file: %w", err)
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Skill written to %s\n", outPath)
		}
	}

	return nil
}

// @implement SPEC-CMD_IDD_CLI-009
func initPackageDocs(cmd *cobra.Command, args []string) error {
	changed, err := collector.InitDocumentPackages(".", args)
	if err != nil {
		return err
	}
	statusTargets := make([]string, 0, len(args))
	for _, packagePath := range args {
		statusTargets = append(
			statusTargets,
			filepath.Join("docs", filepath.FromSlash(packagePath)),
		)
	}
	status, err := collector.InspectDocumentCompletions(statusTargets)
	if err != nil {
		return err
	}
	return writeDocChanges("initialized", status.Targets, changed, status)
}

// @implement SPEC-CMD_IDD_CLI-009
func repairPackageDocs(cmd *cobra.Command, args []string) error {
	changed, err := collector.RepairDocumentTargets(args)
	if err != nil {
		return err
	}
	return writeDocChanges("repaired", cleanCommandTargets(args), changed, nil)
}

// @implement SPEC-CMD_IDD_CLI-009
func statusPackageDocs(cmd *cobra.Command, args []string) error {
	status, err := collector.InspectDocumentCompletions(args)
	if err != nil {
		return err
	}
	return writeDocumentStatus(status)
}

// @implement SPEC-CMD_IDD_CLI-007, SPEC-CMD_IDD_CLI-009
func reviewSpecContext(cmd *cobra.Command, args []string) error {
	cfg, err := loadCLIConfig()
	if err != nil {
		return err
	}
	writeConfigDeprecationWarnings(cfg)
	contextBatch, err := collector.BuildSpecReviewContexts(
		context.Background(),
		cfg,
		reviewDocsPath,
		".",
		args,
	)
	if err != nil {
		return err
	}
	rendered, err := reporter.RenderSpecReviewContexts(contextBatch, format)
	if err != nil {
		return err
	}
	if outPath != "" && outPath != "-" {
		if err := os.WriteFile(outPath, rendered, 0644); err != nil {
			return fmt.Errorf("write review context: %w", err)
		}
		return nil
	}
	_, err = os.Stdout.Write(rendered)
	return err
}

func loadCLIConfig() (*config.Config, error) {
	var cfg *config.Config
	var err error
	if !noConfig {
		configPaths := []string{cfgPath, "./.idd.yaml", "./config/.idd.yaml"}
		if cfgPath != "" {
			configPaths = []string{cfgPath}
		}
		for _, path := range configPaths {
			cfg, err = config.Load(path)
			if err == nil {
				break
			}
		}
		if cfg == nil && cfgPath != "" {
			return nil, fmt.Errorf("failed to load config from %s: %w", cfgPath, err)
		}
	}
	if cfg == nil {
		cfg = config.Default()
	}
	if verbose {
		cfg.Output.Verbose = true
	}
	if outPath != "" {
		cfg.Output.File = outPath
	}
	return cfg, nil
}

func writeConfigDeprecationWarnings(cfg *config.Config) {
	for _, warning := range cfg.DeprecationWarnings() {
		fmt.Fprintf(
			os.Stderr,
			"warning [deprecated-config] %s: %s\n",
			warning.Source,
			warning.Message,
		)
	}
}

func writeDocChanges(
	action string,
	targets []string,
	changed []string,
	status *collector.DocumentCompletionStatus,
) error {
	result := struct {
		Action          string                     `json:"action"`
		Targets         []string                   `json:"targets"`
		Changed         []string                   `json:"changed"`
		Schema          string                     `json:"schema,omitempty"`
		Status          string                     `json:"status,omitempty"`
		IncompleteSlots []collector.IncompleteSlot `json:"incomplete_slots,omitempty"`
	}{
		Action:  action,
		Targets: targets,
		Changed: changed,
	}
	if status != nil {
		result.Schema = status.Schema
		result.Status = status.Status
		result.IncompleteSlots = status.IncompleteSlots
	}

	if format == "json" {
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal document changes: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	if len(changed) == 0 {
		fmt.Println("No document changes needed.")
		return nil
	}
	for _, path := range changed {
		fmt.Println(path)
	}
	if status != nil && len(status.IncompleteSlots) > 0 {
		fmt.Printf("%d document slot(s) still require authored content.\n", len(status.IncompleteSlots))
	}
	return nil
}

func cleanCommandTargets(targets []string) []string {
	seen := make(map[string]bool, len(targets))
	cleaned := make([]string, 0, len(targets))
	for _, target := range targets {
		cleanTarget := filepath.Clean(target)
		key := cleanTarget
		if absolute, err := filepath.Abs(cleanTarget); err == nil {
			key = absolute
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		cleaned = append(cleaned, cleanTarget)
	}
	return cleaned
}

func writeDocumentStatus(status *collector.DocumentCompletionStatus) error {
	if format == "json" {
		data, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal document status: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	if status.Status == "complete" {
		fmt.Println("All required document slots are complete.")
		return nil
	}
	for _, slot := range status.IncompleteSlots {
		fmt.Printf("%s:%d [%s] %s\n", slot.File, slot.Line, slot.Slot, slot.Reason)
	}
	return nil
}
