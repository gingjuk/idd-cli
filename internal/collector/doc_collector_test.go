package collector

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jingxu9x/idd-link-validator/internal/config"
)

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
			content: "# SPEC-BE-001\nSome content",
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

func TestExtractModuleName(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"docs/backend/spec.md", "BACKEND"},
		{"docs/frontend/test.md", "FRONTEND"},
		{"docs/spec.md", ""},
		{"other/spec.md", ""},
	}

	for _, tt := range tests {
		result := ExtractModuleName(tt.path)
		if result != tt.want {
			t.Errorf("ExtractModuleName(%q) = %q, want %q", tt.path, result, tt.want)
		}
	}
}

func TestIsKnownAbbreviation(t *testing.T) {
	if !isKnownAbbreviation("BE", "BACKEND") {
		t.Error("BE should match BACKEND")
	}
	if !isKnownAbbreviation("be", "backend") {
		t.Error("be should match backend (case insensitive)")
	}
	if isKnownAbbreviation("XX", "YYYY") {
		t.Error("XX should not match YYYY")
	}
}

func TestValidateModulePrefix(t *testing.T) {
	tests := []struct {
		id     string
		path   string
		wantOk bool
	}{
		{"SPEC-BE-001", "docs/backend/spec.md", false},
		{"SPEC-FE-001", "docs/backend/spec.md", true},
		{"SPEC-XX-001", "docs/backend/spec.md", true},
		{"SPEC-BE-001", "docs/spec.md", false},
	}

	for _, tt := range tests {
		err := ValidateModulePrefix(tt.id, tt.path)
		hasError := err != nil
		if hasError != tt.wantOk {
			t.Errorf("ValidateModulePrefix(%q, %q) error = %v, want error = %v", tt.id, tt.path, err, tt.wantOk)
		}
	}
}

func TestDocCollector_Collect(t *testing.T) {
	tmpDir := t.TempDir()

	doc := `# SPEC-BE-001
This is a spec.

See TEST-BE-001 for testing.
`
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

func TestDocCollector_extractTitle_NotFound(t *testing.T) {
	coll := NewDocCollector(config.Default())

	content := `# Heading without ID`

	title := coll.extractTitle(content, "SPEC-BE-001")
	if title != "" {
		t.Errorf("extractTitle() = %q, want empty string", title)
	}
}

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
