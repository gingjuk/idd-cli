// Package walk provides file traversal utilities.

// Spec: docs/pkg/walk/spec.md
// Contract: docs/pkg/walk/contract.md
package walk

import (
	"os"
	"path/filepath"
)

// FileVisitor is a callback function type for file visitation.
// @implement SPEC-PKG_WALK-001
type FileVisitor func(path string, info os.FileInfo) error

// Walk walks the filesystem matching files against the given glob patterns.
// @implement SPEC-PKG_WALK-002
func Walk(patterns []string, visitor FileVisitor) error {
	visited := make(map[string]bool)

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}

		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				continue
			}

			if visited[match] {
				continue
			}
			visited[match] = true

			if err := visitor(match, info); err != nil {
				return err
			}

			if info.IsDir() {
				if err := filepath.Walk(match, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return nil
					}
					if visited[path] {
						return nil
					}
					visited[path] = true
					return visitor(path, info)
				}); err != nil {
					continue
				}
			}
		}
	}

	return nil
}

// MatchAnyExtensions checks if a file path has any of the specified extensions.
// @implement SPEC-PKG_WALK-003
func MatchAnyExtensions(path string, extensions []string) bool {
	ext := filepath.Ext(path)
	for _, e := range extensions {
		if ext == e {
			return true
		}
	}
	return false
}
