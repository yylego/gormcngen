[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yylego/gormcngen/release.yml?branch=main&label=BUILD)](https://github.com/yylego/gormcngen/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yylego/gormcngen)](https://pkg.go.dev/github.com/yylego/gormcngen)
[![Coverage Status](https://img.shields.io/coveralls/github/yylego/gormcngen/main.svg)](https://coveralls.io/github/yylego/gormcngen?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.24+-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/yylego/gormcngen.svg)](https://github.com/yylego/gormcngen/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yylego/gormcngen)](https://goreportcard.com/report/github.com/yylego/gormcngen)

# gormcngen

`gormcngen`: 赋予 GORM 模型使用 Columns() 获取列名的函数

就像 Java 生态系统中的 `MyBatis Plus`，它允许开发人员使用像 `Example::getName` 这样的表达式获取列名。

就像 Python 生态系统中的 `SQLAlchemy`，它允许开发人员使用像 `Example.name` 这样的 class 函数获得列名。

`gormcngen` 也为 Go 模型赋予 **类型安全** 的列引用功能。

---

## 生态系统

![GORM Type-Safe Ecosystem](https://github.com/yylego/gormcnm/raw/main/assets/gormcnm-ecosystem.svg)

---

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->

## 英文文档

[ENGLISH README](README.md)

<!-- TEMPLATE (ZH) END: LANGUAGE NAVIGATION -->

---

## 语言生态对比

| 语言       | ORM           | 类型安全列名       | 示例                                    |
| ---------- | ------------- | ------------------ | --------------------------------------- |
| **Java**   | MyBatis Plus  | `Example::getName` | `wrapper.eq(Example::getName, "alice")` |
| **Python** | SQLAlchemy    | `Example.name`     | `query.filter(Example.name == "alice")` |
| **Go**     | **GORMCNGEN** | `cls.Name.Eq()`    | `db.Where(cls.Name.Eq("alice"))`        |

---

## 安装

```bash
go get github.com/yylego/gormcngen
```

## 重要说明

**gormcngen 是一个 Go 包，不是 CLI 应用。** 它需要通过测试代码驱动的生成工作流程来使用。

## 快速开始

### 1. 创建项目结构

设置基本的项目结构，为模型和生成的代码创建专门的 DIR：

```bash
# 创建 models DIR
mkdir -p internal/models
```

### 2. 定义 GORM 模型

定义您的数据模型，gormcngen 将为这些模型生成列访问方法：

创建 `internal/models/models.go`：

```go
package models

type Account struct {
    ID       uint   `gorm:"primaryKey"`
    Username string `gorm:"uniqueIndex;size:100"`
    Mailbox  string `gorm:"index;size:255"`
    Age      int    `gorm:"column:age"`
    IsActive bool   `gorm:"default:true"`
}
```

### 3. 创建生成文件

创建用于存放生成代码的目标文件和包含生成逻辑的测试文件：

```bash
# 创建用于存放生成代码的目标文件，包含 package 声明
echo "package models" > internal/models/ngen.go

# 创建包含生成逻辑的测试文件，包含 package 声明
echo "package models" > internal/models/ngen_test.go
```

### 4. 编写生成逻辑

在测试文件中编写代码生成逻辑，配置生成选项并设置要处理的模型。

**注意**: 在 Go 中，使用测试文件生成源代码是一种常见的做法。

创建 `internal/models/ngen_test.go`：

```go
package models

import (
    "testing"

    "github.com/yylego/gormcngen"
    "github.com/yylego/osexistpath/osmustexist"
    "github.com/yylego/runpath/runtestpath"
)

//go:generate go test -run ^TestGenerate$ -tags=gormcngen_generate
func TestGenerate(t *testing.T) {
    // 获取目标文件的绝对路径（ngen.go）
    absPath := osmustexist.FILE(runtestpath.SrcPath(t))
    t.Log(absPath)

    // 配置生成选项
    options := gormcngen.NewOptions().
        WithColumnClassExportable(true).
        WithColumnsMethodRecvName("c").
        WithColumnsCheckFieldType(true)

    // 定义要处理的模型
    models := []interface{}{
		&Account{},
	}

    // 创建配置并执行生成
    cfg := gormcngen.NewConfigs(models, options, absPath)
    cfg.WithIsGenPreventEdit(true)  // 添加"请勿编辑"警告头部 (默认: true)
    cfg.WithGeneratedFromPos(gormcngen.GetGenPosFuncMark(0))  // 显示生成源位置 (默认: show)
    cfg.WithGoBuildDirective(gormcngen.DirectiveSkipGenerate)  // 添加构建指令，解决重新生成问题
    cfg.Gen()
}
```

### 5. 执行生成

运行测试来触发代码生成，生成的代码将自动写入到目标文件中：

```bash
# 整理依赖
go mod tidy

# 运行生成测试
cd internal/models
go test -v ./...
```

🎉 **生成完成！** 您的 `ngen.go` 文件现在包含了生成的列访问方法。

生成的 `ngen.go` 将包含：

```go
// Code generated using gormcngen. DO NOT EDIT.
// This file was auto generated via github.com/yylego/gormcngen

//go:build !gormcngen_generate

// Generated from: ngen_test.go:20 -> models.TestGenerate
// ========== GORMCNGEN:DO-NOT-EDIT-MARKER:END ==========

func (c *Account) Columns() *AccountColumns {
    return &AccountColumns{
        ID:       gormcnm.Cnm(c.ID, "id"),
        Username: gormcnm.Cnm(c.Username, "username"),
        Mailbox:  gormcnm.Cnm(c.Mailbox, "mailbox"),
        Age:      gormcnm.Cnm(c.Age, "age"),
        IsActive: gormcnm.Cnm(c.IsActive, "is_active"),
    }
}

type AccountColumns struct {
    gormcnm.ColumnOperationClass
    ID       gormcnm.ColumnName[uint]
    Username gormcnm.ColumnName[string]
    Mailbox  gormcnm.ColumnName[string]
    Age      gormcnm.ColumnName[int]
    IsActive gormcnm.ColumnName[bool]
}
```

🚀 **设置完成！** 您现在拥有了模型的类型安全列访问方法。

### 6. 在业务逻辑中使用

现在在编写业务代码时，您可以使用生成的类型安全列方法进行数据库查询：

```go
var account Account
cls := account.Columns()

// 完美的类型安全，零样板代码
err := db.Where(cls.Username.Eq("alice")).
         Where(cls.Age.Gte(18)).
         Where(cls.IsActive.Eq(true)).
         First(&account).Error
```

### 高级用法

```go
// 基础配置（匹配内部示例）
options := gormcngen.NewOptions().
    WithColumnClassExportable(true).           // 生成导出的 ExampleColumns 结构体
    WithEmbedColumnOperations(false)           // 不嵌入操作方法

// 中文字段名支持
chineseOptions := gormcngen.NewOptions().
    WithUseTagName(true).                      // 使用 cnm 标签值作为字段名
    WithTagKeyName("cnm").                     // 指定 'cnm' 作为标签键
    WithColumnClassExportable(true)

// 高级功能（来自 example6）
advancedOptions := gormcngen.NewOptions().
    WithColumnClassExportable(true).           // 导出结构体名称
    WithColumnsMethodRecvName("one").          // 自定义接收器名称
    WithColumnsCheckFieldType(true).           // 类型检查（推荐）
    WithIsGenFuncTableColumns(true)            // 生成 TableColumns 函数

// 批量处理多个模型
allModels := []interface{}{&Account{}, &Product{}, &Item{}, &Client{}}
configs := gormcngen.NewConfigs(allModels, options, "models_gen.go")
configs.WithIsGenPreventEdit(true)  // 为生成的文件添加"请勿编辑"头部
configs.WithGeneratedFromPos(gormcngen.GetGenPosFuncMark(0))  // 显示生成源位置 (默认: show)
configs.WithGoBuildDirective(gormcngen.DirectiveSkipGenerate)  // 添加构建指令（推荐）
configs.Gen()
```

## 构建指令（解决重新生成问题）

当你修改 GORM 模型（如删除或重命名字段）时，`ngen.go` 中的生成代码仍然引用旧字段，导致编译错误。这会造成"先有鸡还是先有蛋"的问题：你无法运行 `TestGenerate` 重新生成代码，因为项目无法编译通过。

### 解决方案

`WithGoBuildDirective(gormcngen.DirectiveSkipGenerate)` 会在生成文件中添加构建约束：

```go
//go:build !gormcngen_generate
```

这会告诉 Go：

- **标准编译** (`go build ./...`)：生成文件正常编译
- **重新生成** (`go test -tags=gormcngen_generate`)：生成文件被**跳过**，允许 `TestGenerate` 运行并重新生成代码

### 使用方法

**在测试文件中：**

```go
//go:generate go test -run ^TestGenerate$ -tags=gormcngen_generate
func TestGenerate(t *testing.T) {
    // ... 配置 ...
    cfg := gormcngen.NewConfigs(models, options, absPath)
    cfg.WithGoBuildDirective(gormcngen.DirectiveSkipGenerate)  // 添加这一行
    cfg.Gen()
}
```

**当模型变更导致编译问题时：**

```bash
# 跳过过时的生成文件并重新生成
go generate ./...
# 等价于：
go test -run ^TestGenerate$ -tags=gormcngen_generate
```

### 推荐

本项目所有示例都使用了 `WithGoBuildDirective`。虽然非强制要求，但我们**推荐**始终添加它，以避免模型变更时陷入困境。

---

## 高级功能

### 多语言字段支持

`cnm` 标签允许您定义中文别名作为字段名，这些别名将被生成为额外的结构体字段：

```go
type Product struct {
    ID          uint          `gorm:"primaryKey"`
    Name        string        `gorm:"size:255;not null" cnm:"V产品名称"`
    Price       decimal.Decimal `gorm:"type:decimal(10,2)"`
    CategoryID  uint          `gorm:"index"`
    CreatedAt   time.Time     `gorm:"autoCreateTime"`
    UpdatedAt   time.Time     `gorm:"autoUpdateTime"`
}
```

**生成结果：**

```go
type ProductColumns struct {
    gormcnm.ColumnOperationClass
    // 模型列的列名和类型
    ID        gormcnm.ColumnName[uint]
    Name      gormcnm.ColumnName[string]      // 映射到 "name"
    V产品名称   gormcnm.ColumnName[string]      // 中文别名，映射到 Name 字段
    Price     gormcnm.ColumnName[decimal.Decimal]
    CategoryID gormcnm.ColumnName[uint]
    CreatedAt gormcnm.ColumnName[time.Time]
    UpdatedAt gormcnm.ColumnName[time.Time]
}

func (*Product) Columns() *ProductColumns {
    return &ProductColumns{
        ID:        "id",
        Name:      "name",
        V产品名称:   "name",           // 中文别名，指向同一列
        Price:     "price",
        CategoryID: "category_id",
        CreatedAt: "created_at",
        UpdatedAt: "updated_at",
    }
}
```

**使用中文字段名进行查询：**

通过生成的中文字段别名，你可以使用母语编写查询语句：

```go
var product Product
var cls = product.Columns()

// 使用中文字段名查询 - 相同的数据库列，不同的 Go 字段名
if err := db.Where(cls.V产品名称.Eq("iPhone")).
    Where(cls.Price.Gte(5000.00)).
    First(&product).Error; err != nil {
    panic(errors.WithMessage(err, "未找到产品"))
}

fmt.Println("找到产品:", product.Name)
```

这让开发者可以用母语编写更易读的代码，同时保持完整的类型安全和数据库兼容性。

### Go Generate 集成

创建生成脚本：

**scripts/generate_columns.go:**

```go
package main

import (
    "github.com/yylego/gormcngen"
    "project-name/models"
)

func main() {
    models := []interface{}{&models.Account{}}
    options := gormcngen.NewOptions()
    configs := gormcngen.NewConfigs(models, options, "models/account_columns_gen.go")
    configs.WithGoBuildDirective(gormcngen.DirectiveSkipGenerate)
    configs.Gen()
}
```

然后在模型文件中使用：

```go
//go:generate go run scripts/generate_columns.go

type Account struct {
    ID       uint   `gorm:"primaryKey"`
    Username string `gorm:"uniqueIndex"`
    Mailbox  string `gorm:"index"`
}
```

## 🔗 配合 gormrepo 使用

将 **gormcngen** 与 **[gormrepo](https://github.com/yylego/gormrepo)** 配合使用，获得类型安全的 CRUD 操作。

### 快速预览

```go
// 创建 repo，传入列定义
repo := gormrepo.NewRepo(&Account{}, (&Account{}).Columns())

// gormrepo/gormclass 简洁写法
repo := gormrepo.NewRepo(gormclass.Use(&Account{}))

// 类型安全查询
account, err := repo.With(ctx, db).First(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Username.Eq("alice"))
})

// 条件查询
accounts, err := repo.With(ctx, db).Find(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Age.Gte(18)).Where(cls.Age.Lte(65))
})

// 类型安全更新
err := repo.With(ctx, db).Updates(
    func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
        return db.Where(cls.ID.Eq(1))
    },
    func(cls *AccountColumns) map[string]interface{} {
        return cls.Kw(cls.Age.Kv(26)).Kw(cls.Nickname.Kv("NewNick")).AsMap()
    },
)
```

👉 查看 **[gormrepo](https://github.com/yylego/gormrepo)** 获取完整文档和更多示例。

---

## 示例

查看 [examples](internal/examples) 和 [demos](internal/demos) 目录。

## 关联项目

探索完整的 GORM 生态系统集成包：

### 核心生态

- **[gormcnm](https://github.com/yylego/gormcnm)** - GORM 基础层，提供类型安全的列操作和条件构建
- **[gormcngen](https://github.com/yylego/gormcngen)** - 使用 AST 的代码生成，实现类型安全的 GORM 操作（本项目）
- **[gormrepo](https://github.com/yylego/gormrepo)** - 仓储模式实现，遵循 GORM 最佳实践
- **[gormmom](https://github.com/yylego/gormmom)** - 原生语言 GORM 标签生成引擎，支持智能列名
- **[gormzhcn](https://github.com/go-zwbc/gormzhcn)** - 完整的 GORM 中文编程接口

每个包针对 GORM 开发的不同方面，从本地化到类型安全和代码生成。

---

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-11-25 03:52:28.131064 +0000 UTC -->

## 📄 许可证类型

MIT 许可证 - 详见 [LICENSE](LICENSE)。

---

## 💬 联系与反馈

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **问题报告？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **新颖思路？** 创建 issue 讨论
- 📖 **文档疑惑？** 报告问题，帮助我们完善文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，协助解决性能问题
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **反馈意见？** 欢迎提出建议和意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/repo-name.git`）
3. **导航**：进入克隆的项目（`cd repo-name`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：面向用户的更改需要更新文档
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Merge Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Merge Request 和报告问题来贡献此项目。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**祝你用这个包编程愉快！** 🎉🎉🎉

<!-- TEMPLATE (ZH) END: STANDARD PROJECT FOOTER -->

---

## 📈 GitHub Stars

[![starring](https://starchart.cc/yylego/gormcngen.svg?variant=adaptive)](https://starchart.cc/yylego/gormcngen)
