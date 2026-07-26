---
idd:
  version: "1.0"
  package: internal/similarity
  document: contract
---

# Contracts: internal/similarity

## Contract: TextScoring

**Guarantees:**

`TextScoring` compares two strings with a local TF-IDF/cosine pipeline.

```go
func Score(docText, codeText string) float64
func (t *TFIDF) Score(docText, codeText string) float64
```

### Inputs and output meaning

Inputs are borrowed immutable strings. The tokenizer lowercases them, retains
ASCII `[a-z0-9]+` tokens, and removes single-character and configured stop-word
tokens. The returned number is lexical-vector similarity for that transformed
pair; it is not a probability, confidence value, or proof that the behaviors
match.

Empty inputs or inputs reduced to zero-norm vectors return zero. The generic
cosine primitive can return values from negative one through one for arbitrary
vectors. Callers should compare the result with a policy threshold rather than
interpret an exact score as a stable semantic grade.

### State and concurrency

The top-level function allocates fresh state for each call. A caller-owned
`TFIDF` reuses its mutable IDF map and must not be mutated concurrently.
Neither form performs I/O, retains the original text, or returns an error.

### Compatibility

Changing token grammar, stop words, IDF formula, zero-IDF substitution, or
pair-corpus construction changes score distribution and therefore the meaning
of configured thresholds. Such changes require engine-level warning tests and
threshold review.

## Contract: WeightingPrimitives

**Guarantees:**

The lower-level primitives expose the current calculation for focused tests and
specialized callers:

- term frequency is occurrence count divided by total token count;
- an empty token slice returns an empty map;
- IDF uses `log((N - df + 0.5) / (df + 0.5))`;
- each document contributes at most one count to document frequency;
- TF-IDF multiplies TF by stored IDF, substituting `1.0` for an exact zero; and
- cosine returns zero if either vector has zero norm.

Returned maps are newly allocated. `ComputeIDF` mutates receiver state and does
not clear keys absent from the new corpus.

## Contract: TextNormalization

**Guarantees:**

`NormalizeText` lowercases Unicode text, replaces non-letter, non-digit, and
non-space runes with spaces, collapses whitespace, and trims the result. It is
pure and Unicode-aware, unlike the scoring tokenizer's ASCII-only grammar. No
contract states that scoring and standalone normalization are interchangeable.
