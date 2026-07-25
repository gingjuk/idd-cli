// Package collector provides tests for self-describing IDD documents.

// Spec: docs/internal/collector/spec.md
// Test: docs/internal/collector/testing.md
package collector

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jingxu9x/idd-cli/internal/config"
)

const validDesignDocument = `---
idd:
  version: "1.0"
  package: internal/auth
  document: design
---

# Design

## Component: AuthModule

AuthModule owns the authentication boundary and keeps transport concerns out of
credential verification.

## Architecture

AuthModule separates authentication from transport concerns.

## Package Layout

One package owns the boundary.

## Function Composition

The handler calls the authenticator.

## Dependencies

Dependencies are injected.

## Testability Hooks

Tests use deterministic fakes.
`

const validContractDocument = `---
idd:
  version: "1.0"
  package: internal/auth
  document: contract
---

# Contracts

## Contract: Authenticator

Authenticator exposes the authentication boundary.
`

const validSpecDocument = `---
idd:
  version: "1.0"
  package: internal/auth
  document: spec
---

# Specifications

## SPEC-INTERNAL_AUTH-001: User authentication

- **Design:** ` + "`AuthModule`" + `
- **Contract:** ` + "`Authenticator`" + `

**Requirement:** Authenticate a user with credentials while exposing one stable
failure for invalid credentials.

Invalid credentials expose one stable failure.
`

const validTestingDocument = `---
idd:
  version: "1.0"
  package: internal/auth
  document: testing
---

# Testing

## TEST-INTERNAL_AUTH-001: Authentication behavior

- **Kind:** ` + "`test`" + `
- **Covers:** ` + "`SPEC-INTERNAL_AUTH-001`" + `

**Purpose:** Verify valid and invalid credentials.

### Scenarios

- Accepted credentials
- Rejected credentials
`

// @test TEST-INTERNAL_COLLECTOR-021
func TestParseIDDDocument(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		wantDocument string
		wantRecords  int
		wantBody     string
		wantNil      bool
		wantErr      bool
	}{
		{
			name:         "human readable spec record",
			content:      validSpecDocument,
			wantDocument: "spec",
			wantRecords:  1,
			wantBody:     "# Specifications",
		},
		{
			name: "multiline testing narrative remains hand editable",
			content: `---
idd:
  version: "1.0"
  package: internal/auth
  document: testing
---

# Testing

## TEST-INTERNAL_AUTH-001: Authentication behavior

- **Kind:** ` + "`contract`" + `
- **Covers:** ` + "`SPEC-INTERNAL_AUTH-001`" + `

**Purpose:** Verify valid and invalid credentials across every implementation.

### Fixture

` + "```go" + `
func newAuthenticator() Authenticator
` + "```" + `
`,
			wantDocument: "testing",
			wantRecords:  1,
			wantBody:     "# Testing",
		},
		{
			name: "legacy frontmatter is not claimed",
			content: `---
markers:
  - id: SPEC-INTERNAL_AUTH-001
    name: Authentication
---

# Specifications
`,
			wantNil: true,
		},
		{
			name: "semantic yaml records are rejected",
			content: `---
idd:
  version: "1.0"
  package: internal/auth
  document: spec
  specs: []
---

# Specifications
`,
			wantErr: true,
		},
		{
			name: "unknown idd field is rejected",
			content: `---
idd:
  version: "1.0"
  package: internal/auth
  document: spec
  unknown: true
---

# Specifications
`,
			wantErr: true,
		},
		{
			name: "malformed idd yaml is rejected",
			content: `---
idd:
  version: [
---

# Specifications
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			document, body, err := ParseIDDDocument([]byte(tt.content))
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseIDDDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if tt.wantNil {
				if document != nil {
					t.Fatalf("ParseIDDDocument() = %#v, want nil", document)
				}
				return
			}
			if document == nil {
				t.Fatal("ParseIDDDocument() = nil, want document")
			}
			if document.Document != tt.wantDocument {
				t.Errorf("Document = %q, want %q", document.Document, tt.wantDocument)
			}
			recordCount := len(document.Specs) + len(document.Tests)
			if recordCount != tt.wantRecords {
				t.Errorf("record count = %d, want %d", recordCount, tt.wantRecords)
			}
			if !strings.Contains(string(body), tt.wantBody) {
				t.Errorf("body = %q, want content %q", body, tt.wantBody)
			}
		})
	}
}

// @test TEST-INTERNAL_COLLECTOR-021
func TestParseIDDDocument_MarkdownRecords(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		wantComponents []string
		wantContracts  []string
		wantSpec       *IDDDocumentSpec
		wantTest       *IDDDocumentTest
	}{
		{
			name:           "component heading owns design declaration",
			content:        validDesignDocument,
			wantComponents: []string{"AuthModule"},
		},
		{
			name:          "contract heading owns contract declaration",
			content:       validContractDocument,
			wantContracts: []string{"Authenticator"},
		},
		{
			name:    "spec fields and wrapped requirement come from record block",
			content: validSpecDocument,
			wantSpec: &IDDDocumentSpec{
				ID:          "SPEC-INTERNAL_AUTH-001",
				Title:       "User authentication",
				Requirement: "Authenticate a user with credentials while exposing one stable failure for invalid credentials.",
				Design:      "AuthModule",
				Contract:    "Authenticator",
			},
		},
		{
			name:    "test fields and purpose come from record block",
			content: validTestingDocument,
			wantTest: &IDDDocumentTest{
				ID:      "TEST-INTERNAL_AUTH-001",
				Title:   "Authentication behavior",
				Purpose: "Verify valid and invalid credentials.",
				Kind:    "test",
				Covers:  []string{"SPEC-INTERNAL_AUTH-001"},
			},
		},
		{
			name: "record-like content in code fence is ignored",
			content: `---
idd:
  version: "1.0"
  package: internal/auth
  document: spec
---

# Specifications

` + "```markdown" + `
## SPEC-INTERNAL_AUTH-999: Not a declaration

- **Design:** ` + "`Fake`" + `
- **Contract:** ` + "`Fake`" + `

**Requirement:** Not a real requirement.
` + "```" + `
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			document, _, err := ParseIDDDocument([]byte(tt.content))
			if err != nil {
				t.Fatalf("ParseIDDDocument() error = %v", err)
			}
			if document == nil {
				t.Fatal("ParseIDDDocument() = nil, want document")
			}
			if got := strings.Join(document.Components, ","); got != strings.Join(tt.wantComponents, ",") {
				t.Errorf("Components = %q, want %q", got, strings.Join(tt.wantComponents, ","))
			}
			if got := strings.Join(document.Contracts, ","); got != strings.Join(tt.wantContracts, ",") {
				t.Errorf("Contracts = %q, want %q", got, strings.Join(tt.wantContracts, ","))
			}
			if tt.wantSpec != nil {
				if len(document.Specs) != 1 || document.Specs[0] != *tt.wantSpec {
					t.Errorf("Specs = %#v, want %#v", document.Specs, *tt.wantSpec)
				}
			} else if len(document.Specs) != 0 {
				t.Errorf("Specs = %#v, want none", document.Specs)
			}
			if tt.wantTest != nil {
				if len(document.Tests) != 1 {
					t.Fatalf("Tests = %#v, want one", document.Tests)
				}
				got := document.Tests[0]
				if got.ID != tt.wantTest.ID ||
					got.Title != tt.wantTest.Title ||
					got.Purpose != tt.wantTest.Purpose ||
					got.Kind != tt.wantTest.Kind ||
					strings.Join(got.Covers, ",") != strings.Join(tt.wantTest.Covers, ",") {
					t.Errorf("Tests[0] = %#v, want %#v", got, *tt.wantTest)
				}
			} else if len(document.Tests) != 0 {
				t.Errorf("Tests = %#v, want none", document.Tests)
			}
		})
	}
}

// @test TEST-INTERNAL_COLLECTOR-021
func TestDocCollector_IDDDocumentValidation(t *testing.T) {
	tests := []struct {
		name       string
		mutate     func(*testing.T, string)
		wantRule   string
		wantCode   string
		wantSource string
	}{
		{name: "valid self-describing document set"},
		{
			name: "package follows docs path",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(t, filepath.Join(docsDir, "spec.md"), "package: internal/auth", "package: internal/wrong")
			},
			wantRule:   "idd-document-identity",
			wantCode:   "package",
			wantSource: "spec.md:",
		},
		{
			name: "role follows fixed filename",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(t, filepath.Join(docsDir, "testing.md"), "document: testing", "document: spec")
			},
			wantRule:   "idd-document-identity",
			wantCode:   "document",
			wantSource: "testing.md:",
		},
		{
			name: "unknown field reports exact Markdown line",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(
					t,
					filepath.Join(docsDir, "spec.md"),
					"  document: spec",
					"  unknown_field: true\n  document: spec",
				)
			},
			wantRule:   "idd-document-parse",
			wantCode:   "spec",
			wantSource: "spec.md:5",
		},
		{
			name: "spec requirement is required",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(
					t,
					filepath.Join(docsDir, "spec.md"),
					"**Requirement:** Authenticate a user with credentials while exposing one stable\nfailure for invalid credentials.\n",
					"",
				)
			},
			wantRule:   "idd-document-schema",
			wantCode:   "requirement",
			wantSource: "spec.md:",
		},
		{
			name: "spec relationship fields use the compact list",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(
					t,
					filepath.Join(docsDir, "spec.md"),
					"- **Design:** `AuthModule`\n",
					"",
				)
			},
			wantRule:   "idd-document-schema",
			wantCode:   "design",
			wantSource: "spec.md:",
		},
		{
			name: "large table is not a canonical spec declaration",
			mutate: func(t *testing.T, docsDir string) {
				writeTestFile(t, filepath.Join(docsDir, "spec.md"), `---
idd:
  version: "1.0"
  package: internal/auth
  document: spec
---

# Specifications

| ID | Title | Design | Contract | Requirement |
| --- | --- | --- | --- | --- |
| SPEC-INTERNAL_AUTH-001 | User authentication | AuthModule | Authenticator | Authenticate credentials |
`)
			},
			wantRule:   "idd-document-markdown",
			wantCode:   "record-format",
			wantSource: "spec.md:",
		},
		{
			name: "spec id uses package-derived module",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(t, filepath.Join(docsDir, "spec.md"), "SPEC-INTERNAL_AUTH-001", "SPEC-WRONG-001")
			},
			wantRule:   "idd-document-schema",
			wantCode:   "id",
			wantSource: "spec.md:",
		},
		{
			name: "placeholder requirement is rejected",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(
					t,
					filepath.Join(docsDir, "spec.md"),
					"Authenticate a user with credentials while exposing one stable\nfailure for invalid credentials.",
					"TBD",
				)
			},
			wantRule:   "idd-document-schema",
			wantCode:   "requirement",
			wantSource: "spec.md:",
		},
		{
			name: "spec design resolves in design document",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(t, filepath.Join(docsDir, "spec.md"), "`AuthModule`", "`MissingModule`")
			},
			wantRule:   "idd-document-reference",
			wantCode:   "design",
			wantSource: "spec.md:",
		},
		{
			name: "declared component has narrative definition",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(
					t,
					filepath.Join(docsDir, "design.md"),
					"AuthModule owns the authentication boundary and keeps transport concerns out of\ncredential verification.",
					"",
				)
			},
			wantRule:   "idd-document-schema",
			wantCode:   "component",
			wantSource: "design.md:",
		},
		{
			name: "declared contract has narrative definition",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(
					t,
					filepath.Join(docsDir, "contract.md"),
					"Authenticator exposes the authentication boundary.",
					"",
				)
			},
			wantRule:   "idd-document-schema",
			wantCode:   "contract",
			wantSource: "contract.md:",
		},
		{
			name: "test kind is constrained",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(t, filepath.Join(docsDir, "testing.md"), "`test`", "`e2e`")
			},
			wantRule:   "idd-document-schema",
			wantCode:   "kind",
			wantSource: "testing.md:",
		},
		{
			name: "test coverage is required",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(
					t,
					filepath.Join(docsDir, "testing.md"),
					"- **Covers:** `SPEC-INTERNAL_AUTH-001`\n",
					"",
				)
			},
			wantRule:   "idd-document-schema",
			wantCode:   "covers",
			wantSource: "testing.md:",
		},
		{
			name: "test coverage has no duplicate identifiers",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(
					t,
					filepath.Join(docsDir, "testing.md"),
					"`SPEC-INTERNAL_AUTH-001`",
					"`SPEC-INTERNAL_AUTH-001`, `SPEC-INTERNAL_AUTH-001`",
				)
			},
			wantRule:   "idd-document-schema",
			wantCode:   "covers",
			wantSource: "testing.md:",
		},
		{
			name: "coverage resolves in spec document",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(t, filepath.Join(docsDir, "testing.md"), "SPEC-INTERNAL_AUTH-001", "SPEC-INTERNAL_AUTH-999")
			},
			wantRule:   "idd-document-reference",
			wantCode:   "covers",
			wantSource: "testing.md:",
		},
		{
			name: "document set is complete",
			mutate: func(t *testing.T, docsDir string) {
				if err := os.Remove(filepath.Join(docsDir, "contract.md")); err != nil {
					t.Fatalf("Remove(contract.md) error = %v", err)
				}
			},
			wantRule:   "idd-document-set",
			wantCode:   "contract",
			wantSource: "contract.md:",
		},
		{
			name: "semantic records are not stored in frontmatter",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(
					t,
					filepath.Join(docsDir, "spec.md"),
					"  document: spec",
					"  document: spec\n  specs: []",
				)
			},
			wantRule:   "idd-document-migration",
			wantCode:   "yaml-semantics",
			wantSource: "spec.md:",
		},
		{
			name: "contract declarations use contract headings",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(
					t,
					filepath.Join(docsDir, "contract.md"),
					"## Contract: Authenticator",
					"## Authenticator",
				)
			},
			wantRule:   "idd-document-reference",
			wantCode:   "contract",
			wantSource: "spec.md:",
		},
		{
			name: "legacy metadata is not mixed with idd metadata",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(
					t,
					filepath.Join(docsDir, "spec.md"),
					"idd:\n",
					"markers: [{id: SPEC-INTERNAL_AUTH-001, name: Duplicate}]\nidd:\n",
				)
			},
			wantRule:   "idd-document-markdown",
			wantCode:   "duplicate-metadata",
			wantSource: "spec.md:",
		},
		{
			name: "central marker catalog requires migration",
			mutate: func(t *testing.T, docsDir string) {
				writeTestFile(t, filepath.Join(docsDir, "idd.yaml"), "version: \"1.0\"\n")
			},
			wantRule:   "idd-document-migration",
			wantCode:   "central-catalog",
			wantSource: "idd.yaml:1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			docsDir := createIDDDocumentFixture(t)
			if tt.mutate != nil {
				tt.mutate(t, docsDir)
			}

			_, errs, err := NewDocCollector(config.Default()).Collect(
				context.Background(),
				filepath.Dir(filepath.Dir(docsDir)),
			)
			if err != nil {
				t.Fatalf("Collect() error = %v", err)
			}
			if tt.wantRule == "" {
				if len(errs) != 0 {
					t.Fatalf("Collect() errors = %v, want none", errs)
				}
				return
			}
			for _, validationErr := range errs {
				if validationErr.Rule != tt.wantRule || validationErr.Code != tt.wantCode {
					continue
				}
				if !strings.Contains(validationErr.Source, tt.wantSource) {
					t.Errorf("Source = %q, want substring %q", validationErr.Source, tt.wantSource)
				}
				return
			}
			t.Fatalf("Collect() errors = %v, want rule %q code %q", errs, tt.wantRule, tt.wantCode)
		})
	}
}

// @test TEST-INTERNAL_COLLECTOR-021
func TestDocCollector_IDDDocumentsBuildDerivedLinks(t *testing.T) {
	docsDir := createIDDDocumentFixture(t)
	collector := NewDocCollector(config.Default())
	set, errs, err := collector.Collect(context.Background(), filepath.Dir(filepath.Dir(docsDir)))
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if len(errs) != 0 {
		t.Fatalf("Collect() errors = %v, want none", errs)
	}

	spec, ok := set.Get("SPEC-INTERNAL_AUTH-001")
	if !ok {
		t.Fatal("SPEC-INTERNAL_AUTH-001 not collected")
	}
	test, ok := set.Get("TEST-INTERNAL_AUTH-001")
	if !ok {
		t.Fatal("TEST-INTERNAL_AUTH-001 not collected")
	}
	if got := strings.Join(spec.Links, ","); got != "TEST-INTERNAL_AUTH-001" {
		t.Errorf("SPEC links = %q, want derived TEST backlink", got)
	}
	if got := strings.Join(test.Links, ","); got != "SPEC-INTERNAL_AUTH-001" {
		t.Errorf("TEST links = %q, want declared SPEC coverage", got)
	}
	if test.Kind != "test" {
		t.Errorf("TEST kind = %q, want test", test.Kind)
	}
	if !strings.HasSuffix(spec.Source, "spec.md") || spec.Line == 0 {
		t.Errorf("SPEC source = %q:%d, want spec.md Markdown record", spec.Source, spec.Line)
	}
	if !strings.HasSuffix(test.Source, "testing.md") || test.Line == 0 {
		t.Errorf("TEST source = %q:%d, want testing.md Markdown record", test.Source, test.Line)
	}
}

// @test TEST-INTERNAL_COLLECTOR-021
func TestDocCollector_IDDDocumentTargetIncludesSiblings(t *testing.T) {
	docsDir := createIDDDocumentFixture(t)
	set, errs, err := NewDocCollector(config.Default()).Collect(
		context.Background(),
		filepath.Join(docsDir, "spec.md"),
	)
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if len(errs) != 0 {
		t.Fatalf("Collect() errors = %v, want none", errs)
	}
	if !set.Has("SPEC-INTERNAL_AUTH-001") || !set.Has("TEST-INTERNAL_AUTH-001") {
		t.Fatalf("Collect(spec.md) identifiers = %v, want SPEC and sibling TEST", set.AllIdentifiers())
	}
}

// @test TEST-INTERNAL_COLLECTOR-021
func TestMarshalIDDDocument(t *testing.T) {
	document, body, err := ParseIDDDocument([]byte(validTestingDocument))
	if err != nil {
		t.Fatalf("ParseIDDDocument() error = %v", err)
	}
	data, err := MarshalIDDDocument(document, body)
	if err != nil {
		t.Fatalf("MarshalIDDDocument() error = %v", err)
	}
	output := string(data)
	if strings.Contains(output, "tests:") || strings.Contains(output, "covers:") {
		t.Errorf("MarshalIDDDocument() copied semantic records into frontmatter:\n%s", output)
	}
	if !strings.Contains(output, "## TEST-INTERNAL_AUTH-001: Authentication behavior") {
		t.Errorf("MarshalIDDDocument() lost Markdown TEST record:\n%s", output)
	}

	roundTrip, roundTripBody, err := ParseIDDDocument(data)
	if err != nil {
		t.Fatalf("round-trip ParseIDDDocument() error = %v", err)
	}
	if len(roundTrip.Tests) != 1 || roundTrip.Tests[0].ID != document.Tests[0].ID {
		t.Errorf("round-trip TEST changed: %#v", roundTrip.Tests)
	}
	if string(roundTripBody) != string(body) {
		t.Errorf("round-trip body changed:\nwant %q\ngot  %q", body, roundTripBody)
	}
}

// @test TEST-INTERNAL_COLLECTOR-021
func TestInitDocuments(t *testing.T) {
	tests := []struct {
		name        string
		packagePath string
		createPkg   bool
		wantErr     bool
	}{
		{name: "creates four self-describing documents", packagePath: "internal/auth", createPkg: true},
		{name: "rejects path traversal", packagePath: "../outside", wantErr: true},
		{name: "requires an existing package", packagePath: "internal/missing", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectRoot := t.TempDir()
			if tt.createPkg {
				if err := os.MkdirAll(filepath.Join(projectRoot, filepath.FromSlash(tt.packagePath)), 0o755); err != nil {
					t.Fatalf("MkdirAll(package) error = %v", err)
				}
			}

			changed, err := InitDocuments(projectRoot, tt.packagePath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("InitDocuments() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if len(changed) != 4 {
				t.Fatalf("InitDocuments() changed %d files, want 4: %v", len(changed), changed)
			}

			docsDir := filepath.Join(projectRoot, "docs", filepath.FromSlash(tt.packagePath))
			guidanceByRole := map[string]string{
				"design":   "## Component: <name>",
				"contract": "## Contract: <name>",
				"spec":     "## SPEC-<MODULE>-<NUMBER>: <title>",
				"testing":  "## TEST-<MODULE>-<NUMBER>: <title>",
			}
			for filename, role := range iddDocumentRoles {
				document, body, parseErr := ParseIDDDocument([]byte(readTestFile(t, filepath.Join(docsDir, filename))))
				if parseErr != nil || document == nil {
					t.Fatalf("generated %s cannot be parsed: document=%#v error=%v", filename, document, parseErr)
				}
				if document.Package != tt.packagePath || document.Document != role {
					t.Errorf("%s identity = package %q document %q", filename, document.Package, document.Document)
				}
				if !strings.Contains(string(body), guidanceByRole[role]) {
					t.Errorf("%s body does not explain its canonical record: %q", filename, body)
				}
			}
			if _, statErr := os.Stat(filepath.Join(docsDir, "idd.yaml")); !os.IsNotExist(statErr) {
				t.Errorf("central idd.yaml exists after initialization: %v", statErr)
			}
			if _, secondErr := InitDocuments(projectRoot, tt.packagePath); secondErr == nil {
				t.Fatal("second InitDocuments() should refuse to overwrite initialized documents")
			}
		})
	}
}

// @test TEST-INTERNAL_COLLECTOR-021
func TestInitDocuments_PrependsPlainNarrative(t *testing.T) {
	projectRoot := t.TempDir()
	packagePath := "internal/auth"
	if err := os.MkdirAll(filepath.Join(projectRoot, packagePath), 0o755); err != nil {
		t.Fatalf("MkdirAll(package) error = %v", err)
	}
	docsDir := filepath.Join(projectRoot, "docs", packagePath)
	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(docs) error = %v", err)
	}
	plainBody := "# Design\n\nA package-specific design narrative.\n"
	writeTestFile(t, filepath.Join(docsDir, "design.md"), plainBody)

	if _, err := InitDocuments(projectRoot, packagePath); err != nil {
		t.Fatalf("InitDocuments() error = %v", err)
	}
	_, body, err := ParseIDDDocument([]byte(readTestFile(t, filepath.Join(docsDir, "design.md"))))
	if err != nil {
		t.Fatalf("ParseIDDDocument(design.md) error = %v", err)
	}
	if string(body) != "\n"+plainBody {
		t.Errorf("design body changed:\nwant %q\ngot  %q", "\n"+plainBody, body)
	}
}

// @test TEST-INTERNAL_COLLECTOR-021
func TestInitDocuments_MergesGenericFrontmatter(t *testing.T) {
	projectRoot := t.TempDir()
	packagePath := "internal/auth"
	if err := os.MkdirAll(filepath.Join(projectRoot, packagePath), 0o755); err != nil {
		t.Fatalf("MkdirAll(package) error = %v", err)
	}
	docsDir := filepath.Join(projectRoot, "docs", packagePath)
	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(docs) error = %v", err)
	}
	existing := "---\ntitle: Custom design\n---\n\n# Design\n\nKeep this body.\n"
	writeTestFile(t, filepath.Join(docsDir, "design.md"), existing)

	if _, err := InitDocuments(projectRoot, packagePath); err != nil {
		t.Fatalf("InitDocuments() error = %v", err)
	}
	content := readTestFile(t, filepath.Join(docsDir, "design.md"))
	if strings.Count(content, "---") != 2 {
		t.Errorf("design.md contains multiple frontmatter blocks:\n%s", content)
	}
	if !strings.Contains(content, "title: Custom design") {
		t.Errorf("design.md lost generic frontmatter:\n%s", content)
	}
	_, body, err := ParseIDDDocument([]byte(content))
	if err != nil {
		t.Fatalf("ParseIDDDocument(design.md) error = %v", err)
	}
	if string(body) != "\n# Design\n\nKeep this body.\n" {
		t.Errorf("design body changed: %q", body)
	}
}

// @test TEST-INTERNAL_COLLECTOR-021
func TestInitDocuments_RefusesLegacyMetadataBeforeWriting(t *testing.T) {
	projectRoot := t.TempDir()
	packagePath := "internal/auth"
	if err := os.MkdirAll(filepath.Join(projectRoot, packagePath), 0o755); err != nil {
		t.Fatalf("MkdirAll(package) error = %v", err)
	}
	docsDir := filepath.Join(projectRoot, "docs", packagePath)
	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(docs) error = %v", err)
	}
	legacySpec := `---
markers:
  - id: SPEC-INTERNAL_AUTH-001
    name: Authentication
---

# Specifications
`
	writeTestFile(t, filepath.Join(docsDir, "spec.md"), legacySpec)

	changed, err := InitDocuments(projectRoot, packagePath)
	if err == nil || !strings.Contains(err.Error(), "legacy metadata") {
		t.Fatalf("InitDocuments() error = %v, want legacy metadata refusal", err)
	}
	if len(changed) != 0 {
		t.Errorf("InitDocuments() changed = %v, want no writes", changed)
	}
	if got := readTestFile(t, filepath.Join(docsDir, "spec.md")); got != legacySpec {
		t.Error("legacy spec changed despite failed initialization")
	}
	if _, statErr := os.Stat(filepath.Join(docsDir, "design.md")); !os.IsNotExist(statErr) {
		t.Errorf("design.md stat error = %v, want no skeleton writes", statErr)
	}
}

// @test TEST-INTERNAL_COLLECTOR-021
func TestHasLegacyDocumentMetadata(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{name: "legacy marker frontmatter", content: "---\nmarkers:\n  - id: SPEC-AUTH-001\n---\n", want: true},
		{name: "legacy relationship field", content: "# Spec\n\n**Tests:** `TEST-AUTH-001`\n", want: true},
		{name: "legacy identifier heading", content: "## SPEC-AUTH-001: Authentication\n", want: true},
		{name: "plain narrative", content: "# Design\n\nArchitecture prose.\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasLegacyDocumentMetadata([]byte(tt.content)); got != tt.want {
				t.Errorf("hasLegacyDocumentMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

// @test TEST-INTERNAL_COLLECTOR-021
func TestInitDocuments_RefusesCentralCatalogBeforeWriting(t *testing.T) {
	projectRoot := t.TempDir()
	packagePath := "internal/auth"
	if err := os.MkdirAll(filepath.Join(projectRoot, packagePath), 0o755); err != nil {
		t.Fatalf("MkdirAll(package) error = %v", err)
	}
	docsDir := filepath.Join(projectRoot, "docs", packagePath)
	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(docs) error = %v", err)
	}
	writeTestFile(t, filepath.Join(docsDir, "idd.yaml"), "version: \"1.0\"\n")

	changed, err := InitDocuments(projectRoot, packagePath)
	if err == nil || !strings.Contains(err.Error(), "central catalog") {
		t.Fatalf("InitDocuments() error = %v, want central catalog refusal", err)
	}
	if len(changed) != 0 {
		t.Errorf("InitDocuments() changed = %v, want no writes", changed)
	}
	if _, statErr := os.Stat(filepath.Join(docsDir, "design.md")); !os.IsNotExist(statErr) {
		t.Errorf("design.md stat error = %v, want no skeleton writes", statErr)
	}
}

// @test TEST-INTERNAL_COLLECTOR-021
func TestRepairDocuments(t *testing.T) {
	docsDir := createIDDDocumentFixture(t)
	replaceTestFile(t, filepath.Join(docsDir, "design.md"), `version: "1.0"`, `version: ""`)
	replaceTestFile(t, filepath.Join(docsDir, "design.md"), "package: internal/auth", "package: wrong/package")

	bodies := make(map[string]string)
	for filename := range iddDocumentRoles {
		_, body, err := ParseIDDDocument([]byte(readTestFile(t, filepath.Join(docsDir, filename))))
		if err != nil {
			t.Fatalf("ParseIDDDocument(%s) error = %v", filename, err)
		}
		bodies[filename] = string(body)
	}

	changed, err := RepairDocuments(docsDir)
	if err != nil {
		t.Fatalf("RepairDocuments() error = %v", err)
	}
	if len(changed) != 1 {
		t.Fatalf("RepairDocuments() changed = %v, want design.md only", changed)
	}

	design, _, err := ParseIDDDocument([]byte(readTestFile(t, filepath.Join(docsDir, "design.md"))))
	if err != nil {
		t.Fatalf("ParseIDDDocument(design.md) error = %v", err)
	}
	if design.Version != "1.0" || design.Package != "internal/auth" || design.Document != "design" {
		t.Errorf("repaired design identity = %#v", design)
	}
	if got := strings.Join(design.Components, ","); got != "AuthModule" {
		t.Errorf("components = %q, want Markdown component declaration", got)
	}
	for filename, wantBody := range bodies {
		_, gotBody, parseErr := ParseIDDDocument([]byte(readTestFile(t, filepath.Join(docsDir, filename))))
		if parseErr != nil {
			t.Fatalf("ParseIDDDocument(%s) after repair error = %v", filename, parseErr)
		}
		if string(gotBody) != wantBody {
			t.Errorf("%s body changed during repair", filename)
		}
	}

	secondChanged, err := RepairDocuments(docsDir)
	if err != nil {
		t.Fatalf("second RepairDocuments() error = %v", err)
	}
	if len(secondChanged) != 0 {
		t.Errorf("second RepairDocuments() changed = %v, want idempotent", secondChanged)
	}
}

// @test TEST-INTERNAL_COLLECTOR-021
func TestRepairDocuments_SingleFileHasSingleWriteTarget(t *testing.T) {
	docsDir := createIDDDocumentFixture(t)
	specPath := filepath.Join(docsDir, "spec.md")
	testingPath := filepath.Join(docsDir, "testing.md")
	replaceTestFile(t, specPath, `version: "1.0"`, `version: ""`)
	testingBefore := readTestFile(t, testingPath)

	changed, err := RepairDocuments(specPath)
	if err != nil {
		t.Fatalf("RepairDocuments(spec.md) error = %v", err)
	}
	if len(changed) != 1 || filepath.Clean(changed[0]) != filepath.Clean(specPath) {
		t.Fatalf("RepairDocuments(spec.md) changed = %v, want only spec.md", changed)
	}
	if got := readTestFile(t, testingPath); got != testingBefore {
		t.Error("single-file repair changed testing.md")
	}
}

// @test TEST-INTERNAL_COLLECTOR-021
func TestRepairDocuments_RefusesCentralCatalog(t *testing.T) {
	docsDir := createIDDDocumentFixture(t)
	specPath := filepath.Join(docsDir, "spec.md")
	replaceTestFile(t, specPath, `version: "1.0"`, `version: ""`)
	specBefore := readTestFile(t, specPath)
	writeTestFile(t, filepath.Join(docsDir, "idd.yaml"), "version: \"1.0\"\n")

	changed, err := RepairDocuments(specPath)
	if err == nil || !strings.Contains(err.Error(), "central catalog") {
		t.Fatalf("RepairDocuments(spec.md) error = %v, want central catalog refusal", err)
	}
	if len(changed) != 0 {
		t.Errorf("RepairDocuments(spec.md) changed = %v, want no writes", changed)
	}
	if got := readTestFile(t, specPath); got != specBefore {
		t.Error("RepairDocuments changed spec.md despite central catalog refusal")
	}
}

// @test TEST-INTERNAL_COLLECTOR-021
func TestRepairDocuments_RefusesYAMLSemanticRecords(t *testing.T) {
	docsDir := createIDDDocumentFixture(t)
	specPath := filepath.Join(docsDir, "spec.md")
	replaceTestFile(t, specPath, "  document: spec", "  document: spec\n  specs: []")
	before := readTestFile(t, specPath)

	changed, err := RepairDocuments(specPath)
	if err == nil || !strings.Contains(err.Error(), "Markdown record") {
		t.Fatalf("RepairDocuments(spec.md) error = %v, want Markdown record migration refusal", err)
	}
	if len(changed) != 0 {
		t.Errorf("RepairDocuments(spec.md) changed = %v, want no writes", changed)
	}
	if got := readTestFile(t, specPath); got != before {
		t.Error("RepairDocuments changed spec.md despite YAML semantic records")
	}
}

func createIDDDocumentFixture(t *testing.T) string {
	t.Helper()
	docsDir := filepath.Join(t.TempDir(), "docs", "internal", "auth")
	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	for filename, content := range map[string]string{
		"design.md":   validDesignDocument,
		"contract.md": validContractDocument,
		"spec.md":     validSpecDocument,
		"testing.md":  validTestingDocument,
	} {
		writeTestFile(t, filepath.Join(docsDir, filename), content)
	}
	return docsDir
}

func replaceTestFile(t *testing.T, path, old, replacement string) {
	t.Helper()
	content := readTestFile(t, path)
	if !strings.Contains(content, old) {
		t.Fatalf("%s does not contain %q", path, old)
	}
	writeTestFile(t, path, strings.Replace(content, old, replacement, 1))
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", path, err)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%s) error = %v", path, err)
	}
	return string(data)
}
