// Package collector provides syntax-tree-backed source declaration collection.

// Spec: docs/internal/collector/spec.md
// Contract: docs/internal/collector/contract.md
package collector

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_cpp "github.com/tree-sitter/tree-sitter-cpp/bindings/go"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
	tree_sitter_java "github.com/tree-sitter/tree-sitter-java/bindings/go"
	tree_sitter_javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	tree_sitter_python "github.com/tree-sitter/tree-sitter-python/bindings/go"
	tree_sitter_typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
)

var supportedSourceExtensions = map[string]bool{
	".go":   true,
	".ts":   true,
	".tsx":  true,
	".js":   true,
	".jsx":  true,
	".cpp":  true,
	".cc":   true,
	".cxx":  true,
	".hpp":  true,
	".hh":   true,
	".hxx":  true,
	".java": true,
	".py":   true,
}

// SourceParseError identifies syntax that prevented trustworthy declaration
// binding.
// @implement SPEC-INTERNAL_COLLECTOR-024
type SourceParseError struct {
	Line    int
	Message string
}

// SourceAnnotation is an IDD annotation read from a real syntax-tree comment.
// @implement SPEC-INTERNAL_COLLECTOR-024
type SourceAnnotation struct {
	Kind            string
	Prefix          string
	Refs            []string
	Line            int
	Raw             string
	Attached        bool
	Declaration     string
	DeclarationLine int
	HasSeparator    bool
	Ignored         bool
	commentEnd      uint
}

// SourceDeclaration is a language-neutral declaration used by collector and
// engine validation.
// @implement SPEC-INTERNAL_COLLECTOR-024
type SourceDeclaration struct {
	Name         string
	Kind         string
	Line         int
	EndLine      int
	Public       bool
	TestKind     string
	Ignored      bool
	Annotations  []SourceAnnotation
	bindingStart uint
}

// SourceAnalysis is the normalized result of parsing one supported source
// file.
// @implement SPEC-INTERNAL_COLLECTOR-024
type SourceAnalysis struct {
	Path         string
	Language     string
	Declarations []SourceDeclaration
	Annotations  []SourceAnnotation
	ParseErrors  []SourceParseError
}

type sourceLanguageProfile struct {
	name     string
	language *tree_sitter.Language
}

type sourceComment struct {
	start   uint
	end     uint
	line    int
	raw     string
	ignored bool
}

type sourceAnnotationMarker struct {
	kind   string
	prefix string
}

// SupportedSourcePath reports whether idd-cli has a pinned Tree-sitter grammar
// for the path's extension.
// @implement SPEC-INTERNAL_COLLECTOR-024
func SupportedSourcePath(path string) bool {
	return supportedSourceExtensions[strings.ToLower(filepath.Ext(path))]
}

// AnalyzeSource parses one supported source file and binds IDD comment
// annotations to actual declarations.
// @implement SPEC-INTERNAL_COLLECTOR-024
func AnalyzeSource(path string, source []byte) (*SourceAnalysis, error) {
	return AnalyzeSourceWithAnnotations(path, source, nil)
}

// AnalyzeSourceWithAnnotations parses source using configured annotation
// prefixes while retaining canonical implement, test, and test-contract kinds.
// @implement SPEC-INTERNAL_COLLECTOR-024
func AnalyzeSourceWithAnnotations(
	path string,
	source []byte,
	annotations map[string]string,
) (*SourceAnalysis, error) {
	profile, ok := sourceProfileForPath(path)
	if !ok {
		return nil, fmt.Errorf("unsupported source extension %q", filepath.Ext(path))
	}
	annotationMarkers := configuredSourceAnnotationMarkers(annotations)

	parser := tree_sitter.NewParser()
	defer parser.Close()
	if err := parser.SetLanguage(profile.language); err != nil {
		return nil, fmt.Errorf("load %s grammar: %w", profile.name, err)
	}
	tree := parser.Parse(source, nil)
	if tree == nil {
		return nil, fmt.Errorf("parse %s: parser returned no syntax tree", path)
	}
	defer tree.Close()

	analysis := &SourceAnalysis{
		Path:         path,
		Language:     profile.name,
		Declarations: make([]SourceDeclaration, 0),
		Annotations:  make([]SourceAnnotation, 0),
		ParseErrors:  make([]SourceParseError, 0),
	}
	root := tree.RootNode()
	if root == nil {
		return nil, fmt.Errorf("parse %s: syntax tree has no root", path)
	}
	if root.HasError() {
		errorNode := firstSourceErrorNode(root)
		line := 1
		message := "source contains syntax errors"
		if errorNode != nil {
			line = int(errorNode.StartPosition().Row) + 1
			if errorNode.IsMissing() {
				message = fmt.Sprintf("source is missing %s", errorNode.Kind())
			}
		}
		analysis.ParseErrors = append(analysis.ParseErrors, SourceParseError{
			Line:    line,
			Message: message,
		})
		return analysis, nil
	}

	var comments []sourceComment
	walkSourceTree(root, func(node *tree_sitter.Node) {
		if isSourceCommentNode(profile.name, node.Kind()) {
			comments = append(comments, sourceComment{
				start: node.StartByte(),
				end:   node.EndByte(),
				line:  int(node.StartPosition().Row) + 1,
				raw:   node.Utf8Text(source),
			})
		}
		if declaration, found := sourceDeclarationForNode(profile.name, path, node, source); found {
			analysis.Declarations = append(analysis.Declarations, declaration)
		}
	})
	sort.Slice(comments, func(i, j int) bool { return comments[i].start < comments[j].start })
	ignoredLines := sourceIgnoreLines(comments, len(strings.Split(string(source), "\n")))
	for index := range comments {
		comments[index].ignored = lineIsIgnored(ignoredLines, comments[index].line)
		for _, annotation := range sourceAnnotationsFromComment(comments[index], annotationMarkers) {
			annotation.Ignored = comments[index].ignored
			analysis.Annotations = append(analysis.Annotations, annotation)
		}
	}
	for index := range analysis.Declarations {
		analysis.Declarations[index].Ignored = lineIsIgnored(ignoredLines, analysis.Declarations[index].Line)
	}

	sort.Slice(analysis.Declarations, func(i, j int) bool {
		if analysis.Declarations[i].bindingStart != analysis.Declarations[j].bindingStart {
			return analysis.Declarations[i].bindingStart < analysis.Declarations[j].bindingStart
		}
		return analysis.Declarations[i].Name < analysis.Declarations[j].Name
	})
	bindSourceAnnotations(analysis, comments, source)
	return analysis, nil
}

func configuredSourceAnnotationMarkers(annotations map[string]string) []sourceAnnotationMarker {
	if len(annotations) == 0 {
		annotations = map[string]string{
			"spec":          "@implement",
			"test":          "@test",
			"test_contract": "@test-contract",
		}
	}
	kinds := map[string]string{
		"spec":          "implement",
		"test":          "test",
		"test_contract": "test-contract",
	}
	markers := make([]sourceAnnotationMarker, 0, len(kinds))
	for role, kind := range kinds {
		if prefix := strings.TrimSpace(annotations[role]); prefix != "" {
			markers = append(markers, sourceAnnotationMarker{kind: kind, prefix: prefix})
		}
	}
	sort.Slice(markers, func(i, j int) bool {
		if len(markers[i].prefix) != len(markers[j].prefix) {
			return len(markers[i].prefix) > len(markers[j].prefix)
		}
		return markers[i].kind < markers[j].kind
	})
	return markers
}

func sourceProfileForPath(path string) (sourceLanguageProfile, bool) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return sourceLanguageProfile{"go", tree_sitter.NewLanguage(tree_sitter_go.Language())}, true
	case ".ts":
		return sourceLanguageProfile{"typescript", tree_sitter.NewLanguage(tree_sitter_typescript.LanguageTypescript())}, true
	case ".tsx":
		return sourceLanguageProfile{"tsx", tree_sitter.NewLanguage(tree_sitter_typescript.LanguageTSX())}, true
	case ".js", ".jsx":
		return sourceLanguageProfile{"javascript", tree_sitter.NewLanguage(tree_sitter_javascript.Language())}, true
	case ".cpp", ".cc", ".cxx", ".hpp", ".hh", ".hxx":
		return sourceLanguageProfile{"cpp", tree_sitter.NewLanguage(tree_sitter_cpp.Language())}, true
	case ".java":
		return sourceLanguageProfile{"java", tree_sitter.NewLanguage(tree_sitter_java.Language())}, true
	case ".py":
		return sourceLanguageProfile{"python", tree_sitter.NewLanguage(tree_sitter_python.Language())}, true
	default:
		return sourceLanguageProfile{}, false
	}
}

func walkSourceTree(node *tree_sitter.Node, visit func(*tree_sitter.Node)) {
	if node == nil {
		return
	}
	visit(node)
	for index := uint(0); index < node.NamedChildCount(); index++ {
		walkSourceTree(node.NamedChild(index), visit)
	}
}

func firstSourceErrorNode(node *tree_sitter.Node) *tree_sitter.Node {
	if node == nil {
		return nil
	}
	if node.IsError() || node.IsMissing() {
		return node
	}
	for index := uint(0); index < node.NamedChildCount(); index++ {
		if found := firstSourceErrorNode(node.NamedChild(index)); found != nil {
			return found
		}
	}
	return nil
}

func isSourceCommentNode(language, kind string) bool {
	switch language {
	case "java":
		return kind == "line_comment" || kind == "block_comment"
	default:
		return kind == "comment"
	}
}

func sourceDeclarationForNode(
	language string,
	path string,
	node *tree_sitter.Node,
	source []byte,
) (SourceDeclaration, bool) {
	executableTestCall := false
	if language == "javascript" || language == "typescript" || language == "tsx" {
		if node.Kind() == "expression_statement" {
			name, _ := javascriptTestCall(node, source)
			executableTestCall = name != ""
		}
	}
	if hasExecutableSourceAncestor(language, node) && !executableTestCall {
		return SourceDeclaration{}, false
	}

	var name, kind string
	bindingNode := node
	switch language {
	case "go":
		switch node.Kind() {
		case "function_declaration":
			name, kind = sourceNodeFieldText(node, "name", source), "function"
		case "method_declaration":
			name, kind = sourceNodeFieldText(node, "name", source), "method"
		case "type_spec":
			name, kind = sourceNodeFieldText(node, "name", source), "type"
			if parent := node.Parent(); parent != nil && parent.Kind() == "type_declaration" {
				if !sourceNodeHasDirectChild(parent, "(") {
					bindingNode = parent
				}
			}
		}
	case "javascript", "typescript", "tsx":
		switch node.Kind() {
		case "function_declaration", "generator_function_declaration", "function_signature":
			name, kind = sourceNodeFieldText(node, "name", source), "function"
		case "class_declaration", "abstract_class_declaration":
			name, kind = sourceNodeFieldText(node, "name", source), "class"
		case "interface_declaration":
			name, kind = sourceNodeFieldText(node, "name", source), "interface"
		case "type_alias_declaration":
			name, kind = sourceNodeFieldText(node, "name", source), "type"
		case "enum_declaration":
			name, kind = sourceNodeFieldText(node, "name", source), "enum"
		case "method_definition", "method_signature", "abstract_method_signature":
			name, kind = sourceNodeFieldText(node, "name", source), "method"
		case "lexical_declaration", "variable_declaration":
			name, kind = javascriptFunctionVariable(node, source)
		case "expression_statement":
			name, kind = javascriptTestCall(node, source)
		}
		if parent := node.Parent(); parent != nil && parent.Kind() == "export_statement" {
			bindingNode = parent
		}
	case "cpp":
		switch node.Kind() {
		case "function_definition":
			name, kind = cppDeclaratorName(node.ChildByFieldName("declarator"), source), "function"
		case "declaration", "field_declaration":
			if declarator := firstDescendantOfKind(node, "function_declarator"); declarator != nil {
				name, kind = cppDeclaratorName(declarator, source), "function"
			}
		case "class_specifier":
			name, kind = sourceNodeFieldText(node, "name", source), "class"
		case "struct_specifier":
			name, kind = sourceNodeFieldText(node, "name", source), "struct"
		case "enum_specifier":
			name, kind = sourceNodeFieldText(node, "name", source), "enum"
		case "type_definition":
			name, kind = sourceNodeFieldText(node, "declarator", source), "type"
		case "alias_declaration":
			name, kind = sourceNodeFieldText(node, "name", source), "type"
		case "concept_definition":
			name, kind = sourceNodeFieldText(node, "name", source), "concept"
		}
		for parent := node.Parent(); parent != nil; parent = parent.Parent() {
			if parent.Kind() == "template_declaration" {
				bindingNode = parent
				break
			}
			if parent.Kind() == "field_declaration_list" ||
				parent.Kind() == "namespace_definition" ||
				parent.Kind() == "translation_unit" {
				break
			}
		}
	case "java":
		switch node.Kind() {
		case "class_declaration":
			name, kind = sourceNodeFieldText(node, "name", source), "class"
		case "interface_declaration":
			name, kind = sourceNodeFieldText(node, "name", source), "interface"
		case "enum_declaration":
			name, kind = sourceNodeFieldText(node, "name", source), "enum"
		case "record_declaration":
			name, kind = sourceNodeFieldText(node, "name", source), "record"
		case "annotation_type_declaration":
			name, kind = sourceNodeFieldText(node, "name", source), "annotation"
		case "method_declaration":
			name, kind = sourceNodeFieldText(node, "name", source), "method"
		case "constructor_declaration":
			name, kind = sourceNodeFieldText(node, "name", source), "constructor"
		}
	case "python":
		switch node.Kind() {
		case "function_definition":
			name, kind = sourceNodeFieldText(node, "name", source), "function"
		case "class_definition":
			name, kind = sourceNodeFieldText(node, "name", source), "class"
		}
		if parent := node.Parent(); parent != nil && parent.Kind() == "decorated_definition" {
			bindingNode = parent
		}
	}
	if name == "" || kind == "" {
		return SourceDeclaration{}, false
	}

	declaration := SourceDeclaration{
		Name:         trimSourceName(name),
		Kind:         kind,
		Line:         int(node.StartPosition().Row) + 1,
		EndLine:      int(node.EndPosition().Row) + 1,
		bindingStart: bindingNode.StartByte(),
	}
	declaration.Public = sourceDeclarationIsPublic(language, node, declaration, source)
	declaration.TestKind = sourceDeclarationTestKind(language, path, node, declaration, source)
	return declaration, true
}

func sourceNodeHasDirectChild(node *tree_sitter.Node, kind string) bool {
	if node == nil {
		return false
	}
	for index := uint(0); index < node.ChildCount(); index++ {
		child := node.Child(index)
		if child != nil && child.Kind() == kind {
			return true
		}
	}
	return false
}

func hasExecutableSourceAncestor(language string, node *tree_sitter.Node) bool {
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		switch language {
		case "go":
			if parent.Kind() == "function_declaration" || parent.Kind() == "method_declaration" ||
				parent.Kind() == "func_literal" {
				return true
			}
		case "javascript", "typescript", "tsx":
			switch parent.Kind() {
			case "function_declaration", "function_expression", "arrow_function",
				"generator_function_declaration", "method_definition":
				return true
			}
		case "cpp":
			if parent.Kind() == "function_definition" || parent.Kind() == "lambda_expression" {
				return true
			}
		case "java":
			if parent.Kind() == "method_declaration" || parent.Kind() == "constructor_declaration" ||
				parent.Kind() == "lambda_expression" {
				return true
			}
		case "python":
			if parent.Kind() == "function_definition" || parent.Kind() == "lambda" {
				return true
			}
		}
	}
	return false
}

func sourceNodeFieldText(node *tree_sitter.Node, field string, source []byte) string {
	if node == nil {
		return ""
	}
	child := node.ChildByFieldName(field)
	if child == nil {
		return ""
	}
	return child.Utf8Text(source)
}

func javascriptFunctionVariable(node *tree_sitter.Node, source []byte) (string, string) {
	for index := uint(0); index < node.NamedChildCount(); index++ {
		child := node.NamedChild(index)
		if child == nil || child.Kind() != "variable_declarator" {
			continue
		}
		value := child.ChildByFieldName("value")
		if value == nil || value.Kind() != "arrow_function" && value.Kind() != "function_expression" {
			continue
		}
		return sourceNodeFieldText(child, "name", source), "function"
	}
	return "", ""
}

func javascriptTestCall(node *tree_sitter.Node, source []byte) (string, string) {
	call := firstDescendantOfKind(node, "call_expression")
	if call == nil {
		return "", ""
	}
	function := call.ChildByFieldName("function")
	if function == nil {
		return "", ""
	}
	callName := function.Utf8Text(source)
	if callName != "test" && callName != "it" &&
		!strings.HasPrefix(callName, "test.") &&
		!strings.HasPrefix(callName, "it.") {
		return "", ""
	}
	arguments := call.ChildByFieldName("arguments")
	if arguments == nil || arguments.NamedChildCount() == 0 {
		return callName, "test"
	}
	name := trimSourceName(arguments.NamedChild(0).Utf8Text(source))
	if name == "" {
		name = callName
	}
	return name, "test"
}

func cppDeclaratorName(node *tree_sitter.Node, source []byte) string {
	if node == nil {
		return ""
	}
	switch node.Kind() {
	case "identifier", "field_identifier", "type_identifier", "operator_name", "destructor_name":
		return node.Utf8Text(source)
	case "qualified_identifier":
		if name := node.ChildByFieldName("name"); name != nil {
			return cppDeclaratorName(name, source)
		}
	}
	if declarator := node.ChildByFieldName("declarator"); declarator != nil {
		if name := cppDeclaratorName(declarator, source); name != "" {
			return name
		}
	}
	for index := uint(0); index < node.NamedChildCount(); index++ {
		if name := cppDeclaratorName(node.NamedChild(index), source); name != "" {
			return name
		}
	}
	return ""
}

func firstDescendantOfKind(node *tree_sitter.Node, kind string) *tree_sitter.Node {
	if node == nil {
		return nil
	}
	if node.Kind() == kind {
		return node
	}
	for index := uint(0); index < node.NamedChildCount(); index++ {
		if found := firstDescendantOfKind(node.NamedChild(index), kind); found != nil {
			return found
		}
	}
	return nil
}

func sourceDeclarationIsPublic(
	language string,
	node *tree_sitter.Node,
	declaration SourceDeclaration,
	source []byte,
) bool {
	switch language {
	case "go":
		first, _ := utf8FirstRune(declaration.Name)
		if !unicode.IsUpper(first) {
			return false
		}
		if declaration.Kind != "method" {
			return true
		}
		receiver := node.ChildByFieldName("receiver")
		receiverType := firstDescendantOfKind(receiver, "type_identifier")
		if receiverType == nil {
			return false
		}
		receiverFirst, _ := utf8FirstRune(receiverType.Utf8Text(source))
		return unicode.IsUpper(receiverFirst)
	case "javascript", "typescript", "tsx":
		if sourceNodeHasAncestor(node, "export_statement") {
			return true
		}
		if declaration.Kind == "method" {
			text := node.Utf8Text(source)
			if strings.Contains(text, "private ") || strings.Contains(text, "protected ") ||
				strings.HasPrefix(declaration.Name, "#") {
				return false
			}
			for parent := node.Parent(); parent != nil; parent = parent.Parent() {
				if parent.Kind() == "class_declaration" ||
					parent.Kind() == "abstract_class_declaration" ||
					parent.Kind() == "interface_declaration" ||
					parent.Kind() == "type_alias_declaration" {
					return sourceNodeHasAncestor(parent, "export_statement")
				}
			}
		}
		return false
	case "cpp":
		return cppDeclarationIsPublic(node, source)
	case "java":
		if sourceNodeHasModifier(node, source, "public") {
			return true
		}
		if declaration.Kind == "method" {
			for parent := node.Parent(); parent != nil; parent = parent.Parent() {
				if parent.Kind() == "interface_declaration" ||
					parent.Kind() == "annotation_type_declaration" {
					return !sourceNodeHasModifier(node, source, "private") &&
						!sourceNodeHasModifier(node, source, "protected")
				}
			}
		}
		return false
	case "python":
		if strings.HasPrefix(declaration.Name, "_") {
			return false
		}
		if sourceNodeIsTopLevel(node, "module") {
			return true
		}
		for parent := node.Parent(); parent != nil; parent = parent.Parent() {
			if parent.Kind() == "class_definition" {
				return !strings.HasPrefix(sourceNodeFieldText(parent, "name", source), "_")
			}
		}
	}
	return false
}

func sourceDeclarationTestKind(
	language string,
	path string,
	node *tree_sitter.Node,
	declaration SourceDeclaration,
	source []byte,
) string {
	isTest := false
	switch language {
	case "go":
		isTest = strings.HasSuffix(strings.ToLower(path), "_test.go") &&
			declaration.Kind == "function" && goTestDeclarationName(declaration.Name)
	case "javascript", "typescript", "tsx":
		isTest = declaration.Kind == "test" ||
			sourcePathLooksLikeTest(path) &&
				(strings.HasPrefix(strings.ToLower(declaration.Name), "test") ||
					strings.HasSuffix(strings.ToLower(declaration.Name), "test"))
	case "cpp":
		isTest = sourcePathLooksLikeTest(path) &&
			(strings.HasPrefix(strings.ToLower(declaration.Name), "test") ||
				strings.HasSuffix(strings.ToLower(declaration.Name), "test"))
	case "java":
		isTest = declaration.Kind == "method" &&
			(javaMethodHasTestAnnotation(node, source) ||
				sourcePathLooksLikeTest(path) && strings.HasPrefix(strings.ToLower(declaration.Name), "test"))
	case "python":
		isTest = sourcePathLooksLikeTest(path) &&
			declaration.Kind == "function" && strings.HasPrefix(declaration.Name, "test_")
	}
	if !isTest {
		return ""
	}
	if sourcePathLooksLikeContractTest(path) {
		return "contract"
	}
	return "test"
}

func goTestDeclarationName(name string) bool {
	if name == "Example" {
		return true
	}
	for _, prefix := range []string{"Test", "Benchmark", "Fuzz", "Example"} {
		if !strings.HasPrefix(name, prefix) || len(name) == len(prefix) {
			continue
		}
		first, _ := utf8FirstRune(name[len(prefix):])
		if !unicode.IsLower(first) {
			return true
		}
	}
	return false
}

func javaMethodHasTestAnnotation(node *tree_sitter.Node, source []byte) bool {
	recognized := map[string]bool{
		"Test":              true,
		"ParameterizedTest": true,
		"RepeatedTest":      true,
		"TestFactory":       true,
		"TestTemplate":      true,
	}
	found := false
	walkSourceTree(sourceNodeModifiers(node), func(child *tree_sitter.Node) {
		if found || child == nil ||
			child.Kind() != "marker_annotation" && child.Kind() != "annotation" {
			return
		}
		name := sourceNodeFieldText(child, "name", source)
		if qualified := strings.LastIndex(name, "."); qualified >= 0 {
			name = name[qualified+1:]
		}
		found = recognized[name]
	})
	return found
}

func sourcePathLooksLikeTest(path string) bool {
	name := strings.ToLower(filepath.Base(path))
	return strings.Contains(name, "_test.") ||
		strings.Contains(name, ".test.") ||
		strings.Contains(name, ".spec.") ||
		strings.HasPrefix(name, "test_") ||
		strings.HasSuffix(strings.TrimSuffix(name, filepath.Ext(name)), "test")
}

func sourcePathLooksLikeContractTest(path string) bool {
	name := strings.ToLower(filepath.Base(path))
	return strings.Contains(name, "contract_test") ||
		strings.Contains(name, "contract.test") ||
		strings.Contains(name, "contract.spec") ||
		strings.Contains(name, "contracttest")
}

func sourceNodeHasAncestor(node *tree_sitter.Node, kind string) bool {
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		if parent.Kind() == kind {
			return true
		}
	}
	return false
}

func sourceNodeIsTopLevel(node *tree_sitter.Node, rootKind string) bool {
	parent := node.Parent()
	if parent != nil && parent.Kind() == "export_statement" {
		parent = parent.Parent()
	}
	return parent != nil && parent.Kind() == rootKind
}

func sourceNodeHasModifier(node *tree_sitter.Node, source []byte, modifier string) bool {
	modifiers := sourceNodeModifiers(node)
	if modifiers == nil {
		return false
	}
	for _, field := range strings.Fields(modifiers.Utf8Text(source)) {
		if field == modifier {
			return true
		}
	}
	return false
}

func sourceNodeModifiers(node *tree_sitter.Node) *tree_sitter.Node {
	if node == nil {
		return nil
	}
	modifiers := node.ChildByFieldName("modifiers")
	if modifiers == nil {
		for index := uint(0); index < node.NamedChildCount(); index++ {
			child := node.NamedChild(index)
			if child != nil && child.Kind() == "modifiers" {
				modifiers = child
				break
			}
		}
	}
	return modifiers
}

func cppDeclarationIsPublic(node *tree_sitter.Node, source []byte) bool {
	scopeNode := node
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		switch parent.Kind() {
		case "field_declaration_list":
			owner := parent.Parent()
			public := owner != nil && owner.Kind() == "struct_specifier"
			for index := uint(0); index < parent.NamedChildCount(); index++ {
				child := parent.NamedChild(index)
				if child == nil || child.StartByte() >= node.StartByte() {
					break
				}
				if child.Kind() == "access_specifier" {
					public = strings.HasPrefix(strings.TrimSpace(child.Utf8Text(source)), "public")
				}
			}
			return public
		case "namespace_definition":
			if sourceNodeFieldText(parent, "name", source) == "" {
				return false
			}
			return !cppHasInternalLinkage(scopeNode.Utf8Text(source))
		case "translation_unit":
			return !cppHasInternalLinkage(scopeNode.Utf8Text(source))
		case "declaration", "template_declaration", "type_definition":
			scopeNode = parent
		}
	}
	return false
}

func cppHasInternalLinkage(text string) bool {
	normalized := " " + strings.Join(strings.Fields(text), " ") + " "
	return strings.Contains(normalized, " static ")
}

func sourceAnnotationsFromComment(
	comment sourceComment,
	markers []sourceAnnotationMarker,
) []SourceAnnotation {
	lines := strings.Split(comment.raw, "\n")
	annotations := make([]SourceAnnotation, 0)
	for index, rawLine := range lines {
		line := normalizeSourceCommentLine(rawLine)
		if !strings.HasPrefix(line, "@") {
			continue
		}
		kind, prefix, rest, separator, found := parseSourceAnnotationLine(line, markers)
		if !found {
			continue
		}
		annotations = append(annotations, SourceAnnotation{
			Kind:         kind,
			Prefix:       prefix,
			Refs:         splitSourceAnnotationRefs(rest),
			Line:         comment.line + index,
			Raw:          strings.TrimSpace(line),
			HasSeparator: separator,
			commentEnd:   comment.end,
		})
	}
	return annotations
}

func normalizeSourceCommentLine(line string) string {
	line = strings.TrimSpace(line)
	for _, prefix := range []string{"//", "/*", "*/", "*", "#"} {
		if strings.HasPrefix(line, prefix) {
			line = strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}
	line = strings.TrimSuffix(line, "*/")
	return strings.TrimSpace(line)
}

func parseSourceAnnotationLine(
	line string,
	markers []sourceAnnotationMarker,
) (kind, prefix, rest string, separator, found bool) {
	lower := strings.ToLower(line)
	for _, marker := range markers {
		prefix := strings.ToLower(marker.prefix)
		if !strings.HasPrefix(lower, prefix) {
			continue
		}
		remainder := line[len(marker.prefix):]
		if remainder == "" {
			return marker.kind, marker.prefix, "", false, true
		}
		if remainder[0] == ' ' || remainder[0] == '\t' {
			return marker.kind, marker.prefix, strings.TrimSpace(remainder), true, true
		}
		expectedPrefix := "SPEC-"
		if marker.kind == "test" || marker.kind == "test-contract" {
			expectedPrefix = "TEST-"
		}
		if !strings.HasPrefix(strings.ToUpper(remainder), expectedPrefix) {
			continue
		}
		return marker.kind, marker.prefix, strings.TrimSpace(remainder), false, true
	}
	return "", "", "", false, false
}

func splitSourceAnnotationRefs(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	refs := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			refs = append(refs, part)
		}
	}
	return refs
}

func bindSourceAnnotations(analysis *SourceAnalysis, comments []sourceComment, source []byte) {
	for annotationIndex := range analysis.Annotations {
		annotation := &analysis.Annotations[annotationIndex]
		if annotation.Ignored {
			continue
		}
		for declarationIndex := range analysis.Declarations {
			declaration := &analysis.Declarations[declarationIndex]
			if declaration.Ignored || declaration.bindingStart < annotation.commentEnd {
				continue
			}
			if !sourceGapIsAttachable(annotation.commentEnd, declaration.bindingStart, comments, source) {
				continue
			}
			annotation.Attached = true
			annotation.Declaration = declaration.Name
			annotation.DeclarationLine = declaration.Line
			declaration.Annotations = append(declaration.Annotations, *annotation)
			break
		}
	}
}

func sourceGapIsAttachable(start, end uint, comments []sourceComment, source []byte) bool {
	if start > end || int(end) > len(source) {
		return false
	}
	cursor := start
	lastCommentEnd := start
	for _, comment := range comments {
		if comment.start < start || comment.end > end {
			continue
		}
		if !sourceSliceIsWhitespace(source, cursor, comment.start) {
			return false
		}
		cursor = comment.end
		lastCommentEnd = comment.end
	}
	if !sourceSliceIsWhitespace(source, cursor, end) {
		return false
	}
	return strings.Count(string(source[lastCommentEnd:end]), "\n") <= 2
}

func sourceSliceIsWhitespace(source []byte, start, end uint) bool {
	if start > end || int(end) > len(source) {
		return false
	}
	return strings.TrimSpace(string(source[start:end])) == ""
}

func sourceIgnoreLines(comments []sourceComment, lineCount int) [][2]int {
	var ranges [][2]int
	activeStart := 0
	for _, comment := range comments {
		for offset, rawLine := range strings.Split(comment.raw, "\n") {
			line := comment.line + offset
			directive := strings.ToLower(normalizeSourceCommentLine(rawLine))
			switch directive {
			case "idd:ignore start", "idd:ignore-start":
				if activeStart == 0 {
					activeStart = line
				}
			case "idd:ignore end", "idd:ignore-end":
				if activeStart > 0 {
					ranges = append(ranges, [2]int{activeStart, line})
					activeStart = 0
				}
			case "idd:ignore":
				ranges = append(ranges, [2]int{line, line})
			}
		}
	}
	if activeStart > 0 {
		ranges = append(ranges, [2]int{activeStart, lineCount})
	}
	return ranges
}

func lineIsIgnored(ranges [][2]int, line int) bool {
	for _, lineRange := range ranges {
		if line >= lineRange[0] && line <= lineRange[1] {
			return true
		}
	}
	return false
}

func trimSourceName(name string) string {
	name = strings.TrimSpace(name)
	if len(name) >= 2 {
		switch name[0] {
		case '\'', '"', '`':
			if name[len(name)-1] == name[0] {
				name = name[1 : len(name)-1]
			}
		}
	}
	return strings.TrimSpace(name)
}

func utf8FirstRune(value string) (rune, int) {
	for _, char := range value {
		return char, len(string(char))
	}
	return 0, 0
}
