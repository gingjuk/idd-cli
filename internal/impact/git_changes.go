// Package impact turns version-control changes into deterministic seeds for
// IDD trace queries. It deliberately knows nothing about IDD relationships;
// those are resolved by graph.TraceIndex.
package impact

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const ChangeSetSchema = "idd.git_changes.v1"

// Hunk is a changed line range in the base and current worktree versions.
type Hunk struct {
	OldStart int `json:"old_start"`
	OldLines int `json:"old_lines"`
	NewStart int `json:"new_start"`
	NewLines int `json:"new_lines"`
}

// FileChange describes one changed repository-relative path.
type FileChange struct {
	Path          string   `json:"path"`
	OldPath       string   `json:"old_path,omitempty"`
	Status        string   `json:"status"`
	Hunks         []Hunk   `json:"hunks"`
	AnnotationIDs []string `json:"annotation_ids,omitempty"`
}

// ChangeSet is the stable Git-derived input to an impacted query.
type ChangeSet struct {
	Schema string       `json:"schema"`
	Base   string       `json:"base"`
	Files  []FileChange `json:"files"`
}

type commandRunner interface {
	Run(ctx context.Context, dir string, args ...string) ([]byte, error)
}

type gitRunner struct{}

func (gitRunner) Run(ctx context.Context, dir string, args ...string) ([]byte, error) {
	command := exec.CommandContext(ctx, "git", args...)
	command.Dir = dir
	output, err := command.Output()
	if err == nil {
		return output, nil
	}
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		message := strings.TrimSpace(string(exitError.Stderr))
		if message != "" {
			return nil, fmt.Errorf("git %s: %s", strings.Join(args, " "), message)
		}
	}
	return nil, fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
}

// CollectGitChanges compares base with the complete current worktree. Tracked
// staged and unstaged changes are included by Git's <base> comparison, while
// untracked files are added explicitly.
// @implement SPEC-CMD_IDD_CLI-012
func CollectGitChanges(ctx context.Context, projectRoot, base string) (*ChangeSet, error) {
	return collectGitChanges(ctx, projectRoot, base, gitRunner{})
}

func collectGitChanges(
	ctx context.Context,
	projectRoot string,
	base string,
	runner commandRunner,
) (*ChangeSet, error) {
	if strings.TrimSpace(base) == "" {
		return nil, fmt.Errorf("git base revision is empty")
	}
	root, err := filepath.Abs(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("resolve project root: %w", err)
	}
	repositoryOutput, err := runner.Run(ctx, root, "rev-parse", "--show-toplevel")
	if err != nil {
		return nil, fmt.Errorf("resolve Git repository: %w", err)
	}
	repositoryRoot := filepath.Clean(strings.TrimSpace(string(repositoryOutput)))
	if repositoryRoot == "." || repositoryRoot == "" {
		return nil, fmt.Errorf("git returned an empty repository root")
	}
	if _, err := runner.Run(ctx, repositoryRoot, "rev-parse", "--verify", base+"^{commit}"); err != nil {
		return nil, fmt.Errorf("resolve base revision %q: %w", base, err)
	}

	nameStatus, err := runner.Run(
		ctx,
		root,
		"diff", "--relative", "--name-status", "-z", "--find-renames", base, "--",
	)
	if err != nil {
		return nil, fmt.Errorf("list tracked changes: %w", err)
	}
	files, err := parseNameStatus(nameStatus)
	if err != nil {
		return nil, err
	}
	for index := range files {
		paths := []string{files[index].Path}
		if files[index].OldPath != "" && files[index].OldPath != files[index].Path {
			paths = append([]string{files[index].OldPath}, paths...)
		}
		arguments := []string{
			"diff", "--relative", "--no-ext-diff", "--no-color", "--unified=0", "--find-renames", base, "--",
		}
		arguments = append(arguments, paths...)
		patch, patchErr := runner.Run(ctx, root, arguments...)
		if patchErr != nil {
			return nil, fmt.Errorf("read patch for %q: %w", files[index].Path, patchErr)
		}
		files[index].Hunks, files[index].AnnotationIDs = parsePatch(patch)
	}

	untrackedOutput, err := runner.Run(
		ctx,
		root,
		"ls-files", "--others", "--exclude-standard", "-z", "--",
	)
	if err != nil {
		return nil, fmt.Errorf("list untracked files: %w", err)
	}
	for _, path := range splitNUL(untrackedOutput) {
		content, readErr := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if readErr != nil {
			return nil, fmt.Errorf("read untracked file %q: %w", path, readErr)
		}
		lineCount := countLines(content)
		change := FileChange{
			Path:          filepath.ToSlash(path),
			Status:        "untracked",
			Hunks:         []Hunk{{OldStart: 0, OldLines: 0, NewStart: 1, NewLines: lineCount}},
			AnnotationIDs: annotationIDs(content),
		}
		files = append(files, change)
	}

	files = normalizeChanges(files)
	return &ChangeSet{Schema: ChangeSetSchema, Base: base, Files: files}, nil
}

func parseNameStatus(data []byte) ([]FileChange, error) {
	fields := splitNUL(data)
	changes := make([]FileChange, 0, len(fields)/2)
	for cursor := 0; cursor < len(fields); {
		status := fields[cursor]
		cursor++
		if status == "" {
			continue
		}
		if cursor >= len(fields) {
			return nil, fmt.Errorf("malformed Git name-status output after %q", status)
		}
		code := status[:1]
		change := FileChange{Status: normalizeStatus(code), Hunks: []Hunk{}}
		if code == "R" || code == "C" {
			if cursor+1 >= len(fields) {
				return nil, fmt.Errorf("malformed Git rename/copy entry %q", status)
			}
			change.OldPath = filepath.ToSlash(fields[cursor])
			change.Path = filepath.ToSlash(fields[cursor+1])
			cursor += 2
		} else {
			change.Path = filepath.ToSlash(fields[cursor])
			cursor++
		}
		changes = append(changes, change)
	}
	return changes, nil
}

var hunkHeaderPattern = regexp.MustCompile(
	`^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@`,
)

var annotationPattern = regexp.MustCompile(
	`(?i)@(implement|test(?:-contract)?)\s+([^\r\n]+)`,
)

var annotationIDPattern = regexp.MustCompile(
	`\b(?:SPEC|TEST)-[A-Z][A-Z0-9_]*-[0-9]+\b`,
)

func parsePatch(patch []byte) ([]Hunk, []string) {
	hunks := make([]Hunk, 0)
	annotationText := make([]byte, 0)
	scanner := bufio.NewScanner(bytes.NewReader(patch))
	insideHunk := false
	for scanner.Scan() {
		line := scanner.Text()
		match := hunkHeaderPattern.FindStringSubmatch(line)
		if len(match) > 0 {
			hunks = append(hunks, Hunk{
				OldStart: parseHunkNumber(match[1]),
				OldLines: parseHunkCount(match[2]),
				NewStart: parseHunkNumber(match[3]),
				NewLines: parseHunkCount(match[4]),
			})
			insideHunk = true
			continue
		}
		if strings.HasPrefix(line, "diff --git ") {
			insideHunk = false
			continue
		}
		if insideHunk && len(line) > 1 && (line[0] == '+' || line[0] == '-') &&
			!strings.HasPrefix(line, "+++") && !strings.HasPrefix(line, "---") {
			annotationText = append(annotationText, line[1:]...)
			annotationText = append(annotationText, '\n')
		}
	}
	return hunks, annotationIDs(annotationText)
}

func annotationIDs(content []byte) []string {
	seen := make(map[string]bool)
	for _, match := range annotationPattern.FindAllSubmatch(content, -1) {
		for _, id := range annotationIDPattern.FindAllString(strings.ToUpper(string(match[2])), -1) {
			seen[id] = true
		}
	}
	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func normalizeChanges(changes []FileChange) []FileChange {
	byPath := make(map[string]FileChange, len(changes))
	for _, change := range changes {
		change.Path = filepath.ToSlash(filepath.Clean(change.Path))
		if change.OldPath != "" {
			change.OldPath = filepath.ToSlash(filepath.Clean(change.OldPath))
		}
		sort.Slice(change.Hunks, func(i, j int) bool {
			if change.Hunks[i].NewStart != change.Hunks[j].NewStart {
				return change.Hunks[i].NewStart < change.Hunks[j].NewStart
			}
			return change.Hunks[i].OldStart < change.Hunks[j].OldStart
		})
		sort.Strings(change.AnnotationIDs)
		byPath[change.Path] = change
	}
	result := make([]FileChange, 0, len(byPath))
	for _, change := range byPath {
		result = append(result, change)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Path < result[j].Path })
	return result
}

func splitNUL(data []byte) []string {
	raw := bytes.Split(data, []byte{0})
	result := make([]string, 0, len(raw))
	for _, field := range raw {
		if len(field) > 0 {
			result = append(result, string(field))
		}
	}
	return result
}

func normalizeStatus(status string) string {
	switch status {
	case "A":
		return "added"
	case "D":
		return "deleted"
	case "R":
		return "renamed"
	case "C":
		return "copied"
	case "T":
		return "type-changed"
	case "U":
		return "unmerged"
	default:
		return "modified"
	}
}

func parseHunkNumber(value string) int {
	number, _ := strconv.Atoi(value)
	return number
}

func parseHunkCount(value string) int {
	if value == "" {
		return 1
	}
	return parseHunkNumber(value)
}

func countLines(content []byte) int {
	if len(content) == 0 {
		return 0
	}
	count := bytes.Count(content, []byte{'\n'})
	if content[len(content)-1] != '\n' {
		count++
	}
	return count
}
