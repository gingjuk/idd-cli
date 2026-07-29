// Package config provides testing utilities for the config module.
package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// @test-contract TEST-INTERNAL_CONFIG-002
func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Version != "1.0" {
		t.Errorf("Version = %q, want 1.0", cfg.Version)
	}

	if cfg.Validation.ConsistencyCheck.Enabled {
		t.Error("deprecated ConsistencyCheck should be zero-valued by default")
	}

	if !cfg.Validation.RequireSpecFields {
		t.Error("RequireSpecFields should be enabled by default")
	}

	if cfg.Validation.ConsistencyCheck.Threshold != 0 {
		t.Errorf("ConsistencyCheck.Threshold = %f, want 0", cfg.Validation.ConsistencyCheck.Threshold)
	}

	if !reflect.DeepEqual(cfg.Code.Patterns, defaultCodePatterns) {
		t.Errorf("Code.Patterns = %#v, want %#v", cfg.Code.Patterns, defaultCodePatterns)
	}
}

// @test-contract TEST-INTERNAL_CONFIG-004
func TestConfig_WorkdirPaths(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "external.go")
	cfg := Default()

	tests := []struct {
		name        string
		root        string
		path        string
		wantResolve string
		wantDisplay string
		wantErr     bool
	}{
		{
			name:        "absolute workdir resolves project path",
			root:        root,
			path:        filepath.Join("internal", "config", "config.go"),
			wantResolve: filepath.Join(root, "internal", "config", "config.go"),
			wantDisplay: filepath.Join("internal", "config", "config.go"),
		},
		{
			name:        "absolute path inside root becomes relative",
			root:        root,
			path:        filepath.Join(root, "docs", "spec.md"),
			wantResolve: filepath.Join(root, "docs", "spec.md"),
			wantDisplay: filepath.Join("docs", "spec.md"),
		},
		{
			name:        "absolute path outside root stays absolute",
			root:        root,
			path:        outside,
			wantResolve: outside,
			wantDisplay: outside,
		},
		{
			name:        "empty workdir preserves caller paths",
			path:        filepath.Join("docs", "spec.md"),
			wantResolve: filepath.Join("docs", "spec.md"),
			wantDisplay: filepath.Join("docs", "spec.md"),
		},
		{
			name:    "relative workdir is rejected",
			root:    "relative/project",
			path:    "spec.md",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := cfg.SetWorkdir(test.root); (err != nil) != test.wantErr {
				t.Fatalf("SetWorkdir(%q) error = %v, wantErr %t", test.root, err, test.wantErr)
			}
			if test.wantErr {
				return
			}
			if got := cfg.ResolvePath(test.path); got != test.wantResolve {
				t.Errorf("ResolvePath(%q) = %q, want %q", test.path, got, test.wantResolve)
			}
			if got := cfg.DisplayPath(test.path); got != test.wantDisplay {
				t.Errorf("DisplayPath(%q) = %q, want %q", test.path, got, test.wantDisplay)
			}
		})
	}
}

// @test-contract TEST-INTERNAL_CONFIG-004
func TestConfig_WorkdirNormalizesConfigurationWarnings(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, ".idd.yaml")
	data := []byte("validation:\n  consistency_check: {}\n")
	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if err := cfg.SetWorkdir(root); err != nil {
		t.Fatalf("SetWorkdir() error = %v", err)
	}
	warnings := cfg.DeprecationWarnings()
	if len(warnings) != 1 || warnings[0].Source != ".idd.yaml" {
		t.Errorf("DeprecationWarnings() = %#v, want project-relative source", warnings)
	}
}

// @test-contract TEST-INTERNAL_CONFIG-002
func TestDefaultMatchesExampleConfiguration(t *testing.T) {
	examplePath := filepath.Join("..", "..", "examples", "idd-config-example.yaml")
	example, err := Load(examplePath)
	if err != nil {
		t.Fatalf("Load(%s) error = %v", examplePath, err)
	}
	if !reflect.DeepEqual(example, Default()) {
		t.Errorf("example configuration does not match Default():\nexample = %#v\ndefault = %#v", example, Default())
	}
}

// @test-contract TEST-INTERNAL_CONFIG-001
func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name          string
		cfg           *Config
		wantThreshold float64
	}{
		{
			name: "valid config",
			cfg: &Config{
				Version: "1.0",
				Validation: ValidationConfig{
					ConsistencyCheck: ConsistencyCheck{
						Enabled:   true,
						Threshold: 0.5,
					},
				},
			},
			wantThreshold: 0.5,
		},
		{
			name: "threshold below zero",
			cfg: &Config{
				Version: "1.0",
				Validation: ValidationConfig{
					ConsistencyCheck: ConsistencyCheck{
						Enabled:   true,
						Threshold: -0.5,
					},
				},
			},
			wantThreshold: -0.5,
		},
		{
			name: "threshold above one",
			cfg: &Config{
				Version: "1.0",
				Validation: ValidationConfig{
					ConsistencyCheck: ConsistencyCheck{
						Enabled:   true,
						Threshold: 1.5,
					},
				},
			},
			wantThreshold: 1.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if err != nil {
				t.Errorf("Validate() error = %v", err)
			}

			if tt.cfg.Validation.ConsistencyCheck.Threshold != tt.wantThreshold {
				t.Errorf(
					"deprecated threshold = %f, want unchanged %f",
					tt.cfg.Validation.ConsistencyCheck.Threshold,
					tt.wantThreshold,
				)
			}
		})
	}
}

// @test-contract TEST-INTERNAL_CONFIG-003
func TestLoad(t *testing.T) {
	tmpDir := t.TempDir()

	yamlContent := `version: "1.0"
docs:
  patterns:
    - "docs/**/*.md"
validation:
  consistency_check:
    enabled: true
    threshold: 0.5
`
	cfgPath := filepath.Join(tmpDir, "idd.yaml")
	err := os.WriteFile(cfgPath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	loaded, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Validation.ConsistencyCheck.Threshold != 0.5 {
		t.Errorf("Threshold = %f, want 0.5", loaded.Validation.ConsistencyCheck.Threshold)
	}

	if !loaded.Validation.ConsistencyCheck.Enabled {
		t.Error("Enabled should be true")
	}
	warnings := loaded.DeprecationWarnings()
	if len(warnings) != 1 {
		t.Fatalf("DeprecationWarnings() = %#v, want one warning", warnings)
	}
	warning := warnings[0]
	if warning.Path != "validation.consistency_check" ||
		warning.Source != cfgPath ||
		!strings.Contains(warning.Message, "deprecated and ignored") ||
		!strings.Contains(warning.Message, "remove it") ||
		!strings.Contains(warning.Message, "docs review-context") {
		t.Errorf("deprecation warning = %#v", warning)
	}
}

// @test-contract TEST-INTERNAL_CONFIG-001, TEST-INTERNAL_CONFIG-003
func TestLoad_DeprecationWarningUsesKeyPresence(t *testing.T) {
	tests := []struct {
		name         string
		validation   string
		wantWarnings int
	}{
		{
			name: "absent key remains quiet",
			validation: `validation:
  require_spec_fields: true
`,
		},
		{
			name: "zero values still warn",
			validation: `validation:
  consistency_check:
    enabled: false
    threshold: 0
`,
			wantWarnings: 1,
		},
		{
			name: "null mapping still warns",
			validation: `validation:
  consistency_check:
`,
			wantWarnings: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfgPath := filepath.Join(t.TempDir(), "idd.yaml")
			if err := os.WriteFile(cfgPath, []byte("version: \"1.0\"\n"+test.validation), 0o644); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}
			cfg, err := Load(cfgPath)
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}
			if got := len(cfg.DeprecationWarnings()); got != test.wantWarnings {
				t.Errorf("warnings = %d, want %d", got, test.wantWarnings)
			}
		})
	}
}

// @test-contract TEST-INTERNAL_CONFIG-001
func TestDeprecationWarnings_ReturnsCopy(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "idd.yaml")
	data := []byte("validation:\n  consistency_check: {}\n")
	if err := os.WriteFile(cfgPath, data, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	warnings := cfg.DeprecationWarnings()
	warnings[0].Path = "changed"
	if got := cfg.DeprecationWarnings()[0].Path; got != "validation.consistency_check" {
		t.Errorf("stored warning path = %q, want immutable copy", got)
	}
}

// @test-contract TEST-INTERNAL_CONFIG-003
func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/idd.yaml")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

// @test-contract TEST-INTERNAL_CONFIG-001
func TestValidateAnnotationKeys(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "all keys present",
			cfg: &Config{
				Code: CodeConfig{
					Annotations: map[string]string{
						"spec":          "@implement",
						"test":          "@test",
						"test_contract": "@test-contract",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing spec key",
			cfg: &Config{
				Code: CodeConfig{
					Annotations: map[string]string{
						"test":          "@test",
						"test_contract": "@test-contract",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing test key",
			cfg: &Config{
				Code: CodeConfig{
					Annotations: map[string]string{
						"spec":          "@implement",
						"test_contract": "@test-contract",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "unknown key",
			cfg: &Config{
				Code: CodeConfig{
					Annotations: map[string]string{
						"spec":          "@implement",
						"test":          "@test",
						"test_contract": "@test-contract",
						"extra":         "@something",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "partial keys",
			cfg: &Config{
				Code: CodeConfig{
					Annotations: map[string]string{
						"spec": "@implement",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty value",
			cfg: &Config{
				Code: CodeConfig{
					Annotations: map[string]string{
						"spec":          "",
						"test":          "@test",
						"test_contract": "@test-contract",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "value without annotation prefix",
			cfg: &Config{
				Code: CodeConfig{
					Annotations: map[string]string{
						"spec":          "implement",
						"test":          "@test",
						"test_contract": "@test-contract",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "duplicate values",
			cfg: &Config{
				Code: CodeConfig{
					Annotations: map[string]string{
						"spec":          "@trace",
						"test":          "@TRACE",
						"test_contract": "@test-contract",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

// @test-contract TEST-INTERNAL_CONFIG-003
func TestLoad_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()

	cfgPath := filepath.Join(tmpDir, "idd.yaml")
	err := os.WriteFile(cfgPath, []byte("invalid: yaml: content:"), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	_, err = Load(cfgPath)
	if err == nil {
		t.Error("Expected error for invalid YAML")
	}
}
