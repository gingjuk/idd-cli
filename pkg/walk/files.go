package walk

import (
	"os"
	"path/filepath"
)

type FileVisitor func(path string, info os.FileInfo) error

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

func MatchAnyExtensions(path string, extensions []string) bool {
	ext := filepath.Ext(path)
	for _, e := range extensions {
		if ext == e {
			return true
		}
	}
	return false
}