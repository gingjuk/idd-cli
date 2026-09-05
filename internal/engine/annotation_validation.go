// Package engine provides IDD validation engine.
package engine

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jingxu9x/idd-cli/internal/model"
)

// validateDuplicateIDs enforces a single canonical documentation declaration.
// Source annotations are evidence occurrences and may legitimately appear in
// several packages belonging to the same documented behavior.
func (e *Engine) validateDuplicateIDs(ids *model.IdentifierSet) {
	e.reportDuplicates(ids.DuplicateDocGroups(), "doc")
}

func (e *Engine) reportDuplicates(groups [][]*model.Identifier, side string) {
	for _, group := range groups {
		dirs := make([]string, len(group))
		for i, id := range group {
			dirs[i] = filepath.Dir(id.Source)
		}
		for i, id := range group {
			others := make([]string, 0, len(dirs)-1)
			for j, d := range dirs {
				if j != i {
					others = append(others, d)
				}
			}
			suggestion := suggestRename(id.ID, id.Source)
			msg := fmt.Sprintf("%s (%s) also defined in: %s", id.ID, side, strings.Join(others, ", "))
			if suggestion != "" {
				msg += fmt.Sprintf("; consider renaming to %s", suggestion)
			}
			e.result.AddError("duplicate-id", msg, id.Source, id.ID, "")
		}
	}
}

// suggestRename derives a rename suggestion by replacing the MODULE segment of id
// with a name derived from the source file's directory path.
func suggestRename(id, sourcePath string) string {
	parts := strings.SplitN(id, "-", 3)
	if len(parts) != 3 {
		return ""
	}
	newModule := moduleFromPath(sourcePath)
	if newModule == "" || newModule == parts[1] {
		return ""
	}
	return parts[0] + "-" + newModule + "-" + parts[2]
}

// moduleFromPath derives an uppercase module name from a file path.
// For doc paths it uses directory components after "docs/".
// For code paths it uses the last 1-2 directory components.
func moduleFromPath(path string) string {
	dir := filepath.ToSlash(filepath.Dir(path))
	segs := strings.Split(dir, "/")
	for i, seg := range segs {
		if seg == "docs" && i+1 < len(segs) {
			first := pathSegToModule(segs[i+1])
			if i+2 < len(segs) && segs[i+2] != "" {
				return first + "_" + pathSegToModule(segs[i+2])
			}
			return first
		}
	}
	var tail []string
	for j := len(segs) - 1; j >= 0 && len(tail) < 2; j-- {
		if segs[j] != "" && segs[j] != "." {
			tail = append([]string{segs[j]}, tail...)
		}
	}
	result := make([]string, len(tail))
	for i, s := range tail {
		result[i] = pathSegToModule(s)
	}
	return strings.Join(result, "_")
}

func pathSegToModule(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, "-", "_"))
}
