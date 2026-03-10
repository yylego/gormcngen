package models

import (
	"testing"

	"github.com/yylego/gormcngen"
	"github.com/yylego/osexistpath/osmustexist"
	"github.com/yylego/runpath/runtestpath"
)

// TestGenerate triggers column code generation using gormcngen
// Configures generation options with custom method argument name
//
// TestGenerate 使用 gormcngen 触发列代码生成
// 配置生成选项使用自定义方法接收者名称
//
//go:generate go test -run ^TestGenerate$ -tags=gormcngen_generate
func TestGenerate(t *testing.T) {
	absPath := osmustexist.FILE(runtestpath.SrcPath(t))
	t.Log(absPath)

	options := gormcngen.NewOptions().
		WithColumnClassExportable(true).  // Exportable class names like PersonColumns // 可导出的类名如 PersonColumns
		WithColumnsMethodRecvName("one"). // Custom method argument name // 自定义方法接收者名称
		WithColumnsCheckFieldType(true)   // Enable field type checking // 启用字段类型检查

	cfg := gormcngen.NewConfigs([]interface{}{&Person{}, &Example{}}, options, absPath)
	cfg.Gen()
}
