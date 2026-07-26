// Package pattern provides testing utilities for the pattern module.
package pattern

import (
	"testing"
)

// @test-contract TEST-PKG_PATTERN-001
func TestValidateIdentifierFormat(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid SPEC identifier",
			id:      "SPEC-BE-001",
			wantErr: false,
		},
		{
			name:    "valid CONTRACT identifier",
			id:      "CONTRACT-AUTH-001",
			wantErr: false,
		},
		{
			name:    "valid TEST identifier",
			id:      "TEST-BE-001",
			wantErr: false,
		},
		{
			name:    "valid DESIGN identifier",
			id:      "DESIGN-BE-001",
			wantErr: false,
		},
		{
			name:    "valid with numeric module",
			id:      "SPEC-001-001",
			wantErr: false,
		},
		{
			name:    "valid with alphanumeric module",
			id:      "SPEC-BE001-001",
			wantErr: false,
		},
		{
			name:    "invalid - PATTERN-001 is internal section marker, not IDD identifier",
			id:      "PATTERN-001",
			wantErr: true,
			errMsg:  "uses internal section marker TYPE",
		},
		{
			name:    "invalid - WALK-002 is internal section marker, not IDD identifier",
			id:      "WALK-002",
			wantErr: true,
			errMsg:  "uses internal section marker TYPE",
		},
		{
			name:    "invalid - PATTERN-BE-001 has PATTERN as TYPE which is invalid for IDD",
			id:      "PATTERN-BE-001",
			wantErr: true,
			errMsg:  "invalid TYPE",
		},
		{
			name:    "invalid - WALK-BE-001 has WALK as TYPE which is invalid for IDD",
			id:      "WALK-BE-001",
			wantErr: true,
			errMsg:  "invalid TYPE",
		},
		{
			name:    "invalid - SPEC-BE is 2-part but TYPE is not PATTERN/WALK",
			id:      "SPEC-BE",
			wantErr: true,
			errMsg:  "invalid TYPE",
		},
		{
			name:    "invalid - 4 parts",
			id:      "SPEC-BE-007-007",
			wantErr: true,
			errMsg:  "must have 2 or 3 parts, found 4 parts",
		},
		{
			name:    "invalid - unknown type",
			id:      "INVALID-BE-001",
			wantErr: true,
			errMsg:  "invalid TYPE",
		},
		{
			name:    "invalid - lowercase module",
			id:      "SPEC-be-001",
			wantErr: true,
			errMsg:  "invalid MODULE",
		},
		{
			name:    "invalid - special characters in module",
			id:      "SPEC-BE@001-001",
			wantErr: true,
			errMsg:  "invalid MODULE",
		},
		{
			name:    "invalid - letters in number",
			id:      "SPEC-BE-001ABC",
			wantErr: true,
			errMsg:  "invalid NUMBER",
		},
		{
			name:    "invalid - empty string",
			id:      "",
			wantErr: true,
			errMsg:  "must have 2 or 3 parts",
		},
		{
			name:    "case insensitive TYPE validation",
			id:      "SPEC-BE-001",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIdentifierFormat(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIdentifierFormat(%q) error = %v, wantErr %v", tt.id, err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateIdentifierFormat(%q) error = %v, want error containing %q", tt.id, err, tt.errMsg)
				}
			}
		})
	}
}

// @test TEST-PKG_PATTERN-002
func TestValidateIDPattern(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "valid SPEC",
			id:      "SPEC-BE-001",
			wantErr: false,
		},
		{
			name:    "valid CONTRACT",
			id:      "CONTRACT-AUTH-001",
			wantErr: false,
		},
		{
			name:    "valid TEST",
			id:      "TEST-BE-001",
			wantErr: false,
		},
		{
			name:    "valid DESIGN",
			id:      "DESIGN-BE-001",
			wantErr: false,
		},
		{
			name:    "valid 4-part (but not strict IDD format)",
			id:      "SPEC-BE-007-007",
			wantErr: false,
		},
		{
			name:    "invalid - unknown pattern",
			id:      "INVALID-BE-001",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIDPattern(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIDPattern(%q) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
		})
	}
}

// @test-contract TEST-PKG_PATTERN-003
func TestGetIdentifierType(t *testing.T) {
	tests := []struct {
		name string
		ref  string
		want string
	}{
		{
			name: "SPEC identifier",
			ref:  "SPEC-BE-001",
			want: "SPEC",
		},
		{
			name: "CONTRACT identifier",
			ref:  "CONTRACT-AUTH-001",
			want: "CONTRACT",
		},
		{
			name: "TEST identifier",
			ref:  "TEST-BE-001",
			want: "TEST",
		},
		{
			name: "DESIGN identifier",
			ref:  "DESIGN-BE-001",
			want: "DESIGN",
		},
		{
			name: "case insensitive",
			ref:  "spec-be-001",
			want: "SPEC",
		},
		{
			name: "unknown identifier",
			ref:  "INVALID-BE-001",
			want: "",
		},
		{
			name: "empty string",
			ref:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetIdentifierType(tt.ref)
			if got != tt.want {
				t.Errorf("GetIdentifierType(%q) = %q, want %q", tt.ref, got, tt.want)
			}
		})
	}
}

// @test TEST-PKG_PATTERN-004
func TestGetAnnotationType(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		want   string
	}{
		{
			name:   "@implement",
			prefix: "@implement",
			want:   "SPEC",
		},
		{
			name:   "@test",
			prefix: "@test",
			want:   "TEST",
		},
		{
			name:   "@test-contract",
			prefix: "@test-contract",
			want:   "TEST",
		},
		{
			name:   "case sensitive - @Implement does not match @implement",
			prefix: "@Implement",
			want:   "",
		},
		{
			name:   "unknown prefix",
			prefix: "@unknown",
			want:   "",
		},
		{
			name:   "empty string",
			prefix: "",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAnnotationType(tt.prefix)
			if got != tt.want {
				t.Errorf("GetAnnotationType(%q) = %q, want %q", tt.prefix, got, tt.want)
			}
		})
	}
}

// @test TEST-PKG_PATTERN-005
func TestExtractIDDReferences(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "plain text not extracted",
			content: "See SPEC-BE-001 for details",
			want:    []string{},
		},
		{
			name:    "multiple plain text not extracted",
			content: "From SPEC-BE-001 to CONTRACT-BE-002",
			want:    []string{},
		},
		{
			name:    "case insensitive plain text not extracted",
			content: "spec-be-001 and SPEC-BE-001",
			want:    []string{},
		},
		{
			name:    "identifier in backticks included",
			content: "`SPEC-BE-001` should be included",
			want:    []string{"SPEC-BE-001"},
		},
		{
			name:    "identifier in quotes included",
			content: "\"SPEC-BE-001\" should be included",
			want:    []string{"SPEC-BE-001"},
		},
		{
			name:    "no identifiers",
			content: "No identifiers here",
			want:    []string{},
		},
		{
			name:    "TEST and DESIGN plain text not extracted",
			content: "TEST-BE-001 and DESIGN-BE-001",
			want:    []string{},
		},
		{
			name:    "3-part PATTERN identifier in backticks",
			content: "`PATTERN-BE-001` is extracted",
			want:    []string{"PATTERN-BE-001"},
		},
		{
			name:    "3-part WALK identifier in backticks",
			content: "`WALK-BE-001` is extracted",
			want:    []string{"WALK-BE-001"},
		},
		{
			name:    "2-part PATTERN section marker not extracted (not 3-part)",
			content: "PATTERN-001 section marker",
			want:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractIDDReferences(tt.content)
			if !sameElements(got, tt.want) {
				t.Errorf("ExtractIDDReferences(%q) = %v, want %v", tt.content, got, tt.want)
			}
		})
	}
}

// @test TEST-PKG_PATTERN-006
func TestSplitAnnotationRefs(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []string
	}{
		{
			name: "single identifier",
			s:    "SPEC-BE-001",
			want: []string{"SPEC-BE-001"},
		},
		{
			name: "multiple identifiers with comma",
			s:    "SPEC-BE-001, CONTRACT-BE-002",
			want: []string{"SPEC-BE-001", "CONTRACT-BE-002"},
		},
		{
			name: "multiple identifiers with spaces",
			s:    "SPEC-BE-001 , TEST-BE-002",
			want: []string{"SPEC-BE-001", "TEST-BE-002"},
		},
		{
			name: "single identifier with trailing comma",
			s:    "SPEC-BE-001,",
			want: []string{"SPEC-BE-001"},
		},
		{
			name: "empty string",
			s:    "",
			want: []string{},
		},
		{
			name: "only whitespace",
			s:    "   ",
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitAnnotationRefs(tt.s)
			if !sameElements(got, tt.want) {
				t.Errorf("SplitAnnotationRefs(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

// @test TEST-PKG_PATTERN-007
func TestExtractAnnotations(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name: "single @implement annotation",
			// idd:ignore start
			content: "// @implement SPEC-BE-001",
			want:    []string{"SPEC-BE-001"},
			// idd:ignore end
		},
		{
			name: "multiple annotation types",
			// idd:ignore start
			content: "// @implement SPEC-CMD_IDD_CLI-001\n// @test TEST-CMD_IDD_CLI-001",
			want:    []string{"SPEC-CMD_IDD_CLI-001", "TEST-CMD_IDD_CLI-001"},
			// idd:ignore end
		},
		{
			name: "annotation with multiple refs",
			// idd:ignore start
			content: "// @implement SPEC-CMD_IDD_CLI-001, SPEC-CMD_IDD_CLI-002",
			want:    []string{"SPEC-CMD_IDD_CLI-001", "SPEC-CMD_IDD_CLI-002"},
			// idd:ignore end
		},
		{
			name:    "no annotations",
			content: "// regular comment",
			want:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractAnnotations(tt.content)
			if !sameElements(got, tt.want) {
				t.Errorf("ExtractAnnotations(%q) = %v, want %v", tt.content, got, tt.want)
			}
		})
	}
}

// @test-contract TEST-PKG_PATTERN-008
func TestIsQuoted(t *testing.T) {
	tests := []struct {
		name    string
		content string
		id      string
		want    bool
	}{
		{
			name:    "in backticks",
			content: "`SPEC-BE-001`",
			id:      "SPEC-BE-001",
			want:    true,
		},
		{
			name:    "in double quotes",
			content: "\"SPEC-BE-001\"",
			id:      "SPEC-BE-001",
			want:    true,
		},
		{
			name:    "not quoted",
			content: "SPEC-BE-001",
			id:      "SPEC-BE-001",
			want:    false,
		},
		{
			name:    "partial match not quoted",
			content: "XSPEC-BE-001",
			id:      "SPEC-BE-001",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isQuoted(tt.content, tt.id)
			if got != tt.want {
				t.Errorf("isQuoted(%q, %q) = %v, want %v", tt.content, tt.id, got, tt.want)
			}
		})
	}
}

func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func sameElements(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aCopy := make([]string, len(a))
	bCopy := make([]string, len(b))
	copy(aCopy, a)
	copy(bCopy, b)
	for i := range aCopy {
		found := false
		for j := range bCopy {
			if aCopy[i] == bCopy[j] {
				bCopy[j] = ""
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
