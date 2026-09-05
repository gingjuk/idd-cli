// Package engine provides configured source and document file traversal.
package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jingxu9x/idd-cli/internal/model"
)

// walkCodeFiles iterates over all code files matching the configured patterns
// and invokes the visitor function for each file with its path and line content.
func (e *Engine) walkCodeFiles(visitor func(path string, lines []string)) {
	for _, globPattern := range e.cfg.Code.Patterns {
		e.walkGlob(globPattern, visitor)
	}
}

// walkDocFiles iterates over all documentation files matching the configured
// patterns and invokes the visitor function for each file with its path and lines.
func (e *Engine) walkDocFiles(visitor func(path string, lines []string)) {
	for _, globPattern := range e.cfg.Docs.Patterns {
		e.walkGlob(globPattern, visitor)
	}
}

// walkGlob walks the directory tree matching the given pattern and invokes
// the visitor function for each matching file. Handles non-recursive patterns.
func (e *Engine) walkGlob(pattern string, visitor func(path string, lines []string)) {
	hasRecursive := strings.Contains(pattern, "**")
	if hasRecursive {
		e.walkGlobRecursive(pattern, visitor)
		return
	}
	dir, file := filepath.Split(pattern)
	if dir == "" {
		dir = "."
	}
	if file == "" {
		return
	}
	if _, err := filepath.Match(file, ""); err != nil {
		e.addScanFinding(pattern, "invalid-pattern", err)
		return
	}
	resolvedDir := e.cfg.ResolvePath(dir)
	if err := filepath.WalkDir(resolvedDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			if !os.IsNotExist(err) {
				e.addScanFinding(e.cfg.DisplayPath(path), scanErrorCode(err, "walk"), err)
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}
		matched, err := filepath.Match(file, d.Name())
		if err != nil {
			e.addScanFinding(pattern, "invalid-pattern", err)
			return nil
		}
		if !matched {
			return nil
		}
		displayPath := e.cfg.DisplayPath(path)
		if e.shouldIgnorePath(displayPath) {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			e.addScanFinding(displayPath, scanErrorCode(err, "read"), err)
			return nil
		}
		lines := strings.Split(string(content), "\n")
		visitor(displayPath, lines)
		return nil
	}); err != nil && !os.IsNotExist(err) {
		e.addScanFinding(e.cfg.DisplayPath(resolvedDir), scanErrorCode(err, "walk"), err)
	}
}

// walkGlobRecursive walks the directory tree matching the given recursive
// pattern (containing **/) and invokes the visitor function for each matching file.
func (e *Engine) walkGlobRecursive(pattern string, visitor func(path string, lines []string)) {
	dir, file := filepath.Split(pattern)
	if dir == "" {
		dir = "."
	}
	if file == "" {
		return
	}
	dir = strings.TrimSuffix(dir, "**/")
	if dir == "" {
		dir = "."
	}
	file = strings.TrimPrefix(file, "**/")
	if _, err := filepath.Match(file, ""); err != nil {
		e.addScanFinding(pattern, "invalid-pattern", err)
		return
	}

	resolvedDir := e.cfg.ResolvePath(dir)
	if err := filepath.WalkDir(resolvedDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			if !os.IsNotExist(err) {
				e.addScanFinding(e.cfg.DisplayPath(path), scanErrorCode(err, "walk"), err)
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}
		matched, err := filepath.Match(file, d.Name())
		if err != nil {
			e.addScanFinding(pattern, "invalid-pattern", err)
			return nil
		}
		if !matched {
			return nil
		}
		displayPath := e.cfg.DisplayPath(path)
		if e.shouldIgnorePath(displayPath) {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			e.addScanFinding(displayPath, scanErrorCode(err, "read"), err)
			return nil
		}
		lines := strings.Split(string(content), "\n")
		visitor(displayPath, lines)
		return nil
	}); err != nil && !os.IsNotExist(err) {
		e.addScanFinding(e.cfg.DisplayPath(resolvedDir), scanErrorCode(err, "walk"), err)
	}
}

func (e *Engine) addScanFinding(path, code string, err error) {
	key := code + "\x00" + path + "\x00" + err.Error()
	if e.scanFindings[key] {
		return
	}
	e.scanFindings[key] = true
	e.result.AddError("filesystem-scan", fmt.Sprintf("%s %s: %v", code, path, err), path, "", code)
}

func scanErrorCode(err error, fallback string) string {
	if os.IsNotExist(err) {
		return "not-found"
	}
	if os.IsPermission(err) {
		return "permission-denied"
	}
	return fallback
}

// shouldIgnorePath checks whether a file path matches any of the configured
// ignore patterns for code files.
func (e *Engine) shouldIgnorePath(path string) bool {
	for _, ignore := range e.cfg.Code.IgnorePaths {
		if matched, _ := filepath.Match(ignore, filepath.Base(path)); matched {
			return true
		}
		if strings.Contains(path, ignore) {
			return true
		}
	}
	return false
}

// BuildReport generates a complete validation report with tool information,
// configuration summary, and the accumulated validation result.
// @implement SPEC-CMD_IDD_CLI-004
func (e *Engine) BuildReport() *model.Report {
	return &model.Report{
		Tool:      "idd-cli",
		Version:   "1.0.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Config: model.ConfigSummary{
			DocPatterns:  e.cfg.Docs.Patterns,
			CodePatterns: e.cfg.Code.Patterns,
			Annotations:  e.cfg.Code.Annotations,
		},
		Result: *e.result,
	}
}
