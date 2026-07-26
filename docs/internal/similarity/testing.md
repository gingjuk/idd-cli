---
idd:
  version: "1.0"
  package: internal/similarity
  document: testing
---

# Testing: internal/similarity

## TEST-INTERNAL_SIMILARITY-001: Tokenization boundaries

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_SIMILARITY-001`
- **Contracts:** `TextScoring`, `cmd/idd-cli#TFIDF`

**Purpose:**

Prove lowercase ASCII token extraction, punctuation splitting, stop-word and
single-character filtering, and empty input. Table-driven expected token counts
make the transformation observable without asserting regular-expression
internals.

**Oracle:** The test passes only when its assertions confirm lowercase
ASCII token extraction, punctuation splitting, stop-word and single-character filtering,
and empty input. Table-driven expected token counts make the transformation observable
without asserting regular-expression internals.

## TEST-INTERNAL_SIMILARITY-002: Term frequency and representative scoring

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_SIMILARITY-002`, `SPEC-INTERNAL_SIMILARITY-004`
- **Contracts:** `WeightingPrimitives`

**Purpose:**

Prove normalized term counts, empty-token behavior, and the practical separation
between related authentication prose and unrelated database prose. Exact TF
values and broad score bounds are the respective oracles.

**Oracle:** The test passes only when its assertions confirm normalized
term counts, empty-token behavior, and the practical separation between related
authentication prose and unrelated database prose. Exact TF values and broad score
bounds are the respective oracles.

## TEST-INTERNAL_SIMILARITY-003: IDF behavior and standalone normalization

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_SIMILARITY-001`, `SPEC-INTERNAL_SIMILARITY-002`

**Purpose:**

Prove rare versus ubiquitous IDF sign behavior and verify lowercase,
punctuation-replacing, whitespace-collapsing normalization. The shared TEST ID
represents text-weight preparation evidence across two focused test functions.

**Oracle:** The test passes only when its assertions confirm rare versus
ubiquitous IDF sign behavior and verify lowercase, punctuation-replacing,
whitespace-collapsing normalization. The shared TEST ID represents text-weight
preparation evidence across two focused test functions.

## TEST-INTERNAL_SIMILARITY-004: TF-IDF multiplication

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_SIMILARITY-002`
- **Contracts:** `TextNormalization`

**Purpose:**

Prove that known term frequencies are multiplied by receiver IDF values into a
fresh vector using simple exact inputs.

**Oracle:** The test passes only when its assertions confirm known
term frequencies are multiplied by receiver IDF values into a fresh vector using simple
exact inputs.

## TEST-INTERNAL_SIMILARITY-005: Cosine vector geometry

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_SIMILARITY-003`

**Purpose:**

Prove identical, orthogonal, opposite, empty, partial-overlap, and one-sided
zero-norm outcomes. Approximate comparisons isolate expected floating-point
rounding while exact zero cases detect invalid division.

**Oracle:** The test passes only when its assertions confirm identical,
orthogonal, opposite, empty, partial-overlap, and one-sided zero-norm outcomes.
Approximate comparisons isolate expected floating-point rounding while exact zero cases
detect invalid division.

## TEST-INTERNAL_SIMILARITY-006: Empty pair scoring

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_SIMILARITY-004`

**Purpose:**

Prove that two empty strings produce a zero score rather than NaN, infinity, or
an error.

**Oracle:** The test passes only when its assertions confirm two
empty strings produce a zero score rather than NaN, infinity, or an error.

## Strategy

The suite tests each mathematical stage independently and then exercises the
composed score. Fixtures are immutable slices, maps, and strings; there is no
filesystem, time, randomness, or network dependency.

The tests intentionally use representative bounds rather than claiming the
heuristic is semantically calibrated. They do not cover multilingual text,
concurrent receiver use, stale IDF entries across reused corpora, threshold
quality, or large-corpus numerical stability.
