// Package similarity provides TF-IDF text similarity scoring.

// Spec: docs/internal/similarity/spec.md
// Contract: docs/internal/similarity/contract.md
package similarity

import (
	"math"
	"regexp"
	"strings"
	"unicode"
)

// TFIDF provides TF-IDF based document similarity analysis.
//
// @implement SPEC-INTERNAL_SIMILARITY-001
type TFIDF struct {
	idf map[string]float64
}

// NewTFIDF creates a new TF-IDF indexer with an empty IDF cache.
//
// @implement SPEC-INTERNAL_SIMILARITY-001
func NewTFIDF() *TFIDF {
	return &TFIDF{idf: make(map[string]float64)}
}

// Tokenize text into lowercase alphanumeric tokens, filtering stop words.
//
// @implement SPEC-INTERNAL_SIMILARITY-001
func (t *TFIDF) Tokenize(text string) []string {
	text = strings.ToLower(text)
	reg := regexp.MustCompile(`[a-z0-9]+`)
	tokens := reg.FindAllString(text, -1)
	var filtered []string
	for _, token := range tokens {
		if len(token) > 1 && !isStopWord(token) {
			filtered = append(filtered, token)
		}
	}
	return filtered
}

func isStopWord(word string) bool {
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
		"for": true, "of": true, "with": true, "by": true, "from": true,
		"is": true, "are": true, "was": true, "were": true, "be": true,
		"been": true, "being": true, "have": true, "has": true, "had": true,
		"do": true, "does": true, "did": true, "will": true, "would": true,
		"could": true, "should": true, "may": true, "might": true, "must": true,
		"shall": true, "can": true, "need": true, "it": true, "its": true,
		"this": true, "that": true, "these": true, "those": true,
	}
	return stopWords[word]
}

// ComputeTF computes term frequency for document tokens.
//
// @implement SPEC-INTERNAL_SIMILARITY-002
func (t *TFIDF) ComputeTF(tokens []string) map[string]float64 {
	tf := make(map[string]float64)
	if len(tokens) == 0 {
		return tf
	}
	for _, token := range tokens {
		tf[token]++
	}
	for token := range tf {
		tf[token] /= float64(len(tokens))
	}
	return tf
}

// ComputeIDF computes inverse document frequency across corpus.
//
// @implement SPEC-INTERNAL_SIMILARITY-002
func (t *TFIDF) ComputeIDF(documents [][]string) {
	df := make(map[string]int)
	numDocs := float64(len(documents))
	for _, doc := range documents {
		seen := make(map[string]bool)
		for _, token := range doc {
			if !seen[token] {
				df[token]++
				seen[token] = true
			}
		}
	}
	for token, docFreq := range df {
		t.idf[token] = math.Log((numDocs - float64(docFreq) + 0.5) / (float64(docFreq) + 0.5))
	}
}

// ComputeTFIDF computes TF-IDF scores by combining term frequency with inverse document frequency.
//
// @implement SPEC-INTERNAL_SIMILARITY-003
func (t *TFIDF) ComputeTFIDF(tf map[string]float64) map[string]float64 {
	tfidf := make(map[string]float64)
	for token, tfVal := range tf {
		idfVal := t.idf[token]
		if idfVal == 0 {
			idfVal = 1.0
		}
		tfidf[token] = tfVal * idfVal
	}
	return tfidf
}

// CosineSimilarity computes cosine similarity between two TF-IDF vectors.
//
// @implement SPEC-INTERNAL_SIMILARITY-003
func CosineSimilarity(vec1, vec2 map[string]float64) float64 {
	var dotProduct, norm1, norm2 float64
	keys := make(map[string]bool)
	for k := range vec1 {
		keys[k] = true
	}
	for k := range vec2 {
		keys[k] = true
	}
	for k := range keys {
		v1 := vec1[k]
		v2 := vec2[k]
		dotProduct += v1 * v2
		norm1 += v1 * v1
		norm2 += v2 * v2
	}
	if norm1 == 0 || norm2 == 0 {
		return 0
	}
	return dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
}

// Score computes similarity score between doc and code text using TF-IDF.
//
// @implement SPEC-INTERNAL_SIMILARITY-004
func (t *TFIDF) Score(docText, codeText string) float64 {
	docTokens := t.Tokenize(docText)
	codeTokens := t.Tokenize(codeText)
	documents := [][]string{docTokens, codeTokens}
	t.ComputeIDF(documents)
	docTF := t.ComputeTF(docTokens)
	codeTF := t.ComputeTF(codeTokens)
	docTFIDF := t.ComputeTFIDF(docTF)
	codeTFIDF := t.ComputeTFIDF(codeTF)
	return CosineSimilarity(docTFIDF, codeTFIDF)
}

// Score computes similarity score between documents using TF-IDF.
//
// @implement SPEC-INTERNAL_SIMILARITY-004
func Score(docText, codeText string) float64 {
	tfidf := NewTFIDF()
	return tfidf.Score(docText, codeText)
}

// NormalizeText normalizes text by converting to lowercase and removing non-alphanumeric characters.
//
// @implement SPEC-INTERNAL_SIMILARITY-001
func NormalizeText(text string) string {
	text = strings.ToLower(text)
	text = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			return r
		}
		return ' '
	}, text)
	text = strings.Join(strings.Fields(text), " ")
	return text
}
