// Package similarity provides testing utilities for the similarity module.

// Spec: docs/internal/similarity/spec.md
// Test: docs/internal/similarity/testing.md
package similarity

import (
	"testing"
)

// @test TEST-INT_SIM-001
func TestTFIDF_Tokenize(t *testing.T) {
	tfidf := NewTFIDF()

	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{"basic", "Hello world hello", 3},
		{"stopwords filtered", "The quick brown fox", 3},
		{"numbers kept", "test123 test456", 2},
		{"short words filtered", "a an the and or but", 0},
		{"empty", "", 0},
		{"punctuation", "hello, world! test-case", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := tfidf.Tokenize(tt.text)
			if len(tokens) != tt.expected {
				t.Errorf("Tokenize(%q) got %d tokens, want %d", tt.text, len(tokens), tt.expected)
			}
		})
	}
}

// @test TEST-INT_SIM-002
func TestTFIDF_ComputeTF(t *testing.T) {
	tfidf := NewTFIDF()

	tokens := []string{"apple", "banana", "apple", "apple"}
	tf := tfidf.ComputeTF(tokens)

	if tf["apple"] != 0.75 {
		t.Errorf("TF for apple = %f, want 0.75", tf["apple"])
	}
	if tf["banana"] != 0.25 {
		t.Errorf("TF for banana = %f, want 0.25", tf["banana"])
	}

	emptyTF := tfidf.ComputeTF([]string{})
	if len(emptyTF) != 0 {
		t.Errorf("Empty TF should be empty map, got %v", emptyTF)
	}
}

// @test TEST-INT_SIM-003
func TestTFIDF_ComputeIDF(t *testing.T) {
	tfidf := NewTFIDF()

	documents := [][]string{
		{"hello", "world"},
		{"hello", "foo"},
		{"hello", "bar"},
	}
	tfidf.ComputeIDF(documents)

	if tfidf.idf["world"] == 0 {
		t.Errorf("IDF for world should not be 0")
	}

	if tfidf.idf["hello"] >= 0 {
		t.Errorf("IDF for hello (appears in all 3 docs) should be <= 0, got %f", tfidf.idf["hello"])
	}
}

// @test TEST-INT_SIM-004
func TestTFIDF_ComputeTFIDF(t *testing.T) {
	tfidf := NewTFIDF()
	tfidf.idf = map[string]float64{
		"hello": 1.0,
		"world": 2.0,
	}

	tf := map[string]float64{
		"hello": 0.5,
		"world": 0.5,
	}

	tfidfResult := tfidf.ComputeTFIDF(tf)

	if tfidfResult["hello"] != 0.5 {
		t.Errorf("TFIDF for hello = %f, want 0.5", tfidfResult["hello"])
	}
	if tfidfResult["world"] != 1.0 {
		t.Errorf("TFIDF for world = %f, want 1.0", tfidfResult["world"])
	}
}

// @test TEST-INT_SIM-005
func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		vec1     map[string]float64
		vec2     map[string]float64
		expected float64
		approx   bool
	}{
		{
			"identical vectors",
			map[string]float64{"a": 1, "b": 2},
			map[string]float64{"a": 1, "b": 2},
			1.0,
			true,
		},
		{
			"orthogonal vectors",
			map[string]float64{"a": 1},
			map[string]float64{"b": 1},
			0.0,
			false,
		},
		{
			"opposite vectors",
			map[string]float64{"a": 1},
			map[string]float64{"a": -1},
			-1.0,
			true,
		},
		{
			"empty vectors",
			map[string]float64{},
			map[string]float64{},
			0.0,
			false,
		},
		{
			"partial overlap",
			map[string]float64{"a": 1, "b": 1},
			map[string]float64{"b": 1, "c": 1},
			0.5,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CosineSimilarity(tt.vec1, tt.vec2)
			if tt.approx {
				diff := result - tt.expected
				if diff < -0.001 || diff > 0.001 {
					t.Errorf("CosineSimilarity = %f, want approx %f", result, tt.expected)
				}
			} else {
				if result != tt.expected {
					t.Errorf("CosineSimilarity = %f, want %f", result, tt.expected)
				}
			}
		})
	}
}

// @test TEST-INT_SIM-005
func TestCosineSimilarity_ZeroNorm(t *testing.T) {
	result := CosineSimilarity(map[string]float64{}, map[string]float64{"a": 1})
	if result != 0 {
		t.Errorf("CosineSimilarity with zero norm should return 0, got %f", result)
	}
}

// @test TEST-INT_SIM-002
func TestScore(t *testing.T) {
	highSimilarity := "validates user credentials and issues JWT tokens"
	codeText := "validates user credentials and issues JWT tokens for session management"

	score := Score(highSimilarity, codeText)
	if score < 0.7 {
		t.Errorf("Similar texts should have high score, got %f", score)
	}

	lowSimilarity := "user authentication and login"
	codeText = "database connection pooling configuration"

	score = Score(lowSimilarity, codeText)
	if score > 0.3 {
		t.Errorf("Dissimilar texts should have low score, got %f", score)
	}
}

// @test TEST-INT_SIM-006
func TestScore_EmptyStrings(t *testing.T) {
	score := Score("", "")
	if score != 0 {
		t.Errorf("Empty strings should return 0, got %f", score)
	}
}

// @test TEST-INT_SIM-003
func TestNormalizeText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello world"},
		{"Test@#$%Case", "test case"},
		{"  multiple   spaces  ", "multiple spaces"},
		{"UPPERCASE", "uppercase"},
	}

	for _, tt := range tests {
		result := NormalizeText(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeText(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
