---
markers:
  - id: CONTRACT-COLLECTOR-001
    name: Collector Interface Contracts
---

# Contracts (collector)

## CONTRACT-COLLECTOR-001: Collector Interface Contracts

**Status:** Done

**Overview:**

Collectors gather IDD identifiers from documentation and source code. This contract defines the interfaces and shared behaviors for all collectors.

### Shared Functionality

#### Frontmatter Parsing

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

#### Document Validation

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

**Related Specs:** `SPEC-COLLECTOR-001`, `SPEC-COLLECTOR-002`, `SPEC-BE-001`
