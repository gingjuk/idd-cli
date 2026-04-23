---
related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Contracts (collector)

**Status:** Done

**Overview:**

Collectors gather IDD identifiers from documentation and source code. This contract defines the interfaces and shared behaviors for all collectors.

## Shared Functionality

### Frontmatter Parsing

```go
type Marker struct {
    ID       string `yaml:"id"`
    Name     string `yaml:"name"`
    Describe string `yaml:"describe"`
}

type Frontmatter struct {
    Markers []Marker `yaml:"markers"`
}

func ParseFrontmatter(content string) (*Frontmatter, error)
func ValidateFrontmatterMarkers(fm *Frontmatter, content string, filePath string) []string
```

### Document Validation

```go
func ValidateDocumentStructure(filePath string, idType string) error
func ValidateModulePrefix(id string, filePath string) error
func GetExpectedFilename(idType string) string
```

### DocCollector Interface

```go
type DocCollector struct {
    cfg *config.Config
}

func NewDocCollector(cfg *config.Config) *DocCollector
func (c *DocCollector) Collect(ctx context.Context, targetPath string) (*model.IdentifierSet, []*model.ValidationError, error)
```

**Collection Process:**

1. **Target Resolution** — Accept file or directory path
2. **File Discovery** — Recursively find `.md` files
3. **Frontmatter Parse** — Extract markers from YAML frontmatter
4. **Content Scan** — Find IDD references in markdown body
5. **Validation** — Check marker consistency and module prefix
6. **Assembly** — Build IdentifierSet with links

### CodeCollector Interface

```go
type CodeCollector struct {
    cfg *config.Config
}

func NewCodeCollector(cfg *config.Config) *CodeCollector
func (c *CodeCollector) Collect(ctx context.Context, targetPath string) (*model.IdentifierSet, error)
```

**Collection Process:**

1. **Target Resolution** — Accept file or directory path
2. **File Discovery** — Find source files by extension
3. **Line Scan** — Search for annotation patterns
4. **Context Extraction** — Capture function name and comments
5. **Identifier Creation** — Convert annotations to identifiers
6. **Origin Set** — Mark as `OriginCode`

### Error Handling Rules

| Error Type | Behavior |
| ---------- | -------- |
| File not found | Return empty set, no error |
| Permission denied | Skip file, continue |
| Parse error | Log warning, skip file |
| Invalid frontmatter | Add ValidationError, continue |
| Module prefix mismatch | Add ValidationError, continue |
| Bare marker (no backticks) | Add ValidationError, continue |

### Interface: ExtractTitle

```go
func (c *DocCollector) ExtractTitle(content string, id string) string
```

Extracts document titles from markdown headings containing identifiers.

### Interface: ExtractAnnotations

```go
// Annotation patterns matched from code
var AnnotationPatterns = []*regexp.Regexp{...}
```

Extracts IDD annotations (@implement, @test, @test-contract) from code.

### Interface: ExtractFunctionContext

```go
func ExtractFunctionComment(lines []string, annotationLine int) string
```

Extracts function name and preceding comments as context for code annotations.

### Interface: SetOrigin

```go
func (id *Identifier) SetOrigin(origin model.OriginType)
```

Sets origin to OriginCode for code-based identifiers.

### Interface: SplitAnnotationRefs

```go
func SplitAnnotationRefs(s string) []string
```

Splits comma-separated IDD references from an annotation.

### Interface: ShouldIgnore

```go
func (c *DocCollector) ShouldIgnore(path string) bool
func (c *CodeCollector) ShouldIgnore(path string) bool
```

Checks if a path should be ignored based on configured ignore patterns.

### Interface: DiscoverFiles

File discovery is handled within the Collect method - it recursively finds source files with supported extensions (.go, .ts, .tsx, .js).

### Interface: GetLanguagePatterns

Language patterns are defined in the annotation patterns configuration.

---

**Related Specs:** `SPEC-INT_COL-001`, `SPEC-INT_COL-002`, `SPEC-INT_COL-006`, `SPEC-INT_COL-009`, `SPEC-INT_COL-010`, `SPEC-INT_COL-011`, `SPEC-INT_COL-012`, `SPEC-INT_COL-013`, `SPEC-INT_COL-014`, `SPEC-INT_COL-015`
