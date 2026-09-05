package collector

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/jingxu9x/idd-cli/internal/config"
)

// @test-contract TEST-INTERNAL_COLLECTOR-031
// @test TEST-INTERNAL_COLLECTOR-033
func TestBuildSpecReviewContext(t *testing.T) {
	tests := []struct {
		name      string
		specID    string
		setup     func(t *testing.T, root string)
		wantError string
		check     func(t *testing.T, context *SpecReviewContext)
	}{
		{
			name:   "collects authored and declaration evidence",
			specID: "SPEC-SAMPLE-001",
			setup: func(t *testing.T, root string) {
				writeReviewFixture(t, root, "sample", "SPEC-SAMPLE-001")
				writeReviewFile(t, filepath.Join(root, "source", "sample.go"), `package sample

// @implement SPEC-SAMPLE-001
func execute(input string) string {
	return input
}

// @test TEST-SAMPLE-001
func TestExecute(t *testing.T) {
	if execute("value") != "value" {
		t.Fatal("value changed")
	}
}
`)
			},
			check: func(t *testing.T, context *SpecReviewContext) {
				if context.Schema != reviewContextSchema {
					t.Fatalf("Schema = %q, want %q", context.Schema, reviewContextSchema)
				}
				if context.Spec.Requirement != "Return the supplied value without mutation." {
					t.Errorf("Requirement = %q", context.Spec.Requirement)
				}
				if !strings.Contains(context.Spec.Markdown, "### Details") ||
					!strings.Contains(context.Spec.Markdown, "Do not allocate") {
					t.Errorf("SPEC Markdown = %q", context.Spec.Markdown)
				}
				if strings.Contains(context.Spec.Markdown, "Package guidance") {
					t.Errorf("SPEC Markdown crossed into the next H2 section: %q", context.Spec.Markdown)
				}
				if len(context.Contracts) != 1 || context.Contracts[0].Guarantees != "Execution preserves the supplied value." {
					t.Errorf("Contracts = %#v", context.Contracts)
				}
				if len(context.Components) == 0 || context.Components[0].Boundary != "Accept one value and return one value." {
					t.Errorf("Components = %#v", context.Components)
				}
				if len(context.Tests) != 1 || context.Tests[0].ID != "TEST-SAMPLE-001" {
					t.Errorf("Tests = %#v", context.Tests)
				}
				if strings.Contains(context.Tests[0].Markdown, "Testing guidance") {
					t.Errorf("TEST Markdown crossed into the next H2 section: %q", context.Tests[0].Markdown)
				}
				if len(context.Declarations) != 2 {
					t.Fatalf("Declarations = %#v", context.Declarations)
				}
				declaration := context.Declarations[0]
				if declaration.Name != "execute" ||
					declaration.Role != "implementation" ||
					!strings.Contains(declaration.Excerpt, "return input") {
					t.Errorf("Declaration = %#v", declaration)
				}
				testDeclaration := context.Declarations[1]
				if testDeclaration.Name != "TestExecute" ||
					testDeclaration.Role != "test" ||
					len(testDeclaration.References) != 1 ||
					testDeclaration.References[0] != "TEST-SAMPLE-001" {
					t.Errorf("Test declaration = %#v", testDeclaration)
				}
			},
		},
		{
			name:   "collects cross-package contract and test evidence",
			specID: "SPEC-SAMPLE-001",
			setup: func(t *testing.T, root string) {
				writeReviewFixture(t, root, "sample", "SPEC-SAMPLE-001")
				writeReviewFixture(t, root, "shared", "SPEC-SHARED-001")

				sampleSpec := filepath.Join(root, "docs", "sample", "spec.md")
				replaceTestFile(t, sampleSpec, "`Execution`", "`shared#SharedExecution`")

				sharedContract := filepath.Join(root, "docs", "shared", "contract.md")
				replaceTestFile(
					t,
					sharedContract,
					"Contract: Execution",
					"Contract: SharedExecution",
				)
				replaceTestFile(
					t,
					sharedContract,
					"Execution preserves the supplied value.",
					"Shared execution preserves the supplied value.",
				)

				sharedTesting := filepath.Join(root, "docs", "shared", "testing.md")
				content := readTestFile(t, sharedTesting)
				content += `

## TEST-SHARED-002: Cross-package evidence

- **Kind:** ` + "`test`" + `
- **Covers:** ` + "`SPEC-SAMPLE-001`" + `

**Purpose:** Exercise the sample behavior from the shared package.

**Oracle:** The shared result equals the supplied sample value.
`
				writeTestFile(t, sharedTesting, content)
			},
			check: func(t *testing.T, context *SpecReviewContext) {
				if len(context.Contracts) != 1 ||
					context.Contracts[0].Name != "SharedExecution" ||
					!strings.Contains(context.Contracts[0].File, filepath.Join("docs", "shared")) {
					t.Errorf("cross-package Contracts = %#v", context.Contracts)
				}
				gotTests := make([]string, 0, len(context.Tests))
				for _, test := range context.Tests {
					gotTests = append(gotTests, test.ID)
				}
				wantTests := []string{"TEST-SAMPLE-001", "TEST-SHARED-002"}
				if !reflect.DeepEqual(gotTests, wantTests) {
					t.Errorf("cross-package Tests = %#v, want %#v", gotTests, wantTests)
				}
			},
		},
		{
			name:      "rejects malformed identifier",
			specID:    "TEST-SAMPLE-001",
			setup:     func(t *testing.T, root string) {},
			wantError: "not a valid SPEC identifier",
		},
		{
			name:   "reports missing canonical owner",
			specID: "SPEC-SAMPLE-999",
			setup: func(t *testing.T, root string) {
				writeReviewFixture(t, root, "sample", "SPEC-SAMPLE-001")
			},
			wantError: "was not found",
		},
		{
			name:   "reports multiple canonical owners",
			specID: "SPEC-SAMPLE-001",
			setup: func(t *testing.T, root string) {
				writeReviewFixture(t, root, "sample", "SPEC-SAMPLE-001")
				writeReviewFixture(t, root, "nested/sample", "SPEC-SAMPLE-001")
			},
			wantError: "multiple document owners",
		},
		{
			name:   "reports duplicate records in one owner file",
			specID: "SPEC-SAMPLE-001",
			setup: func(t *testing.T, root string) {
				writeReviewFixture(t, root, "sample", "SPEC-SAMPLE-001")
				path := filepath.Join(root, "docs", "sample", "spec.md")
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("ReadFile(%s) error = %v", path, err)
				}
				content = append(content, []byte(
					"\n\n## SPEC-SAMPLE-001: Duplicate\n\n"+
						"- **Components:** `Runner`\n"+
						"- **Contracts:** `Execution`\n\n"+
						"**Requirement:** Duplicate the record.\n\n"+
						"**Acceptance:** Duplicate evidence exists.\n",
				)...)
				if err := os.WriteFile(path, content, 0644); err != nil {
					t.Fatalf("WriteFile(%s) error = %v", path, err)
				}
			},
			wantError: "multiple document owners",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			test.setup(t, root)
			cfg := config.Default()
			cfg.Docs.Patterns = []string{"**/*.md"}
			cfg.Docs.IgnorePaths = nil
			cfg.Code.Patterns = []string{"**/*.go"}
			cfg.Code.IgnorePaths = nil

			contextValue, err := BuildSpecReviewContext(
				context.Background(),
				cfg,
				filepath.Join(root, "docs"),
				filepath.Join(root, "source"),
				test.specID,
			)
			if test.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantError) {
					t.Fatalf("error = %v, want substring %q", err, test.wantError)
				}
				return
			}
			if err != nil {
				t.Fatalf("BuildSpecReviewContext() error = %v", err)
			}
			test.check(t, contextValue)
		})
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-031
// @test TEST-INTERNAL_COLLECTOR-033
func TestBuildSpecReviewContexts(t *testing.T) {
	root := t.TempDir()
	writeReviewFixture(t, root, "sample", "SPEC-SAMPLE-001")
	writeReviewFixture(t, root, "other", "SPEC-OTHER-001")
	cfg := config.Default()
	cfg.Docs.Patterns = []string{"**/*.md"}
	cfg.Docs.IgnorePaths = nil
	cfg.Code.Patterns = []string{"**/*.go"}
	cfg.Code.IgnorePaths = nil

	batch, err := BuildSpecReviewContexts(
		context.Background(),
		cfg,
		filepath.Join(root, "docs"),
		filepath.Join(root, "source"),
		[]string{"SPEC-OTHER-001", "SPEC-SAMPLE-001", "SPEC-OTHER-001"},
	)
	if err != nil {
		t.Fatalf("BuildSpecReviewContexts() error = %v", err)
	}
	if batch.Schema != reviewContextBatchSchema {
		t.Errorf("Schema = %q, want %q", batch.Schema, reviewContextBatchSchema)
	}
	wantIDs := []string{"SPEC-OTHER-001", "SPEC-SAMPLE-001"}
	if !reflect.DeepEqual(batch.RequestedSpecIDs, wantIDs) {
		t.Errorf("RequestedSpecIDs = %#v, want %#v", batch.RequestedSpecIDs, wantIDs)
	}
	if len(batch.Contexts) != 2 {
		t.Fatalf("Contexts = %#v, want two contexts", batch.Contexts)
	}
	for index, wantID := range wantIDs {
		if batch.Contexts[index].Spec.ID != wantID {
			t.Errorf("Contexts[%d].Spec.ID = %q, want %q", index, batch.Contexts[index].Spec.ID, wantID)
		}
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-031
func TestNormalizeReviewSpecIDs(t *testing.T) {
	tooMany := make([]string, 0, reviewContextMaxSpecs+1)
	for index := 1; index <= reviewContextMaxSpecs+1; index++ {
		tooMany = append(tooMany, fmt.Sprintf("SPEC-SAMPLE-%03d", index))
	}
	tests := []struct {
		name      string
		input     []string
		want      []string
		wantError string
	}{
		{name: "empty", wantError: "at least one"},
		{
			name:  "deduplicates in first-request order",
			input: []string{"SPEC-SAMPLE-002", "SPEC-SAMPLE-001", "SPEC-SAMPLE-002"},
			want:  []string{"SPEC-SAMPLE-002", "SPEC-SAMPLE-001"},
		},
		{name: "malformed", input: []string{"sample"}, wantError: "not a valid SPEC"},
		{name: "unique limit", input: tooMany, wantError: "at most 10"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := normalizeReviewSpecIDs(test.input)
			if test.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantError) {
					t.Fatalf("error = %v, want substring %q", err, test.wantError)
				}
				return
			}
			if err != nil {
				t.Fatalf("normalizeReviewSpecIDs() error = %v", err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("normalizeReviewSpecIDs() = %#v, want %#v", got, test.want)
			}
		})
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-031
func TestBoundedDeclarationExcerpt(t *testing.T) {
	tests := []struct {
		name          string
		lines         []string
		declaration   SourceDeclaration
		remaining     int
		want          string
		wantTruncated bool
	}{
		{
			name:        "complete declaration",
			lines:       []string{"func run() {", "\twork()", "}"},
			declaration: SourceDeclaration{Line: 1, EndLine: 3},
			remaining:   100,
			want:        "func run() {\n\twork()\n}",
		},
		{
			name:          "remaining budget truncates declaration",
			lines:         []string{"func run() {", "\twork()", "}"},
			declaration:   SourceDeclaration{Line: 1, EndLine: 3},
			remaining:     8,
			want:          "func run",
			wantTruncated: true,
		},
		{
			name:          "exhausted budget remains visible",
			lines:         []string{"func run() {}"},
			declaration:   SourceDeclaration{Line: 1, EndLine: 1},
			remaining:     0,
			want:          "",
			wantTruncated: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, truncated := boundedDeclarationExcerpt(
				test.lines,
				test.declaration,
				test.remaining,
			)
			if got != test.want || truncated != test.wantTruncated {
				t.Fatalf(
					"boundedDeclarationExcerpt() = (%q, %t), want (%q, %t)",
					got,
					truncated,
					test.want,
					test.wantTruncated,
				)
			}
		})
	}
}

func writeReviewFixture(t *testing.T, root string, packagePath string, specID string) {
	t.Helper()
	directory := filepath.Join(root, "docs", filepath.FromSlash(packagePath))
	module := strings.ToUpper(filepath.Base(filepath.FromSlash(packagePath)))
	packageName := filepath.ToSlash(packagePath)
	if err := os.MkdirAll(filepath.Join(root, "source"), 0755); err != nil {
		t.Fatalf("MkdirAll(source) error = %v", err)
	}
	writeReviewFile(t, filepath.Join(directory, "design.md"), `---
idd:
  version: "1.1"
  package: `+packageName+`
  namespace: `+module+`
---

# Design

## Component: Runner

**Purpose:** Execute one sample behavior.

**Ownership:** Own sample execution behavior.

**Boundary:** Accept one value and return one value.

**Decisions:** Keep preservation independent from transport concerns.
`)
	writeReviewFile(t, filepath.Join(directory, "contract.md"), `---
idd:
  version: "1.1"
  package: `+packageName+`
  namespace: `+module+`
---

# Contracts

## Contract: Execution

**Guarantees:** Execution preserves the supplied value.
`)
	writeReviewFile(t, filepath.Join(directory, "spec.md"), `---
idd:
  version: "1.1"
  package: `+packageName+`
  namespace: `+module+`
---

# Specifications

## `+specID+`: Preserve input

- **Components:** `+"`Runner`"+`
- **Contracts:** `+"`Execution`"+`

**Requirement:** Return the supplied value without mutation.

**Acceptance:** Calling execute returns the exact supplied value.

### Details

Do not allocate a replacement value.

## Package guidance

Keep package-wide guidance outside the preceding SPEC record.
`)
	writeReviewFile(t, filepath.Join(directory, "testing.md"), `---
idd:
  version: "1.1"
  package: `+packageName+`
  namespace: `+module+`
---

# Testing

## TEST-`+module+`-001: Preserve input evidence

- **Kind:** `+"`test`"+`
- **Covers:** `+"`"+specID+"`"+`

**Purpose:** Exercise the documented preservation behavior.

**Oracle:** The returned value equals the supplied value.

## Testing guidance

Keep package-wide guidance outside the preceding TEST record.
`)
}

func writeReviewFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("MkdirAll(%s) error = %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", path, err)
	}
}
