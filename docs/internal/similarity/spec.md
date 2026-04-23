---
markers:
  - id: SPEC-INT_SIM-001
    name: TFIDF Structure
  - id: SPEC-INT_SIM-002
    name: Tokenize Function
  - id: SPEC-INT_SIM-003
    name: TF Computation
  - id: SPEC-INT_SIM-004
    name: IDF Computation
  - id: SPEC-INT_SIM-005
    name: TF-IDF Score
  - id: SPEC-INT_SIM-006
    name: TFIDF.Tokenize Method
  - id: SPEC-INT_SIM-007
    name: TFIDF.ComputeTF Method
  - id: SPEC-INT_SIM-008
    name: TFIDF.ComputeIDF Method
  - id: SPEC-INT_SIM-009
    name: TFIDF.Score Method
  - id: SPEC-INT_SIM-010
    name: TFIDF.CosineSimilarity Method

related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Specification (similarity)

## SPEC-INT_SIM-001: TFIDF Structure

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

**Tests:** `TEST-INT_SIM-001`, `TEST-INT_SIM-002`, `TEST-INT_SIM-003`, `TEST-INT_SIM-004`, `TEST-INT_SIM-005`, `TEST-INT_SIM-006`

**Related:** `SPEC-INT_SIM-002` through `SPEC-INT_SIM-005`

---

## SPEC-INT_SIM-002: Tokenize Function

**Status:** Done

**Contract:** `Tokenize`

**Design:** `SimilarityModule`

**Requirement:**

Tokenize text into lowercase alphanumeric tokens.

**Implementation:** `internal/similarity/tfidf.go`

**Tests:** `TEST-INT_SIM-001`

**Public Functions:**

## SPEC-INT_SIM-006: TFIDF.Tokenize

**Function Signature:**
`func (t *TFIDF) Tokenize(text string) []string`

**Purpose:** Tokenizes text into lowercase alphanumeric tokens, filtering stop words.

**Tests:** `TEST-INT_SIM-001`

---

## SPEC-INT_SIM-003: TF Computation

**Status:** Done

**Contract:** `ComputeTF`

**Design:** `SimilarityModule`

**Requirement:**

Compute term frequency for document tokens.

**Implementation:** `internal/similarity/tfidf.go`

**Tests:** `TEST-INT_SIM-002`

**Public Functions:**

## SPEC-INT_SIM-007: TFIDF.ComputeTF

**Function Signature:**
`func (t *TFIDF) ComputeTF(tokens []string) map[string]float64`

**Purpose:** Computes TF = (count of token) / (total tokens).

**Tests:** `TEST-INT_SIM-002`

---

## SPEC-INT_SIM-004: IDF Computation

**Status:** Done

**Contract:** `ComputeIDF`

**Design:** `SimilarityModule`

**Requirement:**

Compute inverse document frequency across corpus.

**Implementation:** `internal/similarity/tfidf.go`

**Tests:** `TEST-INT_SIM-003`

**Public Functions:**

## SPEC-INT_SIM-008: TFIDF.ComputeIDF

**Function Signature:**
`func (t *TFIDF) ComputeIDF(documents [][]string)`

**Purpose:** Computes IDF using formula log((N - df + 0.5) / (df + 0.5)).

**Tests:** `TEST-INT_SIM-003`

---

## SPEC-INT_SIM-005: TF-IDF Score

**Status:** Done

**Contract:** `Score`

**Design:** `SimilarityModule`

**Requirement:**

Compute similarity score between documents using TF-IDF.

**Implementation:** `internal/similarity/tfidf.go`

**Tests:** `TEST-INT_SIM-004`

**Public Functions:**

## SPEC-INT_SIM-009: TFIDF.Score

**Function Signature:**
`func (t *TFIDF) Score(docText, codeText string) float64`

**Purpose:** Computes similarity score between doc and code text.

**Tests:** `TEST-INT_SIM-005`

---

## SPEC-INT_SIM-010: TFIDF.CosineSimilarity Method

**Function Signature:**
`func CosineSimilarity(vec1, vec2 map[string]float64) float64`

**Purpose:** Computes cosine similarity between two TF-IDF vectors.

**Tests:** `TEST-INT_SIM-005`, `TEST-INT_SIM-006`

---

**Related:** `SPEC-INT_SIM-001` through `SPEC-INT_SIM-004`
