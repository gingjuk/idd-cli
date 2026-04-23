---
related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Design (similarity)

**Status:** Done

## Architecture

The similarity module provides text similarity scoring:

1. **TF-IDF Algorithm**: Standard TF-IDF with cosine similarity
2. **IDF Smoothing**: Added 0.5 smoothing to prevent division by zero
3. **Pre-computed IDF**: IDF cached after first computation for efficiency

## Package Layout

```text
internal/similarity/
├── tfidf.go          # TF-IDF implementation
└── similarity.go     # Similarity scoring functions
```

## Function Composition

1. **Score()** - Computes similarity between two text strings
2. **ComputeTF()** - Calculates term frequency for a document
3. **ComputeIDF()** - Calculates inverse document frequency
4. **CosineSimilarity()** - Computes angle between term vectors

## Testability Hooks

- Pure functions with no I/O dependencies
- Scores are deterministic for given inputs
- Threshold-based comparisons enable easy test assertions

## Dependencies

- `internal/model` - For describing document structure
- No external dependencies (self-contained)
