// Package collector provides testing utilities for the collector module.

// Spec: docs/internal/collector/spec.md
// Test: docs/internal/collector/testing.md
package collector

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jingxu9x/idd-cli/internal/config"
)

// @test TEST-INTERNAL_COLLECTOR-010
func TestParseFrontmatter(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantNil bool
		wantErr bool
	}{
		{
			name: "valid frontmatter",
			content: `---
markers:
  - id: SPEC-BE-001
    name: Test
---
# SPEC-BE-001
Content`,
			wantNil: false,
			wantErr: false,
		},
		{
			name:    "no frontmatter",
			content: "# Just a heading\nSome content",
			wantNil: true,
			wantErr: false,
		},
		{
			name:    "empty content",
			content: "",
			wantNil: true,
			wantErr: false,
		},
		{
			name: "code block with dashes",
			content: `---
markers:
  - id: SPEC-BE-001
---
` + "```yaml\n---\nmarkers:\n  - id: WRONG\n---\n```\n# SPEC-BE-001",
			wantNil: false,
			wantErr: false,
		},
		{
			name: "frontmatter with describe",
			content: `---
markers:
  - id: SPEC-BE-001
    name: Test Spec
    describe: Validates user authentication
---
# SPEC-BE-001`,
			wantNil: false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, err := ParseFrontmatter(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFrontmatter() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantNil && fm != nil {
				t.Errorf("ParseFrontmatter() = %v, want nil", fm)
			}
			if !tt.wantNil && fm == nil {
				t.Errorf("ParseFrontmatter() = nil, want non-nil")
			}
		})
	}
}

// @test TEST-INTERNAL_COLLECTOR-011
func TestValidateFrontmatterMarkers(t *testing.T) {
	tests := []struct {
		name    string
		fm      *Frontmatter
		content string
		path    string
		wantLen int
	}{
		{
			name: "marker defined in heading",
			fm: &Frontmatter{
				Markers: []Marker{{ID: "SPEC-BE-001", Name: "Test"}},
			},
			content: "# SPEC-BE-001: Test heading\nSome content",
			path:    "docs/spec.md",
			wantLen: 0,
		},
		{
			name: "marker not in content",
			fm: &Frontmatter{
				Markers: []Marker{{ID: "SPEC-BE-001", Name: "Test"}},
			},
			content: "# Other heading\nSome content",
			path:    "docs/spec.md",
			wantLen: 1,
		},
		{
			name:    "nil frontmatter",
			fm:      nil,
			content: "# SPEC-BE-001",
			path:    "docs/spec.md",
			wantLen: 0,
		},
		{
			name: "marker referenced but not as heading",
			fm: &Frontmatter{
				Markers: []Marker{{ID: "SPEC-BE-001", Name: "Test"}},
			},
			content: "See `SPEC-BE-001` for details",
			path:    "docs/spec.md",
			wantLen: 1,
		},
		{
			name: "body reference not in frontmatter but in backticks",
			fm: &Frontmatter{
				Markers: []Marker{{ID: "SPEC-BE-001", Name: "Test"}},
			},
			content: "# SPEC-BE-001: Test heading\nSee `SPEC-BE-002` for details",
			path:    "docs/spec.md",
			wantLen: 0,
		},
		{
			name: "body reference in frontmatter and as heading is valid",
			fm: &Frontmatter{
				Markers: []Marker{
					{ID: "SPEC-BE-001", Name: "Test 1"},
					{ID: "SPEC-BE-002", Name: "Test 2"},
				},
			},
			content: "# SPEC-BE-001: Test 1 heading\n# SPEC-BE-002: Test 2 heading\nSee `SPEC-BE-002` for details",
			path:    "docs/spec.md",
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateFrontmatterMarkers(tt.fm, tt.content, tt.path)
			if len(errs) != tt.wantLen {
				t.Errorf("ValidateFrontmatterMarkers() got %d errors, want %d: %v", len(errs), tt.wantLen, errs)
			}
		})
	}
}

// @test TEST-INTERNAL_COLLECTOR-012
func TestExtractDefinedMarkers(t *testing.T) {
	content := `# SPEC-BE-001
# CONTRACT-BE-001
# Other heading
No ID here`
	result := extractDefinedMarkers(content)

	if !result["SPEC-BE-001"] {
		t.Error("Should extract SPEC-BE-001 from heading")
	}
	if !result["CONTRACT-BE-001"] {
		t.Error("Should extract CONTRACT-BE-001 from heading")
	}
	if result["NONEXISTENT"] {
		t.Error("Should not contain NONEXISTENT")
	}
}

// @test TEST-INTERNAL_COLLECTOR-013
// @test TEST-INTERNAL_COLLECTOR-014
func TestExtractReferencedMarkers(t *testing.T) {
	content := `# Heading
See SPEC-BE-001 and TEST-BE-001 for details.
Also see CONTRACT-BE-002.`
	result := extractReferencedMarkers(content)

	if !result["SPEC-BE-001"] {
		t.Error("Should extract SPEC-BE-001")
	}
	if !result["TEST-BE-001"] {
		t.Error("Should extract TEST-BE-001")
	}
	if !result["CONTRACT-BE-002"] {
		t.Error("Should extract CONTRACT-BE-002")
	}
}

// @test TEST-INTERNAL_COLLECTOR-019
func TestExtractIDDRefs(t *testing.T) {
	tests := []struct {
		line    string
		wantLen int
	}{
		{"# SPEC-BE-001", 1},
		{"# SPEC-BE-001 and TEST-BE-001", 2},
		{"No ID here", 0},
		{"See @spec SPEC-BE-001", 1},
		{"text SPEC-BE-001 text CONTRACT-BE-002 text", 2},
	}

	for _, tt := range tests {
		result := extractIDDRefs(tt.line)
		if len(result) != tt.wantLen {
			t.Errorf("extractIDDRefs(%q) got %d refs, want %d", tt.line, len(result), tt.wantLen)
		}
	}
}

// @test TEST-INTERNAL_COLLECTOR-015
func TestGetExpectedFilename(t *testing.T) {
	tests := []struct {
		idType string
		want   string
	}{
		{"SPEC", "spec.md"},
		{"spec", "spec.md"},
		{"CONTRACT", "contract.md"},
		{"TEST", "testing.md"},
		{"DESIGN", "design.md"},
		{"UNKNOWN", ""},
	}

	for _, tt := range tests {
		result := GetExpectedFilename(tt.idType)
		if result != tt.want {
			t.Errorf("GetExpectedFilename(%q) = %q, want %q", tt.idType, result, tt.want)
		}
	}
}

// @test TEST-INTERNAL_COLLECTOR-017
func TestValidateDocumentStructure(t *testing.T) {
	tests := []struct {
		path   string
		idType string
		want   bool
	}{
		{"docs/spec.md", "SPEC", false},
		{"docs/backend/spec.md", "SPEC", false},
		{"SPEC.md", "SPEC", false},
		{"docs/backend/testing.md", "TEST", false},
	}

	for _, tt := range tests {
		err := ValidateDocumentStructure(tt.path, tt.idType)
		hasError := err != nil
		if hasError != tt.want {
			t.Errorf("ValidateDocumentStructure(%q, %q) error = %v, want error = %v", tt.path, tt.idType, err, tt.want)
		}
	}
}

// @test TEST-INTERNAL_COLLECTOR-016
func TestIsRootDocFile(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"docs/spec.md", true},
		{"docs/", true},
		{"docs/backend/spec.md", false},
		{"other/spec.md", false},
	}

	for _, tt := range tests {
		result := isRootDocFile(tt.path)
		if result != tt.want {
			t.Errorf("isRootDocFile(%q) = %v, want %v", tt.path, result, tt.want)
		}
	}
}

// @test TEST-INTERNAL_COLLECTOR-001
func TestDocCollector_Collect(t *testing.T) {
	tmpDir := t.TempDir()

	doc := "# `SPEC-BE-001`\nThis is a spec.\n\nSee `TEST-BE-001` for testing.\n"
	err := os.WriteFile(filepath.Join(tmpDir, "spec.md"), []byte(doc), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cfg := config.Default()
	coll := NewDocCollector(cfg)
	set, collErrors, err := coll.Collect(nil, tmpDir)
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if len(collErrors) != 0 {
		t.Errorf("Expected 0 collector errors, got %d: %v", len(collErrors), collErrors)
	}

	if !set.Has("SPEC-BE-001") {
		t.Error("Should have SPEC-BE-001")
	}
	if !set.Has("TEST-BE-001") {
		t.Error("Should have TEST-BE-001")
	}
}

// @test TEST-INTERNAL_COLLECTOR-002
func TestDocCollector_Collect_WithFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()

	doc := `---
markers:
  - id: SPEC-BE-001
    name: Test Spec
    describe: Validates authentication
---
# SPEC-BE-001

See TEST-BE-001.
`
	err := os.WriteFile(filepath.Join(tmpDir, "spec.md"), []byte(doc), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cfg := config.Default()
	coll := NewDocCollector(cfg)
	set, _, err := coll.Collect(nil, tmpDir)
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	specID, ok := set.Get("SPEC-BE-001")
	if !ok {
		t.Fatal("SPEC-BE-001 not found")
	}
	if specID.Describe != "Validates authentication" {
		t.Errorf("Describe = %q, want %q", specID.Describe, "Validates authentication")
	}
}

// @test TEST-INTERNAL_COLLECTOR-003
func TestDocCollector_Collect_FileNotFound(t *testing.T) {
	cfg := config.Default()
	coll := NewDocCollector(cfg)
	set, _, err := coll.Collect(nil, "/nonexistent/path")
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}
	if set.Count() != 0 {
		t.Errorf("Expected 0 identifiers for nonexistent path, got %d", set.Count())
	}
}

// @test TEST-INTERNAL_COLLECTOR-004
func TestDocCollector_Collect_NonMarkdownFile(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("# SPEC-BE-001"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cfg := config.Default()
	coll := NewDocCollector(cfg)
	set, _, _ := coll.Collect(nil, tmpDir)

	if set.Has("SPEC-BE-001") {
		t.Error("Should not extract from .txt files")
	}
}

// @test TEST-INTERNAL_COLLECTOR-005
func TestDocCollector_extractTitle(t *testing.T) {
	coll := NewDocCollector(config.Default())

	content := `# SPEC-BE-001: User Authentication
Some description

## Details
More content`

	title := coll.extractTitle(content, "SPEC-BE-001")
	if title != "SPEC-BE-001: User Authentication" {
		t.Errorf("extractTitle() = %q, want %q", title, "SPEC-BE-001: User Authentication")
	}
}

// @test TEST-INTERNAL_COLLECTOR-006
func TestDocCollector_extractTitle_NotFound(t *testing.T) {
	coll := NewDocCollector(config.Default())

	content := `# Heading without ID`

	title := coll.extractTitle(content, "SPEC-BE-001")
	if title != "" {
		t.Errorf("extractTitle() = %q, want empty string", title)
	}
}

// @test TEST-INTERNAL_COLLECTOR-007
func TestDocCollector_Collect_DirectoryWithNoSpec(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "readme.md"), []byte("# Readme\nNo IDs here"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cfg := config.Default()
	coll := NewDocCollector(cfg)
	set, _, _ := coll.Collect(nil, tmpDir)

	if set.Count() != 0 {
		t.Errorf("Expected 0 identifiers, got %d", set.Count())
	}
}

// @test TEST-INTERNAL_COLLECTOR-021
func TestNewDocCollector(t *testing.T) {
	cfg := config.Default()
	coll := NewDocCollector(cfg)
	if coll == nil {
		t.Error("NewDocCollector returned nil")
	}
	if coll.cfg != cfg {
		t.Error("Collector config not set correctly")
	}
}

// @test TEST-INTERNAL_COLLECTOR-008
func TestMarker_Describe(t *testing.T) {
	m := Marker{
		ID:       "SPEC-BE-001",
		Name:     "Test",
		Describe: "Description here",
	}
	if m.Describe != "Description here" {
		t.Errorf("Describe = %q, want %q", m.Describe, "Description here")
	}
}

// @test TEST-INTERNAL_COLLECTOR-009
func TestFrontmatter_Markers(t *testing.T) {
	fm := &Frontmatter{
		Markers: []Marker{
			{ID: "SPEC-BE-001", Name: "Spec 1"},
			{ID: "TEST-BE-001", Name: "Test 1"},
		},
	}
	if len(fm.Markers) != 2 {
		t.Errorf("Expected 2 markers, got %d", len(fm.Markers))
	}
}

// @test TEST-INTERNAL_COLLECTOR-022
func TestExtractSectionContent(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		markerID  string
		wantEmpty bool
		wantLen   int
	}{
		{
			name: "simple section",
			content: "## SPEC-BE-001: Test\n\n**Tests:** `TEST-BE-001`\n\n---\n\n## SPEC-BE-002: Other",
			markerID:  "SPEC-BE-001",
			wantEmpty: false,
			wantLen:   40,
		},
		{
			name:      "marker not found",
			content:   "# Just a heading",
			markerID:  "SPEC-BE-999",
			wantEmpty: true,
			wantLen:   0,
		},
		{
			name: "section with multiple paragraphs",
			content: "## SPEC-BE-001: Test\n\nSome content here.\n\nMore content.\n\n**Tests:** `TEST-BE-001`\n\n---\n\n## SPEC-BE-002: Other",
			markerID:  "SPEC-BE-001",
			wantEmpty: false,
			wantLen:   80,
		},
		{
			name: "H3 heading section",
			content: "### SPEC-BE-001\n\nContent here.\n\n---\n\n## SPEC-BE-002",
			markerID:  "SPEC-BE-001",
			wantEmpty: false,
			wantLen:   20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractSectionContent(tt.content, tt.markerID)
			if tt.wantEmpty && result != "" {
				t.Errorf("extractSectionContent() = %q, want empty", result)
			}
			if !tt.wantEmpty && result == "" {
				t.Errorf("extractSectionContent() returned empty, want non-empty")
			}
			if len(result) < tt.wantLen {
				t.Errorf("extractSectionContent() length = %d, want at least %d", len(result), tt.wantLen)
			}
		})
	}
}

// @test TEST-INTERNAL_COLLECTOR-023
func TestExtractSpecCoverage(t *testing.T) {
	tests := []struct {
		name   string
		section string
		want   []string
	}{
		{
			name:    "single spec",
			section: "**Spec Coverage:** `SPEC-BE-001`",
			want:    []string{"SPEC-BE-001"},
		},
		{
			name:    "multiple specs",
			section: "**Spec Coverage:** `SPEC-BE-001`, `SPEC-BE-002`",
			want:    []string{"SPEC-BE-001", "SPEC-BE-002"},
		},
		{
			name:    "no spec coverage field",
			section: "Some content without spec coverage",
			want:    nil,
		},
		{
			name:    "with spaces",
			section: "**Spec Coverage:** `SPEC-BE-001`, `SPEC-BE-002`, `SPEC-BE-003`",
			want:    []string{"SPEC-BE-001", "SPEC-BE-002", "SPEC-BE-003"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractSpecCoverage(tt.section)
			if len(result) != len(tt.want) {
				t.Errorf("extractSpecCoverage() got %d items, want %d", len(result), len(tt.want))
				return
			}
			for i, v := range result {
				if v != tt.want[i] {
					t.Errorf("extractSpecCoverage()[%d] = %q, want %q", i, v, tt.want[i])
				}
			}
		})
	}
}

// @test TEST-INTERNAL_COLLECTOR-024
func TestExtractTestsField(t *testing.T) {
	tests := []struct {
		name    string
		section string
		want    []string
	}{
		{
			name:    "single test",
			section: "**Tests:** `TEST-BE-001`",
			want:    []string{"TEST-BE-001"},
		},
		{
			name:    "multiple tests",
			section: "**Tests:** `TEST-BE-001`, `TEST-BE-002`",
			want:    []string{"TEST-BE-001", "TEST-BE-002"},
		},
		{
			name:    "no tests field",
			section: "Some content without tests",
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTestsField(tt.section)
			if len(result) != len(tt.want) {
				t.Errorf("extractTestsField() got %d items, want %d", len(result), len(tt.want))
				return
			}
			for i, v := range result {
				if v != tt.want[i] {
					t.Errorf("extractTestsField()[%d] = %q, want %q", i, v, tt.want[i])
				}
			}
		})
	}
}

// @test TEST-INTERNAL_COLLECTOR-025
func TestExtractContractField(t *testing.T) {
	tests := []struct {
		name    string
		section string
		want    []string
	}{
		{
			name:    "single contract",
			section: "**Contract:** `DocCollector`",
			want:    []string{"DocCollector"},
		},
		{
			name:    "no contract field",
			section: "Some content without contract",
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractContractField(tt.section)
			if len(result) != len(tt.want) {
				t.Errorf("extractContractField() got %d items, want %d", len(result), len(tt.want))
				return
			}
			for i, v := range result {
				if v != tt.want[i] {
					t.Errorf("extractContractField()[%d] = %q, want %q", i, v, tt.want[i])
				}
			}
		})
	}
}

// @test TEST-INTERNAL_COLLECTOR-026
func TestExtractImplementsField(t *testing.T) {
	tests := []struct {
		name    string
		section string
		want    []string
	}{
		{
			name:    "single implements",
			section: "**Implements:** `SPEC-INTERNAL_COLLECTOR-001`",
			want:    []string{"SPEC-INTERNAL_COLLECTOR-001"},
		},
		{
			name:    "multiple implements",
			section: "**Implements:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-002`",
			want:    []string{"SPEC-INTERNAL_COLLECTOR-001", "SPEC-INTERNAL_COLLECTOR-002"},
		},
		{
			name:    "no implements field",
			section: "Some content without implements",
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractImplementsField(tt.section)
			if len(result) != len(tt.want) {
				t.Errorf("extractImplementsField() got %d items, want %d", len(result), len(tt.want))
				return
			}
			for i, v := range result {
				if v != tt.want[i] {
					t.Errorf("extractImplementsField()[%d] = %q, want %q", i, v, tt.want[i])
				}
			}
		})
	}
}

// @test TEST-INTERNAL_COLLECTOR-027
func TestDocCollector_Collect_WithSpecAndTestLinks(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a testing.md with Spec Coverage
	testingDoc := "---\nmarkers:\n  - id: TEST-BE-001\n    name: Test Case\n---\n## TEST-BE-001\n\n**Spec Coverage:** `SPEC-BE-001`\n\n---\n"
	err := os.WriteFile(filepath.Join(tmpDir, "testing.md"), []byte(testingDoc), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cfg := config.Default()
	coll := NewDocCollector(cfg)
	set, _, err := coll.Collect(nil, tmpDir)
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if !set.Has("TEST-BE-001") {
		t.Error("Should have TEST-BE-001")
	}

	// Verify the test has links to the spec
	testID, ok := set.Get("TEST-BE-001")
	if !ok {
		t.Fatal("TEST-BE-001 not found")
	}
	if len(testID.Links) == 0 {
		t.Error("TEST-BE-001 should have links to SPEC")
	}
}

// @test TEST-INTERNAL_COLLECTOR-028
func TestDocCollector_Collect_WithContractLink(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a spec.md with Contract field
	specDoc := "---\nmarkers:\n  - id: SPEC-BE-001\n    name: Test Spec\n---\n## SPEC-BE-001\n\n**Contract:** `DocCollector`\n\n**Tests:** `TEST-BE-001`\n\n---\n"
	err := os.WriteFile(filepath.Join(tmpDir, "spec.md"), []byte(specDoc), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cfg := config.Default()
	coll := NewDocCollector(cfg)
	set, _, err := coll.Collect(nil, tmpDir)
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if !set.Has("SPEC-BE-001") {
		t.Error("Should have SPEC-BE-001")
	}

	specID, ok := set.Get("SPEC-BE-001")
	if !ok {
		t.Fatal("SPEC-BE-001 not found")
	}
	// Contract link should be tracked
	foundContract := false
	for _, link := range specID.Links {
		if link == "DocCollector" {
			foundContract = true
			break
		}
	}
	if !foundContract {
		t.Error("SPEC-BE-001 should have DocCollector contract link")
	}
}

// @test TEST-INTERNAL_COLLECTOR-029
func TestDocCollector_shouldIgnore(t *testing.T) {
	cfg := config.Default()
	cfg.Docs.IgnorePaths = []string{"docs/internal/**"}
	coll := NewDocCollector(cfg)

	if !coll.shouldIgnore("docs/internal/collector/spec.md") {
		t.Error("Should ignore paths matching internal/**")
	}

	if coll.shouldIgnore("docs/cmd/spec.md") {
		t.Error("Should not ignore paths not matching pattern")
	}
}
