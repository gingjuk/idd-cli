---
related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Contract (similarity)

**Status:** Done

**Requirement:**

The similarity module provides TF-IDF based document similarity analysis.

**Key Contracts:**

- Tokenize must filter stop words and lowercase tokens
- ComputeTF must return normalized term frequencies
- ComputeIDF must use proper IDF formula with smoothing
- CosineSimilarity must return value between 0.0 and 1.0

**Implementation:** `internal/similarity/tfidf.go`

**Acceptance Criteria:**

- [x] Documents are vectorized using TF-IDF
- [x] Similarity scores range from 0.0 to 1.0
- [x] Configurable similarity threshold
- [x] Similar files are flagged for review

**Related:** `SPEC-INT_SIM-001` through `SPEC-INT_SIM-005`
