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

## 文档规范

### SPEC 规格文档

```markdown
## SPEC-BE-001: 功能名称

**需求描述:** 功能需求说明

**实现位置:** `internal/xxx/xxx.go`

**测试用例:** TEST-BE-001, TEST-BE-002
```

### CONTRACT 契约文档

```markdown
## CONTRACT-BE-001: 接口名称

**函数签名:** `func DoSomething(input Type) (output Type, error)`

**错误处理:**
- `ErrInvalidInput` - 输入参数无效
- `ErrNotFound` - 资源不存在
```

### TEST 测试文档

```markdown
## TEST-BE-001: 验证功能名称

**测试目的:** 验证功能符合 SPEC-BE-001 规范

**测试步骤:**
1. 准备测试数据
2. 执行操作
3. 验证结果
```

### DESIGN 设计文档

```markdown
## DESIGN-BE-001: 架构决策名称

**背景:** 为什么要做这个架构决策

**决策:** 最终采用的方案

**后果:** 这个决策的影响
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
