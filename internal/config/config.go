// Package config provides configuration loading and validation.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var defaultCodePatterns = []string{
	"**/*.go",
	"**/*.ts",
	"**/*.tsx",
	"**/*.js",
	"**/*.jsx",
	"**/*.cpp",
	"**/*.cc",
	"**/*.cxx",
	"**/*.hpp",
	"**/*.hh",
	"**/*.hxx",
	"**/*.java",
	"**/*.py",
}

func copyDefaultCodePatterns() []string {
	return append([]string(nil), defaultCodePatterns...)
}

// Config is the root configuration structure that holds all settings for the IDD CLI validation tool.
// @implement SPEC-INTERNAL_CONFIG-001
type Config struct {
	Version             string           `yaml:"version"`
	Docs                DocsConfig       `yaml:"docs"`
	Code                CodeConfig       `yaml:"code"`
	Validation          ValidationConfig `yaml:"validation"`
	Output              OutputConfig     `yaml:"output"`
	workdir             string
	deprecationWarnings []DeprecationWarning
}

// SetWorkdir installs the absolute project root used to resolve runtime paths.
// It is process-local state and is never decoded from or written to YAML.
// @implement SPEC-INTERNAL_CONFIG-010
func (c *Config) SetWorkdir(root string) error {
	if root == "" {
		c.workdir = ""
		return nil
	}
	if !filepath.IsAbs(root) {
		return fmt.Errorf("workdir must be absolute: %s", root)
	}
	c.workdir = filepath.Clean(root)
	for index := range c.deprecationWarnings {
		c.deprecationWarnings[index].Source = c.DisplayPath(c.deprecationWarnings[index].Source)
	}
	return nil
}

// ResolvePath resolves a project-relative runtime path without changing the
// process working directory.
// @implement SPEC-INTERNAL_CONFIG-010
func (c *Config) ResolvePath(path string) string {
	if c == nil || c.workdir == "" || filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Clean(filepath.Join(c.workdir, path))
}

// DisplayPath converts a runtime path inside the configured project root back
// to a stable project-relative path for identifiers and findings.
// @implement SPEC-INTERNAL_CONFIG-010
func (c *Config) DisplayPath(path string) string {
	clean := filepath.Clean(path)
	if c == nil || c.workdir == "" {
		return clean
	}
	absolute := clean
	if !filepath.IsAbs(absolute) {
		absolute = filepath.Join(c.workdir, absolute)
	}
	relative, err := filepath.Rel(c.workdir, absolute)
	if err != nil ||
		relative == ".." ||
		strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return clean
	}
	return filepath.Clean(relative)
}

// DeprecationWarning describes one explicitly configured obsolete YAML path.
// @implement SPEC-INTERNAL_CONFIG-006
type DeprecationWarning struct {
	Path         string
	Source       string
	Message      string
	SuggestedFix string
}

// DocsConfig holds documentation-related configuration including patterns and ignore paths.
// @implement SPEC-INTERNAL_CONFIG-002
type DocsConfig struct {
	Patterns    []string            `yaml:"patterns"`
	Units       []DocumentationUnit `yaml:"units"`
	IgnorePaths []string            `yaml:"ignore_paths"`
}

// DocumentationUnit maps one stable IDD documentation package to one or more
// source trees. Sources are project-relative glob patterns and may span
// several physical language packages.
// @implement SPEC-INTERNAL_CONFIG-002
type DocumentationUnit struct {
	Package string   `yaml:"package"`
	Sources []string `yaml:"sources"`
}

// DocumentationUnitsForSource returns all explicitly configured units whose
// source patterns match a project-relative source path, in configuration order.
func (c *Config) DocumentationUnitsForSource(path string) []DocumentationUnit {
	if c == nil {
		return nil
	}
	path = filepath.ToSlash(filepath.Clean(c.DisplayPath(path)))
	var matches []DocumentationUnit
	for _, unit := range c.Docs.Units {
		for _, source := range unit.Sources {
			if matchProjectGlob(source, path) {
				copyUnit := DocumentationUnit{Package: unit.Package, Sources: append([]string(nil), unit.Sources...)}
				matches = append(matches, copyUnit)
				break
			}
		}
	}
	return matches
}

func matchProjectGlob(pattern, path string) bool {
	patternParts := strings.Split(filepath.ToSlash(filepath.Clean(pattern)), "/")
	pathParts := strings.Split(filepath.ToSlash(filepath.Clean(path)), "/")
	var match func(int, int) bool
	match = func(patternIndex, pathIndex int) bool {
		if patternIndex == len(patternParts) {
			return pathIndex == len(pathParts)
		}
		if patternParts[patternIndex] == "**" {
			for next := pathIndex; next <= len(pathParts); next++ {
				if match(patternIndex+1, next) {
					return true
				}
			}
			return false
		}
		if pathIndex == len(pathParts) {
			return false
		}
		matched, err := filepath.Match(patternParts[patternIndex], pathParts[pathIndex])
		return err == nil && matched && match(patternIndex+1, pathIndex+1)
	}
	return match(0, 0)
}

// CodeConfig holds code-related configuration including patterns and annotations.
// @implement SPEC-INTERNAL_CONFIG-004
type CodeConfig struct {
	Patterns    []string          `yaml:"patterns"`
	Annotations map[string]string `yaml:"annotations"`
	IgnorePaths []string          `yaml:"ignore_paths"`
}

// ValidationConfig holds validation rule settings for the IDD CLI.
// @implement SPEC-INTERNAL_CONFIG-005
type ValidationConfig struct {
	RequireDocLinkConsistency    bool `yaml:"require_doc_link_consistency"`
	AllowOrphans                 bool `yaml:"allow_orphans"`
	RequireSpecFields            bool `yaml:"require_spec_fields"`
	RequireSpecTestCoverage      bool `yaml:"require_spec_test_coverage"`
	RequireContractTestCoverage  bool `yaml:"require_contract_test_coverage"`
	RequireDesignSections        bool `yaml:"require_design_sections"`
	RequireDocCodeCorrespondence bool `yaml:"require_doc_code_correspondence"`
	RequirePublicFuncAnnotation  bool `yaml:"require_public_func_annotation"`
	RequireRelatedFiles          bool `yaml:"require_related_files"`
	RequireTestAnnotation        bool `yaml:"require_test_annotation"`
	RequireAnnotationIdentifier  bool `yaml:"require_annotation_identifier"`
	RequireAnnotationOnSameLine  bool `yaml:"require_annotation_on_same_line"`
	RequirePkgDocFiles           bool `yaml:"require_pkg_doc_files"`
	// Deprecated: retained only so existing YAML remains decodable.
	ConsistencyCheck ConsistencyCheck `yaml:"consistency_check"`
}

// ConsistencyCheck preserves the deprecated validation.consistency_check YAML
// shape. Its fields are ignored by validation and review-context.
//
// Deprecated: use `idd-cli docs review-context <SPEC-ID>` for evidence-based
// human or LLM review.
// @implement SPEC-INTERNAL_CONFIG-006
type ConsistencyCheck struct {
	Enabled   bool    `yaml:"enabled"`
	Threshold float64 `yaml:"threshold"`
}

// OutputConfig holds output-related configuration settings.
// @implement SPEC-INTERNAL_CONFIG-001
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

	var document yaml.Node
	if err := yaml.Unmarshal(data, &document); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	var cfg Config
	if err := document.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	cfg.deprecationWarnings = collectDeprecationWarnings(path, &document)

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// DeprecationWarnings returns a copy of load-time obsolete-key diagnostics.
// @implement SPEC-INTERNAL_CONFIG-006
func (c *Config) DeprecationWarnings() []DeprecationWarning {
	if c == nil {
		return nil
	}
	return append([]DeprecationWarning(nil), c.deprecationWarnings...)
}

func collectDeprecationWarnings(source string, document *yaml.Node) []DeprecationWarning {
	var warnings []DeprecationWarning
	if yamlPathExists(document, "validation", "consistency_check") {
		warnings = append(warnings, DeprecationWarning{
			Path:   "validation.consistency_check",
			Source: source,
			Message: "validation.consistency_check is deprecated and ignored; " +
				"remove it from the configuration. Use idd-cli docs review-context " +
				"<SPEC-ID>... when semantic review is needed",
			SuggestedFix: "Remove the complete validation.consistency_check mapping. " +
				"Use idd-cli docs review-context <SPEC-ID>... when semantic review is needed.",
		})
	}
	if yamlPathExists(document, "docs", "identifier_patterns") {
		warnings = append(warnings, DeprecationWarning{
			Path:         "docs.identifier_patterns",
			Source:       source,
			Message:      "docs.identifier_patterns is deprecated and ignored; identifier syntax is owned by the versioned IDD protocol",
			SuggestedFix: "Remove the complete docs.identifier_patterns mapping.",
		})
	}
	return warnings
}

func yamlPathExists(node *yaml.Node, path ...string) bool {
	if node == nil || len(path) == 0 {
		return false
	}
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 0 {
			return false
		}
		node = node.Content[0]
	}
	for _, segment := range path {
		if node.Kind != yaml.MappingNode {
			return false
		}
		var next *yaml.Node
		for index := 0; index+1 < len(node.Content); index += 2 {
			if node.Content[index].Value == segment {
				next = node.Content[index+1]
				break
			}
		}
		if next == nil {
			return false
		}
		node = next
	}
	return true
}

// Default returns a Config with sensible default values for the IDD CLI.
// @implement SPEC-INTERNAL_CONFIG-008
func Default() *Config {
	return &Config{
		Version: "1.0",
		Docs: DocsConfig{
			Patterns: []string{"docs/**/*.md"},
			IgnorePaths: []string{
				"examples/**",
				"cmd/idd-cli/skills/**",
				"cmd/**",
				"docs/architecture/**",
				".planning/**",
			},
		},
		Code: CodeConfig{
			Patterns:    copyDefaultCodePatterns(),
			IgnorePaths: []string{},
			Annotations: map[string]string{
				"spec":          "@implement",
				"test":          "@test",
				"test_contract": "@test-contract",
			},
		},
		Validation: ValidationConfig{
			AllowOrphans:                 false,
			RequireDocLinkConsistency:    true,
			RequireSpecFields:            true,
			RequireSpecTestCoverage:      true,
			RequireContractTestCoverage:  true,
			RequireDesignSections:        true,
			RequireDocCodeCorrespondence: true,
			RequirePublicFuncAnnotation:  false,
			RequireRelatedFiles:          true,
			RequireTestAnnotation:        false,
			RequireAnnotationIdentifier:  true,
			RequireAnnotationOnSameLine:  true,
			RequirePkgDocFiles:           true,
		},
		Output: OutputConfig{
			IncludeGraph: false,
			Verbose:      false,
		},
	}
}

// expectedAnnotationKeys returns the canonical protocol annotation roles.
func expectedAnnotationKeys() map[string]bool {
	return map[string]bool{"spec": true, "test": true, "test_contract": true}
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
		c.Code.Patterns = copyDefaultCodePatterns()
	}
	if len(c.Code.Annotations) == 0 {
		c.Code.Annotations = map[string]string{
			"spec":          "@implement",
			"test":          "@test",
			"test_contract": "@test-contract",
		}
	}
	if err := c.validateAnnotationKeys(); err != nil {
		return err
	}
	if err := c.validateAnnotationValues(); err != nil {
		return err
	}
	if err := c.validateDocumentationUnits(); err != nil {
		return err
	}
	for label, patterns := range map[string][]string{
		"docs.patterns":     c.Docs.Patterns,
		"docs.ignore_paths": c.Docs.IgnorePaths,
		"code.patterns":     c.Code.Patterns,
		"code.ignore_paths": c.Code.IgnorePaths,
	} {
		for index, pattern := range patterns {
			if err := validateGlobPattern(pattern); err != nil {
				return fmt.Errorf("%s[%d] is invalid: %w", label, index, err)
			}
		}
	}
	return nil
}

// validateAnnotationKeys checks that code.annotations uses the protocol roles.
func (c *Config) validateAnnotationKeys() error {
	expected := expectedAnnotationKeys()
	for key := range c.Code.Annotations {
		if !expected[key] {
			return fmt.Errorf("code.annotations has unknown key %q; expected keys: spec, test, test_contract", key)
		}
	}
	for key := range expected {
		if _, ok := c.Code.Annotations[key]; !ok {
			return fmt.Errorf("code.annotations is missing required protocol key %q", key)
		}
	}
	return nil
}

func (c *Config) validateDocumentationUnits() error {
	packages := make(map[string]bool, len(c.Docs.Units))
	for index := range c.Docs.Units {
		unit := &c.Docs.Units[index]
		unit.Package = filepath.ToSlash(filepath.Clean(strings.TrimSpace(unit.Package)))
		if unit.Package == "." || unit.Package == "" || filepath.IsAbs(unit.Package) ||
			unit.Package == ".." || strings.HasPrefix(unit.Package, "../") {
			return fmt.Errorf("docs.units[%d].package must be a normalized project-relative documentation package", index)
		}
		if packages[unit.Package] {
			return fmt.Errorf("docs.units repeats package %q", unit.Package)
		}
		packages[unit.Package] = true
		if len(unit.Sources) == 0 {
			return fmt.Errorf("docs.units[%d].sources must contain at least one project-relative source pattern", index)
		}
		seenSources := make(map[string]bool, len(unit.Sources))
		for sourceIndex, source := range unit.Sources {
			normalized := filepath.ToSlash(filepath.Clean(strings.TrimSpace(source)))
			if normalized == "." || normalized == "" || filepath.IsAbs(normalized) ||
				normalized == ".." || strings.HasPrefix(normalized, "../") {
				return fmt.Errorf("docs.units[%d].sources[%d] must be a project-relative source pattern", index, sourceIndex)
			}
			if seenSources[normalized] {
				return fmt.Errorf("docs.units[%d] repeats source pattern %q", index, normalized)
			}
			if err := validateGlobPattern(normalized); err != nil {
				return fmt.Errorf("docs.units[%d].sources[%d] is invalid: %w", index, sourceIndex, err)
			}
			seenSources[normalized] = true
			unit.Sources[sourceIndex] = normalized
		}
	}
	return nil
}

func validateGlobPattern(pattern string) error {
	if strings.TrimSpace(pattern) == "" {
		return fmt.Errorf("glob pattern must not be empty")
	}
	for _, segment := range strings.FieldsFunc(filepath.ToSlash(pattern), func(r rune) bool { return r == '/' }) {
		if segment == "**" {
			continue
		}
		if strings.Contains(segment, "**") {
			return fmt.Errorf("** must occupy a complete path segment in %q", pattern)
		}
		if _, err := filepath.Match(segment, ""); err != nil {
			return fmt.Errorf("glob pattern %q: %w", pattern, err)
		}
	}
	return nil
}

func (c *Config) validateAnnotationValues() error {
	seen := make(map[string]string, len(c.Code.Annotations))
	for _, key := range []string{"spec", "test", "test_contract"} {
		value := c.Code.Annotations[key]
		if value == "" || strings.TrimSpace(value) != value ||
			!strings.HasPrefix(value, "@") || strings.ContainsAny(value, " \t\r\n") {
			return fmt.Errorf(
				"code.annotations.%s must be a non-empty @-prefixed token without whitespace",
				key,
			)
		}
		normalized := strings.ToLower(value)
		if previous, exists := seen[normalized]; exists {
			return fmt.Errorf(
				"code.annotations.%s duplicates code.annotations.%s value %q",
				key,
				previous,
				value,
			)
		}
		seen[normalized] = key
	}
	return nil
}
