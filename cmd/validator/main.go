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

	"github.com/yourorg/idd-cli/internal/collector"
	"github.com/yourorg/idd-cli/internal/config"
	"github.com/yourorg/idd-cli/internal/engine"
	"github.com/yourorg/idd-cli/internal/reporter"
)

var (
	cfgPath   string
	outPath   string
	format    string
	verbose   bool
	noConfig  bool
	version   = "1.0.0"
)

var rootCmd = &cobra.Command{
	Use:   "idd-verify",
	Short: "IDD Link Validator - validates bidirectional linkage between IDD identifiers",
	Long: `IDD Link Validator scans documentation and source code to build a linkage graph,
then validates that all references are bidirectional (spec→test→code consistency).

Example usage:
  idd-verify run ./docs
  idd-verify run ./docs --config idd.yaml
  idd-verify run ./docs --format json
  idd-verify lint ./docs --format json -o report.json`,
	Version: version,
	SilenceUsage: true,
}

var runCmd = &cobra.Command{
	Use:   "run [path]",
	Short: "Run IDD linkage validation",
	Long: `Run IDD linkage validation on the specified path.

Example:
  idd-verify run ./docs
  idd-verify run ./docs --config idd.yaml`,
	Args: cobra.MaximumNArgs(1),
	RunE: run,
}

var lintCmd = &cobra.Command{
	Use:   "lint [path]",
	Short: "Lint IDD linkage (alias for 'run')",
	Long: `Lint IDD linkage validation. This is an alias for 'run'.

Example:
  idd-verify lint ./docs
  idd-verify lint ./docs --format json -o report.json`,
	Args: cobra.MaximumNArgs(1),
	RunE: run,
}

var skillsCmd = &cobra.Command{
	Use:   "skills",
	Short: "List IDD skills defined in this project",
	Long: `List IDD skills from the skills/ directory.

Parses frontmatter from skill markdown files and outputs skill definitions.

Example:
  idd-verify skills
  idd-verify skills --format json`,
	RunE: listSkills,
}

type SkillInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	License     string `json:"license,omitempty"`
	Compatibility string `json:"compatibility,omitempty"`
	Audience    string `json:"audience,omitempty"`
	Workflow    string `json:"workflow,omitempty"`
	Protected   bool   `json:"protected,omitempty"`
	Module      string `json:"module,omitempty"`
	Path        string `json:"path"`
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "", "Path to idd.yaml config file (default: ./idd.yaml)")
	rootCmd.PersistentFlags().StringVarP(&outPath, "output", "o", "", "Output file path (default: stdout)")
	rootCmd.PersistentFlags().StringVar(&format, "format", "json", "Output format (json, markdown)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVar(&noConfig, "no-config", false, "Disable config file loading")

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(lintCmd)
	rootCmd.AddCommand(skillsCmd)

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

	var cfg *config.Config
	var err error

	if !noConfig {
		configPaths := []string{cfgPath, "./idd.yaml", "./config/idd.yaml"}
		if cfgPath != "" {
			configPaths = []string{cfgPath}
		}

		for _, p := range configPaths {
			cfg, err = config.Load(p)
			if err == nil {
				break
			}
		}
		if cfg == nil && cfgPath != "" {
			return fmt.Errorf("failed to load config from %s: %w", cfgPath, err)
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

	if verbose {
		fmt.Printf("IDD Link Validator v%s\n", version)
		fmt.Printf("Validating: %s\n", targetPath)
		if !noConfig {
			fmt.Printf("Config: %s\n", cfgPath)
		}
		fmt.Println()
	}

	start := time.Now()

	docColl := collector.NewDocCollector(cfg)
	docSet, docErrors, err := docColl.Collect(ctx, targetPath)
	if err != nil {
		return fmt.Errorf("failed to collect doc identifiers: %w", err)
	}

	codeColl := collector.NewCodeCollector(cfg)
	codeSet, err := codeColl.Collect(ctx, targetPath)
	if err != nil {
		return fmt.Errorf("failed to collect code identifiers: %w", err)
	}

	docSet.Merge(codeSet)

	if verbose {
		fmt.Printf("Collected %d identifiers in %v\n", docSet.Count(), time.Since(start))
		fmt.Printf("  SPECs: %d\n", len(docSet.Specs))
		fmt.Printf("  TESTs: %d\n", len(docSet.Tests))
		fmt.Printf("  CONTRACTS: %d\n", len(docSet.Contracts))
		fmt.Printf("  DESIGNs: %d\n", len(docSet.Designs))
		fmt.Println()
	}

	eng := engine.New(cfg)
	eng.AddStructuralErrors(docErrors)
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
		fmt.Printf("Validation completed in %v\n", time.Since(start))
		if result.Valid {
			fmt.Println("✓ Validation passed")
		} else {
			fmt.Printf("✗ Validation failed with %d errors\n", len(result.Errors))
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
			fmt.Println("No skills found")
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
		Audience   string `yaml:"audience"`
		Workflow   string `yaml:"workflow"`
		Protected  bool   `yaml:"protected"`
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
		yaml.Unmarshal([]byte(yamlContent), &fm)
	}

	return SkillInfo{
		Name:          fm.Name,
		Description:   fm.Description,
		License:       fm.License,
		Compatibility: fm.Compatibility,
		Audience:     fm.Metadata.Audience,
		Workflow:      fm.Metadata.Workflow,
		Protected:    fm.Metadata.Protected,
	}
}

func extractModuleName(filename string) string {
	name := strings.TrimSuffix(filename, ".md")
	if strings.HasPrefix(name, "SKILL-") {
		return strings.TrimPrefix(name, "SKILL-")
	}
	return ""
}
