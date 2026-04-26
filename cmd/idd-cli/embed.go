package main

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed skills/*.md
var SkillsFS embed.FS

// @implement SPEC-CMD_IDD_CLI-008
func ListEmbeddedSkills() []string {
	var skills []string
	_ = fs.WalkDir(SkillsFS, "skills", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		skills = append(skills, path)
		return nil
	})
	return skills
}

// @implement SPEC-CMD_IDD_CLI-008
func ReadEmbeddedSkill(path string) ([]byte, error) {
	return SkillsFS.ReadFile(path)
}
