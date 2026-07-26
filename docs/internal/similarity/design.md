---
idd:
  version: "1.0"
  package: internal/similarity
---

# Design: internal/similarity

## Component: SemanticSimilarity

**Purpose:**

`SemanticSimilarity` supplies the engine's advisory comparison between a
documentation description and the comment attached to an implementation
annotation. It reduces both texts to weighted token vectors and returns cosine
similarity; the engine decides whether the configured threshold warrants a
warning.

This component is intentionally heuristic. It can highlight vocabulary drift,
but it cannot establish behavioral equivalence, documentation completeness, or
implementation correctness. Keeping the score advisory prevents a short source
comment or domain synonym from becoming a false structural failure.

### Responsibilities

- normalize input into lowercase ASCII word and number tokens;
- remove single-character tokens and a fixed set of common English stop words;
- calculate normalized term frequency;
- calculate pair-corpus inverse document frequency;
- combine weights and compute cosine similarity; and
- offer both a reusable stateful scorer and a fresh one-shot helper.

`NormalizeText` provides a separate Unicode-aware punctuation normalization
utility. The scoring pipeline does not currently call it; tokenization uses its
own ASCII regular expression.

### State, concurrency, and failure boundaries

A `TFIDF` instance owns a mutable `idf` map. `ComputeIDF` updates entries for
the supplied corpus but does not clear entries from earlier corpora. The score
path only reads tokens present in its current pair, so stale unrelated entries
normally have no effect, but the instance is not safe for concurrent mutation.
The package performs no I/O and returns no errors; empty or zero-norm inputs
produce a score of zero.

### Decisions and trade-offs

The implementation uses a two-document corpus for each comparison. This keeps
the result local and deterministic without maintaining a project-wide index,
at the cost of weaker statistical meaning than corpus-trained TF-IDF. The IDF
formula includes `0.5` smoothing and can produce negative weights for terms
present in most documents. When an IDF value is exactly zero,
`ComputeTFIDF` substitutes `1.0` so a term is not erased entirely.

A compact built-in stop-word set avoids an external language dependency. It
also makes the heuristic English-oriented and unsuitable as a language-neutral
semantic model.

## Architecture

```text
doc text ── Tokenize ── ComputeTF ─┐
                                   ├─ pair IDF ─ TF-IDF ─┐
code text ─ Tokenize ─ ComputeTF ──┘                     ├─ cosine score
                                                         ┘
```

The engine consumes only the numeric score and configured threshold. It owns
warning severity, source locations, ignore rules, and the decision to skip
missing or function-locator-only descriptions.

## Package Layout

`tfidf.go` contains the full scoring pipeline because its stages share the
private IDF state and are small pure numerical operations. Moving threshold
policy or validation findings into this package would couple mathematical
scoring to IDD workflow decisions.

## Function Composition

The top-level `Score` creates a fresh `TFIDF` and delegates to its `Score`
method. The method tokenizes both texts, computes IDF over that pair, computes
term-frequency and weighted vectors, then calls `CosineSimilarity`.
`NormalizeText` remains an independently callable normalization operation.

## Dependencies

Only standard-library math, regular-expression, string, and Unicode packages
are used. There is no model download, network access, locale data, or persistent
index.

## Testability Hooks

Each numerical stage is directly callable, allowing tests to isolate
tokenization, TF, IDF, weighting, cosine geometry, end-to-end scoring, and
normalization. Table-driven vector tests use approximate comparisons where
floating-point arithmetic is involved.

Tests do not establish multilingual quality, threshold precision/recall,
concurrent safety, or semantic correctness. Those are explicit limitations of
the heuristic and require evaluation data rather than more unit assertions.
