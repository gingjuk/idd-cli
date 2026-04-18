# IDD Link Validator Documentation

> 本项目遵循 IDD (Intent-Driven Development) 框架进行文档管理。

## 项目结构

```
idd-link-validator/
├── cmd/validator/         # CLI 入口
├── internal/
│   ├── collector/         # 文档/代码标识符收集
│   ├── config/           # 配置加载
│   ├── engine/          # 验证引擎
│   ├── graph/           # 链接图谱
│   ├── model/           # 数据模型
│   ├── reporter/        # 报告生成
│   └── auth/             # 示例代码
├── pkg/
│   ├── pattern/          # IDD 标识符正则模式
│   ├── walk/            # 文件遍历
│   └── annotation/       # 代码注解解析
├── examples/            # IDD 标识符使用示例
├── docs/                # 本项目 IDD 文档
├── skills/              # IDD skill 定义
└── .github/workflows/    # CI 配置
```

## IDD 标识符格式

**格式:** `<TYPE>-<MODULE>-<NUMBER>`

| 前缀 | 含义 | 示例 |
|------|------|------|
| `SPEC-` | 功能规格说明 | `SPEC-BE-001` |
| `CONTRACT-` | 接口/行为契约 | `CONTRACT-BE-001` |
| `TEST-` | 测试用例 | `TEST-BE-001` |
| `DESIGN-` | 架构设计决策 | `DESIGN-BE-001` |

**正则模式:**
```yaml
identifier_patterns:
  spec: "SPEC-[A-Z]+-[0-9]+"
  contract: "CONTRACT-[A-Z]+-[0-9]+"
  test: "TEST-[A-Z]+-[0-9]+"
  design: "DESIGN-[A-Z]+-[0-9]+"
```

## 文档结构

文档按模块组织存储在 `docs/` 目录下：

```
docs/
├── backend/
│   ├── spec.md       # SPEC-BE-001, SPEC-BE-002
│   ├── testing.md    # TEST-BE-001
│   ├── contract.md   # CONTRACT-BE-001
│   └── design.md     # DESIGN-BE-001
└── IDD.md            # 本文档
```

每份 IDD 文档必须包含 YAML frontmatter，用于工具快速解析：

```yaml
---
markers:
  - id: SPEC-BE-001
    name: 功能名称描述
  - id: CONTRACT-BE-001
    name: 契约接口描述
---
```

### 标识符索引表格格式

```markdown
| ID | Title | Status | Tests |
|----|-------|--------|-------|
| [SPEC-BE-001](#spec-be-001) | 功能名称 | Done | TEST-BE-001 |
```

## 代码注解格式

```go
// @spec SPEC-BE-001: 功能描述
// @contract CONTRACT-BE-001: 接口规范
// @test TEST-BE-001: 测试验证
// @design DESIGN-BE-001: 架构设计
func DoSomething() {
    // 实现
}
```

## 验证规则

1. **Completeness** — 每个 SPEC 必须有至少一个 TEST 链接
2. **Bidirectional** — SPEC→TEST 存在时，TEST→SPEC 也必须存在
3. **Orphan Detection** — 不能有孤立（无任何连接）的标识符
4. **Consistency** — 跨引用链一致性

## 配置示例

```yaml
version: "1.0"

docs:
  patterns:
    - "docs/**/*.md"
  identifier_patterns:
    spec: "SPEC-[A-Z]+-[0-9]+"
    contract: "CONTRACT-[A-Z]+-[0-9]+"
    test: "TEST-[A-Z]+-[0-9]+"
    design: "DESIGN-[A-Z]+-[0-9]+"

code:
  patterns:
    - "**/*.go"
  annotations:
    - "@spec"
    - "@contract"
    - "@test"
    - "@design"

validation:
  require_bidirectional: true
  allow_orphans: false
  require_spec_test_coverage: true

output:
  file: "idd-report.json"
  include_graph: true
  verbose: false
```

## 使用方法

```bash
# 构建
go build -o idd-verify ./cmd/validator

# 运行验证
./idd-verify run --config idd.yaml

# 使用 Make
make build && make run
```
