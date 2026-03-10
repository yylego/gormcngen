// Package gormcngen provides constants used in GORM code generation
//
// gormcngen 提供 GORM 代码生成中使用的常量
package gormcngen

// DirectiveSkipGenerate is the recommended build directive value
// Use this with WithGoBuildDirective() to solve the regeneration problem
// When models change and generated code causes compilation issues, use:
//
//	go test -tags=gormcngen_generate
//
// DirectiveSkipGenerate 是推荐的构建指令值
// 与 WithGoBuildDirective() 配合使用，解决重新生成问题
// 当模型变更导致生成代码编译问题时，使用：
//
//	go test -tags=gormcngen_generate
const DirectiveSkipGenerate = "!gormcngen_generate"
