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

- `@implement SPEC-INTERNAL_AUTH-001`;
- `@test TEST-INTERNAL_AUTH-001`;
- `@test-contract TEST-INTERNAL_AUTH-002`.

They intentionally contain no file-level `Spec:`, `Contract:`, or `Test:`
paths. The identifier is the join key used by idd-cli.

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
