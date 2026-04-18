package main

import (
	"os"
	"path/filepath"
	"strings"
)

func ListEmbeddedSkills() []string {
	var skills []string
	entries, err := os.ReadDir("skills")
	if err != nil {
		return skills
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			skills = append(skills, filepath.Join("skills", entry.Name()))
		}
	}
	return skills
}

func ReadEmbeddedSkill(path string) ([]byte, error) {
	return os.ReadFile(path)
}
