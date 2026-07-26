---
idd:
  version: "1.0"
  package: internal/engine
---

# Testing: internal/engine

## TEST-INTERNAL_ENGINE-001: Headerless source and engine state

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-004`

**Purpose:**

Prove that production, behavior-test, and contract-test declarations join their
matching document records through syntax-tree-bound identifiers without
repeated IDD document-path headers, and that the engine returns structured
state rather than an operational error.

**Oracle:** Table-driven source collection returns the expected SPEC or TEST
evidence for `@implement`, `@test`, and `@test-contract`; adding the matching
document observation makes every engine result valid without `Spec`,
`Contract`, or `Test` header paths. The smaller production-source fixture also
passes with only its useful language package description.

## TEST-INTERNAL_ENGINE-002: Link inference and path-header independence

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-004`

**Purpose:**

Prove the source/target identifier-type mapping used during graph construction
and prove that a source file without a package comment or IDD path header does
not acquire a second code-to-document association mechanism.

**Oracle:** The test passes only when its assertions confirm the
source/target identifier-type mapping and confirm the headerless fixture
remains valid because correspondence is joined by identifier evidence.

## TEST-INTERNAL_ENGINE-003: Relationship, correspondence, design, and annotation-kind rules

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-004`

**Purpose:**

Exercise several focused rule tests sharing one historical TEST identifier:
legacy relationship consistency, identifier-derived doc/code correspondence,
ordinary leading comments without IDD semantics, non-empty self-describing
design sections, and agreement between documented TEST kind and source
annotation kind. Each function isolates one fixture and asserts its exact rule
rather than treating the combined identifier as one scenario.

**Oracle:** The test passes only when its assertions exercise several
focused rule tests sharing one historical TEST identifier: legacy relationship
consistency, paired and unpaired identifier evidence, ordinary leading
comments, non-empty self-describing design sections, and documented/source
TEST-kind agreement. Each fixture asserts its exact rule independently.

## TEST-INTERNAL_ENGINE-004: Structural errors, source locations, and language comments

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-004`

**Purpose:**

Prove collector structural findings survive injection, graph completeness uses
document evidence for missing fields, and an ordinary language package comment
remains outside IDD association policy.

**Oracle:** The test passes only when its assertions confirm collector
structural findings survive injection, graph completeness uses document
evidence for missing fields, and the package-comment fixture passes without
repeated document paths.

## TEST-INTERNAL_ENGINE-005: Report construction and removed path-header policy

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-004`

**Purpose:**

Prove the engine report envelope contains accumulated state and that source
files without legacy path headers remain valid without a package-comment policy
switch.

**Oracle:** The test passes only when its assertions confirm the engine
report envelope contains accumulated state and the headerless fixture produces
no obsolete path-header finding.

## TEST-INTERNAL_ENGINE-006: Headerless test source

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-004`

**Purpose:**

Prove that test files do not repeat SPEC and TEST document paths and rely on
their declaration-level TEST annotations for association.

**Oracle:** The test passes only when its assertions confirm the headerless test
source remains valid while test annotations continue to carry the non-derivable
TEST relationship.

## TEST-INTERNAL_ENGINE-007: Consistency warnings and missing TEST document links

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-004`
- **Contracts:** `SemanticConsistencyWarning`

**Purpose:**

Prove description-similarity warning behavior across low/high similarity,
disabled checks, missing descriptions, function-locator descriptions,
code-only observations, and threshold changes. A separate contract-test fixture
using this identifier proves that contract tests need no repeated SPEC, TEST,
or Contract header paths.

**Oracle:** The test passes only when its assertions prove
description-similarity warning behavior across low/high similarity, disabled checks,
missing descriptions, function-locator descriptions, code-only observations,
and threshold changes. The separate contract-test fixture remains valid
without file-level document paths.

## TEST-INTERNAL_ENGINE-013: Engine storage contract

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_ENGINE-002`
- **Contracts:** `EngineLifecycle`

**Purpose:**

Verify that construction returns a non-nil engine with non-nil retained config,
graph, and result fields. Same-package inspection is the oracle.

**Oracle:** The test passes only when its assertions confirm
construction returns a non-nil engine with non-nil retained config, graph, and result
fields. Same-package inspection is the oracle.

## TEST-INTERNAL_ENGINE-014: Initial engine result contract

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_ENGINE-002`
- **Contracts:** `EngineLifecycle`

**Purpose:**

Verify the constructor's single-run starting state: invalid result and no
errors before validation executes.

**Oracle:** The test passes only when its assertions confirm the
constructor's single-run starting state: invalid result and no errors before validation
executes.

## TEST-INTERNAL_ENGINE-015: Run builds graph statistics

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-004`
- **Contracts:** `EngineLifecycle`, `RepositoryValidation`

**Purpose:**

Verify that a reciprocal SPEC/TEST identifier set produces a non-nil result and
the expected node and link statistics through `Run`.

**Oracle:** The test passes only when its assertions confirm a
reciprocal SPEC/TEST identifier set produces a non-nil result and the expected node and
link statistics through `Run`.

## TEST-INTERNAL_ENGINE-016: Unidirectional links remain processable

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-004`
- **Contracts:** `EngineLifecycle`, `RepositoryValidation`

**Purpose:**

Verify that one forward link is retained and no obsolete
`missing-backlink` warning is introduced merely because the reverse edge is
absent.

**Oracle:** The test passes only when its assertions confirm one
forward link is retained and no obsolete `missing-backlink` warning is introduced merely
because the reverse edge is absent.

## TEST-INTERNAL_ENGINE-017: Orphan validation participates in run

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-004`
- **Contracts:** `EngineLifecycle`, `RepositoryValidation`

**Purpose:**

Verify that an unlinked SPEC produces the named orphan error when orphan
validation is enabled.

**Oracle:** The test passes only when its assertions confirm an
unlinked SPEC produces the named orphan error when orphan validation is enabled.

## TEST-INTERNAL_ENGINE-020: Duplicate SPEC heading rejection

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove two level-two declarations with the same identifier in one legacy
document produce a duplicate-heading finding at a useful source location.

**Oracle:** The test passes only when its assertions confirm two
level-two declarations with the same identifier in one legacy document produce a
duplicate-heading finding at a useful source location.

## TEST-INTERNAL_ENGINE-021: Unique headings remain valid

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove distinct level-two identifier headings do not trigger the duplicate rule.

**Oracle:** The test passes only when its assertions confirm distinct
level-two identifier headings do not trigger the duplicate rule.

## TEST-INTERNAL_ENGINE-022: Legacy contract markers are accepted

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove a legacy contract document with frontmatter markers satisfies the
contract/design marker compatibility check.

**Oracle:** The test passes only when its assertions confirm a legacy
contract document with frontmatter markers satisfies the contract/design marker
compatibility check.

## TEST-INTERNAL_ENGINE-023: Missing legacy contract markers are reported

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove a legacy contract document lacking marker declarations produces the
expected marker finding while self-describing packages use collector-owned
records instead.

**Oracle:** The test passes only when its assertions confirm a legacy
contract document lacking marker declarations produces the expected marker finding while
self-describing packages use collector-owned records instead.

## TEST-INTERNAL_ENGINE-024: Design files are excluded from contract-marker requirement

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove the legacy marker rule does not incorrectly require contract-style
markers in a design document.

**Oracle:** The test passes only when its assertions confirm the legacy
marker rule does not incorrectly require contract-style markers in a design document.

## TEST-INTERNAL_ENGINE-025: Documented package path existence

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove package-path checks distinguish existing implementation directories from
missing ones in temporary repository fixtures.

**Oracle:** The test passes only when its assertions confirm package-path
checks distinguish existing implementation directories from missing ones in temporary
repository fixtures.

## TEST-INTERNAL_ENGINE-026: Ignored documentation paths

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove configured document ignore patterns suppress implementation-path checks
for matching documentation.

**Oracle:** The test passes only when its assertions confirm configured
document ignore patterns suppress implementation-path checks for matching documentation.

## TEST-INTERNAL_ENGINE-027: Existing documented paths produce no warning

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove a valid documented implementation path does not produce a false
path-existence finding.

**Oracle:** The test passes only when its assertions confirm a valid
documented implementation path does not produce a false path-existence finding.

## TEST-INTERNAL_ENGINE-028: Missing legacy related files

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove a nested legacy document without `related_files` frontmatter is rejected
by the compatibility rule.

**Oracle:** The test passes only when its assertions confirm a nested
legacy document without `related_files` frontmatter is rejected by the compatibility
rule.

## TEST-INTERNAL_ENGINE-029: Related files and IDD metadata detection

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove valid legacy `related_files` passes and leading self-describing `idd`
metadata is detected so legacy-only rules can skip the package. Table-driven
metadata cases cover valid placement and misleading body text.

**Oracle:** The test passes only when its assertions confirm valid legacy
`related_files` passes and leading self-describing `idd` metadata is detected so
legacy-only rules can skip the package. Table-driven metadata cases cover valid
placement and misleading body text.

## TEST-INTERNAL_ENGINE-030: Main package path-header independence

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove command `main` source also has no IDD document-path header requirement.

**Oracle:** The test passes only when its assertions confirm the headerless
`main` fixture returns a valid result without special-case path-comment
handling.

## TEST-INTERNAL_ENGINE-031: Incomplete package document set

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove both a partially documented package and a nested source sub-package with
no documentation directory produce findings for their missing canonical role
files.

**Oracle:** The test passes only when its assertions confirm the partial
package reports only its absent files, while the nested source package reports
all four paths under its exact nested `docs/<package>/` directory. The complete
parent package must produce no finding and must not satisfy the child by prefix
matching.

## TEST-INTERNAL_ENGINE-032: Complete package document set

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove a directory containing design, contract, spec, and testing documents
passes the package-set rule.

**Oracle:** The test passes only when its assertions confirm a directory
containing design, contract, spec, and testing documents passes the package-set rule.

## TEST-INTERNAL_ENGINE-033: Docs-root files do not form package sets

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove narrative files directly under the configured docs root are not mistaken
for package directories requiring four role files.

**Oracle:** The test passes only when its assertions confirm narrative
files directly under the configured docs root are not mistaken for package directories
requiring four role files.

## TEST-INTERNAL_ENGINE-034: Duplicate documentation IDs across packages

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-004`

**Purpose:**

Prove same-origin documentation observations with one ID in different package
directories produce duplicate-ID findings for each conflicting location.

**Oracle:** The test passes only when its assertions confirm same-origin
documentation observations with one ID in different package directories produce
duplicate-ID findings for each conflicting location.

## TEST-INTERNAL_ENGINE-035: Same-package ID observations are not package conflicts

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-004`

**Purpose:**

Prove multiple observations of one ID inside the same directory remain valid
doc/code evidence rather than cross-package duplicate errors.

**Oracle:** The test passes only when its assertions confirm multiple
observations of one ID inside the same directory remain valid doc/code evidence rather
than cross-package duplicate errors.

## TEST-INTERNAL_ENGINE-036: Duplicate code IDs across packages

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-004`

**Purpose:**

Prove source annotations with one ID in different code package directories
produce code-side duplicate findings.

**Oracle:** The test passes only when its assertions confirm source
annotations with one ID in different code package directories produce code-side
duplicate findings.

## TEST-INTERNAL_ENGINE-037: Package-derived duplicate rename guidance

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove duplicate diagnostics suggest a replacement module segment derived from
the conflicting source path when that suggestion differs from the current ID.

**Oracle:** The test passes only when its assertions confirm duplicate
diagnostics suggest a replacement module segment derived from the conflicting source
path when that suggestion differs from the current ID.

## TEST-INTERNAL_ENGINE-038: Annotated private functions are allowed

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove a private function may carry a valid implementation annotation without
triggering the public-declaration rule.

**Oracle:** The test passes only when its assertions confirm a private
function may carry a valid implementation annotation without triggering the
public-declaration rule.

## TEST-INTERNAL_ENGINE-039: Annotated private methods are allowed

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove the same optional annotation policy applies to unexported receiver
methods.

**Oracle:** The test passes only when its assertions confirm the same
optional annotation policy applies to unexported receiver methods.

## TEST-INTERNAL_ENGINE-040: Annotated private types are allowed

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove an unexported type may be traced to a SPEC without being treated as an
invalid placement.

**Oracle:** The test passes only when its assertions confirm an
unexported type may be traced to a SPEC without being treated as an invalid placement.

## TEST-INTERNAL_ENGINE-041: Public declarations still require annotations

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove an exported function without `@implement` produces the named
public-declaration finding.

**Oracle:** The test passes only when its assertions confirm an exported
function without `@implement` produces the named public-declaration finding.

## TEST-INTERNAL_ENGINE-042: Annotation must precede a supported declaration

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove an implementation annotation followed by an unrelated declaration is
reported as invalid placement rather than satisfying a later public function.

**Oracle:** The test passes only when its assertions confirm an
implementation annotation followed by an unrelated declaration is reported as invalid
placement rather than satisfying a later public function.

## TEST-INTERNAL_ENGINE-043: Private annotations still require documented SPECs

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove optional annotations on private declarations remain subject to doc/code
correspondence and cannot reference an undocumented SPEC.

**Oracle:** The test passes only when its assertions confirm optional
annotations on private declarations remain subject to doc/code correspondence and cannot
reference an undocumented SPEC.

## TEST-INTERNAL_ENGINE-044: Ignore scope suppresses missing annotation identifiers

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove annotation syntax inside an `idd:ignore` range does not create a missing
identifier finding, while scanning resumes after the range.

**Oracle:** The test passes only when its assertions confirm annotation
syntax inside an `idd:ignore` range does not create a missing identifier finding, while
scanning resumes after the range.

## TEST-INTERNAL_ENGINE-045: Ignore scope suppresses placement checks

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove an otherwise invalid annotation placement inside an ignore range is not
reported.

**Oracle:** The test passes only when its assertions confirm an otherwise
invalid annotation placement inside an ignore range is not reported.

## TEST-INTERNAL_ENGINE-046: Ignore scope suppresses consecutive annotation checks

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:**

Prove consecutive same-kind annotations inside an ignore range do not trigger
the one-line comma-list rule.

**Oracle:** The test passes only when its assertions confirm consecutive
same-kind annotations inside an ignore range do not trigger the one-line comma-list
rule.

## TEST-INTERNAL_ENGINE-047: Reciprocal graph construction contract

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-004`
- **Contracts:** `EngineLifecycle`, `RepositoryValidation`

**Purpose:**

Verify two reciprocal identifier references are both retained in graph
statistics without depending on the removed general backlink-warning feature.

**Oracle:** The test passes only when its assertions confirm two
reciprocal identifier references are both retained in graph statistics without depending
on the removed general backlink-warning feature.

## TEST-INTERNAL_ENGINE-048: Orphan rule contract

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-004`
- **Contracts:** `EngineLifecycle`, `RepositoryValidation`

**Purpose:**

Verify repository validation exposes the documented orphan rule for an
unlinked SPEC.

**Oracle:** The test passes only when its assertions confirm repository
validation exposes the documented orphan rule for an unlinked SPEC.

## TEST-INTERNAL_ENGINE-049: Preexisting errors survive validation

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-004`
- **Contracts:** `EngineLifecycle`, `RepositoryValidation`

**Purpose:**

Verify a structural error injected before validation keeps the result invalid
and remains present after the rule suite executes.

**Oracle:** The test passes only when its assertions confirm a
structural error injected before validation keeps the result invalid and remains present
after the rule suite executes.

## TEST-INTERNAL_ENGINE-050: Engine report envelope contract

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_ENGINE-004`
- **Contracts:** `EngineLifecycle`

**Purpose:**

Verify `BuildReport` returns non-nil output with stable tool/version identity
and a non-empty timestamp after an engine run.

**Oracle:** The test passes only when its assertions verify
`BuildReport` returns non-nil output with stable tool/version identity and a non-empty
timestamp after an engine run.

## TEST-INTERNAL_ENGINE-051: Link inference contract

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`
- **Contracts:** `GraphInterpretation`, `RepositoryValidation`

**Purpose:**

Verify SPEC, TEST, CONTRACT, and DESIGN source types map representative target
references to the documented relationship kinds.

**Oracle:** The test passes only when its assertions confirm SPEC, TEST,
CONTRACT, and DESIGN source types map representative target references to the documented
relationship kinds.

## TEST-INTERNAL_ENGINE-052: Public declarations across source languages

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:** Verify the public-declaration annotation rule consumes normalized
syntax-tree declarations for Go, TypeScript, TSX, JavaScript, C++, Java, and
Python.

**Oracle:** An annotated public declaration in every supported language
produces no public-annotation finding, while the corresponding unannotated
declaration produces that exact finding at its declaration line.

## TEST-INTERNAL_ENGINE-053: Test declarations across source languages

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`

**Purpose:** Verify behavior and contract TEST annotations bind to real test
declarations across all supported language profiles.

**Oracle:** Every correctly bound language fixture passes the test-annotation
rule, while a contract-test path using only `@test` is rejected for missing
`@test-contract`.

## TEST-INTERNAL_ENGINE-054: Derived Contract and Component graph policy

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-004`

**Purpose:** Verify named Contract coverage and typed Component dependency
edges participate in graph-wide validation without conflating lifecycle links
with architecture dependencies.

**Oracle:** A named Contract edge satisfies coverage, an absent edge does not,
and only reciprocal `depends_on` edges—not `supersedes` edges—produce a
Component dependency-cycle finding.

## Strategy

Engine tests use two evidence styles. Pure graph tests construct
`IdentifierSet` values and tune configuration flags to isolate a rule.
Filesystem rules create temporary Go and Markdown trees with absolute patterns,
then assert exact rule names, source paths, severities, and selected messages.
Contract tests verify lifecycle and externally meaningful orchestration without
claiming a production rule-plugin interface.

The suite intentionally aggregates several related test functions under
historical TEST identifiers `001` through `007`; each function remains focused
even when the record summarizes a rule family. Contract TEST identifiers are
separate from behavioral identifiers so `@test` and `@test-contract` never
compete for one declared kind.

Known exclusions are context cancellation, reuse of one engine for multiple
runs, concurrent access, unreadable-file propagation, recursive glob corner
cases, every combination of validation flags, graph aliasing, and semantic
adequacy of human prose. Repository-level validation and reporter tests cover
the full composition after these focused cases.
