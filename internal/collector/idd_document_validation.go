// Package collector validates self-describing IDD document records and graph relationships.
package collector

import (
	stderrors "errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jingxu9x/idd-cli/internal/model"
	"github.com/jingxu9x/idd-cli/pkg/pattern"
	"github.com/yuin/goldmark"
	goldmarkast "github.com/yuin/goldmark/ast"
	goldmarktext "github.com/yuin/goldmark/text"
)

func (c *DocCollector) collectIDDDocumentSet(directory string, set *model.IdentifierSet) []*model.ValidationError {
	documents := make(map[string]*parsedIDDDocument, len(iddDocumentOrder))
	var errors []*model.ValidationError

	for _, filename := range iddDocumentOrder {
		role := iddDocumentRoles[filename]
		path := filepath.Join(directory, filename)
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-set",
					fmt.Sprintf("self-describing document set is missing %s", filename),
					path,
					1,
					"",
					role,
				))
				continue
			}
			errors = append(errors, iddDocumentValidationError(
				"idd-document-parse",
				fmt.Sprintf("read %s: %v", filename, err),
				path,
				1,
				"",
				role,
			))
			continue
		}

		parsed, detected, parseErr := parseIDDDocument(path, data)
		if parseErr != nil {
			var semanticError *iddSemanticFrontmatterError
			if stderrors.As(parseErr, &semanticError) {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-migration",
					semanticError.Error(),
					path,
					semanticError.Line,
					"",
					"yaml-semantics",
				))
				continue
			}
			errors = append(errors, iddDocumentValidationError(
				"idd-document-parse",
				parseErr.Error(),
				path,
				yamlErrorLine(parseErr),
				"",
				role,
			))
			continue
		}
		if !detected {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-set",
				fmt.Sprintf("%s must contain an idd frontmatter block because another package document is self-describing", filename),
				path,
				1,
				"",
				role,
			))
			continue
		}
		documents[role] = parsed
		errors = append(errors, validateIDDDocument(parsed, role)...)
	}

	errors = append(errors, validateIDDDocumentReferences(documents)...)
	addIDDDocumentIdentifiers(documents, set)
	declared := declaredIDDIdentifiers(documents)
	for _, document := range documents {
		errors = append(errors, validateIDDDocumentMarkdown(document, declared)...)
	}
	return errors
}

func validateIDDDocument(document *parsedIDDDocument, expectedRole string) []*model.ValidationError {
	metadata := document.Metadata
	index := document.Index
	var errors []*model.ValidationError
	addError := func(rule, message, link, code string, line int) {
		errors = append(errors, iddDocumentValidationError(rule, message, document.Path, line, link, code))
	}

	if metadata.Version != iddDocumentVersion {
		addError("idd-document-identity", `version must be "1.0"`, "", "version", index.topFieldLine("version"))
	}
	expectedPackage := packageFromDocumentPath(document.Path)
	if expectedPackage == "" {
		addError(
			"idd-document-identity",
			"self-describing documents must be located at docs/<package>/<role>.md",
			"",
			"package",
			index.topFieldLine("package"),
		)
	} else if metadata.Package == "" {
		addError("idd-document-identity", "package is required", "", "package", index.topFieldLine("package"))
	} else if metadata.Package != expectedPackage {
		addError(
			"idd-document-identity",
			fmt.Sprintf("package %q does not match document path; expected %q", metadata.Package, expectedPackage),
			"",
			"package",
			index.topFieldLine("package"),
		)
	}
	if metadata.Document != expectedRole {
		addError(
			"idd-document-identity",
			fmt.Sprintf("document must be %q in %s", expectedRole, filepath.Base(document.Path)),
			"",
			"document",
			index.topFieldLine("document"),
		)
	}

	switch expectedRole {
	case "design":
		errors = append(errors, validateIDDComponents(document)...)
	case "contract":
		errors = append(errors, validateIDDContracts(document)...)
	case "spec":
		errors = append(errors, validateIDDSpecs(document, expectedPackage)...)
	case "testing":
		errors = append(errors, validateIDDTests(document, expectedPackage)...)
	}
	for _, slot := range inspectParsedDocumentCompletion(document, expectedRole, false) {
		errors = append(errors, iddDocumentValidationError(
			"idd-document-incomplete",
			fmt.Sprintf("%s: %s", slot.Slot, slot.Reason),
			slot.File,
			slot.Line,
			"",
			slot.Slot,
		))
	}

	if mappingValue(document.Root, "markers") != nil || mappingValue(document.Root, "related_files") != nil {
		addError(
			"idd-document-markdown",
			"self-describing documents must not repeat legacy markers or related_files metadata",
			"",
			"duplicate-metadata",
			1,
		)
	}
	for _, issue := range document.Issues {
		errors = append(errors, iddDocumentValidationError(
			issue.Rule,
			issue.Message,
			document.Path,
			issue.Line,
			issue.Link,
			issue.Code,
		))
	}
	return errors
}

func validateIDDComponents(document *parsedIDDDocument) []*model.ValidationError {
	var errors []*model.ValidationError
	errors = append(errors, validateIDDNames(document, "components", document.Metadata.Components)...)
	declared := stringSet(document.Metadata.Components)
	for recordIndex, component := range document.Metadata.ComponentRecords {
		validateIDDRequiredRecordFields(&errors, document, "design", recordIndex)
		validateIDDRecordLifecycle(
			&errors,
			document,
			"components",
			recordIndex,
			"Component",
			component.Name,
			component.IDDRecordLifecycle,
			declared,
		)

		seen := make(map[string]bool)
		for _, dependency := range component.DependsOn {
			line := document.Index.recordFieldLine("components", recordIndex, "depends-on")
			if seen[dependency] {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-schema",
					fmt.Sprintf("Component %s repeats dependency %q", component.Name, dependency),
					document.Path,
					line,
					"",
					"depends-on",
				))
				continue
			}
			seen[dependency] = true
			if err := validateIDDScopedName(dependency); err != nil {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-schema",
					fmt.Sprintf(
						"Component %s has invalid dependency %q: %v",
						component.Name,
						dependency,
						err,
					),
					document.Path,
					line,
					"",
					"depends-on",
				))
				continue
			}
			targetPackage, targetName := splitIDDScopedName(document.Metadata.Package, dependency)
			if targetPackage == document.Metadata.Package {
				if targetName == component.Name {
					errors = append(errors, iddDocumentValidationError(
						"idd-document-schema",
						fmt.Sprintf("Component %s cannot depend on itself", component.Name),
						document.Path,
						line,
						"",
						"depends-on",
					))
				} else if !declared[targetName] {
					errors = append(errors, iddDocumentValidationError(
						"idd-document-reference",
						fmt.Sprintf("Component %s depends on undeclared Component %q", component.Name, targetName),
						document.Path,
						line,
						"",
						"depends-on",
					))
				}
			}
		}
	}
	errors = append(errors, validateIDDLifecycleRelationships(
		document,
		"components",
		"Component",
		componentLifecycleRecords(document.Metadata.ComponentRecords),
	)...)
	return errors
}

func validateIDDContracts(document *parsedIDDDocument) []*model.ValidationError {
	var errors []*model.ValidationError
	errors = append(errors, validateIDDNames(document, "contracts", document.Metadata.Contracts)...)
	declared := stringSet(document.Metadata.Contracts)
	for recordIndex, contract := range document.Metadata.ContractRecords {
		validateIDDRequiredRecordFields(&errors, document, "contract", recordIndex)
		validateIDDRecordLifecycle(
			&errors,
			document,
			"contracts",
			recordIndex,
			"Contract",
			contract.Name,
			contract.IDDRecordLifecycle,
			declared,
		)
	}
	errors = append(errors, validateIDDLifecycleRelationships(
		document,
		"contracts",
		"Contract",
		contractLifecycleRecords(document.Metadata.ContractRecords),
	)...)
	return errors
}

func validateIDDNames(document *parsedIDDDocument, field string, values []string) []*model.ValidationError {
	var errors []*model.ValidationError
	seen := make(map[string]bool, len(values))
	for valueIndex, value := range values {
		line := document.Index.recordLine(field, valueIndex)
		if strings.TrimSpace(value) == "" || isIDDPlaceholder(value, "", field) {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf("%s entries must contain concrete, meaningful names", field),
				document.Path,
				line,
				"",
				field,
			))
			continue
		}
		if strings.ContainsAny(value, "#,") {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf(
					"%s entry %q contains a reserved reference separator (# or ,)",
					field,
					value,
				),
				document.Path,
				line,
				"",
				field,
			))
			continue
		}
		if seen[value] {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf("%s repeats %q", field, value),
				document.Path,
				line,
				"",
				field,
			))
			continue
		}
		seen[value] = true
		if !strings.Contains(string(document.Body), value) {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-reference",
				fmt.Sprintf("%s entry %q has no narrative definition in %s", field, value, filepath.Base(document.Path)),
				document.Path,
				line,
				"",
				field,
			))
		}
	}
	return errors
}

var (
	iddLifecycleStatuses = map[string]bool{
		"active":     true,
		"deprecated": true,
		"superseded": true,
	}
	iddConcernNames = map[string]bool{
		"security":      true,
		"concurrency":   true,
		"persistence":   true,
		"performance":   true,
		"compatibility": true,
	}
)

func validateIDDRecordLifecycle(
	errors *[]*model.ValidationError,
	document *parsedIDDDocument,
	section string,
	recordIndex int,
	recordType string,
	recordID string,
	lifecycle IDDRecordLifecycle,
	declared map[string]bool,
) {
	if lifecycle.Status != "" && !iddLifecycleStatuses[lifecycle.Status] {
		*errors = append(*errors, iddDocumentValidationError(
			"idd-document-schema",
			fmt.Sprintf("%s %s status must be active, deprecated, or superseded", recordType, recordID),
			document.Path,
			document.Index.recordFieldLine(section, recordIndex, "status"),
			recordID,
			"status",
		))
	}
	seenSupersedes := make(map[string]bool)
	for _, target := range lifecycle.Supersedes {
		line := document.Index.recordFieldLine(section, recordIndex, "supersedes")
		switch {
		case target == recordID:
			*errors = append(*errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf("%s %s cannot supersede itself", recordType, recordID),
				document.Path,
				line,
				recordID,
				"supersedes",
			))
		case seenSupersedes[target]:
			*errors = append(*errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf("%s %s repeats superseded record %q", recordType, recordID, target),
				document.Path,
				line,
				recordID,
				"supersedes",
			))
		case !declared[target]:
			*errors = append(*errors, iddDocumentValidationError(
				"idd-document-reference",
				fmt.Sprintf("%s %s supersedes undeclared same-kind record %q", recordType, recordID, target),
				document.Path,
				line,
				recordID,
				"supersedes",
			))
		}
		seenSupersedes[target] = true
	}
	if lifecycle.DeprecatedBy != "" {
		line := document.Index.recordFieldLine(section, recordIndex, "deprecated-by")
		if lifecycle.DeprecatedBy == recordID {
			*errors = append(*errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf("%s %s cannot be deprecated by itself", recordType, recordID),
				document.Path,
				line,
				recordID,
				"deprecated-by",
			))
		} else if !declared[lifecycle.DeprecatedBy] {
			*errors = append(*errors, iddDocumentValidationError(
				"idd-document-reference",
				fmt.Sprintf("%s %s is deprecated by undeclared same-kind record %q", recordType, recordID, lifecycle.DeprecatedBy),
				document.Path,
				line,
				recordID,
				"deprecated-by",
			))
		}
	}

	seenConcerns := make(map[string]bool)
	for _, concern := range lifecycle.Concerns {
		line := document.Index.recordFieldLine(section, recordIndex, "concerns")
		switch {
		case !iddConcernNames[concern]:
			*errors = append(*errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf("%s %s uses unsupported concern %q", recordType, recordID, concern),
				document.Path,
				line,
				recordID,
				"concerns",
			))
		case seenConcerns[concern]:
			*errors = append(*errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf("%s %s repeats concern %q", recordType, recordID, concern),
				document.Path,
				line,
				recordID,
				"concerns",
			))
		case !iddRecordHasConcernSection(document, section, recordIndex, concern):
			*errors = append(*errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf(
					"%s %s selects concern %q but has no non-empty ### %s section",
					recordType,
					recordID,
					concern,
					iddConcernHeading(concern),
				),
				document.Path,
				line,
				recordID,
				"concerns",
			))
		}
		seenConcerns[concern] = true
	}
}

func iddConcernHeading(concern string) string {
	if concern == "" {
		return concern
	}
	return strings.ToUpper(concern[:1]) + concern[1:]
}

type iddLifecycleRecord struct {
	ID        string
	Index     int
	Lifecycle IDDRecordLifecycle
}

type iddReplacementLocation struct {
	recordID string
	field    string
	line     int
}

func componentLifecycleRecords(records []IDDDocumentComponent) []iddLifecycleRecord {
	result := make([]iddLifecycleRecord, 0, len(records))
	for index, record := range records {
		result = append(result, iddLifecycleRecord{record.Name, index, record.IDDRecordLifecycle})
	}
	return result
}

func contractLifecycleRecords(records []IDDDocumentContract) []iddLifecycleRecord {
	result := make([]iddLifecycleRecord, 0, len(records))
	for index, record := range records {
		result = append(result, iddLifecycleRecord{record.Name, index, record.IDDRecordLifecycle})
	}
	return result
}

func specLifecycleRecords(records []IDDDocumentSpec) []iddLifecycleRecord {
	result := make([]iddLifecycleRecord, 0, len(records))
	for index, record := range records {
		result = append(result, iddLifecycleRecord{record.ID, index, record.IDDRecordLifecycle})
	}
	return result
}

func testLifecycleRecords(records []IDDDocumentTest) []iddLifecycleRecord {
	result := make([]iddLifecycleRecord, 0, len(records))
	for index, record := range records {
		result = append(result, iddLifecycleRecord{record.ID, index, record.IDDRecordLifecycle})
	}
	return result
}

func validateIDDLifecycleRelationships(
	document *parsedIDDDocument,
	section string,
	recordType string,
	records []iddLifecycleRecord,
) []*model.ValidationError {
	var errors []*model.ValidationError
	byID := make(map[string]iddLifecycleRecord, len(records))
	for _, record := range records {
		if record.ID != "" {
			byID[record.ID] = record
		}
	}

	successors := make(map[string]map[string]iddReplacementLocation)
	addSuccessor := func(oldID, newID string, location iddReplacementLocation) {
		if oldID == "" || newID == "" || oldID == newID ||
			byID[oldID].ID == "" || byID[newID].ID == "" {
			return
		}
		if successors[oldID] == nil {
			successors[oldID] = make(map[string]iddReplacementLocation)
		}
		if _, exists := successors[oldID][newID]; !exists {
			successors[oldID][newID] = location
		}
	}
	for _, record := range records {
		for _, oldID := range record.Lifecycle.Supersedes {
			addSuccessor(oldID, record.ID, iddReplacementLocation{
				recordID: record.ID,
				field:    "supersedes",
				line:     document.Index.recordFieldLine(section, record.Index, "supersedes"),
			})
		}
		if record.Lifecycle.DeprecatedBy != "" {
			addSuccessor(record.ID, record.Lifecycle.DeprecatedBy, iddReplacementLocation{
				recordID: record.ID,
				field:    "deprecated-by",
				line:     document.Index.recordFieldLine(section, record.Index, "deprecated-by"),
			})
		}
	}

	for _, record := range records {
		targets := sortedReplacementTargets(successors[record.ID])
		if record.Lifecycle.Status == "superseded" && len(targets) == 0 {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf(
					"%s %s has superseded status but no replacement; add Deprecated by here or Supersedes on its replacement",
					recordType,
					record.ID,
				),
				document.Path,
				document.Index.recordFieldLine(section, record.Index, "status"),
				record.ID,
				"replacement",
			))
		}
		if len(targets) > 1 {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf(
					"%s %s has multiple replacements: %s",
					recordType,
					record.ID,
					strings.Join(targets, ", "),
				),
				document.Path,
				successors[record.ID][targets[1]].line,
				record.ID,
				"replacement",
			))
		}

		if record.Lifecycle.DeprecatedBy != "" {
			replacement, exists := byID[record.Lifecycle.DeprecatedBy]
			if exists && len(replacement.Lifecycle.Supersedes) > 0 &&
				!stringSet(replacement.Lifecycle.Supersedes)[record.ID] {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-schema",
					fmt.Sprintf(
						"%s %s names Deprecated by %s, but that record's Supersedes field does not include %s",
						recordType,
						record.ID,
						replacement.ID,
						record.ID,
					),
					document.Path,
					document.Index.recordFieldLine(section, replacement.Index, "supersedes"),
					record.ID,
					"replacement",
				))
			}
		}
		for _, oldID := range record.Lifecycle.Supersedes {
			oldRecord, exists := byID[oldID]
			if !exists || oldRecord.Lifecycle.DeprecatedBy == "" ||
				oldRecord.Lifecycle.DeprecatedBy == record.ID {
				continue
			}
			errors = append(errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf(
					"%s %s supersedes %s, but that record names Deprecated by %s",
					recordType,
					record.ID,
					oldID,
					oldRecord.Lifecycle.DeprecatedBy,
				),
				document.Path,
				document.Index.recordFieldLine(section, record.Index, "supersedes"),
				record.ID,
				"replacement",
			))
		}
	}

	state := make(map[string]uint8, len(records))
	var stack []string
	reported := make(map[string]bool)
	var visit func(string)
	visit = func(recordID string) {
		state[recordID] = 1
		stack = append(stack, recordID)
		for _, target := range sortedReplacementTargets(successors[recordID]) {
			switch state[target] {
			case 0:
				visit(target)
			case 1:
				cycleKey := lifecycleCycleKey(stack, target)
				if reported[cycleKey] {
					continue
				}
				reported[cycleKey] = true
				location := successors[recordID][target]
				errors = append(errors, iddDocumentValidationError(
					"idd-document-schema",
					fmt.Sprintf(
						"%s replacement cycle detected: %s -> %s",
						recordType,
						strings.Join(stack, " -> "),
						target,
					),
					document.Path,
					location.line,
					location.recordID,
					"replacement-cycle",
				))
			}
		}
		stack = stack[:len(stack)-1]
		state[recordID] = 2
	}
	for _, record := range records {
		if record.ID != "" && state[record.ID] == 0 {
			visit(record.ID)
		}
	}
	return errors
}

func sortedReplacementTargets(targets map[string]iddReplacementLocation) []string {
	result := make([]string, 0, len(targets))
	for target := range targets {
		result = append(result, target)
	}
	sort.Strings(result)
	return result
}

func lifecycleCycleKey(stack []string, target string) string {
	for index, recordID := range stack {
		if recordID == target {
			cycle := append([]string(nil), stack[index:]...)
			sort.Strings(cycle)
			return strings.Join(cycle, "\x00")
		}
	}
	return target
}

func iddRecordHasConcernSection(
	document *parsedIDDDocument,
	section string,
	recordIndex int,
	concern string,
) bool {
	root := goldmark.DefaultParser().Parse(goldmarktext.NewReader(document.Body))
	records := splitIDDMarkdownRecords(document, root)
	recordLine := document.Index.recordLine(section, recordIndex)
	for _, record := range records {
		if record.line != recordLine {
			continue
		}
		for blockIndex, block := range record.blocks {
			heading, ok := block.(*goldmarkast.Heading)
			if !ok || heading.Level != 3 ||
				!strings.EqualFold(markdownNodeText(heading, document.Body), concern) {
				continue
			}
			var content []goldmarkast.Node
			for _, next := range record.blocks[blockIndex+1:] {
				if nextHeading, headingOK := next.(*goldmarkast.Heading); headingOK && nextHeading.Level <= 3 {
					break
				}
				content = append(content, next)
			}
			return completionRecordHasEffectiveContent(
				iddMarkdownRecord{blocks: content},
				document.Body,
			)
		}
	}
	return false
}

func splitIDDScopedName(currentPackage, value string) (string, string) {
	targetPackage, targetName, qualified := strings.Cut(value, "#")
	if !qualified {
		return currentPackage, value
	}
	return filepath.ToSlash(filepath.Clean(targetPackage)), targetName
}

func validateIDDScopedName(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("reference must not be empty")
	}
	hashCount := strings.Count(value, "#")
	if hashCount == 0 {
		return nil
	}
	if hashCount != 1 {
		return fmt.Errorf("use exactly one # separator in <package>#<name>")
	}
	targetPackage, targetName, _ := strings.Cut(value, "#")
	if strings.TrimSpace(targetPackage) == "" || strings.TrimSpace(targetName) == "" {
		return fmt.Errorf("both package and name are required in <package>#<name>")
	}
	cleanPackage, err := validateDocumentPackagePath(targetPackage)
	if err != nil {
		return fmt.Errorf("invalid package in <package>#<name>: %w", err)
	}
	if filepath.ToSlash(cleanPackage) != targetPackage {
		return fmt.Errorf("package must use a normalized project-relative slash path")
	}
	return nil
}

func validateIDDSpecs(document *parsedIDDDocument, packagePath string) []*model.ValidationError {
	var errors []*model.ValidationError
	module := moduleFromPackage(packagePath)
	seen := make(map[string]bool, len(document.Metadata.Specs))
	declared := make(map[string]bool, len(document.Metadata.Specs))
	for _, spec := range document.Metadata.Specs {
		if spec.ID != "" {
			declared[spec.ID] = true
		}
	}
	for recordIndex, spec := range document.Metadata.Specs {
		line := document.Index.recordLine("specs", recordIndex)
		validateIDDRequiredRecordFields(&errors, document, "spec", recordIndex)
		validateIDDRecordLifecycle(
			&errors,
			document,
			"specs",
			recordIndex,
			"SPEC",
			spec.ID,
			spec.IDDRecordLifecycle,
			declared,
		)

		if spec.ID == "" {
			continue
		}
		if err := pattern.ValidateIdentifierFormat(spec.ID); err != nil || !strings.HasPrefix(spec.ID, "SPEC-") {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf("SPEC record has invalid id %q; expected SPEC-<MODULE>-<NUMBER>", spec.ID),
				document.Path,
				document.Index.recordFieldLine("specs", recordIndex, "id"),
				spec.ID,
				"id",
			))
		} else if module != "" && !strings.HasPrefix(spec.ID, "SPEC-"+module+"-") {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf("SPEC id %s does not use package-derived module %s", spec.ID, module),
				document.Path,
				document.Index.recordFieldLine("specs", recordIndex, "id"),
				spec.ID,
				"id",
			))
		}
		if seen[spec.ID] {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf("duplicate SPEC id %s", spec.ID),
				document.Path,
				line,
				spec.ID,
				"id",
			))
		}
		seen[spec.ID] = true
	}
	errors = append(errors, validateIDDLifecycleRelationships(
		document,
		"specs",
		"SPEC",
		specLifecycleRecords(document.Metadata.Specs),
	)...)
	return errors
}

func validateIDDTests(document *parsedIDDDocument, packagePath string) []*model.ValidationError {
	var errors []*model.ValidationError
	module := moduleFromPackage(packagePath)
	seen := make(map[string]bool, len(document.Metadata.Tests))
	declared := make(map[string]bool, len(document.Metadata.Tests))
	for _, test := range document.Metadata.Tests {
		if test.ID != "" {
			declared[test.ID] = true
		}
	}
	for recordIndex, test := range document.Metadata.Tests {
		line := document.Index.recordLine("tests", recordIndex)
		validateIDDRequiredRecordFields(&errors, document, "testing", recordIndex)
		validateIDDRecordLifecycle(
			&errors,
			document,
			"tests",
			recordIndex,
			"TEST",
			test.ID,
			test.IDDRecordLifecycle,
			declared,
		)

		if test.ID != "" {
			if err := pattern.ValidateIdentifierFormat(test.ID); err != nil || !strings.HasPrefix(test.ID, "TEST-") {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-schema",
					fmt.Sprintf("TEST record has invalid id %q; expected TEST-<MODULE>-<NUMBER>", test.ID),
					document.Path,
					document.Index.recordFieldLine("tests", recordIndex, "id"),
					test.ID,
					"id",
				))
			} else if module != "" && !strings.HasPrefix(test.ID, "TEST-"+module+"-") {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-schema",
					fmt.Sprintf("TEST id %s does not use package-derived module %s", test.ID, module),
					document.Path,
					document.Index.recordFieldLine("tests", recordIndex, "id"),
					test.ID,
					"id",
				))
			}
			if seen[test.ID] {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-schema",
					fmt.Sprintf("duplicate TEST id %s", test.ID),
					document.Path,
					line,
					test.ID,
					"id",
				))
			}
			seen[test.ID] = true
		}

		if test.Kind != "" && test.Kind != "test" && test.Kind != "contract" {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf("TEST %s kind must be test or contract", test.ID),
				document.Path,
				document.Index.recordFieldLine("tests", recordIndex, "kind"),
				test.ID,
				"kind",
			))
		}
		if test.Kind == "test" && len(test.Contracts) > 0 {
			errors = append(errors, iddDocumentValidationError(
				"idd-document-schema",
				fmt.Sprintf("behavior TEST %s must not declare Contracts; use kind contract", test.ID),
				document.Path,
				document.Index.recordFieldLine("tests", recordIndex, "contracts"),
				test.ID,
				"contracts",
			))
		}
		covered := make(map[string]bool, len(test.Covers))
		for _, specID := range test.Covers {
			if err := pattern.ValidateIdentifierFormat(specID); err != nil ||
				!strings.HasPrefix(specID, "SPEC-") {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-schema",
					fmt.Sprintf(
						"TEST %s Covers entry %q must be a valid SPEC identifier",
						test.ID,
						specID,
					),
					document.Path,
					document.Index.recordFieldLine("tests", recordIndex, "covers"),
					test.ID,
					"covers",
				))
				continue
			}
			if covered[specID] {
				errors = append(errors, iddDocumentValidationError(
					"idd-document-schema",
					fmt.Sprintf("TEST %s repeats covered SPEC %s", test.ID, specID),
					document.Path,
					document.Index.recordFieldLine("tests", recordIndex, "covers"),
					test.ID,
					"covers",
				))
			}
			covered[specID] = true
		}
	}
	errors = append(errors, validateIDDLifecycleRelationships(
		document,
		"tests",
		"TEST",
		testLifecycleRecords(document.Metadata.Tests),
	)...)
	return errors
}

func validateIDDRequiredRecordFields(
	errors *[]*model.ValidationError,
	document *parsedIDDDocument,
	role string,
	recordIndex int,
) {
	schema, ok := iddDocumentSchemas[role]
	if !ok {
		return
	}
	records := iddDocumentRecordValuesForRole(document, role)
	if recordIndex < 0 || recordIndex >= len(records) {
		return
	}
	record := records[recordIndex]
	for _, field := range schema.Record.Fields {
		if field.Name == "name" ||
			!iddDocumentRecordFieldRequired(field, record.Kind) {
			continue
		}
		value := record.scalarValue(field.Name)
		if field.List {
			value = strings.Join(record.Lists[field.Name], ", ")
		}
		validateIDDRecordField(
			errors,
			document,
			schema.Record.Type,
			record.Identifier,
			field.Name,
			value,
			iddDocumentRecordFieldLine(
				document,
				schema.Record.Section,
				record.Index,
				field.Name,
			),
		)
	}
}

func validateIDDRecordField(
	errors *[]*model.ValidationError,
	document *parsedIDDDocument,
	recordType string,
	recordID string,
	field string,
	value string,
	line int,
) {
	value = strings.TrimSpace(value)
	if value == "" {
		label := recordID
		if label == "" {
			label = recordType + " record"
		}
		*errors = append(*errors, iddDocumentValidationError(
			"idd-document-schema",
			fmt.Sprintf("%s field %s is required", label, field),
			document.Path,
			line,
			recordID,
			field,
		))
		return
	}
	if iddDocumentRecordFieldComplete(value, recordID, field) {
		return
	}
	*errors = append(*errors, iddDocumentValidationError(
		"idd-document-schema",
		fmt.Sprintf("%s field %s must contain concrete, meaningful content instead of %q", recordID, field, value),
		document.Path,
		line,
		recordID,
		field,
	))
}

func isIDDPlaceholder(value, recordID, field string) bool {
	normalized := strings.ToLower(strings.TrimSpace(value))
	return (field != "id" && recordID != "" && value == recordID) ||
		normalized == "tbd" ||
		normalized == "todo" ||
		normalized == "auto-generated" ||
		normalized == "autogenerated"
}

func validateIDDDocumentReferences(documents map[string]*parsedIDDDocument) []*model.ValidationError {
	var errors []*model.ValidationError
	designs := make(map[string]bool)
	contracts := make(map[string]bool)
	specs := make(map[string]bool)
	if document := documents["design"]; document != nil {
		designs = stringSet(document.Metadata.Components)
	}
	if document := documents["contract"]; document != nil {
		contracts = stringSet(document.Metadata.Contracts)
	}
	if document := documents["spec"]; document != nil {
		for _, spec := range document.Metadata.Specs {
			if spec.ID != "" {
				specs[spec.ID] = true
			}
		}
		for recordIndex, spec := range document.Metadata.Specs {
			if spec.Design != "" {
				line := document.Index.recordFieldLine("specs", recordIndex, "design")
				if err := validateIDDScopedName(spec.Design); err != nil {
					errors = append(errors, iddDocumentValidationError(
						"idd-document-schema",
						fmt.Sprintf("SPEC %s has invalid design reference %q: %v", spec.ID, spec.Design, err),
						document.Path,
						line,
						spec.ID,
						"design",
					))
				} else {
					targetPackage, targetName := splitIDDScopedName(document.Metadata.Package, spec.Design)
					if targetPackage == document.Metadata.Package && !designs[targetName] {
						errors = append(errors, iddDocumentValidationError(
							"idd-document-reference",
							fmt.Sprintf("SPEC %s references undeclared design %q", spec.ID, targetName),
							document.Path,
							line,
							spec.ID,
							"design",
						))
					}
				}
			}
			if spec.Contract != "" {
				line := document.Index.recordFieldLine("specs", recordIndex, "contract")
				if err := validateIDDScopedName(spec.Contract); err != nil {
					errors = append(errors, iddDocumentValidationError(
						"idd-document-schema",
						fmt.Sprintf("SPEC %s has invalid contract reference %q: %v", spec.ID, spec.Contract, err),
						document.Path,
						line,
						spec.ID,
						"contract",
					))
				} else {
					targetPackage, targetName := splitIDDScopedName(document.Metadata.Package, spec.Contract)
					if targetPackage == document.Metadata.Package && !contracts[targetName] {
						errors = append(errors, iddDocumentValidationError(
							"idd-document-reference",
							fmt.Sprintf("SPEC %s references undeclared contract %q", spec.ID, targetName),
							document.Path,
							line,
							spec.ID,
							"contract",
						))
					}
				}
			}
		}
	}
	if document := documents["testing"]; document != nil {
		for recordIndex, test := range document.Metadata.Tests {
			for _, specID := range test.Covers {
				if pattern.ValidateIdentifierFormat(specID) != nil ||
					!strings.HasPrefix(specID, "SPEC-") {
					continue
				}
				if specs[specID] {
					continue
				}
				errors = append(errors, iddDocumentValidationError(
					"idd-document-reference",
					fmt.Sprintf("TEST %s covers undefined SPEC %s", test.ID, specID),
					document.Path,
					document.Index.recordFieldLine("tests", recordIndex, "covers"),
					test.ID,
					"covers",
				))
			}
			seenContracts := make(map[string]bool)
			for _, contractRef := range test.Contracts {
				line := document.Index.recordFieldLine("tests", recordIndex, "contracts")
				if seenContracts[contractRef] {
					errors = append(errors, iddDocumentValidationError(
						"idd-document-schema",
						fmt.Sprintf("TEST %s repeats Contract reference %q", test.ID, contractRef),
						document.Path,
						line,
						test.ID,
						"contracts",
					))
					continue
				}
				seenContracts[contractRef] = true
				if err := validateIDDScopedName(contractRef); err != nil {
					errors = append(errors, iddDocumentValidationError(
						"idd-document-schema",
						fmt.Sprintf(
							"TEST %s has invalid Contract reference %q: %v",
							test.ID,
							contractRef,
							err,
						),
						document.Path,
						line,
						test.ID,
						"contracts",
					))
					continue
				}
				targetPackage, targetName := splitIDDScopedName(document.Metadata.Package, contractRef)
				if targetPackage == document.Metadata.Package && !contracts[targetName] {
					errors = append(errors, iddDocumentValidationError(
						"idd-document-reference",
						fmt.Sprintf("TEST %s references undeclared Contract %q", test.ID, targetName),
						document.Path,
						line,
						test.ID,
						"contracts",
					))
				}
			}
		}
	}
	return errors
}

func addIDDDocumentIdentifiers(documents map[string]*parsedIDDDocument, set *model.IdentifierSet) {
	if document := documents["design"]; document != nil {
		for recordIndex, component := range document.Metadata.ComponentRecords {
			if component.Name == "" {
				continue
			}
			identifier := model.NewIdentifierWithDescribe(
				derivedComponentID(document.Metadata.Package, component.Name),
				model.TypeDesign,
				component.Name,
				component.Purpose,
				document.Path,
				document.Index.recordLine("components", recordIndex),
			)
			identifier.RawRef = component.Name
			identifier.Kind = "component"
			identifier.Derived = true
			for _, dependency := range component.DependsOn {
				targetPackage, targetName := splitIDDScopedName(document.Metadata.Package, dependency)
				if targetName != "" {
					identifier.AddTypedLink(
						derivedComponentID(targetPackage, targetName),
						model.LinkDependsOn,
					)
				}
			}
			addIDDLifecycleLinks(identifier, component.IDDRecordLifecycle, func(name string) string {
				return derivedComponentID(document.Metadata.Package, name)
			})
			set.Add(identifier)
		}
	}
	if document := documents["contract"]; document != nil {
		for recordIndex, contract := range document.Metadata.ContractRecords {
			if contract.Name == "" {
				continue
			}
			identifier := model.NewIdentifierWithDescribe(
				derivedContractID(document.Metadata.Package, contract.Name),
				model.TypeContract,
				contract.Name,
				contract.Guarantees,
				document.Path,
				document.Index.recordLine("contracts", recordIndex),
			)
			identifier.RawRef = contract.Name
			identifier.Kind = "contract"
			identifier.Derived = true
			addIDDLifecycleLinks(identifier, contract.IDDRecordLifecycle, func(name string) string {
				return derivedContractID(document.Metadata.Package, name)
			})
			set.Add(identifier)
		}
	}

	specIdentifiers := make(map[string]*model.Identifier)
	if document := documents["spec"]; document != nil {
		for recordIndex, spec := range document.Metadata.Specs {
			if spec.ID == "" || !strings.HasPrefix(spec.ID, "SPEC-") {
				continue
			}
			identifier := model.NewIdentifierWithDescribe(
				spec.ID,
				model.TypeSpec,
				spec.Title,
				spec.Requirement,
				document.Path,
				document.Index.recordLine("specs", recordIndex),
			)
			identifier.RawRef = spec.ID
			if spec.Design != "" {
				targetPackage, targetName := splitIDDScopedName(document.Metadata.Package, spec.Design)
				if targetName != "" {
					identifier.AddTypedLink(
						derivedComponentID(targetPackage, targetName),
						model.LinkReferences,
					)
				}
			}
			if spec.Contract != "" {
				targetPackage, targetName := splitIDDScopedName(document.Metadata.Package, spec.Contract)
				if targetName != "" {
					identifier.AddTypedLink(
						derivedContractID(targetPackage, targetName),
						model.LinkContract,
					)
				}
			}
			addIDDLifecycleLinks(identifier, spec.IDDRecordLifecycle, func(id string) string {
				return id
			})
			set.Add(identifier)
			specIdentifiers[spec.ID] = identifier
		}
	}
	if document := documents["testing"]; document != nil {
		for recordIndex, test := range document.Metadata.Tests {
			if test.ID == "" || !strings.HasPrefix(test.ID, "TEST-") {
				continue
			}
			identifier := model.NewIdentifierWithDescribe(
				test.ID,
				model.TypeTest,
				test.Title,
				test.Purpose,
				document.Path,
				document.Index.recordLine("tests", recordIndex),
			)
			identifier.RawRef = test.ID
			identifier.Kind = test.Kind
			for _, specID := range test.Covers {
				if pattern.ValidateIdentifierFormat(specID) != nil ||
					!strings.HasPrefix(specID, "SPEC-") {
					continue
				}
				identifier.AddTypedLink(specID, model.LinkImplements)
				if spec := specIdentifiers[specID]; spec != nil {
					spec.AddTypedLink(test.ID, model.LinkTests)
				}
			}
			for _, contractRef := range test.Contracts {
				targetPackage, targetName := splitIDDScopedName(document.Metadata.Package, contractRef)
				if targetName != "" {
					identifier.AddTypedLink(
						derivedContractID(targetPackage, targetName),
						model.LinkContractTests,
					)
				}
			}
			addIDDLifecycleLinks(identifier, test.IDDRecordLifecycle, func(id string) string {
				return id
			})
			set.Add(identifier)
		}
	}
}

func addIDDLifecycleLinks(
	identifier *model.Identifier,
	lifecycle IDDRecordLifecycle,
	resolve func(string) string,
) {
	for _, target := range lifecycle.Supersedes {
		identifier.AddTypedLink(resolve(target), model.LinkSupersedes)
	}
	if lifecycle.DeprecatedBy != "" {
		identifier.AddTypedLink(resolve(lifecycle.DeprecatedBy), model.LinkDeprecatedBy)
	}
}

func derivedComponentID(packagePath, name string) string {
	return "component:" + filepath.ToSlash(filepath.Clean(packagePath)) + "#" + name
}

func derivedContractID(packagePath, name string) string {
	return "contract:" + filepath.ToSlash(filepath.Clean(packagePath)) + "#" + name
}

func declaredIDDIdentifiers(documents map[string]*parsedIDDDocument) map[string]bool {
	declared := make(map[string]bool)
	for _, document := range documents {
		for _, spec := range document.Metadata.Specs {
			if spec.ID != "" {
				declared[spec.ID] = true
			}
		}
		for _, test := range document.Metadata.Tests {
			if test.ID != "" {
				declared[test.ID] = true
			}
		}
	}
	return declared
}
