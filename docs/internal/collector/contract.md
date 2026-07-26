---
idd:
  version: "1.0"
  package: internal/collector
  document: contract
---

# Contracts: internal/collector

## Contract: IDDDocumentSet

**Guarantees:**

`ParseIDDDocument` accepts an `idd` frontmatter block containing only
`version`, `package`, and `document`. Unknown fields, semantic YAML catalogs,
and malformed syntax are rejected. A document without an `idd` block is left
to legacy parsing.

The fixed filenames own disjoint declarations:

- `design.md` owns `## Component: <name>` records;
- `contract.md` owns `## Contract: <name>` records;
- `spec.md` owns `## SPEC-...: <title>` records with `Design`, `Contract`, and
  `Requirement` and `Acceptance` fields;
- `testing.md` owns `## TEST-...: <title>` records with `Kind`, `Covers`, and
  `Purpose` and `Oracle` fields. Contract TEST records additionally own
  `Contracts`.

Every Component has a `Purpose` paragraph and every Contract has a
`Guarantees` paragraph. Optional `Status`, `Supersedes`, `Deprecated by`,
`Concerns`, and Component `Depends on` fields remain part of the owning
Markdown record. Selected concerns require same-named level-three sections.
Component and Contract names may contain ordinary human-readable words but
must not contain `#` or `,`, which are reserved reference separators.

Records are bounded by level-two headings. Their fixed identity and
relationship fields remain small, while the rest of each section is
unrestricted human-authored Markdown whose depth follows the subject rather
than a line-count target.

Validation reports the exact Markdown record and field line for:

- package or document-role values that disagree with the Markdown path;
- identifiers that do not use the package-derived module;
- missing self-describing documents;
- incomplete or placeholder role fields and generated scaffold slots;
- unresolved design, contract, coverage, Component dependency, and lifecycle
  references;
- duplicate identifiers or fixed fields;
- invalid TEST kinds, lifecycle graphs, concern sections, or contract TEST
  mappings.

`MarshalIDDDocument` emits canonical minimal identity frontmatter and preserves
its Markdown body. `RepairDocuments` normalizes only structural identity. A
file target writes only that file; a directory target may create missing
structural documents. Repair never reflows prose or synthesizes semantic
records.

Generated skeletons contain stable `idd:scaffold` slot markers and are invalid
until authored. `InspectDocumentCompletion` applies the same role schema used
by generation and normal validation. It accepts one role file or a directory
tree and returns schema `idd.document_status.v1`, status, and a deterministic
`incomplete_slots` list containing file, line, role, slot, and reason. A slot
is complete only when its marker is absent and its bounded content is
non-empty and non-placeholder. Once a record exists, the list also contains
stable record-field slots for every required scalar or list field that is
missing or still a placeholder. Required-field validation and status
inspection consume the same role schema.

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

**Guarantees:**

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

**Guarantees:**

`NewCodeCollector` returns a concrete `CodeCollector`.
`CodeCollector.CollectWithErrors` scans configured files whose extensions have
pinned Tree-sitter grammars:

- Go `.go`;
- TypeScript/TSX `.ts`, `.tsx`;
- JavaScript/JSX `.js`, `.jsx`;
- C++ `.cpp`, `.cc`, `.cxx`, `.hpp`, `.hh`, `.hxx`;
- Java `.java`;
- Python `.py`.

It reads only real syntax-tree comment nodes, honors standalone ignore ranges,
and binds `@implement`, `@test`, and `@test-contract` to adjacent declarations.
Annotation-like strings, comments inside function bodies, and detached comments
do not become code-origin identifiers. A configured extension without a pinned
grammar and a supported file whose syntax tree contains errors both return a
`source-parse` finding. Neither path uses regex fallback or partial annotation
evidence.

The configured `spec`, `test`, and `test_contract` prefix values replace those
canonical spellings at the lexical boundary while the normalized analysis
still exposes stable `implement`, `test`, and `test-contract` semantic kinds.

`CodeCollector.Collect` is the compatibility wrapper that returns identifiers
and traversal errors. `Analyses` returns a copy of the normalized per-file
syntax-tree analyses from the latest collection so the engine can validate
declaration visibility, test classification, annotation placement, identifier
syntax, and duplicate adjacent annotations without reparsing.

TEST annotations preserve the semantic kind:

- `@test` becomes kind `test`;
- `@test-contract` becomes kind `contract`.

The engine compares this kind with the `testing.md` record after documentation
and code identifiers are merged.

Public declarations are language-specific: Go export rules include the
receiver type for methods, Java uses access modifiers, TypeScript and
JavaScript use export visibility, and Python treats underscore-prefixed names
as private. C++ declarations are considered part of the public source surface
unless they have internal linkage, live in an anonymous namespace, or are under
a non-public class access section. Java interface methods retain their
language-defined public default. Recognized test
declarations follow each ecosystem's filename, name, annotation, decorator, or
test-call syntax as represented by its grammar.

## Contract: LegacyFrontmatter

**Guarantees:**

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
| Generated slot remains | `idd-document-incomplete` finding and `docs status` work item |
| Configured source has no pinned grammar | `source-parse` finding; remove the pattern or select a supported language |
| Supported source has invalid syntax | `source-parse` finding; no source evidence from that file |
| Invalid legacy frontmatter | Legacy structural finding; other files continue |
| Missing role document | IDD document-set finding |
| Package-local central catalog | Migration finding; it is never treated as a marker source |
| Unsafe initialization target | Error before any file is written |
| Malformed document repair input | Error without replacing the document |

## Semantic review boundary

The collector can verify that a record has the required fields, a meaningful
non-placeholder narrative, valid references, and stable source locations. It
does not decide whether an architecture explains enough trade-offs, a contract
captures every externally observable effect, or a SPEC expresses the right
requirement.

No minimum word count is part of this contract. Generated guidance and the
paired Skill define the questions authors must answer, while human review
judges whether the answers are sufficient for the actual complexity.
