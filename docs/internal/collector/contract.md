---
idd:
  version: "1.0"
  package: internal/collector
  document: contract
---

# Contracts: internal/collector

## Contract: IDDDocumentSet

`ParseIDDDocument` accepts an `idd` frontmatter block containing only
`version`, `package`, and `document`. Unknown fields, semantic YAML catalogs,
and malformed syntax are rejected. A document without an `idd` block is left
to legacy parsing.

The fixed filenames own disjoint declarations:

- `design.md` owns `## Component: <name>` records;
- `contract.md` owns `## Contract: <name>` records;
- `spec.md` owns `## SPEC-...: <title>` records with `Design`, `Contract`, and
  `Requirement` fields;
- `testing.md` owns `## TEST-...: <title>` records with `Kind`, `Covers`, and
  `Purpose` fields.

Records are bounded by level-two headings. Their required fields remain short,
while the rest of each section is unrestricted human-authored Markdown.

Validation reports the exact Markdown record and field line for:

- package or document-role values that disagree with the Markdown path;
- identifiers that do not use the package-derived module;
- missing self-describing documents;
- incomplete or placeholder SPEC and TEST fields;
- unresolved design, contract, and coverage references;
- duplicate identifiers or fixed fields;
- invalid TEST kinds.

`MarshalIDDDocument` emits canonical minimal identity frontmatter and preserves
its Markdown body. `RepairDocuments` normalizes only structural identity. A
file target writes only that file; a directory target may create missing
structural documents. Repair never reflows prose or synthesizes semantic
records.

## Document initialization

`InitDocuments` requires an existing project-relative package directory. It
creates the four self-describing Markdown documents or prepends structural
metadata to existing plain narratives. Existing non-IDD frontmatter keys are
merged into the same frontmatter block and preserved.

Initialization performs a complete preflight before writing. If existing
Markdown contains legacy markers or `related_files`, it returns an error so a
new document set cannot silently change that package's parsing mode. A
package-local central `idd.yaml` is also rejected for explicit migration.

## Contract: DocCollector

`NewDocCollector` returns a concrete `DocCollector`; there is no shared
collector interface. `DocCollector.Collect` accepts a file or directory and
returns:

1. all documentation-origin identifiers;
2. structural validation findings that do not stop collection;
3. a traversal error only when collection itself cannot continue.

Self-describing sets are discovered before legacy Markdown. Their bodies may
use declared identifiers as detail headings but may not repeat legacy metadata.
If none of the four fixed files has an `idd` block, the existing frontmatter
parser remains authoritative.

TEST `Covers` fields generate both TEST-to-SPEC and reverse SPEC-to-TEST graph
links. Authors declare the relationship only once.

## Contract: CodeCollector

`NewCodeCollector` returns a concrete `CodeCollector`.
`CodeCollector.Collect` scans supported Go, TypeScript, and JavaScript files,
honors ignore scopes, and emits code-origin identifiers from `@implement`,
`@test`, and `@test-contract`.

TEST annotations preserve the semantic kind:

- `@test` becomes kind `test`;
- `@test-contract` becomes kind `contract`.

The engine compares this kind with the `testing.md` record after documentation
and code identifiers are merged.

## Contract: LegacyFrontmatter

Legacy mode parses `markers` and `related_files`, validates marker headings and
backtick formatting, and maps SPEC, TEST, CONTRACT, and DESIGN identifiers to
their fixed narrative filenames. Fenced code does not define frontmatter or
document markers.

## Failure behavior

| Condition | Result |
| --- | --- |
| Missing target | Empty identifier set |
| Ignored path during directory discovery | No collected identifiers or findings |
| Invalid IDD frontmatter | Structured finding with Markdown line; other inputs continue |
| Invalid legacy frontmatter | Legacy structural finding; other files continue |
| Missing role document | IDD document-set finding |
| Package-local central catalog | Migration finding; it is never treated as a marker source |
| Unsafe initialization target | Error before any file is written |
| Malformed document repair input | Error without replacing the document |

**Related Specs:** `SPEC-INTERNAL_COLLECTOR-001`,
`SPEC-INTERNAL_COLLECTOR-002`, `SPEC-INTERNAL_COLLECTOR-006`,
`SPEC-INTERNAL_COLLECTOR-009`, `SPEC-INTERNAL_COLLECTOR-010`,
`SPEC-INTERNAL_COLLECTOR-011`, `SPEC-INTERNAL_COLLECTOR-012`,
`SPEC-INTERNAL_COLLECTOR-013`
