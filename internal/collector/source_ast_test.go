// Package collector tests syntax-tree-backed source binding.
package collector

import (
	"strings"
	"testing"
)

// @test-contract TEST-INTERNAL_COLLECTOR-024
func TestAnalyzeSource_BindsAnnotationsAcrossSupportedLanguages(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		source      string
		wantLang    string
		wantName    string
		wantKind    string
		wantPublic  bool
		wantLine    int
		wantAnnLine int
	}{
		{
			name:        "go function",
			path:        "service.go",
			source:      "package service\n// @implement SPEC-SERVICE-001\nfunc Run() {}\n",
			wantLang:    "go",
			wantName:    "Run",
			wantKind:    "function",
			wantPublic:  true,
			wantLine:    3,
			wantAnnLine: 2,
		},
		{
			name:        "typescript interface",
			path:        "service.ts",
			source:      "// @implement SPEC-SERVICE-001\nexport interface Service { run(): void }\n",
			wantLang:    "typescript",
			wantName:    "Service",
			wantKind:    "interface",
			wantPublic:  true,
			wantLine:    2,
			wantAnnLine: 1,
		},
		{
			name:        "tsx function",
			path:        "view.tsx",
			source:      "// @implement SPEC-VIEW-001\nexport function View() { return <div/> }\n",
			wantLang:    "tsx",
			wantName:    "View",
			wantKind:    "function",
			wantPublic:  true,
			wantLine:    2,
			wantAnnLine: 1,
		},
		{
			name:        "javascript exported arrow",
			path:        "service.js",
			source:      "// @implement SPEC-SERVICE-001\nexport const run = () => 1\n",
			wantLang:    "javascript",
			wantName:    "run",
			wantKind:    "function",
			wantPublic:  true,
			wantLine:    2,
			wantAnnLine: 1,
		},
		{
			name:        "cpp function",
			path:        "service.cpp",
			source:      "// @implement SPEC-SERVICE-001\nint run() { return 0; }\n",
			wantLang:    "cpp",
			wantName:    "run",
			wantKind:    "function",
			wantPublic:  true,
			wantLine:    2,
			wantAnnLine: 1,
		},
		{
			name:        "java class",
			path:        "Service.java",
			source:      "// @implement SPEC-SERVICE-001\npublic class Service {}\n",
			wantLang:    "java",
			wantName:    "Service",
			wantKind:    "class",
			wantPublic:  true,
			wantLine:    2,
			wantAnnLine: 1,
		},
		{
			name:        "python function",
			path:        "service.py",
			source:      "# @implement SPEC-SERVICE-001\ndef run():\n    pass\n",
			wantLang:    "python",
			wantName:    "run",
			wantKind:    "function",
			wantPublic:  true,
			wantLine:    2,
			wantAnnLine: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			analysis, err := AnalyzeSource(test.path, []byte(test.source))
			if err != nil {
				t.Fatalf("AnalyzeSource() error = %v", err)
			}
			if analysis.Language != test.wantLang {
				t.Errorf("Language = %q, want %q", analysis.Language, test.wantLang)
			}
			if len(analysis.ParseErrors) != 0 {
				t.Fatalf("ParseErrors = %#v, want none", analysis.ParseErrors)
			}
			if len(analysis.Annotations) != 1 {
				t.Fatalf("Annotations = %#v, want one", analysis.Annotations)
			}
			annotation := analysis.Annotations[0]
			if !annotation.Attached || annotation.Declaration != test.wantName ||
				annotation.Kind != "implement" ||
				annotation.Line != test.wantAnnLine ||
				annotation.DeclarationLine != test.wantLine ||
				strings.Join(annotation.Refs, ",") != "SPEC-SERVICE-001" &&
					strings.Join(annotation.Refs, ",") != "SPEC-VIEW-001" {
				t.Errorf("Annotation = %#v", annotation)
			}

			var declaration *SourceDeclaration
			for index := range analysis.Declarations {
				if analysis.Declarations[index].Name == test.wantName {
					declaration = &analysis.Declarations[index]
					break
				}
			}
			if declaration == nil {
				t.Fatalf("Declarations = %#v, want %s", analysis.Declarations, test.wantName)
			}
			if declaration.Kind != test.wantKind ||
				declaration.Public != test.wantPublic ||
				declaration.Line != test.wantLine {
				t.Errorf("Declaration = %#v", declaration)
			}
		})
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-024
func TestAnalyzeSource_RejectsStringAndBodyBindings(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		source string
	}{
		{
			name: "go",
			path: "service.go",
			source: `package service
var fixture = "// @implement SPEC-DEMO-999"
func run() {
	// @implement SPEC-DEMO-002
}
`,
		},
		{
			name: "typescript",
			path: "service.ts",
			source: `const fixture = "// @implement SPEC-DEMO-999";
export function run() {
  // @implement SPEC-DEMO-002
}
`,
		},
		{
			name: "tsx",
			path: "view.tsx",
			source: `const fixture = "// @implement SPEC-DEMO-999";
export function View() {
  // @implement SPEC-DEMO-002
  return <div/>;
}
`,
		},
		{
			name: "javascript",
			path: "service.js",
			source: `const fixture = "// @implement SPEC-DEMO-999";
export function run() {
  // @implement SPEC-DEMO-002
}
`,
		},
		{
			name: "cpp",
			path: "service.cpp",
			source: `const char* fixture = "// @implement SPEC-DEMO-999";
int run() {
  // @implement SPEC-DEMO-002
  return 1;
}
`,
		},
		{
			name: "java",
			path: "Service.java",
			source: `class Service {
  String fixture = "// @implement SPEC-DEMO-999";
  void run() {
    // @implement SPEC-DEMO-002
  }
}
`,
		},
		{
			name: "python",
			path: "service.py",
			source: `fixture = "# @implement SPEC-DEMO-999"
def run():
    # @implement SPEC-DEMO-002
    return 1
`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			analysis, err := AnalyzeSource(test.path, []byte(test.source))
			if err != nil {
				t.Fatalf("AnalyzeSource() error = %v", err)
			}
			if len(analysis.Annotations) != 1 {
				t.Fatalf("Annotations = %#v, want only the real body comment", analysis.Annotations)
			}
			if analysis.Annotations[0].Attached {
				t.Errorf("body annotation attached to declaration: %#v", analysis.Annotations[0])
			}
			if strings.Contains(analysis.Annotations[0].Raw, "999") {
				t.Errorf("annotation-looking string was collected: %#v", analysis.Annotations)
			}
		})
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-024
func TestAnalyzeSource_ClassifiesTestsAcrossLanguages(t *testing.T) {
	tests := []struct {
		path   string
		source string
		name   string
	}{
		{"service_test.go", "package service\n// @test TEST-SERVICE-001\nfunc TestRun(t *testing.T) {}\n", "TestRun"},
		{"benchmark_test.go", "package service\n// @test TEST-SERVICE-001\nfunc BenchmarkRun(b *testing.B) {}\n", "BenchmarkRun"},
		{"fuzz_test.go", "package service\n// @test TEST-SERVICE-001\nfunc FuzzRun(f *testing.F) {}\n", "FuzzRun"},
		{"example_test.go", "package service\n// @test TEST-SERVICE-001\nfunc ExampleRun() {}\n", "ExampleRun"},
		{"example_package_test.go", "package service\n// @test TEST-SERVICE-001\nfunc Example() {}\n", "Example"},
		{"service.test.ts", "// @test TEST-SERVICE-001\ntest('runs', () => {})\n", "runs"},
		{"service.only.test.ts", "// @test TEST-SERVICE-001\ntest.only('runs', () => {})\n", "runs"},
		{"service.each.test.ts", "// @test TEST-SERVICE-001\ntest.each([[1]])('runs', () => {})\n", "runs"},
		{
			"nested.test.ts",
			"describe('service', () => {\n// @test TEST-SERVICE-001\ntest('runs', () => {})\n})\n",
			"runs",
		},
		{"ServiceTest.java", "class ServiceTest { // @test TEST-SERVICE-001\n@Test public void runs() {} }\n", "runs"},
		{"ServiceParameterizedTest.java", "class ServiceTest { // @test TEST-SERVICE-001\n@ParameterizedTest public void runs() {} }\n", "runs"},
		{"test_service.py", "# @test TEST-SERVICE-001\ndef test_run():\n    pass\n", "test_run"},
		{"service_test.cpp", "// @test TEST-SERVICE-001\nvoid TestRun() {}\n", "TestRun"},
		{"service_macro_test.cpp", "// @test TEST-SERVICE-001\nTEST(Service, Runs) {}\n", "TEST"},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			analysis, err := AnalyzeSource(test.path, []byte(test.source))
			if err != nil {
				t.Fatalf("AnalyzeSource() error = %v", err)
			}
			for _, declaration := range analysis.Declarations {
				if declaration.Name == test.name {
					if declaration.TestKind != "test" {
						t.Errorf("TestKind = %q, want test: %#v", declaration.TestKind, declaration)
					}
					return
				}
			}
			t.Fatalf("Declarations = %#v, want test %s", analysis.Declarations, test.name)
		})
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-024
func TestAnalyzeSource_ReportsSyntaxErrors(t *testing.T) {
	analysis, err := AnalyzeSource("broken.py", []byte("def broken(:\n"))
	if err != nil {
		t.Fatalf("AnalyzeSource() error = %v", err)
	}
	if len(analysis.ParseErrors) == 0 {
		t.Fatalf("ParseErrors = %#v, want syntax error", analysis.ParseErrors)
	}
	if len(analysis.Annotations) != 0 {
		t.Errorf("Annotations = %#v, want none from untrustworthy tree", analysis.Annotations)
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-024
func TestSupportedSourcePath(t *testing.T) {
	for _, path := range []string{
		"a.go", "a.ts", "a.tsx", "a.js", "a.jsx", "a.cpp", "a.cc", "a.cxx",
		"a.hpp", "a.hh", "a.hxx", "a.java", "a.py",
	} {
		if !SupportedSourcePath(path) {
			t.Errorf("SupportedSourcePath(%q) = false", path)
		}
	}
	for _, path := range []string{"a.md", "a.c", "a.rs", "a.json"} {
		if SupportedSourcePath(path) {
			t.Errorf("SupportedSourcePath(%q) = true", path)
		}
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-024
func TestAnalyzeSource_DoesNotTreatLongerWordsAsAnnotations(t *testing.T) {
	analysis, err := AnalyzeSource(
		"service.ts",
		[]byte("// @testing is ordinary prose\nexport function run() {}\n"),
	)
	if err != nil {
		t.Fatalf("AnalyzeSource() error = %v", err)
	}
	if len(analysis.Annotations) != 0 {
		t.Errorf("Annotations = %#v, want none", analysis.Annotations)
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-024
func TestAnalyzeSource_IgnoreDirectiveMustBeStandalone(t *testing.T) {
	source := `package service

// Run documents how an idd:ignore block is handled.
// @implement SPEC-SERVICE-001
func Run() {}

// idd:ignore start
// @implement SPEC-SERVICE-002
func Ignored() {}
// idd:ignore end
`
	analysis, err := AnalyzeSource("service.go", []byte(source))
	if err != nil {
		t.Fatalf("AnalyzeSource() error = %v", err)
	}
	if len(analysis.Annotations) != 2 {
		t.Fatalf("Annotations = %#v, want two real comments", analysis.Annotations)
	}
	if analysis.Annotations[0].Ignored || !analysis.Annotations[0].Attached {
		t.Errorf("prose mention activated ignore behavior: %#v", analysis.Annotations[0])
	}
	if !analysis.Annotations[1].Ignored {
		t.Errorf("standalone ignore directive was not honored: %#v", analysis.Annotations[1])
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-024
func TestAnalyzeSource_GoMethodVisibilityIncludesReceiverType(t *testing.T) {
	source := `package service
type publicService struct{}
type Service struct{}
func (publicService) Run() {}
func (*Service) Run() {}
`
	analysis, err := AnalyzeSource("service.go", []byte(source))
	if err != nil {
		t.Fatalf("AnalyzeSource() error = %v", err)
	}
	var visibility []bool
	for _, declaration := range analysis.Declarations {
		if declaration.Kind == "method" && declaration.Name == "Run" {
			visibility = append(visibility, declaration.Public)
		}
	}
	if len(visibility) != 2 || visibility[0] || !visibility[1] {
		t.Errorf("method visibility = %#v, want [false true]", visibility)
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-024
func TestAnalyzeSource_UsesConfiguredAnnotationPrefixes(t *testing.T) {
	source := []byte(
		"// @fulfills SPEC-SERVICE-001\n" +
			"export function run() {}\n" +
			"// @implement SPEC-SERVICE-002\n" +
			"export function canonicalOnly() {}\n",
	)
	analysis, err := AnalyzeSourceWithAnnotations(
		"service.ts",
		source,
		map[string]string{
			"spec":          "@fulfills",
			"test":          "@checks",
			"test_contract": "@checks-contract",
		},
	)
	if err != nil {
		t.Fatalf("AnalyzeSourceWithAnnotations() error = %v", err)
	}
	if len(analysis.Annotations) != 1 {
		t.Fatalf("Annotations = %#v, want only configured prefix", analysis.Annotations)
	}
	if got := analysis.Annotations[0]; got.Kind != "implement" ||
		got.Prefix != "@fulfills" ||
		!got.Attached || got.Declaration != "run" ||
		strings.Join(got.Refs, ",") != "SPEC-SERVICE-001" {
		t.Errorf("configured annotation = %#v", got)
	}
}

// @test-contract TEST-INTERNAL_COLLECTOR-024
func TestAnalyzeSource_PublicLanguageEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		source string
		target string
	}{
		{
			name:   "C++ top-level class",
			path:   "service.hpp",
			source: "// @implement SPEC-SERVICE-001\nclass Service {};\n",
			target: "Service",
		},
		{
			name: "Go grouped type",
			path: "service.go",
			source: "package service\n" +
				"type (\n" +
				"// @implement SPEC-SERVICE-001\n" +
				"Service struct{}\n" +
				")\n",
			target: "Service",
		},
		{
			name: "TypeScript interface method",
			path: "service.ts",
			source: "export interface Service {\n" +
				"// @implement SPEC-SERVICE-001\n" +
				"run(): void\n" +
				"}\n",
			target: "run",
		},
		{
			name: "C++ template alias",
			path: "service.hpp",
			source: "// @implement SPEC-SERVICE-001\n" +
				"template <typename T>\n" +
				"using Service = T;\n",
			target: "Service",
		},
		{
			name: "Java interface method",
			path: "Service.java",
			source: "public interface Service {\n" +
				"// @implement SPEC-SERVICE-001\n" +
				"void run();\n}\n",
			target: "run",
		},
		{
			name: "Python decorated method",
			path: "service.py",
			source: "class Service:\n" +
				"    # @implement SPEC-SERVICE-001\n" +
				"    @classmethod\n" +
				"    def run(cls):\n" +
				"        pass\n",
			target: "run",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			analysis, err := AnalyzeSource(test.path, []byte(test.source))
			if err != nil {
				t.Fatalf("AnalyzeSource() error = %v", err)
			}
			for _, declaration := range analysis.Declarations {
				if declaration.Name != test.target {
					continue
				}
				if !declaration.Public {
					t.Errorf("declaration = %#v, want public", declaration)
				}
				if len(declaration.Annotations) != 1 ||
					!declaration.Annotations[0].Attached {
					t.Errorf(
						"declaration annotation = %#v, want one attached annotation; analysis = %#v",
						declaration.Annotations,
						analysis,
					)
				}
				return
			}
			t.Fatalf("Declarations = %#v, want %s", analysis.Declarations, test.target)
		})
	}
}
