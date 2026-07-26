// Package walk provides testing utilities for the walk module.
package walk

import (
	"os"
	"path/filepath"
	"testing"
)

// @test TEST-PKG_WALK-007
func TestWalk(t *testing.T) {
	tmpDir := t.TempDir()

	os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "file1.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "subdir", "file3.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("README"), 0644)

	var visited []string
	err := Walk([]string{filepath.Join(tmpDir, "*.go")}, func(path string, info os.FileInfo) error {
		visited = append(visited, path)
		return nil
	})

	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	if len(visited) < 2 {
		t.Errorf("Expected at least 2 .go files, got %d: %v", len(visited), visited)
	}
}

// @test-contract TEST-PKG_WALK-001
func TestWalk_NoMatches(t *testing.T) {
	tmpDir := t.TempDir()

	err := Walk([]string{filepath.Join(tmpDir, "*.go")}, func(path string, info os.FileInfo) error {
		t.Error("Should not visit any files")
		return nil
	})

	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}
}

// @test-contract TEST-PKG_WALK-002
func TestWalk_StopOnError(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "file1.go"), []byte("package main"), 0644)

	expectedErr := os.ErrPermission
	err := Walk([]string{filepath.Join(tmpDir, "*.go")}, func(path string, info os.FileInfo) error {
		return expectedErr
	})

	if err != expectedErr {
		t.Errorf("Walk error = %v, want %v", err, expectedErr)
	}
}

// @test TEST-PKG_WALK-003
func TestWalk_Directory(t *testing.T) {
	tmpDir := t.TempDir()

	os.Mkdir(filepath.Join(tmpDir, "mydir"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "mydir", "file.go"), []byte("package main"), 0644)

	var visited []string
	err := Walk([]string{filepath.Join(tmpDir, "mydir")}, func(path string, info os.FileInfo) error {
		visited = append(visited, path)
		return nil
	})

	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	found := false
	for _, v := range visited {
		if filepath.Base(v) == "file.go" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Should visit file.go inside directory")
	}
}

// @test TEST-PKG_WALK-004
func TestWalk_DuplicatePrevention(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "file.go"), []byte("package main"), 0644)

	var visited int
	err := Walk([]string{filepath.Join(tmpDir, "*.go"), filepath.Join(tmpDir, "file.go")}, func(path string, info os.FileInfo) error {
		visited++
		return nil
	})

	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	if visited != 1 {
		t.Errorf("Should visit file.go only once, visited %d times", visited)
	}
}

// @test TEST-PKG_WALK-008
func TestMatchAnyExtensions(t *testing.T) {
	tests := []struct {
		path   string
		exts   []string
		result bool
	}{
		{"file.go", []string{".go", ".ts"}, true},
		{"file.ts", []string{".go", ".ts"}, true},
		{"file.txt", []string{".go", ".ts"}, false},
		{"file.go.txt", []string{".go"}, false},
		{"file", []string{".go"}, false},
		{"file.go", []string{}, false},
		{"file.go", []string{".cpp", ".h"}, false},
	}

	for _, tt := range tests {
		result := MatchAnyExtensions(tt.path, tt.exts)
		if result != tt.result {
			t.Errorf("MatchAnyExtensions(%q, %v) = %v, want %v", tt.path, tt.exts, result, tt.result)
		}
	}
}

// @test TEST-PKG_WALK-005
func TestWalk_NilPatterns(t *testing.T) {
	var visited int
	err := Walk(nil, func(path string, info os.FileInfo) error {
		visited++
		return nil
	})

	if err != nil {
		t.Fatalf("Walk with nil patterns failed: %v", err)
	}
	if visited != 0 {
		t.Errorf("Should not visit any files with nil patterns, visited %d", visited)
	}
}

// @test TEST-PKG_WALK-006
func TestWalk_EmptyPatterns(t *testing.T) {
	var visited int
	err := Walk([]string{}, func(path string, info os.FileInfo) error {
		visited++
		return nil
	})

	if err != nil {
		t.Fatalf("Walk with empty patterns failed: %v", err)
	}
	if visited != 0 {
		t.Errorf("Should not visit any files with empty patterns, visited %d", visited)
	}
}
