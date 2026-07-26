// Package collector provides tests for generated document completion tracking.

// Spec: docs/internal/collector/spec.md
// Test: docs/internal/collector/testing.md
package collector

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// @test-contract TEST-INTERNAL_COLLECTOR-021
func TestRenderIDDDocumentTemplate_ListsSchemaRequiredFields(t *testing.T) {
	tests := []struct {
		role string
		want []string
	}{
		{"design", []string{"Required fields: Purpose."}},
		{"contract", []string{"Required fields: Guarantees."}},
		{
			"spec",
			[]string{"Required fields: Design, Contract, Requirement, Acceptance."},
		},
		{
			"testing",
			[]string{
				"Required fields: Kind, Covers, Purpose, Oracle.",
				"Conditionally required: Contracts when Kind is `contract`.",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.role, func(t *testing.T) {
			rendered := renderIDDDocumentTemplate("internal/example", test.role)
			if !strings.Contains(rendered, `<!-- idd:scaffold slot="`) {
				t.Fatalf("template lacks scaffold marker:\n%s", rendered)
			}
			for _, want := range test.want {
				if !strings.Contains(rendered, want) {
					t.Errorf("template lacks %q:\n%s", want, rendered)
				}
			}
		})
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-021
func TestInspectDocumentCompletion_GeneratedSlots(t *testing.T) {
	projectRoot := t.TempDir()
	packagePath := "internal/example"
	if err := os.MkdirAll(filepath.Join(projectRoot, packagePath), 0o755); err != nil {
		t.Fatalf("MkdirAll(package) error = %v", err)
	}
	if _, err := InitDocuments(projectRoot, packagePath); err != nil {
		t.Fatalf("InitDocuments() error = %v", err)
	}

	docsDir := filepath.Join(projectRoot, "docs", packagePath)
	status, err := InspectDocumentCompletion(docsDir)
	if err != nil {
		t.Fatalf("InspectDocumentCompletion() error = %v", err)
	}
	if status.Status != "incomplete" {
		t.Fatalf("status = %q, want incomplete", status.Status)
	}

	wantSlots := []string{
		"contract.records",
		"design.architecture",
		"design.components",
		"design.dependencies",
		"design.function-composition",
		"design.package-layout",
		"design.testability-hooks",
		"spec.records",
		"testing.records",
	}
	gotSlots := make([]string, 0, len(status.IncompleteSlots))
	for _, slot := range status.IncompleteSlots {
		gotSlots = append(gotSlots, slot.Slot)
		if slot.File == "" || slot.Line <= 0 || slot.Role == "" || slot.Reason == "" {
			t.Errorf("incomplete slot lacks source detail: %#v", slot)
		}
	}
	sort.Strings(gotSlots)
	sort.Strings(wantSlots)
	if strings.Join(gotSlots, ",") != strings.Join(wantSlots, ",") {
		t.Errorf("slots = %v, want %v", gotSlots, wantSlots)
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-021
func TestInspectDocumentCompletion_TreeSkipsLegacyRoleFiles(t *testing.T) {
	projectRoot := t.TempDir()
	packagePath := "internal/example"
	if err := os.MkdirAll(filepath.Join(projectRoot, packagePath), 0o755); err != nil {
		t.Fatalf("MkdirAll(package) error = %v", err)
	}
	if _, err := InitDocuments(projectRoot, packagePath); err != nil {
		t.Fatalf("InitDocuments() error = %v", err)
	}

	legacyDir := filepath.Join(projectRoot, "docs", "architecture")
	if err := os.MkdirAll(legacyDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(legacy) error = %v", err)
	}
	if err := os.WriteFile(
		filepath.Join(legacyDir, "design.md"),
		[]byte("# Legacy architecture\n"),
		0o644,
	); err != nil {
		t.Fatalf("WriteFile(legacy) error = %v", err)
	}

	status, err := InspectDocumentCompletion(filepath.Join(projectRoot, "docs"))
	if err != nil {
		t.Fatalf("InspectDocumentCompletion() error = %v", err)
	}
	if status.Status != "incomplete" || len(status.IncompleteSlots) != 9 {
		t.Fatalf("status = %#v, want only the generated package's nine slots", status)
	}
	for _, slot := range status.IncompleteSlots {
		if strings.Contains(slot.File, "architecture") {
			t.Errorf("legacy role file became a completion target: %#v", slot)
		}
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-021
func TestInspectDocumentCompletion_MarkerRemovalDoesNotCompleteEmptySlots(t *testing.T) {
	projectRoot := t.TempDir()
	packagePath := "internal/example"
	if err := os.MkdirAll(filepath.Join(projectRoot, packagePath), 0o755); err != nil {
		t.Fatalf("MkdirAll(package) error = %v", err)
	}
	if _, err := InitDocuments(projectRoot, packagePath); err != nil {
		t.Fatalf("InitDocuments() error = %v", err)
	}

	docsDir := filepath.Join(projectRoot, "docs", packagePath)
	for filename := range iddDocumentRoles {
		path := filepath.Join(docsDir, filename)
		content := readTestFile(t, path)
		var filtered []string
		for _, line := range strings.Split(content, "\n") {
			if strings.Contains(line, "idd:scaffold") {
				continue
			}
			filtered = append(filtered, line)
		}
		if err := os.WriteFile(path, []byte(strings.Join(filtered, "\n")), 0o644); err != nil {
			t.Fatalf("WriteFile(%s) error = %v", filename, err)
		}
	}

	status, err := InspectDocumentCompletion(docsDir)
	if err != nil {
		t.Fatalf("InspectDocumentCompletion() error = %v", err)
	}
	if status.Status != "incomplete" || len(status.IncompleteSlots) != 9 {
		t.Fatalf("status after deleting markers = %#v, want nine incomplete slots", status)
	}
	for _, slot := range status.IncompleteSlots {
		if strings.Contains(slot.Reason, "marker remains") {
			t.Errorf("slot %s still reports a removed marker: %q", slot.Slot, slot.Reason)
		}
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-021
func TestInspectDocumentCompletion_IgnoresScaffoldExamplesInCodeFences(t *testing.T) {
	docsDir := createIDDDocumentFixture(t)
	specPath := filepath.Join(docsDir, "spec.md")
	data, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, []byte(`

### Scaffold protocol example

`+"```markdown"+`
<!-- idd:scaffold slot="spec.example-only" -->
`+"```"+`
`)...)
	if err := os.WriteFile(specPath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	status, err := InspectDocumentCompletion(specPath)
	if err != nil {
		t.Fatalf("InspectDocumentCompletion() error = %v", err)
	}
	if status.Status != "complete" || len(status.IncompleteSlots) != 0 {
		t.Fatalf("status = %#v, want complete; fenced marker examples are prose", status)
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-021
func TestInspectDocumentCompletion_CompletedDocuments(t *testing.T) {
	projectRoot := t.TempDir()
	docsDir := filepath.Join(projectRoot, "docs", "internal", "example")
	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(docs) error = %v", err)
	}

	bodies := map[string]string{
		"design.md": `
# Design: internal/example

## Component: Runner

**Purpose:** Runner owns execution state and isolates dependency failures.

## Architecture

Requests flow through a normalized model before execution.

## Package Layout

Collection and validation remain separate responsibilities.

## Function Composition

Initialization precedes collection, graph construction, and reporting.

## Dependencies

The package depends only on normalized configuration and models.

## Testability Hooks

Parsers and filesystem operations are exercised through temporary fixtures.
`,
		"contract.md": `
# Contracts: internal/example

## Contract: Runner

**Guarantees:** Callers provide a project path and receive deterministic
validation findings.
`,
		"spec.md": `
# Specifications: internal/example

## SPEC-INTERNAL_EXAMPLE-001: Deterministic execution

- **Design:** ` + "`Runner`" + `
- **Contract:** ` + "`Runner`" + `

**Requirement:**

Execution returns the same ordered findings for the same project state.

**Acceptance:**

Two runs over the same fixture return byte-for-byte equivalent ordered
findings.
`,
		"testing.md": `
# Testing: internal/example

## TEST-INTERNAL_EXAMPLE-001: Deterministic execution evidence

- **Kind:** ` + "`test`" + `
- **Covers:** ` + "`SPEC-INTERNAL_EXAMPLE-001`" + `

**Purpose:**

Repeated fixture runs compare exact ordered findings.

**Oracle:**

The ordered finding list from the second run exactly equals the first.
`,
	}
	for filename, role := range iddDocumentRoles {
		data, err := MarshalIDDDocument(newIDDDocument("internal/example", role), []byte(bodies[filename]))
		if err != nil {
			t.Fatalf("MarshalIDDDocument(%s) error = %v", filename, err)
		}
		if err := os.WriteFile(filepath.Join(docsDir, filename), data, 0o644); err != nil {
			t.Fatalf("WriteFile(%s) error = %v", filename, err)
		}
	}

	status, err := InspectDocumentCompletion(docsDir)
	if err != nil {
		t.Fatalf("InspectDocumentCompletion() error = %v", err)
	}
	if status.Status != "complete" || len(status.IncompleteSlots) != 0 {
		t.Fatalf("completed status = %#v", status)
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-021
func TestInspectDocumentCompletion_ReportsIncompleteRecordFields(t *testing.T) {
	tests := []struct {
		name       string
		role       string
		filename   string
		body       string
		wantSlot   string
		wantReason string
	}{
		{
			name:     "SPEC Acceptance is absent",
			role:     "spec",
			filename: "spec.md",
			body: `
# Specifications: internal/example

## SPEC-INTERNAL_EXAMPLE-001: Deterministic execution

- **Design:** ` + "`Runner`" + `
- **Contract:** ` + "`Runner`" + `

**Requirement:** Execution is deterministic.
`,
			wantSlot:   "spec.spec-internal-example-001.acceptance",
			wantReason: "Acceptance",
		},
		{
			name:     "contract TEST Contracts is absent",
			role:     "testing",
			filename: "testing.md",
			body: `
# Testing: internal/example

## TEST-INTERNAL_EXAMPLE-001: Runner boundary

- **Kind:** ` + "`contract`" + `
- **Covers:** ` + "`SPEC-INTERNAL_EXAMPLE-001`" + `

**Purpose:** Protect the Runner boundary.

**Oracle:** Every implementation returns the same documented result.
`,
			wantSlot:   "testing.test-internal-example-001.contracts",
			wantReason: "Contracts",
		},
		{
			name:     "TEST Oracle is a placeholder",
			role:     "testing",
			filename: "testing.md",
			body: `
# Testing: internal/example

## TEST-INTERNAL_EXAMPLE-001: Runner behavior

- **Kind:** ` + "`test`" + `
- **Covers:** ` + "`SPEC-INTERNAL_EXAMPLE-001`" + `

**Purpose:** Exercise deterministic execution.

**Oracle:** TBD
`,
			wantSlot:   "testing.test-internal-example-001.oracle",
			wantReason: "Oracle",
		},
		{
			name:     "TEST Covers contains a placeholder item",
			role:     "testing",
			filename: "testing.md",
			body: `
# Testing: internal/example

## TEST-INTERNAL_EXAMPLE-001: Runner behavior

- **Kind:** ` + "`test`" + `
- **Covers:** ` + "`SPEC-INTERNAL_EXAMPLE-001`" + `, TBD

**Purpose:** Exercise deterministic execution.

**Oracle:** The observed result exactly matches the scenario expectation.
`,
			wantSlot:   "testing.test-internal-example-001.covers",
			wantReason: "Covers",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			projectRoot := t.TempDir()
			docsDir := filepath.Join(projectRoot, "docs", "internal", "example")
			if err := os.MkdirAll(docsDir, 0o755); err != nil {
				t.Fatalf("MkdirAll(docs) error = %v", err)
			}
			data, err := MarshalIDDDocument(
				newIDDDocument("internal/example", test.role),
				[]byte(test.body),
			)
			if err != nil {
				t.Fatalf("MarshalIDDDocument(%s) error = %v", test.role, err)
			}
			path := filepath.Join(docsDir, test.filename)
			if err := os.WriteFile(path, data, 0o644); err != nil {
				t.Fatalf("WriteFile(%s) error = %v", test.filename, err)
			}

			status, err := InspectDocumentCompletion(path)
			if err != nil {
				t.Fatalf("InspectDocumentCompletion(%s) error = %v", test.filename, err)
			}
			if status.Status != "incomplete" || len(status.IncompleteSlots) != 1 {
				t.Fatalf("status = %#v, want one incomplete record field", status)
			}
			slot := status.IncompleteSlots[0]
			if slot.Slot != test.wantSlot ||
				slot.Role != test.role ||
				slot.Line <= 0 ||
				!strings.Contains(slot.Reason, test.wantReason) {
				t.Errorf("incomplete field slot = %#v", slot)
			}
		})
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-021
func TestRepairDocuments_PreservesScaffoldMarkers(t *testing.T) {
	projectRoot := t.TempDir()
	packagePath := "internal/example"
	if err := os.MkdirAll(filepath.Join(projectRoot, packagePath), 0o755); err != nil {
		t.Fatalf("MkdirAll(package) error = %v", err)
	}
	if _, err := InitDocuments(projectRoot, packagePath); err != nil {
		t.Fatalf("InitDocuments() error = %v", err)
	}

	docsDir := filepath.Join(projectRoot, "docs", packagePath)
	before := readTestFile(t, filepath.Join(docsDir, "design.md"))
	if _, err := RepairDocuments(docsDir); err != nil {
		t.Fatalf("RepairDocuments() error = %v", err)
	}
	after := readTestFile(t, filepath.Join(docsDir, "design.md"))
	if before != after {
		t.Errorf("RepairDocuments() changed generated body:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if !strings.Contains(after, `<!-- idd:scaffold slot="design.architecture" -->`) {
		t.Error("RepairDocuments() removed scaffold marker")
	}
}
