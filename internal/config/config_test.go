// Package config provides testing utilities for the config module.
package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// @test-contract TEST-INTERNAL_CONFIG-002
func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Version != "1.0" {
		t.Errorf("Version = %q, want 1.0", cfg.Version)
	}

	if cfg.Validation.ConsistencyCheck.Enabled {
		t.Error("ConsistencyCheck should be opt-in by default")
	}

	if !cfg.Validation.RequireSpecFields {
		t.Error("RequireSpecFields should be enabled by default")
	}

	if cfg.Validation.ConsistencyCheck.Threshold != 0.3 {
		t.Errorf("ConsistencyCheck.Threshold = %f, want 0.3", cfg.Validation.ConsistencyCheck.Threshold)
	}

	if !reflect.DeepEqual(cfg.Code.Patterns, defaultCodePatterns) {
		t.Errorf("Code.Patterns = %#v, want %#v", cfg.Code.Patterns, defaultCodePatterns)
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
		name      string
		cfg       *Config
		wantPanic bool
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
			wantPanic: false,
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
			wantPanic: false,
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
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if err != nil {
				t.Errorf("Validate() error = %v", err)
			}

			if tt.cfg.Validation.ConsistencyCheck.Threshold < 0 {
				t.Errorf("Threshold should be normalized to >= 0, got %f", tt.cfg.Validation.ConsistencyCheck.Threshold)
			}
			if tt.cfg.Validation.ConsistencyCheck.Threshold > 1.0 {
				t.Errorf("Threshold should be normalized to <= 1.0, got %f", tt.cfg.Validation.ConsistencyCheck.Threshold)
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
