---
markers:
  - id: SPEC-INTERNAL_SIMILARITY-001
    name: TFIDF Structure
  - id: SPEC-INTERNAL_SIMILARITY-002
    name: Tokenize Function
  - id: SPEC-INTERNAL_SIMILARITY-003
    name: TF Computation
  - id: SPEC-INTERNAL_SIMILARITY-004
    name: IDF Computation
  - id: SPEC-INTERNAL_SIMILARITY-005
    name: TF-IDF Score
  - id: SPEC-INTERNAL_SIMILARITY-006
    name: TFIDF.Tokenize Method
  - id: SPEC-INTERNAL_SIMILARITY-007
    name: TFIDF.ComputeTF Method
  - id: SPEC-INTERNAL_SIMILARITY-008
    name: TFIDF.ComputeIDF Method
  - id: SPEC-INTERNAL_SIMILARITY-009
    name: TFIDF.Score Method
  - id: SPEC-INTERNAL_SIMILARITY-010
    name: TFIDF.CosineSimilarity Method

related_files:
  spec: docs/internal/similarity/spec.md
  contract: docs/internal/similarity/contract.md
  design: docs/internal/similarity/design.md
  testing: docs/internal/similarity/testing.md
---

# Specification (similarity)

## SPEC-INTERNAL_SIMILARITY-001: TFIDF Structure

**Status:** Done

**Contract:** `TFIDF`

**Design:** `SimilarityModule`

**Requirement:**

TFIDF provides TF-IDF based document similarity analysis.

**Implementation:** `internal/similarity/tfidf.go`

**Key Types:**

- `TFIDF` — TF-IDF indexer with IDF cache

**Acceptance Criteria:**

- [x] TFIDF stores corpus statistics
- [x] TFIDF computes TF, IDF, and TF-IDF
- [x] TFIDF scores similarity between documents

**Tests:** `TEST-INTERNAL_SIMILARITY-001`, `TEST-INTERNAL_SIMILARITY-002`, `TEST-INTERNAL_SIMILARITY-003`, `TEST-INTERNAL_SIMILARITY-004`, `TEST-INTERNAL_SIMILARITY-005`, `TEST-INTERNAL_SIMILARITY-006`

**Related:** `SPEC-INTERNAL_SIMILARITY-002` through `SPEC-INTERNAL_SIMILARITY-005`

---

## SPEC-INTERNAL_SIMILARITY-002: Tokenize Function

**Status:** Done

**Contract:** `Tokenize`

**Design:** `SimilarityModule`

**Requirement:**

Tokenize text into lowercase alphanumeric tokens.

**Implementation:** `internal/similarity/tfidf.go`

**Tests:** `TEST-INTERNAL_SIMILARITY-001`

**Public Functions:**

## SPEC-INTERNAL_SIMILARITY-006: TFIDF.Tokenize

**Function Signature:**
`func (t *TFIDF) Tokenize(text string) []string`

**Purpose:** Tokenizes text into lowercase alphanumeric tokens, filtering stop words.

**Tests:** `TEST-INTERNAL_SIMILARITY-001`

---

## SPEC-INTERNAL_SIMILARITY-003: TF Computation

**Status:** Done

**Contract:** `ComputeTF`

**Design:** `SimilarityModule`

**Requirement:**

Compute term frequency for document tokens.

**Implementation:** `internal/similarity/tfidf.go`

**Tests:** `TEST-INTERNAL_SIMILARITY-002`

**Public Functions:**

## SPEC-INTERNAL_SIMILARITY-007: TFIDF.ComputeTF

**Function Signature:**
`func (t *TFIDF) ComputeTF(tokens []string) map[string]float64`

**Purpose:** Computes TF = (count of token) / (total tokens).

**Tests:** `TEST-INTERNAL_SIMILARITY-002`

---

## SPEC-INTERNAL_SIMILARITY-004: IDF Computation

**Status:** Done

**Contract:** `ComputeIDF`

**Design:** `SimilarityModule`

**Requirement:**

Compute inverse document frequency across corpus.

**Implementation:** `internal/similarity/tfidf.go`

**Tests:** `TEST-INTERNAL_SIMILARITY-003`

**Public Functions:**

## SPEC-INTERNAL_SIMILARITY-008: TFIDF.ComputeIDF

**Function Signature:**
`func (t *TFIDF) ComputeIDF(documents [][]string)`

**Purpose:** Computes IDF using formula log((N - df + 0.5) / (df + 0.5)).

**Tests:** `TEST-INTERNAL_SIMILARITY-003`

---

## SPEC-INTERNAL_SIMILARITY-005: TF-IDF Score

**Status:** Done

**Contract:** `Score`

**Design:** `SimilarityModule`

**Requirement:**

Compute similarity score between documents using TF-IDF.

**Implementation:** `internal/similarity/tfidf.go`

**Tests:** `TEST-INTERNAL_SIMILARITY-004`

**Public Functions:**

## SPEC-INTERNAL_SIMILARITY-009: TFIDF.Score

**Function Signature:**
`func (t *TFIDF) Score(docText, codeText string) float64`

**Purpose:** Computes similarity score between doc and code text.

**Tests:** `TEST-INTERNAL_SIMILARITY-005`

---

## SPEC-INTERNAL_SIMILARITY-010: TFIDF.CosineSimilarity Method

**Function Signature:**
`func CosineSimilarity(vec1, vec2 map[string]float64) float64`

**Purpose:** Computes cosine similarity between two TF-IDF vectors.

**Tests:** `TEST-INTERNAL_SIMILARITY-005`, `TEST-INTERNAL_SIMILARITY-006`

---

**Related:** `SPEC-INTERNAL_SIMILARITY-001` through `SPEC-INTERNAL_SIMILARITY-004`
