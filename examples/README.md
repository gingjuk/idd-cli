# IDD CLI examples

This directory contains one normative document example and one full
configuration reference. Both track the current CLI behavior.

## Self-describing module documents

[`self-describing-module-docs/`](self-describing-module-docs/) is a complete
four-file package document set. It is written as the contents of
`docs/internal/auth/`; therefore every file declares `internal/auth` as its
package identity.

The example keeps identity in each document's small `idd` frontmatter and keeps
Components, Contracts, SPECs, TESTs, relationships, guarantees, acceptance
evidence, and oracles in human-readable Markdown records. Source snippets use
declaration annotations as the code-to-document join:

The filename alone selects the role. Frontmatter contains `version`, current
`package`, and stable `namespace`; an `idd.document` field is invalid. Real
projects keep this four-file set per documentation unit. With `docs.units`, one
unit may cover several physical source packages. The files have no line-count limit and are not split into names
such as `design-auth.md` or `spec.part.md`.

If a real project contains one of those split names, the validation finding's
`suggested_fix` supplies a path-specific agent prompt to merge all authored
semantics into the canonical file, verify that nothing was lost, remove the
fragment, and rerun status plus project validation. `docs status` returns
equivalent guidance when the split prevents completion inspection.

- `@implement SPEC-INTERNAL_AUTH-001`;
- `@test TEST-INTERNAL_AUTH-001`;
- `@test-contract TEST-INTERNAL_AUTH-002`.

They intentionally contain no file-level `Spec:`, `Contract:`, or `Test:`
paths. The identifier is the join key used by idd-cli.

The example Component separately states Purpose, Ownership, Boundary, and
Decisions. Contract prose distinguishes durable guarantees from
non-guarantees, known limitations, and compatibility commitments so current
implementation accidents do not become promises.

Check that every required record and field is authored:

```bash
./bin/idd-cli docs status examples/self-describing-module-docs --format json
```

The collector regression test copies the files into their intended
`docs/internal/auth/` location and validates the full self-describing document
schema and reference graph:

```bash
CGO_ENABLED=1 go test ./internal/collector -run TestReferenceExampleDocuments
```

## Configuration reference

[`idd-config-example.yaml`](idd-config-example.yaml) lists every current
configuration key with the same values returned by `config.Default()`.
Configuration owns discovery and validation policy; it is not a marker catalog
and does not store Component, Contract, SPEC, TEST, or source-file path
relationships.

The deep-equality regression test prevents defaults and the example from
drifting:

```bash
CGO_ENABLED=1 go test ./internal/config -run TestDefaultMatchesExampleConfiguration
```
