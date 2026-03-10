package models

import (
	"testing"

	"github.com/yylego/gormcngen"
	"github.com/yylego/osexistpath/osmustexist"
	"github.com/yylego/runpath/runtestpath"
)

// TestGenerate triggers column code generation using gormcngen
// Configures generation options with Chinese tag name support
//
// TestGenerate 使用 gormcngen 触发列代码生成
// 配置生成选项以支持中文标签名
//
//go:generate go test -run ^TestGenerate$ -tags=gormcngen_generate
func TestGenerate(t *testing.T) {
	absPath := osmustexist.FILE(runtestpath.SrcPath(t))
	t.Log(absPath)

	options := gormcngen.NewOptions().
		WithColumnClassExportable(true). // Exportable class names like StudentColumns // 可导出的类名如 StudentColumns
		WithColumnsMethodRecvName("c").  // Set method argument name // 设置方法接收者名称
		WithColumnsCheckFieldType(true). // Enable field type checking // 启用字段类型检查
		WithUseTagName(true)             // Enable cnm tag name usage // 启用 cnm 标签名使用

	cfg := gormcngen.NewConfigs([]interface{}{
		&Student{},
		&Class{},
	}, options, absPath)
	cfg.Gen()
}
