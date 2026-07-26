// Package collector provides self-describing IDD document file operations.

// Spec: docs/internal/collector/spec.md
// Contract: docs/internal/collector/contract.md
package collector

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// InitDocuments creates or adopts the four package-local self-describing
// Markdown documents. It never overwrites existing IDD or legacy metadata.
//
// @implement SPEC-INTERNAL_COLLECTOR-001
func InitDocuments(projectRoot, packagePath string) ([]string, error) {
	cleanPackage, err := validateDocumentPackagePath(packagePath)
	if err != nil {
		return nil, err
	}

	projectRoot = filepath.Clean(projectRoot)
	codePath := filepath.Join(projectRoot, cleanPackage)
	if info, statErr := os.Stat(codePath); statErr != nil || !info.IsDir() {
		return nil, fmt.Errorf("package directory %s does not exist", codePath)
	}

	docsDir := filepath.Join(projectRoot, "docs", cleanPackage)
	if _, statErr := os.Stat(filepath.Join(docsDir, "idd.yaml")); statErr == nil {
		return nil, fmt.Errorf(
			"central catalog found at %s; migrate its records into the four role-owned Markdown documents before initialization",
			filepath.Join(docsDir, "idd.yaml"),
		)
	} else if !os.IsNotExist(statErr) {
		return nil, fmt.Errorf("inspect package documentation: %w", statErr)
	}

	for _, filename := range iddDocumentOrder {
		path := filepath.Join(docsDir, filename)
		data, readErr := os.ReadFile(path)
		if os.IsNotExist(readErr) {
			continue
		}
		if readErr != nil {
			return nil, fmt.Errorf("inspect existing document %s: %w", path, readErr)
		}
		if hasIDDDocumentFrontmatter(data) {
			return nil, fmt.Errorf("self-describing IDD document already exists: %s", path)
		}
		if hasLegacyDocumentMetadata(data) {
			return nil, fmt.Errorf(
				"legacy metadata found in %s; migrate its markers and relationships before initializing self-describing documents",
				path,
			)
		}
	}

	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		return nil, fmt.Errorf("create documentation directory: %w", err)
	}

	packageSlash := filepath.ToSlash(cleanPackage)
	changed := make([]string, 0, len(iddDocumentOrder))
	for _, filename := range iddDocumentOrder {
		path := filepath.Join(docsDir, filename)
		role := iddDocumentRoles[filename]
		metadata := newIDDDocument(packageSlash, role)
		existing, readErr := os.ReadFile(path)
		switch {
		case readErr == nil:
			body := append([]byte("\n"), existing...)
			var root *yaml.Node
			frontmatter, existingBody, _, _, found, splitErr := splitLeadingFrontmatter(existing)
			if splitErr != nil {
				return changed, fmt.Errorf("parse existing frontmatter in %s: %w", path, splitErr)
			}
			if found {
				var yamlDocument yaml.Node
				if unmarshalErr := yaml.Unmarshal(frontmatter, &yamlDocument); unmarshalErr != nil {
					return changed, fmt.Errorf("parse existing frontmatter in %s: %w", path, unmarshalErr)
				}
				root = documentRoot(&yamlDocument)
				if root == nil || root.Kind == 0 {
					root = mappingNode()
				}
				body = existingBody
			}
			rendered, marshalErr := marshalIDDDocument(metadata, body, root)
			if marshalErr != nil {
				return changed, marshalErr
			}
			if writeErr := replaceFileAtomically(path, rendered); writeErr != nil {
				return changed, writeErr
			}
			changed = append(changed, path)
		case os.IsNotExist(readErr):
			rendered, marshalErr := MarshalIDDDocument(metadata, []byte(documentNarrativeTemplate(packageSlash, role)))
			if marshalErr != nil {
				return changed, marshalErr
			}
			if writeErr := writeExclusiveFile(path, rendered); writeErr != nil {
				return changed, writeErr
			}
			changed = append(changed, path)
		default:
			return changed, fmt.Errorf("read existing document %s: %w", path, readErr)
		}
	}
	sort.Strings(changed)
	return changed, nil
}

// RepairDocuments normalizes safe metadata in either one self-describing
// document or all four documents in a package directory. A single-file target
// never writes to sibling documents.
//
// @implement SPEC-INTERNAL_COLLECTOR-001
func RepairDocuments(path string) ([]string, error) {
	cleanPath := filepath.Clean(path)
	info, err := os.Stat(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("inspect document path: %w", err)
	}

	var docsDir string
	var targets []string
	if info.IsDir() {
		docsDir = cleanPath
		targets = append(targets, iddDocumentOrder...)
	} else {
		filename := filepath.Base(cleanPath)
		if _, ok := iddDocumentRoles[filename]; !ok {
			return nil, fmt.Errorf("document path must name design.md, contract.md, spec.md, testing.md, or their containing directory")
		}
		docsDir = filepath.Dir(cleanPath)
		targets = []string{filename}
	}
	centralCatalogPath := filepath.Join(docsDir, "idd.yaml")
	if _, statErr := os.Stat(centralCatalogPath); statErr == nil {
		return nil, fmt.Errorf(
			"central catalog found at %s; migrate its records into the four role-owned Markdown documents before repair",
			centralCatalogPath,
		)
	} else if !os.IsNotExist(statErr) {
		return nil, fmt.Errorf("inspect package documentation: %w", statErr)
	}

	packagePath := packageFromDocumentPath(filepath.Join(docsDir, "spec.md"))
	if packagePath == "" {
		return nil, fmt.Errorf("documents must be located at docs/<package>/")
	}

	type pendingWrite struct {
		path string
		data []byte
		new  bool
	}
	pending := make([]pendingWrite, 0, len(targets))
	for _, filename := range targets {
		role := iddDocumentRoles[filename]
		documentPath := filepath.Join(docsDir, filename)
		data, readErr := os.ReadFile(documentPath)
		if os.IsNotExist(readErr) {
			if !info.IsDir() {
				return nil, fmt.Errorf("document does not exist: %s", documentPath)
			}
			rendered, marshalErr := MarshalIDDDocument(
				newIDDDocument(packagePath, role),
				[]byte(documentNarrativeTemplate(packagePath, role)),
			)
			if marshalErr != nil {
				return nil, marshalErr
			}
			pending = append(pending, pendingWrite{path: documentPath, data: rendered, new: true})
			continue
		}
		if readErr != nil {
			return nil, fmt.Errorf("read document %s: %w", documentPath, readErr)
		}

		parsed, detected, parseErr := parseIDDDocument(documentPath, data)
		if parseErr != nil {
			return nil, fmt.Errorf("repair %s: %w", documentPath, parseErr)
		}
		if !detected {
			return nil, fmt.Errorf(
				"%s is not self-describing; migrate legacy metadata before repair",
				documentPath,
			)
		}

		normalizeIDDDocument(parsed.Metadata, packagePath, role)
		rendered, marshalErr := marshalIDDDocument(parsed.Metadata, parsed.Body, parsed.Root)
		if marshalErr != nil {
			return nil, marshalErr
		}
		if !bytes.Equal(data, rendered) {
			pending = append(pending, pendingWrite{path: documentPath, data: rendered})
		}
	}

	changed := make([]string, 0, len(pending))
	for _, write := range pending {
		var writeErr error
		if write.new {
			writeErr = writeExclusiveFile(write.path, write.data)
		} else {
			writeErr = replaceFileAtomically(write.path, write.data)
		}
		if writeErr != nil {
			return changed, writeErr
		}
		changed = append(changed, write.path)
	}
	sort.Strings(changed)
	return changed, nil
}

func newIDDDocument(packagePath, role string) *IDDDocument {
	return &IDDDocument{
		Version:  iddDocumentVersion,
		Package:  filepath.ToSlash(packagePath),
		Document: role,
	}
}

func normalizeIDDDocument(document *IDDDocument, packagePath, role string) {
	document.Version = iddDocumentVersion
	document.Package = filepath.ToSlash(packagePath)
	document.Document = role
}

func validateDocumentPackagePath(packagePath string) (string, error) {
	if packagePath == "" || filepath.IsAbs(packagePath) {
		return "", fmt.Errorf("package path must be a non-empty project-relative path")
	}
	clean := filepath.Clean(filepath.FromSlash(packagePath))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("package path must stay inside the project root")
	}
	return clean, nil
}

func hasLegacyDocumentMetadata(data []byte) bool {
	content := string(data)
	if hasLeadingFrontmatter(content) {
		frontmatter, err := ParseFrontmatter(content)
		if err != nil {
			return true
		}
		if frontmatter != nil && (len(frontmatter.Markers) > 0 || frontmatter.RelatedFiles != nil) {
			return true
		}
	}
	for _, field := range []string{
		"**Design:**",
		"**Contract:**",
		"**Requirement:**",
		"**Tests:**",
		"**Spec Coverage:**",
		"**Implements:**",
	} {
		if strings.Contains(content, field) {
			return true
		}
	}
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "#") {
			continue
		}
		if len(extractIDDRefs(trimmed)) > 0 {
			return true
		}
	}
	return false
}

func documentNarrativeTemplate(packagePath, role string) string {
	return renderIDDDocumentTemplate(packagePath, role)
}

func writeExclusiveFile(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return fmt.Errorf("write %s: %w", path, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close %s: %w", path, err)
	}
	return nil
}

func replaceFileAtomically(path string, data []byte) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("inspect %s: %w", path, err)
	}

	temp, err := os.CreateTemp(filepath.Dir(path), ".idd-document-*")
	if err != nil {
		return fmt.Errorf("create temporary document: %w", err)
	}
	tempPath := temp.Name()
	defer func() {
		_ = os.Remove(tempPath)
	}()

	if err := temp.Chmod(info.Mode().Perm()); err != nil {
		_ = temp.Close()
		return fmt.Errorf("set temporary document permissions: %w", err)
	}
	if _, err := temp.Write(data); err != nil {
		_ = temp.Close()
		return fmt.Errorf("write temporary document: %w", err)
	}
	if err := temp.Sync(); err != nil {
		_ = temp.Close()
		return fmt.Errorf("sync temporary document: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close temporary document: %w", err)
	}
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("replace document: %w", err)
	}
	return nil
}
