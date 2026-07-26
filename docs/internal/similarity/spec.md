---
idd:
  version: "1.0"
  package: internal/similarity
  document: spec
---

# Specifications: internal/similarity

## SPEC-INTERNAL_SIMILARITY-001: Deterministic text preparation

- **Design:** `SemanticSimilarity`
- **Contract:** `TextNormalization`

**Requirement:**

Similarity processing must provide deterministic construction and text
preparation: scorer state starts with an initialized IDF map; scoring tokens are
lowercase ASCII words or numbers longer than one character with common English
stop words removed; standalone normalization preserves Unicode letters and
digits while replacing punctuation and collapsing whitespace.

### Boundaries and rationale

The two normalization paths are deliberately documented separately because
they do not produce identical text. Tokenization is optimized for the current
English-oriented similarity heuristic, while `NormalizeText` is a general
cleanup helper. Neither performs stemming, synonym expansion, language
detection, or semantic embedding.

### Acceptance evidence

**Acceptance:**

Tokenization cases prove case folding, punctuation splitting, stop-word and
single-character filtering, and empty input. Normalization cases prove
punctuation replacement, whitespace collapse, and lowercase output.

## SPEC-INTERNAL_SIMILARITY-002: Local TF-IDF weighting

- **Design:** `SemanticSimilarity`
- **Contract:** `WeightingPrimitives`

**Requirement:**

The scorer must convert tokens into normalized term frequencies, compute
smoothed inverse document frequencies across the supplied documents, and
multiply those values into a fresh weighted vector.

### Required numerical behavior

Repeated tokens contribute proportionally to total token count. A token is
counted once per document for document frequency. The smoothed logarithmic IDF
may be negative for ubiquitous terms. An exact stored zero is treated as one
during weighting so it does not eliminate the term.

### State and edge cases

Empty term input yields an empty map. Computing IDF mutates the receiver and
does not reset unrelated historical keys. No guarantee is made for concurrent
calls on one instance.

### Acceptance evidence

**Acceptance:**

Unit tests assert exact simple TF and TF-IDF values and the sign behavior of
common versus rare IDF terms.

## SPEC-INTERNAL_SIMILARITY-003: Cosine comparison of weighted vectors

- **Design:** `SemanticSimilarity`
- **Contract:** `WeightingPrimitives`

**Requirement:**

Vector comparison must calculate the dot product over the union of keys,
normalize by both vector magnitudes, and return zero when either magnitude is
zero.

### Edge cases

The primitive accepts arbitrary signed weights, so opposite vectors may return
negative similarity. Missing keys behave as zero and input maps are not
modified.

### Acceptance evidence

**Acceptance:**

Table-driven cases cover identical, orthogonal, opposite, empty, partially
overlapping, and one-sided zero-norm vectors with floating-point tolerances.

## SPEC-INTERNAL_SIMILARITY-004: Advisory pairwise similarity score

- **Design:** `SemanticSimilarity`
- **Contract:** `TextScoring`

**Requirement:**

End-to-end scoring must build a two-text corpus, apply the documented
tokenization and weighting stages, and return cosine similarity. The top-level
helper must isolate calls with a fresh scorer.

### Implementation boundary

The score is advisory input to the engine's consistency warning. This
specification does not define a universal passing threshold, guarantee
multilingual behavior, or allow a score to replace human review. Threshold
configuration and the cases skipped by consistency validation belong to the
engine and config packages.

### Acceptance evidence

**Acceptance:**

Related English phrases must score substantially higher than unrelated phrases
in the representative fixtures, and two empty inputs must produce zero.
