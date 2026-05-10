// Package config provides configuration loading and validation.

// Spec: docs/internal/config/spec.md
// Contract: docs/internal/config/contract.md

package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the root configuration structure that holds all settings for the IDD CLI validation tool.
// @implement SPEC-INTERNAL_CONFIG-001
type Config struct {
	Version    string           `yaml:"version"`
	Docs       DocsConfig       `yaml:"docs"`
	Code       CodeConfig       `yaml:"code"`
	Validation ValidationConfig `yaml:"validation"`
	Output     OutputConfig     `yaml:"output"`
}

// DocsConfig holds documentation-related configuration including patterns and ignore paths.
// @implement SPEC-INTERNAL_CONFIG-002
type DocsConfig struct {
	Patterns           []string           `yaml:"patterns"`
	IdentifierPatterns IdentifierPatterns `yaml:"identifier_patterns"`
	IgnorePaths        []string           `yaml:"ignore_paths"`
}

// IdentifierPatterns defines regex patterns for matching SPEC, TEST, and other IDD identifiers.
// @implement SPEC-INTERNAL_CONFIG-003
type IdentifierPatterns struct {
	Spec         string `yaml:"spec"`
	Test         string `yaml:"test"`
	TestContract string `yaml:"test_contract"`
}

// CodeConfig holds code-related configuration including patterns and annotations.
// @implement SPEC-INTERNAL_CONFIG-004
type CodeConfig struct {
	Patterns    []string `yaml:"patterns"`
	Annotations []string `yaml:"annotations"`
	IgnorePaths []string `yaml:"ignore_paths"`
}

// ValidationConfig holds validation rule settings for the IDD CLI.
// @implement SPEC-INTERNAL_CONFIG-005
type ValidationConfig struct {
	RequireDocLinkConsistency    bool             `yaml:"require_doc_link_consistency"`
	AllowOrphans                 bool             `yaml:"allow_orphans"`
	RequireSpecTestCoverage      bool             `yaml:"require_spec_test_coverage"`
	RequireContractTestCoverage  bool             `yaml:"require_contract_test_coverage"`
	RequireDesignSections        bool             `yaml:"require_design_sections"`
	RequireDocCodeCorrespondence bool             `yaml:"require_doc_code_correspondence"`
	RequirePublicFuncAnnotation  bool             `yaml:"require_public_func_annotation"`
	RequirePackageDocComment     bool             `yaml:"require_package_doc_comment"`
	RequireRelatedFiles          bool             `yaml:"require_related_files"`
	RequireTestAnnotation        bool             `yaml:"require_test_annotation"`
	RequireAnnotationIdentifier  bool             `yaml:"require_annotation_identifier"`
	RequireAnnotationOnSameLine  bool             `yaml:"require_annotation_on_same_line"`
	RequirePkgDocFiles           bool             `yaml:"require_pkg_doc_files"`
	ConsistencyCheck             ConsistencyCheck `yaml:"consistency_check"`
}

// ConsistencyCheck validates semantic consistency between doc describe and code comments.
// @implement SPEC-INTERNAL_CONFIG-006
type ConsistencyCheck struct {
	Enabled   bool    `yaml:"enabled"`
	Threshold float64 `yaml:"threshold"` // 0.0-1.0, similarity score below this triggers warning
}

// OutputConfig holds output-related configuration settings.
// @implement SPEC-INTERNAL_CONFIG-007
type OutputConfig struct {
	File         string `yaml:"file"`
	IncludeGraph bool   `yaml:"include_graph"`
	Verbose      bool   `yaml:"verbose"`
}

// Load reads and validates configuration from a YAML file at the given path.
// @implement SPEC-INTERNAL_CONFIG-007
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// Default returns a Config with sensible default values for the IDD CLI.
// @implement SPEC-INTERNAL_CONFIG-008
func Default() *Config {
	return &Config{
		Version: "1.0",
		Docs: DocsConfig{
			Patterns: []string{"docs/**/*.md"},
			IdentifierPatterns: IdentifierPatterns{
				Spec:         `SPEC-[A-Z]+-[0-9]+`,
				Test:         `TEST-[A-Z]+-[0-9]+`,
				TestContract: `TEST-[A-Z]+-[0-9]+`,
			},
			IgnorePaths: []string{
				"examples/**",
				"cmd/idd-cli/skills/**",
				"cmd/**",
				"docs/architecture/**",
				".planning/**",
			},
		},
		Code: CodeConfig{
			Patterns:    []string{"**/*.go"},
			Annotations: []string{"@implement", "@test", "@test-contract"},
		},
		Validation: ValidationConfig{
			AllowOrphans:                 false,
			RequireDocLinkConsistency:    true,
			RequireSpecTestCoverage:      true,
			RequireContractTestCoverage:  true,
			RequireDesignSections:        true,
			RequireDocCodeCorrespondence: true,
			RequirePublicFuncAnnotation:  true,
			RequirePackageDocComment:     true,
			RequireRelatedFiles:          true,
			RequireTestAnnotation:        true,
			RequireAnnotationIdentifier:  true,
			RequireAnnotationOnSameLine:  true,
			RequirePkgDocFiles:           true,
			ConsistencyCheck: ConsistencyCheck{
				Enabled:   true,
				Threshold: 0.3,
			},
		},
		Output: OutputConfig{
			IncludeGraph: false,
			Verbose:      false,
		},
	}
}

// Validate checks that configuration values are correct and sets defaults where appropriate.
// @implement SPEC-INTERNAL_CONFIG-009
func (c *Config) Validate() error {
	if c.Version == "" {
		c.Version = "1.0"
	}
	if len(c.Docs.Patterns) == 0 {
		c.Docs.Patterns = []string{"docs/**/*.md"}
	}
	if len(c.Code.Patterns) == 0 {
		c.Code.Patterns = []string{"**/*.go"}
	}
	if len(c.Code.Annotations) == 0 {
		c.Code.Annotations = []string{"@implement", "@test", "@test-contract"}
	}
	if c.Validation.ConsistencyCheck.Threshold <= 0 {
		c.Validation.ConsistencyCheck.Threshold = 0.3
	}
	if c.Validation.ConsistencyCheck.Threshold > 1.0 {
		c.Validation.ConsistencyCheck.Threshold = 1.0
	}
	return nil
}
