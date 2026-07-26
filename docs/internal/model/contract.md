---
idd:
  version: "1.0"
  package: internal/model
  document: contract
---

# Contracts: internal/model

## Contract: IdentifierVocabulary

**Guarantees:**

`IdentifierType` supports `SPEC`, `CONTRACT`, `TEST`, and `DESIGN`.
`ParseIdentifierType` accepts those names case-insensitively and returns an
error containing the unknown input for every other value.

`Origin` distinguishes `doc` and `code`. It is descriptive metadata rather than
an access-control boundary. Model constructors default identifiers to document
origin; collectors must explicitly set code origin for source annotations.

## Contract: IdentifierRecord

**Guarantees:**

An `Identifier` stores canonical ID and type plus human title/description,
source path and line, raw lexical reference, legacy forward links, explicitly
typed links, origin, optional TEST kind, and whether the node was derived from
a package-scoped Component or Contract. Constructors return a non-nil pointer
with non-nil empty `Links` and `TypedLinks` slices, raw reference equal to the
ID, document origin, and `Derived == false`.

The describe-aware constructor additionally stores the supplied description.
`AddLink` appends an untyped target without validation or deduplication.
`AddTypedLink` appends an `IdentifierLink` containing target and exact
`LinkType`, also without validation or deduplication. `SetOrigin` replaces the
origin. No method validates identifier grammar, source existence, link targets,
derived-key shape, or TEST kind.

## Contract: IdentifierCollection

**Guarantees:**

`IdentifierSet` indexes a slice of observations per ID and appends each
observation to the public slice for its type.

- `Get` returns the first observation and false when absent.
- `GetAll` returns the stored observation slice or nil when absent.
- `Has` reports whether at least one observation exists.
- `Count` counts unique ID strings, not observations.
- `All` returns one first observation per ID sorted by ID.
- `AllIdentifiers` returns every observation sorted by ID.
- `ByOrigin` returns every matching observation sorted by ID.
- `HasOrigin` searches observations for one ID and origin.
- `Merge` appends all observations from another set without cloning.

Duplicate document/code group detection reports one representative per source
directory only when an ID spans multiple directories. Returned pointers and
slices are borrowed mutable data. The set is not concurrency-safe and rejects
no nil identifiers; callers must supply valid pointers.

## Contract: SourceAnnotation

**Guarantees:**

`Annotation` retains parsed type, target reference, source line, raw annotation,
nearby context, and optional declaration comment. Constructors copy those
values without validation.

`ToIdentifier` creates a document-origin identifier with the annotation's ID,
type, source, and line; replaces `RawRef` with the full raw annotation; and
copies a non-empty function comment into `Describe`. It does not copy context
into the identifier or set code origin.

## Contract: ValidationResult

**Guarantees:**

`ValidationError.Error` formats only rule and message as `[rule] message`.
Structured source, link, and code fields remain available for reporters.

`NewValidationResult` returns non-nil empty error and warning slices with
`Valid == false`. `AddError` appends all evidence fields and sets valid false.
`AddWarning` appends without changing validity. `Sort` orders errors and
warnings independently by rule and then message; ties have no further ordering
guarantee.

Stats, complete reports, config summaries, and graph snapshots are plain JSON
values. The model does not calculate or validate them.

## Contract: Relationship

**Guarantees:**

`LinkType` defines implementation, testing, contract-testing, reference,
annotation, contract, reverse-contract, Component dependency, reverse
dependency, supersession, and deprecation relationships.
`ReverseLinkType` swaps `tests`/`implements`,
`contract`/`contract_implements`, `depends_on`/`depended_by`, and
`supersedes`/`deprecated_by`. Contract-test, reference, and annotation links
are self-reversing; unknown values are returned unchanged.

`IdentifierLink` is an explicitly typed outgoing target owned by an
`Identifier`. It carries no source position of its own because the owning
identifier supplies that evidence when the graph materializes the edge.

`NewLink` returns a directed value with supplied endpoints, type, source, and
line. It does not require endpoint existence, normalize direction, or verify a
reverse link.

## Contract: LLMFindingReport

**Guarantees:**

The LLM wire model contains schema and status, summary counts and rule groups,
and self-contained findings with location, primary identifier, expected and
actual state, repair guidance, and related identifiers.

JSON tags define the compatibility surface. Optional values use `omitempty`
where declared; required fields may still contain zero values because model
types perform no validation. Reporter owns schema version selection, grouping,
enrichment, ordering, and output formatting.
