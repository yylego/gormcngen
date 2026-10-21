// Package gormcngen_test validates SchemaConfig with one GORM struct code generation
// Auto tests schema parsing, options configuration, and output generation
// Demonstrates various generation modes including tag handling and field output
//
// gormcngen_test 验证 SchemaConfig 的单模型代码生成
// 自动测试 schema 解析、选项配置和输出生成
// 演示包括标签处理和字段过滤在内的各种生成模式
package gormcngen_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yylego/gormcngen"
	"github.com/yylego/neatjson/neatjsons"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestGen(t *testing.T) {
	type Example struct {
		ID        int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
		Name      string    `gorm:"not null,type:text"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}

	config := gormcngen.NewSchemaConfig(&Example{}, gormcngen.NewOptions())
	output := config.Gen()
	t.Log(output.GetMethodCode())
	t.Log(output.GetStructCode())
	t.Log(neatjsons.S(output.GetPkgImports()))
}

func TestGen_UseTagName(t *testing.T) {
	type Example struct {
		ID        int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
		Name      string    `gorm:"not null,type:text" form:"name" json:"name" validate:"required" cnm:"V名称"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}

	options := gormcngen.NewOptions().
		WithUseTagName(true)

	config := gormcngen.NewSchemaConfig(&Example{}, options)
	output := config.Gen()
	t.Log(output.GetMethodCode())
	t.Log(output.GetStructCode())
	t.Log(neatjsons.S(output.GetPkgImports()))
}

func TestGen_ExcludeUntaggedFields(t *testing.T) {
	type Example struct {
		ID        int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
		Name      string    `gorm:"not null,type:text" form:"name" json:"name" validate:"required" cnm:"V名称"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}

	options := gormcngen.NewOptions().
		WithUseTagName(true).
		WithExcludeUntaggedFields(true)

	config := gormcngen.NewSchemaConfig(&Example{}, options)
	output := config.Gen()
	t.Log(output.GetMethodCode())
	t.Log(output.GetStructCode())
	t.Log(neatjsons.S(output.GetPkgImports()))
}

func TestGen_HasOneAssociation(t *testing.T) {
	type Profile struct {
		ID        uint
		Bio       string
		ExampleID uint // 外键
	}

	type Example struct {
		ID      uint
		Name    string
		Profile Profile `gorm:"foreignKey:ExampleID"`
	}

	options := gormcngen.NewOptions().
		WithUseTagName(true)

	config := gormcngen.NewSchemaConfig(&Example{}, options)
	output := config.Gen()
	t.Log(output.GetMethodCode())
	t.Log(output.GetStructCode())
	t.Log(neatjsons.S(output.GetPkgImports()))
}

func TestGen_ColumnsMethodRecvName(t *testing.T) {
	type Example struct {
		ID        int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
		Name      string    `gorm:"not null,type:text" form:"name" json:"name" validate:"required" cnm:"V名称"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}

	options := gormcngen.NewOptions().
		WithUseTagName(true).
		WithColumnsMethodRecvName("example")

	config := gormcngen.NewSchemaConfig(&Example{}, options)
	output := config.Gen()
	t.Log(output.GetMethodCode())
	t.Log(output.GetStructCode())
	t.Log(neatjsons.S(output.GetPkgImports()))
}

func TestGen_ColumnsCheckFieldType(t *testing.T) {
	type Example struct {
		ID        int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
		Name      string    `gorm:"not null,type:text" form:"name" json:"name" validate:"required" cnm:"V名称"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}

	options := gormcngen.NewOptions().
		WithUseTagName(true).
		WithColumnsMethodRecvName("example").
		WithColumnsCheckFieldType(true)

	config := gormcngen.NewSchemaConfig(&Example{}, options)
	output := config.Gen()
	t.Log(output.GetMethodCode())
	t.Log(output.GetStructCode())
	t.Log(neatjsons.S(output.GetPkgImports()))
}

type compositeTag struct{ Name string }

type compositeBox[T any] struct{ Value T }

type compositeRow struct {
	ID     uint
	Pinned *compositeTag              `gorm:"serializer:json"`
	Tags   []compositeTag             `gorm:"serializer:json"`
	Times  []time.Time                `gorm:"serializer:json"`
	Named  map[string]compositeTag    `gorm:"serializer:json"`
	Grid   [2]compositeTag            `gorm:"serializer:json"`
	Box    compositeBox[compositeTag] `gorm:"serializer:json"`
	Expire sql.Null[time.Time]
	Picked sql.Null[compositeTag]
}

// TestGen_CompositeFieldTypes checks that composite and generic field types come out as valid Go source in the struct's own package.
// Types of that package keep bare names, and each package the types touch lands in the imports.
//
// TestGen_CompositeFieldTypes 验证复合类型和泛型类型的字段，生成的是模型包里合法的 Go 源码。
// 模型包内的类型不加包名，类型用到的每个包都进入导入列表。
func TestGen_CompositeFieldTypes(t *testing.T) {
	options := gormcngen.NewOptions().
		WithColumnsMethodRecvName("c").
		WithColumnsCheckFieldType(true)

	output := gormcngen.NewSchemaConfig(&compositeRow{}, options).Gen()
	t.Log(output.GetStructCode())

	for _, want := range []string{
		"Pinned gormcnm.ColumnName[*compositeTag]",
		"Tags gormcnm.ColumnName[[]compositeTag]",
		"Times gormcnm.ColumnName[[]time.Time]",
		"Named gormcnm.ColumnName[map[string]compositeTag]",
		"Grid gormcnm.ColumnName[[2]compositeTag]",
		"Box gormcnm.ColumnName[compositeBox[compositeTag]]",
		"Expire gormcnm.ColumnName[sql.Null[time.Time]]",
		"Picked gormcnm.ColumnName[sql.Null[compositeTag]]",
	} {
		require.Contains(t, output.GetStructCode(), want)
	}
	require.Equal(t, map[string]bool{
		"github.com/yylego/gormcnm": true,
		"time":                      true,
		"database/sql":              true,
	}, output.GetPkgImports())
}

type embedContact struct {
	Name  string
	Phone string
}

type embedProfile struct {
	ID     uint
	PostID uint
}

type embedPost struct {
	gorm.Model
	Name    string
	Contact embedContact  `gorm:"embedded;embeddedPrefix:contact_"`
	Backup  *embedContact `gorm:"embedded;embeddedPrefix:backup_"`
	Profile embedProfile  `gorm:"foreignKey:PostID"`
}

// TestGen_EmbeddedFields checks the embedding cases: anonymous embedding keeps bare field names,
// a named embedded struct adds its name in front and reads through it, and a *T one writes the column name as-is.
// Associations map to no column and get skipped.
//
// TestGen_EmbeddedFields 验证嵌入的几种情况：匿名嵌入保持原字段名，
// 具名嵌入结构体在前面加上它的名字并经由它取值，指针嵌入直接写列名。
// 关联关系没有对应的列，不生成。
func TestGen_EmbeddedFields(t *testing.T) {
	options := gormcngen.NewOptions().
		WithColumnsMethodRecvName("c").
		WithColumnsCheckFieldType(true)

	output := gormcngen.NewSchemaConfig(&embedPost{}, options).Gen()
	t.Log(output.GetMethodCode())

	for _, want := range []string{
		`ID:gormcnm.Cnm(c.ID, "id"),`,
		`Name:gormcnm.Cnm(c.Name, "name"),`,
		`ContactName:gormcnm.Cnm(c.Contact.Name, "contact_name"),`,
		`ContactPhone:gormcnm.Cnm(c.Contact.Phone, "contact_phone"),`,
		`BackupName:"backup_name",`,
		`BackupPhone:"backup_phone",`,
	} {
		require.Contains(t, output.GetMethodCode(), want)
	}
	require.NotContains(t, output.GetMethodCode(), "Profile")

	tableOutput := gormcngen.NewSchemaConfig(&embedPost{}, options.WithIsGenFuncTableColumns(true)).Gen()
	require.Contains(t, tableOutput.GetMethodTableColumnsCode(), `ContactName:gormcnm.Cmn(c.Contact.Name, "contact_name", decoration),`)
	require.Contains(t, tableOutput.GetMethodTableColumnsCode(), `BackupName:gormcnm.New[string](decoration.DecorateColumnName("backup_name")),`)
}

func TestGen_IsGenFuncTableColumns(t *testing.T) {
	type Example struct {
		ID        int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
		Name      string    `gorm:"not null,type:text" form:"name" json:"name" validate:"required" cnm:"V名称"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}

	options := gormcngen.NewOptions().
		WithUseTagName(true).
		WithColumnsMethodRecvName("example").
		WithColumnsCheckFieldType(true).
		WithIsGenFuncTableColumns(true)

	config := gormcngen.NewSchemaConfig(&Example{}, options)
	output := config.Gen()
	t.Log(output.GetMethodCode())
	t.Log(output.GetMethodTableColumnsCode())
	t.Log(output.GetStructCode())
	t.Log(neatjsons.S(output.GetPkgImports()))
}
