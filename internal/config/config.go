package config

// @spec SPEC-BE-003
// @contract CONTRACT-BE-001

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version    string           `yaml:"version"`
	Docs       DocsConfig       `yaml:"docs"`
	Code       CodeConfig       `yaml:"code"`
	Validation ValidationConfig `yaml:"validation"`
	Output     OutputConfig     `yaml:"output"`
}

type DocsConfig struct {
	Patterns           []string           `yaml:"patterns"`
	IdentifierPatterns IdentifierPatterns `yaml:"identifier_patterns"`
	IgnorePaths        []string           `yaml:"ignore_paths"`
}

type IdentifierPatterns struct {
	Spec     string `yaml:"spec"`
	Contract string `yaml:"contract"`
	Test     string `yaml:"test"`
	Design   string `yaml:"design"`
}

type CodeConfig struct {
	Patterns    []string `yaml:"patterns"`
	Annotations []string `yaml:"annotations"`
	IgnorePaths []string `yaml:"ignore_paths"`
}

type ValidationConfig struct {
	RequireBidirectional         bool             `yaml:"require_bidirectional"`
	AllowOrphans                 bool             `yaml:"allow_orphans"`
	RequireSpecTestCoverage      bool             `yaml:"require_spec_test_coverage"`
	RequireDocCodeCorrespondence bool             `yaml:"require_doc_code_correspondence"`
	ConsistencyCheck             ConsistencyCheck `yaml:"consistency_check"`
}

// ConsistencyCheck validates semantic consistency between doc describe and code comments.
type ConsistencyCheck struct {
	Enabled   bool    `yaml:"enabled"`
	Threshold float64 `yaml:"threshold"` // 0.0-1.0, similarity score below this triggers warning
}

type OutputConfig struct {
	File         string `yaml:"file"`
	IncludeGraph bool   `yaml:"include_graph"`
	Verbose      bool   `yaml:"verbose"`
}

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

func Default() *Config {
	return &Config{
		Version: "1.0",
		Docs: DocsConfig{
			Patterns: []string{"docs/**/*.md"},
			IdentifierPatterns: IdentifierPatterns{
				Spec:     `SPEC-[A-Z]+-[0-9]+`,
				Contract: `CONTRACT-[A-Z]+-[0-9]+`,
				Test:     `TEST-[A-Z]+-[0-9]+`,
				Design:   `DESIGN-[A-Z]+-[0-9]+`,
			},
		},
		Code: CodeConfig{
			Patterns:    []string{"**/*.go", "**/*.ts", "**/*.tsx", "**/*.js"},
			Annotations: []string{"@spec", "@contract", "@test", "@design"},
		},
		Validation: ValidationConfig{
			RequireBidirectional:         true,
			AllowOrphans:                 false,
			RequireSpecTestCoverage:      true,
			RequireDocCodeCorrespondence: true,
			ConsistencyCheck: ConsistencyCheck{
				Enabled:   true,
				Threshold: 0.3, // default threshold for similarity score
			},
		},
		Output: OutputConfig{
			IncludeGraph: false,
			Verbose:      false,
		},
	}
}

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
		c.Code.Annotations = []string{"@spec", "@contract", "@test", "@design"}
	}
	if c.Validation.ConsistencyCheck.Threshold <= 0 {
		c.Validation.ConsistencyCheck.Threshold = 0.3
	}
	if c.Validation.ConsistencyCheck.Threshold > 1.0 {
		c.Validation.ConsistencyCheck.Threshold = 1.0
	}
	return nil
}
