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

**Design:** `SimilarityModule`

**Contract:** `TFIDF`

**Requirement:**

TFIDF provides TF-IDF based document similarity analysis.

**Tests:** `TEST-INTERNAL_SIMILARITY-001`, `TEST-INTERNAL_SIMILARITY-002`, `TEST-INTERNAL_SIMILARITY-003`, `TEST-INTERNAL_SIMILARITY-004`, `TEST-INTERNAL_SIMILARITY-005`, `TEST-INTERNAL_SIMILARITY-006`

**Status:** Done

**Implementation:** `internal/similarity/tfidf.go`

**Key Types:**

- `TFIDF` — TF-IDF indexer with IDF cache

**Acceptance Criteria:**

- [x] TFIDF stores corpus statistics
- [x] TFIDF computes TF, IDF, and TF-IDF
- [x] TFIDF scores similarity between documents

**Related:** `SPEC-INTERNAL_SIMILARITY-002` through `SPEC-INTERNAL_SIMILARITY-005`

## SPEC-INTERNAL_SIMILARITY-002: Tokenize Function

**Design:** `SimilarityModule`

**Contract:** `Tokenize`

**Requirement:**

Tokenize text into lowercase alphanumeric tokens.

**Tests:** `TEST-INTERNAL_SIMILARITY-001`

**Status:** Done

**Implementation:** `internal/similarity/tfidf.go`

**Public Functions:**

---

## SPEC-INTERNAL_SIMILARITY-006: TFIDF.Tokenize

**Design:** `SimilarityModule`

**Contract:** `Tokenize`

**Requirement:**

Tokenize text into lowercase alphanumeric tokens while filtering stop words.

**Tests:** `TEST-INTERNAL_SIMILARITY-001`

**Function Signature:**
`func (t *TFIDF) Tokenize(text string) []string`

**Purpose:** Tokenizes text into lowercase alphanumeric tokens, filtering stop words.

---

## SPEC-INTERNAL_SIMILARITY-003: TF Computation

**Design:** `SimilarityModule`

**Contract:** `ComputeTF`

**Requirement:**

Compute term frequency for document tokens.

**Tests:** `TEST-INTERNAL_SIMILARITY-002`

**Status:** Done

**Implementation:** `internal/similarity/tfidf.go`

**Public Functions:**

## SPEC-INTERNAL_SIMILARITY-007: TFIDF.ComputeTF

**Design:** `SimilarityModule`

**Contract:** `ComputeTF`

**Requirement:**

Compute term frequency as count of token divided by total token count.

**Tests:** `TEST-INTERNAL_SIMILARITY-002`

**Function Signature:**
`func (t *TFIDF) ComputeTF(tokens []string) map[string]float64`

**Purpose:** Computes TF = (count of token) / (total tokens).

---

## SPEC-INTERNAL_SIMILARITY-004: IDF Computation

**Design:** `SimilarityModule`

**Contract:** `ComputeIDF`

**Requirement:**

Compute inverse document frequency across corpus.

**Tests:** `TEST-INTERNAL_SIMILARITY-003`

**Status:** Done

**Implementation:** `internal/similarity/tfidf.go`

**Public Functions:**

## SPEC-INTERNAL_SIMILARITY-008: TFIDF.ComputeIDF

**Design:** `SimilarityModule`

**Contract:** `ComputeIDF`

**Requirement:**

Compute inverse document frequency across tokenized documents.

**Tests:** `TEST-INTERNAL_SIMILARITY-003`

**Function Signature:**
`func (t *TFIDF) ComputeIDF(documents [][]string)`

**Purpose:** Computes IDF using formula log((N - df + 0.5) / (df + 0.5)).

---

## SPEC-INTERNAL_SIMILARITY-005: TF-IDF Score

**Design:** `SimilarityModule`

**Contract:** `Score`

**Requirement:**

Compute similarity score between documents using TF-IDF.

**Tests:** `TEST-INTERNAL_SIMILARITY-004`

**Status:** Done

**Implementation:** `internal/similarity/tfidf.go`

**Public Functions:**

## SPEC-INTERNAL_SIMILARITY-009: TFIDF.Score

**Design:** `SimilarityModule`

**Contract:** `Score`

**Requirement:**

Compute a similarity score between documentation text and code text.

**Tests:** `TEST-INTERNAL_SIMILARITY-005`

**Function Signature:**
`func (t *TFIDF) Score(docText, codeText string) float64`

**Purpose:** Computes similarity score between doc and code text.

---

## SPEC-INTERNAL_SIMILARITY-010: TFIDF.CosineSimilarity Method

**Design:** `SimilarityModule`

**Contract:** `CosineSimilarity`

**Requirement:**

Compute cosine similarity between two TF-IDF vectors.

**Tests:** `TEST-INTERNAL_SIMILARITY-005`, `TEST-INTERNAL_SIMILARITY-006`

---

**Function Signature:**
`func CosineSimilarity(vec1, vec2 map[string]float64) float64`

**Purpose:** Computes cosine similarity between two TF-IDF vectors.

**Related:** `SPEC-INTERNAL_SIMILARITY-001` through `SPEC-INTERNAL_SIMILARITY-004`
