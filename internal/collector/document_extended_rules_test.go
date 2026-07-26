// Package collector tests deterministic fields added to self-describing IDD records.

// Spec: docs/internal/collector/spec.md
// Test: docs/internal/collector/testing.md
package collector

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jingxu9x/idd-cli/internal/config"
)

// @test-contract TEST-INTERNAL_COLLECTOR-021
func TestParseIDDDocument_ExtendedRecordFields(t *testing.T) {
	design, _, err := ParseIDDDocument([]byte(strings.Replace(
		validDesignDocument,
		"- **Concerns:** `security`",
		"- **Depends on:** `internal/store#CredentialStore`, `AuditTrail`\n- **Concerns:** `security`",
		1,
	)))
	if err != nil {
		t.Fatalf("ParseIDDDocument(design) error = %v", err)
	}
	if len(design.ComponentRecords) != 1 {
		t.Fatalf("ComponentRecords = %#v, want one", design.ComponentRecords)
	}
	component := design.ComponentRecords[0]
	if component.Name != "AuthModule" ||
		component.Purpose == "" ||
		component.Status != "active" ||
		strings.Join(component.DependsOn, ",") != "internal/store#CredentialStore,AuditTrail" ||
		strings.Join(component.Concerns, ",") != "security" {
		t.Errorf("Component record = %#v", component)
	}

	contract, _, err := ParseIDDDocument([]byte(validContractDocument))
	if err != nil {
		t.Fatalf("ParseIDDDocument(contract) error = %v", err)
	}
	if len(contract.ContractRecords) != 1 {
		t.Fatalf("ContractRecords = %#v, want one", contract.ContractRecords)
	}
	if got := contract.ContractRecords[0]; got.Name != "Authenticator" ||
		got.Guarantees == "" ||
		got.Status != "active" ||
		strings.Join(got.Concerns, ",") != "compatibility" {
		t.Errorf("Contract record = %#v", got)
	}

	spec, _, err := ParseIDDDocument([]byte(strings.Replace(
		validSpecDocument,
		"- **Status:** `active`",
		"- **Status:** `superseded`\n- **Deprecated by:** `SPEC-INTERNAL_AUTH-002`",
		1,
	)))
	if err != nil {
		t.Fatalf("ParseIDDDocument(spec) error = %v", err)
	}
	if len(spec.Specs) != 1 {
		t.Fatalf("Specs = %#v, want one", spec.Specs)
	}
	if got := spec.Specs[0]; got.Acceptance == "" ||
		got.Status != "superseded" ||
		got.DeprecatedBy != "SPEC-INTERNAL_AUTH-002" ||
		strings.Join(got.Concerns, ",") != "security" {
		t.Errorf("SPEC record = %#v", got)
	}

	contractTest := strings.Replace(validTestingDocument, "`test`", "`contract`", 1)
	contractTest = strings.Replace(
		contractTest,
		"- **Status:** `active`",
		"- **Contracts:** `Authenticator`\n- **Status:** `active`\n- **Supersedes:** `TEST-INTERNAL_AUTH-000`\n- **Concerns:** `compatibility`",
		1,
	)
	testingDocument, _, err := ParseIDDDocument([]byte(contractTest))
	if err != nil {
		t.Fatalf("ParseIDDDocument(testing) error = %v", err)
	}
	if len(testingDocument.Tests) != 1 {
		t.Fatalf("Tests = %#v, want one", testingDocument.Tests)
	}
	if got := testingDocument.Tests[0]; got.Oracle == "" ||
		strings.Join(got.Contracts, ",") != "Authenticator" ||
		strings.Join(got.Supersedes, ",") != "TEST-INTERNAL_AUTH-000" ||
		strings.Join(got.Concerns, ",") != "compatibility" {
		t.Errorf("TEST record = %#v", got)
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-021
func TestDocCollector_IDDExtendedRuleValidation(t *testing.T) {
	tests := []struct {
		name     string
		mutate   func(*testing.T, string)
		wantCode string
	}{
		{name: "valid extended records"},
		{
			name: "component purpose is required",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(t, filepath.Join(docsDir, "design.md"), "**Purpose:**", "**Missing purpose:**")
			},
			wantCode: "purpose",
		},
		{
			name: "contract guarantees are required",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(t, filepath.Join(docsDir, "contract.md"), "**Guarantees:**", "**Missing guarantees:**")
			},
			wantCode: "guarantees",
		},
		{
			name: "record names reject reference separators",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(
					t,
					filepath.Join(docsDir, "contract.md"),
					"Contract: Authenticator",
					"Contract: Auth#enticator",
				)
			},
			wantCode: "contracts",
		},
		{
			name: "spec acceptance is required",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(t, filepath.Join(docsDir, "spec.md"), "**Acceptance:**", "**Missing acceptance:**")
			},
			wantCode: "acceptance",
		},
		{
			name: "test oracle is required",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(t, filepath.Join(docsDir, "testing.md"), "**Oracle:**", "**Missing oracle:**")
			},
			wantCode: "oracle",
		},
		{
			name: "test covers rejects unparsed list values",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(
					t,
					filepath.Join(docsDir, "testing.md"),
					"`SPEC-INTERNAL_AUTH-001`",
					"`SPEC-INTERNAL_AUTH-001`, `not-a-spec`",
				)
			},
			wantCode: "covers",
		},
		{
			name: "contract test names contracts",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(t, filepath.Join(docsDir, "testing.md"), "`test`", "`contract`")
			},
			wantCode: "contracts",
		},
		{
			name: "status uses supported lifecycle",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(t, filepath.Join(docsDir, "spec.md"), "`active`", "`draft`")
			},
			wantCode: "status",
		},
		{
			name: "concern uses supported value",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(t, filepath.Join(docsDir, "spec.md"), "`security`", "`availability`")
			},
			wantCode: "concerns",
		},
		{
			name: "selected concern requires narrative section",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(
					t,
					filepath.Join(docsDir, "spec.md"),
					"### Security\n\nFailures must not reveal which credential field was incorrect.\n",
					"",
				)
			},
			wantCode: "concerns",
		},
		{
			name: "test concern is accepted with narrative section",
			mutate: func(t *testing.T, docsDir string) {
				path := filepath.Join(docsDir, "testing.md")
				replaceTestFile(
					t,
					path,
					"- **Status:** `active`",
					"- **Status:** `active`\n- **Concerns:** `security`",
				)
				replaceTestFile(
					t,
					path,
					"### Scenarios",
					"### Security\n\nTest diagnostics never include submitted credentials.\n\n### Scenarios",
				)
			},
		},
		{
			name: "component dependency resolves",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(
					t,
					filepath.Join(docsDir, "design.md"),
					"- **Status:** `active`",
					"- **Status:** `active`\n- **Depends on:** `MissingComponent`",
				)
			},
			wantCode: "depends-on",
		},
		{
			name: "component dependency uses complete scoped name",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(
					t,
					filepath.Join(docsDir, "design.md"),
					"- **Status:** `active`",
					"- **Status:** `active`\n- **Depends on:** `internal/shared#`",
				)
			},
			wantCode: "depends-on",
		},
		{
			name: "spec contract uses complete scoped name",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(
					t,
					filepath.Join(docsDir, "spec.md"),
					"- **Contract:** `Authenticator`",
					"- **Contract:** `internal/shared#`",
				)
			},
			wantCode: "contract",
		},
		{
			name: "lifecycle does not reference itself",
			mutate: func(t *testing.T, docsDir string) {
				replaceTestFile(
					t,
					filepath.Join(docsDir, "spec.md"),
					"- **Status:** `active`",
					"- **Status:** `active`\n- **Supersedes:** `SPEC-INTERNAL_AUTH-001`",
				)
			},
			wantCode: "supersedes",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			docsDir := createIDDDocumentFixture(t)
			if test.mutate != nil {
				test.mutate(t, docsDir)
			}
			_, errs, err := NewDocCollector(config.Default()).Collect(
				context.Background(),
				filepath.Dir(filepath.Dir(docsDir)),
			)
			if err != nil {
				t.Fatalf("Collect() error = %v", err)
			}
			if test.wantCode == "" {
				if len(errs) != 0 {
					t.Fatalf("Collect() errors = %v, want none", errs)
				}
				return
			}
			for _, validationErr := range errs {
				if validationErr.Rule == "idd-document-schema" && validationErr.Code == test.wantCode {
					return
				}
				if validationErr.Rule == "idd-document-reference" && validationErr.Code == test.wantCode {
					return
				}
			}
			t.Fatalf("Collect() errors = %v, want code %q", errs, test.wantCode)
		})
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-021
func TestValidateIDDLifecycleRelationships(t *testing.T) {
	tests := []struct {
		name     string
		records  []iddLifecycleRecord
		wantCode string
	}{
		{
			name: "replacement can be owned by successor only",
			records: []iddLifecycleRecord{
				{
					ID:        "SPEC-DEMO-001",
					Index:     0,
					Lifecycle: IDDRecordLifecycle{Status: "superseded"},
				},
				{
					ID:        "SPEC-DEMO-002",
					Index:     1,
					Lifecycle: IDDRecordLifecycle{Supersedes: []string{"SPEC-DEMO-001"}},
				},
			},
		},
		{
			name: "matching fields agree",
			records: []iddLifecycleRecord{
				{
					ID: "SPEC-DEMO-001",
					Lifecycle: IDDRecordLifecycle{
						Status:       "superseded",
						DeprecatedBy: "SPEC-DEMO-002",
					},
				},
				{
					ID:        "SPEC-DEMO-002",
					Index:     1,
					Lifecycle: IDDRecordLifecycle{Supersedes: []string{"SPEC-DEMO-001"}},
				},
			},
		},
		{
			name: "superseded requires a replacement",
			records: []iddLifecycleRecord{
				{
					ID:        "SPEC-DEMO-001",
					Lifecycle: IDDRecordLifecycle{Status: "superseded"},
				},
			},
			wantCode: "replacement",
		},
		{
			name: "two successors conflict",
			records: []iddLifecycleRecord{
				{
					ID: "SPEC-DEMO-001",
					Lifecycle: IDDRecordLifecycle{
						DeprecatedBy: "SPEC-DEMO-003",
					},
				},
				{
					ID:        "SPEC-DEMO-002",
					Index:     1,
					Lifecycle: IDDRecordLifecycle{Supersedes: []string{"SPEC-DEMO-001"}},
				},
				{ID: "SPEC-DEMO-003", Index: 2},
			},
			wantCode: "replacement",
		},
		{
			name: "replacement cycle",
			records: []iddLifecycleRecord{
				{
					ID:        "SPEC-DEMO-001",
					Lifecycle: IDDRecordLifecycle{Supersedes: []string{"SPEC-DEMO-002"}},
				},
				{
					ID:        "SPEC-DEMO-002",
					Index:     1,
					Lifecycle: IDDRecordLifecycle{Supersedes: []string{"SPEC-DEMO-001"}},
				},
			},
			wantCode: "replacement-cycle",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			document := &parsedIDDDocument{
				Path: "docs/demo/spec.md",
				Index: iddDocumentNodeIndex{
					recordLines:      map[string][]int{"specs": {10, 20, 30}},
					recordFieldLines: make(map[string][]map[string]int),
				},
			}
			errs := validateIDDLifecycleRelationships(
				document,
				"specs",
				"SPEC",
				test.records,
			)
			if test.wantCode == "" {
				if len(errs) != 0 {
					t.Fatalf("errors = %#v, want none", errs)
				}
				return
			}
			for _, validationErr := range errs {
				if validationErr.Code == test.wantCode {
					return
				}
			}
			t.Fatalf("errors = %#v, want code %q", errs, test.wantCode)
		})
	}
}
