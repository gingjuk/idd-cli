package collector

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// @test-contract TEST-INTERNAL_COLLECTOR-032
func TestInspectDocumentCompletions_DeduplicatesTargetsAndOverlappingSlots(t *testing.T) {
	projectRoot := t.TempDir()
	packages := []string{"internal/auth", "internal/config"}
	for _, packagePath := range packages {
		if err := os.MkdirAll(filepath.Join(projectRoot, packagePath), 0o755); err != nil {
			t.Fatalf("MkdirAll(%s) error = %v", packagePath, err)
		}
	}
	if _, err := InitDocumentPackages(projectRoot, packages); err != nil {
		t.Fatalf("InitDocumentPackages() error = %v", err)
	}

	authDocs := filepath.Join(projectRoot, "docs", "internal", "auth")
	configDocs := filepath.Join(projectRoot, "docs", "internal", "config")
	targets := []string{
		authDocs,
		filepath.Join(authDocs, "spec.md"),
		configDocs,
		authDocs,
	}
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	relativeAuthDocs, err := filepath.Rel(workingDirectory, authDocs)
	if err != nil {
		t.Fatalf("Rel(%s) error = %v", authDocs, err)
	}
	targets = append(targets, relativeAuthDocs)
	status, err := InspectDocumentCompletions(targets)
	if err != nil {
		t.Fatalf("InspectDocumentCompletions() error = %v", err)
	}
	if len(status.Targets) != 3 {
		t.Fatalf("targets = %v, want three normalized first-occurrence targets", status.Targets)
	}
	if status.Targets[0] != authDocs ||
		status.Targets[1] != filepath.Join(authDocs, "spec.md") ||
		status.Targets[2] != configDocs {
		t.Errorf("targets = %v, want first-occurrence order", status.Targets)
	}
	if len(status.IncompleteSlots) != 8 {
		t.Fatalf("incomplete slots = %d, want 8 without overlap duplicates", len(status.IncompleteSlots))
	}
	seen := make(map[string]bool, len(status.IncompleteSlots))
	for index, slot := range status.IncompleteSlots {
		key := strings.Join([]string{slot.File, slot.Slot, slot.Reason}, "\x00")
		if seen[key] {
			t.Errorf("duplicate incomplete slot: %#v", slot)
		}
		seen[key] = true
		if index > 0 {
			previous := status.IncompleteSlots[index-1]
			if previous.File > slot.File ||
				(previous.File == slot.File && previous.Line > slot.Line) {
				t.Errorf("incomplete slots are not deterministically sorted: %#v then %#v", previous, slot)
			}
		}
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-032
func TestInitDocumentPackages_BatchPreflight(t *testing.T) {
	tests := []struct {
		name       string
		requested  []string
		create     []string
		wantChange int
		wantErr    bool
	}{
		{
			name:       "initializes unique packages once",
			requested:  []string{"internal/auth", "internal/config", "internal/auth"},
			create:     []string{"internal/auth", "internal/config"},
			wantChange: 8,
		},
		{
			name:      "invalid later package prevents earlier writes",
			requested: []string{"internal/auth", "internal/missing"},
			create:    []string{"internal/auth"},
			wantErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			projectRoot := t.TempDir()
			for _, packagePath := range test.create {
				if err := os.MkdirAll(filepath.Join(projectRoot, packagePath), 0o755); err != nil {
					t.Fatalf("MkdirAll(%s) error = %v", packagePath, err)
				}
			}

			changed, err := InitDocumentPackages(projectRoot, test.requested)
			if (err != nil) != test.wantErr {
				t.Fatalf("InitDocumentPackages() error = %v, wantErr %v", err, test.wantErr)
			}
			if len(changed) != test.wantChange {
				t.Errorf("changed = %v, want %d paths", changed, test.wantChange)
			}
			if test.wantErr {
				authDocs := filepath.Join(projectRoot, "docs", "internal", "auth")
				if _, statErr := os.Stat(authDocs); !os.IsNotExist(statErr) {
					t.Errorf("valid earlier package was written before batch failure: %v", statErr)
				}
			}
		})
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-032
func TestRepairDocumentTargets_DeduplicatesAndPreflightsBatch(t *testing.T) {
	tests := []struct {
		name    string
		invalid bool
	}{
		{name: "overlapping targets write once"},
		{name: "invalid later target prevents earlier repair", invalid: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			projectRoot := t.TempDir()
			authDocs := writeValidDocumentSet(t, projectRoot, "internal/auth")
			authDesign := filepath.Join(authDocs, "design.md")
			replaceTestFile(t, authDesign, `version: "1.1"`, `version: ""`)
			before := readTestFile(t, authDesign)

			workingDirectory, getwdErr := os.Getwd()
			if getwdErr != nil {
				t.Fatalf("Getwd() error = %v", getwdErr)
			}
			relativeAuthDocs, relErr := filepath.Rel(workingDirectory, authDocs)
			if relErr != nil {
				t.Fatalf("Rel(%s) error = %v", authDocs, relErr)
			}
			targets := []string{relativeAuthDocs, authDesign, authDocs}
			if test.invalid {
				configDocs := writeValidDocumentSet(t, projectRoot, "internal/config")
				writeTestFile(t, filepath.Join(configDocs, "idd.yaml"), "version: \"1.0\"\n")
				targets = append(targets, configDocs)
			}

			changed, err := RepairDocumentTargets(targets)
			if test.invalid {
				if err == nil || !strings.Contains(err.Error(), "central catalog") {
					t.Fatalf("RepairDocumentTargets() error = %v, want central catalog refusal", err)
				}
				if len(changed) != 0 {
					t.Errorf("changed = %v, want no writes", changed)
				}
				if got := readTestFile(t, authDesign); got != before {
					t.Error("earlier valid repair target changed before later preflight failure")
				}
				return
			}

			if err != nil {
				t.Fatalf("RepairDocumentTargets() error = %v", err)
			}
			if len(changed) != 1 ||
				documentPathIdentity(changed[0]) != documentPathIdentity(authDesign) {
				t.Errorf("changed = %v, want design.md once", changed)
			}
		})
	}
}

func writeValidDocumentSet(t *testing.T, projectRoot, packagePath string) string {
	t.Helper()
	docsDir := filepath.Join(projectRoot, "docs", filepath.FromSlash(packagePath))
	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%s) error = %v", docsDir, err)
	}
	for filename, content := range map[string]string{
		"design.md":   validDesignDocument,
		"contract.md": validContractDocument,
		"spec.md":     validSpecDocument,
		"testing.md":  validTestingDocument,
	} {
		content = strings.ReplaceAll(content, "internal/auth", packagePath)
		writeTestFile(t, filepath.Join(docsDir, filename), content)
	}
	return docsDir
}
